package store

import (
	"github.com/dealense7/go-rate-app/internal/enum"
)

type StoreAgrohub struct {
	Glovo
	Store
}

type StoreCarrefour struct {
	Glovo
	Store
}

type StoreMagniti struct {
	Glovo
	Store
}

func NewStoreAgrohub() *StoreAgrohub {
	return &StoreAgrohub{
		Store: Store{
			Id:    enum.AGROHUB,
			Name:  "Agrohub",
			Route: "https://glovoapp.com/ge/en/tbilisi/agrohubtbi/",
		},
	}
}

func NewStoreCarrefour() *StoreCarrefour {
	return &StoreCarrefour{
		Store: Store{
			Id:    enum.CARREFOUR,
			Name:  "Carrefour",
			Route: "https://glovoapp.com/ge/en/tbilisi/1carrefour-tbi/",
		},
	}
}

func NewStoreMagniti() *StoreMagniti {
	return &StoreMagniti{
		Store: Store{
			Id:    enum.MAGNITI,
			Name:  "Magniti",
			Route: "https://glovoapp.com/ge/en/tbilisi/magniti-tbi/",
		},
	}
}
