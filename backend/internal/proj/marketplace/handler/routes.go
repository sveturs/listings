// TEMPORARY: Will be moved to microservice
package handler

import (
	"time"

	"backend/internal/middleware"

	"github.com/gofiber/fiber/v2"
	authMiddleware "github.com/sveturs/auth/pkg/http/fiber/middleware"
)

// RegisterRoutes регистрирует маршруты marketplace модуля
func (h *Handler) RegisterRoutes(app *fiber.App, mw *middleware.Middleware) error {
	// Публичные эндпоинты (без аутентификации)
	app.Get("/api/v1/marketplace/categories", h.GetCategories)
	app.Get("/api/v1/marketplace/popular-categories", h.GetPopularCategories)
	app.Get("/api/v1/marketplace/categories/:id/attributes", h.GetCategoryAttributes)
	app.Get("/api/v1/marketplace/categories/:slug/variant-attributes", h.GetVariantAttributes)
	app.Get("/api/v1/marketplace/neighborhood-stats", h.GetNeighborhoodStats)

	// Storefronts (B2C) - публичные эндпоинты
	app.Get("/api/v1/b2c", h.GetStorefronts)
	app.Get("/api/v1/b2c/:slug", h.GetStorefrontBySlug)
	app.Get("/api/v1/b2c/:id/products", h.GetStorefrontProducts)

	// Storefronts CRUD (защищенные)
	// ВАЖНО: /storefronts/my должен быть ПЕРЕД /storefronts/:id чтобы не матчился как параметр
	app.Get("/api/v1/marketplace/storefronts/my", h.jwtParserMW, authMiddleware.RequireAuth(), h.GetMyStorefronts)
	app.Post("/api/v1/marketplace/storefronts", h.jwtParserMW, authMiddleware.RequireAuth(), h.CreateStorefront)

	// Защищенные эндпоинты (требуют аутентификацию)
	// ВАЖНО: НЕ используем Group - это создает middleware leak для публичных роутов!
	app.Get("/api/v1/marketplace/favorites", h.jwtParserMW, authMiddleware.RequireAuth(), h.GetFavorites)
	app.Post("/api/v1/marketplace/favorites/:id", h.jwtParserMW, authMiddleware.RequireAuth(), h.AddToFavorites)
	app.Delete("/api/v1/marketplace/favorites/:id", h.jwtParserMW, authMiddleware.RequireAuth(), h.RemoveFromFavorites)
	app.Get("/api/v1/marketplace/chat", h.jwtParserMW, authMiddleware.RequireAuth(), h.GetChats)

	// Listings CRUD (TEMPORARY: direct DB until microservice migration complete)
	app.Post("/api/v1/marketplace/listings", h.jwtParserMW, authMiddleware.RequireAuth(), h.CreateListing)
	app.Get("/api/v1/marketplace/listings/:id", h.GetListing)
	app.Get("/api/v1/marketplace/listings/:id/similar", h.GetSimilarListings)
	app.Post("/api/v1/marketplace/listings/:id/images", h.jwtParserMW, authMiddleware.RequireAuth(), mw.RateLimitByUser(10, time.Minute), h.UploadListingImages)
	app.Delete("/api/v1/marketplace/listings/:id/images/:imageId", h.jwtParserMW, authMiddleware.RequireAuth(), h.DeleteListingImage)
	app.Patch("/api/v1/marketplace/listings/:id/images/reorder", h.jwtParserMW, authMiddleware.RequireAuth(), h.ReorderListingImages)

	return nil
}

// GetPrefix возвращает префикс проекта для логирования
func (h *Handler) GetPrefix() string {
	return "/api/v1/marketplace"
}
