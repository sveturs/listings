-- Rollback: Drop Dead Letter Queue table
DROP INDEX IF EXISTS idx_opensearch_dlq_listing_id;
DROP INDEX IF EXISTS idx_opensearch_dlq_attempts;
DROP INDEX IF EXISTS idx_opensearch_dlq_created_at;
DROP TABLE IF EXISTS opensearch_indexing_dlq;
