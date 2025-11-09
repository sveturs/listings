-- Add attributes JSONB column to listings table
-- This provides backward compatibility with fixtures that use attributes column

ALTER TABLE listings
ADD COLUMN IF NOT EXISTS attributes JSONB DEFAULT '{}'::jsonb;

-- Create GIN index for efficient JSONB queries
CREATE INDEX IF NOT EXISTS idx_listings_attributes ON listings USING GIN (attributes);

-- Add comment
COMMENT ON COLUMN listings.attributes IS 'JSONB attributes for flexible product metadata. Note: listing_attributes table also exists for key-value storage.';
