package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"

	"backend/internal/middleware"
	"backend/internal/proj/gis/service"
)

// RegisterRoutes регистрация маршрутов GIS модуля
func RegisterRoutes(app *fiber.App, db *sqlx.DB, authMiddleware *middleware.Middleware) {
	// Создаем сервис и обработчик
	spatialService := service.NewSpatialService(db)
	spatialHandler := NewSpatialHandler(spatialService)

	// Группа маршрутов для GIS
	gis := app.Group("/api/v1/gis")

	// Публичные маршруты (не требуют авторизации)
	gis.Get("/search", spatialHandler.SearchListings)
	gis.Get("/clusters", spatialHandler.GetClusters)
	gis.Get("/nearby", spatialHandler.GetNearbyListings)
	gis.Get("/listings/:id/location", spatialHandler.GetListingLocation)

	// Защищенные маршруты (требуют авторизации)
	gis.Put("/listings/:id/location", authMiddleware.AuthRequiredJWT, spatialHandler.UpdateListingLocation)
}

// RegisterPublicRoutes регистрация только публичных маршрутов
func RegisterPublicRoutes(router fiber.Router, spatialHandler *SpatialHandler) {
	gis := router.Group("/gis")

	gis.Get("/search", spatialHandler.SearchListings)
	gis.Get("/clusters", spatialHandler.GetClusters)
	gis.Get("/nearby", spatialHandler.GetNearbyListings)
	gis.Get("/listings/:id/location", spatialHandler.GetListingLocation)
}

// RegisterProtectedRoutes регистрация защищенных маршрутов
func RegisterProtectedRoutes(router fiber.Router, spatialHandler *SpatialHandler) {
	gis := router.Group("/gis")

	gis.Put("/listings/:id/location", spatialHandler.UpdateListingLocation)
}
