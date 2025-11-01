// backend/internal/proj/unified/handler/routes.go
package handler

import (
	"backend/internal/middleware"

	"github.com/gofiber/fiber/v2"
	authMiddleware "github.com/sveturs/auth/pkg/http/fiber/middleware"
)

// RegisterMarketplaceRoutes регистрирует unified marketplace routes
func (h *MarketplaceHandler) RegisterMarketplaceRoutes(app *fiber.App, mw *middleware.Middleware, jwtParserMW fiber.Handler) error {
	// Middleware для добавления X-Served-By header (Phase 3 traffic measurement)
	servedByMiddleware := func(c *fiber.Ctx) error {
		// Устанавливаем default header (monolith)
		// Microservice может override через свой response
		c.Set("X-Served-By", "monolith")
		return c.Next()
	}

	// Public routes (без аутентификации)
	app.Get("/api/v1/marketplace/search", servedByMiddleware, h.SearchListings)
	app.Get("/api/v1/marketplace/listings/:id", servedByMiddleware, h.GetListing)

	// Protected routes (требуют аутентификацию)
	protected := app.Group("/api/v1/marketplace", servedByMiddleware, jwtParserMW, authMiddleware.RequireAuth())
	protected.Post("/listings", h.CreateListing)
	protected.Put("/listings/:id", h.UpdateListing)
	protected.Delete("/listings/:id", h.DeleteListing)

	h.logger.Info().Msg("Unified marketplace routes registered successfully")
	return nil
}

// RegisterRoutes регистрирует роуты для unified listings (legacy support)
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
