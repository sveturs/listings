package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	"backend/internal/proj/delivery/models"
)

const (
	ProductTypeListing = "listing"
)

// GetProductAttributes получает атрибуты доставки товара
func (s *Storage) GetProductAttributes(ctx context.Context, productID int, productType string) (json.RawMessage, error) {
	var jsonData json.RawMessage
	var query string

	if productType == ProductTypeListing {
		query = `
			SELECT
				COALESCE(
					ml.metadata->'delivery_attributes',
					(SELECT jsonb_build_object(
						'weight_kg', dcd.default_weight_kg,
						'dimensions', jsonb_build_object(
							'length_cm', dcd.default_length_cm,
							'width_cm', dcd.default_width_cm,
							'height_cm', dcd.default_height_cm
						),
						'packaging_type', dcd.default_packaging_type,
						'is_fragile', dcd.is_typically_fragile,
						'stackable', true,
						'requires_special_handling', false
					)
					FROM delivery_category_defaults dcd
					WHERE dcd.category_id = ml.category_id),
					'{}'::jsonb
				) as attributes
			FROM c2c_listings ml
			WHERE ml.id = $1`
	} else {
		query = `
			SELECT
				COALESCE(
					sp.attributes->'delivery_attributes',
					(SELECT jsonb_build_object(
						'weight_kg', dcd.default_weight_kg,
						'dimensions', jsonb_build_object(
							'length_cm', dcd.default_length_cm,
							'width_cm', dcd.default_width_cm,
							'height_cm', dcd.default_height_cm
						),
						'packaging_type', dcd.default_packaging_type,
						'is_fragile', dcd.is_typically_fragile,
						'stackable', true,
						'requires_special_handling', false
					)
					FROM delivery_category_defaults dcd
					WHERE dcd.category_id = sp.category_id),
					'{}'::jsonb
				) as attributes
			FROM b2c_products sp
			WHERE sp.id = $1`
	}

	if err := s.db.GetContext(ctx, &jsonData, query, productID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("product not found")
		}
		return nil, err
	}

	return jsonData, nil
}

// UpdateProductAttributes обновляет атрибуты доставки товара
func (s *Storage) UpdateProductAttributes(ctx context.Context, productID int, productType string, attrs json.RawMessage) error {
	var query string

	if productType == ProductTypeListing {
		query = `
			UPDATE c2c_listings
			SET metadata = jsonb_set(
				COALESCE(metadata, '{}'),
				'{delivery_attributes}',
				$1::jsonb
			)
			WHERE id = $2`
	} else {
		query = `
			UPDATE b2c_products
			SET attributes = jsonb_set(
				COALESCE(attributes, '{}'),
				'{delivery_attributes}',
				$1::jsonb
			)
			WHERE id = $2`
	}

	result, err := s.db.ExecContext(ctx, query, attrs, productID)
	if err != nil {
		return fmt.Errorf("failed to update attributes: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("product not found")
	}

	return nil
}

// UpdateProductAttributesTx обновляет атрибуты доставки товара в рамках транзакции
func (s *Storage) UpdateProductAttributesTx(ctx context.Context, tx *sqlx.Tx, productID int, productType string, attrs json.RawMessage) error {
	var query string

	if productType == ProductTypeListing {
		query = `
			UPDATE c2c_listings
			SET metadata = jsonb_set(
				COALESCE(metadata, '{}'),
				'{delivery_attributes}',
				$1::jsonb
			)
			WHERE id = $2`
	} else {
		query = `
			UPDATE b2c_products
			SET attributes = jsonb_set(
				COALESCE(attributes, '{}'),
				'{delivery_attributes}',
				$1::jsonb
			)
			WHERE id = $2`
	}

	_, err := tx.ExecContext(ctx, query, attrs, productID)
	return err
}

// GetCategoryDefaults получает дефолтные атрибуты категории
func (s *Storage) GetCategoryDefaults(ctx context.Context, categoryID int) (*models.CategoryDefaults, error) {
	var defaults models.CategoryDefaults
	query := `
		SELECT * FROM delivery_category_defaults
		WHERE category_id = $1`

	if err := s.db.GetContext(ctx, &defaults, query, categoryID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Возвращаем пустые дефолты если не найдены
			return &models.CategoryDefaults{
				CategoryID: categoryID,
			}, nil
		}
		return nil, err
	}

	return &defaults, nil
}

// UpdateCategoryDefaults обновляет дефолтные атрибуты категории
func (s *Storage) UpdateCategoryDefaults(ctx context.Context, defaults *models.CategoryDefaults) error {
	query := `
		INSERT INTO delivery_category_defaults (
			category_id, default_weight_kg, default_length_cm,
			default_width_cm, default_height_cm, default_packaging_type,
			is_typically_fragile
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (category_id) DO UPDATE SET
			default_weight_kg = EXCLUDED.default_weight_kg,
			default_length_cm = EXCLUDED.default_length_cm,
			default_width_cm = EXCLUDED.default_width_cm,
			default_height_cm = EXCLUDED.default_height_cm,
			default_packaging_type = EXCLUDED.default_packaging_type,
			is_typically_fragile = EXCLUDED.is_typically_fragile,
			updated_at = NOW()
		RETURNING id`

	return s.db.GetContext(ctx, &defaults.ID, query,
		defaults.CategoryID,
		defaults.DefaultWeightKg,
		defaults.DefaultLengthCm,
		defaults.DefaultWidthCm,
		defaults.DefaultHeightCm,
		defaults.DefaultPackagingType,
		defaults.IsTypicallyFragile,
	)
}

// ApplyCategoryDefaultsToListings применяет дефолтные атрибуты к товарам C2C
func (s *Storage) ApplyCategoryDefaultsToListings(ctx context.Context, categoryID int, attrs json.RawMessage) (int64, error) {
	query := `
		UPDATE c2c_listings
		SET metadata = jsonb_set(
			COALESCE(metadata, '{}'),
			'{delivery_attributes}',
			$1::jsonb
		)
		WHERE category_id = $2
		AND (metadata->'delivery_attributes' IS NULL OR metadata->'delivery_attributes' = '{}'::jsonb)`

	result, err := s.db.ExecContext(ctx, query, attrs, categoryID)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

// ApplyCategoryDefaultsToProducts применяет дефолтные атрибуты к товарам B2C
func (s *Storage) ApplyCategoryDefaultsToProducts(ctx context.Context, categoryID int, attrs json.RawMessage) (int64, error) {
	query := `
		UPDATE b2c_products
		SET attributes = jsonb_set(
			COALESCE(attributes, '{}'),
			'{delivery_attributes}',
			$1::jsonb
		)
		WHERE category_id = $2
		AND (attributes->'delivery_attributes' IS NULL OR attributes->'delivery_attributes' = '{}'::jsonb)`

	result, err := s.db.ExecContext(ctx, query, attrs, categoryID)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}
