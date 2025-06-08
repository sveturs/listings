// backend/internal/proj/payments/handler/handler.go
package handler

import (
	"log"

	"github.com/gofiber/fiber/v2"

	globalService "backend/internal/proj/global/service"
	paymentService "backend/internal/proj/payments/service"
)

type Handler struct {
	payment paymentService.PaymentServiceInterface
}

func NewHandler(services globalService.ServicesInterface) *Handler {
	return &Handler{
		payment: services.Payment(),
	}
}

// HandleWebhook обрабатывает webhook от платежной системы (Stripe)
// @Summary Webhook платежной системы
// @Description Обрабатывает webhook уведомления от Stripe о статусе платежей
// @Tags payments
// @Accept json
// @Produce json
// @Param Stripe-Signature header string true "Подпись Stripe для верификации webhook"
// @Param payload body string true "Тело запроса от Stripe"
// @Success 200 {string} string "OK"
// @Failure 400 {string} string "Bad Request"
// @Router /api/v1/payments/stripe/webhook [post]
func (h *Handler) HandleWebhook(c *fiber.Ctx) error {
	payload := c.Body()
	signature := c.Get("Stripe-Signature")

	err := h.payment.HandleWebhook(c.Context(), payload, signature)
	if err != nil {
		log.Printf("Webhook error: %v", err)
		return c.SendStatus(fiber.StatusBadRequest)
	}

	return c.SendStatus(fiber.StatusOK)
}