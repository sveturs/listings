// Package handler
// backend/internal/proj/payments/handler/routes.go
package handler

import (
	"github.com/gofiber/fiber/v2"

	"backend/internal/middleware"
)

// RegisterRoutes регистрирует все маршруты для проекта payments
func (h *Handler) RegisterRoutes(app *fiber.App, mw *middleware.Middleware) error {
	// Webhook маршруты с rate limiting для внешних сервисов
	webhooks := app.Group("/api/v1", mw.WebhookRateLimit())
	webhooks.Post("/payments/stripe/webhook", h.HandleWebhook)
	if h.webhook != nil {
		webhooks.Post("/webhooks/allsecure", h.webhook.HandleAllSecureWebhook)
	}

	// AllSecure routes (с аутентификацией и rate limiting)
	if h.allsecure != nil {
		// Группа для обычных платежных операций
		authenticated := app.Group("/api/v1/payments",
			mw.AuthRequiredJWT,
			mw.PaymentAPIRateLimit())
		authenticated.Post("/create", h.allsecure.CreatePayment)
		authenticated.Get("/:id/status", h.allsecure.GetPaymentStatus)

		// Группа для критических операций с более строгим rate limiting
		criticalOps := app.Group("/api/v1/payments",
			mw.AuthRequiredJWT,
			mw.StrictPaymentRateLimit())
		criticalOps.Post("/:id/capture", h.allsecure.CapturePayment)
		criticalOps.Post("/:id/refund", h.allsecure.RefundPayment)
	}

	return nil
}

// GetPrefix возвращает префикс проекта для логирования
func (h *Handler) GetPrefix() string {
	return "/api/v1/payments"
}
