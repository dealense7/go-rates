package gas

import (
	"crypto/tls"
	"github.com/dealense7/go-rate-app/internal/enum"
	"github.com/gocolly/colly"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type GasLukoili struct {
	Station
}

func NewGasLukoili() *GasLukoili {
	return &GasLukoili{
		Station: Station{
			Id:   enum.LUKOILI,
			Name: "Lukoili",
		},
	}
}

func (g GasLukoili) GetData() ([]Item, error) {
	type FuelInfo struct {
		Name string `json:"name"`
		Tag  string `json:"tag"`
	}

	// map of product name → tag
	interesting := map[string]FuelInfo{
		"Super Ecto 100": {
			Name: "სუპერ ექტო 100",
			Tag:  "სუპერი",
		},
		"Euro Regular": {
			Name: "ევრო რეგულარი",
			Tag:  "რეგულარი",
		},
		"Super Ecto": {
			Name: "სუპერ ექტო",
			Tag:  "სუპერი",
		},
		"Premium Avangard": {
			Name: "პრემიუმ ავანგარდი",
			Tag:  "პრემიუმი",
		},
		"Euro Diesel": {
			Name: "ევრო დიზელი",
			Tag:  "დიზელი",
		},
	}

	var results []Item
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (compatible; Colly/2.1; +https://github.com/gocolly/colly)"),
	)
	c.WithTransport(&http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // ⚠️ Not secure
	})

	// run on each product tile
	c.OnHTML("body > div.w-screen > div.lg\\:grid span", func(e *colly.HTMLElement) {
		name := strings.TrimSpace(e.ChildText("div:nth-child(2) p:nth-child(2)"))

		tag, ok := interesting[name]
		if !ok {
			return
		}
		priceText := strings.TrimSpace(e.ChildText("div:nth-child(2) p:nth-child(1)"))

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
	if err := c.Visit("https://www.lukoil.ge/"); err != nil {
		return nil, err
	}

	return results, nil
}
