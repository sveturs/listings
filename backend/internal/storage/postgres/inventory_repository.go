package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"backend/internal/domain/models"
)

// InventoryRepositoryInterface определяет интерфейс для работы с инвентарем
type InventoryRepositoryInterface interface {
	GetStock(ctx context.Context, productID int64, variantID *int64) (*models.InventoryStock, error)
	UpdateStock(ctx context.Context, productID int64, variantID *int64, quantity int) error
	ReserveStock(ctx context.Context, reservation *models.InventoryReservation) error
	ReleaseReservation(ctx context.Context, reservationID int64) error
	ConfirmReservation(ctx context.Context, reservationID int64) error
	GetExpiredReservations(ctx context.Context) ([]models.InventoryReservation, error)
	GetReservationsByOrder(ctx context.Context, orderID int64) ([]models.InventoryReservation, error)
	RecordMovement(ctx context.Context, movement *models.InventoryMovement) error
	GetLowStockItems(ctx context.Context, storefrontID int) ([]models.LowStockItem, error)
}

// inventoryRepository реализует интерфейс для работы с инвентарем
type inventoryRepository struct {
	pool *pgxpool.Pool
}

// NewInventoryRepository создает новый репозиторий инвентаря
func NewInventoryRepository(pool *pgxpool.Pool) InventoryRepositoryInterface {
	return &inventoryRepository{pool: pool}
}

// GetStock получает информацию о запасах товара
func (r *inventoryRepository) GetStock(ctx context.Context, productID int64, variantID *int64) (*models.InventoryStock, error) {
	query := `
		SELECT product_id, variant_id, quantity, reserved_quantity, 
			   available_quantity, low_stock_threshold, updated_at
		FROM inventory_stock
		WHERE product_id = $1 AND COALESCE(variant_id, 0) = COALESCE($2, 0)`

	var stock models.InventoryStock
	var dbVariantID sql.NullInt64

	err := r.pool.QueryRow(ctx, query, productID, variantID).Scan(
		&stock.ProductID,
		&dbVariantID,
		&stock.Quantity,
		&stock.ReservedQuantity,
		&stock.AvailableQuantity,
		&stock.LowStockThreshold,
		&stock.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			// Если записи нет, возвращаем пустой stock
			return &models.InventoryStock{
				ProductID:         productID,
				VariantID:         variantID,
				Quantity:          0,
				ReservedQuantity:  0,
				AvailableQuantity: 0,
				LowStockThreshold: 0,
			}, nil
		}
		return nil, fmt.Errorf("failed to get stock: %w", err)
	}

	if dbVariantID.Valid {
		stock.VariantID = &dbVariantID.Int64
	}

	return &stock, nil
}

// UpdateStock обновляет количество товара на складе
func (r *inventoryRepository) UpdateStock(ctx context.Context, productID int64, variantID *int64, quantity int) error {
	// Используем UPSERT для создания записи если её нет
	query := `
		INSERT INTO inventory_stock (product_id, variant_id, quantity, reserved_quantity, available_quantity)
		VALUES ($1, $2, $3, 0, $3)
		ON CONFLICT (product_id, COALESCE(variant_id, 0))
		DO UPDATE SET 
			quantity = $3,
			available_quantity = $3 - inventory_stock.reserved_quantity,
			updated_at = CURRENT_TIMESTAMP`

	_, err := r.pool.Exec(ctx, query, productID, variantID, quantity)
	if err != nil {
		return fmt.Errorf("failed to update stock: %w", err)
	}

	return nil
}

// ReserveStock резервирует товар для заказа
func (r *inventoryRepository) ReserveStock(ctx context.Context, reservation *models.InventoryReservation) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Проверяем доступность товара и обновляем резерв
	updateQuery := `
		UPDATE inventory_stock 
		SET reserved_quantity = reserved_quantity + $1,
			available_quantity = quantity - (reserved_quantity + $1),
			updated_at = CURRENT_TIMESTAMP
		WHERE product_id = $2 
			AND COALESCE(variant_id, 0) = COALESCE($3, 0)
			AND available_quantity >= $1
		RETURNING available_quantity`

	var newAvailable int
	err = tx.QueryRow(ctx, updateQuery,
		reservation.Quantity,
		reservation.ProductID,
		reservation.VariantID,
	).Scan(&newAvailable)
	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("insufficient stock")
		}
		return fmt.Errorf("failed to reserve stock: %w", err)
	}

	// Создаем запись о резервировании
	insertQuery := `
		INSERT INTO inventory_reservations (
			order_id, product_id, variant_id, quantity, 
			status, expires_at
		) VALUES (
			$1, $2, $3, $4, $5, $6
		) RETURNING id, created_at`

	err = tx.QueryRow(ctx, insertQuery,
		reservation.OrderID,
		reservation.ProductID,
		reservation.VariantID,
		reservation.Quantity,
		reservation.Status,
		reservation.ExpiresAt,
	).Scan(&reservation.ID, &reservation.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create reservation: %w", err)
	}

	// Записываем движение товара
	movementQuery := `
		INSERT INTO inventory_movements (
			product_id, variant_id, type, quantity, 
			reference_type, reference_id, notes
		) VALUES (
			$1, $2, 'reservation', $3, 'order', $4, 'Stock reserved for order'
		)`

	_, err = tx.Exec(ctx, movementQuery,
		reservation.ProductID,
		reservation.VariantID,
		-reservation.Quantity, // Отрицательное для резервирования
		reservation.OrderID,
	)
	if err != nil {
		return fmt.Errorf("failed to record movement: %w", err)
	}

	return tx.Commit(ctx)
}

// ReleaseReservation освобождает зарезервированный товар
func (r *inventoryRepository) ReleaseReservation(ctx context.Context, reservationID int64) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Получаем информацию о резервировании
	var reservation models.InventoryReservation
	var variantID sql.NullInt64

	getQuery := `
		SELECT product_id, variant_id, quantity, order_id, status
		FROM inventory_reservations
		WHERE id = $1 AND status = 'reserved'`

	err = tx.QueryRow(ctx, getQuery, reservationID).Scan(
		&reservation.ProductID,
		&variantID,
		&reservation.Quantity,
		&reservation.OrderID,
		&reservation.Status,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("reservation not found or already processed")
		}
		return fmt.Errorf("failed to get reservation: %w", err)
	}

	if variantID.Valid {
		reservation.VariantID = &variantID.Int64
	}

	// Обновляем статус резервирования
	updateReservationQuery := `
		UPDATE inventory_reservations 
		SET status = 'released', updated_at = CURRENT_TIMESTAMP
		WHERE id = $1`

	_, err = tx.Exec(ctx, updateReservationQuery, reservationID)
	if err != nil {
		return fmt.Errorf("failed to update reservation: %w", err)
	}

	// Возвращаем товар в доступные
	updateStockQuery := `
		UPDATE inventory_stock 
		SET reserved_quantity = reserved_quantity - $1,
			available_quantity = available_quantity + $1,
			updated_at = CURRENT_TIMESTAMP
		WHERE product_id = $2 AND COALESCE(variant_id, 0) = COALESCE($3, 0)`

	_, err = tx.Exec(ctx, updateStockQuery,
		reservation.Quantity,
		reservation.ProductID,
		reservation.VariantID,
	)
	if err != nil {
		return fmt.Errorf("failed to update stock: %w", err)
	}

	// Записываем движение товара
	movementQuery := `
		INSERT INTO inventory_movements (
			product_id, variant_id, type, quantity,
			reference_type, reference_id, notes
		) VALUES (
			$1, $2, 'release', $3, 'order', $4, 'Reservation released'
		)`

	_, err = tx.Exec(ctx, movementQuery,
		reservation.ProductID,
		reservation.VariantID,
		reservation.Quantity, // Положительное для возврата
		reservation.OrderID,
	)
	if err != nil {
		return fmt.Errorf("failed to record movement: %w", err)
	}

	return tx.Commit(ctx)
}

// ConfirmReservation подтверждает резервирование при оплате заказа
func (r *inventoryRepository) ConfirmReservation(ctx context.Context, reservationID int64) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Получаем информацию о резервировании
	var reservation models.InventoryReservation
	var variantID sql.NullInt64

	getQuery := `
		SELECT product_id, variant_id, quantity, order_id
		FROM inventory_reservations
		WHERE id = $1 AND status = 'reserved'`

	err = tx.QueryRow(ctx, getQuery, reservationID).Scan(
		&reservation.ProductID,
		&variantID,
		&reservation.Quantity,
		&reservation.OrderID,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("reservation not found or already processed")
		}
		return fmt.Errorf("failed to get reservation: %w", err)
	}

	if variantID.Valid {
		reservation.VariantID = &variantID.Int64
	}

	// Обновляем статус резервирования
	updateReservationQuery := `
		UPDATE inventory_reservations 
		SET status = 'confirmed', updated_at = CURRENT_TIMESTAMP
		WHERE id = $1`

	_, err = tx.Exec(ctx, updateReservationQuery, reservationID)
	if err != nil {
		return fmt.Errorf("failed to update reservation: %w", err)
	}

	// Списываем товар со склада
	updateStockQuery := `
		UPDATE inventory_stock 
		SET quantity = quantity - $1,
			reserved_quantity = reserved_quantity - $1,
			updated_at = CURRENT_TIMESTAMP
		WHERE product_id = $2 AND COALESCE(variant_id, 0) = COALESCE($3, 0)`

	_, err = tx.Exec(ctx, updateStockQuery,
		reservation.Quantity,
		reservation.ProductID,
		reservation.VariantID,
	)
	if err != nil {
		return fmt.Errorf("failed to update stock: %w", err)
	}

	// Записываем движение товара
	movementQuery := `
		INSERT INTO inventory_movements (
			product_id, variant_id, type, quantity,
			reference_type, reference_id, notes
		) VALUES (
			$1, $2, 'sale', $3, 'order', $4, 'Stock sold'
		)`

	_, err = tx.Exec(ctx, movementQuery,
		reservation.ProductID,
		reservation.VariantID,
		-reservation.Quantity, // Отрицательное для продажи
		reservation.OrderID,
	)
	if err != nil {
		return fmt.Errorf("failed to record movement: %w", err)
	}

	return tx.Commit(ctx)
}

// GetExpiredReservations получает истекшие резервирования
func (r *inventoryRepository) GetExpiredReservations(ctx context.Context) ([]models.InventoryReservation, error) {
	query := `
		SELECT id, order_id, product_id, variant_id, quantity,
			   status, expires_at, created_at, updated_at
		FROM inventory_reservations
		WHERE status = 'reserved' AND expires_at < $1`

	rows, err := r.pool.Query(ctx, query, time.Now())
	if err != nil {
		return nil, fmt.Errorf("failed to get expired reservations: %w", err)
	}
	defer rows.Close()

	var reservations []models.InventoryReservation
	for rows.Next() {
		var res models.InventoryReservation
		var variantID sql.NullInt64

		err := rows.Scan(
			&res.ID,
			&res.OrderID,
			&res.ProductID,
			&variantID,
			&res.Quantity,
			&res.Status,
			&res.ExpiresAt,
			&res.CreatedAt,
			&res.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan reservation: %w", err)
		}

		if variantID.Valid {
			res.VariantID = &variantID.Int64
		}

		reservations = append(reservations, res)
	}

	return reservations, nil
}

// GetReservationsByOrder получает все резервирования для заказа
func (r *inventoryRepository) GetReservationsByOrder(ctx context.Context, orderID int64) ([]models.InventoryReservation, error) {
	query := `
		SELECT id, order_id, product_id, variant_id, quantity,
			   status, expires_at, created_at, updated_at
		FROM inventory_reservations
		WHERE order_id = $1`

	rows, err := r.pool.Query(ctx, query, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get reservations: %w", err)
	}
	defer rows.Close()

	var reservations []models.InventoryReservation
	for rows.Next() {
		var res models.InventoryReservation
		var variantID sql.NullInt64

		err := rows.Scan(
			&res.ID,
			&res.OrderID,
			&res.ProductID,
			&variantID,
			&res.Quantity,
			&res.Status,
			&res.ExpiresAt,
			&res.CreatedAt,
			&res.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan reservation: %w", err)
		}

		if variantID.Valid {
			res.VariantID = &variantID.Int64
		}

		reservations = append(reservations, res)
	}

	return reservations, nil
}

// RecordMovement записывает движение товара
func (r *inventoryRepository) RecordMovement(ctx context.Context, movement *models.InventoryMovement) error {
	metadataJSON, _ := json.Marshal(movement.Metadata)

	query := `
		INSERT INTO inventory_movements (
			product_id, variant_id, type, quantity,
			reference_type, reference_id, notes, metadata
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8::jsonb
		) RETURNING id, created_at`

	err := r.pool.QueryRow(ctx, query,
		movement.ProductID,
		movement.VariantID,
		movement.Type,
		movement.Quantity,
		movement.ReferenceType,
		movement.ReferenceID,
		movement.Notes,
		metadataJSON,
	).Scan(&movement.ID, &movement.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to record movement: %w", err)
	}

	return nil
}

// GetLowStockItems получает товары с низким запасом
func (r *inventoryRepository) GetLowStockItems(ctx context.Context, storefrontID int) ([]models.LowStockItem, error) {
	query := `
		SELECT 
			p.id as product_id,
			p.name as product_name,
			pv.id as variant_id,
			pv.name as variant_name,
			s.quantity,
			s.available_quantity,
			s.low_stock_threshold
		FROM inventory_stock s
		JOIN storefront_products p ON p.id = s.product_id
		LEFT JOIN storefront_product_variants pv ON pv.id = s.variant_id
		WHERE p.storefront_id = $1
			AND s.available_quantity <= s.low_stock_threshold
			AND s.low_stock_threshold > 0
		ORDER BY s.available_quantity ASC`

	rows, err := r.pool.Query(ctx, query, storefrontID)
	if err != nil {
		return nil, fmt.Errorf("failed to get low stock items: %w", err)
	}
	defer rows.Close()

	var items []models.LowStockItem
	for rows.Next() {
		var item models.LowStockItem
		var variantID sql.NullInt64
		var variantName sql.NullString

		err := rows.Scan(
			&item.ProductID,
			&item.ProductName,
			&variantID,
			&variantName,
			&item.Quantity,
			&item.AvailableQuantity,
			&item.LowStockThreshold,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan low stock item: %w", err)
		}

		if variantID.Valid {
			item.VariantID = &variantID.Int64
		}
		if variantName.Valid {
			item.VariantName = &variantName.String
		}

		items = append(items, item)
	}

	return items, nil
}
