package service

import (
	"backend/internal/domain/models"
	"context"
	"database/sql"
	"fmt"
	"time"
)

// CreateAttribute создает новый атрибут категории
func (s *MarketplaceService) CreateAttribute(ctx context.Context, attribute *models.CategoryAttribute) (int, error) {
	// Преобразуем Options в JSON, если они представлены как структура
	var optionsJSON []byte
	var err error
	if attribute.Options != nil {
		optionsJSON = attribute.Options
	}

	// Преобразуем ValidRules в JSON, если они представлены как структура
	var validRulesJSON []byte
	if attribute.ValidRules != nil {
		validRulesJSON = attribute.ValidRules
	}

	// Создаем атрибут в БД
	query := `
		INSERT INTO category_attributes (
			name, display_name, attribute_type, icon, options, validation_rules, 
			is_searchable, is_filterable, is_required, sort_order, custom_component
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, created_at
	`

	var id int
	var createdAt time.Time

	err = s.storage.QueryRow(
		ctx, query,
		attribute.Name,
		attribute.DisplayName,
		attribute.AttributeType,
		attribute.Icon,
		optionsJSON,
		validRulesJSON,
		attribute.IsSearchable,
		attribute.IsFilterable,
		attribute.IsRequired,
		attribute.SortOrder,
		attribute.CustomComponent,
	).Scan(&id, &createdAt)

	if err != nil {
		return 0, fmt.Errorf("не удалось создать атрибут: %w", err)
	}

	// Устанавливаем ID и CreatedAt в структуре
	attribute.ID = id
	attribute.CreatedAt = createdAt

	// Сохраняем переводы для атрибута
	if attribute.Translations != nil && len(attribute.Translations) > 0 {
		for lang, text := range attribute.Translations {
			translation := &models.Translation{
				EntityType:     "attribute",
				EntityID:       id,
				Language:       lang,
				FieldName:      "display_name",
				TranslatedText: text,
				IsVerified:     true,
			}
			if err := s.UpdateTranslation(ctx, translation); err != nil {
				return id, fmt.Errorf("не удалось сохранить перевод для %s: %w", lang, err)
			}
		}
	}

	// Сохраняем переводы для опций атрибута
	if attribute.OptionTranslations != nil && len(attribute.OptionTranslations) > 0 {
		for lang, options := range attribute.OptionTranslations {
			for optionKey, optionValue := range options {
				translation := &models.Translation{
					EntityType:     "attribute_option",
					EntityID:       id,
					Language:       lang,
					FieldName:      optionKey,
					TranslatedText: optionValue,
					IsVerified:     true,
				}
				if err := s.UpdateTranslation(ctx, translation); err != nil {
					return id, fmt.Errorf("не удалось сохранить перевод опции для %s: %w", lang, err)
				}
			}
		}
	}

	return id, nil
}

// UpdateAttribute обновляет существующий атрибут
func (s *MarketplaceService) UpdateAttribute(ctx context.Context, attribute *models.CategoryAttribute) error {
	// Преобразуем Options в JSON, если они представлены как структура
	var optionsJSON []byte
	var err error
	if attribute.Options != nil {
		optionsJSON = attribute.Options
	}

	// Преобразуем ValidRules в JSON, если они представлены как структура
	var validRulesJSON []byte
	if attribute.ValidRules != nil {
		validRulesJSON = attribute.ValidRules
	}

	// Обновляем атрибут в БД
	query := `
		UPDATE category_attributes
		SET 
			name = $1, 
			display_name = $2, 
			attribute_type = $3, 
			icon = $4,
			options = $5, 
			validation_rules = $6, 
			is_searchable = $7, 
			is_filterable = $8, 
			is_required = $9, 
			sort_order = $10,
			custom_component = $11
		WHERE id = $12
	`

	_, err = s.storage.Exec(
		ctx, query,
		attribute.Name,
		attribute.DisplayName,
		attribute.AttributeType,
		attribute.Icon,
		optionsJSON,
		validRulesJSON,
		attribute.IsSearchable,
		attribute.IsFilterable,
		attribute.IsRequired,
		attribute.SortOrder,
		attribute.CustomComponent,
		attribute.ID,
	)

	if err != nil {
		return fmt.Errorf("не удалось обновить атрибут: %w", err)
	}

	// Обновляем переводы для атрибута
	if attribute.Translations != nil && len(attribute.Translations) > 0 {
		for lang, text := range attribute.Translations {
			translation := &models.Translation{
				EntityType:     "attribute",
				EntityID:       attribute.ID,
				Language:       lang,
				FieldName:      "display_name",
				TranslatedText: text,
				IsVerified:     true,
			}
			if err := s.UpdateTranslation(ctx, translation); err != nil {
				return fmt.Errorf("не удалось обновить перевод для %s: %w", lang, err)
			}
		}
	}

	// Обновляем переводы для опций атрибута
	if attribute.OptionTranslations != nil && len(attribute.OptionTranslations) > 0 {
		for lang, options := range attribute.OptionTranslations {
			for optionKey, optionValue := range options {
				translation := &models.Translation{
					EntityType:     "attribute_option",
					EntityID:       attribute.ID,
					Language:       lang,
					FieldName:      optionKey,
					TranslatedText: optionValue,
					IsVerified:     true,
				}
				if err := s.UpdateTranslation(ctx, translation); err != nil {
					return fmt.Errorf("не удалось обновить перевод опции для %s: %w", lang, err)
				}
			}
		}
	}

	// Инвалидируем кеш атрибутов для всех категорий, связанных с этим атрибутом
	query = `
		SELECT DISTINCT category_id 
		FROM category_attribute_mapping 
		WHERE attribute_id = $1
	`
	rows, err := s.storage.Query(ctx, query, attribute.ID)
	if err != nil {
		return fmt.Errorf("не удалось получить связанные категории: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var categoryID int
		if err := rows.Scan(&categoryID); err != nil {
			return fmt.Errorf("не удалось прочитать ID категории: %w", err)
		}

		// Инвалидируем кеш для каждой категории
		if err := s.InvalidateAttributeCache(ctx, categoryID); err != nil {
			return fmt.Errorf("не удалось инвалидировать кеш для категории %d: %w", categoryID, err)
		}
	}

	return nil
}

// DeleteAttribute удаляет атрибут по ID
func (s *MarketplaceService) DeleteAttribute(ctx context.Context, id int) error {
	// Проверяем, используется ли атрибут в объявлениях
	var count int
	err := s.storage.QueryRow(ctx, "SELECT COUNT(*) FROM listing_attribute_values WHERE attribute_id = $1", id).Scan(&count)
	if err != nil {
		return fmt.Errorf("не удалось проверить использование атрибута: %w", err)
	}

	if count > 0 {
		return fmt.Errorf("атрибут используется в %d объявлениях и не может быть удален", count)
	}

	// Получаем список категорий, связанных с атрибутом
	query := `
		SELECT DISTINCT category_id 
		FROM category_attribute_mapping 
		WHERE attribute_id = $1
	`
	rows, err := s.storage.Query(ctx, query, id)
	if err != nil {
		return fmt.Errorf("не удалось получить связанные категории: %w", err)
	}

	categoryIDs := make([]int, 0)
	for rows.Next() {
		var categoryID int
		if err := rows.Scan(&categoryID); err != nil {
			rows.Close()
			return fmt.Errorf("не удалось прочитать ID категории: %w", err)
		}
		categoryIDs = append(categoryIDs, categoryID)
	}
	rows.Close()

	// Начинаем транзакцию
	tx, err := s.storage.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("не удалось начать транзакцию: %w", err)
	}
	defer tx.Rollback()

	// Удаляем связи с категориями
	_, err = tx.Exec(ctx, "DELETE FROM category_attribute_mapping WHERE attribute_id = $1", id)
	if err != nil {
		return fmt.Errorf("не удалось удалить связи с категориями: %w", err)
	}

	// Удаляем переводы атрибута
	_, err = tx.Exec(ctx, "DELETE FROM translations WHERE (entity_type = 'attribute' OR entity_type = 'attribute_option') AND entity_id = $1", id)
	if err != nil {
		return fmt.Errorf("не удалось удалить переводы: %w", err)
	}

	// Удаляем сам атрибут
	_, err = tx.Exec(ctx, "DELETE FROM category_attributes WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("не удалось удалить атрибут: %w", err)
	}

	// Подтверждаем транзакцию
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("не удалось завершить транзакцию: %w", err)
	}

	// Инвалидируем кеш для всех затронутых категорий
	for _, categoryID := range categoryIDs {
		if err := s.InvalidateAttributeCache(ctx, categoryID); err != nil {
			return fmt.Errorf("не удалось инвалидировать кеш для категории %d: %w", categoryID, err)
		}
	}

	return nil
}

// GetAttributeByID получает атрибут по ID
func (s *MarketplaceService) GetAttributeByID(ctx context.Context, id int) (*models.CategoryAttribute, error) {
	query := `
		SELECT 
			id, name, display_name, attribute_type, icon, options, validation_rules, 
			is_searchable, is_filterable, is_required, sort_order, created_at, custom_component
		FROM category_attributes
		WHERE id = $1
	`

	var attribute models.CategoryAttribute
	var optionsJSON, validRulesJSON []byte

	err := s.storage.QueryRow(ctx, query, id).Scan(
		&attribute.ID,
		&attribute.Name,
		&attribute.DisplayName,
		&attribute.AttributeType,
		&attribute.Icon,
		&optionsJSON,
		&validRulesJSON,
		&attribute.IsSearchable,
		&attribute.IsFilterable,
		&attribute.IsRequired,
		&attribute.SortOrder,
		&attribute.CreatedAt,
		&attribute.CustomComponent,
	)

	if err != nil {
		return nil, fmt.Errorf("не удалось получить атрибут: %w", err)
	}

	// Устанавливаем Options и ValidRules как json.RawMessage
	attribute.Options = optionsJSON
	attribute.ValidRules = validRulesJSON

	// Получаем переводы для атрибута
	translationsQuery := `
		SELECT language, field_name, translated_text
		FROM translations
		WHERE entity_type = 'attribute' AND entity_id = $1 AND field_name = 'display_name'
	`
	rows, err := s.storage.Query(ctx, translationsQuery, id)
	if err != nil {
		return &attribute, fmt.Errorf("не удалось получить переводы: %w", err)
	}
	defer rows.Close()

	attribute.Translations = make(map[string]string)
	for rows.Next() {
		var lang, field, text string
		if err := rows.Scan(&lang, &field, &text); err != nil {
			return &attribute, fmt.Errorf("не удалось прочитать перевод: %w", err)
		}
		attribute.Translations[lang] = text
	}

	// Получаем переводы для опций атрибута
	optionTranslationsQuery := `
		SELECT language, field_name, translated_text
		FROM translations
		WHERE entity_type = 'attribute_option' AND entity_id = $1
	`
	rows, err = s.storage.Query(ctx, optionTranslationsQuery, id)
	if err != nil {
		return &attribute, fmt.Errorf("не удалось получить переводы опций: %w", err)
	}
	defer rows.Close()

	attribute.OptionTranslations = make(map[string]map[string]string)
	for rows.Next() {
		var lang, option, text string
		if err := rows.Scan(&lang, &option, &text); err != nil {
			return &attribute, fmt.Errorf("не удалось прочитать перевод опции: %w", err)
		}

		if attribute.OptionTranslations[lang] == nil {
			attribute.OptionTranslations[lang] = make(map[string]string)
		}
		attribute.OptionTranslations[lang][option] = text
	}

	return &attribute, nil
}

// AddAttributeToCategory привязывает атрибут к категории
func (s *MarketplaceService) AddAttributeToCategory(ctx context.Context, categoryID int, attributeID int, isRequired bool) error {
	// Используем новый метод с sortOrder=0 (будет использовано значение из атрибута)
	return s.AddAttributeToCategoryWithOrder(ctx, categoryID, attributeID, isRequired, 0)
}

// AddAttributeToCategoryWithOrder привязывает атрибут к категории с указанием порядка сортировки
func (s *MarketplaceService) AddAttributeToCategoryWithOrder(ctx context.Context, categoryID int, attributeID int, isRequired bool, sortOrder int) error {
	// Проверяем, что категория и атрибут существуют
	var categoryExists, attributeExists bool

	err := s.storage.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM marketplace_categories WHERE id = $1)", categoryID).Scan(&categoryExists)
	if err != nil {
		return fmt.Errorf("не удалось проверить существование категории: %w", err)
	}

	if !categoryExists {
		return fmt.Errorf("категория с ID %d не существует", categoryID)
	}

	err = s.storage.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM category_attributes WHERE id = $1)", attributeID).Scan(&attributeExists)
	if err != nil {
		return fmt.Errorf("не удалось проверить существование атрибута: %w", err)
	}

	if !attributeExists {
		return fmt.Errorf("атрибут с ID %d не существует", attributeID)
	}

	// Добавляем связь с учетом sort_order и custom_component
	_, err = s.storage.Exec(ctx, `
		INSERT INTO category_attribute_mapping (category_id, attribute_id, is_enabled, is_required, sort_order, custom_component)
		VALUES ($1, $2, true, $3, $4, NULL)
		ON CONFLICT (category_id, attribute_id)
		DO UPDATE SET is_enabled = true, is_required = $3, sort_order = $4
	`, categoryID, attributeID, isRequired, sortOrder)

	if err != nil {
		return fmt.Errorf("не удалось привязать атрибут к категории: %w", err)
	}

	// Инвалидируем кеш атрибутов для категории
	return s.InvalidateAttributeCache(ctx, categoryID)
}

// RemoveAttributeFromCategory отвязывает атрибут от категории
func (s *MarketplaceService) RemoveAttributeFromCategory(ctx context.Context, categoryID int, attributeID int) error {
	// Проверяем, есть ли объявления, использующие этот атрибут в данной категории
	var count int
	err := s.storage.QueryRow(ctx, `
		SELECT COUNT(*) FROM listing_attribute_values lav
		JOIN marketplace_listings ml ON lav.listing_id = ml.id
		WHERE ml.category_id = $1 AND lav.attribute_id = $2
	`, categoryID, attributeID).Scan(&count)

	if err != nil {
		return fmt.Errorf("не удалось проверить использование атрибута: %w", err)
	}

	if count > 0 {
		return fmt.Errorf("атрибут используется в %d объявлениях в данной категории и не может быть отвязан", count)
	}

	// Удаляем связь
	_, err = s.storage.Exec(ctx, `
		DELETE FROM category_attribute_mapping
		WHERE category_id = $1 AND attribute_id = $2
	`, categoryID, attributeID)

	if err != nil {
		return fmt.Errorf("не удалось отвязать атрибут от категории: %w", err)
	}

	// Инвалидируем кеш атрибутов для категории
	return s.InvalidateAttributeCache(ctx, categoryID)
}

// UpdateAttributeCategory обновляет настройки связи атрибута с категорией
func (s *MarketplaceService) UpdateAttributeCategory(ctx context.Context, categoryID int, attributeID int, isRequired bool, isEnabled bool) error {
	// Обновляем связь
	_, err := s.storage.Exec(ctx, `
		UPDATE category_attribute_mapping
		SET is_required = $1, is_enabled = $2
		WHERE category_id = $3 AND attribute_id = $4
	`, isRequired, isEnabled, categoryID, attributeID)

	if err != nil {
		return fmt.Errorf("не удалось обновить связь атрибута с категорией: %w", err)
	}

	// Инвалидируем кеш атрибутов для категории
	return s.InvalidateAttributeCache(ctx, categoryID)
}

// UpdateAttributeCategoryExtended обновляет расширенные настройки связи атрибута с категорией
func (s *MarketplaceService) UpdateAttributeCategoryExtended(
	ctx context.Context,
	categoryID int,
	attributeID int,
	isRequired bool,
	isEnabled bool,
	sortOrder int,
	customComponent string,
) error {
	// Обновляем связь с дополнительными полями
	_, err := s.storage.Exec(ctx, `
		UPDATE category_attribute_mapping
		SET is_required = $1, is_enabled = $2, sort_order = $3, custom_component = $4
		WHERE category_id = $5 AND attribute_id = $6
	`, isRequired, isEnabled, sortOrder, customComponent, categoryID, attributeID)

	if err != nil {
		return fmt.Errorf("не удалось обновить связь атрибута с категорией: %w", err)
	}

	// Инвалидируем кеш атрибутов для категории
	return s.InvalidateAttributeCache(ctx, categoryID)
}

// InvalidateAttributeCache инвалидирует кеш атрибутов для указанной категории
func (s *MarketplaceService) InvalidateAttributeCache(ctx context.Context, categoryID int) error {
	// В текущей реализации кеш не используется,
	// но метод оставлен для будущих расширений
	return nil
}

// GetCategoryByID получает информацию о категории по ID
func (s *MarketplaceService) GetCategoryByID(ctx context.Context, id int) (*models.MarketplaceCategory, error) {
	query := `
		SELECT id, name, slug, parent_id, icon, created_at, has_custom_ui, custom_ui_component, 
                       0 as listing_count, sort_order, COALESCE(level, 0) as level, COALESCE(count, 0) as count, 
                       COALESCE(external_id, '') as external_id
		FROM marketplace_categories
		WHERE id = $1
	`

	var category models.MarketplaceCategory
	err := s.storage.QueryRow(ctx, query, id).Scan(
		&category.ID,
		&category.Name,
		&category.Slug,
		&category.ParentID,
		&category.Icon,
		&category.CreatedAt,
		&category.HasCustomUI,
		&category.CustomUIComponent,
		&category.ListingCount,
		&category.SortOrder,
		&category.Level,
		&category.Count,
		&category.ExternalID,
	)

	if err != nil {
		return nil, fmt.Errorf("не удалось получить категорию: %w", err)
	}

	// Получаем переводы для категории
	translationsQuery := `
		SELECT language, field_name, translated_text
		FROM translations
		WHERE entity_type = 'category' AND entity_id = $1 AND field_name = 'name'
	`
	rows, err := s.storage.Query(ctx, translationsQuery, id)
	if err != nil {
		return &category, fmt.Errorf("не удалось получить переводы: %w", err)
	}
	defer rows.Close()

	category.Translations = make(map[string]string)
	for rows.Next() {
		var lang, field, text string
		if err := rows.Scan(&lang, &field, &text); err != nil {
			return &category, fmt.Errorf("не удалось прочитать перевод: %w", err)
		}
		category.Translations[lang] = text
	}

	return &category, nil
}

// GetCategoryAttributes получает все атрибуты для указанной категории
func (s *MarketplaceService) GetCategoryAttributes(ctx context.Context, categoryID int) ([]models.CategoryAttribute, error) {
	// Обновленный запрос, который получает атрибуты из обоих источников:
	// 1. Из прямого маппинга (category_attribute_mapping)
	// 2. Из групп атрибутов (через category_attribute_groups -> attribute_group_items)
	query := `
		WITH combined_attributes AS (
			-- Атрибуты из прямого маппинга
			SELECT 
				a.id, a.name, a.display_name, a.attribute_type, a.options, a.validation_rules,
				a.is_searchable, a.is_filterable, cam.is_required, a.sort_order, a.created_at,
				COALESCE(cam.custom_component, a.custom_component) as custom_component,
				cam.sort_order as effective_sort_order,
				'direct' as source
			FROM category_attributes a
			JOIN category_attribute_mapping cam ON a.id = cam.attribute_id
			WHERE cam.category_id = $1 AND cam.is_enabled = true
			
			UNION
			
			-- Атрибуты из групп
			SELECT 
				a.id, a.name, 
				COALESCE(agi.custom_display_name, a.display_name) as display_name,
				a.attribute_type, a.options, a.validation_rules,
				a.is_searchable, a.is_filterable, a.is_required, a.sort_order, a.created_at,
				a.custom_component,
				-- Сортировка для групп: сначала по группе, потом по атрибуту внутри группы
				(cag.sort_order * 1000 + agi.sort_order) as effective_sort_order,
				'group' as source
			FROM category_attributes a
			JOIN attribute_group_items agi ON a.id = agi.attribute_id
			JOIN attribute_groups ag ON agi.group_id = ag.id
			JOIN category_attribute_groups cag ON ag.id = cag.group_id
			WHERE cag.category_id = $1 
				AND cag.is_active = true 
				AND ag.is_active = true
		)
		SELECT DISTINCT id, name, display_name, attribute_type, options, validation_rules,
			is_searchable, is_filterable, is_required, sort_order, created_at, 
			custom_component, MIN(effective_sort_order) as final_sort_order
		FROM combined_attributes
		GROUP BY id, name, display_name, attribute_type, options, validation_rules,
			is_searchable, is_filterable, is_required, sort_order, created_at, custom_component
		ORDER BY final_sort_order, id
	`

	rows, err := s.storage.Query(ctx, query, categoryID)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить атрибуты категории: %w", err)
	}
	defer rows.Close()

	attributes := make([]models.CategoryAttribute, 0)
	for rows.Next() {
		var attribute models.CategoryAttribute
		var optionsJSON, validRulesJSON []byte
		var finalSortOrder int
		var customComponent sql.NullString

		err := rows.Scan(
			&attribute.ID,
			&attribute.Name,
			&attribute.DisplayName,
			&attribute.AttributeType,
			&optionsJSON,
			&validRulesJSON,
			&attribute.IsSearchable,
			&attribute.IsFilterable,
			&attribute.IsRequired,
			&attribute.SortOrder,
			&attribute.CreatedAt,
			&customComponent,
			&finalSortOrder,
		)
		if err != nil {
			return nil, fmt.Errorf("не удалось прочитать атрибут: %w", err)
		}

		// Используем финальный sort_order из запроса
		attribute.SortOrder = finalSortOrder

		// Устанавливаем CustomComponent, обрабатывая NULL-значения
		attribute.CustomComponent = ""
		if customComponent.Valid {
			attribute.CustomComponent = customComponent.String
		}

		// Устанавливаем Options и ValidRules
		attribute.Options = optionsJSON
		attribute.ValidRules = validRulesJSON

		// Получаем переводы для атрибута
		translationsQuery := `
			SELECT language, field_name, translated_text
			FROM translations
			WHERE entity_type = 'attribute' AND entity_id = $1 AND field_name = 'display_name'
		`
		tRows, err := s.storage.Query(ctx, translationsQuery, attribute.ID)
		if err == nil {
			attribute.Translations = make(map[string]string)
			for tRows.Next() {
				var lang, field, text string
				if err := tRows.Scan(&lang, &field, &text); err == nil {
					attribute.Translations[lang] = text
				}
			}
			tRows.Close()
		}

		// Получаем переводы для опций атрибута
		optionTranslationsQuery := `
			SELECT language, field_name, translated_text
			FROM translations
			WHERE entity_type = 'attribute_option' AND entity_id = $1
		`
		oRows, err := s.storage.Query(ctx, optionTranslationsQuery, attribute.ID)
		if err == nil {
			attribute.OptionTranslations = make(map[string]map[string]string)
			for oRows.Next() {
				var lang, option, text string
				if err := oRows.Scan(&lang, &option, &text); err == nil {
					if attribute.OptionTranslations[lang] == nil {
						attribute.OptionTranslations[lang] = make(map[string]string)
					}
					attribute.OptionTranslations[lang][option] = text
				}
			}
			oRows.Close()
		}

		attributes = append(attributes, attribute)
	}

	return attributes, nil
}

// GetAttributeRanges получает минимальные и максимальные значения числовых атрибутов для категории
func (s *MarketplaceService) GetAttributeRanges(ctx context.Context, categoryID int) (map[string]map[string]interface{}, error) {
	// Получаем все атрибуты категории
	attributes, err := s.GetCategoryAttributes(ctx, categoryID)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить атрибуты категории: %w", err)
	}

	// Фильтруем числовые атрибуты и атрибуты с диапазоном
	numericAttributeIDs := make([]int, 0)
	numericAttributeNames := make(map[int]string)
	for _, attr := range attributes {
		if attr.AttributeType == "number" || attr.AttributeType == "range" {
			numericAttributeIDs = append(numericAttributeIDs, attr.ID)
			numericAttributeNames[attr.ID] = attr.Name
		}
	}

	if len(numericAttributeIDs) == 0 {
		return make(map[string]map[string]interface{}), nil
	}

	// Создаем результат
	result := make(map[string]map[string]interface{})

	// Вычисляем минимальные и максимальные значения для каждого числового атрибута
	for _, attrID := range numericAttributeIDs {
		attrName := numericAttributeNames[attrID]

		// Запрос для получения минимального и максимального значения
		query := `
			SELECT MIN(numeric_value), MAX(numeric_value)
			FROM listing_attribute_values
			JOIN marketplace_listings ON listing_attribute_values.listing_id = marketplace_listings.id
			WHERE attribute_id = $1 AND marketplace_listings.category_id = $2 AND marketplace_listings.status = 'active'
		`

		var min, max *float64
		err := s.storage.QueryRow(ctx, query, attrID, categoryID).Scan(&min, &max)
		if err != nil {
			// Если значения не найдены, просто пропускаем этот атрибут
			continue
		}

		if min != nil && max != nil {
			result[attrName] = map[string]interface{}{
				"min": *min,
				"max": *max,
			}
		}
	}

	return result, nil
}

// SaveListingAttributes сохраняет значения атрибутов для объявления
func (s *MarketplaceService) SaveListingAttributes(ctx context.Context, listingID int, attributes []models.ListingAttributeValue) error {
	// Начинаем транзакцию
	tx, err := s.storage.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("не удалось начать транзакцию: %w", err)
	}
	defer tx.Rollback()

	// Удаляем существующие атрибуты для данного объявления
	_, err = tx.Exec(ctx, "DELETE FROM listing_attribute_values WHERE listing_id = $1", listingID)
	if err != nil {
		return fmt.Errorf("не удалось удалить существующие атрибуты: %w", err)
	}

	// Добавляем новые атрибуты
	for _, attr := range attributes {
		var valueType string
		var textValue *string
		var numValue *float64
		var boolValue *bool
		var jsonValue []byte

		// Определяем тип значения и устанавливаем соответствующее поле
		switch attr.AttributeType {
		case "text", "textarea", "select", "multiselect":
			valueType = "text"
			textValue = attr.TextValue
		case "number", "range":
			valueType = "numeric"
			numValue = attr.NumericValue
		case "boolean", "checkbox":
			valueType = "boolean"
			boolValue = attr.BooleanValue
		case "json", "complex":
			valueType = "json"
			jsonValue = attr.JSONValue
		default:
			valueType = "text"
			textValue = attr.TextValue
		}

		// Добавляем запись в базу данных
		_, err = tx.Exec(ctx, `
			INSERT INTO listing_attribute_values (
				listing_id, attribute_id, value_type, text_value, numeric_value, boolean_value, json_value, unit
			)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		`, listingID, attr.AttributeID, valueType, textValue, numValue, boolValue, jsonValue, attr.Unit)

		if err != nil {
			return fmt.Errorf("не удалось сохранить значение атрибута %d: %w", attr.AttributeID, err)
		}
	}

	// Фиксируем транзакцию
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("не удалось зафиксировать транзакцию: %w", err)
	}

	return nil
}
