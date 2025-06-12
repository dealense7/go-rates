package enum

import (
	"fmt"
	"github.com/dealense7/go-rate-app/internal/helpers"
)

type StoreProvider int
type StoreCategoryProvider int

const (
	GOODWILL  StoreProvider = 1
	CARREFOUR StoreProvider = 2
	ORINABIJI StoreProvider = 3
	AGROHUB   StoreProvider = 4
	SPAR      StoreProvider = 5
	MAGNITI   StoreProvider = 6
)

const (
	Grocery StoreCategoryProvider = 1
	Drinks  StoreCategoryProvider = 2
	Dairy   StoreCategoryProvider = 3
	Sweet   StoreCategoryProvider = 4
)

var storeNames = map[StoreProvider]string{
	GOODWILL:  "Goodwill",
	ORINABIJI: "OriNabiji",
	CARREFOUR: "Carrefour",
	AGROHUB:   "AgroHub",
	SPAR:      "Spar",
	MAGNITI:   "Magniti",
}

var categoryNames = map[StoreCategoryProvider]string{
	Grocery: "სურსათი",
	Drinks:  "სასმელი",
}

func (p StoreCategoryProvider) String() string {
	if s, ok := categoryNames[p]; ok {
		return s
	}
	return fmt.Sprintf("StoreCategoryProvider(%d)", int(p))
}

func (p StoreProvider) String() string {
	if s, ok := storeNames[p]; ok {
		return s
	}
	return fmt.Sprintf("StoreProvider(%d)", int(p))
}

func (p StoreProvider) Logo() string {
	var imagePath = map[StoreProvider]string{
		GOODWILL:  "static/img/logos/store/goodwill.webp",
		ORINABIJI: "static/img/logos/store/orinabiji.webp",
		CARREFOUR: "static/img/logos/store/goodwill.webp",
		AGROHUB:   "static/img/logos/store/orinabiji.webp",
		SPAR:      "static/img/logos/store/orinabiji.webp",
		MAGNITI:   "static/img/logos/store/magniti.webp",
	}

	if s, ok := imagePath[p]; ok {
		return s
	}
	return fmt.Sprintf("StoreProvider(%d)", int(p))
}

func (p StoreProvider) Slug() string {
	if s, ok := storeNames[p]; ok {
		return helpers.Slugify(s)
	}
	return fmt.Sprintf("StoreProvider(%d)", int(p))
}
