// Package handler
// backend/internal/proj/payments/handler/handler.go
package handler

import (
	"log"

	"github.com/gofiber/fiber/v2"

	globalService "backend/internal/proj/global/service"
	paymentService "backend/internal/proj/payments/service"
	"backend/pkg/logger"
	"backend/pkg/utils"
)

type Handler struct {
	payment   paymentService.PaymentServiceInterface
	allsecure *PaymentHandler // AllSecure payment handler
	webhook   *WebhookHandler // AllSecure webhook handler
}

func NewHandler(services globalService.ServicesInterface) *Handler {
	return &Handler{
		payment:   services.Payment(),
		allsecure: nil, // Будет инициализирован через InitAllSecure
		webhook:   nil, // Будет инициализирован через InitAllSecure
	}
}

// InitAllSecure инициализирует AllSecure handlers
// Вызывается после создания основного handler'а, когда доступна конфигурация
func (h *Handler) InitAllSecure(allsecureService *paymentService.AllSecureService) {
	if allsecureService == nil {
		log.Printf("Warning: AllSecure service is nil, skipping AllSecure handlers initialization")
		return
	}

	logger := logger.New()

	h.allsecure = NewPaymentHandler(allsecureService, *logger)
	h.webhook = NewWebhookHandler(allsecureService, "", *logger) // webhook secret из конфига

	log.Printf("AllSecure handlers initialized successfully")
}

// HandleWebhook processes webhook from payment system (Stripe)
// @Summary Process payment webhook
// @Description Processes webhook notifications from Stripe about payment status
// @Tags payments
// @Accept json
// @Produce json
// @Param Stripe-Signature header string true "Stripe signature for webhook verification"
// @Param payload body StripeWebhookRequest true "Stripe webhook payload"
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag{data=WebhookResponse} "Webhook processed successfully"
// @Failure 400 {object} backend_pkg_utils.ErrorResponseSwag "payments.webhook.invalid_signature or payments.webhook.processing_error"
// @Router /api/v1/payments/stripe/webhook [post]
func (h *Handler) HandleWebhook(c *fiber.Ctx) error {
	payload := c.Body()
	signature := c.Get("Stripe-Signature")

	err := h.payment.HandleWebhook(c.Context(), payload, signature)
	if err != nil {
		log.Printf("Webhook error: %v", err)
		// Check if it's a signature verification error
		if err.Error() == "invalid signature" {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "payments.webhook.invalid_signature")
		}
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "payments.webhook.processing_error")
	}

	response := WebhookResponse{
		Status:  "success",
		Message: "payments.webhook.processed",
	}
	return utils.SuccessResponse(c, response)
}
