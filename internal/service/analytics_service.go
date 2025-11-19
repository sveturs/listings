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

	"github.com/sveturs/listings/internal/domain"
	listingssvcv1 "github.com/sveturs/listings/api/proto/listings/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ============================================================================
// INTERFACES
// ============================================================================

// AnalyticsRepository defines minimal repository interface for analytics data access
type AnalyticsRepository interface {
	// GetOverviewStats retrieves platform-wide analytics statistics
	GetOverviewStats(ctx context.Context, filter *domain.GetOverviewStatsFilter) (*domain.OverviewStats, error)

	// GetListingStats retrieves analytics for a specific listing
	GetListingStats(ctx context.Context, filter *domain.GetListingStatsFilter) (*domain.ListingStats, error)
}

// AnalyticsService defines the service interface for analytics operations
type AnalyticsService interface {
	// GetOverviewStats retrieves platform-wide analytics (admin only)
	GetOverviewStats(ctx context.Context, req *listingssvcv1.GetOverviewStatsRequest, userID int64, isAdmin bool) (*listingssvcv1.GetOverviewStatsResponse, error)

	// GetListingStats retrieves analytics for a specific listing (owner or admin)
	GetListingStats(ctx context.Context, req *listingssvcv1.GetListingStatsRequest, userID int64, isAdmin bool) (*listingssvcv1.GetListingStatsResponse, error)
}

// ============================================================================
// CACHE CONFIGURATION
// ============================================================================

const (
	// Cache keys
	analyticsCacheKeyOverview = "analytics:overview:%s" // %s = MD5 hash of request
	analyticsCacheKeyListing  = "analytics:listing:%s"  // %s = MD5 hash of request

	// Cache TTLs
	overviewStatsCacheTTL = 1 * time.Hour  // Overview stats cached for 1 hour
	listingStatsCacheTTL  = 15 * time.Minute // Listing stats cached for 15 minutes

	// Validation constants
	maxDateRangeDays = 365 // Maximum 365 days for analytics queries
)

// ============================================================================
// IMPLEMENTATION
// ============================================================================

// analyticsServiceImpl implements AnalyticsService interface
type analyticsServiceImpl struct {
	repo   AnalyticsRepository
	cache  *AnalyticsCache
	logger zerolog.Logger
}

// AnalyticsCache provides caching functionality for analytics
type AnalyticsCache struct {
	client redis.UniversalClient
	logger zerolog.Logger
}

// NewAnalyticsCache creates a new analytics cache service
func NewAnalyticsCache(client redis.UniversalClient, logger zerolog.Logger) *AnalyticsCache {
	return &AnalyticsCache{
		client: client,
		logger: logger.With().Str("component", "analytics_cache").Logger(),
	}
}

// NewAnalyticsService creates a new analytics service
func NewAnalyticsService(
	repo AnalyticsRepository,
	cacheClient redis.UniversalClient,
	logger zerolog.Logger,
) AnalyticsService {
	return &analyticsServiceImpl{
		repo:   repo,
		cache:  NewAnalyticsCache(cacheClient, logger),
		logger: logger.With().Str("component", "analytics_service").Logger(),
	}
}

// ============================================================================
// PUBLIC SERVICE METHODS
// ============================================================================

// GetOverviewStats retrieves platform-wide analytics statistics
func (s *analyticsServiceImpl) GetOverviewStats(
	ctx context.Context,
	req *listingssvcv1.GetOverviewStatsRequest,
	userID int64,
	isAdmin bool,
) (*listingssvcv1.GetOverviewStatsResponse, error) {
	// Authorization: only admins can access overview stats
	if err := s.requireAdmin(userID, isAdmin); err != nil {
		s.logger.Warn().
			Int64("user_id", userID).
			Bool("is_admin", isAdmin).
			Msg("unauthorized access to overview stats")
		return nil, err
	}

	// Validate request
	if err := s.validateOverviewStatsRequest(req); err != nil {
		s.logger.Warn().
			Err(err).
			Msg("invalid overview stats request")
		return nil, fmt.Errorf("%w: %v", ErrInvalidInput, err)
	}

	// Try cache first
	cacheKey := s.generateOverviewCacheKey(req)
	if cached, err := s.cache.GetOverviewStats(ctx, cacheKey); err == nil && cached != nil {
		s.logger.Debug().
			Str("cache_key", cacheKey).
			Msg("overview stats retrieved from cache")
		return cached, nil
	}

	s.logger.Debug().
		Str("cache_key", cacheKey).
		Msg("cache miss, fetching from repository")

	// Build domain filter
	filter := s.buildOverviewStatsFilter(req)

	// Fetch from repository
	stats, err := s.repo.GetOverviewStats(ctx, filter)
	if err != nil {
		s.logger.Error().
			Err(err).
			Interface("filter", filter).
			Msg("failed to get overview stats from repository")
		return nil, fmt.Errorf("%w: failed to retrieve analytics", ErrInternal)
	}

	// Enrich with calculated fields
	stats.EnrichWithCalculatedFields()

	// Convert to proto response
	response := s.convertOverviewStatsToProto(stats)

	// Cache the result
	if err := s.cache.SetOverviewStats(ctx, cacheKey, response, overviewStatsCacheTTL); err != nil {
		s.logger.Warn().
			Err(err).
			Str("cache_key", cacheKey).
			Msg("failed to cache overview stats")
		// Don't fail on cache error
	}

	s.logger.Info().
		Int64("user_id", userID).
		Time("date_from", req.DateFrom.AsTime()).
		Time("date_to", req.DateTo.AsTime()).
		Int64("total_views", stats.TotalViews).
		Int64("total_orders", stats.TotalOrders).
		Msg("overview stats retrieved successfully")

	return response, nil
}

// GetListingStats retrieves analytics for a specific listing
func (s *analyticsServiceImpl) GetListingStats(
	ctx context.Context,
	req *listingssvcv1.GetListingStatsRequest,
	userID int64,
	isAdmin bool,
) (*listingssvcv1.GetListingStatsResponse, error) {
	// Validate request
	if err := s.validateListingStatsRequest(req); err != nil {
		s.logger.Warn().
			Err(err).
			Msg("invalid listing stats request")
		return nil, fmt.Errorf("%w: %v", ErrInvalidInput, err)
	}

	// Extract listing ID
	listingID := s.extractListingID(req)
	if listingID == 0 {
		return nil, fmt.Errorf("%w: listing_id or product_id is required", ErrInvalidInput)
	}

	// Authorization: owner or admin only
	// Note: In production, we'd fetch listing owner from DB and compare
	// For now, we rely on the caller to pass correct userID and isAdmin
	if err := s.requireListingAccess(userID, isAdmin, listingID); err != nil {
		s.logger.Warn().
			Int64("user_id", userID).
			Bool("is_admin", isAdmin).
			Int64("listing_id", listingID).
			Msg("unauthorized access to listing stats")
		return nil, err
	}

	// Try cache first
	cacheKey := s.generateListingCacheKey(req)
	if cached, err := s.cache.GetListingStats(ctx, cacheKey); err == nil && cached != nil {
		s.logger.Debug().
			Str("cache_key", cacheKey).
			Int64("listing_id", listingID).
			Msg("listing stats retrieved from cache")
		return cached, nil
	}

	s.logger.Debug().
		Str("cache_key", cacheKey).
		Int64("listing_id", listingID).
		Msg("cache miss, fetching from repository")

	// Build domain filter
	filter := s.buildListingStatsFilter(req)

	// Fetch from repository
	stats, err := s.repo.GetListingStats(ctx, filter)
	if err != nil {
		s.logger.Error().
			Err(err).
			Interface("filter", filter).
			Int64("listing_id", listingID).
			Msg("failed to get listing stats from repository")
		return nil, fmt.Errorf("%w: failed to retrieve listing analytics", ErrInternal)
	}

	// Enrich with calculated fields
	stats.EnrichWithCalculatedFields()

	// Convert to proto response
	response := s.convertListingStatsToProto(stats)

	// Cache the result
	if err := s.cache.SetListingStats(ctx, cacheKey, response, listingStatsCacheTTL); err != nil {
		s.logger.Warn().
			Err(err).
			Str("cache_key", cacheKey).
			Msg("failed to cache listing stats")
		// Don't fail on cache error
	}

	s.logger.Info().
		Int64("user_id", userID).
		Int64("listing_id", stats.ListingID).
		Time("date_from", req.DateFrom.AsTime()).
		Time("date_to", req.DateTo.AsTime()).
		Int64("views", stats.ViewsCount).
		Int64("orders", stats.OrdersCount).
		Msg("listing stats retrieved successfully")

	return response, nil
}

// ============================================================================
// VALIDATION METHODS
// ============================================================================

// validateOverviewStatsRequest validates GetOverviewStatsRequest
func (s *analyticsServiceImpl) validateOverviewStatsRequest(req *listingssvcv1.GetOverviewStatsRequest) error {
	if req == nil {
		return fmt.Errorf("request cannot be nil")
	}

	// Validate date range
	if req.DateFrom == nil {
		return fmt.Errorf("date_from is required")
	}
	if req.DateTo == nil {
		return fmt.Errorf("date_to is required")
	}

	dateFrom := req.DateFrom.AsTime()
	dateTo := req.DateTo.AsTime()

	if dateFrom.After(dateTo) {
		return fmt.Errorf("date_from must be before or equal to date_to")
	}

	// Check maximum date range (365 days)
	daysDiff := dateTo.Sub(dateFrom).Hours() / 24
	if daysDiff > float64(maxDateRangeDays) {
		return fmt.Errorf("date range cannot exceed %d days", maxDateRangeDays)
	}

	// Validate optional fields
	if req.ListingType != nil && *req.ListingType != "" {
		listingType := *req.ListingType
		if listingType != "b2c" && listingType != "c2c" {
			return fmt.Errorf("listing_type must be either 'b2c' or 'c2c'")
		}
	}

	return nil
}

// validateListingStatsRequest validates GetListingStatsRequest
func (s *analyticsServiceImpl) validateListingStatsRequest(req *listingssvcv1.GetListingStatsRequest) error {
	if req == nil {
		return fmt.Errorf("request cannot be nil")
	}

	// Must have listing_id or product_id
	if req.GetListingId() == 0 && req.GetProductId() == 0 {
		return fmt.Errorf("either listing_id or product_id must be provided")
	}

	// Validate date range
	if req.DateFrom == nil {
		return fmt.Errorf("date_from is required")
	}
	if req.DateTo == nil {
		return fmt.Errorf("date_to is required")
	}

	dateFrom := req.DateFrom.AsTime()
	dateTo := req.DateTo.AsTime()

	if dateFrom.After(dateTo) {
		return fmt.Errorf("date_from must be before or equal to date_to")
	}

	// Check maximum date range (365 days)
	daysDiff := dateTo.Sub(dateFrom).Hours() / 24
	if daysDiff > float64(maxDateRangeDays) {
		return fmt.Errorf("date range cannot exceed %d days", maxDateRangeDays)
	}

	return nil
}

// ============================================================================
// AUTHORIZATION HELPERS
// ============================================================================

// requireAdmin checks if user is admin
func (s *analyticsServiceImpl) requireAdmin(userID int64, isAdmin bool) error {
	if !isAdmin {
		return fmt.Errorf("%w: admin access required", ErrUnauthorized)
	}
	return nil
}

// requireListingAccess checks if user has access to listing stats (owner or admin)
func (s *analyticsServiceImpl) requireListingAccess(userID int64, isAdmin bool, listingID int64) error {
	// Admin has access to all listings
	if isAdmin {
		return nil
	}

	// TODO: In production, fetch listing from DB and check if userID matches owner_id
	// For now, we assume the caller has already verified ownership
	// This is a placeholder that should be replaced with actual DB lookup
	if userID <= 0 {
		return fmt.Errorf("%w: authentication required", ErrUnauthorized)
	}

	return nil
}

// ============================================================================
// FILTER BUILDERS
// ============================================================================

// buildOverviewStatsFilter converts proto request to domain filter
func (s *analyticsServiceImpl) buildOverviewStatsFilter(req *listingssvcv1.GetOverviewStatsRequest) *domain.GetOverviewStatsFilter {
	filter := &domain.GetOverviewStatsFilter{
		StartDate:   req.DateFrom.AsTime(),
		EndDate:     req.DateTo.AsTime(),
		Granularity: s.convertMetricPeriodToGranularity(req.GetPeriod()),
	}

	// Optional filters
	if req.StorefrontId != nil {
		storefrontID := *req.StorefrontId
		filter.StorefrontID = &storefrontID
	}

	if req.CategoryId != nil {
		categoryID := *req.CategoryId
		filter.CategoryID = &categoryID
	}

	if req.ListingType != nil && *req.ListingType != "" {
		listingType := *req.ListingType
		filter.SourceType = &listingType
	}

	// Set reasonable defaults for pagination
	filter.Limit = 100
	filter.Offset = 0

	return filter
}

// buildListingStatsFilter converts proto request to domain filter
func (s *analyticsServiceImpl) buildListingStatsFilter(req *listingssvcv1.GetListingStatsRequest) *domain.GetListingStatsFilter {
	filter := &domain.GetListingStatsFilter{
		Granularity:       s.convertMetricPeriodToGranularity(req.GetPeriod()),
		IncludeTimeSeries: true,
		TimeSeriesLimit:   100,
	}

	// Set listing ID
	listingID := s.extractListingID(req)
	filter.ListingID = &listingID

	// Set date range
	startDate := req.DateFrom.AsTime()
	endDate := req.DateTo.AsTime()
	filter.StartDate = &startDate
	filter.EndDate = &endDate

	return filter
}

// ============================================================================
// PROTO CONVERTERS
// ============================================================================

// convertOverviewStatsToProto converts domain.OverviewStats to proto response
func (s *analyticsServiceImpl) convertOverviewStatsToProto(stats *domain.OverviewStats) *listingssvcv1.GetOverviewStatsResponse {
	response := &listingssvcv1.GetOverviewStatsResponse{
		Listings: &listingssvcv1.ListingsStats{
			TotalListings:  int32(stats.ActiveListings),
			ActiveListings: int32(stats.ActiveListings),
		},
		Revenue: &listingssvcv1.RevenueStats{
			TotalRevenue:    stats.TotalRevenue,
			AvgOrderValue:   stats.AverageOrderValue,
			Transactions:    int32(stats.TotalOrders),
		},
		Users: &listingssvcv1.UsersStats{
			ActiveUsers: int32(stats.ActiveUsers),
		},
		Orders: &listingssvcv1.OrdersStats{
			TotalOrders:     int32(stats.TotalOrders),
			CompletedOrders: int32(stats.TotalOrders),
		},
		Engagement: &listingssvcv1.EngagementMetrics{
			TotalViews:     int32(stats.TotalViews),
			FavoritesAdded: int32(stats.TotalFavorites),
		},
		ConversionFunnel: &listingssvcv1.ConversionFunnel{
			StageViews:             int32(stats.TotalViews),
			StageCompleted:         int32(stats.TotalOrders),
			OverallConversionRate:  stats.ConversionRate,
		},
		GeneratedAt: timestamppb.New(time.Now()),
		DataFrom:    timestamppb.New(stats.PeriodStart),
		DataTo:      timestamppb.New(stats.PeriodEnd),
	}

	// Convert time series if present
	if stats.TimeSeries != nil && len(stats.TimeSeries) > 0 {
		response.TimeSeries = make([]*listingssvcv1.TimeSeriesPoint, len(stats.TimeSeries))
		for i, point := range stats.TimeSeries {
			response.TimeSeries[i] = &listingssvcv1.TimeSeriesPoint{
				Timestamp:      timestamppb.New(point.Timestamp),
				Views:          int32(point.Views),
				Orders:         int32(point.Orders),
				Revenue:        point.Revenue,
				Favorites:      int32(point.Favorites),
				ConversionRate: point.ConversionRate,
			}
		}
	}

	return response
}

// convertListingStatsToProto converts domain.ListingStats to proto response
func (s *analyticsServiceImpl) convertListingStatsToProto(stats *domain.ListingStats) *listingssvcv1.GetListingStatsResponse {
	response := &listingssvcv1.GetListingStatsResponse{
		ListingId:       stats.ListingID,
		ListingName:     "", // TODO: Fetch from listing table
		ListingType:     "b2c", // TODO: Determine from listing
		TotalViews:      int32(stats.ViewsCount),
		FavoriteCount:   int32(stats.FavoritesCount),
		TotalSales:      int32(stats.OrdersCount),
		TotalRevenue:    stats.TotalRevenue,
		AvgOrderValue:   stats.AverageOrderValue,
		ConversionRate:  stats.ConversionRate,
		Engagement: &listingssvcv1.EngagementMetrics{
			TotalViews:     int32(stats.ViewsCount),
			FavoritesAdded: int32(stats.FavoritesCount),
		},
		CreatedAt:   timestamppb.New(stats.CreatedAt),
		UpdatedAt:   timestamppb.New(stats.UpdatedAt),
		GeneratedAt: timestamppb.New(time.Now()),
		DataFrom:    timestamppb.New(stats.PeriodStart),
		DataTo:      timestamppb.New(stats.PeriodEnd),
	}

	// Convert time series if present
	if stats.TimeSeries != nil && len(stats.TimeSeries) > 0 {
		response.TimeSeries = make([]*listingssvcv1.ListingTimeSeriesPoint, len(stats.TimeSeries))
		for i, point := range stats.TimeSeries {
			response.TimeSeries[i] = &listingssvcv1.ListingTimeSeriesPoint{
				Timestamp:      timestamppb.New(point.Timestamp),
				Views:          int32(point.Views),
				Sales:          int32(point.Orders),
				Revenue:        point.Revenue,
				Favorites:      int32(point.Favorites),
				ConversionRate: point.ConversionRate,
			}
		}
	}

	return response
}

// ============================================================================
// CACHE KEY GENERATION
// ============================================================================

// generateOverviewCacheKey generates MD5-based cache key for overview stats
func (s *analyticsServiceImpl) generateOverviewCacheKey(req *listingssvcv1.GetOverviewStatsRequest) string {
	// Build cache key components
	keyData := struct {
		DateFrom     time.Time
		DateTo       time.Time
		Period       listingssvcv1.MetricPeriod
		StorefrontID *int64
		CategoryID   *int64
		ListingType  *string
	}{
		DateFrom:     req.DateFrom.AsTime(),
		DateTo:       req.DateTo.AsTime(),
		Period:       req.GetPeriod(),
		StorefrontID: req.StorefrontId,
		CategoryID:   req.CategoryId,
		ListingType:  req.ListingType,
	}

	// Convert to JSON and hash
	jsonData, _ := json.Marshal(keyData)
	hash := md5.Sum(jsonData)
	hashStr := fmt.Sprintf("%x", hash)

	return fmt.Sprintf(analyticsCacheKeyOverview, hashStr)
}

// generateListingCacheKey generates MD5-based cache key for listing stats
func (s *analyticsServiceImpl) generateListingCacheKey(req *listingssvcv1.GetListingStatsRequest) string {
	// Build cache key components
	keyData := struct {
		ListingID       int64
		ProductID       int64
		DateFrom        time.Time
		DateTo          time.Time
		Period          listingssvcv1.MetricPeriod
		IncludeVariants *bool
		IncludeGeo      *bool
	}{
		ListingID:       req.GetListingId(),
		ProductID:       req.GetProductId(),
		DateFrom:        req.DateFrom.AsTime(),
		DateTo:          req.DateTo.AsTime(),
		Period:          req.GetPeriod(),
		IncludeVariants: req.IncludeVariants,
		IncludeGeo:      req.IncludeGeo,
	}

	// Convert to JSON and hash
	jsonData, _ := json.Marshal(keyData)
	hash := md5.Sum(jsonData)
	hashStr := fmt.Sprintf("%x", hash)

	return fmt.Sprintf(analyticsCacheKeyListing, hashStr)
}

// ============================================================================
// HELPER METHODS
// ============================================================================

// extractListingID extracts listing ID from oneof identifier
func (s *analyticsServiceImpl) extractListingID(req *listingssvcv1.GetListingStatsRequest) int64 {
	switch v := req.Identifier.(type) {
	case *listingssvcv1.GetListingStatsRequest_ListingId:
		return v.ListingId
	case *listingssvcv1.GetListingStatsRequest_ProductId:
		return v.ProductId
	default:
		return 0
	}
}

// convertMetricPeriodToGranularity converts proto MetricPeriod to granularity string
func (s *analyticsServiceImpl) convertMetricPeriodToGranularity(period listingssvcv1.MetricPeriod) string {
	switch period {
	case listingssvcv1.MetricPeriod_METRIC_PERIOD_HOURLY:
		return domain.GranularityHourly
	case listingssvcv1.MetricPeriod_METRIC_PERIOD_DAILY:
		return domain.GranularityDaily
	default:
		return domain.GranularityDaily // Default to daily
	}
}

// ============================================================================
// CACHE OPERATIONS (AnalyticsCache)
// ============================================================================

// GetOverviewStats retrieves overview stats from cache
func (c *AnalyticsCache) GetOverviewStats(ctx context.Context, key string) (*listingssvcv1.GetOverviewStatsResponse, error) {
	data, err := c.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil, nil // Cache miss
	}
	if err != nil {
		c.logger.Warn().Err(err).Str("key", key).Msg("cache get error")
		return nil, err
	}

	var response listingssvcv1.GetOverviewStatsResponse
	if err := json.Unmarshal(data, &response); err != nil {
		c.logger.Error().Err(err).Str("key", key).Msg("failed to unmarshal cached overview stats")
		return nil, err
	}

	return &response, nil
}

// SetOverviewStats stores overview stats in cache
func (c *AnalyticsCache) SetOverviewStats(ctx context.Context, key string, response *listingssvcv1.GetOverviewStatsResponse, ttl time.Duration) error {
	data, err := json.Marshal(response)
	if err != nil {
		c.logger.Error().Err(err).Str("key", key).Msg("failed to marshal overview stats for cache")
		return err
	}

	if err := c.client.Set(ctx, key, data, ttl).Err(); err != nil {
		c.logger.Warn().Err(err).Str("key", key).Msg("cache set error")
		return err
	}

	c.logger.Debug().Str("key", key).Dur("ttl", ttl).Msg("overview stats cached successfully")
	return nil
}

// GetListingStats retrieves listing stats from cache
func (c *AnalyticsCache) GetListingStats(ctx context.Context, key string) (*listingssvcv1.GetListingStatsResponse, error) {
	data, err := c.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil, nil // Cache miss
	}
	if err != nil {
		c.logger.Warn().Err(err).Str("key", key).Msg("cache get error")
		return nil, err
	}

	var response listingssvcv1.GetListingStatsResponse
	if err := json.Unmarshal(data, &response); err != nil {
		c.logger.Error().Err(err).Str("key", key).Msg("failed to unmarshal cached listing stats")
		return nil, err
	}

	return &response, nil
}

// SetListingStats stores listing stats in cache
func (c *AnalyticsCache) SetListingStats(ctx context.Context, key string, response *listingssvcv1.GetListingStatsResponse, ttl time.Duration) error {
	data, err := json.Marshal(response)
	if err != nil {
		c.logger.Error().Err(err).Str("key", key).Msg("failed to marshal listing stats for cache")
		return err
	}

	if err := c.client.Set(ctx, key, data, ttl).Err(); err != nil {
		c.logger.Warn().Err(err).Str("key", key).Msg("cache set error")
		return err
	}

	c.logger.Debug().Str("key", key).Dur("ttl", ttl).Msg("listing stats cached successfully")
	return nil
}
