-- Rollback real AI translations for marketplace listings
-- Created: 2025-08-28

-- Revert translations back to placeholder versions for these listings
UPDATE translations 
SET 
  translated_text = '[' || UPPER(language) || '] ' || (
    SELECT CASE 
      WHEN field_name = 'title' THEN ml.title
      WHEN field_name = 'description' THEN ml.description
    END
    FROM marketplace_listings ml
    WHERE ml.id = translations.entity_id
  ),
  is_machine_translated = true,
  is_verified = false,
  updated_at = CURRENT_TIMESTAMP
WHERE entity_type = 'listing'
  AND entity_id IN (183, 250, 251, 252, 253, 254, 270)
  AND language IN ('ru', 'en')
  AND field_name IN ('title', 'description');