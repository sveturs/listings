package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/shopspring/decimal"

	"backend/internal/domain/models"
	"backend/internal/proj/payments/service"
	"backend/pkg/logger"
	"backend/pkg/utils"
)

// OrderPaymentHandler обрабатывает HTTP запросы для платежей заказов
type OrderPaymentHandler struct {
	service service.PaymentServiceInterface
	logger  logger.Logger
}

// NewOrderPaymentHandler создает новый order payment handler
func NewOrderPaymentHandler(service service.PaymentServiceInterface, logger logger.Logger) *OrderPaymentHandler {
	return &OrderPaymentHandler{
		service: service,
		logger:  logger,
	}
}

// CreateOrderPaymentRequest представляет запрос на создание платежа для заказа
type CreateOrderPaymentRequest struct {
	OrderID     int    `json:"order_id" validate:"required"`
	Amount      string `json:"amount" validate:"required"`
	Currency    string `json:"currency" validate:"required,len=3"`
	ReturnURL   string `json:"return_url" validate:"required,url"`
	CancelURL   string `json:"cancel_url" validate:"required,url"`
	Description string `json:"description"`
}

// CreateOrderPaymentResponse представляет ответ на создание платежа для заказа
type CreateOrderPaymentResponse struct {
	SessionID      string `json:"session_id"`
	PaymentURL     string `json:"payment_url"`
	Status         string `json:"status"`
	ExpiresAt      string `json:"expires_at,omitempty"`
	RequiresAction bool   `json:"requires_action"`
}

// CreateOrderPayment создает новый платеж для заказа
// @Summary Create Order Payment
// @Description Создает новый платеж для заказа через платежную систему
// @Tags order-payments
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body CreateOrderPaymentRequest true "Order payment details"
// @Success 200 {object} utils.SuccessResponseSwag{data=CreateOrderPaymentResponse} "Payment session created successfully"
// @Failure 400 {object} utils.ErrorResponseSwag "Invalid request"
// @Failure 401 {object} utils.ErrorResponseSwag "Unauthorized"
// @Failure 404 {object} utils.ErrorResponseSwag "Order not found"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Router /api/v1/orders/{id}/payment [post]
func (h *OrderPaymentHandler) CreateOrderPayment(c *fiber.Ctx) error {
	// Получаем ID пользователя из контекста
	userID, ok := c.Locals("user_id").(int)
	if !ok {
		return utils.ErrorResponse(c, 401, "unauthorized")
	}

	// Получаем order ID из параметров
	orderIDStr := c.Params("id")
	orderID, err := strconv.Atoi(orderIDStr)
	if err != nil {
		return utils.ErrorResponse(c, 400, "invalid order ID")
	}

	var req CreateOrderPaymentRequest
	if err := c.BodyParser(&req); err != nil {
		h.logger.Error("Failed to parse request body: %v", err)
		return utils.ErrorResponse(c, 400, "invalid request body")
	}

	// Валидация
	if req.Amount == "" {
		return utils.ErrorResponse(c, 400, "amount is required")
	}
	if req.Currency == "" {
		return utils.ErrorResponse(c, 400, "currency is required")
	}
	if req.ReturnURL == "" {
		return utils.ErrorResponse(c, 400, "return_url is required")
	}

	// Конвертируем amount в float64
	amount, err := decimal.NewFromString(req.Amount)
	if err != nil {
		return utils.ErrorResponse(c, 400, "invalid amount format")
	}
	amountFloat, _ := amount.Float64()

	// Создаем платежную сессию для заказа
	session, err := h.service.CreateOrderPayment(c.Context(), orderID, userID, amountFloat, req.Currency, "allsecure")
	if err != nil {
		h.logger.Error("Failed to create order payment: %v (userID: %d, orderID: %d)", err, userID, orderID)
		return utils.ErrorResponse(c, 500, "order payment creation failed")
	}

	// Формируем ответ
	response := CreateOrderPaymentResponse{
		SessionID:      "", // session.SessionID,
		PaymentURL:     session.PaymentURL,
		Status:         string(session.Status),
		RequiresAction: session.PaymentURL != "",
	}

	h.logger.Info("Order payment session created successfully (orderID: %d, userID: %d)", orderID, userID)

	return utils.SuccessResponse(c, response)
}

// GetOrderPaymentStatus получает статус платежа заказа
// @Summary Get Order Payment Status
// @Description Получает текущий статус платежа заказа
// @Tags order-payments
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "Order ID"
// @Success 200 {object} utils.SuccessResponseSwag{data=models.PaymentSession} "Payment status"
// @Failure 400 {object} utils.ErrorResponseSwag "Invalid order ID"
// @Failure 401 {object} utils.ErrorResponseSwag "Unauthorized"
// @Failure 404 {object} utils.ErrorResponseSwag "Order or payment not found"
// @Router /api/v1/orders/{id}/payment/status [get]
func (h *OrderPaymentHandler) GetOrderPaymentStatus(c *fiber.Ctx) error {
	// Получаем пользователя из контекста
	userID := c.Locals("user_id")
	if userID == nil {
		return utils.ErrorResponse(c, 401, "user not authenticated")
	}

	// Получаем order ID из параметров
	orderIDStr := c.Params("id")
	orderID, err := strconv.Atoi(orderIDStr)
	if err != nil {
		return utils.ErrorResponse(c, 400, "invalid order ID")
	}

	h.logger.Info("Order payment status requested (orderID: %d, userID: %v)", orderID, userID)

	// TODO: Реализовать получение статуса платежа заказа из базы данных
	return utils.SuccessResponse(c, map[string]interface{}{
		"order_id": orderID,
		"status":   "pending", // Временное значение
	})
}

// CancelOrderPayment отменяет платеж заказа
// @Summary Cancel Order Payment
// @Description Отменяет платеж заказа
// @Tags order-payments
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "Order ID"
// @Success 200 {object} utils.SuccessResponseSwag "Payment cancelled successfully"
// @Failure 400 {object} utils.ErrorResponseSwag "Invalid order ID"
// @Failure 401 {object} utils.ErrorResponseSwag "Unauthorized"
// @Failure 404 {object} utils.ErrorResponseSwag "Order or payment not found"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Router /api/v1/orders/{id}/payment/cancel [post]
func (h *OrderPaymentHandler) CancelOrderPayment(c *fiber.Ctx) error {
	// Получаем пользователя из контекста
	userID := c.Locals("user_id")
	if userID == nil {
		return utils.ErrorResponse(c, 401, "user not authenticated")
	}

	// Получаем order ID из параметров
	orderIDStr := c.Params("id")
	orderID, err := strconv.Atoi(orderIDStr)
	if err != nil {
		return utils.ErrorResponse(c, 400, "invalid order ID")
	}

	h.logger.Info("Order payment cancellation requested (orderID: %d, userID: %v)", orderID, userID)

	// TODO: Реализовать отмену платежа заказа
	return utils.SuccessResponse(c, "payment cancelled")
}
