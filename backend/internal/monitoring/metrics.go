package monitoring

import (
	"context"
	"sync"
	"time"

	"backend/pkg/logger"
)

// MetricsCollector collects and reports metrics
type MetricsCollector struct {
	mu                   sync.RWMutex
	rateLimitHits        map[string]int64
	rateLimitExceeded    map[string]int64
	webhookRetries       map[string]int64
	paymentAttempts      map[string]int64
	suspiciousActivities []SuspiciousActivity
	logger               *logger.Logger
}

// SuspiciousActivity represents a suspicious activity event
type SuspiciousActivity struct {
	Timestamp   time.Time              `json:"timestamp"`
	UserID      int                    `json:"user_id,omitempty"`
	IP          string                 `json:"ip"`
	Type        string                 `json:"type"`
	Description string                 `json:"description"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// NewMetricsCollector creates a new metrics collector
func NewMetricsCollector(logger *logger.Logger) *MetricsCollector {
	return &MetricsCollector{
		rateLimitHits:        make(map[string]int64),
		rateLimitExceeded:    make(map[string]int64),
		webhookRetries:       make(map[string]int64),
		paymentAttempts:      make(map[string]int64),
		suspiciousActivities: make([]SuspiciousActivity, 0),
		logger:               logger,
	}
}

// RecordRateLimitHit records a rate limit check
func (m *MetricsCollector) RecordRateLimitHit(endpoint string, userID int, ip string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	key := endpoint
	m.rateLimitHits[key]++

	// Log for monitoring
	m.logger.Info("Rate limit hit: endpoint=%s, user_id=%d, ip=%s, total_hits=%d",
		endpoint, userID, ip, m.rateLimitHits[key])
}

// RecordRateLimitExceeded records when rate limit is exceeded
func (m *MetricsCollector) RecordRateLimitExceeded(endpoint string, userID int, ip string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	key := endpoint
	m.rateLimitExceeded[key]++

	// Log warning for monitoring
	m.logger.Error("RATE LIMIT EXCEEDED: endpoint=%s, user_id=%d, ip=%s, total_exceeded=%d",
		endpoint, userID, ip, m.rateLimitExceeded[key])

	// Check for suspicious patterns
	if m.rateLimitExceeded[key] > 10 {
		m.addSuspiciousActivity(SuspiciousActivity{
			Timestamp:   time.Now(),
			UserID:      userID,
			IP:          ip,
			Type:        "rate_limit_abuse",
			Description: "Multiple rate limit violations",
			Metadata: map[string]interface{}{
				"endpoint":   endpoint,
				"violations": m.rateLimitExceeded[key],
			},
		})
	}
}

// RecordWebhookRetry records webhook retry attempts
func (m *MetricsCollector) RecordWebhookRetry(webhookType string, attemptNumber int) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.webhookRetries[webhookType]++

	if attemptNumber > 3 {
		m.logger.Error("HIGH WEBHOOK RETRY COUNT: type=%s, attempt=%d",
			webhookType, attemptNumber)
	}
}

// RecordPaymentAttempt records payment attempts
func (m *MetricsCollector) RecordPaymentAttempt(userID int, amount float64, success bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	key := "failed"
	if success {
		key = "successful"
	}

	m.paymentAttempts[key]++
	m.paymentAttempts["total"]++

	// Check for suspicious payment patterns
	failedKey := "failed"
	if m.paymentAttempts[failedKey] > 5 {
		m.addSuspiciousActivity(SuspiciousActivity{
			Timestamp:   time.Now(),
			UserID:      userID,
			Type:        "payment_fraud_attempt",
			Description: "Multiple failed payment attempts",
			Metadata: map[string]interface{}{
				"failed_attempts": m.paymentAttempts[failedKey],
				"amount":          amount,
			},
		})
	}
}

// addSuspiciousActivity adds a suspicious activity
func (m *MetricsCollector) addSuspiciousActivity(activity SuspiciousActivity) {
	m.suspiciousActivities = append(m.suspiciousActivities, activity)

	// Alert immediately for critical activities
	m.logger.Error("🚨 SUSPICIOUS ACTIVITY DETECTED: type=%s, user_id=%d, ip=%s, description=%s",
		activity.Type, activity.UserID, activity.IP, activity.Description)

	// TODO: Send alert to administrators (email, Slack, etc.)
}

// GetMetrics returns current metrics
func (m *MetricsCollector) GetMetrics() map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return map[string]interface{}{
		"rate_limit_hits":       m.rateLimitHits,
		"rate_limit_exceeded":   m.rateLimitExceeded,
		"webhook_retries":       m.webhookRetries,
		"payment_attempts":      m.paymentAttempts,
		"suspicious_activities": m.suspiciousActivities,
	}
}

// StartMetricsReporter starts periodic metrics reporting
func (m *MetricsCollector) StartMetricsReporter(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for {
			select {
			case <-ctx.Done():
				ticker.Stop()
				return
			case <-ticker.C:
				m.reportMetrics()
			}
		}
	}()
}

// reportMetrics reports metrics periodically
func (m *MetricsCollector) reportMetrics() {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Calculate totals
	totalRateLimitHits := int64(0)
	for _, count := range m.rateLimitHits {
		totalRateLimitHits += count
	}

	totalRateLimitExceeded := int64(0)
	for _, count := range m.rateLimitExceeded {
		totalRateLimitExceeded += count
	}

	totalWebhookRetries := int64(0)
	for _, count := range m.webhookRetries {
		totalWebhookRetries += count
	}

	// Log metrics summary
	m.logger.Info("📊 METRICS REPORT: rate_limit_hits=%d, rate_limit_exceeded=%d, webhook_retries=%d, suspicious_activities=%d",
		totalRateLimitHits, totalRateLimitExceeded, totalWebhookRetries, len(m.suspiciousActivities))

	// Alert if thresholds exceeded
	if totalRateLimitExceeded > 100 {
		m.logger.Error("⚠️ HIGH RATE LIMIT VIOLATIONS: %d exceeds threshold", totalRateLimitExceeded)
	}

	if totalWebhookRetries > 50 {
		m.logger.Error("⚠️ HIGH WEBHOOK RETRY RATE: %d retries", totalWebhookRetries)
	}

	if len(m.suspiciousActivities) > 10 {
		m.logger.Error("⚠️ MULTIPLE SUSPICIOUS ACTIVITIES: %d events detected", len(m.suspiciousActivities))
	}
}

// ClearOldMetrics clears metrics older than the specified duration
func (m *MetricsCollector) ClearOldMetrics(olderThan time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()

	cutoff := time.Now().Add(-olderThan)

	// Clear old suspicious activities
	newActivities := make([]SuspiciousActivity, 0)
	for _, activity := range m.suspiciousActivities {
		if activity.Timestamp.After(cutoff) {
			newActivities = append(newActivities, activity)
		}
	}
	m.suspiciousActivities = newActivities

	// Reset counters periodically (e.g., daily)
	// This is a simple implementation; in production, use time-series data
	if time.Now().Hour() == 0 && time.Now().Minute() < 5 {
		m.rateLimitHits = make(map[string]int64)
		m.rateLimitExceeded = make(map[string]int64)
		m.webhookRetries = make(map[string]int64)
		m.paymentAttempts = make(map[string]int64)
	}
}
