-- Rollback migration of coordinates
-- This will remove all migrated data from listings_geo

-- First, log what will be deleted
DO $$
DECLARE
    records_to_delete INTEGER;
BEGIN
    SELECT COUNT(*) INTO records_to_delete
    FROM listings_geo lg
    INNER JOIN marketplace_listings l ON lg.listing_id = l.id
    WHERE l.latitude IS NOT NULL 
    AND l.longitude IS NOT NULL
    AND l.latitude != 0 
    AND l.longitude != 0;
    
    RAISE NOTICE 'Will delete % migrated records from listings_geo', records_to_delete;
END $$;

-- Delete only the records that were migrated from existing coordinates
DELETE FROM listings_geo
WHERE listing_id IN (
    SELECT l.id
    FROM marketplace_listings l
    WHERE l.latitude IS NOT NULL 
    AND l.longitude IS NOT NULL
    AND l.latitude != 0 
    AND l.longitude != 0
);

-- Log results
DO $$
DECLARE
    remaining_count INTEGER;
BEGIN
    SELECT COUNT(*) INTO remaining_count FROM listings_geo;
    RAISE NOTICE 'Rollback completed. % records remain in listings_geo', remaining_count;
END $$;