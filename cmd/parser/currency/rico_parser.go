package currency

import (
	"context"
	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
	"github.com/dealense7/go-rate-app/internal/enum"
	"log"
	"strconv"
	"strings"
	"time"
)

type RicoParser struct {
	Currency
}

func NewRicoParser() *RicoParser {
	return &RicoParser{
		Currency: Currency{
			Id:    enum.RICO,
			Route: "https://www.rico.ge/ka",
		},
	}
}

// GetData fetches the first 50 products and returns them.
func (g *RicoParser) GetName() string {
	return g.Id.String()
}

func (g *RicoParser) GetData() ([]Item, error) {
	var items []Item

	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var html string
	url := g.Route

	currencies := map[string]enum.CurrencyCode{
		enum.USD.String(): enum.USD,
		enum.EUR.String(): enum.EUR,
		enum.GBP.String(): enum.GBP,
	}

	// navigate to page and get full HTML
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.Sleep(2*time.Second), // optional: wait for JS to fully render
		chromedp.WaitReady("body", chromedp.ByQuery),
		chromedp.OuterHTML("html", &html),
	)
	if err != nil {
		log.Fatal(err)
	}

	// parse HTML using goquery
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		log.Fatal(err)
	}

	// find exchange rate components
	doc.Find("body > main > div > section.calculators-section > div > section > div.currencies > div.currencies-table > div > table > tbody").Each(func(i int, s *goquery.Selection) {
		s.Find("tr").Each(func(i int, g *goquery.Selection) {
			name := g.Find(".flag-title").Text()
			buy := g.Find(".currency-value").Eq(0).Text()
			sell := g.Find(".currency-value ").Eq(1).Text()
			name = strings.TrimSpace(name)
			buyRate, _ := strconv.ParseFloat(strings.TrimSpace(buy), 64)   // First
			sellRate, _ := strconv.ParseFloat(strings.TrimSpace(sell), 64) // First
			if currency, ok := currencies[name]; ok {
				items = append(items, Item{
					Currency: currency,
					BuyRate:  int(buyRate * 1000),
					SellRate: int(sellRate * 1000),
				})
			}
		})
	})

	return items, nil
}
