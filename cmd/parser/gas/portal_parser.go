package gas

import (
	"github.com/dealense7/go-rate-app/internal/enum"
	"github.com/gocolly/colly"
	"log"
	"strconv"
	"strings"
	"time"
)

type GasPortal struct {
	Station
}

func NewGasPortal() *GasPortal {
	return &GasPortal{
		Station: Station{
			Id:   enum.PORTAL,
			Name: "Portal",
		},
	}
}

func (g GasPortal) GetData() ([]Item, error) {
	type FuelInfo struct {
		Name string `json:"name"`
		Tag  string `json:"tag"`
	}

	interesting := map[string]FuelInfo{
		"SUPER": {
			Name: "სუპერი",
			Tag:  "სუპერი",
		},
		"PREMIUM": {
			Name: "პრემიუმი",
			Tag:  "პრემიუმი",
		},
		"EURO REGULAR": {
			Name: "ევრო რეგულარი",
			Tag:  "რეგულარი",
		},
		"EURO DIESEL": {
			Name: "ევრო დიზელი",
			Tag:  "დიზელი",
		},
		"EFFECT DIESEL": {
			Name: "ეფექტ დიზელი",
			Tag:  "დიზელი",
		},
	}

	var results []Item
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (compatible; Colly/2.1; +https://github.com/gocolly/colly)"),
	)

	// run on each product tile
	c.OnHTML("body > section > div.content_div > div > div", func(e *colly.HTMLElement) {
		name := strings.TrimSpace(e.ChildText("h3"))
		tag, ok := interesting[name]
		if !ok {
			return
		}
		priceText := strings.TrimSpace(e.ChildText(".fuel_price"))

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
			Name:  tag.Name,
			Tag:   tag.Tag,
			Price: price,
			Date:  time.Now().Format("2006-01-02"),
		})
	})

	// visit the page
	if err := c.Visit("https://portal.com.ge/georgian/newfuel"); err != nil {
		return nil, err
	}

	return results, nil
}
