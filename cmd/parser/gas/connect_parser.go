package gas

import (
	"github.com/dealense7/go-rate-app/internal/enum"
	"github.com/gocolly/colly"
	"log"
	"strconv"
	"strings"
	"time"
)

type GasConnect struct {
	Station
}

func NewGasConnect() *GasConnect {
	return &GasConnect{
		Station: Station{
			Id:   enum.CONNECT,
			Name: "Connect",
		},
	}
}

func (g GasConnect) GetData() ([]Item, error) {
	// map of product name → tag
	interesting := map[string]string{
		"SUPER RON 98":       "სუპერი",
		"PREMIUM RON 95":     "პრემიუმი",
		"REGULAR RON 93":     "რეგულარი",
		"EURO DIESEL 10 PPM": "დიზელი",
		"DIESEL 10 PPM":      "დიზელი",
		"LPG":                "გაზი",
	}

	var results []Item
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (compatible; Colly/2.1; +https://github.com/gocolly/colly)"),
	)

	// run on each product tile
	c.OnHTML("#wpdtSimpleTable-1 > tbody > tr", func(e *colly.HTMLElement) {
		name := strings.TrimSpace(e.ChildText("td:nth-child(1)"))
		tag, ok := interesting[name]
		if !ok {
			return
		}
		priceText := strings.TrimSpace(e.ChildText("td:nth-child(2)"))

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
	if err := c.Visit("https://connect.com.ge/"); err != nil {
		return nil, err
	}

	return results, nil
}
