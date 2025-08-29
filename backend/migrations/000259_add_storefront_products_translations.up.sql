-- Migration to add Russian and English translations for storefront products
-- Created: 2025-08-28

-- Insert translations for storefront products names and descriptions
-- that don't have translations yet

-- Function to translate text using DeepL API (placeholder implementation)
-- In production, this should be integrated with actual translation service
CREATE OR REPLACE FUNCTION translate_product_content(original_text TEXT, source_lang TEXT, target_lang TEXT) 
RETURNS TEXT AS $$
BEGIN
    -- Temporary placeholder - will be replaced with actual translation service call
    -- For now, return a prefixed version to indicate it's a translation
    RETURN '[' || upper(target_lang) || '] ' || original_text;
END;
$$ LANGUAGE plpgsql;

-- Insert Russian translations for storefront products names
INSERT INTO translations (entity_type, entity_id, language, field_name, translated_text, is_machine_translated, is_verified)
SELECT 
    'storefront_product'::text,
    sp.id,
    'ru'::text,
    'name'::text,
    translate_product_content(sp.name, 'sr', 'ru'),
    true,
    false
FROM storefront_products sp
WHERE sp.is_active = true
  AND NOT EXISTS (
    SELECT 1 FROM translations t 
    WHERE t.entity_type = 'storefront_product' 
      AND t.entity_id = sp.id 
      AND t.language = 'ru' 
      AND t.field_name = 'name'
  );

-- Insert Russian translations for storefront products descriptions
INSERT INTO translations (entity_type, entity_id, language, field_name, translated_text, is_machine_translated, is_verified)
SELECT 
    'storefront_product'::text,
    sp.id,
    'ru'::text,
    'description'::text,
    translate_product_content(sp.description, 'sr', 'ru'),
    true,
    false
FROM storefront_products sp
WHERE sp.is_active = true
  AND sp.description IS NOT NULL
  AND sp.description != ''
  AND NOT EXISTS (
    SELECT 1 FROM translations t 
    WHERE t.entity_type = 'storefront_product' 
      AND t.entity_id = sp.id 
      AND t.language = 'ru' 
      AND t.field_name = 'description'
  );

-- Insert English translations for storefront products names
INSERT INTO translations (entity_type, entity_id, language, field_name, translated_text, is_machine_translated, is_verified)
SELECT 
    'storefront_product'::text,
    sp.id,
    'en'::text,
    'name'::text,
    translate_product_content(sp.name, 'sr', 'en'),
    true,
    false
FROM storefront_products sp
WHERE sp.is_active = true
  AND NOT EXISTS (
    SELECT 1 FROM translations t 
    WHERE t.entity_type = 'storefront_product' 
      AND t.entity_id = sp.id 
      AND t.language = 'en' 
      AND t.field_name = 'name'
  );

-- Insert English translations for storefront products descriptions
INSERT INTO translations (entity_type, entity_id, language, field_name, translated_text, is_machine_translated, is_verified)
SELECT 
    'storefront_product'::text,
    sp.id,
    'en'::text,
    'description'::text,
    translate_product_content(sp.description, 'sr', 'en'),
    true,
    false
FROM storefront_products sp
WHERE sp.is_active = true
  AND sp.description IS NOT NULL
  AND sp.description != ''
  AND NOT EXISTS (
    SELECT 1 FROM translations t 
    WHERE t.entity_type = 'storefront_product' 
      AND t.entity_id = sp.id 
      AND t.language = 'en' 
      AND t.field_name = 'description'
  );

-- Clean up the temporary function
DROP FUNCTION translate_product_content(TEXT, TEXT, TEXT);