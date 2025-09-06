package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"backend/pkg/logger"
)

// WebhookRepository handles database operations for webhooks
type WebhookRepository struct {
	db     *sqlx.DB
	logger *logger.Logger
}

// NewWebhookRepository creates a new webhook repository
func NewWebhookRepository(db *sqlx.DB, logger *logger.Logger) *WebhookRepository {
	return &WebhookRepository{
		db:     db,
		logger: logger,
	}
}

// FailedWebhook represents a failed webhook in the database
type FailedWebhook struct {
	ID            int64           `db:"id"`
	WebhookID     string          `db:"webhook_id"`
	WebhookType   string          `db:"webhook_type"`
	Endpoint      sql.NullString  `db:"endpoint"`
	Payload       json.RawMessage `db:"payload"`
	Signature     sql.NullString  `db:"signature"`
	RetryCount    int             `db:"retry_count"`
	MaxRetries    int             `db:"max_retries"`
	LastError     sql.NullString  `db:"last_error"`
	Status        string          `db:"status"`
	NextRetryAt   sql.NullTime    `db:"next_retry_at"`
	Metadata      json.RawMessage `db:"metadata"`
	CreatedAt     time.Time       `db:"created_at"`
	LastAttemptAt sql.NullTime    `db:"last_attempt_at"`
	CompletedAt   sql.NullTime    `db:"completed_at"`
	UpdatedAt     time.Time       `db:"updated_at"`
}

// WebhookAuditLog represents an audit log entry
type WebhookAuditLog struct {
	ID             int64           `db:"id"`
	WebhookID      sql.NullString  `db:"webhook_id"`
	WebhookType    string          `db:"webhook_type"`
	Action         string          `db:"action"`
	Status         string          `db:"status"`
	Details        json.RawMessage `db:"details"`
	ErrorMessage   sql.NullString  `db:"error_message"`
	ProcessingTime sql.NullInt32   `db:"processing_time_ms"`
	CreatedAt      time.Time       `db:"created_at"`
}

// SaveFailedWebhook saves a failed webhook for retry
func (r *WebhookRepository) SaveFailedWebhook(ctx context.Context, webhook FailedWebhook) error {
	query := `
		INSERT INTO failed_webhooks (
			webhook_id, webhook_type, endpoint, payload, signature,
			retry_count, max_retries, last_error, status, next_retry_at,
			metadata, last_attempt_at
		) VALUES (
			:webhook_id, :webhook_type, :endpoint, :payload, :signature,
			:retry_count, :max_retries, :last_error, :status, :next_retry_at,
			:metadata, :last_attempt_at
		)
		ON CONFLICT (webhook_id) 
		DO UPDATE SET
			retry_count = EXCLUDED.retry_count,
			last_error = EXCLUDED.last_error,
			status = EXCLUDED.status,
			next_retry_at = EXCLUDED.next_retry_at,
			last_attempt_at = EXCLUDED.last_attempt_at,
			updated_at = CURRENT_TIMESTAMP
	`

	_, err := r.db.NamedExecContext(ctx, query, webhook)
	if err != nil {
		r.logger.Error("Failed to save failed webhook: %v", err)
		return fmt.Errorf("failed to save failed webhook: %w", err)
	}

	return nil
}

// GetPendingWebhooks retrieves webhooks that need to be retried
func (r *WebhookRepository) GetPendingWebhooks(ctx context.Context, limit int) ([]FailedWebhook, error) {
	query := `
		SELECT * FROM failed_webhooks
		WHERE status IN ('pending', 'retrying')
		AND (next_retry_at IS NULL OR next_retry_at <= CURRENT_TIMESTAMP)
		AND retry_count < max_retries
		ORDER BY next_retry_at ASC NULLS FIRST, created_at ASC
		LIMIT $1
	`

	var webhooks []FailedWebhook
	err := r.db.SelectContext(ctx, &webhooks, query, limit)
	if err != nil {
		r.logger.Error("Failed to get pending webhooks: %v", err)
		return nil, fmt.Errorf("failed to get pending webhooks: %w", err)
	}

	return webhooks, nil
}

// UpdateWebhookStatus updates the status of a webhook
func (r *WebhookRepository) UpdateWebhookStatus(ctx context.Context, webhookID string, status string, error string) error {
	query := `
		UPDATE failed_webhooks
		SET status = $1,
		    last_error = $2,
		    updated_at = CURRENT_TIMESTAMP,
		    completed_at = CASE WHEN $1 = 'completed' THEN CURRENT_TIMESTAMP ELSE completed_at END
		WHERE webhook_id = $3
	`

	var lastError sql.NullString
	if error != "" {
		lastError = sql.NullString{String: error, Valid: true}
	}

	_, err := r.db.ExecContext(ctx, query, status, lastError, webhookID)
	if err != nil {
		r.logger.Error("Failed to update webhook status: %v", err)
		return fmt.Errorf("failed to update webhook status: %w", err)
	}

	return nil
}

// IncrementRetryCount increments the retry count and sets next retry time
func (r *WebhookRepository) IncrementRetryCount(ctx context.Context, webhookID string, nextRetryAt time.Time) error {
	query := `
		UPDATE failed_webhooks
		SET retry_count = retry_count + 1,
		    next_retry_at = $1,
		    last_attempt_at = CURRENT_TIMESTAMP,
		    status = 'retrying',
		    updated_at = CURRENT_TIMESTAMP
		WHERE webhook_id = $2
	`

	_, err := r.db.ExecContext(ctx, query, nextRetryAt, webhookID)
	if err != nil {
		r.logger.Error("Failed to increment retry count: %v", err)
		return fmt.Errorf("failed to increment retry count: %w", err)
	}

	return nil
}

// LogWebhookAudit logs webhook processing activity
func (r *WebhookRepository) LogWebhookAudit(ctx context.Context, audit WebhookAuditLog) error {
	query := `
		INSERT INTO webhook_audit_log (
			webhook_id, webhook_type, action, status,
			details, error_message, processing_time_ms
		) VALUES (
			:webhook_id, :webhook_type, :action, :status,
			:details, :error_message, :processing_time_ms
		)
	`

	_, err := r.db.NamedExecContext(ctx, query, audit)
	if err != nil {
		r.logger.Error("Failed to log webhook audit: %v", err)
		return fmt.Errorf("failed to log webhook audit: %w", err)
	}

	return nil
}

// GetFailedWebhookStats returns statistics about failed webhooks
func (r *WebhookRepository) GetFailedWebhookStats(ctx context.Context) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Get counts by status
	statusQuery := `
		SELECT status, COUNT(*) as count
		FROM failed_webhooks
		GROUP BY status
	`

	rows, err := r.db.QueryContext(ctx, statusQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get status stats: %w", err)
	}
	defer func() { _ = rows.Close() }()

	statusCounts := make(map[string]int)
	for rows.Next() {
		var status string
		var count int
		if err := rows.Scan(&status, &count); err != nil {
			continue
		}
		statusCounts[status] = count
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating status rows: %w", err)
	}
	stats["status_counts"] = statusCounts

	// Get counts by webhook type
	typeQuery := `
		SELECT webhook_type, COUNT(*) as count
		FROM failed_webhooks
		GROUP BY webhook_type
	`

	rows2, err := r.db.QueryContext(ctx, typeQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get type stats: %w", err)
	}
	defer func() { _ = rows2.Close() }()

	typeCounts := make(map[string]int)
	for rows2.Next() {
		var webhookType string
		var count int
		if err := rows2.Scan(&webhookType, &count); err != nil {
			continue
		}
		typeCounts[webhookType] = count
	}

	if err := rows2.Err(); err != nil {
		return nil, fmt.Errorf("error iterating type rows: %w", err)
	}
	stats["type_counts"] = typeCounts

	// Get recent failures
	recentQuery := `
		SELECT COUNT(*) 
		FROM failed_webhooks 
		WHERE created_at > CURRENT_TIMESTAMP - INTERVAL '1 hour'
		AND status = 'failed'
	`

	var recentFailures int
	err = r.db.GetContext(ctx, &recentFailures, recentQuery)
	if err == nil {
		stats["recent_failures_1h"] = recentFailures
	}

	return stats, nil
}

// CleanupOldWebhooks removes old completed or permanently failed webhooks
func (r *WebhookRepository) CleanupOldWebhooks(ctx context.Context, olderThan time.Duration) (int64, error) {
	cutoff := time.Now().Add(-olderThan)

	query := `
		DELETE FROM failed_webhooks
		WHERE (status = 'completed' OR (status = 'failed' AND retry_count >= max_retries))
		AND updated_at < $1
	`

	result, err := r.db.ExecContext(ctx, query, cutoff)
	if err != nil {
		r.logger.Error("Failed to cleanup old webhooks: %v", err)
		return 0, fmt.Errorf("failed to cleanup old webhooks: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %w", err)
	}

	return rowsAffected, nil
}
