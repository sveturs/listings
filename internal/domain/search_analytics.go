// Package domain defines core business entities for search analytics
package domain

import (
	"fmt"
	"strings"
	"time"
)

// ============================================================================
// SEARCH QUERY TRACKING
// ============================================================================

// SearchQuery represents a tracked search query for analytics
type SearchQuery struct {
	ID                int64      `json:"id" db:"id"`
	QueryText         string     `json:"query_text" db:"query_text"`
	CategoryID        *int64     `json:"category_id,omitempty" db:"category_id"`
	UserID            *int64     `json:"user_id,omitempty" db:"user_id"`
	SessionID         *string    `json:"session_id,omitempty" db:"session_id"`
	ResultsCount      int32      `json:"results_count" db:"results_count"`
	ClickedListingID  *int64     `json:"clicked_listing_id,omitempty" db:"clicked_listing_id"`
	CreatedAt         time.Time  `json:"created_at" db:"created_at"`
}

// CreateSearchQueryInput represents input for logging a search query
type CreateSearchQueryInput struct {
	QueryText        string  `json:"query_text" validate:"required,min=1,max=500"`
	CategoryID       *int64  `json:"category_id,omitempty"`
	UserID           *int64  `json:"user_id,omitempty"`
	SessionID        *string `json:"session_id,omitempty" validate:"omitempty,uuid4"`
	ResultsCount     int32   `json:"results_count" validate:"gte=0"`
	ClickedListingID *int64  `json:"clicked_listing_id,omitempty"`
}

// Validate validates CreateSearchQueryInput
func (input *CreateSearchQueryInput) Validate() error {
	// Sanitize query text
	input.QueryText = strings.TrimSpace(input.QueryText)

	// Check query text length
	if len(input.QueryText) == 0 {
		return fmt.Errorf("query_text cannot be empty")
	}
	if len(input.QueryText) > 500 {
		return fmt.Errorf("query_text cannot exceed 500 characters")
	}

	// Must have either user_id or session_id (or both for migration)
	if input.UserID == nil && input.SessionID == nil {
		return fmt.Errorf("either user_id or session_id must be provided")
	}

	// Validate session_id format if provided (should be UUID)
	if input.SessionID != nil && *input.SessionID != "" {
		sessionID := strings.TrimSpace(*input.SessionID)
		if len(sessionID) == 0 {
			return fmt.Errorf("session_id cannot be empty string")
		}
		// Simple UUID validation (length check)
		if len(sessionID) != 36 {
			return fmt.Errorf("session_id must be a valid UUID (36 chars)")
		}
	}

	// Validate results_count
	if input.ResultsCount < 0 {
		return fmt.Errorf("results_count cannot be negative")
	}

	return nil
}

// ============================================================================
// TRENDING SEARCH QUERIES
// ============================================================================

// TrendingSearch represents a trending search query with aggregated stats
type TrendingSearch struct {
	QueryText    string    `json:"query_text"`
	SearchCount  int64     `json:"search_count"`
	LastSearched time.Time `json:"last_searched"`
	CategoryID   *int64    `json:"category_id,omitempty"` // NULL for global trending
}

// GetTrendingQueriesFilter represents filters for trending queries
type GetTrendingQueriesFilter struct {
	CategoryID       *int64 `json:"category_id,omitempty"`
	Limit            int32  `json:"limit" validate:"required,gte=1,lte=100"`
	DaysAgo          int32  `json:"days_ago" validate:"gte=1,lte=90"` // Time range: last N days
	MinResultsCount  int32  `json:"min_results_count"`                // Filter: only searches with results >= N
	IncludeZeroResults bool `json:"include_zero_results"`             // Include searches with 0 results?
}

// Validate validates GetTrendingQueriesFilter
func (filter *GetTrendingQueriesFilter) Validate() error {
	// Validate limit
	if filter.Limit < 1 || filter.Limit > 100 {
		return fmt.Errorf("limit must be between 1 and 100")
	}

	// Validate days_ago
	if filter.DaysAgo < 1 || filter.DaysAgo > 90 {
		return fmt.Errorf("days_ago must be between 1 and 90")
	}

	// Validate min_results_count
	if filter.MinResultsCount < 0 {
		return fmt.Errorf("min_results_count cannot be negative")
	}

	return nil
}

// DefaultTrendingFilter returns default filter for trending queries
func DefaultTrendingFilter() *GetTrendingQueriesFilter {
	return &GetTrendingQueriesFilter{
		CategoryID:         nil,   // Global trending
		Limit:              10,    // Top 10
		DaysAgo:            7,     // Last 7 days
		MinResultsCount:    1,     // Only successful searches
		IncludeZeroResults: false, // Exclude empty searches
	}
}

// ============================================================================
// USER SEARCH HISTORY
// ============================================================================

// GetUserHistoryFilter represents filters for user search history
type GetUserHistoryFilter struct {
	UserID     *int64  `json:"user_id,omitempty"`
	SessionID  *string `json:"session_id,omitempty" validate:"omitempty,uuid4"`
	CategoryID *int64  `json:"category_id,omitempty"` // Filter by category (optional)
	Limit      int32   `json:"limit" validate:"required,gte=1,lte=100"`
}

// Validate validates GetUserHistoryFilter
func (filter *GetUserHistoryFilter) Validate() error {
	// Must have either user_id or session_id
	if filter.UserID == nil && filter.SessionID == nil {
		return fmt.Errorf("either user_id or session_id must be provided")
	}

	// Validate limit
	if filter.Limit < 1 || filter.Limit > 100 {
		return fmt.Errorf("limit must be between 1 and 100")
	}

	// Validate session_id format if provided
	if filter.SessionID != nil && *filter.SessionID != "" {
		sessionID := strings.TrimSpace(*filter.SessionID)
		if len(sessionID) != 36 {
			return fmt.Errorf("session_id must be a valid UUID")
		}
	}

	return nil
}

// ============================================================================
// POPULAR SEARCHES (All-Time)
// ============================================================================

// GetPopularQueriesFilter represents filters for all-time popular searches
type GetPopularQueriesFilter struct {
	CategoryID      *int64 `json:"category_id,omitempty"`
	Limit           int32  `json:"limit" validate:"required,gte=1,lte=100"`
	MinSearchCount  int64  `json:"min_search_count"` // Filter: only queries with count >= N
}

// Validate validates GetPopularQueriesFilter
func (filter *GetPopularQueriesFilter) Validate() error {
	// Validate limit
	if filter.Limit < 1 || filter.Limit > 100 {
		return fmt.Errorf("limit must be between 1 and 100")
	}

	// Validate min_search_count
	if filter.MinSearchCount < 0 {
		return fmt.Errorf("min_search_count cannot be negative")
	}

	return nil
}

// ============================================================================
// CTR ANALYSIS
// ============================================================================

// SearchQueryCTR represents click-through rate analysis for a query
type SearchQueryCTR struct {
	QueryText     string  `json:"query_text"`
	TotalSearches int64   `json:"total_searches"`
	TotalClicks   int64   `json:"total_clicks"`
	CTRPercent    float64 `json:"ctr_percent"` // (clicks / searches) * 100
	CategoryID    *int64  `json:"category_id,omitempty"`
}

// GetCTRAnalysisFilter represents filters for CTR analysis
type GetCTRAnalysisFilter struct {
	QueryText  *string `json:"query_text,omitempty"`  // Specific query or NULL for all
	CategoryID *int64  `json:"category_id,omitempty"` // Filter by category
	DaysAgo    int32   `json:"days_ago" validate:"gte=1,lte=90"`
	Limit      int32   `json:"limit" validate:"required,gte=1,lte=100"`
}

// Validate validates GetCTRAnalysisFilter
func (filter *GetCTRAnalysisFilter) Validate() error {
	// Validate days_ago
	if filter.DaysAgo < 1 || filter.DaysAgo > 90 {
		return fmt.Errorf("days_ago must be between 1 and 90")
	}

	// Validate limit
	if filter.Limit < 1 || filter.Limit > 100 {
		return fmt.Errorf("limit must be between 1 and 100")
	}

	// Validate query_text if provided
	if filter.QueryText != nil && *filter.QueryText != "" {
		queryText := strings.TrimSpace(*filter.QueryText)
		if len(queryText) == 0 {
			return fmt.Errorf("query_text cannot be empty string")
		}
		if len(queryText) > 500 {
			return fmt.Errorf("query_text cannot exceed 500 characters")
		}
	}

	return nil
}

// ============================================================================
// CONSTANTS
// ============================================================================

const (
	// MaxQueryTextLength is the maximum length for search query text
	MaxQueryTextLength = 500

	// DefaultTrendingDays is the default time range for trending queries
	DefaultTrendingDays = 7

	// DefaultHistoryLimit is the default limit for user history
	DefaultHistoryLimit = 20

	// DefaultTrendingLimit is the default limit for trending queries
	DefaultTrendingLimit = 10

	// RetentionPolicyDays is the number of days to keep search query data
	RetentionPolicyDays = 90
)
