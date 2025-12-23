// Package service implements business logic for the listings microservice.
package service

import (
	"context"
	"errors"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/vondi-global/listings/internal/domain"
)

// MockCategoryRepository is a mock for category repository operations
type MockCategoryRepository struct {
	mock.Mock
}

func (m *MockCategoryRepository) GetCategoriesWithPagination(ctx context.Context, parentID *string, isActive *bool, limit, offset int32) ([]*domain.Category, int32, error) {
	args := m.Called(ctx, parentID, isActive, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int32), args.Error(2)
	}
	return args.Get(0).([]*domain.Category), args.Get(1).(int32), args.Error(2)
}

func (m *MockCategoryRepository) GetCategoryByID(ctx context.Context, id string) (*domain.Category, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Category), args.Error(1)
}

func (m *MockCategoryRepository) GetCategoryBySlug(ctx context.Context, slug string) (*domain.Category, error) {
	args := m.Called(ctx, slug)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Category), args.Error(1)
}

func (m *MockCategoryRepository) GetCategoryTree(ctx context.Context, categoryID string) (*domain.CategoryTreeNode, error) {
	args := m.Called(ctx, categoryID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.CategoryTreeNode), args.Error(1)
}

func (m *MockCategoryRepository) CreateCategory(ctx context.Context, cat *domain.Category) (*domain.Category, error) {
	args := m.Called(ctx, cat)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Category), args.Error(1)
}

func (m *MockCategoryRepository) UpdateCategory(ctx context.Context, cat *domain.Category) (*domain.Category, error) {
	args := m.Called(ctx, cat)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Category), args.Error(1)
}

func (m *MockCategoryRepository) DeleteCategory(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// Helper function to create a test service with minimal Redis mock
func setupTestCategoryService(t *testing.T) *CategoryServiceImpl {
	// Use a simple Redis client (will gracefully handle connection errors for cache)
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Will fail gracefully in tests
	})

	logger := zerolog.New(nil).Level(zerolog.Disabled)

	return &CategoryServiceImpl{
		repo:   nil, // Will be set by individual tests
		cache:  NewCategoryCache(redisClient, logger),
		logger: logger,
	}
}

// =============================================================================
// Test: GetCategories
// =============================================================================

func TestCategoryService_GetCategories_Success(t *testing.T) {
	service := setupTestCategoryService(t)
	mockRepo := new(MockCategoryRepository)
	service.repo = mockRepo

	ctx := context.Background()
	var parentID *string
	isActive := true

	expectedCategories := []*domain.Category{
		{ID: "1", Name: "Electronics", Slug: "electronics", IsActive: true},
		{ID: "2", Name: "Books", Slug: "books", IsActive: true},
	}

	mockRepo.On("GetCategoriesWithPagination", ctx, parentID, &isActive, int32(10), int32(0)).
		Return(expectedCategories, int32(2), nil)

	categories, total, err := service.GetCategories(ctx, parentID, &isActive, 10, 0)

	assert.NoError(t, err)
	assert.Equal(t, int32(2), total)
	assert.Len(t, categories, 2)
	assert.Equal(t, "Electronics", categories[0].Name)
	mockRepo.AssertExpectations(t)
}

func TestCategoryService_GetCategories_RepositoryError(t *testing.T) {
	service := setupTestCategoryService(t)
	mockRepo := new(MockCategoryRepository)
	service.repo = mockRepo

	ctx := context.Background()
	var parentID *string
	isActive := true

	mockRepo.On("GetCategoriesWithPagination", ctx, parentID, &isActive, int32(10), int32(0)).
		Return(nil, int32(0), errors.New("database error"))

	categories, total, err := service.GetCategories(ctx, parentID, &isActive, 10, 0)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get categories")
	assert.Nil(t, categories)
	assert.Equal(t, int32(0), total)
	mockRepo.AssertExpectations(t)
}

// =============================================================================
// Test: GetCategory
// =============================================================================

func TestCategoryService_GetCategory_Success(t *testing.T) {
	service := setupTestCategoryService(t)
	mockRepo := new(MockCategoryRepository)
	service.repo = mockRepo

	ctx := context.Background()
	expectedCategory := &domain.Category{
		ID:       "1",
		Name:     "Electronics",
		Slug:     "electronics",
		IsActive: true,
	}

	mockRepo.On("GetCategoryByID", ctx, "1").Return(expectedCategory, nil)

	category, err := service.GetCategory(ctx, "1")

	assert.NoError(t, err)
	assert.NotNil(t, category)
	assert.Equal(t, "Electronics", category.Name)
	mockRepo.AssertExpectations(t)
}

func TestCategoryService_GetCategory_NotFound(t *testing.T) {
	service := setupTestCategoryService(t)
	mockRepo := new(MockCategoryRepository)
	service.repo = mockRepo

	ctx := context.Background()

	mockRepo.On("GetCategoryByID", ctx, "999").Return(nil, errors.New("not found"))

	category, err := service.GetCategory(ctx, "999")

	assert.Error(t, err)
	assert.Nil(t, category)
	mockRepo.AssertExpectations(t)
}

// =============================================================================
// Test: GetCategoryBySlug
// =============================================================================

func TestCategoryService_GetCategoryBySlug_Success(t *testing.T) {
	service := setupTestCategoryService(t)
	mockRepo := new(MockCategoryRepository)
	service.repo = mockRepo

	ctx := context.Background()
	expectedCategory := &domain.Category{
		ID:       "1",
		Name:     "Electronics",
		Slug:     "electronics",
		IsActive: true,
	}

	mockRepo.On("GetCategoryBySlug", ctx, "electronics").Return(expectedCategory, nil)

	category, err := service.GetCategoryBySlug(ctx, "electronics")

	assert.NoError(t, err)
	assert.NotNil(t, category)
	assert.Equal(t, "Electronics", category.Name)
	mockRepo.AssertExpectations(t)
}

// =============================================================================
// Test: CreateCategory
// =============================================================================

func TestCategoryService_CreateCategory_Success(t *testing.T) {
	service := setupTestCategoryService(t)
	mockRepo := new(MockCategoryRepository)
	service.repo = mockRepo

	ctx := context.Background()
	newCategory := &domain.Category{
		Name:     "Electronics",
		Slug:     "electronics",
		IsActive: true,
	}

	createdCategory := &domain.Category{
		ID:       "1",
		Name:     "Electronics",
		Slug:     "electronics",
		IsActive: true,
	}

	// Check slug uniqueness
	mockRepo.On("GetCategoryBySlug", ctx, "electronics").Return(nil, errors.New("not found"))
	// Create category
	mockRepo.On("CreateCategory", ctx, newCategory).Return(createdCategory, nil)

	created, err := service.CreateCategory(ctx, newCategory)

	assert.NoError(t, err)
	assert.NotNil(t, created)
	assert.Equal(t, int64(1), created.ID)
	assert.Equal(t, "Electronics", created.Name)
	mockRepo.AssertExpectations(t)
}

func TestCategoryService_CreateCategory_ValidationError_ShortName(t *testing.T) {
	service := setupTestCategoryService(t)
	mockRepo := new(MockCategoryRepository)
	service.repo = mockRepo

	ctx := context.Background()
	newCategory := &domain.Category{
		Name:     "E", // Too short
		Slug:     "e",
		IsActive: true,
	}

	created, err := service.CreateCategory(ctx, newCategory)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "at least 2 characters")
	assert.Nil(t, created)
	mockRepo.AssertNotCalled(t, "CreateCategory")
}

func TestCategoryService_CreateCategory_DuplicateSlug(t *testing.T) {
	service := setupTestCategoryService(t)
	mockRepo := new(MockCategoryRepository)
	service.repo = mockRepo

	ctx := context.Background()
	newCategory := &domain.Category{
		Name:     "Electronics",
		Slug:     "electronics",
		IsActive: true,
	}

	existingCategory := &domain.Category{
		ID:   "1",
		Name: "Electronics",
		Slug: "electronics",
	}

	mockRepo.On("GetCategoryBySlug", ctx, "electronics").Return(existingCategory, nil)

	created, err := service.CreateCategory(ctx, newCategory)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already exists")
	assert.Nil(t, created)
	mockRepo.AssertNotCalled(t, "CreateCategory")
	mockRepo.AssertExpectations(t)
}

func TestCategoryService_CreateCategory_InvalidParent(t *testing.T) {
	service := setupTestCategoryService(t)
	mockRepo := new(MockCategoryRepository)
	service.repo = mockRepo

	ctx := context.Background()
	parentID := "999"
	newCategory := &domain.Category{
		Name:     "Laptops",
		Slug:     "laptops",
		ParentID: &parentID,
		IsActive: true,
	}

	mockRepo.On("GetCategoryBySlug", ctx, "laptops").Return(nil, errors.New("not found"))
	mockRepo.On("GetCategoryByID", ctx, "999").Return(nil, errors.New("not found"))

	created, err := service.CreateCategory(ctx, newCategory)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "parent category")
	assert.Nil(t, created)
	mockRepo.AssertNotCalled(t, "CreateCategory")
	mockRepo.AssertExpectations(t)
}

// =============================================================================
// Test: UpdateCategory
// =============================================================================

func TestCategoryService_UpdateCategory_Success(t *testing.T) {
	service := setupTestCategoryService(t)
	mockRepo := new(MockCategoryRepository)
	service.repo = mockRepo

	ctx := context.Background()
	existingCategory := &domain.Category{
		ID:       "1",
		Name:     "Electronics",
		Slug:     "electronics",
		IsActive: true,
	}

	updatedCategory := &domain.Category{
		ID:       "1",
		Name:     "Consumer Electronics",
		Slug:     "consumer-electronics",
		IsActive: true,
	}

	// First call: get existing category
	mockRepo.On("GetCategoryByID", ctx, "1").Return(existingCategory, nil).Once()
	// Check slug uniqueness
	mockRepo.On("GetCategoryBySlug", ctx, "consumer-electronics").Return(nil, errors.New("not found"))
	// Update category
	mockRepo.On("UpdateCategory", ctx, updatedCategory).Return(updatedCategory, nil)
	// Cache invalidation: get category again to get slug
	mockRepo.On("GetCategoryByID", ctx, "1").Return(updatedCategory, nil).Once()

	updated, err := service.UpdateCategory(ctx, updatedCategory)

	assert.NoError(t, err)
	assert.NotNil(t, updated)
	assert.Equal(t, "Consumer Electronics", updated.Name)
	mockRepo.AssertExpectations(t)
}

func TestCategoryService_UpdateCategory_NotFound(t *testing.T) {
	service := setupTestCategoryService(t)
	mockRepo := new(MockCategoryRepository)
	service.repo = mockRepo

	ctx := context.Background()
	updatedCategory := &domain.Category{
		ID:   "999",
		Name: "Non-existent",
		Slug: "non-existent",
	}

	mockRepo.On("GetCategoryByID", ctx, "999").Return(nil, errors.New("not found"))

	updated, err := service.UpdateCategory(ctx, updatedCategory)

	assert.Error(t, err)
	assert.Nil(t, updated)
	mockRepo.AssertNotCalled(t, "UpdateCategory")
	mockRepo.AssertExpectations(t)
}

func TestCategoryService_UpdateCategory_CircularDependency(t *testing.T) {
	service := setupTestCategoryService(t)
	mockRepo := new(MockCategoryRepository)
	service.repo = mockRepo

	ctx := context.Background()
	categoryID := "1"
	existingCategory := &domain.Category{
		ID:       "1",
		Name:     "Electronics",
		Slug:     "electronics",
		IsActive: true,
	}

	updatedCategory := &domain.Category{
		ID:       "1",
		Name:     "Electronics",
		Slug:     "electronics",
		ParentID: &categoryID, // Self as parent
		IsActive: true,
	}

	mockRepo.On("GetCategoryByID", ctx, "1").Return(existingCategory, nil)

	updated, err := service.UpdateCategory(ctx, updatedCategory)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot be its own parent")
	assert.Nil(t, updated)
	mockRepo.AssertNotCalled(t, "UpdateCategory")
	mockRepo.AssertExpectations(t)
}

// =============================================================================
// Test: DeleteCategory
// =============================================================================

func TestCategoryService_DeleteCategory_Success(t *testing.T) {
	service := setupTestCategoryService(t)
	mockRepo := new(MockCategoryRepository)
	service.repo = mockRepo

	ctx := context.Background()
	existingCategory := &domain.Category{
		ID:       "1",
		Name:     "Electronics",
		Slug:     "electronics",
		IsActive: true,
	}

	mockRepo.On("GetCategoryByID", ctx, "1").Return(existingCategory, nil)
	mockRepo.On("DeleteCategory", ctx, "1").Return(nil)

	err := service.DeleteCategory(ctx, "1")

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestCategoryService_DeleteCategory_NotFound(t *testing.T) {
	service := setupTestCategoryService(t)
	mockRepo := new(MockCategoryRepository)
	service.repo = mockRepo

	ctx := context.Background()

	mockRepo.On("GetCategoryByID", ctx, "999").Return(nil, errors.New("not found"))

	err := service.DeleteCategory(ctx, "999")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get category")
	mockRepo.AssertNotCalled(t, "DeleteCategory")
	mockRepo.AssertExpectations(t)
}

// =============================================================================
// Test: Validation
// =============================================================================

func TestCategoryService_Validation_NameTooShort(t *testing.T) {
	service := setupTestCategoryService(t)

	cat := &domain.Category{
		Name: "A",
		Slug: "a",
	}

	err := service.validateCategory(cat)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "at least 2 characters")
}

func TestCategoryService_Validation_NameTooLong(t *testing.T) {
	service := setupTestCategoryService(t)

	longName := string(make([]byte, 101))
	for i := range longName {
		longName = longName[:i] + "A" + longName[i+1:]
	}

	cat := &domain.Category{
		Name: longName,
		Slug: "test",
	}

	err := service.validateCategory(cat)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "at most 100 characters")
}

func TestCategoryService_Validation_Success(t *testing.T) {
	service := setupTestCategoryService(t)

	cat := &domain.Category{
		Name: "Valid Category Name",
		Slug: "valid-category",
	}

	err := service.validateCategory(cat)
	assert.NoError(t, err)
}
