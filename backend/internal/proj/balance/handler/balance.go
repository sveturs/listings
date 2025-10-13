// Package handler
// backend/internal/proj/balance/handler/balance.go
package handler

import (
	"context"
	"log"
	"strconv"

	authMiddleware "github.com/sveturs/auth/pkg/http/fiber/middleware"

	"github.com/gofiber/fiber/v2"

	"backend/internal/config"
	"backend/internal/domain/models"
	balance "backend/internal/proj/balance/service"
	paymentService "backend/internal/proj/payments/service"
	"backend/pkg/utils"
)

// ContextKey is a type for context keys to avoid collisions
type ContextKey string

const (
	// Context keys
	ContextKeyReturnURL ContextKey = "return_url"
)

type BalanceHandler struct {
	balanceService balance.BalanceServiceInterface
	paymentService paymentService.PaymentServiceInterface // Теперь используем импортированный пакет
}

func NewBalanceHandler(balanceService balance.BalanceServiceInterface, paymentService paymentService.PaymentServiceInterface) *BalanceHandler {
	return &BalanceHandler{
		balanceService: balanceService,
		paymentService: paymentService,
	}
}

// GetBalance returns user balance information
// @Summary Get user balance
// @Description Returns current balance and frozen balance for authenticated user
// @Tags balance
// @Accept json
// @Produce json
// @Success 200 {object} utils.SuccessResponseSwag{data=models.UserBalance} "User balance information"
// @Failure 401 {object} utils.ErrorResponseSwag "Unauthorized"
// @Failure 500 {object} utils.ErrorResponseSwag "balance.getError"
// @Security BearerAuth
// @Router /api/v1/balance [get]
func (h *BalanceHandler) GetBalance(c *fiber.Ctx) error {
	userID, _ := authMiddleware.GetUserID(c)

	log.Printf("Getting balance for user %d", userID)

	balance, err := h.balanceService.GetBalance(c.Context(), userID)
	if err != nil {
		log.Printf("Error getting balance for user %d: %v", userID, err)
		// Если записи нет, возвращаем нулевой баланс
		if err.Error() == "no rows in result set" {
			log.Printf("No balance record found for user %d, returning zero balance", userID)
			return utils.SuccessResponse(c, &models.UserBalance{
				UserID:   userID,
				Balance:  0,
				Currency: config.GetGlobalDefaultCurrency(),
			})
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "balance.getError")
	}

	log.Printf("Balance for user %d: amount=%f, currency=%s", userID, balance.Balance, balance.Currency)
	return utils.SuccessResponse(c, balance)
}

// DepositRequest represents deposit creation request
type DepositRequest struct {
	Amount        float64 `json:"amount" example:"1000.50"`
	PaymentMethod string  `json:"payment_method" example:"card"`
	ReturnURL     string  `json:"return_url,omitempty" example:"http://localhost:3001/en/balance/deposit/success"`
}

// CreateDeposit creates a new deposit payment session
// @Summary Create deposit
// @Description Creates a new payment session for balance deposit
// @Tags balance
// @Accept json
// @Produce json
// @Param request body DepositRequest true "Deposit details"
// @Success 200 {object} utils.SuccessResponseSwag{data=models.PaymentSession} "Payment session created"
// @Failure 400 {object} utils.ErrorResponseSwag "balance.invalidRequest"
// @Failure 401 {object} utils.ErrorResponseSwag "Unauthorized"
// @Failure 500 {object} utils.ErrorResponseSwag "balance.createDepositError"
// @Security BearerAuth
// @Router /api/v1/balance/deposit [post]
func (h *BalanceHandler) CreateDeposit(c *fiber.Ctx) error {
	userID, _ := authMiddleware.GetUserID(c)

	var request DepositRequest

	if err := c.BodyParser(&request); err != nil {
		log.Printf("Error parsing deposit request: %v", err)
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "balance.invalidRequest")
	}

	log.Printf("Processing deposit request: amount=%f, method=%s, return_url=%s", request.Amount, request.PaymentMethod, request.ReturnURL)

	// Добавляем return_url в контекст для мок-сервиса
	ctx := context.WithValue(c.Context(), ContextKeyReturnURL, request.ReturnURL)

	// Создаем платежную сессию вместо прямого создания депозита
	session, err := h.paymentService.CreatePaymentSession(
		ctx,
		userID,
		request.Amount,
		"rsd",
		request.PaymentMethod,
	)
	if err != nil {
		log.Printf("Error creating payment session: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "balance.createDepositError")
	}

	log.Printf("Created payment session: %+v", session)
	return utils.SuccessResponse(c, session)
}

// GetTransactions returns user transaction history
// @Summary Get transaction history
// @Description Returns paginated list of user balance transactions
// @Tags balance
// @Accept json
// @Produce json
// @Param limit query int false "Limit" default(20)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} utils.SuccessResponseSwag{data=[]models.BalanceTransaction} "List of transactions"
// @Failure 401 {object} utils.ErrorResponseSwag "Unauthorized"
// @Failure 500 {object} utils.ErrorResponseSwag "balance.getTransactionsError"
// @Security BearerAuth
// @Router /api/v1/balance/transactions [get]
func (h *BalanceHandler) GetTransactions(c *fiber.Ctx) error {
	userID, _ := authMiddleware.GetUserID(c)

	// Заменяем utils.QueryInt на собственную реализацию
	limit := 20
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	offset := 0
	if offsetStr := c.Query("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	transactions, err := h.balanceService.GetTransactions(c.Context(), userID, limit, offset)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "balance.getTransactionsError")
	}

	return utils.SuccessResponse(c, transactions)
}

// GetPaymentMethods returns available payment methods
// @Summary Get payment methods
// @Description Returns list of available payment methods
// @Tags balance
// @Accept json
// @Produce json
// @Success 200 {object} utils.SuccessResponseSwag{data=[]models.PaymentMethod} "List of payment methods"
// @Failure 500 {object} utils.ErrorResponseSwag "balance.getPaymentMethodsError"
// @Router /api/v1/balance/payment-methods [get]
func (h *BalanceHandler) GetPaymentMethods(c *fiber.Ctx) error {
	methods, err := h.balanceService.GetPaymentMethods(c.Context())
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "balance.getPaymentMethodsError")
	}

	return utils.SuccessResponse(c, methods)
}

// CompleteMockPaymentRequest represents mock payment completion request
type CompleteMockPaymentRequest struct {
	SessionID string  `json:"session_id" example:"mock_session_7_1751361491"`
	Amount    float64 `json:"amount,omitempty" example:"2000"`
}

// CompleteMockPayment completes a mock payment session
// @Summary Complete mock payment
// @Description Completes a mock payment session for testing purposes
// @Tags balance
// @Accept json
// @Produce json
// @Param request body CompleteMockPaymentRequest true "Payment session ID"
// @Success 200 {object} utils.SuccessResponseSwag{data=object} "Payment completed successfully"
// @Failure 400 {object} utils.ErrorResponseSwag "balance.invalidRequest"
// @Failure 401 {object} utils.ErrorResponseSwag "Unauthorized"
// @Failure 500 {object} utils.ErrorResponseSwag "balance.completeMockPaymentError"
// @Security BearerAuth
// @Router /api/v1/balance/mock/complete [post]
func (h *BalanceHandler) CompleteMockPayment(c *fiber.Ctx) error {
	userID, _ := authMiddleware.GetUserID(c)

	var request CompleteMockPaymentRequest
	if err := c.BodyParser(&request); err != nil {
		log.Printf("Error parsing mock payment request: %v", err)
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "balance.invalidRequest")
	}

	log.Printf("Completing mock payment for session: %s, user: %d", request.SessionID, userID)

	// В реальной системе здесь бы была проверка сессии и обновление статуса
	// Для mock системы сразу создаем транзакцию пополнения

	// Используем сумму из запроса, если она передана
	amount := request.Amount
	if amount == 0 {
		// Если сумма не передана, используем дефолтное значение
		amount = 2000.0
	}

	// Создаем транзакцию пополнения
	transaction, err := h.balanceService.CreateDeposit(c.Context(), userID, amount, "mock_payment")
	if err != nil {
		log.Printf("Error creating deposit transaction: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "balance.completeMockPaymentError")
	}

	log.Printf("Mock payment completed successfully, transaction: %+v", transaction)

	// Получаем обновленный баланс
	balance, err := h.balanceService.GetBalance(c.Context(), userID)
	if err != nil {
		log.Printf("Error getting updated balance: %v", err)
		balance = &models.UserBalance{Balance: amount}
	}

	return utils.SuccessResponse(c, map[string]interface{}{
		"success":     true,
		"transaction": transaction,
		"new_balance": balance.Balance,
	})
}
