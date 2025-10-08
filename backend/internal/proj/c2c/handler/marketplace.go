package handler

import (
	"backend/internal/logger"
	"backend/internal/proj/c2c/service"
	"backend/internal/storage/postgres"

	"github.com/rs/zerolog"
)

// MarketplaceHandler represents the main handler for marketplace operations
type MarketplaceHandler struct {
	storage *postgres.Storage
	service service.MarketplaceServiceInterface
	logger  zerolog.Logger
}

// NewMarketplaceHandler creates a new marketplace handler
func NewMarketplaceHandler(storage *postgres.Storage, marketplaceService service.MarketplaceServiceInterface) *MarketplaceHandler {
	return &MarketplaceHandler{
		storage: storage,
		service: marketplaceService,
		logger:  logger.Get().With().Str("handler", "marketplace").Logger(),
	}
}
