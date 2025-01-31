DROP TABLE IF EXISTS translations;
ALTER TABLE marketplace_listings DROP COLUMN IF EXISTS original_language;
DROP FUNCTION IF EXISTS update_translations_updated_at();