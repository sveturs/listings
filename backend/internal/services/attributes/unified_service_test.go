package attributes

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"backend/internal/domain/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUnifiedAttributeStorage - мок для storage интерфейса
type MockUnifiedAttributeStorage struct {
	mock.Mock
}

func (m *MockUnifiedAttributeStorage) CreateAttribute(ctx context.Context, attr *models.UnifiedAttribute) (int, error) {
	args := m.Called(ctx, attr)
	return args.Int(0), args.Error(1)
}

func (m *MockUnifiedAttributeStorage) GetAttribute(ctx context.Context, id int) (*models.UnifiedAttribute, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.UnifiedAttribute), args.Error(1)
}

func (m *MockUnifiedAttributeStorage) GetAttributeByCode(ctx context.Context, code string) (*models.UnifiedAttribute, error) {
	args := m.Called(ctx, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.UnifiedAttribute), args.Error(1)
}

func (m *MockUnifiedAttributeStorage) ListAttributes(ctx context.Context, filter *models.UnifiedAttributeFilter) ([]*models.UnifiedAttribute, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.UnifiedAttribute), args.Error(1)
}

func (m *MockUnifiedAttributeStorage) UpdateAttribute(ctx context.Context, id int, updates map[string]interface{}) error {
	args := m.Called(ctx, id, updates)
	return args.Error(0)
}

func (m *MockUnifiedAttributeStorage) DeleteAttribute(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUnifiedAttributeStorage) AttachAttributeToCategory(ctx context.Context, categoryID, attributeID int, settings *models.UnifiedCategoryAttribute) error {
	args := m.Called(ctx, categoryID, attributeID, settings)
	return args.Error(0)
}

func (m *MockUnifiedAttributeStorage) DetachAttributeFromCategory(ctx context.Context, categoryID, attributeID int) error {
	args := m.Called(ctx, categoryID, attributeID)
	return args.Error(0)
}

func (m *MockUnifiedAttributeStorage) UpdateCategoryAttribute(ctx context.Context, categoryID, attributeID int, isRequired, isFilter *bool, displayOrder *int, groupID *int) error {
	args := m.Called(ctx, categoryID, attributeID, isRequired, isFilter, displayOrder, groupID)
	return args.Error(0)
}

func (m *MockUnifiedAttributeStorage) GetCategoryAttributes(ctx context.Context, categoryID int) ([]*models.UnifiedAttribute, error) {
	args := m.Called(ctx, categoryID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.UnifiedAttribute), args.Error(1)
}

func (m *MockUnifiedAttributeStorage) GetCategoryAttributesWithSettings(ctx context.Context, categoryID int) ([]*models.UnifiedCategoryAttribute, error) {
	args := m.Called(ctx, categoryID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.UnifiedCategoryAttribute), args.Error(1)
}

func (m *MockUnifiedAttributeStorage) SaveAttributeValue(ctx context.Context, value *models.UnifiedAttributeValue) error {
	args := m.Called(ctx, value)
	return args.Error(0)
}

func (m *MockUnifiedAttributeStorage) GetAttributeValues(ctx context.Context, entityType models.AttributeEntityType, entityID int) ([]*models.UnifiedAttributeValue, error) {
	args := m.Called(ctx, entityType, entityID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.UnifiedAttributeValue), args.Error(1)
}

func (m *MockUnifiedAttributeStorage) DeleteAttributeValues(ctx context.Context, entityType models.AttributeEntityType, entityID int) error {
	args := m.Called(ctx, entityType, entityID)
	return args.Error(0)
}

func (m *MockUnifiedAttributeStorage) GetAttributeByLegacyID(ctx context.Context, legacyID int, isProductVariant bool) (*models.UnifiedAttribute, error) {
	args := m.Called(ctx, legacyID, isProductVariant)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.UnifiedAttribute), args.Error(1)
}

func (m *MockUnifiedAttributeStorage) MigrateFromLegacySystem(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockUnifiedAttributeStorage) InvalidateCache(categoryID int) {
	m.Called(categoryID)
}

// Тесты для UnifiedAttributeService

func TestNewUnifiedAttributeService(t *testing.T) {
	mockStorage := new(MockUnifiedAttributeStorage)

	service := NewUnifiedAttributeService(mockStorage, true, true)

	assert.NotNil(t, service)
	assert.Equal(t, mockStorage, service.storage)
	assert.True(t, service.useLegacyFallback)
	assert.True(t, service.dualWrite)
}

func TestGetCategoryAttributes(t *testing.T) {
	ctx := context.Background()
	mockStorage := new(MockUnifiedAttributeStorage)
	service := NewUnifiedAttributeService(mockStorage, false, false)

	expectedAttrs := []*models.UnifiedAttribute{
		{
			ID:            1,
			Code:          "size",
			Name:          "Size",
			AttributeType: "select",
			Purpose:       models.PurposeRegular,
		},
		{
			ID:            2,
			Code:          "color",
			Name:          "Color",
			AttributeType: "select",
			Purpose:       models.PurposeRegular,
		},
	}

	mockStorage.On("GetCategoryAttributes", ctx, 1).Return(expectedAttrs, nil)

	attrs, err := service.GetCategoryAttributes(ctx, 1)

	assert.NoError(t, err)
	assert.Equal(t, expectedAttrs, attrs)
	mockStorage.AssertExpectations(t)
}

func TestGetCategoryAttributes_WithError(t *testing.T) {
	ctx := context.Background()
	mockStorage := new(MockUnifiedAttributeStorage)
	service := NewUnifiedAttributeService(mockStorage, false, false)

	expectedErr := errors.New("database error")
	mockStorage.On("GetCategoryAttributes", ctx, 1).Return(nil, expectedErr)

	attrs, err := service.GetCategoryAttributes(ctx, 1)

	assert.Error(t, err)
	assert.Nil(t, attrs)
	assert.Contains(t, err.Error(), "database error")
	mockStorage.AssertExpectations(t)
}

func TestSaveAttributeValues(t *testing.T) {
	ctx := context.Background()
	mockStorage := new(MockUnifiedAttributeStorage)
	service := NewUnifiedAttributeService(mockStorage, false, false)

	values := map[int]interface{}{
		1: "Large",
		2: "Red",
	}

	// Mock для сохранения новых значений
	mockStorage.On("SaveAttributeValue", ctx, mock.AnythingOfType("*models.UnifiedAttributeValue")).Return(nil).Twice()

	err := service.SaveAttributeValues(ctx, models.AttributeEntityType("listing"), 100, values)

	assert.NoError(t, err)
	mockStorage.AssertExpectations(t)
}

func TestValidateAttributeValue_Required(t *testing.T) {
	ctx := context.Background()
	mockStorage := new(MockUnifiedAttributeStorage)
	service := NewUnifiedAttributeService(mockStorage, false, false)

	attr := &models.UnifiedAttribute{
		ID:            1,
		Name:          "Size",
		IsRequired:    true,
		AttributeType: "text",
	}

	mockStorage.On("GetAttribute", ctx, 1).Return(attr, nil)

	// Тест с пустым значением для обязательного атрибута
	err := service.ValidateAttributeValue(ctx, 1, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "required")

	// Тест с заполненным значением
	mockStorage.On("GetAttribute", ctx, 1).Return(attr, nil)
	err = service.ValidateAttributeValue(ctx, 1, "Large")
	assert.NoError(t, err)
	mockStorage.AssertExpectations(t)
}

func TestValidateAttributeValue_SelectOptions(t *testing.T) {
	ctx := context.Background()
	mockStorage := new(MockUnifiedAttributeStorage)
	service := NewUnifiedAttributeService(mockStorage, false, false)

	attr := &models.UnifiedAttribute{
		ID:            1,
		Name:          "Size",
		AttributeType: "select",
		Options:       json.RawMessage(`["Small", "Medium", "Large"]`),
	}

	mockStorage.On("GetAttribute", ctx, 1).Return(attr, nil)

	// Тест с валидным значением
	err := service.ValidateAttributeValue(ctx, 1, "Medium")
	assert.NoError(t, err)

	// Тест с невалидным значением
	mockStorage.On("GetAttribute", ctx, 1).Return(attr, nil)
	err = service.ValidateAttributeValue(ctx, 1, "Extra Large")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not in allowed options")
	mockStorage.AssertExpectations(t)
}

func TestValidateAttributeValue_NumberRange(t *testing.T) {
	ctx := context.Background()
	mockStorage := new(MockUnifiedAttributeStorage)
	service := NewUnifiedAttributeService(mockStorage, false, false)

	attr := &models.UnifiedAttribute{
		ID:              1,
		Name:            "Price",
		AttributeType:   "number",
		ValidationRules: json.RawMessage(`{"min": 0, "max": 1000}`),
	}

	mockStorage.On("GetAttribute", ctx, 1).Return(attr, nil)

	// Тест с валидным значением
	err := service.ValidateAttributeValue(ctx, 1, float64(500))
	assert.NoError(t, err)

	// Тест со значением меньше минимума
	mockStorage.On("GetAttribute", ctx, 1).Return(attr, nil)
	err = service.ValidateAttributeValue(ctx, 1, float64(-10))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "must be at least")

	// Тест со значением больше максимума
	mockStorage.On("GetAttribute", ctx, 1).Return(attr, nil)
	err = service.ValidateAttributeValue(ctx, 1, float64(2000))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "must not exceed")
	mockStorage.AssertExpectations(t)
}

func TestAttachAttributeToCategory(t *testing.T) {
	ctx := context.Background()
	mockStorage := new(MockUnifiedAttributeStorage)
	service := NewUnifiedAttributeService(mockStorage, false, false)

	attr := &models.UnifiedAttribute{
		ID:   10,
		Code: "test_attr",
		Name: "Test Attribute",
	}

	settings := &models.UnifiedCategoryAttribute{
		IsEnabled:  true,
		IsRequired: true,
		SortOrder:  1,
	}

	// Mock для проверки существования атрибута
	mockStorage.On("GetAttribute", ctx, 10).Return(attr, nil)
	// Mock для привязки атрибута к категории
	mockStorage.On("AttachAttributeToCategory", ctx, 1, 10, mock.AnythingOfType("*models.UnifiedCategoryAttribute")).Return(nil)
	// Mock для инвалидации кеша
	mockStorage.On("InvalidateCache", 1).Return()

	err := service.AttachAttributeToCategory(ctx, 1, 10, settings)

	assert.NoError(t, err)
	mockStorage.AssertExpectations(t)
}

func TestDetachAttributeFromCategory(t *testing.T) {
	ctx := context.Background()
	mockStorage := new(MockUnifiedAttributeStorage)
	service := NewUnifiedAttributeService(mockStorage, false, false)

	mockStorage.On("DetachAttributeFromCategory", ctx, 1, 10).Return(nil)
	mockStorage.On("InvalidateCache", 1).Return()

	err := service.DetachAttributeFromCategory(ctx, 1, 10)

	assert.NoError(t, err)
	mockStorage.AssertExpectations(t)
}

func TestUpdateCategoryAttribute(t *testing.T) {
	ctx := context.Background()
	mockStorage := new(MockUnifiedAttributeStorage)
	service := NewUnifiedAttributeService(mockStorage, false, false)

	settings := &models.UnifiedCategoryAttribute{
		IsRequired: true,
		IsFilter:   true,
		SortOrder:  5,
	}

	isRequired := settings.IsRequired
	isFilter := settings.IsFilter
	sortOrder := settings.SortOrder

	mockStorage.On("UpdateCategoryAttribute", ctx, 1, 10, &isRequired, &isFilter, &sortOrder, (*int)(nil)).Return(nil)

	err := service.UpdateCategoryAttribute(ctx, 1, 10, settings)

	assert.NoError(t, err)
	mockStorage.AssertExpectations(t)
}

func TestMigrateFromLegacySystem(t *testing.T) {
	ctx := context.Background()
	mockStorage := new(MockUnifiedAttributeStorage)
	service := NewUnifiedAttributeService(mockStorage, false, false)

	mockStorage.On("MigrateFromLegacySystem", ctx).Return(nil)

	err := service.MigrateFromLegacySystem(ctx)

	assert.NoError(t, err)
	mockStorage.AssertExpectations(t)
}

func TestGetMigrationStatus(t *testing.T) {
	ctx := context.Background()
	service := &UnifiedAttributeService{}

	status, err := service.GetMigrationStatus(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, status)
	assert.Equal(t, "completed", status["status"])

	details, ok := status["details"].(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, 85, details["attributes_migrated"])
	assert.Equal(t, 14, details["categories_processed"])
	assert.Equal(t, 15, details["values_migrated"])
}
