package listings

import (
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TestMapGRPCError проверяет маппинг gRPC ошибок на доменные ошибки
func TestMapGRPCError(t *testing.T) {
	tests := []struct {
		name     string
		input    error
		expected error
	}{
		{
			name:     "nil error",
			input:    nil,
			expected: nil,
		},
		{
			name:     "not found",
			input:    status.Error(codes.NotFound, "not found"),
			expected: ErrListingNotFound,
		},
		{
			name:     "invalid argument",
			input:    status.Error(codes.InvalidArgument, "invalid"),
			expected: ErrInvalidInput,
		},
		{
			name:     "permission denied",
			input:    status.Error(codes.PermissionDenied, "denied"),
			expected: ErrUnauthorized,
		},
		{
			name:     "unauthenticated",
			input:    status.Error(codes.Unauthenticated, "unauth"),
			expected: ErrUnauthorized,
		},
		{
			name:     "already exists",
			input:    status.Error(codes.AlreadyExists, "exists"),
			expected: ErrAlreadyExists,
		},
		{
			name:     "unavailable",
			input:    status.Error(codes.Unavailable, "unavailable"),
			expected: ErrServiceUnavailable,
		},
		{
			name:     "deadline exceeded",
			input:    status.Error(codes.DeadlineExceeded, "timeout"),
			expected: ErrServiceUnavailable,
		},
		{
			name:     "internal error",
			input:    status.Error(codes.Internal, "internal"),
			expected: ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MapGRPCError(tt.input)
			if result != tt.expected {
				t.Errorf("MapGRPCError() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestShouldRetry проверяет логику повторных попыток
func TestShouldRetry(t *testing.T) {
	client := &Client{}

	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "should retry on unavailable",
			err:      status.Error(codes.Unavailable, "unavailable"),
			expected: true,
		},
		{
			name:     "should retry on deadline exceeded",
			err:      status.Error(codes.DeadlineExceeded, "timeout"),
			expected: true,
		},
		{
			name:     "should not retry on invalid argument",
			err:      status.Error(codes.InvalidArgument, "invalid"),
			expected: false,
		},
		{
			name:     "should not retry on not found",
			err:      status.Error(codes.NotFound, "not found"),
			expected: false,
		},
		{
			name:     "should not retry on permission denied",
			err:      status.Error(codes.PermissionDenied, "denied"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := client.shouldRetry(tt.err)
			if result != tt.expected {
				t.Errorf("shouldRetry() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestIsErrorHelpers проверяет вспомогательные функции для проверки ошибок
func TestIsErrorHelpers(t *testing.T) {
	t.Run("IsNotFound", func(t *testing.T) {
		if !IsNotFound(ErrListingNotFound) {
			t.Error("IsNotFound() should return true for ErrListingNotFound")
		}
		if IsNotFound(ErrInvalidInput) {
			t.Error("IsNotFound() should return false for non-NotFound error")
		}
	})

	t.Run("IsInvalidInput", func(t *testing.T) {
		if !IsInvalidInput(ErrInvalidInput) {
			t.Error("IsInvalidInput() should return true for ErrInvalidInput")
		}
		if IsInvalidInput(ErrListingNotFound) {
			t.Error("IsInvalidInput() should return false for non-InvalidInput error")
		}
	})

	t.Run("IsUnauthorized", func(t *testing.T) {
		if !IsUnauthorized(ErrUnauthorized) {
			t.Error("IsUnauthorized() should return true for ErrUnauthorized")
		}
		if IsUnauthorized(ErrInvalidInput) {
			t.Error("IsUnauthorized() should return false for non-Unauthorized error")
		}
	})

	t.Run("IsServiceUnavailable", func(t *testing.T) {
		if !IsServiceUnavailable(ErrServiceUnavailable) {
			t.Error("IsServiceUnavailable() should return true for ErrServiceUnavailable")
		}
		if IsServiceUnavailable(ErrInvalidInput) {
			t.Error("IsServiceUnavailable() should return false for non-ServiceUnavailable error")
		}
	})
}
