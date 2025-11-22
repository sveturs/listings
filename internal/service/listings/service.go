package listings

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"

	"github.com/sveturs/listings/internal/domain"
	"github.com/sveturs/listings/internal/repository/postgres"
)

// Repository defines the interface for listing data access
type Repository interface {
	CreateListing(ctx context.Context, input *domain.CreateListingInput) (*domain.Listing, error)
	GetListingByID(ctx context.Context, id int64) (*domain.Listing, error)
	GetListingByUUID(ctx context.Context, uuid string) (*domain.Listing, error)
	GetListingBySlug(ctx context.Context, slug string) (*domain.Listing, error)
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
	ReorderImages(ctx context.Context, listingID int64, orders []postgres.ImageOrder) error

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

	// Variant operations (old ListingVariant - deprecated)
	CreateVariants(ctx context.Context, variants []*domain.ListingVariant) error
	GetVariants(ctx context.Context, listingID int64) ([]*domain.ListingVariant, error)
	UpdateVariant(ctx context.Context, variant *domain.ListingVariant) error
	DeleteVariant(ctx context.Context, variantID int64) error

	// B2C Variant operations (new domain.Variant for b2c_product_variants table)
	// Repository methods for postgres implementation
	CreateVariant(ctx context.Context, variant *domain.Variant) (*domain.Variant, error)
	GetVariant(ctx context.Context, id int64) (*domain.Variant, error)
	UpdateB2CVariant(ctx context.Context, id int64, update *domain.VariantUpdate) (*domain.Variant, error)
	DeleteB2CVariant(ctx context.Context, id int64) error
	ListVariants(ctx context.Context, filters *domain.VariantFilters) ([]*domain.Variant, error)

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

	// Product Images operations (B2C)
	GetProductImageByID(ctx context.Context, imageID int64) (*domain.ProductImage, error)
	AddProductImage(ctx context.Context, image *domain.ProductImage) (*domain.ProductImage, error)
	GetProductImages(ctx context.Context, productID int64) ([]*domain.ProductImage, error)
	DeleteProductImage(ctx context.Context, imageID int64) error
	ReorderProductImages(ctx context.Context, productID int64, orders []postgres.ProductImageOrder) error

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
	GetSimilarListings(ctx context.Context, listingID int64, limit int32) ([]*domain.Listing, int32, error)
}

// Service implements business logic for listings
type Service struct {
	repo          Repository
	cache         CacheRepository
	indexer       IndexingService
	validator     *Validator
	slugGenerator *SlugGenerator
	stdValidator  *validator.Validate
	logger        zerolog.Logger
}

// NewService creates a new listings service
func NewService(repo Repository, cache CacheRepository, indexer IndexingService, logger zerolog.Logger) *Service {
	return &Service{
		repo:          repo,
		cache:         cache,
		indexer:       indexer,
		validator:     NewValidator(repo),
		slugGenerator: NewSlugGenerator(repo),
		stdValidator:  validator.New(),
		logger:        logger.With().Str("component", "listings_service").Logger(),
	}
}

// CreateListing creates a new listing with validation
func (s *Service) CreateListing(ctx context.Context, input *domain.CreateListingInput) (*domain.Listing, error) {
	// 1. Comprehensive validation using custom validator
	if err := s.validator.ValidateCreateInput(ctx, input); err != nil {
		s.logger.Warn().Err(err).Msg("validation failed for create listing")
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// 2. Generate unique slug from title
	slug, err := s.slugGenerator.Generate(ctx, input.Title)
	if err != nil {
		s.logger.Error().Err(err).Str("title", input.Title).Msg("failed to generate slug")
		return nil, fmt.Errorf("failed to generate slug: %w", err)
	}

	// 3. Create listing entity
	listing := &domain.Listing{
		UserID:       input.UserID,
		StorefrontID: input.StorefrontID,
		Title:        input.Title,
		Description:  input.Description,
		Price:        input.Price,
		Currency:     input.Currency,
		CategoryID:   input.CategoryID,
		Slug:         slug,
		Quantity:     input.Quantity,
		SKU:          input.SKU,
		SourceType:   input.SourceType,
		Status:       domain.StatusDraft, // Default status
		Visibility:   domain.VisibilityPublic,
	}

	// 4. Set expiration date for C2C listings (30 days)
	if input.SourceType == domain.SourceTypeC2C {
		expiresAt := time.Now().AddDate(0, 0, 30) // 30 days from now
		listing.ExpiresAt = &expiresAt
	}

	// 5. Create in database (using CreateListingInput)
	created, err := s.repo.CreateListing(ctx, input)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to create listing in repository")
		return nil, fmt.Errorf("failed to create listing: %w", err)
	}

	// Note: We need to set the slug on the created listing since repo doesn't know about it yet
	// This is a temporary workaround until we update the repository layer
	created.Slug = slug
	if listing.ExpiresAt != nil {
		created.ExpiresAt = listing.ExpiresAt
	}

	// 6. Enqueue for async indexing
	if err := s.repo.EnqueueIndexing(ctx, created.ID, domain.IndexOpIndex); err != nil {
		s.logger.Warn().Err(err).Int64("listing_id", created.ID).Msg("failed to enqueue indexing (non-critical)")
	}

	s.logger.Info().
		Int64("listing_id", created.ID).
		Int64("user_id", created.UserID).
		Str("slug", slug).
		Str("source_type", created.SourceType).
		Msg("listing created successfully")

	return created, nil
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

	// Load images
	images, err := s.repo.GetImages(ctx, id)
	if err != nil {
		// Log error but don't fail the request
		s.logger.Warn().Err(err).Int64("listing_id", id).Msg("failed to load images")
	} else {
		listing.Images = images
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

	// Load images
	images, err := s.repo.GetImages(ctx, listing.ID)
	if err != nil {
		// Log error but don't fail the request
		s.logger.Warn().Err(err).Int64("listing_id", listing.ID).Msg("failed to load images")
	} else {
		listing.Images = images
	}

	return listing, nil
}

// UpdateListing updates an existing listing with validation
func (s *Service) UpdateListing(ctx context.Context, id int64, userID int64, input *domain.UpdateListingInput) (*domain.Listing, error) {
	// 1. Validate input using custom validator
	if err := s.validator.ValidateUpdateInput(input); err != nil {
		s.logger.Warn().Err(err).Msg("validation failed for update listing")
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// 2. Get existing listing
	existing, err := s.repo.GetListingByID(ctx, id)
	if err != nil {
		s.logger.Error().Err(err).Int64("listing_id", id).Msg("failed to get existing listing")
		return nil, fmt.Errorf("listing not found: %w", err)
	}

	if existing == nil {
		return nil, fmt.Errorf("listing not found")
	}

	// 3. Verify ownership
	if existing.UserID != userID {
		s.logger.Warn().
			Int64("listing_id", id).
			Int64("user_id", userID).
			Int64("owner_id", existing.UserID).
			Msg("unauthorized update attempt")
		return nil, fmt.Errorf("unauthorized: user does not own this listing")
	}

	// 4. Validate status transition if status is being updated
	if input.Status != nil {
		if err := s.validator.ValidateStatusTransition(existing.Status, *input.Status); err != nil {
			s.logger.Warn().
				Err(err).
				Str("from_status", existing.Status).
				Str("to_status", *input.Status).
				Msg("invalid status transition")
			return nil, fmt.Errorf("invalid status transition: %w", err)
		}
	}

	// 5. Update listing in database
	updated, err := s.repo.UpdateListing(ctx, id, input)
	if err != nil {
		s.logger.Error().Err(err).Int64("listing_id", id).Msg("failed to update listing")
		return nil, fmt.Errorf("failed to update listing: %w", err)
	}

	// 6. Load images
	images, err := s.repo.GetImages(ctx, id)
	if err != nil {
		// Log error but don't fail the request
		s.logger.Warn().Err(err).Int64("listing_id", id).Msg("failed to load images")
	} else {
		updated.Images = images
	}

	// 7. Invalidate cache (if available)
	if s.cache != nil {
		cacheKey := fmt.Sprintf("listing:%d", id)
		if err := s.cache.Delete(ctx, cacheKey); err != nil {
			s.logger.Warn().Err(err).Int64("listing_id", id).Msg("failed to invalidate cache")
		}
	}

	// 8. Enqueue for async re-indexing
	if err := s.repo.EnqueueIndexing(ctx, updated.ID, domain.IndexOpUpdate); err != nil {
		s.logger.Warn().Err(err).Int64("listing_id", updated.ID).Msg("failed to enqueue indexing (non-critical)")
	}

	s.logger.Info().
		Int64("listing_id", updated.ID).
		Int64("user_id", userID).
		Msg("listing updated successfully")

	return updated, nil
}

// DeleteListing soft-deletes a listing with ownership check and cascade cleanup
func (s *Service) DeleteListing(ctx context.Context, id int64, userID int64) error {
	// 1. Get existing listing and verify ownership
	existing, err := s.repo.GetListingByID(ctx, id)
	if err != nil {
		s.logger.Error().Err(err).Int64("listing_id", id).Msg("failed to get listing for deletion")
		return fmt.Errorf("listing not found: %w", err)
	}

	if existing == nil {
		return fmt.Errorf("listing not found")
	}

	if existing.UserID != userID {
		s.logger.Warn().
			Int64("listing_id", id).
			Int64("user_id", userID).
			Int64("owner_id", existing.UserID).
			Msg("unauthorized delete attempt")
		return fmt.Errorf("unauthorized: user does not own this listing")
	}

	// 2. Soft delete the listing
	if err := s.repo.DeleteListing(ctx, id); err != nil {
		s.logger.Error().Err(err).Int64("listing_id", id).Msg("failed to delete listing")
		return fmt.Errorf("failed to delete listing: %w", err)
	}

	// 3. Cascade delete: Remove associated images
	// Note: This is best-effort - we log errors but don't fail the operation
	images, err := s.repo.GetImages(ctx, id)
	if err == nil && len(images) > 0 {
		for _, img := range images {
			if err := s.repo.DeleteImage(ctx, img.ID); err != nil {
				s.logger.Warn().
					Err(err).
					Int64("listing_id", id).
					Int64("image_id", img.ID).
					Msg("failed to delete associated image")
			}
		}
		s.logger.Info().
			Int64("listing_id", id).
			Int("images_count", len(images)).
			Msg("cascade deleted associated images")
	}

	// 4. Invalidate all related caches
	if s.cache != nil {
		// Invalidate listing cache
		cacheKey := fmt.Sprintf("listing:%d", id)
		if err := s.cache.Delete(ctx, cacheKey); err != nil {
			s.logger.Warn().Err(err).Int64("listing_id", id).Msg("failed to invalidate listing cache")
		}

		// Invalidate favorites count cache
		favCountKey := fmt.Sprintf("favorites:listing:%d:count", id)
		if err := s.cache.Delete(ctx, favCountKey); err != nil {
			s.logger.Warn().Err(err).Int64("listing_id", id).Msg("failed to invalidate favorites count cache")
		}

		// Invalidate user's listings cache (if we cached user listings)
		userListingsKey := fmt.Sprintf("user:%d:listings", existing.UserID)
		if err := s.cache.Delete(ctx, userListingsKey); err != nil {
			s.logger.Warn().Err(err).Int64("user_id", existing.UserID).Msg("failed to invalidate user listings cache")
		}
	}

	// 5. Enqueue for async index deletion
	if err := s.repo.EnqueueIndexing(ctx, id, domain.IndexOpDelete); err != nil {
		s.logger.Warn().Err(err).Int64("listing_id", id).Msg("failed to enqueue indexing deletion (non-critical)")
	}

	s.logger.Info().
		Int64("listing_id", id).
		Int64("user_id", userID).
		Msg("listing deleted successfully with cascade cleanup")

	return nil
}

// ListListings returns a filtered list of listings
func (s *Service) ListListings(ctx context.Context, filter *domain.ListListingsFilter) ([]*domain.Listing, int32, error) {
	// Validate filter
	if err := s.stdValidator.Struct(filter); err != nil {
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

	// Load images for each listing (eager loading)
	for _, listing := range listings {
		images, err := s.repo.GetImages(ctx, listing.ID)
		if err != nil {
			// Log error but don't fail the request
			s.logger.Warn().Err(err).Int64("listing_id", listing.ID).Msg("failed to load images for listing")
		} else {
			listing.Images = images
		}
	}

	s.logger.Debug().Int("count", len(listings)).Int32("total", total).Msg("listings retrieved")
	return listings, total, nil
}

// SearchListings performs full-text search on listings
func (s *Service) SearchListings(ctx context.Context, query *domain.SearchListingsQuery) ([]*domain.Listing, int32, error) {
	// Validate query
	if err := s.stdValidator.Struct(query); err != nil {
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

	// Load images for each listing (eager loading)
	for _, listing := range listings {
		images, err := s.repo.GetImages(ctx, listing.ID)
		if err != nil {
			// Log error but don't fail the request
			s.logger.Warn().Err(err).Int64("listing_id", listing.ID).Msg("failed to load images for listing")
		} else {
			listing.Images = images
		}
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

// GetSimilarListings finds listings similar to the given listing
func (s *Service) GetSimilarListings(ctx context.Context, listingID int64, limit int32) ([]*domain.Listing, int32, error) {
	// Apply business rules for limit
	if limit <= 0 {
		limit = 10 // Default to 10 similar items
	}
	if limit > 20 {
		limit = 20 // Cap at 20 items max
	}

	// Try cache first (1-hour TTL for similar listings)
	cacheKey := fmt.Sprintf("similar:%d:%d", listingID, limit)

	if s.cache != nil {
		var cached struct {
			Listings []*domain.Listing
			Total    int32
		}
		if err := s.cache.Get(ctx, cacheKey, &cached); err == nil {
			s.logger.Debug().Int64("listing_id", listingID).Msg("similar listings found in cache")
			return cached.Listings, cached.Total, nil
		}
	}

	// Cache miss - query OpenSearch via indexer
	listings, total, err := s.indexer.GetSimilarListings(ctx, listingID, limit)
	if err != nil {
		s.logger.Error().Err(err).Int64("listing_id", listingID).Msg("failed to get similar listings")
		// Graceful degradation: return empty array instead of error
		return []*domain.Listing{}, 0, nil
	}

	// Load images for each similar listing (eager loading)
	for _, listing := range listings {
		images, err := s.repo.GetImages(ctx, listing.ID)
		if err != nil {
			// Log error but don't fail the request
			s.logger.Warn().Err(err).Int64("listing_id", listing.ID).Msg("failed to load images for similar listing")
		} else {
			listing.Images = images
		}
	}

	// Cache results (non-blocking, 1-hour TTL)
	if s.cache != nil {
		go func() {
			cached := struct {
				Listings []*domain.Listing
				Total    int32
			}{
				Listings: listings,
				Total:    total,
			}
			if err := s.cache.Set(context.Background(), cacheKey, cached); err != nil {
				s.logger.Warn().Err(err).Int64("listing_id", listingID).Msg("failed to cache similar listings")
			}
		}()
	}

	s.logger.Debug().
		Int64("listing_id", listingID).
		Int("count", len(listings)).
		Int32("total", total).
		Msg("similar listings search completed")

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

	// Load images
	images, err := s.repo.GetImages(ctx, id)
	if err != nil {
		// Log error but don't fail the request
		s.logger.Warn().Err(err).Int64("listing_id", id).Msg("failed to load images")
	} else {
		listing.Images = images
	}

	s.logger.Info().Int64("listing_id", id).Msg("admin retrieved listing")
	return listing, nil
}

// AdminUpdateListing updates any listing (admin operation)
func (s *Service) AdminUpdateListing(ctx context.Context, id int64, input *domain.UpdateListingInput) (*domain.Listing, error) {
	// Validate input
	if err := s.stdValidator.Struct(input); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	listing, err := s.repo.UpdateListing(ctx, id, input)
	if err != nil {
		s.logger.Error().Err(err).Int64("listing_id", id).Msg("admin failed to update listing")
		return nil, fmt.Errorf("failed to update listing: %w", err)
	}

	// Load images
	images, err := s.repo.GetImages(ctx, id)
	if err != nil {
		// Log error but don't fail the request
		s.logger.Warn().Err(err).Int64("listing_id", id).Msg("failed to load images")
	} else {
		listing.Images = images
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
	// Get image first to know which listing to invalidate cache for
	image, err := s.repo.GetImageByID(ctx, imageID)
	if err != nil {
		return err
	}

	listingID := image.ListingID

	// Delete image
	err = s.repo.DeleteImage(ctx, imageID)
	if err != nil {
		return err
	}

	// Invalidate cache for this listing (image was deleted)
	if s.cache != nil {
		cacheKey := fmt.Sprintf("listing:%d", listingID)
		if err := s.cache.Delete(ctx, cacheKey); err != nil {
			// Log but don't fail - cache invalidation failure is not critical
			s.logger.Warn().Err(err).Int64("listing_id", listingID).Msg("failed to invalidate listing cache after deleting image")
		} else {
			s.logger.Debug().Int64("listing_id", listingID).Msg("listing cache invalidated after deleting image")
		}
	}

	// Reindex in OpenSearch (async, non-blocking)
	if s.indexer != nil {
		go func() {
			defer func() {
				if r := recover(); r != nil {
					s.logger.Warn().Interface("panic", r).Int64("listing_id", listingID).
						Msg("recovered from panic during OpenSearch reindexing after image delete")
				}
			}()

			indexCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			// Get full listing with remaining images
			listing, err := s.repo.GetListingByID(indexCtx, listingID)
			if err != nil {
				s.logger.Warn().Err(err).Int64("listing_id", listingID).
					Msg("failed to get listing for reindexing after image delete (non-blocking)")
				return
			}

			// Load remaining images
			images, err := s.repo.GetImages(indexCtx, listingID)
			if err == nil {
				listing.Images = images
			}

			// Update in OpenSearch
			if err := s.indexer.UpdateListing(indexCtx, listing); err != nil {
				s.logger.Warn().Err(err).Int64("listing_id", listingID).
					Msg("OpenSearch reindexing failed after image delete (non-blocking)")
			} else {
				s.logger.Debug().Int64("listing_id", listingID).Msg("listing reindexed in OpenSearch after image delete")
			}
		}()
	}

	return nil
}

func (s *Service) AddImage(ctx context.Context, image *domain.ListingImage) (*domain.ListingImage, error) {
	result, err := s.repo.AddImage(ctx, image)
	if err != nil {
		return nil, err
	}

	// Invalidate cache for this listing (images were added)
	if s.cache != nil {
		cacheKey := fmt.Sprintf("listing:%d", image.ListingID)
		if err := s.cache.Delete(ctx, cacheKey); err != nil {
			// Log but don't fail - cache invalidation failure is not critical
			s.logger.Warn().Err(err).Int64("listing_id", image.ListingID).Msg("failed to invalidate listing cache after adding image")
		} else {
			s.logger.Debug().Int64("listing_id", image.ListingID).Msg("listing cache invalidated after adding image")
		}
	}

	// Reindex in OpenSearch (async, non-blocking)
	if s.indexer != nil {
		go func() {
			defer func() {
				if r := recover(); r != nil {
					s.logger.Warn().Interface("panic", r).Int64("listing_id", image.ListingID).
						Msg("recovered from panic during OpenSearch reindexing after image add")
				}
			}()

			indexCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			// Get full listing with all images
			listing, err := s.repo.GetListingByID(indexCtx, image.ListingID)
			if err != nil {
				s.logger.Warn().Err(err).Int64("listing_id", image.ListingID).
					Msg("failed to get listing for reindexing after image add (non-blocking)")
				return
			}

			// Load all images
			images, err := s.repo.GetImages(indexCtx, image.ListingID)
			if err == nil {
				listing.Images = images
			}

			// Update in OpenSearch
			if err := s.indexer.UpdateListing(indexCtx, listing); err != nil {
				s.logger.Warn().Err(err).Int64("listing_id", image.ListingID).
					Msg("OpenSearch reindexing failed after image add (non-blocking)")
			} else {
				s.logger.Debug().Int64("listing_id", image.ListingID).Msg("listing reindexed in OpenSearch after image add")
			}
		}()
	}

	return result, nil
}

func (s *Service) GetImages(ctx context.Context, listingID int64) ([]*domain.ListingImage, error) {
	return s.repo.GetImages(ctx, listingID)
}

func (s *Service) ReorderImages(ctx context.Context, listingID int64, orders []postgres.ImageOrder) error {
	// Validate listing_id
	if listingID <= 0 {
		return fmt.Errorf("invalid listing_id: %d", listingID)
	}

	// Validate orders
	if len(orders) == 0 {
		return fmt.Errorf("no image orders provided")
	}

	// Validate each order
	for i, order := range orders {
		if order.ImageID <= 0 {
			return fmt.Errorf("invalid image_id at index %d: %d", i, order.ImageID)
		}
		if order.DisplayOrder < 0 {
			return fmt.Errorf("invalid display_order at index %d: %d", i, order.DisplayOrder)
		}
	}

	s.logger.Debug().
		Int64("listing_id", listingID).
		Int("count", len(orders)).
		Msg("reordering images")

	err := s.repo.ReorderImages(ctx, listingID, orders)
	if err != nil {
		return err
	}

	// Invalidate cache for this listing (images were reordered)
	if s.cache != nil {
		cacheKey := fmt.Sprintf("listing:%d", listingID)
		if err := s.cache.Delete(ctx, cacheKey); err != nil {
			// Log but don't fail - cache invalidation failure is not critical
			s.logger.Warn().Err(err).Int64("listing_id", listingID).Msg("failed to invalidate listing cache after reordering images")
		} else {
			s.logger.Debug().Int64("listing_id", listingID).Msg("listing cache invalidated after reordering images")
		}
	}

	// Reindex in OpenSearch (async, non-blocking)
	if s.indexer != nil {
		go func() {
			defer func() {
				if r := recover(); r != nil {
					s.logger.Warn().Interface("panic", r).Int64("listing_id", listingID).
						Msg("recovered from panic during OpenSearch reindexing after image reorder")
				}
			}()

			indexCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			// Get full listing with reordered images
			listing, err := s.repo.GetListingByID(indexCtx, listingID)
			if err != nil {
				s.logger.Warn().Err(err).Int64("listing_id", listingID).
					Msg("failed to get listing for reindexing after image reorder (non-blocking)")
				return
			}

			// Load images in new order
			images, err := s.repo.GetImages(indexCtx, listingID)
			if err == nil {
				listing.Images = images
			}

			// Update in OpenSearch
			if err := s.indexer.UpdateListing(indexCtx, listing); err != nil {
				s.logger.Warn().Err(err).Int64("listing_id", listingID).
					Msg("OpenSearch reindexing failed after image reorder (non-blocking)")
			} else {
				s.logger.Debug().Int64("listing_id", listingID).Msg("listing reindexed in OpenSearch after image reorder")
			}
		}()
	}

	return nil
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

	// Invalidate related caches
	if s.cache != nil {
		// Invalidate user's favorites list cache
		userFavKey := fmt.Sprintf("favorites:user:%d", userID)
		if err := s.cache.Delete(ctx, userFavKey); err != nil {
			s.logger.Warn().Err(err).Int64("user_id", userID).Msg("failed to invalidate user favorites cache")
		}

		// Invalidate favorites count cache for this listing
		countKey := fmt.Sprintf("favorites:listing:%d:count", listingID)
		if err := s.cache.Delete(ctx, countKey); err != nil {
			s.logger.Warn().Err(err).Int64("listing_id", listingID).Msg("failed to invalidate favorites count cache")
		}

		// Invalidate IsFavorite cache for this user-listing pair
		isFavKey := fmt.Sprintf("favorites:user:%d:listing:%d", userID, listingID)
		if err := s.cache.Delete(ctx, isFavKey); err != nil {
			s.logger.Warn().Err(err).Int64("user_id", userID).Int64("listing_id", listingID).Msg("failed to invalidate IsFavorite cache")
		}
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

	// Invalidate related caches
	if s.cache != nil {
		// Invalidate user's favorites list cache
		userFavKey := fmt.Sprintf("favorites:user:%d", userID)
		if err := s.cache.Delete(ctx, userFavKey); err != nil {
			s.logger.Warn().Err(err).Int64("user_id", userID).Msg("failed to invalidate user favorites cache")
		}

		// Invalidate favorites count cache for this listing
		countKey := fmt.Sprintf("favorites:listing:%d:count", listingID)
		if err := s.cache.Delete(ctx, countKey); err != nil {
			s.logger.Warn().Err(err).Int64("listing_id", listingID).Msg("failed to invalidate favorites count cache")
		}

		// Invalidate IsFavorite cache for this user-listing pair
		isFavKey := fmt.Sprintf("favorites:user:%d:listing:%d", userID, listingID)
		if err := s.cache.Delete(ctx, isFavKey); err != nil {
			s.logger.Warn().Err(err).Int64("user_id", userID).Int64("listing_id", listingID).Msg("failed to invalidate IsFavorite cache")
		}
	}

	s.logger.Info().Int64("user_id", userID).Int64("listing_id", listingID).Msg("removed from favorites")
	return nil
}

func (s *Service) GetUserFavorites(ctx context.Context, userID int64) ([]int64, error) {
	// Validate user ID
	if userID <= 0 {
		return nil, fmt.Errorf("invalid user ID")
	}

	cacheKey := fmt.Sprintf("favorites:user:%d", userID)

	// 1. Try to get from cache
	if s.cache != nil {
		var cachedIDs []int64
		if err := s.cache.Get(ctx, cacheKey, &cachedIDs); err == nil {
			s.logger.Debug().Int64("user_id", userID).Int("count", len(cachedIDs)).Msg("user favorites from cache")
			return cachedIDs, nil
		}
	}

	// 2. Cache miss - get from repository
	listingIDs, err := s.repo.GetUserFavorites(ctx, userID)
	if err != nil {
		s.logger.Error().Err(err).Int64("user_id", userID).Msg("failed to get user favorites")
		return nil, fmt.Errorf("failed to get user favorites: %w", err)
	}

	// 3. Cache the result for 5 minutes
	if s.cache != nil {
		if err := s.cache.Set(ctx, cacheKey, listingIDs); err != nil {
			s.logger.Warn().Err(err).Int64("user_id", userID).Msg("failed to cache user favorites")
			// Don't fail the request if caching fails
		}
	}

	s.logger.Debug().Int64("user_id", userID).Int("count", len(listingIDs)).Msg("user favorites retrieved from database")
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

	cacheKey := fmt.Sprintf("favorites:user:%d:listing:%d", userID, listingID)

	// 1. Try to get from cache
	if s.cache != nil {
		var cachedStatus bool
		if err := s.cache.Get(ctx, cacheKey, &cachedStatus); err == nil {
			s.logger.Debug().Int64("user_id", userID).Int64("listing_id", listingID).Bool("is_favorite", cachedStatus).Msg("favorite status from cache")
			return cachedStatus, nil
		}
	}

	// 2. Cache miss - get from repository
	isFavorite, err := s.repo.IsFavorite(ctx, userID, listingID)
	if err != nil {
		s.logger.Error().Err(err).Int64("user_id", userID).Int64("listing_id", listingID).Msg("failed to check favorite status")
		return false, fmt.Errorf("failed to check favorite status: %w", err)
	}

	// 3. Cache the result for 5 minutes
	if s.cache != nil {
		if err := s.cache.Set(ctx, cacheKey, isFavorite); err != nil {
			s.logger.Warn().Err(err).Int64("user_id", userID).Int64("listing_id", listingID).Msg("failed to cache favorite status")
			// Don't fail the request if caching fails
		}
	}

	s.logger.Debug().Int64("user_id", userID).Int64("listing_id", listingID).Bool("is_favorite", isFavorite).Msg("favorite status from database")
	return isFavorite, nil
}

// GetFavoritesCount returns the number of times a listing has been favorited with caching
func (s *Service) GetFavoritesCount(ctx context.Context, listingID int64) (int64, error) {
	// Validate listing ID
	if listingID <= 0 {
		return 0, fmt.Errorf("invalid listing ID")
	}

	cacheKey := fmt.Sprintf("favorites:listing:%d:count", listingID)

	// 1. Try to get from cache
	if s.cache != nil {
		var cachedCount int64
		if err := s.cache.Get(ctx, cacheKey, &cachedCount); err == nil {
			s.logger.Debug().Int64("listing_id", listingID).Int64("count", cachedCount).Msg("favorites count from cache")
			return cachedCount, nil
		}
	}

	// 2. Cache miss - get favorited users from repository
	users, err := s.repo.GetFavoritedUsers(ctx, listingID)
	if err != nil {
		s.logger.Error().Err(err).Int64("listing_id", listingID).Msg("failed to get favorited users")
		return 0, fmt.Errorf("failed to get favorites count: %w", err)
	}

	count := int64(len(users))

	// 3. Cache the result for 5 minutes
	if s.cache != nil {
		if err := s.cache.Set(ctx, cacheKey, count); err != nil {
			s.logger.Warn().Err(err).Int64("listing_id", listingID).Msg("failed to cache favorites count")
			// Don't fail the request if caching fails
		}
	}

	s.logger.Debug().Int64("listing_id", listingID).Int64("count", count).Msg("favorites count from database")
	return count, nil
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
	listings, err := s.repo.GetListingsForReindex(ctx, limit)
	if err != nil {
		return nil, err
	}

	// Load images for each listing (eager loading)
	for _, listing := range listings {
		images, err := s.repo.GetImages(ctx, listing.ID)
		if err != nil {
			// Log error but don't fail the request
			s.logger.Warn().Err(err).Int64("listing_id", listing.ID).Msg("failed to load images for listing")
		} else {
			listing.Images = images
		}
	}

	return listings, nil
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
	if err := s.stdValidator.Struct(input); err != nil {
		s.logger.Error().Err(err).Msg("product validation failed")
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Create product in repository
	product, err := s.repo.CreateProduct(ctx, input)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to create product")
		return nil, err // Return as-is to preserve error placeholders
	}

	// Async indexing in OpenSearch (non-blocking, graceful degradation)
	if s.indexer != nil {
		go func() {
			defer func() {
				if r := recover(); r != nil {
					s.logger.Warn().Interface("panic", r).Int64("product_id", product.ID).
						Msg("recovered from panic during OpenSearch indexing")
				}
			}()

			indexCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			// Get full product with images for indexing
			fullProduct, err := s.repo.GetProductByID(indexCtx, product.ID, nil)
			if err != nil {
				s.logger.Warn().Err(err).Int64("product_id", product.ID).
					Msg("failed to get product for indexing (non-blocking)")
				return
			}

			// Get images
			// Note: Images are loaded via repository's GetProductByID if it supports it
			// Or we can load them separately if needed
			// For now, assume fullProduct may have images loaded

			// Convert domain.Product to domain.Listing for indexing
			listing := convertProductToListing(fullProduct)

			if err := s.indexer.IndexListing(indexCtx, listing); err != nil {
				s.logger.Warn().Err(err).Int64("product_id", product.ID).
					Msg("OpenSearch indexing failed (non-blocking)")
			} else {
				s.logger.Debug().Int64("product_id", product.ID).Msg("product indexed in OpenSearch")
			}
		}()
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

	// Basic nil check and ensure storefront_id matches
	for i, input := range inputs {
		if input == nil {
			return nil, nil, fmt.Errorf("product at index %d is nil", i)
		}
		// Ensure storefront_id matches
		input.StorefrontID = storefrontID
	}

	// Create products in repository (detailed validation happens gracefully there)
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
	if err := s.stdValidator.Struct(input); err != nil {
		s.logger.Error().Err(err).Msg("product update validation failed")
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Update product in repository
	product, err := s.repo.UpdateProduct(ctx, productID, storefrontID, input)
	if err != nil {
		s.logger.Error().Err(err).Int64("product_id", productID).Msg("failed to update product")
		return nil, err // Return as-is to preserve error placeholders
	}

	// Async re-indexing in OpenSearch (non-blocking, graceful degradation)
	if s.indexer != nil {
		go func() {
			defer func() {
				if r := recover(); r != nil {
					s.logger.Warn().Interface("panic", r).Int64("product_id", product.ID).
						Msg("recovered from panic during OpenSearch re-indexing")
				}
			}()

			indexCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			// Get full product with images for indexing
			fullProduct, err := s.repo.GetProductByID(indexCtx, product.ID, nil)
			if err != nil {
				s.logger.Warn().Err(err).Int64("product_id", product.ID).
					Msg("failed to get product for re-indexing (non-blocking)")
				return
			}

			// Convert domain.Product to domain.Listing for indexing
			listing := convertProductToListing(fullProduct)

			if err := s.indexer.UpdateListing(indexCtx, listing); err != nil {
				s.logger.Warn().Err(err).Int64("product_id", product.ID).
					Msg("OpenSearch re-indexing failed (non-blocking)")
			} else {
				s.logger.Debug().Int64("product_id", product.ID).Msg("product re-indexed in OpenSearch")
			}
		}()
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
		if err := s.stdValidator.Struct(update); err != nil {
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

	// Async deletion from OpenSearch (non-blocking, graceful degradation)
	if s.indexer != nil {
		go func() {
			defer func() {
				if r := recover(); r != nil {
					s.logger.Warn().Interface("panic", r).Int64("product_id", productID).
						Msg("recovered from panic during OpenSearch deletion")
				}
			}()

			indexCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			if err := s.indexer.DeleteListing(indexCtx, productID); err != nil {
				s.logger.Warn().Err(err).Int64("product_id", productID).
					Msg("OpenSearch deletion failed (non-blocking)")
			} else {
				s.logger.Debug().Int64("product_id", productID).Msg("product removed from OpenSearch index")
			}
		}()
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

	// Async deletion from OpenSearch (non-blocking, graceful degradation)
	if s.indexer != nil && successCount > 0 {
		go func(productIDs []int64) {
			defer func() {
				if r := recover(); r != nil {
					s.logger.Warn().Interface("panic", r).
						Msg("recovered from panic during bulk OpenSearch deletion")
				}
			}()

			indexCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			deletedCount := 0
			for _, productID := range productIDs {
				if err := s.indexer.DeleteListing(indexCtx, productID); err != nil {
					s.logger.Warn().Err(err).Int64("product_id", productID).
						Msg("OpenSearch deletion failed (non-blocking)")
				} else {
					deletedCount++
				}
			}
			s.logger.Info().Int("deleted_from_index", deletedCount).Int("total", len(productIDs)).
				Msg("bulk OpenSearch deletion completed")
		}(deduplicatedIDs)
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
	if err := s.stdValidator.Struct(input); err != nil {
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
	if err := s.stdValidator.Struct(input); err != nil {
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
		if err := s.stdValidator.Struct(input); err != nil {
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

// ============================================================================
// Helper Functions for OpenSearch Integration
// ============================================================================

// ReindexAll performs full reindexing of all products to OpenSearch
// This is an administrative operation used for rebuilding search index
func (s *Service) ReindexAll(ctx context.Context, sourceType string, batchSize int) (int32, int32, int, []string, error) {
	s.logger.Info().
		Str("source_type", sourceType).
		Int("batch_size", batchSize).
		Msg("starting full reindexing")

	startTime := time.Now()

	// Validate inputs
	if s.indexer == nil {
		return 0, 0, 0, nil, fmt.Errorf("opensearch_not_configured")
	}

	if batchSize <= 0 || batchSize > 10000 {
		batchSize = 1000 // Default batch size
	}

	// Validate source_type filter
	if sourceType != "" && sourceType != "b2c" && sourceType != "c2c" {
		return 0, 0, 0, nil, fmt.Errorf("invalid source_type: must be 'b2c', 'c2c', or empty")
	}

	var totalIndexed int32
	var totalFailed int32
	var errors []string
	offset := 0

	// Fetch and index products in batches
	for {
		// Check context for cancellation
		select {
		case <-ctx.Done():
			duration := int(time.Since(startTime).Seconds())
			s.logger.Warn().
				Int32("indexed", totalIndexed).
				Int32("failed", totalFailed).
				Int("duration_seconds", duration).
				Msg("reindexing cancelled")
			return totalIndexed, totalFailed, duration, errors, ctx.Err()
		default:
		}

		// Fetch batch of products
		// Note: Need to implement GetProductsBatch in repository
		// For now, use a simpler approach with ListProducts
		s.logger.Debug().Int("offset", offset).Int("batch_size", batchSize).Msg("fetching batch")

		// Build filter for listing query
		filter := &domain.ListListingsFilter{
			Limit:  int32(batchSize),
			Offset: int32(offset),
		}
		if sourceType != "" {
			filter.SourceType = &sourceType
		}

		// Fetch listings (products)
		listings, total, err := s.repo.ListListings(ctx, filter)
		if err != nil {
			s.logger.Error().Err(err).Msg("failed to fetch products batch")
			errors = append(errors, fmt.Sprintf("batch fetch error at offset %d: %v", offset, err))
			totalFailed += int32(batchSize) // Approximate
			break
		}

		// No more products
		if len(listings) == 0 {
			s.logger.Debug().Msg("no more products to index")
			break
		}

		s.logger.Debug().Int("count", len(listings)).Int32("total", total).Msg("fetched batch")

		// Load images for each listing
		for _, listing := range listings {
			images, err := s.repo.GetImages(ctx, listing.ID)
			if err != nil {
				s.logger.Warn().Err(err).Int64("listing_id", listing.ID).Msg("failed to load images")
			} else {
				listing.Images = images
			}
		}

		// Bulk index the batch using OpenSearch client
		// We need to access the underlying OpenSearch client for bulk operations
		// For now, index one by one (can be optimized later with BulkIndexProducts)
		for _, listing := range listings {
			if err := s.indexer.IndexListing(ctx, listing); err != nil {
				s.logger.Warn().Err(err).Int64("listing_id", listing.ID).Msg("failed to index listing")
				totalFailed++
				if len(errors) < 10 { // Keep max 10 sample errors
					errors = append(errors, fmt.Sprintf("listing %d: %v", listing.ID, err))
				}
			} else {
				totalIndexed++
			}
		}

		s.logger.Info().
			Int32("indexed_so_far", totalIndexed).
			Int32("failed_so_far", totalFailed).
			Int("batch_size", len(listings)).
			Msg("batch indexed")

		// Move to next batch
		offset += batchSize

		// Check if we've processed all
		if offset >= int(total) {
			break
		}
	}

	duration := int(time.Since(startTime).Seconds())

	s.logger.Info().
		Int32("total_indexed", totalIndexed).
		Int32("total_failed", totalFailed).
		Int("duration_seconds", duration).
		Msg("reindexing completed")

	return totalIndexed, totalFailed, duration, errors, nil
}

// convertProductToListing converts domain.Product to domain.Listing for OpenSearch indexing
// This is necessary because OpenSearch Client expects domain.Listing format
// Note: Product and Listing are different entities, this is a mapping for search indexing only
func convertProductToListing(product *domain.Product) *domain.Listing {
	if product == nil {
		return nil
	}

	// Create a minimal listing representation for search indexing
	listing := &domain.Listing{
		ID:         product.ID,
		Title:      product.Name,
		Price:      product.Price,
		Currency:   product.Currency,
		CategoryID: product.CategoryID,
		Quantity:   product.StockQuantity,
		ViewsCount: product.ViewCount,
		CreatedAt:  product.CreatedAt,
		UpdatedAt:  product.UpdatedAt,
	}

	// Set source type as b2c
	listing.SourceType = "b2c"

	// Storefront ID
	listing.StorefrontID = &product.StorefrontID

	// Optional description
	desc := product.Description
	if desc != "" {
		listing.Description = &desc
	}

	// Optional SKU
	if product.SKU != nil {
		listing.SKU = product.SKU
	}

	// Stock status
	stockStatus := product.StockStatus
	listing.StockStatus = &stockStatus

	// Set default status and visibility for B2C products
	if product.IsActive {
		listing.Status = "active"
		listing.Visibility = "public"
	} else {
		listing.Status = "inactive"
		listing.Visibility = "hidden"
	}

	// TODO: Convert attributes to JSON if present
	// Attributes are already map[string]interface{}, we need to store as JSONB
	// For now, we skip attribute conversion as it needs proper JSONB handling

	// Convert images
	if len(product.Images) > 0 {
		listing.Images = make([]*domain.ListingImage, 0, len(product.Images))
		for _, img := range product.Images {
			listing.Images = append(listing.Images, &domain.ListingImage{
				ID:        img.ID,
				URL:       img.URL,
				IsPrimary: img.IsPrimary,
			})
		}
	}

	// Convert location
	if product.HasIndividualLocation {
		listing.Location = &domain.ListingLocation{}
		if product.IndividualAddress != nil {
			listing.Location.AddressLine1 = product.IndividualAddress
		}
		if product.IndividualLatitude != nil {
			listing.Location.Latitude = product.IndividualLatitude
		}
		if product.IndividualLongitude != nil {
			listing.Location.Longitude = product.IndividualLongitude
		}
	}

	return listing
}

// =============================================================================
// Product Images (B2C)
// =============================================================================

// GetProductImageByID retrieves a single product image by ID
func (s *Service) GetProductImageByID(ctx context.Context, imageID int64) (*domain.ProductImage, error) {
	return s.repo.GetProductImageByID(ctx, imageID)
}

// AddProductImage adds a new image to a B2C product
func (s *Service) AddProductImage(ctx context.Context, image *domain.ProductImage) (*domain.ProductImage, error) {
	result, err := s.repo.AddProductImage(ctx, image)
	if err != nil {
		return nil, err
	}

	// Reindex product in OpenSearch after image add (async)
	if s.indexer != nil && result.ProductID != nil && *result.ProductID > 0 {
		productID := *result.ProductID
		go func() {
			indexCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			// Get full product with images
			fullProduct, err := s.repo.GetProductByID(indexCtx, productID, nil)
			if err != nil {
				s.logger.Warn().Err(err).Int64("product_id", productID).
					Msg("failed to get product for reindexing after image add (non-blocking)")
				return
			}

			// Convert to listing for indexing
			listing := convertProductToListing(fullProduct)

			// Update in OpenSearch
			if err := s.indexer.UpdateListing(indexCtx, listing); err != nil {
				s.logger.Warn().Err(err).Int64("product_id", productID).
					Msg("OpenSearch reindexing failed after product image add (non-blocking)")
			} else {
				s.logger.Debug().Int64("product_id", productID).Msg("product reindexed in OpenSearch after image add")
			}
		}()
	}

	return result, nil
}

// GetProductImages retrieves all images for a B2C product
func (s *Service) GetProductImages(ctx context.Context, productID int64) ([]*domain.ProductImage, error) {
	return s.repo.GetProductImages(ctx, productID)
}

// DeleteProductImage removes a product image
func (s *Service) DeleteProductImage(ctx context.Context, imageID int64) error {
	return s.repo.DeleteProductImage(ctx, imageID)
}

// ReorderProductImages updates display order for product images
func (s *Service) ReorderProductImages(ctx context.Context, productID int64, orders []postgres.ProductImageOrder) error {
	return s.repo.ReorderProductImages(ctx, productID, orders)
}
