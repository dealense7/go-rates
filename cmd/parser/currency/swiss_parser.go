package currency

import (
	"context"
	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
	"github.com/dealense7/go-rate-app/internal/enum"
	"log"
	"strconv"
	"strings"
)

type SwissParser struct {
	Currency
}

func NewSwissParser() *SwissParser {
	return &SwissParser{
		Currency: Currency{
			Id:    enum.SWISS,
			Route: "https://swisscapital.ge/ge/currency",
		},
	}
}

// GetData fetches the first 50 products and returns them.
func (g *SwissParser) GetName() string {
	return g.Id.String()
}

func (g *SwissParser) GetData() ([]Item, error) {
	var items []Item

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"),
	)
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	var html string
	url := g.Route

	currencies := map[string]enum.CurrencyCode{
		"1 USD": enum.USD,
		"1 EUR": enum.EUR,
		"1 GBP": enum.GBP,
	}

	// navigate to page and get full HTML
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitVisible(`#block-currency_table-90`, chromedp.ByID), // wait for table
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
	doc.Find("#block-currency_table-90 tr").Each(func(i int, s *goquery.Selection) {
		name := s.Find("td").Eq(0).Find("span").Text()
		buy := s.Find("td").Eq(1).Find("span").Text()
		sell := s.Find("td").Eq(2).Find("span").Text()
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

	return items, nil
}
