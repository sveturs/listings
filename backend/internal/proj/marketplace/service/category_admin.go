package service

import (
	"context"
	"fmt"
	"time"
	"unicode"

	"backend/internal/cache"
	"backend/internal/domain/models"
)

// Реализации методов для администрирования категорий

// CreateCategory создает новую категорию
func (s *MarketplaceService) CreateCategory(ctx context.Context, category *models.MarketplaceCategory) (int, error) {
	// Создаем категорию в БД
	query := `
		INSERT INTO marketplace_categories (name, slug, parent_id, icon, has_custom_ui, custom_ui_component, description, is_active, seo_title, seo_description, seo_keywords)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id
	`
	var id int
	err := s.storage.QueryRow(ctx, query,
		category.Name,
		category.Slug,
		category.ParentID,
		category.Icon,
		category.HasCustomUI,
		category.CustomUIComponent,
		category.Description,
		category.IsActive,
		category.SEOTitle,
		category.SEODescription,
		category.SEOKeywords,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("не удалось создать категорию: %w", err)
	}

	// Если есть переводы, сохраняем их
	if len(category.Translations) > 0 {
		for lang, text := range category.Translations {
			translation := &models.Translation{
				EntityType:     "category",
				EntityID:       id,
				Language:       lang,
				FieldName:      "name",
				TranslatedText: text,
				IsVerified:     true,
			}
			if err := s.UpdateTranslation(ctx, translation); err != nil {
				return id, fmt.Errorf("не удалось сохранить перевод для %s: %w", lang, err)
			}
		}
	} else {
		// Если переводы не предоставлены, создаем автоматические переводы
		languages := []string{"en", "ru", "sr"}
		for _, targetLang := range languages {
			// Пропускаем, если название уже на целевом языке (простая эвристика)
			if targetLang == "en" && isLikelyEnglish(category.Name) {
				continue
			}
			if targetLang == "ru" && isLikelyCyrillic(category.Name) {
				continue
			}

			// Переводим название категории
			translatedText, err := s.TranslateText(ctx, category.Name, "auto", targetLang)
			if err != nil {
				// Логируем ошибку, но не прерываем создание категории
				fmt.Printf("Не удалось перевести на %s: %v\n", targetLang, err)
				continue
			}

			translation := &models.Translation{
				EntityType:          "category",
				EntityID:            id,
				Language:            targetLang,
				FieldName:           "name",
				TranslatedText:      translatedText,
				IsMachineTranslated: true,
				IsVerified:          false,
			}
			if err := s.UpdateTranslation(ctx, translation); err != nil {
				fmt.Printf("Не удалось сохранить перевод для %s: %v\n", targetLang, err)
			}
		}
	}

	// Обновляем материализованное представление для обновления счетчиков
	_ = s.RefreshCategoryListingCounts(ctx)

	// Инвалидируем кеш категорий
	if s.cache != nil {
		_ = s.cache.DeletePattern(ctx, cache.BuildAllCategoriesInvalidationPattern())
	}

	return id, nil
}

// UpdateCategory обновляет существующую категорию
func (s *MarketplaceService) UpdateCategory(ctx context.Context, category *models.MarketplaceCategory) error {
	// Обновляем категорию в БД
	query := `
		UPDATE marketplace_categories
		SET name = $1, slug = $2, parent_id = $3, icon = $4, has_custom_ui = $5, custom_ui_component = $6, description = $7, is_active = $8, seo_title = $9, seo_description = $10, seo_keywords = $11
		WHERE id = $12
	`
	_, err := s.storage.Exec(ctx, query,
		category.Name,
		category.Slug,
		category.ParentID,
		category.Icon,
		category.HasCustomUI,
		category.CustomUIComponent,
		category.Description,
		category.IsActive,
		category.SEOTitle,
		category.SEODescription,
		category.SEOKeywords,
		category.ID,
	)
	if err != nil {
		return fmt.Errorf("не удалось обновить категорию: %w", err)
	}

	// Если есть переводы, обновляем их
	if len(category.Translations) > 0 {
		for lang, text := range category.Translations {
			translation := &models.Translation{
				EntityType:     "category",
				EntityID:       category.ID,
				Language:       lang,
				FieldName:      "name",
				TranslatedText: text,
				IsVerified:     true,
			}
			if err := s.UpdateTranslation(ctx, translation); err != nil {
				return fmt.Errorf("не удалось обновить перевод для %s: %w", lang, err)
			}
		}
	}

	// Инвалидируем кеш категорий и конкретной категории
	if s.cache != nil {
		_ = s.cache.DeletePattern(ctx, cache.BuildAllCategoriesInvalidationPattern())
		_ = s.cache.Delete(ctx, cache.BuildCategoryKey(int64(category.ID)))
		_ = s.cache.DeletePattern(ctx, cache.BuildCategoryInvalidationPattern(int64(category.ID)))
	}

	return nil
}

// DeleteCategory удаляет категорию по ID
func (s *MarketplaceService) DeleteCategory(ctx context.Context, id int) error {
	// Проверяем наличие объявлений в категории
	var listingCount int
	err := s.storage.QueryRow(ctx, "SELECT COUNT(*) FROM marketplace_listings WHERE category_id = $1", id).Scan(&listingCount)
	if err != nil {
		return fmt.Errorf("не удалось проверить наличие объявлений: %w", err)
	}

	if listingCount > 0 {
		return fmt.Errorf("категория содержит %d объявлений и не может быть удалена", listingCount)
	}

	// Проверяем наличие дочерних категорий
	var childCount int
	err = s.storage.QueryRow(ctx, "SELECT COUNT(*) FROM marketplace_categories WHERE parent_id = $1", id).Scan(&childCount)
	if err != nil {
		return fmt.Errorf("не удалось проверить наличие дочерних категорий: %w", err)
	}

	if childCount > 0 {
		return fmt.Errorf("категория содержит %d дочерних категорий и не может быть удалена", childCount)
	}

	// Удаляем связи с атрибутами
	_, err = s.storage.Exec(ctx, "DELETE FROM category_attribute_mapping WHERE category_id = $1", id)
	if err != nil {
		return fmt.Errorf("не удалось удалить связи с атрибутами: %w", err)
	}

	// Удаляем переводы категории
	_, err = s.storage.Exec(ctx, "DELETE FROM translations WHERE entity_type = 'category' AND entity_id = $1", id)
	if err != nil {
		return fmt.Errorf("не удалось удалить переводы: %w", err)
	}

	// Удаляем саму категорию
	_, err = s.storage.Exec(ctx, "DELETE FROM marketplace_categories WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("не удалось удалить категорию: %w", err)
	}

	// Обновляем материализованное представление для обновления счетчиков
	_ = s.RefreshCategoryListingCounts(ctx)

	// Инвалидируем кеш категорий
	if s.cache != nil {
		_ = s.cache.DeletePattern(ctx, cache.BuildAllCategoriesInvalidationPattern())
		_ = s.cache.Delete(ctx, cache.BuildCategoryKey(int64(id)))
		_ = s.cache.DeletePattern(ctx, cache.BuildCategoryInvalidationPattern(int64(id)))
	}

	return nil
}

// ReorderCategories изменяет порядок категорий
func (s *MarketplaceService) ReorderCategories(ctx context.Context, orderedIDs []int) error {
	// Начинаем транзакцию для атомарной операции
	tx, err := s.storage.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("не удалось начать транзакцию: %w", err)
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			fmt.Printf("Failed to rollback transaction: %v", err)
		}
	}()

	// Обновляем порядок для каждой категории
	for i, id := range orderedIDs {
		_, err = tx.Exec(ctx, "UPDATE marketplace_categories SET sort_order = $1 WHERE id = $2", i, id)
		if err != nil {
			return fmt.Errorf("не удалось обновить порядок для категории %d: %w", id, err)
		}
	}

	// Подтверждаем транзакцию
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("не удалось завершить транзакцию: %w", err)
	}

	return nil
}

// MoveCategory перемещает категорию в иерархии
func (s *MarketplaceService) MoveCategory(ctx context.Context, id int, newParentID int) error {
	// Проверяем, что категории существуют
	var count int
	err := s.storage.QueryRow(ctx, "SELECT COUNT(*) FROM marketplace_categories WHERE id IN ($1, $2)", id, newParentID).Scan(&count)
	if err != nil {
		return fmt.Errorf("не удалось проверить наличие категорий: %w", err)
	}

	if count < 2 {
		return fmt.Errorf("одна из категорий не существует")
	}

	// Проверяем, что новый родитель не является потомком перемещаемой категории
	var isDescendant bool
	err = s.storage.QueryRow(ctx, `
		WITH RECURSIVE category_tree AS (
			SELECT id, parent_id FROM marketplace_categories WHERE id = $1
			UNION ALL
			SELECT c.id, c.parent_id FROM marketplace_categories c
			JOIN category_tree ct ON c.parent_id = ct.id
		)
		SELECT EXISTS(SELECT 1 FROM category_tree WHERE id = $2)
	`, id, newParentID).Scan(&isDescendant)
	if err != nil {
		return fmt.Errorf("не удалось проверить иерархию категорий: %w", err)
	}

	if isDescendant {
		return fmt.Errorf("нельзя переместить категорию внутрь её собственного поддерева")
	}

	// Перемещаем категорию
	_, err = s.storage.Exec(ctx, "UPDATE marketplace_categories SET parent_id = $1 WHERE id = $2", newParentID, id)
	if err != nil {
		return fmt.Errorf("не удалось переместить категорию: %w", err)
	}

	return nil
}

// GetCategoryAttributeGroups получает группы атрибутов, привязанные к категории
func (s *MarketplaceService) GetCategoryAttributeGroups(ctx context.Context, categoryID int) ([]*models.AttributeGroup, error) {
	// Если кеш не настроен, работаем напрямую
	if s.cache == nil {
		return s.getCategoryAttributeGroupsFromDB(ctx, categoryID)
	}

	// Получаем язык из контекста (по умолчанию "en")
	locale := "en"
	if lang, ok := ctx.Value("locale").(string); ok && lang != "" {
		locale = lang
	}

	// Формируем ключ кеша
	cacheKey := cache.BuildAttributeGroupsKey(int64(categoryID), locale)

	// Пытаемся получить из кеша
	var result []*models.AttributeGroup
	err := s.cache.GetOrSet(ctx, cacheKey, &result, 4*time.Hour, func() (interface{}, error) {
		// Загружаем данные из БД
		return s.getCategoryAttributeGroupsFromDB(ctx, categoryID)
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

// getCategoryAttributeGroupsFromDB получает группы атрибутов из БД
func (s *MarketplaceService) getCategoryAttributeGroupsFromDB(ctx context.Context, categoryID int) ([]*models.AttributeGroup, error) {
	query := `
		SELECT 
			ag.id, ag.name, ag.display_name, ag.description, ag.icon, 
			ag.sort_order, ag.is_active, ag.is_system, ag.created_at, ag.updated_at
		FROM attribute_groups ag
		INNER JOIN category_attribute_groups cag ON ag.id = cag.group_id
		WHERE cag.category_id = $1 AND cag.is_active = true
		ORDER BY cag.sort_order, ag.sort_order
	`

	rows, err := s.storage.Query(ctx, query, categoryID)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить группы атрибутов категории: %w", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			// Логирование ошибки закрытия rows
		}
	}()

	var groups []*models.AttributeGroup
	for rows.Next() {
		group := &models.AttributeGroup{}
		err := rows.Scan(
			&group.ID,
			&group.Name,
			&group.DisplayName,
			&group.Description,
			&group.Icon,
			&group.SortOrder,
			&group.IsActive,
			&group.IsSystem,
			&group.CreatedAt,
			&group.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("не удалось прочитать группу атрибутов: %w", err)
		}
		groups = append(groups, group)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при чтении групп атрибутов: %w", err)
	}

	return groups, nil
}

// AttachAttributeGroupToCategory привязывает группу атрибутов к категории
func (s *MarketplaceService) AttachAttributeGroupToCategory(ctx context.Context, categoryID int, groupID int, sortOrder int) (int, error) {
	// Проверяем, что категория существует
	var categoryExists bool
	err := s.storage.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM marketplace_categories WHERE id = $1)", categoryID).Scan(&categoryExists)
	if err != nil {
		return 0, fmt.Errorf("не удалось проверить существование категории: %w", err)
	}
	if !categoryExists {
		return 0, fmt.Errorf("категория с ID %d не найдена", categoryID)
	}

	// Проверяем, что группа существует
	var groupExists bool
	err = s.storage.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM attribute_groups WHERE id = $1)", groupID).Scan(&groupExists)
	if err != nil {
		return 0, fmt.Errorf("не удалось проверить существование группы: %w", err)
	}
	if !groupExists {
		return 0, fmt.Errorf("группа атрибутов с ID %d не найдена", groupID)
	}

	// Проверяем, что связь еще не существует
	var linkExists bool
	err = s.storage.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM category_attribute_groups WHERE category_id = $1 AND group_id = $2)", categoryID, groupID).Scan(&linkExists)
	if err != nil {
		return 0, fmt.Errorf("не удалось проверить существование связи: %w", err)
	}
	if linkExists {
		return 0, fmt.Errorf("группа атрибутов уже привязана к этой категории")
	}

	// Создаем связь
	query := `
		INSERT INTO category_attribute_groups (category_id, group_id, sort_order, is_active)
		VALUES ($1, $2, $3, true)
		RETURNING id
	`
	var id int
	err = s.storage.QueryRow(ctx, query, categoryID, groupID, sortOrder).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("не удалось привязать группу к категории: %w", err)
	}

	return id, nil
}

// DetachAttributeGroupFromCategory отвязывает группу атрибутов от категории
func (s *MarketplaceService) DetachAttributeGroupFromCategory(ctx context.Context, categoryID int, groupID int) error {
	// Проверяем, что связь существует
	var linkExists bool
	err := s.storage.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM category_attribute_groups WHERE category_id = $1 AND group_id = $2)", categoryID, groupID).Scan(&linkExists)
	if err != nil {
		return fmt.Errorf("не удалось проверить существование связи: %w", err)
	}
	if !linkExists {
		return fmt.Errorf("группа атрибутов не привязана к этой категории")
	}

	// Удаляем связь
	_, err = s.storage.Exec(ctx, "DELETE FROM category_attribute_groups WHERE category_id = $1 AND group_id = $2", categoryID, groupID)
	if err != nil {
		return fmt.Errorf("не удалось отвязать группу от категории: %w", err)
	}

	return nil
}

// isLikelyEnglish проверяет, похож ли текст на английский
func isLikelyEnglish(text string) bool {
	latinCount := 0
	totalLetters := 0

	for _, r := range text {
		if unicode.IsLetter(r) {
			totalLetters++
			if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
				latinCount++
			}
		}
	}

	if totalLetters == 0 {
		return false
	}

	// Если более 80% букв - латиница, считаем текст английским
	return float64(latinCount)/float64(totalLetters) > 0.8
}

// isLikelyCyrillic проверяет, похож ли текст на русский/сербский
func isLikelyCyrillic(text string) bool {
	cyrillicCount := 0
	totalLetters := 0

	for _, r := range text {
		if unicode.IsLetter(r) {
			totalLetters++
			if unicode.Is(unicode.Cyrillic, r) {
				cyrillicCount++
			}
		}
	}

	if totalLetters == 0 {
		return false
	}

	// Если более 80% букв - кириллица, считаем текст русским/сербским
	return float64(cyrillicCount)/float64(totalLetters) > 0.8
}
