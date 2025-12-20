// Package service implements business logic for the listings microservice.
package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/vondi-global/listings/internal/domain"
)

// MockAttributeRepository is a mock for AttributeRepository
type MockAttributeRepository struct {
	mock.Mock
}

func (m *MockAttributeRepository) Create(ctx context.Context, input *domain.CreateAttributeInput) (*domain.Attribute, error) {
	args := m.Called(ctx, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Attribute), args.Error(1)
}

func (m *MockAttributeRepository) Update(ctx context.Context, id int32, input *domain.UpdateAttributeInput) (*domain.Attribute, error) {
	args := m.Called(ctx, id, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Attribute), args.Error(1)
}

func (m *MockAttributeRepository) Delete(ctx context.Context, id int32) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockAttributeRepository) GetByID(ctx context.Context, id int32) (*domain.Attribute, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Attribute), args.Error(1)
}

func (m *MockAttributeRepository) GetByCode(ctx context.Context, code string) (*domain.Attribute, error) {
	args := m.Called(ctx, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Attribute), args.Error(1)
}

func (m *MockAttributeRepository) List(ctx context.Context, filter *domain.ListAttributesFilter) ([]*domain.Attribute, int64, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*domain.Attribute), args.Get(1).(int64), args.Error(2)
}

func (m *MockAttributeRepository) LinkToCategory(ctx context.Context, categoryID string, attributeID int32, settings *domain.CategoryAttributeSettings) (*domain.CategoryAttribute, error) {
	args := m.Called(ctx, categoryID, attributeID, settings)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.CategoryAttribute), args.Error(1)
}

func (m *MockAttributeRepository) UpdateCategoryAttribute(ctx context.Context, catAttrID int32, settings *domain.CategoryAttributeSettings) (*domain.CategoryAttribute, error) {
	args := m.Called(ctx, catAttrID, settings)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.CategoryAttribute), args.Error(1)
}

func (m *MockAttributeRepository) UnlinkFromCategory(ctx context.Context, categoryID string, attributeID int32) error {
	args := m.Called(ctx, categoryID, attributeID)
	return args.Error(0)
}

func (m *MockAttributeRepository) GetCategoryAttributes(ctx context.Context, categoryID string, filter *domain.GetCategoryAttributesFilter) ([]*domain.CategoryAttribute, error) {
	args := m.Called(ctx, categoryID, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.CategoryAttribute), args.Error(1)
}

func (m *MockAttributeRepository) GetListingValues(ctx context.Context, listingID int32) ([]*domain.ListingAttributeValue, error) {
	args := m.Called(ctx, listingID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.ListingAttributeValue), args.Error(1)
}

func (m *MockAttributeRepository) SetListingValues(ctx context.Context, listingID int32, values []domain.SetListingAttributeValue) error {
	args := m.Called(ctx, listingID, values)
	return args.Error(0)
}

func (m *MockAttributeRepository) DeleteListingValues(ctx context.Context, listingID int32) error {
	args := m.Called(ctx, listingID)
	return args.Error(0)
}

func (m *MockAttributeRepository) GetCategoryVariantAttributes(ctx context.Context, categoryID string) ([]*domain.VariantAttribute, error) {
	args := m.Called(ctx, categoryID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.VariantAttribute), args.Error(1)
}

func (m *MockAttributeRepository) GetVariantValues(ctx context.Context, variantID int32) ([]*domain.VariantAttributeValue, error) {
	args := m.Called(ctx, variantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.VariantAttributeValue), args.Error(1)
}

// Helper function to create a test service with mocks
func setupTestService(t *testing.T) (*AttributeServiceImpl, *MockAttributeRepository, redis.UniversalClient) {
	mockRepo := new(MockAttributeRepository)

	// Use miniredis for testing (in-memory Redis mock)
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("failed to start miniredis: %v", err)
	}

	// Clean up miniredis when test completes
	t.Cleanup(func() {
		mr.Close()
	})

	redisClient := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	logger := zerolog.New(nil).Level(zerolog.Disabled)

	// Use NewAttributeService constructor to ensure proper initialization
	service := NewAttributeService(mockRepo, redisClient, logger).(*AttributeServiceImpl)

	return service, mockRepo, redisClient
}

// Test CreateAttribute

func TestCreateAttribute_Success(t *testing.T) {
	service, mockRepo, _ := setupTestService(t)
	ctx := context.Background()

	input := &domain.CreateAttributeInput{
		Code:          "test_attr",
		Name:          map[string]string{"en": "Test Attribute"},
		DisplayName:   map[string]string{"en": "Test"},
		AttributeType: domain.AttributeTypeText,
		Purpose:       domain.AttributePurposeRegular,
	}

	expectedAttr := &domain.Attribute{
		ID:            1,
		Code:          "test_attr",
		Name:          map[string]string{"en": "Test Attribute"},
		DisplayName:   map[string]string{"en": "Test"},
		AttributeType: domain.AttributeTypeText,
		Purpose:       domain.AttributePurposeRegular,
		IsActive:      true,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Mock: Check code uniqueness (not found)
	mockRepo.On("GetByCode", ctx, "test_attr").Return(nil, errors.New("not found"))

	// Mock: Create attribute
	mockRepo.On("Create", ctx, input).Return(expectedAttr, nil)

	// Execute
	result, err := service.CreateAttribute(ctx, input)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedAttr.ID, result.ID)
	assert.Equal(t, expectedAttr.Code, result.Code)
	mockRepo.AssertExpectations(t)
}

func TestCreateAttribute_DuplicateCode(t *testing.T) {
	service, mockRepo, _ := setupTestService(t)
	ctx := context.Background()

	input := &domain.CreateAttributeInput{
		Code:          "existing_attr",
		Name:          map[string]string{"en": "Existing"},
		DisplayName:   map[string]string{"en": "Existing"},
		AttributeType: domain.AttributeTypeText,
	}

	existingAttr := &domain.Attribute{
		ID:   1,
		Code: "existing_attr",
	}

	// Mock: Code already exists
	mockRepo.On("GetByCode", ctx, "existing_attr").Return(existingAttr, nil)

	// Execute
	result, err := service.CreateAttribute(ctx, input)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "already exists")
	mockRepo.AssertExpectations(t)
}

func TestCreateAttribute_SelectWithoutOptions(t *testing.T) {
	service, mockRepo, _ := setupTestService(t)
	ctx := context.Background()

	input := &domain.CreateAttributeInput{
		Code:          "test_select",
		Name:          map[string]string{"en": "Test Select"},
		DisplayName:   map[string]string{"en": "Test Select"},
		AttributeType: domain.AttributeTypeSelect,
		Options:       []domain.AttributeOption{}, // Empty options
	}

	// Execute
	result, err := service.CreateAttribute(ctx, input)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "must have at least one option")
	mockRepo.AssertNotCalled(t, "Create")
}

// Test UpdateAttribute

func TestUpdateAttribute_Success(t *testing.T) {
	service, mockRepo, _ := setupTestService(t)
	ctx := context.Background()

	newName := map[string]string{"en": "Updated Name"}
	input := &domain.UpdateAttributeInput{
		Name: &newName,
	}

	existingAttr := &domain.Attribute{
		ID:            1,
		Code:          "test_attr",
		AttributeType: domain.AttributeTypeText,
		Name:          map[string]string{"en": "Old Name"},
	}

	updatedAttr := &domain.Attribute{
		ID:            1,
		Code:          "test_attr",
		AttributeType: domain.AttributeTypeText,
		Name:          newName,
	}

	// Mock: Get existing
	mockRepo.On("GetByID", ctx, int32(1)).Return(existingAttr, nil)

	// Mock: Update
	mockRepo.On("Update", ctx, int32(1), input).Return(updatedAttr, nil)

	// Execute
	result, err := service.UpdateAttribute(ctx, 1, input)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, newName, result.Name)
	mockRepo.AssertExpectations(t)
}

func TestUpdateAttribute_NotFound(t *testing.T) {
	service, mockRepo, _ := setupTestService(t)
	ctx := context.Background()

	newName := map[string]string{"en": "Updated"}
	input := &domain.UpdateAttributeInput{
		Name: &newName,
	}

	// Mock: Attribute not found
	mockRepo.On("GetByID", ctx, int32(999)).Return(nil, errors.New("not found"))

	// Execute
	result, err := service.UpdateAttribute(ctx, 999, input)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
	mockRepo.AssertNotCalled(t, "Update")
}

// Test DeleteAttribute

func TestDeleteAttribute_Success(t *testing.T) {
	service, mockRepo, _ := setupTestService(t)
	ctx := context.Background()

	existingAttr := &domain.Attribute{
		ID:   1,
		Code: "test_attr",
	}

	// Mock: Get existing
	mockRepo.On("GetByID", ctx, int32(1)).Return(existingAttr, nil)

	// Mock: Delete
	mockRepo.On("Delete", ctx, int32(1)).Return(nil)

	// Execute
	err := service.DeleteAttribute(ctx, 1)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// Test GetAttributeByID

func TestGetAttributeByID_CacheMiss_Success(t *testing.T) {
	service, mockRepo, _ := setupTestService(t)
	ctx := context.Background()

	expectedAttr := &domain.Attribute{
		ID:   1,
		Code: "test_attr",
	}

	// Mock: Repository fetch (cache miss)
	mockRepo.On("GetByID", ctx, int32(1)).Return(expectedAttr, nil)

	// Execute
	result, err := service.GetAttributeByID(ctx, 1)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedAttr.ID, result.ID)
	mockRepo.AssertExpectations(t)
}

func TestGetAttributeByCode_Success(t *testing.T) {
	service, mockRepo, _ := setupTestService(t)
	ctx := context.Background()

	expectedAttr := &domain.Attribute{
		ID:   1,
		Code: "test_attr",
	}

	// Mock: Repository fetch (cache will miss and call repository)
	mockRepo.On("GetByCode", ctx, "test_attr").Return(expectedAttr, nil).Once()

	// Execute
	result, err := service.GetAttributeByCode(ctx, "test_attr")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedAttr.Code, result.Code)
	mockRepo.AssertExpectations(t)
}

// Test ListAttributes

func TestListAttributes_Success(t *testing.T) {
	service, mockRepo, _ := setupTestService(t)
	ctx := context.Background()

	filter := &domain.ListAttributesFilter{
		Limit:  10,
		Offset: 0,
	}

	expectedAttrs := []*domain.Attribute{
		{ID: 1, Code: "attr1"},
		{ID: 2, Code: "attr2"},
	}

	// Mock: List
	mockRepo.On("List", ctx, filter).Return(expectedAttrs, int64(2), nil)

	// Execute
	results, total, err := service.ListAttributes(ctx, filter)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, results, 2)
	assert.Equal(t, int64(2), total)
	mockRepo.AssertExpectations(t)
}

func TestListAttributes_WithDefaults(t *testing.T) {
	service, mockRepo, _ := setupTestService(t)
	ctx := context.Background()

	// Mock: List with default filter
	mockRepo.On("List", ctx, mock.MatchedBy(func(f *domain.ListAttributesFilter) bool {
		return f.Limit == 20 && f.Offset == 0
	})).Return([]*domain.Attribute{}, int64(0), nil)

	// Execute with nil filter
	results, total, err := service.ListAttributes(ctx, nil)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, results, 0)
	assert.Equal(t, int64(0), total)
	mockRepo.AssertExpectations(t)
}

// Test LinkAttributeToCategory

func TestLinkAttributeToCategory_Success(t *testing.T) {
	service, mockRepo, _ := setupTestService(t)
	ctx := context.Background()

	settings := &domain.CategoryAttributeSettings{
		IsEnabled: true,
		SortOrder: 10,
	}

	attr := &domain.Attribute{
		ID:   1,
		Code: "test_attr",
	}

	catAttr := &domain.CategoryAttribute{
		ID:          1,
		CategoryID:  "3b4246cc-9970-403c-af01-c142a4178dc6",
		AttributeID: 1,
		IsEnabled:   true,
		SortOrder:   10,
	}

	// Mock: Get attribute (cache miss will call repository)
	// Note: May be called multiple times due to caching attempts
	mockRepo.On("GetByID", ctx, int32(1)).Return(attr, nil)

	// Mock: Link
	mockRepo.On("LinkToCategory", ctx, "3b4246cc-9970-403c-af01-c142a4178dc6", int32(1), settings).Return(catAttr, nil).Once()

	// Execute
	err := service.LinkAttributeToCategory(ctx, "3b4246cc-9970-403c-af01-c142a4178dc6", 1, settings)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestLinkAttributeToCategory_AttributeNotFound(t *testing.T) {
	service, mockRepo, _ := setupTestService(t)
	ctx := context.Background()

	settings := &domain.CategoryAttributeSettings{
		IsEnabled: true,
	}

	// Mock: Attribute not found
	mockRepo.On("GetByID", ctx, int32(999)).Return(nil, errors.New("not found"))

	// Execute
	err := service.LinkAttributeToCategory(ctx, "3b4246cc-9970-403c-af01-c142a4178dc6", 999, settings)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
	mockRepo.AssertExpectations(t)
	mockRepo.AssertNotCalled(t, "LinkToCategory")
}

// Test GetCategoryAttributes

func TestGetCategoryAttributes_CacheMiss_Success(t *testing.T) {
	service, mockRepo, _ := setupTestService(t)
	ctx := context.Background()

	expectedAttrs := []*domain.CategoryAttribute{
		{
			ID:          1,
			CategoryID:  "3b4246cc-9970-403c-af01-c142a4178dc6",
			AttributeID: 1,
			Attribute: &domain.Attribute{
				ID:   1,
				Code: "attr1",
			},
		},
	}

	// Mock: Repository fetch
	mockRepo.On("GetCategoryAttributes", ctx, "3b4246cc-9970-403c-af01-c142a4178dc6", (*domain.GetCategoryAttributesFilter)(nil)).
		Return(expectedAttrs, nil)

	// Execute
	results, err := service.GetCategoryAttributes(ctx, "3b4246cc-9970-403c-af01-c142a4178dc6", nil)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, results, 1)
	mockRepo.AssertExpectations(t)
}

// Test ValidateAttributeValues

func TestValidateAttributeValues_Success(t *testing.T) {
	service, mockRepo, _ := setupTestService(t)
	ctx := context.Background()

	textValue := "Test Value"
	numberValue := 42.5

	values := []domain.SetListingAttributeValue{
		{
			AttributeID: 1,
			ValueText:   &textValue,
		},
		{
			AttributeID: 2,
			ValueNumber: &numberValue,
		},
	}

	catAttrs := []*domain.CategoryAttribute{
		{
			ID:          1,
			CategoryID:  "3b4246cc-9970-403c-af01-c142a4178dc6",
			AttributeID: 1,
			IsEnabled:   true,
			Attribute: &domain.Attribute{
				ID:              1,
				Code:            "text_attr",
				AttributeType:   domain.AttributeTypeText,
				ValidationRules: map[string]interface{}{},
				Options:         []domain.AttributeOption{},
			},
		},
		{
			ID:          2,
			CategoryID:  "3b4246cc-9970-403c-af01-c142a4178dc6",
			AttributeID: 2,
			IsEnabled:   true,
			Attribute: &domain.Attribute{
				ID:              2,
				Code:            "number_attr",
				AttributeType:   domain.AttributeTypeNumber,
				ValidationRules: map[string]interface{}{},
				Options:         []domain.AttributeOption{},
			},
		},
	}

	// Mock: Get category attributes (cache miss will call repository)
	mockRepo.On("GetCategoryAttributes", ctx, "3b4246cc-9970-403c-af01-c142a4178dc6", (*domain.GetCategoryAttributesFilter)(nil)).
		Return(catAttrs, nil).Once()

	// Execute
	err := service.ValidateAttributeValues(ctx, "3b4246cc-9970-403c-af01-c142a4178dc6", values)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestValidateAttributeValues_AttributeNotLinked(t *testing.T) {
	service, mockRepo, _ := setupTestService(t)
	ctx := context.Background()

	textValue := "Test"
	values := []domain.SetListingAttributeValue{
		{
			AttributeID: 999, // Not linked to category
			ValueText:   &textValue,
		},
	}

	catAttrs := []*domain.CategoryAttribute{
		{
			ID:          1,
			CategoryID:  "3b4246cc-9970-403c-af01-c142a4178dc6",
			AttributeID: 1,
			IsEnabled:   true,
			Attribute: &domain.Attribute{
				ID:   1,
				Code: "attr1",
			},
		},
	}

	// Mock: Get category attributes (cache will miss and call repository)
	mockRepo.On("GetCategoryAttributes", ctx, "3b4246cc-9970-403c-af01-c142a4178dc6", (*domain.GetCategoryAttributesFilter)(nil)).
		Return(catAttrs, nil)

	// Execute
	err := service.ValidateAttributeValues(ctx, "3b4246cc-9970-403c-af01-c142a4178dc6", values)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not linked to category")
	mockRepo.AssertExpectations(t)
}

func TestValidateAttributeValues_MissingRequiredAttribute(t *testing.T) {
	service, mockRepo, _ := setupTestService(t)
	ctx := context.Background()

	values := []domain.SetListingAttributeValue{} // Empty - missing required

	isRequired := true
	catAttrs := []*domain.CategoryAttribute{
		{
			ID:          1,
			CategoryID:  "3b4246cc-9970-403c-af01-c142a4178dc6",
			AttributeID: 1,
			IsEnabled:   true,
			IsRequired:  &isRequired, // Required!
			Attribute: &domain.Attribute{
				ID:            1,
				Code:          "required_attr",
				AttributeType: domain.AttributeTypeText,
			},
		},
	}

	// Mock: Get category attributes (cache will miss and call repository)
	mockRepo.On("GetCategoryAttributes", ctx, "3b4246cc-9970-403c-af01-c142a4178dc6", (*domain.GetCategoryAttributesFilter)(nil)).
		Return(catAttrs, nil)

	// Execute
	err := service.ValidateAttributeValues(ctx, "3b4246cc-9970-403c-af01-c142a4178dc6", values)

	// Assert
	assert.Error(t, err)
	if err != nil {
		assert.Contains(t, err.Error(), "required attribute")
		assert.Contains(t, err.Error(), "missing")
	}
	mockRepo.AssertExpectations(t)
}

// Test SetListingAttributes

func TestSetListingAttributes_Success(t *testing.T) {
	service, mockRepo, _ := setupTestService(t)
	ctx := context.Background()

	textValue := "Test Value"
	values := []domain.SetListingAttributeValue{
		{
			AttributeID: 1,
			ValueText:   &textValue,
		},
	}

	attr := &domain.Attribute{
		ID:              1,
		Code:            "test_attr",
		AttributeType:   domain.AttributeTypeText, // IMPORTANT: Must set type!
		ValidationRules: map[string]interface{}{},
		Options:         []domain.AttributeOption{},
	}

	// Mock: Get attribute for validation (cache will miss and call repository)
	mockRepo.On("GetByID", ctx, int32(1)).Return(attr, nil)

	// Mock: Set values
	mockRepo.On("SetListingValues", ctx, int32(100), values).Return(nil).Once()

	// Execute
	err := service.SetListingAttributes(ctx, 100, values)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestSetListingAttributes_InvalidAttributeID(t *testing.T) {
	service, mockRepo, _ := setupTestService(t)
	ctx := context.Background()

	textValue := "Test"
	values := []domain.SetListingAttributeValue{
		{
			AttributeID: 999, // Invalid
			ValueText:   &textValue,
		},
	}

	// Mock: Attribute not found
	mockRepo.On("GetByID", ctx, int32(999)).Return(nil, errors.New("not found"))

	// Execute
	err := service.SetListingAttributes(ctx, 100, values)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid attribute_id")
	mockRepo.AssertExpectations(t)
	mockRepo.AssertNotCalled(t, "SetListingValues")
}

// Test GetCategoryVariantAttributes

func TestGetCategoryVariantAttributes_Success(t *testing.T) {
	service, mockRepo, _ := setupTestService(t)
	ctx := context.Background()

	expectedAttrs := []*domain.VariantAttribute{
		{
			ID:          1,
			CategoryID:  "3b4246cc-9970-403c-af01-c142a4178dc6",
			AttributeID: 1,
		},
	}

	// Mock: Repository fetch
	mockRepo.On("GetCategoryVariantAttributes", ctx, "3b4246cc-9970-403c-af01-c142a4178dc6").Return(expectedAttrs, nil)

	// Execute
	results, err := service.GetCategoryVariantAttributes(ctx, "3b4246cc-9970-403c-af01-c142a4178dc6")

	// Assert
	assert.NoError(t, err)
	assert.Len(t, results, 1)
	mockRepo.AssertExpectations(t)
}

// Test Cache Invalidation

func TestInvalidateAttributeCache_Success(t *testing.T) {
	service, mockRepo, _ := setupTestService(t)
	ctx := context.Background()

	attr := &domain.Attribute{
		ID:   1,
		Code: "test_attr",
	}

	// Mock: Get attribute
	mockRepo.On("GetByID", ctx, int32(1)).Return(attr, nil)

	// Execute
	err := service.InvalidateAttributeCache(ctx, 1)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestInvalidateCategoryCache_Success(t *testing.T) {
	service, _, _ := setupTestService(t)
	ctx := context.Background()

	// Execute
	err := service.InvalidateCategoryCache(ctx, "3b4246cc-9970-403c-af01-c142a4178dc6")

	// Assert
	assert.NoError(t, err)
}

func TestInvalidateListingCache_Success(t *testing.T) {
	service, _, _ := setupTestService(t)
	ctx := context.Background()

	// Execute
	err := service.InvalidateListingCache(ctx, 100)

	// Assert
	assert.NoError(t, err)
}
