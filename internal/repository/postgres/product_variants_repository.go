package postgres

// ⚠️ DEPRECATED: Product variants functionality is deprecated.
// The b2c_product_variants table was removed in Phase 11.5 (migration 000010).
// This file is kept for API compatibility but methods will return errors.
// Consider creating a new unified product_variants table if variants are needed.

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/lib/pq"
	"github.com/vondi-global/listings/internal/domain"
)

// CreateProductVariant creates a new product variant
func (r *Repository) CreateProductVariant(ctx context.Context, input *domain.CreateVariantInput) (*domain.ProductVariant, error) {
	r.logger.Debug().
		Int64("product_id", input.ProductID).
		Interface("variant_attributes", input.VariantAttributes).
		Msg("creating product variant")

	// Validate input
	if input.ProductID <= 0 {
		return nil, fmt.Errorf("variants.invalid_product_id")
	}

	if input.StockQuantity < 0 {
		return nil, fmt.Errorf("variants.invalid_stock_quantity")
	}

	// Start transaction
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		r.logger.Error().Err(err).Msg("failed to begin transaction")
		return nil, fmt.Errorf("variants.create_failed")
	}
	defer tx.Rollback()

	// Check if product exists and get has_variants flag
	var hasVariants bool
	err = tx.QueryRowContext(ctx, `
		SELECT has_variants FROM listings
		WHERE id = $1 AND status = 'active' AND deleted_at IS NULL AND source_type = 'b2c'
	`, input.ProductID).Scan(&hasVariants)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("variants.product_not_found")
		}
		r.logger.Error().Err(err).Msg("failed to check product")
		return nil, fmt.Errorf("variants.create_failed")
	}

	// Business rule: Product must have has_variants=true
	if !hasVariants {
		return nil, fmt.Errorf("variants.product_no_variants")
	}

	// If this is a default variant, unset other defaults
	if input.IsDefault {
		_, err = tx.ExecContext(ctx, `
			UPDATE b2c_product_variants
			SET is_default = false, updated_at = NOW()
			WHERE product_id = $1 AND is_default = true
		`, input.ProductID)

		if err != nil {
			r.logger.Error().Err(err).Msg("failed to unset other defaults")
			return nil, fmt.Errorf("variants.create_failed")
		}
	}

	// Marshal JSONB fields
	var variantAttributesJSON []byte
	if len(input.VariantAttributes) > 0 {
		variantAttributesJSON, err = json.Marshal(input.VariantAttributes)
		if err != nil {
			r.logger.Error().Err(err).Msg("failed to marshal variant attributes")
			return nil, fmt.Errorf("variants.create_failed")
		}
	} else {
		variantAttributesJSON = []byte("{}")
	}

	var dimensionsJSON []byte
	if len(input.Dimensions) > 0 {
		dimensionsJSON, err = json.Marshal(input.Dimensions)
		if err != nil {
			r.logger.Error().Err(err).Msg("failed to marshal dimensions")
			return nil, fmt.Errorf("variants.create_failed")
		}
	} else {
		dimensionsJSON = []byte("{}")
	}

	// Determine stock status
	stockStatus := domain.StockStatusOutOfStock
	if input.StockQuantity > 0 {
		if input.LowStockThreshold != nil && input.StockQuantity <= *input.LowStockThreshold {
			stockStatus = domain.StockStatusLowStock
		} else {
			stockStatus = domain.StockStatusInStock
		}
	}

	// Insert variant
	query := `
		INSERT INTO b2c_product_variants (
			product_id, sku, barcode, price, compare_at_price, cost_price,
			stock_quantity, stock_status, low_stock_threshold,
			variant_attributes, weight, dimensions,
			is_active, is_default, view_count, sold_count
		) VALUES (
			$1, $2, $3, $4, $5, $6,
			$7, $8, $9,
			$10, $11, $12,
			true, $13, 0, 0
		)
		RETURNING
			id, product_id, sku, barcode, price, compare_at_price, cost_price,
			stock_quantity, stock_status, low_stock_threshold,
			variant_attributes, weight, dimensions,
			is_active, is_default, view_count, sold_count,
			created_at, updated_at
	`

	var variant domain.ProductVariant
	var sku, barcode sql.NullString
	var price, compareAtPrice, costPrice, weight sql.NullFloat64
	var lowStockThreshold sql.NullInt32
	var returnedVariantAttributesJSON, returnedDimensionsJSON []byte

	err = tx.QueryRowContext(
		ctx,
		query,
		input.ProductID,
		input.SKU,
		input.Barcode,
		input.Price,
		input.CompareAtPrice,
		input.CostPrice,
		input.StockQuantity,
		stockStatus,
		input.LowStockThreshold,
		variantAttributesJSON,
		input.Weight,
		dimensionsJSON,
		input.IsDefault,
	).Scan(
		&variant.ID,
		&variant.ProductID,
		&sku,
		&barcode,
		&price,
		&compareAtPrice,
		&costPrice,
		&variant.StockQuantity,
		&variant.StockStatus,
		&lowStockThreshold,
		&returnedVariantAttributesJSON,
		&weight,
		&returnedDimensionsJSON,
		&variant.IsActive,
		&variant.IsDefault,
		&variant.ViewCount,
		&variant.SoldCount,
		&variant.CreatedAt,
		&variant.UpdatedAt,
	)

	if err != nil {
		// Check for unique constraint violation (duplicate SKU)
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" { // unique_violation
				r.logger.Error().Err(err).Msg("duplicate SKU")
				return nil, fmt.Errorf("variants.sku_duplicate")
			}
		}
		r.logger.Error().Err(err).Msg("failed to create variant")
		return nil, fmt.Errorf("variants.create_failed")
	}

	// Handle nullable fields
	if sku.Valid {
		variant.SKU = &sku.String
	}
	if barcode.Valid {
		variant.Barcode = &barcode.String
	}
	if price.Valid {
		variant.Price = &price.Float64
	}
	if compareAtPrice.Valid {
		variant.CompareAtPrice = &compareAtPrice.Float64
	}
	if costPrice.Valid {
		variant.CostPrice = &costPrice.Float64
	}
	if weight.Valid {
		variant.Weight = &weight.Float64
	}
	if lowStockThreshold.Valid {
		variant.LowStockThreshold = &lowStockThreshold.Int32
	}

	// Parse JSONB fields
	if len(returnedVariantAttributesJSON) > 0 {
		if err := json.Unmarshal(returnedVariantAttributesJSON, &variant.VariantAttributes); err != nil {
			r.logger.Error().Err(err).Msg("failed to unmarshal variant attributes")
			return nil, fmt.Errorf("variants.create_failed")
		}
	}
	if len(returnedDimensionsJSON) > 0 {
		if err := json.Unmarshal(returnedDimensionsJSON, &variant.Dimensions); err != nil {
			r.logger.Error().Err(err).Msg("failed to unmarshal dimensions")
			return nil, fmt.Errorf("variants.create_failed")
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		r.logger.Error().Err(err).Msg("failed to commit transaction")
		return nil, fmt.Errorf("variants.create_failed")
	}

	r.logger.Info().
		Int64("variant_id", variant.ID).
		Int64("product_id", variant.ProductID).
		Msg("product variant created successfully")

	return &variant, nil
}

// UpdateProductVariant updates an existing product variant
func (r *Repository) UpdateProductVariant(ctx context.Context, variantID int64, productID int64, input *domain.UpdateVariantInput) (*domain.ProductVariant, error) {
	r.logger.Debug().
		Int64("variant_id", variantID).
		Int64("product_id", productID).
		Msg("updating product variant")

	// Validate IDs
	if variantID <= 0 {
		return nil, fmt.Errorf("variants.invalid_variant_id")
	}
	if productID <= 0 {
		return nil, fmt.Errorf("variants.invalid_product_id")
	}

	// Start transaction
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		r.logger.Error().Err(err).Msg("failed to begin transaction")
		return nil, fmt.Errorf("variants.update_failed")
	}
	defer tx.Rollback()

	// Check if variant exists and belongs to product
	var exists bool
	err = tx.QueryRowContext(ctx, `
		SELECT EXISTS(
			SELECT 1 FROM b2c_product_variants
			WHERE id = $1 AND product_id = $2
		)
	`, variantID, productID).Scan(&exists)

	if err != nil {
		r.logger.Error().Err(err).Msg("failed to check variant")
		return nil, fmt.Errorf("variants.update_failed")
	}

	if !exists {
		return nil, fmt.Errorf("variants.not_found")
	}

	// Business rule: Cannot disable last active variant
	if input.IsActive != nil && !*input.IsActive {
		var activeCount int
		err = tx.QueryRowContext(ctx, `
			SELECT COUNT(*) FROM b2c_product_variants
			WHERE product_id = $1 AND is_active = true AND id != $2
		`, productID, variantID).Scan(&activeCount)

		if err != nil {
			r.logger.Error().Err(err).Msg("failed to count active variants")
			return nil, fmt.Errorf("variants.update_failed")
		}

		if activeCount == 0 {
			return nil, fmt.Errorf("variants.last_variant")
		}
	}

	// If setting this as default, unset other defaults
	if input.IsDefault != nil && *input.IsDefault {
		_, err = tx.ExecContext(ctx, `
			UPDATE b2c_product_variants
			SET is_default = false, updated_at = NOW()
			WHERE product_id = $1 AND is_default = true AND id != $2
		`, productID, variantID)

		if err != nil {
			r.logger.Error().Err(err).Msg("failed to unset other defaults")
			return nil, fmt.Errorf("variants.update_failed")
		}
	}

	// Build dynamic update query
	updates := []string{}
	args := []interface{}{}
	argPos := 1

	if input.SKU != nil {
		updates = append(updates, fmt.Sprintf("sku = $%d", argPos))
		args = append(args, *input.SKU)
		argPos++
	}

	if input.Barcode != nil {
		updates = append(updates, fmt.Sprintf("barcode = $%d", argPos))
		args = append(args, *input.Barcode)
		argPos++
	}

	if input.Price != nil {
		updates = append(updates, fmt.Sprintf("price = $%d", argPos))
		args = append(args, *input.Price)
		argPos++
	}

	if input.CompareAtPrice != nil {
		updates = append(updates, fmt.Sprintf("compare_at_price = $%d", argPos))
		args = append(args, *input.CompareAtPrice)
		argPos++
	}

	if input.CostPrice != nil {
		updates = append(updates, fmt.Sprintf("cost_price = $%d", argPos))
		args = append(args, *input.CostPrice)
		argPos++
	}

	if input.StockQuantity != nil {
		updates = append(updates, fmt.Sprintf("stock_quantity = $%d", argPos))
		args = append(args, *input.StockQuantity)
		argPos++
	}

	if input.StockStatus != nil {
		updates = append(updates, fmt.Sprintf("stock_status = $%d", argPos))
		args = append(args, *input.StockStatus)
		argPos++
	}

	if input.LowStockThreshold != nil {
		updates = append(updates, fmt.Sprintf("low_stock_threshold = $%d", argPos))
		args = append(args, *input.LowStockThreshold)
		argPos++
	}

	if input.VariantAttributes != nil {
		attributesJSON, err := json.Marshal(input.VariantAttributes)
		if err != nil {
			r.logger.Error().Err(err).Msg("failed to marshal variant attributes")
			return nil, fmt.Errorf("variants.update_failed")
		}
		updates = append(updates, fmt.Sprintf("variant_attributes = $%d", argPos))
		args = append(args, attributesJSON)
		argPos++
	}

	if input.Weight != nil {
		updates = append(updates, fmt.Sprintf("weight = $%d", argPos))
		args = append(args, *input.Weight)
		argPos++
	}

	if input.Dimensions != nil {
		dimensionsJSON, err := json.Marshal(input.Dimensions)
		if err != nil {
			r.logger.Error().Err(err).Msg("failed to marshal dimensions")
			return nil, fmt.Errorf("variants.update_failed")
		}
		updates = append(updates, fmt.Sprintf("dimensions = $%d", argPos))
		args = append(args, dimensionsJSON)
		argPos++
	}

	if input.IsActive != nil {
		updates = append(updates, fmt.Sprintf("is_active = $%d", argPos))
		args = append(args, *input.IsActive)
		argPos++
	}

	if input.IsDefault != nil {
		updates = append(updates, fmt.Sprintf("is_default = $%d", argPos))
		args = append(args, *input.IsDefault)
		argPos++
	}

	if len(updates) == 0 {
		// No updates requested, return current variant
		return r.GetVariantByID(ctx, variantID, &productID)
	}

	// Always update updated_at
	updates = append(updates, "updated_at = NOW()")

	// Add variant ID and product ID as last parameters
	args = append(args, variantID, productID)

	query := fmt.Sprintf(`
		UPDATE b2c_product_variants
		SET %s
		WHERE id = $%d AND product_id = $%d
		RETURNING
			id, product_id, sku, barcode, price, compare_at_price, cost_price,
			stock_quantity, stock_status, low_stock_threshold,
			variant_attributes, weight, dimensions,
			is_active, is_default, view_count, sold_count,
			created_at, updated_at
	`, strings.Join(updates, ", "), argPos, argPos+1)

	var variant domain.ProductVariant
	var sku, barcode sql.NullString
	var price, compareAtPrice, costPrice, weight sql.NullFloat64
	var lowStockThreshold sql.NullInt32
	var variantAttributesJSON, dimensionsJSON []byte

	err = tx.QueryRowContext(ctx, query, args...).Scan(
		&variant.ID,
		&variant.ProductID,
		&sku,
		&barcode,
		&price,
		&compareAtPrice,
		&costPrice,
		&variant.StockQuantity,
		&variant.StockStatus,
		&lowStockThreshold,
		&variantAttributesJSON,
		&weight,
		&dimensionsJSON,
		&variant.IsActive,
		&variant.IsDefault,
		&variant.ViewCount,
		&variant.SoldCount,
		&variant.CreatedAt,
		&variant.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("variants.not_found")
		}
		// Check for unique constraint violation
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" {
				return nil, fmt.Errorf("variants.sku_duplicate")
			}
		}
		r.logger.Error().Err(err).Msg("failed to update variant")
		return nil, fmt.Errorf("variants.update_failed")
	}

	// Handle nullable fields
	if sku.Valid {
		variant.SKU = &sku.String
	}
	if barcode.Valid {
		variant.Barcode = &barcode.String
	}
	if price.Valid {
		variant.Price = &price.Float64
	}
	if compareAtPrice.Valid {
		variant.CompareAtPrice = &compareAtPrice.Float64
	}
	if costPrice.Valid {
		variant.CostPrice = &costPrice.Float64
	}
	if weight.Valid {
		variant.Weight = &weight.Float64
	}
	if lowStockThreshold.Valid {
		variant.LowStockThreshold = &lowStockThreshold.Int32
	}

	// Parse JSONB fields
	if len(variantAttributesJSON) > 0 {
		if err := json.Unmarshal(variantAttributesJSON, &variant.VariantAttributes); err != nil {
			r.logger.Error().Err(err).Msg("failed to unmarshal variant attributes")
			return nil, fmt.Errorf("variants.update_failed")
		}
	}
	if len(dimensionsJSON) > 0 {
		if err := json.Unmarshal(dimensionsJSON, &variant.Dimensions); err != nil {
			r.logger.Error().Err(err).Msg("failed to unmarshal dimensions")
			return nil, fmt.Errorf("variants.update_failed")
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		r.logger.Error().Err(err).Msg("failed to commit transaction")
		return nil, fmt.Errorf("variants.update_failed")
	}

	r.logger.Info().
		Int64("variant_id", variant.ID).
		Msg("product variant updated successfully")

	return &variant, nil
}

// DeleteProductVariant deletes a product variant
func (r *Repository) DeleteProductVariant(ctx context.Context, variantID int64, productID int64) error {
	r.logger.Debug().
		Int64("variant_id", variantID).
		Int64("product_id", productID).
		Msg("deleting product variant")

	// Validate IDs
	if variantID <= 0 {
		return fmt.Errorf("variants.invalid_variant_id")
	}
	if productID <= 0 {
		return fmt.Errorf("variants.invalid_product_id")
	}

	// Start transaction
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		r.logger.Error().Err(err).Msg("failed to begin transaction")
		return fmt.Errorf("variants.delete_failed")
	}
	defer tx.Rollback()

	// Check if variant exists and belongs to product
	var exists bool
	var isDefault bool
	err = tx.QueryRowContext(ctx, `
		SELECT EXISTS(
			SELECT 1 FROM b2c_product_variants
			WHERE id = $1 AND product_id = $2
		),
		COALESCE((SELECT is_default FROM b2c_product_variants WHERE id = $1), false)
	`, variantID, productID).Scan(&exists, &isDefault)

	if err != nil {
		r.logger.Error().Err(err).Msg("failed to check variant")
		return fmt.Errorf("variants.delete_failed")
	}

	if !exists {
		return fmt.Errorf("variants.not_found")
	}

	// Count active variants (excluding current one)
	var activeCount int
	err = tx.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM b2c_product_variants
		WHERE product_id = $1 AND id != $2
	`, productID, variantID).Scan(&activeCount)

	if err != nil {
		r.logger.Error().Err(err).Msg("failed to count variants")
		return fmt.Errorf("variants.delete_failed")
	}

	// Business rule: If this is the last variant, set product.has_variants=false
	if activeCount == 0 {
		_, err = tx.ExecContext(ctx, `
			UPDATE listings
			SET has_variants = false, updated_at = NOW()
			WHERE id = $1 AND source_type = 'b2c'
		`, productID)

		if err != nil {
			r.logger.Error().Err(err).Msg("failed to update product has_variants")
			return fmt.Errorf("variants.delete_failed")
		}
	}

	// Business rule: If deleted variant was default, assign default to another
	if isDefault && activeCount > 0 {
		_, err = tx.ExecContext(ctx, `
			UPDATE b2c_product_variants
			SET is_default = true, updated_at = NOW()
			WHERE product_id = $1 AND id != $2
			ORDER BY id ASC
			LIMIT 1
		`, productID, variantID)

		if err != nil {
			r.logger.Error().Err(err).Msg("failed to assign new default")
			return fmt.Errorf("variants.delete_failed")
		}
	}

	// Delete variant
	result, err := tx.ExecContext(ctx, `
		DELETE FROM b2c_product_variants
		WHERE id = $1 AND product_id = $2
	`, variantID, productID)

	if err != nil {
		r.logger.Error().Err(err).Msg("failed to delete variant")
		return fmt.Errorf("variants.delete_failed")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.logger.Error().Err(err).Msg("failed to get rows affected")
		return fmt.Errorf("variants.delete_failed")
	}

	if rowsAffected == 0 {
		return fmt.Errorf("variants.not_found")
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		r.logger.Error().Err(err).Msg("failed to commit transaction")
		return fmt.Errorf("variants.delete_failed")
	}

	r.logger.Info().
		Int64("variant_id", variantID).
		Int64("product_id", productID).
		Msg("product variant deleted successfully")

	return nil
}

// BulkCreateProductVariants creates multiple product variants in one atomic transaction
func (r *Repository) BulkCreateProductVariants(ctx context.Context, productID int64, inputs []*domain.CreateVariantInput) ([]*domain.ProductVariant, error) {
	r.logger.Debug().
		Int64("product_id", productID).
		Int("count", len(inputs)).
		Msg("bulk creating product variants")

	// Validate inputs
	if productID <= 0 {
		return nil, fmt.Errorf("variants.invalid_product_id")
	}

	if len(inputs) == 0 {
		return nil, fmt.Errorf("variants.bulk_empty")
	}

	if len(inputs) > 1000 {
		return nil, fmt.Errorf("variants.bulk_too_many")
	}

	// Start transaction
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		r.logger.Error().Err(err).Msg("failed to begin transaction")
		return nil, fmt.Errorf("variants.bulk_create_failed")
	}
	defer tx.Rollback()

	// Check if product exists and get has_variants flag
	var hasVariants bool
	err = tx.QueryRowContext(ctx, `
		SELECT has_variants FROM listings
		WHERE id = $1 AND status = 'active' AND deleted_at IS NULL AND source_type = 'b2c'
	`, productID).Scan(&hasVariants)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("variants.product_not_found")
		}
		r.logger.Error().Err(err).Msg("failed to check product")
		return nil, fmt.Errorf("variants.bulk_create_failed")
	}

	// Business rule: Product must have has_variants=true
	if !hasVariants {
		return nil, fmt.Errorf("variants.product_no_variants")
	}

	// Count how many inputs have is_default=true
	defaultCount := 0
	for _, input := range inputs {
		if input.IsDefault {
			defaultCount++
		}
	}

	// Business rule: Only one default variant allowed
	if defaultCount > 1 {
		return nil, fmt.Errorf("variants.multiple_defaults")
	}

	// If any variant is marked as default, unset existing defaults
	if defaultCount == 1 {
		_, err = tx.ExecContext(ctx, `
			UPDATE b2c_product_variants
			SET is_default = false, updated_at = NOW()
			WHERE product_id = $1 AND is_default = true
		`, productID)

		if err != nil {
			r.logger.Error().Err(err).Msg("failed to unset other defaults")
			return nil, fmt.Errorf("variants.bulk_create_failed")
		}
	}

	// Prepare batch insert
	variants := make([]*domain.ProductVariant, 0, len(inputs))

	// Build values for batch insert
	valueStrings := make([]string, 0, len(inputs))
	valueArgs := make([]interface{}, 0, len(inputs)*13)
	argPos := 1

	for _, input := range inputs {
		// Marshal JSONB fields
		var variantAttributesJSON []byte
		if len(input.VariantAttributes) > 0 {
			variantAttributesJSON, err = json.Marshal(input.VariantAttributes)
			if err != nil {
				r.logger.Error().Err(err).Msg("failed to marshal variant attributes")
				return nil, fmt.Errorf("variants.bulk_create_failed")
			}
		} else {
			variantAttributesJSON = []byte("{}")
		}

		var dimensionsJSON []byte
		if len(input.Dimensions) > 0 {
			dimensionsJSON, err = json.Marshal(input.Dimensions)
			if err != nil {
				r.logger.Error().Err(err).Msg("failed to marshal dimensions")
				return nil, fmt.Errorf("variants.bulk_create_failed")
			}
		} else {
			dimensionsJSON = []byte("{}")
		}

		// Determine stock status
		stockStatus := domain.StockStatusOutOfStock
		if input.StockQuantity > 0 {
			if input.LowStockThreshold != nil && input.StockQuantity <= *input.LowStockThreshold {
				stockStatus = domain.StockStatusLowStock
			} else {
				stockStatus = domain.StockStatusInStock
			}
		}

		// Build value string for this variant
		valueStrings = append(valueStrings, fmt.Sprintf(
			"($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)",
			argPos, argPos+1, argPos+2, argPos+3, argPos+4, argPos+5,
			argPos+6, argPos+7, argPos+8, argPos+9, argPos+10, argPos+11, argPos+12,
		))

		// Add arguments
		valueArgs = append(valueArgs,
			productID,
			input.SKU,
			input.Barcode,
			input.Price,
			input.CompareAtPrice,
			input.CostPrice,
			input.StockQuantity,
			stockStatus,
			input.LowStockThreshold,
			variantAttributesJSON,
			input.Weight,
			dimensionsJSON,
			input.IsDefault,
		)

		argPos += 13
	}

	// Execute batch insert
	query := fmt.Sprintf(`
		INSERT INTO b2c_product_variants (
			product_id, sku, barcode, price, compare_at_price, cost_price,
			stock_quantity, stock_status, low_stock_threshold,
			variant_attributes, weight, dimensions, is_default
		) VALUES %s
		RETURNING
			id, product_id, sku, barcode, price, compare_at_price, cost_price,
			stock_quantity, stock_status, low_stock_threshold,
			variant_attributes, weight, dimensions,
			is_active, is_default, view_count, sold_count,
			created_at, updated_at
	`, strings.Join(valueStrings, ", "))

	rows, err := tx.QueryContext(ctx, query, valueArgs...)
	if err != nil {
		// Check for unique constraint violation (duplicate SKU)
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" { // unique_violation
				r.logger.Error().Err(err).Msg("duplicate SKU in bulk create")
				return nil, fmt.Errorf("variants.sku_duplicate")
			}
		}
		r.logger.Error().Err(err).Msg("failed to bulk create variants")
		return nil, fmt.Errorf("variants.bulk_create_failed")
	}
	defer rows.Close()

	// Scan all returned variants
	for rows.Next() {
		var variant domain.ProductVariant
		var sku, barcode sql.NullString
		var price, compareAtPrice, costPrice, weight sql.NullFloat64
		var lowStockThreshold sql.NullInt32
		var variantAttributesJSON, dimensionsJSON []byte

		err = rows.Scan(
			&variant.ID,
			&variant.ProductID,
			&sku,
			&barcode,
			&price,
			&compareAtPrice,
			&costPrice,
			&variant.StockQuantity,
			&variant.StockStatus,
			&lowStockThreshold,
			&variantAttributesJSON,
			&weight,
			&dimensionsJSON,
			&variant.IsActive,
			&variant.IsDefault,
			&variant.ViewCount,
			&variant.SoldCount,
			&variant.CreatedAt,
			&variant.UpdatedAt,
		)

		if err != nil {
			r.logger.Error().Err(err).Msg("failed to scan variant")
			return nil, fmt.Errorf("variants.bulk_create_failed")
		}

		// Handle nullable fields
		if sku.Valid {
			variant.SKU = &sku.String
		}
		if barcode.Valid {
			variant.Barcode = &barcode.String
		}
		if price.Valid {
			variant.Price = &price.Float64
		}
		if compareAtPrice.Valid {
			variant.CompareAtPrice = &compareAtPrice.Float64
		}
		if costPrice.Valid {
			variant.CostPrice = &costPrice.Float64
		}
		if weight.Valid {
			variant.Weight = &weight.Float64
		}
		if lowStockThreshold.Valid {
			variant.LowStockThreshold = &lowStockThreshold.Int32
		}

		// Parse JSONB fields
		if len(variantAttributesJSON) > 0 {
			if err := json.Unmarshal(variantAttributesJSON, &variant.VariantAttributes); err != nil {
				r.logger.Error().Err(err).Msg("failed to unmarshal variant attributes")
				return nil, fmt.Errorf("variants.bulk_create_failed")
			}
		}
		if len(dimensionsJSON) > 0 {
			if err := json.Unmarshal(dimensionsJSON, &variant.Dimensions); err != nil {
				r.logger.Error().Err(err).Msg("failed to unmarshal dimensions")
				return nil, fmt.Errorf("variants.bulk_create_failed")
			}
		}

		variants = append(variants, &variant)
	}

	if err = rows.Err(); err != nil {
		r.logger.Error().Err(err).Msg("error iterating variants")
		return nil, fmt.Errorf("variants.bulk_create_failed")
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		r.logger.Error().Err(err).Msg("failed to commit transaction")
		return nil, fmt.Errorf("variants.bulk_create_failed")
	}

	r.logger.Info().
		Int("count", len(variants)).
		Int64("product_id", productID).
		Msg("product variants bulk created successfully")

	return variants, nil
}
