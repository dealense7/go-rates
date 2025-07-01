package services

import (
	"context"
	"github.com/dealense7/go-rate-app/internal/interfaces"
	"github.com/dealense7/go-rate-app/internal/models"
)

type StoreService struct {
	repo interfaces.StoreRepository
}

func NewStoreService(repo interfaces.StoreRepository) interfaces.StoreService {
	return &StoreService{repo: repo}
}

func (s StoreService) GetProductById(id int) (models.SingleProductItem, error) {
	return s.repo.GetProductById(id)
}

func (s StoreService) GetForSlider(ctx context.Context) ([]models.ProductItem, error) {
	return s.repo.GetForSlider(ctx)
}

func (s StoreService) GetItemsList(ctx context.Context, page int, category int, name string) ([]models.ProductItem, error) {
	return s.repo.GetItemsList(ctx, page, category, name)
}

func (s StoreService) GetItemsCount(ctx context.Context, category int, name string) (int, error) {
	return s.repo.GetItemsCount(ctx, category, name)
}

func (s StoreService) GetForCategorySlider(ctx context.Context) ([]models.CategorySlider, error) {
	return s.repo.GetForCategorySlider(ctx)

}
