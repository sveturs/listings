-- Remove the index
DROP INDEX IF EXISTS idx_marketplace_listings_address_multilingual;

-- Remove the column
ALTER TABLE marketplace_listings
DROP COLUMN IF EXISTS address_multilingual;