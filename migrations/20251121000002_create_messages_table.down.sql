-- =====================================================
-- Migration: 20251121000002_create_messages_table.down.sql
-- Description: Rollback messages table creation
-- Author: Chat Microservice Phase 1 - Database Setup
-- Date: 2025-11-21
-- =====================================================

-- Drop triggers
DROP TRIGGER IF EXISTS trigger_update_chat_last_message ON messages;
DROP TRIGGER IF EXISTS trigger_messages_read_at ON messages;
DROP TRIGGER IF EXISTS trigger_messages_updated_at ON messages;

-- Drop trigger functions
DROP FUNCTION IF EXISTS update_chat_last_message_at();
DROP FUNCTION IF EXISTS update_messages_read_at();
DROP FUNCTION IF EXISTS update_messages_updated_at();

-- Drop indexes (will be dropped automatically with table, but explicit for clarity)
DROP INDEX IF EXISTS idx_messages_has_attachments;
DROP INDEX IF EXISTS idx_messages_created_at;
DROP INDEX IF EXISTS idx_messages_chat_receiver_unread;
DROP INDEX IF EXISTS idx_messages_receiver_unread;
DROP INDEX IF EXISTS idx_messages_receiver_id;
DROP INDEX IF EXISTS idx_messages_sender_id;
DROP INDEX IF EXISTS idx_messages_chat_id_id;
DROP INDEX IF EXISTS idx_messages_chat_id_created;

-- Drop table
DROP TABLE IF EXISTS messages CASCADE;
