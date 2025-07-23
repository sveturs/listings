package handler

import (
	"backend/internal/logger"
	"backend/internal/proj/marketplace/service"
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
func NewMarketplaceHandler(storage *postgres.Storage) *MarketplaceHandler {
	return &MarketplaceHandler{
		storage: storage,
		logger:  logger.Get().With().Str("handler", "marketplace").Logger(),
	}
}


