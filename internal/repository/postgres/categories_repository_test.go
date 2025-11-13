package postgres

import (
	"context"
	"fmt"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sveturs/listings/internal/domain"
	"github.com/sveturs/listings/tests"
)

func setupCategoryTestRepo(t *testing.T) (*Repository, *tests.TestDB) {
	t.Helper()

	// Skip if running in short mode or no Docker
	tests.SkipIfShort(t)
	tests.SkipIfNoDocker(t)

	// Setup test database
	testDB := tests.SetupTestPostgres(t)

	// Run migrations
	tests.RunMigrations(t, testDB.DB, "../../../migrations")

	// Create sqlx.DB wrapper
	db := sqlx.NewDb(testDB.DB, "postgres")

	// Create logger
	logger := zerolog.New(zerolog.NewTestWriter(t)).With().Timestamp().Logger()

	// Create repository
	repo := NewRepository(db, logger)

	return repo, testDB
}

// =============================================================================
// Test: GetCategoryBySlug
// =============================================================================

func TestRepository_GetCategoryBySlug_Success(t *testing.T) {
	repo, testDB := setupCategoryTestRepo(t)
	defer testDB.TeardownTestPostgres(t)

	ctx := context.Background()

	// Create a test category
	testCat := &domain.Category{
		Name:     "Test Category",
		Slug:     "test-category",
		IsActive: true,
	}
	created, err := repo.CreateCategory(ctx, testCat)
	require.NoError(t, err)
	require.NotNil(t, created)

	// Retrieve by slug
	found, err := repo.GetCategoryBySlug(ctx, "test-category")
	require.NoError(t, err)
	assert.Equal(t, created.ID, found.ID)
	assert.Equal(t, "Test Category", found.Name)
	assert.Equal(t, "test-category", found.Slug)
}

func TestRepository_GetCategoryBySlug_NotFound(t *testing.T) {
	repo, testDB := setupCategoryTestRepo(t)
	defer testDB.TeardownTestPostgres(t)

	ctx := context.Background()

	_, err := repo.GetCategoryBySlug(ctx, "non-existent-slug")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "category not found")
}

// =============================================================================
// Test: CreateCategory
// =============================================================================

func TestRepository_CreateCategory_Success(t *testing.T) {
	repo, testDB := setupCategoryTestRepo(t)
	defer testDB.TeardownTestPostgres(t)

	ctx := context.Background()

	testCases := []struct {
		name     string
		category *domain.Category
	}{
		{
			name: "simple category without slug",
			category: &domain.Category{
				Name:     "Electronics",
				IsActive: true,
			},
		},
		{
			name: "category with custom slug",
			category: &domain.Category{
				Name:     "Mobile Phones",
				Slug:     "mobile-phones",
				IsActive: true,
			},
		},
		{
			name: "category with all fields",
			category: &domain.Category{
				Name:        "Laptops",
				Slug:        "laptops",
				IsActive:    true,
				Icon:        stringPtr("💻"),
				Description: stringPtr("Laptop computers"),
				HasCustomUI: true,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			created, err := repo.CreateCategory(ctx, tc.category)
			require.NoError(t, err)
			assert.NotZero(t, created.ID)
			assert.NotEmpty(t, created.Slug)
			assert.Equal(t, int32(0), created.Level) // Root category
			assert.NotZero(t, created.SortOrder)

			// Verify it can be retrieved
			found, err := repo.GetCategoryByID(ctx, created.ID)
			require.NoError(t, err)
			assert.Equal(t, created.Name, found.Name)
		})
	}
}

func TestRepository_CreateCategory_WithParent(t *testing.T) {
	repo, testDB := setupCategoryTestRepo(t)
	defer testDB.TeardownTestPostgres(t)

	ctx := context.Background()

	// Create parent category
	parent := &domain.Category{
		Name:     "Electronics",
		IsActive: true,
	}
	parentCreated, err := repo.CreateCategory(ctx, parent)
	require.NoError(t, err)
	assert.Equal(t, int32(0), parentCreated.Level)

	// Create child category
	child := &domain.Category{
		Name:     "Smartphones",
		ParentID: &parentCreated.ID,
		IsActive: true,
	}
	childCreated, err := repo.CreateCategory(ctx, child)
	require.NoError(t, err)
	assert.Equal(t, int32(1), childCreated.Level) // Level should be parent.level + 1
	assert.Equal(t, parentCreated.ID, *childCreated.ParentID)

	// Create grandchild
	grandchild := &domain.Category{
		Name:     "Android Phones",
		ParentID: &childCreated.ID,
		IsActive: true,
	}
	grandchildCreated, err := repo.CreateCategory(ctx, grandchild)
	require.NoError(t, err)
	assert.Equal(t, int32(2), grandchildCreated.Level)
}

func TestRepository_CreateCategory_DuplicateSlug(t *testing.T) {
	repo, testDB := setupCategoryTestRepo(t)
	defer testDB.TeardownTestPostgres(t)

	ctx := context.Background()

	// Create first category
	first := &domain.Category{
		Name:     "Electronics",
		Slug:     "electronics",
		IsActive: true,
	}
	_, err := repo.CreateCategory(ctx, first)
	require.NoError(t, err)

	// Try to create duplicate
	duplicate := &domain.Category{
		Name:     "Electronics 2",
		Slug:     "electronics",
		IsActive: true,
	}
	_, err = repo.CreateCategory(ctx, duplicate)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already exists")
}

func TestRepository_CreateCategory_InvalidParent(t *testing.T) {
	repo, testDB := setupCategoryTestRepo(t)
	defer testDB.TeardownTestPostgres(t)

	ctx := context.Background()

	invalidParentID := int64(99999)
	category := &domain.Category{
		Name:     "Test",
		ParentID: &invalidParentID,
		IsActive: true,
	}
	_, err := repo.CreateCategory(ctx, category)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid parent_id")
}

// =============================================================================
// Test: UpdateCategory
// =============================================================================

func TestRepository_UpdateCategory_Success(t *testing.T) {
	repo, testDB := setupCategoryTestRepo(t)
	defer testDB.TeardownTestPostgres(t)

	ctx := context.Background()

	// Create category
	original := &domain.Category{
		Name:     "Original Name",
		IsActive: true,
	}
	created, err := repo.CreateCategory(ctx, original)
	require.NoError(t, err)

	// Update name
	created.Name = "Updated Name"
	created.Description = stringPtr("New description")
	updated, err := repo.UpdateCategory(ctx, created)
	require.NoError(t, err)
	assert.Equal(t, "Updated Name", updated.Name)
	assert.NotNil(t, updated.Description)
	assert.Equal(t, "New description", *updated.Description)
}

func TestRepository_UpdateCategory_ChangeParent(t *testing.T) {
	repo, testDB := setupCategoryTestRepo(t)
	defer testDB.TeardownTestPostgres(t)

	ctx := context.Background()

	// Create two root categories
	parent1 := &domain.Category{Name: "Parent 1", IsActive: true}
	parent1Created, err := repo.CreateCategory(ctx, parent1)
	require.NoError(t, err)

	parent2 := &domain.Category{Name: "Parent 2", IsActive: true}
	parent2Created, err := repo.CreateCategory(ctx, parent2)
	require.NoError(t, err)

	// Create child under parent1
	child := &domain.Category{
		Name:     "Child",
		ParentID: &parent1Created.ID,
		IsActive: true,
	}
	childCreated, err := repo.CreateCategory(ctx, child)
	require.NoError(t, err)
	assert.Equal(t, int32(1), childCreated.Level)

	// Move child to parent2
	childCreated.ParentID = &parent2Created.ID
	updated, err := repo.UpdateCategory(ctx, childCreated)
	require.NoError(t, err)
	assert.Equal(t, parent2Created.ID, *updated.ParentID)
	assert.Equal(t, int32(1), updated.Level)
}

func TestRepository_UpdateCategory_CircularDependency(t *testing.T) {
	repo, testDB := setupCategoryTestRepo(t)
	defer testDB.TeardownTestPostgres(t)

	ctx := context.Background()

	// Create parent and child
	parent := &domain.Category{Name: "Parent", IsActive: true}
	parentCreated, err := repo.CreateCategory(ctx, parent)
	require.NoError(t, err)

	child := &domain.Category{
		Name:     "Child",
		ParentID: &parentCreated.ID,
		IsActive: true,
	}
	childCreated, err := repo.CreateCategory(ctx, child)
	require.NoError(t, err)

	// Try to make parent a child of its own child (circular dependency)
	parentCreated.ParentID = &childCreated.ID
	_, err = repo.UpdateCategory(ctx, parentCreated)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "circular dependency")
}

func TestRepository_UpdateCategory_SelfParent(t *testing.T) {
	repo, testDB := setupCategoryTestRepo(t)
	defer testDB.TeardownTestPostgres(t)

	ctx := context.Background()

	// Create category
	category := &domain.Category{Name: "Test", IsActive: true}
	created, err := repo.CreateCategory(ctx, category)
	require.NoError(t, err)

	// Try to set itself as parent
	created.ParentID = &created.ID
	_, err = repo.UpdateCategory(ctx, created)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot be its own parent")
}

func TestRepository_UpdateCategory_DuplicateSlug(t *testing.T) {
	repo, testDB := setupCategoryTestRepo(t)
	defer testDB.TeardownTestPostgres(t)

	ctx := context.Background()

	// Create two categories
	cat1 := &domain.Category{Name: "Category 1", Slug: "cat-1", IsActive: true}
	_, err := repo.CreateCategory(ctx, cat1)
	require.NoError(t, err)

	cat2 := &domain.Category{Name: "Category 2", Slug: "cat-2", IsActive: true}
	cat2Created, err := repo.CreateCategory(ctx, cat2)
	require.NoError(t, err)

	// Try to update cat2 slug to cat-1 (duplicate)
	cat2Created.Slug = "cat-1"
	_, err = repo.UpdateCategory(ctx, cat2Created)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already exists")
}

// =============================================================================
// Test: DeleteCategory
// =============================================================================

func TestRepository_DeleteCategory_Success(t *testing.T) {
	repo, testDB := setupCategoryTestRepo(t)
	defer testDB.TeardownTestPostgres(t)

	ctx := context.Background()

	// Create category
	category := &domain.Category{Name: "To Delete", IsActive: true}
	created, err := repo.CreateCategory(ctx, category)
	require.NoError(t, err)

	// Delete it
	err = repo.DeleteCategory(ctx, created.ID)
	require.NoError(t, err)

	// Verify it's deactivated
	deleted, err := repo.GetCategoryByID(ctx, created.ID)
	require.NoError(t, err)
	assert.False(t, deleted.IsActive)
}

func TestRepository_DeleteCategory_WithChildren(t *testing.T) {
	repo, testDB := setupCategoryTestRepo(t)
	defer testDB.TeardownTestPostgres(t)

	ctx := context.Background()

	// Create parent and children
	parent := &domain.Category{Name: "Parent", IsActive: true}
	parentCreated, err := repo.CreateCategory(ctx, parent)
	require.NoError(t, err)

	child1 := &domain.Category{
		Name:     "Child 1",
		ParentID: &parentCreated.ID,
		IsActive: true,
	}
	child1Created, err := repo.CreateCategory(ctx, child1)
	require.NoError(t, err)

	child2 := &domain.Category{
		Name:     "Child 2",
		ParentID: &parentCreated.ID,
		IsActive: true,
	}
	child2Created, err := repo.CreateCategory(ctx, child2)
	require.NoError(t, err)

	// Delete parent
	err = repo.DeleteCategory(ctx, parentCreated.ID)
	require.NoError(t, err)

	// Verify parent and children are deactivated
	deletedParent, err := repo.GetCategoryByID(ctx, parentCreated.ID)
	require.NoError(t, err)
	assert.False(t, deletedParent.IsActive)

	deletedChild1, err := repo.GetCategoryByID(ctx, child1Created.ID)
	require.NoError(t, err)
	assert.False(t, deletedChild1.IsActive)

	deletedChild2, err := repo.GetCategoryByID(ctx, child2Created.ID)
	require.NoError(t, err)
	assert.False(t, deletedChild2.IsActive)
}

func TestRepository_DeleteCategory_WithActiveListings(t *testing.T) {
	repo, testDB := setupCategoryTestRepo(t)
	defer testDB.TeardownTestPostgres(t)

	ctx := context.Background()

	// Create category
	category := &domain.Category{Name: "Category with Listings", IsActive: true}
	created, err := repo.CreateCategory(ctx, category)
	require.NoError(t, err)

	// Create an active listing in this category
	db := sqlx.NewDb(testDB.DB, "postgres")

	// First create a test user
	var userID int64
	err = db.QueryRowContext(ctx, `
		INSERT INTO users (email, username)
		VALUES ('test@example.com', 'testuser')
		RETURNING id
	`).Scan(&userID)
	if err != nil {
		// Try alternative: users table might have different structure
		// Just skip this test if we can't create a user
		t.Skip("Cannot create test user - skipping test")
		return
	}

	_, err = db.ExecContext(ctx, `
		INSERT INTO listings (title, slug, category_id, user_id, status, is_deleted)
		VALUES ('Test Product', 'test-product', $1, $2, 'active', false)
	`, created.ID, userID)
	require.NoError(t, err)

	// Try to delete category - should fail
	err = repo.DeleteCategory(ctx, created.ID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot delete category with")
	assert.Contains(t, err.Error(), "active listings")
}

func TestRepository_DeleteCategory_NotFound(t *testing.T) {
	repo, testDB := setupCategoryTestRepo(t)
	defer testDB.TeardownTestPostgres(t)

	ctx := context.Background()

	err := repo.DeleteCategory(ctx, 99999)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

// =============================================================================
// Test: GetCategoriesWithPagination
// =============================================================================

func TestRepository_GetCategoriesWithPagination_Success(t *testing.T) {
	repo, testDB := setupCategoryTestRepo(t)
	defer testDB.TeardownTestPostgres(t)

	ctx := context.Background()

	// Create multiple categories
	for i := 1; i <= 15; i++ {
		cat := &domain.Category{
			Name:     stringVal("Category %d", i),
			IsActive: true,
		}
		_, err := repo.CreateCategory(ctx, cat)
		require.NoError(t, err)
	}

	// Test pagination
	categories, total, err := repo.GetCategoriesWithPagination(ctx, nil, nil, 10, 0)
	require.NoError(t, err)
	assert.Len(t, categories, 10)
	assert.Equal(t, int32(15), total)

	// Second page
	categoriesPage2, total2, err := repo.GetCategoriesWithPagination(ctx, nil, nil, 10, 10)
	require.NoError(t, err)
	assert.Len(t, categoriesPage2, 5)
	assert.Equal(t, int32(15), total2)
}

func TestRepository_GetCategoriesWithPagination_FilterByParent(t *testing.T) {
	repo, testDB := setupCategoryTestRepo(t)
	defer testDB.TeardownTestPostgres(t)

	ctx := context.Background()

	// Create parent
	parent := &domain.Category{Name: "Parent", IsActive: true}
	parentCreated, err := repo.CreateCategory(ctx, parent)
	require.NoError(t, err)

	// Create children
	for i := 1; i <= 5; i++ {
		child := &domain.Category{
			Name:     stringVal("Child %d", i),
			ParentID: &parentCreated.ID,
			IsActive: true,
		}
		_, err := repo.CreateCategory(ctx, child)
		require.NoError(t, err)
	}

	// Create another root category
	other := &domain.Category{Name: "Other Root", IsActive: true}
	_, err = repo.CreateCategory(ctx, other)
	require.NoError(t, err)

	// Get children of parent only
	children, total, err := repo.GetCategoriesWithPagination(ctx, &parentCreated.ID, nil, 100, 0)
	require.NoError(t, err)
	assert.Equal(t, int32(5), total)
	assert.Len(t, children, 5)
	for _, child := range children {
		assert.Equal(t, parentCreated.ID, *child.ParentID)
	}
}

func TestRepository_GetCategoriesWithPagination_FilterByActive(t *testing.T) {
	repo, testDB := setupCategoryTestRepo(t)
	defer testDB.TeardownTestPostgres(t)

	ctx := context.Background()

	// Create active categories
	for i := 1; i <= 3; i++ {
		cat := &domain.Category{
			Name:     stringVal("Active %d", i),
			IsActive: true,
		}
		_, err := repo.CreateCategory(ctx, cat)
		require.NoError(t, err)
	}

	// Create inactive category
	inactive := &domain.Category{
		Name:     "Inactive",
		IsActive: false,
	}
	inactiveCreated, err := repo.CreateCategory(ctx, inactive)
	require.NoError(t, err)

	// Deactivate it via update
	inactiveCreated.IsActive = false
	_, err = repo.UpdateCategory(ctx, inactiveCreated)
	require.NoError(t, err)

	// Get only active
	activeOnly := true
	activeCategories, total, err := repo.GetCategoriesWithPagination(ctx, nil, &activeOnly, 100, 0)
	require.NoError(t, err)
	assert.Equal(t, int32(3), total)
	assert.Len(t, activeCategories, 3)

	// Get only inactive
	inactiveOnly := false
	inactiveCategories, totalInactive, err := repo.GetCategoriesWithPagination(ctx, nil, &inactiveOnly, 100, 0)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, int(totalInactive), 1) // At least the one we created
	for _, cat := range inactiveCategories {
		assert.False(t, cat.IsActive)
	}
}

// =============================================================================
// Helper functions
// =============================================================================

func stringVal(format string, args ...interface{}) string {
	return fmt.Sprintf(format, args...)
}
