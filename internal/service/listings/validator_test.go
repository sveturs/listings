package listings

import (
	"context"
	"database/sql"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/sveturs/listings/internal/domain"
)

// MockRepository is a mock implementation of Repository for testing
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) GetCategoryByID(ctx context.Context, categoryID int64) (*domain.Category, error) {
	args := m.Called(ctx, categoryID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Category), args.Error(1)
}

func (m *MockRepository) GetListingBySlug(ctx context.Context, slug string) (*domain.Listing, error) {
	args := m.Called(ctx, slug)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Listing), args.Error(1)
}

// Implement other required Repository methods as no-ops for testing
func (m *MockRepository) CreateListing(ctx context.Context, input *domain.CreateListingInput) (*domain.Listing, error) {
	return nil, nil
}
func (m *MockRepository) GetListingByID(ctx context.Context, id int64) (*domain.Listing, error) {
	return nil, nil
}
func (m *MockRepository) GetListingByUUID(ctx context.Context, uuid string) (*domain.Listing, error) {
	return nil, nil
}
func (m *MockRepository) UpdateListing(ctx context.Context, id int64, input *domain.UpdateListingInput) (*domain.Listing, error) {
	return nil, nil
}
func (m *MockRepository) DeleteListing(ctx context.Context, id int64) error { return nil }
func (m *MockRepository) ListListings(ctx context.Context, filter *domain.ListListingsFilter) ([]*domain.Listing, int32, error) {
	return nil, 0, nil
}
func (m *MockRepository) SearchListings(ctx context.Context, query *domain.SearchListingsQuery) ([]*domain.Listing, int32, error) {
	return nil, 0, nil
}
func (m *MockRepository) EnqueueIndexing(ctx context.Context, listingID int64, operation string) error {
	return nil
}
func (m *MockRepository) GetImageByID(ctx context.Context, imageID int64) (*domain.ListingImage, error) {
	return nil, nil
}
func (m *MockRepository) DeleteImage(ctx context.Context, imageID int64) error { return nil }
func (m *MockRepository) AddImage(ctx context.Context, image *domain.ListingImage) (*domain.ListingImage, error) {
	return nil, nil
}
func (m *MockRepository) GetImages(ctx context.Context, listingID int64) ([]*domain.ListingImage, error) {
	return nil, nil
}
func (m *MockRepository) GetRootCategories(ctx context.Context) ([]*domain.Category, error) {
	return nil, nil
}
func (m *MockRepository) GetAllCategories(ctx context.Context) ([]*domain.Category, error) {
	return nil, nil
}
func (m *MockRepository) GetPopularCategories(ctx context.Context, limit int) ([]*domain.Category, error) {
	return nil, nil
}
func (m *MockRepository) GetCategoryTree(ctx context.Context, categoryID int64) (*domain.CategoryTreeNode, error) {
	return nil, nil
}
func (m *MockRepository) GetFavoritedUsers(ctx context.Context, listingID int64) ([]int64, error) {
	return nil, nil
}
func (m *MockRepository) AddToFavorites(ctx context.Context, userID, listingID int64) error {
	return nil
}
func (m *MockRepository) RemoveFromFavorites(ctx context.Context, userID, listingID int64) error {
	return nil
}
func (m *MockRepository) GetUserFavorites(ctx context.Context, userID int64) ([]int64, error) {
	return nil, nil
}
func (m *MockRepository) IsFavorite(ctx context.Context, userID, listingID int64) (bool, error) {
	return false, nil
}
func (m *MockRepository) CreateVariants(ctx context.Context, variants []*domain.ListingVariant) error {
	return nil
}
func (m *MockRepository) GetVariants(ctx context.Context, listingID int64) ([]*domain.ListingVariant, error) {
	return nil, nil
}
func (m *MockRepository) UpdateVariant(ctx context.Context, variant *domain.ListingVariant) error {
	return nil
}
func (m *MockRepository) DeleteVariant(ctx context.Context, variantID int64) error { return nil }
func (m *MockRepository) GetListingsForReindex(ctx context.Context, limit int) ([]*domain.Listing, error) {
	return nil, nil
}
func (m *MockRepository) ResetReindexFlags(ctx context.Context, listingIDs []int64) error {
	return nil
}
func (m *MockRepository) SyncDiscounts(ctx context.Context) error { return nil }

// Products operations
func (m *MockRepository) GetProductByID(ctx context.Context, productID int64, storefrontID *int64) (*domain.Product, error) {
	return nil, nil
}
func (m *MockRepository) GetProductsBySKUs(ctx context.Context, skus []string, storefrontID *int64) ([]*domain.Product, error) {
	return nil, nil
}
func (m *MockRepository) GetProductsByIDs(ctx context.Context, productIDs []int64, storefrontID *int64) ([]*domain.Product, error) {
	return nil, nil
}
func (m *MockRepository) ListProducts(ctx context.Context, storefrontID int64, page, pageSize int, isActiveOnly bool) ([]*domain.Product, int, error) {
	return nil, 0, nil
}
func (m *MockRepository) CreateProduct(ctx context.Context, input *domain.CreateProductInput) (*domain.Product, error) {
	return nil, nil
}
func (m *MockRepository) BulkCreateProducts(ctx context.Context, storefrontID int64, inputs []*domain.CreateProductInput) ([]*domain.Product, []domain.BulkProductError, error) {
	return nil, nil, nil
}
func (m *MockRepository) UpdateProduct(ctx context.Context, productID int64, storefrontID int64, input *domain.UpdateProductInput) (*domain.Product, error) {
	return nil, nil
}
func (m *MockRepository) DeleteProduct(ctx context.Context, productID, storefrontID int64, hardDelete bool) (int32, error) {
	return 0, nil
}
func (m *MockRepository) BulkDeleteProducts(ctx context.Context, storefrontID int64, productIDs []int64, hardDelete bool) (int32, int32, int32, map[int64]string, error) {
	return 0, 0, 0, nil, nil
}
func (m *MockRepository) BulkUpdateProducts(ctx context.Context, storefrontID int64, updates []*domain.BulkUpdateProductInput) (*domain.BulkUpdateProductsResult, error) {
	return nil, nil
}

// Product Variants operations
func (m *MockRepository) GetVariantByID(ctx context.Context, variantID int64, productID *int64) (*domain.ProductVariant, error) {
	return nil, nil
}
func (m *MockRepository) GetVariantsByProductID(ctx context.Context, productID int64, isActiveOnly bool) ([]*domain.ProductVariant, error) {
	return nil, nil
}
func (m *MockRepository) CreateProductVariant(ctx context.Context, input *domain.CreateVariantInput) (*domain.ProductVariant, error) {
	return nil, nil
}
func (m *MockRepository) UpdateProductVariant(ctx context.Context, variantID int64, productID int64, input *domain.UpdateVariantInput) (*domain.ProductVariant, error) {
	return nil, nil
}
func (m *MockRepository) DeleteProductVariant(ctx context.Context, variantID int64, productID int64) error {
	return nil
}
func (m *MockRepository) BulkCreateProductVariants(ctx context.Context, productID int64, inputs []*domain.CreateVariantInput) ([]*domain.ProductVariant, error) {
	return nil, nil
}

// Inventory Management operations
func (m *MockRepository) UpdateProductInventory(ctx context.Context, storefrontID, productID, variantID int64, movementType string, quantity int32, reason, notes string, userID int64) (int32, int32, error) {
	return 0, 0, nil
}
func (m *MockRepository) GetProductStats(ctx context.Context, storefrontID int64) (*domain.ProductStats, error) {
	return nil, nil
}
func (m *MockRepository) IncrementProductViews(ctx context.Context, productID int64) error {
	return nil
}
func (m *MockRepository) BatchUpdateStock(ctx context.Context, storefrontID int64, items []domain.StockUpdateItem, reason string, userID int64) (int32, int32, []domain.StockUpdateResult, error) {
	return 0, 0, nil, nil
}

// Storefront operations
func (m *MockRepository) GetStorefront(ctx context.Context, storefrontID int64) (*domain.Storefront, error) {
	return nil, nil
}
func (m *MockRepository) GetStorefrontBySlug(ctx context.Context, slug string) (*domain.Storefront, error) {
	return nil, nil
}
func (m *MockRepository) ListStorefronts(ctx context.Context, limit, offset int) ([]*domain.Storefront, int64, error) {
	return nil, 0, nil
}

// Transaction and database operations
func (m *MockRepository) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return nil, nil
}
func (m *MockRepository) GetDB() *sqlx.DB {
	return nil
}

func TestValidator_ValidatePrice(t *testing.T) {
	mockRepo := new(MockRepository)
	validator := NewValidator(mockRepo)

	tests := []struct {
		name    string
		price   float64
		wantErr bool
	}{
		{"valid price", 100.50, false},
		{"zero price", 0.0, false},
		{"negative price", -10.0, true},
		{"max price", MaxPrice, false},
		{"exceeds max price", MaxPrice + 1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidatePrice(tt.price)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidator_ValidateTitle(t *testing.T) {
	mockRepo := new(MockRepository)
	validator := NewValidator(mockRepo)

	tests := []struct {
		name    string
		title   string
		wantErr bool
	}{
		{"valid title", "Valid Product Title", false},
		{"min length title", "abc", false},
		{"too short title", "ab", true},
		{"empty title", "", true},
		{"whitespace only", "   ", true},
		{"max length title", string(make([]rune, MaxTitleLength)), false},
		{"exceeds max length", string(make([]rune, MaxTitleLength+1)), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateTitle(tt.title)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidator_ValidateDescription(t *testing.T) {
	mockRepo := new(MockRepository)
	validator := NewValidator(mockRepo)

	validDesc := "This is a valid description"
	longDesc := string(make([]rune, MaxDescriptionLength+1))

	tests := []struct {
		name    string
		desc    *string
		wantErr bool
	}{
		{"nil description", nil, false},
		{"valid description", &validDesc, false},
		{"exceeds max length", &longDesc, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateDescription(tt.desc)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidator_ValidateQuantity(t *testing.T) {
	mockRepo := new(MockRepository)
	validator := NewValidator(mockRepo)

	tests := []struct {
		name     string
		quantity int32
		wantErr  bool
	}{
		{"valid quantity", 10, false},
		{"zero quantity", 0, false},
		{"negative quantity", -5, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateQuantity(tt.quantity)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidator_ValidateCurrency(t *testing.T) {
	mockRepo := new(MockRepository)
	validator := NewValidator(mockRepo)

	tests := []struct {
		name     string
		currency string
		wantErr  bool
	}{
		{"valid RSD", "RSD", false},
		{"valid EUR", "EUR", false},
		{"valid USD", "USD", false},
		{"lowercase", "usd", true},
		{"too long", "USDD", true},
		{"too short", "US", true},
		{"empty", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateCurrency(tt.currency)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidator_ValidateCategory(t *testing.T) {
	mockRepo := new(MockRepository)
	validator := NewValidator(mockRepo)

	ctx := context.Background()

	// Setup mock for active category
	activeCategory := &domain.Category{
		ID:       1,
		Name:     "Electronics",
		IsActive: true,
	}
	mockRepo.On("GetCategoryByID", ctx, int64(1)).Return(activeCategory, nil)

	// Setup mock for inactive category
	inactiveCategory := &domain.Category{
		ID:       2,
		Name:     "Inactive Category",
		IsActive: false,
	}
	mockRepo.On("GetCategoryByID", ctx, int64(2)).Return(inactiveCategory, nil)

	tests := []struct {
		name       string
		categoryID int64
		wantErr    bool
	}{
		{"valid active category", 1, false},
		{"inactive category", 2, true},
		{"invalid category ID", 0, true},
		{"negative category ID", -1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateCategory(ctx, tt.categoryID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidator_ValidateStatusTransition(t *testing.T) {
	mockRepo := new(MockRepository)
	validator := NewValidator(mockRepo)

	tests := []struct {
		name    string
		from    string
		to      string
		wantErr bool
	}{
		{"draft to active", domain.StatusDraft, domain.StatusActive, false},
		{"draft to inactive", domain.StatusDraft, domain.StatusInactive, false},
		{"active to sold", domain.StatusActive, domain.StatusSold, false},
		{"active to inactive", domain.StatusActive, domain.StatusInactive, false},
		{"inactive to active", domain.StatusInactive, domain.StatusActive, false},
		{"sold to active", domain.StatusSold, domain.StatusActive, false},
		{"invalid: draft to sold", domain.StatusDraft, domain.StatusSold, true},
		{"invalid: sold to inactive", domain.StatusSold, domain.StatusInactive, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateStatusTransition(tt.from, tt.to)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidator_ValidateImages(t *testing.T) {
	mockRepo := new(MockRepository)
	validator := NewValidator(mockRepo)

	validFileSize := int64(1024 * 1024) // 1MB
	largeFileSize := int64(MaxImageSize + 1)
	validWidth := int32(800)
	validHeight := int32(600)
	smallWidth := int32(50)
	largeWidth := int32(MaxImageWidth + 1)
	validMimeType := "image/jpeg"
	invalidMimeType := "image/gif"

	tests := []struct {
		name    string
		images  []*domain.ListingImage
		wantErr bool
	}{
		{
			name: "valid images",
			images: []*domain.ListingImage{
				{FileSize: &validFileSize, Width: &validWidth, Height: &validHeight, MimeType: &validMimeType},
			},
			wantErr: false,
		},
		{
			name: "too many images",
			images: make([]*domain.ListingImage, MaxImagesPerListing+1),
			wantErr: true,
		},
		{
			name: "file size exceeds limit",
			images: []*domain.ListingImage{
				{FileSize: &largeFileSize},
			},
			wantErr: true,
		},
		{
			name: "dimensions too small",
			images: []*domain.ListingImage{
				{Width: &smallWidth, Height: &validHeight},
			},
			wantErr: true,
		},
		{
			name: "dimensions too large",
			images: []*domain.ListingImage{
				{Width: &largeWidth, Height: &validHeight},
			},
			wantErr: true,
		},
		{
			name: "invalid mime type",
			images: []*domain.ListingImage{
				{MimeType: &invalidMimeType},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateImages(tt.images)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateSlug(t *testing.T) {
	tests := []struct {
		name string
		slug string
		want bool
	}{
		{"valid slug", "valid-product-slug", true},
		{"valid with numbers", "product-123", true},
		{"uppercase not allowed", "Product-Slug", false},
		{"spaces not allowed", "product slug", false},
		{"empty slug", "", false},
		{"too long", string(make([]rune, 251)), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidateSlug(tt.slug)
			assert.Equal(t, tt.want, got)
		})
	}
}
