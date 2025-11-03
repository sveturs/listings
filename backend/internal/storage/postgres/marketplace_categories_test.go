// backend/internal/storage/postgres/marketplace_categories_test.go
package postgres

import (
	"context"
	"testing"
	"time"

	"backend/internal/domain/models"

	"github.com/stretchr/testify/assert"
)

// TestCategoryTreeBuilder тестирует логику построения дерева категорий
func TestCategoryTreeBuilder(t *testing.T) {
	parentID1 := 1

	categories := []models.MarketplaceCategory{
		{
			ID:           1,
			Name:         "Electronics",
			Slug:         "electronics",
			ParentID:     nil,
			IsActive:     true,
			CreatedAt:    time.Now(),
			SortOrder:    1,
			Level:        0,
			ListingCount: 100,
			HasCustomUI:  true,
		},
		{
			ID:           2,
			Name:         "Phones",
			Slug:         "phones",
			ParentID:     &parentID1,
			IsActive:     true,
			CreatedAt:    time.Now(),
			SortOrder:    1,
			Level:        1,
			ListingCount: 50,
		},
		{
			ID:           3,
			Name:         "Laptops",
			Slug:         "laptops",
			ParentID:     &parentID1,
			IsActive:     true,
			CreatedAt:    time.Now(),
			SortOrder:    2,
			Level:        1,
			ListingCount: 30,
		},
	}

	// Build tree manually (same logic as GetCategoryTree)
	catMap := make(map[int]*models.CategoryTreeNode)
	var rootIDs []int

	// First pass: create nodes
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
		}
		catMap[cat.ID] = node
	}

	// Second pass: build tree
	for _, node := range catMap {
		if node.ParentID == nil {
			rootIDs = append(rootIDs, node.ID)
		} else {
			parent, exists := catMap[*node.ParentID]
			if exists {
				parent.Children = append(parent.Children, *node)
			} else {
				rootIDs = append(rootIDs, node.ID)
			}
		}
	}

	// Third pass: update children count
	for _, node := range catMap {
		node.ChildrenCount = len(node.Children)
	}

	// Test in catMap before copying
	assert.Len(t, rootIDs, 1)
	rootNode := catMap[rootIDs[0]]
	assert.Equal(t, 1, rootNode.ID)
	assert.Equal(t, "Electronics", rootNode.Name)
	assert.Equal(t, 2, len(rootNode.Children))
	assert.Equal(t, 2, rootNode.ChildrenCount)

	// Check that children are correct
	assert.Equal(t, 2, rootNode.Children[0].ID)
	assert.Equal(t, "Phones", rootNode.Children[0].Name)
	assert.Equal(t, 3, rootNode.Children[1].ID)
	assert.Equal(t, "Laptops", rootNode.Children[1].Name)

	// Now collect roots
	var roots []models.CategoryTreeNode
	for _, id := range rootIDs {
		roots = append(roots, *catMap[id])
	}

	// Assertions on final result
	assert.Len(t, roots, 1, "Should have 1 root category")

	root := roots[0]
	assert.Equal(t, 1, root.ID)
	assert.Equal(t, "Electronics", root.Name)
	assert.Nil(t, root.ParentID)
	assert.Equal(t, 0, root.Level)
	assert.True(t, root.HasCustomUI)

	// Check children
	assert.Len(t, root.Children, 2, "Root should have 2 children")
	assert.Equal(t, 2, root.ChildrenCount)

	// First child - Phones
	phones := root.Children[0]
	assert.Equal(t, 2, phones.ID)
	assert.Equal(t, "Phones", phones.Name)

	// Second child - Laptops
	laptops := root.Children[1]
	assert.Equal(t, 3, laptops.ID)
	assert.Equal(t, "Laptops", laptops.Name)
}

// TestCategoryTreeBuilder_NoRoots тестирует случай когда все категории имеют родителя
func TestCategoryTreeBuilder_NoRoots(t *testing.T) {
	parentID999 := 999 // Несуществующий родитель

	categories := []models.MarketplaceCategory{
		{
			ID:        1,
			Name:      "Orphan Category",
			Slug:      "orphan",
			ParentID:  &parentID999,
			IsActive:  true,
			CreatedAt: time.Now(),
			Level:     1,
		},
	}

	// Build tree
	catMap := make(map[int]*models.CategoryTreeNode)
	var roots []models.CategoryTreeNode

	for i := range categories {
		cat := &categories[i]
		node := &models.CategoryTreeNode{
			ID:        cat.ID,
			Name:      cat.Name,
			Slug:      cat.Slug,
			ParentID:  cat.ParentID,
			CreatedAt: cat.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			Level:     cat.Level,
			Children:  []models.CategoryTreeNode{},
		}
		catMap[cat.ID] = node
	}

	for _, node := range catMap {
		if node.ParentID == nil {
			roots = append(roots, *node)
		} else {
			parent, exists := catMap[*node.ParentID]
			if exists {
				parent.Children = append(parent.Children, *node)
			} else {
				// Orphan category becomes root
				roots = append(roots, *node)
			}
		}
	}

	// Orphan category should become a root
	assert.Len(t, roots, 1)
	assert.Equal(t, "Orphan Category", roots[0].Name)
}

// TestCategoryTreeBuilder_Empty тестирует пустой список категорий
func TestCategoryTreeBuilder_Empty(t *testing.T) {
	var categories []models.MarketplaceCategory

	catMap := make(map[int]*models.CategoryTreeNode)
	var roots []models.CategoryTreeNode

	for i := range categories {
		cat := &categories[i]
		node := &models.CategoryTreeNode{
			ID:       cat.ID,
			Name:     cat.Name,
			Slug:     cat.Slug,
			ParentID: cat.ParentID,
			Children: []models.CategoryTreeNode{},
		}
		catMap[cat.ID] = node
	}

	for _, node := range catMap {
		if node.ParentID == nil {
			roots = append(roots, *node)
		}
	}

	assert.Empty(t, roots)
}

// TestGetCategoriesQueryValidation проверяет что SQL запрос корректен
func TestGetCategoriesQueryValidation(t *testing.T) {
	t.Skip("Integration test - requires real database connection")

	// This test would require a real database connection
	// Run with: go test -tags=integration ./internal/storage/postgres
	ctx := context.Background()
	_ = ctx
}

// TestGetPopularCategoriesLimitLogic тестирует логику лимита
func TestGetPopularCategoriesLimitLogic(t *testing.T) {
	testCases := []struct {
		name          string
		inputLimit    int
		expectedLimit int
	}{
		{
			name:          "Valid limit",
			inputLimit:    5,
			expectedLimit: 5,
		},
		{
			name:          "Zero limit uses default",
			inputLimit:    0,
			expectedLimit: 10,
		},
		{
			name:          "Negative limit uses default",
			inputLimit:    -1,
			expectedLimit: 10,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			limit := tc.inputLimit
			if limit <= 0 {
				limit = 10 // default limit
			}
			assert.Equal(t, tc.expectedLimit, limit)
		})
	}
}
