-- Migrate existing coordinates from marketplace_listings to listings_geo
-- This migration handles existing latitude/longitude data

-- First, let's check if there are any listings with coordinates
DO $$
DECLARE
    listings_with_coords INTEGER;
BEGIN
    -- Count listings that have both latitude and longitude
    SELECT COUNT(*) INTO listings_with_coords
    FROM marketplace_listings
    WHERE latitude IS NOT NULL 
    AND longitude IS NOT NULL
    AND latitude != 0 
    AND longitude != 0;
    
    IF listings_with_coords > 0 THEN
        RAISE NOTICE 'Found % listings with coordinates to migrate', listings_with_coords;
    ELSE
        RAISE NOTICE 'No listings with coordinates found to migrate';
    END IF;
END $$;

-- Insert geographic data for listings that have coordinates
INSERT INTO listings_geo (
    listing_id,
    location,
    geohash,
    is_precise,
    blur_radius,
    created_at,
    updated_at
)
SELECT 
    l.id AS listing_id,
    ST_SetSRID(ST_MakePoint(l.longitude, l.latitude), 4326) AS location,
    -- Generate geohash with precision 8 (approximately 19m x 19m at equator)
    ST_GeoHash(ST_SetSRID(ST_MakePoint(l.longitude, l.latitude), 4326), 8) AS geohash,
    true AS is_precise, -- Assume existing coordinates are precise
    0 AS blur_radius,   -- No obfuscation for existing data
    l.created_at,
    l.updated_at
FROM marketplace_listings l
WHERE l.latitude IS NOT NULL 
    AND l.longitude IS NOT NULL
    AND l.latitude != 0 
    AND l.longitude != 0
    -- Only insert if not already exists (in case migration is run multiple times)
    AND NOT EXISTS (
        SELECT 1 FROM listings_geo lg WHERE lg.listing_id = l.id
    );

-- Log migration results
DO $$
DECLARE
    migrated_count INTEGER;
    total_listings INTEGER;
BEGIN
    SELECT COUNT(*) INTO migrated_count FROM listings_geo;
    SELECT COUNT(*) INTO total_listings FROM marketplace_listings;
    
    RAISE NOTICE 'Migration completed: % out of % listings now have geographic data', 
        migrated_count, total_listings;
END $$;

-- Optionally, add a check constraint to ensure data integrity
-- This ensures that if a listing has coordinates in the old columns,
-- it must have a corresponding record in listings_geo
-- (Commented out to avoid breaking existing code that might still use old columns)
-- ALTER TABLE marketplace_listings
-- ADD CONSTRAINT check_geo_migration
-- CHECK (
--     (latitude IS NULL AND longitude IS NULL) OR
--     (latitude = 0 AND longitude = 0) OR
--     EXISTS (SELECT 1 FROM listings_geo WHERE listing_id = marketplace_listings.id)
-- );