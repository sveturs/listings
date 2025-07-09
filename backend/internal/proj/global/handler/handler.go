// backend/internal/proj/global/handler/handler.go
package handler

import (
	"github.com/gofiber/fiber/v2"

	"backend/internal/config"
	"backend/internal/middleware"
	globalService "backend/internal/proj/global/service"
)

// Handler объединяет все глобальные обработчики
type Handler struct {
	UnifiedSearch *UnifiedSearchHandler
	service       globalService.ServicesInterface
	searchWeights *config.SearchWeights
}

// NewHandler создает новый глобальный обработчик
func NewHandler(services globalService.ServicesInterface, searchWeights *config.SearchWeights) *Handler {
	return &Handler{
		UnifiedSearch: NewUnifiedSearchHandler(services),
		service:       services,
		searchWeights: searchWeights,
	}
}

// GetPrefix возвращает префикс для глобальных API
func (h *Handler) GetPrefix() string {
	return "/api/v1"
}

// RegisterRoutes регистрирует все глобальные маршруты
func (h *Handler) RegisterRoutes(app *fiber.App, mw *middleware.Middleware) error {
	// Регистрируем унифицированный поиск напрямую в app,
	// чтобы избежать конфликтов с другими middleware
	app.Get("/api/v1/search", h.UnifiedSearch.UnifiedSearch)

	return nil
}
