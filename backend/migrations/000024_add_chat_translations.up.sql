-- Add translations JSONB column to marketplace_messages
-- This column will store translations in format: {"en": "Hello", "ru": "Привет", "sr": "Здраво"}
ALTER TABLE marketplace_messages
ADD COLUMN IF NOT EXISTS translations JSONB DEFAULT '{}'::jsonb;

-- Create GIN index for fast JSON queries
CREATE INDEX IF NOT EXISTS idx_marketplace_messages_translations
ON marketplace_messages USING gin(translations);

-- Update original_language to support more languages (not just 2 chars)
ALTER TABLE marketplace_messages
ALTER COLUMN original_language TYPE VARCHAR(10);

-- Add comments for documentation
COMMENT ON COLUMN marketplace_messages.translations IS
'JSON object storing message translations: {"en": "Hello", "ru": "Привет", "sr": "Здраво"}';

COMMENT ON COLUMN marketplace_messages.original_language IS
'ISO 639-1 language code detected from message content (ru, en, sr, etc.)';
