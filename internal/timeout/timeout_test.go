package timeout

import (
	"context"
	"testing"
	"time"
)

func TestWithTimeout(t *testing.T) {
	tests := []struct {
		name             string
		existingDeadline time.Duration // 0 means no existing deadline
		newTimeout       time.Duration
		expectNewTimeout bool
	}{
		{
			name:             "no existing deadline - sets new timeout",
			existingDeadline: 0,
			newTimeout:       5 * time.Second,
			expectNewTimeout: true,
		},
		{
			name:             "existing deadline is sooner - keeps existing",
			existingDeadline: 2 * time.Second,
			newTimeout:       5 * time.Second,
			expectNewTimeout: false,
		},
		{
			name:             "existing deadline is later - sets new timeout",
			existingDeadline: 10 * time.Second,
			newTimeout:       5 * time.Second,
			expectNewTimeout: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var baseCtx context.Context
			if tt.existingDeadline > 0 {
				var cancel context.CancelFunc
				baseCtx, cancel = context.WithTimeout(context.Background(), tt.existingDeadline)
				defer cancel()
			} else {
				baseCtx = context.Background()
			}

			ctx, cancel := WithTimeout(baseCtx, tt.newTimeout)
			defer cancel()

			deadline, ok := ctx.Deadline()
			if !ok {
				t.Fatal("expected deadline to be set")
			}

			remaining := time.Until(deadline)

			// Allow 100ms tolerance for test execution time
			tolerance := 100 * time.Millisecond

			if tt.expectNewTimeout {
				expected := tt.newTimeout
				if remaining < expected-tolerance || remaining > expected+tolerance {
					t.Errorf("expected timeout ~%v, got %v", expected, remaining)
				}
			} else {
				expected := tt.existingDeadline
				if remaining < expected-tolerance || remaining > expected+tolerance {
					t.Errorf("expected to preserve existing deadline ~%v, got %v", expected, remaining)
				}
			}
		})
	}
}

func TestRemainingTime(t *testing.T) {
	tests := []struct {
		name     string
		setup    func() context.Context
		validate func(time.Duration) bool
	}{
		{
			name: "context with deadline",
			setup: func() context.Context {
				ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
				return ctx
			},
			validate: func(d time.Duration) bool {
				return d > 4*time.Second && d <= 5*time.Second
			},
		},
		{
			name: "context without deadline",
			setup: func() context.Context {
				return context.Background()
			},
			validate: func(d time.Duration) bool {
				return d == 0
			},
		},
		{
			name: "expired context",
			setup: func() context.Context {
				ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
				defer cancel()
				time.Sleep(10 * time.Millisecond)
				return ctx
			},
			validate: func(d time.Duration) bool {
				return d == 0
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := tt.setup()
			remaining := RemainingTime(ctx)

			if !tt.validate(remaining) {
				t.Errorf("unexpected remaining time: %v", remaining)
			}
		})
	}
}

func TestIsDeadlineExceeded(t *testing.T) {
	tests := []struct {
		name   string
		err    error
		expect bool
	}{
		{
			name:   "deadline exceeded error",
			err:    context.DeadlineExceeded,
			expect: true,
		},
		{
			name:   "cancelled error",
			err:    context.Canceled,
			expect: false,
		},
		{
			name:   "nil error",
			err:    nil,
			expect: false,
		},
		{
			name:   "other error",
			err:    context.TODO().Err(),
			expect: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsDeadlineExceeded(tt.err)
			if result != tt.expect {
				t.Errorf("expected %v, got %v", tt.expect, result)
			}
		})
	}
}

func TestHasSufficientTime(t *testing.T) {
	tests := []struct {
		name     string
		setup    func() context.Context
		required time.Duration
		expect   bool
	}{
		{
			name: "sufficient time remaining",
			setup: func() context.Context {
				ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
				return ctx
			},
			required: 5 * time.Second,
			expect:   true,
		},
		{
			name: "insufficient time remaining",
			setup: func() context.Context {
				ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
				return ctx
			},
			required: 5 * time.Second,
			expect:   false,
		},
		{
			name: "no deadline - assume sufficient",
			setup: func() context.Context {
				return context.Background()
			},
			required: 5 * time.Second,
			expect:   true,
		},
		{
			name: "expired context",
			setup: func() context.Context {
				ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
				defer cancel()
				time.Sleep(10 * time.Millisecond)
				return ctx
			},
			required: 1 * time.Second,
			expect:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := tt.setup()
			result := HasSufficientTime(ctx, tt.required)
			if result != tt.expect {
				t.Errorf("expected %v, got %v (remaining: %v)", tt.expect, result, RemainingTime(ctx))
			}
		})
	}
}

func TestCheckDeadline(t *testing.T) {
	tests := []struct {
		name      string
		setup     func() context.Context
		expectErr bool
	}{
		{
			name: "context not expired",
			setup: func() context.Context {
				ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
				return ctx
			},
			expectErr: false,
		},
		{
			name: "context expired",
			setup: func() context.Context {
				ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
				defer cancel()
				time.Sleep(10 * time.Millisecond)
				return ctx
			},
			expectErr: true,
		},
		{
			name: "context cancelled",
			setup: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			},
			expectErr: true,
		},
		{
			name: "no deadline",
			setup: func() context.Context {
				return context.Background()
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := tt.setup()
			err := CheckDeadline(ctx)

			if tt.expectErr && err == nil {
				t.Error("expected error, got nil")
			}
			if !tt.expectErr && err != nil {
				t.Errorf("expected no error, got: %v", err)
			}
		})
	}
}
