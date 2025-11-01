// backend/internal/proj/c2c/storage/opensearch/repository.go
package opensearch

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"backend/internal/config"
	"backend/internal/logger"
	"backend/internal/storage"
	osClient "backend/internal/storage/opensearch"
	"backend/pkg/transliteration"
)

const (
	// Field names
	fieldNamePrice     = "price"
	fieldNameCreatedAt = "created_at"

	// Boolean values
	boolValueTrue = "true"

	// Sort orders
	sortOrderDesc = "desc"
	sortOrderAsc  = "asc"
)

// Repository реализует интерфейс MarketplaceSearchRepository
type Repository struct {
	client         *osClient.OpenSearchClient
	indexName      string
	storage        storage.Storage
	transliterator *transliteration.SerbianTransliterator
	boostWeights   *config.OpenSearchBoostWeights
	asyncIndexer   *AsyncIndexer // Async indexer for non-blocking indexing
	useAsync       bool          // Flag to enable/disable async indexing
}

// NewRepository создает новый репозиторий
func NewRepository(client *osClient.OpenSearchClient, indexName string, storage storage.Storage, searchWeights *config.SearchWeights) *Repository {
	var boostWeights *config.OpenSearchBoostWeights
	if searchWeights != nil {
		boostWeights = &searchWeights.OpenSearchBoosts
	}

	repo := &Repository{
		client:         client,
		indexName:      indexName,
		storage:        storage,
		transliterator: transliteration.NewSerbianTransliterator(),
		boostWeights:   boostWeights,
		useAsync:       false, // По умолчанию выключен (для обратной совместимости)
	}

	return repo
}

// EnableAsyncIndexing включает асинхронную индексацию с указанными параметрами
func (r *Repository) EnableAsyncIndexing(db interface{}, workers int, queueSize int) error {
	// Проверяем тип db - должен быть *sqlx.DB
	sqlxDB, ok := db.(*sqlx.DB)
	if !ok {
		return fmt.Errorf("db must be *sqlx.DB, got %T", db)
	}

	r.asyncIndexer = NewAsyncIndexer(r, sqlxDB, workers, queueSize)
	r.useAsync = true

	logger.Info().
		Int("workers", workers).
		Int("queueSize", queueSize).
		Msg("Async indexing enabled for repository")

	return nil
}

// DisableAsyncIndexing выключает асинхронную индексацию
func (r *Repository) DisableAsyncIndexing() {
	r.useAsync = false
	logger.Info().Msg("Async indexing disabled for repository")
}

// ShutdownAsyncIndexer gracefully останавливает async indexer
func (r *Repository) ShutdownAsyncIndexer(timeout time.Duration) error {
	if r.asyncIndexer != nil {
		return r.asyncIndexer.Shutdown(timeout)
	}
	return nil
}

// GetAsyncIndexer возвращает async indexer (для тестов и мониторинга)
func (r *Repository) GetAsyncIndexer() *AsyncIndexer {
	return r.asyncIndexer
}

func (r *Repository) GetClient() *osClient.OpenSearchClient {
	return r.client
}

// DBTranslation представляет перевод из БД
type DBTranslation struct {
	Language       string `json:"language"`
	FieldName      string `json:"field_name"`
	TranslatedText string `json:"translated_text"`
}

// SimilarListing представляет похожее объявление с оценкой схожести
type SimilarListing struct {
	ID         int32   `json:"id"`
	CategoryID int32   `json:"category_id"`
	Title      string  `json:"title"`
	Score      float64 `json:"score"`
}
