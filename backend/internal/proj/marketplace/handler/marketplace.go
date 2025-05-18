package handler

import (
	"encoding/json"
	"log"
	"net/http"

	postgres "backend/internal/storage/postgres"
)

// MarketplaceHandler представляет основной обработчик для маркетплейса
type MarketplaceHandler struct {
	storage *postgres.Storage
}

// NewMarketplaceHandler создает новый основной обработчик маркетплейса
func NewMarketplaceHandler(storage *postgres.Storage) *MarketplaceHandler {
	return &MarketplaceHandler{
		storage: storage,
	}
}

// respondWithJSON отправляет JSON-ответ
func (h *MarketplaceHandler) respondWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Ошибка сериализации в JSON: %v", err)
		h.respondWithError(w, http.StatusInternalServerError, "Ошибка обработки ответа")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(response)
}

// respondWithError отправляет ошибку в формате JSON
func (h *MarketplaceHandler) respondWithError(w http.ResponseWriter, statusCode int, message string) {
	h.respondWithJSON(w, statusCode, map[string]string{"error": message})
}