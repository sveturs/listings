package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"backend/internal/proj/marketplace/service"
	postgres "backend/internal/storage/postgres"
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
		log.Printf("JSON serialization error: %v", err)
		h.respondWithError(w, http.StatusInternalServerError, "marketplace.serializationError")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(response)
}

// respondWithError sends an error response in JSON format
func (h *MarketplaceHandler) respondWithError(w http.ResponseWriter, statusCode int, message string) {
	h.respondWithJSON(w, statusCode, map[string]string{"error": message})
}
