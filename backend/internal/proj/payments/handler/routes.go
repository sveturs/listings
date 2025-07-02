// Package handler
// backend/internal/proj/payments/handler/routes.go
package handler

import (
	"github.com/gofiber/fiber/v2"

	"backend/internal/middleware"
)

// RegisterRoutes регистрирует все маршруты для проекта payments
func (h *Handler) RegisterRoutes(app *fiber.App, mw *middleware.Middleware) error {
	// Stripe webhook маршрут для обработки платежей (без аутентификации)
	app.Post("/api/v1/payments/stripe/webhook", h.HandleWebhook)

	// AllSecure routes (с аутентификацией)
	// TODO: Добавить проверку что allsecure handler инициализирован
	if h.allsecure != nil {
		authenticated := app.Group("/api/v1/payments", mw.AuthRequiredJWT)
		authenticated.Post("/create", h.allsecure.CreatePayment)
		authenticated.Post("/:id/capture", h.allsecure.CapturePayment)
		authenticated.Post("/:id/refund", h.allsecure.RefundPayment)
		authenticated.Get("/:id/status", h.allsecure.GetPaymentStatus)
	}

	// AllSecure webhook (без аутентификации)
	if h.webhook != nil {
		app.Post("/api/v1/webhooks/allsecure", h.webhook.HandleAllSecureWebhook)
	}

	return nil
}

// GetPrefix возвращает префикс проекта для логирования
func (h *Handler) GetPrefix() string {
	return "/api/v1/payments"
}
