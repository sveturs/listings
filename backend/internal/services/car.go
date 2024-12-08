package services

import (
	"backend/internal/domain/models"
    "backend/internal/storage"
	"context"

)
type CarService struct {
    storage storage.Storage
}

func NewCarService(storage storage.Storage) *CarService {
    return &CarService{storage: storage}
}

func (s *CarService) AddCar(ctx context.Context, car *models.Car) (int, error) {
    return s.storage.AddCar(ctx, car)
}

func (s *CarService) GetAvailableCars(ctx context.Context) ([]models.Car, error) {
    return s.storage.GetAvailableCars(ctx)
}
