// backend/internal/proj/marketplace/service/auto.go
package service

import (
	"backend/internal/domain/models"
	"backend/internal/storage"
	"context"
	"fmt"
	"log"
)

// AutoService реализует интерфейс AutoServiceInterface
type AutoService struct {
	storage storage.Storage
}

// NewAutoService создает новый экземпляр сервиса автомобилей
func NewAutoService(storage storage.Storage) AutoServiceInterface {
	return &AutoService{
		storage: storage,
	}
}

// CreateAutoListing создает новое автомобильное объявление
func (s *AutoService) CreateAutoListing(ctx context.Context, listing *models.MarketplaceListing, autoProps *models.AutoProperties) (int, error) {
	// Проверяем, что категория относится к автомобилям
	isAuto, err := s.IsAutoCategory(ctx, listing.CategoryID)
	if err != nil {
		return 0, fmt.Errorf("ошибка проверки категории: %w", err)
	}

	if !isAuto {
		return 0, fmt.Errorf("категория %d не является автомобильной", listing.CategoryID)
	}

	// Начинаем транзакцию
	tx, err := s.storage.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("ошибка начала транзакции: %w", err)
	}
	defer tx.Rollback()

	// Создаем базовое объявление
	listingID, err := s.storage.CreateListing(ctx, listing)
	if err != nil {
		return 0, fmt.Errorf("ошибка создания базового объявления: %w", err)
	}

	// Устанавливаем ID объявления для автомобильных свойств
	autoProps.ListingID = listingID

	// Создаем автомобильные свойства
	err = s.storage.CreateAutoProperties(ctx, autoProps)
	if err != nil {
		return 0, fmt.Errorf("ошибка создания автомобильных свойств: %w", err)
	}

	// Фиксируем транзакцию
	if err = tx.Commit(); err != nil {
		return 0, fmt.Errorf("ошибка фиксации транзакции: %w", err)
	}

	return listingID, nil
}

// GetAutoListingByID получает автомобильное объявление по ID
func (s *AutoService) GetAutoListingByID(ctx context.Context, id int) (*models.AutoListing, error) {
	return s.storage.GetAutoListingByID(ctx, id)
}

// UpdateAutoListing обновляет автомобильное объявление
func (s *AutoService) UpdateAutoListing(ctx context.Context, listing *models.MarketplaceListing, autoProps *models.AutoProperties) error {
	// Проверяем, что категория относится к автомобилям
	isAuto, err := s.IsAutoCategory(ctx, listing.CategoryID)
	if err != nil {
		return fmt.Errorf("ошибка проверки категории: %w", err)
	}

	if !isAuto {
		return fmt.Errorf("категория %d не является автомобильной", listing.CategoryID)
	}

	// Начинаем транзакцию
	tx, err := s.storage.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("ошибка начала транзакции: %w", err)
	}
	defer tx.Rollback()

	// Обновляем базовое объявление
	err = s.storage.UpdateListing(ctx, listing)
	if err != nil {
		return fmt.Errorf("ошибка обновления базового объявления: %w", err)
	}

	// Устанавливаем ID объявления для автомобильных свойств
	autoProps.ListingID = listing.ID

	// Обновляем автомобильные свойства
	err = s.storage.UpdateAutoProperties(ctx, autoProps)
	if err != nil {
		return fmt.Errorf("ошибка обновления автомобильных свойств: %w", err)
	}

	// Фиксируем транзакцию
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("ошибка фиксации транзакции: %w", err)
	}

	return nil
}

// GetAutoListings получает список автомобильных объявлений с фильтрацией
func (s *AutoService) GetAutoListings(ctx context.Context, filters map[string]string, limit int, offset int) ([]models.AutoListing, int64, error) {
	return s.storage.GetAutoListings(ctx, filters, limit, offset)
}

// SearchAutoListings выполняет расширенный поиск автомобилей
func (s *AutoService) SearchAutoListings(ctx context.Context, autoFilters *models.AutoFilter, baseFilters map[string]string, limit int, offset int) ([]models.AutoListing, int64, error) {
	return s.storage.SearchAutoListings(ctx, autoFilters, baseFilters, limit, offset)
}

// GetAutoConstants возвращает константы для автомобильных свойств
func (s *AutoService) GetAutoConstants() models.AutoConstants {
	return models.GetAutoConstants()
}

// Исправленная версия функции IsAutoCategory
func (s *AutoService) IsAutoCategory(ctx context.Context, categoryID int) (bool, error) {
    // Список ID автомобильных категорий
    autoCategories := []int{2000, 2100, 2200, 2210, 2220, 2230, 2240, 2300, 2310, 2315, 2320, 2325, 2330, 2335, 2340, 2345, 2350, 2355, 2360, 2365}
    
    // Проверяем, есть ли ID в списке автомобильных категорий
    for _, id := range autoCategories {
        if id == categoryID {
            return true, nil
        }
    }
    
    // Запрос для проверки, является ли категория дочерней или внучатой для категории 2000 (Автомобили)
    query := `
        WITH RECURSIVE category_tree AS (
            -- Базовый случай: категория 2000 (Автомобили)
            SELECT id, parent_id FROM marketplace_categories WHERE id = 2000
            UNION ALL
            -- Рекурсивный случай: все дочерние категории
            SELECT c.id, c.parent_id
            FROM marketplace_categories c
            JOIN category_tree ct ON c.parent_id = ct.id
        )
        SELECT EXISTS(SELECT 1 FROM category_tree WHERE id = $1)
    `
    
    var isAuto bool
    err := s.storage.QueryRow(ctx, query, categoryID).Scan(&isAuto)
    if err != nil {
        return false, fmt.Errorf("ошибка проверки категории: %w", err)
    }
    
    return isAuto, nil
}
// GetModelsByBrand возвращает список моделей для указанной марки
func (s *AutoService) GetModelsByBrand(ctx context.Context, brand string) ([]string, error) {
	query := `
		SELECT DISTINCT model 
		FROM auto_properties 
		WHERE brand = $1 
		ORDER BY model
	`

	rows, err := s.storage.Query(ctx, query, brand)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения моделей: %w", err)
	}
	defer rows.Close()

	var models []string
	for rows.Next() {
		var model string
		if err := rows.Scan(&model); err != nil {
			log.Printf("Ошибка сканирования модели: %v", err)
			continue
		}
		models = append(models, model)
	}

	return models, nil
}

// GetAvailableBrands возвращает список доступных марок автомобилей
func (s *AutoService) GetAvailableBrands(ctx context.Context) ([]string, error) {
	query := `
		SELECT DISTINCT brand 
		FROM auto_properties 
		WHERE brand != '' 
		ORDER BY brand
	`

	rows, err := s.storage.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения марок: %w", err)
	}
	defer rows.Close()

	var brands []string
	for rows.Next() {
		var brand string
		if err := rows.Scan(&brand); err != nil {
			log.Printf("Ошибка сканирования марки: %v", err)
			continue
		}
		brands = append(brands, brand)
	}

	return brands, nil
}