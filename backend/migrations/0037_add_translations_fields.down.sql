-- /backend/migrations/000X_add_translations_fields.down.sql
ALTER TABLE reviews DROP COLUMN IF EXISTS original_language;
ALTER TABLE marketplace_messages DROP COLUMN IF EXISTS original_language;