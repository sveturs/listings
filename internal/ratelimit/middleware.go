package ratelimit

import (
	"context"
	"fmt"
	"strings"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

// MetricsRecorder is an interface for recording rate limit metrics
type MetricsRecorder interface {
	RecordRateLimitEvaluation(method, identifierType string, allowed bool)
}

// UnaryServerInterceptor creates a gRPC unary interceptor for rate limiting
func UnaryServerInterceptor(limiter RateLimiter, config *Config, logger zerolog.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// Get endpoint configuration
		endpointConfig := config.GetEndpointConfig(info.FullMethod)

		// Skip if rate limiting is disabled for this endpoint
		if !endpointConfig.Enabled {
			return handler(ctx, req)
		}

		// Extract identifier based on configuration
		identifier, err := extractIdentifier(ctx, endpointConfig.Identifier)
		if err != nil {
			logger.Error().
				Err(err).
				Str("method", info.FullMethod).
				Msg("failed to extract identifier for rate limiting, allowing request")
			// Fail open - allow the request if we can't extract identifier
			return handler(ctx, req)
		}

		// Build rate limit key: method:identifier
		key := buildRateLimitKey(info.FullMethod, identifier)

		// Check rate limit
		allowed, err := limiter.Allow(ctx, key, endpointConfig.Limit, endpointConfig.Window)
		if err != nil {
			logger.Error().
				Err(err).
				Str("method", info.FullMethod).
				Str("identifier", identifier).
				Msg("rate limit check failed, allowing request")
			// Fail open - allow the request on error
			return handler(ctx, req)
		}

		// If rate limit exceeded, return error
		if !allowed {
			logger.Warn().
				Str("method", info.FullMethod).
				Str("identifier", identifier).
				Int("limit", endpointConfig.Limit).
				Dur("window", endpointConfig.Window).
				Msg("rate limit exceeded")

			return nil, status.Errorf(
				codes.ResourceExhausted,
				"rate limit exceeded: maximum %d requests per %s",
				endpointConfig.Limit,
				endpointConfig.Window,
			)
		}

		logger.Debug().
			Str("method", info.FullMethod).
			Str("identifier", identifier).
			Msg("rate limit check passed")

		// Continue with the request
		return handler(ctx, req)
	}
}

// StreamServerInterceptor creates a gRPC stream interceptor for rate limiting
func StreamServerInterceptor(limiter RateLimiter, config *Config, logger zerolog.Logger) grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		// Get endpoint configuration
		endpointConfig := config.GetEndpointConfig(info.FullMethod)

		// Skip if rate limiting is disabled for this endpoint
		if !endpointConfig.Enabled {
			return handler(srv, ss)
		}

		// Extract identifier based on configuration
		identifier, err := extractIdentifier(ss.Context(), endpointConfig.Identifier)
		if err != nil {
			logger.Error().
				Err(err).
				Str("method", info.FullMethod).
				Msg("failed to extract identifier for rate limiting, allowing request")
			// Fail open
			return handler(srv, ss)
		}

		// Build rate limit key
		key := buildRateLimitKey(info.FullMethod, identifier)

		// Check rate limit
		allowed, err := limiter.Allow(ss.Context(), key, endpointConfig.Limit, endpointConfig.Window)
		if err != nil {
			logger.Error().
				Err(err).
				Str("method", info.FullMethod).
				Str("identifier", identifier).
				Msg("rate limit check failed, allowing request")
			// Fail open
			return handler(srv, ss)
		}

		// If rate limit exceeded, return error
		if !allowed {
			logger.Warn().
				Str("method", info.FullMethod).
				Str("identifier", identifier).
				Int("limit", endpointConfig.Limit).
				Dur("window", endpointConfig.Window).
				Msg("rate limit exceeded")

			return status.Errorf(
				codes.ResourceExhausted,
				"rate limit exceeded: maximum %d requests per %s",
				endpointConfig.Limit,
				endpointConfig.Window,
			)
		}

		logger.Debug().
			Str("method", info.FullMethod).
			Str("identifier", identifier).
			Msg("rate limit check passed")

		// Continue with the stream
		return handler(srv, ss)
	}
}

// extractIdentifier extracts the client identifier from context based on identifier type
func extractIdentifier(ctx context.Context, identifierType IdentifierType) (string, error) {
	switch identifierType {
	case ByIP:
		return getClientIP(ctx)
	case ByUserID:
		return getUserID(ctx)
	case ByIPAndUserID:
		ip, err := getClientIP(ctx)
		if err != nil {
			return "", err
		}
		userID, err := getUserID(ctx)
		if err != nil {
			// If user ID is not available, fall back to IP only
			return ip, nil
		}
		return fmt.Sprintf("%s:%s", ip, userID), nil
	default:
		return "", fmt.Errorf("unknown identifier type: %d", identifierType)
	}
}

// getClientIP extracts the client IP from gRPC context
func getClientIP(ctx context.Context) (string, error) {
	// First, try to get from metadata (set by proxy/load balancer)
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		// Check X-Forwarded-For header
		if xff := md.Get("x-forwarded-for"); len(xff) > 0 {
			// X-Forwarded-For can contain multiple IPs, use the first one
			ips := strings.Split(xff[0], ",")
			if len(ips) > 0 {
				return strings.TrimSpace(ips[0]), nil
			}
		}

		// Check X-Real-IP header
		if xri := md.Get("x-real-ip"); len(xri) > 0 {
			return xri[0], nil
		}
	}

	// Fallback to peer address
	p, ok := peer.FromContext(ctx)
	if !ok {
		return "", fmt.Errorf("no peer info in context")
	}

	return p.Addr.String(), nil
}

// getUserID extracts the user ID from gRPC context metadata
// This assumes the caller sets "user-id" metadata
func getUserID(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", fmt.Errorf("no metadata in context")
	}

	userIDs := md.Get("user-id")
	if len(userIDs) == 0 {
		return "", fmt.Errorf("no user-id in metadata")
	}

	return userIDs[0], nil
}

// buildRateLimitKey builds a unique key for rate limiting
func buildRateLimitKey(method, identifier string) string {
	// Clean method name (remove leading slash)
	method = strings.TrimPrefix(method, "/")
	// Replace slashes with colons for better Redis key structure
	method = strings.ReplaceAll(method, "/", ":")
	return fmt.Sprintf("%s:%s", method, identifier)
}

// UnaryServerInterceptorWithMetrics creates a gRPC unary interceptor for rate limiting with metrics support
func UnaryServerInterceptorWithMetrics(limiter RateLimiter, config *Config, metrics MetricsRecorder, logger zerolog.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// Get endpoint configuration
		endpointConfig := config.GetEndpointConfig(info.FullMethod)

		// Skip if rate limiting is disabled for this endpoint
		if !endpointConfig.Enabled {
			return handler(ctx, req)
		}

		// Extract identifier based on configuration
		identifier, err := extractIdentifier(ctx, endpointConfig.Identifier)
		if err != nil {
			logger.Error().
				Err(err).
				Str("method", info.FullMethod).
				Msg("failed to extract identifier for rate limiting, allowing request")
			// Fail open - allow the request if we can't extract identifier
			return handler(ctx, req)
		}

		// Build rate limit key: method:identifier
		key := buildRateLimitKey(info.FullMethod, identifier)

		// Check rate limit
		allowed, err := limiter.Allow(ctx, key, endpointConfig.Limit, endpointConfig.Window)
		if err != nil {
			logger.Error().
				Err(err).
				Str("method", info.FullMethod).
				Str("identifier", identifier).
				Msg("rate limit check failed, allowing request")
			// Fail open - allow the request on error
			return handler(ctx, req)
		}

		// Record metrics
		if metrics != nil {
			metrics.RecordRateLimitEvaluation(info.FullMethod, endpointConfig.Identifier.String(), allowed)
		}

		// If rate limit exceeded, return error
		if !allowed {
			logger.Warn().
				Str("method", info.FullMethod).
				Str("identifier", identifier).
				Int("limit", endpointConfig.Limit).
				Dur("window", endpointConfig.Window).
				Msg("rate limit exceeded")

			return nil, status.Errorf(
				codes.ResourceExhausted,
				"rate limit exceeded: maximum %d requests per %s",
				endpointConfig.Limit,
				endpointConfig.Window,
			)
		}

		logger.Debug().
			Str("method", info.FullMethod).
			Str("identifier", identifier).
			Msg("rate limit check passed")

		// Continue with the request
		return handler(ctx, req)
	}
}
