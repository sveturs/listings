-- =====================================================
-- Migration: 20251121000002_create_messages_table.up.sql
-- Description: Create messages table for Chat microservice
-- Author: Chat Microservice Phase 1 - Database Setup
-- Date: 2025-11-21
-- =====================================================

-- =====================================================
-- Table: messages
-- Description: Stores individual messages in chats
-- =====================================================
CREATE TABLE IF NOT EXISTS messages (
    -- Identification
    id BIGSERIAL PRIMARY KEY,
    chat_id BIGINT NOT NULL,
    sender_id BIGINT NOT NULL,
    receiver_id BIGINT NOT NULL,

    -- Content
    content TEXT NOT NULL,
    original_language VARCHAR(2) NOT NULL DEFAULT 'en',

    -- Context (optional - inherited from chat if not provided)
    listing_id BIGINT,
    storefront_product_id BIGINT,

    -- Status
    status VARCHAR(20) NOT NULL DEFAULT 'sent',
    is_read BOOLEAN NOT NULL DEFAULT false,

    -- Attachments
    has_attachments BOOLEAN NOT NULL DEFAULT false,
    attachments_count INT NOT NULL DEFAULT 0,

    -- Timestamps
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    read_at TIMESTAMP,

    -- =====================================================
    -- Foreign Keys
    -- =====================================================

    CONSTRAINT fk_messages_chat FOREIGN KEY (chat_id)
        REFERENCES chats(id)
        ON DELETE CASCADE,

    -- =====================================================
    -- Constraints
    -- =====================================================

    -- Ensure only one context type is set
    CONSTRAINT check_message_context CHECK (
        (listing_id IS NOT NULL AND storefront_product_id IS NULL) OR
        (listing_id IS NULL AND storefront_product_id IS NOT NULL) OR
        (listing_id IS NULL AND storefront_product_id IS NULL)
    ),

    -- Content must be between 1 and 10000 characters
    CONSTRAINT check_content_length CHECK (
        length(content) >= 1 AND length(content) <= 10000
    ),

    -- Status validation
    CONSTRAINT check_message_status CHECK (
        status IN ('sent', 'delivered', 'read', 'failed')
    ),

    -- Language code validation (ISO 639-1)
    CONSTRAINT check_original_language CHECK (
        original_language ~ '^[a-z]{2}$'
    ),

    -- Attachments count consistency
    CONSTRAINT check_attachments_consistency CHECK (
        (has_attachments = true AND attachments_count > 0) OR
        (has_attachments = false AND attachments_count = 0)
    )
);

-- =====================================================
-- Indexes for Performance
-- =====================================================

-- Get messages in chat ordered by creation time (most common query)
CREATE INDEX idx_messages_chat_id_created ON messages(chat_id, created_at DESC);

-- Get messages in chat ordered by ID (for cursor pagination)
CREATE INDEX idx_messages_chat_id_id ON messages(chat_id, id DESC);

-- Find messages by sender
CREATE INDEX idx_messages_sender_id ON messages(sender_id);

-- Find messages by receiver
CREATE INDEX idx_messages_receiver_id ON messages(receiver_id);

-- Get unread messages for a user (notification badge)
CREATE INDEX idx_messages_receiver_unread ON messages(receiver_id, is_read)
    WHERE NOT is_read;

-- Get unread count per chat for a user
CREATE INDEX idx_messages_chat_receiver_unread ON messages(chat_id, receiver_id)
    WHERE NOT is_read;

-- Time-based queries for analytics
CREATE INDEX idx_messages_created_at ON messages(created_at);

-- Find messages with attachments
CREATE INDEX idx_messages_has_attachments ON messages(chat_id, has_attachments)
    WHERE has_attachments = true;

-- =====================================================
-- Comments for Documentation
-- =====================================================

COMMENT ON TABLE messages IS
'Individual messages in chats between buyers and sellers';

COMMENT ON COLUMN messages.id IS 'Primary key, used for cursor pagination';
COMMENT ON COLUMN messages.chat_id IS 'Reference to parent chat';
COMMENT ON COLUMN messages.sender_id IS 'User who sent the message (references users in Auth Service)';
COMMENT ON COLUMN messages.receiver_id IS 'User who receives the message (references users in Auth Service)';
COMMENT ON COLUMN messages.content IS 'Message text (1-10000 characters)';
COMMENT ON COLUMN messages.original_language IS 'ISO 639-1 language code (en, ru, sr, etc.)';
COMMENT ON COLUMN messages.listing_id IS 'Optional: Listing being discussed (if different from chat context)';
COMMENT ON COLUMN messages.storefront_product_id IS 'Optional: Product being discussed (if different from chat context)';
COMMENT ON COLUMN messages.status IS 'Delivery status: sent, delivered, read, or failed';
COMMENT ON COLUMN messages.is_read IS 'Whether message has been read by receiver';
COMMENT ON COLUMN messages.has_attachments IS 'Quick check if message has attachments';
COMMENT ON COLUMN messages.attachments_count IS 'Number of attachments (for UI display)';
COMMENT ON COLUMN messages.created_at IS 'When message was created';
COMMENT ON COLUMN messages.updated_at IS 'Last update timestamp';
COMMENT ON COLUMN messages.read_at IS 'When message was read by receiver (null if unread)';

-- =====================================================
-- Trigger: Update updated_at on row change
-- =====================================================

CREATE OR REPLACE FUNCTION update_messages_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_messages_updated_at
    BEFORE UPDATE ON messages
    FOR EACH ROW
    EXECUTE FUNCTION update_messages_updated_at();

-- =====================================================
-- Trigger: Update read_at when is_read changes to true
-- =====================================================

CREATE OR REPLACE FUNCTION update_messages_read_at()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.is_read = true AND OLD.is_read = false THEN
        NEW.read_at = CURRENT_TIMESTAMP;
        NEW.status = 'read';
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_messages_read_at
    BEFORE UPDATE ON messages
    FOR EACH ROW
    WHEN (NEW.is_read IS DISTINCT FROM OLD.is_read)
    EXECUTE FUNCTION update_messages_read_at();

-- =====================================================
-- Trigger: Update parent chat's last_message_at
-- =====================================================

CREATE OR REPLACE FUNCTION update_chat_last_message_at()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE chats
    SET last_message_at = NEW.created_at,
        updated_at = CURRENT_TIMESTAMP
    WHERE id = NEW.chat_id;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_chat_last_message
    AFTER INSERT ON messages
    FOR EACH ROW
    EXECUTE FUNCTION update_chat_last_message_at();

COMMENT ON FUNCTION update_messages_updated_at IS
'Automatically update updated_at timestamp when message is modified';

COMMENT ON FUNCTION update_messages_read_at IS
'Automatically set read_at timestamp and update status when message is marked as read';

COMMENT ON FUNCTION update_chat_last_message_at IS
'Update parent chat last_message_at when new message is inserted';
