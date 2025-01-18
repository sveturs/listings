-- /backend/migrations/0036_add_translations.down.sql
DROP TRIGGER IF EXISTS update_translations_timestamp ON translations;
DROP FUNCTION IF EXISTS update_translations_updated_at();
DROP TABLE IF EXISTS translations;
ALTER TABLE marketplace_listings DROP COLUMN IF EXISTS original_language;
ALTER TABLE reviews DROP COLUMN IF EXISTS original_language;
ALTER TABLE marketplace_messages DROP COLUMN IF EXISTS original_language;