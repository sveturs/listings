package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/shopspring/decimal"

	"backend/internal/proj/payments/service"
	"backend/pkg/logger"
	"backend/pkg/utils"
)

// PaymentHandler обрабатывает HTTP запросы для платежей
type PaymentHandler struct {
	service *service.AllSecureService
	logger  logger.Logger
}

// NewPaymentHandler создает новый payment handler
func NewPaymentHandler(service *service.AllSecureService, logger logger.Logger) *PaymentHandler {
	return &PaymentHandler{
		service: service,
		logger:  logger,
	}
}

// CreatePaymentRequest представляет запрос на создание платежа
type CreatePaymentRequest struct {
	ListingID   int    `json:"listing_id" validate:"required"`
	Amount      string `json:"amount" validate:"required"`
	Currency    string `json:"currency" validate:"required,len=3"`
	Description string `json:"description"`
	ReturnURL   string `json:"return_url" validate:"required,url"`
}

// CreatePaymentResponse представляет ответ на создание платежа
type CreatePaymentResponse struct {
	TransactionID  int64  `json:"transaction_id"`
	GatewayUUID    string `json:"gateway_uuid"`
	Status         string `json:"status"`
	RedirectURL    string `json:"redirect_url,omitempty"`
	RequiresAction bool   `json:"requires_action"`
}

// CreatePayment создает новый платеж
// @Summary Create Payment
// @Description Создает новый платеж через AllSecure
// @Tags payments
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body CreatePaymentRequest true "Payment details"
// @Success 200 {object} utils.SuccessResponseSwag{data=CreatePaymentResponse} "Payment created successfully"
// @Failure 400 {object} utils.ErrorResponseSwag "Invalid request"
// @Failure 401 {object} utils.ErrorResponseSwag "Unauthorized"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Router /payments/create [post]
func (h *PaymentHandler) CreatePayment(c *fiber.Ctx) error {
	// Получаем ID пользователя из контекста
	userID, ok := c.Locals("user_id").(int)
	if !ok {
		return utils.ErrorResponse(c, 401, "unauthorized")
	}

	var req CreatePaymentRequest
	if err := c.BodyParser(&req); err != nil {
		h.logger.Error("Failed to parse request body: %v", err)
		return utils.ErrorResponse(c, 400, "invalid request body")
	}

	// Простая валидация
	if req.ListingID <= 0 {
		return utils.ErrorResponse(c, 400, "listing_id is required")
	}
	if req.Amount == "" {
		return utils.ErrorResponse(c, 400, "amount is required")
	}
	if req.Currency == "" {
		return utils.ErrorResponse(c, 400, "currency is required")
	}
	if req.ReturnURL == "" {
		return utils.ErrorResponse(c, 400, "return_url is required")
	}

	// Конвертируем amount в decimal
	amount, err := decimal.NewFromString(req.Amount)
	if err != nil {
		return utils.ErrorResponse(c, 400, "invalid amount format")
	}

	// Создаем запрос к сервису
	serviceReq := service.CreatePaymentRequest{
		UserID:      userID,
		ListingID:   req.ListingID,
		Amount:      amount,
		Currency:    req.Currency,
		Description: req.Description,
		ReturnURL:   req.ReturnURL,
	}

	// Создаем платеж
	result, err := h.service.CreatePayment(c.Context(), serviceReq)
	if err != nil {
		h.logger.Error("Failed to create payment: %v (userID: %d, listingID: %d)", err, userID, req.ListingID)
		return utils.ErrorResponse(c, 500, "payment creation failed")
	}

	// Формируем ответ
	response := CreatePaymentResponse{
		TransactionID:  result.TransactionID,
		GatewayUUID:    result.GatewayUUID,
		Status:         result.Status,
		RedirectURL:    result.RedirectURL,
		RequiresAction: result.RequiresAction,
	}

	h.logger.Info("Payment created successfully (transactionID: %d, userID: %d)", result.TransactionID, userID)

	return utils.SuccessResponse(c, response)
}

// CapturePayment захватывает авторизованный платеж
// @Summary Capture Payment
// @Description Захватывает авторизованный платеж
// @Tags payments
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "Transaction ID"
// @Success 200 {object} utils.SuccessResponseSwag "Payment captured successfully"
// @Failure 400 {object} utils.ErrorResponseSwag "Invalid transaction ID"
// @Failure 401 {object} utils.ErrorResponseSwag "Unauthorized"
// @Failure 404 {object} utils.ErrorResponseSwag "Transaction not found"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Router /payments/{id}/capture [post]
func (h *PaymentHandler) CapturePayment(c *fiber.Ctx) error {
	// Получаем пользователя из контекста
	userID := c.Locals("user_id")
	if userID == nil {
		return utils.ErrorResponse(c, 401, "user not authenticated")
	}

	// Получаем transaction ID из параметров
	transactionIDStr := c.Params("id")
	transactionID, err := strconv.ParseInt(transactionIDStr, 10, 64)
	if err != nil {
		return utils.ErrorResponse(c, 400, "invalid transaction ID")
	}

	// Захватываем платеж
	err = h.service.CapturePayment(c.Context(), transactionID)
	if err != nil {
		h.logger.Error("Failed to capture payment: %v (transactionID: %d, userID: %v)", err, transactionID, userID)
		return utils.ErrorResponse(c, 500, "payment capture failed")
	}

	h.logger.Info("Payment captured successfully (transactionID: %d, userID: %v)", transactionID, userID)

	return utils.SuccessResponse(c, "payment captured")
}

// RefundPaymentRequest представляет запрос на возврат средств
type RefundPaymentRequest struct {
	Amount string `json:"amount" validate:"required"`
	Reason string `json:"reason"`
}

// RefundPayment возвращает средства
// @Summary Refund Payment
// @Description Возвращает средства по платежу
// @Tags payments
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "Transaction ID"
// @Param request body RefundPaymentRequest true "Refund details"
// @Success 200 {object} utils.SuccessResponseSwag "Refund processed successfully"
// @Failure 400 {object} utils.ErrorResponseSwag "Invalid request"
// @Failure 401 {object} utils.ErrorResponseSwag "Unauthorized"
// @Failure 404 {object} utils.ErrorResponseSwag "Transaction not found"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Router /payments/{id}/refund [post]
func (h *PaymentHandler) RefundPayment(c *fiber.Ctx) error {
	// Получаем пользователя из контекста
	userID := c.Locals("user_id")
	if userID == nil {
		return utils.ErrorResponse(c, 401, "user not authenticated")
	}

	// Получаем transaction ID из параметров
	transactionIDStr := c.Params("id")
	transactionID, err := strconv.ParseInt(transactionIDStr, 10, 64)
	if err != nil {
		return utils.ErrorResponse(c, 400, "invalid transaction ID")
	}

	var req RefundPaymentRequest
	if err := c.BodyParser(&req); err != nil {
		h.logger.Error("Failed to parse request body: %v", err)
		return utils.ErrorResponse(c, 400, "invalid request body")
	}

	// Простая валидация
	if req.Amount == "" {
		return utils.ErrorResponse(c, 400, "amount is required")
	}

	// Конвертируем amount в decimal
	amount, err := decimal.NewFromString(req.Amount)
	if err != nil {
		return utils.ErrorResponse(c, 400, "invalid amount format")
	}

	// Возвращаем средства
	err = h.service.RefundPayment(c.Context(), transactionID, amount)
	if err != nil {
		h.logger.Error("Failed to refund payment: %v (transactionID: %d, userID: %v)", err, transactionID, userID)
		return utils.ErrorResponse(c, 500, "payment refund failed")
	}

	h.logger.Info("Payment refunded successfully (transactionID: %d, amount: %s, userID: %v)", transactionID, amount.String(), userID)

	return utils.SuccessResponse(c, "payment refunded")
}

// GetPaymentStatus получает статус платежа
// @Summary Get Payment Status
// @Description Получает текущий статус платежа
// @Tags payments
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "Transaction ID"
// @Success 200 {object} utils.SuccessResponseSwag{data=backend_internal_proj_payments_models.PaymentTransaction} "Payment status"
// @Failure 400 {object} utils.ErrorResponseSwag "Invalid transaction ID"
// @Failure 401 {object} utils.ErrorResponseSwag "Unauthorized"
// @Failure 404 {object} utils.ErrorResponseSwag "Transaction not found"
// @Router /payments/{id}/status [get]
func (h *PaymentHandler) GetPaymentStatus(c *fiber.Ctx) error {
	// Получаем пользователя из контекста
	userID := c.Locals("user_id")
	if userID == nil {
		return utils.ErrorResponse(c, 401, "user not authenticated")
	}

	// Получаем transaction ID из параметров
	transactionIDStr := c.Params("id")
	transactionID, err := strconv.ParseInt(transactionIDStr, 10, 64)
	if err != nil {
		return utils.ErrorResponse(c, 400, "invalid transaction ID")
	}

	// Здесь нужно добавить метод в сервис для получения статуса
	// transaction, err := h.service.GetPaymentStatus(c.Context(), transactionID)
	// if err != nil {
	//     return utils.ErrorResponse(c, 404, "transaction not found")
	// }

	// Временно возвращаем простой ответ
	h.logger.Info("Payment status requested (transactionID: %d, userID: %v)", transactionID, userID)

	return utils.SuccessResponse(c, map[string]interface{}{
		"transaction_id": transactionID,
		"status":         "pending", // Временное значение
	})
}
