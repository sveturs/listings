package listings

import (
	"context"
	"testing"
	"time"

	"github.com/rs/zerolog"

	"github.com/sveturs/listings/internal/domain"
	"github.com/sveturs/listings/internal/service/listings/mocks"
)

// TestContext creates a context with timeout for testing
func TestContext() context.Context {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	_ = cancel // Store cancel in test cleanup if needed
	return ctx
}

// SetupServiceTest creates a service instance with all mocks for testing
func SetupServiceTest(t *testing.T) (*Service, *mocks.MockRepository, *mocks.MockCacheRepository, *mocks.MockIndexingService) {
	t.Helper()

	mockRepo := new(mocks.MockRepository)
	mockCache := new(mocks.MockCacheRepository)
	mockIndexer := new(mocks.MockIndexingService)

	// Create zerolog logger that writes to test output
	logger := zerolog.New(zerolog.NewTestWriter(t)).With().Timestamp().Logger()

	service := NewService(mockRepo, mockCache, mockIndexer, logger)

	return service, mockRepo, mockCache, mockIndexer
}

// Helper functions for creating test data

// NewTestListing creates a test listing with default values
func NewTestListing(id int64, userID int64, title string) *domain.Listing {
	now := time.Now()
	description := "Test listing description"
	sku := "TEST-SKU-001"

	return &domain.Listing{
		ID:             id,
		UUID:           "test-uuid-" + time.Now().Format("20060102150405"),
		UserID:         userID,
		StorefrontID:   nil,
		Title:          title,
		Description:    &description,
		Price:          99.99,
		Currency:       "USD",
		CategoryID:     1,
		Status:         domain.StatusActive,
		Visibility:     domain.VisibilityPublic,
		Quantity:       10,
		SKU:            &sku,
		ViewsCount:     0,
		FavoritesCount: 0,
		CreatedAt:      now,
		UpdatedAt:      now,
		PublishedAt:    &now,
		DeletedAt:      nil,
		IsDeleted:      false,
	}
}

// NewCreateListingInput creates a test CreateListingInput
func NewCreateListingInput(userID int64, title string) *domain.CreateListingInput {
	description := "Test listing description"
	sku := "TEST-SKU-001"

	return &domain.CreateListingInput{
		UserID:      userID,
		Title:       title,
		Description: &description,
		Price:       99.99,
		Currency:    "USD",
		CategoryID:  1,
		Quantity:    10,
		SKU:         &sku,
		SourceType:  domain.SourceTypeC2C, // Required field added
	}
}

// NewUpdateListingInput creates a test UpdateListingInput
func NewUpdateListingInput(title string, price float64) *domain.UpdateListingInput {
	quantity := int32(15)
	status := domain.StatusActive

	return &domain.UpdateListingInput{
		Title:    &title,
		Price:    &price,
		Quantity: &quantity,
		Status:   &status,
	}
}

// NewListListingsFilter creates a test ListListingsFilter
func NewListListingsFilter(limit int32) *domain.ListListingsFilter {
	return &domain.ListListingsFilter{
		Limit:  limit,
		Offset: 0,
	}
}

// NewSearchListingsQuery creates a test SearchListingsQuery
func NewSearchListingsQuery(query string, limit int32) *domain.SearchListingsQuery {
	return &domain.SearchListingsQuery{
		Query:  query,
		Limit:  limit,
		Offset: 0,
	}
}

// NewTestCategory creates a test category
func NewTestCategory(id int64, name string) *domain.Category {
	slug := "test-category-" + time.Now().Format("20060102150405")

	return &domain.Category{
		ID:           id,
		Name:         name,
		Slug:         slug,
		ParentID:     nil,
		Icon:         nil,
		Description:  nil,
		IsActive:     true,
		ListingCount: 0,
		SortOrder:    0,
		Level:        0,
		HasCustomUI:  false,
		CreatedAt:    time.Now(),
	}
}

// NewTestListingImage creates a test listing image
func NewTestListingImage(id int64, listingID int64) *domain.ListingImage {
	thumbURL := "https://example.com/thumb.jpg"
	storagePath := "/storage/images/test.jpg"
	mimeType := "image/jpeg"
	width := int32(800)
	height := int32(600)
	fileSize := int64(102400)

	return &domain.ListingImage{
		ID:           id,
		ListingID:    listingID,
		URL:          "https://example.com/image.jpg",
		StoragePath:  &storagePath,
		ThumbnailURL: &thumbURL,
		DisplayOrder: 0,
		IsPrimary:    true,
		Width:        &width,
		Height:       &height,
		FileSize:     &fileSize,
		MimeType:     &mimeType,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}

// NewTestProduct creates a test product
func NewTestProduct(id int64, storefrontID int64, name string) *domain.Product {
	sku := "PROD-SKU-001"
	barcode := "1234567890123"

	return &domain.Product{
		ID:            id,
		StorefrontID:  storefrontID,
		Name:          name,
		Description:   "Test product description",
		Price:         149.99,
		Currency:      "USD",
		CategoryID:    1,
		SKU:           &sku,
		Barcode:       &barcode,
		StockQuantity: 50,
		StockStatus:   domain.StockStatusInStock,
		IsActive:      true,
		Attributes:    map[string]interface{}{"color": "blue", "size": "M"},
		ViewCount:     0,
		SoldCount:     0,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		HasVariants:   false,
		ShowOnMap:     false,
	}
}

// NewCreateProductInput creates a test CreateProductInput
func NewCreateProductInput(storefrontID int64, name string) *domain.CreateProductInput {
	sku := "PROD-SKU-001"
	barcode := "1234567890123"

	return &domain.CreateProductInput{
		StorefrontID:  storefrontID,
		Name:          name,
		Description:   "Test product description",
		Price:         149.99,
		Currency:      "USD",
		CategoryID:    1,
		SKU:           &sku,
		Barcode:       &barcode,
		StockQuantity: 50,
		Attributes:    map[string]interface{}{"color": "blue", "size": "M"},
	}
}

// NewUpdateProductInput creates a test UpdateProductInput
func NewUpdateProductInput(name string, price float64) *domain.UpdateProductInput {
	stockQuantity := int32(75)
	isActive := true

	return &domain.UpdateProductInput{
		Name:          &name,
		Price:         &price,
		StockQuantity: &stockQuantity,
		IsActive:      &isActive,
	}
}

// NewTestProductVariant creates a test product variant
func NewTestProductVariant(id int64, productID int64) *domain.ProductVariant {
	sku := "VAR-SKU-001"
	price := 159.99
	compareAtPrice := 199.99
	lowStockThreshold := int32(5)

	return &domain.ProductVariant{
		ID:                id,
		ProductID:         productID,
		SKU:               &sku,
		Price:             &price,
		CompareAtPrice:    &compareAtPrice,
		StockQuantity:     25,
		StockStatus:       domain.StockStatusInStock,
		LowStockThreshold: &lowStockThreshold,
		VariantAttributes: map[string]interface{}{"size": "L", "color": "red"},
		IsActive:          true,
		IsDefault:         false,
		ViewCount:         0,
		SoldCount:         0,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}
}

// NewCreateVariantInput creates a test CreateVariantInput
func NewCreateVariantInput(productID int64) *domain.CreateVariantInput {
	sku := "VAR-SKU-001"
	price := 159.99
	compareAtPrice := 199.99
	lowStockThreshold := int32(5)

	return &domain.CreateVariantInput{
		ProductID:         productID,
		SKU:               &sku,
		Price:             &price,
		CompareAtPrice:    &compareAtPrice,
		StockQuantity:     25,
		LowStockThreshold: &lowStockThreshold,
		VariantAttributes: map[string]interface{}{"size": "L", "color": "red"},
		IsDefault:         false,
	}
}

// NewUpdateVariantInput creates a test UpdateVariantInput
func NewUpdateVariantInput(price float64, stockQuantity int32) *domain.UpdateVariantInput {
	isActive := true

	return &domain.UpdateVariantInput{
		Price:         &price,
		StockQuantity: &stockQuantity,
		IsActive:      &isActive,
	}
}

// NewTestProductStats creates a test ProductStats
func NewTestProductStats() *domain.ProductStats {
	return &domain.ProductStats{
		TotalProducts:  100,
		ActiveProducts: 85,
		OutOfStock:     5,
		LowStock:       10,
		TotalValue:     14999.50,
		TotalSold:      250,
	}
}

// NewStockUpdateItem creates a test StockUpdateItem
func NewStockUpdateItem(productID int64, quantity int32) domain.StockUpdateItem {
	reason := "test stock update"

	return domain.StockUpdateItem{
		ProductID: productID,
		VariantID: nil,
		Quantity:  quantity,
		Reason:    &reason,
	}
}

// NewStockUpdateResult creates a test StockUpdateResult
func NewStockUpdateResult(productID int64, stockBefore, stockAfter int32, success bool) domain.StockUpdateResult {
	return domain.StockUpdateResult{
		ProductID:   productID,
		VariantID:   nil,
		StockBefore: stockBefore,
		StockAfter:  stockAfter,
		Success:     success,
		Error:       nil,
	}
}

// NewBulkUpdateProductInput creates a test BulkUpdateProductInput
func NewBulkUpdateProductInput(productID int64, name string, price float64) *domain.BulkUpdateProductInput {
	return &domain.BulkUpdateProductInput{
		ProductID: productID,
		Name:      &name,
		Price:     &price,
	}
}

// NewBulkUpdateProductsResult creates a test BulkUpdateProductsResult
func NewBulkUpdateProductsResult(products []*domain.Product, errors []domain.BulkUpdateError) *domain.BulkUpdateProductsResult {
	return &domain.BulkUpdateProductsResult{
		SuccessfulProducts: products,
		FailedUpdates:      errors,
	}
}
