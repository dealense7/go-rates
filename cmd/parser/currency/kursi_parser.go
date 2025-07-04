package currency

import (
	"encoding/json"
	"fmt"
	"github.com/dealense7/go-rate-app/internal/enum"
	"net/http"
)

type KursiParser struct {
	Currency
}

func NewKursiParser() *KursiParser {
	return &KursiParser{
		Currency: Currency{
			Id:    enum.KURSI,
			Route: "https://api.kursi.ge/api/public/currencies",
		},
	}
}

// GetData fetches the first 50 products and returns them.
func (g *KursiParser) GetName() string {
	return g.Id.String()
}

func (g *KursiParser) GetData() ([]Item, error) {
	var items []Item
	currencies := [3]enum.CurrencyCode{
		enum.USD,
		enum.EUR,
		enum.GBP,
	}

	resp, _ := http.Get(g.Route)
	defer resp.Body.Close()

	var envelope []interface{}
	if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
		fmt.Println("decode envelope: %w", err)
	}

	for _, currency := range currencies {
		for _, item := range envelope {
			element := item.(map[string]interface{})
			baseCurrencyCode := element["baseCurrencyCode"]
			secondaryCurrencyCode := element["secondaryCurrencyCode"]

			if baseCurrencyCode == "GEL" && secondaryCurrencyCode == currency.String() {
				items = append(items, Item{
					Currency: currency,
					BuyRate:  int(element["buyRate"].(float64) * 1000),
					SellRate: int(element["sellRate"].(float64) * 1000),
				})
			}
		}
	}

	return items, nil
}
