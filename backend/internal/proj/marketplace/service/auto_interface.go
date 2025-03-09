// backend/internal/proj/marketplace/service/auto_interface.go
package service

import (
	"backend/internal/domain/models"
	"context"
)

// AutoServiceInterface определяет методы для работы с автомобильными объявлениями
type AutoServiceInterface interface {
	// Создание автомобильного объявления
	CreateAutoListing(ctx context.Context, listing *models.MarketplaceListing, autoProps *models.AutoProperties) (int, error)
	
	// Получение автомобильного объявления по ID
	GetAutoListingByID(ctx context.Context, id int) (*models.AutoListing, error)
	
	// Обновление автомобильного объявления
	UpdateAutoListing(ctx context.Context, listing *models.MarketplaceListing, autoProps *models.AutoProperties) error
	
	// Получение списка автомобильных объявлений с фильтрацией
	GetAutoListings(ctx context.Context, filters map[string]string, limit int, offset int) ([]models.AutoListing, int64, error)
	
	// Расширенный поиск автомобилей
	SearchAutoListings(ctx context.Context, autoFilters *models.AutoFilter, baseFilters map[string]string, limit int, offset int) ([]models.AutoListing, int64, error)
	
	// Получение констант для автомобильных свойств
	GetAutoConstants() models.AutoConstants
	
	// Проверка, является ли объявление автомобильным
	IsAutoCategory(ctx context.Context, categoryID int) (bool, error)
	
	// Получение всех популярных моделей для указанной марки
	GetModelsByBrand(ctx context.Context, brand string) ([]string, error)
	
	// Получение списка доступных марок автомобилей
	GetAvailableBrands(ctx context.Context) ([]string, error)
}