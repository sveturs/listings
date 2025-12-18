package postgres

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vondi-global/listings/internal/domain"
)

// TestGetBySlugV2Integration tests fetching category by slug from real DB
func TestGetBySlugV2Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	repo, testDB := setupCategoryTestRepo(t)
	defer testDB.TeardownTestPostgres(t)
	ctx := context.Background()

	// Test with a known category from seed data
	cat, err := repo.GetBySlugV2(ctx, "elektronika")
	require.NoError(t, err, "should fetch elektronika category")
	require.NotNil(t, cat, "category should not be nil")

	assert.Equal(t, "elektronika", cat.Slug)
	assert.NotEqual(t, uuid.Nil, cat.ID)
	assert.Equal(t, int32(1), cat.Level)
	assert.True(t, cat.IsActive)

	// Check JSONB fields are populated
	assert.NotEmpty(t, cat.Name, "name map should not be empty")
	assert.Contains(t, cat.Name, "sr", "should have Serbian name")
	assert.Equal(t, "Elektronika", cat.Name["sr"])

	if enName, ok := cat.Name["en"]; ok {
		assert.Equal(t, "Electronics", enName)
	}

	t.Logf("✅ Category: %s (ID: %s, Level: %d)", cat.Slug, cat.ID, cat.Level)
}

// TestLocalizeCategory tests the Localize method
func TestLocalizeCategory(t *testing.T) {
	cat := &domain.CategoryV2{
		ID:    uuid.New(),
		Slug:  "test-category",
		Level: 1,
		Name: map[string]string{
			"sr": "Тест категорија",
			"en": "Test Category",
			"ru": "Тестовая категория",
		},
		Description: map[string]string{
			"sr": "Опис на српском",
			"en": "Description in English",
		},
	}

	// Test Serbian localization
	locSr := cat.Localize("sr")
	assert.Equal(t, "Тест категорија", locSr.Name)
	assert.Equal(t, "Опис на српском", locSr.Description)

	// Test English localization
	locEn := cat.Localize("en")
	assert.Equal(t, "Test Category", locEn.Name)
	assert.Equal(t, "Description in English", locEn.Description)

	// Test fallback to Serbian when Russian description is missing
	locRu := cat.Localize("ru")
	assert.Equal(t, "Тестовая категория", locRu.Name)
	assert.Equal(t, "Опис на српском", locRu.Description) // Falls back to sr

	// Test fallback for unknown locale
	locUnknown := cat.Localize("fr")
	assert.Equal(t, "Тест категорија", locUnknown.Name) // Falls back to sr
}

// TestGetTreeV2Integration tests fetching category tree with localization
func TestGetTreeV2Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	repo, testDB := setupCategoryTestRepo(t)
	defer testDB.TeardownTestPostgres(t)
	ctx := context.Background()

	// Fetch root categories in Serbian
	filter := &domain.GetCategoryTreeFilterV2{
		RootID:     nil, // Root level
		Locale:     "sr",
		ActiveOnly: true,
		MaxDepth:   nil, // No depth limit
	}

	tree, err := repo.GetTreeV2(ctx, filter)
	require.NoError(t, err, "should fetch category tree")
	require.NotEmpty(t, tree, "tree should not be empty")

	// Check first category
	firstCat := tree[0].Category
	assert.NotNil(t, firstCat)
	assert.NotEmpty(t, firstCat.Name, "localized name should not be empty")
	assert.NotEqual(t, uuid.Nil, firstCat.ID)

	t.Logf("✅ Fetched %d root categories", len(tree))
	t.Logf("   First category: %s (%s)", firstCat.Name, firstCat.Slug)
}

// TestGetBreadcrumbIntegration tests breadcrumb generation
func TestGetBreadcrumbIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	repo, testDB := setupCategoryTestRepo(t)
	defer testDB.TeardownTestPostgres(t)
	ctx := context.Background()

	// First, fetch a category to get its ID
	cat, err := repo.GetBySlugV2(ctx, "elektronika")
	require.NoError(t, err)

	// Get breadcrumb
	breadcrumbs, err := repo.GetBreadcrumb(ctx, cat.ID.String(), "sr")
	require.NoError(t, err)
	require.NotEmpty(t, breadcrumbs, "breadcrumb should not be empty")

	// Root category should have breadcrumb with just itself
	assert.Equal(t, 1, len(breadcrumbs))
	assert.Equal(t, cat.ID, breadcrumbs[0].ID)
	assert.Equal(t, cat.Slug, breadcrumbs[0].Slug)
	assert.NotEmpty(t, breadcrumbs[0].Name)

	t.Logf("✅ Breadcrumb: %s (Level %d)", breadcrumbs[0].Name, breadcrumbs[0].Level)
}

// TestListV2WithPagination tests category listing with pagination
func TestListV2WithPagination(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	repo, testDB := setupCategoryTestRepo(t)
	defer testDB.TeardownTestPostgres(t)
	ctx := context.Background()

	// List root categories (parent_id IS NULL)
	categories, total, err := repo.ListV2(ctx, nil, true, 1, 10)
	require.NoError(t, err)
	require.NotEmpty(t, categories, "should have root categories")
	assert.Greater(t, total, int64(0), "total count should be > 0")
	assert.LessOrEqual(t, len(categories), 10, "should respect page size")

	t.Logf("✅ Listed %d categories (total: %d)", len(categories), total)

	// Verify first category has all required fields
	first := categories[0]
	assert.NotEqual(t, uuid.Nil, first.ID)
	assert.NotEmpty(t, first.Slug)
	assert.NotEmpty(t, first.Name)
	assert.Equal(t, int32(1), first.Level)
}

// Note: setupCategoryTestRepo is defined in categories_repository_test.go
// It initializes a Repository instance connected to the test database
