-- Create Dead Letter Queue table for failed OpenSearch indexing tasks
CREATE TABLE IF NOT EXISTS opensearch_indexing_dlq (
    id SERIAL PRIMARY KEY,
    listing_id INTEGER NOT NULL,
    action VARCHAR(20) NOT NULL CHECK (action IN ('index', 'delete')),
    data JSONB,
    attempts INTEGER NOT NULL DEFAULT 0,
    last_error TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    last_attempt_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE (listing_id, action)
);

-- Index для быстрого поиска failed tasks
CREATE INDEX idx_opensearch_dlq_created_at ON opensearch_indexing_dlq(created_at);
CREATE INDEX idx_opensearch_dlq_attempts ON opensearch_indexing_dlq(attempts);
CREATE INDEX idx_opensearch_dlq_listing_id ON opensearch_indexing_dlq(listing_id);

-- Комментарии для документации
COMMENT ON TABLE opensearch_indexing_dlq IS 'Dead Letter Queue for failed OpenSearch indexing tasks';
COMMENT ON COLUMN opensearch_indexing_dlq.listing_id IS 'ID of the listing that failed to index';
COMMENT ON COLUMN opensearch_indexing_dlq.action IS 'Type of indexing action: index or delete';
COMMENT ON COLUMN opensearch_indexing_dlq.data IS 'JSON data of the listing (for index action)';
COMMENT ON COLUMN opensearch_indexing_dlq.attempts IS 'Number of retry attempts';
COMMENT ON COLUMN opensearch_indexing_dlq.last_error IS 'Last error message';
COMMENT ON COLUMN opensearch_indexing_dlq.created_at IS 'When the task first failed';
COMMENT ON COLUMN opensearch_indexing_dlq.last_attempt_at IS 'Last retry attempt timestamp';
