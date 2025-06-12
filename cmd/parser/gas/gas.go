package gas

import (
	"fmt"
	"github.com/dealense7/go-rate-app/cmd/parser/root"
	"github.com/dealense7/go-rate-app/internal/enum"
	"github.com/dealense7/go-rate-app/internal/utils"
	"github.com/spf13/cobra"
)

type Item struct {
	Name  string `json:"name"`
	Tag   string `json:"tag"`
	Price int    `json:"price"` // in tetri (₾×100)
	Date  string `json:"date"`  // YYYY-MM-DD
}

type Station struct {
	Id    enum.GasProvider `json:"id"`
	Name  string           `json:"name"`
	Items []Item           `json:"items"`
}

func (s Station) GetName() string {
	return s.Name
}
func (s Station) GetProvider() int {
	return int(s.Id)
}

type GasProvider interface {
	GetName() string
	GetProvider() int
	GetData() ([]Item, error)
}

func init() {
	root.ParseCmd.AddCommand(currencyCmd)
}

var currencyCmd = &cobra.Command{
	Use:   "gas",
	Short: "Parse a gas prices",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("parsing gas data")
		parseData()
		return nil
	},
}

func parseData() {
	var items = []GasProvider{
		NewGasConnect(),
		NewGasGulf(),
		NewGasLukoili(),
		NewGasPortal(),
		NewGasRompetrol(),
		NewGasSocar(),
		NewGasWissol(),
	}

	db := utils.NewDB()

	for _, val := range items {
		fmt.Println("parsing data for", val.GetName())
		data, err := val.GetData()

		if err != nil {
			fmt.Println("error fetching:", err.Error())
			fmt.Println("------")
			continue
		}

		for _, item := range data {
			tx := db.MustBegin()

			const query = `INSERT INTO gas_rates (provider_id, name, tag, price, date) VALUES (?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE tag = VALUES(tag), price = VALUES(price);`

			tx.MustExec(query, val.GetProvider(), item.Name, item.Tag, item.Price, item.Date)

			err = tx.Commit()
			if err != nil {
				fmt.Println(err.Error())
				fmt.Println("==========")
			}
		}
	}
}
