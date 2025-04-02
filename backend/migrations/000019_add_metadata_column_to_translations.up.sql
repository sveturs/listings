-- Add metadata column to translations table
ALTER TABLE translations ADD COLUMN metadata JSONB DEFAULT '{}';

-- Update existing records to have empty metadata
UPDATE translations SET metadata = '{}' WHERE metadata IS NULL;

-- Add an index on the metadata column for faster querying
CREATE INDEX idx_translations_metadata ON translations USING GIN (metadata);

COMMENT ON COLUMN translations.metadata IS 'JSON field to store metadata about the translation, such as provider used and other settings';