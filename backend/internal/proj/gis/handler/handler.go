package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"

	"backend/internal/middleware"
	"backend/internal/proj/gis/service"
)

// Handler основной обработчик GIS модуля
type Handler struct {
	spatialHandler *SpatialHandler
}

// NewHandler создает новый обработчик GIS модуля
func NewHandler(db *sqlx.DB) *Handler {
	spatialService := service.NewSpatialService(db)
	spatialHandler := NewSpatialHandler(spatialService)

	return &Handler{
		spatialHandler: spatialHandler,
	}
}

// RegisterRoutes регистрирует все маршруты GIS модуля
func (h *Handler) RegisterRoutes(app *fiber.App, authMiddleware *middleware.Middleware) error {
	// Группа маршрутов для GIS
	gis := app.Group("/api/v1/gis")

	// Публичные маршруты (не требуют авторизации)
	gis.Get("/search", h.spatialHandler.SearchListings)
	gis.Get("/nearby", h.spatialHandler.GetNearbyListings)
	gis.Get("/listings/:id/location", h.spatialHandler.GetListingLocation)

	// Защищенные маршруты (требуют авторизации)
	gis.Put("/listings/:id/location", authMiddleware.AuthRequiredJWT, h.spatialHandler.UpdateListingLocation)

	return nil
}

// GetPrefix возвращает префикс модуля для логирования
func (h *Handler) GetPrefix() string {
	return "gis"
}
