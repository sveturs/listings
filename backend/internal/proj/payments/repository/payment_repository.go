package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/shopspring/decimal"

	"backend/internal/domain/models"
)

// PaymentRepositoryInterface определяет интерфейс для работы с платежными данными
type PaymentRepositoryInterface interface {
	// Payment Transactions
	CreateTransaction(ctx context.Context, req CreateTransactionRequest) (*models.PaymentTransaction, error)
	GetByID(ctx context.Context, id int64) (*models.PaymentTransaction, error)
	GetByGatewayTransactionID(ctx context.Context, gatewayID string) (*models.PaymentTransaction, error)
	UpdateTransaction(ctx context.Context, id int64, req UpdateTransactionRequest) error
	UpdateTransactionStatus(ctx context.Context, id int64, status string, gatewayResponse map[string]interface{}) error

	// Escrow Payments
	CreateEscrowPayment(ctx context.Context, req CreateEscrowRequest) (*models.EscrowPayment, error)
	GetEscrowByTransactionID(ctx context.Context, transactionID int64) (*models.EscrowPayment, error)
	ReleaseEscrow(ctx context.Context, escrowID int64) error

	// Merchant Payouts
	CreatePayout(ctx context.Context, req CreatePayoutRequest) (*models.MerchantPayout, error)
	GetPayoutsBySellerID(ctx context.Context, sellerID int) ([]*models.MerchantPayout, error)
}

// PaymentRepository реализует интерфейс для работы с платежными данными
type PaymentRepository struct {
	db *sql.DB
}

// NewPaymentRepository создает новый репозиторий
func NewPaymentRepository(db *sql.DB) *PaymentRepository {
	return &PaymentRepository{db: db}
}

// CreateTransactionRequest представляет запрос на создание транзакции
type CreateTransactionRequest struct {
	UserID                int
	ListingID             int
	Amount                decimal.Decimal
	Currency              string
	MarketplaceCommission decimal.Decimal
	SellerAmount          decimal.Decimal
	Description           string
	Status                string
}

// UpdateTransactionRequest представляет запрос на обновление транзакции
type UpdateTransactionRequest struct {
	GatewayTransactionID string
	GatewayReferenceID   string
	Status               string
	GatewayResponse      interface{}
}

// CreateEscrowRequest представляет запрос на создание escrow платежа
type CreateEscrowRequest struct {
	PaymentTransactionID  int64
	SellerID              int
	BuyerID               int
	ListingID             int
	Amount                decimal.Decimal
	MarketplaceCommission decimal.Decimal
	SellerAmount          decimal.Decimal
	Status                string
}

// CreatePayoutRequest представляет запрос на создание выплаты
type CreatePayoutRequest struct {
	SellerID        int
	GatewayID       int
	Amount          decimal.Decimal
	Currency        string
	BankAccountInfo map[string]interface{}
	Status          string
}

// CreateTransaction создает новую платежную транзакцию
func (r *PaymentRepository) CreateTransaction(ctx context.Context, req CreateTransactionRequest) (*models.PaymentTransaction, error) {
	query := `
		INSERT INTO payment_transactions (
			gateway_id, user_id, listing_id, order_reference,
			amount, currency, marketplace_commission, seller_amount,
			description, status, created_at, updated_at
		) VALUES (
			1, $1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW()
		) RETURNING id, created_at, updated_at`

	orderRef := generateOrderReference(req.UserID, req.ListingID)

	var transaction models.PaymentTransaction
	err := r.db.QueryRowContext(ctx, query,
		req.UserID, req.ListingID, orderRef,
		req.Amount, req.Currency, req.MarketplaceCommission, req.SellerAmount,
		req.Description, req.Status,
	).Scan(&transaction.ID, &transaction.CreatedAt, &transaction.UpdatedAt)
	if err != nil {
		return nil, err
	}

	// Заполняем остальные поля
	transaction.UserID = req.UserID
	transaction.ListingID = &req.ListingID
	transaction.OrderReference = orderRef
	transaction.Amount = req.Amount
	transaction.Currency = req.Currency
	transaction.MarketplaceCommission = &req.MarketplaceCommission
	transaction.SellerAmount = &req.SellerAmount
	transaction.Description = &req.Description
	transaction.Status = req.Status

	return &transaction, nil
}

// GetByID получает транзакцию по ID
func (r *PaymentRepository) GetByID(ctx context.Context, id int64) (*models.PaymentTransaction, error) {
	query := `
		SELECT id, gateway_id, user_id, listing_id, order_reference,
			gateway_transaction_id, gateway_reference_id,
			amount, currency, marketplace_commission, seller_amount,
			status, gateway_status, payment_method, customer_email, description,
			gateway_response, error_details,
			created_at, updated_at, authorized_at, captured_at, failed_at
		FROM payment_transactions
		WHERE id = $1`

	var transaction models.PaymentTransaction
	var gatewayResponse, errorDetails sql.NullString

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&transaction.ID, &transaction.GatewayID, &transaction.UserID, &transaction.ListingID,
		&transaction.OrderReference, &transaction.GatewayTransactionID, &transaction.GatewayReferenceID,
		&transaction.Amount, &transaction.Currency, &transaction.MarketplaceCommission, &transaction.SellerAmount,
		&transaction.Status, &transaction.GatewayStatus, &transaction.PaymentMethod,
		&transaction.CustomerEmail, &transaction.Description,
		&gatewayResponse, &errorDetails,
		&transaction.CreatedAt, &transaction.UpdatedAt, &transaction.AuthorizedAt,
		&transaction.CapturedAt, &transaction.FailedAt,
	)
	if err != nil {
		return nil, err
	}

	// Парсим JSON поля
	if gatewayResponse.Valid {
		if err := json.Unmarshal([]byte(gatewayResponse.String), &transaction.GatewayResponse); err != nil {
			// Логируем ошибку, но не прерываем выполнение
		}
	}
	if errorDetails.Valid {
		if err := json.Unmarshal([]byte(errorDetails.String), &transaction.ErrorDetails); err != nil {
			// Логируем ошибку, но не прерываем выполнение
		}
	}

	return &transaction, nil
}

// GetByGatewayTransactionID получает транзакцию по ID от шлюза
func (r *PaymentRepository) GetByGatewayTransactionID(ctx context.Context, gatewayID string) (*models.PaymentTransaction, error) {
	query := `
		SELECT id, gateway_id, user_id, listing_id, order_reference,
			gateway_transaction_id, gateway_reference_id,
			amount, currency, marketplace_commission, seller_amount,
			status, gateway_status, payment_method, customer_email, description,
			gateway_response, error_details,
			created_at, updated_at, authorized_at, captured_at, failed_at
		FROM payment_transactions
		WHERE gateway_transaction_id = $1`

	var transaction models.PaymentTransaction
	var gatewayResponse, errorDetails sql.NullString

	err := r.db.QueryRowContext(ctx, query, gatewayID).Scan(
		&transaction.ID, &transaction.GatewayID, &transaction.UserID, &transaction.ListingID,
		&transaction.OrderReference, &transaction.GatewayTransactionID, &transaction.GatewayReferenceID,
		&transaction.Amount, &transaction.Currency, &transaction.MarketplaceCommission, &transaction.SellerAmount,
		&transaction.Status, &transaction.GatewayStatus, &transaction.PaymentMethod,
		&transaction.CustomerEmail, &transaction.Description,
		&gatewayResponse, &errorDetails,
		&transaction.CreatedAt, &transaction.UpdatedAt, &transaction.AuthorizedAt,
		&transaction.CapturedAt, &transaction.FailedAt,
	)
	if err != nil {
		return nil, err
	}

	// Парсим JSON поля
	if gatewayResponse.Valid {
		if err := json.Unmarshal([]byte(gatewayResponse.String), &transaction.GatewayResponse); err != nil {
			// Логируем ошибку, но не прерываем выполнение
		}
	}
	if errorDetails.Valid {
		if err := json.Unmarshal([]byte(errorDetails.String), &transaction.ErrorDetails); err != nil {
			// Логируем ошибку, но не прерываем выполнение
		}
	}

	return &transaction, nil
}

// UpdateTransaction обновляет транзакцию
func (r *PaymentRepository) UpdateTransaction(ctx context.Context, id int64, req UpdateTransactionRequest) error {
	gatewayResponseJSON, _ := json.Marshal(req.GatewayResponse)

	query := `
		UPDATE payment_transactions SET
			gateway_transaction_id = $2,
			gateway_reference_id = $3,
			status = $4,
			gateway_response = $5,
			updated_at = NOW()
		WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query,
		id, req.GatewayTransactionID, req.GatewayReferenceID,
		req.Status, string(gatewayResponseJSON),
	)

	return err
}

// UpdateTransactionStatus обновляет статус транзакции
func (r *PaymentRepository) UpdateTransactionStatus(ctx context.Context, id int64, status string, gatewayResponse map[string]interface{}) error {
	gatewayResponseJSON, _ := json.Marshal(gatewayResponse)

	query := `
		UPDATE payment_transactions SET
			status = $2,
			gateway_response = $3,
			updated_at = NOW()
		WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, id, status, string(gatewayResponseJSON))
	return err
}

// CreateEscrowPayment создает escrow платеж
func (r *PaymentRepository) CreateEscrowPayment(ctx context.Context, req CreateEscrowRequest) (*models.EscrowPayment, error) {
	query := `
		INSERT INTO escrow_payments (
			payment_transaction_id, seller_id, buyer_id, listing_id,
			amount, marketplace_commission, seller_amount, status,
			created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW()
		) RETURNING id, created_at, updated_at`

	var escrow models.EscrowPayment
	err := r.db.QueryRowContext(ctx, query,
		req.PaymentTransactionID, req.SellerID, req.BuyerID, req.ListingID,
		req.Amount, req.MarketplaceCommission, req.SellerAmount, req.Status,
	).Scan(&escrow.ID, &escrow.CreatedAt, &escrow.UpdatedAt)
	if err != nil {
		return nil, err
	}

	// Заполняем остальные поля
	escrow.PaymentTransactionID = req.PaymentTransactionID
	escrow.SellerID = req.SellerID
	escrow.BuyerID = req.BuyerID
	escrow.ListingID = req.ListingID
	escrow.Amount = req.Amount
	escrow.MarketplaceCommission = req.MarketplaceCommission
	escrow.SellerAmount = req.SellerAmount
	escrow.Status = req.Status

	return &escrow, nil
}

// GetEscrowByTransactionID получает escrow платеж по ID транзакции
func (r *PaymentRepository) GetEscrowByTransactionID(ctx context.Context, transactionID int64) (*models.EscrowPayment, error) {
	query := `
		SELECT id, payment_transaction_id, seller_id, buyer_id, listing_id,
			amount, marketplace_commission, seller_amount, status, release_date,
			created_at, updated_at
		FROM escrow_payments
		WHERE payment_transaction_id = $1`

	var escrow models.EscrowPayment
	err := r.db.QueryRowContext(ctx, query, transactionID).Scan(
		&escrow.ID, &escrow.PaymentTransactionID, &escrow.SellerID, &escrow.BuyerID,
		&escrow.ListingID, &escrow.Amount, &escrow.MarketplaceCommission,
		&escrow.SellerAmount, &escrow.Status, &escrow.ReleaseDate,
		&escrow.CreatedAt, &escrow.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &escrow, nil
}

// ReleaseEscrow освобождает escrow платеж
func (r *PaymentRepository) ReleaseEscrow(ctx context.Context, escrowID int64) error {
	query := `
		UPDATE escrow_payments SET
			status = 'released',
			release_date = NOW(),
			updated_at = NOW()
		WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, escrowID)
	return err
}

// CreatePayout создает выплату продавцу
func (r *PaymentRepository) CreatePayout(ctx context.Context, req CreatePayoutRequest) (*models.MerchantPayout, error) {
	bankAccountJSON, _ := json.Marshal(req.BankAccountInfo)

	query := `
		INSERT INTO merchant_payouts (
			seller_id, gateway_id, amount, currency,
			bank_account_info, status, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, NOW(), NOW()
		) RETURNING id, created_at, updated_at`

	var payout models.MerchantPayout
	err := r.db.QueryRowContext(ctx, query,
		req.SellerID, req.GatewayID, req.Amount, req.Currency,
		string(bankAccountJSON), req.Status,
	).Scan(&payout.ID, &payout.CreatedAt, &payout.UpdatedAt)
	if err != nil {
		return nil, err
	}

	// Заполняем остальные поля
	payout.SellerID = req.SellerID
	payout.GatewayID = req.GatewayID
	payout.Amount = req.Amount
	payout.Currency = req.Currency
	payout.BankAccountInfo = req.BankAccountInfo
	payout.Status = req.Status

	return &payout, nil
}

// GetPayoutsBySellerID получает выплаты для продавца
func (r *PaymentRepository) GetPayoutsBySellerID(ctx context.Context, sellerID int) ([]*models.MerchantPayout, error) {
	query := `
		SELECT id, seller_id, gateway_id, amount, currency,
			gateway_payout_id, gateway_reference_id, status,
			bank_account_info, gateway_response, error_details,
			created_at, updated_at, processed_at
		FROM merchant_payouts
		WHERE seller_id = $1
		ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, sellerID)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			// Логирование ошибки закрытия rows
		}
	}()

	var payouts []*models.MerchantPayout
	for rows.Next() {
		var payout models.MerchantPayout
		var bankAccountInfo, gatewayResponse, errorDetails sql.NullString

		err := rows.Scan(
			&payout.ID, &payout.SellerID, &payout.GatewayID, &payout.Amount,
			&payout.Currency, &payout.GatewayPayoutID, &payout.GatewayReferenceID,
			&payout.Status, &bankAccountInfo, &gatewayResponse, &errorDetails,
			&payout.CreatedAt, &payout.UpdatedAt, &payout.ProcessedAt,
		)
		if err != nil {
			return nil, err
		}

		// Парсим JSON поля
		if bankAccountInfo.Valid {
			json.Unmarshal([]byte(bankAccountInfo.String), &payout.BankAccountInfo)
		}
		if gatewayResponse.Valid {
			json.Unmarshal([]byte(gatewayResponse.String), &payout.GatewayResponse)
		}
		if errorDetails.Valid {
			json.Unmarshal([]byte(errorDetails.String), &payout.ErrorDetails)
		}

		payouts = append(payouts, &payout)
	}

	return payouts, nil
}

// generateOrderReference генерирует уникальную ссылку на заказ
func generateOrderReference(userID, listingID int) string {
	timestamp := time.Now().Unix()
	return fmt.Sprintf("ORD-%d-%d-%d", userID, listingID, timestamp)
}
