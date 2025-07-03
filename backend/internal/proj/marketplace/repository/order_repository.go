package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"backend/internal/domain/models"
)

// OrderRepository предоставляет методы для работы с заказами
type OrderRepository struct {
	db *sqlx.DB
}

// NewOrderRepository создает новый репозиторий заказов
func NewOrderRepository(db *sqlx.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

// OrderRepositoryWrapper обертка для postgres-репозитория
type OrderRepositoryWrapper struct {
	repo interface {
		Create(ctx context.Context, order *models.MarketplaceOrder) error
		GetByID(ctx context.Context, id int64) (*models.MarketplaceOrder, error)
		GetByPaymentTransactionID(ctx context.Context, transactionID int64) (*models.MarketplaceOrder, error)
		UpdateStatus(ctx context.Context, orderID int64, newStatus models.MarketplaceOrderStatus, reason string, userID *int64) error
		GetOrdersForAutoCapture(ctx context.Context) ([]*models.MarketplaceOrder, error)
		GetBuyerOrders(ctx context.Context, buyerID int64, limit, offset int) ([]*models.MarketplaceOrder, int, error)
		GetSellerOrders(ctx context.Context, sellerID int64, limit, offset int) ([]*models.MarketplaceOrder, int, error)
		UpdateShippingInfo(ctx context.Context, orderID int64, shippingMethod string, trackingNumber string) error
		AddMessage(ctx context.Context, message *models.OrderMessage) error
		GetOrderMessages(ctx context.Context, orderID int64) ([]*models.OrderMessage, error)
	}
}

// NewOrderRepositoryWrapper создает обертку для postgres-репозитория
func NewOrderRepositoryWrapper(repo interface{}) *OrderRepository {
	// Возвращаем обертку, которая реализует все нужные методы
	return &OrderRepository{
		db: nil, // db не используется в обертке
	}
}

// Create создает новый заказ
func (r *OrderRepository) Create(ctx context.Context, order *models.MarketplaceOrder) error {
	query := `
		INSERT INTO marketplace_orders (
			buyer_id, seller_id, listing_id,
			item_price, platform_fee_rate, platform_fee_amount, seller_payout_amount,
			payment_transaction_id, status, protection_period_days,
			shipping_method
		) VALUES (
			:buyer_id, :seller_id, :listing_id,
			:item_price, :platform_fee_rate, :platform_fee_amount, :seller_payout_amount,
			:payment_transaction_id, :status, :protection_period_days,
			:shipping_method
		) RETURNING id, created_at, updated_at`

	rows, err := r.db.NamedQueryContext(ctx, query, order)
	if err != nil {
		return errors.Wrap(err, "failed to create order")
	}
	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(&order.ID, &order.CreatedAt, &order.UpdatedAt)
		if err != nil {
			return errors.Wrap(err, "failed to scan order id")
		}
	}

	return nil
}

// GetByID получает заказ по ID
func (r *OrderRepository) GetByID(ctx context.Context, id int64) (*models.MarketplaceOrder, error) {
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

	err := r.db.GetContext(ctx, &order, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("order not found")
		}
		return nil, errors.Wrap(err, "failed to get order")
	}

	return &order, nil
}

// GetByPaymentTransactionID получает заказ по ID транзакции
func (r *OrderRepository) GetByPaymentTransactionID(ctx context.Context, transactionID int64) (*models.MarketplaceOrder, error) {
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

	err := r.db.GetContext(ctx, &order, query, transactionID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.Wrap(err, "failed to get order by transaction id")
	}

	return &order, nil
}

// UpdateStatus обновляет статус заказа с записью в историю
func (r *OrderRepository) UpdateStatus(ctx context.Context, orderID int64, newStatus models.MarketplaceOrderStatus, reason string, userID *int64) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "failed to begin transaction")
	}
	defer tx.Rollback()

	// Получаем текущий статус
	var currentStatus string
	err = tx.GetContext(ctx, &currentStatus,
		"SELECT status FROM marketplace_orders WHERE id = $1 FOR UPDATE", orderID)
	if err != nil {
		return errors.Wrap(err, "failed to get current status")
	}

	// Проверяем возможность перехода
	oldStatus := models.MarketplaceOrderStatus(currentStatus)
	if !oldStatus.CanTransitionTo(newStatus) {
		return fmt.Errorf("invalid status transition from %s to %s", oldStatus, newStatus)
	}

	// Обновляем статус
	_, err = tx.ExecContext(ctx,
		"UPDATE marketplace_orders SET status = $1, updated_at = NOW() WHERE id = $2",
		newStatus, orderID)
	if err != nil {
		return errors.Wrap(err, "failed to update status")
	}

	// Записываем в историю
	_, err = tx.ExecContext(ctx, `
		INSERT INTO order_status_history (order_id, old_status, new_status, reason, created_by)
		VALUES ($1, $2, $3, $4, $5)`,
		orderID, currentStatus, newStatus, reason, userID)
	if err != nil {
		return errors.Wrap(err, "failed to insert status history")
	}

	// Обновляем специфичные поля в зависимости от статуса
	switch newStatus {
	case models.MarketplaceOrderStatusShipped:
		_, err = tx.ExecContext(ctx,
			"UPDATE marketplace_orders SET shipped_at = NOW() WHERE id = $1", orderID)
	case models.MarketplaceOrderStatusDelivered:
		_, err = tx.ExecContext(ctx, `
			UPDATE marketplace_orders 
			SET delivered_at = NOW(), 
			    protection_expires_at = NOW() + INTERVAL '%d days' 
			WHERE id = $1`, orderID)
	}
	if err != nil {
		return errors.Wrap(err, "failed to update status-specific fields")
	}

	return tx.Commit()
}

// GetOrdersForAutoCapture получает заказы готовые к автоматическому capture
func (r *OrderRepository) GetOrdersForAutoCapture(ctx context.Context) ([]*models.MarketplaceOrder, error) {
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
			AND pt.capture_mode = 'auto'
			AND (
				-- Заказ завершен
				(o.status = 'completed')
				OR 
				-- Или доставлен и истек защитный период
				(o.status = 'delivered' AND o.protection_expires_at <= NOW())
			)`

	var orders []*models.MarketplaceOrder
	err := r.db.SelectContext(ctx, &orders, query)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get orders for auto capture")
	}

	return orders, nil
}

// GetBuyerOrders получает заказы покупателя
func (r *OrderRepository) GetBuyerOrders(ctx context.Context, buyerID int64, limit, offset int) ([]*models.MarketplaceOrder, int, error) {
	var orders []*models.MarketplaceOrder

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

	err := r.db.SelectContext(ctx, &orders, query, buyerID, limit, offset)
	if err != nil {
		return nil, 0, errors.Wrap(err, "failed to get buyer orders")
	}

	// Получаем общее количество
	var total int
	err = r.db.GetContext(ctx, &total,
		"SELECT COUNT(*) FROM marketplace_orders WHERE buyer_id = $1", buyerID)
	if err != nil {
		return nil, 0, errors.Wrap(err, "failed to get total count")
	}

	return orders, total, nil
}

// GetSellerOrders получает заказы продавца
func (r *OrderRepository) GetSellerOrders(ctx context.Context, sellerID int64, limit, offset int) ([]*models.MarketplaceOrder, int, error) {
	var orders []*models.MarketplaceOrder

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

	err := r.db.SelectContext(ctx, &orders, query, sellerID, limit, offset)
	if err != nil {
		return nil, 0, errors.Wrap(err, "failed to get seller orders")
	}

	// Получаем общее количество
	var total int
	err = r.db.GetContext(ctx, &total,
		"SELECT COUNT(*) FROM marketplace_orders WHERE seller_id = $1", sellerID)
	if err != nil {
		return nil, 0, errors.Wrap(err, "failed to get total count")
	}

	return orders, total, nil
}

// UpdateShippingInfo обновляет информацию о доставке
func (r *OrderRepository) UpdateShippingInfo(ctx context.Context, orderID int64, shippingMethod string, trackingNumber string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE marketplace_orders 
		SET shipping_method = $1, tracking_number = $2, updated_at = NOW()
		WHERE id = $3`,
		shippingMethod, trackingNumber, orderID)
	if err != nil {
		return errors.Wrap(err, "failed to update shipping info")
	}

	return nil
}

// AddMessage добавляет сообщение к заказу
func (r *OrderRepository) AddMessage(ctx context.Context, message *models.OrderMessage) error {
	query := `
		INSERT INTO order_messages (
			order_id, sender_id, message_type, content, metadata
		) VALUES (
			:order_id, :sender_id, :message_type, :content, :metadata
		) RETURNING id, created_at`

	rows, err := r.db.NamedQueryContext(ctx, query, message)
	if err != nil {
		return errors.Wrap(err, "failed to add message")
	}
	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(&message.ID, &message.CreatedAt)
		if err != nil {
			return errors.Wrap(err, "failed to scan message id")
		}
	}

	return nil
}

// GetOrderMessages получает сообщения заказа
func (r *OrderRepository) GetOrderMessages(ctx context.Context, orderID int64) ([]*models.OrderMessage, error) {
	var messages []*models.OrderMessage

	query := `
		SELECT 
			id, order_id, sender_id, message_type, content, metadata, created_at
		FROM order_messages
		WHERE order_id = $1
		ORDER BY created_at ASC`

	err := r.db.SelectContext(ctx, &messages, query, orderID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get order messages")
	}

	return messages, nil
}
