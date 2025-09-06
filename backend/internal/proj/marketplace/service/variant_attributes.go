package service

import (
	"context"
	"fmt"

	"backend/internal/domain/models"

	"github.com/lib/pq"
)

// categoryVariantAttributesMap определяет какие вариативные атрибуты доступны для каждой категории
var categoryVariantAttributesMap = map[string][]string{
	// Одежда
	"womens-clothing": {"color", "size", "material", "pattern", "style"},
	"mens-clothing":   {"color", "size", "material", "pattern", "style"},
	"kids-clothing":   {"color", "size", "material", "pattern"},
	"sports-clothing": {"color", "size", "material"},

	// Обувь
	"shoes": {"color", "size", "material", "style"},

	// Аксессуары
	"bags":        {"color", "size", "material", "style", "pattern"},
	"accessories": {"color", "size", "material", "style", "pattern"},

	// Электроника
	"smartphones":             {"color", "memory", "storage"},
	"computers":               {"color", "memory", "storage", "connectivity"},
	"gaming-consoles":         {"color", "storage", "bundle"},
	"electronics-accessories": {"color", "connectivity", "bundle"},

	// Бытовая техника
	"home-appliances": {"color", "capacity", "power"},

	// Мебель
	"furniture": {"color", "material", "style"},

	// Кухонная утварь
	"kitchenware": {"color", "capacity", "material"},
}

// GetProductVariantAttributes возвращает все доступные вариативные атрибуты
func (s *MarketplaceService) GetProductVariantAttributes(ctx context.Context) ([]*models.ProductVariantAttribute, error) {
	query := `
		SELECT 
			id, name, display_name, type, is_required, 
			sort_order, affects_stock, created_at, updated_at
		FROM product_variant_attributes
		ORDER BY sort_order, id
	`

	rows, err := s.storage.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get product variant attributes: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var attributes []*models.ProductVariantAttribute
	for rows.Next() {
		var attr models.ProductVariantAttribute
		err := rows.Scan(
			&attr.ID, &attr.Name, &attr.DisplayName, &attr.Type, &attr.IsRequired,
			&attr.SortOrder, &attr.AffectsStock, &attr.CreatedAt, &attr.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan product variant attribute: %w", err)
		}
		attributes = append(attributes, &attr)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate product variant attributes: %w", err)
	}

	return attributes, nil
}

// GetCategoryVariantAttributes возвращает вариативные атрибуты для конкретной категории
func (s *MarketplaceService) GetCategoryVariantAttributes(ctx context.Context, categorySlug string) ([]*models.ProductVariantAttribute, error) {
	// Получаем список имен атрибутов для категории
	attributeNames, exists := categoryVariantAttributesMap[categorySlug]
	if !exists {
		// Если категория не найдена в маппинге, возвращаем пустой список
		return []*models.ProductVariantAttribute{}, nil
	}

	// Получаем полную информацию об атрибутах из БД
	query := `
		SELECT 
			id, name, display_name, type, is_required, 
			sort_order, affects_stock, created_at, updated_at
		FROM product_variant_attributes
		WHERE name = ANY($1)
		ORDER BY sort_order, id
	`

	rows, err := s.storage.Query(ctx, query, pq.Array(attributeNames))
	if err != nil {
		return nil, fmt.Errorf("failed to get category variant attributes: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var attributes []*models.ProductVariantAttribute
	for rows.Next() {
		var attr models.ProductVariantAttribute
		err := rows.Scan(
			&attr.ID, &attr.Name, &attr.DisplayName, &attr.Type, &attr.IsRequired,
			&attr.SortOrder, &attr.AffectsStock, &attr.CreatedAt, &attr.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan category variant attribute: %w", err)
		}
		attributes = append(attributes, &attr)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate category variant attributes: %w", err)
	}

	return attributes, nil
}
