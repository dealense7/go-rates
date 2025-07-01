package services

import (
	"github.com/dealense7/go-rate-app/internal/interfaces"
	"github.com/dealense7/go-rate-app/internal/models"
)

type CategoryService struct {
	repo interfaces.CategoryRepository
}

func NewCategoryService(repo interfaces.CategoryRepository) interfaces.CategoryService {
	return &CategoryService{repo: repo}
}

func (s CategoryService) GetItems() ([]models.Category, error) {
	return s.repo.GetItems()
}
