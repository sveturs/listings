package grpc

import (
	"context"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"
)

func TestLoggingInterceptor(t *testing.T) {
	logger := zerolog.New(nil).Level(zerolog.Disabled)
	interceptor := LoggingInterceptor(logger)

	if interceptor == nil {
		t.Fatal("expected non-nil interceptor")
	}

	// Create a mock invoker
	mockInvoker := func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
		return nil
	}

	err := interceptor(context.Background(), "/test.Method", nil, nil, nil, mockInvoker)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
}

func TestMetricsInterceptor(t *testing.T) {
	interceptor := MetricsInterceptor()

	if interceptor == nil {
		t.Fatal("expected non-nil interceptor")
	}

	mockInvoker := func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
		return nil
	}

	err := interceptor(context.Background(), "/test.Method", nil, nil, nil, mockInvoker)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
}

func TestAuthInterceptor(t *testing.T) {
	token := "test-token"
	interceptor := AuthInterceptor(token)

	if interceptor == nil {
		t.Fatal("expected non-nil interceptor")
	}

	mockInvoker := func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
		// In real scenario, we would verify that metadata contains the token
		return nil
	}

	err := interceptor(context.Background(), "/test.Method", nil, nil, nil, mockInvoker)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
}

func TestRetryInterceptor(t *testing.T) {
	logger := zerolog.New(nil).Level(zerolog.Disabled)
	interceptor := RetryInterceptor(3, 10*time.Millisecond, logger)

	if interceptor == nil {
		t.Fatal("expected non-nil interceptor")
	}

	// Test successful call (no retries needed)
	mockInvoker := func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
		return nil
	}

	err := interceptor(context.Background(), "/test.Method", nil, nil, nil, mockInvoker)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
}

func TestTimeoutInterceptor(t *testing.T) {
	timeout := 100 * time.Millisecond
	interceptor := TimeoutInterceptor(timeout)

	if interceptor == nil {
		t.Fatal("expected non-nil interceptor")
	}

	mockInvoker := func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
		// Check that context has a deadline
		_, hasDeadline := ctx.Deadline()
		if !hasDeadline {
			t.Error("expected context to have deadline")
		}
		return nil
	}

	err := interceptor(context.Background(), "/test.Method", nil, nil, nil, mockInvoker)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
}

func TestChainUnaryClient(t *testing.T) {
	logger := zerolog.New(nil).Level(zerolog.Disabled)

	interceptor1 := LoggingInterceptor(logger)
	interceptor2 := MetricsInterceptor()
	interceptor3 := TimeoutInterceptor(5 * time.Second)

	chained := ChainUnaryClient(interceptor1, interceptor2, interceptor3)

	if chained == nil {
		t.Fatal("expected non-nil chained interceptor")
	}

	mockInvoker := func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
		return nil
	}

	err := chained(context.Background(), "/test.Method", nil, nil, nil, mockInvoker)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
}

func TestIsRetryableError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "context canceled - not retryable",
			err:      context.Canceled,
			expected: false,
		},
		{
			name:     "context deadline exceeded - not retryable",
			err:      context.DeadlineExceeded,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isRetryableError(tt.err)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}
