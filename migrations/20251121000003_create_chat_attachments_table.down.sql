-- =====================================================
-- Migration: 20251121000003_create_chat_attachments_table.down.sql
-- Description: Rollback chat_attachments table creation
-- Author: Chat Microservice Phase 1 - Database Setup
-- Date: 2025-11-21
-- =====================================================

-- Drop helper functions
DROP FUNCTION IF EXISTS is_valid_file_size(TEXT, BIGINT);
DROP FUNCTION IF EXISTS get_file_type_from_content_type(TEXT);

-- Drop triggers
DROP TRIGGER IF EXISTS trigger_update_message_attachments_count_delete ON chat_attachments;
DROP TRIGGER IF EXISTS trigger_update_message_attachments_count_insert ON chat_attachments;

-- Drop trigger function
DROP FUNCTION IF EXISTS update_message_attachments_count();

-- Drop indexes (will be dropped automatically with table, but explicit for clarity)
DROP INDEX IF EXISTS idx_attachments_metadata;
DROP INDEX IF EXISTS idx_attachments_file_size;
DROP INDEX IF EXISTS idx_attachments_storage;
DROP INDEX IF EXISTS idx_attachments_created_at;
DROP INDEX IF EXISTS idx_attachments_file_type;
DROP INDEX IF EXISTS idx_attachments_message_id;

-- Drop table
DROP TABLE IF EXISTS chat_attachments CASCADE;
