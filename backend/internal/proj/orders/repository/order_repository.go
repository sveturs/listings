package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"

	"backend/internal/domain/models"
)

// OrderRepositoryInterface определяет интерфейс для работы с заказами
type OrderRepositoryInterface interface {
	Create(ctx context.Context, order *models.StorefrontOrder) (*models.StorefrontOrder, error)
	Update(ctx context.Context, order *models.StorefrontOrder) error
	Delete(ctx context.Context, orderID int64) error
	GetByID(ctx context.Context, orderID int64) (*models.StorefrontOrder, error)
	GetByIDWithDetails(ctx context.Context, orderID int64) (*models.StorefrontOrder, error)
	GetByFilter(ctx context.Context, filter *models.OrderFilter) ([]models.StorefrontOrder, int, error)
	AddItem(ctx context.Context, item *models.StorefrontOrderItem) error
	GetItems(ctx context.Context, orderID int64) ([]models.StorefrontOrderItem, error)
	GetReservations(ctx context.Context, orderID int64) ([]models.InventoryReservation, error)
	GetOrderSummaries(ctx context.Context, filter *models.OrderFilter) ([]models.OrderSummary, int, error)
}

// OrderRepository реализует интерфейс для работы с заказами
type OrderRepository struct {
	db *sqlx.DB
}

// NewOrderRepository создает новый репозиторий заказов
func NewOrderRepository(db *sqlx.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

// Create создает новый заказ
func (r *OrderRepository) Create(ctx context.Context, order *models.StorefrontOrder) (*models.StorefrontOrder, error) {
	query := `
		INSERT INTO storefront_orders (
			storefront_id, customer_id, subtotal_amount, shipping_amount, 
			tax_amount, total_amount, commission_amount, seller_amount, 
			currency, status, escrow_days, shipping_address, shipping_method, 
			customer_notes
		) VALUES (
			:storefront_id, :customer_id, :subtotal_amount, :shipping_amount,
			:tax_amount, :total_amount, :commission_amount, :seller_amount,
			:currency, :status, :escrow_days, :shipping_address, :shipping_method,
			:customer_notes
		) RETURNING id, order_number, created_at, updated_at, escrow_release_date`

	rows, err := r.db.NamedQueryContext(ctx, query, order)
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			// Логирование ошибки закрытия rows
		}
	}()

	if rows.Next() {
		err = rows.Scan(&order.ID, &order.OrderNumber, &order.CreatedAt, &order.UpdatedAt, &order.EscrowReleaseDate)
		if err != nil {
			return nil, fmt.Errorf("failed to scan created order: %w", err)
		}
	}

	return order, nil
}

// Update обновляет заказ
func (r *OrderRepository) Update(ctx context.Context, order *models.StorefrontOrder) error {
	query := `
		UPDATE storefront_orders SET
			payment_transaction_id = :payment_transaction_id,
			subtotal_amount = :subtotal_amount,
			shipping_amount = :shipping_amount,
			tax_amount = :tax_amount,
			total_amount = :total_amount,
			commission_amount = :commission_amount,
			seller_amount = :seller_amount,
			status = :status,
			escrow_release_date = :escrow_release_date,
			escrow_days = :escrow_days,
			shipping_address = :shipping_address,
			shipping_method = :shipping_method,
			shipping_provider = :shipping_provider,
			tracking_number = :tracking_number,
			customer_notes = :customer_notes,
			seller_notes = :seller_notes,
			confirmed_at = :confirmed_at,
			shipped_at = :shipped_at,
			delivered_at = :delivered_at,
			cancelled_at = :cancelled_at,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = :id`

	_, err := r.db.NamedExecContext(ctx, query, order)
	if err != nil {
		return fmt.Errorf("failed to update order: %w", err)
	}

	return nil
}

// Delete удаляет заказ
func (r *OrderRepository) Delete(ctx context.Context, orderID int64) error {
	query := `DELETE FROM storefront_orders WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, orderID)
	if err != nil {
		return fmt.Errorf("failed to delete order: %w", err)
	}
	return nil
}

// GetByID получает заказ по ID
func (r *OrderRepository) GetByID(ctx context.Context, orderID int64) (*models.StorefrontOrder, error) {
	query := `
		SELECT id, order_number, storefront_id, customer_id, payment_transaction_id,
			   subtotal_amount, shipping_amount, tax_amount, total_amount,
			   commission_amount, seller_amount, currency, status,
			   escrow_release_date, escrow_days, shipping_address, shipping_method,
			   shipping_provider, tracking_number, customer_notes, seller_notes,
			   confirmed_at, shipped_at, delivered_at, cancelled_at,
			   created_at, updated_at
		FROM storefront_orders 
		WHERE id = $1`

	var order models.StorefrontOrder
	err := r.db.GetContext(ctx, &order, query, orderID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("order not found")
		}
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	return &order, nil
}

// GetByIDWithDetails получает заказ со всеми связанными данными
func (r *OrderRepository) GetByIDWithDetails(ctx context.Context, orderID int64) (*models.StorefrontOrder, error) {
	order, err := r.GetByID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	// Получаем позиции заказа
	items, err := r.GetItems(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order items: %w", err)
	}
	order.Items = items

	// TODO: Получить связанные данные (Storefront, Customer, PaymentTransaction)
	// Это можно сделать через JOIN запросы или отдельные вызовы

	return order, nil
}

// GetByFilter получает заказы по фильтру
func (r *OrderRepository) GetByFilter(ctx context.Context, filter *models.OrderFilter) ([]models.StorefrontOrder, int, error) {
	// Строим WHERE условия
	whereConditions := []string{}
	args := map[string]interface{}{}

	if filter.StorefrontID != nil {
		whereConditions = append(whereConditions, "storefront_id = :storefront_id")
		args["storefront_id"] = *filter.StorefrontID
	}

	if filter.CustomerID != nil {
		whereConditions = append(whereConditions, "customer_id = :customer_id")
		args["customer_id"] = *filter.CustomerID
	}

	if filter.Status != nil {
		whereConditions = append(whereConditions, "status = :status")
		args["status"] = *filter.Status
	}

	if filter.DateFrom != nil {
		whereConditions = append(whereConditions, "created_at >= :date_from")
		args["date_from"] = *filter.DateFrom
	}

	if filter.DateTo != nil {
		whereConditions = append(whereConditions, "created_at <= :date_to")
		args["date_to"] = *filter.DateTo
	}

	if filter.MinAmount != nil {
		whereConditions = append(whereConditions, "total_amount >= :min_amount")
		args["min_amount"] = *filter.MinAmount
	}

	if filter.MaxAmount != nil {
		whereConditions = append(whereConditions, "total_amount <= :max_amount")
		args["max_amount"] = *filter.MaxAmount
	}

	if filter.OrderNumber != nil {
		whereConditions = append(whereConditions, "order_number ILIKE :order_number")
		args["order_number"] = "%" + *filter.OrderNumber + "%"
	}

	if filter.TrackingNumber != nil {
		whereConditions = append(whereConditions, "tracking_number ILIKE :tracking_number")
		args["tracking_number"] = "%" + *filter.TrackingNumber + "%"
	}

	whereClause := ""
	if len(whereConditions) > 0 {
		whereClause = "WHERE " + strings.Join(whereConditions, " AND ")
	}

	// Считаем общее количество
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM storefront_orders %s", whereClause)
	var total int
	countRows, err := r.db.NamedQueryContext(ctx, countQuery, args)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count orders: %w", err)
	}
	defer func() {
		if err := countRows.Close(); err != nil {
			// Логирование ошибки закрытия rows
		}
	}()
	if countRows.Next() {
		if err := countRows.Scan(&total); err != nil {
			return nil, 0, fmt.Errorf("failed to scan total count: %w", err)
		}
	}

	// Сортировка
	sortBy := "created_at"
	if filter.SortBy != "" {
		sortBy = filter.SortBy
	}

	sortOrder := "DESC"
	if filter.SortOrder == "asc" {
		sortOrder = "ASC"
	}

	// Пагинация
	limit := 20
	if filter.Limit > 0 && filter.Limit <= 100 {
		limit = filter.Limit
	}

	offset := 0
	if filter.Offset > 0 {
		offset = filter.Offset
	}

	args["limit"] = limit
	args["offset"] = offset

	// Основной запрос
	query := fmt.Sprintf(`
		SELECT id, order_number, storefront_id, customer_id, payment_transaction_id,
			   subtotal_amount, shipping_amount, tax_amount, total_amount,
			   commission_amount, seller_amount, currency, status,
			   escrow_release_date, escrow_days, shipping_address, shipping_method,
			   shipping_provider, tracking_number, customer_notes, seller_notes,
			   confirmed_at, shipped_at, delivered_at, cancelled_at,
			   created_at, updated_at
		FROM storefront_orders 
		%s
		ORDER BY %s %s
		LIMIT :limit OFFSET :offset`,
		whereClause, sortBy, sortOrder)

	rows, err := r.db.NamedQueryContext(ctx, query, args)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get orders: %w", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			// Логирование ошибки закрытия rows
		}
	}()

	var orders []models.StorefrontOrder
	for rows.Next() {
		var order models.StorefrontOrder
		err = rows.StructScan(&order)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan order: %w", err)
		}
		orders = append(orders, order)
	}

	return orders, total, nil
}

// AddItem добавляет позицию к заказу
func (r *OrderRepository) AddItem(ctx context.Context, item *models.StorefrontOrderItem) error {
	query := `
		INSERT INTO storefront_order_items (
			order_id, product_id, variant_id, product_name, product_sku,
			variant_name, quantity, price_per_unit, total_price, product_attributes
		) VALUES (
			:order_id, :product_id, :variant_id, :product_name, :product_sku,
			:variant_name, :quantity, :price_per_unit, :total_price, :product_attributes
		) RETURNING id, created_at`

	rows, err := r.db.NamedQueryContext(ctx, query, item)
	if err != nil {
		return fmt.Errorf("failed to add order item: %w", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			// Логирование ошибки закрытия rows
		}
	}()

	if rows.Next() {
		err = rows.Scan(&item.ID, &item.CreatedAt)
		if err != nil {
			return fmt.Errorf("failed to scan created order item: %w", err)
		}
	}

	return nil
}

// GetItems получает позиции заказа
func (r *OrderRepository) GetItems(ctx context.Context, orderID int64) ([]models.StorefrontOrderItem, error) {
	query := `
		SELECT id, order_id, product_id, variant_id, product_name, product_sku,
			   variant_name, quantity, price_per_unit, total_price, 
			   product_attributes, created_at
		FROM storefront_order_items 
		WHERE order_id = $1
		ORDER BY id`

	var items []models.StorefrontOrderItem
	err := r.db.SelectContext(ctx, &items, query, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order items: %w", err)
	}

	return items, nil
}

// GetReservations получает резервирования для заказа
func (r *OrderRepository) GetReservations(ctx context.Context, orderID int64) ([]models.InventoryReservation, error) {
	query := `
		SELECT id, product_id, variant_id, order_id, reserved_quantity,
			   status, expires_at, created_at, released_at
		FROM inventory_reservations 
		WHERE order_id = $1
		ORDER BY id`

	var reservations []models.InventoryReservation
	err := r.db.SelectContext(ctx, &reservations, query, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get reservations: %w", err)
	}

	return reservations, nil
}

// GetOrderSummaries получает краткую информацию о заказах
func (r *OrderRepository) GetOrderSummaries(ctx context.Context, filter *models.OrderFilter) ([]models.OrderSummary, int, error) {
	// Аналогично GetByFilter, но возвращаем только основную информацию
	// TODO: Реализовать с JOIN к пользователям и витринам для получения имен
	return nil, 0, fmt.Errorf("not implemented")
}
