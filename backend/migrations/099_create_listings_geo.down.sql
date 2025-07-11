-- Drop trigger and function
DROP TRIGGER IF EXISTS trigger_update_listings_geo_updated_at ON listings_geo;
DROP FUNCTION IF EXISTS update_listings_geo_updated_at();

-- Drop indexes
DROP INDEX IF EXISTS idx_listings_geo_location;
DROP INDEX IF EXISTS idx_listings_geo_geohash;
DROP INDEX IF EXISTS idx_listings_geo_is_precise;
DROP INDEX IF EXISTS idx_listings_geo_geohash_precise;

-- Drop table
DROP TABLE IF EXISTS listings_geo;