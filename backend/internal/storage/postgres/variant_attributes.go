package postgres

import (
	"context"
	"fmt"

	"backend/internal/domain/models"
)

// VariantAttributeStorage интерфейс для работы с вариативными атрибутами
type VariantAttributeStorage interface {
	GetCategoryVariantMappings(ctx context.Context, categoryID int) ([]*models.VariantAttributeMapping, error)
	CreateVariantMapping(ctx context.Context, mapping *models.VariantAttributeMappingCreateRequest) (*models.VariantAttributeMapping, error)
	UpdateVariantMapping(ctx context.Context, id int, update *models.VariantAttributeMappingUpdateRequest) error
	DeleteVariantMapping(ctx context.Context, id int) error
	DeleteCategoryVariantMappings(ctx context.Context, categoryID int) error
}

// GetVariantCompatibleAttributes получает атрибуты, которые могут быть вариантами
func (s *Storage) GetVariantCompatibleAttributes(ctx context.Context) ([]*models.UnifiedAttribute, error) {
	query := `
		SELECT 
			id, code, name, display_name, attribute_type, purpose,
			options, validation_rules, ui_settings,
			is_searchable, is_filterable, is_required,
			is_variant_compatible, affects_stock, affects_price,
			sort_order, is_active, created_at, updated_at
		FROM unified_attributes
		WHERE is_variant_compatible = true AND is_active = true
		ORDER BY sort_order, name`

	var attributes []*models.UnifiedAttribute
	rows, err := s.GetPool().Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query variant compatible attributes: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var attr models.UnifiedAttribute
		err := rows.Scan(
			&attr.ID, &attr.Code, &attr.Name, &attr.DisplayName,
			&attr.AttributeType, &attr.Purpose, &attr.Options,
			&attr.ValidationRules, &attr.UISettings,
			&attr.IsSearchable, &attr.IsFilterable, &attr.IsRequired,
			&attr.IsVariantCompatible, &attr.AffectsStock, &attr.AffectsPrice,
			&attr.SortOrder, &attr.IsActive, &attr.CreatedAt, &attr.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan attribute: %w", err)
		}
		attributes = append(attributes, &attr)
	}

	return attributes, nil
}

// GetCategoryVariantMappings получает вариативные атрибуты для категории
func (s *Storage) GetCategoryVariantMappings(ctx context.Context, categoryID int) ([]*models.VariantAttributeMapping, error) {
	query := `
		SELECT 
			vam.id, vam.variant_attribute_id, vam.category_id,
			vam.sort_order, vam.is_required, vam.created_at, vam.updated_at,
			ua.id as "attribute.id", ua.code as "attribute.code",
			ua.name as "attribute.name", ua.display_name as "attribute.display_name",
			ua.attribute_type as "attribute.attribute_type",
			ua.options as "attribute.options",
			ua.affects_stock as "attribute.affects_stock",
			ua.affects_price as "attribute.affects_price"
		FROM variant_attribute_mappings vam
		JOIN unified_attributes ua ON vam.variant_attribute_id = ua.id
		WHERE vam.category_id = $1
		ORDER BY vam.sort_order, ua.name`

	rows, err := s.GetPool().Query(ctx, query, categoryID)
	if err != nil {
		return nil, fmt.Errorf("failed to get category variant mappings: %w", err)
	}
	defer rows.Close()

	var mappings []*models.VariantAttributeMapping
	for rows.Next() {
		var m models.VariantAttributeMapping
		var attr models.UnifiedAttribute

		err := rows.Scan(
			&m.ID, &m.VariantAttributeID, &m.CategoryID,
			&m.SortOrder, &m.IsRequired, &m.CreatedAt, &m.UpdatedAt,
			&attr.ID, &attr.Code, &attr.Name, &attr.DisplayName,
			&attr.AttributeType, &attr.Options,
			&attr.AffectsStock, &attr.AffectsPrice,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan mapping: %w", err)
		}

		m.Attribute = &attr
		mappings = append(mappings, &m)
	}

	return mappings, nil
}

// CreateVariantMapping создает новую связь
func (s *Storage) CreateVariantMapping(ctx context.Context, mapping *models.VariantAttributeMappingCreateRequest) (*models.VariantAttributeMapping, error) {
	query := `
		INSERT INTO variant_attribute_mappings 
		(variant_attribute_id, category_id, sort_order, is_required)
		VALUES ($1, $2, $3, $4)
		RETURNING id, variant_attribute_id, category_id, sort_order, is_required, created_at, updated_at`

	var m models.VariantAttributeMapping
	err := s.GetPool().QueryRow(ctx, query,
		mapping.VariantAttributeID,
		mapping.CategoryID,
		mapping.SortOrder,
		mapping.IsRequired,
	).Scan(
		&m.ID, &m.VariantAttributeID, &m.CategoryID,
		&m.SortOrder, &m.IsRequired, &m.CreatedAt, &m.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create variant mapping: %w", err)
	}

	return &m, nil
}

// UpdateVariantMapping обновляет связь
func (s *Storage) UpdateVariantMapping(ctx context.Context, id int, update *models.VariantAttributeMappingUpdateRequest) error {
	query := `
		UPDATE variant_attribute_mappings 
		SET sort_order = COALESCE($2, sort_order),
		    is_required = COALESCE($3, is_required),
		    updated_at = NOW()
		WHERE id = $1`

	_, err := s.GetPool().Exec(ctx, query, id, update.SortOrder, update.IsRequired)
	if err != nil {
		return fmt.Errorf("failed to update variant mapping: %w", err)
	}

	return nil
}

// DeleteVariantMapping удаляет связь
func (s *Storage) DeleteVariantMapping(ctx context.Context, id int) error {
	query := `DELETE FROM variant_attribute_mappings WHERE id = $1`

	_, err := s.GetPool().Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete variant mapping: %w", err)
	}

	return nil
}

// DeleteCategoryVariantMappings удаляет все вариативные атрибуты категории
func (s *Storage) DeleteCategoryVariantMappings(ctx context.Context, categoryID int) error {
	query := `DELETE FROM variant_attribute_mappings WHERE category_id = $1`

	_, err := s.GetPool().Exec(ctx, query, categoryID)
	if err != nil {
		return fmt.Errorf("failed to delete category variant mappings: %w", err)
	}

	return nil
}
