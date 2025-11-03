package store

import (
	"github.com/dealense7/go-rate-app/internal/enum"
)

type StoreItem struct {
	Glovo
	Store
}

func NewStoreAgrohub() *StoreItem {
	return &StoreItem{
		Store: Store{
			Id:    enum.AGROHUB,
			Name:  "Agrohub",
			Route: "https://glovoapp.com/ge/en/tbilisi/agrohubtbi/",
		},
	}
}

func NewStoreEuroproduct() *StoreItem {
	return &StoreItem{
		Store: Store{
			Id:    enum.Europroduct,
			Name:  "Europroduct",
			Route: "https://glovoapp.com/ge/en/tbilisi/europroduct-c-tbi/",
		},
	}
}

func NewStoreFresco() *StoreItem {
	return &StoreItem{
		Store: Store{
			Id:    enum.Fresco,
			Name:  "Fresco",
			Route: "https://glovoapp.com/ge/en/tbilisi/fresco-tbi/",
		},
	}
}

func NewStoreCarrefour() *StoreItem {
	return &StoreItem{
		Store: Store{
			Id:    enum.CARREFOUR,
			Name:  "Carrefour",
			Route: "https://glovoapp.com/ge/en/tbilisi/1carrefour-tbi/",
		},
	}
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

func NewStoreMagniti() *StoreItem {
	return &StoreItem{
		Store: Store{
			Id:    enum.MAGNITI,
			Name:  "Magniti",
			Route: "https://glovoapp.com/ge/en/tbilisi/magniti-tbi/",
		},
	}
}
