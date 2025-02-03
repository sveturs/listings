--  backend/migrations/0037_add_translations_fields.up.sql
ALTER TABLE reviews ADD COLUMN IF NOT EXISTS original_language VARCHAR(2) DEFAULT 'en';
ALTER TABLE marketplace_messages ADD COLUMN IF NOT EXISTS original_language VARCHAR(2) DEFAULT 'en';
