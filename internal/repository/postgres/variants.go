package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/sveturs/listings/internal/domain"
)

// CreateVariant creates a new product variant
func (r *Repository) CreateVariant(ctx context.Context, variant *domain.Variant) (*domain.Variant, error) {
	query := `
		INSERT INTO b2c_product_variants (
			product_id, sku, barcode, price, compare_at_price, cost_price,
			stock_quantity, stock_status, low_stock_threshold,
			variant_attributes, weight, dimensions, is_active, is_default
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10::jsonb, $11, $12::jsonb, $13, $14)
		RETURNING id, product_id, sku, barcode, price, compare_at_price, cost_price,
		          stock_quantity, stock_status, low_stock_threshold,
		          variant_attributes, weight, dimensions, is_active, is_default,
		          view_count, sold_count, created_at, updated_at
	`

	var created domain.Variant
	err := r.db.QueryRowContext(
		ctx,
		query,
		variant.ProductID,
		variant.SKU,
		variant.Barcode,
		variant.Price,
		variant.CompareAtPrice,
		variant.CostPrice,
		variant.StockQuantity,
		variant.StockStatus,
		variant.LowStockThreshold,
		variant.VariantAttributes,
		variant.Weight,
		variant.Dimensions,
		variant.IsActive,
		variant.IsDefault,
	).Scan(
		&created.ID,
		&created.ProductID,
		&created.SKU,
		&created.Barcode,
		&created.Price,
		&created.CompareAtPrice,
		&created.CostPrice,
		&created.StockQuantity,
		&created.StockStatus,
		&created.LowStockThreshold,
		&created.VariantAttributes,
		&created.Weight,
		&created.Dimensions,
		&created.IsActive,
		&created.IsDefault,
		&created.ViewCount,
		&created.SoldCount,
		&created.CreatedAt,
		&created.UpdatedAt,
	)

	if err != nil {
		r.logger.Error().Err(err).Int64("product_id", variant.ProductID).Msg("failed to create variant")
		return nil, fmt.Errorf("failed to create variant: %w", err)
	}

	r.logger.Info().Int64("variant_id", created.ID).Int64("product_id", variant.ProductID).Msg("variant created")
	return &created, nil
}

// GetVariant retrieves a variant by ID
func (r *Repository) GetVariant(ctx context.Context, id int64) (*domain.Variant, error) {
	query := `
		SELECT id, product_id, sku, barcode, price, compare_at_price, cost_price,
		       stock_quantity, stock_status, low_stock_threshold,
		       variant_attributes, weight, dimensions, is_active, is_default,
		       view_count, sold_count, created_at, updated_at
		FROM b2c_product_variants
		WHERE id = $1
	`

	var variant domain.Variant
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&variant.ID,
		&variant.ProductID,
		&variant.SKU,
		&variant.Barcode,
		&variant.Price,
		&variant.CompareAtPrice,
		&variant.CostPrice,
		&variant.StockQuantity,
		&variant.StockStatus,
		&variant.LowStockThreshold,
		&variant.VariantAttributes,
		&variant.Weight,
		&variant.Dimensions,
		&variant.IsActive,
		&variant.IsDefault,
		&variant.ViewCount,
		&variant.SoldCount,
		&variant.CreatedAt,
		&variant.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("variant not found: %w", err)
		}
		r.logger.Error().Err(err).Int64("variant_id", id).Msg("failed to get variant")
		return nil, fmt.Errorf("failed to get variant: %w", err)
	}

	return &variant, nil
}

// UpdateB2CVariant updates an existing B2C variant with partial updates
// Uses COALESCE pattern to handle optional updates (**T pattern)
func (r *Repository) UpdateB2CVariant(ctx context.Context, id int64, update *domain.VariantUpdate) (*domain.Variant, error) {
	// Build dynamic update query
	updates := []string{}
	args := []interface{}{}
	argPos := 1

	// Handle **string fields (nil = no update, *nil = set NULL, *value = update)
	if update.SKU != nil {
		if *update.SKU == nil {
			updates = append(updates, fmt.Sprintf("sku = NULL"))
		} else {
			updates = append(updates, fmt.Sprintf("sku = $%d", argPos))
			args = append(args, **update.SKU)
			argPos++
		}
	}

	if update.Barcode != nil {
		if *update.Barcode == nil {
			updates = append(updates, fmt.Sprintf("barcode = NULL"))
		} else {
			updates = append(updates, fmt.Sprintf("barcode = $%d", argPos))
			args = append(args, **update.Barcode)
			argPos++
		}
	}

	if update.Price != nil {
		if *update.Price == nil {
			updates = append(updates, fmt.Sprintf("price = NULL"))
		} else {
			updates = append(updates, fmt.Sprintf("price = $%d", argPos))
			args = append(args, **update.Price)
			argPos++
		}
	}

	if update.CompareAtPrice != nil {
		if *update.CompareAtPrice == nil {
			updates = append(updates, fmt.Sprintf("compare_at_price = NULL"))
		} else {
			updates = append(updates, fmt.Sprintf("compare_at_price = $%d", argPos))
			args = append(args, **update.CompareAtPrice)
			argPos++
		}
	}

	if update.CostPrice != nil {
		if *update.CostPrice == nil {
			updates = append(updates, fmt.Sprintf("cost_price = NULL"))
		} else {
			updates = append(updates, fmt.Sprintf("cost_price = $%d", argPos))
			args = append(args, **update.CostPrice)
			argPos++
		}
	}

	if update.StockQuantity != nil {
		updates = append(updates, fmt.Sprintf("stock_quantity = $%d", argPos))
		args = append(args, *update.StockQuantity)
		argPos++
	}

	if update.StockStatus != nil {
		updates = append(updates, fmt.Sprintf("stock_status = $%d", argPos))
		args = append(args, *update.StockStatus)
		argPos++
	}

	if update.LowStockThreshold != nil {
		if *update.LowStockThreshold == nil {
			updates = append(updates, fmt.Sprintf("low_stock_threshold = NULL"))
		} else {
			updates = append(updates, fmt.Sprintf("low_stock_threshold = $%d", argPos))
			args = append(args, **update.LowStockThreshold)
			argPos++
		}
	}

	if update.VariantAttributes != nil {
		updates = append(updates, fmt.Sprintf("variant_attributes = $%d::jsonb", argPos))
		args = append(args, *update.VariantAttributes)
		argPos++
	}

	if update.Weight != nil {
		if *update.Weight == nil {
			updates = append(updates, fmt.Sprintf("weight = NULL"))
		} else {
			updates = append(updates, fmt.Sprintf("weight = $%d", argPos))
			args = append(args, **update.Weight)
			argPos++
		}
	}

	if update.Dimensions != nil {
		if *update.Dimensions == nil {
			updates = append(updates, fmt.Sprintf("dimensions = NULL"))
		} else {
			updates = append(updates, fmt.Sprintf("dimensions = $%d::jsonb", argPos))
			args = append(args, **update.Dimensions)
			argPos++
		}
	}

	if update.IsActive != nil {
		updates = append(updates, fmt.Sprintf("is_active = $%d", argPos))
		args = append(args, *update.IsActive)
		argPos++
	}

	if update.IsDefault != nil {
		updates = append(updates, fmt.Sprintf("is_default = $%d", argPos))
		args = append(args, *update.IsDefault)
		argPos++
	}

	// If no updates, just return current variant
	if len(updates) == 0 {
		return r.GetVariant(ctx, id)
	}

	// Add variant ID as last parameter
	args = append(args, id)

	query := fmt.Sprintf(`
		UPDATE b2c_product_variants
		SET %s, updated_at = CURRENT_TIMESTAMP
		WHERE id = $%d
		RETURNING id, product_id, sku, barcode, price, compare_at_price, cost_price,
		          stock_quantity, stock_status, low_stock_threshold,
		          variant_attributes, weight, dimensions, is_active, is_default,
		          view_count, sold_count, created_at, updated_at
	`, strings.Join(updates, ", "), argPos)

	var updated domain.Variant
	err := r.db.QueryRowContext(ctx, query, args...).Scan(
		&updated.ID,
		&updated.ProductID,
		&updated.SKU,
		&updated.Barcode,
		&updated.Price,
		&updated.CompareAtPrice,
		&updated.CostPrice,
		&updated.StockQuantity,
		&updated.StockStatus,
		&updated.LowStockThreshold,
		&updated.VariantAttributes,
		&updated.Weight,
		&updated.Dimensions,
		&updated.IsActive,
		&updated.IsDefault,
		&updated.ViewCount,
		&updated.SoldCount,
		&updated.CreatedAt,
		&updated.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("variant not found: %w", err)
		}
		r.logger.Error().Err(err).Int64("variant_id", id).Msg("failed to update variant")
		return nil, fmt.Errorf("failed to update variant: %w", err)
	}

	r.logger.Info().Int64("variant_id", updated.ID).Msg("variant updated")
	return &updated, nil
}

// DeleteB2CVariant deletes a B2C variant by ID
func (r *Repository) DeleteB2CVariant(ctx context.Context, id int64) error {
	query := `DELETE FROM b2c_product_variants WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		r.logger.Error().Err(err).Int64("variant_id", id).Msg("failed to delete variant")
		return fmt.Errorf("failed to delete variant: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("variant not found")
	}

	r.logger.Info().Int64("variant_id", id).Msg("variant deleted")
	return nil
}

// ListVariants retrieves all variants for a product with optional filters
func (r *Repository) ListVariants(ctx context.Context, filters *domain.VariantFilters) ([]*domain.Variant, error) {
	// Build WHERE clause
	whereConditions := []string{"product_id = $1"}
	args := []interface{}{filters.ProductID}
	argPos := 2

	if filters.ActiveOnly != nil && *filters.ActiveOnly {
		whereConditions = append(whereConditions, fmt.Sprintf("is_active = $%d", argPos))
		args = append(args, true)
		argPos++
	}

	if filters.StockStatus != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("stock_status = $%d", argPos))
		args = append(args, *filters.StockStatus)
		argPos++
	}

	whereClause := strings.Join(whereConditions, " AND ")

	query := fmt.Sprintf(`
		SELECT id, product_id, sku, barcode, price, compare_at_price, cost_price,
		       stock_quantity, stock_status, low_stock_threshold,
		       variant_attributes, weight, dimensions, is_active, is_default,
		       view_count, sold_count, created_at, updated_at
		FROM b2c_product_variants
		WHERE %s
		ORDER BY is_default DESC, created_at ASC
	`, whereClause)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		r.logger.Error().Err(err).Int64("product_id", filters.ProductID).Msg("failed to list variants")
		return nil, fmt.Errorf("failed to list variants: %w", err)
	}
	defer rows.Close()

	var variants []*domain.Variant
	for rows.Next() {
		var v domain.Variant
		err := rows.Scan(
			&v.ID,
			&v.ProductID,
			&v.SKU,
			&v.Barcode,
			&v.Price,
			&v.CompareAtPrice,
			&v.CostPrice,
			&v.StockQuantity,
			&v.StockStatus,
			&v.LowStockThreshold,
			&v.VariantAttributes,
			&v.Weight,
			&v.Dimensions,
			&v.IsActive,
			&v.IsDefault,
			&v.ViewCount,
			&v.SoldCount,
			&v.CreatedAt,
			&v.UpdatedAt,
		)
		if err != nil {
			r.logger.Error().Err(err).Msg("failed to scan variant")
			return nil, fmt.Errorf("failed to scan variant: %w", err)
		}
		variants = append(variants, &v)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return variants, nil
}
