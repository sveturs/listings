package handler

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"

	"github.com/gofiber/fiber/v2"

	"backend/internal/proj/payments/service"
	"backend/pkg/logger"
	"backend/pkg/utils"
)

// WebhookHandler обрабатывает webhook'и от AllSecure
type WebhookHandler struct {
	service       *service.AllSecureService
	webhookSecret string
	logger        logger.Logger
}

// NewWebhookHandler создает новый webhook handler
func NewWebhookHandler(service *service.AllSecureService, webhookSecret string, logger logger.Logger) *WebhookHandler {
	return &WebhookHandler{
		service:       service,
		webhookSecret: webhookSecret,
		logger:        logger,
	}
}

// HandleAllSecureWebhook обрабатывает webhook от AllSecure
// @Summary AllSecure Webhook Handler
// @Description Обработка уведомлений от AllSecure о статусе платежей
// @Tags payments,webhooks
// @Accept json
// @Produce json
// @Param X-Signature header string true "Webhook signature"
// @Param payload body object true "Webhook payload"
// @Success 200 {object} utils.SuccessResponseSwag "Webhook processed successfully"
// @Failure 400 {object} utils.ErrorResponseSwag "Invalid signature or payload"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Router /webhooks/allsecure [post]
func (h *WebhookHandler) HandleAllSecureWebhook(c *fiber.Ctx) error {
	// Получаем подпись из заголовка
	signature := c.Get("X-Signature")
	if signature == "" {
		h.logger.Error("Missing webhook signature")
		return utils.ErrorResponse(c, 400, "missing signature")
	}

	// Получаем тело запроса
	payload := c.Body()
	if len(payload) == 0 {
		h.logger.Error("Empty webhook payload")
		return utils.ErrorResponse(c, 400, "empty payload")
	}

	// Проверяем подпись
	if !h.verifySignature(payload, signature) {
		// Don't log the actual signature value for security
		h.logger.Error("Invalid webhook signature received")
		return utils.ErrorResponse(c, 400, "invalid signature")
	}

	// Логируем webhook без sensitive данных
	h.logger.Info("Received AllSecure webhook (payload_size: %d, signature_present: %v)",
		len(payload), signature != "")

	// Обрабатываем webhook
	err := h.service.ProcessWebhook(c.Context(), payload)
	if err != nil {
		// Don't log the actual payload which may contain sensitive payment data
		h.logger.Error("Failed to process AllSecure webhook: %v", err)
		return utils.ErrorResponse(c, 500, "webhook processing failed")
	}

	h.logger.Info("AllSecure webhook processed successfully")
	return utils.SuccessResponse(c, "webhook processed")
}

// verifySignature проверяет подпись webhook'а
func (h *WebhookHandler) verifySignature(payload []byte, signature string) bool {
	if h.webhookSecret == "" {
		// В sandbox режиме можем пропустить проверку подписи
		h.logger.Info("Webhook secret not configured, skipping signature verification")
		return true
	}

	expectedSignature := h.calculateSignature(payload)
	return hmac.Equal([]byte(signature), []byte(expectedSignature))
}

// calculateSignature вычисляет ожидаемую подпись
func (h *WebhookHandler) calculateSignature(payload []byte) string {
	mac := hmac.New(sha256.New, []byte(h.webhookSecret))
	mac.Write(payload)
	return hex.EncodeToString(mac.Sum(nil))
}
