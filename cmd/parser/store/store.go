package store

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/dealense7/go-rate-app/cmd/parser/root"
	dto "github.com/dealense7/go-rate-app/internal/DTO"
	"github.com/dealense7/go-rate-app/internal/enum"
	"github.com/dealense7/go-rate-app/internal/interfaces"
	"github.com/dealense7/go-rate-app/internal/repositories"
	"github.com/dealense7/go-rate-app/internal/utils"
	"github.com/spf13/cobra"
)

type Item struct {
	BarCode  string             `json:"bar_code"`
	Name     string             `json:"name"`
	Image    string             `json:"image"`
	Meta     *map[string]string `json:"meta"`
	Price    int                `json:"price"`     // in tetri (₾×100)
	OldPrice int                `json:"old_price"` // in tetri (₾×100)
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

	// Init Repository
	var repo interfaces.StoreRepository
	repo = repositories.NewMySQLStoreRepository(db)

	var items = []StoreProvider{
		NewStoreOrinabiji(),
		NewStoreGoodwill(),
		NewStoreCarrefour(),
		NewStoreAgrohub(),
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
			valid++

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
		}

		fmt.Println("valid data count: ", len(data))

		if err := tx.Commit(); err != nil {
			fmt.Println("commit error:", err)
		}
	}

	err := repo.DisableOldPrices()
	if err != nil {
		fmt.Println("commit error:", err)
	}

	fmt.Println("done")
}
