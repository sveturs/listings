package metrics

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// UnaryServerInterceptor returns a gRPC unary interceptor that records metrics for all handler calls
func (m *Metrics) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// Extract method name from info
		method := info.FullMethod

		// Track active requests
		m.GRPCHandlerRequestsActive.WithLabelValues(method).Inc()
		defer m.GRPCHandlerRequestsActive.WithLabelValues(method).Dec()

		// Record start time
		start := time.Now()

		// Call the actual handler
		resp, err := handler(ctx, req)

		// Record duration
		duration := time.Since(start).Seconds()

		// Determine status
		statusCode := "success"
		if err != nil {
			st, _ := status.FromError(err)
			statusCode = st.Code().String()
		}

		// Record metrics
		m.RecordGRPCRequest(method, statusCode, duration)

		return resp, err
	}
}
