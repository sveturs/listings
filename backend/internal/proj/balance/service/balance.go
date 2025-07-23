// backend/internal/proj/balance/service/balance.go

package balance

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"backend/internal/domain/models"
	"backend/internal/storage"
)

type BalanceService struct {
	storage storage.Storage
	logger  *slog.Logger
}

func NewBalanceService(storage storage.Storage) *BalanceService {
	return &BalanceService{
		storage: storage,
		logger:  slog.Default(),
	}
}

func (s *BalanceService) GetBalance(ctx context.Context, userID int) (*models.UserBalance, error) {
	var balance models.UserBalance
	err := s.storage.QueryRow(ctx, `
        WITH new_balance AS (
            INSERT INTO user_balances (user_id, balance, frozen_balance, currency)
            VALUES ($1, 0, 0, 'RSD')
            ON CONFLICT (user_id) DO NOTHING
            RETURNING user_id, balance, frozen_balance, currency, updated_at
        )
        SELECT user_id, balance, frozen_balance, currency, updated_at
        FROM new_balance
        UNION ALL
        SELECT user_id, balance, frozen_balance, currency, updated_at
        FROM user_balances
        WHERE user_id = $1
        LIMIT 1
    `, userID).Scan(
		&balance.UserID,
		&balance.Balance,
		&balance.FrozenBalance,
		&balance.Currency,
		&balance.UpdatedAt,
	)
	if err != nil {
		return &models.UserBalance{
			UserID:   userID,
			Balance:  0,
			Currency: "RSD",
		}, nil
	}

	return &balance, nil
}

// backend/internal/proj/balance/service/balance.go

func (s *BalanceService) CreateDeposit(ctx context.Context, userID int, amount float64, method string) (*models.BalanceTransaction, error) {
	// Проверяем метод оплаты
	var paymentMethod models.PaymentMethod
	err := s.storage.QueryRow(ctx, `
        SELECT id, name, minimum_amount, maximum_amount, fee_percentage, fixed_fee
        FROM payment_methods
        WHERE code = $1 AND is_active = true
    `, method).Scan(
		&paymentMethod.ID,
		&paymentMethod.Name,
		&paymentMethod.MinimumAmount,
		&paymentMethod.MaximumAmount,
		&paymentMethod.FeePercentage,
		&paymentMethod.FixedFee,
	)
	if err != nil {
		return nil, fmt.Errorf("invalid payment method: %w", err)
	}

	// Проверяем лимиты
	if amount < paymentMethod.MinimumAmount {
		return nil, fmt.Errorf("amount %v is below minimum allowed %v", amount, paymentMethod.MinimumAmount)
	}
	if paymentMethod.MaximumAmount > 0 && amount > paymentMethod.MaximumAmount {
		return nil, fmt.Errorf("amount %v is above maximum allowed %v", amount, paymentMethod.MaximumAmount)
	}

	// Начинаем транзакцию
	tx, err := s.storage.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
			s.logger.Error("Failed to rollback transaction", "error", err)
		}
	}()

	// Рассчитываем комиссию
	fee := paymentMethod.FixedFee + (amount * paymentMethod.FeePercentage / 100)
	totalAmount := amount + fee

	// Создаем транзакцию сразу в статусе completed

	now := time.Now() // Создаем переменную
	transaction := &models.BalanceTransaction{
		UserID:        userID,
		Type:          "deposit",
		Amount:        totalAmount,
		Currency:      "RSD",
		Status:        "completed",
		PaymentMethod: method,
		Description:   fmt.Sprintf("Пополнение через %s", paymentMethod.Name),
		CreatedAt:     now,
		CompletedAt:   &now, // Используем адрес переменной
	}

	// Создаем запись о транзакции
	var transactionID int
	err = tx.QueryRow(ctx, `
        INSERT INTO balance_transactions
        (user_id, type, amount, currency, status, payment_method, description, created_at, completed_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
        RETURNING id
    `,
		transaction.UserID,
		transaction.Type,
		transaction.Amount,
		transaction.Currency,
		transaction.Status,
		transaction.PaymentMethod,
		transaction.Description,
		transaction.CreatedAt,
		transaction.CompletedAt,
	).Scan(&transactionID)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	// Обновляем баланс пользователя в той же транзакции
	_, err = tx.Exec(ctx, `
        INSERT INTO user_balances (user_id, balance, currency)
        VALUES ($1, $2, $3)
        ON CONFLICT (user_id) 
        DO UPDATE SET balance = user_balances.balance + EXCLUDED.balance
    `, userID, totalAmount, "RSD")
	if err != nil {
		return nil, fmt.Errorf("failed to update balance: %w", err)
	}

	// Фиксируем транзакцию
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	transaction.ID = transactionID
	return transaction, nil
}

func (s *BalanceService) ProcessDeposit(ctx context.Context, transactionID int) error {
	tx, err := s.storage.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
			s.logger.Error("Failed to rollback transaction", "error", err)
		}
	}()

	// Получаем транзакцию
	var transaction models.BalanceTransaction
	err = tx.QueryRow(ctx, `
        SELECT id, user_id, amount, status 
        FROM balance_transactions 
        WHERE id = $1 FOR UPDATE
    `, transactionID).Scan(
		&transaction.ID,
		&transaction.UserID,
		&transaction.Amount,
		&transaction.Status,
	)
	if err != nil {
		return fmt.Errorf("failed to get transaction: %w", err)
	}

	if transaction.Status != "pending" {
		return fmt.Errorf("invalid transaction status: %s", transaction.Status)
	}

	// Обновляем баланс и статус в одной транзакции
	_, err = tx.Exec(ctx, `
        WITH balance_update AS (
            UPDATE user_balances 
            SET balance = balance + $1, 
                updated_at = NOW()
            WHERE user_id = $2
        )
        UPDATE balance_transactions
        SET status = 'completed',
            completed_at = NOW()
        WHERE id = $3
    `, transaction.Amount, transaction.UserID, transactionID)
	if err != nil {
		return fmt.Errorf("failed to process transaction: %w", err)
	}

	return tx.Commit()
}

func (s *BalanceService) GetTransactions(ctx context.Context, userID int, limit, offset int) ([]models.BalanceTransaction, error) {
	rows, err := s.storage.Query(ctx, `
        SELECT id, user_id, type, amount, currency, status, payment_method, 
               description, created_at, completed_at
        FROM balance_transactions
        WHERE user_id = $1
        ORDER BY created_at DESC
        LIMIT $2 OFFSET $3
    `, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions: %w", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			s.logger.Error("Failed to close rows", "error", err)
		}
	}()

	var transactions []models.BalanceTransaction
	for rows.Next() {
		var tx models.BalanceTransaction
		err := rows.Scan(
			&tx.ID, &tx.UserID, &tx.Type, &tx.Amount, &tx.Currency,
			&tx.Status, &tx.PaymentMethod, &tx.Description,
			&tx.CreatedAt, &tx.CompletedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan transaction: %w", err)
		}
		transactions = append(transactions, tx)
	}

	return transactions, nil
}

func (s *BalanceService) GetPaymentMethods(ctx context.Context) ([]models.PaymentMethod, error) {
	rows, err := s.storage.Query(ctx, `
        SELECT id, name, code, type, is_active, minimum_amount, 
               maximum_amount, fee_percentage, fixed_fee, created_at
        FROM payment_methods
        WHERE is_active = true
        ORDER BY name
    `)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment methods: %w", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			s.logger.Error("Failed to close rows", "error", err)
		}
	}()

	var methods []models.PaymentMethod
	for rows.Next() {
		var method models.PaymentMethod
		err := rows.Scan(
			&method.ID, &method.Name, &method.Code, &method.Type,
			&method.IsActive, &method.MinimumAmount, &method.MaximumAmount,
			&method.FeePercentage, &method.FixedFee, &method.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan payment method: %w", err)
		}
		methods = append(methods, method)
	}

	return methods, nil
}
