// Package service
// backend/internal/proj/marketplace/service/router_with_timeout.go
package service

import (
	"context"
	"errors"
	"time"

	"github.com/rs/zerolog"

	"backend/internal/config"
	"backend/internal/metrics"
)

// TimeoutRouter wraps TrafficRouter with timeout functionality for microservice calls
//
// Architecture:
// 1. Creates per-request context with deadline
// 2. Tracks timeout events via Prometheus metrics
// 3. Handles fallback to monolith on timeout
// 4. Preserves request ID and metadata in context
//
// Usage:
//
//	router := NewTimeoutRouter(trafficRouter, cfg, logger)
//	result, err := router.ExecuteWithTimeout(ctx, "get", func(ctx context.Context) (interface{}, error) {
//	    return microserviceClient.GetListing(ctx, id)
//	})
type TimeoutRouter struct {
	router *TrafficRouter
	config *config.MarketplaceConfig
	logger zerolog.Logger
}

// NewTimeoutRouter creates a new timeout router
func NewTimeoutRouter(router *TrafficRouter, cfg *config.MarketplaceConfig, logger zerolog.Logger) *TimeoutRouter {
	return &TimeoutRouter{
		router: router,
		config: cfg,
		logger: logger.With().Str("component", "timeout_router").Logger(),
	}
}

// ExecuteWithTimeout executes a microservice operation with timeout and fallback
//
// Parameters:
//   - ctx: parent context (preserves request ID, trace, etc)
//   - operation: operation name for metrics ("get", "search", "create", etc)
//   - microserviceFunc: function that calls microservice (with timeout context)
//   - fallbackFunc: function that calls monolith (with original context, no timeout)
//
// Returns:
//   - result: operation result (from microservice or fallback)
//   - err: error if both microservice and fallback failed
//   - usedFallback: true if fallback was used
//
// Example:
//
//	result, err, usedFallback := router.ExecuteWithTimeout(ctx, "get",
//	    func(ctx context.Context) (interface{}, error) {
//	        return microserviceClient.GetListing(ctx, id)
//	    },
//	    func(ctx context.Context) (interface{}, error) {
//	        return monolithService.GetListing(ctx, id)
//	    },
//	)
func (tr *TimeoutRouter) ExecuteWithTimeout(
	ctx context.Context,
	operation string,
	microserviceFunc func(context.Context) (interface{}, error),
	fallbackFunc func(context.Context) (interface{}, error),
) (interface{}, error, bool) {
	// Create timeout context for microservice call
	timeoutCtx, cancel := context.WithTimeout(ctx, tr.config.MicroserviceTimeout)
	defer cancel()

	// Record start time for latency tracking
	startTime := time.Now()

	// Try microservice with timeout
	result, err := microserviceFunc(timeoutCtx)
	duration := time.Since(startTime)

	// Handle timeout
	if errors.Is(err, context.DeadlineExceeded) {
		tr.logger.Warn().
			Str("operation", operation).
			Dur("timeout", tr.config.MicroserviceTimeout).
			Dur("actual_duration", duration).
			Msg("Microservice request timeout")

		// Record timeout metrics
		metrics.RecordMicroserviceError("timeout")
		metrics.RecordTimeout(operation)

		// Fallback to monolith if enabled
		if tr.config.FallbackToMonolith && fallbackFunc != nil {
			tr.logger.Info().
				Str("operation", operation).
				Msg("Falling back to monolith after timeout")

			metrics.RecordFallback("timeout")
			metrics.RecordTimeoutFallback(operation)

			// Use original context (without timeout) for fallback
			fallbackStartTime := time.Now()
			result, err = fallbackFunc(ctx)
			fallbackDuration := time.Since(fallbackStartTime)

			// Record fallback metrics
			metrics.ObserveRouteDuration("monolith", operation, fallbackDuration.Seconds())

			if err != nil {
				tr.logger.Error().
					Err(err).
					Str("operation", operation).
					Msg("Fallback to monolith also failed")
				return nil, err, true
			}

			tr.logger.Info().
				Str("operation", operation).
				Dur("fallback_duration", fallbackDuration).
				Msg("Successfully fell back to monolith")

			return result, nil, true
		}

		// No fallback - return timeout error
		return nil, err, false
	}

	// Handle other errors (connection, internal, etc)
	if err != nil {
		tr.logger.Error().
			Err(err).
			Str("operation", operation).
			Dur("duration", duration).
			Msg("Microservice request failed")

		// Classify and record error
		errorType := metrics.ClassifyGRPCError(err)
		metrics.RecordMicroserviceError(errorType)

		// Fallback to monolith if enabled
		if tr.config.FallbackToMonolith && fallbackFunc != nil {
			tr.logger.Info().
				Str("operation", operation).
				Str("error_type", errorType).
				Msg("Falling back to monolith after error")

			metrics.RecordFallback("microservice_error")

			fallbackStartTime := time.Now()
			result, err = fallbackFunc(ctx)
			fallbackDuration := time.Since(fallbackStartTime)

			metrics.ObserveRouteDuration("monolith", operation, fallbackDuration.Seconds())

			if err != nil {
				tr.logger.Error().
					Err(err).
					Str("operation", operation).
					Msg("Fallback to monolith also failed")
				return nil, err, true
			}

			tr.logger.Info().
				Str("operation", operation).
				Dur("fallback_duration", fallbackDuration).
				Msg("Successfully fell back to monolith")

			return result, nil, true
		}

		// No fallback - return error
		return nil, err, false
	}

	// Success - record metrics
	metrics.ObserveRouteDuration("microservice", operation, duration.Seconds())

	tr.logger.Debug().
		Str("operation", operation).
		Dur("duration", duration).
		Msg("Microservice request succeeded")

	return result, nil, false
}

// ExecuteWithTimeoutOrMonolith executes operation based on routing decision
//
// This is a higher-level wrapper that combines routing decision + timeout execution.
//
// Parameters:
//   - ctx: parent context
//   - userID: user ID for routing decision
//   - isAdmin: whether user is admin
//   - operation: operation name for metrics
//   - microserviceFunc: function that calls microservice
//   - monolithFunc: function that calls monolith
//
// Returns:
//   - result: operation result
//   - err: error if operation failed
//
// Example:
//
//	result, err := router.ExecuteWithTimeoutOrMonolith(ctx, userID, isAdmin, "get",
//	    func(ctx context.Context) (interface{}, error) {
//	        return microserviceClient.GetListing(ctx, id)
//	    },
//	    func(ctx context.Context) (interface{}, error) {
//	        return monolithService.GetListing(ctx, id)
//	    },
//	)
func (tr *TimeoutRouter) ExecuteWithTimeoutOrMonolith(
	ctx context.Context,
	userID string,
	isAdmin bool,
	operation string,
	microserviceFunc func(context.Context) (interface{}, error),
	monolithFunc func(context.Context) (interface{}, error),
) (interface{}, error) {
	// Make routing decision
	decision := tr.router.ShouldUseMicroservice(userID, isAdmin)

	// If routing to monolith - call monolith directly (no timeout)
	if !decision.UseМicroservice {
		startTime := time.Now()
		result, err := monolithFunc(ctx)
		duration := time.Since(startTime)

		metrics.ObserveRouteDuration("monolith", operation, duration.Seconds())

		if err != nil {
			tr.logger.Error().
				Err(err).
				Str("operation", operation).
				Str("backend", "monolith").
				Msg("Monolith request failed")
		}

		return result, err
	}

	// Routing to microservice - execute with timeout and fallback
	result, err, usedFallback := tr.ExecuteWithTimeout(ctx, operation, microserviceFunc, monolithFunc)

	if usedFallback {
		tr.logger.Warn().
			Str("operation", operation).
			Str("user_id", userID).
			Bool("is_admin", isAdmin).
			Msg("Used fallback to monolith")
	}

	return result, err
}

// GetConfig returns current configuration (for debugging)
func (tr *TimeoutRouter) GetConfig() *config.MarketplaceConfig {
	return tr.config
}

// UpdateConfig updates configuration (for hot reload)
func (tr *TimeoutRouter) UpdateConfig(cfg *config.MarketplaceConfig) {
	tr.config = cfg
	tr.router.UpdateConfig(cfg)

	tr.logger.Info().
		Dur("microservice_timeout", cfg.MicroserviceTimeout).
		Bool("fallback_enabled", cfg.FallbackToMonolith).
		Msg("Timeout router config updated")
}
