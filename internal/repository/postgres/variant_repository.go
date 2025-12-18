// Package postgres implements PostgreSQL repositories for the listings microservice.
// This file contains the VariantRepository implementation for product variants.
package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
	"github.com/vondi-global/listings/internal/domain"
)

// VariantRepository handles database operations for product variants
type VariantRepository struct {
	db     *sqlx.DB
	logger zerolog.Logger
}

// NewVariantRepository creates a new variant repository
func NewVariantRepository(db *sqlx.DB, logger zerolog.Logger) *VariantRepository {
	return &VariantRepository{
		db:     db,
		logger: logger,
	}
}

// Create creates a new product variant with attributes
func (r *VariantRepository) Create(ctx context.Context, input *domain.CreateVariantInputV2) (*domain.ProductVariantV2, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Create variant
	variant := &domain.ProductVariantV2{
		ID:              uuid.New(),
		ProductID:       input.ProductID,
		SKU:             input.SKU,
		Price:           input.Price,
		CompareAtPrice:  input.CompareAtPrice,
		StockQuantity:   input.StockQuantity,
		ReservedQuantity: 0,
		LowStockAlert:   input.LowStockAlert,
		WeightGrams:     input.WeightGrams,
		Barcode:         input.Barcode,
		IsDefault:       input.IsDefault,
		Position:        input.Position,
		Status:          domain.VariantStatusActive,
	}

	query := `
		INSERT INTO product_variants (
			id, product_id, sku, price, compare_at_price,
			stock_quantity, reserved_quantity, low_stock_alert,
			weight_grams, barcode, is_default, position, status
		) VALUES (
			$1, $2, $3, $4, $5,
			$6, $7, $8,
			$9, $10, $11, $12, $13
		)
		RETURNING created_at, updated_at
	`

	err = tx.QueryRowContext(ctx, query,
		variant.ID, variant.ProductID, variant.SKU, variant.Price, variant.CompareAtPrice,
		variant.StockQuantity, variant.ReservedQuantity, variant.LowStockAlert,
		variant.WeightGrams, variant.Barcode, variant.IsDefault, variant.Position, variant.Status,
	).Scan(&variant.CreatedAt, &variant.UpdatedAt)

	if err != nil {
		r.logger.Error().Err(err).Str("sku", input.SKU).Msg("failed to create variant")
		return nil, fmt.Errorf("failed to create variant: %w", err)
	}

	// Create variant attributes
	if len(input.Attributes) > 0 {
		attrQuery := `
			INSERT INTO variant_attribute_values (
				id, variant_id, attribute_id,
				value_text, value_number, value_boolean, value_date, value_json
			) VALUES (
				$1, $2, $3, $4, $5, $6, $7, $8
			)
		`

		for _, attr := range input.Attributes {
			attrID := uuid.New()
			_, err = tx.ExecContext(ctx, attrQuery,
				attrID, variant.ID, attr.AttributeID,
				attr.ValueText, attr.ValueNumber, attr.ValueBoolean, attr.ValueDate, attr.ValueJSON,
			)
			if err != nil {
				r.logger.Error().Err(err).Int32("attribute_id", attr.AttributeID).Msg("failed to create variant attribute")
				return nil, fmt.Errorf("failed to create variant attribute: %w", err)
			}
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Load created variant with attributes
	return r.GetByID(ctx, variant.ID.String())
}

// GetByID retrieves a variant by ID
func (r *VariantRepository) GetByID(ctx context.Context, id string) (*domain.ProductVariantV2, error) {
	variantID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid variant ID: %w", err)
	}

	query := `
		SELECT
			id, product_id, sku, price, compare_at_price,
			stock_quantity, reserved_quantity, low_stock_alert,
			weight_grams, barcode, is_default, position, status,
			created_at, updated_at
		FROM product_variants
		WHERE id = $1
	`

	variant := &domain.ProductVariantV2{}

	err = r.db.QueryRowContext(ctx, query, variantID).Scan(
		&variant.ID, &variant.ProductID, &variant.SKU, &variant.Price, &variant.CompareAtPrice,
		&variant.StockQuantity, &variant.ReservedQuantity, &variant.LowStockAlert,
		&variant.WeightGrams, &variant.Barcode, &variant.IsDefault, &variant.Position, &variant.Status,
		&variant.CreatedAt, &variant.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrVariantNotFound
	}
	if err != nil {
		r.logger.Error().Err(err).Str("id", id).Msg("failed to get variant by ID")
		return nil, fmt.Errorf("failed to get variant: %w", err)
	}

	// Load attributes
	if err = r.loadVariantAttributes(ctx, variant); err != nil {
		r.logger.Warn().Err(err).Msg("failed to load variant attributes")
	}

	return variant, nil
}

// GetBySKU retrieves a variant by SKU
func (r *VariantRepository) GetBySKU(ctx context.Context, sku string) (*domain.ProductVariantV2, error) {
	query := `
		SELECT
			id, product_id, sku, price, compare_at_price,
			stock_quantity, reserved_quantity, low_stock_alert,
			weight_grams, barcode, is_default, position, status,
			created_at, updated_at
		FROM product_variants
		WHERE sku = $1
	`

	variant := &domain.ProductVariantV2{}

	err := r.db.QueryRowContext(ctx, query, sku).Scan(
		&variant.ID, &variant.ProductID, &variant.SKU, &variant.Price, &variant.CompareAtPrice,
		&variant.StockQuantity, &variant.ReservedQuantity, &variant.LowStockAlert,
		&variant.WeightGrams, &variant.Barcode, &variant.IsDefault, &variant.Position, &variant.Status,
		&variant.CreatedAt, &variant.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrVariantNotFound
	}
	if err != nil {
		r.logger.Error().Err(err).Str("sku", sku).Msg("failed to get variant by SKU")
		return nil, fmt.Errorf("failed to get variant: %w", err)
	}

	// Load attributes
	if err = r.loadVariantAttributes(ctx, variant); err != nil {
		r.logger.Warn().Err(err).Msg("failed to load variant attributes")
	}

	return variant, nil
}

// ListByProduct retrieves all variants for a product
func (r *VariantRepository) ListByProduct(ctx context.Context, filter *domain.ListVariantsFilter) ([]*domain.ProductVariantV2, error) {
	query := `
		SELECT
			id, product_id, sku, price, compare_at_price,
			stock_quantity, reserved_quantity, low_stock_alert,
			weight_grams, barcode, is_default, position, status,
			created_at, updated_at
		FROM product_variants
		WHERE product_id = $1
	`

	args := []interface{}{filter.ProductID}

	if filter.ActiveOnly {
		query += " AND status = $2"
		args = append(args, domain.VariantStatusActive)
	}

	if filter.InStockOnly {
		if filter.ActiveOnly {
			query += " AND (stock_quantity - reserved_quantity) > 0"
		} else {
			query += " AND (stock_quantity - reserved_quantity) > 0"
		}
	}

	query += " ORDER BY position ASC, created_at ASC"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		r.logger.Error().Err(err).Interface("filter", filter).Msg("failed to list variants")
		return nil, fmt.Errorf("failed to list variants: %w", err)
	}
	defer rows.Close()

	variants := make([]*domain.ProductVariantV2, 0)

	for rows.Next() {
		variant := &domain.ProductVariantV2{}
		err = rows.Scan(
			&variant.ID, &variant.ProductID, &variant.SKU, &variant.Price, &variant.CompareAtPrice,
			&variant.StockQuantity, &variant.ReservedQuantity, &variant.LowStockAlert,
			&variant.WeightGrams, &variant.Barcode, &variant.IsDefault, &variant.Position, &variant.Status,
			&variant.CreatedAt, &variant.UpdatedAt,
		)
		if err != nil {
			r.logger.Error().Err(err).Msg("failed to scan variant row")
			continue
		}

		// Load attributes if requested
		if filter.IncludeAttributes {
			if err = r.loadVariantAttributes(ctx, variant); err != nil {
				r.logger.Warn().Err(err).Msg("failed to load variant attributes")
			}
		}

		variants = append(variants, variant)
	}

	return variants, nil
}

// FindByAttributes finds a variant by attribute combination
func (r *VariantRepository) FindByAttributes(ctx context.Context, filter *domain.FindVariantByAttributesFilter) (*domain.ProductVariantV2, error) {
	// First, get all variants for the product with attributes
	variants, err := r.ListByProduct(ctx, &domain.ListVariantsFilter{
		ProductID:         filter.ProductID,
		ActiveOnly:        true,
		IncludeAttributes: true,
	})

	if err != nil {
		return nil, err
	}

	// Match variant by attributes
	for _, variant := range variants {
		if variant.MatchesAttributes(filter.Attributes) {
			return variant, nil
		}
	}

	return nil, domain.ErrVariantNotFound
}

// Update updates a variant
func (r *VariantRepository) Update(ctx context.Context, id string, input *domain.UpdateVariantInputV2) (*domain.ProductVariantV2, error) {
	variantID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid variant ID: %w", err)
	}

	// Build dynamic update query
	updates := make([]string, 0)
	args := make([]interface{}, 0)
	argIdx := 1

	if input.SKU != nil {
		updates = append(updates, fmt.Sprintf("sku = $%d", argIdx))
		args = append(args, *input.SKU)
		argIdx++
	}

	if input.Price != nil {
		updates = append(updates, fmt.Sprintf("price = $%d", argIdx))
		args = append(args, *input.Price)
		argIdx++
	}

	if input.CompareAtPrice != nil {
		updates = append(updates, fmt.Sprintf("compare_at_price = $%d", argIdx))
		args = append(args, *input.CompareAtPrice)
		argIdx++
	}

	if input.StockQuantity != nil {
		updates = append(updates, fmt.Sprintf("stock_quantity = $%d", argIdx))
		args = append(args, *input.StockQuantity)
		argIdx++
	}

	if input.LowStockAlert != nil {
		updates = append(updates, fmt.Sprintf("low_stock_alert = $%d", argIdx))
		args = append(args, *input.LowStockAlert)
		argIdx++
	}

	if input.WeightGrams != nil {
		updates = append(updates, fmt.Sprintf("weight_grams = $%d", argIdx))
		args = append(args, *input.WeightGrams)
		argIdx++
	}

	if input.Barcode != nil {
		updates = append(updates, fmt.Sprintf("barcode = $%d", argIdx))
		args = append(args, *input.Barcode)
		argIdx++
	}

	if input.IsDefault != nil {
		updates = append(updates, fmt.Sprintf("is_default = $%d", argIdx))
		args = append(args, *input.IsDefault)
		argIdx++
	}

	if input.Position != nil {
		updates = append(updates, fmt.Sprintf("position = $%d", argIdx))
		args = append(args, *input.Position)
		argIdx++
	}

	if input.Status != nil {
		updates = append(updates, fmt.Sprintf("status = $%d", argIdx))
		args = append(args, *input.Status)
		argIdx++
	}

	if len(updates) == 0 {
		return r.GetByID(ctx, id)
	}

	args = append(args, variantID)

	query := fmt.Sprintf(`
		UPDATE product_variants
		SET %s, updated_at = CURRENT_TIMESTAMP
		WHERE id = $%d
	`, joinUpdates(updates), argIdx)

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		r.logger.Error().Err(err).Str("id", id).Msg("failed to update variant")
		return nil, fmt.Errorf("failed to update variant: %w", err)
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		return nil, domain.ErrVariantNotFound
	}

	return r.GetByID(ctx, id)
}

// Delete deletes a variant
func (r *VariantRepository) Delete(ctx context.Context, id string) error {
	variantID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid variant ID: %w", err)
	}

	query := `DELETE FROM product_variants WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, variantID)
	if err != nil {
		r.logger.Error().Err(err).Str("id", id).Msg("failed to delete variant")
		return fmt.Errorf("failed to delete variant: %w", err)
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		return domain.ErrVariantNotFound
	}

	return nil
}

// GetForUpdate retrieves a variant with row-level lock for update (used in stock operations)
func (r *VariantRepository) GetForUpdate(ctx context.Context, tx *sqlx.Tx, id string) (*domain.ProductVariantV2, error) {
	variantID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid variant ID: %w", err)
	}

	query := `
		SELECT
			id, product_id, sku, price, compare_at_price,
			stock_quantity, reserved_quantity, low_stock_alert,
			weight_grams, barcode, is_default, position, status,
			created_at, updated_at
		FROM product_variants
		WHERE id = $1
		FOR UPDATE
	`

	variant := &domain.ProductVariantV2{}

	err = tx.QueryRowContext(ctx, query, variantID).Scan(
		&variant.ID, &variant.ProductID, &variant.SKU, &variant.Price, &variant.CompareAtPrice,
		&variant.StockQuantity, &variant.ReservedQuantity, &variant.LowStockAlert,
		&variant.WeightGrams, &variant.Barcode, &variant.IsDefault, &variant.Position, &variant.Status,
		&variant.CreatedAt, &variant.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrVariantNotFound
	}
	if err != nil {
		r.logger.Error().Err(err).Str("id", id).Msg("failed to get variant for update")
		return nil, fmt.Errorf("failed to get variant for update: %w", err)
	}

	return variant, nil
}

// loadVariantAttributes loads attributes for a variant
func (r *VariantRepository) loadVariantAttributes(ctx context.Context, variant *domain.ProductVariantV2) error {
	query := `
		SELECT
			id, variant_id, attribute_id,
			value_text, value_number, value_boolean, value_date, value_json,
			created_at, updated_at
		FROM variant_attribute_values
		WHERE variant_id = $1
		ORDER BY attribute_id
	`

	rows, err := r.db.QueryContext(ctx, query, variant.ID)
	if err != nil {
		return err
	}
	defer rows.Close()

	attributes := make([]*domain.VariantAttributeValueV2, 0)

	for rows.Next() {
		attr := &domain.VariantAttributeValueV2{}
		err = rows.Scan(
			&attr.ID, &attr.VariantID, &attr.AttributeID,
			&attr.ValueText, &attr.ValueNumber, &attr.ValueBoolean, &attr.ValueDate, &attr.ValueJSON,
			&attr.CreatedAt, &attr.UpdatedAt,
		)
		if err != nil {
			r.logger.Error().Err(err).Msg("failed to scan attribute row")
			continue
		}

		attributes = append(attributes, attr)
	}

	variant.Attributes = attributes

	return nil
}

// Helper function to join SQL updates
func joinUpdates(updates []string) string {
	result := ""
	for i, u := range updates {
		if i > 0 {
			result += ", "
		}
		result += u
	}
	return result
}
