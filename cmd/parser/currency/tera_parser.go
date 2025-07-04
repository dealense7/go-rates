package currency

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/dealense7/go-rate-app/internal/enum"
	"net/http"
)

type TeraParser struct {
	Currency
}

func NewTeraParser() *TeraParser {
	return &TeraParser{
		Currency: Currency{
			Id:    enum.TERA,
			Route: "https://terabank.ge/_mvcapi/CurrencyRatesApi/GetTeraCrossRates",
		},
	}
}

// GetData fetches the first 50 products and returns them.
func (g *TeraParser) GetName() string {
	return g.Id.String()
}

func (g *TeraParser) GetData() ([]Item, error) {
	var items []Item

	currencies := map[string]enum.CurrencyCode{
		enum.USD.String(): enum.USD,
		enum.EUR.String(): enum.EUR,
		enum.GBP.String(): enum.GBP,
	}

	resp, _ := http.Post(g.Route, "application/json", bytes.NewBuffer([]byte{}))

	defer resp.Body.Close()

	var envelope map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
		fmt.Println("decode envelope: %w", err)
	}

	for _, item := range envelope["data"].([]interface{}) {
		data := item.(map[string]interface{})

		if currency, ok := currencies[data["iso"].(string)]; ok {
			items = append(items, Item{
				Currency: currency,
				BuyRate:  int(data["teraCrossRateBuy"].(float64) * 1000),
				SellRate: int(data["teraCrossRateSell"].(float64) * 1000),
			})
		}
	}

	return items, nil
}
