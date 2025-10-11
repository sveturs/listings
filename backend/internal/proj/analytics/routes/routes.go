package routes

import (
	"backend/internal/middleware"
	"backend/internal/proj/analytics/handler"

	"github.com/gofiber/fiber/v2"
	authMiddleware "github.com/sveturs/auth/pkg/http/fiber/middleware"
)

// RegisterAnalyticsRoutes регистрирует маршруты для аналитики
func RegisterAnalyticsRoutes(app *fiber.App, h *handler.AnalyticsHandler, mw *middleware.Middleware, jwtParserMW fiber.Handler) {
	// Публичный маршрут - регистрируем напрямую БЕЗ группы чтобы избежать наследования middleware
	// Запись событий (может быть анонимной)
	app.Post("/api/v1/analytics/event", h.RecordEvent)

	// Защищенные admin маршруты (требуют авторизации и роли admin)
	// БЕЗ CSRF - используем BFF proxy архитектуру
	adminMetrics := app.Group("/api/v1/analytics/metrics", jwtParserMW, authMiddleware.RequireAuthString("admin"))
	{
		// Метрики поиска (только для админов)
		adminMetrics.Get("/search", h.GetSearchMetrics)
		// Производительность товаров (только для админов)
		adminMetrics.Get("/items", h.GetItemsPerformance)
	}
}
