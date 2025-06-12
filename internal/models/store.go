package models

import "encoding/json"

type Product struct {
	ID         int64           `json:"id" db:"id"`
	Name       string          `json:"name" db:"name"`
	Image      string          `json:"image" db:"image_url"`
	Meta       json.RawMessage `json:"meta" db:"meta"`
	Volume     *string         `json:"volume" db:"volume"`
	Origin     *string         `json:"origin" db:"origin"`
	CategoryId *int            `json:"category" db:"category_id"`
	Status     bool            `json:"status" db:"status"`
	CreatedAt  string          `json:"created_at" db:"created_at"`
}

type ProductPrice struct {
	ID        string `json:"id" db:"id"`
	ProductId int    `json:"category" db:"product_id"`
	StoreId   int    `json:"store_id" db:"store_id"`
	Price     int    `json:"price" db:"price_id"`
	OldPrice  int    `json:"old_price" db:"old_price"`
	Status    bool   `json:"status" db:"status"`
	CreatedAT string `json:"date" db:"created_at"`
}

type PriceEntry struct {
	Price     float64 `json:"price"`
	CreatedAt string  `json:"created_at"`
	OldPrice  float64 `json:"old_price"`
}

type ProductItem struct {
	ID       int64    `json:"id" db:"id"`
	Name     string   `json:"name" db:"name"`
	Company  string   `json:"company" db:"company"`
	Image    string   `json:"image" db:"image_url"`
	Volume   *string  `json:"volume" db:"volume"`
	Origin   *string  `json:"origin" db:"origin"`
	MinPrice *float64 `json:"min_price" db:"min_price"`
	MaxPrice *float64 `json:"max_price" db:"max_price"`
}

type SingleProductItem struct {
	ID      int64  `json:"id" db:"id"`
	Name    string `json:"name" db:"name"`
	Company string `json:"company" db:"company"`
	Image   string `json:"image" db:"image_url"`
	Prices  string `json:"prices" db:"prices"`
}

type CategorySlider struct {
	Name     string `json:"name" db:"name"`
	Products string `json:"products" db:"products"`
}
