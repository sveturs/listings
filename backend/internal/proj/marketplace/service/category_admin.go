package service

import (
	"context"
	"fmt"

	"backend/internal/domain/models"
)

// Реализации методов для администрирования категорий

// CreateCategory создает новую категорию
func (s *MarketplaceService) CreateCategory(ctx context.Context, category *models.MarketplaceCategory) (int, error) {
	// Создаем категорию в БД
	query := `
		INSERT INTO marketplace_categories (name, slug, parent_id, icon, has_custom_ui, custom_ui_component)
		VALUES ($1, $2, $3, $4, $5, $6)
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
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("не удалось создать категорию: %w", err)
	}

	// Если есть переводы, сохраняем их
	if category.Translations != nil && len(category.Translations) > 0 {
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
	}

	// Обновляем материализованное представление для обновления счетчиков
	_ = s.RefreshCategoryListingCounts(ctx)

	return id, nil
}

// UpdateCategory обновляет существующую категорию
func (s *MarketplaceService) UpdateCategory(ctx context.Context, category *models.MarketplaceCategory) error {
	// Обновляем категорию в БД
	query := `
		UPDATE marketplace_categories
		SET name = $1, slug = $2, parent_id = $3, icon = $4, has_custom_ui = $5, custom_ui_component = $6
		WHERE id = $7
	`
	_, err := s.storage.Exec(ctx, query,
		category.Name,
		category.Slug,
		category.ParentID,
		category.Icon,
		category.HasCustomUI,
		category.CustomUIComponent,
		category.ID,
	)
	if err != nil {
		return fmt.Errorf("не удалось обновить категорию: %w", err)
	}

	// Если есть переводы, обновляем их
	if category.Translations != nil && len(category.Translations) > 0 {
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

	return nil
}

// ReorderCategories изменяет порядок категорий
func (s *MarketplaceService) ReorderCategories(ctx context.Context, orderedIDs []int) error {
	// Начинаем транзакцию для атомарной операции
	tx, err := s.storage.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("не удалось начать транзакцию: %w", err)
	}
	defer tx.Rollback()

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
