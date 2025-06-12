package interfaces

import "github.com/dealense7/go-rate-app/internal/models"

type GasRepository interface {
	FindAll() ([]models.Gas, error)
}

type GasService interface {
	GetAll() ([]models.Gas, error)
}
