// backend/internal/proj/notifications/routes.go
package handler

import (
	"github.com/gofiber/fiber/v2"

	"backend/internal/middleware"
)

// RegisterRoutes регистрирует все маршруты для проекта notifications
func (h *Handler) RegisterRoutes(app *fiber.App, mw *middleware.Middleware) error {
	// Публичные маршруты
	app.Post("/api/v1/notifications/telegram/webhook", h.HandleTelegramWebhook)
	app.Get("/api/v1/notifications/telegram", h.GetTelegramStatus)
	app.Post("/api/v1/public/send-email", h.SendPublicEmail)

	// Защищенные маршруты
	protected := app.Group("/api/v1/notifications", mw.AuthRequiredJWT)

	protected.Get("/", h.GetNotifications)
	protected.Get("/settings", h.GetSettings)
	protected.Put("/settings", h.UpdateSettings)
	protected.Get("/telegram", h.GetTelegramStatus)
	protected.Get("/telegram/token", h.GetTelegramToken)
	protected.Put("/:id/read", h.MarkAsRead)
	protected.Post("/telegram/token", h.GetTelegramToken)

	return nil
}

// GetPrefix возвращает префикс проекта для логирования
func (h *Handler) GetPrefix() string {
	return "/api/v1/notifications"
}
