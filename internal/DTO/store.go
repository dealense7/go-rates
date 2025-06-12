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
	Price    int   `json:"price"`
	OldPrice int   `json:"old_price"`
}
