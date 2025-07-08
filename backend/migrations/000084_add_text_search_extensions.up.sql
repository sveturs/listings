-- Enable PostgreSQL extensions for text search
-- pg_trgm: for trigram-based similarity search (fuzzy matching)
-- unaccent: for removing accents from text (useful for multi-language search)

-- Enable pg_trgm extension for trigram similarity search
CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- Enable unaccent extension for accent-insensitive search
CREATE EXTENSION IF NOT EXISTS unaccent;

-- Create immutable wrapper for unaccent function
-- PostgreSQL's unaccent function is not marked as IMMUTABLE, but we need it for indexes
CREATE OR REPLACE FUNCTION f_unaccent(text) 
RETURNS text AS $$
    SELECT unaccent($1)
$$ LANGUAGE sql IMMUTABLE PARALLEL SAFE;

-- Create unaccent dictionary for text search
-- This allows us to use unaccent in full-text search configurations
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_ts_dict WHERE dictname = 'unaccent_dict') THEN
        CREATE TEXT SEARCH DICTIONARY unaccent_dict (
            TEMPLATE = unaccent,
            RULES = 'unaccent'
        );
    END IF;
END$$;

-- Create custom text search configuration for Russian language with unaccent
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_ts_config WHERE cfgname = 'russian_unaccent') THEN
        CREATE TEXT SEARCH CONFIGURATION russian_unaccent (COPY = russian);
        ALTER TEXT SEARCH CONFIGURATION russian_unaccent
            ALTER MAPPING FOR hword, hword_part, word
            WITH unaccent_dict, russian_stem;
    END IF;
END$$;

-- Create custom text search configuration for English language with unaccent
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_ts_config WHERE cfgname = 'english_unaccent') THEN
        CREATE TEXT SEARCH CONFIGURATION english_unaccent (COPY = english);
        ALTER TEXT SEARCH CONFIGURATION english_unaccent
            ALTER MAPPING FOR hword, hword_part, word
            WITH unaccent_dict, english_stem;
    END IF;
END$$;

-- Create GIN indexes for trigram search on marketplace_listings
CREATE INDEX IF NOT EXISTS idx_marketplace_listings_title_trgm 
    ON marketplace_listings USING gin (title gin_trgm_ops);

CREATE INDEX IF NOT EXISTS idx_marketplace_listings_description_trgm 
    ON marketplace_listings USING gin (description gin_trgm_ops);

-- Create functional GIN indexes for unaccented search
CREATE INDEX IF NOT EXISTS idx_marketplace_listings_title_unaccent_trgm 
    ON marketplace_listings USING gin (f_unaccent(title) gin_trgm_ops);

CREATE INDEX IF NOT EXISTS idx_marketplace_listings_description_unaccent_trgm 
    ON marketplace_listings USING gin (f_unaccent(description) gin_trgm_ops);

-- Create full-text search indexes with unaccent support
CREATE INDEX IF NOT EXISTS idx_marketplace_listings_title_fts_ru 
    ON marketplace_listings USING gin (to_tsvector('russian_unaccent', title));

CREATE INDEX IF NOT EXISTS idx_marketplace_listings_description_fts_ru 
    ON marketplace_listings USING gin (to_tsvector('russian_unaccent', description));

CREATE INDEX IF NOT EXISTS idx_marketplace_listings_title_fts_en 
    ON marketplace_listings USING gin (to_tsvector('english_unaccent', title));

CREATE INDEX IF NOT EXISTS idx_marketplace_listings_description_fts_en 
    ON marketplace_listings USING gin (to_tsvector('english_unaccent', description));

-- Create combined search vector column for better performance
ALTER TABLE marketplace_listings 
    ADD COLUMN IF NOT EXISTS search_vector tsvector;

-- Create trigger to update search vector automatically
CREATE OR REPLACE FUNCTION update_listing_search_vector() 
RETURNS trigger AS $$
BEGIN
    NEW.search_vector := 
        setweight(to_tsvector('russian_unaccent', COALESCE(NEW.title, '')), 'A') ||
        setweight(to_tsvector('russian_unaccent', COALESCE(NEW.description, '')), 'B') ||
        setweight(to_tsvector('english_unaccent', COALESCE(NEW.title, '')), 'A') ||
        setweight(to_tsvector('english_unaccent', COALESCE(NEW.description, '')), 'B');
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger
DROP TRIGGER IF EXISTS trigger_update_listing_search_vector ON marketplace_listings;
DROP TRIGGER IF EXISTS trigger_update_listing_search_vector ON trigger_update_listing_search_vector;
CREATE TRIGGER trigger_update_listing_search_vector
    BEFORE INSERT OR UPDATE OF title, description ON marketplace_listings
    FOR EACH ROW
    EXECUTE FUNCTION update_listing_search_vector();

-- Update existing records
UPDATE marketplace_listings 
SET search_vector = 
    setweight(to_tsvector('russian_unaccent', COALESCE(title, '')), 'A') ||
    setweight(to_tsvector('russian_unaccent', COALESCE(description, '')), 'B') ||
    setweight(to_tsvector('english_unaccent', COALESCE(title, '')), 'A') ||
    setweight(to_tsvector('english_unaccent', COALESCE(description, '')), 'B');

-- Create index on search vector
CREATE INDEX IF NOT EXISTS idx_marketplace_listings_search_vector 
    ON marketplace_listings USING gin (search_vector);

-- Add similarity threshold setting
-- This can be adjusted later for tuning search sensitivity
-- Note: This sets the default similarity threshold for trigram matching
-- Can be overridden per session with: SET pg_trgm.similarity_threshold = 0.3;
DO $$
BEGIN
    EXECUTE format('ALTER DATABASE %I SET pg_trgm.similarity_threshold = 0.3', current_database());
END$$;