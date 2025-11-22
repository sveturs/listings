-- Remove is_system column from messages table
DROP INDEX IF EXISTS idx_messages_is_system;
ALTER TABLE messages DROP COLUMN IF EXISTS is_system;
