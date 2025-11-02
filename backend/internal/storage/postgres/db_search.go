// backend/internal/storage/postgres/db_search.go
package postgres

import (
	"context"
	"fmt"

	"backend/internal/config"
	"backend/internal/domain/search"
)

// GetSearchWeights возвращает веса для поиска
func (db *Database) GetSearchWeights() *config.SearchWeights {
	return db.searchWeights
}

// GetOpenSearchClient возвращает клиент OpenSearch для прямого выполнения запросов
func (db *Database) GetOpenSearchClient() (interface {
	Execute(ctx context.Context, method, path string, body []byte) ([]byte, error)
}, error,
) {
	if db.osClient == nil {
		return nil, fmt.Errorf("OpenSearch клиент не настроен")
	}
	return db.osClient, nil
}

// PrepareIndex подготавливает индекс OpenSearch
func (db *Database) PrepareIndex(ctx context.Context) error {
	if true { // OpenSearch disabled after removing c2c
		// Если репозиторий OpenSearch не инициализирован, просто возвращаем nil
		// Поиск будет работать без OpenSearch
		return nil
	}

	// TODO: OpenSearch temporarily disabled during refactoring
	return nil
}

// SearchListingsOpenSearch - TODO: temporarily disabled during refactoring
func (db *Database) SearchListingsOpenSearch(ctx context.Context, params *search.SearchParams) (*search.SearchResult, error) {
	return nil, fmt.Errorf("OpenSearch temporarily disabled")
}
