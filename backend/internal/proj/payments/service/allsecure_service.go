package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/shopspring/decimal"

	"backend/internal/domain/models"
	"backend/internal/pkg/allsecure"
	"backend/internal/proj/payments/repository"
	"backend/pkg/logger"
)

// AllSecureService представляет сервис для работы с AllSecure
type AllSecureService struct {
	client         *allsecure.Client
	repository     repository.PaymentRepositoryInterface
	userRepo       UserRepositoryInterface
	listingRepo    ListingRepositoryInterface
	config         *AllSecureConfig
	logger         logger.Logger
	commissionRate decimal.Decimal
}

// AllSecureConfig содержит конфигурацию для AllSecure
type AllSecureConfig struct {
	BaseURL                   string
	Username                  string
	Password                  string
	WebhookURL                string
	WebhookSecret             string
	MarketplaceCommissionRate decimal.Decimal
	EscrowReleaseDays         int
}

// UserRepositoryInterface определяет интерфейс для работы с пользователями
type UserRepositoryInterface interface {
	GetByID(ctx context.Context, id int) (*models.User, error)
}

// ListingRepositoryInterface определяет интерфейс для работы с объявлениями
type ListingRepositoryInterface interface {
	GetByID(ctx context.Context, id int) (*models.MarketplaceListing, error)
}

// NewAllSecureService создает новый сервис AllSecure
func NewAllSecureService(
	client *allsecure.Client,
	repository repository.PaymentRepositoryInterface,
	userRepo UserRepositoryInterface,
	listingRepo ListingRepositoryInterface,
	config *AllSecureConfig,
	logger logger.Logger,
) *AllSecureService {
	return &AllSecureService{
		client:         client,
		repository:     repository,
		userRepo:       userRepo,
		listingRepo:    listingRepo,
		config:         config,
		logger:         logger,
		commissionRate: config.MarketplaceCommissionRate,
	}
}

// CreatePaymentRequest представляет запрос на создание платежа
type CreatePaymentRequest struct {
	UserID      int             `json:"user_id" validate:"required"`
	ListingID   int             `json:"listing_id" validate:"required"`
	Amount      decimal.Decimal `json:"amount" validate:"required,gt=0"`
	Currency    string          `json:"currency" validate:"required,len=3"`
	Description string          `json:"description"`
	ReturnURL   string          `json:"return_url" validate:"required,url"`
}

// PaymentResult представляет результат создания платежа
type PaymentResult struct {
	TransactionID  int64  `json:"transaction_id"`
	GatewayUUID    string `json:"gateway_uuid"`
	Status         string `json:"status"`
	RedirectURL    string `json:"redirect_url,omitempty"`
	RequiresAction bool   `json:"requires_action"`
}

// CreatePayment создает новый платеж через AllSecure
func (s *AllSecureService) CreatePayment(ctx context.Context, req CreatePaymentRequest) (*PaymentResult, error) {
	// 1. Валидация запроса
	if err := s.validatePaymentRequest(ctx, req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// 2. Получение пользователя
	user, err := s.userRepo.GetByID(ctx, req.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Валидируем что объявление существует
	_, err = s.listingRepo.GetByID(ctx, req.ListingID)
	if err != nil {
		return nil, fmt.Errorf("failed to get listing: %w", err)
	}

	// 3. Расчет комиссий
	marketplaceCommission := s.calculateCommission(req.Amount)
	sellerAmount := req.Amount.Sub(marketplaceCommission)

	// 4. Создание записи транзакции
	transaction, err := s.repository.CreateTransaction(ctx, repository.CreateTransactionRequest{
		UserID:                req.UserID,
		ListingID:             req.ListingID,
		Amount:                req.Amount,
		Currency:              req.Currency,
		MarketplaceCommission: marketplaceCommission,
		SellerAmount:          sellerAmount,
		Description:           req.Description,
		Status:                models.PaymentStatusPending,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	// 5. Создание запроса к AllSecure
	allsecureReq := allsecure.TransactionRequest{
		Amount:       req.Amount.String(),
		Currency:     req.Currency,
		Description:  req.Description,
		MerchantTxID: fmt.Sprintf("SVT-%d", transaction.ID),
		Customer: allsecure.Customer{
			Identification: fmt.Sprintf("user-%d", user.ID),
			FirstName:      user.Name,
			Email:          user.Email,
		},
		CallbackURL: s.config.WebhookURL,
		SuccessURL:  req.ReturnURL + "?status=success",
		CancelURL:   req.ReturnURL + "?status=canceled",
		ErrorURL:    req.ReturnURL + "?status=error",
	}

	// 6. Отправка запроса в AllSecure
	response, err := s.client.Preauthorize(ctx, allsecureReq)
	if err != nil {
		// Обновляем статус транзакции
		if updateErr := s.repository.UpdateTransactionStatus(ctx, transaction.ID, models.PaymentStatusFailed, map[string]interface{}{
			"error": err.Error(),
		}); updateErr != nil {
			// Логируем ошибку обновления, но не прерываем выполнение
			_ = updateErr // Explicitly ignore error
		}
		return nil, fmt.Errorf("AllSecure request failed: %w", err)
	}

	// 7. Обновление транзакции данными от AllSecure
	err = s.repository.UpdateTransaction(ctx, transaction.ID, repository.UpdateTransactionRequest{
		GatewayTransactionID: response.UUID,
		GatewayReferenceID:   response.PurchaseID,
		Status:               s.mapAllSecureStatus(response.Status),
		GatewayResponse:      response,
	})
	if err != nil {
		s.logger.Error("Failed to update transaction: %v (transactionID: %d)", err, transaction.ID)
	}

	return &PaymentResult{
		TransactionID:  transaction.ID,
		GatewayUUID:    response.UUID,
		Status:         s.mapAllSecureStatus(response.Status),
		RedirectURL:    response.RedirectURL,
		RequiresAction: response.ReturnType == "REDIRECT",
	}, nil
}

// CapturePayment захватывает авторизованный платеж
func (s *AllSecureService) CapturePayment(ctx context.Context, transactionID int64) error {
	transaction, err := s.repository.GetByID(ctx, transactionID)
	if err != nil {
		return fmt.Errorf("failed to get transaction: %w", err)
	}

	if transaction.Status != models.PaymentStatusAuthorized {
		return fmt.Errorf("transaction not in authorized status: %s", transaction.Status)
	}

	if transaction.GatewayTransactionID == nil {
		return fmt.Errorf("gateway transaction ID is missing")
	}

	// Отправляем capture запрос в AllSecure
	response, err := s.client.Capture(ctx, *transaction.GatewayTransactionID, transaction.Amount.String())
	if err != nil {
		return fmt.Errorf("AllSecure capture failed: %w", err)
	}

	// Обновляем статус транзакции
	now := time.Now()
	err = s.repository.UpdateTransactionStatus(ctx, transactionID, s.mapAllSecureStatus(response.Status), map[string]interface{}{
		"capture_response": response,
		"captured_at":      now,
	})
	if err != nil {
		return fmt.Errorf("failed to update transaction status: %w", err)
	}

	// Создаем escrow платеж
	if response.Success {
		err = s.createEscrowPayment(ctx, transaction)
		if err != nil {
			s.logger.Error("Failed to create escrow payment: %v (transactionID: %d)", err, transactionID)
		}
	}

	return nil
}

// RefundPayment возвращает средства
func (s *AllSecureService) RefundPayment(ctx context.Context, transactionID int64, amount decimal.Decimal) error {
	transaction, err := s.repository.GetByID(ctx, transactionID)
	if err != nil {
		return fmt.Errorf("failed to get transaction: %w", err)
	}

	if transaction.Status != models.PaymentStatusCaptured {
		return fmt.Errorf("transaction not in captured status: %s", transaction.Status)
	}

	if transaction.GatewayTransactionID == nil {
		return fmt.Errorf("gateway transaction ID is missing")
	}

	// Отправляем refund запрос в AllSecure
	response, err := s.client.Refund(ctx, *transaction.GatewayTransactionID, amount.String())
	if err != nil {
		return fmt.Errorf("AllSecure refund failed: %w", err)
	}

	// Обновляем статус транзакции
	err = s.repository.UpdateTransactionStatus(ctx, transactionID, models.PaymentStatusRefunded, map[string]interface{}{
		"refund_response": response,
		"refunded_at":     time.Now(),
	})
	if err != nil {
		return fmt.Errorf("failed to update transaction status: %w", err)
	}

	return nil
}

// ProcessWebhook обрабатывает webhook от AllSecure
func (s *AllSecureService) ProcessWebhook(ctx context.Context, payload []byte) error {
	var webhook allsecure.WebhookPayload
	if err := json.Unmarshal(payload, &webhook); err != nil {
		return fmt.Errorf("failed to unmarshal webhook: %w", err)
	}

	// Находим транзакцию по merchant transaction ID
	merchantTxID := webhook.MerchantTxID
	if !strings.HasPrefix(merchantTxID, "SVT-") {
		return fmt.Errorf("invalid merchant transaction ID format: %s", merchantTxID)
	}

	transactionIDStr := strings.TrimPrefix(merchantTxID, "SVT-")
	transactionID, err := strconv.ParseInt(transactionIDStr, 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse transaction ID: %w", err)
	}

	transaction, err := s.repository.GetByID(ctx, transactionID)
	if err != nil {
		return fmt.Errorf("failed to get transaction: %w", err)
	}

	// Обновляем статус транзакции
	newStatus := s.mapAllSecureStatus(webhook.Status)
	err = s.repository.UpdateTransactionStatus(ctx, transactionID, newStatus, map[string]interface{}{
		"webhook_data": webhook,
		"updated_at":   time.Now(),
	})
	if err != nil {
		return fmt.Errorf("failed to update transaction status: %w", err)
	}

	// Обрабатываем специфичную логику по статусам
	switch newStatus {
	case models.PaymentStatusAuthorized:
		// Создаем escrow платеж
		err = s.createEscrowPayment(ctx, transaction)
		if err != nil {
			s.logger.Error("Failed to create escrow payment: %v (transactionID: %d)", err, transactionID)
		}

	case models.PaymentStatusCaptured:
		// Освобождаем escrow платеж
		err = s.releaseEscrowPayment(ctx, transaction)
		if err != nil {
			s.logger.Error("Failed to release escrow payment: %v (transactionID: %d)", err, transactionID)
		}

	case models.PaymentStatusFailed:
		// Обрабатываем неудачный платеж
		s.handleFailedPayment(ctx, transaction)
	}

	return nil
}

// validatePaymentRequest валидирует запрос на создание платежа
func (s *AllSecureService) validatePaymentRequest(ctx context.Context, req CreatePaymentRequest) error {
	if req.Amount.LessThanOrEqual(decimal.Zero) {
		return fmt.Errorf("amount must be positive")
	}

	if req.Currency == "" {
		return fmt.Errorf("currency is required")
	}

	if req.Currency != "RSD" && req.Currency != "EUR" && req.Currency != "USD" {
		return fmt.Errorf("unsupported currency: %s", req.Currency)
	}

	if req.UserID <= 0 {
		return fmt.Errorf("invalid user ID")
	}

	if req.ListingID <= 0 {
		return fmt.Errorf("invalid listing ID")
	}

	if req.ReturnURL == "" {
		return fmt.Errorf("return URL is required")
	}

	return nil
}

// calculateCommission рассчитывает комиссию маркетплейса
func (s *AllSecureService) calculateCommission(amount decimal.Decimal) decimal.Decimal {
	commission := amount.Mul(s.commissionRate)
	// Округляем до 4 знаков после запятой для точности малых комиссий
	return commission.Round(4)
}

// mapAllSecureStatus маппит статус AllSecure в наш статус
func (s *AllSecureService) mapAllSecureStatus(allsecureStatus string) string {
	switch allsecureStatus {
	case "FINISHED":
		return models.PaymentStatusCaptured
	case "PENDING":
		return models.PaymentStatusAuthorized
	case "ERROR":
		return models.PaymentStatusFailed
	default:
		return models.PaymentStatusPending
	}
}

// createEscrowPayment создает escrow платеж
func (s *AllSecureService) createEscrowPayment(ctx context.Context, transaction *models.PaymentTransaction) error {
	if transaction.ListingID == nil {
		return fmt.Errorf("listing ID is required for escrow payment")
	}

	listing, err := s.listingRepo.GetByID(ctx, *transaction.ListingID)
	if err != nil {
		return fmt.Errorf("failed to get listing: %w", err)
	}

	escrowReq := repository.CreateEscrowRequest{
		PaymentTransactionID:  transaction.ID,
		SellerID:              listing.UserID,
		BuyerID:               transaction.UserID,
		ListingID:             *transaction.ListingID,
		Amount:                transaction.Amount,
		MarketplaceCommission: *transaction.MarketplaceCommission,
		SellerAmount:          *transaction.SellerAmount,
		Status:                models.EscrowStatusHeld,
	}

	_, err = s.repository.CreateEscrowPayment(ctx, escrowReq)
	if err != nil {
		return fmt.Errorf("failed to create escrow payment: %w", err)
	}

	return nil
}

// releaseEscrowPayment освобождает escrow платеж
func (s *AllSecureService) releaseEscrowPayment(ctx context.Context, transaction *models.PaymentTransaction) error {
	// Здесь будет логика освобождения escrow платежа
	// Возможно, с автоматической выплатой продавцу
	s.logger.Info("Escrow payment released (transactionID: %d)", transaction.ID)
	return nil
}

// handleFailedPayment обрабатывает неудачный платеж
func (s *AllSecureService) handleFailedPayment(ctx context.Context, transaction *models.PaymentTransaction) {
	s.logger.Info("Payment failed (transactionID: %d)", transaction.ID)
	// Здесь может быть логика уведомлений, очистки данных и т.д.
}
