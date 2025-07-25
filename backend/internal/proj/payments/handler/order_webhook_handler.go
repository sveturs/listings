package handler

import (
	"github.com/gofiber/fiber/v2"

	"backend/internal/proj/payments/service"
	"backend/pkg/logger"
	"backend/pkg/utils"
)

// OrderWebhookHandler обрабатывает webhooks для платежей заказов
type OrderWebhookHandler struct {
	service service.PaymentServiceInterface
	logger  logger.Logger
}

// NewOrderWebhookHandler создает новый order webhook handler
func NewOrderWebhookHandler(service service.PaymentServiceInterface, logger logger.Logger) *OrderWebhookHandler {
	return &OrderWebhookHandler{
		service: service,
		logger:  logger,
	}
}

// HandleOrderPaymentWebhook обрабатывает webhook для платежа заказа
// @Summary Handle Order Payment Webhook
// @Description Обрабатывает webhook уведомления о статусе платежа заказа
// @Tags order-payments
// @Accept json
// @Produce json
// @Param signature header string false "Webhook signature"
// @Param payload body object true "Webhook payload"
// @Success 200 {object} utils.SuccessResponseSwag "Webhook processed successfully"
// @Failure 400 {object} utils.ErrorResponseSwag "Invalid webhook payload"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Router /api/v1/webhooks/orders/payment [post]
func (h *OrderWebhookHandler) HandleOrderPaymentWebhook(c *fiber.Ctx) error {
	// Получаем сигнатуру из заголовков
	signature := c.Get("X-Webhook-Signature")
	if signature == "" {
		signature = c.Get("X-AllSecure-Signature") // Альтернативное имя заголовка
	}

	// Получаем тело запроса
	payload := c.Body()
	if len(payload) == 0 {
		h.logger.Error("Empty webhook payload received")
		return utils.ErrorResponse(c, 400, "empty payload")
	}

	h.logger.Info("Order payment webhook received (signature: %s, payloadSize: %d)", signature, len(payload))

	// Обрабатываем webhook
	err := h.service.HandleOrderPaymentWebhook(c.Context(), payload, signature)
	if err != nil {
		h.logger.Error("Failed to process order payment webhook: %v", err)
		return utils.ErrorResponse(c, 500, "webhook processing failed")
	}

	h.logger.Info("Order payment webhook processed successfully")
	return utils.SuccessResponse(c, "webhook processed")
}

// WebhookResponse представляет ответ webhook
type WebhookResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}
