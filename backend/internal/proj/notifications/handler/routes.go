// Package handler
// backend/internal/proj/notifications/routes.go
package handler

import (
	"github.com/gofiber/fiber/v2"
	authMiddleware "github.com/sveturs/auth/pkg/http/fiber/middleware"

	"backend/internal/middleware"
)

// RegisterRoutes регистрирует все маршруты для проекта notifications
func (h *Handler) RegisterRoutes(app *fiber.App, mw *middleware.Middleware) error {
	// Публичные маршруты
	app.Post("/api/v1/notifications/telegram/webhook", h.HandleTelegramWebhook)
	app.Post("/api/v1/notifications/email/public", h.SendPublicEmail)

	// Защищенные маршруты
	protected := app.Group("/api/v1/notifications", mw.JWTParser(), authMiddleware.RequireAuth())

	protected.Get("/", h.GetNotifications)
	protected.Get("/settings", h.GetSettings)
	protected.Put("/settings", h.UpdateSettings)
	protected.Get("/telegram/status", h.GetTelegramStatus)
	protected.Get("/telegram/token", h.GetTelegramToken)
	protected.Post("/telegram/connect", h.ConnectTelegram)
	protected.Put("/:id/read", h.MarkAsRead)

	return nil
}

// GetPrefix возвращает префикс проекта для логирования
func (h *Handler) GetPrefix() string {
	return "/api/v1/notifications"
}
