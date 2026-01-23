-- Remove UNIQUE constraint on listings.sku

DROP INDEX IF EXISTS unique_listing_sku;
