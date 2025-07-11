-- Create table for geographic data of listings
CREATE TABLE IF NOT EXISTS listings_geo (
    id BIGSERIAL PRIMARY KEY,
    listing_id BIGINT NOT NULL,
    location geography(Point, 4326) NOT NULL,
    geohash VARCHAR(12) NOT NULL,
    is_precise BOOLEAN NOT NULL DEFAULT true,
    blur_radius NUMERIC(10, 2) DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    -- Foreign key constraint
    CONSTRAINT fk_listings_geo_listing_id 
        FOREIGN KEY (listing_id) 
        REFERENCES marketplace_listings(id) 
        ON DELETE CASCADE,
    
    -- Ensure one geo record per listing
    CONSTRAINT uk_listings_geo_listing_id UNIQUE (listing_id)
);

-- Create spatial index for geographic queries
CREATE INDEX idx_listings_geo_location ON listings_geo USING GIST (location);

-- Create index for geohash-based queries (for proximity searches)
CREATE INDEX idx_listings_geo_geohash ON listings_geo (geohash);

-- Create index for filtering by precision
CREATE INDEX idx_listings_geo_is_precise ON listings_geo (is_precise);

-- Create composite index for common query patterns
CREATE INDEX idx_listings_geo_geohash_precise ON listings_geo (geohash, is_precise);

-- Add trigger to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_listings_geo_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_listings_geo_updated_at
    BEFORE UPDATE ON listings_geo
    FOR EACH ROW
    EXECUTE FUNCTION update_listings_geo_updated_at();

-- Add comment to table
COMMENT ON TABLE listings_geo IS 'Geographic location data for marketplace listings';
COMMENT ON COLUMN listings_geo.location IS 'Point geometry in WGS84 (EPSG:4326) coordinate system';
COMMENT ON COLUMN listings_geo.geohash IS 'Geohash string for efficient proximity searches';
COMMENT ON COLUMN listings_geo.is_precise IS 'Whether the location is precise or has been obfuscated for privacy';
COMMENT ON COLUMN listings_geo.blur_radius IS 'Radius in meters for location obfuscation (0 for precise locations)';