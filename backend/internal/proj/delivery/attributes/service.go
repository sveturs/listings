package attributes

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"

	"backend/internal/proj/delivery/models"
	"backend/internal/proj/delivery/storage"
)

const (
	ProductTypeListing = "listing"
)

// Service - сервис управления атрибутами доставки
type Service struct {
	db      *sqlx.DB
	storage *storage.Storage
}

// NewService создает новый экземпляр сервиса атрибутов
func NewService(db *sqlx.DB, storage *storage.Storage) *Service {
	return &Service{
		db:      db,
		storage: storage,
	}
}

// GetProductAttributes - получает атрибуты доставки товара
func (s *Service) GetProductAttributes(ctx context.Context, productID int, productType string) (*models.DeliveryAttributes, error) {
	var attrs models.DeliveryAttributes

	// Получаем JSON данные через storage
	jsonData, err := s.storage.GetProductAttributes(ctx, productID, productType)
	if err != nil {
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

	return s.storage.UpdateProductAttributes(ctx, productID, productType, jsonData)
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
	return s.storage.GetCategoryDefaults(ctx, categoryID)
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

	return s.storage.UpdateCategoryDefaults(ctx, defaults)
}

// BatchUpdateProductAttributes - массовое обновление атрибутов товаров
func (s *Service) BatchUpdateProductAttributes(ctx context.Context, updates []ProductAttributesUpdate) error {
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			log.Error().Err(rollbackErr).Msg("Failed to rollback transaction")
		}
	}()

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

	return s.storage.UpdateProductAttributesTx(ctx, tx, productID, productType, jsonData)
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

	// Обновляем товары из c2c_listings
	count1, err := s.storage.ApplyCategoryDefaultsToListings(ctx, categoryID, jsonData)
	if err != nil {
		return 0, err
	}

	// Обновляем товары из b2c_products
	count2, err := s.storage.ApplyCategoryDefaultsToProducts(ctx, categoryID, jsonData)
	if err != nil {
		return 0, err
	}

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
