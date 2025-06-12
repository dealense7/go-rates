package gas

import (
	"github.com/dealense7/go-rate-app/internal/enum"
	"github.com/gocolly/colly"
	"log"
	"strconv"
	"strings"
	"time"
)

type GasGulf struct {
	Station
}

func NewGasGulf() *GasGulf {
	return &GasGulf{
		Station: Station{
			Id:   enum.GULF,
			Name: "Gulf",
		},
	}
}

func (g GasGulf) GetData() ([]Item, error) {
	// map of product name → tag
	interesting := map[string]string{
		"G-Force სუპერი":        "სუპერი",
		"G-Force პრემიუმი":      "პრემიუმი",
		"G-Force ევრო რეგულარი": "რეგულარი",
		"ევრო რეგულარი":         "რეგულარი",
		"G-Force ევრო დიზელი":   "დიზელი",
		"ევრო დიზელი":           "დიზელი",
		"გაზი":                  "გაზი",
	}

	var results []Item
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (compatible; Colly/2.1; +https://github.com/gocolly/colly)"),
	)

	// run on each product tile
	c.OnHTML(".price_entry", func(e *colly.HTMLElement) {
		name := strings.TrimSpace(e.ChildText(".product_name"))
		tag, ok := interesting[name]
		if !ok {
			return
		}
		priceText := strings.TrimSpace(e.ChildText(".product_price"))

		priceFloat, err := strconv.ParseFloat(strings.ReplaceAll(priceText, ",", "."), 64)
		if err != nil {
			log.Printf("warning: could not parse price %q: %v", priceText, err)
			return
		}

		price := int(priceFloat * 100)
		if price == 0 {
			return
		}

		results = append(results, Item{
			Name:  name,
			Tag:   tag,
			Price: price,
			Date:  time.Now().Format("2006-01-02"),
		})
	})

	// visit the page
	if err := c.Visit("https://gulf.ge/"); err != nil {
		return nil, err
	}

	return results, nil
}
