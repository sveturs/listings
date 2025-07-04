package service

import (
	"context"
	"fmt"

	"github.com/shopspring/decimal"

	"backend/internal/proj/payments/repository"
)

// PaymentService адаптер для интеграции с платежной системой
type PaymentService struct {
	paymentRepo repository.PaymentRepositoryInterface
}

// NewPaymentService создает новый платежный сервис
func NewPaymentService() *PaymentService {
	// TODO: Добавить инжекцию настоящего репозитория
	return &PaymentService{
		paymentRepo: nil, // Пока nil, интеграция будет позже
	}
}

// CreatePaymentRequest запрос на создание платежа
type CreatePaymentRequest struct {
	UserID        int
	Amount        decimal.Decimal
	Currency      string
	PaymentMethod string
	Description   string
	ReturnURL     string
	MetaData      map[string]interface{}
	CaptureMode   string
	AutoCaptureIn int
}

// CreatePayment создает новый платеж с поддержкой delayed capture
func (p *PaymentService) CreatePayment(ctx context.Context, req CreatePaymentRequest) (*PaymentResult, error) {
	if p.paymentRepo == nil {
		// Временная заглушка пока репозиторий не интегрирован
		// Извлекаем order_id из метаданных для формирования URL
		var orderID string
		if oid, ok := req.MetaData["order_id"].(int64); ok {
			orderID = fmt.Sprintf("%d", oid)
		}

		// Генерируем session_id для mock платежа
		sessionID := fmt.Sprintf("mock_%d_%d", req.UserID, req.Amount.IntPart())

		// Формируем URL с параметрами для mock страницы
		paymentURL := fmt.Sprintf("/payment/mock?session_id=%s&amount=%.2f&currency=%s",
			sessionID, req.Amount.InexactFloat64(), req.Currency)

		if orderID != "" {
			paymentURL += "&order_id=" + orderID
		}

		return &PaymentResult{
			TransactionID: 1,          // временный ID
			PaymentURL:    paymentURL, // URL с параметрами
			Status:        "pending",
		}, nil
	}

	// Извлекаем listing_id из метаданных
	var listingID int
	if lid, ok := req.MetaData["listing_id"].(int64); ok {
		listingID = int(lid)
	}

	// Рассчитываем комиссию маркетплейса (по умолчанию 5%)
	commissionRate := decimal.NewFromFloat(0.05) // 5%
	marketplaceCommission := req.Amount.Mul(commissionRate)
	sellerAmount := req.Amount.Sub(marketplaceCommission)

	// Создаем транзакцию
	createReq := repository.CreateTransactionRequest{
		UserID:                req.UserID,
		ListingID:             listingID,
		Amount:                req.Amount,
		Currency:              req.Currency,
		MarketplaceCommission: marketplaceCommission,
		SellerAmount:          sellerAmount,
		Description:           req.Description,
		Status:                "pending", // Начинаем с pending статуса
	}

	transaction, err := p.paymentRepo.CreateTransaction(ctx, createReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create payment transaction: %w", err)
	}

	// Формируем URL для оплаты
	paymentURL := fmt.Sprintf("/payment/redirect/%d", transaction.ID)
	if req.ReturnURL != "" {
		paymentURL += "?return_url=" + req.ReturnURL
	}

	return &PaymentResult{
		TransactionID: transaction.ID,
		PaymentURL:    paymentURL,
		Status:        transaction.Status,
	}, nil
}

// CapturePayment захватывает предавторизованный платеж
func (p *PaymentService) CapturePayment(ctx context.Context, transactionID int64) error {
	if p.paymentRepo == nil {
		// Временная заглушка
		return nil
	}

	// Обновляем статус транзакции на captured
	err := p.paymentRepo.UpdateTransactionStatus(ctx, transactionID, "captured", map[string]interface{}{
		"capture_action": "manual_capture",
		"captured_at":    "now",
	})
	if err != nil {
		return fmt.Errorf("failed to capture payment: %w", err)
	}

	return nil
}

// GetPaymentStatus получает статус платежа
func (p *PaymentService) GetPaymentStatus(ctx context.Context, transactionID int64) (string, error) {
	if p.paymentRepo == nil {
		return "pending", nil
	}

	transaction, err := p.paymentRepo.GetByID(ctx, transactionID)
	if err != nil {
		return "", fmt.Errorf("failed to get payment status: %w", err)
	}

	return transaction.Status, nil
}
