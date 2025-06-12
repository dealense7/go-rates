package models

type Gas struct {
	ID           int    `json:"id" db:"id"`
	Name         string `json:"name" db:"name"`
	Tag          string `json:"tag" db:"tag"`
	Price        int32  `json:"price" db:"price"`
	ProviderName string `json:"providerName" db:"provider_name"`
	ProviderLogo string `json:"providerLogo" db:"provider_logo"`
}
