-- Drop functions
DROP FUNCTION IF EXISTS refresh_category_stats();
DROP FUNCTION IF EXISTS get_category_path_names(INTEGER, VARCHAR);
DROP FUNCTION IF EXISTS search_categories_fuzzy(TEXT, VARCHAR, FLOAT);

-- Drop materialized view
DROP MATERIALIZED VIEW IF EXISTS category_listing_stats;

-- Drop indexes on marketplace_listing_attributes
DROP INDEX IF EXISTS idx_listing_attributes_value_text_unaccent_trgm;
DROP INDEX IF EXISTS idx_listing_attributes_value_text_trgm;

-- Drop indexes on translations
DROP INDEX IF EXISTS idx_translations_value_fts_en;
DROP INDEX IF EXISTS idx_translations_value_fts_ru;
DROP INDEX IF EXISTS idx_translations_value_unaccent_trgm;
DROP INDEX IF EXISTS idx_translations_value_trgm;