package interfaces

import (
	"github.com/dealense7/go-rate-app/internal/models"
)

type CategoryRepository interface {
	GetItems() ([]models.Category, error)
}

type CategoryService interface {
	GetItems() ([]models.Category, error)
}
