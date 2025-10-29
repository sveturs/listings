package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	pkgErrors "github.com/pkg/errors"

	"backend/internal/domain/models"
)

// ErrMarketplaceOrderNotFound возвращается когда заказ маркетплейса не найден
var ErrMarketplaceOrderNotFound = errors.New("marketplace order not found")

// MarketplaceOrderRepository предоставляет методы для работы с заказами маркетплейса
type MarketplaceOrderRepository struct {
	pool *pgxpool.Pool
}

// NewMarketplaceOrderRepository создает новый репозиторий заказов маркетплейса
func NewMarketplaceOrderRepository(pool *pgxpool.Pool) *MarketplaceOrderRepository {
	return &MarketplaceOrderRepository{pool: pool}
}

// Create создает новый заказ
func (r *MarketplaceOrderRepository) Create(ctx context.Context, order *models.MarketplaceOrder) error {
	query := `
		INSERT INTO marketplace_orders (
			buyer_id, seller_id, listing_id,
			item_price, platform_fee_rate, platform_fee_amount, seller_payout_amount,
			payment_transaction_id, status, protection_period_days,
			shipping_method
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
		) RETURNING id, created_at, updated_at`

	err := r.pool.QueryRow(ctx, query,
		order.BuyerID, order.SellerID, order.ListingID,
		order.ItemPrice, order.PlatformFeeRate, order.PlatformFeeAmount, order.SellerPayoutAmount,
		order.PaymentTransactionID, order.Status, order.ProtectionPeriodDays,
		order.ShippingMethod,
	).Scan(&order.ID, &order.CreatedAt, &order.UpdatedAt)
	if err != nil {
		return pkgErrors.Wrap(err, "failed to create marketplace order")
	}

	return nil
}

// GetByID получает заказ по ID
func (r *MarketplaceOrderRepository) GetByID(ctx context.Context, id int64) (*models.MarketplaceOrder, error) {
	var order models.MarketplaceOrder
	query := `
		SELECT 
			id, buyer_id, seller_id, listing_id,
			item_price, platform_fee_rate, platform_fee_amount, seller_payout_amount,
			payment_transaction_id, status, protection_period_days, protection_expires_at,
			shipping_method, tracking_number, shipped_at, delivered_at,
			created_at, updated_at
		FROM marketplace_orders
		WHERE id = $1`

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&order.ID, &order.BuyerID, &order.SellerID, &order.ListingID,
		&order.ItemPrice, &order.PlatformFeeRate, &order.PlatformFeeAmount, &order.SellerPayoutAmount,
		&order.PaymentTransactionID, &order.Status, &order.ProtectionPeriodDays, &order.ProtectionExpiresAt,
		&order.ShippingMethod, &order.TrackingNumber, &order.ShippedAt, &order.DeliveredAt,
		&order.CreatedAt, &order.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("marketplace order not found")
		}
		return nil, pkgErrors.Wrap(err, "failed to get marketplace order")
	}

	return &order, nil
}

// GetByPaymentTransactionID получает заказ по ID транзакции
func (r *MarketplaceOrderRepository) GetByPaymentTransactionID(ctx context.Context, transactionID int64) (*models.MarketplaceOrder, error) {
	var order models.MarketplaceOrder
	query := `
		SELECT 
			id, buyer_id, seller_id, listing_id,
			item_price, platform_fee_rate, platform_fee_amount, seller_payout_amount,
			payment_transaction_id, status, protection_period_days, protection_expires_at,
			shipping_method, tracking_number, shipped_at, delivered_at,
			created_at, updated_at
		FROM marketplace_orders
		WHERE payment_transaction_id = $1`

	err := r.pool.QueryRow(ctx, query, transactionID).Scan(
		&order.ID, &order.BuyerID, &order.SellerID, &order.ListingID,
		&order.ItemPrice, &order.PlatformFeeRate, &order.PlatformFeeAmount, &order.SellerPayoutAmount,
		&order.PaymentTransactionID, &order.Status, &order.ProtectionPeriodDays, &order.ProtectionExpiresAt,
		&order.ShippingMethod, &order.TrackingNumber, &order.ShippedAt, &order.DeliveredAt,
		&order.CreatedAt, &order.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrMarketplaceOrderNotFound
		}
		return nil, pkgErrors.Wrap(err, "failed to get marketplace order by transaction id")
	}

	return &order, nil
}

// UpdateStatus обновляет статус заказа с записью в историю
func (r *MarketplaceOrderRepository) UpdateStatus(ctx context.Context, orderID int64, newStatus models.MarketplaceOrderStatus, reason string, userID *int64) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return pkgErrors.Wrap(err, "failed to begin transaction")
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil {
			// Игнорируем ошибку если транзакция уже была завершена
			_ = err // Explicitly ignore error
		}
	}()

	// Получаем текущий статус
	var currentStatus string
	err = tx.QueryRow(ctx,
		"SELECT status FROM marketplace_orders WHERE id = $1 FOR UPDATE", orderID).Scan(&currentStatus)
	if err != nil {
		return pkgErrors.Wrap(err, "failed to get current status")
	}

	// Проверяем возможность перехода
	oldStatus := models.MarketplaceOrderStatus(currentStatus)
	if !oldStatus.CanTransitionTo(newStatus) {
		return fmt.Errorf("invalid status transition from %s to %s", oldStatus, newStatus)
	}

	// Обновляем статус
	_, err = tx.Exec(ctx,
		"UPDATE marketplace_orders SET status = $1, updated_at = NOW() WHERE id = $2",
		newStatus, orderID)
	if err != nil {
		return pkgErrors.Wrap(err, "failed to update status")
	}

	// Записываем в историю
	_, err = tx.Exec(ctx, `
		INSERT INTO order_status_history (order_id, old_status, new_status, reason, created_by)
		VALUES ($1, $2, $3, $4, $5)`,
		orderID, currentStatus, newStatus, reason, userID)
	if err != nil {
		return pkgErrors.Wrap(err, "failed to insert status history")
	}

	// Обновляем специфичные поля в зависимости от статуса
	switch newStatus {
	case models.MarketplaceOrderStatusShipped:
		_, err = tx.Exec(ctx,
			"UPDATE marketplace_orders SET shipped_at = NOW() WHERE id = $1", orderID)
	case models.MarketplaceOrderStatusDelivered:
		_, err = tx.Exec(ctx, `
			UPDATE marketplace_orders 
			SET delivered_at = NOW(), 
			    protection_expires_at = NOW() + INTERVAL '7 days' 
			WHERE id = $1`, orderID)
	case models.MarketplaceOrderStatusPending,
		models.MarketplaceOrderStatusPaid,
		models.MarketplaceOrderStatusCompleted,
		models.MarketplaceOrderStatusDisputed,
		models.MarketplaceOrderStatusCancelled,
		models.MarketplaceOrderStatusRefunded:
		// Для этих статусов нет специфичных полей для обновления
	}
	if err != nil {
		return pkgErrors.Wrap(err, "failed to update status-specific fields")
	}

	return tx.Commit(ctx)
}

// GetOrdersForAutoCapture получает заказы готовые к автоматическому capture
func (r *MarketplaceOrderRepository) GetOrdersForAutoCapture(ctx context.Context) ([]*models.MarketplaceOrder, error) {
	query := `
		SELECT 
			o.id, o.buyer_id, o.seller_id, o.listing_id,
			o.item_price, o.platform_fee_rate, o.platform_fee_amount, o.seller_payout_amount,
			o.payment_transaction_id, o.status, o.protection_period_days, o.protection_expires_at,
			o.shipping_method, o.tracking_number, o.shipped_at, o.delivered_at,
			o.created_at, o.updated_at
		FROM marketplace_orders o
		INNER JOIN payment_transactions pt ON o.payment_transaction_id = pt.id
		WHERE 
			pt.status = 'authorized' 
			AND pt.capture_mode = 'manual'
			AND (
				-- Заказ завершен
				(o.status = 'completed')
				OR 
				-- Или доставлен и истек защитный период
				(o.status = 'delivered' AND o.protection_expires_at <= NOW())
			)`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, pkgErrors.Wrap(err, "failed to get orders for auto capture")
	}
	defer rows.Close()

	var orders []*models.MarketplaceOrder
	for rows.Next() {
		var order models.MarketplaceOrder
		err := rows.Scan(
			&order.ID, &order.BuyerID, &order.SellerID, &order.ListingID,
			&order.ItemPrice, &order.PlatformFeeRate, &order.PlatformFeeAmount, &order.SellerPayoutAmount,
			&order.PaymentTransactionID, &order.Status, &order.ProtectionPeriodDays, &order.ProtectionExpiresAt,
			&order.ShippingMethod, &order.TrackingNumber, &order.ShippedAt, &order.DeliveredAt,
			&order.CreatedAt, &order.UpdatedAt,
		)
		if err != nil {
			return nil, pkgErrors.Wrap(err, "failed to scan order")
		}
		orders = append(orders, &order)
	}

	if err = rows.Err(); err != nil {
		return nil, pkgErrors.Wrap(err, "error during rows iteration")
	}

	return orders, nil
}

// GetBuyerOrders получает заказы покупателя
func (r *MarketplaceOrderRepository) GetBuyerOrders(ctx context.Context, buyerID int64, limit, offset int) ([]*models.MarketplaceOrder, int, error) {
	query := `
		SELECT 
			id, buyer_id, seller_id, listing_id,
			item_price, platform_fee_rate, platform_fee_amount, seller_payout_amount,
			payment_transaction_id, status, protection_period_days, protection_expires_at,
			shipping_method, tracking_number, shipped_at, delivered_at,
			created_at, updated_at
		FROM marketplace_orders
		WHERE buyer_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.pool.Query(ctx, query, buyerID, limit, offset)
	if err != nil {
		return nil, 0, pkgErrors.Wrap(err, "failed to get buyer orders")
	}
	defer rows.Close()

	var orders []*models.MarketplaceOrder
	for rows.Next() {
		var order models.MarketplaceOrder
		err := rows.Scan(
			&order.ID, &order.BuyerID, &order.SellerID, &order.ListingID,
			&order.ItemPrice, &order.PlatformFeeRate, &order.PlatformFeeAmount, &order.SellerPayoutAmount,
			&order.PaymentTransactionID, &order.Status, &order.ProtectionPeriodDays, &order.ProtectionExpiresAt,
			&order.ShippingMethod, &order.TrackingNumber, &order.ShippedAt, &order.DeliveredAt,
			&order.CreatedAt, &order.UpdatedAt,
		)
		if err != nil {
			return nil, 0, pkgErrors.Wrap(err, "failed to scan order")
		}
		orders = append(orders, &order)
	}

	// Получаем общее количество
	var total int
	err = r.pool.QueryRow(ctx,
		"SELECT COUNT(*) FROM marketplace_orders WHERE buyer_id = $1", buyerID).Scan(&total)
	if err != nil {
		return nil, 0, pkgErrors.Wrap(err, "failed to get total count")
	}

	return orders, total, nil
}

// GetSellerOrders получает заказы продавца
func (r *MarketplaceOrderRepository) GetSellerOrders(ctx context.Context, sellerID int64, limit, offset int) ([]*models.MarketplaceOrder, int, error) {
	query := `
		SELECT 
			id, buyer_id, seller_id, listing_id,
			item_price, platform_fee_rate, platform_fee_amount, seller_payout_amount,
			payment_transaction_id, status, protection_period_days, protection_expires_at,
			shipping_method, tracking_number, shipped_at, delivered_at,
			created_at, updated_at
		FROM marketplace_orders
		WHERE seller_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.pool.Query(ctx, query, sellerID, limit, offset)
	if err != nil {
		return nil, 0, pkgErrors.Wrap(err, "failed to get seller orders")
	}
	defer rows.Close()

	var orders []*models.MarketplaceOrder
	for rows.Next() {
		var order models.MarketplaceOrder
		err := rows.Scan(
			&order.ID, &order.BuyerID, &order.SellerID, &order.ListingID,
			&order.ItemPrice, &order.PlatformFeeRate, &order.PlatformFeeAmount, &order.SellerPayoutAmount,
			&order.PaymentTransactionID, &order.Status, &order.ProtectionPeriodDays, &order.ProtectionExpiresAt,
			&order.ShippingMethod, &order.TrackingNumber, &order.ShippedAt, &order.DeliveredAt,
			&order.CreatedAt, &order.UpdatedAt,
		)
		if err != nil {
			return nil, 0, pkgErrors.Wrap(err, "failed to scan order")
		}
		orders = append(orders, &order)
	}

	// Получаем общее количество
	var total int
	err = r.pool.QueryRow(ctx,
		"SELECT COUNT(*) FROM marketplace_orders WHERE seller_id = $1", sellerID).Scan(&total)
	if err != nil {
		return nil, 0, pkgErrors.Wrap(err, "failed to get total count")
	}

	return orders, total, nil
}

// UpdateShippingInfo обновляет информацию о доставке (только shipping_method, tracking_number устанавливается через delivery service)
func (r *MarketplaceOrderRepository) UpdateShippingInfo(ctx context.Context, orderID int64, shippingMethod string) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE marketplace_orders
		SET shipping_method = $1, updated_at = NOW()
		WHERE id = $2`,
		shippingMethod, orderID)
	if err != nil {
		return pkgErrors.Wrap(err, "failed to update shipping info")
	}

	return nil
}

// AddMessage добавляет сообщение к заказу
func (r *MarketplaceOrderRepository) AddMessage(ctx context.Context, message *models.OrderMessage) error {
	query := `
		INSERT INTO order_messages (
			order_id, sender_id, message_type, content, metadata
		) VALUES (
			$1, $2, $3, $4, $5
		) RETURNING id, created_at`

	err := r.pool.QueryRow(ctx, query,
		message.OrderID, message.SenderID, message.MessageType,
		message.Content, message.Metadata,
	).Scan(&message.ID, &message.CreatedAt)
	if err != nil {
		return pkgErrors.Wrap(err, "failed to add message")
	}

	return nil
}

// GetOrderMessages получает сообщения заказа
func (r *MarketplaceOrderRepository) GetOrderMessages(ctx context.Context, orderID int64) ([]*models.OrderMessage, error) {
	query := `
		SELECT 
			id, order_id, sender_id, message_type, content, metadata, created_at
		FROM order_messages
		WHERE order_id = $1
		ORDER BY created_at ASC`

	rows, err := r.pool.Query(ctx, query, orderID)
	if err != nil {
		return nil, pkgErrors.Wrap(err, "failed to get order messages")
	}
	defer rows.Close()

	var messages []*models.OrderMessage
	for rows.Next() {
		var message models.OrderMessage
		err := rows.Scan(
			&message.ID, &message.OrderID, &message.SenderID,
			&message.MessageType, &message.Content, &message.Metadata,
			&message.CreatedAt,
		)
		if err != nil {
			return nil, pkgErrors.Wrap(err, "failed to scan message")
		}
		messages = append(messages, &message)
	}

	if err = rows.Err(); err != nil {
		return nil, pkgErrors.Wrap(err, "error during rows iteration")
	}

	return messages, nil
}
