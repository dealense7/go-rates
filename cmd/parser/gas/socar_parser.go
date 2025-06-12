package gas

import (
	"encoding/json"
	"github.com/dealense7/go-rate-app/internal/enum"
	"io"
	"log"
	"net/http"
	"time"
)

type GasSocar struct {
	Station
}

func NewGasSocar() *GasSocar {
	return &GasSocar{
		Station: Station{
			Id:   enum.SOCAR,
			Name: "Socar",
		},
	}
}

func (g GasSocar) GetData() ([]Item, error) {
	type Price struct {
		ActionDate    string  `json:"ActionDate"`
		FuelCode      string  `json:"FuelCode"`
		FuelNameEng   string  `json:"FuelNameEng"`
		FuelNameGeo   string  `json:"FuelNameGeo"`
		FuelUnitPrice float64 `json:"FuelUnitPrice"`
	}

	type GetCurrentPricesResp struct {
		Results []Price `json:"Results"`
	}

	type Envelope struct {
		GetCurrentPrices GetCurrentPricesResp `json:"GetCurrentPrices"`
	}

	// map of product name → tag
	interesting := map[string]string{
		"ნანო სუპერი":        "სუპერი",
		"ნანო პრემიუმი":      "პრემიუმი",
		"ნანო ევრო რეგულარი": "რეგულარი",
		"ევრო 5 დიზელი":      "დიზელი",
		"ნანო ევრო 5 დიზელი": "დიზელი",
		"თხევადი აირი":       "გაზი",
		"ბუნებრივი აირი":     "გაზი",
	}

	resp, err := http.Get("https://sgp.ge/sgp-backend/api/integration/info/current-prices")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var results []Item
	var env Envelope
	if err := json.Unmarshal(body, &env); err != nil {
		log.Fatal(err)
	}

	// now you have a []Price
	for _, p := range env.GetCurrentPrices.Results {
		results = append(results, Item{
			Name:  p.FuelNameGeo,
			Tag:   interesting[p.FuelNameGeo],
			Price: int(p.FuelUnitPrice * 100),
			Date:  time.Now().Format("2006-01-02"),
		})
	}

	return results, nil
}
