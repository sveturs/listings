package service

import (
	"context"
	"fmt"

	"backend/internal/domain/models"
)

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
	// Получаем ID категории по slug
	var categoryID int
	categoryQuery := `SELECT id FROM marketplace_categories WHERE slug = $1`
	err := s.storage.QueryRow(ctx, categoryQuery, categorySlug).Scan(&categoryID)
	if err != nil {
		return nil, fmt.Errorf("failed to get category by slug: %w", err)
	}

	// Получаем вариативные атрибуты из variant_attribute_mappings
	query := `
		SELECT
			ua.id, ua.name, ua.display_name, ua.attribute_type, vam.is_required,
			vam.sort_order, ua.affects_stock, ua.created_at, ua.updated_at
		FROM variant_attribute_mappings vam
		JOIN unified_attributes ua ON ua.id = vam.variant_attribute_id
		WHERE vam.category_id = $1
		ORDER BY vam.sort_order, ua.id
	`

	rows, err := s.storage.Query(ctx, query, categoryID)
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
