package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"

	"backend/internal/middleware"
	"backend/internal/proj/gis/service"
)

// RegisterRoutes регистрация маршрутов GIS модуля
func RegisterRoutes(app *fiber.App, db *sqlx.DB, authMiddleware *middleware.Middleware) {
	// Создаем сервисы и обработчики
	spatialService := service.NewSpatialService(db)
	geocodingService := service.NewGeocodingService(db)

	spatialHandler := NewSpatialHandler(spatialService)
	geocodingHandler := NewGeocodingHandler(geocodingService)
	clusterHandler := NewClusterHandler(db)

	// Группа маршрутов для GIS
	gis := app.Group("/api/v1/gis")

	// ========== Публичные маршруты геокодирования (Phase 2) ==========
	gis.Post("/geocode/validate", geocodingHandler.ValidateGeocode)
	gis.Get("/geocode/suggestions", geocodingHandler.SearchAddressSuggestions)
	gis.Get("/geocode/reverse", geocodingHandler.ReverseGeocode)
	gis.Post("/geocode/multilingual", geocodingHandler.MultilingualReverseGeocode)
	gis.Get("/geocode/cache/stats", geocodingHandler.GetCacheStats)

	// ========== Публичные маршруты поиска ==========
	gis.Get("/search", spatialHandler.SearchListings)
	gis.Get("/search/radius", spatialHandler.RadiusSearch)
	gis.Get("/nearby", spatialHandler.GetNearbyListings)
	gis.Get("/listings/:id/location", spatialHandler.GetListingLocation)

	// ========== Маршруты кластеризации и визуализации ==========
	gis.Get("/clusters", clusterHandler.GetClusters)
	gis.Get("/heatmap", clusterHandler.GetHeatmap)

	// ========== Защищенные маршруты (требуют авторизации) ==========
	protected := gis.Group("/", authMiddleware.AuthRequiredJWT)

	// Старые endpoints
	protected.Put("/listings/:id/location", spatialHandler.UpdateListingLocation)

	// Новые endpoints Phase 2
	protected.Put("/listings/:id/address", spatialHandler.UpdateListingAddress)
	protected.Post("/geocode/cache/cleanup", geocodingHandler.CleanupExpiredCache)
}

// RegisterPublicRoutes регистрация только публичных маршрутов
func RegisterPublicRoutes(router fiber.Router, spatialHandler *SpatialHandler) {
	gis := router.Group("/gis")

	gis.Get("/search", spatialHandler.SearchListings)
	gis.Get("/search/radius", spatialHandler.RadiusSearch)
	gis.Get("/nearby", spatialHandler.GetNearbyListings)
	gis.Get("/listings/:id/location", spatialHandler.GetListingLocation)
}

// RegisterProtectedRoutes регистрация защищенных маршрутов
func RegisterProtectedRoutes(router fiber.Router, spatialHandler *SpatialHandler) {
	gis := router.Group("/gis")

	gis.Put("/listings/:id/location", spatialHandler.UpdateListingLocation)
}
