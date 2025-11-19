package cache

import "time"

// SearchCacheConfig holds cache TTL configurations for different search types
type SearchCacheConfig struct {
	// SearchTTL is the TTL for basic search results (Phase 21.1)
	// Default: 5 minutes - search results change frequently
	SearchTTL time.Duration

	// FacetsTTL is the TTL for facets/aggregations (Phase 21.2)
	// Default: 5 minutes - facet counts change as listings are added/removed
	FacetsTTL time.Duration

	// SuggestionsTTL is the TTL for autocomplete suggestions (Phase 21.2)
	// Default: 1 hour - suggestions are relatively stable
	SuggestionsTTL time.Duration

	// PopularTTL is the TTL for popular searches (Phase 21.2)
	// Default: 15 minutes - trending searches need frequent updates
	PopularTTL time.Duration

	// FilteredSearchTTL is the TTL for filtered/sorted search results (Phase 21.2)
	// Default: 5 minutes - same as basic search
	FilteredSearchTTL time.Duration

	// HistoryTTL is the TTL for user search history (Phase 28)
	// Default: 5 minutes - personal history changes with each search
	HistoryTTL time.Duration
}

// DefaultSearchCacheConfig returns default cache configuration
func DefaultSearchCacheConfig() SearchCacheConfig {
	return SearchCacheConfig{
		SearchTTL:         5 * time.Minute,
		FacetsTTL:         5 * time.Minute,
		SuggestionsTTL:    1 * time.Hour,
		PopularTTL:        15 * time.Minute,
		FilteredSearchTTL: 5 * time.Minute,
		HistoryTTL:        5 * time.Minute,
	}
}

// Validate ensures all TTL values are reasonable
func (c *SearchCacheConfig) Validate() error {
	// Minimum TTL is 1 minute
	minTTL := 1 * time.Minute
	// Maximum TTL is 24 hours
	maxTTL := 24 * time.Hour

	ttls := map[string]time.Duration{
		"SearchTTL":         c.SearchTTL,
		"FacetsTTL":         c.FacetsTTL,
		"SuggestionsTTL":    c.SuggestionsTTL,
		"PopularTTL":        c.PopularTTL,
		"FilteredSearchTTL": c.FilteredSearchTTL,
		"HistoryTTL":        c.HistoryTTL,
	}

	for name, ttl := range ttls {
		if ttl < minTTL {
			// Auto-fix to default instead of error
			switch name {
			case "SearchTTL":
				c.SearchTTL = 5 * time.Minute
			case "FacetsTTL":
				c.FacetsTTL = 5 * time.Minute
			case "SuggestionsTTL":
				c.SuggestionsTTL = 1 * time.Hour
			case "PopularTTL":
				c.PopularTTL = 15 * time.Minute
			case "FilteredSearchTTL":
				c.FilteredSearchTTL = 5 * time.Minute
			case "HistoryTTL":
				c.HistoryTTL = 5 * time.Minute
			}
		}
		if ttl > maxTTL {
			// Cap at maximum
			switch name {
			case "SearchTTL":
				c.SearchTTL = maxTTL
			case "FacetsTTL":
				c.FacetsTTL = maxTTL
			case "SuggestionsTTL":
				c.SuggestionsTTL = maxTTL
			case "PopularTTL":
				c.PopularTTL = maxTTL
			case "FilteredSearchTTL":
				c.FilteredSearchTTL = maxTTL
			case "HistoryTTL":
				c.HistoryTTL = maxTTL
			}
		}
	}

	return nil
}
