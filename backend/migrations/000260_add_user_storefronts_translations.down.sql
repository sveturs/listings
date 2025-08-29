-- Rollback migration for user storefronts translations
-- Created: 2025-08-28

-- Remove all translations for user storefronts that were added by this migration
-- We identify them by being machine translated and created during this migration

DELETE FROM translations 
WHERE entity_type = 'storefront'
  AND language IN ('ru', 'en')
  AND field_name IN ('name', 'description')
  AND is_machine_translated = true
  AND created_at >= (
    SELECT CURRENT_TIMESTAMP - INTERVAL '1 hour'
  );