package services

import (
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

func (s StoreService) GetForSlider() ([]models.ProductItem, error) {
	return s.repo.GetForSlider()
}

func (s StoreService) GetItemsList(offset int) ([]models.ProductItem, error) {
	return s.repo.GetItemsList(offset)
}

func (s StoreService) GetItemsCount() (int, error) {
	return s.repo.GetItemsCount()
}

func (s StoreService) GetForCategorySlider() ([]models.CategorySlider, error) {
	return s.repo.GetForCategorySlider()

}
