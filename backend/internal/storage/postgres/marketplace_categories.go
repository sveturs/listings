// backend/internal/storage/postgres/marketplace_categories.go
package postgres

import (
	"context"
	"fmt"

	"backend/internal/domain/models"
)

// GetCategories возвращает все активные категории с базовой информацией
func (db *Database) GetCategories(ctx context.Context) ([]models.MarketplaceCategory, error) {
	categoriesPtrs, err := db.grpcClient.GetAllCategories(ctx)
	if err != nil {
		return nil, err
	}

	// Convert []*models.MarketplaceCategory to []models.MarketplaceCategory
	categories := make([]models.MarketplaceCategory, len(categoriesPtrs))
	for i, cat := range categoriesPtrs {
		categories[i] = *cat
	}

	return categories, nil
}

// GetAllCategories - алиас для GetCategories для обратной совместимости
func (db *Database) GetAllCategories(ctx context.Context) ([]models.MarketplaceCategory, error) {
	return db.GetCategories(ctx)
}

// GetCategoryByID возвращает категорию по ID
func (db *Database) GetCategoryByID(ctx context.Context, id int) (*models.MarketplaceCategory, error) {
	return db.grpcClient.GetCategoryByID(ctx, int64(id))
}

// GetPopularCategories возвращает топ N категорий по количеству активных объявлений
func (db *Database) GetPopularCategories(ctx context.Context, limit int) ([]models.MarketplaceCategory, error) {
	if limit <= 0 {
		limit = 10 // default limit
	}

	categoriesPtrs, err := db.grpcClient.GetPopularCategories(ctx, limit)
	if err != nil {
		return nil, err
	}

	// Convert []*models.MarketplaceCategory to []models.MarketplaceCategory
	categories := make([]models.MarketplaceCategory, len(categoriesPtrs))
	for i, cat := range categoriesPtrs {
		categories[i] = *cat
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
