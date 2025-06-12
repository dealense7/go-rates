package interfaces

import (
	dto "github.com/dealense7/go-rate-app/internal/DTO"
	"github.com/dealense7/go-rate-app/internal/models"
)

type StoreRepository interface {
	GetProductById(id int) (models.SingleProductItem, error)
	GetForSlider() ([]models.ProductItem, error)
	GetItemsList(offset int) ([]models.ProductItem, error)
	GetItemsCount() (int, error)
	GetForCategorySlider() ([]models.CategorySlider, error)
	GetProductByBarCode(barCode string) (int64, error)
	CreateItem(data dto.Product) (int64, error)
	AddOrUpdatePrice(itemId int64, data dto.ProductPrice) error
	DisableOldPrices() error
}

type StoreService interface {
	GetProductById(id int) (models.SingleProductItem, error)
	GetForSlider() ([]models.ProductItem, error)
	GetItemsList(offset int) ([]models.ProductItem, error)
	GetItemsCount() (int, error)
	GetForCategorySlider() ([]models.CategorySlider, error)
}
