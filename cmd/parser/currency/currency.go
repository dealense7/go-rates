package currency

import (
	"fmt"
	"github.com/dealense7/go-rate-app/cmd/parser/root"
	"github.com/dealense7/go-rate-app/internal/enum"
	"github.com/spf13/cobra"
)

type Currency struct {
	Id    enum.CurrencyProvider
	Route string
	Items []Item
}

type Item struct {
	Currency enum.CurrencyCode
	BuyRate  int
	SellRate int
}

func init() {
	root.ParseCmd.AddCommand(currencyCmd)
}

var currencyCmd = &cobra.Command{
	Use:   "currency",
	Short: "Parse a currency amount",
	Long:  `Parses the given currency amount (e.g. "USD 123.45").`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("parsing currency data")
		parseData()
		return nil
	},
}

type CurrencyProvider interface {
	GetData() ([]Item, error)
	GetName() string
}

func parseData() {
	// Init DB Connection
	//db := utils.NewDB()

	// Init Repository
	//var repo interfaces.StoreRepository
	//repo = repositories.NewMySQLStoreRepository(db)

	var items = []CurrencyProvider{
		//NewBogParser(),
		//NewBasisParser(),
		//NewCredoParser(),
		//NewCrystalParser(),
		//NewKursiParser(),
		//NewLibertyParser(),
		//NewRicoParser(),
		//NewSwissParser(),
		//NewTBCParser(),
		NewTeraParser(),
	}

	for _, item := range items {
		fmt.Println("parsing item", item.GetName())
		fetchedItems, _ := item.GetData()
		fmt.Println(fetchedItems)
	}
}
