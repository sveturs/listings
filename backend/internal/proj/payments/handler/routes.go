// backend/internal/proj/payments/handler/routes.go
package handler

import (
	"github.com/gofiber/fiber/v2"

	"backend/internal/middleware"
)

// RegisterRoutes регистрирует все маршруты для проекта payments
func (h *Handler) RegisterRoutes(app *fiber.App, mw *middleware.Middleware) error {
	// Webhook маршрут для обработки платежей (без аутентификации)
	app.Post("/webhook/stripe", h.HandleWebhook)

	return nil
}

// GetPrefix возвращает префикс проекта для логирования
func (h *Handler) GetPrefix() string {
	return "/webhook"
}