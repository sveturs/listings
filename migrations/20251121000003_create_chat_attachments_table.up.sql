-- =====================================================
-- Migration: 20251121000003_create_chat_attachments_table.up.sql
-- Description: Create chat_attachments table for Chat microservice
-- Author: Chat Microservice Phase 1 - Database Setup
-- Date: 2025-11-21
-- =====================================================

-- =====================================================
-- Table: chat_attachments
-- Description: Stores file attachments for messages
-- =====================================================
CREATE TABLE IF NOT EXISTS chat_attachments (
    -- Identification
    id BIGSERIAL PRIMARY KEY,
    message_id BIGINT NOT NULL,

    -- File metadata
    file_type VARCHAR(20) NOT NULL,
    file_name VARCHAR(255) NOT NULL,
    file_size BIGINT NOT NULL,
    content_type VARCHAR(100) NOT NULL,

    -- Storage information
    storage_type VARCHAR(20) NOT NULL DEFAULT 'minio',
    storage_bucket VARCHAR(100) NOT NULL DEFAULT 'chat-files',
    file_path VARCHAR(500) NOT NULL,

    -- URLs (generated for quick access)
    public_url TEXT,
    thumbnail_url TEXT,

    -- Additional metadata (dimensions, duration, etc.)
    metadata JSONB,

    -- Timestamps
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    -- =====================================================
    -- Foreign Keys
    -- =====================================================

    CONSTRAINT fk_attachments_message FOREIGN KEY (message_id)
        REFERENCES messages(id)
        ON DELETE CASCADE,

    -- =====================================================
    -- Constraints
    -- =====================================================

    -- File type validation
    CONSTRAINT check_file_type CHECK (
        file_type IN ('image', 'video', 'document')
    ),

    -- Storage type validation
    CONSTRAINT check_storage_type CHECK (
        storage_type IN ('minio', 's3', 'local')
    ),

    -- File size limits per type
    CONSTRAINT check_file_size CHECK (
        (file_type = 'image' AND file_size > 0 AND file_size <= 10485760) OR     -- 10MB for images
        (file_type = 'video' AND file_size > 0 AND file_size <= 52428800) OR     -- 50MB for videos
        (file_type = 'document' AND file_size > 0 AND file_size <= 20971520)     -- 20MB for documents
    ),

    -- File name must not be empty
    CONSTRAINT check_file_name CHECK (
        length(trim(file_name)) > 0 AND length(file_name) <= 255
    ),

    -- File path must not be empty
    CONSTRAINT check_file_path CHECK (
        length(trim(file_path)) > 0 AND length(file_path) <= 500
    ),

    -- Content type must follow format (e.g., image/jpeg, video/mp4)
    CONSTRAINT check_content_type CHECK (
        content_type ~ '^[a-z]+/[a-z0-9\+\-\.]+$'
    )
);

-- =====================================================
-- Indexes for Performance
-- =====================================================

-- Get attachments for a message (most common query)
CREATE INDEX idx_attachments_message_id ON chat_attachments(message_id);

-- Filter by file type (for media galleries)
CREATE INDEX idx_attachments_file_type ON chat_attachments(file_type);

-- Time-based queries (for cleanup, analytics)
CREATE INDEX idx_attachments_created_at ON chat_attachments(created_at);

-- Storage management queries
CREATE INDEX idx_attachments_storage ON chat_attachments(storage_type, storage_bucket);

-- Find large files (for storage optimization)
CREATE INDEX idx_attachments_file_size ON chat_attachments(file_size DESC)
    WHERE file_size > 1048576; -- Index only files > 1MB

-- GIN index for metadata JSONB queries
CREATE INDEX idx_attachments_metadata ON chat_attachments USING GIN (metadata);

-- =====================================================
-- Comments for Documentation
-- =====================================================

COMMENT ON TABLE chat_attachments IS
'File attachments for chat messages (images, videos, documents)';

COMMENT ON COLUMN chat_attachments.id IS 'Primary key';
COMMENT ON COLUMN chat_attachments.message_id IS 'Reference to parent message';
COMMENT ON COLUMN chat_attachments.file_type IS 'Type of file: image, video, or document';
COMMENT ON COLUMN chat_attachments.file_name IS 'Original filename (max 255 chars)';
COMMENT ON COLUMN chat_attachments.file_size IS 'File size in bytes (limits: 10MB images, 50MB videos, 20MB docs)';
COMMENT ON COLUMN chat_attachments.content_type IS 'MIME type (e.g., image/jpeg, application/pdf)';
COMMENT ON COLUMN chat_attachments.storage_type IS 'Storage backend: minio, s3, or local';
COMMENT ON COLUMN chat_attachments.storage_bucket IS 'Storage bucket/container name';
COMMENT ON COLUMN chat_attachments.file_path IS 'Path to file in storage (max 500 chars)';
COMMENT ON COLUMN chat_attachments.public_url IS 'Public URL to access file (if available)';
COMMENT ON COLUMN chat_attachments.thumbnail_url IS 'Thumbnail URL for images/videos (if generated)';
COMMENT ON COLUMN chat_attachments.metadata IS 'Additional metadata in JSON format (dimensions, duration, etc.)';
COMMENT ON COLUMN chat_attachments.created_at IS 'When attachment was uploaded';

-- =====================================================
-- Trigger: Update parent message attachments count
-- =====================================================

CREATE OR REPLACE FUNCTION update_message_attachments_count()
RETURNS TRIGGER AS $$
DECLARE
    attachment_count INT;
BEGIN
    -- Count attachments for the message
    SELECT COUNT(*) INTO attachment_count
    FROM chat_attachments
    WHERE message_id = COALESCE(NEW.message_id, OLD.message_id);

    -- Update message
    UPDATE messages
    SET
        attachments_count = attachment_count,
        has_attachments = (attachment_count > 0),
        updated_at = CURRENT_TIMESTAMP
    WHERE id = COALESCE(NEW.message_id, OLD.message_id);

    RETURN COALESCE(NEW, OLD);
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_message_attachments_count_insert
    AFTER INSERT ON chat_attachments
    FOR EACH ROW
    EXECUTE FUNCTION update_message_attachments_count();

CREATE TRIGGER trigger_update_message_attachments_count_delete
    AFTER DELETE ON chat_attachments
    FOR EACH ROW
    EXECUTE FUNCTION update_message_attachments_count();

COMMENT ON FUNCTION update_message_attachments_count IS
'Update parent message attachments_count and has_attachments when attachments are added/removed';

-- =====================================================
-- Helper Functions
-- =====================================================

-- Function to get file type from content_type
CREATE OR REPLACE FUNCTION get_file_type_from_content_type(content_type TEXT)
RETURNS VARCHAR(20) AS $$
BEGIN
    CASE
        WHEN content_type LIKE 'image/%' THEN RETURN 'image';
        WHEN content_type LIKE 'video/%' THEN RETURN 'video';
        ELSE RETURN 'document';
    END CASE;
END;
$$ LANGUAGE plpgsql IMMUTABLE;

COMMENT ON FUNCTION get_file_type_from_content_type IS
'Determine file_type from MIME content_type';

-- Function to validate file size against type limits
CREATE OR REPLACE FUNCTION is_valid_file_size(file_type TEXT, file_size BIGINT)
RETURNS BOOLEAN AS $$
BEGIN
    RETURN (
        (file_type = 'image' AND file_size > 0 AND file_size <= 10485760) OR
        (file_type = 'video' AND file_size > 0 AND file_size <= 52428800) OR
        (file_type = 'document' AND file_size > 0 AND file_size <= 20971520)
    );
END;
$$ LANGUAGE plpgsql IMMUTABLE;

COMMENT ON FUNCTION is_valid_file_size IS
'Check if file size is within allowed limits for the file type';
