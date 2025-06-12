package services

import (
	"github.com/dealense7/go-rate-app/internal/interfaces"
	"github.com/dealense7/go-rate-app/internal/models"
)

type GasService struct {
	repo interfaces.GasRepository
}

func NewGasService(repo interfaces.GasRepository) interfaces.GasService {
	return &GasService{repo: repo}
}

func (s *GasService) GetAll() ([]models.Gas, error) {
	return s.repo.FindAll()
}
