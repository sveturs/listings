-- =====================================================
-- Migration: 20251121000005_add_delivered_at_to_messages.up.sql
-- Description: Add delivered_at timestamp for message delivery tracking
-- Author: Phase 32 - Presence & Message Status System
-- Date: 2025-11-21
-- =====================================================

-- Add delivered_at column for tracking when message was delivered to recipient
ALTER TABLE messages
ADD COLUMN IF NOT EXISTS delivered_at TIMESTAMP;

-- Create index for efficient delivery status queries
CREATE INDEX IF NOT EXISTS idx_messages_delivered_at
    ON messages(delivered_at)
    WHERE delivered_at IS NOT NULL;

-- Create composite index for getting undelivered messages per user
CREATE INDEX IF NOT EXISTS idx_messages_receiver_undelivered
    ON messages(receiver_id, status)
    WHERE status = 'sent';

-- =====================================================
-- Trigger: Update delivered_at when status changes to delivered
-- =====================================================

CREATE OR REPLACE FUNCTION update_messages_delivered_at()
RETURNS TRIGGER AS $$
BEGIN
    -- When status changes to 'delivered' and delivered_at is not set
    IF NEW.status = 'delivered' AND OLD.status = 'sent' AND NEW.delivered_at IS NULL THEN
        NEW.delivered_at = CURRENT_TIMESTAMP;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Drop trigger if exists (for idempotency)
DROP TRIGGER IF EXISTS trigger_messages_delivered_at ON messages;

CREATE TRIGGER trigger_messages_delivered_at
    BEFORE UPDATE ON messages
    FOR EACH ROW
    WHEN (NEW.status IS DISTINCT FROM OLD.status)
    EXECUTE FUNCTION update_messages_delivered_at();

-- =====================================================
-- Comments for Documentation
-- =====================================================

COMMENT ON COLUMN messages.delivered_at IS 'Timestamp when message was delivered to recipient device (null if not yet delivered)';

COMMENT ON FUNCTION update_messages_delivered_at IS
'Automatically set delivered_at timestamp when message status changes to delivered';
