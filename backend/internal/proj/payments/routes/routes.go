package routes

import (
	"github.com/gofiber/fiber/v2"

	"backend/internal/proj/payments/handler"
)

// RegisterPaymentRoutes регистрирует маршруты для платежей
func RegisterPaymentRoutes(app *fiber.App, paymentHandler *handler.PaymentHandler, webhookHandler *handler.WebhookHandler) {
	// Публичные webhook endpoints (без авторизации)
	webhooks := app.Group("/api/v1/webhooks")
	webhooks.Post("/allsecure", webhookHandler.HandleAllSecureWebhook)

	// Защищенные endpoints для платежей (требуют авторизации)
	payments := app.Group("/api/v1/payments")
	// Здесь должен быть middleware для авторизации
	// payments.Use(middleware.JWTAuth())

	payments.Post("/create", paymentHandler.CreatePayment)
	payments.Post("/:id/capture", paymentHandler.CapturePayment)
	payments.Post("/:id/refund", paymentHandler.RefundPayment)
	payments.Get("/:id/status", paymentHandler.GetPaymentStatus)
}
