package store

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/dealense7/go-rate-app/internal/enum"
	"io"
	"net/http"
	"time"
)

type StoreOrinabiji struct {
	Store
}

func NewStoreOrinabiji() *StoreOrinabiji {
	return &StoreOrinabiji{
		Store: Store{
			Id:    enum.ORINABIJI,
			Name:  "OriNabiji",
			Route: "https://catalog-api.orinabiji.ge/catalog/api/products/search?lang=ge&sortField=isInStock&sortDirection=-1",
		},
	}
}

// GetData fetches the first 50 products and returns them.
func (g *StoreOrinabiji) GetData(route string) ([]Item, error) {
	var items []Item

	var interestingCategories = []struct {
		category enum.StoreCategoryProvider
		ids      []string
	}{
		{
			category: enum.Grocery,
			ids: []string{
				"60224a31a8e27e0010eaffb4", "6221fc8b241c190011299152", "6221fca4241c190011299153",
				"6221fcc2241c190011299154", "6221fcf6241c190011299155", "6222023d241c19001129915d",
				"62220258241c19001129915e", "62a7464cf6d4b9001477a549", "6221fd61241c190011299156",
				"6222026f241c19001129915f", "62220291241c190011299160", "622202f8241c190011299161",
				"62220328241c190011299162", "6222033e241c190011299163", "6221fe86241c190011299157",
				"622203ce241c190011299167", "622203ee241c190011299168", "623c5b48e0a53f00163895d3",
				"6221fec8241c190011299159", "62220406241c190011299169", "623886b98c55b100145f49a6",
				"6222012c241c19001129915a", "6222015c241c19001129915b", "622204a4241c19001129916c",
				"622204c1241c19001129916d", "622204e1241c19001129916e", "6222050e241c190011299170",
				"62220533241c190011299171", "6222054d241c190011299172", "622201c1241c19001129915c",
				"6222056d241c190011299173", "622205a5241c190011299174", "626910eee3b90f0015c8c940",
				"62275e07241c1900112992dc", "60227d9cc51f2700106a943a", "6022873ca8e27e0010eaffe2",
				"60228755a8e27e0010eaffe3", "6022876dc51f2700106a9452", "623885d7e42ba80016883e54",
			},
		},
		{
			category: enum.Grocery,
			ids: []string{
				"6221f6ac241c190011299149", "602289bba8e27e0010eaffec", "602289e9c51f2700106a945d",
				"60228a04a8e27e0010eaffed", "60228b4aa8e27e0010eaffee", "60228b65a8e27e0010eaffef",
				"60228b76c51f2700106a945f", "60228bb2c51f2700106a9460", "6221f824241c19001129914a",
				"66bb9313f31739062e475d74",
			},
		},
		{
			category: enum.Drinks,
			ids: []string{
				"602248f8a8e27e0010eaffb2", "6220d81b241c1900112990c4", "6220d9df241c1900112990ce",
				"6220da8d241c1900112990d1", "6220d838241c1900112990c5", "6220d84f241c1900112990c6",
				"6221c6e1241c1900112990f8", "6221c6f8241c1900112990f9", "6220d867241c1900112990c7",
				"6221c724241c1900112990fa", "6221c73d241c1900112990fb", "6221c760241c1900112990fc",
				"6220d8df241c1900112990c9", "6221c7dc241c1900112990fe", "6221c8a8241c1900112990ff",
				"6221c8cd241c190011299100", "6221c8ea241c190011299101", "6220d944241c1900112990cb",
				"6221c906241c190011299102", "6221c921241c190011299103", "6221c93d241c190011299104",
				"6220da07241c1900112990cf", "6221c6ca241c1900112990f7", "6220d8a7241c1900112990c8",
				"6221c7aa241c1900112990fd", "6220d906241c1900112990ca", "62823152344af000159dd75a",
				"6220d974241c1900112990cc", "6221c96d241c190011299106", "6221c99d241c190011299107",
				"6220d9b7241c1900112990cd",
			},
		},
		{
			category: enum.Drinks,
			ids: []string{
				"60224a5bc51f2700106a942c", "60227e6ec51f2700106a943b", "637cb2bb8c3cc07140570dfc",
				"60228e1ba8e27e0010eafff9", "60228e38c51f2700106a9466", "60228e4cc51f2700106a9467",
				"623888908c55b100145f49a7", "60227e8aa8e27e0010eaffcb", "60228ec1a8e27e0010eafffa",
				"60228ed4c51f2700106a9468", "60228eefc51f2700106a9469", "6221f528241c190011299146",
				"6221f678241c190011299148", "60227ea0a8e27e0010eaffcc",
			},
		},
	}

	for _, val := range interestingCategories {
		g.fetchData(route, &items, val.ids, val.category, 0)
	}

	return items, nil
}

func (g *StoreOrinabiji) fetchData(route string, items *[]Item, data []string, categoryId enum.StoreCategoryProvider, skip int) {
	url := route

	limit := 500

	payload := map[string]interface{}{
		"skip":        skip,
		"limit":       limit,
		"categoryIds": data,
	}

	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("failed to encode payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		fmt.Printf("failed to create request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Project-Key", "CATALOG")

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("request failed: %w", err)
	}
	defer resp.Body.Close()
	reader, err := g.getReader(resp)
	if err != nil {
		fmt.Println("decode gzip: %w", err)
	}

	var envelope map[string]interface{}
	if err := json.NewDecoder(reader).Decode(&envelope); err != nil {
		fmt.Printf("decode envelope: %v\n", err)
		return
	}

	dataMap, ok := envelope["data"].(map[string]interface{})
	if !ok {
		fmt.Println("data is not a map")
		return
	}

	prodSlice, ok := dataMap["products"].([]interface{})
	if !ok {
		fmt.Println("products is not a slice")
		return
	}

	if len(prodSlice) == limit {
		g.fetchData(route, items, data, categoryId, skip+limit)
	}

	g.transformItems(items, prodSlice, categoryId)
}

func (g *StoreOrinabiji) getReader(resp *http.Response) (io.Reader, error) {
	if resp.Header.Get("Content-Encoding") == "gzip" {
		gz, err := gzip.NewReader(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("gzip.NewReader: %w", err)
		}
		return gz, nil
	}
	return resp.Body, nil
}

func (g *StoreOrinabiji) transformItems(items *[]Item, data []interface{}, categoryId enum.StoreCategoryProvider) {

	for _, raw := range data {
		m, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		var imageId string

		stock, _ := m["stock"].(map[string]interface{})
		price := int(stock["price"].(float64) * 100)
		images, _ := m["images"].([]interface{})
		image, _ := images[0].(map[string]interface{})
		imageId = image["imageId"].(string)

		var oldPrice int
		if discountRaw, ok := m["discount"]; ok && discountRaw != nil {
			discountMap, ok := discountRaw.(map[string]interface{})
			if ok {
				if priceVal, ok := discountMap["price"].(float64); ok {
					oldPrice = price
					price = int(priceVal * 100)
				}
			}
		}
		volume := fmt.Sprintf("%v", m["productNetWeight"])
		description := m["description"].(string)

		*items = append(*items, Item{
			BarCode:  m["barCode"].(string),
			Name:     m["title"].(string),
			Image:    fmt.Sprintf("https://first.media.2nabiji.ge/api/files/resize/500/500/%v/.webp", imageId),
			Price:    price,
			Meta:     &map[string]string{"description": description},
			OldPrice: oldPrice,
			Date:     time.Now().String(),
			Volume:   &volume,
		})
	}
}
