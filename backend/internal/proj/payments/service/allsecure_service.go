package service

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
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

// Payment status constants
const (
	PaymentStatusPending   = "pending"
	PaymentStatusCompleted = "completed"
	PaymentStatusFailed    = "failed"
	PaymentStatusCanceled  = "canceled"
)

// AllSecureClientInterface определяет интерфейс для AllSecure клиента
type AllSecureClientInterface interface {
	Preauthorize(ctx context.Context, req allsecure.TransactionRequest) (*allsecure.TransactionResponse, error)
	Capture(ctx context.Context, uuid string, amount string) (*allsecure.TransactionResponse, error)
	Refund(ctx context.Context, uuid string, amount string) (*allsecure.TransactionResponse, error)
	Debit(ctx context.Context, req allsecure.TransactionRequest) (*allsecure.TransactionResponse, error)
	Void(ctx context.Context, uuid string) (*allsecure.TransactionResponse, error)
}

// AllSecureService представляет сервис для работы с AllSecure
type AllSecureService struct {
	client         AllSecureClientInterface
	repository     repository.PaymentRepositoryInterface
	userRepo       UserRepositoryInterface
	listingRepo    ListingRepositoryInterface
	config         *AllSecureConfig
	logger         *logger.Logger
	commissionRate decimal.Decimal
	retryManager   *WebhookRetryManager
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
	logger *logger.Logger,
) *AllSecureService {
	// Create retry manager
	retryConfig := DefaultWebhookRetryConfig()
	retryManager := NewWebhookRetryManager(retryConfig, logger)

	// Start retry worker in background
	ctx := context.Background()
	retryManager.Start(ctx)

	return &AllSecureService{
		client:         client,
		repository:     repository,
		userRepo:       userRepo,
		listingRepo:    listingRepo,
		config:         config,
		logger:         logger,
		commissionRate: config.MarketplaceCommissionRate,
		retryManager:   retryManager,
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

// ValidateWebhookSignature validates webhook signature from AllSecure
func (s *AllSecureService) ValidateWebhookSignature(payload []byte, signature string) bool {
	if s.config.WebhookSecret == "" {
		s.logger.Info("Webhook secret not configured, skipping signature validation")
		return true // For backward compatibility, but should be fixed in production
	}

	// AllSecure typically uses HMAC-SHA256 or similar
	expectedSignature := s.calculateWebhookSignature(payload)
	return subtle.ConstantTimeCompare([]byte(signature), []byte(expectedSignature)) == 1
}

// calculateWebhookSignature calculates expected webhook signature
func (s *AllSecureService) calculateWebhookSignature(payload []byte) string {
	h := hmac.New(sha256.New, []byte(s.config.WebhookSecret))
	h.Write(payload)
	return hex.EncodeToString(h.Sum(nil))
}

// HandleWebhook processes webhook from AllSecure payment gateway
func (s *AllSecureService) HandleWebhook(ctx context.Context, payload []byte, signature string) error {
	// Validate webhook signature
	if !s.ValidateWebhookSignature(payload, signature) {
		s.logger.Error("Invalid webhook signature")
		return fmt.Errorf("invalid webhook signature")
	}

	// Parse webhook payload
	var webhookData map[string]interface{}
	if err := json.Unmarshal(payload, &webhookData); err != nil {
		s.logger.Error("Failed to parse webhook payload: %v", err)

		// Add to retry queue for transient parsing errors
		if s.retryManager != nil {
			retryJob := &WebhookRetryJob{
				ID:          fmt.Sprintf("webhook_%d", time.Now().UnixNano()),
				WebhookType: "allsecure_payment",
				Payload:     payload,
				Signature:   signature,
				Endpoint:    "internal",
				RetryCount:  0,
				LastError:   err.Error(),
				CreatedAt:   time.Now(),
				Metadata:    map[string]interface{}{"error_type": "parse_error"},
			}
			if retryErr := s.retryManager.AddRetryJob(retryJob); retryErr != nil {
				s.logger.Error("Failed to add webhook to retry queue: %v", retryErr)
			}
		}

		return fmt.Errorf("failed to parse webhook payload: %w", err)
	}

	// Extract transaction information
	transactionID, ok := webhookData["merchantTxId"].(string)
	if !ok {
		return fmt.Errorf("missing merchantTxId in webhook payload")
	}

	status, ok := webhookData["status"].(string)
	if !ok {
		return fmt.Errorf("missing status in webhook payload")
	}

	// Process the webhook based on status
	s.logger.Info("Processing webhook for transaction %s with status %s", transactionID, status)

	// Update transaction status in database
	// Parse transaction ID from format "SVT-123"
	parts := strings.Split(transactionID, "-")
	if len(parts) != 2 || parts[0] != "SVT" {
		return fmt.Errorf("invalid transaction ID format: %s", transactionID)
	}

	_, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse transaction ID: %w", err)
	}

	// Map AllSecure status to our internal status
	var internalStatus string
	switch status {
	case "SUCCESS", "CAPTURED":
		internalStatus = PaymentStatusCompleted
	case "PENDING", "PREAUTHORIZED":
		internalStatus = PaymentStatusPending
	case "FAILED", "DECLINED", "ERROR":
		internalStatus = PaymentStatusFailed
	case "CANCELLED", "VOIDED": //nolint:misspell // CANCELLED is the correct spelling used by AllSecure API
		internalStatus = PaymentStatusCanceled
	default:
		s.logger.Info("Unknown webhook status: %s", status)
		internalStatus = PaymentStatusPending
	}

	// Store webhook data in repository for audit
	// TODO: Implement proper transaction status update when repository methods are available
	s.logger.Info("Webhook processed - Transaction: %s, Status: %s", transactionID, internalStatus)

	// Handle specific status actions
	switch internalStatus {
	case "completed":
		s.logger.Info("Payment completed for transaction %s", transactionID)
		// TODO: Trigger order fulfillment, inventory update, etc.
	case "failed":
		s.logger.Info("Payment failed for transaction %s", transactionID)
		// TODO: Trigger failure handling, notification, etc.
	}

	return nil
}

// HandleOrderPaymentWebhook processes webhook specifically for order payments
func (s *AllSecureService) HandleOrderPaymentWebhook(ctx context.Context, payload []byte, signature string) error {
	// Validate webhook signature first
	if !s.ValidateWebhookSignature(payload, signature) {
		s.logger.Error("Invalid webhook signature for order payment")
		return fmt.Errorf("invalid webhook signature")
	}

	// For order payments, we can add additional order-specific logic here
	// For now, delegate to the main webhook handler
	return s.HandleWebhook(ctx, payload, signature)
}
