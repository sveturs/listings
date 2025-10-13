// backend/internal/proj/c2c/storage/postgres/listings_variants.go
package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"backend/internal/domain/models"
)

// CreateListingVariants создает варианты для листинга
func (s *Storage) CreateListingVariants(ctx context.Context, listingID int, variants []models.MarketplaceListingVariant) error {
	if len(variants) == 0 {
		return nil
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil && !errors.Is(rollbackErr, sql.ErrTxDone) {
			log.Printf("Failed to rollback transaction: %v", rollbackErr)
		}
	}()

	for _, variant := range variants {
		attributesJSON, err := json.Marshal(variant.Attributes)
		if err != nil {
			return fmt.Errorf("failed to marshal variant attributes: %w", err)
		}

		_, err = tx.Exec(ctx, `
			INSERT INTO c2c_listing_variants (listing_id, sku, price, stock, attributes, image_url, is_active)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
		`, listingID, variant.SKU, variant.Price, variant.Stock, attributesJSON, variant.ImageURL, true)
		if err != nil {
			return fmt.Errorf("failed to insert variant: %w", err)
		}
	}

	return tx.Commit(ctx)
}

// GetListingVariants получает все активные варианты листинга
func (s *Storage) GetListingVariants(ctx context.Context, listingID int) ([]models.MarketplaceListingVariant, error) {
	query := `
		SELECT id, listing_id, sku, price, stock, attributes, image_url, is_active,
		       created_at::text, updated_at::text
		FROM c2c_listing_variants
		WHERE listing_id = $1 AND is_active = true
		ORDER BY id
	`

	rows, err := s.pool.Query(ctx, query, listingID)
	if err != nil {
		return nil, fmt.Errorf("failed to query variants: %w", err)
	}
	defer rows.Close()

	var variants []models.MarketplaceListingVariant
	for rows.Next() {
		var variant models.MarketplaceListingVariant
		var attributesJSON []byte

		err := rows.Scan(
			&variant.ID, &variant.ListingID, &variant.SKU, &variant.Price, &variant.Stock,
			&attributesJSON, &variant.ImageURL, &variant.IsActive,
			&variant.CreatedAt, &variant.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan variant: %w", err)
		}

		if len(attributesJSON) > 0 {
			err = json.Unmarshal(attributesJSON, &variant.Attributes)
			if err != nil {
				return nil, fmt.Errorf("failed to unmarshal variant attributes: %w", err)
			}
		}

		variants = append(variants, variant)
	}

	return variants, rows.Err()
}

// UpdateListingVariant обновляет вариант листинга
func (s *Storage) UpdateListingVariant(ctx context.Context, variant *models.MarketplaceListingVariant) error {
	attributesJSON, err := json.Marshal(variant.Attributes)
	if err != nil {
		return fmt.Errorf("failed to marshal variant attributes: %w", err)
	}

	query := `
		UPDATE c2c_listing_variants
		SET sku = $1, price = $2, stock = $3, attributes = $4, image_url = $5, is_active = $6
		WHERE id = $7
	`

	result, err := s.pool.Exec(ctx, query,
		variant.SKU, variant.Price, variant.Stock, attributesJSON,
		variant.ImageURL, variant.IsActive, variant.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update variant: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("variant not found")
	}

	return nil
}

// DeleteListingVariant удаляет вариант (soft delete - помечает как неактивный)
func (s *Storage) DeleteListingVariant(ctx context.Context, variantID int) error {
	// Soft delete - просто помечаем как неактивный
	query := `UPDATE c2c_listing_variants SET is_active = false WHERE id = $1`

	result, err := s.pool.Exec(ctx, query, variantID)
	if err != nil {
		return fmt.Errorf("failed to delete variant: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("variant not found")
	}

	return nil
}
