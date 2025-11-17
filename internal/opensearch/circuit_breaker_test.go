package opensearch

import (
	"errors"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCircuitBreaker(t *testing.T) {
	logger := zerolog.Nop()

	tests := []struct {
		name         string
		maxFailures  int
		resetTimeout time.Duration
		wantMaxFail  int
		wantReset    time.Duration
	}{
		{
			name:         "valid config",
			maxFailures:  5,
			resetTimeout: 60 * time.Second,
			wantMaxFail:  5,
			wantReset:    60 * time.Second,
		},
		{
			name:         "zero maxFailures uses default",
			maxFailures:  0,
			resetTimeout: 60 * time.Second,
			wantMaxFail:  5,
			wantReset:    60 * time.Second,
		},
		{
			name:         "zero resetTimeout uses default",
			maxFailures:  3,
			resetTimeout: 0,
			wantMaxFail:  3,
			wantReset:    60 * time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cb := NewCircuitBreaker(tt.maxFailures, tt.resetTimeout, logger)
			require.NotNil(t, cb)
			assert.Equal(t, tt.wantMaxFail, cb.maxFailures)
			assert.Equal(t, tt.wantReset, cb.resetTimeout)
			assert.Equal(t, StateClosed, cb.state)
			assert.Equal(t, 0, cb.failures)
		})
	}
}

func TestCircuitBreaker_Execute_Success(t *testing.T) {
	logger := zerolog.Nop()
	cb := NewCircuitBreaker(3, 100*time.Millisecond, logger)

	// Successful execution
	err := cb.Execute(func() error {
		return nil
	})

	assert.NoError(t, err)
	assert.Equal(t, StateClosed, cb.GetState())
	assert.Equal(t, 0, cb.GetFailures())
}

func TestCircuitBreaker_Execute_Failure(t *testing.T) {
	logger := zerolog.Nop()
	cb := NewCircuitBreaker(3, 100*time.Millisecond, logger)

	testErr := errors.New("test error")

	// First failure
	err := cb.Execute(func() error {
		return testErr
	})

	assert.Error(t, err)
	assert.Equal(t, testErr, err)
	assert.Equal(t, StateClosed, cb.GetState())
	assert.Equal(t, 1, cb.GetFailures())
}

func TestCircuitBreaker_TransitionToOpen(t *testing.T) {
	logger := zerolog.Nop()
	cb := NewCircuitBreaker(3, 100*time.Millisecond, logger)

	testErr := errors.New("test error")

	// Execute failures until circuit opens
	for i := 0; i < 3; i++ {
		err := cb.Execute(func() error {
			return testErr
		})
		assert.Error(t, err)
	}

	assert.Equal(t, StateOpen, cb.GetState())
	assert.Equal(t, 3, cb.GetFailures())

	// Next execution should fail fast
	err := cb.Execute(func() error {
		t.Fatal("should not execute")
		return nil
	})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "circuit breaker is open")
}

func TestCircuitBreaker_TransitionToHalfOpen(t *testing.T) {
	logger := zerolog.Nop()
	resetTimeout := 100 * time.Millisecond
	cb := NewCircuitBreaker(2, resetTimeout, logger)

	testErr := errors.New("test error")

	// Trigger circuit to open
	for i := 0; i < 2; i++ {
		_ = cb.Execute(func() error {
			return testErr
		})
	}

	assert.Equal(t, StateOpen, cb.GetState())

	// Wait for reset timeout
	time.Sleep(resetTimeout + 10*time.Millisecond)

	// Next execution should transition to half-open
	executed := false
	err := cb.Execute(func() error {
		executed = true
		return nil
	})

	assert.NoError(t, err)
	assert.True(t, executed)
	assert.Equal(t, StateClosed, cb.GetState()) // Successful test closes circuit
	assert.Equal(t, 0, cb.GetFailures())
}

func TestCircuitBreaker_HalfOpenToOpen(t *testing.T) {
	logger := zerolog.Nop()
	resetTimeout := 100 * time.Millisecond
	cb := NewCircuitBreaker(2, resetTimeout, logger)

	testErr := errors.New("test error")

	// Trigger circuit to open
	for i := 0; i < 2; i++ {
		_ = cb.Execute(func() error {
			return testErr
		})
	}

	assert.Equal(t, StateOpen, cb.GetState())

	// Wait for reset timeout
	time.Sleep(resetTimeout + 10*time.Millisecond)

	// Execute with failure - should transition half-open -> open
	err := cb.Execute(func() error {
		return testErr
	})

	assert.Error(t, err)
	assert.Equal(t, StateOpen, cb.GetState())
}

func TestCircuitBreaker_Reset(t *testing.T) {
	logger := zerolog.Nop()
	cb := NewCircuitBreaker(2, 100*time.Millisecond, logger)

	testErr := errors.New("test error")

	// Trigger circuit to open
	for i := 0; i < 2; i++ {
		_ = cb.Execute(func() error {
			return testErr
		})
	}

	assert.Equal(t, StateOpen, cb.GetState())
	assert.Equal(t, 2, cb.GetFailures())

	// Manual reset
	cb.Reset()

	assert.Equal(t, StateClosed, cb.GetState())
	assert.Equal(t, 0, cb.GetFailures())

	// Should be able to execute again
	err := cb.Execute(func() error {
		return nil
	})

	assert.NoError(t, err)
}

func TestCircuitBreaker_GetStats(t *testing.T) {
	logger := zerolog.Nop()
	cb := NewCircuitBreaker(3, 60*time.Second, logger)

	stats := cb.GetStats()

	assert.Equal(t, "closed", stats["state"])
	assert.Equal(t, 0, stats["failures"])
	assert.Equal(t, 3, stats["max_failures"])
	assert.Equal(t, "1m0s", stats["reset_timeout"])
}

func TestCircuitBreaker_ConcurrentAccess(t *testing.T) {
	logger := zerolog.Nop()
	cb := NewCircuitBreaker(10, 100*time.Millisecond, logger)

	// Execute multiple operations concurrently
	done := make(chan bool)
	for i := 0; i < 20; i++ {
		go func() {
			_ = cb.Execute(func() error {
				time.Sleep(1 * time.Millisecond)
				return nil
			})
			done <- true
		}()
	}

	// Wait for all goroutines
	for i := 0; i < 20; i++ {
		<-done
	}

	assert.Equal(t, StateClosed, cb.GetState())
	assert.Equal(t, 0, cb.GetFailures())
}
