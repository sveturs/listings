-- Rollback: Remove b2c_products compatibility columns

-- Drop indexes
DROP INDEX IF EXISTS idx_listings_has_variants;
DROP INDEX IF EXISTS idx_listings_location;
DROP INDEX IF EXISTS idx_listings_sold_count;
DROP INDEX IF EXISTS idx_listings_view_count;

-- Drop columns (in reverse order)
ALTER TABLE listings DROP COLUMN IF EXISTS has_variants;
ALTER TABLE listings DROP COLUMN IF EXISTS show_on_map;
ALTER TABLE listings DROP COLUMN IF EXISTS location_privacy;
ALTER TABLE listings DROP COLUMN IF EXISTS individual_longitude;
ALTER TABLE listings DROP COLUMN IF EXISTS individual_latitude;
ALTER TABLE listings DROP COLUMN IF EXISTS individual_address;
ALTER TABLE listings DROP COLUMN IF EXISTS has_individual_location;
ALTER TABLE listings DROP COLUMN IF EXISTS sold_count;

-- Rename back view_count → views_count
ALTER TABLE listings RENAME COLUMN view_count TO views_count;
