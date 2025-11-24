-- =====================================================
-- Migration: 20251121000001_create_chats_table.down.sql
-- Description: Rollback chats table creation
-- Author: Chat Microservice Phase 1 - Database Setup
-- Date: 2025-11-21
-- =====================================================

-- Drop trigger and function
DROP TRIGGER IF EXISTS trigger_chats_updated_at ON chats;
DROP FUNCTION IF EXISTS update_chats_updated_at();

-- Drop indexes (will be dropped automatically with table, but explicit for clarity)
DROP INDEX IF EXISTS idx_chats_seller_active_recent;
DROP INDEX IF EXISTS idx_chats_buyer_active_recent;
DROP INDEX IF EXISTS idx_chats_is_archived;
DROP INDEX IF EXISTS idx_chats_last_message_at;
DROP INDEX IF EXISTS idx_chats_storefront_product_id;
DROP INDEX IF EXISTS idx_chats_listing_id;
DROP INDEX IF EXISTS idx_chats_seller_id;
DROP INDEX IF EXISTS idx_chats_buyer_id;

-- Drop table
DROP TABLE IF EXISTS chats CASCADE;
