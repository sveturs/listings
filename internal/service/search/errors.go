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
)
