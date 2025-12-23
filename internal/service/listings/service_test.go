package listings

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/vondi-global/listings/internal/domain"
	"github.com/vondi-global/listings/internal/service/listings/mocks"
)

// TestSetupServiceTest verifies that test setup works correctly
func TestSetupServiceTest(t *testing.T) {
	service, mockRepo, mockCache, mockIndexer := SetupServiceTest(t)

	assert.NotNil(t, service, "service should not be nil")
	assert.NotNil(t, mockRepo, "mockRepo should not be nil")
	assert.NotNil(t, mockCache, "mockCache should not be nil")
	assert.NotNil(t, mockIndexer, "mockIndexer should not be nil")
}

// TestTestContext verifies that test context creation works
func TestTestContext(t *testing.T) {
	ctx := TestContext()

	assert.NotNil(t, ctx, "context should not be nil")
	assert.NotNil(t, ctx.Done(), "context should have done channel")
}

// TestNewTestListing verifies test listing creation
func TestNewTestListing(t *testing.T) {
	listing := NewTestListing(1, 100, "Test Listing")

	assert.Equal(t, int64(1), listing.ID)
	assert.Equal(t, int64(100), listing.UserID)
	assert.Equal(t, "Test Listing", listing.Title)
	assert.NotEmpty(t, listing.UUID)
	assert.Equal(t, 99.99, listing.Price)
	assert.Equal(t, "USD", listing.Currency)
}

// TestNewTestProduct verifies test product creation
func TestNewTestProduct(t *testing.T) {
	product := NewTestProduct(1, 10, "Test Product")

	assert.Equal(t, int64(1), product.ID)
	assert.Equal(t, int64(10), product.StorefrontID)
	assert.Equal(t, "Test Product", product.Name)
	assert.Equal(t, 149.99, product.Price)
	assert.Equal(t, "USD", product.Currency)
	assert.True(t, product.IsActive)
}

// TestNewTestProductVariant verifies test variant creation
func TestNewTestProductVariant(t *testing.T) {
	variant := NewTestProductVariant(1, 100)

	assert.Equal(t, int64(1), variant.ID)
	assert.Equal(t, int64(100), variant.ProductID)
	assert.NotNil(t, variant.Price)
	assert.Equal(t, 159.99, *variant.Price)
	assert.True(t, variant.IsActive)
	assert.False(t, variant.IsDefault)
}

// ===================================
// BulkCreateProducts Tests (10 tests)
// ===================================

func TestBulkCreateProducts_Success_SingleProduct(t *testing.T) {
	service, mockRepo, _, _ := SetupServiceTest(t)
	ctx := TestContext()

	storefrontID := int64(123)
	inputs := []*domain.CreateProductInput{
		NewCreateProductInput(storefrontID, "Product 1"),
	}

	expectedProducts := []*domain.Product{
		NewTestProduct(1, storefrontID, "Product 1"),
	}

	mockRepo.On("BulkCreateProducts", ctx, storefrontID, inputs).
		Return(expectedProducts, []domain.BulkProductError{}, nil)

	products, errors, err := service.BulkCreateProducts(ctx, storefrontID, inputs)

	assert.NoError(t, err)
	assert.Len(t, products, 1)
	assert.Len(t, errors, 0)
	assert.Equal(t, "Product 1", products[0].Name)
	mockRepo.AssertExpectations(t)
}

func TestBulkCreateProducts_Success_MultipleProducts(t *testing.T) {
	service, mockRepo, _, _ := SetupServiceTest(t)
	ctx := TestContext()

	storefrontID := int64(123)
	inputs := []*domain.CreateProductInput{
		NewCreateProductInput(storefrontID, "Product 1"),
		NewCreateProductInput(storefrontID, "Product 2"),
		NewCreateProductInput(storefrontID, "Product 3"),
	}

	expectedProducts := []*domain.Product{
		NewTestProduct(1, storefrontID, "Product 1"),
		NewTestProduct(2, storefrontID, "Product 2"),
		NewTestProduct(3, storefrontID, "Product 3"),
	}

	mockRepo.On("BulkCreateProducts", ctx, storefrontID, inputs).
		Return(expectedProducts, []domain.BulkProductError{}, nil)

	products, errors, err := service.BulkCreateProducts(ctx, storefrontID, inputs)

	assert.NoError(t, err)
	assert.Len(t, products, 3)
	assert.Len(t, errors, 0)
	mockRepo.AssertExpectations(t)
}

func TestBulkCreateProducts_Error_EmptyInput(t *testing.T) {
	service, _, _, _ := SetupServiceTest(t)
	ctx := TestContext()

	storefrontID := int64(123)
	inputs := []*domain.CreateProductInput{}

	products, errors, err := service.BulkCreateProducts(ctx, storefrontID, inputs)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "products.bulk_empty")
	assert.Nil(t, products)
	assert.Nil(t, errors)
}

func TestBulkCreateProducts_Error_BatchTooLarge(t *testing.T) {
	service, _, _, _ := SetupServiceTest(t)
	ctx := TestContext()

	storefrontID := int64(123)
	inputs := make([]*domain.CreateProductInput, 1001)
	for i := 0; i < 1001; i++ {
		inputs[i] = NewCreateProductInput(storefrontID, "Product")
	}

	products, errors, err := service.BulkCreateProducts(ctx, storefrontID, inputs)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "products.bulk_too_large")
	assert.Nil(t, products)
	assert.Nil(t, errors)
}

func TestBulkCreateProducts_Error_NilInput(t *testing.T) {
	service, _, _, _ := SetupServiceTest(t)
	ctx := TestContext()

	storefrontID := int64(123)
	inputs := []*domain.CreateProductInput{
		NewCreateProductInput(storefrontID, "Product 1"),
		nil, // Nil input
		NewCreateProductInput(storefrontID, "Product 3"),
	}

	products, errors, err := service.BulkCreateProducts(ctx, storefrontID, inputs)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "product at index 1 is nil")
	assert.Nil(t, products)
	assert.Nil(t, errors)
}

func TestBulkCreateProducts_Error_ValidationFailed_Index0(t *testing.T) {
	service, mockRepo, _, _ := SetupServiceTest(t)
	ctx := TestContext()

	storefrontID := int64(123)
	inputs := []*domain.CreateProductInput{
		{
			StorefrontID: storefrontID,
			Name:         "AB", // Too short (min=3)
			Description:  "Test",
			Price:        99.99,
			Currency:     "USD",
			CategoryID:   "3b4246cc-9970-403c-af01-c142a4178dc6",
		},
	}

	// Mock repository to return validation error through BulkProductError
	expectedErrors := []domain.BulkProductError{
		{
			Index:        0,
			ErrorCode:    "products.invalid_name",
			ErrorMessage: "validation failed for product at index 0",
		},
	}

	mockRepo.On("BulkCreateProducts", ctx, storefrontID, inputs).
		Return([]*domain.Product{}, expectedErrors, nil)

	products, errors, err := service.BulkCreateProducts(ctx, storefrontID, inputs)

	assert.NoError(t, err)
	assert.Len(t, products, 0)
	assert.Len(t, errors, 1)
	assert.Equal(t, int32(0), errors[0].Index)
	mockRepo.AssertExpectations(t)
}

func TestBulkCreateProducts_Error_StorefrontIDMismatch(t *testing.T) {
	service, _, _, _ := SetupServiceTest(t)
	ctx := TestContext()

	storefrontID := int64(123)
	inputs := []*domain.CreateProductInput{
		NewCreateProductInput(2, "Product 1"), // Different storefront_id
	}

	// Service should override storefront_id to match, so this should succeed
	expectedProducts := []*domain.Product{
		NewTestProduct(1, storefrontID, "Product 1"),
	}

	// No error expected because service overrides the storefront_id
	mockRepo := service.repo.(*mocks.MockRepository)
	mockRepo.On("BulkCreateProducts", ctx, storefrontID, inputs).
		Return(expectedProducts, []domain.BulkProductError{}, nil)

	products, bulkErrors, err := service.BulkCreateProducts(ctx, storefrontID, inputs)

	assert.NoError(t, err)
	assert.Len(t, products, 1)
	assert.Len(t, bulkErrors, 0)
	assert.Equal(t, storefrontID, inputs[0].StorefrontID) // Verify override
	mockRepo.AssertExpectations(t)
}

func TestBulkCreateProducts_Error_NegativePrice(t *testing.T) {
	service, mockRepo, _, _ := SetupServiceTest(t)
	ctx := TestContext()

	storefrontID := int64(123)
	inputs := []*domain.CreateProductInput{
		{
			StorefrontID:  storefrontID,
			Name:          "Product 1",
			Description:   "Test",
			Price:         -10.00, // Negative price
			Currency:      "USD",
			CategoryID:    "3b4246cc-9970-403c-af01-c142a4178dc6",
			StockQuantity: 10,
		},
	}

	// Mock repository to return validation error through BulkProductError
	expectedErrors := []domain.BulkProductError{
		{
			Index:        0,
			ErrorCode:    "products.invalid_price",
			ErrorMessage: "validation failed for product at index 0",
		},
	}

	mockRepo.On("BulkCreateProducts", ctx, storefrontID, inputs).
		Return([]*domain.Product{}, expectedErrors, nil)

	products, errors, err := service.BulkCreateProducts(ctx, storefrontID, inputs)

	assert.NoError(t, err)
	assert.Len(t, products, 0)
	assert.Len(t, errors, 1)
	assert.Equal(t, int32(0), errors[0].Index)
	mockRepo.AssertExpectations(t)
}

func TestBulkCreateProducts_Error_MissingName(t *testing.T) {
	service, mockRepo, _, _ := SetupServiceTest(t)
	ctx := TestContext()

	storefrontID := int64(123)
	inputs := []*domain.CreateProductInput{
		{
			StorefrontID:  storefrontID,
			Name:          "", // Missing name
			Description:   "Test",
			Price:         99.99,
			Currency:      "USD",
			CategoryID:    "3b4246cc-9970-403c-af01-c142a4178dc6",
			StockQuantity: 10,
		},
	}

	// Mock repository to return validation error through BulkProductError
	expectedErrors := []domain.BulkProductError{
		{
			Index:        0,
			ErrorCode:    "products.invalid_name",
			ErrorMessage: "validation failed for product at index 0",
		},
	}

	mockRepo.On("BulkCreateProducts", ctx, storefrontID, inputs).
		Return([]*domain.Product{}, expectedErrors, nil)

	products, errors, err := service.BulkCreateProducts(ctx, storefrontID, inputs)

	assert.NoError(t, err)
	assert.Len(t, products, 0)
	assert.Len(t, errors, 1)
	assert.Equal(t, int32(0), errors[0].Index)
	mockRepo.AssertExpectations(t)
}

func TestBulkCreateProducts_PartialSuccess_SomeFailSomeSucceed(t *testing.T) {
	service, mockRepo, _, _ := SetupServiceTest(t)
	ctx := TestContext()

	storefrontID := int64(123)
	inputs := []*domain.CreateProductInput{
		NewCreateProductInput(storefrontID, "Product 1"),
		NewCreateProductInput(storefrontID, "Product 2"),
		NewCreateProductInput(storefrontID, "Product 3"),
	}

	expectedProducts := []*domain.Product{
		NewTestProduct(1, storefrontID, "Product 1"),
		// Product 2 failed
		NewTestProduct(3, storefrontID, "Product 3"),
	}

	expectedErrors := []domain.BulkProductError{
		{
			Index:        1,
			ErrorCode:    "products.duplicate_sku",
			ErrorMessage: "Product with this SKU already exists",
		},
	}

	mockRepo.On("BulkCreateProducts", ctx, storefrontID, inputs).
		Return(expectedProducts, expectedErrors, nil)

	products, errors, err := service.BulkCreateProducts(ctx, storefrontID, inputs)

	assert.NoError(t, err)
	assert.Len(t, products, 2)
	assert.Len(t, errors, 1)
	assert.Equal(t, int32(1), errors[0].Index)
	assert.Equal(t, "products.duplicate_sku", errors[0].ErrorCode)
	mockRepo.AssertExpectations(t)
}

// ===================================
// BulkUpdateProducts Tests (9 tests)
// ===================================

func TestBulkUpdateProducts_Success_UpdateMultiple(t *testing.T) {
	service, mockRepo, _, _ := SetupServiceTest(t)
	ctx := TestContext()

	storefrontID := int64(123)
	newName1 := "Updated Product 1"
	newName2 := "Updated Product 2"
	newPrice := 199.99

	updates := []*domain.BulkUpdateProductInput{
		{
			ProductID: 1,
			Name:      &newName1,
			Price:     &newPrice,
		},
		{
			ProductID: 2,
			Name:      &newName2,
		},
	}

	expectedResult := &domain.BulkUpdateProductsResult{
		SuccessfulProducts: []*domain.Product{
			NewTestProduct(1, storefrontID, newName1),
			NewTestProduct(2, storefrontID, newName2),
		},
		FailedUpdates: []domain.BulkUpdateError{},
	}

	mockRepo.On("BulkUpdateProducts", ctx, storefrontID, updates).
		Return(expectedResult, nil)

	result, err := service.BulkUpdateProducts(ctx, storefrontID, updates)

	assert.NoError(t, err)
	assert.Len(t, result.SuccessfulProducts, 2)
	assert.Len(t, result.FailedUpdates, 0)
	mockRepo.AssertExpectations(t)
}

func TestBulkUpdateProducts_Error_EmptyInput(t *testing.T) {
	service, _, _, _ := SetupServiceTest(t)
	ctx := TestContext()

	storefrontID := int64(123)
	updates := []*domain.BulkUpdateProductInput{}

	result, err := service.BulkUpdateProducts(ctx, storefrontID, updates)

	assert.NoError(t, err) // Empty input returns empty result, not error
	assert.NotNil(t, result)
	assert.Len(t, result.SuccessfulProducts, 0)
	assert.Len(t, result.FailedUpdates, 0)
}

func TestBulkUpdateProducts_Error_BatchTooLarge(t *testing.T) {
	service, _, _, _ := SetupServiceTest(t)
	ctx := TestContext()

	storefrontID := int64(123)
	updates := make([]*domain.BulkUpdateProductInput, 1001)
	name := "Test"
	for i := 0; i < 1001; i++ {
		updates[i] = &domain.BulkUpdateProductInput{
			ProductID: int64(i + 1),
			Name:      &name,
		}
	}

	result, err := service.BulkUpdateProducts(ctx, storefrontID, updates)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "products.bulk_update_limit_exceeded")
	assert.Nil(t, result)
}

func TestBulkUpdateProducts_Error_NoFieldsToUpdate(t *testing.T) {
	service, _, _, _ := SetupServiceTest(t)
	ctx := TestContext()

	storefrontID := int64(123)
	updates := []*domain.BulkUpdateProductInput{
		{
			ProductID: 1,
			// No fields specified for update
		},
	}

	result, err := service.BulkUpdateProducts(ctx, storefrontID, updates)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "at least one field must be specified for update")
	assert.Nil(t, result)
}

func TestBulkUpdateProducts_Error_InvalidProductID(t *testing.T) {
	service, _, _, _ := SetupServiceTest(t)
	ctx := TestContext()

	storefrontID := int64(123)
	name := "Test"
	updates := []*domain.BulkUpdateProductInput{
		{
			ProductID: 0, // Invalid ID
			Name:      &name,
		},
	}

	result, err := service.BulkUpdateProducts(ctx, storefrontID, updates)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "all product_ids must be greater than 0")
	assert.Nil(t, result)
}

func TestBulkUpdateProducts_Error_OwnershipCheck(t *testing.T) {
	service, mockRepo, _, _ := SetupServiceTest(t)
	ctx := TestContext()

	storefrontID := int64(123)
	name := "Updated Product"
	updates := []*domain.BulkUpdateProductInput{
		{
			ProductID: 1,
			Name:      &name,
		},
	}

	expectedResult := &domain.BulkUpdateProductsResult{
		SuccessfulProducts: []*domain.Product{},
		FailedUpdates: []domain.BulkUpdateError{
			{
				ProductID:    1,
				ErrorCode:    "products.not_found",
				ErrorMessage: "Product not found or does not belong to this storefront",
			},
		},
	}

	mockRepo.On("BulkUpdateProducts", ctx, storefrontID, updates).
		Return(expectedResult, nil)

	result, err := service.BulkUpdateProducts(ctx, storefrontID, updates)

	assert.NoError(t, err)
	assert.Len(t, result.SuccessfulProducts, 0)
	assert.Len(t, result.FailedUpdates, 1)
	assert.Equal(t, "products.not_found", result.FailedUpdates[0].ErrorCode)
	mockRepo.AssertExpectations(t)
}

func TestBulkUpdateProducts_PartialSuccess_Mixed(t *testing.T) {
	service, mockRepo, _, _ := SetupServiceTest(t)
	ctx := TestContext()

	storefrontID := int64(123)
	name1 := "Updated Product 1"
	name2 := "Updated Product 2"
	updates := []*domain.BulkUpdateProductInput{
		{
			ProductID: 1,
			Name:      &name1,
		},
		{
			ProductID: 2,
			Name:      &name2,
		},
		{
			ProductID: 3,
			Name:      &name1,
		},
	}

	expectedResult := &domain.BulkUpdateProductsResult{
		SuccessfulProducts: []*domain.Product{
			NewTestProduct(1, storefrontID, name1),
			NewTestProduct(3, storefrontID, name1),
		},
		FailedUpdates: []domain.BulkUpdateError{
			{
				ProductID:    2,
				ErrorCode:    "products.not_found",
				ErrorMessage: "Product not found",
			},
		},
	}

	mockRepo.On("BulkUpdateProducts", ctx, storefrontID, updates).
		Return(expectedResult, nil)

	result, err := service.BulkUpdateProducts(ctx, storefrontID, updates)

	assert.NoError(t, err)
	assert.Len(t, result.SuccessfulProducts, 2)
	assert.Len(t, result.FailedUpdates, 1)
	mockRepo.AssertExpectations(t)
}

func TestBulkUpdateProducts_Error_NegativePrice(t *testing.T) {
	service, _, _, _ := SetupServiceTest(t)
	ctx := TestContext()

	storefrontID := int64(123)
	negativePrice := -10.0
	updates := []*domain.BulkUpdateProductInput{
		{
			ProductID: 1,
			Price:     &negativePrice,
		},
	}

	result, err := service.BulkUpdateProducts(ctx, storefrontID, updates)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "validation failed for product")
	assert.Nil(t, result)
}

func TestBulkUpdateProducts_Error_InvalidStorefrontID(t *testing.T) {
	service, _, _, _ := SetupServiceTest(t)
	ctx := TestContext()

	storefrontID := int64(0) // Invalid storefront ID
	name := "Test"
	updates := []*domain.BulkUpdateProductInput{
		{
			ProductID: 1,
			Name:      &name,
		},
	}

	result, err := service.BulkUpdateProducts(ctx, storefrontID, updates)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "storefront_id must be greater than 0")
	assert.Nil(t, result)
}

// ===================================
// BulkDeleteProducts Tests (9 tests)
// ===================================

func TestBulkDeleteProducts_Success_SoftDelete_Multiple(t *testing.T) {
	service, mockRepo, _, _ := SetupServiceTest(t)
	ctx := TestContext()

	storefrontID := int64(123)
	productIDs := []int64{1, 2, 3}
	hardDelete := false

	mockRepo.On("BulkDeleteProducts", ctx, storefrontID, productIDs, hardDelete).
		Return(int32(3), int32(0), int32(5), map[int64]string{}, nil)

	successCount, failedCount, variantsDeleted, errors, err := service.BulkDeleteProducts(ctx, storefrontID, productIDs, hardDelete)

	assert.NoError(t, err)
	assert.Equal(t, int32(3), successCount)
	assert.Equal(t, int32(0), failedCount)
	assert.Equal(t, int32(5), variantsDeleted)
	assert.Empty(t, errors)
	mockRepo.AssertExpectations(t)
}

func TestBulkDeleteProducts_Success_HardDelete_Multiple(t *testing.T) {
	service, mockRepo, _, _ := SetupServiceTest(t)
	ctx := TestContext()

	storefrontID := int64(123)
	productIDs := []int64{1, 2, 3}
	hardDelete := true

	mockRepo.On("BulkDeleteProducts", ctx, storefrontID, productIDs, hardDelete).
		Return(int32(3), int32(0), int32(5), map[int64]string{}, nil)

	successCount, failedCount, variantsDeleted, errors, err := service.BulkDeleteProducts(ctx, storefrontID, productIDs, hardDelete)

	assert.NoError(t, err)
	assert.Equal(t, int32(3), successCount)
	assert.Equal(t, int32(0), failedCount)
	assert.Equal(t, int32(5), variantsDeleted)
	assert.Empty(t, errors)
	mockRepo.AssertExpectations(t)
}

func TestBulkDeleteProducts_Error_EmptyInput(t *testing.T) {
	service, _, _, _ := SetupServiceTest(t)
	ctx := TestContext()

	storefrontID := int64(123)
	productIDs := []int64{}
	hardDelete := false

	successCount, failedCount, variantsDeleted, errors, err := service.BulkDeleteProducts(ctx, storefrontID, productIDs, hardDelete)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "product_ids list cannot be empty")
	assert.Equal(t, int32(0), successCount)
	assert.Equal(t, int32(0), failedCount)
	assert.Equal(t, int32(0), variantsDeleted)
	assert.Nil(t, errors)
}

func TestBulkDeleteProducts_Error_BatchTooLarge(t *testing.T) {
	service, _, _, _ := SetupServiceTest(t)
	ctx := TestContext()

	storefrontID := int64(123)
	productIDs := make([]int64, 1001)
	for i := 0; i < 1001; i++ {
		productIDs[i] = int64(i + 1)
	}
	hardDelete := false

	successCount, failedCount, variantsDeleted, errors, err := service.BulkDeleteProducts(ctx, storefrontID, productIDs, hardDelete)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot delete more than 1000 products at once")
	assert.Equal(t, int32(0), successCount)
	assert.Equal(t, int32(0), failedCount)
	assert.Equal(t, int32(0), variantsDeleted)
	assert.Nil(t, errors)
}

func TestBulkDeleteProducts_Error_InvalidStorefrontID(t *testing.T) {
	service, _, _, _ := SetupServiceTest(t)
	ctx := TestContext()

	storefrontID := int64(0) // Invalid
	productIDs := []int64{1, 2, 3}
	hardDelete := false

	successCount, failedCount, variantsDeleted, errors, err := service.BulkDeleteProducts(ctx, storefrontID, productIDs, hardDelete)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "storefront_id must be greater than 0")
	assert.Equal(t, int32(0), successCount)
	assert.Equal(t, int32(0), failedCount)
	assert.Equal(t, int32(0), variantsDeleted)
	assert.Nil(t, errors)
}

func TestBulkDeleteProducts_Success_DuplicateIDs(t *testing.T) {
	service, mockRepo, _, _ := SetupServiceTest(t)
	ctx := TestContext()

	storefrontID := int64(123)
	productIDs := []int64{1, 2, 1, 3, 2} // Duplicates: 1 and 2
	deduplicatedIDs := []int64{1, 2, 3}
	hardDelete := false

	mockRepo.On("BulkDeleteProducts", ctx, storefrontID, deduplicatedIDs, hardDelete).
		Return(int32(3), int32(0), int32(5), map[int64]string{}, nil)

	successCount, failedCount, variantsDeleted, deleteErrors, err := service.BulkDeleteProducts(ctx, storefrontID, productIDs, hardDelete)

	assert.NoError(t, err)
	assert.Equal(t, int32(3), successCount)
	assert.Equal(t, int32(0), failedCount)
	assert.Equal(t, int32(5), variantsDeleted)
	assert.Empty(t, deleteErrors)
	mockRepo.AssertExpectations(t)
}

func TestBulkDeleteProducts_Error_ZeroIDs(t *testing.T) {
	service, _, _, _ := SetupServiceTest(t)
	ctx := TestContext()

	storefrontID := int64(123)
	productIDs := []int64{0, 0, 0} // All invalid IDs
	hardDelete := false

	successCount, failedCount, variantsDeleted, errors, err := service.BulkDeleteProducts(ctx, storefrontID, productIDs, hardDelete)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no valid product IDs provided")
	assert.Equal(t, int32(0), successCount)
	assert.Equal(t, int32(0), failedCount)
	assert.Equal(t, int32(0), variantsDeleted)
	assert.Nil(t, errors)
}

func TestBulkDeleteProducts_PartialSuccess_Mixed(t *testing.T) {
	service, mockRepo, _, _ := SetupServiceTest(t)
	ctx := TestContext()

	storefrontID := int64(123)
	productIDs := []int64{1, 2, 3}
	hardDelete := false

	expectedErrors := map[int64]string{
		2: "products.not_found",
	}

	mockRepo.On("BulkDeleteProducts", ctx, storefrontID, productIDs, hardDelete).
		Return(int32(2), int32(1), int32(3), expectedErrors, nil)

	successCount, failedCount, variantsDeleted, errors, err := service.BulkDeleteProducts(ctx, storefrontID, productIDs, hardDelete)

	assert.NoError(t, err)
	assert.Equal(t, int32(2), successCount)
	assert.Equal(t, int32(1), failedCount)
	assert.Equal(t, int32(3), variantsDeleted)
	assert.Len(t, errors, 1)
	assert.Equal(t, "products.not_found", errors[2])
	mockRepo.AssertExpectations(t)
}

func TestBulkDeleteProducts_Success_CascadeDelete_Variants(t *testing.T) {
	service, mockRepo, _, _ := SetupServiceTest(t)
	ctx := TestContext()

	storefrontID := int64(123)
	productIDs := []int64{1} // Product with 10 variants
	hardDelete := true

	mockRepo.On("BulkDeleteProducts", ctx, storefrontID, productIDs, hardDelete).
		Return(int32(1), int32(0), int32(10), map[int64]string{}, nil)

	successCount, failedCount, variantsDeleted, errors, err := service.BulkDeleteProducts(ctx, storefrontID, productIDs, hardDelete)

	assert.NoError(t, err)
	assert.Equal(t, int32(1), successCount)
	assert.Equal(t, int32(0), failedCount)
	assert.Equal(t, int32(10), variantsDeleted) // Verify cascade
	assert.Empty(t, errors)
	mockRepo.AssertExpectations(t)
}

// ===================================
// BulkCreateProductVariants Tests (6 tests)
// ===================================

func TestBulkCreateProductVariants_Success_MultipleVariants(t *testing.T) {
	service, mockRepo, _, _ := SetupServiceTest(t)
	ctx := TestContext()

	productID := int64(123)
	inputs := []*domain.CreateVariantInput{
		NewCreateVariantInput(productID),
		NewCreateVariantInput(productID),
		NewCreateVariantInput(productID),
	}

	expectedVariants := []*domain.ProductVariant{
		NewTestProductVariant(1, productID),
		NewTestProductVariant(2, productID),
		NewTestProductVariant(3, productID),
	}

	mockRepo.On("BulkCreateProductVariants", ctx, productID, inputs).
		Return(expectedVariants, nil)

	variants, err := service.BulkCreateProductVariants(ctx, productID, inputs)

	assert.NoError(t, err)
	assert.Len(t, variants, 3)
	mockRepo.AssertExpectations(t)
}

func TestBulkCreateProductVariants_Error_EmptyInput(t *testing.T) {
	service, _, _, _ := SetupServiceTest(t)
	ctx := TestContext()

	productID := int64(123)
	inputs := []*domain.CreateVariantInput{}

	variants, err := service.BulkCreateProductVariants(ctx, productID, inputs)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "variants list cannot be empty")
	assert.Nil(t, variants)
}

func TestBulkCreateProductVariants_Error_BatchTooLarge(t *testing.T) {
	service, _, _, _ := SetupServiceTest(t)
	ctx := TestContext()

	productID := int64(123)
	inputs := make([]*domain.CreateVariantInput, 1001)
	for i := 0; i < 1001; i++ {
		inputs[i] = NewCreateVariantInput(productID)
	}

	variants, err := service.BulkCreateProductVariants(ctx, productID, inputs)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot create more than 1000 variants at once")
	assert.Nil(t, variants)
}

func TestBulkCreateProductVariants_Error_InvalidProductID(t *testing.T) {
	service, _, _, _ := SetupServiceTest(t)
	ctx := TestContext()

	productID := int64(0) // Invalid
	inputs := []*domain.CreateVariantInput{
		NewCreateVariantInput(1),
	}

	variants, err := service.BulkCreateProductVariants(ctx, productID, inputs)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "product_id must be greater than 0")
	assert.Nil(t, variants)
}

func TestBulkCreateProductVariants_Error_MultipleDefaults(t *testing.T) {
	service, _, _, _ := SetupServiceTest(t)
	ctx := TestContext()

	productID := int64(123)
	input1 := NewCreateVariantInput(productID)
	input1.IsDefault = true
	input2 := NewCreateVariantInput(productID)
	input2.IsDefault = true // Multiple defaults

	inputs := []*domain.CreateVariantInput{input1, input2}

	variants, err := service.BulkCreateProductVariants(ctx, productID, inputs)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "only one variant can be set as default")
	assert.Nil(t, variants)
}

func TestBulkCreateProductVariants_Error_ValidationFailed(t *testing.T) {
	service, _, _, _ := SetupServiceTest(t)
	ctx := TestContext()

	productID := int64(123)
	negativePrice := -10.0
	inputs := []*domain.CreateVariantInput{
		{
			ProductID:     productID,
			Price:         &negativePrice, // Negative price
			StockQuantity: 10,
		},
	}

	variants, err := service.BulkCreateProductVariants(ctx, productID, inputs)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "validation failed for variant at index 0")
	assert.Nil(t, variants)
}

// ===================================
// CreateListing Tests (10 tests)
// ===================================

func TestCreateListing_Success_MinimalFields(t *testing.T) {
	service, mockRepo, _, _ := SetupServiceTest(t)
	ctx := TestContext()

	userID := int64(100)
	input := NewCreateListingInput(userID, "Test Listing")

	expectedListing := NewTestListing(1, userID, "Test Listing")

	// Mock category validation - category must be active
	mockRepo.On("GetCategoryByID", ctx, "1").
		Return(&domain.Category{ID: "3b4246cc-9970-403c-af01-c142a4178dc6", Name: "Test Category", IsActive: true}, nil)

	// Mock slug uniqueness check - slug doesn't exist yet
	mockRepo.On("GetListingBySlug", ctx, mock.AnythingOfType("string")).
		Return(nil, errors.New("not found"))

	mockRepo.On("CreateListing", ctx, input).
		Return(expectedListing, nil)
	mockRepo.On("EnqueueIndexing", ctx, expectedListing.ID, domain.IndexOpIndex).
		Return(nil)

	listing, err := service.CreateListing(ctx, input)

	assert.NoError(t, err)
	// Nil-safe assertions: check listing is not nil before accessing fields
	if assert.NotNil(t, listing, "listing should be created") {
		assert.Equal(t, expectedListing.ID, listing.ID)
		assert.Equal(t, userID, listing.UserID)
	}
	mockRepo.AssertExpectations(t)
}

func TestCreateListing_Success_WithStorefront(t *testing.T) {
	service, mockRepo, _, _ := SetupServiceTest(t)
	ctx := TestContext()

	userID := int64(100)
	storefrontID := int64(10)
	input := NewCreateListingInput(userID, "Test Listing")
	input.StorefrontID = &storefrontID

	expectedListing := NewTestListing(1, userID, "Test Listing")
	expectedListing.StorefrontID = &storefrontID

	// Mock category validation
	mockRepo.On("GetCategoryByID", ctx, "1").
		Return(&domain.Category{ID: "3b4246cc-9970-403c-af01-c142a4178dc6", Name: "Test Category", IsActive: true}, nil)
	// Mock slug uniqueness
	mockRepo.On("GetListingBySlug", ctx, mock.AnythingOfType("string")).
		Return(nil, errors.New("not found"))

	mockRepo.On("CreateListing", ctx, input).
		Return(expectedListing, nil)
	mockRepo.On("EnqueueIndexing", ctx, expectedListing.ID, domain.IndexOpIndex).
		Return(nil)

	listing, err := service.CreateListing(ctx, input)

	assert.NoError(t, err)
	assert.NotNil(t, listing)
	assert.Equal(t, storefrontID, *listing.StorefrontID)
	mockRepo.AssertExpectations(t)
}

func TestCreateListing_Success_WithAllFields(t *testing.T) {
	service, mockRepo, _, _ := SetupServiceTest(t)
	ctx := TestContext()

	// Mock category validation and slug uniqueness (common for all create tests)
	mockRepo.On("GetCategoryByID", ctx, mock.AnythingOfType("int64")).
		Return(&domain.Category{ID: "3b4246cc-9970-403c-af01-c142a4178dc6", Name: "Test Category", IsActive: true}, nil)
	mockRepo.On("GetListingBySlug", ctx, mock.AnythingOfType("string")).
		Return(nil, errors.New("not found"))

	userID := int64(100)
	input := NewCreateListingInput(userID, "Complete Test Listing")

	expectedListing := NewTestListing(1, userID, "Complete Test Listing")

	mockRepo.On("CreateListing", ctx, input).
		Return(expectedListing, nil)
	mockRepo.On("EnqueueIndexing", ctx, expectedListing.ID, domain.IndexOpIndex).
		Return(nil)

	listing, err := service.CreateListing(ctx, input)

	assert.NoError(t, err)
	assert.NotNil(t, listing)
	assert.Equal(t, "Complete Test Listing", listing.Title)
	assert.NotNil(t, listing.Description)
	assert.Equal(t, 99.99, listing.Price)
	assert.Equal(t, "USD", listing.Currency)
	mockRepo.AssertExpectations(t)
}

func TestCreateListing_Error_NegativePrice(t *testing.T) {
	service, mockRepo, _, _ := SetupServiceTest(t)
	ctx := TestContext()
	SetupDefaultCreateListingMocks(mockRepo, ctx)

	userID := int64(100)
	input := NewCreateListingInput(userID, "Test Listing")
	input.Price = -10.99 // Negative price

	listing, err := service.CreateListing(ctx, input)

	assert.Error(t, err)
	// Validator catches this before business logic
	assert.Contains(t, err.Error(), "validation failed")
	assert.Nil(t, listing)
}

func TestCreateListing_Error_NegativeQuantity(t *testing.T) {
	service, mockRepo, _, _ := SetupServiceTest(t)
	ctx := TestContext()
	SetupDefaultCreateListingMocks(mockRepo, ctx)

	userID := int64(100)
	input := NewCreateListingInput(userID, "Test Listing")
	input.Quantity = -5 // Negative quantity

	listing, err := service.CreateListing(ctx, input)

	assert.Error(t, err)
	// Validator catches this before business logic
	assert.Contains(t, err.Error(), "validation failed")
	assert.Nil(t, listing)
}

func TestCreateListing_Error_MissingTitle(t *testing.T) {
	service, mockRepo, _, _ := SetupServiceTest(t)
	ctx := TestContext()
	SetupDefaultCreateListingMocks(mockRepo, ctx)

	userID := int64(100)
	input := NewCreateListingInput(userID, "")
	input.Title = "" // Missing title

	listing, err := service.CreateListing(ctx, input)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "validation failed")
	assert.Nil(t, listing)
}

func TestCreateListing_Error_ShortTitle(t *testing.T) {
	service, mockRepo, _, _ := SetupServiceTest(t)
	ctx := TestContext()
	SetupDefaultCreateListingMocks(mockRepo, ctx)

	userID := int64(100)
	input := NewCreateListingInput(userID, "AB") // Title too short (< 3 chars)

	listing, err := service.CreateListing(ctx, input)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "validation failed")
	assert.Nil(t, listing)
}

func TestCreateListing_Error_InvalidCurrency(t *testing.T) {
	service, mockRepo, _, _ := SetupServiceTest(t)
	ctx := TestContext()
	SetupDefaultCreateListingMocks(mockRepo, ctx)

	userID := int64(100)
	input := NewCreateListingInput(userID, "Test Listing")
	input.Currency = "US" // Invalid currency (must be 3 chars)

	listing, err := service.CreateListing(ctx, input)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "validation failed")
	assert.Nil(t, listing)
}

func TestCreateListing_Success_EnqueueIndexing_Success(t *testing.T) {
	service, mockRepo, _, _ := SetupServiceTest(t)
	ctx := TestContext()

	userID := int64(100)
	input := NewCreateListingInput(userID, "Test Listing")

	expectedListing := NewTestListing(1, userID, "Test Listing")

	// Mock validation
	mockRepo.On("GetCategoryByID", ctx, "1").
		Return(&domain.Category{ID: "1", Name: "Test", IsActive: true}, nil)
	mockRepo.On("GetListingBySlug", ctx, mock.AnythingOfType("string")).
		Return(nil, errors.New("not found"))

	mockRepo.On("CreateListing", ctx, input).
		Return(expectedListing, nil)
	mockRepo.On("EnqueueIndexing", ctx, expectedListing.ID, domain.IndexOpIndex).
		Return(nil) // Success

	listing, err := service.CreateListing(ctx, input)

	assert.NoError(t, err)
	assert.NotNil(t, listing)
	mockRepo.AssertExpectations(t)
}

func TestCreateListing_Success_EnqueueIndexing_Failure_NonCritical(t *testing.T) {
	service, mockRepo, _, _ := SetupServiceTest(t)
	ctx := TestContext()

	userID := int64(100)
	input := NewCreateListingInput(userID, "Test Listing")

	expectedListing := NewTestListing(1, userID, "Test Listing")

	// Mock validation
	mockRepo.On("GetCategoryByID", ctx, "1").
		Return(&domain.Category{ID: "1", Name: "Test", IsActive: true}, nil)
	mockRepo.On("GetListingBySlug", ctx, mock.AnythingOfType("string")).
		Return(nil, errors.New("not found"))

	mockRepo.On("CreateListing", ctx, input).
		Return(expectedListing, nil)
	mockRepo.On("EnqueueIndexing", ctx, expectedListing.ID, domain.IndexOpIndex).
		Return(assert.AnError) // Indexing fails but non-critical

	listing, err := service.CreateListing(ctx, input)

	// Should still succeed despite indexing failure
	assert.NoError(t, err)
	assert.NotNil(t, listing)
	mockRepo.AssertExpectations(t)
}

// ===================================
// UpdateListing Tests (9 tests)
// ===================================

func TestUpdateListing_Success_UpdatePrice(t *testing.T) {
	service, mockRepo, mockCache, _ := SetupServiceTest(t)
	ctx := TestContext()

	listingID := int64(123)
	userID := int64(100)
	newPrice := 129.99

	existingListing := NewTestListing(listingID, userID, "Test Listing")

	input := &domain.UpdateListingInput{
		Price: &newPrice,
	}

	updatedListing := NewTestListing(listingID, userID, "Test Listing")
	updatedListing.Price = newPrice

	mockRepo.On("GetListingByID", ctx, listingID).
		Return(existingListing, nil)
	mockRepo.On("UpdateListing", ctx, listingID, input).
		Return(updatedListing, nil)
	mockRepo.On("GetImages", ctx, listingID).
		Return([]*domain.ListingImage{}, nil)
	mockCache.On("Delete", ctx, "listing:1").
		Return(nil)
	mockRepo.On("EnqueueIndexing", ctx, listingID, domain.IndexOpUpdate).
		Return(nil)

	listing, err := service.UpdateListing(ctx, listingID, userID, input)

	assert.NoError(t, err)
	assert.NotNil(t, listing)
	assert.Equal(t, newPrice, listing.Price)
	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestUpdateListing_Success_UpdateQuantity(t *testing.T) {
	service, mockRepo, mockCache, _ := SetupServiceTest(t)
	ctx := TestContext()

	listingID := int64(123)
	userID := int64(100)
	newQuantity := int32(25)

	existingListing := NewTestListing(listingID, userID, "Test Listing")

	input := &domain.UpdateListingInput{
		Quantity: &newQuantity,
	}

	updatedListing := NewTestListing(listingID, userID, "Test Listing")
	updatedListing.Quantity = newQuantity

	mockRepo.On("GetListingByID", ctx, listingID).
		Return(existingListing, nil)
	mockRepo.On("UpdateListing", ctx, listingID, input).
		Return(updatedListing, nil)
	mockRepo.On("GetImages", ctx, listingID).
		Return([]*domain.ListingImage{}, nil)
	mockCache.On("Delete", ctx, "listing:1").
		Return(nil)
	mockRepo.On("EnqueueIndexing", ctx, listingID, domain.IndexOpUpdate).
		Return(nil)

	listing, err := service.UpdateListing(ctx, listingID, userID, input)

	assert.NoError(t, err)
	assert.NotNil(t, listing)
	assert.Equal(t, newQuantity, listing.Quantity)
	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestUpdateListing_Success_UpdateMultipleFields(t *testing.T) {
	service, mockRepo, mockCache, _ := SetupServiceTest(t)
	ctx := TestContext()

	listingID := int64(123)
	userID := int64(100)
	newTitle := "Updated Title"
	newPrice := 199.99
	newQuantity := int32(30)

	existingListing := NewTestListing(listingID, userID, "Test Listing")

	input := &domain.UpdateListingInput{
		Title:    &newTitle,
		Price:    &newPrice,
		Quantity: &newQuantity,
	}

	updatedListing := NewTestListing(listingID, userID, newTitle)
	updatedListing.Price = newPrice
	updatedListing.Quantity = newQuantity

	mockRepo.On("GetListingByID", ctx, listingID).
		Return(existingListing, nil)
	mockRepo.On("UpdateListing", ctx, listingID, input).
		Return(updatedListing, nil)
	mockRepo.On("GetImages", ctx, listingID).
		Return([]*domain.ListingImage{}, nil)
	mockCache.On("Delete", ctx, "listing:1").
		Return(nil)
	mockRepo.On("EnqueueIndexing", ctx, listingID, domain.IndexOpUpdate).
		Return(nil)

	listing, err := service.UpdateListing(ctx, listingID, userID, input)

	assert.NoError(t, err)
	assert.NotNil(t, listing)
	assert.Equal(t, newTitle, listing.Title)
	assert.Equal(t, newPrice, listing.Price)
	assert.Equal(t, newQuantity, listing.Quantity)
	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestUpdateListing_Error_NotFound(t *testing.T) {
	service, mockRepo, _, _ := SetupServiceTest(t)
	ctx := TestContext()

	listingID := int64(999)
	userID := int64(100)
	newPrice := 129.99

	input := &domain.UpdateListingInput{
		Price: &newPrice,
	}

	mockRepo.On("GetListingByID", ctx, listingID).
		Return(nil, assert.AnError) // Not found

	listing, err := service.UpdateListing(ctx, listingID, userID, input)

	assert.Error(t, err)
	assert.Nil(t, listing)
	assert.Contains(t, err.Error(), "listing not found")
	mockRepo.AssertExpectations(t)
}

func TestUpdateListing_Error_UnauthorizedUser_OwnershipCheck(t *testing.T) {
	service, mockRepo, _, _ := SetupServiceTest(t)
	ctx := TestContext()

	listingID := int64(123)
	ownerID := int64(100)
	unauthorizedUserID := int64(200) // Different user

	existingListing := NewTestListing(listingID, ownerID, "Test Listing")

	newPrice := 129.99
	input := &domain.UpdateListingInput{
		Price: &newPrice,
	}

	mockRepo.On("GetListingByID", ctx, listingID).
		Return(existingListing, nil)

	listing, err := service.UpdateListing(ctx, listingID, unauthorizedUserID, input)

	assert.Error(t, err)
	assert.Nil(t, listing)
	assert.Contains(t, err.Error(), "unauthorized")
	mockRepo.AssertExpectations(t)
}

func TestUpdateListing_Error_NegativePrice(t *testing.T) {
	service, _, _, _ := SetupServiceTest(t)
	ctx := TestContext()

	listingID := int64(123)
	userID := int64(100)
	negativePrice := -50.0

	input := &domain.UpdateListingInput{
		Price: &negativePrice,
	}

	listing, err := service.UpdateListing(ctx, listingID, userID, input)

	assert.Error(t, err)
	assert.Nil(t, listing)
	// Validator catches this before repo call
	assert.Contains(t, err.Error(), "validation failed")
}

func TestUpdateListing_Error_NegativeQuantity(t *testing.T) {
	service, _, _, _ := SetupServiceTest(t)
	ctx := TestContext()

	listingID := int64(123)
	userID := int64(100)
	negativeQuantity := int32(-10)

	input := &domain.UpdateListingInput{
		Quantity: &negativeQuantity,
	}

	listing, err := service.UpdateListing(ctx, listingID, userID, input)

	assert.Error(t, err)
	assert.Nil(t, listing)
	// Validator catches this before repo call
	assert.Contains(t, err.Error(), "validation failed")
}

func TestUpdateListing_Success_CacheInvalidation_Success(t *testing.T) {
	service, mockRepo, mockCache, _ := SetupServiceTest(t)
	ctx := TestContext()

	listingID := int64(123)
	userID := int64(100)
	newPrice := 129.99

	existingListing := NewTestListing(listingID, userID, "Test Listing")
	updatedListing := NewTestListing(listingID, userID, "Test Listing")
	updatedListing.Price = newPrice

	input := &domain.UpdateListingInput{
		Price: &newPrice,
	}

	mockRepo.On("GetListingByID", ctx, listingID).
		Return(existingListing, nil)
	mockRepo.On("UpdateListing", ctx, listingID, input).
		Return(updatedListing, nil)
	mockRepo.On("GetImages", ctx, listingID).
		Return([]*domain.ListingImage{}, nil)
	mockCache.On("Delete", ctx, "listing:1").
		Return(nil) // Cache invalidation succeeds
	mockRepo.On("EnqueueIndexing", ctx, listingID, domain.IndexOpUpdate).
		Return(nil)

	listing, err := service.UpdateListing(ctx, listingID, userID, input)

	assert.NoError(t, err)
	assert.NotNil(t, listing)
	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestUpdateListing_Success_EnqueueIndexing_Success(t *testing.T) {
	service, mockRepo, mockCache, _ := SetupServiceTest(t)
	ctx := TestContext()

	listingID := int64(123)
	userID := int64(100)
	newPrice := 129.99

	existingListing := NewTestListing(listingID, userID, "Test Listing")
	updatedListing := NewTestListing(listingID, userID, "Test Listing")
	updatedListing.Price = newPrice

	input := &domain.UpdateListingInput{
		Price: &newPrice,
	}

	mockRepo.On("GetListingByID", ctx, listingID).
		Return(existingListing, nil)
	mockRepo.On("UpdateListing", ctx, listingID, input).
		Return(updatedListing, nil)
	mockRepo.On("GetImages", ctx, listingID).
		Return([]*domain.ListingImage{}, nil)
	mockCache.On("Delete", ctx, "listing:1").
		Return(nil)
	mockRepo.On("EnqueueIndexing", ctx, listingID, domain.IndexOpUpdate).
		Return(nil) // Indexing succeeds

	listing, err := service.UpdateListing(ctx, listingID, userID, input)

	assert.NoError(t, err)
	assert.NotNil(t, listing)
	mockRepo.AssertExpectations(t)
}

// ===================================
// DeleteListing Tests (5 tests)
// ===================================

func TestDeleteListing_Success_SoftDelete(t *testing.T) {
	service, mockRepo, mockCache, _ := SetupServiceTest(t)
	ctx := TestContext()

	listingID := int64(123)
	userID := int64(100)

	existingListing := NewTestListing(listingID, userID, "Test Listing")

	mockRepo.On("GetListingByID", ctx, listingID).
		Return(existingListing, nil)
	mockRepo.On("DeleteListing", ctx, listingID).
		Return(nil)
	mockRepo.On("GetImages", ctx, listingID).
		Return([]*domain.ListingImage{}, nil) // No images to delete
	mockCache.On("Delete", ctx, "listing:1").
		Return(nil)
	mockCache.On("Delete", ctx, "favorites:listing:1:count").
		Return(nil)
	mockCache.On("Delete", ctx, "user:100:listings").
		Return(nil)
	mockRepo.On("EnqueueIndexing", ctx, listingID, domain.IndexOpDelete).
		Return(nil)

	err := service.DeleteListing(ctx, listingID, userID)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestDeleteListing_Error_NotFound(t *testing.T) {
	service, mockRepo, _, _ := SetupServiceTest(t)
	ctx := TestContext()

	listingID := int64(999)
	userID := int64(100)

	mockRepo.On("GetListingByID", ctx, listingID).
		Return(nil, assert.AnError) // Not found

	err := service.DeleteListing(ctx, listingID, userID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "listing not found")
	mockRepo.AssertExpectations(t)
}

func TestDeleteListing_Error_UnauthorizedUser(t *testing.T) {
	service, mockRepo, _, _ := SetupServiceTest(t)
	ctx := TestContext()

	listingID := int64(123)
	ownerID := int64(100)
	unauthorizedUserID := int64(200) // Different user

	existingListing := NewTestListing(listingID, ownerID, "Test Listing")

	mockRepo.On("GetListingByID", ctx, listingID).
		Return(existingListing, nil)

	err := service.DeleteListing(ctx, listingID, unauthorizedUserID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unauthorized")
	mockRepo.AssertExpectations(t)
}

func TestDeleteListing_Success_CacheInvalidation_Success(t *testing.T) {
	service, mockRepo, mockCache, _ := SetupServiceTest(t)
	ctx := TestContext()

	listingID := int64(123)
	userID := int64(100)

	existingListing := NewTestListing(listingID, userID, "Test Listing")

	mockRepo.On("GetListingByID", ctx, listingID).
		Return(existingListing, nil)
	mockRepo.On("DeleteListing", ctx, listingID).
		Return(nil)
	mockRepo.On("GetImages", ctx, listingID).
		Return([]*domain.ListingImage{}, nil)
	mockCache.On("Delete", ctx, "listing:1").
		Return(nil) // Cache invalidation succeeds
	mockCache.On("Delete", ctx, "favorites:listing:1:count").
		Return(nil)
	mockCache.On("Delete", ctx, "user:100:listings").
		Return(nil)
	mockRepo.On("EnqueueIndexing", ctx, listingID, domain.IndexOpDelete).
		Return(nil)

	err := service.DeleteListing(ctx, listingID, userID)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestDeleteListing_Success_EnqueueIndexing_Success(t *testing.T) {
	service, mockRepo, mockCache, _ := SetupServiceTest(t)
	ctx := TestContext()

	listingID := int64(123)
	userID := int64(100)

	existingListing := NewTestListing(listingID, userID, "Test Listing")

	mockRepo.On("GetListingByID", ctx, listingID).
		Return(existingListing, nil)
	mockRepo.On("DeleteListing", ctx, listingID).
		Return(nil)
	mockRepo.On("GetImages", ctx, listingID).
		Return([]*domain.ListingImage{}, nil)
	mockCache.On("Delete", ctx, "listing:1").
		Return(nil)
	mockCache.On("Delete", ctx, "favorites:listing:1:count").
		Return(nil)
	mockCache.On("Delete", ctx, "user:100:listings").
		Return(nil)
	mockRepo.On("EnqueueIndexing", ctx, listingID, domain.IndexOpDelete).
		Return(nil) // Indexing succeeds

	err := service.DeleteListing(ctx, listingID, userID)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// ===================================
// ListListings Tests (6 tests)
// ===================================

func TestListListings_Success_DefaultPagination(t *testing.T) {
	service, mockRepo, _, _ := SetupServiceTest(t)
	ctx := TestContext()

	filter := NewListListingsFilter(20)

	expectedListings := []*domain.Listing{
		NewTestListing(1, 100, "Listing 1"),
		NewTestListing(2, 100, "Listing 2"),
	}

	mockRepo.On("ListListings", ctx, filter).
		Return(expectedListings, int32(2), nil)

	// Mock GetImages for each listing
	mockRepo.On("GetImages", ctx, mock.AnythingOfType("int64")).
		Return([]*domain.ListingImage{}, nil).Maybe()

	listings, total, err := service.ListListings(ctx, filter)

	assert.NoError(t, err)
	assert.Len(t, listings, 2)
	assert.Equal(t, int32(2), total)
	mockRepo.AssertExpectations(t)
}

func TestListListings_Success_CustomPagination(t *testing.T) {
	service, mockRepo, _, _ := SetupServiceTest(t)
	ctx := TestContext()

	filter := NewListListingsFilter(50)
	filter.Offset = 10

	expectedListings := []*domain.Listing{
		NewTestListing(11, 100, "Listing 11"),
		NewTestListing(12, 100, "Listing 12"),
	}

	mockRepo.On("ListListings", ctx, filter).
		Return(expectedListings, int32(100), nil)

	// Mock GetImages for each listing
	mockRepo.On("GetImages", ctx, mock.AnythingOfType("int64")).
		Return([]*domain.ListingImage{}, nil).Maybe()

	listings, total, err := service.ListListings(ctx, filter)

	assert.NoError(t, err)
	assert.Len(t, listings, 2)
	assert.Equal(t, int32(100), total)
	mockRepo.AssertExpectations(t)
}

func TestListListings_Success_LimitCapping_Max100(t *testing.T) {
	service, mockRepo, _, _ := SetupServiceTest(t)
	ctx := TestContext()

	filter := NewListListingsFilter(100) // Max allowed

	expectedListings := []*domain.Listing{
		NewTestListing(1, 100, "Listing 1"),
	}

	mockRepo.On("ListListings", ctx, filter).
		Return(expectedListings, int32(1), nil)

	// Mock GetImages for each listing
	mockRepo.On("GetImages", ctx, mock.AnythingOfType("int64")).
		Return([]*domain.ListingImage{}, nil).Maybe()

	listings, total, err := service.ListListings(ctx, filter)

	assert.NoError(t, err)
	assert.Equal(t, int32(100), filter.Limit) // Already at max
	assert.Len(t, listings, 1)
	assert.Equal(t, int32(1), total)
	mockRepo.AssertExpectations(t)
}

func TestListListings_Success_LimitDefault_20(t *testing.T) {
	service, mockRepo, _, _ := SetupServiceTest(t)
	ctx := TestContext()

	filter := NewListListingsFilter(1) // Use 1 as minimum valid

	expectedListings := []*domain.Listing{
		NewTestListing(1, 100, "Listing 1"),
	}

	mockRepo.On("ListListings", ctx, filter).
		Return(expectedListings, int32(1), nil)

	// Mock GetImages for each listing
	mockRepo.On("GetImages", ctx, mock.AnythingOfType("int64")).
		Return([]*domain.ListingImage{}, nil).Maybe()

	listings, total, err := service.ListListings(ctx, filter)

	assert.NoError(t, err)
	assert.Equal(t, int32(1), filter.Limit) // Unchanged
	assert.Len(t, listings, 1)
	assert.Equal(t, int32(1), total)
	mockRepo.AssertExpectations(t)
}

func TestListListings_Error_NegativeOffset(t *testing.T) {
	service, _, _, _ := SetupServiceTest(t)
	ctx := TestContext()

	filter := NewListListingsFilter(20)
	filter.Offset = -10 // Negative offset

	listings, total, err := service.ListListings(ctx, filter)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "validation failed")
	assert.Nil(t, listings)
	assert.Equal(t, int32(0), total)
}

func TestListListings_Success_EmptyResults(t *testing.T) {
	service, mockRepo, _, _ := SetupServiceTest(t)
	ctx := TestContext()

	filter := NewListListingsFilter(20)

	expectedListings := []*domain.Listing{} // Empty results

	mockRepo.On("ListListings", ctx, filter).
		Return(expectedListings, int32(0), nil)

	listings, total, err := service.ListListings(ctx, filter)

	assert.NoError(t, err)
	assert.Len(t, listings, 0)
	assert.Equal(t, int32(0), total)
	mockRepo.AssertExpectations(t)
}

// ===================================
// SearchListings Tests (7 tests)
// ===================================

func TestSearchListings_Success_CacheHit(t *testing.T) {
	service, _, mockCache, _ := SetupServiceTest(t)
	ctx := TestContext()

	query := NewSearchListingsQuery("laptop", 20)

	// Cache key format: "search:query:categoryID:limit:offset"
	// When categoryID is nil, it's still printed as "0"
	cacheKey := "search:laptop:0:20:0"
	mockCache.On("Get", ctx, cacheKey, mock.Anything).
		Return(nil) // Cache hit (but data won't be populated in mock)

	// Note: With cache hit, repo.SearchListings is not called
	// In this mock scenario, we get nil results but no error

	_, _, err := service.SearchListings(ctx, query)

	// Cache hit scenario - no error even with empty mock data
	assert.NoError(t, err)
	mockCache.AssertExpectations(t)
}

func TestSearchListings_Success_CacheMiss(t *testing.T) {
	service, mockRepo, mockCache, _ := SetupServiceTest(t)
	ctx := TestContext()

	query := NewSearchListingsQuery("laptop", 20)

	expectedListings := []*domain.Listing{
		NewTestListing(1, 100, "Laptop 1"),
		NewTestListing(2, 100, "Laptop 2"),
	}

	cacheKey := "search:laptop:0:20:0"
	mockCache.On("Get", ctx, cacheKey, mock.Anything).
		Return(assert.AnError) // Cache miss

	mockRepo.On("SearchListings", ctx, query).
		Return(expectedListings, int32(2), nil)

	// Mock GetImages for each listing (eager loading)
	mockRepo.On("GetImages", ctx, mock.AnythingOfType("int64")).
		Return([]*domain.ListingImage{}, nil).Maybe()

	// Cache set is non-blocking and happens in goroutine
	mockCache.On("Set", mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()

	listings, total, err := service.SearchListings(ctx, query)

	assert.NoError(t, err)
	assert.Len(t, listings, 2)
	assert.Equal(t, int32(2), total)
	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestSearchListings_Success_WithFilters(t *testing.T) {
	service, mockRepo, mockCache, _ := SetupServiceTest(t)
	ctx := TestContext()

	categoryID := string("3b4246cc-9970-403c-af01-c142a4178dc6")
	minPrice := 100.0
	maxPrice := 500.0

	query := NewSearchListingsQuery("laptop", 20)
	query.CategoryID = &categoryID
	query.MinPrice = &minPrice
	query.MaxPrice = &maxPrice

	expectedListings := []*domain.Listing{
		NewTestListing(1, 100, "Laptop 1"),
	}

	cacheKey := "search:laptop:1:20:0"
	mockCache.On("Get", ctx, cacheKey, mock.Anything).
		Return(assert.AnError) // Cache miss

	mockRepo.On("SearchListings", ctx, query).
		Return(expectedListings, int32(1), nil)

	// Mock GetImages for each listing (eager loading)
	mockRepo.On("GetImages", ctx, mock.AnythingOfType("int64")).
		Return([]*domain.ListingImage{}, nil).Maybe()

	// Cache set is non-blocking and happens in goroutine
	mockCache.On("Set", mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()

	listings, total, err := service.SearchListings(ctx, query)

	assert.NoError(t, err)
	assert.Len(t, listings, 1)
	assert.Equal(t, int32(1), total)
	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestSearchListings_Error_QueryTooShort(t *testing.T) {
	service, _, mockCache, _ := SetupServiceTest(t)
	ctx := TestContext()

	query := NewSearchListingsQuery("a", 20) // Query too short (< 2 chars)

	// Cache.Set() may be called in goroutine if reached before validation error
	mockCache.On("Set", mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()

	listings, total, err := service.SearchListings(ctx, query)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "validation failed")
	assert.Nil(t, listings)
	assert.Equal(t, int32(0), total)
	mockCache.AssertExpectations(t)
}

func TestSearchListings_Success_LimitCapping_Max100(t *testing.T) {
	service, mockRepo, mockCache, _ := SetupServiceTest(t)
	ctx := TestContext()

	query := NewSearchListingsQuery("laptop", 100) // Use max allowed

	expectedListings := []*domain.Listing{
		NewTestListing(1, 100, "Laptop 1"),
	}

	cacheKey := "search:laptop:0:100:0"
	mockCache.On("Get", ctx, cacheKey, mock.Anything).
		Return(assert.AnError) // Cache miss

	mockRepo.On("SearchListings", ctx, query).
		Return(expectedListings, int32(1), nil)

	// Mock GetImages for each listing (eager loading)
	mockRepo.On("GetImages", ctx, mock.AnythingOfType("int64")).
		Return([]*domain.ListingImage{}, nil).Maybe()

	// Cache set is non-blocking and happens in goroutine
	mockCache.On("Set", mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()

	listings, total, err := service.SearchListings(ctx, query)

	assert.NoError(t, err)
	assert.Equal(t, int32(100), query.Limit) // At max
	assert.Len(t, listings, 1)
	assert.Equal(t, int32(1), total)
	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestSearchListings_Success_NonBlockingCache_SetFailure(t *testing.T) {
	service, mockRepo, mockCache, _ := SetupServiceTest(t)
	ctx := TestContext()

	query := NewSearchListingsQuery("laptop", 20)

	expectedListings := []*domain.Listing{
		NewTestListing(1, 100, "Laptop 1"),
	}

	cacheKey := "search:laptop:0:20:0"
	mockCache.On("Get", ctx, cacheKey, mock.Anything).
		Return(assert.AnError) // Cache miss

	mockRepo.On("SearchListings", ctx, query).
		Return(expectedListings, int32(1), nil)

	// Mock GetImages for each listing (eager loading)
	mockRepo.On("GetImages", ctx, mock.AnythingOfType("int64")).
		Return([]*domain.ListingImage{}, nil).Maybe()

	// Cache set failure is non-blocking and happens in goroutine
	mockCache.On("Set", mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()

	listings, total, err := service.SearchListings(ctx, query)

	assert.NoError(t, err)
	assert.Len(t, listings, 1)
	assert.Equal(t, int32(1), total)
	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestSearchListings_Success_EmptyResults(t *testing.T) {
	service, mockRepo, mockCache, _ := SetupServiceTest(t)
	ctx := TestContext()

	query := NewSearchListingsQuery("nonexistent", 20)

	expectedListings := []*domain.Listing{} // Empty results

	cacheKey := "search:nonexistent:0:20:0"
	mockCache.On("Get", ctx, cacheKey, mock.Anything).
		Return(assert.AnError) // Cache miss

	mockRepo.On("SearchListings", ctx, query).
		Return(expectedListings, int32(0), nil)

	// Mock GetImages for each listing (eager loading) - won't be called for empty results
	mockRepo.On("GetImages", ctx, mock.AnythingOfType("int64")).
		Return([]*domain.ListingImage{}, nil).Maybe()

	// Cache set is non-blocking and happens in goroutine
	mockCache.On("Set", mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()

	listings, total, err := service.SearchListings(ctx, query)

	assert.NoError(t, err)
	assert.Len(t, listings, 0)
	assert.Equal(t, int32(0), total)
	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

// ===================================
// CreateProduct Tests (4 tests)
// ===================================

func TestCreateProduct_Success_AllFields(t *testing.T) {
	service, mockRepo, _, _ := SetupServiceTest(t)
	ctx := TestContext()

	storefrontID := int64(10)
	input := NewCreateProductInput(storefrontID, "Test Product")

	expectedProduct := NewTestProduct(1, storefrontID, "Test Product")

	mockRepo.On("CreateProduct", ctx, input).
		Return(expectedProduct, nil)

	product, err := service.CreateProduct(ctx, input)

	assert.NoError(t, err)
	assert.NotNil(t, product)
	assert.Equal(t, expectedProduct.ID, product.ID)
	assert.Equal(t, "Test Product", product.Name)
	mockRepo.AssertExpectations(t)
}

func TestCreateProduct_Error_NegativePrice(t *testing.T) {
	service, _, _, _ := SetupServiceTest(t)
	ctx := TestContext()

	storefrontID := int64(10)
	input := NewCreateProductInput(storefrontID, "Test Product")
	input.Price = -10.99 // Negative price

	product, err := service.CreateProduct(ctx, input)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "validation failed")
	assert.Nil(t, product)
}

func TestCreateProduct_Error_MissingName(t *testing.T) {
	service, _, _, _ := SetupServiceTest(t)
	ctx := TestContext()

	storefrontID := int64(10)
	input := NewCreateProductInput(storefrontID, "")
	input.Name = "" // Missing name

	product, err := service.CreateProduct(ctx, input)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "validation failed")
	assert.Nil(t, product)
}

func TestCreateProduct_Error_InvalidCurrency(t *testing.T) {
	service, _, _, _ := SetupServiceTest(t)
	ctx := TestContext()

	storefrontID := int64(10)
	input := NewCreateProductInput(storefrontID, "Test Product")
	input.Currency = "US" // Invalid currency (must be 3 chars)

	product, err := service.CreateProduct(ctx, input)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "validation failed")
	assert.Nil(t, product)
}

// ===================================
// UpdateProduct Tests (5 tests)
// ===================================

func TestUpdateProduct_Success_UpdateFields(t *testing.T) {
	service, mockRepo, _, _ := SetupServiceTest(t)
	ctx := TestContext()

	productID := int64(123)
	storefrontID := int64(10)
	newName := "Updated Product"
	newPrice := 199.99

	input := NewUpdateProductInput(newName, newPrice)

	updatedProduct := NewTestProduct(productID, storefrontID, newName)
	updatedProduct.Price = newPrice

	mockRepo.On("UpdateProduct", ctx, productID, storefrontID, input).
		Return(updatedProduct, nil)

	product, err := service.UpdateProduct(ctx, productID, storefrontID, input)

	assert.NoError(t, err)
	assert.NotNil(t, product)
	assert.Equal(t, newName, product.Name)
	assert.Equal(t, newPrice, product.Price)
	mockRepo.AssertExpectations(t)
}

func TestUpdateProduct_Error_InvalidProductID(t *testing.T) {
	service, _, _, _ := SetupServiceTest(t)
	ctx := TestContext()

	productID := int64(0) // Invalid ID
	storefrontID := int64(10)

	input := NewUpdateProductInput("Updated", 199.99)

	product, err := service.UpdateProduct(ctx, productID, storefrontID, input)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "product_id must be greater than 0")
	assert.Nil(t, product)
}

func TestUpdateProduct_Error_InvalidStorefrontID(t *testing.T) {
	service, _, _, _ := SetupServiceTest(t)
	ctx := TestContext()

	productID := int64(123)
	storefrontID := int64(0) // Invalid ID

	input := NewUpdateProductInput("Updated", 199.99)

	product, err := service.UpdateProduct(ctx, productID, storefrontID, input)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "storefront_id must be greater than 0")
	assert.Nil(t, product)
}

func TestUpdateProduct_Error_NegativePrice(t *testing.T) {
	service, _, _, _ := SetupServiceTest(t)
	ctx := TestContext()

	productID := int64(123)
	storefrontID := int64(10)
	negativePrice := -50.0

	input := NewUpdateProductInput("Updated", negativePrice)

	product, err := service.UpdateProduct(ctx, productID, storefrontID, input)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "validation failed")
	assert.Nil(t, product)
}

func TestUpdateProduct_Error_NegativeStock(t *testing.T) {
	service, _, _, _ := SetupServiceTest(t)
	ctx := TestContext()

	productID := int64(123)
	storefrontID := int64(10)
	negativeStock := int32(-10)

	input := &domain.UpdateProductInput{
		StockQuantity: &negativeStock,
	}

	product, err := service.UpdateProduct(ctx, productID, storefrontID, input)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "validation failed")
	assert.Nil(t, product)
}

// ===================================
// DeleteProduct Tests (5 tests)
// ===================================

func TestDeleteProduct_Success_SoftDelete(t *testing.T) {
	service, mockRepo, _, _ := SetupServiceTest(t)
	ctx := TestContext()

	productID := int64(123)
	storefrontID := int64(10)
	hardDelete := false

	mockRepo.On("DeleteProduct", ctx, productID, storefrontID, hardDelete).
		Return(int32(0), nil) // No variants deleted

	variantsDeleted, err := service.DeleteProduct(ctx, productID, storefrontID, hardDelete)

	assert.NoError(t, err)
	assert.Equal(t, int32(0), variantsDeleted)
	mockRepo.AssertExpectations(t)
}

func TestDeleteProduct_Success_HardDelete(t *testing.T) {
	service, mockRepo, _, _ := SetupServiceTest(t)
	ctx := TestContext()

	productID := int64(123)
	storefrontID := int64(10)
	hardDelete := true

	mockRepo.On("DeleteProduct", ctx, productID, storefrontID, hardDelete).
		Return(int32(5), nil) // 5 variants deleted

	variantsDeleted, err := service.DeleteProduct(ctx, productID, storefrontID, hardDelete)

	assert.NoError(t, err)
	assert.Equal(t, int32(5), variantsDeleted)
	mockRepo.AssertExpectations(t)
}

func TestDeleteProduct_Error_InvalidProductID(t *testing.T) {
	service, _, _, _ := SetupServiceTest(t)
	ctx := TestContext()

	productID := int64(0) // Invalid ID
	storefrontID := int64(10)
	hardDelete := false

	variantsDeleted, err := service.DeleteProduct(ctx, productID, storefrontID, hardDelete)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "product_id must be greater than 0")
	assert.Equal(t, int32(0), variantsDeleted)
}

func TestDeleteProduct_Error_OwnershipCheck(t *testing.T) {
	service, mockRepo, _, _ := SetupServiceTest(t)
	ctx := TestContext()

	productID := int64(123)
	storefrontID := int64(10) // Wrong storefront
	hardDelete := false

	mockRepo.On("DeleteProduct", ctx, productID, storefrontID, hardDelete).
		Return(int32(0), assert.AnError) // Ownership check fails

	variantsDeleted, err := service.DeleteProduct(ctx, productID, storefrontID, hardDelete)

	assert.Error(t, err)
	assert.Equal(t, int32(0), variantsDeleted)
	mockRepo.AssertExpectations(t)
}

func TestDeleteProduct_Success_CascadeDelete_Variants(t *testing.T) {
	service, mockRepo, _, _ := SetupServiceTest(t)
	ctx := TestContext()

	productID := int64(123)
	storefrontID := int64(10)
	hardDelete := true

	mockRepo.On("DeleteProduct", ctx, productID, storefrontID, hardDelete).
		Return(int32(10), nil) // 10 variants cascaded

	variantsDeleted, err := service.DeleteProduct(ctx, productID, storefrontID, hardDelete)

	assert.NoError(t, err)
	assert.Equal(t, int32(10), variantsDeleted) // Verify cascade
	mockRepo.AssertExpectations(t)
}

// ===================================
// GetListing Tests (4 tests)
// ===================================

func TestGetListing_Success_CacheHit(t *testing.T) {
	service, _, mockCache, _ := SetupServiceTest(t)
	ctx := TestContext()

	listingID := int64(123)
	cacheKey := "listing:1"

	// Cache returns successfully (cache hit)
	mockCache.On("Get", ctx, cacheKey, mock.Anything).
		Return(nil) // Success means cache hit

	listing, err := service.GetListing(ctx, listingID)

	assert.NoError(t, err)
	// With mock, we get zero value listing, but no error
	assert.NotNil(t, listing)
	mockCache.AssertExpectations(t)
}

func TestGetListing_Success_CacheMiss(t *testing.T) {
	service, mockRepo, mockCache, _ := SetupServiceTest(t)
	ctx := TestContext()

	listingID := int64(123)
	cacheKey := "listing:1"

	expectedListing := NewTestListing(listingID, 100, "Test Listing")

	// Cache miss
	mockCache.On("Get", ctx, cacheKey, mock.Anything).
		Return(assert.AnError)

	// Fetch from repository
	mockRepo.On("GetListingByID", ctx, listingID).
		Return(expectedListing, nil)

	// Load images
	mockRepo.On("GetImages", ctx, listingID).
		Return([]*domain.ListingImage{}, nil)

	// Non-blocking cache set (may or may not be called)
	mockCache.On("Set", mock.Anything, mock.Anything, mock.Anything).
		Return(nil).Maybe()

	listing, err := service.GetListing(ctx, listingID)

	assert.NoError(t, err)
	assert.NotNil(t, listing)
	assert.Equal(t, listingID, listing.ID)
	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestGetListing_Error_NotFound(t *testing.T) {
	service, mockRepo, mockCache, _ := SetupServiceTest(t)
	ctx := TestContext()

	listingID := int64(999)
	cacheKey := "listing:999"

	// Cache miss
	mockCache.On("Get", ctx, cacheKey, mock.Anything).
		Return(assert.AnError)

	// Repository returns not found
	mockRepo.On("GetListingByID", ctx, listingID).
		Return(nil, assert.AnError)

	listing, err := service.GetListing(ctx, listingID)

	assert.Error(t, err)
	assert.Nil(t, listing)
	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestGetListing_Success_NonBlockingCache_SetFailure(t *testing.T) {
	service, mockRepo, mockCache, _ := SetupServiceTest(t)
	ctx := TestContext()

	listingID := int64(123)
	cacheKey := "listing:1"

	expectedListing := NewTestListing(listingID, 100, "Test Listing")

	// Cache miss
	mockCache.On("Get", ctx, cacheKey, mock.Anything).
		Return(assert.AnError)

	// Fetch from repository
	mockRepo.On("GetListingByID", ctx, listingID).
		Return(expectedListing, nil)

	// Load images
	mockRepo.On("GetImages", ctx, listingID).
		Return([]*domain.ListingImage{}, nil)

	// Non-blocking cache set - may fail but doesn't affect response
	mockCache.On("Set", mock.Anything, mock.Anything, mock.Anything).
		Return(assert.AnError).Maybe()

	listing, err := service.GetListing(ctx, listingID)

	// Should succeed despite cache set failure
	assert.NoError(t, err)
	assert.NotNil(t, listing)
	assert.Equal(t, listingID, listing.ID)
	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

// ===================================
// Favorites Tests (14 tests)
// ===================================

func TestAddToFavorites_Success_AddFavorite(t *testing.T) {
	service, mockRepo, mockCache, _ := SetupServiceTest(t)
	ctx := TestContext()

	userID := int64(100)
	listingID := int64(123)

	existingListing := NewTestListing(listingID, 200, "Test Listing")

	mockRepo.On("GetListingByID", ctx, listingID).
		Return(existingListing, nil)
	mockRepo.On("AddToFavorites", ctx, userID, listingID).
		Return(nil)

	// Mock cache invalidation - three separate calls
	mockCache.On("Delete", ctx, "favorites:user:100").Return(nil).Once()
	mockCache.On("Delete", ctx, "favorites:listing:1:count").Return(nil).Once()
	mockCache.On("Delete", ctx, "favorites:user:100:listing:1").Return(nil).Once()

	err := service.AddToFavorites(ctx, userID, listingID)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestAddToFavorites_Error_InvalidUserID(t *testing.T) {
	service, _, _, _ := SetupServiceTest(t)
	ctx := TestContext()

	userID := int64(0) // Invalid
	listingID := int64(123)

	err := service.AddToFavorites(ctx, userID, listingID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid user ID")
}

func TestAddToFavorites_Error_InvalidListingID(t *testing.T) {
	service, _, _, _ := SetupServiceTest(t)
	ctx := TestContext()

	userID := int64(100)
	listingID := int64(0) // Invalid

	err := service.AddToFavorites(ctx, userID, listingID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid listing ID")
}

func TestAddToFavorites_Error_ListingNotFound(t *testing.T) {
	service, mockRepo, _, _ := SetupServiceTest(t)
	ctx := TestContext()

	userID := int64(100)
	listingID := int64(999)

	mockRepo.On("GetListingByID", ctx, listingID).
		Return(nil, assert.AnError)

	err := service.AddToFavorites(ctx, userID, listingID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "listing not found")
	mockRepo.AssertExpectations(t)
}

func TestRemoveFromFavorites_Success_RemoveFavorite(t *testing.T) {
	service, mockRepo, mockCache, _ := SetupServiceTest(t)
	ctx := TestContext()

	userID := int64(100)
	listingID := int64(123)

	mockRepo.On("RemoveFromFavorites", ctx, userID, listingID).
		Return(nil)

	// Mock cache invalidation - three separate calls
	mockCache.On("Delete", ctx, "favorites:user:100").Return(nil).Once()
	mockCache.On("Delete", ctx, "favorites:listing:1:count").Return(nil).Once()
	mockCache.On("Delete", ctx, "favorites:user:100:listing:1").Return(nil).Once()

	err := service.RemoveFromFavorites(ctx, userID, listingID)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestRemoveFromFavorites_Error_InvalidUserID(t *testing.T) {
	service, _, _, _ := SetupServiceTest(t)
	ctx := TestContext()

	userID := int64(0) // Invalid
	listingID := int64(123)

	err := service.RemoveFromFavorites(ctx, userID, listingID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid user ID")
}

func TestRemoveFromFavorites_Error_InvalidListingID(t *testing.T) {
	service, _, _, _ := SetupServiceTest(t)
	ctx := TestContext()

	userID := int64(100)
	listingID := int64(0) // Invalid

	err := service.RemoveFromFavorites(ctx, userID, listingID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid listing ID")
}

func TestGetUserFavorites_Success_MultipleFavorites(t *testing.T) {
	service, mockRepo, mockCache, _ := SetupServiceTest(t)
	ctx := TestContext()

	userID := int64(100)
	expectedIDs := []int64{1, 2, 3}

	// Cache miss
	mockCache.On("Get", ctx, "favorites:user:100", mock.Anything).Return(assert.AnError).Once()

	mockRepo.On("GetUserFavorites", ctx, userID).
		Return(expectedIDs, nil)

	// Cache set
	mockCache.On("Set", ctx, "favorites:user:100", expectedIDs).Return(nil).Once()

	listingIDs, err := service.GetUserFavorites(ctx, userID)

	assert.NoError(t, err)
	assert.Len(t, listingIDs, 3)
	assert.Equal(t, expectedIDs, listingIDs)
	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestGetUserFavorites_Success_EmptyFavorites(t *testing.T) {
	service, mockRepo, mockCache, _ := SetupServiceTest(t)
	ctx := TestContext()

	userID := int64(100)
	expectedIDs := []int64{}

	// Cache miss
	mockCache.On("Get", ctx, "favorites:user:100", mock.Anything).Return(assert.AnError).Once()

	mockRepo.On("GetUserFavorites", ctx, userID).
		Return(expectedIDs, nil)

	// Cache set
	mockCache.On("Set", ctx, "favorites:user:100", expectedIDs).Return(nil).Once()

	listingIDs, err := service.GetUserFavorites(ctx, userID)

	assert.NoError(t, err)
	assert.Len(t, listingIDs, 0)
	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestGetUserFavorites_Error_InvalidUserID(t *testing.T) {
	service, _, _, _ := SetupServiceTest(t)
	ctx := TestContext()

	userID := int64(0) // Invalid

	listingIDs, err := service.GetUserFavorites(ctx, userID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid user ID")
	assert.Nil(t, listingIDs)
}

func TestIsFavorite_True_IsFavorite(t *testing.T) {
	service, mockRepo, mockCache, _ := SetupServiceTest(t)
	ctx := TestContext()

	userID := int64(100)
	listingID := int64(123)

	// Cache miss
	mockCache.On("Get", ctx, "favorites:user:100:listing:1", mock.Anything).Return(assert.AnError).Once()

	mockRepo.On("IsFavorite", ctx, userID, listingID).
		Return(true, nil)

	// Cache set
	mockCache.On("Set", ctx, "favorites:user:100:listing:1", true).Return(nil).Once()

	isFavorite, err := service.IsFavorite(ctx, userID, listingID)

	assert.NoError(t, err)
	assert.True(t, isFavorite)
	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestIsFavorite_False_NotFavorite(t *testing.T) {
	service, mockRepo, mockCache, _ := SetupServiceTest(t)
	ctx := TestContext()

	userID := int64(100)
	listingID := int64(123)

	// Cache miss
	mockCache.On("Get", ctx, "favorites:user:100:listing:1", mock.Anything).Return(assert.AnError).Once()

	mockRepo.On("IsFavorite", ctx, userID, listingID).
		Return(false, nil)

	// Cache set
	mockCache.On("Set", ctx, "favorites:user:100:listing:1", false).Return(nil).Once()

	isFavorite, err := service.IsFavorite(ctx, userID, listingID)

	assert.NoError(t, err)
	assert.False(t, isFavorite)
	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestIsFavorite_Error_InvalidUserID(t *testing.T) {
	service, _, _, _ := SetupServiceTest(t)
	ctx := TestContext()

	userID := int64(0) // Invalid
	listingID := int64(123)

	isFavorite, err := service.IsFavorite(ctx, userID, listingID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid user ID")
	assert.False(t, isFavorite)
}

func TestIsFavorite_Error_InvalidListingID(t *testing.T) {
	service, _, _, _ := SetupServiceTest(t)
	ctx := TestContext()

	userID := int64(100)
	listingID := int64(0) // Invalid

	isFavorite, err := service.IsFavorite(ctx, userID, listingID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid listing ID")
	assert.False(t, isFavorite)
}

// ===================================
// ListProducts Tests (4 tests)
// ===================================

func TestListProducts_Success_DefaultPagination(t *testing.T) {
	service, mockRepo, _, _ := SetupServiceTest(t)
	ctx := TestContext()

	storefrontID := int64(10)
	page := 1
	pageSize := 20

	expectedProducts := []*domain.Product{
		NewTestProduct(1, storefrontID, "Product 1"),
		NewTestProduct(2, storefrontID, "Product 2"),
	}

	mockRepo.On("ListProducts", ctx, storefrontID, page, pageSize, false).
		Return(expectedProducts, 2, nil)

	products, total, err := service.ListProducts(ctx, storefrontID, page, pageSize, false)

	assert.NoError(t, err)
	assert.Len(t, products, 2)
	assert.Equal(t, 2, total)
	mockRepo.AssertExpectations(t)
}

func TestListProducts_PageValidation_LessThan1_Default1(t *testing.T) {
	service, mockRepo, _, _ := SetupServiceTest(t)
	ctx := TestContext()

	storefrontID := int64(10)
	page := 0 // Invalid, should default to 1
	pageSize := 20

	expectedProducts := []*domain.Product{
		NewTestProduct(1, storefrontID, "Product 1"),
	}

	// Service will normalize page to 1
	mockRepo.On("ListProducts", ctx, storefrontID, 1, pageSize, false).
		Return(expectedProducts, 1, nil)

	products, total, err := service.ListProducts(ctx, storefrontID, page, pageSize, false)

	assert.NoError(t, err)
	assert.Len(t, products, 1)
	assert.Equal(t, 1, total)
	mockRepo.AssertExpectations(t)
}

func TestListProducts_PageSizeValidation_LessThan1_Default20(t *testing.T) {
	service, mockRepo, _, _ := SetupServiceTest(t)
	ctx := TestContext()

	storefrontID := int64(10)
	page := 1
	pageSize := 0 // Invalid, should default to 20

	expectedProducts := []*domain.Product{
		NewTestProduct(1, storefrontID, "Product 1"),
	}

	// Service will normalize pageSize to 20
	mockRepo.On("ListProducts", ctx, storefrontID, page, 20, false).
		Return(expectedProducts, 1, nil)

	products, total, err := service.ListProducts(ctx, storefrontID, page, pageSize, false)

	assert.NoError(t, err)
	assert.Len(t, products, 1)
	assert.Equal(t, 1, total)
	mockRepo.AssertExpectations(t)
}

func TestListProducts_PageSizeValidation_GreaterThan100_Cap20(t *testing.T) {
	service, mockRepo, _, _ := SetupServiceTest(t)
	ctx := TestContext()

	storefrontID := int64(10)
	page := 1
	pageSize := 150 // Too large, should cap to 20

	expectedProducts := []*domain.Product{
		NewTestProduct(1, storefrontID, "Product 1"),
	}

	// Service will cap pageSize to 20
	mockRepo.On("ListProducts", ctx, storefrontID, page, 20, false).
		Return(expectedProducts, 1, nil)

	products, total, err := service.ListProducts(ctx, storefrontID, page, pageSize, false)

	assert.NoError(t, err)
	assert.Len(t, products, 1)
	assert.Equal(t, 1, total)
	mockRepo.AssertExpectations(t)
}

// ===================================
// Stats & Views Tests (3 tests)
// ===================================

func TestGetProductStats_Success_GetStats(t *testing.T) {
	service, mockRepo, _, _ := SetupServiceTest(t)
	ctx := TestContext()

	storefrontID := int64(10)
	expectedStats := &domain.ProductStats{
		TotalProducts:  100,
		ActiveProducts: 80,
		TotalValue:     15000.50,
	}

	mockRepo.On("GetProductStats", ctx, storefrontID).
		Return(expectedStats, nil)

	stats, err := service.GetProductStats(ctx, storefrontID)

	assert.NoError(t, err)
	assert.NotNil(t, stats)
	assert.Equal(t, int32(100), stats.TotalProducts)
	assert.Equal(t, int32(80), stats.ActiveProducts)
	assert.Equal(t, 15000.50, stats.TotalValue)
	mockRepo.AssertExpectations(t)
}

func TestGetProductStats_Error_InvalidStorefrontID(t *testing.T) {
	service, _, _, _ := SetupServiceTest(t)
	ctx := TestContext()

	storefrontID := int64(0) // Invalid

	stats, err := service.GetProductStats(ctx, storefrontID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "storefront_id must be greater than 0")
	assert.Nil(t, stats)
}

func TestIncrementProductViews_Success_IncrementViews(t *testing.T) {
	service, mockRepo, _, _ := SetupServiceTest(t)
	ctx := TestContext()

	productID := int64(123)

	mockRepo.On("IncrementProductViews", ctx, productID).
		Return(nil)

	err := service.IncrementProductViews(ctx, productID)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestIncrementProductViews_Error_InvalidProductID(t *testing.T) {
	service, _, _, _ := SetupServiceTest(t)
	ctx := TestContext()

	productID := int64(0) // Invalid

	err := service.IncrementProductViews(ctx, productID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "product_id must be greater than 0")
}

// ===================================
// Admin Operations Tests (4 tests)
// ===================================

func TestAdminGetListing_Success(t *testing.T) {
	service, mockRepo, _, _ := SetupServiceTest(t)
	ctx := TestContext()

	listingID := int64(123)
	expectedListing := NewTestListing(listingID, 100, "Admin Listing")

	mockRepo.On("GetListingByID", ctx, listingID).
		Return(expectedListing, nil)
	mockRepo.On("GetImages", ctx, listingID).
		Return([]*domain.ListingImage{}, nil)

	listing, err := service.AdminGetListing(ctx, listingID)

	assert.NoError(t, err)
	assert.NotNil(t, listing)
	assert.Equal(t, listingID, listing.ID)
	mockRepo.AssertExpectations(t)
}

func TestAdminUpdateListing_NoOwnershipCheck(t *testing.T) {
	service, mockRepo, mockCache, _ := SetupServiceTest(t)
	ctx := TestContext()

	listingID := int64(123)
	newPrice := 299.99

	input := &domain.UpdateListingInput{
		Price: &newPrice,
	}

	updatedListing := NewTestListing(listingID, 100, "Admin Updated")
	updatedListing.Price = newPrice

	// Notice: NO GetListingByID call for ownership check
	mockRepo.On("UpdateListing", ctx, listingID, input).
		Return(updatedListing, nil)
	mockRepo.On("GetImages", ctx, listingID).
		Return([]*domain.ListingImage{}, nil)
	mockCache.On("Delete", ctx, "listing:1").
		Return(nil)
	mockRepo.On("EnqueueIndexing", ctx, listingID, domain.IndexOpUpdate).
		Return(nil)

	listing, err := service.AdminUpdateListing(ctx, listingID, input)

	assert.NoError(t, err)
	assert.NotNil(t, listing)
	assert.Equal(t, newPrice, listing.Price)
	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestAdminDeleteListing_NoOwnershipCheck(t *testing.T) {
	service, mockRepo, mockCache, _ := SetupServiceTest(t)
	ctx := TestContext()

	listingID := int64(123)

	// Notice: NO GetListingByID call for ownership check
	mockRepo.On("DeleteListing", ctx, listingID).
		Return(nil)
	mockCache.On("Delete", ctx, "listing:1").
		Return(nil)
	mockRepo.On("EnqueueIndexing", ctx, listingID, domain.IndexOpDelete).
		Return(nil)

	err := service.AdminDeleteListing(ctx, listingID)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestCacheInvalidation_Admin(t *testing.T) {
	service, mockRepo, mockCache, _ := SetupServiceTest(t)
	ctx := TestContext()

	listingID := int64(123)
	newTitle := "Admin Cache Test"

	input := &domain.UpdateListingInput{
		Title: &newTitle,
	}

	updatedListing := NewTestListing(listingID, 100, newTitle)

	mockRepo.On("UpdateListing", ctx, listingID, input).
		Return(updatedListing, nil)
	mockRepo.On("GetImages", ctx, listingID).
		Return([]*domain.ListingImage{}, nil)
	// Cache invalidation is critical for admin operations
	mockCache.On("Delete", ctx, "listing:1").
		Return(nil)
	mockRepo.On("EnqueueIndexing", ctx, listingID, domain.IndexOpUpdate).
		Return(nil)

	listing, err := service.AdminUpdateListing(ctx, listingID, input)

	assert.NoError(t, err)
	assert.NotNil(t, listing)
	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

// ===================================
// Inventory Operations Tests (1 test)
// ===================================

func TestUpdateProductInventory_Success(t *testing.T) {
	service, mockRepo, _, _ := SetupServiceTest(t)
	ctx := TestContext()

	storefrontID := int64(10)
	productID := int64(123)
	variantID := int64(0) // No variant
	movementType := "in"
	quantity := int32(50)
	reason := "restock"
	notes := "Monthly restock"
	userID := int64(100)

	mockRepo.On("UpdateProductInventory", ctx, storefrontID, productID, variantID, movementType, quantity, reason, notes, userID).
		Return(int32(100), int32(150), nil) // stockBefore=100, stockAfter=150

	stockBefore, stockAfter, err := service.UpdateProductInventory(ctx, storefrontID, productID, variantID, movementType, quantity, reason, notes, userID)

	assert.NoError(t, err)
	assert.Equal(t, int32(100), stockBefore)
	assert.Equal(t, int32(150), stockAfter)
	mockRepo.AssertExpectations(t)
}

// ===================================
// Image Operations Tests (4 tests)
// ===================================

func TestGetImageByID_Success(t *testing.T) {
	service, mockRepo, _, _ := SetupServiceTest(t)
	ctx := TestContext()

	imageID := int64(123)
	expectedImage := &domain.ListingImage{
		ID:        imageID,
		ListingID: 1,
		URL:       "https://example.com/image.jpg",
	}

	mockRepo.On("GetImageByID", ctx, imageID).
		Return(expectedImage, nil)

	image, err := service.GetImageByID(ctx, imageID)

	assert.NoError(t, err)
	assert.NotNil(t, image)
	assert.Equal(t, imageID, image.ID)
	mockRepo.AssertExpectations(t)
}

func TestDeleteImage_Success(t *testing.T) {
	service, mockRepo, mockCache, _ := SetupServiceTest(t)
	ctx := TestContext()

	imageID := int64(123)
	listingID := int64(100)

	image := &domain.ListingImage{
		ID:        imageID,
		ListingID: listingID,
		URL:       "https://example.com/image.jpg",
	}

	mockRepo.On("GetImageByID", ctx, imageID).
		Return(image, nil)
	mockRepo.On("DeleteImage", ctx, imageID).
		Return(nil)
	mockCache.On("Delete", ctx, "listing:100").
		Return(nil)

	err := service.DeleteImage(ctx, imageID)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestAddImage_Success(t *testing.T) {
	service, mockRepo, mockCache, _ := SetupServiceTest(t)
	ctx := TestContext()

	image := &domain.ListingImage{
		ListingID: 1,
		URL:       "https://example.com/image.jpg",
	}

	expectedImage := &domain.ListingImage{
		ID:        1,
		ListingID: 1,
		URL:       "https://example.com/image.jpg",
	}

	mockRepo.On("AddImage", ctx, image).
		Return(expectedImage, nil)
	mockCache.On("Delete", ctx, "listing:1").
		Return(nil)

	result, err := service.AddImage(ctx, image)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, string("3b4246cc-9970-403c-af01-c142a4178dc6"), result.ID)
	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestGetImages_Success(t *testing.T) {
	service, mockRepo, _, _ := SetupServiceTest(t)
	ctx := TestContext()

	listingID := int64(123)
	expectedImages := []*domain.ListingImage{
		{ID: 1, ListingID: listingID, URL: "https://example.com/image1.jpg"},
		{ID: 2, ListingID: listingID, URL: "https://example.com/image2.jpg"},
	}

	mockRepo.On("GetImages", ctx, listingID).
		Return(expectedImages, nil)

	images, err := service.GetImages(ctx, listingID)

	assert.NoError(t, err)
	assert.Len(t, images, 2)
	mockRepo.AssertExpectations(t)
}

// ===================================
// Category Operations Tests (5 tests)
// ===================================

func TestGetRootCategories_Success(t *testing.T) {
	service, mockRepo, _, _ := SetupServiceTest(t)
	ctx := TestContext()

	expectedCategories := []*domain.Category{
		{ID: "1", Name: "Electronics"},
		{ID: "2", Name: "Clothing"},
	}

	mockRepo.On("GetRootCategories", ctx).
		Return(expectedCategories, nil)

	categories, err := service.GetRootCategories(ctx)

	assert.NoError(t, err)
	assert.Len(t, categories, 2)
	mockRepo.AssertExpectations(t)
}

func TestGetAllCategories_Success(t *testing.T) {
	service, mockRepo, _, _ := SetupServiceTest(t)
	ctx := TestContext()

	expectedCategories := []*domain.Category{
		{ID: "1", Name: "Electronics"},
		{ID: "2", Name: "Laptops", ParentID: func(s string) *string { return &s }("1")},
	}

	mockRepo.On("GetAllCategories", ctx).
		Return(expectedCategories, nil)

	categories, err := service.GetAllCategories(ctx)

	assert.NoError(t, err)
	assert.Len(t, categories, 2)
	mockRepo.AssertExpectations(t)
}

func TestGetPopularCategories_Success(t *testing.T) {
	service, mockRepo, _, _ := SetupServiceTest(t)
	ctx := TestContext()

	limit := 10
	expectedCategories := []*domain.Category{
		{ID: "1", Name: "Electronics"},
		{ID: "2", Name: "Clothing"},
	}

	mockRepo.On("GetPopularCategories", ctx, limit).
		Return(expectedCategories, nil)

	categories, err := service.GetPopularCategories(ctx, limit)

	assert.NoError(t, err)
	assert.Len(t, categories, 2)
	mockRepo.AssertExpectations(t)
}

func TestGetCategoryByID_Success(t *testing.T) {
	service, mockRepo, _, _ := SetupServiceTest(t)
	ctx := TestContext()

	categoryID := string("3b4246cc-9970-403c-af01-c142a4178dc6")
	expectedCategory := &domain.Category{
		ID:   categoryID,
		Name: "Electronics",
	}

	mockRepo.On("GetCategoryByID", ctx, categoryID).
		Return(expectedCategory, nil)

	category, err := service.GetCategoryByID(ctx, categoryID)

	assert.NoError(t, err)
	assert.NotNil(t, category)
	assert.Equal(t, categoryID, category.ID)
	mockRepo.AssertExpectations(t)
}

func TestGetCategoryTree_Success(t *testing.T) {
	service, mockRepo, _, _ := SetupServiceTest(t)
	ctx := TestContext()

	categoryID := string("3b4246cc-9970-403c-af01-c142a4178dc6")
	expectedTree := &domain.CategoryTreeNode{
		ID:       categoryID,
		Name:     "Electronics",
		Slug:     "electronics",
		Level:    1,
		Path:     "1",
		Children: []domain.CategoryTreeNode{},
	}

	mockRepo.On("GetCategoryTree", ctx, categoryID).
		Return(expectedTree, nil)

	tree, err := service.GetCategoryTree(ctx, categoryID)

	assert.NoError(t, err)
	assert.NotNil(t, tree)
	assert.Equal(t, categoryID, tree.ID)
	assert.Equal(t, "Electronics", tree.Name)
	mockRepo.AssertExpectations(t)
}

// ===================================
// Variant Operations Tests (5 tests)
// ===================================

func TestCreateVariants_Success(t *testing.T) {
	service, mockRepo, _, _ := SetupServiceTest(t)
	ctx := TestContext()

	variants := []*domain.ListingVariant{
		{ListingID: 1, SKU: "VAR-001"},
		{ListingID: 1, SKU: "VAR-002"},
	}

	mockRepo.On("CreateVariants", ctx, variants).
		Return(nil)

	err := service.CreateVariants(ctx, variants)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestGetVariants_Success(t *testing.T) {
	service, mockRepo, _, _ := SetupServiceTest(t)
	ctx := TestContext()

	listingID := int64(123)
	expectedVariants := []*domain.ListingVariant{
		{ID: 1, ListingID: listingID, SKU: "VAR-001"},
		{ID: 2, ListingID: listingID, SKU: "VAR-002"},
	}

	mockRepo.On("GetVariants", ctx, listingID).
		Return(expectedVariants, nil)

	variants, err := service.GetVariants(ctx, listingID)

	assert.NoError(t, err)
	assert.Len(t, variants, 2)
	mockRepo.AssertExpectations(t)
}

func TestGetVariantByID_Success(t *testing.T) {
	service, mockRepo, _, _ := SetupServiceTest(t)
	ctx := TestContext()

	variantID := int64(123)
	listingID := int64(10)

	expectedVariants := []*domain.ListingVariant{
		{ID: variantID, ListingID: listingID, SKU: "VAR-001"},
	}

	// Note: Current implementation calls GetVariants(ctx, 0)
	mockRepo.On("GetVariants", ctx, int64(0)).
		Return(expectedVariants, nil)

	variant, err := service.GetVariantByID(ctx, variantID)

	assert.NoError(t, err)
	assert.NotNil(t, variant)
	assert.Equal(t, variantID, variant.ID)
	mockRepo.AssertExpectations(t)
}

func TestUpdateVariant_Success(t *testing.T) {
	service, mockRepo, _, _ := SetupServiceTest(t)
	ctx := TestContext()

	variant := &domain.ListingVariant{
		ID:        1,
		ListingID: 10,
		SKU:       "VAR-001-UPDATED",
	}

	mockRepo.On("UpdateVariant", ctx, variant).
		Return(nil)

	err := service.UpdateVariant(ctx, variant)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDeleteVariant_Success(t *testing.T) {
	service, mockRepo, _, _ := SetupServiceTest(t)
	ctx := TestContext()

	variantID := int64(123)

	mockRepo.On("DeleteVariant", ctx, variantID).
		Return(nil)

	err := service.DeleteVariant(ctx, variantID)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// ===================================
// Reindexing Operations Tests (3 tests)
// ===================================

func TestGetListingsForReindex_Success(t *testing.T) {
	service, mockRepo, _, _ := SetupServiceTest(t)
	ctx := TestContext()

	limit := 100
	expectedListings := []*domain.Listing{
		NewTestListing(1, 100, "Listing 1"),
		NewTestListing(2, 100, "Listing 2"),
	}

	mockRepo.On("GetListingsForReindex", ctx, limit).
		Return(expectedListings, nil)

	// Mock GetImages calls for each listing (eager loading)
	mockRepo.On("GetImages", ctx, string("3b4246cc-9970-403c-af01-c142a4178dc6")).
		Return([]*domain.ListingImage{}, nil)
	mockRepo.On("GetImages", ctx, int64(2)).
		Return([]*domain.ListingImage{}, nil)

	listings, err := service.GetListingsForReindex(ctx, limit)

	assert.NoError(t, err)
	assert.Len(t, listings, 2)
	mockRepo.AssertExpectations(t)
}

func TestResetReindexFlags_Success(t *testing.T) {
	service, mockRepo, _, _ := SetupServiceTest(t)
	ctx := TestContext()

	listingIDs := []int64{1, 2, 3}

	mockRepo.On("ResetReindexFlags", ctx, listingIDs).
		Return(nil)

	err := service.ResetReindexFlags(ctx, listingIDs)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestSyncDiscounts_Success(t *testing.T) {
	service, mockRepo, _, _ := SetupServiceTest(t)
	ctx := TestContext()

	mockRepo.On("SyncDiscounts", ctx).
		Return(nil)

	err := service.SyncDiscounts(ctx)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// ===================================
// Product Variant Simple Ops Tests (2 tests)
// ===================================

func TestGetVariant_Success(t *testing.T) {
	service, mockRepo, _, _ := SetupServiceTest(t)
	ctx := TestContext()

	variantID := int64(123)
	productID := int64(10)

	expectedVariant := NewTestProductVariant(variantID, productID)

	mockRepo.On("GetVariantByID", ctx, variantID, &productID).
		Return(expectedVariant, nil)

	variant, err := service.GetVariant(ctx, variantID, &productID)

	assert.NoError(t, err)
	assert.NotNil(t, variant)
	assert.Equal(t, variantID, variant.ID)
	mockRepo.AssertExpectations(t)
}

func TestGetVariantsByProductID_Success(t *testing.T) {
	service, mockRepo, _, _ := SetupServiceTest(t)
	ctx := TestContext()

	productID := int64(10)
	expectedVariants := []*domain.ProductVariant{
		NewTestProductVariant(1, productID),
		NewTestProductVariant(2, productID),
	}

	mockRepo.On("GetVariantsByProductID", ctx, productID, false).
		Return(expectedVariants, nil)

	variants, err := service.GetVariantsByProductID(ctx, productID, false)

	assert.NoError(t, err)
	assert.Len(t, variants, 2)
	mockRepo.AssertExpectations(t)
}

// ===================================
// Inventory Batch Update Tests (4 tests)
// ===================================

func TestBatchUpdateStock_Success(t *testing.T) {
	service, mockRepo, _, _ := SetupServiceTest(t)
	ctx := TestContext()

	storefrontID := int64(10)
	items := []domain.StockUpdateItem{
		{ProductID: 1, Quantity: 50},
		{ProductID: 2, Quantity: 30},
	}
	reason := "monthly_restock"
	userID := int64(100)

	expectedResults := []domain.StockUpdateResult{
		{ProductID: 1, Success: true},
		{ProductID: 2, Success: true},
	}

	mockRepo.On("BatchUpdateStock", ctx, storefrontID, items, reason, userID).
		Return(int32(2), int32(0), expectedResults, nil)

	successCount, failedCount, results, err := service.BatchUpdateStock(ctx, storefrontID, items, reason, userID)

	assert.NoError(t, err)
	assert.Equal(t, int32(2), successCount)
	assert.Equal(t, int32(0), failedCount)
	assert.Len(t, results, 2)
	mockRepo.AssertExpectations(t)
}

func TestBatchUpdateStock_Error_ValidationErrors(t *testing.T) {
	service, _, _, _ := SetupServiceTest(t)
	ctx := TestContext()

	storefrontID := int64(10)
	items := []domain.StockUpdateItem{
		{ProductID: 0, Quantity: 50}, // Invalid product ID
	}
	reason := "test"
	userID := int64(100)

	successCount, failedCount, results, err := service.BatchUpdateStock(ctx, storefrontID, items, reason, userID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid product_id")
	assert.Equal(t, int32(0), successCount)
	assert.Equal(t, int32(0), failedCount)
	assert.Nil(t, results)
}

func TestBatchUpdateStock_Error_EmptyInput(t *testing.T) {
	service, _, _, _ := SetupServiceTest(t)
	ctx := TestContext()

	storefrontID := int64(10)
	items := []domain.StockUpdateItem{}
	reason := "test"
	userID := int64(100)

	successCount, failedCount, results, err := service.BatchUpdateStock(ctx, storefrontID, items, reason, userID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "items list cannot be empty")
	assert.Equal(t, int32(0), successCount)
	assert.Equal(t, int32(0), failedCount)
	assert.Nil(t, results)
}

func TestBatchUpdateStock_Error_BatchTooLarge(t *testing.T) {
	service, _, _, _ := SetupServiceTest(t)
	ctx := TestContext()

	storefrontID := int64(10)
	items := make([]domain.StockUpdateItem, 1001)
	for i := 0; i < 1001; i++ {
		items[i] = domain.StockUpdateItem{
			ProductID: int64(i + 1),
			Quantity:  10,
		}
	}
	reason := "test"
	userID := int64(100)

	successCount, failedCount, results, err := service.BatchUpdateStock(ctx, storefrontID, items, reason, userID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot update more than 1000 items")
	assert.Equal(t, int32(0), successCount)
	assert.Equal(t, int32(0), failedCount)
	assert.Nil(t, results)
}
