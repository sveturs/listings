-- Rollback: Remove attributes column from listings table

ALTER TABLE listings
DROP COLUMN IF EXISTS attributes;

DROP INDEX IF EXISTS idx_listings_attributes;
