-- Rollback: Remove stock_status column from listings table

-- Drop index
DROP INDEX IF EXISTS idx_listings_stock_status;

-- Drop column
ALTER TABLE listings
DROP COLUMN IF EXISTS stock_status;
