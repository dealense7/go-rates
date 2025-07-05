package currency

import (
	"encoding/json"
	"fmt"
	"github.com/dealense7/go-rate-app/internal/enum"
	"io"
	"log"
	"net/http"
)

type TBCParser struct {
	Currency
}

func NewTBCParser() *TBCParser {
	return &TBCParser{
		Currency: Currency{
			Id:    enum.TBC,
			Route: "https://apigw.tbc.ge/api/v1/exchangeRates/getExchangeRate",
		},
	}
}

// GetData fetches the first 50 products and returns them.
func (g *TBCParser) GetName() string {
	return g.Id.String()
}

func (g *TBCParser) GetData() ([]Item, error) {
	var items []Item

	currencies := map[string]enum.CurrencyCode{
		enum.USD.String(): enum.USD,
		enum.EUR.String(): enum.EUR,
		enum.GBP.String(): enum.GBP,
	}

	for _, currency := range currencies {
		url := g.Route + fmt.Sprintf("?Iso1=%s&Iso2=GEL", currency.String())
		resp, _ := http.Get(url)

		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				log.Fatal(err)
			}
		}(resp.Body)

		var envelope map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
			fmt.Println("decode envelope: %w", err)
		}
		buyRate, _ := envelope["buyRate"].(float64)
		sellRate, _ := envelope["sellRate"].(float64)
		items = append(items, Item{
			Currency: currency,
			BuyRate:  int(buyRate * 1000),
			SellRate: int(sellRate * 1000),
		})

	}

	return items, nil
}
