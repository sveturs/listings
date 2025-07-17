-- Migration down: Remove unified geo system

-- Drop triggers
DROP TRIGGER IF EXISTS trigger_unified_geo_cache_refresh ON unified_geo;
DROP TRIGGER IF EXISTS trigger_storefront_products_cache_refresh ON storefront_products;
DROP TRIGGER IF EXISTS trigger_marketplace_listings_cache_refresh ON marketplace_listings;
DROP TRIGGER IF EXISTS trigger_cleanup_storefront_product_geo ON storefront_products;
DROP TRIGGER IF EXISTS trigger_auto_geocode_storefront_product ON storefront_products;
DROP TRIGGER IF EXISTS trigger_update_unified_geo_updated_at ON unified_geo;

-- Drop functions
DROP FUNCTION IF EXISTS trigger_refresh_map_cache();
DROP FUNCTION IF EXISTS refresh_map_items_cache();
DROP FUNCTION IF EXISTS cleanup_unified_geo();
DROP FUNCTION IF EXISTS auto_geocode_storefront_product();
DROP FUNCTION IF EXISTS update_unified_geo_updated_at();

-- Drop materialized view
DROP MATERIALIZED VIEW IF EXISTS map_items_cache;

-- Drop views
DROP VIEW IF EXISTS storefront_products_geo;
DROP VIEW IF EXISTS listings_geo_view;

-- Drop indexes
DROP INDEX IF EXISTS idx_unified_geo_status;
DROP INDEX IF EXISTS idx_unified_geo_geohash;
DROP INDEX IF EXISTS idx_unified_geo_category;
DROP INDEX IF EXISTS idx_unified_geo_type;
DROP INDEX IF EXISTS idx_unified_geo_location;
DROP INDEX IF EXISTS idx_unified_geo_storefront_active;
DROP INDEX IF EXISTS idx_unified_geo_marketplace_active;
DROP INDEX IF EXISTS idx_unified_geo_location_bounds;
DROP INDEX IF EXISTS idx_unified_geo_composite;
DROP INDEX IF EXISTS idx_unified_geo_source_id;
DROP INDEX IF EXISTS idx_unified_geo_source_type;
DROP INDEX IF EXISTS idx_unified_geo_geohash;
DROP INDEX IF EXISTS idx_unified_geo_location;

-- Drop table
DROP TABLE IF EXISTS unified_geo;

-- Drop enum
DROP TYPE IF EXISTS geo_source_type;
