// backend/internal/proj/geocode/handler/routes.go
package handler

import (
	"github.com/gofiber/fiber/v2"

	"backend/internal/middleware"
)

// RegisterRoutes регистрирует все маршруты для проекта geocode
func (h *Handler) RegisterRoutes(app *fiber.App, mw *middleware.Middleware) error {
	// Публичные маршруты для геокодинга
	geocodeApi := app.Group("/api/v1/geocode")
	geocodeApi.Get("/reverse", h.Geocode.ReverseGeocode)

	// Публичные маршруты для городов
	citiesApi := app.Group("/api/v1/cities")
	citiesApi.Get("/suggest", h.Geocode.GetCitySuggestions)

	return nil
}

// GetPrefix возвращает префикс проекта для логирования
func (h *Handler) GetPrefix() string {
	return "/api/v1/geocode"
}