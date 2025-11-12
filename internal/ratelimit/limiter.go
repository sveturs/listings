package ratelimit

import (
	"context"
	"time"
)

// RateLimiter defines the interface for rate limiting
type RateLimiter interface {
	// Allow checks if a request is allowed under the rate limit
	// Returns true if allowed, false if rate limit exceeded
	Allow(ctx context.Context, key string, limit int, window time.Duration) (bool, error)

	// Remaining returns the number of remaining requests in the current window
	Remaining(ctx context.Context, key string) (int, error)

	// Reset clears the rate limit for a specific key (useful for testing)
	Reset(ctx context.Context, key string) error
}

// IdentifierType defines how to identify the client for rate limiting
type IdentifierType int

const (
	// ByIP identifies clients by IP address
	ByIP IdentifierType = iota

	// ByUserID identifies clients by user ID
	ByUserID

	// ByIPAndUserID combines IP and user ID for stricter limits
	ByIPAndUserID
)

// String returns the string representation of IdentifierType
func (i IdentifierType) String() string {
	switch i {
	case ByIP:
		return "ip"
	case ByUserID:
		return "user_id"
	case ByIPAndUserID:
		return "ip_and_user_id"
	default:
		return "unknown"
	}
}
