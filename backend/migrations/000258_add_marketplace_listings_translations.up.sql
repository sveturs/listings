-- Migration to add Russian and English translations for marketplace listings
-- Created: 2025-08-28

-- Insert translations for marketplace listings titles and descriptions
-- that don't have translations yet

-- Function to translate text using DeepL API (placeholder implementation)
-- In production, this should be integrated with actual translation service
CREATE OR REPLACE FUNCTION translate_listing_content(original_text TEXT, source_lang TEXT, target_lang TEXT) 
RETURNS TEXT AS $$
BEGIN
    -- Temporary placeholder - will be replaced with actual translation service call
    -- For now, return a prefixed version to indicate it's a translation
    RETURN '[' || upper(target_lang) || '] ' || original_text;
END;
$$ LANGUAGE plpgsql;

-- Insert Russian translations for marketplace listings titles
INSERT INTO translations (entity_type, entity_id, language, field_name, translated_text, is_machine_translated, is_verified)
SELECT 
    'listing'::text,
    ml.id,
    'ru'::text,
    'title'::text,
    translate_listing_content(ml.title, COALESCE(ml.original_language, 'sr'), 'ru'),
    true,
    false
FROM marketplace_listings ml
WHERE ml.status = 'active'
  AND NOT EXISTS (
    SELECT 1 FROM translations t 
    WHERE t.entity_type = 'listing' 
      AND t.entity_id = ml.id 
      AND t.language = 'ru' 
      AND t.field_name = 'title'
  );

-- Insert Russian translations for marketplace listings descriptions
INSERT INTO translations (entity_type, entity_id, language, field_name, translated_text, is_machine_translated, is_verified)
SELECT 
    'listing'::text,
    ml.id,
    'ru'::text,
    'description'::text,
    translate_listing_content(ml.description, COALESCE(ml.original_language, 'sr'), 'ru'),
    true,
    false
FROM marketplace_listings ml
WHERE ml.status = 'active'
  AND ml.description IS NOT NULL
  AND ml.description != ''
  AND NOT EXISTS (
    SELECT 1 FROM translations t 
    WHERE t.entity_type = 'listing' 
      AND t.entity_id = ml.id 
      AND t.language = 'ru' 
      AND t.field_name = 'description'
  );

-- Insert English translations for marketplace listings titles
INSERT INTO translations (entity_type, entity_id, language, field_name, translated_text, is_machine_translated, is_verified)
SELECT 
    'listing'::text,
    ml.id,
    'en'::text,
    'title'::text,
    translate_listing_content(ml.title, COALESCE(ml.original_language, 'sr'), 'en'),
    true,
    false
FROM marketplace_listings ml
WHERE ml.status = 'active'
  AND NOT EXISTS (
    SELECT 1 FROM translations t 
    WHERE t.entity_type = 'listing' 
      AND t.entity_id = ml.id 
      AND t.language = 'en' 
      AND t.field_name = 'title'
  );

-- Insert English translations for marketplace listings descriptions
INSERT INTO translations (entity_type, entity_id, language, field_name, translated_text, is_machine_translated, is_verified)
SELECT 
    'listing'::text,
    ml.id,
    'en'::text,
    'description'::text,
    translate_listing_content(ml.description, COALESCE(ml.original_language, 'sr'), 'en'),
    true,
    false
FROM marketplace_listings ml
WHERE ml.status = 'active'
  AND ml.description IS NOT NULL
  AND ml.description != ''
  AND NOT EXISTS (
    SELECT 1 FROM translations t 
    WHERE t.entity_type = 'listing' 
      AND t.entity_id = ml.id 
      AND t.language = 'en' 
      AND t.field_name = 'description'
  );

-- Clean up the temporary function
DROP FUNCTION translate_listing_content(TEXT, TEXT, TEXT);