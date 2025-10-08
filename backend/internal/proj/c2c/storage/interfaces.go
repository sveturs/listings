// backend/internal/proj/c2c/storage/interfaces.go
package storage

import (
	"context"

	"backend/internal/domain/models"
	"backend/internal/proj/c2c/service"
)

// MarketplaceStorageExtended расширенный интерфейс для marketplace storage
type MarketplaceStorageExtended interface {
	// Методы для работы с поисковыми запросами
	GetPopularSearchQueries(ctx context.Context, query string, limit int) ([]service.SearchQuery, error)
	SaveSearchQuery(ctx context.Context, query, normalizedQuery string, resultsCount int, language string) error
	SearchCategories(ctx context.Context, query string, limit int) ([]models.MarketplaceCategory, error)
}
