package currency

import (
	"encoding/json"
	"fmt"
	"github.com/dealense7/go-rate-app/internal/enum"
	"log"
	"net/http"
	"strconv"
)

type BasisParser struct {
	Currency
}

func NewBasisParser() *BasisParser {
	return &BasisParser{
		Currency: Currency{
			Id:    enum.BASIS,
			Route: "https://static.bb.ge/source/api/view/main/getXrates",
		},
	}
}

// GetData fetches the first 50 products and returns them.
func (g *BasisParser) GetName() string {
	return g.Id.String()
}

func (g *BasisParser) GetData() ([]Item, error) {
	var items []Item
	currencies := map[string]enum.CurrencyCode{
		enum.USD.String(): enum.USD,
		enum.EUR.String(): enum.EUR,
		enum.GBP.String(): enum.GBP,
	}

	resp, _ := http.Get(g.Route)
	defer resp.Body.Close()

	var envelope []interface{}
	if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
		fmt.Println("decode envelope: %w", err)
	}

	data := envelope[0].(map[string]interface{})
	xratesStr, ok := data["xrates"].(string)
	if !ok {
		log.Fatal("xrates is not a string")
	}

	// Step 2: Parse the JSON string into a map
	var rates map[string]interface{}
	err := json.Unmarshal([]byte(xratesStr), &rates)
	if err != nil {
		log.Fatal("error decoding xrates:", err)
	}
	buyElements := rates["kursBuy"].(map[string]interface{})
	sellElements := rates["kursSell"].(map[string]interface{})

	for index, val := range buyElements {
		if currency, ok := currencies[index]; ok {
			floatVal, _ := strconv.ParseFloat(val.(string), 64)
			sellVal, _ := strconv.ParseFloat(sellElements[index].(string), 64)

			items = append(items, Item{
				Currency: currency,
				BuyRate:  int(floatVal * 1000),
				SellRate: int(sellVal * 1000),
			})

			continue
		}
	}

	return items, nil
}
