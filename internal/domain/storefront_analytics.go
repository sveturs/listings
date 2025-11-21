// Package domain defines core business entities and domain models for the listings microservice.
package domain

import (
	"errors"
	"fmt"
	"time"
)

// ============================================================================
// STOREFRONT ANALYTICS
// ============================================================================

// StorefrontStats represents performance analytics for a storefront
type StorefrontStats struct {
	// Basic Info
	StorefrontID   int64  `json:"storefront_id" db:"storefront_id"`
	StorefrontName string `json:"storefront_name" db:"storefront_name"`
	OwnerID        int64  `json:"owner_id" db:"owner_id"`

	// Sales metrics
	TotalSales        int64   `json:"total_sales" db:"total_sales"`               // Number of completed orders (delivered)
	TotalRevenue      float64 `json:"total_revenue" db:"total_revenue"`           // Sum of order amounts
	AverageOrderValue float64 `json:"average_order_value" db:"average_order_value"` // Avg order value

	// Listings metrics
	ActiveListings int32 `json:"active_listings" db:"active_listings"` // Active listings count
	TotalListings  int32 `json:"total_listings" db:"total_listings"`   // All listings

	// Performance metrics
	TotalViews     int64   `json:"total_views" db:"total_views"`         // Sum of all listing views
	TotalFavorites int64   `json:"total_favorites" db:"total_favorites"` // Sum of all favorites
	ConversionRate float64 `json:"conversion_rate" db:"conversion_rate"` // Orders / views ratio (%)

	// Top listings (populated separately)
	TopListings []*TopListingInfo `json:"top_listings,omitempty" db:"-"`

	// Filter info
	Period string `json:"period" db:"-"` // Applied period filter

	// Timestamps
	GeneratedAt   time.Time `json:"generated_at" db:"-"`
	LastUpdatedAt time.Time `json:"last_updated_at" db:"last_updated_at"`
}

// TopListingInfo represents a top-performing listing for a storefront
type TopListingInfo struct {
	ListingID      int64   `json:"listing_id" db:"listing_id"`
	Title          string  `json:"title" db:"title"`
	Revenue        float64 `json:"revenue" db:"revenue"`
	OrderCount     int32   `json:"order_count" db:"order_count"`
	ViewCount      int32   `json:"view_count" db:"view_count"`
	ConversionRate float64 `json:"conversion_rate" db:"conversion_rate"` // Calculated: (orders / views) * 100
}

// StorefrontStatsRequest represents a request for storefront analytics
type StorefrontStatsRequest struct {
	StorefrontID int64    `json:"storefront_id"`
	Period       string   `json:"period"` // "7d", "30d", "90d", "all"
	UserID       int64    `json:"user_id"`
	Roles        []string `json:"roles"`
}

// ============================================================================
// VALIDATION
// ============================================================================

// ValidPeriods defines allowed period values
var ValidPeriods = []string{"7d", "30d", "90d", "all"}

// Validate validates the storefront stats request
func (r *StorefrontStatsRequest) Validate() error {
	if r.StorefrontID <= 0 {
		return errors.New("storefront_id must be positive")
	}

	if r.UserID <= 0 {
		return errors.New("user_id is required")
	}

	// Validate period if provided
	if r.Period != "" {
		valid := false
		for _, p := range ValidPeriods {
			if r.Period == p {
				valid = true
				break
			}
		}
		if !valid {
			return fmt.Errorf("invalid period: %s (allowed: 7d, 30d, 90d, all)", r.Period)
		}
	}

	return nil
}

// SetDefaults sets default values for optional fields
func (r *StorefrontStatsRequest) SetDefaults() {
	if r.Period == "" {
		r.Period = "30d" // Default to 30 days
	}
}

// IsAdmin checks if the user has admin role
func (r *StorefrontStatsRequest) IsAdmin() bool {
	for _, role := range r.Roles {
		if role == "admin" {
			return true
		}
	}
	return false
}

// ============================================================================
// HELPER METHODS
// ============================================================================

// GetPeriodDuration returns the time duration for the given period string
func GetPeriodDuration(period string) (time.Duration, error) {
	switch period {
	case "7d":
		return 7 * 24 * time.Hour, nil
	case "30d":
		return 30 * 24 * time.Hour, nil
	case "90d":
		return 90 * 24 * time.Hour, nil
	case "all":
		return 0, nil // Special case: no time filter
	default:
		return 0, fmt.Errorf("invalid period: %s", period)
	}
}

// GetPeriodStartTime returns the start time for the given period
func GetPeriodStartTime(period string, now time.Time) (time.Time, error) {
	duration, err := GetPeriodDuration(period)
	if err != nil {
		return time.Time{}, err
	}

	if duration == 0 {
		// "all" period: use epoch (or a reasonable minimum date)
		return time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), nil
	}

	return now.Add(-duration), nil
}
