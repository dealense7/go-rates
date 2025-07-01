package dto

type Product struct {
	Name     string             `json:"name"`
	BarCode  string             `json:"bar_code"`
	Volume   *string            `json:"volume,omitempty"`
	Image    string             `json:"image"`
	Meta     *map[string]string `json:"meta,omitempty"`
	Category int                `json:"category"`
}

type ProductPrice struct {
	StoreId  int64 `json:"name"`
	Price    int64 `json:"price"`
	OldPrice int64 `json:"old_price"`
}
