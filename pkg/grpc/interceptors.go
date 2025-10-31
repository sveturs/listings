// Package grpc provides gRPC client utilities for listings service.
// This includes interceptors, connection pooling, and helper functions.
package grpc

import (
	"context"
	"time"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// LoggingInterceptor creates a unary client interceptor that logs all gRPC calls.
// It logs the method name, duration, and any errors that occur.
//
// Example:
//
//	conn, err := grpc.Dial(
//	    addr,
//	    grpc.WithUnaryInterceptor(grpcpkg.LoggingInterceptor(logger)),
//	)
func LoggingInterceptor(logger zerolog.Logger) grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		start := time.Now()

		// Invoke the RPC
		err := invoker(ctx, method, req, reply, cc, opts...)

		// Log the result
		duration := time.Since(start)

		logEvent := logger.Info().
			Str("method", method).
			Dur("duration_ms", duration)

		if err != nil {
			logEvent = logEvent.Err(err).
				Str("level", "error")
		}

		logEvent.Msg("gRPC call")

		return err
	}
}

// MetricsInterceptor creates a unary client interceptor that collects metrics.
// This is a placeholder - integrate with your metrics system (Prometheus, etc.)
//
// Example:
//
//	conn, err := grpc.Dial(
//	    addr,
//	    grpc.WithUnaryInterceptor(grpcpkg.MetricsInterceptor()),
//	)
func MetricsInterceptor() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		start := time.Now()

		// Invoke the RPC
		err := invoker(ctx, method, req, reply, cc, opts...)

		// Record metrics
		duration := time.Since(start)

		// TODO: Integrate with your metrics system
		// Example: prometheus.RecordRPCDuration(method, duration, err)
		_ = duration

		return err
	}
}

// AuthInterceptor creates a unary client interceptor that adds authentication token to metadata.
// This is used for service-to-service authentication.
//
// Example:
//
//	conn, err := grpc.Dial(
//	    addr,
//	    grpc.WithUnaryInterceptor(grpcpkg.AuthInterceptor(token)),
//	)
func AuthInterceptor(token string) grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		// Add auth token to metadata
		if token != "" {
			ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+token)
		}

		// Invoke the RPC
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

// RetryInterceptor creates a unary client interceptor that implements retry logic.
// It retries failed requests with exponential backoff.
//
// Example:
//
//	conn, err := grpc.Dial(
//	    addr,
//	    grpc.WithUnaryInterceptor(grpcpkg.RetryInterceptor(3, 100*time.Millisecond, logger)),
//	)
func RetryInterceptor(maxRetries int, initialDelay time.Duration, logger zerolog.Logger) grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		var err error
		delay := initialDelay

		for attempt := 0; attempt <= maxRetries; attempt++ {
			err = invoker(ctx, method, req, reply, cc, opts...)
			if err == nil {
				return nil
			}

			// Check if error is retryable
			if !isRetryableError(err) {
				return err
			}

			// Check if we should retry
			if attempt == maxRetries {
				logger.Warn().
					Err(err).
					Str("method", method).
					Int("attempts", attempt+1).
					Msg("Max retries exceeded")
				return err
			}

			// Wait before retry with exponential backoff
			logger.Debug().
				Err(err).
				Str("method", method).
				Int("attempt", attempt+1).
				Dur("delay", delay).
				Msg("Retrying gRPC call")

			select {
			case <-time.After(delay):
				delay *= 2 // Exponential backoff
			case <-ctx.Done():
				return ctx.Err()
			}
		}

		return err
	}
}

// isRetryableError checks if an error is retryable.
// Currently, this is a simple check - expand as needed.
func isRetryableError(err error) bool {
	// TODO: Check specific gRPC error codes
	// For now, retry all errors except context errors
	if err == context.Canceled || err == context.DeadlineExceeded {
		return false
	}
	return true
}

// ChainUnaryClient chains multiple unary client interceptors into one.
// This is useful when you want to apply multiple interceptors.
//
// Example:
//
//	interceptor := grpcpkg.ChainUnaryClient(
//	    grpcpkg.LoggingInterceptor(logger),
//	    grpcpkg.MetricsInterceptor(),
//	    grpcpkg.AuthInterceptor(token),
//	)
//	conn, err := grpc.Dial(addr, grpc.WithUnaryInterceptor(interceptor))
func ChainUnaryClient(interceptors ...grpc.UnaryClientInterceptor) grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		// Build a chain of interceptors
		chain := func(currentInvoker grpc.UnaryInvoker, currentInterceptor grpc.UnaryClientInterceptor) grpc.UnaryInvoker {
			return func(
				currentCtx context.Context,
				currentMethod string,
				currentReq, currentReply interface{},
				currentConn *grpc.ClientConn,
				currentOpts ...grpc.CallOption,
			) error {
				return currentInterceptor(
					currentCtx,
					currentMethod,
					currentReq,
					currentReply,
					currentConn,
					currentInvoker,
					currentOpts...,
				)
			}
		}

		chainedInvoker := invoker
		for i := len(interceptors) - 1; i >= 0; i-- {
			chainedInvoker = chain(chainedInvoker, interceptors[i])
		}

		return chainedInvoker(ctx, method, req, reply, cc, opts...)
	}
}

// TimeoutInterceptor creates a unary client interceptor that enforces a timeout.
// This wraps each RPC call with a context timeout.
//
// Example:
//
//	conn, err := grpc.Dial(
//	    addr,
//	    grpc.WithUnaryInterceptor(grpcpkg.TimeoutInterceptor(5*time.Second)),
//	)
func TimeoutInterceptor(timeout time.Duration) grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		// Create context with timeout
		timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		// Invoke with timeout context
		return invoker(timeoutCtx, method, req, reply, cc, opts...)
	}
}
