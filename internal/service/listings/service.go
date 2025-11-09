package listings

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"

	"github.com/sveturs/listings/internal/domain"
)

// Repository defines the interface for listing data access
type Repository interface {
	CreateListing(ctx context.Context, input *domain.CreateListingInput) (*domain.Listing, error)
	GetListingByID(ctx context.Context, id int64) (*domain.Listing, error)
	GetListingByUUID(ctx context.Context, uuid string) (*domain.Listing, error)
	UpdateListing(ctx context.Context, id int64, input *domain.UpdateListingInput) (*domain.Listing, error)
	DeleteListing(ctx context.Context, id int64) error
	ListListings(ctx context.Context, filter *domain.ListListingsFilter) ([]*domain.Listing, int32, error)
	SearchListings(ctx context.Context, query *domain.SearchListingsQuery) ([]*domain.Listing, int32, error)
	EnqueueIndexing(ctx context.Context, listingID int64, operation string) error

	// Image operations
	GetImageByID(ctx context.Context, imageID int64) (*domain.ListingImage, error)
	DeleteImage(ctx context.Context, imageID int64) error
	AddImage(ctx context.Context, image *domain.ListingImage) (*domain.ListingImage, error)
	GetImages(ctx context.Context, listingID int64) ([]*domain.ListingImage, error)

	// Category operations
	GetRootCategories(ctx context.Context) ([]*domain.Category, error)
	GetAllCategories(ctx context.Context) ([]*domain.Category, error)
	GetPopularCategories(ctx context.Context, limit int) ([]*domain.Category, error)
	GetCategoryByID(ctx context.Context, categoryID int64) (*domain.Category, error)
	GetCategoryTree(ctx context.Context, categoryID int64) (*domain.CategoryTreeNode, error)

	// Favorites operations
	GetFavoritedUsers(ctx context.Context, listingID int64) ([]int64, error)
	AddToFavorites(ctx context.Context, userID, listingID int64) error
	RemoveFromFavorites(ctx context.Context, userID, listingID int64) error
	GetUserFavorites(ctx context.Context, userID int64) ([]int64, error)
	IsFavorite(ctx context.Context, userID, listingID int64) (bool, error)

	// Variant operations
	CreateVariants(ctx context.Context, variants []*domain.ListingVariant) error
	GetVariants(ctx context.Context, listingID int64) ([]*domain.ListingVariant, error)
	UpdateVariant(ctx context.Context, variant *domain.ListingVariant) error
	DeleteVariant(ctx context.Context, variantID int64) error

	// Reindexing operations
	GetListingsForReindex(ctx context.Context, limit int) ([]*domain.Listing, error)
	ResetReindexFlags(ctx context.Context, listingIDs []int64) error
	SyncDiscounts(ctx context.Context) error

	// Products operations
	GetProductByID(ctx context.Context, productID int64, storefrontID *int64) (*domain.Product, error)
	GetProductsBySKUs(ctx context.Context, skus []string, storefrontID *int64) ([]*domain.Product, error)
	GetProductsByIDs(ctx context.Context, productIDs []int64, storefrontID *int64) ([]*domain.Product, error)
	ListProducts(ctx context.Context, storefrontID int64, page, pageSize int, isActiveOnly bool) ([]*domain.Product, int, error)
	CreateProduct(ctx context.Context, input *domain.CreateProductInput) (*domain.Product, error)
	BulkCreateProducts(ctx context.Context, storefrontID int64, inputs []*domain.CreateProductInput) ([]*domain.Product, []domain.BulkProductError, error)
	UpdateProduct(ctx context.Context, productID int64, storefrontID int64, input *domain.UpdateProductInput) (*domain.Product, error)
	DeleteProduct(ctx context.Context, productID, storefrontID int64, hardDelete bool) (int32, error)
	BulkDeleteProducts(ctx context.Context, storefrontID int64, productIDs []int64, hardDelete bool) (int32, int32, int32, map[int64]string, error)
	BulkUpdateProducts(ctx context.Context, storefrontID int64, updates []*domain.BulkUpdateProductInput) (*domain.BulkUpdateProductsResult, error)

	// Product Variants operations
	GetVariantByID(ctx context.Context, variantID int64, productID *int64) (*domain.ProductVariant, error)
	GetVariantsByProductID(ctx context.Context, productID int64, isActiveOnly bool) ([]*domain.ProductVariant, error)
	CreateProductVariant(ctx context.Context, input *domain.CreateVariantInput) (*domain.ProductVariant, error)
	UpdateProductVariant(ctx context.Context, variantID int64, productID int64, input *domain.UpdateVariantInput) (*domain.ProductVariant, error)
	DeleteProductVariant(ctx context.Context, variantID int64, productID int64) error
	BulkCreateProductVariants(ctx context.Context, productID int64, inputs []*domain.CreateVariantInput) ([]*domain.ProductVariant, error)

	// Inventory Management operations
	UpdateProductInventory(ctx context.Context, storefrontID, productID, variantID int64, movementType string, quantity int32, reason, notes string, userID int64) (int32, int32, error)
	GetProductStats(ctx context.Context, storefrontID int64) (*domain.ProductStats, error)
	IncrementProductViews(ctx context.Context, productID int64) error
	BatchUpdateStock(ctx context.Context, storefrontID int64, items []domain.StockUpdateItem, reason string, userID int64) (int32, int32, []domain.StockUpdateResult, error)

	// Transaction and database operations
	BeginTx(ctx context.Context) (*sql.Tx, error)
	GetDB() *sqlx.DB
}

// CacheRepository defines the interface for caching operations
type CacheRepository interface {
	Get(ctx context.Context, key string, dest interface{}) error
	Set(ctx context.Context, key string, value interface{}) error
	Delete(ctx context.Context, key string) error
}

// IndexingService defines the interface for search indexing
type IndexingService interface {
	IndexListing(ctx context.Context, listing *domain.Listing) error
	UpdateListing(ctx context.Context, listing *domain.Listing) error
	DeleteListing(ctx context.Context, listingID int64) error
}

// Service implements business logic for listings
type Service struct {
	repo      Repository
	cache     CacheRepository
	indexer   IndexingService
	validator *validator.Validate
	logger    zerolog.Logger
}

// NewService creates a new listings service
func NewService(repo Repository, cache CacheRepository, indexer IndexingService, logger zerolog.Logger) *Service {
	return &Service{
		repo:      repo,
		cache:     cache,
		indexer:   indexer,
		validator: validator.New(),
		logger:    logger.With().Str("component", "listings_service").Logger(),
	}
}

// CreateListing creates a new listing with validation
func (s *Service) CreateListing(ctx context.Context, input *domain.CreateListingInput) (*domain.Listing, error) {
	// Validate input
	if err := s.validator.Struct(input); err != nil {
		s.logger.Warn().Err(err).Msg("invalid create listing input")
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Business validation
	if input.Price < 0 {
		return nil, fmt.Errorf("price cannot be negative")
	}

	if input.Quantity < 0 {
		return nil, fmt.Errorf("quantity cannot be negative")
	}

	// Create listing in database
	listing, err := s.repo.CreateListing(ctx, input)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to create listing in repository")
		return nil, fmt.Errorf("failed to create listing: %w", err)
	}

	// Enqueue for async indexing
	if err := s.repo.EnqueueIndexing(ctx, listing.ID, domain.IndexOpIndex); err != nil {
		s.logger.Warn().Err(err).Int64("listing_id", listing.ID).Msg("failed to enqueue indexing (non-critical)")
	}

	s.logger.Info().Int64("listing_id", listing.ID).Int64("user_id", listing.UserID).Msg("listing created successfully")
	return listing, nil
}

// GetListing retrieves a listing by ID with caching
func (s *Service) GetListing(ctx context.Context, id int64) (*domain.Listing, error) {
	// Try cache first (if available)
	if s.cache != nil {
		cacheKey := fmt.Sprintf("listing:%d", id)
		var cachedListing domain.Listing

		if err := s.cache.Get(ctx, cacheKey, &cachedListing); err == nil {
			s.logger.Debug().Int64("listing_id", id).Msg("listing found in cache")
			return &cachedListing, nil
		}
	}

	// Cache miss - fetch from database
	listing, err := s.repo.GetListingByID(ctx, id)
	if err != nil {
		s.logger.Error().Err(err).Int64("listing_id", id).Msg("failed to get listing")
		return nil, fmt.Errorf("failed to get listing: %w", err)
	}

	// Check for nil listing (defensive programming)
	if listing == nil {
		s.logger.Warn().Int64("listing_id", id).Msg("listing not found")
		return nil, fmt.Errorf("listing not found")
	}

	// Store in cache (non-blocking, ignore errors, if cache available)
	if s.cache != nil {
		go func() {
			cacheKey := fmt.Sprintf("listing:%d", id)
			if err := s.cache.Set(context.Background(), cacheKey, listing); err != nil {
				s.logger.Warn().Err(err).Int64("listing_id", id).Msg("failed to cache listing")
			}
		}()
	}

	return listing, nil
}

// GetListingByUUID retrieves a listing by UUID
func (s *Service) GetListingByUUID(ctx context.Context, uuid string) (*domain.Listing, error) {
	listing, err := s.repo.GetListingByUUID(ctx, uuid)
	if err != nil {
		s.logger.Error().Err(err).Str("uuid", uuid).Msg("failed to get listing by UUID")
		return nil, fmt.Errorf("failed to get listing: %w", err)
	}

	return listing, nil
}

// UpdateListing updates an existing listing with validation
func (s *Service) UpdateListing(ctx context.Context, id int64, userID int64, input *domain.UpdateListingInput) (*domain.Listing, error) {
	// Validate input
	if err := s.validator.Struct(input); err != nil {
		s.logger.Warn().Err(err).Msg("invalid update listing input")
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Check ownership
	existing, err := s.repo.GetListingByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("listing not found: %w", err)
	}

	if existing.UserID != userID {
		s.logger.Warn().Int64("listing_id", id).Int64("user_id", userID).Int64("owner_id", existing.UserID).Msg("unauthorized update attempt")
		return nil, fmt.Errorf("unauthorized: user does not own this listing")
	}

	// Business validation
	if input.Price != nil && *input.Price < 0 {
		return nil, fmt.Errorf("price cannot be negative")
	}

	if input.Quantity != nil && *input.Quantity < 0 {
		return nil, fmt.Errorf("quantity cannot be negative")
	}

	// Update listing
	listing, err := s.repo.UpdateListing(ctx, id, input)
	if err != nil {
		s.logger.Error().Err(err).Int64("listing_id", id).Msg("failed to update listing")
		return nil, fmt.Errorf("failed to update listing: %w", err)
	}

	// Invalidate cache (if available)
	if s.cache != nil {
		cacheKey := fmt.Sprintf("listing:%d", id)
		if err := s.cache.Delete(ctx, cacheKey); err != nil {
			s.logger.Warn().Err(err).Int64("listing_id", id).Msg("failed to invalidate cache")
		}
	}

	// Enqueue for async re-indexing
	if err := s.repo.EnqueueIndexing(ctx, listing.ID, domain.IndexOpUpdate); err != nil {
		s.logger.Warn().Err(err).Int64("listing_id", listing.ID).Msg("failed to enqueue indexing (non-critical)")
	}

	s.logger.Info().Int64("listing_id", listing.ID).Msg("listing updated successfully")
	return listing, nil
}

// DeleteListing soft-deletes a listing with ownership check
func (s *Service) DeleteListing(ctx context.Context, id int64, userID int64) error {
	// Check ownership
	existing, err := s.repo.GetListingByID(ctx, id)
	if err != nil {
		return fmt.Errorf("listing not found: %w", err)
	}

	if existing.UserID != userID {
		s.logger.Warn().Int64("listing_id", id).Int64("user_id", userID).Int64("owner_id", existing.UserID).Msg("unauthorized delete attempt")
		return fmt.Errorf("unauthorized: user does not own this listing")
	}

	// Delete listing
	if err := s.repo.DeleteListing(ctx, id); err != nil {
		s.logger.Error().Err(err).Int64("listing_id", id).Msg("failed to delete listing")
		return fmt.Errorf("failed to delete listing: %w", err)
	}

	// Invalidate cache (if available)
	if s.cache != nil {
		cacheKey := fmt.Sprintf("listing:%d", id)
		if err := s.cache.Delete(ctx, cacheKey); err != nil {
			s.logger.Warn().Err(err).Int64("listing_id", id).Msg("failed to invalidate cache")
		}
	}

	// Enqueue for async index deletion
	if err := s.repo.EnqueueIndexing(ctx, id, domain.IndexOpDelete); err != nil {
		s.logger.Warn().Err(err).Int64("listing_id", id).Msg("failed to enqueue indexing deletion (non-critical)")
	}

	s.logger.Info().Int64("listing_id", id).Msg("listing deleted successfully")
	return nil
}

// ListListings returns a filtered list of listings
func (s *Service) ListListings(ctx context.Context, filter *domain.ListListingsFilter) ([]*domain.Listing, int32, error) {
	// Validate filter
	if err := s.validator.Struct(filter); err != nil {
		s.logger.Warn().Err(err).Msg("invalid list filter")
		return nil, 0, fmt.Errorf("validation failed: %w", err)
	}

	// Apply business rules
	if filter.Limit > 100 {
		filter.Limit = 100 // Max 100 items per page
	}

	if filter.Limit <= 0 {
		filter.Limit = 20 // Default to 20
	}

	listings, total, err := s.repo.ListListings(ctx, filter)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to list listings")
		return nil, 0, fmt.Errorf("failed to list listings: %w", err)
	}

	s.logger.Debug().Int("count", len(listings)).Int32("total", total).Msg("listings retrieved")
	return listings, total, nil
}

// SearchListings performs full-text search on listings
func (s *Service) SearchListings(ctx context.Context, query *domain.SearchListingsQuery) ([]*domain.Listing, int32, error) {
	// Validate query
	if err := s.validator.Struct(query); err != nil {
		s.logger.Warn().Err(err).Msg("invalid search query")
		return nil, 0, fmt.Errorf("validation failed: %w", err)
	}

	// Apply business rules
	if query.Limit > 100 {
		query.Limit = 100
	}

	if query.Limit <= 0 {
		query.Limit = 20
	}

	if len(query.Query) < 2 {
		return nil, 0, fmt.Errorf("search query must be at least 2 characters")
	}

	// Try cache for search results (if cache is available)
	categoryID := int64(0)
	if query.CategoryID != nil {
		categoryID = *query.CategoryID
	}
	cacheKey := fmt.Sprintf("search:%s:%d:%d:%d", query.Query, categoryID, query.Limit, query.Offset)
	var cachedResults []*domain.Listing
	var cachedTotal int32

	if s.cache != nil {
		if err := s.cache.Get(ctx, cacheKey, &cachedResults); err == nil {
			s.logger.Debug().Str("query", query.Query).Msg("search results found in cache")
			return cachedResults, cachedTotal, nil
		}
	}

	// Cache miss - search in database
	listings, total, err := s.repo.SearchListings(ctx, query)
	if err != nil {
		s.logger.Error().Err(err).Str("query", query.Query).Msg("failed to search listings")
		return nil, 0, fmt.Errorf("failed to search listings: %w", err)
	}

	// Cache search results (non-blocking, if cache is available)
	if s.cache != nil {
		go func() {
			if err := s.cache.Set(context.Background(), cacheKey, listings); err != nil {
				s.logger.Warn().Err(err).Str("query", query.Query).Msg("failed to cache search results")
			}
		}()
	}

	s.logger.Debug().Str("query", query.Query).Int("count", len(listings)).Int32("total", total).Msg("search completed")
	return listings, total, nil
}

// Admin operations (no ownership check)

// AdminGetListing retrieves any listing (admin operation)
func (s *Service) AdminGetListing(ctx context.Context, id int64) (*domain.Listing, error) {
	listing, err := s.repo.GetListingByID(ctx, id)
	if err != nil {
		s.logger.Error().Err(err).Int64("listing_id", id).Msg("admin failed to get listing")
		return nil, fmt.Errorf("failed to get listing: %w", err)
	}

	s.logger.Info().Int64("listing_id", id).Msg("admin retrieved listing")
	return listing, nil
}

// AdminUpdateListing updates any listing (admin operation)
func (s *Service) AdminUpdateListing(ctx context.Context, id int64, input *domain.UpdateListingInput) (*domain.Listing, error) {
	// Validate input
	if err := s.validator.Struct(input); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	listing, err := s.repo.UpdateListing(ctx, id, input)
	if err != nil {
		s.logger.Error().Err(err).Int64("listing_id", id).Msg("admin failed to update listing")
		return nil, fmt.Errorf("failed to update listing: %w", err)
	}

	// Invalidate cache (if available)
	if s.cache != nil {
		cacheKey := fmt.Sprintf("listing:%d", id)
		if err := s.cache.Delete(ctx, cacheKey); err != nil {
			s.logger.Warn().Err(err).Int64("listing_id", id).Msg("failed to invalidate cache")
		}
	}

	// Enqueue for async re-indexing
	if err := s.repo.EnqueueIndexing(ctx, listing.ID, domain.IndexOpUpdate); err != nil {
		s.logger.Warn().Err(err).Int64("listing_id", listing.ID).Msg("failed to enqueue indexing")
	}

	s.logger.Info().Int64("listing_id", listing.ID).Msg("admin updated listing")
	return listing, nil
}

// AdminDeleteListing deletes any listing (admin operation)
func (s *Service) AdminDeleteListing(ctx context.Context, id int64) error {
	if err := s.repo.DeleteListing(ctx, id); err != nil {
		s.logger.Error().Err(err).Int64("listing_id", id).Msg("admin failed to delete listing")
		return fmt.Errorf("failed to delete listing: %w", err)
	}

	// Invalidate cache (if available)
	if s.cache != nil {
		cacheKey := fmt.Sprintf("listing:%d", id)
		if err := s.cache.Delete(ctx, cacheKey); err != nil {
			s.logger.Warn().Err(err).Int64("listing_id", id).Msg("failed to invalidate cache")
		}
	}

	// Enqueue for async index deletion
	if err := s.repo.EnqueueIndexing(ctx, id, domain.IndexOpDelete); err != nil {
		s.logger.Warn().Err(err).Int64("listing_id", id).Msg("failed to enqueue indexing deletion")
	}

	s.logger.Info().Int64("listing_id", id).Msg("admin deleted listing")
	return nil
}

// Image operations

func (s *Service) GetImageByID(ctx context.Context, imageID int64) (*domain.ListingImage, error) {
	return s.repo.GetImageByID(ctx, imageID)
}

func (s *Service) DeleteImage(ctx context.Context, imageID int64) error {
	return s.repo.DeleteImage(ctx, imageID)
}

func (s *Service) AddImage(ctx context.Context, image *domain.ListingImage) (*domain.ListingImage, error) {
	return s.repo.AddImage(ctx, image)
}

func (s *Service) GetImages(ctx context.Context, listingID int64) ([]*domain.ListingImage, error) {
	return s.repo.GetImages(ctx, listingID)
}

// Category operations

func (s *Service) GetRootCategories(ctx context.Context) ([]*domain.Category, error) {
	return s.repo.GetRootCategories(ctx)
}

func (s *Service) GetAllCategories(ctx context.Context) ([]*domain.Category, error) {
	return s.repo.GetAllCategories(ctx)
}

func (s *Service) GetPopularCategories(ctx context.Context, limit int) ([]*domain.Category, error) {
	return s.repo.GetPopularCategories(ctx, limit)
}

func (s *Service) GetCategoryByID(ctx context.Context, categoryID int64) (*domain.Category, error) {
	return s.repo.GetCategoryByID(ctx, categoryID)
}

func (s *Service) GetCategoryTree(ctx context.Context, categoryID int64) (*domain.CategoryTreeNode, error) {
	return s.repo.GetCategoryTree(ctx, categoryID)
}

// Favorites operations

func (s *Service) GetFavoritedUsers(ctx context.Context, listingID int64) ([]int64, error) {
	return s.repo.GetFavoritedUsers(ctx, listingID)
}

func (s *Service) AddToFavorites(ctx context.Context, userID, listingID int64) error {
	// Validate IDs
	if userID <= 0 {
		return fmt.Errorf("invalid user ID")
	}
	if listingID <= 0 {
		return fmt.Errorf("invalid listing ID")
	}

	// Check if listing exists
	_, err := s.repo.GetListingByID(ctx, listingID)
	if err != nil {
		s.logger.Error().Err(err).Int64("listing_id", listingID).Msg("listing not found")
		return fmt.Errorf("listing not found: %w", err)
	}

	// Add to favorites
	err = s.repo.AddToFavorites(ctx, userID, listingID)
	if err != nil {
		s.logger.Error().Err(err).Int64("user_id", userID).Int64("listing_id", listingID).Msg("failed to add to favorites")
		return fmt.Errorf("failed to add to favorites: %w", err)
	}

	s.logger.Info().Int64("user_id", userID).Int64("listing_id", listingID).Msg("added to favorites")
	return nil
}

func (s *Service) RemoveFromFavorites(ctx context.Context, userID, listingID int64) error {
	// Validate IDs
	if userID <= 0 {
		return fmt.Errorf("invalid user ID")
	}
	if listingID <= 0 {
		return fmt.Errorf("invalid listing ID")
	}

	// Remove from favorites
	err := s.repo.RemoveFromFavorites(ctx, userID, listingID)
	if err != nil {
		s.logger.Error().Err(err).Int64("user_id", userID).Int64("listing_id", listingID).Msg("failed to remove from favorites")
		return fmt.Errorf("failed to remove from favorites: %w", err)
	}

	s.logger.Info().Int64("user_id", userID).Int64("listing_id", listingID).Msg("removed from favorites")
	return nil
}

func (s *Service) GetUserFavorites(ctx context.Context, userID int64) ([]int64, error) {
	// Validate user ID
	if userID <= 0 {
		return nil, fmt.Errorf("invalid user ID")
	}

	listingIDs, err := s.repo.GetUserFavorites(ctx, userID)
	if err != nil {
		s.logger.Error().Err(err).Int64("user_id", userID).Msg("failed to get user favorites")
		return nil, fmt.Errorf("failed to get user favorites: %w", err)
	}

	s.logger.Debug().Int64("user_id", userID).Int("count", len(listingIDs)).Msg("user favorites retrieved")
	return listingIDs, nil
}

func (s *Service) IsFavorite(ctx context.Context, userID, listingID int64) (bool, error) {
	// Validate IDs
	if userID <= 0 {
		return false, fmt.Errorf("invalid user ID")
	}
	if listingID <= 0 {
		return false, fmt.Errorf("invalid listing ID")
	}

	isFavorite, err := s.repo.IsFavorite(ctx, userID, listingID)
	if err != nil {
		s.logger.Error().Err(err).Int64("user_id", userID).Int64("listing_id", listingID).Msg("failed to check favorite status")
		return false, fmt.Errorf("failed to check favorite status: %w", err)
	}

	return isFavorite, nil
}

// Variant operations

func (s *Service) CreateVariants(ctx context.Context, variants []*domain.ListingVariant) error {
	return s.repo.CreateVariants(ctx, variants)
}

func (s *Service) GetVariants(ctx context.Context, listingID int64) ([]*domain.ListingVariant, error) {
	return s.repo.GetVariants(ctx, listingID)
}

func (s *Service) GetVariantByID(ctx context.Context, variantID int64) (*domain.ListingVariant, error) {
	// Get all variants for listing and find the one we need
	// This is a simple implementation - in production you'd want a dedicated query
	variants, err := s.repo.GetVariants(ctx, 0) // Pass 0 to get all or implement proper method
	if err != nil {
		return nil, err
	}

	for _, v := range variants {
		if v.ID == variantID {
			return v, nil
		}
	}

	return nil, fmt.Errorf("variant not found")
}

func (s *Service) UpdateVariant(ctx context.Context, variant *domain.ListingVariant) error {
	return s.repo.UpdateVariant(ctx, variant)
}

func (s *Service) DeleteVariant(ctx context.Context, variantID int64) error {
	return s.repo.DeleteVariant(ctx, variantID)
}

// Reindexing operations

func (s *Service) GetListingsForReindex(ctx context.Context, limit int) ([]*domain.Listing, error) {
	return s.repo.GetListingsForReindex(ctx, limit)
}

func (s *Service) ResetReindexFlags(ctx context.Context, listingIDs []int64) error {
	return s.repo.ResetReindexFlags(ctx, listingIDs)
}

func (s *Service) SyncDiscounts(ctx context.Context) error {
	return s.repo.SyncDiscounts(ctx)
}

// Products operations

func (s *Service) GetProduct(ctx context.Context, productID int64, storefrontID *int64) (*domain.Product, error) {
	s.logger.Debug().Int64("product_id", productID).Interface("storefront_id", storefrontID).Msg("getting product")

	product, err := s.repo.GetProductByID(ctx, productID, storefrontID)
	if err != nil {
		s.logger.Error().Err(err).Int64("product_id", productID).Msg("failed to get product")
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	return product, nil
}

func (s *Service) GetProductsBySKUs(ctx context.Context, skus []string, storefrontID *int64) ([]*domain.Product, error) {
	if len(skus) == 0 {
		return []*domain.Product{}, nil
	}

	s.logger.Debug().Int("sku_count", len(skus)).Interface("storefront_id", storefrontID).Msg("getting products by SKUs")

	products, err := s.repo.GetProductsBySKUs(ctx, skus, storefrontID)
	if err != nil {
		s.logger.Error().Err(err).Int("sku_count", len(skus)).Msg("failed to get products by SKUs")
		return nil, fmt.Errorf("failed to get products by SKUs: %w", err)
	}

	return products, nil
}

func (s *Service) GetProductsByIDs(ctx context.Context, productIDs []int64, storefrontID *int64) ([]*domain.Product, error) {
	if len(productIDs) == 0 {
		return []*domain.Product{}, nil
	}

	s.logger.Debug().Int("id_count", len(productIDs)).Interface("storefront_id", storefrontID).Msg("getting products by IDs")

	products, err := s.repo.GetProductsByIDs(ctx, productIDs, storefrontID)
	if err != nil {
		s.logger.Error().Err(err).Int("id_count", len(productIDs)).Msg("failed to get products by IDs")
		return nil, fmt.Errorf("failed to get products by IDs: %w", err)
	}

	return products, nil
}

func (s *Service) ListProducts(ctx context.Context, storefrontID int64, page, pageSize int, isActiveOnly bool) ([]*domain.Product, int, error) {
	// Validate pagination parameters
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20 // default
	}

	s.logger.Debug().Int64("storefront_id", storefrontID).Int("page", page).Int("page_size", pageSize).Bool("is_active_only", isActiveOnly).Msg("listing products")

	products, total, err := s.repo.ListProducts(ctx, storefrontID, page, pageSize, isActiveOnly)
	if err != nil {
		s.logger.Error().Err(err).Int64("storefront_id", storefrontID).Msg("failed to list products")
		return nil, 0, fmt.Errorf("failed to list products: %w", err)
	}

	s.logger.Debug().Int("count", len(products)).Int("total", total).Msg("products listed successfully")
	return products, total, nil
}

func (s *Service) CreateProduct(ctx context.Context, input *domain.CreateProductInput) (*domain.Product, error) {
	s.logger.Debug().
		Int64("storefront_id", input.StorefrontID).
		Str("name", input.Name).
		Msg("creating product")

	// Validate input using validator
	if err := s.validator.Struct(input); err != nil {
		s.logger.Error().Err(err).Msg("product validation failed")
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Create product in repository
	product, err := s.repo.CreateProduct(ctx, input)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to create product")
		return nil, err // Return as-is to preserve error placeholders
	}

	s.logger.Info().Int64("product_id", product.ID).Msg("product created successfully")
	return product, nil
}

// BulkCreateProducts creates multiple products in a single atomic transaction
func (s *Service) BulkCreateProducts(ctx context.Context, storefrontID int64, inputs []*domain.CreateProductInput) ([]*domain.Product, []domain.BulkProductError, error) {
	s.logger.Debug().
		Int64("storefront_id", storefrontID).
		Int("product_count", len(inputs)).
		Msg("bulk creating products")

	// Validate storefront_id
	if storefrontID <= 0 {
		return nil, nil, fmt.Errorf("storefront_id must be greater than 0")
	}

	// Validate batch size
	if len(inputs) == 0 {
		return nil, nil, fmt.Errorf("products.bulk_empty")
	}
	if len(inputs) > 1000 {
		return nil, nil, fmt.Errorf("products.bulk_too_large")
	}

	// Validate each input
	for i, input := range inputs {
		if input == nil {
			return nil, nil, fmt.Errorf("product at index %d is nil", i)
		}
		// Ensure storefront_id matches
		input.StorefrontID = storefrontID

		// Validate using validator
		if err := s.validator.Struct(input); err != nil {
			s.logger.Error().Err(err).Int("index", i).Msg("product validation failed")
			return nil, nil, fmt.Errorf("validation failed for product at index %d: %w", i, err)
		}
	}

	// Create products in repository
	products, errors, err := s.repo.BulkCreateProducts(ctx, storefrontID, inputs)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to bulk create products")
		return products, errors, err // Return partial results if available
	}

	s.logger.Info().
		Int("successful", len(products)).
		Int("failed", len(errors)).
		Msg("bulk product creation completed")

	return products, errors, nil
}

// UpdateProduct updates an existing product with ownership validation
func (s *Service) UpdateProduct(ctx context.Context, productID int64, storefrontID int64, input *domain.UpdateProductInput) (*domain.Product, error) {
	s.logger.Debug().
		Int64("product_id", productID).
		Int64("storefront_id", storefrontID).
		Msg("updating product")

	// Validate product_id and storefront_id
	if productID <= 0 {
		return nil, fmt.Errorf("product_id must be greater than 0")
	}
	if storefrontID <= 0 {
		return nil, fmt.Errorf("storefront_id must be greater than 0")
	}

	// Validate input using validator
	if err := s.validator.Struct(input); err != nil {
		s.logger.Error().Err(err).Msg("product update validation failed")
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Update product in repository
	product, err := s.repo.UpdateProduct(ctx, productID, storefrontID, input)
	if err != nil {
		s.logger.Error().Err(err).Int64("product_id", productID).Msg("failed to update product")
		return nil, err // Return as-is to preserve error placeholders
	}

	s.logger.Info().Int64("product_id", product.ID).Msg("product updated successfully")
	return product, nil
}

// BulkUpdateProducts updates multiple products in a single atomic operation
func (s *Service) BulkUpdateProducts(ctx context.Context, storefrontID int64, updates []*domain.BulkUpdateProductInput) (*domain.BulkUpdateProductsResult, error) {
	s.logger.Debug().
		Int64("storefront_id", storefrontID).
		Int("update_count", len(updates)).
		Msg("bulk updating products via service")

	// Validate inputs
	if storefrontID <= 0 {
		return nil, fmt.Errorf("storefront_id must be greater than 0")
	}

	if len(updates) == 0 {
		return &domain.BulkUpdateProductsResult{
			SuccessfulProducts: []*domain.Product{},
			FailedUpdates:      []domain.BulkUpdateError{},
		}, nil
	}

	if len(updates) > 1000 {
		return nil, fmt.Errorf("products.bulk_update_limit_exceeded")
	}

	// Validate each update input
	for _, update := range updates {
		if update.ProductID <= 0 {
			return nil, fmt.Errorf("all product_ids must be greater than 0")
		}

		// Validate using validator
		if err := s.validator.Struct(update); err != nil {
			s.logger.Error().Err(err).Int64("product_id", update.ProductID).Msg("bulk update validation failed")
			return nil, fmt.Errorf("validation failed for product %d: %w", update.ProductID, err)
		}

		// Check that at least one field is set for update
		hasUpdate := update.Name != nil || update.Description != nil || update.Price != nil ||
			update.SKU != nil || update.Barcode != nil || update.StockQuantity != nil ||
			update.StockStatus != nil || update.IsActive != nil || update.Attributes != nil ||
			update.HasIndividualLocation != nil || update.IndividualAddress != nil ||
			update.IndividualLatitude != nil || update.IndividualLongitude != nil ||
			update.LocationPrivacy != nil || update.ShowOnMap != nil

		if !hasUpdate {
			return nil, fmt.Errorf("at least one field must be specified for update for product %d", update.ProductID)
		}
	}

	// Call repository bulk update
	result, err := s.repo.BulkUpdateProducts(ctx, storefrontID, updates)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to bulk update products")
		return nil, err // Return as-is to preserve error placeholders
	}

	s.logger.Info().
		Int("successful_count", len(result.SuccessfulProducts)).
		Int("failed_count", len(result.FailedUpdates)).
		Msg("bulk update products completed successfully")

	return result, nil
}

// DeleteProduct deletes a product (soft or hard) with ownership validation
func (s *Service) DeleteProduct(ctx context.Context, productID, storefrontID int64, hardDelete bool) (int32, error) {
	s.logger.Debug().
		Int64("product_id", productID).
		Int64("storefront_id", storefrontID).
		Bool("hard_delete", hardDelete).
		Msg("deleting product")

	// Validate inputs
	if productID <= 0 {
		return 0, fmt.Errorf("product_id must be greater than 0")
	}
	if storefrontID <= 0 {
		return 0, fmt.Errorf("storefront_id must be greater than 0")
	}

	// Delete product in repository (includes ownership check)
	variantsDeleted, err := s.repo.DeleteProduct(ctx, productID, storefrontID, hardDelete)
	if err != nil {
		s.logger.Error().Err(err).Int64("product_id", productID).Msg("failed to delete product")
		return 0, err // Return as-is to preserve error placeholders
	}

	deleteType := "soft deleted"
	if hardDelete {
		deleteType = "hard deleted"
	}

	s.logger.Info().
		Int64("product_id", productID).
		Int32("variants_deleted", variantsDeleted).
		Str("delete_type", deleteType).
		Msg("product deleted successfully")

	return variantsDeleted, nil
}

// BulkDeleteProducts deletes multiple products in a single atomic operation
func (s *Service) BulkDeleteProducts(ctx context.Context, storefrontID int64, productIDs []int64, hardDelete bool) (int32, int32, int32, map[int64]string, error) {
	s.logger.Debug().
		Int64("storefront_id", storefrontID).
		Int("product_count", len(productIDs)).
		Bool("hard_delete", hardDelete).
		Msg("bulk deleting products via service")

	// Validate inputs
	if storefrontID <= 0 {
		return 0, 0, 0, nil, fmt.Errorf("storefront_id must be greater than 0")
	}

	if len(productIDs) == 0 {
		return 0, 0, 0, nil, fmt.Errorf("product_ids list cannot be empty")
	}

	if len(productIDs) > 1000 {
		return 0, 0, 0, nil, fmt.Errorf("cannot delete more than 1000 products at once")
	}

	// Deduplicate product IDs
	uniqueIDs := make(map[int64]bool)
	deduplicatedIDs := make([]int64, 0, len(productIDs))
	for _, id := range productIDs {
		if id > 0 && !uniqueIDs[id] {
			uniqueIDs[id] = true
			deduplicatedIDs = append(deduplicatedIDs, id)
		}
	}

	if len(deduplicatedIDs) == 0 {
		return 0, 0, 0, nil, fmt.Errorf("no valid product IDs provided")
	}

	// Call repository method
	successCount, failedCount, variantsDeleted, errors, err := s.repo.BulkDeleteProducts(ctx, storefrontID, deduplicatedIDs, hardDelete)
	if err != nil {
		s.logger.Error().Err(err).Int64("storefront_id", storefrontID).Msg("failed to bulk delete products")
		return 0, 0, 0, nil, err // Return as-is to preserve error details
	}

	deleteType := "soft deleted"
	if hardDelete {
		deleteType = "hard deleted"
	}

	s.logger.Info().
		Int32("success_count", successCount).
		Int32("failed_count", failedCount).
		Int32("variants_deleted", variantsDeleted).
		Str("delete_type", deleteType).
		Msg("bulk delete products completed successfully")

	return successCount, failedCount, variantsDeleted, errors, nil
}

// Product Variants operations

func (s *Service) GetVariant(ctx context.Context, variantID int64, productID *int64) (*domain.ProductVariant, error) {
	s.logger.Debug().Int64("variant_id", variantID).Interface("product_id", productID).Msg("getting product variant")

	variant, err := s.repo.GetVariantByID(ctx, variantID, productID)
	if err != nil {
		s.logger.Error().Err(err).Int64("variant_id", variantID).Msg("failed to get variant")
		return nil, fmt.Errorf("failed to get variant: %w", err)
	}

	return variant, nil
}

func (s *Service) GetVariantsByProductID(ctx context.Context, productID int64, isActiveOnly bool) ([]*domain.ProductVariant, error) {
	s.logger.Debug().Int64("product_id", productID).Bool("is_active_only", isActiveOnly).Msg("getting variants by product ID")

	variants, err := s.repo.GetVariantsByProductID(ctx, productID, isActiveOnly)
	if err != nil {
		s.logger.Error().Err(err).Int64("product_id", productID).Msg("failed to get variants by product ID")
		return nil, fmt.Errorf("failed to get variants by product ID: %w", err)
	}

	s.logger.Debug().Int("count", len(variants)).Msg("variants retrieved successfully")
	return variants, nil
}

// CreateProductVariant creates a new product variant with validation
func (s *Service) CreateProductVariant(ctx context.Context, input *domain.CreateVariantInput) (*domain.ProductVariant, error) {
	s.logger.Debug().
		Int64("product_id", input.ProductID).
		Msg("creating product variant")

	// Validate input using validator
	if err := s.validator.Struct(input); err != nil {
		s.logger.Error().Err(err).Msg("variant validation failed")
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Create variant in repository
	variant, err := s.repo.CreateProductVariant(ctx, input)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to create variant")
		return nil, err // Return as-is to preserve error placeholders
	}

	s.logger.Info().Int64("variant_id", variant.ID).Msg("product variant created successfully")
	return variant, nil
}

// UpdateProductVariant updates an existing product variant with validation
func (s *Service) UpdateProductVariant(ctx context.Context, variantID int64, productID int64, input *domain.UpdateVariantInput) (*domain.ProductVariant, error) {
	s.logger.Debug().
		Int64("variant_id", variantID).
		Int64("product_id", productID).
		Msg("updating product variant")

	// Validate input using validator
	if err := s.validator.Struct(input); err != nil {
		s.logger.Error().Err(err).Msg("variant validation failed")
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Update variant in repository
	variant, err := s.repo.UpdateProductVariant(ctx, variantID, productID, input)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to update variant")
		return nil, err // Return as-is to preserve error placeholders
	}

	s.logger.Info().Int64("variant_id", variant.ID).Msg("product variant updated successfully")
	return variant, nil
}

// DeleteProductVariant deletes a product variant with business rules enforcement
func (s *Service) DeleteProductVariant(ctx context.Context, variantID int64, productID int64) error {
	s.logger.Debug().
		Int64("variant_id", variantID).
		Int64("product_id", productID).
		Msg("deleting product variant")

	// Delete variant in repository (business rules enforced there)
	err := s.repo.DeleteProductVariant(ctx, variantID, productID)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to delete variant")
		return err // Return as-is to preserve error placeholders
	}

	s.logger.Info().Int64("variant_id", variantID).Msg("product variant deleted successfully")
	return nil
}

// BulkCreateProductVariants creates multiple product variants with validation
func (s *Service) BulkCreateProductVariants(ctx context.Context, productID int64, inputs []*domain.CreateVariantInput) ([]*domain.ProductVariant, error) {
	s.logger.Debug().
		Int64("product_id", productID).
		Int("count", len(inputs)).
		Msg("bulk creating product variants")

	// Validate product ID
	if productID <= 0 {
		return nil, fmt.Errorf("product_id must be greater than 0")
	}

	// Validate inputs count
	if len(inputs) == 0 {
		return nil, fmt.Errorf("variants list cannot be empty")
	}

	if len(inputs) > 1000 {
		return nil, fmt.Errorf("cannot create more than 1000 variants at once")
	}

	// Validate each input using validator
	for i, input := range inputs {
		if err := s.validator.Struct(input); err != nil {
			s.logger.Error().Err(err).Int("index", i).Msg("variant validation failed")
			return nil, fmt.Errorf("validation failed for variant at index %d: %w", i, err)
		}
	}

	// Business rule: Only one variant can be default
	defaultCount := 0
	for _, input := range inputs {
		if input.IsDefault {
			defaultCount++
		}
	}

	if defaultCount > 1 {
		return nil, fmt.Errorf("only one variant can be set as default")
	}

	// Create variants in repository (transaction handled there)
	variants, err := s.repo.BulkCreateProductVariants(ctx, productID, inputs)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to bulk create variants")
		return nil, err // Return as-is to preserve error placeholders
	}

	s.logger.Info().
		Int("count", len(variants)).
		Int64("product_id", productID).
		Msg("product variants bulk created successfully")

	return variants, nil
}

// UpdateProductInventory updates product inventory with movement tracking
func (s *Service) UpdateProductInventory(ctx context.Context, storefrontID, productID, variantID int64, movementType string, quantity int32, reason, notes string, userID int64) (int32, int32, error) {
	s.logger.Debug().
		Int64("storefront_id", storefrontID).
		Int64("product_id", productID).
		Int64("variant_id", variantID).
		Str("movement_type", movementType).
		Int32("quantity", quantity).
		Msg("updating product inventory")

	// Validate inputs
	if storefrontID <= 0 {
		return 0, 0, fmt.Errorf("storefront_id must be greater than 0")
	}
	if productID <= 0 {
		return 0, 0, fmt.Errorf("product_id must be greater than 0")
	}
	if movementType != "in" && movementType != "out" && movementType != "adjustment" {
		return 0, 0, fmt.Errorf("invalid movement_type: must be 'in', 'out', or 'adjustment'")
	}
	if quantity < 0 {
		return 0, 0, fmt.Errorf("quantity cannot be negative")
	}
	if userID <= 0 {
		return 0, 0, fmt.Errorf("user_id must be greater than 0")
	}

	// Call repository to update inventory
	stockBefore, stockAfter, err := s.repo.UpdateProductInventory(ctx, storefrontID, productID, variantID, movementType, quantity, reason, notes, userID)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to update product inventory")
		return 0, 0, err
	}

	s.logger.Info().
		Int64("product_id", productID).
		Int64("variant_id", variantID).
		Int32("stock_before", stockBefore).
		Int32("stock_after", stockAfter).
		Msg("product inventory updated successfully")

	return stockBefore, stockAfter, nil
}

// GetProductStats retrieves product statistics for a storefront
func (s *Service) GetProductStats(ctx context.Context, storefrontID int64) (*domain.ProductStats, error) {
	s.logger.Debug().Int64("storefront_id", storefrontID).Msg("getting product stats")

	// Validate storefront ID
	if storefrontID <= 0 {
		return nil, fmt.Errorf("storefront_id must be greater than 0")
	}

	// Fetch stats from repository
	stats, err := s.repo.GetProductStats(ctx, storefrontID)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to get product stats")
		return nil, err
	}

	s.logger.Info().
		Int32("total_products", stats.TotalProducts).
		Int32("active_products", stats.ActiveProducts).
		Float64("total_value", stats.TotalValue).
		Msg("product stats retrieved successfully")

	return stats, nil
}

// IncrementProductViews increments the view counter for a product
func (s *Service) IncrementProductViews(ctx context.Context, productID int64) error {
	s.logger.Debug().Int64("product_id", productID).Msg("incrementing product views")

	// Validate product ID
	if productID <= 0 {
		return fmt.Errorf("product_id must be greater than 0")
	}

	// Call repository to increment views
	if err := s.repo.IncrementProductViews(ctx, productID); err != nil {
		s.logger.Error().Err(err).Msg("failed to increment product views")
		return err
	}

	s.logger.Debug().Int64("product_id", productID).Msg("product views incremented successfully")
	return nil
}

// BatchUpdateStock updates stock for multiple products/variants with validation
func (s *Service) BatchUpdateStock(ctx context.Context, storefrontID int64, items []domain.StockUpdateItem, reason string, userID int64) (int32, int32, []domain.StockUpdateResult, error) {
	s.logger.Debug().
		Int64("storefront_id", storefrontID).
		Int("item_count", len(items)).
		Msg("batch updating stock")

	// Validate inputs
	if storefrontID <= 0 {
		return 0, 0, nil, fmt.Errorf("storefront_id must be greater than 0")
	}
	if len(items) == 0 {
		return 0, 0, nil, fmt.Errorf("items list cannot be empty")
	}
	if len(items) > 1000 {
		return 0, 0, nil, fmt.Errorf("cannot update more than 1000 items at once")
	}
	if userID <= 0 {
		return 0, 0, nil, fmt.Errorf("user_id must be greater than 0")
	}

	// Validate each item
	for i, item := range items {
		if item.ProductID <= 0 {
			return 0, 0, nil, fmt.Errorf("invalid product_id at index %d", i)
		}
		if item.Quantity < 0 {
			return 0, 0, nil, fmt.Errorf("invalid quantity at index %d: cannot be negative", i)
		}
	}

	// Call repository to perform batch update
	successCount, failedCount, results, err := s.repo.BatchUpdateStock(ctx, storefrontID, items, reason, userID)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to batch update stock")
		return 0, 0, nil, err
	}

	s.logger.Info().
		Int32("success_count", successCount).
		Int32("failed_count", failedCount).
		Msg("batch stock update completed")

	return successCount, failedCount, results, nil
}
