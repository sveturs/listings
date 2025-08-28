-- Create table for storing failed webhooks for retry and audit
CREATE TABLE IF NOT EXISTS failed_webhooks (
    id BIGSERIAL PRIMARY KEY,
    webhook_id VARCHAR(255) UNIQUE NOT NULL,
    webhook_type VARCHAR(100) NOT NULL,
    endpoint VARCHAR(500),
    payload JSONB NOT NULL,
    signature VARCHAR(500),
    retry_count INTEGER DEFAULT 0,
    max_retries INTEGER DEFAULT 5,
    last_error TEXT,
    status VARCHAR(50) DEFAULT 'pending', -- pending, retrying, failed, completed
    next_retry_at TIMESTAMPTZ,
    metadata JSONB,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    last_attempt_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for efficient querying
CREATE INDEX idx_failed_webhooks_status ON failed_webhooks(status);
CREATE INDEX idx_failed_webhooks_webhook_type ON failed_webhooks(webhook_type);
CREATE INDEX idx_failed_webhooks_next_retry_at ON failed_webhooks(next_retry_at) WHERE status IN ('pending', 'retrying');
CREATE INDEX idx_failed_webhooks_created_at ON failed_webhooks(created_at);

-- Trigger to update updated_at
CREATE OR REPLACE FUNCTION update_failed_webhooks_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_failed_webhooks_updated_at
    BEFORE UPDATE ON failed_webhooks
    FOR EACH ROW
    EXECUTE FUNCTION update_failed_webhooks_updated_at();

-- Table for webhook processing audit log
CREATE TABLE IF NOT EXISTS webhook_audit_log (
    id BIGSERIAL PRIMARY KEY,
    webhook_id VARCHAR(255),
    webhook_type VARCHAR(100) NOT NULL,
    action VARCHAR(50) NOT NULL, -- received, validated, processed, failed, retried
    status VARCHAR(50) NOT NULL, -- success, failure
    details JSONB,
    error_message TEXT,
    processing_time_ms INTEGER,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Index for efficient querying
CREATE INDEX idx_webhook_audit_webhook_id ON webhook_audit_log(webhook_id);
CREATE INDEX idx_webhook_audit_created_at ON webhook_audit_log(created_at);
CREATE INDEX idx_webhook_audit_webhook_type ON webhook_audit_log(webhook_type);
CREATE INDEX idx_webhook_audit_action_status ON webhook_audit_log(action, status);

-- Comments for documentation
COMMENT ON TABLE failed_webhooks IS 'Stores failed webhooks for retry mechanism and audit trail';
COMMENT ON COLUMN failed_webhooks.webhook_id IS 'Unique identifier for the webhook';
COMMENT ON COLUMN failed_webhooks.webhook_type IS 'Type of webhook (allsecure_payment, order_status, etc.)';
COMMENT ON COLUMN failed_webhooks.payload IS 'Original webhook payload in JSON format';
COMMENT ON COLUMN failed_webhooks.signature IS 'Webhook signature for validation';
COMMENT ON COLUMN failed_webhooks.retry_count IS 'Number of retry attempts made';
COMMENT ON COLUMN failed_webhooks.status IS 'Current status: pending, retrying, failed, completed';
COMMENT ON COLUMN failed_webhooks.next_retry_at IS 'Scheduled time for next retry attempt';

COMMENT ON TABLE webhook_audit_log IS 'Audit log for all webhook processing activities';
COMMENT ON COLUMN webhook_audit_log.action IS 'Action performed: received, validated, processed, failed, retried';
COMMENT ON COLUMN webhook_audit_log.processing_time_ms IS 'Time taken to process the webhook in milliseconds';