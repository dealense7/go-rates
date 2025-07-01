package interfaces

import (
	"context"
	dto "github.com/dealense7/go-rate-app/internal/DTO"
	"github.com/dealense7/go-rate-app/internal/models"
)

type StoreRepository interface {
	GetProductById(id int) (models.SingleProductItem, error)
	GetForSlider(ctx context.Context) ([]models.ProductItem, error)
	GetItemsList(ctx context.Context, page int, category int, name string) ([]models.ProductItem, error)
	GetItemsCount(ctx context.Context, category int, name string) (int, error)
	GetForCategorySlider(ctx context.Context) ([]models.CategorySlider, error)
	GetProductByBarCode(barCode string) (int64, error)
	CreateItem(data dto.Product) (int64, error)
	AddOrUpdatePrice(itemId int64, data dto.ProductPrice) error
	DisableOldPrices() error
}

type StoreService interface {
	GetProductById(id int) (models.SingleProductItem, error)
	GetForSlider(ctx context.Context) ([]models.ProductItem, error)
	GetItemsList(ctx context.Context, page int, category int, name string) ([]models.ProductItem, error)
	GetItemsCount(ctx context.Context, category int, name string) (int, error)
	GetForCategorySlider(ctx context.Context) ([]models.CategorySlider, error)
}
