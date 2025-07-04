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

type LibertyParser struct {
	Currency
}

func NewLibertyParser() *LibertyParser {
	return &LibertyParser{
		Currency: Currency{
			Id:    enum.LIBERTY,
			Route: "https://www.libertybank.ge/en/",
		},
	}
}

// GetData fetches the first 50 products and returns them.
func (g *LibertyParser) GetName() string {
	return g.Id.String()
}

func (g *LibertyParser) GetData() ([]Item, error) {
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
	seen := make(map[string]bool)

	// find exchange rate components
	doc.Find(".currency-rates__row").Each(func(i int, s *goquery.Selection) {
		name := s.Find(".currency-rates__item").Eq(0).Find(".currency-rates__currency-name").Eq(0).Text()
		buy := s.Find(".currency-rates__item").Eq(1).Find(".currency-rates__currency").Eq(0).Text()
		sell := s.Find(".currency-rates__item").Eq(1).Find(".currency-rates__currency").Eq(1).Text()

		name = strings.TrimSpace(name)
		buyRate, _ := strconv.ParseFloat(strings.TrimSpace(buy), 64)   // First
		sellRate, _ := strconv.ParseFloat(strings.TrimSpace(sell), 64) // First

		if currency, ok := currencies[name]; ok {
			if !seen[name] {
				seen[name] = true
				items = append(items, Item{
					Currency: currency,
					BuyRate:  int(buyRate * 1000),
					SellRate: int(sellRate * 1000),
				})
			}
		}
	})

	return items, nil
}
