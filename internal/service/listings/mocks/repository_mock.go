// Package mocks provides test mocks for the listings service layer
package mocks

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/mock"

	"github.com/sveturs/listings/internal/domain"
)

// MockRepository is a mock implementation of listings.Repository interface
type MockRepository struct {
	mock.Mock
}

// CreateListing mocks creating a new listing
func (m *MockRepository) CreateListing(ctx context.Context, input *domain.CreateListingInput) (*domain.Listing, error) {
	args := m.Called(ctx, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Listing), args.Error(1)
}

// GetListingByID mocks getting a listing by ID
func (m *MockRepository) GetListingByID(ctx context.Context, id int64) (*domain.Listing, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Listing), args.Error(1)
}

// GetListingByUUID mocks getting a listing by UUID
func (m *MockRepository) GetListingByUUID(ctx context.Context, uuid string) (*domain.Listing, error) {
	args := m.Called(ctx, uuid)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Listing), args.Error(1)
}

// UpdateListing mocks updating a listing
func (m *MockRepository) UpdateListing(ctx context.Context, id int64, input *domain.UpdateListingInput) (*domain.Listing, error) {
	args := m.Called(ctx, id, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Listing), args.Error(1)
}

// DeleteListing mocks deleting a listing
func (m *MockRepository) DeleteListing(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// ListListings mocks listing listings with filter
func (m *MockRepository) ListListings(ctx context.Context, filter *domain.ListListingsFilter) ([]*domain.Listing, int32, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int32), args.Error(2)
	}
	return args.Get(0).([]*domain.Listing), args.Get(1).(int32), args.Error(2)
}

// SearchListings mocks searching listings
func (m *MockRepository) SearchListings(ctx context.Context, query *domain.SearchListingsQuery) ([]*domain.Listing, int32, error) {
	args := m.Called(ctx, query)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int32), args.Error(2)
	}
	return args.Get(0).([]*domain.Listing), args.Get(1).(int32), args.Error(2)
}

// EnqueueIndexing mocks enqueueing indexing operation
func (m *MockRepository) EnqueueIndexing(ctx context.Context, listingID int64, operation string) error {
	args := m.Called(ctx, listingID, operation)
	return args.Error(0)
}

// Image operations

// GetImageByID mocks getting an image by ID
func (m *MockRepository) GetImageByID(ctx context.Context, imageID int64) (*domain.ListingImage, error) {
	args := m.Called(ctx, imageID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ListingImage), args.Error(1)
}

// DeleteImage mocks deleting an image
func (m *MockRepository) DeleteImage(ctx context.Context, imageID int64) error {
	args := m.Called(ctx, imageID)
	return args.Error(0)
}

// AddImage mocks adding an image
func (m *MockRepository) AddImage(ctx context.Context, image *domain.ListingImage) (*domain.ListingImage, error) {
	args := m.Called(ctx, image)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ListingImage), args.Error(1)
}

// GetImages mocks getting images for a listing
func (m *MockRepository) GetImages(ctx context.Context, listingID int64) ([]*domain.ListingImage, error) {
	args := m.Called(ctx, listingID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.ListingImage), args.Error(1)
}

// Category operations

// GetRootCategories mocks getting root categories
func (m *MockRepository) GetRootCategories(ctx context.Context) ([]*domain.Category, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Category), args.Error(1)
}

// GetAllCategories mocks getting all categories
func (m *MockRepository) GetAllCategories(ctx context.Context) ([]*domain.Category, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Category), args.Error(1)
}

// GetPopularCategories mocks getting popular categories
func (m *MockRepository) GetPopularCategories(ctx context.Context, limit int) ([]*domain.Category, error) {
	args := m.Called(ctx, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Category), args.Error(1)
}

// GetCategoryByID mocks getting a category by ID
func (m *MockRepository) GetCategoryByID(ctx context.Context, categoryID int64) (*domain.Category, error) {
	args := m.Called(ctx, categoryID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Category), args.Error(1)
}

// GetCategoryTree mocks getting category tree
func (m *MockRepository) GetCategoryTree(ctx context.Context, categoryID int64) (*domain.CategoryTreeNode, error) {
	args := m.Called(ctx, categoryID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.CategoryTreeNode), args.Error(1)
}

// Favorites operations

// GetFavoritedUsers mocks getting users who favorited a listing
func (m *MockRepository) GetFavoritedUsers(ctx context.Context, listingID int64) ([]int64, error) {
	args := m.Called(ctx, listingID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]int64), args.Error(1)
}

// AddToFavorites mocks adding to favorites
func (m *MockRepository) AddToFavorites(ctx context.Context, userID, listingID int64) error {
	args := m.Called(ctx, userID, listingID)
	return args.Error(0)
}

// RemoveFromFavorites mocks removing from favorites
func (m *MockRepository) RemoveFromFavorites(ctx context.Context, userID, listingID int64) error {
	args := m.Called(ctx, userID, listingID)
	return args.Error(0)
}

// GetUserFavorites mocks getting user favorites
func (m *MockRepository) GetUserFavorites(ctx context.Context, userID int64) ([]int64, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]int64), args.Error(1)
}

// IsFavorite mocks checking if a listing is favorited
func (m *MockRepository) IsFavorite(ctx context.Context, userID, listingID int64) (bool, error) {
	args := m.Called(ctx, userID, listingID)
	return args.Get(0).(bool), args.Error(1)
}

// Variant operations

// CreateVariants mocks creating variants
func (m *MockRepository) CreateVariants(ctx context.Context, variants []*domain.ListingVariant) error {
	args := m.Called(ctx, variants)
	return args.Error(0)
}

// GetVariants mocks getting variants
func (m *MockRepository) GetVariants(ctx context.Context, listingID int64) ([]*domain.ListingVariant, error) {
	args := m.Called(ctx, listingID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.ListingVariant), args.Error(1)
}

// UpdateVariant mocks updating a variant
func (m *MockRepository) UpdateVariant(ctx context.Context, variant *domain.ListingVariant) error {
	args := m.Called(ctx, variant)
	return args.Error(0)
}

// DeleteVariant mocks deleting a variant
func (m *MockRepository) DeleteVariant(ctx context.Context, variantID int64) error {
	args := m.Called(ctx, variantID)
	return args.Error(0)
}

// Reindexing operations

// GetListingsForReindex mocks getting listings for reindexing
func (m *MockRepository) GetListingsForReindex(ctx context.Context, limit int) ([]*domain.Listing, error) {
	args := m.Called(ctx, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Listing), args.Error(1)
}

// ResetReindexFlags mocks resetting reindex flags
func (m *MockRepository) ResetReindexFlags(ctx context.Context, listingIDs []int64) error {
	args := m.Called(ctx, listingIDs)
	return args.Error(0)
}

// SyncDiscounts mocks syncing discounts
func (m *MockRepository) SyncDiscounts(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// Products operations

// GetProductByID mocks getting a product by ID
func (m *MockRepository) GetProductByID(ctx context.Context, productID int64, storefrontID *int64) (*domain.Product, error) {
	args := m.Called(ctx, productID, storefrontID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Product), args.Error(1)
}

// GetProductsBySKUs mocks getting products by SKUs
func (m *MockRepository) GetProductsBySKUs(ctx context.Context, skus []string, storefrontID *int64) ([]*domain.Product, error) {
	args := m.Called(ctx, skus, storefrontID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Product), args.Error(1)
}

// GetProductsByIDs mocks getting products by IDs
func (m *MockRepository) GetProductsByIDs(ctx context.Context, productIDs []int64, storefrontID *int64) ([]*domain.Product, error) {
	args := m.Called(ctx, productIDs, storefrontID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Product), args.Error(1)
}

// ListProducts mocks listing products
func (m *MockRepository) ListProducts(ctx context.Context, storefrontID int64, page, pageSize int, isActiveOnly bool) ([]*domain.Product, int, error) {
	args := m.Called(ctx, storefrontID, page, pageSize, isActiveOnly)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int), args.Error(2)
	}
	return args.Get(0).([]*domain.Product), args.Get(1).(int), args.Error(2)
}

// CreateProduct mocks creating a product
func (m *MockRepository) CreateProduct(ctx context.Context, input *domain.CreateProductInput) (*domain.Product, error) {
	args := m.Called(ctx, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Product), args.Error(1)
}

// BulkCreateProducts mocks bulk creating products
func (m *MockRepository) BulkCreateProducts(ctx context.Context, storefrontID int64, inputs []*domain.CreateProductInput) ([]*domain.Product, []domain.BulkProductError, error) {
	args := m.Called(ctx, storefrontID, inputs)
	var products []*domain.Product
	var errors []domain.BulkProductError

	if args.Get(0) != nil {
		products = args.Get(0).([]*domain.Product)
	}
	if args.Get(1) != nil {
		errors = args.Get(1).([]domain.BulkProductError)
	}

	return products, errors, args.Error(2)
}

// UpdateProduct mocks updating a product
func (m *MockRepository) UpdateProduct(ctx context.Context, productID int64, storefrontID int64, input *domain.UpdateProductInput) (*domain.Product, error) {
	args := m.Called(ctx, productID, storefrontID, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Product), args.Error(1)
}

// DeleteProduct mocks deleting a product
func (m *MockRepository) DeleteProduct(ctx context.Context, productID, storefrontID int64, hardDelete bool) (int32, error) {
	args := m.Called(ctx, productID, storefrontID, hardDelete)
	return args.Get(0).(int32), args.Error(1)
}

// BulkDeleteProducts mocks bulk deleting products
func (m *MockRepository) BulkDeleteProducts(ctx context.Context, storefrontID int64, productIDs []int64, hardDelete bool) (int32, int32, int32, map[int64]string, error) {
	args := m.Called(ctx, storefrontID, productIDs, hardDelete)

	var errors map[int64]string
	if args.Get(3) != nil {
		errors = args.Get(3).(map[int64]string)
	}

	return args.Get(0).(int32), args.Get(1).(int32), args.Get(2).(int32), errors, args.Error(4)
}

// BulkUpdateProducts mocks bulk updating products
func (m *MockRepository) BulkUpdateProducts(ctx context.Context, storefrontID int64, updates []*domain.BulkUpdateProductInput) (*domain.BulkUpdateProductsResult, error) {
	args := m.Called(ctx, storefrontID, updates)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.BulkUpdateProductsResult), args.Error(1)
}

// Product Variants operations

// GetVariantByID mocks getting a variant by ID
func (m *MockRepository) GetVariantByID(ctx context.Context, variantID int64, productID *int64) (*domain.ProductVariant, error) {
	args := m.Called(ctx, variantID, productID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ProductVariant), args.Error(1)
}

// GetVariantsByProductID mocks getting variants by product ID
func (m *MockRepository) GetVariantsByProductID(ctx context.Context, productID int64, isActiveOnly bool) ([]*domain.ProductVariant, error) {
	args := m.Called(ctx, productID, isActiveOnly)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.ProductVariant), args.Error(1)
}

// CreateProductVariant mocks creating a product variant
func (m *MockRepository) CreateProductVariant(ctx context.Context, input *domain.CreateVariantInput) (*domain.ProductVariant, error) {
	args := m.Called(ctx, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ProductVariant), args.Error(1)
}

// UpdateProductVariant mocks updating a product variant
func (m *MockRepository) UpdateProductVariant(ctx context.Context, variantID int64, productID int64, input *domain.UpdateVariantInput) (*domain.ProductVariant, error) {
	args := m.Called(ctx, variantID, productID, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ProductVariant), args.Error(1)
}

// DeleteProductVariant mocks deleting a product variant
func (m *MockRepository) DeleteProductVariant(ctx context.Context, variantID int64, productID int64) error {
	args := m.Called(ctx, variantID, productID)
	return args.Error(0)
}

// BulkCreateProductVariants mocks bulk creating product variants
func (m *MockRepository) BulkCreateProductVariants(ctx context.Context, productID int64, inputs []*domain.CreateVariantInput) ([]*domain.ProductVariant, error) {
	args := m.Called(ctx, productID, inputs)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.ProductVariant), args.Error(1)
}

// Inventory Management operations

// UpdateProductInventory mocks updating product inventory
func (m *MockRepository) UpdateProductInventory(ctx context.Context, storefrontID, productID, variantID int64, movementType string, quantity int32, reason, notes string, userID int64) (int32, int32, error) {
	args := m.Called(ctx, storefrontID, productID, variantID, movementType, quantity, reason, notes, userID)
	return args.Get(0).(int32), args.Get(1).(int32), args.Error(2)
}

// GetProductStats mocks getting product stats
func (m *MockRepository) GetProductStats(ctx context.Context, storefrontID int64) (*domain.ProductStats, error) {
	args := m.Called(ctx, storefrontID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ProductStats), args.Error(1)
}

// IncrementProductViews mocks incrementing product views
func (m *MockRepository) IncrementProductViews(ctx context.Context, productID int64) error {
	args := m.Called(ctx, productID)
	return args.Error(0)
}

// BatchUpdateStock mocks batch updating stock
func (m *MockRepository) BatchUpdateStock(ctx context.Context, storefrontID int64, items []domain.StockUpdateItem, reason string, userID int64) (int32, int32, []domain.StockUpdateResult, error) {
	args := m.Called(ctx, storefrontID, items, reason, userID)

	var results []domain.StockUpdateResult
	if args.Get(2) != nil {
		results = args.Get(2).([]domain.StockUpdateResult)
	}

	return args.Get(0).(int32), args.Get(1).(int32), results, args.Error(3)
}

// Transaction and database operations

// BeginTx mocks beginning a transaction
func (m *MockRepository) BeginTx(ctx context.Context) (*sql.Tx, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*sql.Tx), args.Error(1)
}

// GetDB mocks getting the database connection
func (m *MockRepository) GetDB() *sqlx.DB {
	args := m.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*sqlx.DB)
}
