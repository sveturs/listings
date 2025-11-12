-- Migration: 000018_add_indexing_queue_unique_constraint (ROLLBACK)
-- Description: Remove partial unique index from indexing_queue
-- Date: 2025-11-09

-- Drop partial unique index
DROP INDEX IF EXISTS idx_indexing_queue_listing_id_pending;
