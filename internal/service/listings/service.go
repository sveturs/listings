package listings

import (
	"context"
	"fmt"

	"github.com/go-playground/validator/v10"
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
	// Try cache first
	cacheKey := fmt.Sprintf("listing:%d", id)
	var cachedListing domain.Listing

	if err := s.cache.Get(ctx, cacheKey, &cachedListing); err == nil {
		s.logger.Debug().Int64("listing_id", id).Msg("listing found in cache")
		return &cachedListing, nil
	}

	// Cache miss - fetch from database
	listing, err := s.repo.GetListingByID(ctx, id)
	if err != nil {
		s.logger.Error().Err(err).Int64("listing_id", id).Msg("failed to get listing")
		return nil, fmt.Errorf("failed to get listing: %w", err)
	}

	// Store in cache (non-blocking, ignore errors)
	go func() {
		if err := s.cache.Set(context.Background(), cacheKey, listing); err != nil {
			s.logger.Warn().Err(err).Int64("listing_id", id).Msg("failed to cache listing")
		}
	}()

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

	// Invalidate cache
	cacheKey := fmt.Sprintf("listing:%d", id)
	if err := s.cache.Delete(ctx, cacheKey); err != nil {
		s.logger.Warn().Err(err).Int64("listing_id", id).Msg("failed to invalidate cache")
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

	// Invalidate cache
	cacheKey := fmt.Sprintf("listing:%d", id)
	if err := s.cache.Delete(ctx, cacheKey); err != nil {
		s.logger.Warn().Err(err).Int64("listing_id", id).Msg("failed to invalidate cache")
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

	// Try cache for search results
	cacheKey := fmt.Sprintf("search:%s:%d:%d:%d", query.Query, query.CategoryID, query.Limit, query.Offset)
	var cachedResults []*domain.Listing
	var cachedTotal int32

	if err := s.cache.Get(ctx, cacheKey, &cachedResults); err == nil {
		s.logger.Debug().Str("query", query.Query).Msg("search results found in cache")
		return cachedResults, cachedTotal, nil
	}

	// Cache miss - search in database
	listings, total, err := s.repo.SearchListings(ctx, query)
	if err != nil {
		s.logger.Error().Err(err).Str("query", query.Query).Msg("failed to search listings")
		return nil, 0, fmt.Errorf("failed to search listings: %w", err)
	}

	// Cache search results (non-blocking)
	go func() {
		if err := s.cache.Set(context.Background(), cacheKey, listings); err != nil {
			s.logger.Warn().Err(err).Str("query", query.Query).Msg("failed to cache search results")
		}
	}()

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

	// Invalidate cache
	cacheKey := fmt.Sprintf("listing:%d", id)
	if err := s.cache.Delete(ctx, cacheKey); err != nil {
		s.logger.Warn().Err(err).Int64("listing_id", id).Msg("failed to invalidate cache")
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

	// Invalidate cache
	cacheKey := fmt.Sprintf("listing:%d", id)
	if err := s.cache.Delete(ctx, cacheKey); err != nil {
		s.logger.Warn().Err(err).Int64("listing_id", id).Msg("failed to invalidate cache")
	}

	// Enqueue for async index deletion
	if err := s.repo.EnqueueIndexing(ctx, id, domain.IndexOpDelete); err != nil {
		s.logger.Warn().Err(err).Int64("listing_id", id).Msg("failed to enqueue indexing deletion")
	}

	s.logger.Info().Int64("listing_id", id).Msg("admin deleted listing")
	return nil
}
