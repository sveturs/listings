-- Add is_system column to messages table for system/marketplace notifications
ALTER TABLE messages ADD COLUMN IF NOT EXISTS is_system BOOLEAN NOT NULL DEFAULT false;

-- Create index for filtering system messages
CREATE INDEX IF NOT EXISTS idx_messages_is_system ON messages(is_system) WHERE is_system = true;

-- Add comment
COMMENT ON COLUMN messages.is_system IS 'Flag indicating if message is from the marketplace system (e.g., order notifications)';
