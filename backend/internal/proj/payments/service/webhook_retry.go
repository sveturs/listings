package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"backend/pkg/logger"
)

// WebhookRetryConfig holds configuration for webhook retry mechanism
type WebhookRetryConfig struct {
	MaxRetries     int              // Maximum number of retry attempts
	InitialDelay   time.Duration    // Initial delay before first retry
	MaxDelay       time.Duration    // Maximum delay between retries
	BackoffFactor  float64          // Exponential backoff factor
	RetryableError func(error) bool // Function to determine if error is retryable
}

// DefaultWebhookRetryConfig returns default retry configuration
func DefaultWebhookRetryConfig() WebhookRetryConfig {
	return WebhookRetryConfig{
		MaxRetries:    5,
		InitialDelay:  1 * time.Second,
		MaxDelay:      5 * time.Minute,
		BackoffFactor: 2.0,
		RetryableError: func(err error) bool {
			// By default, retry all errors except validation errors
			// You can customize this based on error types
			return true
		},
	}
}

// WebhookRetryManager manages webhook retry logic
type WebhookRetryManager struct {
	config WebhookRetryConfig
	logger *logger.Logger
	queue  chan *WebhookRetryJob
}

// WebhookRetryJob represents a webhook that needs to be retried
type WebhookRetryJob struct {
	ID            string                 `json:"id"`
	WebhookType   string                 `json:"webhook_type"`
	Payload       []byte                 `json:"payload"`
	Signature     string                 `json:"signature"`
	Endpoint      string                 `json:"endpoint"`
	RetryCount    int                    `json:"retry_count"`
	LastError     string                 `json:"last_error,omitempty"`
	CreatedAt     time.Time              `json:"created_at"`
	LastAttemptAt time.Time              `json:"last_attempt_at"`
	NextRetryAt   time.Time              `json:"next_retry_at"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
}

// NewWebhookRetryManager creates a new webhook retry manager
func NewWebhookRetryManager(config WebhookRetryConfig, logger *logger.Logger) *WebhookRetryManager {
	return &WebhookRetryManager{
		config: config,
		logger: logger,
		queue:  make(chan *WebhookRetryJob, 1000), // Buffer for 1000 jobs
	}
}

// Start starts the retry manager worker
func (m *WebhookRetryManager) Start(ctx context.Context) {
	go m.worker(ctx)
}

// worker processes retry jobs from the queue
func (m *WebhookRetryManager) worker(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			m.logger.Info("Webhook retry worker shutting down")
			return
		case job := <-m.queue:
			m.processRetryJob(ctx, job)
		}
	}
}

// AddRetryJob adds a webhook to the retry queue
func (m *WebhookRetryManager) AddRetryJob(job *WebhookRetryJob) error {
	// Calculate next retry time with exponential backoff
	if job.RetryCount == 0 {
		job.NextRetryAt = time.Now().Add(m.config.InitialDelay)
	} else {
		delay := m.calculateBackoffDelay(job.RetryCount)
		job.NextRetryAt = time.Now().Add(delay)
	}

	// Check if we've exceeded max retries
	if job.RetryCount >= m.config.MaxRetries {
		m.logger.Error("Webhook retry limit exceeded for job %s after %d attempts",
			job.ID, job.RetryCount)
		// Here you might want to send an alert or store in a dead letter queue
		return fmt.Errorf("max retries exceeded")
	}

	// Add to queue with delay
	go func() {
		time.Sleep(time.Until(job.NextRetryAt))
		select {
		case m.queue <- job:
			m.logger.Info("Added webhook job %s to retry queue (attempt %d/%d)",
				job.ID, job.RetryCount+1, m.config.MaxRetries)
		default:
			m.logger.Error("Retry queue is full, dropping job %s", job.ID)
		}
	}()

	return nil
}

// processRetryJob processes a single retry job
func (m *WebhookRetryManager) processRetryJob(ctx context.Context, job *WebhookRetryJob) {
	m.logger.Info("Processing webhook retry job %s (attempt %d/%d)",
		job.ID, job.RetryCount+1, m.config.MaxRetries)

	job.LastAttemptAt = time.Now()
	job.RetryCount++

	// Here you would call the actual webhook processing function
	// For example:
	err := m.executeWebhook(ctx, job)

	if err != nil {
		job.LastError = err.Error()

		// Check if error is retryable
		if m.config.RetryableError(err) && job.RetryCount < m.config.MaxRetries {
			m.logger.Info("Webhook %s failed (attempt %d/%d): %v. Will retry.",
				job.ID, job.RetryCount, m.config.MaxRetries, err)

			// Re-add to retry queue
			if retryErr := m.AddRetryJob(job); retryErr != nil {
				m.logger.Error("Failed to re-queue webhook %s: %v", job.ID, retryErr)
			}
		} else {
			m.logger.Error("Webhook %s failed permanently after %d attempts: %v",
				job.ID, job.RetryCount, err)
			// Store in dead letter queue or alert administrators
			m.handlePermanentFailure(job)
		}
	} else {
		m.logger.Info("Webhook %s processed successfully after %d attempts",
			job.ID, job.RetryCount)
	}
}

// executeWebhook executes the actual webhook
func (m *WebhookRetryManager) executeWebhook(ctx context.Context, job *WebhookRetryJob) error {
	// This is a placeholder - implement actual webhook execution
	// based on webhook type and endpoint

	m.logger.Info("Executing webhook %s to endpoint %s", job.ID, job.Endpoint)

	// Simulate webhook execution
	// In real implementation, you would:
	// 1. Make HTTP request to the endpoint
	// 2. Validate response
	// 3. Handle specific webhook types differently

	return nil
}

// calculateBackoffDelay calculates exponential backoff delay
func (m *WebhookRetryManager) calculateBackoffDelay(retryCount int) time.Duration {
	delay := float64(m.config.InitialDelay)

	for i := 0; i < retryCount; i++ {
		delay *= m.config.BackoffFactor
	}

	// Cap at max delay
	if time.Duration(delay) > m.config.MaxDelay {
		return m.config.MaxDelay
	}

	return time.Duration(delay)
}

// handlePermanentFailure handles webhooks that have failed permanently
func (m *WebhookRetryManager) handlePermanentFailure(job *WebhookRetryJob) {
	// Log detailed failure information
	m.logger.Error("Webhook permanent failure - ID: %s, Type: %s, Endpoint: %s, Attempts: %d, Error: %s",
		job.ID, job.WebhookType, job.Endpoint, job.RetryCount, job.LastError)

	// Store in database for manual review
	failureRecord := map[string]interface{}{
		"webhook_id":   job.ID,
		"webhook_type": job.WebhookType,
		"endpoint":     job.Endpoint,
		"retry_count":  job.RetryCount,
		"last_error":   job.LastError,
		"payload_size": len(job.Payload),
		"created_at":   job.CreatedAt,
		"failed_at":    time.Now(),
		"metadata":     job.Metadata,
	}

	// Convert to JSON for logging
	if data, err := json.MarshalIndent(failureRecord, "", "  "); err == nil {
		m.logger.Error("Failed webhook details:\n%s", string(data))
	}

	// TODO: Implement actual storage to database or dead letter queue
	// TODO: Send alert to administrators
	// TODO: Potentially trigger compensating transaction
}

// GetQueueSize returns the current size of the retry queue
func (m *WebhookRetryManager) GetQueueSize() int {
	return len(m.queue)
}

// GetQueueStats returns statistics about the retry queue
func (m *WebhookRetryManager) GetQueueStats() map[string]interface{} {
	return map[string]interface{}{
		"queue_size":     len(m.queue),
		"queue_capacity": cap(m.queue),
		"config": map[string]interface{}{
			"max_retries":    m.config.MaxRetries,
			"initial_delay":  m.config.InitialDelay.String(),
			"max_delay":      m.config.MaxDelay.String(),
			"backoff_factor": m.config.BackoffFactor,
		},
	}
}
