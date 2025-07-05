package currency

import (
	"encoding/json"
	"fmt"
	"github.com/dealense7/go-rate-app/internal/enum"
	"io"
	"log"
	"net/http"
)

type BOGParser struct {
	Currency
}

func NewBogParser() *BOGParser {
	return &BOGParser{
		Currency: Currency{
			Id:    enum.BOG,
			Route: "https://bankofgeorgia.ge/api/currencies/history",
		},
	}
}

// GetData fetches the first 50 products and returns them.
func (g *BOGParser) GetName() string {
	return g.Id.String()
}

func (g *BOGParser) GetData() ([]Item, error) {
	var items []Item
	currencies := [3]enum.CurrencyCode{
		enum.USD,
		enum.EUR,
		enum.GBP,
	}

	resp, _ := http.Get(g.Route)
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

	elements := envelope["data"].([]interface{})

	for _, currency := range currencies {
		for _, d := range elements {
			element := d.(map[string]interface{})
			if element["ccy"] == currency.String() {
				items = append(items, Item{
					Currency: currency,
					BuyRate:  int(element["buyRate"].(float64) * 1000),
					SellRate: int(element["sellRate"].(float64) * 1000),
				})
				break
			}
		}
	}

	return items, nil
}
