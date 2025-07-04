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

type CredoParser struct {
	Currency
}

func NewCredoParser() *CredoParser {
	return &CredoParser{
		Currency: Currency{
			Id:    enum.CREDO,
			Route: "https://credobank.ge/currency",
		},
	}
}

func (g *CredoParser) GetName() string {
	return g.Id.String()
}

func (g *CredoParser) GetData() ([]Item, error) {
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
	doc.Find(".exchange-rate-component").Each(func(i int, s *goquery.Selection) {
		name := strings.TrimSpace(s.Find(".currency-description").Text())
		buySell := s.Find(".buy-sell-title")

		buyRate, _ := strconv.ParseFloat(strings.TrimSpace(buySell.Eq(0).Text()), 64)  // First
		sellRate, _ := strconv.ParseFloat(strings.TrimSpace(buySell.Eq(1).Text()), 64) // First
		if currency, ok := currencies[name]; ok {
			items = append(items, Item{
				Currency: currency,
				BuyRate:  int(buyRate * 1000),
				SellRate: int(sellRate * 1000),
			})
		}
	})

	return items, nil
}
