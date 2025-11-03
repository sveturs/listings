// backend/internal/storage/postgres/marketplace_categories.go
package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"backend/internal/domain/models"

	"github.com/jackc/pgx/v5"
)

// GetCategories возвращает все активные категории с базовой информацией
func (db *Database) GetCategories(ctx context.Context) ([]models.MarketplaceCategory, error) {
	query := `
		SELECT
			id, name, slug, parent_id, icon, description,
			is_active, created_at, sort_order, level,
			has_custom_ui, custom_ui_component, external_id,
			seo_title, seo_description, seo_keywords,
			COALESCE((SELECT COUNT(*) FROM c2c_listings WHERE category_id = mc.id AND status = 'active'), 0) as listing_count
		FROM marketplace_categories mc
		WHERE is_active = true
		ORDER BY sort_order ASC, name ASC
	`

	rows, err := db.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query categories: %w", err)
	}
	defer rows.Close()

	var categories []models.MarketplaceCategory
	for rows.Next() {
		cat := models.MarketplaceCategory{}
		var icon, description, customUIComponent, externalID, seoTitle, seoDescription, seoKeywords sql.NullString

		err := rows.Scan(
			&cat.ID,
			&cat.Name,
			&cat.Slug,
			&cat.ParentID,
			&icon,
			&description,
			&cat.IsActive,
			&cat.CreatedAt,
			&cat.SortOrder,
			&cat.Level,
			&cat.HasCustomUI,
			&customUIComponent,
			&externalID,
			&seoTitle,
			&seoDescription,
			&seoKeywords,
			&cat.ListingCount,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan category: %w", err)
		}

		// Handle nullable fields
		if icon.Valid {
			cat.Icon = &icon.String
		}
		if description.Valid {
			cat.Description = &description.String
		}
		if customUIComponent.Valid {
			cat.CustomUIComponent = &customUIComponent.String
		}
		if externalID.Valid {
			cat.ExternalID = &externalID.String
		}
		if seoTitle.Valid {
			cat.SEOTitle = &seoTitle.String
		}
		if seoDescription.Valid {
			cat.SEODescription = &seoDescription.String
		}
		if seoKeywords.Valid {
			cat.SEOKeywords = &seoKeywords.String
		}

		categories = append(categories, cat)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return categories, nil
}

// GetAllCategories - алиас для GetCategories для обратной совместимости
func (db *Database) GetAllCategories(ctx context.Context) ([]models.MarketplaceCategory, error) {
	return db.GetCategories(ctx)
}

// GetCategoryByID возвращает категорию по ID
func (db *Database) GetCategoryByID(ctx context.Context, id int) (*models.MarketplaceCategory, error) {
	query := `
		SELECT
			id, name, slug, parent_id, icon, description,
			is_active, created_at, sort_order, level,
			has_custom_ui, custom_ui_component, external_id,
			seo_title, seo_description, seo_keywords,
			COALESCE((SELECT COUNT(*) FROM c2c_listings WHERE category_id = mc.id AND status = 'active'), 0) as listing_count
		FROM marketplace_categories mc
		WHERE id = $1
	`

	cat := &models.MarketplaceCategory{}
	var icon, description, customUIComponent, externalID, seoTitle, seoDescription, seoKeywords sql.NullString

	err := db.pool.QueryRow(ctx, query, id).Scan(
		&cat.ID,
		&cat.Name,
		&cat.Slug,
		&cat.ParentID,
		&icon,
		&description,
		&cat.IsActive,
		&cat.CreatedAt,
		&cat.SortOrder,
		&cat.Level,
		&cat.HasCustomUI,
		&customUIComponent,
		&externalID,
		&seoTitle,
		&seoDescription,
		&seoKeywords,
		&cat.ListingCount,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("category not found")
		}
		return nil, fmt.Errorf("failed to get category: %w", err)
	}

	// Handle nullable fields
	if icon.Valid {
		cat.Icon = &icon.String
	}
	if description.Valid {
		cat.Description = &description.String
	}
	if customUIComponent.Valid {
		cat.CustomUIComponent = &customUIComponent.String
	}
	if externalID.Valid {
		cat.ExternalID = &externalID.String
	}
	if seoTitle.Valid {
		cat.SEOTitle = &seoTitle.String
	}
	if seoDescription.Valid {
		cat.SEODescription = &seoDescription.String
	}
	if seoKeywords.Valid {
		cat.SEOKeywords = &seoKeywords.String
	}

	return cat, nil
}

// GetPopularCategories возвращает топ N категорий по количеству активных объявлений
func (db *Database) GetPopularCategories(ctx context.Context, limit int) ([]models.MarketplaceCategory, error) {
	if limit <= 0 {
		limit = 10 // default limit
	}

	query := `
		SELECT
			mc.id, mc.name, mc.slug, mc.parent_id, mc.icon, mc.description,
			mc.is_active, mc.created_at, mc.sort_order, mc.level,
			mc.has_custom_ui, mc.custom_ui_component, mc.external_id,
			mc.seo_title, mc.seo_description, mc.seo_keywords,
			COUNT(l.id) as listing_count
		FROM marketplace_categories mc
		LEFT JOIN c2c_listings l ON l.category_id = mc.id AND l.status = 'active'
		WHERE mc.is_active = true
		GROUP BY mc.id, mc.name, mc.slug, mc.parent_id, mc.icon, mc.description,
				 mc.is_active, mc.created_at, mc.sort_order, mc.level,
				 mc.has_custom_ui, mc.custom_ui_component, mc.external_id,
				 mc.seo_title, mc.seo_description, mc.seo_keywords
		ORDER BY listing_count DESC, mc.name ASC
		LIMIT $1
	`

	rows, err := db.pool.Query(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query popular categories: %w", err)
	}
	defer rows.Close()

	var categories []models.MarketplaceCategory
	for rows.Next() {
		cat := models.MarketplaceCategory{}
		var icon, description, customUIComponent, externalID, seoTitle, seoDescription, seoKeywords sql.NullString

		err := rows.Scan(
			&cat.ID,
			&cat.Name,
			&cat.Slug,
			&cat.ParentID,
			&icon,
			&description,
			&cat.IsActive,
			&cat.CreatedAt,
			&cat.SortOrder,
			&cat.Level,
			&cat.HasCustomUI,
			&customUIComponent,
			&externalID,
			&seoTitle,
			&seoDescription,
			&seoKeywords,
			&cat.ListingCount,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan popular category: %w", err)
		}

		// Handle nullable fields
		if icon.Valid {
			cat.Icon = &icon.String
		}
		if description.Valid {
			cat.Description = &description.String
		}
		if customUIComponent.Valid {
			cat.CustomUIComponent = &customUIComponent.String
		}
		if externalID.Valid {
			cat.ExternalID = &externalID.String
		}
		if seoTitle.Valid {
			cat.SEOTitle = &seoTitle.String
		}
		if seoDescription.Valid {
			cat.SEODescription = &seoDescription.String
		}
		if seoKeywords.Valid {
			cat.SEOKeywords = &seoKeywords.String
		}

		categories = append(categories, cat)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return categories, nil
}

// GetCategoryTree строит иерархическое дерево категорий
func (db *Database) GetCategoryTree(ctx context.Context) ([]models.CategoryTreeNode, error) {
	// Получаем все категории
	categories, err := db.GetCategories(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get categories: %w", err)
	}

	// Строим карту для быстрого поиска
	catMap := make(map[int]*models.CategoryTreeNode)
	var roots []models.CategoryTreeNode

	// Первый проход: создаем узлы и копируем данные из MarketplaceCategory
	for i := range categories {
		cat := &categories[i]
		node := &models.CategoryTreeNode{
			ID:           cat.ID,
			Name:         cat.Name,
			Slug:         cat.Slug,
			ParentID:     cat.ParentID,
			CreatedAt:    cat.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			Level:        cat.Level,
			ListingCount: cat.ListingCount,
			Children:     []models.CategoryTreeNode{},
			HasCustomUI:  cat.HasCustomUI,
			Translations: make(map[string]string),
		}

		// Handle nullable fields
		if cat.Icon != nil {
			node.Icon = *cat.Icon
		}
		if cat.CustomUIComponent != nil {
			node.CustomUIComponent = *cat.CustomUIComponent
		}

		catMap[cat.ID] = node
	}

	// Второй проход: строим дерево (связываем детей с родителями)
	var rootIDs []int
	for _, node := range catMap {
		if node.ParentID == nil {
			// Это корневая категория
			rootIDs = append(rootIDs, node.ID)
		} else {
			// Находим родителя и добавляем к нему
			parent, exists := catMap[*node.ParentID]
			if exists {
				parent.Children = append(parent.Children, *node)
			} else {
				// Если родитель не найден (не активен), делаем корневой
				rootIDs = append(rootIDs, node.ID)
			}
		}
	}

	// Третий проход: обновляем ChildrenCount и собираем корневые узлы
	for _, node := range catMap {
		node.ChildrenCount = len(node.Children)
	}
	for _, id := range rootIDs {
		roots = append(roots, *catMap[id])
	}

	return roots, nil
}
