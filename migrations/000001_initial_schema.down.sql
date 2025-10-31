-- Rollback initial schema migration

-- Drop triggers first
DROP TRIGGER IF EXISTS update_listings_updated_at ON listings;
DROP TRIGGER IF EXISTS update_listing_images_updated_at ON listing_images;
DROP TRIGGER IF EXISTS update_listing_locations_updated_at ON listing_locations;
DROP TRIGGER IF EXISTS update_listing_stats_updated_at ON listing_stats;
DROP TRIGGER IF EXISTS update_indexing_queue_updated_at ON indexing_queue;

-- Drop function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop tables (in reverse order of creation)
DROP TABLE IF EXISTS indexing_queue;
DROP TABLE IF EXISTS listing_stats;
DROP TABLE IF EXISTS listing_locations;
DROP TABLE IF EXISTS listing_tags;
DROP TABLE IF EXISTS listing_images;
DROP TABLE IF EXISTS listing_attributes;
DROP TABLE IF EXISTS listings;

-- Drop extension (only if not used by other tables)
-- DROP EXTENSION IF EXISTS "uuid-ossp";
