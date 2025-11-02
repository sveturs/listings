// backend/internal/storage/postgres/search_queries.go
package postgres

import (
	"context"
	"errors"

	"backend/internal/domain/models"
)

// GetPopularSearchQueries возвращает популярные поисковые запросы
// TODO: Migrate to marketplace microservice
func (db *Database) GetPopularSearchQueries(ctx context.Context, query string, limit int) ([]interface{}, error) {
	return nil, errors.New("marketplace service removed - use microservice")
}

// SaveSearchQuery сохраняет или обновляет поисковый запрос
// TODO: Migrate to marketplace microservice
func (db *Database) SaveSearchQuery(ctx context.Context, query, normalizedQuery string, resultsCount int, language string) error {
	return errors.New("marketplace service removed - use microservice")
}

// SearchCategories ищет категории по названию
// TODO: Migrate to marketplace microservice
func (db *Database) SearchCategories(ctx context.Context, query string, limit int) ([]models.MarketplaceCategory, error) {
	return nil, errors.New("marketplace service removed - use microservice")
}

// GetPopularSearchQueriesTyped возвращает типизированные популярные поисковые запросы для сервисов
// TODO: Migrate to marketplace microservice
func (db *Database) GetPopularSearchQueriesTyped(ctx context.Context, query string, limit int) ([]interface{}, error) {
	return nil, errors.New("marketplace service removed - use microservice")
}
