package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"backend/internal/domain/models"
	"backend/internal/logger"
)

// GetCategoryAttributesImpl возвращает атрибуты категории (schema)
// Используется для получения списка доступных атрибутов для категории
func (db *Database) GetCategoryAttributesImpl(ctx context.Context, categoryID int) ([]models.CategoryAttribute, error) {
	query := `
		SELECT
			ua.id,
			ua.code as name,
			ua.display_name,
			ua.attribute_type,
			ua.icon,
			ua.options,
			ua.validation_rules,
			ua.is_searchable,
			ua.is_filterable,
			uca.is_required,
			ua.show_in_card,
			uca.sort_order,
			ua.created_at,
			ua.is_variant_compatible,
			ua.affects_stock
		FROM unified_attributes ua
		JOIN unified_category_attributes uca ON uca.attribute_id = ua.id
		WHERE uca.category_id = $1
		  AND uca.is_enabled = true
		  AND ua.is_active = true
		ORDER BY uca.sort_order ASC, ua.display_name ASC
	`

	rows, err := db.QueryContext(ctx, query, categoryID)
	if err != nil {
		logger.Error().Err(err).Int("category_id", categoryID).Msg("Failed to query category attributes")
		return nil, fmt.Errorf("failed to get category attributes: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var attributes []models.CategoryAttribute
	for rows.Next() {
		var attr models.CategoryAttribute
		var options, validRules sql.NullString
		var icon sql.NullString

		err := rows.Scan(
			&attr.ID,
			&attr.Name,
			&attr.DisplayName,
			&attr.AttributeType,
			&icon,
			&options,
			&validRules,
			&attr.IsSearchable,
			&attr.IsFilterable,
			&attr.IsRequired,
			&attr.ShowInCard,
			&attr.SortOrder,
			&attr.CreatedAt,
			&attr.IsVariantCompatible,
			&attr.AffectsStock,
		)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to scan category attribute")
			return nil, fmt.Errorf("failed to scan attribute: %w", err)
		}

		// Заполняем JSON поля
		if icon.Valid {
			attr.Icon = icon.String
		}
		if options.Valid {
			attr.Options = json.RawMessage(options.String)
		}
		if validRules.Valid {
			attr.ValidRules = json.RawMessage(validRules.String)
		}

		attributes = append(attributes, attr)
	}

	if err = rows.Err(); err != nil {
		logger.Error().Err(err).Msg("Error iterating category attributes")
		return nil, fmt.Errorf("error iterating attributes: %w", err)
	}

	logger.Debug().Int("category_id", categoryID).Int("count", len(attributes)).Msg("Retrieved category attributes")
	return attributes, nil
}

// SaveListingAttributesImpl сохраняет значения атрибутов листинга
// Удаляет старые значения и вставляет новые
func (db *Database) SaveListingAttributesImpl(ctx context.Context, listingID int, attributes []models.ListingAttributeValue) error {
	if len(attributes) == 0 {
		logger.Debug().Int("listing_id", listingID).Msg("No attributes to save")
		return nil
	}

	// Начинаем транзакцию
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to begin transaction")
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	// Удаляем старые значения атрибутов для listing
	deleteQuery := `
		DELETE FROM unified_attribute_values
		WHERE entity_type = $1 AND entity_id = $2
	`
	_, err = tx.Exec(ctx, deleteQuery, models.AttributeEntityTypeListing, listingID)
	if err != nil {
		logger.Error().Err(err).Int("listing_id", listingID).Msg("Failed to delete old attributes")
		return fmt.Errorf("failed to delete old attributes: %w", err)
	}

	// Вставляем новые значения
	insertQuery := `
		INSERT INTO unified_attribute_values
		(entity_type, entity_id, attribute_id, text_value, numeric_value, boolean_value, date_value, json_value)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	for _, attr := range attributes {
		_, err = tx.Exec(ctx, insertQuery,
			models.AttributeEntityTypeListing,
			listingID,
			attr.AttributeID,
			attr.TextValue,
			attr.NumericValue,
			attr.BooleanValue,
			nil, // date_value не используется в ListingAttributeValue
			attr.JSONValue,
		)
		if err != nil {
			logger.Error().
				Err(err).
				Int("listing_id", listingID).
				Int("attribute_id", attr.AttributeID).
				Msg("Failed to insert attribute value")
			return fmt.Errorf("failed to save attribute %d: %w", attr.AttributeID, err)
		}
	}

	// Коммитим транзакцию
	if err = tx.Commit(); err != nil {
		logger.Error().Err(err).Msg("Failed to commit transaction")
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	logger.Debug().Int("listing_id", listingID).Int("count", len(attributes)).Msg("Saved listing attributes")
	return nil
}

// GetListingAttributesImpl получает значения атрибутов листинга
// Возвращает полную информацию о каждом атрибуте со значением
func (db *Database) GetListingAttributesImpl(ctx context.Context, listingID int) ([]models.ListingAttributeValue, error) {
	query := `
		SELECT
			uav.attribute_id,
			ua.code as attribute_name,
			ua.display_name,
			ua.attribute_type,
			ua.icon,
			uav.text_value,
			uav.numeric_value,
			uav.boolean_value,
			uav.json_value,
			uca.is_required,
			ua.show_in_card
		FROM unified_attribute_values uav
		JOIN unified_attributes ua ON ua.id = uav.attribute_id
		LEFT JOIN unified_category_attributes uca ON uca.attribute_id = ua.id
		LEFT JOIN c2c_listings cl ON cl.id = uav.entity_id AND uav.entity_type = $2
		WHERE uav.entity_type = $2
		  AND uav.entity_id = $1
		  AND ua.is_active = true
		  AND (uca.category_id = cl.category_id OR uca.category_id IS NULL)
		ORDER BY uca.sort_order ASC, ua.display_name ASC
	`

	rows, err := db.QueryContext(ctx, query, listingID, models.AttributeEntityTypeListing)
	if err != nil {
		logger.Error().Err(err).Int("listing_id", listingID).Msg("Failed to query listing attributes")
		return nil, fmt.Errorf("failed to get listing attributes: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var attributes []models.ListingAttributeValue
	for rows.Next() {
		var attr models.ListingAttributeValue
		var icon sql.NullString
		var textVal sql.NullString
		var numVal sql.NullFloat64
		var boolVal sql.NullBool
		var jsonVal sql.NullString

		err := rows.Scan(
			&attr.AttributeID,
			&attr.AttributeName,
			&attr.DisplayName,
			&attr.AttributeType,
			&icon,
			&textVal,
			&numVal,
			&boolVal,
			&jsonVal,
			&attr.IsRequired,
			&attr.ShowInCard,
		)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to scan listing attribute")
			return nil, fmt.Errorf("failed to scan attribute: %w", err)
		}

		attr.ListingID = listingID

		// Обрабатываем nullable поля
		if textVal.Valid {
			attr.TextValue = &textVal.String
			attr.DisplayValue = textVal.String
		}
		if numVal.Valid {
			attr.NumericValue = &numVal.Float64
			if attr.DisplayValue == "" {
				attr.DisplayValue = fmt.Sprintf("%.2f", numVal.Float64)
			}
		}
		if boolVal.Valid {
			attr.BooleanValue = &boolVal.Bool
			if attr.DisplayValue == "" {
				if boolVal.Bool {
					attr.DisplayValue = "Yes"
				} else {
					attr.DisplayValue = "No"
				}
			}
		}
		if jsonVal.Valid {
			attr.JSONValue = json.RawMessage(jsonVal.String)
			if attr.DisplayValue == "" {
				attr.DisplayValue = jsonVal.String
			}
		}

		attributes = append(attributes, attr)
	}

	if err = rows.Err(); err != nil {
		logger.Error().Err(err).Msg("Error iterating listing attributes")
		return nil, fmt.Errorf("error iterating attributes: %w", err)
	}

	logger.Debug().Int("listing_id", listingID).Int("count", len(attributes)).Msg("Retrieved listing attributes")
	return attributes, nil
}

// GetAttributeRangesImpl возвращает диапазоны значений для фильтруемых атрибутов
// Используется для построения фильтров (min/max цена, год выпуска и т.д.)
func (db *Database) GetAttributeRangesImpl(ctx context.Context, categoryID int) (map[string]map[string]interface{}, error) {
	query := `
		SELECT
			ua.code as attribute_name,
			ua.display_name,
			ua.attribute_type,
			MIN(uav.numeric_value) as min_value,
			MAX(uav.numeric_value) as max_value,
			COUNT(DISTINCT uav.entity_id) as count
		FROM unified_attributes ua
		JOIN unified_category_attributes uca ON uca.attribute_id = ua.id
		JOIN unified_attribute_values uav ON uav.attribute_id = ua.id
		JOIN c2c_listings cl ON cl.id = uav.entity_id AND uav.entity_type = $2
		WHERE uca.category_id = $1
		  AND uca.is_enabled = true
		  AND ua.is_active = true
		  AND ua.is_filterable = true
		  AND ua.attribute_type = 'number'
		  AND uav.numeric_value IS NOT NULL
		  AND cl.status = 'active'
		GROUP BY ua.code, ua.display_name, ua.attribute_type
		ORDER BY ua.display_name ASC
	`

	rows, err := db.QueryContext(ctx, query, categoryID, models.AttributeEntityTypeListing)
	if err != nil {
		logger.Error().Err(err).Int("category_id", categoryID).Msg("Failed to query attribute ranges")
		return nil, fmt.Errorf("failed to get attribute ranges: %w", err)
	}
	defer func() { _ = rows.Close() }()

	ranges := make(map[string]map[string]interface{})
	for rows.Next() {
		var attrName, displayName, attrType string
		var minVal, maxVal sql.NullFloat64
		var count int

		err := rows.Scan(&attrName, &displayName, &attrType, &minVal, &maxVal, &count)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to scan attribute range")
			return nil, fmt.Errorf("failed to scan range: %w", err)
		}

		rangeInfo := map[string]interface{}{
			"display_name": displayName,
			"type":         attrType,
			"count":        count,
		}

		if minVal.Valid {
			rangeInfo["min"] = minVal.Float64
		}
		if maxVal.Valid {
			rangeInfo["max"] = maxVal.Float64
		}

		ranges[attrName] = rangeInfo
	}

	if err = rows.Err(); err != nil {
		logger.Error().Err(err).Msg("Error iterating attribute ranges")
		return nil, fmt.Errorf("error iterating ranges: %w", err)
	}

	logger.Debug().Int("category_id", categoryID).Int("count", len(ranges)).Msg("Retrieved attribute ranges")
	return ranges, nil
}
