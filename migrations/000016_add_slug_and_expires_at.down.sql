-- Migration Rollback: Remove slug and expires_at fields from listings table
-- Phase: 13.1.15.4 - Repository Layer Merge

-- Remove indexes first (to avoid dependency errors)
DROP INDEX IF EXISTS idx_listings_status_expires_at;
DROP INDEX IF EXISTS idx_listings_expires_at;
DROP INDEX IF EXISTS idx_listings_slug_all;
DROP INDEX IF EXISTS idx_listings_slug;

-- Remove columns
ALTER TABLE listings DROP COLUMN IF EXISTS expires_at;
ALTER TABLE listings DROP COLUMN IF EXISTS slug;

-- Note: This migration is DESTRUCTIVE
-- If rolled back, all slug data will be lost
-- Consider backing up slug values before rollback if needed
