package routes

import (
	"backend/internal/middleware"
	"backend/internal/proj/analytics/handler"

	"github.com/gofiber/fiber/v2"
)

// RegisterAnalyticsRoutes регистрирует маршруты для аналитики
func RegisterAnalyticsRoutes(app *fiber.App, h *handler.AnalyticsHandler, authMiddleware *middleware.Middleware) {
	api := app.Group("/api/v1/analytics")

	// Публичные маршруты (не требуют авторизации)
	public := api.Group("")
	{
		// Запись событий (может быть анонимной)
		public.Post("/event", h.RecordEvent)
	}

	// Защищенные маршруты (требуют авторизации)
	protected := api.Group("", authMiddleware.AuthRequiredJWT)
	{
		// Метрики поиска (только для админов)
		protected.Get("/metrics/search", h.GetSearchMetrics)
		// Производительность товаров (только для админов)
		protected.Get("/metrics/items", h.GetItemsPerformance)
	}
}
