package timeout

import (
	"context"
	"errors"
	"time"
)

// WithTimeout wraps context with timeout and returns cancel func.
// If the context already has a deadline that is sooner than the provided duration,
// the existing deadline is preserved.
func WithTimeout(ctx context.Context, duration time.Duration) (context.Context, context.CancelFunc) {
	// Check if context already has deadline
	if deadline, ok := ctx.Deadline(); ok {
		remaining := time.Until(deadline)
		if remaining < duration {
			// Use existing deadline (it's tighter)
			return context.WithCancel(ctx)
		}
	}

	// Set new timeout
	return context.WithTimeout(ctx, duration)
}

// RemainingTime returns time until context deadline.
// Returns 0 if no deadline is set.
func RemainingTime(ctx context.Context) time.Duration {
	if deadline, ok := ctx.Deadline(); ok {
		remaining := time.Until(deadline)
		if remaining < 0 {
			return 0
		}
		return remaining
	}
	return 0
}

// IsDeadlineExceeded checks if error is deadline exceeded.
func IsDeadlineExceeded(err error) bool {
	return errors.Is(err, context.DeadlineExceeded)
}

// HasSufficientTime checks if context has at least the specified duration remaining.
// Returns true if no deadline is set (assume we have time).
// Returns false if deadline exists and insufficient time remains.
func HasSufficientTime(ctx context.Context, required time.Duration) bool {
	deadline, ok := ctx.Deadline()
	if !ok {
		// No deadline set - assume we have time
		return true
	}

	remaining := time.Until(deadline)
	if remaining <= 0 {
		// Deadline already passed
		return false
	}

	return remaining >= required
}

// CheckDeadline returns an error if the context deadline has been exceeded.
func CheckDeadline(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}
