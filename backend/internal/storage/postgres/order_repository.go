package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"backend/internal/domain/models"
)

// OrderRepositoryInterface определяет интерфейс для работы с заказами
type OrderRepositoryInterface interface {
	Create(ctx context.Context, order *models.StorefrontOrder) (*models.StorefrontOrder, error)
	GetByID(ctx context.Context, orderID int64) (*models.StorefrontOrder, error)
	GetByOrderNumber(ctx context.Context, orderNumber string) (*models.StorefrontOrder, error)
	Update(ctx context.Context, order *models.StorefrontOrder) error
	List(ctx context.Context, filter models.OrderFilter) ([]models.StorefrontOrder, int, error)
	UpdateStatus(ctx context.Context, orderID int64, status models.OrderStatus, metadata map[string]interface{}) error
	GetUserOrders(ctx context.Context, userID int, limit, offset int) ([]models.StorefrontOrder, int, error)
	GetStorefrontOrders(ctx context.Context, storefrontID int, filter models.OrderFilter) ([]models.StorefrontOrder, int, error)
}

// orderRepository реализует интерфейс для работы с заказами
type orderRepository struct {
	pool *pgxpool.Pool
}

// NewOrderRepository создает новый репозиторий заказов
func NewOrderRepository(pool *pgxpool.Pool) OrderRepositoryInterface {
	return &orderRepository{pool: pool}
}

// Create создает новый заказ
func (r *orderRepository) Create(ctx context.Context, order *models.StorefrontOrder) (*models.StorefrontOrder, error) {
	// Конвертируем shipping и billing address в JSON
	shippingJSON, err := json.Marshal(order.ShippingAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal shipping address: %w", err)
	}

	billingJSON, err := json.Marshal(order.BillingAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal billing address: %w", err)
	}

	// Конвертируем metadata в JSON
	metadataJSON, err := json.Marshal(order.Metadata)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal metadata: %w", err)
	}

	// Начинаем транзакцию
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Создаем заказ (без items - они будут в отдельной таблице)
	query := `
		INSERT INTO storefront_orders (
			storefront_id, customer_id, status,
			subtotal_amount, tax_amount, shipping_amount, discount, total_amount,
			commission_amount, seller_amount, currency, 
			shipping_address, billing_address,
			payment_method, payment_status, payment_transaction_id,
			customer_notes, metadata
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11,
			$12::jsonb, $13::jsonb, $14, $15, $16, $17, $18::jsonb
		) RETURNING id, order_number, created_at, updated_at`

	err = tx.QueryRow(ctx, query,
		order.StorefrontID,
		order.UserID,
		order.Status,
		order.Subtotal.InexactFloat64(),
		order.Tax.InexactFloat64(),
		order.Shipping.InexactFloat64(),
		order.Discount.InexactFloat64(),
		order.Total.InexactFloat64(),
		order.CommissionAmount.InexactFloat64(),
		order.SellerAmount.InexactFloat64(),
		order.Currency,
		shippingJSON,
		billingJSON,
		order.PaymentMethod,
		order.PaymentStatus,
		order.PaymentTransactionID,
		order.CustomerNotes,
		metadataJSON,
	).Scan(&order.ID, &order.OrderNumber, &order.CreatedAt, &order.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	// Создаем позиции заказа в отдельной таблице
	if len(order.Items) > 0 {
		for _, item := range order.Items {
			itemQuery := `
				INSERT INTO storefront_order_items (
					order_id, product_id, variant_id,
					product_name, product_sku, variant_name,
					quantity, price_per_unit, total_price,
					product_attributes
				) VALUES (
					$1, $2, $3, $4, $5, $6, $7, $8, $9, $10::jsonb
				)`

			// Подготавливаем данные для вставки
			productAttrs, _ := json.Marshal(item.ProductAttributes)

			_, err = tx.Exec(ctx, itemQuery,
				order.ID,
				item.ProductID,
				item.VariantID,
				item.ProductName,
				item.ProductSKU,
				item.VariantName,
				item.Quantity,
				item.PricePerUnit.InexactFloat64(),
				item.TotalPrice.InexactFloat64(),
				productAttrs,
			)
			if err != nil {
				return nil, fmt.Errorf("failed to create order item: %w", err)
			}
		}
	}

	// Коммитим транзакцию
	if err = tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return order, nil
}

// GetByID получает заказ по ID
func (r *orderRepository) GetByID(ctx context.Context, orderID int64) (*models.StorefrontOrder, error) {
	// Получаем основную информацию о заказе
	query := `
		SELECT 
			id, order_number, storefront_id, customer_id, status,
			subtotal_amount, tax_amount, shipping_amount, discount, total_amount, 
			commission_amount, seller_amount, currency,
			shipping_address, billing_address,
			payment_method, payment_status, payment_transaction_id,
			customer_notes, metadata, created_at, updated_at,
			confirmed_at, shipped_at, delivered_at, cancelled_at
		FROM storefront_orders
		WHERE id = $1`

	var order models.StorefrontOrder
	var shippingJSON, billingJSON, metadataJSON json.RawMessage
	var customerNotes, paymentTransactionID sql.NullString
	var confirmedAt, shippedAt, deliveredAt, cancelledAt sql.NullTime

	err := r.pool.QueryRow(ctx, query, orderID).Scan(
		&order.ID,
		&order.OrderNumber,
		&order.StorefrontID,
		&order.UserID,
		&order.Status,
		&order.Subtotal,
		&order.Tax,
		&order.Shipping,
		&order.Discount,
		&order.Total,
		&order.CommissionAmount,
		&order.SellerAmount,
		&order.Currency,
		&shippingJSON,
		&billingJSON,
		&order.PaymentMethod,
		&order.PaymentStatus,
		&paymentTransactionID,
		&customerNotes,
		&metadataJSON,
		&order.CreatedAt,
		&order.UpdatedAt,
		&confirmedAt,
		&shippedAt,
		&deliveredAt,
		&cancelledAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("order not found")
		}
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	// Обработка NULL значений
	if customerNotes.Valid {
		order.CustomerNotes = &customerNotes.String
	}
	if paymentTransactionID.Valid {
		order.PaymentTransactionID = &paymentTransactionID.String
	}
	if confirmedAt.Valid {
		order.ConfirmedAt = &confirmedAt.Time
	}
	if shippedAt.Valid {
		order.ShippedAt = &shippedAt.Time
	}
	if deliveredAt.Valid {
		order.DeliveredAt = &deliveredAt.Time
	}
	if cancelledAt.Valid {
		order.CancelledAt = &cancelledAt.Time
	}

	// Парсим JSON поля
	if err := json.Unmarshal(shippingJSON, &order.ShippingAddress); err != nil {
		return nil, fmt.Errorf("failed to unmarshal shipping address: %w", err)
	}
	if err := json.Unmarshal(billingJSON, &order.BillingAddress); err != nil {
		return nil, fmt.Errorf("failed to unmarshal billing address: %w", err)
	}
	if metadataJSON != nil {
		if err := json.Unmarshal(metadataJSON, &order.Metadata); err != nil {
			return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
		}
	}

	// Загружаем позиции заказа из отдельной таблицы
	itemsQuery := `
		SELECT 
			id, product_id, variant_id,
			product_name, product_sku, variant_name,
			quantity, price_per_unit, total_price,
			product_attributes
		FROM storefront_order_items
		WHERE order_id = $1
		ORDER BY id`

	rows, err := r.pool.Query(ctx, itemsQuery, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order items: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var item models.StorefrontOrderItem
		var productAttrsJSON json.RawMessage
		var variantID sql.NullInt64
		var productSKU, variantName sql.NullString

		err := rows.Scan(
			&item.ID,
			&item.ProductID,
			&variantID,
			&item.ProductName,
			&productSKU,
			&variantName,
			&item.Quantity,
			&item.PricePerUnit,
			&item.TotalPrice,
			&productAttrsJSON,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order item: %w", err)
		}

		if variantID.Valid {
			id := variantID.Int64
			item.VariantID = &id
		}
		if productSKU.Valid {
			item.ProductSKU = &productSKU.String
		}
		if variantName.Valid {
			item.VariantName = &variantName.String
		}
		if productAttrsJSON != nil {
			if err := json.Unmarshal(productAttrsJSON, &item.ProductAttributes); err != nil {
				return nil, fmt.Errorf("failed to unmarshal product attributes: %w", err)
			}
		}

		order.Items = append(order.Items, item)
	}

	return &order, nil
}

// GetByOrderNumber получает заказ по номеру
func (r *orderRepository) GetByOrderNumber(ctx context.Context, orderNumber string) (*models.StorefrontOrder, error) {
	// Просто переиспользуем GetByID после получения ID по номеру
	var orderID int64
	err := r.pool.QueryRow(ctx, "SELECT id FROM storefront_orders WHERE order_number = $1", orderNumber).Scan(&orderID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("order not found")
		}
		return nil, fmt.Errorf("failed to get order by number: %w", err)
	}

	return r.GetByID(ctx, orderID)
}

// Update обновляет заказ
func (r *orderRepository) Update(ctx context.Context, order *models.StorefrontOrder) error {
	// Конвертируем JSON поля
	itemsJSON, _ := json.Marshal(order.Items)
	shippingJSON, _ := json.Marshal(order.ShippingAddress)
	billingJSON, _ := json.Marshal(order.BillingAddress)
	metadataJSON, _ := json.Marshal(order.Metadata)

	query := `
		UPDATE storefront_orders SET
			status = $1,
			subtotal_amount = $2,
			tax_amount = $3,
			shipping_amount = $4,
			discount = $5,
			total_amount = $6,
			commission_amount = $7,
			seller_amount = $8,
			items = $9::jsonb,
			shipping_address = $10::jsonb,
			billing_address = $11::jsonb,
			payment_method = $12,
			payment_status = $13,
			payment_transaction_id = $14,
			notes = $15,
			metadata = $16::jsonb,
			confirmed_at = $17,
			shipped_at = $18,
			delivered_at = $19,
			cancelled_at = $20,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $21`

	_, err := r.pool.Exec(ctx, query,
		order.Status,
		order.Subtotal.InexactFloat64(),
		order.Tax.InexactFloat64(),
		order.Shipping.InexactFloat64(),
		order.Discount.InexactFloat64(),
		order.Total.InexactFloat64(),
		order.CommissionAmount.InexactFloat64(),
		order.SellerAmount.InexactFloat64(),
		itemsJSON,
		shippingJSON,
		billingJSON,
		order.PaymentMethod,
		order.PaymentStatus,
		order.PaymentTransactionID,
		order.Notes,
		metadataJSON,
		order.ConfirmedAt,
		order.ShippedAt,
		order.DeliveredAt,
		order.CancelledAt,
		order.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update order: %w", err)
	}

	return nil
}

// List возвращает список заказов с фильтрацией и пагинацией
func (r *orderRepository) List(ctx context.Context, filter models.OrderFilter) ([]models.StorefrontOrder, int, error) {
	// Базовый запрос
	baseQuery := `FROM storefront_orders WHERE 1=1`
	countQuery := `SELECT COUNT(*) ` + baseQuery
	selectQuery := `SELECT id, order_number, storefront_id, customer_id, status,
		subtotal_amount, tax_amount, shipping_amount, discount, total_amount, currency,
		created_at, updated_at ` + baseQuery

	var args []interface{}
	var conditions []string
	argCount := 1

	// Добавляем условия фильтрации
	if filter.StorefrontID != nil {
		conditions = append(conditions, fmt.Sprintf("storefront_id = $%d", argCount))
		args = append(args, *filter.StorefrontID)
		argCount++
	}

	if filter.UserID != nil {
		conditions = append(conditions, fmt.Sprintf("customer_id = $%d", argCount))
		args = append(args, *filter.UserID)
		argCount++
	}

	if filter.Status != nil {
		conditions = append(conditions, fmt.Sprintf("status = $%d", argCount))
		args = append(args, *filter.Status)
		argCount++
	}

	if filter.PaymentStatus != nil {
		conditions = append(conditions, fmt.Sprintf("payment_status = $%d", argCount))
		args = append(args, *filter.PaymentStatus)
		argCount++
	}

	if filter.DateFrom != nil {
		conditions = append(conditions, fmt.Sprintf("created_at >= $%d", argCount))
		args = append(args, *filter.DateFrom)
		argCount++
	}

	if filter.DateTo != nil {
		conditions = append(conditions, fmt.Sprintf("created_at <= $%d", argCount))
		args = append(args, *filter.DateTo)
		argCount++
	}

	// Добавляем условия к запросам
	if len(conditions) > 0 {
		whereClause := " AND " + strings.Join(conditions, " AND ")
		countQuery += whereClause
		selectQuery += whereClause
	}

	// Получаем общее количество
	var total int
	err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count orders: %w", err)
	}

	// Добавляем сортировку и пагинацию
	selectQuery += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argCount, argCount+1)
	args = append(args, filter.Limit, filter.Offset)

	// Выполняем запрос
	rows, err := r.pool.Query(ctx, selectQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list orders: %w", err)
	}
	defer rows.Close()

	var orders []models.StorefrontOrder
	for rows.Next() {
		var order models.StorefrontOrder
		err := rows.Scan(
			&order.ID,
			&order.OrderNumber,
			&order.StorefrontID,
			&order.UserID,
			&order.Status,
			&order.Subtotal,
			&order.Tax,
			&order.Shipping,
			&order.Discount,
			&order.Total,
			&order.Currency,
			&order.CreatedAt,
			&order.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan order: %w", err)
		}
		orders = append(orders, order)
	}

	return orders, total, nil
}

// UpdateStatus обновляет статус заказа
func (r *orderRepository) UpdateStatus(ctx context.Context, orderID int64, status models.OrderStatus, metadata map[string]interface{}) error {
	metadataJSON, _ := json.Marshal(metadata)

	query := `UPDATE storefront_orders SET status = $1, metadata = metadata || $2::jsonb`
	args := []interface{}{status, metadataJSON}

	// Обновляем временные метки в зависимости от статуса
	switch status {
	case models.OrderStatusConfirmed:
		query += ", confirmed_at = CURRENT_TIMESTAMP"
	case models.OrderStatusShipped:
		query += ", shipped_at = CURRENT_TIMESTAMP"
	case models.OrderStatusDelivered:
		query += ", delivered_at = CURRENT_TIMESTAMP"
	case models.OrderStatusCancelled:
		query += ", cancelled_at = CURRENT_TIMESTAMP"
	}

	query += ", updated_at = CURRENT_TIMESTAMP WHERE id = $3"
	args = append(args, orderID)

	_, err := r.pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}

	return nil
}

// GetUserOrders получает заказы пользователя
func (r *orderRepository) GetUserOrders(ctx context.Context, userID int, limit, offset int) ([]models.StorefrontOrder, int, error) {
	filter := models.OrderFilter{
		UserID: &userID,
		Limit:  limit,
		Offset: offset,
	}
	return r.List(ctx, filter)
}

// GetStorefrontOrders получает заказы витрины
func (r *orderRepository) GetStorefrontOrders(ctx context.Context, storefrontID int, filter models.OrderFilter) ([]models.StorefrontOrder, int, error) {
	filter.StorefrontID = &storefrontID
	return r.List(ctx, filter)
}
