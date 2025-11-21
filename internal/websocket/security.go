package websocket

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

// RateLimiter implements rate limiting for WebSocket messages
type RateLimiter struct {
	// Map of user_id -> message timestamps
	userMessages map[int64][]time.Time
	mutex        sync.RWMutex
	logger       zerolog.Logger

	// Configuration
	maxMessages   int           // Max messages per window
	windowSeconds int           // Time window in seconds
	cleanupTicker *time.Ticker  // For periodic cleanup
	stopChan      chan struct{} // For stopping cleanup goroutine
}

// NewRateLimiter creates a new rate limiter for WebSocket messages
func NewRateLimiter(maxMessages, windowSeconds int, logger zerolog.Logger) *RateLimiter {
	rl := &RateLimiter{
		userMessages:  make(map[int64][]time.Time),
		maxMessages:   maxMessages,
		windowSeconds: windowSeconds,
		logger:        logger.With().Str("component", "ws_rate_limiter").Logger(),
		stopChan:      make(chan struct{}),
	}

	// Start cleanup goroutine to prevent memory leaks
	rl.cleanupTicker = time.NewTicker(time.Duration(windowSeconds) * time.Second)
	go rl.cleanupLoop()

	return rl
}

// CheckLimit checks if a user can send a message
// Returns true if allowed, false if rate limit exceeded
func (rl *RateLimiter) CheckLimit(userID int64) bool {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	now := time.Now()
	cutoff := now.Add(-time.Duration(rl.windowSeconds) * time.Second)

	// Get user's recent messages
	messages, exists := rl.userMessages[userID]
	if !exists {
		// First message from this user
		rl.userMessages[userID] = []time.Time{now}
		return true
	}

	// Filter out old messages
	var recent []time.Time
	for _, ts := range messages {
		if ts.After(cutoff) {
			recent = append(recent, ts)
		}
	}

	// Check if limit exceeded
	if len(recent) >= rl.maxMessages {
		rl.logger.Warn().
			Int64("user_id", userID).
			Int("message_count", len(recent)).
			Int("max_allowed", rl.maxMessages).
			Msg("rate limit exceeded")
		return false
	}

	// Add current message
	recent = append(recent, now)
	rl.userMessages[userID] = recent

	return true
}

// cleanupLoop periodically removes old entries to prevent memory leaks
func (rl *RateLimiter) cleanupLoop() {
	for {
		select {
		case <-rl.cleanupTicker.C:
			rl.cleanup()
		case <-rl.stopChan:
			rl.cleanupTicker.Stop()
			return
		}
	}
}

// cleanup removes old message timestamps
func (rl *RateLimiter) cleanup() {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	now := time.Now()
	cutoff := now.Add(-time.Duration(rl.windowSeconds*2) * time.Second) // Keep 2x window for safety

	for userID, messages := range rl.userMessages {
		// Filter out old messages
		var recent []time.Time
		for _, ts := range messages {
			if ts.After(cutoff) {
				recent = append(recent, ts)
			}
		}

		if len(recent) == 0 {
			// No recent messages, remove user entirely
			delete(rl.userMessages, userID)
		} else {
			rl.userMessages[userID] = recent
		}
	}

	rl.logger.Debug().
		Int("active_users", len(rl.userMessages)).
		Msg("rate limiter cleanup completed")
}

// Stop stops the cleanup goroutine
func (rl *RateLimiter) Stop() {
	close(rl.stopChan)
}

// OriginChecker checks if the origin is allowed
type OriginChecker struct {
	allowedOrigins map[string]bool
	logger         zerolog.Logger
}

// NewOriginChecker creates a new origin checker
func NewOriginChecker(allowedOrigins []string, logger zerolog.Logger) *OriginChecker {
	origins := make(map[string]bool)
	for _, origin := range allowedOrigins {
		origins[origin] = true
	}

	return &OriginChecker{
		allowedOrigins: origins,
		logger:         logger.With().Str("component", "origin_checker").Logger(),
	}
}

// IsAllowed checks if an origin is allowed
func (oc *OriginChecker) IsAllowed(origin string) bool {
	if origin == "" {
		return false
	}

	allowed := oc.allowedOrigins[origin]
	if !allowed {
		oc.logger.Warn().
			Str("origin", origin).
			Msg("origin not allowed")
	}

	return allowed
}

// ConnectionTracker tracks connection statistics
type ConnectionTracker struct {
	// Metrics
	totalConnections   int64
	currentConnections int64
	messagesReceived   int64
	messagesSent       int64
	errorsTotal        int64

	mutex  sync.RWMutex
	logger zerolog.Logger
}

// NewConnectionTracker creates a new connection tracker
func NewConnectionTracker(logger zerolog.Logger) *ConnectionTracker {
	return &ConnectionTracker{
		logger: logger.With().Str("component", "connection_tracker").Logger(),
	}
}

// OnConnect increments connection counters
func (ct *ConnectionTracker) OnConnect() {
	ct.mutex.Lock()
	defer ct.mutex.Unlock()

	ct.totalConnections++
	ct.currentConnections++

	ct.logger.Debug().
		Int64("total", ct.totalConnections).
		Int64("current", ct.currentConnections).
		Msg("connection established")
}

// OnDisconnect decrements current connections
func (ct *ConnectionTracker) OnDisconnect() {
	ct.mutex.Lock()
	defer ct.mutex.Unlock()

	ct.currentConnections--

	ct.logger.Debug().
		Int64("current", ct.currentConnections).
		Msg("connection closed")
}

// OnMessageReceived increments received message counter
func (ct *ConnectionTracker) OnMessageReceived() {
	ct.mutex.Lock()
	defer ct.mutex.Unlock()

	ct.messagesReceived++
}

// OnMessageSent increments sent message counter
func (ct *ConnectionTracker) OnMessageSent() {
	ct.mutex.Lock()
	defer ct.mutex.Unlock()

	ct.messagesSent++
}

// OnError increments error counter
func (ct *ConnectionTracker) OnError() {
	ct.mutex.Lock()
	defer ct.mutex.Unlock()

	ct.errorsTotal++
}

// GetStats returns current statistics
func (ct *ConnectionTracker) GetStats() map[string]int64 {
	ct.mutex.RLock()
	defer ct.mutex.RUnlock()

	return map[string]int64{
		"total_connections":   ct.totalConnections,
		"current_connections": ct.currentConnections,
		"messages_received":   ct.messagesReceived,
		"messages_sent":       ct.messagesSent,
		"errors_total":        ct.errorsTotal,
	}
}

// SecurityMiddleware provides security features for WebSocket
type SecurityMiddleware struct {
	rateLimiter       *RateLimiter
	originChecker     *OriginChecker
	connectionTracker *ConnectionTracker
	logger            zerolog.Logger
}

// NewSecurityMiddleware creates a new security middleware
func NewSecurityMiddleware(
	maxMessagesPerMinute int,
	allowedOrigins []string,
	logger zerolog.Logger,
) *SecurityMiddleware {
	return &SecurityMiddleware{
		rateLimiter:       NewRateLimiter(maxMessagesPerMinute, 60, logger),
		originChecker:     NewOriginChecker(allowedOrigins, logger),
		connectionTracker: NewConnectionTracker(logger),
		logger:            logger.With().Str("component", "security_middleware").Logger(),
	}
}

// CheckRateLimit checks if user can send message
func (sm *SecurityMiddleware) CheckRateLimit(userID int64) error {
	if !sm.rateLimiter.CheckLimit(userID) {
		return fmt.Errorf("rate limit exceeded: too many messages")
	}
	return nil
}

// CheckOrigin checks if origin is allowed
func (sm *SecurityMiddleware) CheckOrigin(origin string) error {
	if !sm.originChecker.IsAllowed(origin) {
		return fmt.Errorf("origin not allowed: %s", origin)
	}
	return nil
}

// TrackConnection tracks connection lifecycle
func (sm *SecurityMiddleware) TrackConnection(ctx context.Context) func() {
	sm.connectionTracker.OnConnect()

	return func() {
		sm.connectionTracker.OnDisconnect()
	}
}

// TrackMessage tracks message activity
func (sm *SecurityMiddleware) TrackMessage(received bool) {
	if received {
		sm.connectionTracker.OnMessageReceived()
	} else {
		sm.connectionTracker.OnMessageSent()
	}
}

// TrackError tracks errors
func (sm *SecurityMiddleware) TrackError() {
	sm.connectionTracker.OnError()
}

// GetStats returns all statistics
func (sm *SecurityMiddleware) GetStats() map[string]interface{} {
	stats := sm.connectionTracker.GetStats()

	return map[string]interface{}{
		"connections": stats,
	}
}

// Stop stops the middleware
func (sm *SecurityMiddleware) Stop() {
	sm.rateLimiter.Stop()
}
