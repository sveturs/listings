// backend/internal/storage/postgres/search_queries.go
package postgres

import (
	"context"

	"backend/internal/domain/models"
	marketplaceService "backend/internal/proj/marketplace/service"
)

// GetPopularSearchQueries возвращает популярные поисковые запросы
func (db *Database) GetPopularSearchQueries(ctx context.Context, query string, limit int) ([]interface{}, error) {
	queries, err := db.marketplaceDB.GetPopularSearchQueries(ctx, query, limit)
	if err != nil {
		return nil, err
	}

	// Конвертируем в []interface{}
	result := make([]interface{}, len(queries))
	for i, q := range queries {
		result[i] = q
	}

	return result, nil
}

// SaveSearchQuery сохраняет или обновляет поисковый запрос
func (db *Database) SaveSearchQuery(ctx context.Context, query, normalizedQuery string, resultsCount int, language string) error {
	return db.marketplaceDB.SaveSearchQuery(ctx, query, normalizedQuery, resultsCount, language)
}

// SearchCategories ищет категории по названию
func (db *Database) SearchCategories(ctx context.Context, query string, limit int) ([]models.MarketplaceCategory, error) {
	return db.marketplaceDB.SearchCategories(ctx, query, limit)
}

// GetPopularSearchQueriesTyped возвращает типизированные популярные поисковые запросы для сервисов
func (db *Database) GetPopularSearchQueriesTyped(ctx context.Context, query string, limit int) ([]marketplaceService.SearchQuery, error) {
	return db.marketplaceDB.GetPopularSearchQueries(ctx, query, limit)
}
