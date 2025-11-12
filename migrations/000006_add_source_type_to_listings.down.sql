-- Rollback migration: Remove source_type field from listings table

-- Drop index
DROP INDEX IF EXISTS idx_listings_source_type;

-- Drop check constraint
ALTER TABLE listings
DROP CONSTRAINT IF EXISTS listings_source_type_check;

-- Drop column
ALTER TABLE listings
DROP COLUMN IF EXISTS source_type;
