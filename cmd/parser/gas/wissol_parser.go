package gas

import (
	"github.com/dealense7/go-rate-app/internal/enum"
	"github.com/gocolly/colly"
	"log"
	"strconv"
	"strings"
	"time"
)

type GasWissol struct {
	Station
}

func NewGasWissol() *GasWissol {
	return &GasWissol{
		Station: Station{
			Id:   enum.WISSOL,
			Name: "Wissol",
		},
	}
}

func (g GasWissol) GetData() ([]Item, error) {
	// map of product name → tag
	interesting := map[string]string{
		"ეკო სუპერი":    "სუპერი",
		"ეკო პრემიუმი":  "პრემიუმი",
		"ეკო დიზელი":    "დიზელი",
		"დიზელ ენერჯი":  "დიზელი",
		"ევრო რეგულარი": "რეგულარი",
		"ვისოლ გაზი":    "გაზი",
	}

	var results []Item
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (compatible; Colly/2.1; +https://github.com/gocolly/colly)"),
	)

	// run on each product tile
	c.OnHTML(".prices_wrapper > ul > li", func(e *colly.HTMLElement) {
		name := strings.TrimSpace(e.ChildText("span p"))
		tag, ok := interesting[name]
		if !ok {
			return
		}
		priceText := strings.TrimSpace(e.ChildText("button > p:nth-child(2)"))

		priceFloat, err := strconv.ParseFloat(strings.ReplaceAll(priceText, " ₾", ""), 64)
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
	if err := c.Visit("https://wissol.ge/ka/fuel-prices"); err != nil {
		return nil, err
	}

	return results, nil
}
