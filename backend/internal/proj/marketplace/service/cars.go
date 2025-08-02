package service

import (
	"context"

	"backend/internal/domain/models"
)

// GetCarMakes возвращает список марок автомобилей с фильтрацией
func (s *MarketplaceService) GetCarMakes(ctx context.Context, country string, isDomestic bool, isMotorcycle bool, activeOnly bool) ([]models.CarMake, error) {
	return s.storage.GetCarMakes(ctx, country, isDomestic, isMotorcycle, activeOnly)
}

// GetCarModelsByMake возвращает модели автомобилей для конкретной марки
func (s *MarketplaceService) GetCarModelsByMake(ctx context.Context, makeSlug string, activeOnly bool) ([]models.CarModel, error) {
	return s.storage.GetCarModelsByMake(ctx, makeSlug, activeOnly)
}

// GetCarGenerationsByModel возвращает поколения для конкретной модели
func (s *MarketplaceService) GetCarGenerationsByModel(ctx context.Context, modelID int, activeOnly bool) ([]models.CarGeneration, error) {
	return s.storage.GetCarGenerationsByModel(ctx, modelID, activeOnly)
}

// SearchCarMakes выполняет поиск марок автомобилей по названию
func (s *MarketplaceService) SearchCarMakes(ctx context.Context, query string, limit int) ([]models.CarMake, error) {
	return s.storage.SearchCarMakes(ctx, query, limit)
}
