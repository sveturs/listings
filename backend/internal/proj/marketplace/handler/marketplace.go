package handler

import (
	"encoding/json"
	"net/http"

	"backend/internal/logger"
	"backend/internal/proj/marketplace/service"
	"backend/internal/storage/postgres"
)

// MarketplaceHandler represents the main handler for marketplace operations
type MarketplaceHandler struct {
	storage *postgres.Storage
	service service.MarketplaceServiceInterface
}

// NewMarketplaceHandler creates a new marketplace handler
func NewMarketplaceHandler(storage *postgres.Storage) *MarketplaceHandler {
	return &MarketplaceHandler{
		storage: storage,
	}
}

// respondWithJSON sends a JSON response
func (h *MarketplaceHandler) respondWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		logger.Error().Err(err).Msg("JSON serialization error")
		h.respondWithError(w, http.StatusInternalServerError, "marketplace.serializationError")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if _, err := w.Write(response); err != nil {
		h.logger.Error().Err(err).Msg("Failed to write response")
	}
}

// respondWithError sends an error response in JSON format
func (h *MarketplaceHandler) respondWithError(w http.ResponseWriter, statusCode int, message string) {
	h.respondWithJSON(w, statusCode, ErrorResponse{
		Success: false,
		Error:   message,
	})
}
