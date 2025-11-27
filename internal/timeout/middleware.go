package timeout

import (
	"context"
	"time"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/vondi-global/listings/internal/metrics"
)

// UnaryServerInterceptor creates a gRPC interceptor that enforces timeouts on handlers.
// It wraps each request with a timeout context based on the endpoint configuration.
func UnaryServerInterceptor(m *metrics.Metrics, logger zerolog.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// Get timeout for this method
		timeout := GetTimeout(info.FullMethod)

		// Create context with timeout
		ctx, cancel := WithTimeout(ctx, timeout)
		defer cancel()

		// Track start time
		start := time.Now()

		// Call handler with timeout context
		resp, err := handler(ctx, req)

		// Calculate duration
		duration := time.Since(start)

		// Check if timeout was exceeded
		if IsDeadlineExceeded(err) || IsDeadlineExceeded(ctx.Err()) {
			logger.Warn().
				Str("method", info.FullMethod).
				Dur("timeout", timeout).
				Dur("elapsed", duration).
				Msg("Request timed out")

			// Increment timeout metric
			m.TimeoutsTotal.WithLabelValues(info.FullMethod).Inc()
			m.TimeoutDuration.WithLabelValues(info.FullMethod).Observe(duration.Seconds())

			return nil, status.Errorf(
				codes.DeadlineExceeded,
				"request timeout after %v (limit: %v)",
				duration.Round(time.Millisecond), timeout,
			)
		}

		// Check if request was close to timing out (warn if > 80% of timeout used)
		if timeout > 0 {
			usagePercent := float64(duration) / float64(timeout) * 100
			if usagePercent > 80 {
				logger.Warn().
					Str("method", info.FullMethod).
					Dur("timeout", timeout).
					Dur("elapsed", duration).
					Float64("usage_percent", usagePercent).
					Msg("Request approached timeout threshold")

				// Track near-timeouts
				m.NearTimeoutsTotal.WithLabelValues(info.FullMethod).Inc()
			}
		}

		return resp, err
	}
}

// StreamServerInterceptor creates a gRPC interceptor for streaming endpoints.
// Currently streaming is not heavily used, but this provides future-proofing.
func StreamServerInterceptor(m *metrics.Metrics, logger zerolog.Logger) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		// Get timeout for this method
		timeout := GetTimeout(info.FullMethod)

		// Create context with timeout
		ctx, cancel := WithTimeout(ss.Context(), timeout)
		defer cancel()

		// Wrap the stream with timeout context
		wrappedStream := &timeoutServerStream{
			ServerStream: ss,
			ctx:          ctx,
		}

		// Track start time
		start := time.Now()

		// Call handler
		err := handler(srv, wrappedStream)

		// Calculate duration
		duration := time.Since(start)

		// Check if timeout was exceeded
		if IsDeadlineExceeded(err) || IsDeadlineExceeded(ctx.Err()) {
			logger.Warn().
				Str("method", info.FullMethod).
				Dur("timeout", timeout).
				Dur("elapsed", duration).
				Msg("Stream timed out")

			// Increment timeout metric
			m.TimeoutsTotal.WithLabelValues(info.FullMethod).Inc()
			m.TimeoutDuration.WithLabelValues(info.FullMethod).Observe(duration.Seconds())

			return status.Errorf(
				codes.DeadlineExceeded,
				"stream timeout after %v (limit: %v)",
				duration.Round(time.Millisecond), timeout,
			)
		}

		return err
	}
}

// timeoutServerStream wraps grpc.ServerStream to inject timeout context
type timeoutServerStream struct {
	grpc.ServerStream
	ctx context.Context
}

// Context returns the timeout-wrapped context
func (s *timeoutServerStream) Context() context.Context {
	return s.ctx
}
