-- Drop trigger and function
DROP TRIGGER IF EXISTS trigger_update_listing_search_vector ON marketplace_listings;
DROP FUNCTION IF EXISTS update_listing_search_vector();

-- Drop search vector column
ALTER TABLE marketplace_listings DROP COLUMN IF EXISTS search_vector;

-- Drop indexes
DROP INDEX IF EXISTS idx_marketplace_listings_search_vector;
DROP INDEX IF EXISTS idx_marketplace_listings_description_fts_en;
DROP INDEX IF EXISTS idx_marketplace_listings_title_fts_en;
DROP INDEX IF EXISTS idx_marketplace_listings_description_fts_ru;
DROP INDEX IF EXISTS idx_marketplace_listings_title_fts_ru;
DROP INDEX IF EXISTS idx_marketplace_listings_description_unaccent_trgm;
DROP INDEX IF EXISTS idx_marketplace_listings_title_unaccent_trgm;
DROP INDEX IF EXISTS idx_marketplace_listings_description_trgm;
DROP INDEX IF EXISTS idx_marketplace_listings_title_trgm;

-- Drop text search configurations
DROP TEXT SEARCH CONFIGURATION IF EXISTS english_unaccent;
DROP TEXT SEARCH CONFIGURATION IF EXISTS russian_unaccent;

-- Drop unaccent dictionary
DROP TEXT SEARCH DICTIONARY IF EXISTS unaccent_dict;

-- Note: We don't drop the extensions as they might be used by other parts of the system
-- If you need to drop them, uncomment these lines:
-- DROP EXTENSION IF EXISTS unaccent;
-- DROP EXTENSION IF EXISTS pg_trgm;