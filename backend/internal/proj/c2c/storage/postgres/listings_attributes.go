// backend/internal/proj/c2c/storage/postgres/listings_attributes.go
package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"

	"backend/internal/common"
	"backend/internal/domain/models"

	"github.com/jackc/pgx/v5"
)

// ============================================================================
// UTILITY FUNCTIONS
// ============================================================================

// sanitizeAttributeValue очищает и нормализует значения атрибутов
func sanitizeAttributeValue(attr *models.ListingAttributeValue) {
	// Ограничение длины текстовых атрибутов
	if attr.TextValue != nil {
		if len(*attr.TextValue) > 1000 {
			truncated := (*attr.TextValue)[:1000]
			attr.TextValue = &truncated
			log.Printf("Attribute value truncated for attribute %s (ID: %d)",
				attr.AttributeName, attr.AttributeID)
		}
	}

	// Проверка на NaN и Inf для числовых атрибутов
	if attr.NumericValue != nil {
		numVal := *attr.NumericValue
		if math.IsNaN(numVal) || math.IsInf(numVal, 0) {
			defaultVal := 0.0
			attr.NumericValue = &defaultVal
			log.Printf("Invalid numeric value (NaN/Inf) replaced with 0 for attribute %s (ID: %d)",
				attr.AttributeName, attr.AttributeID)
		}
	}

	// Стандартизация обработки пустых значений
	if attr.TextValue != nil && *attr.TextValue == "" {
		attr.TextValue = nil // Пустые строки -> NULL
	}

	if attr.NumericValue != nil && *attr.NumericValue == 0 {
		// Для некоторых атрибутов нуль может быть валидным значением
		if !isZeroValidValue(attr.AttributeName) {
			attr.NumericValue = nil
		}
	}

	// Если все значения NULL, устанавливаем DisplayValue в пустую строку
	if attr.TextValue == nil && attr.NumericValue == nil &&
		attr.BooleanValue == nil && attr.JSONValue == nil {
		attr.DisplayValue = ""
	}
}

// isZeroValidValue определяет, является ли нулевое значение допустимым для атрибута
func isZeroValidValue(attrName string) bool {
	// Для этих атрибутов ноль - допустимое значение
	zeroValidAttrs := map[string]bool{
		"floor":         true, // Например, цокольный этаж
		attrNameMileage: true, // Для новых автомобилей
		"price":         true, // Для бесплатных объявлений
	}
	return zeroValidAttrs[attrName]
}

// ============================================================================
// SAVE & FORMAT METHODS
// ============================================================================

// SaveListingAttributes сохраняет значения атрибутов для объявления
func (s *Storage) SaveListingAttributes(ctx context.Context, listingID int, attributes []models.ListingAttributeValue) error {
	log.Printf("Saving %d attributes for listing %d", len(attributes), listingID)

	// Начинаем транзакцию
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil {
			_ = err // Explicitly ignore error if transaction was already committed
		}
	}()

	// Удаляем старые атрибуты
	_, err = tx.Exec(ctx, `DELETE FROM listing_attribute_values WHERE listing_id = $1`, listingID)
	if err != nil {
		return fmt.Errorf("error deleting old attributes: %w", err)
	}

	// Проверяем, есть ли атрибуты для сохранения
	if len(attributes) == 0 {
		log.Printf("Storage: No attributes to save for listing %d", listingID)
		return tx.Commit(ctx)
	}

	// Подготовка данных для bulk insert
	valueStrings, valueArgs := s.buildAttributeInsertData(attributes, listingID)

	// Если нет атрибутов для вставки, завершаем транзакцию
	if len(valueStrings) == 0 {
		log.Printf("Storage: No valid attributes found for listing %d after filtering", listingID)
		return tx.Commit(ctx)
	}

	// Выполняем bulk insert
	if err := s.executeBulkAttributeInsert(ctx, tx, valueStrings, valueArgs); err != nil {
		return err
	}

	// Фиксируем транзакцию
	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}

	log.Printf("Storage: Successfully saved %d unique attributes for listing %d", len(valueStrings), listingID)
	return nil
}

// buildAttributeInsertData подготавливает данные для bulk insert атрибутов
func (s *Storage) buildAttributeInsertData(attributes []models.ListingAttributeValue, listingID int) ([]string, []interface{}) {
	seen := make(map[int]bool)
	valueStrings := make([]string, 0, len(attributes))
	valueArgs := make([]interface{}, 0, len(attributes)*7)
	counter := 1

	for i, attr := range attributes {
		// Санитизация значений атрибутов
		sanitizeAttributeValue(&attr)
		attributes[i] = attr

		// Проверка на нулевые или некорректные attribute_id
		if attr.AttributeID <= 0 {
			log.Printf("Storage: Invalid attribute ID: %d, skipping", attr.AttributeID)
			continue
		}

		// Проверка на дубликаты по attribute_id
		if seen[attr.AttributeID] {
			log.Printf("Storage: Duplicate attribute ID %d for listing %d, skipping", attr.AttributeID, listingID)
			continue
		}
		seen[attr.AttributeID] = true

		// Проверяем, что есть хотя бы одно значение для сохранения
		hasValue := attr.TextValue != nil || attr.NumericValue != nil ||
			attr.BooleanValue != nil || attr.JSONValue != nil ||
			attr.DisplayValue != ""
		if !hasValue {
			log.Printf("Storage: No value provided for attribute %d, skipping", attr.AttributeID)
			continue
		}

		// Определяем единицу измерения
		unit := s.determineAttributeUnit(attr)

		// Числовые атрибуты - дополнительная обработка (конвертация текста в число)
		s.convertTextToNumeric(&attr)

		// Подготавливаем часть запроса для этого атрибута
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d)",
			counter, counter+1, counter+2, counter+3, counter+4, counter+5, counter+6))

		// Добавляем параметры
		valueArgs = append(valueArgs, listingID, attr.AttributeID)

		// Текстовое значение
		if attr.TextValue != nil && *attr.TextValue != "" {
			valueArgs = append(valueArgs, *attr.TextValue)
		} else {
			valueArgs = append(valueArgs, nil)
		}

		// Числовое значение с проверками
		if attr.NumericValue != nil {
			numericVal := *attr.NumericValue
			if math.IsNaN(numericVal) || math.IsInf(numericVal, 0) {
				log.Printf("Storage: Invalid numeric value (NaN/Inf) for attribute %d, using 0", attr.AttributeID)
				numericVal = 0.0
			}
			valueArgs = append(valueArgs, numericVal)
		} else {
			valueArgs = append(valueArgs, nil)
		}

		// Логическое значение
		if attr.BooleanValue != nil {
			valueArgs = append(valueArgs, *attr.BooleanValue)
		} else {
			valueArgs = append(valueArgs, nil)
		}

		// JSON значение
		if len(attr.JSONValue) > 0 {
			valueArgs = append(valueArgs, string(attr.JSONValue))
		} else {
			valueArgs = append(valueArgs, nil)
		}

		// Единица измерения
		valueArgs = append(valueArgs, unit)

		counter += 7
	}

	return valueStrings, valueArgs
}

// determineAttributeUnit определяет единицу измерения для атрибута
func (s *Storage) determineAttributeUnit(attr models.ListingAttributeValue) string {
	if attr.Unit != "" {
		return attr.Unit
	}

	// Определяем единицу измерения на основе имени атрибута
	switch attr.AttributeName {
	case attrNameArea:
		return "m²"
	case attrNameLandArea:
		return "ar"
	case attrNameMileage:
		return "km"
	case attrNameEngineCapacity:
		return "l"
	case attrNamePower:
		return "ks"
	case "screen_size":
		return "inč"
	case "rooms":
		return "soba"
	case "floor", "total_floors":
		return "sprat"
	default:
		return ""
	}
}

// convertTextToNumeric конвертирует текстовое значение в число для числовых атрибутов
func (s *Storage) convertTextToNumeric(attr *models.ListingAttributeValue) {
	if attr.NumericValue != nil || attr.TextValue == nil || *attr.TextValue == "" {
		return
	}

	// Список числовых атрибутов для конвертации
	numericAttrs := map[string]bool{
		"rooms": true, "floor": true, "total_floors": true,
		attrNameArea: true, attrNameLandArea: true, attrNameMileage: true,
		attrNameYear: true, attrNameEngineCapacity: true, attrNamePower: true,
		"screen_size": true,
	}

	if numericAttrs[attr.AttributeName] {
		// Преобразуем текст в число
		clean := regexp.MustCompile(`[^\d\.-]`).ReplaceAllString(*attr.TextValue, "")
		if numVal, err := strconv.ParseFloat(clean, 64); err == nil {
			attr.NumericValue = &numVal
			log.Printf("Converted text value '%s' to numeric: %f for attribute %s",
				*attr.TextValue, numVal, attr.AttributeName)
		}
	}
}

// executeBulkAttributeInsert выполняет bulk insert атрибутов
func (s *Storage) executeBulkAttributeInsert(ctx context.Context, tx pgx.Tx, valueStrings []string, valueArgs []interface{}) error {
	query := fmt.Sprintf(`
        INSERT INTO listing_attribute_values (
            listing_id, attribute_id, text_value, numeric_value, boolean_value, json_value, unit
        ) VALUES %s
        ON CONFLICT (listing_id, attribute_id) DO UPDATE SET
            text_value = EXCLUDED.text_value,
            numeric_value = EXCLUDED.numeric_value,
            boolean_value = EXCLUDED.boolean_value,
            json_value = EXCLUDED.json_value,
            unit = EXCLUDED.unit
    `, strings.Join(valueStrings, ","))

	_, err := tx.Exec(ctx, query, valueArgs...)
	if err != nil {
		log.Printf("Storage: Error executing bulk insert: %v", err)
		log.Printf("Storage: Query: %s", query)
		log.Printf("Storage: Args: %+v", valueArgs)
		return fmt.Errorf("error inserting attribute values: %w", err)
	}

	return nil
}

// GetFormattedAttributeValue возвращает отформатированное значение атрибута с переводом
func (s *Storage) GetFormattedAttributeValue(ctx context.Context, attr models.ListingAttributeValue, language string) string {
	// Для числовых атрибутов с единицей измерения
	if attr.NumericValue != nil && attr.Unit != "" {
		// Получаем перевод единицы измерения
		var displayFormat string
		err := s.pool.QueryRow(ctx, `
            SELECT display_format FROM unit_translations
            WHERE unit = $1 AND language = $2
        `, attr.Unit, language).Scan(&displayFormat)

		if err == nil && displayFormat != "" {
			return fmt.Sprintf(displayFormat, *attr.NumericValue)
		}

		// Если не нашли перевод, используем стандартный формат
		return fmt.Sprintf("%g %s", *attr.NumericValue, attr.Unit)
	}

	// Для других типов атрибутов возвращаем DisplayValue
	return attr.DisplayValue
}

// ============================================================================
// GET LISTING ATTRIBUTES (OPTIMIZED)
// ============================================================================

// GetListingAttributes получает значения атрибутов для объявления
func (s *Storage) GetListingAttributes(ctx context.Context, listingID int) ([]models.ListingAttributeValue, error) {
	rows, err := s.queryListingAttributeValues(ctx, listingID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	log.Printf("Запрос атрибутов для объявления %d", listingID)
	return s.processListingAttributeRows(ctx, rows)
}

// queryListingAttributeValues выполняет SQL запрос атрибутов листинга
func (s *Storage) queryListingAttributeValues(ctx context.Context, listingID int) (pgx.Rows, error) {
	query := `
        WITH attribute_translations AS (
            SELECT
                entity_id,
                jsonb_object_agg(language, translated_text) as translations
            FROM translations
            WHERE entity_type = 'attribute' AND field_name = 'display_name'
            GROUP BY entity_id
        ),
        option_translations AS (
            SELECT
                attribute_name,
                jsonb_object_agg(lang, options_json) as option_translations
            FROM (
                SELECT
                    attribute_name,
                    lang,
                    jsonb_object_agg(option_value, translation) as options_json
                FROM (
                    SELECT
                        attribute_name, 'ru' as lang, option_value, ru_translation as translation
                    FROM attribute_option_translations
                    UNION ALL
                    SELECT
                        attribute_name, 'sr' as lang, option_value, sr_translation as translation
                    FROM attribute_option_translations
                ) o
                GROUP BY attribute_name, lang
            ) grouped
            GROUP BY attribute_name
        )
        SELECT DISTINCT ON (a.id)
            v.listing_id, v.attribute_id, a.name AS attribute_name, a.display_name,
            a.attribute_type, v.text_value, v.numeric_value, v.boolean_value,
            v.json_value, v.unit, a.is_required as is_required, a.show_in_card as show_in_card,
            false as show_in_list,
            COALESCE(at.translations, '{}'::jsonb) as translations,
            COALESCE(ot.option_translations, '{}'::jsonb) as option_translations
        FROM listing_attribute_values v
        JOIN unified_attributes a ON v.attribute_id = a.id
        JOIN c2c_listings ml ON ml.id = v.listing_id
        LEFT JOIN attribute_translations at ON at.entity_id = a.id
        LEFT JOIN option_translations ot ON ot.attribute_name = a.name
        WHERE v.listing_id = $1
        ORDER BY a.id, a.sort_order, a.display_name
    `

	rows, err := s.pool.Query(ctx, query, listingID)
	if err != nil {
		return nil, fmt.Errorf("error querying listing attributes: %w", err)
	}

	return rows, nil
}

// processListingAttributeRows обрабатывает результаты запроса атрибутов
func (s *Storage) processListingAttributeRows(ctx context.Context, rows pgx.Rows) ([]models.ListingAttributeValue, error) {
	var allAttributes []models.ListingAttributeValue
	seen := make(map[int]bool) // Защита от дубликатов

	// Получаем язык из контекста
	locale := "sr"
	if lang, ok := ctx.Value(common.ContextKeyLocale).(string); ok && lang != "" {
		locale = lang
	}

	for rows.Next() {
		attr, err := s.scanListingAttributeRow(rows, locale)
		if err != nil {
			return nil, err
		}

		// Проверяем дубликаты
		if seen[attr.AttributeID] {
			continue
		}
		seen[attr.AttributeID] = true

		allAttributes = append(allAttributes, attr)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating listing attributes: %w", err)
	}

	return allAttributes, nil
}

// scanListingAttributeRow сканирует одну строку атрибута
func (s *Storage) scanListingAttributeRow(rows pgx.Rows, locale string) (models.ListingAttributeValue, error) {
	var attr models.ListingAttributeValue
	var textValue sql.NullString
	var numericValue sql.NullFloat64
	var boolValue sql.NullBool
	var jsonValue sql.NullString
	var unit sql.NullString
	var translationsJson []byte
	var optionTranslationsJson []byte

	if err := rows.Scan(
		&attr.ListingID, &attr.AttributeID, &attr.AttributeName, &attr.DisplayName,
		&attr.AttributeType, &textValue, &numericValue, &boolValue,
		&jsonValue, &unit, &attr.IsRequired, &attr.ShowInCard,
		&attr.ShowInList, &translationsJson, &optionTranslationsJson,
	); err != nil {
		log.Printf("Error scanning attribute: %v", err)
		return attr, fmt.Errorf("error scanning listing attribute: %w", err)
	}

	// Парсинг переводов
	if err := json.Unmarshal(translationsJson, &attr.Translations); err != nil {
		attr.Translations = make(map[string]string)
	}

	if err := json.Unmarshal(optionTranslationsJson, &attr.OptionTranslations); err != nil {
		attr.OptionTranslations = make(map[string]map[string]string)
	}

	// Обработка значений атрибута
	if unit.Valid {
		attr.Unit = unit.String
	}

	s.populateAttributeValue(&attr, textValue, numericValue, boolValue, jsonValue, locale)

	return attr, nil
}

// populateAttributeValue заполняет значение атрибута в зависимости от типа
func (s *Storage) populateAttributeValue(attr *models.ListingAttributeValue, textValue sql.NullString,
	numericValue sql.NullFloat64, boolValue sql.NullBool, jsonValue sql.NullString, locale string) {

	// Текстовое значение (для select ищем перевод)
	if textValue.Valid {
		attr.TextValue = &textValue.String
		if attr.AttributeType == "select" && attr.OptionTranslations != nil {
			attr.DisplayValue = s.translateSelectOption(textValue.String, attr.OptionTranslations, locale)
		} else {
			attr.DisplayValue = textValue.String
		}
	}

	// Числовое значение
	if numericValue.Valid {
		attr.NumericValue = &numericValue.Float64
		attr.DisplayValue = s.formatNumericValue(numericValue.Float64, attr.AttributeName, attr.Unit)
	}

	// Логическое значение
	if boolValue.Valid {
		attr.BooleanValue = &boolValue.Bool
		if boolValue.Bool {
			attr.DisplayValue = "Да"
		} else {
			attr.DisplayValue = "Нет"
		}
	}

	// JSON значение (для multiselect форматируем массив)
	if jsonValue.Valid {
		attr.JSONValue = json.RawMessage(jsonValue.String)
		if attr.AttributeType == "multiselect" {
			var values []string
			if err := json.Unmarshal(attr.JSONValue, &values); err == nil {
				attr.DisplayValue = strings.Join(values, ", ")
			}
		} else {
			attr.DisplayValue = jsonValue.String
		}
	}
}

// translateSelectOption переводит значение select атрибута
func (s *Storage) translateSelectOption(value string, optionTranslations map[string]map[string]string, locale string) string {
	if langTranslations, ok := optionTranslations[locale]; ok {
		// Сначала проверяем точное совпадение
		if translation, ok := langTranslations[value]; ok {
			return translation
		}
		// Пробуем найти перевод в нижнем регистре
		lowerValue := strings.ToLower(value)
		if translation, ok := langTranslations[lowerValue]; ok {
			return translation
		}
	}
	return value // Если перевода нет, используем оригинальное значение
}

// formatNumericValue форматирует числовое значение с учетом единиц измерения
func (s *Storage) formatNumericValue(value float64, attrName, unit string) string {
	// Для года выводим без дробной части
	if attrName == attrNameYear {
		return fmt.Sprintf("%d", int(value))
	}

	// Определяем единицу измерения если не указана
	unitStr := unit
	if unitStr == "" {
		switch attrName {
		case attrNameArea:
			unitStr = "m²"
		case attrNameLandArea:
			unitStr = "ar"
		case attrNameMileage:
			unitStr = "km"
		case attrNameEngineCapacity:
			unitStr = "l"
		case attrNamePower:
			unitStr = "ks"
		case "screen_size":
			unitStr = "inč"
		}
	}

	if unitStr != "" {
		return fmt.Sprintf("%g %s", value, unitStr)
	}
	return fmt.Sprintf("%g", value)
}

// ============================================================================
// GET CATEGORY ATTRIBUTES (OPTIMIZED)
// ============================================================================

// GetCategoryAttributes получает атрибуты категории с кэшированием
func (s *Storage) GetCategoryAttributes(ctx context.Context, categoryID int) ([]models.CategoryAttribute, error) {
	// Проверяем кэш
	if cached, ok := s.getCachedCategoryAttributes(categoryID); ok {
		log.Printf("Using cached attributes for category %d", categoryID)
		return cached, nil
	}

	log.Printf("GetCategoryAttributes: Получение атрибутов для категории %d", categoryID)

	// Загружаем из БД
	attrs, err := s.fetchCategoryAttributesFromDB(ctx, categoryID)
	if err != nil {
		return nil, err
	}

	// Обогащаем атрибуты реальными диапазонами
	if err := s.enrichAttributesWithRanges(ctx, categoryID, attrs); err != nil {
		log.Printf("GetCategoryAttributes: Ошибка обогащения диапазонами: %v", err)
		// Продолжаем без диапазонов
	}

	// Сохраняем в кэш
	s.setCachedCategoryAttributes(categoryID, attrs)

	log.Printf("GetCategoryAttributes: Успешно получено %d атрибутов для категории %d", len(attrs), categoryID)
	return attrs, nil
}

// getCachedCategoryAttributes получает атрибуты из кэша (thread-safe)
func (s *Storage) getCachedCategoryAttributes(categoryID int) ([]models.CategoryAttribute, bool) {
	s.attributeCacheMutex.RLock()
	defer s.attributeCacheMutex.RUnlock()

	cached, ok := s.attributeCache[categoryID]
	if !ok {
		return nil, false
	}

	cacheTime, ok := s.attributeCacheTime[categoryID]
	if !ok || time.Since(cacheTime) > s.cacheTTL {
		return nil, false
	}

	return cached, true
}

// setCachedCategoryAttributes сохраняет атрибуты в кэш (thread-safe)
func (s *Storage) setCachedCategoryAttributes(categoryID int, attrs []models.CategoryAttribute) {
	s.attributeCacheMutex.Lock()
	defer s.attributeCacheMutex.Unlock()

	s.attributeCache[categoryID] = attrs
	s.attributeCacheTime[categoryID] = time.Now()
}

// fetchCategoryAttributesFromDB загружает атрибуты категории из БД
func (s *Storage) fetchCategoryAttributesFromDB(ctx context.Context, categoryID int) ([]models.CategoryAttribute, error) {
	query := `
    WITH RECURSIVE category_hierarchy AS (
        WITH RECURSIVE parents AS (
            SELECT id, parent_id FROM c2c_categories WHERE id = $1
            UNION
            SELECT c.id, c.parent_id FROM c2c_categories c
            INNER JOIN parents p ON c.id = p.parent_id
        )
        SELECT id FROM parents
    ),
    attribute_translations AS (
        SELECT entity_id, jsonb_object_agg(language, translated_text) as translations
        FROM translations
        WHERE entity_type = 'attribute' AND field_name = 'display_name'
        GROUP BY entity_id
    ),
    option_translations AS (
        SELECT entity_id, language, jsonb_object_agg(field_name, translated_text) as field_translations
        FROM translations
        WHERE entity_type = 'attribute_option'
        GROUP BY entity_id, language
    ),
    option_lang_agg AS (
        SELECT entity_id, jsonb_object_agg(language, field_translations) as option_translations
        FROM option_translations
        GROUP BY entity_id
    )
    SELECT DISTINCT ON (a.id)
        a.id, a.name, a.display_name, a.icon, a.attribute_type, a.options,
        a.validation_rules, a.is_searchable, a.is_filterable,
        COALESCE(m.is_required, a.is_required) as is_required,
        a.sort_order, a.created_at,
        COALESCE(m.custom_component, a.custom_component) as custom_component,
        COALESCE(at.translations, '{}'::jsonb) as translations,
        COALESCE(ol.option_translations, '{}'::jsonb) as option_translations
    FROM category_attribute_mapping m
    JOIN category_attributes a ON m.attribute_id = a.id
    JOIN category_hierarchy h ON m.category_id = h.id
    LEFT JOIN attribute_translations at ON a.id = at.entity_id
    LEFT JOIN option_lang_agg ol ON a.id = ol.entity_id
    WHERE m.is_enabled = true
    ORDER BY a.id, m.category_id = $1 DESC, a.sort_order, a.display_name
    `

	rows, err := s.pool.Query(ctx, query, categoryID)
	if err != nil {
		log.Printf("GetCategoryAttributes: Ошибка запроса: %v", err)
		return nil, fmt.Errorf("error querying category attributes: %w", err)
	}
	defer rows.Close()

	var attributes []models.CategoryAttribute
	for rows.Next() {
		attr, err := s.scanCategoryAttributeRow(rows)
		if err != nil {
			return nil, err
		}
		attributes = append(attributes, attr)
	}

	if err := rows.Err(); err != nil {
		log.Printf("GetCategoryAttributes: Ошибка при итерации результатов: %v", err)
		return nil, fmt.Errorf("error iterating category attributes: %w", err)
	}

	return attributes, nil
}

// scanCategoryAttributeRow сканирует одну строку атрибута категории
func (s *Storage) scanCategoryAttributeRow(rows pgx.Rows) (models.CategoryAttribute, error) {
	var attr models.CategoryAttribute
	var options, validRules, customComponent sql.NullString
	var translationsJson, optionTranslationsJson []byte

	if err := rows.Scan(
		&attr.ID, &attr.Name, &attr.DisplayName, &attr.Icon, &attr.AttributeType,
		&options, &validRules, &attr.IsSearchable, &attr.IsFilterable,
		&attr.IsRequired, &attr.SortOrder, &attr.CreatedAt,
		&customComponent, &translationsJson, &optionTranslationsJson,
	); err != nil {
		log.Printf("GetCategoryAttributes: Ошибка при сканировании результата: %v", err)
		return attr, fmt.Errorf("error scanning category attribute: %w", err)
	}

	// Обработка опциональных JSON полей
	if options.Valid && len(options.String) > 0 {
		attr.Options = json.RawMessage(options.String)
	} else {
		attr.Options = json.RawMessage(`{}`)
	}

	if validRules.Valid && len(validRules.String) > 0 {
		attr.ValidRules = json.RawMessage(validRules.String)
	} else {
		attr.ValidRules = json.RawMessage(`{}`)
	}

	attr.CustomComponent = ""
	if customComponent.Valid {
		attr.CustomComponent = customComponent.String
	}

	// Парсинг переводов
	attr.Translations = make(map[string]string)
	if err := json.Unmarshal(translationsJson, &attr.Translations); err != nil {
		log.Printf("GetCategoryAttributes: Ошибка парсинга переводов для атрибута %d: %v", attr.ID, err)
	}

	attr.OptionTranslations = make(map[string]map[string]string)
	if err := json.Unmarshal(optionTranslationsJson, &attr.OptionTranslations); err != nil {
		log.Printf("GetCategoryAttributes: Ошибка парсинга переводов опций для атрибута %d: %v", attr.ID, err)
	}

	return attr, nil
}

// enrichAttributesWithRanges обогащает числовые атрибуты реальными диапазонами
func (s *Storage) enrichAttributesWithRanges(ctx context.Context, categoryID int, attributes []models.CategoryAttribute) error {
	attributeRanges, err := s.GetAttributeRanges(ctx, categoryID)
	if err != nil {
		return err
	}

	for i, attr := range attributes {
		if attr.AttributeType == "number" {
			if ranges, ok := attributeRanges[attr.Name]; ok {
				var options map[string]interface{}

				// Используем существующие options или создаем новые
				if len(attr.Options) > 0 {
					if err := json.Unmarshal(attr.Options, &options); err != nil {
						options = make(map[string]interface{})
					}
				} else {
					options = make(map[string]interface{})
				}

				// Обновляем значения диапазонов
				options["min"] = ranges["min"]
				options["max"] = ranges["max"]
				options["step"] = ranges["step"]
				options["real_data"] = true

				// Сериализуем обратно в JSON
				if optionsJSON, err := json.Marshal(options); err == nil {
					attributes[i].Options = optionsJSON
					log.Printf("GetCategoryAttributes: Обновлены диапазоны для атрибута %s: min=%.2f, max=%.2f",
						attr.Name, ranges["min"], ranges["max"])
				}
			}
		}
	}

	return nil
}

// ============================================================================
// ATTRIBUTE RANGES & CACHE
// ============================================================================

// GetAttributeRanges получает минимальные и максимальные значения для числовых атрибутов
func (s *Storage) GetAttributeRanges(ctx context.Context, categoryID int) (map[string]map[string]interface{}, error) {
	// Проверяем кэш
	s.rangesCacheMutex.RLock()
	cachedRanges, hasCached := s.rangesCache[categoryID]
	cacheTime, hasTime := s.rangesCacheTime[categoryID]
	s.rangesCacheMutex.RUnlock()

	if hasCached && hasTime && time.Since(cacheTime) < s.cacheTTL {
		log.Printf("Using cached attribute ranges for category %d", categoryID)
		return cachedRanges, nil
	}

	// Получаем ID всех подкатегорий
	categoryIDs, err := s.getCategoryTreeIDs(ctx, categoryID)
	if err != nil {
		return nil, err
	}

	// Запрашиваем диапазоны
	ranges, err := s.queryAttributeRanges(ctx, categoryIDs)
	if err != nil {
		return nil, err
	}

	// Заполняем отсутствующие атрибуты значениями по умолчанию
	s.fillDefaultRanges(ranges)

	// Кешируем результат
	s.rangesCacheMutex.Lock()
	s.rangesCache[categoryID] = ranges
	s.rangesCacheTime[categoryID] = time.Now()
	s.rangesCacheMutex.Unlock()

	return ranges, nil
}

// getCategoryTreeIDs получает все ID категорий в дереве
func (s *Storage) getCategoryTreeIDs(ctx context.Context, categoryID int) (string, error) {
	query := `
    WITH RECURSIVE category_tree AS (
        SELECT id FROM c2c_categories WHERE id = $1
        UNION ALL
        SELECT c.id FROM c2c_categories c
        JOIN category_tree t ON c.parent_id = t.id
    )
    SELECT string_agg(id::text, ',') FROM category_tree
    `

	var categoryIDs string
	err := s.pool.QueryRow(ctx, query, categoryID).Scan(&categoryIDs)
	if err != nil {
		return "", fmt.Errorf("error getting category tree: %w", err)
	}

	if categoryIDs == "" {
		categoryIDs = strconv.Itoa(categoryID)
	}

	return categoryIDs, nil
}

// queryAttributeRanges запрашивает границы числовых атрибутов
func (s *Storage) queryAttributeRanges(ctx context.Context, categoryIDs string) (map[string]map[string]interface{}, error) {
	rangesQuery := `
    SELECT a.name, MIN(v.numeric_value) as min_value, MAX(v.numeric_value) as max_value,
           COUNT(DISTINCT v.numeric_value) as value_count
    FROM listing_attribute_values v
    JOIN category_attributes a ON v.attribute_id = a.id
    JOIN c2c_listings l ON v.listing_id = l.id
    WHERE l.category_id IN (` + categoryIDs + `)
        AND l.status = 'active'
        AND v.numeric_value IS NOT NULL
        AND a.attribute_type = 'number'
    GROUP BY a.name
    `

	rows, err := s.pool.Query(ctx, rangesQuery)
	if err != nil {
		return nil, fmt.Errorf("error querying attribute ranges: %w", err)
	}
	defer rows.Close()

	ranges := make(map[string]map[string]interface{})
	for rows.Next() {
		var attrName string
		var minValue, maxValue float64
		var valueCount int

		if err := rows.Scan(&attrName, &minValue, &maxValue, &valueCount); err != nil {
			return nil, fmt.Errorf("error scanning attribute range: %w", err)
		}

		// Округляем для целочисленных параметров
		if attrName == attrNameYear || attrName == "rooms" || attrName == "floor" || attrName == "total_floors" {
			minValue = float64(int(minValue))
			maxValue = float64(int(maxValue))
		}

		// Для года добавляем запас +1 год
		if attrName == attrNameYear && maxValue >= float64(time.Now().Year()-1) {
			maxValue = float64(time.Now().Year() + 1)
		}

		// Устанавливаем шаг
		var step float64
		switch attrName {
		case "engine_capacity":
			step = 0.1
		case "area", "land_area":
			step = 0.5
		default:
			step = 1.0
		}

		ranges[attrName] = map[string]interface{}{
			"min":   minValue,
			"max":   maxValue,
			"step":  step,
			"count": valueCount,
		}

		log.Printf("Attribute %s range: min=%.2f, max=%.2f, values=%d",
			attrName, minValue, maxValue, valueCount)
	}

	return ranges, nil
}

// fillDefaultRanges заполняет отсутствующие атрибуты значениями по умолчанию
func (s *Storage) fillDefaultRanges(ranges map[string]map[string]interface{}) {
	defaultRanges := map[string]map[string]interface{}{
		attrNameYear:       {"min": float64(time.Now().Year() - 30), "max": float64(time.Now().Year() + 1), "step": 1.0},
		attrNameMileage:    {"min": 0.0, "max": 500000.0, "step": 1000.0},
		"engine_capacity":  {"min": 0.5, "max": 8.0, "step": 0.1},
		attrNamePower:      {"min": 50.0, "max": 500.0, "step": 10.0},
		"rooms":            {"min": 1.0, "max": 10.0, "step": 1.0},
		"floor":            {"min": 1.0, "max": 25.0, "step": 1.0},
		"total_floors":     {"min": 1.0, "max": 30.0, "step": 1.0},
		"area":             {"min": 10.0, "max": 300.0, "step": 0.5},
		"land_area":        {"min": 1.0, "max": 100.0, "step": 0.5},
	}

	for attr, defaultRange := range defaultRanges {
		if _, exists := ranges[attr]; !exists {
			ranges[attr] = defaultRange
			log.Printf("No data for attribute %s, using defaults: min=%.2f, max=%.2f",
				attr, defaultRange["min"], defaultRange["max"])
		}
	}
}

// InvalidateAttributesCache очищает кэш атрибутов для категории
func (s *Storage) InvalidateAttributesCache(categoryID int) {
	s.attributeCacheMutex.Lock()
	delete(s.attributeCache, categoryID)
	delete(s.attributeCacheTime, categoryID)
	s.attributeCacheMutex.Unlock()

	s.rangesCacheMutex.Lock()
	delete(s.rangesCache, categoryID)
	delete(s.rangesCacheTime, categoryID)
	s.rangesCacheMutex.Unlock()

	log.Printf("Invalidated attributes cache for category %d", categoryID)
}
