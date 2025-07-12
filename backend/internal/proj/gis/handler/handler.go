package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"

	"backend/internal/middleware"
	"backend/internal/proj/gis/service"
)

// Handler основной обработчик GIS модуля
type Handler struct {
	spatialHandler   *SpatialHandler
	geocodingHandler *GeocodingHandler
}

// NewHandler создает новый обработчик GIS модуля
func NewHandler(db *sqlx.DB) *Handler {
	spatialService := service.NewSpatialService(db)
	geocodingService := service.NewGeocodingService(db)

	spatialHandler := NewSpatialHandler(spatialService)
	geocodingHandler := NewGeocodingHandler(geocodingService)

	return &Handler{
		spatialHandler:   spatialHandler,
		geocodingHandler: geocodingHandler,
	}
}

// RegisterRoutes регистрирует все маршруты GIS модуля
func (h *Handler) RegisterRoutes(app *fiber.App, authMiddleware *middleware.Middleware) error {
	// Группа маршрутов для GIS
	gis := app.Group("/api/v1/gis")

	// ========== Публичные маршруты геокодирования (Phase 2) ==========
	gis.Post("/geocode/validate", h.geocodingHandler.ValidateGeocode)
	gis.Get("/geocode/suggestions", h.geocodingHandler.SearchAddressSuggestions)
	gis.Get("/geocode/reverse", h.geocodingHandler.ReverseGeocode)
	gis.Get("/geocode/cache/stats", h.geocodingHandler.GetCacheStats)

	// ========== Публичные маршруты поиска ==========
	gis.Get("/search", h.spatialHandler.SearchListings)
	gis.Get("/search/radius", h.spatialHandler.RadiusSearch)
	gis.Get("/nearby", h.spatialHandler.GetNearbyListings)
	gis.Get("/listings/:id/location", h.spatialHandler.GetListingLocation)

	// ========== Защищенные маршруты (требуют авторизации) ==========
	protected := gis.Group("/", authMiddleware.AuthRequiredJWT)

	// Старые endpoints
	protected.Put("/listings/:id/location", h.spatialHandler.UpdateListingLocation)

	// Новые endpoints Phase 2
	protected.Put("/listings/:id/address", h.spatialHandler.UpdateListingAddress)
	protected.Post("/geocode/cache/cleanup", h.geocodingHandler.CleanupExpiredCache)

	return nil
}

// GetPrefix возвращает префикс модуля для логирования
func (h *Handler) GetPrefix() string {
	return "gis"
}
