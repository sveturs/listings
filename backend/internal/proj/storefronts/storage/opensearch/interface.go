package opensearch

import (
	"context"

	"backend/internal/domain/models"
)

// StorefrontSearchRepository интерфейс для работы с поиском витрин в OpenSearch
type StorefrontSearchRepository interface {
	// PrepareIndex подготавливает индекс (создает если не существует)
	PrepareIndex(ctx context.Context) error

	// Index индексирует одну витрину
	Index(ctx context.Context, storefront *models.Storefront) error

	// BulkIndex индексирует несколько витрин
	BulkIndex(ctx context.Context, storefronts []*models.Storefront) error

	// Delete удаляет витрину из индекса
	Delete(ctx context.Context, storefrontID int) error

	// Search выполняет поиск витрин
	Search(ctx context.Context, params *StorefrontSearchParams) (*StorefrontSearchResult, error)

	// ReindexAll переиндексирует все витрины
	ReindexAll(ctx context.Context) error
}
