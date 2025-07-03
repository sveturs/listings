package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"

	"backend/internal/domain/models"
)

// InventoryRepositoryInterface определяет интерфейс для работы с инвентарем
type InventoryRepositoryInterface interface {
	CreateReservation(ctx context.Context, reservation *models.InventoryReservation) (*models.InventoryReservation, error)
	UpdateReservation(ctx context.Context, reservation *models.InventoryReservation) error
	GetReservation(ctx context.Context, reservationID int64) (*models.InventoryReservation, error)
	GetReservationsByOrder(ctx context.Context, orderID int64) ([]models.InventoryReservation, error)
	GetExpiredReservations(ctx context.Context) ([]models.InventoryReservation, error)
	GetReservedQuantity(ctx context.Context, productID int64, variantID *int64) (int, error)

	GetLowStockProducts(ctx context.Context, storefrontID int, threshold int) ([]models.StorefrontProduct, error)
	GetStockMovements(ctx context.Context, productID int64, limit int) ([]models.StorefrontInventoryMovement, error)
	CreateStockMovement(ctx context.Context, movement *models.StorefrontInventoryMovement) error
}

// InventoryRepository реализует интерфейс для работы с инвентарем
type InventoryRepository struct {
	db *sqlx.DB
}

// NewInventoryRepository создает новый репозиторий инвентаря
func NewInventoryRepository(db *sqlx.DB) *InventoryRepository {
	return &InventoryRepository{db: db}
}

// CreateReservation создает новое резервирование
func (r *InventoryRepository) CreateReservation(ctx context.Context, reservation *models.InventoryReservation) (*models.InventoryReservation, error) {
	query := `
		INSERT INTO inventory_reservations (
			product_id, variant_id, order_id, reserved_quantity, 
			status, expires_at
		) VALUES (
			:product_id, :variant_id, :order_id, :reserved_quantity,
			:status, :expires_at
		) RETURNING id, created_at`

	rows, err := r.db.NamedQueryContext(ctx, query, reservation)
	if err != nil {
		return nil, fmt.Errorf("failed to create reservation: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(&reservation.ID, &reservation.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan created reservation: %w", err)
		}
	}

	return reservation, nil
}

// UpdateReservation обновляет резервирование
func (r *InventoryRepository) UpdateReservation(ctx context.Context, reservation *models.InventoryReservation) error {
	query := `
		UPDATE inventory_reservations SET
			status = :status,
			released_at = :released_at
		WHERE id = :id`

	_, err := r.db.NamedExecContext(ctx, query, reservation)
	if err != nil {
		return fmt.Errorf("failed to update reservation: %w", err)
	}

	return nil
}

// GetReservation получает резервирование по ID
func (r *InventoryRepository) GetReservation(ctx context.Context, reservationID int64) (*models.InventoryReservation, error) {
	query := `
		SELECT id, product_id, variant_id, order_id, reserved_quantity,
			   status, expires_at, created_at, released_at
		FROM inventory_reservations 
		WHERE id = $1`

	var reservation models.InventoryReservation
	err := r.db.GetContext(ctx, &reservation, query, reservationID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("reservation not found")
		}
		return nil, fmt.Errorf("failed to get reservation: %w", err)
	}

	return &reservation, nil
}

// GetReservationsByOrder получает все резервирования для заказа
func (r *InventoryRepository) GetReservationsByOrder(ctx context.Context, orderID int64) ([]models.InventoryReservation, error) {
	query := `
		SELECT id, product_id, variant_id, order_id, reserved_quantity,
			   status, expires_at, created_at, released_at
		FROM inventory_reservations 
		WHERE order_id = $1
		ORDER BY created_at`

	var reservations []models.InventoryReservation
	err := r.db.SelectContext(ctx, &reservations, query, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get reservations by order: %w", err)
	}

	return reservations, nil
}

// GetExpiredReservations получает истекшие активные резервирования
func (r *InventoryRepository) GetExpiredReservations(ctx context.Context) ([]models.InventoryReservation, error) {
	query := `
		SELECT id, product_id, variant_id, order_id, reserved_quantity,
			   status, expires_at, created_at, released_at
		FROM inventory_reservations 
		WHERE status = 'active' AND expires_at < CURRENT_TIMESTAMP
		ORDER BY expires_at`

	var reservations []models.InventoryReservation
	err := r.db.SelectContext(ctx, &reservations, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get expired reservations: %w", err)
	}

	return reservations, nil
}

// GetReservedQuantity получает общее количество зарезервированного товара
func (r *InventoryRepository) GetReservedQuantity(ctx context.Context, productID int64, variantID *int64) (int, error) {
	query := `
		SELECT COALESCE(SUM(reserved_quantity), 0)
		FROM inventory_reservations 
		WHERE product_id = $1 
		AND COALESCE(variant_id, 0) = COALESCE($2, 0)
		AND status = 'active' 
		AND expires_at > CURRENT_TIMESTAMP`

	var reservedQuantity int
	err := r.db.GetContext(ctx, &reservedQuantity, query, productID, variantID)
	if err != nil {
		return 0, fmt.Errorf("failed to get reserved quantity: %w", err)
	}

	return reservedQuantity, nil
}

// GetLowStockProducts получает товары с низким остатком
func (r *InventoryRepository) GetLowStockProducts(ctx context.Context, storefrontID int, threshold int) ([]models.StorefrontProduct, error) {
	query := `
		SELECT p.id, p.storefront_id, p.name, p.description, p.price, p.currency,
			   p.category_id, p.sku, p.barcode, p.stock_quantity, p.stock_status,
			   p.is_active, p.attributes, p.view_count, p.sold_count,
			   p.external_id, p.created_at, p.updated_at
		FROM storefront_products p
		WHERE p.storefront_id = $1 
		AND p.is_active = true
		AND p.stock_quantity <= $2
		ORDER BY p.stock_quantity ASC, p.name`

	var products []models.StorefrontProduct
	err := r.db.SelectContext(ctx, &products, query, storefrontID, threshold)
	if err != nil {
		return nil, fmt.Errorf("failed to get low stock products: %w", err)
	}

	return products, nil
}

// GetStockMovements получает историю движений товара
func (r *InventoryRepository) GetStockMovements(ctx context.Context, productID int64, limit int) ([]models.StorefrontInventoryMovement, error) {
	query := `
		SELECT id, storefront_product_id, variant_id, type, quantity, reason,
			   order_id, notes, user_id, created_at
		FROM storefront_inventory_movements 
		WHERE storefront_product_id = $1
		ORDER BY created_at DESC
		LIMIT $2`

	var movements []models.StorefrontInventoryMovement
	err := r.db.SelectContext(ctx, &movements, query, productID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get stock movements: %w", err)
	}

	return movements, nil
}

// CreateStockMovement создает запись о движении товара
func (r *InventoryRepository) CreateStockMovement(ctx context.Context, movement *models.StorefrontInventoryMovement) error {
	query := `
		INSERT INTO storefront_inventory_movements (
			storefront_product_id, variant_id, type, quantity, reason,
			order_id, notes, user_id
		) VALUES (
			:storefront_product_id, :variant_id, :type, :quantity, :reason,
			:order_id, :notes, :user_id
		) RETURNING id, created_at`

	rows, err := r.db.NamedQueryContext(ctx, query, movement)
	if err != nil {
		return fmt.Errorf("failed to create stock movement: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(&movement.ID, &movement.CreatedAt)
		if err != nil {
			return fmt.Errorf("failed to scan created stock movement: %w", err)
		}
	}

	return nil
}
