package store

import (
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dealense7/go-rate-app/internal/enum"
	"io"
	"net/http"
	"regexp"
	"time"
)

type StoreGoodwill struct {
	Store
}

func NewStoreGoodwill() *StoreGoodwill {
	return &StoreGoodwill{
		Store: Store{
			Id:    enum.GOODWILL,
			Name:  "Goodwill",
			Route: "https://api.goodwill.ge/v1/Products/v3?ShopId=1&Page=%d&Limit=%d",
		},
	}
}

// GetData fetches the first 50 products and returns them.
func (g *StoreGoodwill) GetData(route string) ([]Item, error) {
	// 1) Retrieve access token
	token, err := g.getToken()
	if err != nil {
		return nil, fmt.Errorf("fetch token: %w", err)
	}

	var items []Item

	g.fetchData(route, &items, token, 1)

	return items, nil
}

func (g *StoreGoodwill) fetchData(route string, items *[]Item, token string, page int) {
	limit := 500
	url := fmt.Sprintf(route, page, limit)

	resp, err := g.getProducts(url, token)
	if err != nil {
		fmt.Println("fetch products: %w", err)
	}
	defer resp.Body.Close()

	reader, err := g.getReader(resp)
	if err != nil {
		fmt.Println("decode gzip: %w", err)
	}

	// get a result as an array to map
	var envelope map[string]interface{}
	if err := json.NewDecoder(reader).Decode(&envelope); err != nil {
		fmt.Println("decode envelope: %w", err)
	}

	prodSlice, ok := envelope["products"].([]interface{})
	if !ok {
		fmt.Println("expected products to be []interface{}, got %T", envelope["products"])
	}

	if len(prodSlice) == limit {
		g.fetchData(route, items, token, page+1)
	}

	g.transformItems(items, prodSlice)
}

func (g *StoreGoodwill) getToken() (string, error) {
	resp, err := http.Get("https://goodwill.ge/")
	if err != nil {
		return "", fmt.Errorf("http.Get token page: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read token page: %w", err)
	}

	re := regexp.MustCompile(`"accessToken"\s*:\s*"([^"]+)"`)
	match := re.FindStringSubmatch(string(body))
	if len(match) < 2 {
		return "", errors.New("accessToken not found in page HTML")
	}
	return match[1], nil
}

// getProducts performs the authenticated GET and checks status.
func (g *StoreGoodwill) getProducts(url, token string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Encoding", "gzip")
	req.Header.Set("User-Agent", "Go-http-client/1.1")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("client.Do: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("unexpected status %d", resp.StatusCode)
	}
	return resp, nil
}

func (g *StoreGoodwill) getReader(resp *http.Response) (io.Reader, error) {
	if resp.Header.Get("Content-Encoding") == "gzip" {
		gz, err := gzip.NewReader(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("gzip.NewReader: %w", err)
		}
		return gz, nil
	}
	return resp.Body, nil
}

func (g *StoreGoodwill) transformItems(items *[]Item, data []interface{}) {

	for _, raw := range data {
		m, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}

		price := int(m["price"].(float64) * 100)
		oldPrice := 0
		if f, ok := m["previousPrice"].(float64); ok {
			oldPrice = int(f * 100)
		}

		image := m["imageUrl"]

		if image == nil {
			continue
		}

		*items = append(*items, Item{
			BarCode:  m["barCode"].(string),
			Name:     m["name"].(string),
			Image:    image.(string),
			Price:    price,
			OldPrice: oldPrice,
			Date:     time.Now().String(),
		})
	}
}
