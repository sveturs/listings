package grpc

import (
	"context"
	"runtime/debug"
	"time"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/vondi-global/listings/internal/metrics"
)

// OrderServiceMiddleware provides gRPC interceptors for order service
type OrderServiceMiddleware struct {
	logger  zerolog.Logger
	metrics *metrics.Metrics
}

// NewOrderServiceMiddleware creates a new middleware instance
func NewOrderServiceMiddleware(logger zerolog.Logger, m *metrics.Metrics) *OrderServiceMiddleware {
	return &OrderServiceMiddleware{
		logger:  logger.With().Str("component", "grpc_orders_middleware").Logger(),
		metrics: m,
	}
}

// ============================================================================
// UNARY INTERCEPTORS
// ============================================================================

// AuthInterceptor extracts user_id from gRPC metadata and adds it to context
// This allows downstream handlers to access authenticated user information
func (m *OrderServiceMiddleware) AuthInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// Extract metadata
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			m.logger.Debug().Str("method", info.FullMethod).Msg("no metadata in request")
			// Continue without auth - handlers will validate if needed
			return handler(ctx, req)
		}

		// Extract user_id from metadata (if present)
		userIDs := md.Get("user_id")
		if len(userIDs) > 0 {
			// Add user_id to context for handlers to use
			ctx = context.WithValue(ctx, "user_id", userIDs[0])
			m.logger.Debug().
				Str("method", info.FullMethod).
				Str("user_id", userIDs[0]).
				Msg("user authenticated")
		}

		// Extract session_id from metadata (for anonymous carts)
		sessionIDs := md.Get("session_id")
		if len(sessionIDs) > 0 {
			ctx = context.WithValue(ctx, "session_id", sessionIDs[0])
			m.logger.Debug().
				Str("method", info.FullMethod).
				Str("session_id", sessionIDs[0]).
				Msg("session identified")
		}

		return handler(ctx, req)
	}
}

// LoggingInterceptor logs all RPC calls with timing information
func (m *OrderServiceMiddleware) LoggingInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		startTime := time.Now()

		m.logger.Info().
			Str("method", info.FullMethod).
			Msg("gRPC request started")

		// Call handler
		resp, err := handler(ctx, req)

		// Calculate duration
		duration := time.Since(startTime)

		// Log result
		if err != nil {
			// Extract gRPC status code
			st, ok := status.FromError(err)
			code := codes.Unknown
			if ok {
				code = st.Code()
			}

			m.logger.Error().
				Err(err).
				Str("method", info.FullMethod).
				Str("code", code.String()).
				Dur("duration_ms", duration).
				Msg("gRPC request failed")
		} else {
			m.logger.Info().
				Str("method", info.FullMethod).
				Dur("duration_ms", duration).
				Msg("gRPC request completed")
		}

		return resp, err
	}
}

// RecoveryInterceptor catches panics and converts them to gRPC Internal errors
func (m *OrderServiceMiddleware) RecoveryInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		// Recover from panic
		defer func() {
			if r := recover(); r != nil {
				m.logger.Error().
					Str("method", info.FullMethod).
					Interface("panic", r).
					Str("stack", string(debug.Stack())).
					Msg("panic recovered in gRPC handler")

				// Return internal error
				err = status.Errorf(codes.Internal, "internal server error: %v", r)
			}
		}()

		// Call handler
		return handler(ctx, req)
	}
}

// MetricsInterceptor records metrics for all RPC calls
func (m *OrderServiceMiddleware) MetricsInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		startTime := time.Now()

		// Call handler
		resp, err := handler(ctx, req)

		// Calculate duration
		duration := time.Since(startTime)

		// Extract status code
		code := codes.OK
		if err != nil {
			st, ok := status.FromError(err)
			if ok {
				code = st.Code()
			} else {
				code = codes.Unknown
			}
		}

		// Record metrics (if metrics instance exists)
		if m.metrics != nil {
			// Extract method name from full method (e.g., "/OrderService/CreateOrder" -> "CreateOrder")
			method := info.FullMethod
			if len(method) > 0 && method[0] == '/' {
				parts := splitFullMethod(method)
				if len(parts) == 2 {
					method = parts[1]
				}
			}

			// Record gRPC request
			m.metrics.RecordGRPCRequest(method, code.String(), duration.Seconds())

			// Increment request counter
			if code == codes.OK {
				m.logger.Debug().
					Str("method", method).
					Float64("duration_seconds", duration.Seconds()).
					Msg("metrics recorded")
			} else {
				m.logger.Debug().
					Str("method", method).
					Str("code", code.String()).
					Float64("duration_seconds", duration.Seconds()).
					Msg("metrics recorded (error)")
			}
		}

		return resp, err
	}
}

// ============================================================================
// CHAIN INTERCEPTORS
// ============================================================================

// ChainUnaryInterceptors chains multiple unary interceptors into one
// Interceptors are executed in order: first -> last
func ChainUnaryInterceptors(interceptors ...grpc.UnaryServerInterceptor) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// Build chain from last to first
		chain := handler
		for i := len(interceptors) - 1; i >= 0; i-- {
			interceptor := interceptors[i]
			next := chain
			chain = func(currentCtx context.Context, currentReq interface{}) (interface{}, error) {
				return interceptor(currentCtx, currentReq, info, next)
			}
		}

		return chain(ctx, req)
	}
}

// GetDefaultInterceptors returns default interceptor chain for order service
// Order: Recovery -> Logging -> Auth -> Metrics -> Handler
func (m *OrderServiceMiddleware) GetDefaultInterceptors() grpc.UnaryServerInterceptor {
	return ChainUnaryInterceptors(
		m.RecoveryInterceptor(), // 1. Catch panics first
		m.LoggingInterceptor(),  // 2. Log all requests
		m.AuthInterceptor(),     // 3. Extract auth info
		m.MetricsInterceptor(),  // 4. Record metrics
	)
}

// ============================================================================
// HELPER FUNCTIONS
// ============================================================================

// splitFullMethod splits a full gRPC method name into service and method
// Example: "/OrderService/CreateOrder" -> ["OrderService", "CreateOrder"]
func splitFullMethod(fullMethod string) []string {
	// Remove leading slash
	if len(fullMethod) > 0 && fullMethod[0] == '/' {
		fullMethod = fullMethod[1:]
	}

	// Split by slash
	parts := make([]string, 0, 2)
	slashIndex := -1
	for i, ch := range fullMethod {
		if ch == '/' {
			slashIndex = i
			break
		}
	}

	if slashIndex > 0 {
		parts = append(parts, fullMethod[:slashIndex])
		parts = append(parts, fullMethod[slashIndex+1:])
	} else {
		parts = append(parts, fullMethod)
	}

	return parts
}

// ============================================================================
// STREAM INTERCEPTORS (for future use)
// ============================================================================

// LoggingStreamInterceptor logs all streaming RPC calls
func (m *OrderServiceMiddleware) LoggingStreamInterceptor() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		startTime := time.Now()

		m.logger.Info().
			Str("method", info.FullMethod).
			Msg("gRPC stream started")

		// Call handler
		err := handler(srv, ss)

		// Calculate duration
		duration := time.Since(startTime)

		// Log result
		if err != nil {
			m.logger.Error().
				Err(err).
				Str("method", info.FullMethod).
				Dur("duration_ms", duration).
				Msg("gRPC stream failed")
		} else {
			m.logger.Info().
				Str("method", info.FullMethod).
				Dur("duration_ms", duration).
				Msg("gRPC stream completed")
		}

		return err
	}
}

// RecoveryStreamInterceptor catches panics in streaming RPCs
func (m *OrderServiceMiddleware) RecoveryStreamInterceptor() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) (err error) {
		// Recover from panic
		defer func() {
			if r := recover(); r != nil {
				m.logger.Error().
					Str("method", info.FullMethod).
					Interface("panic", r).
					Str("stack", string(debug.Stack())).
					Msg("panic recovered in gRPC stream handler")

				// Return internal error
				err = status.Errorf(codes.Internal, "internal server error: %v", r)
			}
		}()

		// Call handler
		return handler(srv, ss)
	}
}
