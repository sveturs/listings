// Package postgres implements PostgreSQL repository layer for attributes system.
// This file contains listing attribute values operations.
package postgres

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/vondi-global/listings/internal/domain"
)

// =============================================================================
// Listing Attribute Values Methods
// =============================================================================

// GetListingValues retrieves all attribute values for a listing
func (r *AttributeRepository) GetListingValues(ctx context.Context, listingID int32) ([]*domain.ListingAttributeValue, error) {
	query := `
		SELECT lav.id, lav.listing_id, lav.attribute_id,
		       lav.value_text, lav.value_number, lav.value_boolean, lav.value_date, lav.value_json,
		       lav.created_at, lav.updated_at,
		       a.id, a.code, a.name, a.display_name, a.attribute_type, a.purpose,
		       a.options, a.validation_rules, a.ui_settings,
		       a.is_searchable, a.is_filterable, a.is_required, a.is_variant_compatible,
		       a.affects_stock, a.affects_price, a.show_in_card,
		       a.is_active, a.sort_order, a.icon, a.created_at, a.updated_at
		FROM listing_attribute_values lav
		INNER JOIN attributes a ON lav.attribute_id = a.id
		WHERE lav.listing_id = $1 AND a.is_active = true
		ORDER BY a.sort_order ASC, (a.name->>'en') ASC
	`

	rows, err := r.db.QueryContext(ctx, query, listingID)
	if err != nil {
		r.logger.Error().Err(err).Int32("listing_id", listingID).Msg("failed to get listing attribute values")
		return nil, fmt.Errorf("failed to get listing attribute values: %w", err)
	}
	defer rows.Close()

	var listingValues []*domain.ListingAttributeValue
	for rows.Next() {
		var lav domain.ListingAttributeValue
		var attr domain.Attribute
		var valueJSONBytes []byte
		var attrNameBytes, attrDisplayNameBytes, attrOptionsBytes, attrValidationRulesBytes, attrUISettingsBytes []byte

		err := rows.Scan(
			&lav.ID,
			&lav.ListingID,
			&lav.AttributeID,
			&lav.ValueText,
			&lav.ValueNumber,
			&lav.ValueBoolean,
			&lav.ValueDate,
			&valueJSONBytes,
			&lav.CreatedAt,
			&lav.UpdatedAt,
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
			r.logger.Error().Err(err).Msg("failed to scan listing attribute value row")
			return nil, fmt.Errorf("failed to scan listing attribute value: %w", err)
		}

		// Unmarshal value_json
		if len(valueJSONBytes) > 0 && string(valueJSONBytes) != "null" {
			if err := json.Unmarshal(valueJSONBytes, &lav.ValueJSON); err != nil {
				r.logger.Error().Err(err).Msg("failed to unmarshal value_json")
				return nil, fmt.Errorf("failed to unmarshal value_json: %w", err)
			}
		}

		// Unmarshal attribute JSONB fields
		if err := r.unmarshalAttributeJSONB(&attr, attrNameBytes, attrDisplayNameBytes, attrOptionsBytes, attrValidationRulesBytes, attrUISettingsBytes); err != nil {
			return nil, err
		}

		lav.Attribute = &attr
		listingValues = append(listingValues, &lav)
	}

	if err := rows.Err(); err != nil {
		r.logger.Error().Err(err).Msg("error iterating listing attribute value rows")
		return nil, fmt.Errorf("failed to iterate listing attribute values: %w", err)
	}

	return listingValues, nil
}

// SetListingValues sets or updates attribute values for a listing
// Uses transaction to ensure atomicity
func (r *AttributeRepository) SetListingValues(ctx context.Context, listingID int32, values []domain.SetListingAttributeValue) error {
	if len(values) == 0 {
		return nil // Nothing to do
	}

	// Start transaction
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		r.logger.Error().Err(err).Msg("failed to begin transaction")
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	// Prepare statement for upsert
	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO listing_attribute_values (
			listing_id, attribute_id, value_text, value_number, value_boolean, value_date, value_json
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (listing_id, attribute_id)
		DO UPDATE SET
			value_text = EXCLUDED.value_text,
			value_number = EXCLUDED.value_number,
			value_boolean = EXCLUDED.value_boolean,
			value_date = EXCLUDED.value_date,
			value_json = EXCLUDED.value_json,
			updated_at = CURRENT_TIMESTAMP
	`)
	if err != nil {
		_ = tx.Rollback()
		r.logger.Error().Err(err).Msg("failed to prepare statement")
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	// Execute for each value
	for _, val := range values {
		// Marshal value_json if present - use interface{} to properly handle NULL
		var valueJSONBytes interface{} = nil
		if len(val.ValueJSON) > 0 {
			bytes, err := json.Marshal(val.ValueJSON)
			if err != nil {
				_ = tx.Rollback()
				return fmt.Errorf("failed to marshal value_json: %w", err)
			}
			valueJSONBytes = bytes
		}

		_, err = stmt.ExecContext(
			ctx,
			listingID,
			val.AttributeID,
			val.ValueText,
			val.ValueNumber,
			val.ValueBoolean,
			val.ValueDate,
			valueJSONBytes,
		)
		if err != nil {
			_ = tx.Rollback()
			r.logger.Error().Err(err).Int32("attribute_id", val.AttributeID).Msg("failed to set listing attribute value")
			return fmt.Errorf("failed to set listing attribute value: %w", err)
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		r.logger.Error().Err(err).Msg("failed to commit transaction")
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	r.logger.Info().Int32("listing_id", listingID).Int("count", len(values)).Msg("listing attribute values set")
	return nil
}

// DeleteListingValues deletes all attribute values for a listing
func (r *AttributeRepository) DeleteListingValues(ctx context.Context, listingID int32) error {
	query := `
		DELETE FROM listing_attribute_values
		WHERE listing_id = $1
	`

	result, err := r.db.ExecContext(ctx, query, listingID)
	if err != nil {
		r.logger.Error().Err(err).Int32("listing_id", listingID).Msg("failed to delete listing attribute values")
		return fmt.Errorf("failed to delete listing attribute values: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	r.logger.Info().Int32("listing_id", listingID).Int64("count", rowsAffected).Msg("listing attribute values deleted")
	return nil
}

// =============================================================================
// Variant Attribute Methods
// =============================================================================

// GetCategoryVariantAttributes retrieves variant attribute definitions for a category
func (r *AttributeRepository) GetCategoryVariantAttributes(ctx context.Context, categoryID string) ([]*domain.VariantAttribute, error) {
	query := `
		SELECT cva.id, cva.category_id, cva.attribute_id,
		       cva.is_required, cva.affects_price, cva.affects_stock, cva.sort_order, cva.display_as,
		       cva.is_active, cva.created_at, cva.updated_at,
		       a.id, a.code, a.name, a.display_name, a.attribute_type, a.purpose,
		       a.options, a.validation_rules, a.ui_settings,
		       a.is_searchable, a.is_filterable, a.is_required, a.is_variant_compatible,
		       a.affects_stock, a.affects_price, a.show_in_card,
		       a.is_active, a.sort_order, a.icon, a.created_at, a.updated_at
		FROM category_variant_attributes cva
		INNER JOIN attributes a ON cva.attribute_id = a.id
		WHERE cva.category_id = $1 AND cva.is_active = true AND a.is_active = true
		ORDER BY cva.sort_order ASC, (a.name->>'en') ASC
	`

	rows, err := r.db.QueryContext(ctx, query, categoryID)
	if err != nil {
		r.logger.Error().Err(err).Str("category_id", categoryID).Msg("failed to get category variant attributes")
		return nil, fmt.Errorf("failed to get category variant attributes: %w", err)
	}
	defer rows.Close()

	var variantAttrs []*domain.VariantAttribute
	for rows.Next() {
		var va domain.VariantAttribute
		var attr domain.Attribute
		var attrNameBytes, attrDisplayNameBytes, attrOptionsBytes, attrValidationRulesBytes, attrUISettingsBytes []byte

		err := rows.Scan(
			&va.ID,
			&va.CategoryID,
			&va.AttributeID,
			&va.IsRequired,
			&va.AffectsPrice,
			&va.AffectsStock,
			&va.SortOrder,
			&va.DisplayAs,
			&va.IsActive,
			&va.CreatedAt,
			&va.UpdatedAt,
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
			r.logger.Error().Err(err).Msg("failed to scan variant attribute row")
			return nil, fmt.Errorf("failed to scan variant attribute: %w", err)
		}

		// Unmarshal attribute JSONB fields
		if err := r.unmarshalAttributeJSONB(&attr, attrNameBytes, attrDisplayNameBytes, attrOptionsBytes, attrValidationRulesBytes, attrUISettingsBytes); err != nil {
			return nil, err
		}

		va.Attribute = &attr
		variantAttrs = append(variantAttrs, &va)
	}

	if err := rows.Err(); err != nil {
		r.logger.Error().Err(err).Msg("error iterating variant attribute rows")
		return nil, fmt.Errorf("failed to iterate variant attributes: %w", err)
	}

	return variantAttrs, nil
}

// GetVariantValues retrieves attribute values for a product variant
func (r *AttributeRepository) GetVariantValues(ctx context.Context, variantID int32) ([]*domain.VariantAttributeValue, error) {
	query := `
		SELECT vav.id, vav.variant_id, vav.attribute_id,
		       vav.value_text, vav.value_number, vav.value_boolean, vav.value_date, vav.value_json,
		       vav.price_modifier, vav.price_modifier_type,
		       vav.created_at, vav.updated_at,
		       a.id, a.code, a.name, a.display_name, a.attribute_type, a.purpose,
		       a.options, a.validation_rules, a.ui_settings,
		       a.is_searchable, a.is_filterable, a.is_required, a.is_variant_compatible,
		       a.affects_stock, a.affects_price, a.show_in_card,
		       a.is_active, a.sort_order, a.icon, a.created_at, a.updated_at
		FROM variant_attribute_values vav
		INNER JOIN attributes a ON vav.attribute_id = a.id
		WHERE vav.variant_id = $1 AND a.is_active = true
		ORDER BY a.sort_order ASC, (a.name->>'en') ASC
	`

	rows, err := r.db.QueryContext(ctx, query, variantID)
	if err != nil {
		r.logger.Error().Err(err).Int32("variant_id", variantID).Msg("failed to get variant attribute values")
		return nil, fmt.Errorf("failed to get variant attribute values: %w", err)
	}
	defer rows.Close()

	var variantValues []*domain.VariantAttributeValue
	for rows.Next() {
		var vav domain.VariantAttributeValue
		var attr domain.Attribute
		var valueJSONBytes []byte
		var attrNameBytes, attrDisplayNameBytes, attrOptionsBytes, attrValidationRulesBytes, attrUISettingsBytes []byte

		err := rows.Scan(
			&vav.ID,
			&vav.VariantID,
			&vav.AttributeID,
			&vav.ValueText,
			&vav.ValueNumber,
			&vav.ValueBoolean,
			&vav.ValueDate,
			&valueJSONBytes,
			&vav.PriceModifier,
			&vav.PriceModifierType,
			&vav.CreatedAt,
			&vav.UpdatedAt,
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
			r.logger.Error().Err(err).Msg("failed to scan variant attribute value row")
			return nil, fmt.Errorf("failed to scan variant attribute value: %w", err)
		}

		// Unmarshal value_json
		if len(valueJSONBytes) > 0 && string(valueJSONBytes) != "null" {
			if err := json.Unmarshal(valueJSONBytes, &vav.ValueJSON); err != nil {
				r.logger.Error().Err(err).Msg("failed to unmarshal value_json")
				return nil, fmt.Errorf("failed to unmarshal value_json: %w", err)
			}
		}

		// Unmarshal attribute JSONB fields
		if err := r.unmarshalAttributeJSONB(&attr, attrNameBytes, attrDisplayNameBytes, attrOptionsBytes, attrValidationRulesBytes, attrUISettingsBytes); err != nil {
			return nil, err
		}

		vav.Attribute = &attr
		variantValues = append(variantValues, &vav)
	}

	if err := rows.Err(); err != nil {
		r.logger.Error().Err(err).Msg("error iterating variant attribute value rows")
		return nil, fmt.Errorf("failed to iterate variant attribute values: %w", err)
	}

	return variantValues, nil
}
