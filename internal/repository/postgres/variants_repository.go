package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/sveturs/listings/internal/domain"
)

// CreateVariants creates multiple variants for a listing
func (r *Repository) CreateVariants(ctx context.Context, variants []*domain.ListingVariant) error {
	if len(variants) == 0 {
		return nil
	}

	// Use transaction for batch insert (WithTransaction accepts *sqlx.Tx)
	return r.WithTransaction(ctx, func(tx *sqlx.Tx) error {
		query := `
			INSERT INTO marketplace_listing_variants (
				listing_id, sku, price, stock, attributes, image_url, is_active
			)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
		`

		stmt, err := tx.PreparexContext(ctx, query)
		if err != nil {
			return fmt.Errorf("failed to prepare statement: %w", err)
		}
		defer stmt.Close()

		for _, variant := range variants {
			// Marshal attributes to JSON
			var attributesJSON sql.NullString
			if len(variant.Attributes) > 0 {
				attrBytes, err := json.Marshal(variant.Attributes)
				if err != nil {
					return fmt.Errorf("failed to marshal attributes: %w", err)
				}
				attributesJSON = sql.NullString{String: string(attrBytes), Valid: true}
			}

			var price sql.NullFloat64
			if variant.Price != nil {
				price = sql.NullFloat64{Float64: *variant.Price, Valid: true}
			}

			var stock sql.NullInt32
			if variant.Stock != nil {
				stock = sql.NullInt32{Int32: *variant.Stock, Valid: true}
			}

			var imageURL sql.NullString
			if variant.ImageURL != nil {
				imageURL = sql.NullString{String: *variant.ImageURL, Valid: true}
			}

			_, err = stmt.ExecContext(ctx,
				variant.ListingID,
				variant.SKU,
				price,
				stock,
				attributesJSON,
				imageURL,
				variant.IsActive,
			)
			if err != nil {
				r.logger.Error().Err(err).Int64("listing_id", variant.ListingID).Str("sku", variant.SKU).Msg("failed to insert variant")
				return fmt.Errorf("failed to insert variant: %w", err)
			}
		}

		r.logger.Info().Int("count", len(variants)).Msg("variants created")
		return nil
	})
}

// GetVariants retrieves all variants for a listing
func (r *Repository) GetVariants(ctx context.Context, listingID int64) ([]*domain.ListingVariant, error) {
	query := `
		SELECT id, listing_id, sku, price, stock, attributes, image_url, is_active, created_at, updated_at
		FROM marketplace_listing_variants
		WHERE listing_id = $1
		ORDER BY id ASC
	`

	rows, err := r.db.QueryxContext(ctx, query, listingID)
	if err != nil {
		r.logger.Error().Err(err).Int64("listing_id", listingID).Msg("failed to query variants")
		return nil, fmt.Errorf("failed to query variants: %w", err)
	}
	defer rows.Close()

	var variants []*domain.ListingVariant
	for rows.Next() {
		variant := &domain.ListingVariant{}
		var price sql.NullFloat64
		var stock sql.NullInt32
		var attributesJSON sql.NullString
		var imageURL sql.NullString

		err := rows.Scan(
			&variant.ID,
			&variant.ListingID,
			&variant.SKU,
			&price,
			&stock,
			&attributesJSON,
			&imageURL,
			&variant.IsActive,
			&variant.CreatedAt,
			&variant.UpdatedAt,
		)
		if err != nil {
			r.logger.Error().Err(err).Msg("failed to scan variant")
			return nil, fmt.Errorf("failed to scan variant: %w", err)
		}

		// Handle nullable fields
		if price.Valid {
			variant.Price = &price.Float64
		}
		if stock.Valid {
			variant.Stock = &stock.Int32
		}
		if attributesJSON.Valid && attributesJSON.String != "" {
			var attrs map[string]string
			if err := json.Unmarshal([]byte(attributesJSON.String), &attrs); err == nil {
				variant.Attributes = attrs
			}
		}
		if imageURL.Valid {
			variant.ImageURL = &imageURL.String
		}

		variants = append(variants, variant)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return variants, nil
}

// UpdateVariant updates a specific variant
func (r *Repository) UpdateVariant(ctx context.Context, variant *domain.ListingVariant) error {
	// Marshal attributes to JSON
	var attributesJSON sql.NullString
	if len(variant.Attributes) > 0 {
		attrBytes, err := json.Marshal(variant.Attributes)
		if err != nil {
			return fmt.Errorf("failed to marshal attributes: %w", err)
		}
		attributesJSON = sql.NullString{String: string(attrBytes), Valid: true}
	}

	var price sql.NullFloat64
	if variant.Price != nil {
		price = sql.NullFloat64{Float64: *variant.Price, Valid: true}
	}

	var stock sql.NullInt32
	if variant.Stock != nil {
		stock = sql.NullInt32{Int32: *variant.Stock, Valid: true}
	}

	var imageURL sql.NullString
	if variant.ImageURL != nil {
		imageURL = sql.NullString{String: *variant.ImageURL, Valid: true}
	}

	query := `
		UPDATE marketplace_listing_variants
		SET sku = $1, price = $2, stock = $3, attributes = $4, image_url = $5, is_active = $6, updated_at = CURRENT_TIMESTAMP
		WHERE id = $7
	`

	result, err := r.db.ExecContext(ctx, query,
		variant.SKU,
		price,
		stock,
		attributesJSON,
		imageURL,
		variant.IsActive,
		variant.ID,
	)
	if err != nil {
		r.logger.Error().Err(err).Int64("variant_id", variant.ID).Msg("failed to update variant")
		return fmt.Errorf("failed to update variant: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("variant not found")
	}

	r.logger.Info().Int64("variant_id", variant.ID).Msg("variant updated")
	return nil
}

// DeleteVariant removes a variant from a listing
func (r *Repository) DeleteVariant(ctx context.Context, variantID int64) error {
	query := `
		DELETE FROM marketplace_listing_variants
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query, variantID)
	if err != nil {
		r.logger.Error().Err(err).Int64("variant_id", variantID).Msg("failed to delete variant")
		return fmt.Errorf("failed to delete variant: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("variant not found")
	}

	r.logger.Info().Int64("variant_id", variantID).Msg("variant deleted")
	return nil
}
