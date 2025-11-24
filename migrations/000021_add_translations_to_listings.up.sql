-- Add translations support to listings table
-- This migration adds JSONB columns for storing field-level translations

-- Add translation columns for individual fields
ALTER TABLE listings
ADD COLUMN title_translations JSONB DEFAULT '{}',
ADD COLUMN description_translations JSONB DEFAULT '{}',
ADD COLUMN location_translations JSONB DEFAULT '{}',
ADD COLUMN city_translations JSONB DEFAULT '{}',
ADD COLUMN country_translations JSONB DEFAULT '{}',
ADD COLUMN original_language VARCHAR(10) DEFAULT 'sr';

-- Create GIN indexes for efficient JSONB queries
CREATE INDEX idx_listings_title_translations ON listings USING gin (title_translations);
CREATE INDEX idx_listings_description_translations ON listings USING gin (description_translations);
CREATE INDEX idx_listings_original_language ON listings (original_language);

-- Add comment for documentation
COMMENT ON COLUMN listings.title_translations IS 'JSONB map of language codes to translated titles, e.g., {"en": "Title", "ru": "Заголовок", "sr": "Наслов"}';
COMMENT ON COLUMN listings.description_translations IS 'JSONB map of language codes to translated descriptions';
COMMENT ON COLUMN listings.location_translations IS 'JSONB map of language codes to translated location';
COMMENT ON COLUMN listings.city_translations IS 'JSONB map of language codes to translated city names';
COMMENT ON COLUMN listings.country_translations IS 'JSONB map of language codes to translated country names';
COMMENT ON COLUMN listings.original_language IS 'Original language code (sr, en, ru) in which the listing was created';

-- Add CHECK constraint to ensure original_language is valid
ALTER TABLE listings ADD CONSTRAINT chk_original_language CHECK (original_language IN ('sr', 'en', 'ru'));
