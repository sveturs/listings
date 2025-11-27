// Package middleware provides reusable Fiber middleware for listings service integration.
// This package can be imported by other services in the ecosystem.
package middleware

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"

	"github.com/vondi-global/listings/pkg/service"
)

const (
	// LocalsKeyListingsClient is the key for storing listings client in fiber context
	LocalsKeyListingsClient = "listingsClient"

	// LocalsKeyListing is the key for storing a listing in fiber context
	LocalsKeyListing = "listing"
)

// InjectListingsClient creates a middleware that injects the listings client into the context.
// This allows handlers to access the client via c.Locals(LocalsKeyListingsClient).
//
// Example:
//
//	app.Use(middleware.InjectListingsClient(client))
//	app.Get("/listings/:id", func(c *fiber.Ctx) error {
//	    client := c.Locals(middleware.LocalsKeyListingsClient).(*service.Client)
//	    // use client...
//	})
func InjectListingsClient(client *service.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Locals(LocalsKeyListingsClient, client)
		return c.Next()
	}
}

// GetListingsClient retrieves the listings client from fiber context.
// Returns nil if the client is not found.
func GetListingsClient(c *fiber.Ctx) *service.Client {
	client, ok := c.Locals(LocalsKeyListingsClient).(*service.Client)
	if !ok {
		return nil
	}
	return client
}

// RequireListingOwnership creates a middleware that verifies the authenticated user owns the listing.
// This middleware requires:
//   - InjectListingsClient middleware to be registered first
//   - User authentication (userID must be available in context)
//   - Listing ID in route parameter (default: "id", can be customized with paramName)
//
// Example:
//
//	// Assuming auth middleware sets "userID" in locals
//	app.Put("/listings/:id",
//	    authMiddleware.RequireAuth(),
//	    listingsMiddleware.RequireListingOwnership(client),
//	    handler.UpdateListing,
//	)
func RequireListingOwnership(client *service.Client, paramName ...string) fiber.Handler {
	idParam := "id"
	if len(paramName) > 0 {
		idParam = paramName[0]
	}

	return func(c *fiber.Ctx) error {
		// Get authenticated user ID from context
		userID, ok := c.Locals("userID").(int64)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "listings.unauthorized",
			})
		}

		// Parse listing ID from route parameter
		listingIDStr := c.Params(idParam)
		listingID, err := strconv.ParseInt(listingIDStr, 10, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "listings.invalid_id",
			})
		}

		// Fetch listing to verify ownership
		listing, err := client.GetListing(c.Context(), listingID)
		if err != nil {
			if err == service.ErrNotFound {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"error": "listings.not_found",
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "listings.fetch_failed",
			})
		}

		// Verify ownership
		if listing.UserID != userID {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "listings.not_owner",
			})
		}

		// Store listing in context for handler use
		c.Locals(LocalsKeyListing, listing)

		return c.Next()
	}
}

// GetListing retrieves the listing from fiber context (set by RequireListingOwnership).
// Returns nil if the listing is not found.
func GetListing(c *fiber.Ctx) *service.Listing {
	listing, ok := c.Locals(LocalsKeyListing).(*service.Listing)
	if !ok {
		return nil
	}
	return listing
}

// CacheListingsConfig holds configuration for listings caching middleware.
type CacheListingsConfig struct {
	// TTL is the time-to-live for cached listings
	TTL time.Duration

	// Logger for structured logging
	Logger zerolog.Logger

	// KeyGenerator generates cache keys (optional, defaults to request path)
	KeyGenerator func(c *fiber.Ctx) string
}

// CacheListings creates a middleware that caches listing responses.
// This is useful for read-heavy endpoints like GetListing or SearchListings.
//
// Note: This is a simple in-memory cache. For production, consider using Redis.
//
// Example:
//
//	app.Get("/listings/:id",
//	    middleware.CacheListings(middleware.CacheListingsConfig{
//	        TTL:    5 * time.Minute,
//	        Logger: logger,
//	    }),
//	    handler.GetListing,
//	)
func CacheListings(config CacheListingsConfig) fiber.Handler {
	if config.TTL == 0 {
		config.TTL = 5 * time.Minute
	}

	if config.KeyGenerator == nil {
		config.KeyGenerator = func(c *fiber.Ctx) string {
			return c.Path()
		}
	}

	// Simple in-memory cache (not recommended for production)
	cache := make(map[string]cacheEntry)

	return func(c *fiber.Ctx) error {
		// Only cache GET requests
		if c.Method() != fiber.MethodGet {
			return c.Next()
		}

		key := config.KeyGenerator(c)

		// Check cache
		if entry, found := cache[key]; found {
			if time.Now().Before(entry.expiresAt) {
				config.Logger.Debug().
					Str("key", key).
					Msg("Cache hit")

				c.Set("X-Cache", "HIT")
				return c.JSON(entry.data)
			}
			// Expired, remove from cache
			delete(cache, key)
		}

		// Cache miss, continue to handler
		config.Logger.Debug().
			Str("key", key).
			Msg("Cache miss")

		c.Set("X-Cache", "MISS")

		// TODO: Implement cache storage after response
		// This requires capturing the response body, which is non-trivial in Fiber
		// For now, this is a placeholder

		return c.Next()
	}
}

// cacheEntry represents a cached response
type cacheEntry struct {
	data      interface{}
	expiresAt time.Time
}

// RateLimitByUserID creates a middleware that rate limits requests per user.
// This is useful for preventing abuse of write operations.
//
// Example:
//
//	app.Post("/listings",
//	    authMiddleware.RequireAuth(),
//	    middleware.RateLimitByUserID(10, time.Minute), // 10 requests per minute
//	    handler.CreateListing,
//	)
func RateLimitByUserID(maxRequests int, window time.Duration) fiber.Handler {
	// Simple in-memory rate limiter (not recommended for production)
	// For production, use Redis-based rate limiting
	rateLimits := make(map[int64]*rateLimitEntry)

	return func(c *fiber.Ctx) error {
		// Get authenticated user ID from context
		userID, ok := c.Locals("userID").(int64)
		if !ok {
			// No user ID, skip rate limiting
			return c.Next()
		}

		now := time.Now()

		// Get or create rate limit entry
		entry, found := rateLimits[userID]
		if !found || now.After(entry.resetAt) {
			// Create new entry
			rateLimits[userID] = &rateLimitEntry{
				count:   1,
				resetAt: now.Add(window),
			}
			return c.Next()
		}

		// Check if user exceeded rate limit
		if entry.count >= maxRequests {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error":       "listings.rate_limit_exceeded",
				"reset_at":    entry.resetAt.Unix(),
				"retry_after": int(time.Until(entry.resetAt).Seconds()),
			})
		}

		// Increment counter
		entry.count++

		return c.Next()
	}
}

// rateLimitEntry represents rate limit state for a user
type rateLimitEntry struct {
	count   int
	resetAt time.Time
}

// LogListingOperations creates a middleware that logs listing operations.
// This is useful for audit trails and debugging.
//
// Example:
//
//	app.Use("/listings", middleware.LogListingOperations(logger))
func LogListingOperations(logger zerolog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Process request
		err := c.Next()

		// Log after request completes
		duration := time.Since(start)

		logEvent := logger.Info().
			Str("method", c.Method()).
			Str("path", c.Path()).
			Int("status", c.Response().StatusCode()).
			Dur("duration_ms", duration)

		// Add user ID if available
		if userID, ok := c.Locals("userID").(int64); ok {
			logEvent = logEvent.Int64("user_id", userID)
		}

		// Add listing ID if available
		if listingIDStr := c.Params("id"); listingIDStr != "" {
			if listingID, parseErr := strconv.ParseInt(listingIDStr, 10, 64); parseErr == nil {
				logEvent = logEvent.Int64("listing_id", listingID)
			}
		}

		if err != nil {
			logEvent = logEvent.Err(err)
		}

		logEvent.Msg("Listing operation")

		return err
	}
}
