package listings

import (
	"context"
	"database/sql"
	"fmt"
)

// StockItem represents a product or variant with quantity to modify
type StockItem struct {
	ProductID int64
	VariantID *int64 // nil for products, non-nil for variants
	Quantity  int32
}

// StockResult represents the result of a stock operation
type StockResult struct {
	ProductID   int64
	VariantID   *int64
	StockBefore int32
	StockAfter  int32
	Success     bool
	Error       *string
}

// DecrementStock atomically decrements stock for multiple products/variants
// Used during order creation. Returns error if any item has insufficient stock.
func (s *Service) DecrementStock(ctx context.Context, items []StockItem, orderID *string) ([]StockResult, error) {
	logger := s.logger.With().
		Str("method", "DecrementStock").
		Int("items_count", len(items)).
		Logger()

	if orderID != nil {
		logger = logger.With().Str("order_id", *orderID).Logger()
	}

	logger.Debug().Msg("Starting stock decrement")

	// Start transaction
	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to begin transaction")
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	results := make([]StockResult, 0, len(items))

	// Process each item
	for _, item := range items {
		result, err := s.decrementStockItem(ctx, tx, item)
		if err != nil {
			logger.Error().
				Err(err).
				Int64("product_id", item.ProductID).
				Msg("Failed to decrement stock")

			// Rollback transaction
			tx.Rollback()

			errMsg := err.Error()
			result.Success = false
			result.Error = &errMsg
			results = append(results, result)

			return results, fmt.Errorf("failed to decrement stock for product %d: %w", item.ProductID, err)
		}

		results = append(results, result)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		logger.Error().Err(err).Msg("Failed to commit transaction")
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	logger.Info().
		Int("items_processed", len(results)).
		Msg("Stock decremented successfully")

	return results, nil
}

// decrementStockItem decrements stock for a single item
func (s *Service) decrementStockItem(ctx context.Context, tx *sql.Tx, item StockItem) (StockResult, error) {
	if item.VariantID != nil {
		// Decrement variant stock
		return s.decrementVariantStock(ctx, tx, item.ProductID, *item.VariantID, item.Quantity)
	}

	// Decrement product stock
	return s.decrementProductStock(ctx, tx, item.ProductID, item.Quantity)
}

// decrementProductStock decrements stock for a product
func (s *Service) decrementProductStock(ctx context.Context, tx *sql.Tx, productID int64, quantity int32) (StockResult, error) {
	result := StockResult{
		ProductID: productID,
		VariantID: nil,
	}

	// Lock row and get current stock
	var currentStock int32
	query := `
		SELECT stock_quantity
		FROM b2c_products
		WHERE id = $1
		FOR UPDATE
	`
	err := tx.QueryRowContext(ctx, query, productID).Scan(&currentStock)
	if err != nil {
		if err == sql.ErrNoRows {
			errMsg := "product not found"
			result.Error = &errMsg
			return result, fmt.Errorf("product %d not found", productID)
		}
		return result, fmt.Errorf("failed to get product stock: %w", err)
	}

	result.StockBefore = currentStock

	// Check if sufficient stock
	if currentStock < quantity {
		errMsg := fmt.Sprintf("insufficient stock: have %d, need %d", currentStock, quantity)
		result.Error = &errMsg
		return result, fmt.Errorf("insufficient stock for product %d: have %d, need %d", productID, currentStock, quantity)
	}

	// Decrement stock
	updateQuery := `
		UPDATE b2c_products
		SET stock_quantity = stock_quantity - $1,
		    updated_at = NOW()
		WHERE id = $2 AND stock_quantity >= $1
	`
	execResult, err := tx.ExecContext(ctx, updateQuery, quantity, productID)
	if err != nil {
		return result, fmt.Errorf("failed to decrement stock: %w", err)
	}

	rowsAffected, err := execResult.RowsAffected()
	if err != nil {
		return result, fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		errMsg := "stock changed during transaction"
		result.Error = &errMsg
		return result, fmt.Errorf("stock changed during transaction for product %d", productID)
	}

	result.StockAfter = currentStock - quantity
	result.Success = true

	return result, nil
}

// decrementVariantStock decrements stock for a variant
func (s *Service) decrementVariantStock(ctx context.Context, tx *sql.Tx, productID, variantID int64, quantity int32) (StockResult, error) {
	result := StockResult{
		ProductID: productID,
		VariantID: &variantID,
	}

	// Lock row and get current stock
	var currentStock int32
	query := `
		SELECT stock_quantity
		FROM b2c_product_variants
		WHERE id = $1 AND product_id = $2
		FOR UPDATE
	`
	err := tx.QueryRowContext(ctx, query, variantID, productID).Scan(&currentStock)
	if err != nil {
		if err == sql.ErrNoRows {
			errMsg := "variant not found"
			result.Error = &errMsg
			return result, fmt.Errorf("variant %d for product %d not found", variantID, productID)
		}
		return result, fmt.Errorf("failed to get variant stock: %w", err)
	}

	result.StockBefore = currentStock

	// Check if sufficient stock
	if currentStock < quantity {
		errMsg := fmt.Sprintf("insufficient stock: have %d, need %d", currentStock, quantity)
		result.Error = &errMsg
		return result, fmt.Errorf("insufficient stock for variant %d: have %d, need %d", variantID, currentStock, quantity)
	}

	// Decrement stock
	updateQuery := `
		UPDATE b2c_product_variants
		SET stock_quantity = stock_quantity - $1,
		    updated_at = NOW()
		WHERE id = $2 AND product_id = $3 AND stock_quantity >= $1
	`
	execResult, err := tx.ExecContext(ctx, updateQuery, quantity, variantID, productID)
	if err != nil {
		return result, fmt.Errorf("failed to decrement variant stock: %w", err)
	}

	rowsAffected, err := execResult.RowsAffected()
	if err != nil {
		return result, fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		errMsg := "stock changed during transaction"
		result.Error = &errMsg
		return result, fmt.Errorf("stock changed during transaction for variant %d", variantID)
	}

	result.StockAfter = currentStock - quantity
	result.Success = true

	return result, nil
}

// RollbackStock restores stock (compensating transaction for failed orders)
// IDEMPOTENCY PROTECTION: Multiple calls with same orderID will NOT duplicate rollback
func (s *Service) RollbackStock(ctx context.Context, items []StockItem, orderID *string) ([]StockResult, error) {
	logger := s.logger.With().
		Str("method", "RollbackStock").
		Int("items_count", len(items)).
		Logger()

	if orderID != nil {
		logger = logger.With().Str("order_id", *orderID).Logger()
	}

	logger.Debug().Msg("Starting stock rollback")

	// Validate: orderID is REQUIRED for idempotency protection
	if orderID == nil || *orderID == "" {
		logger.Error().Msg("order_id is required for rollback")
		return nil, fmt.Errorf("order_id is required for idempotency protection")
	}

	// Start transaction
	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to begin transaction")
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	results := make([]StockResult, 0, len(items))

	// Process each item - increment stock back
	for _, item := range items {
		result, err := s.rollbackStockItem(ctx, tx, item, *orderID)
		if err != nil {
			logger.Error().
				Err(err).
				Int64("product_id", item.ProductID).
				Msg("Failed to rollback stock")

			// Continue with other items even if one fails
			errMsg := err.Error()
			result.Success = false
			result.Error = &errMsg
		}

		results = append(results, result)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		logger.Error().Err(err).Msg("Failed to commit transaction")
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	logger.Info().
		Int("items_processed", len(results)).
		Msg("Stock rollback completed")

	return results, nil
}

// rollbackStockItem increments stock for a single item (reverse of decrement)
func (s *Service) rollbackStockItem(ctx context.Context, tx *sql.Tx, item StockItem, orderID string) (StockResult, error) {
	if item.VariantID != nil {
		// Rollback variant stock
		return s.rollbackVariantStock(ctx, tx, item.ProductID, *item.VariantID, item.Quantity, orderID)
	}

	// Rollback product stock
	return s.rollbackProductStock(ctx, tx, item.ProductID, item.Quantity, orderID)
}

// rollbackProductStock increments product stock with idempotency protection
func (s *Service) rollbackProductStock(ctx context.Context, tx *sql.Tx, productID int64, quantity int32, orderID string) (StockResult, error) {
	result := StockResult{
		ProductID: productID,
		VariantID: nil,
	}

	// IDEMPOTENCY CHECK: Check if rollback already performed for this order + product
	exists, err := s.checkRollbackExists(ctx, tx, orderID, productID, nil)
	if err != nil {
		return result, fmt.Errorf("failed to check rollback existence: %w", err)
	}

	if exists {
		// Rollback already done - return success with current stock (idempotent)
		var currentStock int32
		query := `SELECT stock_quantity FROM b2c_products WHERE id = $1`
		if err := tx.QueryRowContext(ctx, query, productID).Scan(&currentStock); err != nil {
			return result, fmt.Errorf("failed to get product stock: %w", err)
		}

		s.logger.Info().
			Str("order_id", orderID).
			Int64("product_id", productID).
			Msg("Rollback already performed (idempotent)")

		result.StockBefore = currentStock
		result.StockAfter = currentStock // No change
		result.Success = true
		return result, nil
	}

	// Get current stock (no need for lock on rollback)
	var currentStock int32
	query := `SELECT stock_quantity FROM b2c_products WHERE id = $1`
	err = tx.QueryRowContext(ctx, query, productID).Scan(&currentStock)
	if err != nil {
		if err == sql.ErrNoRows {
			errMsg := "product not found"
			result.Error = &errMsg
			return result, fmt.Errorf("product %d not found during rollback", productID)
		}
		return result, fmt.Errorf("failed to get product stock: %w", err)
	}

	result.StockBefore = currentStock

	// Increment stock back
	updateQuery := `
		UPDATE b2c_products
		SET stock_quantity = stock_quantity + $1,
		    updated_at = NOW()
		WHERE id = $2
	`
	_, err = tx.ExecContext(ctx, updateQuery, quantity, productID)
	if err != nil {
		return result, fmt.Errorf("failed to rollback product stock: %w", err)
	}

	// Record rollback in audit table for idempotency
	if err := s.recordRollback(ctx, tx, orderID, productID, nil, quantity); err != nil {
		return result, fmt.Errorf("failed to record rollback: %w", err)
	}

	result.StockAfter = currentStock + quantity
	result.Success = true

	return result, nil
}

// rollbackVariantStock increments variant stock with idempotency protection
func (s *Service) rollbackVariantStock(ctx context.Context, tx *sql.Tx, productID, variantID int64, quantity int32, orderID string) (StockResult, error) {
	result := StockResult{
		ProductID: productID,
		VariantID: &variantID,
	}

	// IDEMPOTENCY CHECK: Check if rollback already performed for this order + variant
	exists, err := s.checkRollbackExists(ctx, tx, orderID, productID, &variantID)
	if err != nil {
		return result, fmt.Errorf("failed to check rollback existence: %w", err)
	}

	if exists {
		// Rollback already done - return success with current stock (idempotent)
		var currentStock int32
		query := `SELECT stock_quantity FROM b2c_product_variants WHERE id = $1 AND product_id = $2`
		if err := tx.QueryRowContext(ctx, query, variantID, productID).Scan(&currentStock); err != nil {
			return result, fmt.Errorf("failed to get variant stock: %w", err)
		}

		s.logger.Info().
			Str("order_id", orderID).
			Int64("product_id", productID).
			Int64("variant_id", variantID).
			Msg("Rollback already performed (idempotent)")

		result.StockBefore = currentStock
		result.StockAfter = currentStock // No change
		result.Success = true
		return result, nil
	}

	// Get current stock (no need for lock on rollback)
	var currentStock int32
	query := `SELECT stock_quantity FROM b2c_product_variants WHERE id = $1 AND product_id = $2`
	err = tx.QueryRowContext(ctx, query, variantID, productID).Scan(&currentStock)
	if err != nil {
		if err == sql.ErrNoRows {
			errMsg := "variant not found"
			result.Error = &errMsg
			return result, fmt.Errorf("variant %d for product %d not found during rollback", variantID, productID)
		}
		return result, fmt.Errorf("failed to get variant stock: %w", err)
	}

	result.StockBefore = currentStock

	// Increment stock back
	updateQuery := `
		UPDATE b2c_product_variants
		SET stock_quantity = stock_quantity + $1,
		    updated_at = NOW()
		WHERE id = $2 AND product_id = $3
	`
	_, err = tx.ExecContext(ctx, updateQuery, quantity, variantID, productID)
	if err != nil {
		return result, fmt.Errorf("failed to rollback variant stock: %w", err)
	}

	// Record rollback in audit table for idempotency
	if err := s.recordRollback(ctx, tx, orderID, productID, &variantID, quantity); err != nil {
		return result, fmt.Errorf("failed to record rollback: %w", err)
	}

	result.StockAfter = currentStock + quantity
	result.Success = true

	return result, nil
}

// CheckStockAvailability checks if requested quantities are available
func (s *Service) CheckStockAvailability(ctx context.Context, items []StockItem) (bool, []StockAvailability, error) {
	logger := s.logger.With().
		Str("method", "CheckStockAvailability").
		Int("items_count", len(items)).
		Logger()

	logger.Debug().Msg("Checking stock availability")

	availabilities := make([]StockAvailability, 0, len(items))
	allAvailable := true

	for _, item := range items {
		availability, err := s.checkItemAvailability(ctx, item)
		if err != nil {
			logger.Error().
				Err(err).
				Int64("product_id", item.ProductID).
				Msg("Failed to check item availability")
			return false, nil, err
		}

		availabilities = append(availabilities, availability)

		if !availability.IsAvailable {
			allAvailable = false
		}
	}

	logger.Debug().
		Bool("all_available", allAvailable).
		Msg("Stock availability check completed")

	return allAvailable, availabilities, nil
}

// StockAvailability represents availability status for a single item
type StockAvailability struct {
	ProductID         int64
	VariantID         *int64
	RequestedQuantity int32
	AvailableQuantity int32
	IsAvailable       bool
}

// checkItemAvailability checks availability for a single item
func (s *Service) checkItemAvailability(ctx context.Context, item StockItem) (StockAvailability, error) {
	availability := StockAvailability{
		ProductID:         item.ProductID,
		VariantID:         item.VariantID,
		RequestedQuantity: item.Quantity,
	}

	var currentStock int32
	var err error

	if item.VariantID != nil {
		// Check variant stock
		query := `SELECT stock_quantity FROM b2c_product_variants WHERE id = $1 AND product_id = $2`
		err = s.repo.GetDB().QueryRowContext(ctx, query, *item.VariantID, item.ProductID).Scan(&currentStock)
	} else {
		// Check product stock
		query := `SELECT stock_quantity FROM b2c_products WHERE id = $1`
		err = s.repo.GetDB().QueryRowContext(ctx, query, item.ProductID).Scan(&currentStock)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			availability.AvailableQuantity = 0
			availability.IsAvailable = false
			return availability, nil
		}
		return availability, fmt.Errorf("failed to get stock: %w", err)
	}

	availability.AvailableQuantity = currentStock
	availability.IsAvailable = currentStock >= item.Quantity

	return availability, nil
}

// ============================================================================
// IDEMPOTENCY HELPER METHODS
// ============================================================================

// checkRollbackExists checks if rollback already performed for given order + product/variant
// Returns true if rollback record exists (idempotency check)
func (s *Service) checkRollbackExists(ctx context.Context, tx *sql.Tx, orderID string, productID int64, variantID *int64) (bool, error) {
	var count int
	var query string
	var args []interface{}

	if variantID != nil {
		// Check variant rollback
		query = `
			SELECT COUNT(*)
			FROM b2c_inventory_movements
			WHERE order_id = $1
			  AND variant_id = $2
			  AND movement_type = 'rollback'
		`
		args = []interface{}{orderID, *variantID}
	} else {
		// Check product rollback
		query = `
			SELECT COUNT(*)
			FROM b2c_inventory_movements
			WHERE order_id = $1
			  AND storefront_product_id = $2
			  AND variant_id IS NULL
			  AND movement_type = 'rollback'
		`
		args = []interface{}{orderID, productID}
	}

	err := tx.QueryRowContext(ctx, query, args...).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check rollback existence: %w", err)
	}

	return count > 0, nil
}

// recordRollback records rollback operation in audit table for idempotency tracking
// This ensures duplicate rollbacks can be detected and prevented
func (s *Service) recordRollback(ctx context.Context, tx *sql.Tx, orderID string, productID int64, variantID *int64, quantity int32) error {
	// Get user_id from context (default to 0 if not present)
	// TODO: Extract user_id from context when auth is integrated
	userID := int64(0)

	query := `
		INSERT INTO b2c_inventory_movements (
			storefront_product_id,
			variant_id,
			type,
			quantity,
			reason,
			notes,
			user_id,
			order_id,
			movement_type,
			created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW())
	`

	notes := fmt.Sprintf("Stock rollback for order %s", orderID)

	_, err := tx.ExecContext(ctx, query,
		productID,
		variantID,
		"in",             // type: stock is coming back IN
		quantity,         // quantity restored
		"rollback",       // reason
		notes,            // notes
		userID,           // user_id (system operation)
		orderID,          // order_id for idempotency
		"rollback",       // movement_type for filtering
	)

	if err != nil {
		// Check if it's a unique constraint violation (duplicate rollback attempt)
		// This can happen in concurrent scenarios - the UNIQUE index will prevent corruption
		if err.Error() == "pq: duplicate key value violates unique constraint" {
			return fmt.Errorf("rollback already recorded for order %s (concurrent protection)", orderID)
		}
		return fmt.Errorf("failed to record rollback: %w", err)
	}

	s.logger.Debug().
		Str("order_id", orderID).
		Int64("product_id", productID).
		Int32("quantity", quantity).
		Msg("Rollback recorded in audit table")

	return nil
}
