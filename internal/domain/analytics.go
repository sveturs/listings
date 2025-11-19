// Package domain defines core business entities and domain models for the listings microservice.
package domain

import (
	"errors"
	"fmt"
	"time"
)

// ============================================================================
// ANALYTICS AGGREGATION
// ============================================================================

// OverviewStats represents aggregate analytics statistics across multiple listings
type OverviewStats struct {
	// Time range
	PeriodStart time.Time `json:"period_start" db:"period_start"`
	PeriodEnd   time.Time `json:"period_end" db:"period_end"`

	// Aggregate metrics
	TotalViews           int64   `json:"total_views" db:"total_views"`
	TotalFavorites       int64   `json:"total_favorites" db:"total_favorites"`
	TotalOrders          int64   `json:"total_orders" db:"total_orders"`
	TotalRevenue         float64 `json:"total_revenue" db:"total_revenue"`
	ActiveListings       int64   `json:"active_listings" db:"active_listings"`
	ActiveUsers          int64   `json:"active_users" db:"active_users"`
	ConversionRate       float64 `json:"conversion_rate" db:"conversion_rate"` // (orders / views) * 100
	AverageOrderValue    float64 `json:"average_order_value"`                  // total_revenue / total_orders
	AverageFavoritesRate float64 `json:"average_favorites_rate"`               // (favorites / views) * 100

	// Time series data for trends
	TimeSeries []*TimeSeriesDataPoint `json:"time_series,omitempty" db:"-"`
}

// TimeSeriesDataPoint represents a single data point in a time series (daily/hourly aggregation)
type TimeSeriesDataPoint struct {
	// Identification
	Timestamp time.Time `json:"timestamp" db:"timestamp"` // Day or hour start

	// Metrics
	Views       int64   `json:"views" db:"views"`
	Favorites   int64   `json:"favorites" db:"favorites"`
	Orders      int64   `json:"orders" db:"orders"`
	Revenue     float64 `json:"revenue" db:"revenue"`
	ActiveUsers int64   `json:"active_users" db:"active_users"`

	// Calculated fields
	ConversionRate float64 `json:"conversion_rate" db:"-"` // (orders / views) * 100
}

// ListingStats represents engagement, conversion, and revenue metrics for a single listing
type ListingStats struct {
	// Identification
	ListingID int64  `json:"listing_id" db:"listing_id"`
	UUID      string `json:"uuid" db:"uuid"`

	// Engagement metrics
	ViewsCount     int64 `json:"views_count" db:"views_count"`
	FavoritesCount int64 `json:"favorites_count" db:"favorites_count"`
	InquiriesCount int64 `json:"inquiries_count" db:"inquiries_count"`

	// Conversion & revenue metrics
	OrdersCount      int64   `json:"orders_count" db:"orders_count"`
	TotalRevenue     float64 `json:"total_revenue" db:"total_revenue"`
	ConversionRate   float64 `json:"conversion_rate" db:"-"`       // (orders / views) * 100
	FavoritesRate    float64 `json:"favorites_rate" db:"-"`        // (favorites / views) * 100
	AverageOrderValue float64 `json:"average_order_value" db:"-"`  // total_revenue / orders_count

	// Engagement tracking
	FirstViewedAt   *time.Time `json:"first_viewed_at,omitempty" db:"first_viewed_at"`
	LastViewedAt    *time.Time `json:"last_viewed_at,omitempty" db:"last_viewed_at"`
	LastFavoritedAt *time.Time `json:"last_favorited_at,omitempty" db:"last_favorited_at"`
	LastOrderedAt   *time.Time `json:"last_ordered_at,omitempty" db:"last_ordered_at"`

	// Period tracking
	PeriodStart time.Time `json:"period_start" db:"period_start"`
	PeriodEnd   time.Time `json:"period_end" db:"period_end"`

	// Timestamps
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`

	// Time series data for trends
	TimeSeries []*ListingTimeSeriesPoint `json:"time_series,omitempty" db:"-"`
}

// ListingTimeSeriesPoint represents a single data point in a listing's time series
type ListingTimeSeriesPoint struct {
	// Identification
	ListingID int64     `json:"listing_id" db:"listing_id"`
	Timestamp time.Time `json:"timestamp" db:"timestamp"` // Day or hour start

	// Metrics
	Views     int64   `json:"views" db:"views"`
	Favorites int64   `json:"favorites" db:"favorites"`
	Orders    int64   `json:"orders" db:"orders"`
	Revenue   float64 `json:"revenue" db:"revenue"`

	// Calculated fields
	ConversionRate float64 `json:"conversion_rate" db:"-"` // (orders / views) * 100
}

// ============================================================================
// QUERY FILTERS & INPUT TYPES
// ============================================================================

// GetOverviewStatsFilter represents filters for overview statistics queries
type GetOverviewStatsFilter struct {
	// Time range (required)
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	Granularity string    `json:"granularity" validate:"required,oneof=hourly daily"` // hourly or daily aggregation

	// Optional filters
	UserID       *int64  `json:"user_id,omitempty"`
	CategoryID   *int64  `json:"category_id,omitempty"`
	StorefrontID *int64  `json:"storefront_id,omitempty"`
	SourceType   *string `json:"source_type,omitempty" validate:"omitempty,oneof=c2c b2c"`

	// Pagination (optional, for time series)
	Limit  int32 `json:"limit" validate:"omitempty,gte=1,lte=1000"`
	Offset int32 `json:"offset" validate:"gte=0"`
}

// Validate validates GetOverviewStatsFilter
func (filter *GetOverviewStatsFilter) Validate() error {
	if filter == nil {
		return errors.New("filter cannot be nil")
	}

	// Validate time range
	if filter.StartDate.IsZero() {
		return errors.New("start_date is required")
	}
	if filter.EndDate.IsZero() {
		return errors.New("end_date is required")
	}

	if filter.StartDate.After(filter.EndDate) {
		return errors.New("start_date must be before or equal to end_date")
	}

	// Maximum date range validation (e.g., no more than 2 years)
	maxDuration := 730 * 24 * time.Hour
	if filter.EndDate.Sub(filter.StartDate) > maxDuration {
		return errors.New("date range cannot exceed 730 days")
	}

	// Validate granularity
	if filter.Granularity != "hourly" && filter.Granularity != "daily" {
		return errors.New("granularity must be either 'hourly' or 'daily'")
	}

	// Validate pagination
	if filter.Limit == 0 {
		filter.Limit = 100 // Default limit
	}
	if filter.Limit > 1000 {
		return errors.New("limit cannot exceed 1000")
	}

	return nil
}

// GetListingStatsFilter represents filters for listing statistics queries
type GetListingStatsFilter struct {
	// Identification (one required)
	ListingID *int64  `json:"listing_id,omitempty"`
	ListingUUID *string `json:"listing_uuid,omitempty"`

	// Time range (optional, defaults to all-time)
	StartDate *time.Time `json:"start_date,omitempty"`
	EndDate   *time.Time `json:"end_date,omitempty"`
	Granularity string    `json:"granularity" validate:"omitempty,oneof=hourly daily"` // For time series

	// Include time series data
	IncludeTimeSeries bool  `json:"include_time_series"`
	TimeSeriesLimit   int32 `json:"time_series_limit" validate:"omitempty,gte=1,lte=1000"`
}

// Validate validates GetListingStatsFilter
func (filter *GetListingStatsFilter) Validate() error {
	if filter == nil {
		return errors.New("filter cannot be nil")
	}

	// Must have listing_id or listing_uuid
	if (filter.ListingID == nil || *filter.ListingID <= 0) && (filter.ListingUUID == nil || *filter.ListingUUID == "") {
		return errors.New("either listing_id or listing_uuid must be provided")
	}

	// Validate time range if provided
	if filter.StartDate != nil && filter.EndDate != nil {
		if filter.StartDate.After(*filter.EndDate) {
			return errors.New("start_date must be before or equal to end_date")
		}
	}

	// Validate granularity if time series is requested
	if filter.IncludeTimeSeries && filter.Granularity == "" {
		filter.Granularity = "daily" // Default to daily
	}

	if filter.TimeSeriesLimit == 0 {
		filter.TimeSeriesLimit = 100 // Default limit
	}
	if filter.TimeSeriesLimit > 1000 {
		return errors.New("time_series_limit cannot exceed 1000")
	}

	return nil
}

// ============================================================================
// HELPER METHODS FOR CALCULATIONS
// ============================================================================

// CalculateConversionRate calculates conversion rate from views and orders
func CalculateConversionRate(views, orders int64) float64 {
	if views == 0 {
		return 0
	}
	return (float64(orders) / float64(views)) * 100
}

// CalculateFavoritesRate calculates favorites rate from views and favorites
func CalculateFavoritesRate(views, favorites int64) float64 {
	if views == 0 {
		return 0
	}
	return (float64(favorites) / float64(views)) * 100
}

// CalculateAverageOrderValue calculates average order value
func CalculateAverageOrderValue(revenue float64, orderCount int64) float64 {
	if orderCount == 0 {
		return 0
	}
	return revenue / float64(orderCount)
}

// CalculateCTR calculates click-through rate
func CalculateCTR(impressions, clicks int64) float64 {
	if impressions == 0 {
		return 0
	}
	return (float64(clicks) / float64(impressions)) * 100
}

// ============================================================================
// VALIDATION METHODS
// ============================================================================

// Validate validates OverviewStats
func (s *OverviewStats) Validate() error {
	if s == nil {
		return errors.New("overview stats cannot be nil")
	}

	if s.PeriodStart.IsZero() {
		return errors.New("period_start is required")
	}
	if s.PeriodEnd.IsZero() {
		return errors.New("period_end is required")
	}

	if s.PeriodStart.After(s.PeriodEnd) {
		return errors.New("period_start must be before or equal to period_end")
	}

	if s.TotalViews < 0 {
		return errors.New("total_views cannot be negative")
	}
	if s.TotalFavorites < 0 {
		return errors.New("total_favorites cannot be negative")
	}
	if s.TotalOrders < 0 {
		return errors.New("total_orders cannot be negative")
	}
	if s.TotalRevenue < 0 {
		return errors.New("total_revenue cannot be negative")
	}
	if s.ActiveListings < 0 {
		return errors.New("active_listings cannot be negative")
	}
	if s.ActiveUsers < 0 {
		return errors.New("active_users cannot be negative")
	}

	return nil
}

// Validate validates ListingStats
func (s *ListingStats) Validate() error {
	if s == nil {
		return errors.New("listing stats cannot be nil")
	}

	if s.ListingID <= 0 {
		return errors.New("listing_id must be greater than 0")
	}
	if s.UUID == "" {
		return errors.New("uuid is required")
	}

	if s.ViewsCount < 0 {
		return errors.New("views_count cannot be negative")
	}
	if s.FavoritesCount < 0 {
		return errors.New("favorites_count cannot be negative")
	}
	if s.InquiriesCount < 0 {
		return errors.New("inquiries_count cannot be negative")
	}
	if s.OrdersCount < 0 {
		return errors.New("orders_count cannot be negative")
	}
	if s.TotalRevenue < 0 {
		return errors.New("total_revenue cannot be negative")
	}

	if s.PeriodStart.IsZero() {
		return errors.New("period_start is required")
	}
	if s.PeriodEnd.IsZero() {
		return errors.New("period_end is required")
	}

	if s.PeriodStart.After(s.PeriodEnd) {
		return errors.New("period_start must be before or equal to period_end")
	}

	return nil
}

// Validate validates TimeSeriesDataPoint
func (p *TimeSeriesDataPoint) Validate() error {
	if p == nil {
		return errors.New("time series data point cannot be nil")
	}

	if p.Timestamp.IsZero() {
		return errors.New("timestamp is required")
	}

	if p.Views < 0 {
		return errors.New("views cannot be negative")
	}
	if p.Favorites < 0 {
		return errors.New("favorites cannot be negative")
	}
	if p.Orders < 0 {
		return errors.New("orders cannot be negative")
	}
	if p.Revenue < 0 {
		return errors.New("revenue cannot be negative")
	}
	if p.ActiveUsers < 0 {
		return errors.New("active_users cannot be negative")
	}

	return nil
}

// Validate validates ListingTimeSeriesPoint
func (p *ListingTimeSeriesPoint) Validate() error {
	if p == nil {
		return errors.New("listing time series point cannot be nil")
	}

	if p.ListingID <= 0 {
		return errors.New("listing_id must be greater than 0")
	}
	if p.Timestamp.IsZero() {
		return errors.New("timestamp is required")
	}

	if p.Views < 0 {
		return errors.New("views cannot be negative")
	}
	if p.Favorites < 0 {
		return errors.New("favorites cannot be negative")
	}
	if p.Orders < 0 {
		return errors.New("orders cannot be negative")
	}
	if p.Revenue < 0 {
		return errors.New("revenue cannot be negative")
	}

	return nil
}

// ============================================================================
// ENRICHMENT METHODS
// ============================================================================

// EnrichWithCalculatedFields populates calculated fields based on raw metrics
func (s *OverviewStats) EnrichWithCalculatedFields() {
	s.ConversionRate = CalculateConversionRate(s.TotalViews, s.TotalOrders)
	s.AverageFavoritesRate = CalculateFavoritesRate(s.TotalViews, s.TotalFavorites)

	if s.TotalOrders > 0 {
		s.AverageOrderValue = s.TotalRevenue / float64(s.TotalOrders)
	} else {
		s.AverageOrderValue = 0
	}

	// Enrich time series points
	if s.TimeSeries != nil && len(s.TimeSeries) > 0 {
		for _, point := range s.TimeSeries {
			point.ConversionRate = CalculateConversionRate(point.Views, point.Orders)
		}
	}
}

// EnrichWithCalculatedFields populates calculated fields based on raw metrics
func (s *ListingStats) EnrichWithCalculatedFields() {
	s.ConversionRate = CalculateConversionRate(s.ViewsCount, s.OrdersCount)
	s.FavoritesRate = CalculateFavoritesRate(s.ViewsCount, s.FavoritesCount)

	if s.OrdersCount > 0 {
		s.AverageOrderValue = s.TotalRevenue / float64(s.OrdersCount)
	} else {
		s.AverageOrderValue = 0
	}

	// Enrich time series points
	if s.TimeSeries != nil && len(s.TimeSeries) > 0 {
		for _, point := range s.TimeSeries {
			point.ConversionRate = CalculateConversionRate(point.Views, point.Orders)
		}
	}
}

// ============================================================================
// AGGREGATION & COMPARISON METHODS
// ============================================================================

// AggregateTimeSeriesPoints combines multiple time series points into summary stats
func AggregateTimeSeriesPoints(points []*TimeSeriesDataPoint) *OverviewStats {
	if points == nil || len(points) == 0 {
		return &OverviewStats{
			PeriodStart: time.Now(),
			PeriodEnd:   time.Now(),
		}
	}

	stats := &OverviewStats{
		PeriodStart: points[0].Timestamp,
		PeriodEnd:   points[len(points)-1].Timestamp,
		TimeSeries:  points,
	}

	// Sum all metrics across time series
	for _, point := range points {
		stats.TotalViews += point.Views
		stats.TotalFavorites += point.Favorites
		stats.TotalOrders += point.Orders
		stats.TotalRevenue += point.Revenue
		stats.ActiveUsers += point.ActiveUsers
	}

	// Calculate rates
	stats.EnrichWithCalculatedFields()

	return stats
}

// AggregateListingTimeSeriesPoints combines listing time series points into summary stats
func AggregateListingTimeSeriesPoints(listingID int64, points []*ListingTimeSeriesPoint) *ListingStats {
	if points == nil || len(points) == 0 {
		return &ListingStats{
			ListingID:   listingID,
			PeriodStart: time.Now(),
			PeriodEnd:   time.Now(),
		}
	}

	stats := &ListingStats{
		ListingID:   listingID,
		PeriodStart: points[0].Timestamp,
		PeriodEnd:   points[len(points)-1].Timestamp,
		TimeSeries:  points,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Sum all metrics across time series
	for _, point := range points {
		stats.ViewsCount += point.Views
		stats.FavoritesCount += point.Favorites
		stats.OrdersCount += point.Orders
		stats.TotalRevenue += point.Revenue
	}

	// Calculate rates
	stats.EnrichWithCalculatedFields()

	return stats
}

// ComparePerformance compares two listing stats and returns growth percentages
func ComparePerformance(previous, current *ListingStats) map[string]float64 {
	result := make(map[string]float64)

	if previous == nil {
		// If no previous data, return 100% growth for all metrics
		return map[string]float64{
			"views_growth":     100.0,
			"favorites_growth": 100.0,
			"orders_growth":    100.0,
			"revenue_growth":   100.0,
		}
	}

	// Calculate percentage changes
	if previous.ViewsCount > 0 {
		result["views_growth"] = ((float64(current.ViewsCount) - float64(previous.ViewsCount)) / float64(previous.ViewsCount)) * 100
	}
	if previous.FavoritesCount > 0 {
		result["favorites_growth"] = ((float64(current.FavoritesCount) - float64(previous.FavoritesCount)) / float64(previous.FavoritesCount)) * 100
	}
	if previous.OrdersCount > 0 {
		result["orders_growth"] = ((float64(current.OrdersCount) - float64(previous.OrdersCount)) / float64(previous.OrdersCount)) * 100
	}
	if previous.TotalRevenue > 0 {
		result["revenue_growth"] = ((current.TotalRevenue - previous.TotalRevenue) / previous.TotalRevenue) * 100
	}

	return result
}

// ============================================================================
// FORMAT & PRESENTATION METHODS
// ============================================================================

// String returns a string representation of OverviewStats
func (s *OverviewStats) String() string {
	return fmt.Sprintf(
		"OverviewStats{Period: %s to %s, Views: %d, Favorites: %d, Orders: %d, Revenue: %.2f, ConversionRate: %.2f%%}",
		s.PeriodStart.Format("2006-01-02"),
		s.PeriodEnd.Format("2006-01-02"),
		s.TotalViews,
		s.TotalFavorites,
		s.TotalOrders,
		s.TotalRevenue,
		s.ConversionRate,
	)
}

// String returns a string representation of ListingStats
func (s *ListingStats) String() string {
	return fmt.Sprintf(
		"ListingStats{ListingID: %d, Views: %d, Favorites: %d, Orders: %d, Revenue: %.2f, ConversionRate: %.2f%%}",
		s.ListingID,
		s.ViewsCount,
		s.FavoritesCount,
		s.OrdersCount,
		s.TotalRevenue,
		s.ConversionRate,
	)
}

// ============================================================================
// CONSTANTS
// ============================================================================

const (
	// GranularityHourly represents hourly time series granularity
	GranularityHourly = "hourly"

	// GranularityDaily represents daily time series granularity
	GranularityDaily = "daily"

	// MaxAnalyticsPeriodDays is the maximum number of days for analytics queries
	MaxAnalyticsPeriodDays = 730

	// DefaultTimeSeriesLimit is the default limit for time series results
	DefaultTimeSeriesLimit = 100

	// MaxTimeSeriesLimit is the maximum limit for time series results
	MaxTimeSeriesLimit = 1000
)
