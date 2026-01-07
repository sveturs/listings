-- Add UNIQUE constraint on listings.sku to prevent duplicate SKUs
-- This constraint allows NULL values (multiple listings can have NULL sku)

CREATE UNIQUE INDEX unique_listing_sku ON listings(sku) WHERE sku IS NOT NULL;
