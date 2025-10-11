// backend/internal/proj/unified/handler/routes.go
package handler

import (
	"backend/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

// RegisterRoutes регистрирует роуты для unified listings
func (h *UnifiedHandler) RegisterRoutes(app *fiber.App, mw *middleware.Middleware) error {
	// Публичные эндпоинты (без аутентификации)
	unified := app.Group("/api/v1/unified")
	unified.Get("/listings", h.GetUnifiedListings)
	unified.Get("/listings/:id", h.GetUnifiedListingByID)

	h.log.Info().Msg("Unified routes registered successfully")
	return nil
}

// GetPrefix возвращает префикс проекта для логирования
func (h *UnifiedHandler) GetPrefix() string {
	return "/api/v1/unified"
}
