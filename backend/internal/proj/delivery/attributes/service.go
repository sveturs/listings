package attributes

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/jmoiron/sqlx"

	"backend/internal/proj/delivery/models"
)

// Service - сервис управления атрибутами доставки
type Service struct {
	db *sqlx.DB
}

// NewService создает новый экземпляр сервиса атрибутов
func NewService(db *sqlx.DB) *Service {
	return &Service{
		db: db,
	}
}

// GetProductAttributes - получает атрибуты доставки товара
func (s *Service) GetProductAttributes(ctx context.Context, productID int, productType string) (*models.DeliveryAttributes, error) {
	var attrs models.DeliveryAttributes
	var query string
	var jsonData json.RawMessage

	if productType == "listing" {
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
			FROM marketplace_listings ml
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
			FROM storefront_products sp
			WHERE sp.id = $1`
	}

	if err := s.db.GetContext(ctx, &jsonData, query, productID); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product not found")
		}
		return nil, err
	}

	if err := json.Unmarshal(jsonData, &attrs); err != nil {
		return nil, fmt.Errorf("failed to unmarshal attributes: %w", err)
	}

	// Рассчитываем объем если не указан
	if attrs.VolumeM3 == 0 && attrs.Dimensions != nil {
		attrs.VolumeM3 = attrs.Dimensions.CalculateVolume()
	}

	// Устанавливаем дефолтные значения
	if attrs.PackagingType == "" {
		attrs.PackagingType = "box"
	}

	return &attrs, nil
}

// UpdateProductAttributes - обновляет атрибуты доставки товара
func (s *Service) UpdateProductAttributes(ctx context.Context, productID int, productType string, attrs *models.DeliveryAttributes) error {
	// Валидируем атрибуты
	if err := s.validateAttributes(attrs); err != nil {
		return fmt.Errorf("invalid attributes: %w", err)
	}

	// Рассчитываем объем
	if attrs.Dimensions != nil {
		attrs.VolumeM3 = attrs.Dimensions.CalculateVolume()
	}

	jsonData, err := json.Marshal(attrs)
	if err != nil {
		return fmt.Errorf("failed to marshal attributes: %w", err)
	}

	var query string
	if productType == "listing" {
		query = `
			UPDATE marketplace_listings
			SET metadata = jsonb_set(
				COALESCE(metadata, '{}'),
				'{delivery_attributes}',
				$1::jsonb
			)
			WHERE id = $2`
	} else {
		query = `
			UPDATE storefront_products
			SET attributes = jsonb_set(
				COALESCE(attributes, '{}'),
				'{delivery_attributes}',
				$1::jsonb
			)
			WHERE id = $2`
	}

	result, err := s.db.ExecContext(ctx, query, jsonData, productID)
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

// validateAttributes - валидирует атрибуты доставки
func (s *Service) validateAttributes(attrs *models.DeliveryAttributes) error {
	if attrs.WeightKg < 0 {
		return fmt.Errorf("weight cannot be negative")
	}

	if attrs.WeightKg > 1000 {
		return fmt.Errorf("weight exceeds maximum (1000 kg)")
	}

	if attrs.Dimensions != nil {
		if attrs.Dimensions.LengthCm < 0 || attrs.Dimensions.WidthCm < 0 || attrs.Dimensions.HeightCm < 0 {
			return fmt.Errorf("dimensions cannot be negative")
		}

		if attrs.Dimensions.LengthCm > 500 || attrs.Dimensions.WidthCm > 500 || attrs.Dimensions.HeightCm > 500 {
			return fmt.Errorf("dimensions exceed maximum (500 cm)")
		}
	}

	if attrs.MaxStackWeightKg < 0 {
		return fmt.Errorf("max stack weight cannot be negative")
	}

	validPackagingTypes := map[string]bool{
		"box":      true,
		"envelope": true,
		"pallet":   true,
		"custom":   true,
	}

	if attrs.PackagingType != "" && !validPackagingTypes[attrs.PackagingType] {
		return fmt.Errorf("invalid packaging type: %s", attrs.PackagingType)
	}

	return nil
}

// GetCategoryDefaults - получает дефолтные атрибуты для категории
func (s *Service) GetCategoryDefaults(ctx context.Context, categoryID int) (*models.CategoryDefaults, error) {
	var defaults models.CategoryDefaults
	query := `
		SELECT * FROM delivery_category_defaults
		WHERE category_id = $1`

	if err := s.db.GetContext(ctx, &defaults, query, categoryID); err != nil {
		if err == sql.ErrNoRows {
			// Возвращаем пустые дефолты если не найдены
			return &models.CategoryDefaults{
				CategoryID: categoryID,
			}, nil
		}
		return nil, err
	}

	return &defaults, nil
}

// UpdateCategoryDefaults - обновляет дефолтные атрибуты для категории
func (s *Service) UpdateCategoryDefaults(ctx context.Context, defaults *models.CategoryDefaults) error {
	// Валидируем дефолтные значения
	if defaults.DefaultWeightKg != nil && (*defaults.DefaultWeightKg < 0 || *defaults.DefaultWeightKg > 1000) {
		return fmt.Errorf("invalid default weight")
	}

	if defaults.DefaultLengthCm != nil && (*defaults.DefaultLengthCm < 0 || *defaults.DefaultLengthCm > 500) {
		return fmt.Errorf("invalid default length")
	}

	if defaults.DefaultWidthCm != nil && (*defaults.DefaultWidthCm < 0 || *defaults.DefaultWidthCm > 500) {
		return fmt.Errorf("invalid default width")
	}

	if defaults.DefaultHeightCm != nil && (*defaults.DefaultHeightCm < 0 || *defaults.DefaultHeightCm > 500) {
		return fmt.Errorf("invalid default height")
	}

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

// BatchUpdateProductAttributes - массовое обновление атрибутов товаров
func (s *Service) BatchUpdateProductAttributes(ctx context.Context, updates []ProductAttributesUpdate) error {
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	for _, update := range updates {
		if err := s.updateProductAttributesTx(ctx, tx, update.ProductID, update.ProductType, update.Attributes); err != nil {
			return fmt.Errorf("failed to update product %d: %w", update.ProductID, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// ProductAttributesUpdate - структура для массового обновления
type ProductAttributesUpdate struct {
	ProductID   int                        `json:"product_id"`
	ProductType string                     `json:"product_type"`
	Attributes  *models.DeliveryAttributes `json:"attributes"`
}

// updateProductAttributesTx - обновляет атрибуты в транзакции
func (s *Service) updateProductAttributesTx(ctx context.Context, tx *sqlx.Tx, productID int, productType string, attrs *models.DeliveryAttributes) error {
	if err := s.validateAttributes(attrs); err != nil {
		return err
	}

	if attrs.Dimensions != nil {
		attrs.VolumeM3 = attrs.Dimensions.CalculateVolume()
	}

	jsonData, err := json.Marshal(attrs)
	if err != nil {
		return err
	}

	var query string
	if productType == "listing" {
		query = `
			UPDATE marketplace_listings
			SET metadata = jsonb_set(
				COALESCE(metadata, '{}'),
				'{delivery_attributes}',
				$1::jsonb
			)
			WHERE id = $2`
	} else {
		query = `
			UPDATE storefront_products
			SET attributes = jsonb_set(
				COALESCE(attributes, '{}'),
				'{delivery_attributes}',
				$1::jsonb
			)
			WHERE id = $2`
	}

	_, err = tx.ExecContext(ctx, query, jsonData, productID)
	return err
}

// ApplyCategoryDefaultsToProducts - применяет дефолтные атрибуты категории к товарам без атрибутов
func (s *Service) ApplyCategoryDefaultsToProducts(ctx context.Context, categoryID int) (int, error) {
	// Получаем дефолтные атрибуты
	defaults, err := s.GetCategoryDefaults(ctx, categoryID)
	if err != nil {
		return 0, err
	}

	if defaults.DefaultWeightKg == nil {
		return 0, fmt.Errorf("no defaults configured for category %d", categoryID)
	}

	// Создаем JSON с дефолтными атрибутами
	defaultAttrs := models.DeliveryAttributes{
		WeightKg: *defaults.DefaultWeightKg,
		Dimensions: &models.Dimensions{
			LengthCm: *defaults.DefaultLengthCm,
			WidthCm:  *defaults.DefaultWidthCm,
			HeightCm: *defaults.DefaultHeightCm,
		},
		PackagingType: *defaults.DefaultPackagingType,
		IsFragile:     defaults.IsTypicallyFragile,
		Stackable:     true,
	}

	jsonData, err := json.Marshal(defaultAttrs)
	if err != nil {
		return 0, err
	}

	// Обновляем товары из marketplace_listings
	query1 := `
		UPDATE marketplace_listings
		SET metadata = jsonb_set(
			COALESCE(metadata, '{}'),
			'{delivery_attributes}',
			$1::jsonb
		)
		WHERE category_id = $2
		AND (metadata->'delivery_attributes' IS NULL OR metadata->'delivery_attributes' = '{}'::jsonb)`

	result1, err := s.db.ExecContext(ctx, query1, jsonData, categoryID)
	if err != nil {
		return 0, err
	}

	count1, _ := result1.RowsAffected()

	// Обновляем товары из storefront_products
	query2 := `
		UPDATE storefront_products
		SET attributes = jsonb_set(
			COALESCE(attributes, '{}'),
			'{delivery_attributes}',
			$1::jsonb
		)
		WHERE category_id = $2
		AND (attributes->'delivery_attributes' IS NULL OR attributes->'delivery_attributes' = '{}'::jsonb)`

	result2, err := s.db.ExecContext(ctx, query2, jsonData, categoryID)
	if err != nil {
		return 0, err
	}

	count2, _ := result2.RowsAffected()

	return int(count1 + count2), nil
}

// CalculateVolumetricWeight - рассчитывает объемный вес
func (s *Service) CalculateVolumetricWeight(dims *models.Dimensions, divisor float64) float64 {
	if dims == nil || divisor == 0 {
		return 0
	}
	return dims.CalculateVolumetricWeight(divisor)
}

// GetEffectiveWeight - возвращает эффективный вес (максимум из реального и объемного)
func (s *Service) GetEffectiveWeight(attrs *models.DeliveryAttributes, volumetricDivisor float64) float64 {
	if attrs == nil {
		return 0
	}

	realWeight := attrs.WeightKg
	volumetricWeight := s.CalculateVolumetricWeight(attrs.Dimensions, volumetricDivisor)

	if volumetricWeight > realWeight {
		return volumetricWeight
	}
	return realWeight
}