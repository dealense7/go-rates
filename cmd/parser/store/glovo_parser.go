package store

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	url2 "net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Glovo struct {
	Store
}

func (g *Glovo) GetData(route string) ([]Item, error) {
	allRoutes := g.GetLinks(route)

	seen := make(map[string]bool)
	var items []Item

	for _, val := range allRoutes {
		if !seen[val] {
			g.fetchData(&items, val)
			seen[val] = true
		}
	}

	return items, nil
}

func (g *Glovo) fetchData(items *[]Item, path string) {
	url := "https://api.glovoapp.com/v3/" + path

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Add necessary headers
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Accept-Language", "en-GB,en-US;q=0.9,en;q=0.8,ka;q=0.7")
	req.Header.Add("Glovo-Api-Version", "14")
	req.Header.Add("Glovo-App-Platform", "web")
	req.Header.Add("Glovo-App-Type", "customer")
	req.Header.Add("Glovo-Location-City-Code", "TBI")
	req.Header.Add("Glovo-Location-Country-Code", "GE")
	req.Header.Add("Origin", "https://glovoapp.com")
	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/134.0.0.0 Safari/537.36")

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	var envelope map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
		fmt.Println("decode envelope: %w", err)
	}

	data, _ := envelope["data"].(map[string]interface{})

	products, _ := data["body"].([]interface{})

	for _, i2 := range products {
		i2 := i2.(map[string]interface{})
		content := i2["data"].(map[string]interface{})
		if content["elements"] == nil {
			continue
		}
		elements := content["elements"].([]interface{})
		g.transformItems(items, elements)
	}
}

func (g *Glovo) transformItems(items *[]Item, data []interface{}) {

	for _, raw := range data {
		m, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		m = m["data"].(map[string]interface{})

		price := int(m["price"].(float64) * 100)
		oldPrice := 0
		if strings.Contains(m["description"].(string), "Old Price/ძველი ფასი") {
			// Regular expression to find the price (assuming the price is a number with a decimal)
			re := regexp.MustCompile(`\d+\.\d+`)

			// Find the price in the description
			oldPriceDesc := re.FindString(m["description"].(string))
			floatVal, _ := strconv.ParseFloat(oldPriceDesc, 64)

			oldPrice = int(floatVal * 100)
		}

		barCode := m["externalId"].(string)

		allowedBarCodeLengths := map[int]bool{
			13: true,
			12: true,
			8:  true,
		}
		if !allowedBarCodeLengths[len(barCode)] {
			re := regexp.MustCompile(`(?:^|\D)(\d{8}|\d{12}|\d{13})(?:\D|$)`)

			match := re.FindStringSubmatch(m["name"].(string))
			if len(match) >= 2 {
				barCode = match[1]
			}
		}

		if oldPrice < price {
			oldPrice = 0
		}

		if m["imageUrl"] != nil {
			*items = append(*items, Item{
				BarCode:  barCode,
				Name:     m["name"].(string),
				Image:    m["imageUrl"].(string),
				Price:    price,
				OldPrice: oldPrice,
				Date:     time.Now().String(),
			})
		}
	}
}

func (g *Glovo) GetLinks(url string) []string {
	var matches []string

	for _, cookieValue := range getLocations() {
		// Create a new HTTP request
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Println("Error creating request:", err)
			return nil
		}

		// Set only essential headers
		req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/134.0.0.0 Safari/537.36")
		req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8")
		req.Header.Set("Accept-Language", "en-GB,en;q=0.9")
		req.Header.Set("Cache-Control", "no-cache")
		req.Header.Set("Pragma", "no-cache")

		// Set cookies (important for location-based content)

		req.Header.Set("Cookie", cookieValue)

		// Perform the HTTP request
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error performing request:", err)
			return nil
		}
		defer resp.Body.Close()

		// Read the response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading response body:", err)
			return nil
		}

		// Define the regular expression pattern
		re := regexp.MustCompile(`[^"]*nodeType=DEEP_LINK[^"]*-sc[^"]*`)

		// Find all matching links
		matches = append(matches, re.FindAllString(string(body), -1)...)
	}

	return matches
}

func getLocations() []string {
	deliveryAddress := map[string]interface{}{
		"geo": map[string]float64{
			"lat": 41.7118588,
			"lng": 44.7956066,
		},
		"city": map[string]string{
			"code":        "TBI",
			"name":        "T'bilisi",
			"countryCode": "GE",
		},
		"placeId": "ChIJA2YoWSgNREAR-9mpmkmMTAs",
		"text":    "Davit Aghmashenebeli Avenue",
		"details": "11",
	}

	jsonBytes, err := json.Marshal(deliveryAddress)
	if err != nil {
		panic(err)
	}

	encoded := url2.QueryEscape(string(jsonBytes))

	cookieValue := fmt.Sprintf(
		"glovo_user_lang=en; glovo_last_visited_city=TBI; glovo_delivery_address=%s",
		encoded,
	)

	emptyAddresscookieValue := "glovo_user_lang=en; glovo_last_visited_city=TBI; glovo_delivery_address="

	return []string{cookieValue, emptyAddresscookieValue}
}
