package service

import (
	"backend/internal/domain/models"
	"context"
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
			name, display_name, attribute_type, options, validation_rules, 
			is_searchable, is_filterable, is_required, sort_order, custom_component
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, created_at
	`

	var id int
	var createdAt time.Time

	err = s.storage.QueryRow(
		ctx, query,
		attribute.Name,
		attribute.DisplayName,
		attribute.AttributeType,
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
			options = $4, 
			validation_rules = $5, 
			is_searchable = $6, 
			is_filterable = $7, 
			is_required = $8, 
			sort_order = $9,
			custom_component = $10
		WHERE id = $11
	`

	_, err = s.storage.Exec(
		ctx, query,
		attribute.Name,
		attribute.DisplayName,
		attribute.AttributeType,
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
			id, name, display_name, attribute_type, options, validation_rules, 
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

	// Добавляем связь
	_, err = s.storage.Exec(ctx, `
		INSERT INTO category_attribute_mapping (category_id, attribute_id, is_enabled, is_required)
		VALUES ($1, $2, true, $3)
		ON CONFLICT (category_id, attribute_id)
		DO UPDATE SET is_enabled = true, is_required = $3
	`, categoryID, attributeID, isRequired)

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

// InvalidateAttributeCache инвалидирует кеш атрибутов для указанной категории
func (s *MarketplaceService) InvalidateAttributeCache(ctx context.Context, categoryID int) error {
	// В текущей реализации кеш не используется,
	// но метод оставлен для будущих расширений
	return nil
}
