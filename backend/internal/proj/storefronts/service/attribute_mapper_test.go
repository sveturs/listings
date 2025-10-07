// backend/internal/proj/storefronts/service/attribute_mapper_test.go
package service

import (
	"context"
	"testing"

	"backend/internal/domain/models"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

// mockStorage заглушка для Storage (реализует все методы интерфейса)
type mockStorageForAttributeMapper struct{}

// Все методы интерфейса Storage (заглушки)
func (m *mockStorageForAttributeMapper) GetStorefrontProducts(ctx context.Context, filter models.ProductFilter) ([]*models.StorefrontProduct, error) {
	return nil, nil
}

func (m *mockStorageForAttributeMapper) GetStorefrontProduct(ctx context.Context, storefrontID, productID int) (*models.StorefrontProduct, error) {
	return nil, nil
}

func (m *mockStorageForAttributeMapper) GetStorefrontProductByID(ctx context.Context, productID int) (*models.StorefrontProduct, error) {
	return nil, nil
}

func (m *mockStorageForAttributeMapper) GetStorefrontProductBySKU(ctx context.Context, storefrontID int, sku string) (*models.StorefrontProduct, error) {
	return nil, nil
}

func (m *mockStorageForAttributeMapper) GetStorefrontProductsBySKUs(ctx context.Context, storefrontID int, skus []string) (map[string]*models.StorefrontProduct, error) {
	return nil, nil
}

func (m *mockStorageForAttributeMapper) CreateStorefrontProduct(ctx context.Context, storefrontID int, req *models.CreateProductRequest) (*models.StorefrontProduct, error) {
	return nil, nil
}

func (m *mockStorageForAttributeMapper) BatchCreateStorefrontProducts(ctx context.Context, storefrontID int, requests []*models.CreateProductRequest) ([]*models.StorefrontProduct, error) {
	return nil, nil
}

func (m *mockStorageForAttributeMapper) UpdateStorefrontProduct(ctx context.Context, storefrontID, productID int, req *models.UpdateProductRequest) error {
	return nil
}

func (m *mockStorageForAttributeMapper) DeleteStorefrontProduct(ctx context.Context, storefrontID, productID int) error {
	return nil
}

func (m *mockStorageForAttributeMapper) HardDeleteStorefrontProduct(ctx context.Context, storefrontID, productID int) error {
	return nil
}

func (m *mockStorageForAttributeMapper) UpdateProductInventory(ctx context.Context, storefrontID, productID int, userID int, req *models.UpdateInventoryRequest) error {
	return nil
}

func (m *mockStorageForAttributeMapper) GetProductStats(ctx context.Context, storefrontID int) (*models.ProductStats, error) {
	return nil, nil
}

func (m *mockStorageForAttributeMapper) CreateProductVariant(ctx context.Context, variant *models.CreateProductVariantRequest) (*models.StorefrontProductVariant, error) {
	return nil, nil
}

func (m *mockStorageForAttributeMapper) BatchCreateProductVariants(ctx context.Context, variants []*models.CreateProductVariantRequest) ([]*models.StorefrontProductVariant, error) {
	return nil, nil
}

func (m *mockStorageForAttributeMapper) CreateProductVariantImage(ctx context.Context, image *models.CreateProductVariantImageRequest) (*models.StorefrontProductVariantImage, error) {
	return nil, nil
}

func (m *mockStorageForAttributeMapper) BatchCreateProductVariantImages(ctx context.Context, images []*models.CreateProductVariantImageRequest) ([]*models.StorefrontProductVariantImage, error) {
	return nil, nil
}

func (m *mockStorageForAttributeMapper) GetProductVariants(ctx context.Context, productID int) ([]*models.StorefrontProductVariant, error) {
	return nil, nil
}

func (m *mockStorageForAttributeMapper) IncrementProductViews(ctx context.Context, productID int) error {
	return nil
}

func (m *mockStorageForAttributeMapper) BulkCreateProducts(ctx context.Context, storefrontID int, products []models.CreateProductRequest) ([]int, []error) {
	return nil, nil
}

func (m *mockStorageForAttributeMapper) BulkUpdateProducts(ctx context.Context, storefrontID int, updates []models.BulkUpdateItem) ([]int, []error) {
	return nil, nil
}

func (m *mockStorageForAttributeMapper) BulkDeleteProducts(ctx context.Context, storefrontID int, productIDs []int) ([]int, []error) {
	return nil, nil
}

func (m *mockStorageForAttributeMapper) BulkUpdateStatus(ctx context.Context, storefrontID int, productIDs []int, isActive bool) ([]int, []error) {
	return nil, nil
}

func (m *mockStorageForAttributeMapper) GetStorefrontByID(ctx context.Context, id int) (*models.Storefront, error) {
	return nil, nil
}

func (m *mockStorageForAttributeMapper) SaveTranslation(ctx context.Context, translation *models.Translation) error {
	return nil
}

func (m *mockStorageForAttributeMapper) BeginTx(ctx context.Context) (Transaction, error) {
	return nil, nil
}

func (m *mockStorageForAttributeMapper) CreateStorefrontProductTx(ctx context.Context, tx Transaction, storefrontID int, req *models.CreateProductRequest) (*models.StorefrontProduct, error) {
	return nil, nil
}

// newTestAttributeMapper создает тестовый AttributeMapper
func newTestAttributeMapper() *AttributeMapper {
	logger := zerolog.Nop()
	storage := &mockStorageForAttributeMapper{}

	mapper := NewAttributeMapper(storage, logger)

	// Инициализируем кэш вручную для тестов
	mapper.attributesCache = map[string]*AttributeTemplate{
		"brand": {
			ID:            1,
			Code:          "brand",
			Name:          "Brand",
			DisplayName:   "Бренд",
			AttributeType: "text",
			Purpose:       "regular",
			IsSearchable:  true,
			IsFilterable:  true,
		},
		"color": {
			ID:            2,
			Code:          "color",
			Name:          "Color",
			DisplayName:   "Цвет",
			AttributeType: "select",
			Purpose:       "variant",
			IsFilterable:  true,
		},
		"price": {
			ID:            3,
			Code:          "price",
			Name:          "Price",
			DisplayName:   "Цена",
			AttributeType: "number",
			Purpose:       "regular",
		},
		"warranty": {
			ID:            4,
			Code:          "warranty",
			Name:          "Warranty",
			DisplayName:   "Гарантия",
			AttributeType: "boolean",
			Purpose:       "regular",
		},
	}

	return mapper
}

// TestNormalizeAttributeName тестирует нормализацию имени атрибута
func TestNormalizeAttributeName(t *testing.T) {
	mapper := newTestAttributeMapper()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Lowercase conversion",
			input:    "BRAND",
			expected: "brand",
		},
		{
			name:     "Trim spaces",
			input:    "  color  ",
			expected: "color",
		},
		{
			name:     "Replace underscores",
			input:    "product_weight",
			expected: "product weight",
		},
		{
			name:     "Multiple spaces",
			input:    "product   name",
			expected: "product name",
		},
		{
			name:     "Mixed",
			input:    "  Product_NAME  ",
			expected: "product name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mapper.normalizeAttributeName(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestMapExternalAttribute_DirectMatch тестирует прямое совпадение атрибута
func TestMapExternalAttribute_DirectMatch(t *testing.T) {
	mapper := newTestAttributeMapper()
	ctx := context.Background()

	result, err := mapper.MapExternalAttribute(ctx, "brand", "Samsung", nil)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, result.AttributeID)
	assert.Equal(t, "brand", result.Code)
	assert.Equal(t, "Samsung", result.Value)
	assert.Equal(t, 1.0, result.Confidence)
	assert.False(t, result.IsNewAttribute)
}

// TestMapExternalAttribute_CaseInsensitive тестирует case-insensitive поиск
func TestMapExternalAttribute_CaseInsensitive(t *testing.T) {
	mapper := newTestAttributeMapper()
	ctx := context.Background()

	result, err := mapper.MapExternalAttribute(ctx, "BRAND", "Sony", nil)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, result.AttributeID)
	assert.Equal(t, "brand", result.Code)
	assert.Equal(t, "Sony", result.Value)
	assert.Equal(t, 1.0, result.Confidence)
}

// TestMapExternalAttribute_NotFound тестирует ненайденный атрибут
func TestMapExternalAttribute_NotFound(t *testing.T) {
	mapper := newTestAttributeMapper()
	ctx := context.Background()

	result, err := mapper.MapExternalAttribute(ctx, "unknown_attribute", "value", nil)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.IsNewAttribute)
	assert.Equal(t, "unknown_attribute", result.SuggestedCode)
	assert.Equal(t, "unknown_attribute", result.Name)
	assert.Equal(t, 0.0, result.Confidence)
}

// TestTransformValue_Number тестирует трансформацию в число
func TestTransformValue_Number(t *testing.T) {
	mapper := newTestAttributeMapper()

	tests := []struct {
		name     string
		input    interface{}
		expected interface{}
	}{
		{
			name:     "Integer",
			input:    42,
			expected: 42,
		},
		{
			name:     "Float",
			input:    42.5,
			expected: 42.5,
		},
		{
			name:     "String number",
			input:    "123.45",
			expected: 123.45,
		},
		{
			name:     "Invalid string",
			input:    "not a number",
			expected: "not a number",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mapper.transformValue(tt.input, "number")
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestTransformValue_Boolean тестирует трансформацию в boolean
func TestTransformValue_Boolean(t *testing.T) {
	mapper := newTestAttributeMapper()

	tests := []struct {
		name     string
		input    interface{}
		expected interface{}
	}{
		{
			name:     "Boolean true",
			input:    true,
			expected: true,
		},
		{
			name:     "Boolean false",
			input:    false,
			expected: false,
		},
		{
			name:     "String true",
			input:    "true",
			expected: true,
		},
		{
			name:     "String yes",
			input:    "yes",
			expected: true,
		},
		{
			name:     "String да (Russian)",
			input:    "да",
			expected: true,
		},
		{
			name:     "String 1",
			input:    "1",
			expected: true,
		},
		{
			name:     "String false",
			input:    "false",
			expected: false,
		},
		{
			name:     "Integer 1",
			input:    1,
			expected: true,
		},
		{
			name:     "Integer 0",
			input:    0,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mapper.transformValue(tt.input, "boolean")
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestTransformValue_Text тестирует трансформацию в текст
func TestTransformValue_Text(t *testing.T) {
	mapper := newTestAttributeMapper()

	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{
			name:     "String",
			input:    "Hello",
			expected: "Hello",
		},
		{
			name:     "Integer",
			input:    123,
			expected: "123",
		},
		{
			name:     "Float",
			input:    123.45,
			expected: "123.45",
		},
		{
			name:     "Boolean",
			input:    true,
			expected: "true",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mapper.transformValue(tt.input, "text")
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestGenerateAttributeCode тестирует генерацию кода атрибута
func TestGenerateAttributeCode(t *testing.T) {
	mapper := newTestAttributeMapper()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Simple name",
			input:    "Product Name",
			expected: "product_name",
		},
		{
			name:     "With special chars",
			input:    "Product-Name (2023)",
			expected: "product_name_2023",
		},
		{
			name:     "Multiple spaces",
			input:    "Product   Name",
			expected: "product_name",
		},
		{
			name:     "Cyrillic",
			input:    "Название товара",
			expected: "",
		},
		{
			name:     "Mixed",
			input:    "  Product_NAME-123  ",
			expected: "product_name_123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mapper.generateAttributeCode(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestBatchMapAttributes тестирует пакетный маппинг
func TestBatchMapAttributes(t *testing.T) {
	mapper := newTestAttributeMapper()
	ctx := context.Background()

	requests := []AttributeMappingRequest{
		{
			ExternalName:  "brand",
			ExternalValue: "Samsung",
		},
		{
			ExternalName:  "color",
			ExternalValue: "black",
		},
		{
			ExternalName:  "price",
			ExternalValue: "999.99",
		},
		{
			ExternalName:  "unknown_attr",
			ExternalValue: "some value",
		},
	}

	results, err := mapper.BatchMapAttributes(ctx, requests)

	assert.NoError(t, err)
	assert.NotNil(t, results)
	assert.Equal(t, 4, len(results))

	// Проверяем первые 3 атрибута (должны быть найдены)
	assert.Equal(t, "brand", results[0].Code)
	assert.False(t, results[0].IsNewAttribute)

	assert.Equal(t, "color", results[1].Code)
	assert.False(t, results[1].IsNewAttribute)

	assert.Equal(t, "price", results[2].Code)
	assert.False(t, results[2].IsNewAttribute)
	assert.Equal(t, 999.99, results[2].Value) // должно быть преобразовано в число

	// Проверяем неизвестный атрибут
	assert.True(t, results[3].IsNewAttribute)
	assert.Equal(t, "unknown_attr", results[3].SuggestedCode)
}

// TestCalculateConfidence тестирует вычисление уверенности
func TestCalculateConfidence(t *testing.T) {
	mapper := newTestAttributeMapper()

	template := &AttributeTemplate{
		ID:            1,
		Code:          "brand",
		Name:          "Brand",
		AttributeType: "text",
	}

	tests := []struct {
		name           string
		normalizedName string
		expectedMin    float64
		expectedMax    float64
	}{
		{
			name:           "Exact code match",
			normalizedName: "brand",
			expectedMin:    1.0,
			expectedMax:    1.0,
		},
		{
			name:           "Name match",
			normalizedName: "brand",
			expectedMin:    0.95,
			expectedMax:    1.0,
		},
		{
			name:           "Partial match",
			normalizedName: "product brand",
			expectedMin:    0.7,
			expectedMax:    0.9,
		},
		{
			name:           "No match",
			normalizedName: "color",
			expectedMin:    0.0,
			expectedMax:    0.7,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			confidence := mapper.calculateConfidence(tt.normalizedName, template, "test value")
			assert.GreaterOrEqual(t, confidence, tt.expectedMin)
			assert.LessOrEqual(t, confidence, tt.expectedMax)
		})
	}
}

// TestMappingCache тестирует работу кэша маппинга
func TestMappingCache(t *testing.T) {
	mapper := newTestAttributeMapper()
	ctx := context.Background()

	// Первый вызов - маппинг и сохранение в кэш
	result1, err := mapper.MapExternalAttribute(ctx, "brand", "Samsung", nil)
	assert.NoError(t, err)
	assert.NotNil(t, result1)

	// Проверяем что кэш заполнен
	assert.NotEmpty(t, mapper.mappingCache)

	// Второй вызов - должен использовать кэш
	result2, err := mapper.MapExternalAttribute(ctx, "brand", "Sony", nil)
	assert.NoError(t, err)
	assert.NotNil(t, result2)

	// Результаты должны быть одинаковыми (кроме значения)
	assert.Equal(t, result1.AttributeID, result2.AttributeID)
	assert.Equal(t, result1.Code, result2.Code)
	assert.Equal(t, result1.Name, result2.Name)
}
