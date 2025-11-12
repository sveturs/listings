-- Add source_type field to listings table
-- This migration adds source_type column to distinguish between C2C (consumer-to-consumer)
-- and B2C (business-to-consumer) listings

-- Add source_type column with default value 'c2c' for existing records
ALTER TABLE listings
ADD COLUMN source_type VARCHAR(10) NOT NULL DEFAULT 'c2c';

-- Add check constraint to ensure only valid values
ALTER TABLE listings
ADD CONSTRAINT listings_source_type_check
CHECK (source_type IN ('c2c', 'b2c'));

-- Create index for efficient filtering by source_type
CREATE INDEX idx_listings_source_type ON listings(source_type) WHERE is_deleted = false;

-- Add comment for documentation
COMMENT ON COLUMN listings.source_type IS 'Type of listing: c2c (consumer-to-consumer) or b2c (business-to-consumer)';
