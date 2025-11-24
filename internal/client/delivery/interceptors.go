package delivery

import (
	"context"
	"time"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// loggingInterceptor creates a unary client interceptor that logs all gRPC calls
func loggingInterceptor(logger zerolog.Logger) grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		start := time.Now()

		err := invoker(ctx, method, req, reply, cc, opts...)

		duration := time.Since(start)

		logEvent := logger.Debug().
			Str("service", "delivery").
			Str("method", method).
			Dur("duration", duration)

		if err != nil {
			st, _ := status.FromError(err)
			logEvent.
				Str("status", st.Code().String()).
				Str("error", st.Message()).
				Msg("gRPC call failed")
		} else {
			logEvent.Msg("gRPC call completed")
		}

		return err
	}
}

// retryInterceptor creates a unary client interceptor that retries on transient errors
func retryInterceptor(maxRetries int, initialDelay time.Duration, logger zerolog.Logger) grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		var lastErr error
		delay := initialDelay

		for attempt := 0; attempt <= maxRetries; attempt++ {
			lastErr = invoker(ctx, method, req, reply, cc, opts...)
			if lastErr == nil {
				return nil
			}

			// Check if error is retryable
			if !isRetryableError(lastErr) {
				return lastErr
			}

			// Don't retry if this was the last attempt
			if attempt == maxRetries {
				break
			}

			logger.Debug().
				Err(lastErr).
				Str("method", method).
				Int("attempt", attempt+1).
				Int("max_retries", maxRetries).
				Dur("delay", delay).
				Msg("Retrying delivery gRPC call")

			// Wait before retrying
			select {
			case <-time.After(delay):
				delay *= 2 // Exponential backoff
				if delay > 5*time.Second {
					delay = 5 * time.Second // Cap at 5 seconds
				}
			case <-ctx.Done():
				return ctx.Err()
			}
		}

		logger.Warn().
			Err(lastErr).
			Str("method", method).
			Int("attempts", maxRetries+1).
			Msg("All retry attempts exhausted")

		return lastErr
	}
}

// isRetryableError determines if a gRPC error should be retried
func isRetryableError(err error) bool {
	if err == nil {
		return false
	}

	// Context errors are not retryable
	if err == context.Canceled || err == context.DeadlineExceeded {
		return false
	}

	st, ok := status.FromError(err)
	if !ok {
		// Non-gRPC errors are generally retryable (network issues)
		return true
	}

	switch st.Code() {
	case codes.Unavailable:
		// Service temporarily unavailable - retry
		return true
	case codes.ResourceExhausted:
		// Rate limited - retry with backoff
		return true
	case codes.Aborted:
		// Operation aborted, can retry
		return true
	case codes.DeadlineExceeded:
		// Timeout - context was exhausted, don't retry
		return false
	case codes.Internal:
		// Internal server error - might be transient
		return true
	case codes.Unknown:
		// Unknown errors might be transient
		return true
	default:
		// Other errors (InvalidArgument, NotFound, etc.) are not retryable
		return false
	}
}
