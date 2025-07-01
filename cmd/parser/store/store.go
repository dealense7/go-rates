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
	"github.com/redis/go-redis/v9"
	"github.com/spf13/cobra"
)

type Item struct {
	BarCode  string             `json:"bar_code"`
	Name     string             `json:"name"`
	Image    string             `json:"image"`
	Meta     *map[string]string `json:"meta"`
	Price    int64              `json:"price"`     // in tetri (₾×100)
	OldPrice int64              `json:"old_price"` // in tetri (₾×100)
	Date     string             `json:"date"`      // YYYY-MM-DD
	Volume   *string            `json:"volume"`
}

type Store struct {
	Id    enum.StoreProvider `json:"id"`
	Name  string             `json:"name"`
	Route string             `json:"route"`
	Items []Item             `json:"items"`
}

type Station struct {
	Id    enum.GasProvider `json:"id"`
	Name  string           `json:"name"`
	Items []Item           `json:"items"`
}

func (s Store) GetName() string {
	return s.Name
}
func (s Store) GetRoute() string {
	return s.Route
}
func (s Store) GetProvider() int64 {
	return int64(s.Id)
}

type StoreProvider interface {
	GetName() string
	GetRoute() string
	GetProvider() int64
	GetData(route string) ([]Item, error)
}

func init() {
	root.ParseCmd.AddCommand(storeCmd)
}

var storeCmd = &cobra.Command{
	Use:   "store",
	Short: "Parse a store prices",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("parsing store data")
		parseData()
		return nil
	},
}

func parseData() {
	// Init DB Connection
	db := utils.NewDB()
	cacheClient := utils.NewCacheClient()
	pattern := "tag:products:*"

	// Init Repository
	var repo interfaces.StoreRepository
	repo = repositories.NewMySQLStoreRepository(db)

	var items = []StoreProvider{
		NewStoreCarrefour(),
		NewStoreGoodwill(),
		NewStoreOrinabiji(),
		NewStoreAgrohub(),
		NewStoreEuroproduct(),
		NewStoreFresco(),
		NewStoreMagniti(),
	}

	allowedBarCodeLengths := map[int]bool{
		13: true,
		12: true,
		8:  true,
	}
	for _, val := range items {
		fmt.Println("parsing data for", val.GetName())

		// Fetch data for a store
		data, err := val.GetData(val.GetRoute())
		if err != nil {
			fmt.Println("error fetching:", err.Error())
			continue
		}

		// Start Transaction and Rollback if something went wrong and did not commit
		tx, err := db.Beginx()
		if err != nil {
			fmt.Println("transaction Not Started: ", err.Error())
			continue
		}
		defer tx.Rollback()

		valid := 0
		fmt.Println("data count: ", len(data))
		for _, item := range data {

			if !allowedBarCodeLengths[len(item.BarCode)] || item.BarCode == "" {
				continue
			}

			var productId int64

			// STEP I - find already existed product
			productId, err := repo.GetProductByBarCode(item.BarCode)

			// STEP II - create a product if not exists
			if err != nil {
				if !errors.Is(err, sql.ErrNoRows) {
					fmt.Println("error while searching by barCode: ", err.Error())
					continue
				}

				data := dto.Product{Name: item.Name, Meta: item.Meta, BarCode: item.BarCode, Volume: item.Volume, Image: item.Image}
				productId, err = repo.CreateItem(data)
				if err != nil {
					fmt.Println("product not created: ", err.Error())
					continue
				}
			}

			// STEP III - create or update today price on the product
			productPriceData := dto.ProductPrice{Price: item.Price, OldPrice: item.OldPrice, StoreId: val.GetProvider()}
			err = repo.AddOrUpdatePrice(productId, productPriceData)
			if err != nil {
				fmt.Println("product price not created: ", err.Error())
				continue
			}
			valid++
		}

		fmt.Println("valid data count: ", valid)

		if err := tx.Commit(); err != nil {
			fmt.Println("commit error:", err)
		}
	}

	err := repo.DisableOldPrices()
	if err != nil {
		fmt.Println("commit error:", err)
	}

	err = deleteKeysByPattern(cacheClient, pattern)
	if err != nil {
		return
	}

	fmt.Println("done")
}
func deleteKeysByPattern(client *redis.Client, pattern string) error {
	ctx := context.Background()
	// Access the underlying redis.Client

	// Create a SCAN iterator for the pattern
	iter := client.Scan(ctx, 0, pattern, 0).Iterator()

	// Iterate over matching keys
	for iter.Next(ctx) {
		key := iter.Val()
		// Delete the key
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
