package opensearch

import (
	"fmt"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

// CircuitBreakerState represents the state of the circuit breaker
type CircuitBreakerState int

const (
	// StateClosed - normal operation, requests pass through
	StateClosed CircuitBreakerState = iota
	// StateOpen - circuit is open, requests fail fast
	StateOpen
	// StateHalfOpen - testing if service recovered
	StateHalfOpen
)

func (s CircuitBreakerState) String() string {
	switch s {
	case StateClosed:
		return "closed"
	case StateOpen:
		return "open"
	case StateHalfOpen:
		return "half-open"
	default:
		return "unknown"
	}
}

// CircuitBreaker implements a simple circuit breaker pattern
type CircuitBreaker struct {
	maxFailures  int           // Maximum consecutive failures before opening
	resetTimeout time.Duration // Time to wait before attempting reset
	state        CircuitBreakerState
	failures     int       // Current consecutive failure count
	lastFailure  time.Time // Time of last failure
	mu           sync.RWMutex
	logger       zerolog.Logger
}

// NewCircuitBreaker creates a new circuit breaker
func NewCircuitBreaker(maxFailures int, resetTimeout time.Duration, logger zerolog.Logger) *CircuitBreaker {
	if maxFailures <= 0 {
		maxFailures = 5
	}
	if resetTimeout <= 0 {
		resetTimeout = 60 * time.Second
	}

	return &CircuitBreaker{
		maxFailures:  maxFailures,
		resetTimeout: resetTimeout,
		state:        StateClosed,
		failures:     0,
		logger:       logger.With().Str("component", "circuit_breaker").Logger(),
	}
}

// Execute runs the given function with circuit breaker protection
func (cb *CircuitBreaker) Execute(fn func() error) error {
	// Check if circuit is open
	if !cb.canExecute() {
		return fmt.Errorf("circuit breaker is open")
	}

	// Execute function
	err := fn()

	// Handle result
	if err != nil {
		cb.recordFailure()
		return err
	}

	cb.recordSuccess()
	return nil
}

// canExecute checks if request can be executed based on circuit state
func (cb *CircuitBreaker) canExecute() bool {
	cb.mu.RLock()
	defer cb.mu.RUnlock()

	switch cb.state {
	case StateClosed:
		return true
	case StateOpen:
		// Check if reset timeout expired
		if time.Since(cb.lastFailure) > cb.resetTimeout {
			// Transition to half-open
			cb.mu.RUnlock()
			cb.mu.Lock()
			cb.state = StateHalfOpen
			cb.failures = 0
			cb.logger.Info().Msg("circuit breaker transitioning to half-open")
			cb.mu.Unlock()
			cb.mu.RLock()
			return true
		}
		return false
	case StateHalfOpen:
		return true
	default:
		return false
	}
}

// recordSuccess records successful execution
func (cb *CircuitBreaker) recordSuccess() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	if cb.state == StateHalfOpen {
		// Transition back to closed
		cb.state = StateClosed
		cb.failures = 0
		cb.logger.Info().Msg("circuit breaker closed after successful test")
	} else if cb.state == StateClosed {
		// Reset failure count on success
		if cb.failures > 0 {
			cb.failures = 0
		}
	}
}

// recordFailure records failed execution
func (cb *CircuitBreaker) recordFailure() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.failures++
	cb.lastFailure = time.Now()

	if cb.state == StateHalfOpen {
		// Transition back to open
		cb.state = StateOpen
		cb.logger.Warn().
			Int("failures", cb.failures).
			Msg("circuit breaker opened after half-open test failed")
	} else if cb.state == StateClosed && cb.failures >= cb.maxFailures {
		// Transition to open
		cb.state = StateOpen
		cb.logger.Warn().
			Int("failures", cb.failures).
			Int("max_failures", cb.maxFailures).
			Msg("circuit breaker opened due to consecutive failures")
	}
}

// GetState returns current circuit breaker state (thread-safe)
func (cb *CircuitBreaker) GetState() CircuitBreakerState {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.state
}

// GetFailures returns current failure count (thread-safe)
func (cb *CircuitBreaker) GetFailures() int {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.failures
}

// Reset manually resets the circuit breaker to closed state
func (cb *CircuitBreaker) Reset() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.state = StateClosed
	cb.failures = 0
	cb.logger.Info().Msg("circuit breaker manually reset")
}

// GetStats returns circuit breaker statistics
func (cb *CircuitBreaker) GetStats() map[string]interface{} {
	cb.mu.RLock()
	defer cb.mu.RUnlock()

	return map[string]interface{}{
		"state":              cb.state.String(),
		"failures":           cb.failures,
		"max_failures":       cb.maxFailures,
		"reset_timeout":      cb.resetTimeout.String(),
		"last_failure":       cb.lastFailure,
		"time_since_failure": time.Since(cb.lastFailure).String(),
	}
}
