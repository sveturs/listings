// backend/internal/proj/c2c/storage/opensearch/repository.go
package opensearch

import (
	"backend/internal/config"
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
}

// NewRepository создает новый репозиторий
func NewRepository(client *osClient.OpenSearchClient, indexName string, storage storage.Storage, searchWeights *config.SearchWeights) *Repository {
	var boostWeights *config.OpenSearchBoostWeights
	if searchWeights != nil {
		boostWeights = &searchWeights.OpenSearchBoosts
	}

	return &Repository{
		client:         client,
		indexName:      indexName,
		storage:        storage,
		transliterator: transliteration.NewSerbianTransliterator(),
		boostWeights:   boostWeights,
	}
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
