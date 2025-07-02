// backend/internal/storage/postgres/balance.go

package postgres

import (
	"context"
	"log"

	"backend/internal/domain/models"
)

func (db *Database) GetUserBalance(ctx context.Context, userID int) (*models.UserBalance, error) {
	log.Printf("Getting balance from database for user %d", userID)

	var balance models.UserBalance
	err := db.pool.QueryRow(ctx, `
        WITH inserted AS (
            INSERT INTO user_balances (user_id, balance, frozen_balance, currency)
            VALUES ($1, 0, 0, 'RSD')
            ON CONFLICT (user_id) DO NOTHING
            RETURNING user_id, balance, frozen_balance, currency, updated_at
        )
        SELECT user_id, balance, frozen_balance, currency, updated_at
        FROM inserted
        UNION ALL
        SELECT user_id, balance, frozen_balance, currency, updated_at
        FROM user_balances
        WHERE user_id = $1
        AND NOT EXISTS (SELECT 1 FROM inserted)
        LIMIT 1
    `, userID).Scan(
		&balance.UserID,
		&balance.Balance,
		&balance.FrozenBalance,
		&balance.Currency,
		&balance.UpdatedAt,
	)
	if err != nil {
		log.Printf("Error getting balance: %v", err)
		return nil, err
	}

	log.Printf("Successfully got balance: %+v", balance)
	return &balance, nil
}

func (db *Database) GetActivePaymentMethods(ctx context.Context) ([]models.PaymentMethod, error) {
	rows, err := db.pool.Query(ctx, `
        SELECT id, name, code, type, minimum_amount, maximum_amount, fee_percentage, fixed_fee 
        FROM payment_methods 
        WHERE is_active = true
        ORDER BY name
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var methods []models.PaymentMethod
	for rows.Next() {
		var method models.PaymentMethod
		if err := rows.Scan(
			&method.ID,
			&method.Name,
			&method.Code,
			&method.Type,
			&method.MinimumAmount,
			&method.MaximumAmount,
			&method.FeePercentage,
			&method.FixedFee,
		); err != nil {
			return nil, err
		}
		methods = append(methods, method)
	}

	return methods, nil
}

func (db *Database) GetUserTransactions(ctx context.Context, userID int, limit, offset int) ([]models.BalanceTransaction, error) {
	rows, err := db.pool.Query(ctx, `
        SELECT id, user_id, type, amount, currency, status, payment_method,
               description, created_at, completed_at
        FROM balance_transactions
        WHERE user_id = $1
        ORDER BY created_at DESC
        LIMIT $2 OFFSET $3
    `, userID, limit, offset)
	if err != nil {
		return []models.BalanceTransaction{}, err // Возвращаем пустой массив при ошибке
	}
	defer rows.Close()

	var transactions []models.BalanceTransaction
	for rows.Next() {
		var tx models.BalanceTransaction
		if err := rows.Scan(
			&tx.ID,
			&tx.UserID,
			&tx.Type,
			&tx.Amount,
			&tx.Currency,
			&tx.Status,
			&tx.PaymentMethod,
			&tx.Description,
			&tx.CreatedAt,
			&tx.CompletedAt,
		); err != nil {
			return []models.BalanceTransaction{}, err
		}
		transactions = append(transactions, tx)
	}

	if err = rows.Err(); err != nil {
		return []models.BalanceTransaction{}, err
	}

	// Если транзакций нет, вернем пустой массив
	if len(transactions) == 0 {
		return []models.BalanceTransaction{}, nil
	}

	return transactions, nil
}

func (db *Database) CreateTransaction(ctx context.Context, transaction *models.BalanceTransaction) (int, error) {
	var id int
	err := db.pool.QueryRow(ctx, `
        INSERT INTO balance_transactions (
            user_id, type, amount, currency, status, 
            payment_method, description, created_at
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
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
	).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (db *Database) UpdateBalance(ctx context.Context, userID int, amount float64) error {
	result, err := db.pool.Exec(ctx, `
        UPDATE user_balances 
        SET balance = balance + $1, 
            updated_at = CURRENT_TIMESTAMP
        WHERE user_id = $2
    `, amount, userID)
	if err != nil {
		return err
	}

	affected := result.RowsAffected()
	if affected == 0 {
		// Если записи нет, создаем новую
		_, err = db.pool.Exec(ctx, `
            INSERT INTO user_balances (user_id, balance, currency)
            VALUES ($1, $2, 'RSD')
        `, userID, amount)
		if err != nil {
			return err
		}
	}

	return nil
}
