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

// GetCarMakeBySlug возвращает марку автомобиля по slug
func (s *MarketplaceService) GetCarMakeBySlug(ctx context.Context, slug string) (*models.CarMake, error) {
	return s.storage.GetCarMakeBySlug(ctx, slug)
}

// GetCarListingsCount возвращает количество автомобильных объявлений
func (s *MarketplaceService) GetCarListingsCount(ctx context.Context) (int, error) {
	return s.storage.GetCarListingsCount(ctx)
}

// GetTotalCarModelsCount возвращает общее количество моделей автомобилей
func (s *MarketplaceService) GetTotalCarModelsCount(ctx context.Context) (int, error) {
	return s.storage.GetTotalCarModelsCount(ctx)
}
