-- Migration: 000018_add_indexing_queue_unique_constraint
-- Description: Add partial unique index for indexing_queue to support ON CONFLICT
-- Date: 2025-11-09
-- Author: Test fix
--
-- This migration adds a partial unique index on (listing_id) WHERE status = 'pending'
-- to support the ON CONFLICT clause in EnqueueIndexing.

-- Create partial unique index for pending indexing jobs
-- This allows ON CONFLICT (listing_id) WHERE status = 'pending' to work
CREATE UNIQUE INDEX idx_indexing_queue_listing_id_pending
ON indexing_queue(listing_id)
WHERE status = 'pending';

-- Add comment
COMMENT ON INDEX idx_indexing_queue_listing_id_pending IS 'Ensures only one pending indexing job per listing at a time';
