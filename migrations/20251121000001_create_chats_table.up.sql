-- =====================================================
-- Migration: 20251121000001_create_chats_table.up.sql
-- Description: Create chats table for Chat microservice
-- Author: Chat Microservice Phase 1 - Database Setup
-- Date: 2025-11-21
-- =====================================================

-- =====================================================
-- Table: chats
-- Description: Stores conversations between buyers and sellers
-- =====================================================
CREATE TABLE IF NOT EXISTS chats (
    -- Identification
    id BIGSERIAL PRIMARY KEY,
    buyer_id BIGINT NOT NULL,
    seller_id BIGINT NOT NULL,

    -- Context (what is being discussed)
    listing_id BIGINT,
    storefront_product_id BIGINT,

    -- Status
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    is_archived BOOLEAN NOT NULL DEFAULT false,

    -- Metadata
    last_message_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    -- =====================================================
    -- Constraints
    -- =====================================================

    -- Ensure only one context type is set (listing OR product OR direct chat)
    CONSTRAINT check_chat_context CHECK (
        (listing_id IS NOT NULL AND storefront_product_id IS NULL) OR
        (listing_id IS NULL AND storefront_product_id IS NOT NULL) OR
        (listing_id IS NULL AND storefront_product_id IS NULL)
    ),

    -- Prevent self-chat
    CONSTRAINT check_participants CHECK (buyer_id != seller_id),

    -- Unique constraints to prevent duplicate chats
    CONSTRAINT chats_listing_participants_unique UNIQUE (listing_id, buyer_id, seller_id),
    CONSTRAINT chats_product_participants_unique UNIQUE (storefront_product_id, buyer_id, seller_id),

    -- Status validation
    CONSTRAINT check_chat_status CHECK (status IN ('active', 'archived', 'blocked'))
);

-- =====================================================
-- Indexes for Performance
-- =====================================================

-- Buyer's chats list
CREATE INDEX idx_chats_buyer_id ON chats(buyer_id);

-- Seller's chats list
CREATE INDEX idx_chats_seller_id ON chats(seller_id);

-- Find chat by listing
CREATE INDEX idx_chats_listing_id ON chats(listing_id) WHERE listing_id IS NOT NULL;

-- Find chat by storefront product
CREATE INDEX idx_chats_storefront_product_id ON chats(storefront_product_id) WHERE storefront_product_id IS NOT NULL;

-- Sort chats by last message (exclude archived for performance)
CREATE INDEX idx_chats_last_message_at ON chats(last_message_at DESC) WHERE NOT is_archived;

-- Filter archived chats efficiently
CREATE INDEX idx_chats_is_archived ON chats(is_archived) WHERE NOT is_archived;

-- Composite index for buyer's active chats sorted by last message
CREATE INDEX idx_chats_buyer_active_recent ON chats(buyer_id, last_message_at DESC)
    WHERE NOT is_archived AND status = 'active';

-- Composite index for seller's active chats sorted by last message
CREATE INDEX idx_chats_seller_active_recent ON chats(seller_id, last_message_at DESC)
    WHERE NOT is_archived AND status = 'active';

-- Unique index for direct chats (prevents duplicate direct chats between same users)
CREATE UNIQUE INDEX idx_chats_direct_participants_unique ON chats (
    LEAST(buyer_id, seller_id),
    GREATEST(buyer_id, seller_id)
) WHERE listing_id IS NULL AND storefront_product_id IS NULL;

-- =====================================================
-- Comments for Documentation
-- =====================================================

COMMENT ON TABLE chats IS
'Conversations between buyers and sellers. Can be about a specific listing, storefront product, or direct chat';

COMMENT ON COLUMN chats.id IS 'Primary key';
COMMENT ON COLUMN chats.buyer_id IS 'User who initiated the chat (references users in Auth Service)';
COMMENT ON COLUMN chats.seller_id IS 'User who receives the chat (references users in Auth Service)';
COMMENT ON COLUMN chats.listing_id IS 'Optional: Marketplace listing being discussed (references listings table)';
COMMENT ON COLUMN chats.storefront_product_id IS 'Optional: Storefront product being discussed (references storefront_products table)';
COMMENT ON COLUMN chats.status IS 'Chat status: active, archived, or blocked';
COMMENT ON COLUMN chats.is_archived IS 'Whether chat is archived by current user (UI state)';
COMMENT ON COLUMN chats.last_message_at IS 'Timestamp of last message (for sorting)';
COMMENT ON COLUMN chats.created_at IS 'When chat was created';
COMMENT ON COLUMN chats.updated_at IS 'Last update timestamp';

-- =====================================================
-- Trigger: Update updated_at on row change
-- =====================================================

CREATE OR REPLACE FUNCTION update_chats_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_chats_updated_at
    BEFORE UPDATE ON chats
    FOR EACH ROW
    EXECUTE FUNCTION update_chats_updated_at();

COMMENT ON FUNCTION update_chats_updated_at IS
'Automatically update updated_at timestamp when chat is modified';
