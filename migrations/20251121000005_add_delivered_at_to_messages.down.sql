-- =====================================================
-- Migration: 20251121000005_add_delivered_at_to_messages.down.sql
-- Description: Remove delivered_at timestamp from messages
-- Author: Phase 32 - Presence & Message Status System
-- Date: 2025-11-21
-- =====================================================

-- Drop trigger
DROP TRIGGER IF EXISTS trigger_messages_delivered_at ON messages;

-- Drop function
DROP FUNCTION IF EXISTS update_messages_delivered_at();

-- Drop indexes
DROP INDEX IF EXISTS idx_messages_receiver_undelivered;
DROP INDEX IF EXISTS idx_messages_delivered_at;

-- Remove column
ALTER TABLE messages DROP COLUMN IF EXISTS delivered_at;
