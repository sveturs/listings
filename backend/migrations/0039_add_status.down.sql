DROP INDEX IF EXISTS idx_marketplace_listings_status;
ALTER TABLE marketplace_listings DROP COLUMN IF EXISTS status;