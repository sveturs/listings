-- Rollback: Remove translations support from listings table

-- Drop CHECK constraint
ALTER TABLE listings DROP CONSTRAINT IF EXISTS chk_original_language;

-- Drop indexes
DROP INDEX IF EXISTS idx_listings_title_translations;
DROP INDEX IF EXISTS idx_listings_description_translations;
DROP INDEX IF EXISTS idx_listings_original_language;

-- Drop columns
ALTER TABLE listings
DROP COLUMN IF EXISTS title_translations,
DROP COLUMN IF EXISTS description_translations,
DROP COLUMN IF EXISTS location_translations,
DROP COLUMN IF EXISTS city_translations,
DROP COLUMN IF EXISTS country_translations,
DROP COLUMN IF EXISTS original_language;
