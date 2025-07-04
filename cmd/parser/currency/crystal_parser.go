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

type CrystalParser struct {
	Currency
}

func NewCrystalParser() *CrystalParser {
	return &CrystalParser{
		Currency: Currency{
			Id:    enum.CRYSTAL,
			Route: "https://crystal.ge/valutis-kursebi",
		},
	}
}

func (g *CrystalParser) GetName() string {
	return g.Id.String()
}

func (g *CrystalParser) GetData() ([]Item, error) {
	var items []Item

	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var html string
	url := g.Route

	currencies := map[string]enum.CurrencyCode{
		"1დოლარი": enum.USD,
		"1ევრო":   enum.EUR,
		"1ფუნტი":  enum.GBP,
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
	doc.Find("#currencyData > div").Each(func(i int, s *goquery.Selection) {
		name := s.Find("p").Eq(0).Text()

		buyRate, _ := strconv.ParseFloat(strings.TrimSpace(s.Find(".uk-flex-middle").Eq(3).Find("span").Eq(0).Text()), 64)  // First
		sellRate, _ := strconv.ParseFloat(strings.TrimSpace(s.Find(".uk-flex-middle").Eq(5).Find("span").Eq(0).Text()), 64) // First
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
