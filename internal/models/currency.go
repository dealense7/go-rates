package models

type CurrencyProvider struct {
	ID   int    `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
	Logo string `json:"logo" db:"logo_url"`
}
