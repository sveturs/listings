// backend/internal/proj/c2c/service/marketplace_base.go
package service

import (
	"context"
	"fmt"

	"backend/internal/config"
	"backend/internal/domain/models"
	"backend/internal/domain/search"
	"backend/internal/logger"
	"backend/internal/storage"

	"github.com/rs/zerolog"
)

const (
	// Attribute names
	attributeNameModel = "model"

	// Attribute types
	attributeTypeText = "text"

	// Languages
	languageAuto = "auto"

	// Field names
	fieldNameName = "name"

	// SQL queries
	insertTranslationQuery = `
        INSERT INTO translations (
            entity_type, entity_id, language, field_name,
            translated_text, is_machine_translated, is_verified, metadata,
            last_modified_by
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
        ON CONFLICT (entity_type, entity_id, language, field_name)
        DO UPDATE SET
            translated_text = EXCLUDED.translated_text,
            is_machine_translated = EXCLUDED.is_machine_translated,
            is_verified = EXCLUDED.is_verified,
            metadata = EXCLUDED.metadata,
            last_modified_by = EXCLUDED.last_modified_by,
            updated_at = CURRENT_TIMESTAMP
    `
)

type MarketplaceService struct {
	storage            storage.Storage
	translationService TranslationServiceInterface
	OrderService       OrderServiceInterface
	searchWeights      *config.SearchWeights
	cache              CacheInterface
	logger             zerolog.Logger
}

func NewMarketplaceService(storage storage.Storage, translationService TranslationServiceInterface, searchWeights *config.SearchWeights, cache CacheInterface) MarketplaceServiceInterface {
	ms := &MarketplaceService{
		storage:            storage,
		translationService: translationService,
		searchWeights:      searchWeights,
		cache:              cache,
		logger:             logger.Get().With().Str("service", "marketplace").Logger(),
	}

	// Создаем сервис заказов напрямую, избегая циклической зависимости
	ms.OrderService = NewSimpleOrderService(storage)

	return ms
}

// SetTranslationService allows injecting a translation service after creation
func (s *MarketplaceService) SetTranslationService(svc TranslationServiceInterface) {
	s.translationService = svc
}

func (s *MarketplaceService) GetOpenSearchRepository() (interface {
	Index(ctx context.Context, listing *models.MarketplaceListing) error
	Delete(ctx context.Context, listingID int) error
	SearchListings(ctx context.Context, params *search.ServiceParams) (*search.ServiceResult, error)
	GetSuggestions(ctx context.Context, query string, size int) ([]string, error)
	IndexAll(ctx context.Context, listings []*models.MarketplaceListing) error
	Exists(ctx context.Context, listingID int) (bool, error)
}, error,
) {
	repo, ok := s.storage.(interface {
		GetOpenSearchRepository() (interface {
			Index(ctx context.Context, listing *models.MarketplaceListing) error
			Delete(ctx context.Context, listingID int) error
			SearchListings(ctx context.Context, params *search.ServiceParams) (*search.ServiceResult, error)
			GetSuggestions(ctx context.Context, query string, size int) ([]string, error)
			IndexAll(ctx context.Context, listings []*models.MarketplaceListing) error
			Exists(ctx context.Context, listingID int) (bool, error)
		}, error)
	})
	if !ok {
		return nil, fmt.Errorf("storage does not implement GetOpenSearchRepository")
	}
	return repo.GetOpenSearchRepository()
}

func (s *MarketplaceService) Storage() storage.Storage {
	return s.storage
}

func (s *MarketplaceService) Service() *Service {
	// This method provides access to the full service, but we only have marketplace service here
	// Return a minimal Service struct with just the Marketplace field populated
	return &Service{
		Marketplace: s,
	}
}

func (s *MarketplaceService) SaveSearchQuery(ctx context.Context, query string, resultsCount int, language string) error {
	return s.storage.SaveSearchQuery(ctx, "", query, resultsCount, language)
}
