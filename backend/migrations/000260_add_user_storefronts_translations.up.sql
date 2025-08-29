-- Migration to add Russian and English translations for user storefronts
-- Created: 2025-08-28

-- Insert translations for user storefronts names and descriptions
-- that don't have translations yet

-- Function to translate text using DeepL API (placeholder implementation)
-- In production, this should be integrated with actual translation service
CREATE OR REPLACE FUNCTION translate_storefront_content(original_text TEXT, source_lang TEXT, target_lang TEXT) 
RETURNS TEXT AS $$
BEGIN
    -- Temporary placeholder - will be replaced with actual translation service call
    -- For now, return a prefixed version to indicate it's a translation
    RETURN '[' || upper(target_lang) || '] ' || original_text;
END;
$$ LANGUAGE plpgsql;

-- Insert Russian translations for user storefronts names
INSERT INTO translations (entity_type, entity_id, language, field_name, translated_text, is_machine_translated, is_verified)
SELECT 
    'storefront'::text,
    us.id,
    'ru'::text,
    'name'::text,
    translate_storefront_content(us.name, 'sr', 'ru'),
    true,
    false
FROM user_storefronts us
WHERE us.status = 'active'
  AND NOT EXISTS (
    SELECT 1 FROM translations t 
    WHERE t.entity_type = 'storefront' 
      AND t.entity_id = us.id 
      AND t.language = 'ru' 
      AND t.field_name = 'name'
  );

-- Insert Russian translations for user storefronts descriptions
INSERT INTO translations (entity_type, entity_id, language, field_name, translated_text, is_machine_translated, is_verified)
SELECT 
    'storefront'::text,
    us.id,
    'ru'::text,
    'description'::text,
    translate_storefront_content(us.description, 'sr', 'ru'),
    true,
    false
FROM user_storefronts us
WHERE us.status = 'active'
  AND us.description IS NOT NULL
  AND us.description != ''
  AND NOT EXISTS (
    SELECT 1 FROM translations t 
    WHERE t.entity_type = 'storefront' 
      AND t.entity_id = us.id 
      AND t.language = 'ru' 
      AND t.field_name = 'description'
  );

-- Insert English translations for user storefronts names
INSERT INTO translations (entity_type, entity_id, language, field_name, translated_text, is_machine_translated, is_verified)
SELECT 
    'storefront'::text,
    us.id,
    'en'::text,
    'name'::text,
    translate_storefront_content(us.name, 'sr', 'en'),
    true,
    false
FROM user_storefronts us
WHERE us.status = 'active'
  AND NOT EXISTS (
    SELECT 1 FROM translations t 
    WHERE t.entity_type = 'storefront' 
      AND t.entity_id = us.id 
      AND t.language = 'en' 
      AND t.field_name = 'name'
  );

-- Insert English translations for user storefronts descriptions
INSERT INTO translations (entity_type, entity_id, language, field_name, translated_text, is_machine_translated, is_verified)
SELECT 
    'storefront'::text,
    us.id,
    'en'::text,
    'description'::text,
    translate_storefront_content(us.description, 'sr', 'en'),
    true,
    false
FROM user_storefronts us
WHERE us.status = 'active'
  AND us.description IS NOT NULL
  AND us.description != ''
  AND NOT EXISTS (
    SELECT 1 FROM translations t 
    WHERE t.entity_type = 'storefront' 
      AND t.entity_id = us.id 
      AND t.language = 'en' 
      AND t.field_name = 'description'
  );

-- Clean up the temporary function
DROP FUNCTION translate_storefront_content(TEXT, TEXT, TEXT);