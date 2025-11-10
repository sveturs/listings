-- Add deleted_at column for soft delete support
ALTER TABLE storefronts
ADD COLUMN deleted_at TIMESTAMP WITH TIME ZONE DEFAULT NULL;

-- Add index for efficient soft delete queries
CREATE INDEX idx_storefronts_deleted_at
ON storefronts(deleted_at)
WHERE deleted_at IS NULL;

-- Comment
COMMENT ON COLUMN storefronts.deleted_at IS 'Timestamp when storefront was soft deleted (NULL = active)';
