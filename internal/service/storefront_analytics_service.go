// Package service implements business logic for the listings microservice.
package service

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"google.golang.org/protobuf/types/known/timestamppb"

	listingssvcv1 "github.com/sveturs/listings/api/proto/listings/v1"
	"github.com/sveturs/listings/internal/domain"
)

// ============================================================================
// INTERFACES
// ============================================================================

// StorefrontAnalyticsRepository defines repository interface for storefront analytics
type StorefrontAnalyticsRepository interface {
	// GetStorefrontStats retrieves performance analytics for a storefront
	GetStorefrontStats(ctx context.Context, storefrontID int64, period string) (*domain.StorefrontStats, error)

	// GetTopListings retrieves top-performing listings for a storefront
	GetTopListings(ctx context.Context, storefrontID int64, period string, limit int) ([]*domain.TopListingInfo, error)

	// GetStorefrontOwnerID retrieves the owner user_id for a storefront
	GetStorefrontOwnerID(ctx context.Context, storefrontID int64) (int64, error)
}

// StorefrontAnalyticsService defines the service interface for storefront analytics operations
type StorefrontAnalyticsService interface {
	// GetStorefrontStats retrieves analytics for a storefront (owner or admin)
	GetStorefrontStats(ctx context.Context, req *listingssvcv1.GetStorefrontStatsRequest) (*listingssvcv1.GetStorefrontStatsResponse, error)
}

// ============================================================================
// CACHE CONFIGURATION
// ============================================================================

const (
	// Cache keys
	storefrontStatsCacheKey = "analytics:storefront:%d:%s" // %d = storefront_id, %s = period

	// Cache TTL
	storefrontStatsCacheTTL = 15 * time.Minute // Storefront stats cached for 15 minutes

	// Top listings limit
	topListingsLimit = 10
)

// ============================================================================
// IMPLEMENTATION
// ============================================================================

// storefrontAnalyticsServiceImpl implements StorefrontAnalyticsService interface
type storefrontAnalyticsServiceImpl struct {
	repo   StorefrontAnalyticsRepository
	cache  *StorefrontAnalyticsCache
	logger zerolog.Logger
}

// StorefrontAnalyticsCache provides caching functionality for storefront analytics
type StorefrontAnalyticsCache struct {
	client redis.UniversalClient
	logger zerolog.Logger
}

// NewStorefrontAnalyticsCache creates a new storefront analytics cache service
func NewStorefrontAnalyticsCache(client redis.UniversalClient, logger zerolog.Logger) *StorefrontAnalyticsCache {
	return &StorefrontAnalyticsCache{
		client: client,
		logger: logger.With().Str("component", "storefront_analytics_cache").Logger(),
	}
}

// NewStorefrontAnalyticsService creates a new storefront analytics service
func NewStorefrontAnalyticsService(
	repo StorefrontAnalyticsRepository,
	cacheClient redis.UniversalClient,
	logger zerolog.Logger,
) StorefrontAnalyticsService {
	return &storefrontAnalyticsServiceImpl{
		repo:   repo,
		cache:  NewStorefrontAnalyticsCache(cacheClient, logger),
		logger: logger.With().Str("component", "storefront_analytics_service").Logger(),
	}
}

// ============================================================================
// PUBLIC SERVICE METHODS
// ============================================================================

// GetStorefrontStats retrieves analytics for a storefront with authorization
func (s *storefrontAnalyticsServiceImpl) GetStorefrontStats(
	ctx context.Context,
	req *listingssvcv1.GetStorefrontStatsRequest,
) (*listingssvcv1.GetStorefrontStatsResponse, error) {
	// Convert request to domain model
	domainReq := &domain.StorefrontStatsRequest{
		StorefrontID: req.StorefrontId,
		Period:       req.GetPeriod(),
		UserID:       req.UserId,
		Roles:        req.Roles,
	}

	// Set defaults
	domainReq.SetDefaults()

	// Validate request
	if err := domainReq.Validate(); err != nil {
		s.logger.Warn().
			Err(err).
			Int64("storefront_id", req.StorefrontId).
			Msg("invalid storefront stats request")
		return nil, fmt.Errorf("%w: %v", ErrInvalidInput, err)
	}

	// Authorization check: admin OR owner
	if err := s.authorizeStorefrontAccess(ctx, domainReq); err != nil {
		s.logger.Warn().
			Err(err).
			Int64("storefront_id", req.StorefrontId).
			Int64("user_id", req.UserId).
			Strs("roles", req.Roles).
			Msg("unauthorized access to storefront stats")
		return nil, err
	}

	// Try cache first
	cacheKey := fmt.Sprintf(storefrontStatsCacheKey, req.StorefrontId, domainReq.Period)
	if cached, err := s.cache.Get(ctx, cacheKey); err == nil && cached != nil {
		s.logger.Debug().
			Str("cache_key", cacheKey).
			Msg("storefront stats retrieved from cache")
		return cached, nil
	}

	s.logger.Debug().
		Str("cache_key", cacheKey).
		Msg("cache miss, fetching from repository")

	// Fetch from repository
	stats, err := s.repo.GetStorefrontStats(ctx, req.StorefrontId, domainReq.Period)
	if err != nil {
		s.logger.Error().
			Err(err).
			Int64("storefront_id", req.StorefrontId).
			Str("period", domainReq.Period).
			Msg("failed to get storefront stats from repository")
		return nil, fmt.Errorf("%w: failed to retrieve storefront analytics", ErrInternal)
	}

	// Fetch top listings
	topListings, err := s.repo.GetTopListings(ctx, req.StorefrontId, domainReq.Period, topListingsLimit)
	if err != nil {
		s.logger.Error().
			Err(err).
			Int64("storefront_id", req.StorefrontId).
			Msg("failed to get top listings")
		// Don't fail the whole request, just log and continue with empty list
		topListings = []*domain.TopListingInfo{}
	}

	// Convert to proto response
	response := s.domainToProto(stats, topListings)

	// Cache the result
	if err := s.cache.Set(ctx, cacheKey, response, storefrontStatsCacheTTL); err != nil {
		s.logger.Warn().
			Err(err).
			Str("cache_key", cacheKey).
			Msg("failed to cache storefront stats")
		// Don't fail the request, just log the warning
	}

	s.logger.Info().
		Int64("storefront_id", req.StorefrontId).
		Str("period", domainReq.Period).
		Int64("total_sales", stats.TotalSales).
		Float64("total_revenue", stats.TotalRevenue).
		Msg("storefront stats retrieved successfully")

	return response, nil
}

// ============================================================================
// AUTHORIZATION
// ============================================================================

// authorizeStorefrontAccess checks if user is admin OR storefront owner
func (s *storefrontAnalyticsServiceImpl) authorizeStorefrontAccess(
	ctx context.Context,
	req *domain.StorefrontStatsRequest,
) error {
	// Admin can access any storefront
	if req.IsAdmin() {
		s.logger.Debug().
			Int64("user_id", req.UserID).
			Int64("storefront_id", req.StorefrontID).
			Msg("admin access granted")
		return nil
	}

	// Check if user is the storefront owner
	ownerID, err := s.repo.GetStorefrontOwnerID(ctx, req.StorefrontID)
	if err != nil {
		s.logger.Error().
			Err(err).
			Int64("storefront_id", req.StorefrontID).
			Msg("failed to get storefront owner")
		return fmt.Errorf("%w: failed to verify storefront ownership", ErrInternal)
	}

	if req.UserID != ownerID {
		s.logger.Warn().
			Int64("user_id", req.UserID).
			Int64("owner_id", ownerID).
			Int64("storefront_id", req.StorefrontID).
			Msg("user is not storefront owner")
		return fmt.Errorf("%w: you don't have permission to view this storefront's analytics", ErrUnauthorized)
	}

	s.logger.Debug().
		Int64("user_id", req.UserID).
		Int64("storefront_id", req.StorefrontID).
		Msg("owner access granted")

	return nil
}

// ============================================================================
// CONVERSION METHODS
// ============================================================================

// domainToProto converts domain model to proto response
func (s *storefrontAnalyticsServiceImpl) domainToProto(
	stats *domain.StorefrontStats,
	topListings []*domain.TopListingInfo,
) *listingssvcv1.GetStorefrontStatsResponse {
	// Convert top listings
	protoTopListings := make([]*listingssvcv1.TopListingInfo, 0, len(topListings))
	for _, listing := range topListings {
		protoTopListings = append(protoTopListings, &listingssvcv1.TopListingInfo{
			ListingId:      listing.ListingID,
			Title:          listing.Title,
			Revenue:        listing.Revenue,
			OrderCount:     listing.OrderCount,
			ViewCount:      listing.ViewCount,
			ConversionRate: listing.ConversionRate,
		})
	}

	return &listingssvcv1.GetStorefrontStatsResponse{
		StorefrontId:      stats.StorefrontID,
		StorefrontName:    stats.StorefrontName,
		TotalSales:        stats.TotalSales,
		TotalRevenue:      stats.TotalRevenue,
		AverageOrderValue: stats.AverageOrderValue,
		ActiveListings:    stats.ActiveListings,
		TotalListings:     stats.TotalListings,
		TotalViews:        stats.TotalViews,
		TotalFavorites:    stats.TotalFavorites,
		ConversionRate:    stats.ConversionRate,
		TopListings:       protoTopListings,
		Period:            stats.Period,
		GeneratedAt:       timestamppb.New(stats.GeneratedAt),
	}
}

// ============================================================================
// CACHE METHODS
// ============================================================================

// Get retrieves storefront stats from cache
func (c *StorefrontAnalyticsCache) Get(ctx context.Context, key string) (*listingssvcv1.GetStorefrontStatsResponse, error) {
	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("cache miss")
		}
		c.logger.Warn().Err(err).Str("key", key).Msg("cache get error")
		return nil, err
	}

	var stats listingssvcv1.GetStorefrontStatsResponse
	if err := json.Unmarshal([]byte(val), &stats); err != nil {
		c.logger.Error().Err(err).Str("key", key).Msg("failed to unmarshal cached stats")
		return nil, err
	}

	return &stats, nil
}

// Set stores storefront stats in cache
func (c *StorefrontAnalyticsCache) Set(ctx context.Context, key string, stats *listingssvcv1.GetStorefrontStatsResponse, ttl time.Duration) error {
	data, err := json.Marshal(stats)
	if err != nil {
		c.logger.Error().Err(err).Str("key", key).Msg("failed to marshal stats for cache")
		return err
	}

	if err := c.client.Set(ctx, key, data, ttl).Err(); err != nil {
		c.logger.Warn().Err(err).Str("key", key).Msg("failed to set cache")
		return err
	}

	c.logger.Debug().
		Str("key", key).
		Dur("ttl", ttl).
		Msg("stats cached successfully")

	return nil
}

// generateCacheKey generates a deterministic cache key from request parameters
func generateCacheKey(prefix string, params interface{}) string {
	data, _ := json.Marshal(params)
	hash := md5.Sum(data)
	return fmt.Sprintf("%s:%x", prefix, hash)
}
