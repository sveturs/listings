package search

import "errors"

var (
	// ErrQueryTooLong is returned when search query exceeds maximum length
	ErrQueryTooLong = errors.New("search query too long (max 500 characters)")

	// ErrInvalidLimit is returned when limit is out of range
	ErrInvalidLimit = errors.New("limit must be between 1 and 100")

	// ErrInvalidOffset is returned when offset is negative
	ErrInvalidOffset = errors.New("offset must be >= 0")

	// ErrSearchFailed is returned when OpenSearch query fails
	ErrSearchFailed = errors.New("search query failed")

	// ErrCacheUnavailable is returned when cache is unavailable (non-critical)
	ErrCacheUnavailable = errors.New("cache unavailable")

	// Phase 21.2: Advanced Search Errors

	// ErrInvalidPriceRange is returned when price range is invalid
	ErrInvalidPriceRange = errors.New("invalid price range: min and max must be >= 0, min must be <= max")

	// ErrInvalidLatitude is returned when latitude is out of range
	ErrInvalidLatitude = errors.New("invalid latitude: must be between -90 and 90")

	// ErrInvalidLongitude is returned when longitude is out of range
	ErrInvalidLongitude = errors.New("invalid longitude: must be between -180 and 180")

	// ErrInvalidRadius is returned when radius is out of range
	ErrInvalidRadius = errors.New("invalid radius: must be between 0 and 1000 km")

	// ErrInvalidSourceType is returned when source type is invalid
	ErrInvalidSourceType = errors.New("invalid source type: must be 'c2c' or 'b2c'")

	// ErrInvalidStockStatus is returned when stock status is invalid
	ErrInvalidStockStatus = errors.New("invalid stock status: must be 'in_stock', 'out_of_stock', or 'low_stock'")

	// ErrTooManyAttributeFilters is returned when too many attribute filters
	ErrTooManyAttributeFilters = errors.New("too many attribute filters: maximum 10 allowed")

	// ErrEmptyAttributeKey is returned when attribute key is empty
	ErrEmptyAttributeKey = errors.New("attribute key cannot be empty")

	// ErrEmptyAttributeValues is returned when attribute values list is empty
	ErrEmptyAttributeValues = errors.New("attribute values list cannot be empty")

	// ErrTooManyAttributeValues is returned when too many values for attribute
	ErrTooManyAttributeValues = errors.New("too many values for attribute: maximum 20 allowed per attribute")

	// ErrInvalidSortField is returned when sort field is invalid
	ErrInvalidSortField = errors.New("invalid sort field: must be 'price', 'created_at', 'relevance', 'views_count', or 'favorites_count'")

	// ErrInvalidSortOrder is returned when sort order is invalid
	ErrInvalidSortOrder = errors.New("invalid sort order: must be 'asc' or 'desc'")

	// ErrPrefixTooShort is returned when autocomplete prefix is too short
	ErrPrefixTooShort = errors.New("prefix too short: minimum 2 characters required")

	// ErrPrefixTooLong is returned when autocomplete prefix is too long
	ErrPrefixTooLong = errors.New("prefix too long: maximum 100 characters allowed")

	// ErrInvalidTimeRange is returned when time range is invalid
	ErrInvalidTimeRange = errors.New("invalid time range: must be '24h', '7d', or '30d'")
)
