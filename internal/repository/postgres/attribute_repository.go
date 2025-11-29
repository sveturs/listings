// Package postgres implements PostgreSQL repository layer for attributes system.
package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"

	"github.com/vondi-global/listings/internal/domain"
)

// AttributeRepository implements PostgreSQL data access for attributes
type AttributeRepository struct {
	db     *sqlx.DB
	logger zerolog.Logger
}

// NewAttributeRepository creates a new PostgreSQL attribute repository
func NewAttributeRepository(db *sqlx.DB, logger zerolog.Logger) *AttributeRepository {
	return &AttributeRepository{
		db:     db,
		logger: logger.With().Str("component", "attribute_repository").Logger(),
	}
}

// Create creates a new attribute
func (r *AttributeRepository) Create(ctx context.Context, input *domain.CreateAttributeInput) (*domain.Attribute, error) {
	// Validate input
	if input == nil {
		return nil, fmt.Errorf("input cannot be nil")
	}
	if input.Code == "" {
		return nil, fmt.Errorf("code is required")
	}
	if len(input.Name) == 0 {
		return nil, fmt.Errorf("name is required")
	}
	if input.AttributeType == "" {
		return nil, fmt.Errorf("attribute_type is required")
	}

	// Marshal JSONB fields
	nameJSON, err := json.Marshal(input.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal name: %w", err)
	}

	displayNameJSON, err := json.Marshal(input.DisplayName)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal display_name: %w", err)
	}

	optionsJSON, err := json.Marshal(input.Options)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal options: %w", err)
	}

	validationRulesJSON, err := json.Marshal(input.ValidationRules)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal validation_rules: %w", err)
	}

	uiSettingsJSON, err := json.Marshal(input.UISettings)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal ui_settings: %w", err)
	}

	// Default purpose to 'regular' if not provided
	purpose := input.Purpose
	if purpose == "" {
		purpose = domain.AttributePurposeRegular
	}

	query := `
		INSERT INTO attributes (
			code, name, display_name, attribute_type, purpose,
			options, validation_rules, ui_settings,
			is_searchable, is_filterable, is_required, is_variant_compatible,
			affects_stock, affects_price, show_in_card,
			sort_order, icon
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
		RETURNING id, code, name, display_name, attribute_type, purpose,
		          options, validation_rules, ui_settings,
		          is_searchable, is_filterable, is_required, is_variant_compatible,
		          affects_stock, affects_price, show_in_card,
		          is_active, sort_order, icon, created_at, updated_at
	`

	var attr domain.Attribute
	var nameBytes, displayNameBytes, optionsBytes, validationRulesBytes, uiSettingsBytes []byte

	err = r.db.QueryRowContext(
		ctx,
		query,
		input.Code,
		nameJSON,
		displayNameJSON,
		input.AttributeType,
		purpose,
		optionsJSON,
		validationRulesJSON,
		uiSettingsJSON,
		input.IsSearchable,
		input.IsFilterable,
		input.IsRequired,
		input.IsVariantCompatible,
		input.AffectsStock,
		input.AffectsPrice,
		input.ShowInCard,
		input.SortOrder,
		input.Icon,
	).Scan(
		&attr.ID,
		&attr.Code,
		&nameBytes,
		&displayNameBytes,
		&attr.AttributeType,
		&attr.Purpose,
		&optionsBytes,
		&validationRulesBytes,
		&uiSettingsBytes,
		&attr.IsSearchable,
		&attr.IsFilterable,
		&attr.IsRequired,
		&attr.IsVariantCompatible,
		&attr.AffectsStock,
		&attr.AffectsPrice,
		&attr.ShowInCard,
		&attr.IsActive,
		&attr.SortOrder,
		&attr.Icon,
		&attr.CreatedAt,
		&attr.UpdatedAt,
	)

	if err != nil {
		r.logger.Error().Err(err).Str("code", input.Code).Msg("failed to create attribute")
		return nil, fmt.Errorf("failed to create attribute: %w", err)
	}

	// Unmarshal JSONB fields
	if err := r.unmarshalAttributeJSONB(&attr, nameBytes, displayNameBytes, optionsBytes, validationRulesBytes, uiSettingsBytes); err != nil {
		return nil, err
	}

	r.logger.Info().Int32("id", attr.ID).Str("code", attr.Code).Msg("attribute created")
	return &attr, nil
}

// Update updates an existing attribute
func (r *AttributeRepository) Update(ctx context.Context, id int32, input *domain.UpdateAttributeInput) (*domain.Attribute, error) {
	if input == nil {
		return nil, fmt.Errorf("input cannot be nil")
	}

	// Build dynamic update query
	updates := []string{}
	args := []interface{}{}
	argPos := 1

	if input.Name != nil {
		nameJSON, err := json.Marshal(*input.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal name: %w", err)
		}
		updates = append(updates, fmt.Sprintf("name = $%d", argPos))
		args = append(args, nameJSON)
		argPos++
	}

	if input.DisplayName != nil {
		displayNameJSON, err := json.Marshal(*input.DisplayName)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal display_name: %w", err)
		}
		updates = append(updates, fmt.Sprintf("display_name = $%d", argPos))
		args = append(args, displayNameJSON)
		argPos++
	}

	if input.AttributeType != nil {
		updates = append(updates, fmt.Sprintf("attribute_type = $%d", argPos))
		args = append(args, *input.AttributeType)
		argPos++
	}

	if input.Purpose != nil {
		updates = append(updates, fmt.Sprintf("purpose = $%d", argPos))
		args = append(args, *input.Purpose)
		argPos++
	}

	if input.Options != nil {
		optionsJSON, err := json.Marshal(*input.Options)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal options: %w", err)
		}
		updates = append(updates, fmt.Sprintf("options = $%d", argPos))
		args = append(args, optionsJSON)
		argPos++
	}

	if input.ValidationRules != nil {
		validationRulesJSON, err := json.Marshal(*input.ValidationRules)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal validation_rules: %w", err)
		}
		updates = append(updates, fmt.Sprintf("validation_rules = $%d", argPos))
		args = append(args, validationRulesJSON)
		argPos++
	}

	if input.UISettings != nil {
		uiSettingsJSON, err := json.Marshal(*input.UISettings)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal ui_settings: %w", err)
		}
		updates = append(updates, fmt.Sprintf("ui_settings = $%d", argPos))
		args = append(args, uiSettingsJSON)
		argPos++
	}

	if input.IsSearchable != nil {
		updates = append(updates, fmt.Sprintf("is_searchable = $%d", argPos))
		args = append(args, *input.IsSearchable)
		argPos++
	}

	if input.IsFilterable != nil {
		updates = append(updates, fmt.Sprintf("is_filterable = $%d", argPos))
		args = append(args, *input.IsFilterable)
		argPos++
	}

	if input.IsRequired != nil {
		updates = append(updates, fmt.Sprintf("is_required = $%d", argPos))
		args = append(args, *input.IsRequired)
		argPos++
	}

	if input.IsVariantCompatible != nil {
		updates = append(updates, fmt.Sprintf("is_variant_compatible = $%d", argPos))
		args = append(args, *input.IsVariantCompatible)
		argPos++
	}

	if input.AffectsStock != nil {
		updates = append(updates, fmt.Sprintf("affects_stock = $%d", argPos))
		args = append(args, *input.AffectsStock)
		argPos++
	}

	if input.AffectsPrice != nil {
		updates = append(updates, fmt.Sprintf("affects_price = $%d", argPos))
		args = append(args, *input.AffectsPrice)
		argPos++
	}

	if input.ShowInCard != nil {
		updates = append(updates, fmt.Sprintf("show_in_card = $%d", argPos))
		args = append(args, *input.ShowInCard)
		argPos++
	}

	if input.IsActive != nil {
		updates = append(updates, fmt.Sprintf("is_active = $%d", argPos))
		args = append(args, *input.IsActive)
		argPos++
	}

	if input.SortOrder != nil {
		updates = append(updates, fmt.Sprintf("sort_order = $%d", argPos))
		args = append(args, *input.SortOrder)
		argPos++
	}

	if input.Icon != nil {
		updates = append(updates, fmt.Sprintf("icon = $%d", argPos))
		args = append(args, *input.Icon)
		argPos++
	}

	if len(updates) == 0 {
		// Nothing to update, return current attribute
		return r.GetByID(ctx, id)
	}

	// Add updated_at
	updates = append(updates, "updated_at = CURRENT_TIMESTAMP")

	// Add ID as last parameter
	args = append(args, id)

	query := fmt.Sprintf(`
		UPDATE attributes
		SET %s
		WHERE id = $%d AND is_active = true
		RETURNING id, code, name, display_name, attribute_type, purpose,
		          options, validation_rules, ui_settings,
		          is_searchable, is_filterable, is_required, is_variant_compatible,
		          affects_stock, affects_price, show_in_card,
		          is_active, sort_order, icon, created_at, updated_at
	`, strings.Join(updates, ", "), argPos)

	var attr domain.Attribute
	var nameBytes, displayNameBytes, optionsBytes, validationRulesBytes, uiSettingsBytes []byte

	err := r.db.QueryRowContext(ctx, query, args...).Scan(
		&attr.ID,
		&attr.Code,
		&nameBytes,
		&displayNameBytes,
		&attr.AttributeType,
		&attr.Purpose,
		&optionsBytes,
		&validationRulesBytes,
		&uiSettingsBytes,
		&attr.IsSearchable,
		&attr.IsFilterable,
		&attr.IsRequired,
		&attr.IsVariantCompatible,
		&attr.AffectsStock,
		&attr.AffectsPrice,
		&attr.ShowInCard,
		&attr.IsActive,
		&attr.SortOrder,
		&attr.Icon,
		&attr.CreatedAt,
		&attr.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("attribute not found or already deleted")
		}
		r.logger.Error().Err(err).Int32("id", id).Msg("failed to update attribute")
		return nil, fmt.Errorf("failed to update attribute: %w", err)
	}

	// Unmarshal JSONB fields
	if err := r.unmarshalAttributeJSONB(&attr, nameBytes, displayNameBytes, optionsBytes, validationRulesBytes, uiSettingsBytes); err != nil {
		return nil, err
	}

	r.logger.Info().Int32("id", attr.ID).Msg("attribute updated")
	return &attr, nil
}

// Delete soft-deletes an attribute by setting is_active to false
func (r *AttributeRepository) Delete(ctx context.Context, id int32) error {
	query := `
		UPDATE attributes
		SET is_active = false, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND is_active = true
	`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		r.logger.Error().Err(err).Int32("id", id).Msg("failed to delete attribute")
		return fmt.Errorf("failed to delete attribute: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("attribute not found or already deleted")
	}

	r.logger.Info().Int32("id", id).Msg("attribute deleted (soft)")
	return nil
}

// GetByID retrieves an attribute by its ID
func (r *AttributeRepository) GetByID(ctx context.Context, id int32) (*domain.Attribute, error) {
	query := `
		SELECT id, code, name, display_name, attribute_type, purpose,
		       options, validation_rules, ui_settings,
		       is_searchable, is_filterable, is_required, is_variant_compatible,
		       affects_stock, affects_price, show_in_card,
		       is_active, sort_order, icon, created_at, updated_at
		FROM attributes
		WHERE id = $1 AND is_active = true
	`

	var attr domain.Attribute
	var nameBytes, displayNameBytes, optionsBytes, validationRulesBytes, uiSettingsBytes []byte

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&attr.ID,
		&attr.Code,
		&nameBytes,
		&displayNameBytes,
		&attr.AttributeType,
		&attr.Purpose,
		&optionsBytes,
		&validationRulesBytes,
		&uiSettingsBytes,
		&attr.IsSearchable,
		&attr.IsFilterable,
		&attr.IsRequired,
		&attr.IsVariantCompatible,
		&attr.AffectsStock,
		&attr.AffectsPrice,
		&attr.ShowInCard,
		&attr.IsActive,
		&attr.SortOrder,
		&attr.Icon,
		&attr.CreatedAt,
		&attr.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("attribute not found")
		}
		r.logger.Error().Err(err).Int32("id", id).Msg("failed to get attribute")
		return nil, fmt.Errorf("failed to get attribute: %w", err)
	}

	// Unmarshal JSONB fields
	if err := r.unmarshalAttributeJSONB(&attr, nameBytes, displayNameBytes, optionsBytes, validationRulesBytes, uiSettingsBytes); err != nil {
		return nil, err
	}

	return &attr, nil
}

// GetByCode retrieves an attribute by its code
func (r *AttributeRepository) GetByCode(ctx context.Context, code string) (*domain.Attribute, error) {
	query := `
		SELECT id, code, name, display_name, attribute_type, purpose,
		       options, validation_rules, ui_settings,
		       is_searchable, is_filterable, is_required, is_variant_compatible,
		       affects_stock, affects_price, show_in_card,
		       is_active, sort_order, icon, created_at, updated_at
		FROM attributes
		WHERE code = $1 AND is_active = true
	`

	var attr domain.Attribute
	var nameBytes, displayNameBytes, optionsBytes, validationRulesBytes, uiSettingsBytes []byte

	err := r.db.QueryRowContext(ctx, query, code).Scan(
		&attr.ID,
		&attr.Code,
		&nameBytes,
		&displayNameBytes,
		&attr.AttributeType,
		&attr.Purpose,
		&optionsBytes,
		&validationRulesBytes,
		&uiSettingsBytes,
		&attr.IsSearchable,
		&attr.IsFilterable,
		&attr.IsRequired,
		&attr.IsVariantCompatible,
		&attr.AffectsStock,
		&attr.AffectsPrice,
		&attr.ShowInCard,
		&attr.IsActive,
		&attr.SortOrder,
		&attr.Icon,
		&attr.CreatedAt,
		&attr.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("attribute not found")
		}
		r.logger.Error().Err(err).Str("code", code).Msg("failed to get attribute by code")
		return nil, fmt.Errorf("failed to get attribute: %w", err)
	}

	// Unmarshal JSONB fields
	if err := r.unmarshalAttributeJSONB(&attr, nameBytes, displayNameBytes, optionsBytes, validationRulesBytes, uiSettingsBytes); err != nil {
		return nil, err
	}

	return &attr, nil
}

// List retrieves a filtered list of attributes with pagination
func (r *AttributeRepository) List(ctx context.Context, filter *domain.ListAttributesFilter) ([]*domain.Attribute, int64, error) {
	if filter == nil {
		filter = &domain.ListAttributesFilter{
			Limit:  10,
			Offset: 0,
		}
	}

	// Build WHERE clause dynamically
	whereConditions := []string{"is_active = true"}
	args := []interface{}{}
	argPos := 1

	if filter.AttributeType != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("attribute_type = $%d", argPos))
		args = append(args, *filter.AttributeType)
		argPos++
	}

	if filter.Purpose != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("purpose = $%d", argPos))
		args = append(args, *filter.Purpose)
		argPos++
	}

	if filter.IsSearchable != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("is_searchable = $%d", argPos))
		args = append(args, *filter.IsSearchable)
		argPos++
	}

	if filter.IsFilterable != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("is_filterable = $%d", argPos))
		args = append(args, *filter.IsFilterable)
		argPos++
	}

	if filter.IsVariantCompatible != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("is_variant_compatible = $%d", argPos))
		args = append(args, *filter.IsVariantCompatible)
		argPos++
	}

	if filter.IsActive != nil {
		// Override default is_active = true if explicitly set
		whereConditions[0] = fmt.Sprintf("is_active = $%d", argPos)
		args = append(args, *filter.IsActive)
		argPos++
	}

	whereClause := strings.Join(whereConditions, " AND ")

	// Get total count
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM attributes WHERE %s", whereClause)
	var total int64
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		r.logger.Error().Err(err).Msg("failed to count attributes")
		return nil, 0, fmt.Errorf("failed to count attributes: %w", err)
	}

	// Get attributes with pagination
	args = append(args, filter.Limit, filter.Offset)
	query := fmt.Sprintf(`
		SELECT id, code, name, display_name, attribute_type, purpose,
		       options, validation_rules, ui_settings,
		       is_searchable, is_filterable, is_required, is_variant_compatible,
		       affects_stock, affects_price, show_in_card,
		       is_active, sort_order, icon, created_at, updated_at
		FROM attributes
		WHERE %s
		ORDER BY sort_order ASC, (name->>'en') ASC
		LIMIT $%d OFFSET $%d
	`, whereClause, argPos, argPos+1)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		r.logger.Error().Err(err).Msg("failed to list attributes")
		return nil, 0, fmt.Errorf("failed to list attributes: %w", err)
	}
	defer rows.Close()

	var attributes []*domain.Attribute
	for rows.Next() {
		var attr domain.Attribute
		var nameBytes, displayNameBytes, optionsBytes, validationRulesBytes, uiSettingsBytes []byte

		err := rows.Scan(
			&attr.ID,
			&attr.Code,
			&nameBytes,
			&displayNameBytes,
			&attr.AttributeType,
			&attr.Purpose,
			&optionsBytes,
			&validationRulesBytes,
			&uiSettingsBytes,
			&attr.IsSearchable,
			&attr.IsFilterable,
			&attr.IsRequired,
			&attr.IsVariantCompatible,
			&attr.AffectsStock,
			&attr.AffectsPrice,
			&attr.ShowInCard,
			&attr.IsActive,
			&attr.SortOrder,
			&attr.Icon,
			&attr.CreatedAt,
			&attr.UpdatedAt,
		)

		if err != nil {
			r.logger.Error().Err(err).Msg("failed to scan attribute row")
			return nil, 0, fmt.Errorf("failed to scan attribute: %w", err)
		}

		// Unmarshal JSONB fields
		if err := r.unmarshalAttributeJSONB(&attr, nameBytes, displayNameBytes, optionsBytes, validationRulesBytes, uiSettingsBytes); err != nil {
			return nil, 0, err
		}

		attributes = append(attributes, &attr)
	}

	if err := rows.Err(); err != nil {
		r.logger.Error().Err(err).Msg("error iterating attribute rows")
		return nil, 0, fmt.Errorf("failed to iterate attributes: %w", err)
	}

	return attributes, total, nil
}

// unmarshalAttributeJSONB unmarshals JSONB fields for an attribute
func (r *AttributeRepository) unmarshalAttributeJSONB(attr *domain.Attribute, nameBytes, displayNameBytes, optionsBytes, validationRulesBytes, uiSettingsBytes []byte) error {
	if len(nameBytes) > 0 {
		if err := json.Unmarshal(nameBytes, &attr.Name); err != nil {
			r.logger.Error().Err(err).Msg("failed to unmarshal name")
			return fmt.Errorf("failed to unmarshal name: %w", err)
		}
	}

	if len(displayNameBytes) > 0 {
		if err := json.Unmarshal(displayNameBytes, &attr.DisplayName); err != nil {
			r.logger.Error().Err(err).Msg("failed to unmarshal display_name")
			return fmt.Errorf("failed to unmarshal display_name: %w", err)
		}
	}

	if len(optionsBytes) > 0 && string(optionsBytes) != "null" {
		if err := json.Unmarshal(optionsBytes, &attr.Options); err != nil {
			r.logger.Error().Err(err).Msg("failed to unmarshal options")
			return fmt.Errorf("failed to unmarshal options: %w", err)
		}
	}

	if len(validationRulesBytes) > 0 && string(validationRulesBytes) != "null" {
		if err := json.Unmarshal(validationRulesBytes, &attr.ValidationRules); err != nil {
			r.logger.Error().Err(err).Msg("failed to unmarshal validation_rules")
			return fmt.Errorf("failed to unmarshal validation_rules: %w", err)
		}
	}

	if len(uiSettingsBytes) > 0 && string(uiSettingsBytes) != "null" {
		if err := json.Unmarshal(uiSettingsBytes, &attr.UISettings); err != nil {
			r.logger.Error().Err(err).Msg("failed to unmarshal ui_settings")
			return fmt.Errorf("failed to unmarshal ui_settings: %w", err)
		}
	}

	return nil
}

// =============================================================================
// Category Linking Methods
// =============================================================================

// LinkToCategory links an attribute to a category with specific settings
func (r *AttributeRepository) LinkToCategory(ctx context.Context, categoryID int32, attributeID int32, settings *domain.CategoryAttributeSettings) (*domain.CategoryAttribute, error) {
	if settings == nil {
		return nil, fmt.Errorf("settings cannot be nil")
	}

	// Marshal JSONB fields - use interface{} to properly handle NULL
	var categorySpecificOptionsJSON interface{} = nil
	var err error
	if settings.CategorySpecificOptions != nil {
		bytes, err := json.Marshal(*settings.CategorySpecificOptions)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal category_specific_options: %w", err)
		}
		categorySpecificOptionsJSON = bytes
	}

	var customValidationRulesJSON interface{} = nil
	if settings.CustomValidationRules != nil {
		bytes, err := json.Marshal(*settings.CustomValidationRules)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal custom_validation_rules: %w", err)
		}
		customValidationRulesJSON = bytes
	}

	var customUISettingsJSON interface{} = nil
	if settings.CustomUISettings != nil {
		bytes, err := json.Marshal(*settings.CustomUISettings)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal custom_ui_settings: %w", err)
		}
		customUISettingsJSON = bytes
	}

	query := `
		INSERT INTO category_attributes (
			category_id, attribute_id, is_enabled, is_required, is_searchable, is_filterable,
			sort_order, category_specific_options, custom_validation_rules, custom_ui_settings
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		ON CONFLICT (category_id, attribute_id)
		DO UPDATE SET
			is_enabled = EXCLUDED.is_enabled,
			is_required = EXCLUDED.is_required,
			is_searchable = EXCLUDED.is_searchable,
			is_filterable = EXCLUDED.is_filterable,
			sort_order = EXCLUDED.sort_order,
			category_specific_options = EXCLUDED.category_specific_options,
			custom_validation_rules = EXCLUDED.custom_validation_rules,
			custom_ui_settings = EXCLUDED.custom_ui_settings,
			updated_at = CURRENT_TIMESTAMP
		RETURNING id, category_id, attribute_id, is_enabled, is_required, is_searchable, is_filterable,
		          sort_order, category_specific_options, custom_validation_rules, custom_ui_settings,
		          is_active, created_at, updated_at
	`

	var catAttr domain.CategoryAttribute
	var categorySpecificOptionsBytes, customValidationRulesBytes, customUISettingsBytes []byte

	err = r.db.QueryRowContext(
		ctx,
		query,
		categoryID,
		attributeID,
		settings.IsEnabled,
		settings.IsRequired,
		settings.IsSearchable,
		settings.IsFilterable,
		settings.SortOrder,
		categorySpecificOptionsJSON,
		customValidationRulesJSON,
		customUISettingsJSON,
	).Scan(
		&catAttr.ID,
		&catAttr.CategoryID,
		&catAttr.AttributeID,
		&catAttr.IsEnabled,
		&catAttr.IsRequired,
		&catAttr.IsSearchable,
		&catAttr.IsFilterable,
		&catAttr.SortOrder,
		&categorySpecificOptionsBytes,
		&customValidationRulesBytes,
		&customUISettingsBytes,
		&catAttr.IsActive,
		&catAttr.CreatedAt,
		&catAttr.UpdatedAt,
	)

	if err != nil {
		r.logger.Error().Err(err).Int32("category_id", categoryID).Int32("attribute_id", attributeID).Msg("failed to link attribute to category")
		return nil, fmt.Errorf("failed to link attribute to category: %w", err)
	}

	// Unmarshal JSONB fields
	if err := r.unmarshalCategoryAttributeJSONB(&catAttr, categorySpecificOptionsBytes, customValidationRulesBytes, customUISettingsBytes); err != nil {
		return nil, err
	}

	r.logger.Info().Int32("id", catAttr.ID).Int32("category_id", categoryID).Int32("attribute_id", attributeID).Msg("attribute linked to category")
	return &catAttr, nil
}

// UpdateCategoryAttribute updates category-specific attribute settings
func (r *AttributeRepository) UpdateCategoryAttribute(ctx context.Context, catAttrID int32, settings *domain.CategoryAttributeSettings) (*domain.CategoryAttribute, error) {
	if settings == nil {
		return nil, fmt.Errorf("settings cannot be nil")
	}

	// Marshal JSONB fields - use interface{} to properly handle NULL
	var categorySpecificOptionsJSON interface{} = nil
	var err error
	if settings.CategorySpecificOptions != nil {
		bytes, err := json.Marshal(*settings.CategorySpecificOptions)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal category_specific_options: %w", err)
		}
		categorySpecificOptionsJSON = bytes
	}

	var customValidationRulesJSON interface{} = nil
	if settings.CustomValidationRules != nil {
		bytes, err := json.Marshal(*settings.CustomValidationRules)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal custom_validation_rules: %w", err)
		}
		customValidationRulesJSON = bytes
	}

	var customUISettingsJSON interface{} = nil
	if settings.CustomUISettings != nil {
		bytes, err := json.Marshal(*settings.CustomUISettings)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal custom_ui_settings: %w", err)
		}
		customUISettingsJSON = bytes
	}

	query := `
		UPDATE category_attributes
		SET is_enabled = $2,
		    is_required = $3,
		    is_searchable = $4,
		    is_filterable = $5,
		    sort_order = $6,
		    category_specific_options = $7,
		    custom_validation_rules = $8,
		    custom_ui_settings = $9,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND is_active = true
		RETURNING id, category_id, attribute_id, is_enabled, is_required, is_searchable, is_filterable,
		          sort_order, category_specific_options, custom_validation_rules, custom_ui_settings,
		          is_active, created_at, updated_at
	`

	var catAttr domain.CategoryAttribute
	var categorySpecificOptionsBytes, customValidationRulesBytes, customUISettingsBytes []byte

	err = r.db.QueryRowContext(
		ctx,
		query,
		catAttrID,
		settings.IsEnabled,
		settings.IsRequired,
		settings.IsSearchable,
		settings.IsFilterable,
		settings.SortOrder,
		categorySpecificOptionsJSON,
		customValidationRulesJSON,
		customUISettingsJSON,
	).Scan(
		&catAttr.ID,
		&catAttr.CategoryID,
		&catAttr.AttributeID,
		&catAttr.IsEnabled,
		&catAttr.IsRequired,
		&catAttr.IsSearchable,
		&catAttr.IsFilterable,
		&catAttr.SortOrder,
		&categorySpecificOptionsBytes,
		&customValidationRulesBytes,
		&customUISettingsBytes,
		&catAttr.IsActive,
		&catAttr.CreatedAt,
		&catAttr.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("category attribute not found or already deleted")
		}
		r.logger.Error().Err(err).Int32("id", catAttrID).Msg("failed to update category attribute")
		return nil, fmt.Errorf("failed to update category attribute: %w", err)
	}

	// Unmarshal JSONB fields
	if err := r.unmarshalCategoryAttributeJSONB(&catAttr, categorySpecificOptionsBytes, customValidationRulesBytes, customUISettingsBytes); err != nil {
		return nil, err
	}

	r.logger.Info().Int32("id", catAttr.ID).Msg("category attribute updated")
	return &catAttr, nil
}

// UnlinkFromCategory removes attribute-category link
func (r *AttributeRepository) UnlinkFromCategory(ctx context.Context, categoryID int32, attributeID int32) error {
	query := `
		UPDATE category_attributes
		SET is_active = false, updated_at = CURRENT_TIMESTAMP
		WHERE category_id = $1 AND attribute_id = $2 AND is_active = true
	`

	result, err := r.db.ExecContext(ctx, query, categoryID, attributeID)
	if err != nil {
		r.logger.Error().Err(err).Int32("category_id", categoryID).Int32("attribute_id", attributeID).Msg("failed to unlink attribute from category")
		return fmt.Errorf("failed to unlink attribute from category: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("category attribute link not found or already deleted")
	}

	r.logger.Info().Int32("category_id", categoryID).Int32("attribute_id", attributeID).Msg("attribute unlinked from category")
	return nil
}

// GetCategoryAttributes retrieves attributes for a specific category with filters
func (r *AttributeRepository) GetCategoryAttributes(ctx context.Context, categoryID int32, filter *domain.GetCategoryAttributesFilter) ([]*domain.CategoryAttribute, error) {
	// Build WHERE clause dynamically
	whereConditions := []string{"ca.category_id = $1", "ca.is_active = true", "a.is_active = true"}
	args := []interface{}{categoryID}
	argPos := 2

	if filter != nil {
		if filter.IsEnabled != nil {
			whereConditions = append(whereConditions, fmt.Sprintf("ca.is_enabled = $%d", argPos))
			args = append(args, *filter.IsEnabled)
			argPos++
		}

		if filter.IsRequired != nil {
			whereConditions = append(whereConditions, fmt.Sprintf("COALESCE(ca.is_required, a.is_required) = $%d", argPos))
			args = append(args, *filter.IsRequired)
			argPos++
		}

		if filter.IsSearchable != nil {
			whereConditions = append(whereConditions, fmt.Sprintf("COALESCE(ca.is_searchable, a.is_searchable) = $%d", argPos))
			args = append(args, *filter.IsSearchable)
			argPos++
		}

		if filter.IsFilterable != nil {
			whereConditions = append(whereConditions, fmt.Sprintf("COALESCE(ca.is_filterable, a.is_filterable) = $%d", argPos))
			args = append(args, *filter.IsFilterable)
		}
	}

	whereClause := strings.Join(whereConditions, " AND ")

	query := fmt.Sprintf(`
		SELECT ca.id, ca.category_id, ca.attribute_id, ca.is_enabled, ca.is_required, ca.is_searchable, ca.is_filterable,
		       ca.sort_order, ca.category_specific_options, ca.custom_validation_rules, ca.custom_ui_settings,
		       ca.is_active, ca.created_at, ca.updated_at,
		       a.id, a.code, a.name, a.display_name, a.attribute_type, a.purpose,
		       a.options, a.validation_rules, a.ui_settings,
		       a.is_searchable, a.is_filterable, a.is_required, a.is_variant_compatible,
		       a.affects_stock, a.affects_price, a.show_in_card,
		       a.is_active, a.sort_order, a.icon, a.created_at, a.updated_at
		FROM category_attributes ca
		INNER JOIN attributes a ON ca.attribute_id = a.id
		WHERE %s
		ORDER BY ca.sort_order ASC, (a.name->>'en') ASC
	`, whereClause)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		r.logger.Error().Err(err).Int32("category_id", categoryID).Msg("failed to get category attributes")
		return nil, fmt.Errorf("failed to get category attributes: %w", err)
	}
	defer rows.Close()

	var categoryAttributes []*domain.CategoryAttribute
	for rows.Next() {
		var catAttr domain.CategoryAttribute
		var attr domain.Attribute
		var categorySpecificOptionsBytes, customValidationRulesBytes, customUISettingsBytes []byte
		var attrNameBytes, attrDisplayNameBytes, attrOptionsBytes, attrValidationRulesBytes, attrUISettingsBytes []byte

		err := rows.Scan(
			&catAttr.ID,
			&catAttr.CategoryID,
			&catAttr.AttributeID,
			&catAttr.IsEnabled,
			&catAttr.IsRequired,
			&catAttr.IsSearchable,
			&catAttr.IsFilterable,
			&catAttr.SortOrder,
			&categorySpecificOptionsBytes,
			&customValidationRulesBytes,
			&customUISettingsBytes,
			&catAttr.IsActive,
			&catAttr.CreatedAt,
			&catAttr.UpdatedAt,
			&attr.ID,
			&attr.Code,
			&attrNameBytes,
			&attrDisplayNameBytes,
			&attr.AttributeType,
			&attr.Purpose,
			&attrOptionsBytes,
			&attrValidationRulesBytes,
			&attrUISettingsBytes,
			&attr.IsSearchable,
			&attr.IsFilterable,
			&attr.IsRequired,
			&attr.IsVariantCompatible,
			&attr.AffectsStock,
			&attr.AffectsPrice,
			&attr.ShowInCard,
			&attr.IsActive,
			&attr.SortOrder,
			&attr.Icon,
			&attr.CreatedAt,
			&attr.UpdatedAt,
		)

		if err != nil {
			r.logger.Error().Err(err).Msg("failed to scan category attribute row")
			return nil, fmt.Errorf("failed to scan category attribute: %w", err)
		}

		// Unmarshal category attribute JSONB fields
		if err := r.unmarshalCategoryAttributeJSONB(&catAttr, categorySpecificOptionsBytes, customValidationRulesBytes, customUISettingsBytes); err != nil {
			return nil, err
		}

		// Unmarshal attribute JSONB fields
		if err := r.unmarshalAttributeJSONB(&attr, attrNameBytes, attrDisplayNameBytes, attrOptionsBytes, attrValidationRulesBytes, attrUISettingsBytes); err != nil {
			return nil, err
		}

		catAttr.Attribute = &attr
		categoryAttributes = append(categoryAttributes, &catAttr)
	}

	if err := rows.Err(); err != nil {
		r.logger.Error().Err(err).Msg("error iterating category attribute rows")
		return nil, fmt.Errorf("failed to iterate category attributes: %w", err)
	}

	return categoryAttributes, nil
}

// unmarshalCategoryAttributeJSONB unmarshals JSONB fields for a category attribute
func (r *AttributeRepository) unmarshalCategoryAttributeJSONB(catAttr *domain.CategoryAttribute, categorySpecificOptionsBytes, customValidationRulesBytes, customUISettingsBytes []byte) error {
	if len(categorySpecificOptionsBytes) > 0 && string(categorySpecificOptionsBytes) != "null" {
		if err := json.Unmarshal(categorySpecificOptionsBytes, &catAttr.CategorySpecificOptions); err != nil {
			r.logger.Error().Err(err).Msg("failed to unmarshal category_specific_options")
			return fmt.Errorf("failed to unmarshal category_specific_options: %w", err)
		}
	}

	if len(customValidationRulesBytes) > 0 && string(customValidationRulesBytes) != "null" {
		if err := json.Unmarshal(customValidationRulesBytes, &catAttr.CustomValidationRules); err != nil {
			r.logger.Error().Err(err).Msg("failed to unmarshal custom_validation_rules")
			return fmt.Errorf("failed to unmarshal custom_validation_rules: %w", err)
		}
	}

	if len(customUISettingsBytes) > 0 && string(customUISettingsBytes) != "null" {
		if err := json.Unmarshal(customUISettingsBytes, &catAttr.CustomUISettings); err != nil {
			r.logger.Error().Err(err).Msg("failed to unmarshal custom_ui_settings")
			return fmt.Errorf("failed to unmarshal custom_ui_settings: %w", err)
		}
	}

	return nil
}
