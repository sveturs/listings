-- Revert translations column changes
DROP INDEX IF EXISTS idx_marketplace_messages_translations;

ALTER TABLE marketplace_messages
DROP COLUMN IF EXISTS translations;

-- Revert original_language type back to VARCHAR(2)
ALTER TABLE marketplace_messages
ALTER COLUMN original_language TYPE VARCHAR(2);
