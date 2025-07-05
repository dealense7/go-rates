package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/dealense7/go-rate-app/cmd/parser/root"
	dto "github.com/dealense7/go-rate-app/internal/DTO"
	"github.com/dealense7/go-rate-app/internal/enum"
	"github.com/dealense7/go-rate-app/internal/interfaces"
	"github.com/dealense7/go-rate-app/internal/repositories"
	"github.com/dealense7/go-rate-app/internal/utils"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/cobra"
)

// Item represents a product with details like barcode, name, price, and metadata.
type Item struct {
	BarCode  string             `json:"bar_code"`  // Unique product identifier (e.g., EAN-13, UPC)
	Name     string             `json:"name"`      // Product name
	Image    string             `json:"image"`     // Product image URL
	Meta     *map[string]string `json:"meta"`      // Optional metadata (key-value pairs)
	Price    int64              `json:"price"`     // Current price in tetri (₾×100)
	OldPrice int64              `json:"old_price"` // Previous price in tetri (₾×100)
	Date     string             `json:"date"`      // Date in YYYY-MM-DD format
	Volume   *string            `json:"volume"`    // Optional product volume (e.g., "500ml")
}

// Store represents a store with a unique ID, name, data source, and list of products.
type Store struct {
	Id    enum.StoreProvider `json:"id"`    // Unique store identifier (enum)
	Name  string             `json:"name"`  // Store name (e.g., Carrefour)
	Route string             `json:"route"` // URL or endpoint to fetch store data
	Items []Item             `json:"items"` // List of products in the store
}

// Station represents a gas station with a unique ID, name, and list of products (e.g., fuel types).
type Station struct {
	Id    enum.GasProvider `json:"id"`    // Unique gas station identifier (enum)
	Name  string           `json:"name"`  // Gas station name
	Items []Item           `json:"items"` // List of products
}

// StoreProvider defines the interface for interacting with store data.
type StoreProvider interface {
	GetName() string                      // Returns the store's name
	GetRoute() string                     // Returns the store's data source (e.g., URL)
	GetProvider() int64                   // Returns the store's unique ID
	GetData(route string) ([]Item, error) // Fetches product data from the route
}

// GetName returns the store's name.
func (s Store) GetName() string {
	return s.Name
}

// GetRoute returns the store's data source (e.g., URL).
func (s Store) GetRoute() string {
	return s.Route
}

// GetProvider returns the store's unique ID as an int64.
func (s Store) GetProvider() int64 {
	return int64(s.Id)
}

// init registers the "store" command with the root command.
func init() {
	root.ParseCmd.AddCommand(storeCmd)
}

// storeCmd defines the Cobra command for parsing store prices.
var storeCmd = &cobra.Command{
	Use:   "store",
	Short: "Parse a store prices",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("parsing store data")
		parseData()
		return nil
	},
}

// parseData fetches, validates, and stores product data for multiple stores.
func parseData() {
	// Initialize database and Redis connections
	db := utils.NewDB()
	cacheClient := utils.NewCacheClient()
	pattern := "tag:products:*" // Redis key pattern for product cache

	// Initialize repository for database operations
	var repo interfaces.StoreRepository
	repo = repositories.NewMySQLStoreRepository(db)

	// List of stores to process
	var items = []StoreProvider{
		NewStoreCarrefour(),
		NewStoreGoodwill(),
		NewStoreOrinabiji(),
		NewStoreAgrohub(),
		NewStoreEuroproduct(),
		NewStoreFresco(),
		NewStoreMagniti(),
	}

	// Define valid barcode lengths (EAN-13, UPC-A, EAN-8)
	allowedBarCodeLengths := map[int]bool{
		13: true, // EAN-13
		12: true, // UPC-A
		8:  true, // EAN-8
	}

	// Process each store
	for _, val := range items {
		fmt.Println("parsing data for", val.GetName())

		// Fetch product data from the store
		data, err := val.GetData(val.GetRoute())
		if err != nil {
			fmt.Println("error fetching:", err.Error())
			continue
		}

		// Start a database transaction
		tx, err := db.Beginx()
		if err != nil {
			fmt.Println("transaction Not Started: ", err.Error())
			continue
		}
		defer func(tx *sqlx.Tx) {
			err := tx.Rollback()
			if err != nil && !errors.Is(err, sql.ErrTxDone) {
				fmt.Println("rollback error:", err)
			}
		}(tx)

		valid := 0
		fmt.Println("data count: ", len(data))

		// Process each item in the store's data
		for _, item := range data {
			// Validate barcode
			if !allowedBarCodeLengths[len(item.BarCode)] || item.BarCode == "" {
				continue
			}

			var productId int64

			// STEP I: Check if the product exists by barcode
			productId, err := repo.GetProductByBarCode(item.BarCode)
			if err != nil {
				if !errors.Is(err, sql.ErrNoRows) {
					fmt.Println("error while searching by barCode: ", err.Error())
					continue
				}

				// STEP II: Create a product if it doesn't exist
				data := dto.Product{
					Name:    item.Name,
					Meta:    item.Meta,
					BarCode: item.BarCode,
					Volume:  item.Volume,
					Image:   item.Image,
				}
				productId, err = repo.CreateItem(data)
				if err != nil {
					fmt.Println("product not created: ", err.Error())
					continue
				}
			}

			// STEP III: Add or update today's price for the product
			productPriceData := dto.ProductPrice{
				Price:    item.Price,
				OldPrice: item.OldPrice,
				StoreId:  val.GetProvider(),
			}
			err = repo.AddOrUpdatePrice(productId, productPriceData)
			if err != nil {
				fmt.Println("product price not created: ", err.Error())
				continue
			}
			valid++
		}

		fmt.Println("valid data count: ", valid)

		// Commit the transaction
		if err := tx.Commit(); err != nil {
			fmt.Println("commit error:", err)
		}
	}

	// Disable outdated prices in the database
	err := repo.DisableOldPrices()
	if err != nil {
		fmt.Println("commit error:", err)
	}

	// Clear Redis cache for product data
	err = deleteKeysByPattern(cacheClient, pattern)
	if err != nil {
		fmt.Println("cache clear error:", err)
	}

	fmt.Println("done")
}

// deleteKeysByPattern removes Redis keys matching the given pattern.
func deleteKeysByPattern(client *redis.Client, pattern string) error {
	ctx := context.Background()
	// Create a SCAN iterator for the pattern
	iter := client.Scan(ctx, 0, pattern, 0).Iterator()

	// Iterate over matching keys and delete them
	for iter.Next(ctx) {
		key := iter.Val()
		err := client.Del(ctx, key).Err()
		if err != nil {
			return fmt.Errorf("failed to delete key %s: %v", key, err)
		}
	}

	// Check for errors during iteration
	if err := iter.Err(); err != nil {
		return fmt.Errorf("error during scan: %v", err)
	}

	return nil
}
