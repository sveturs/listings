-- Migration down: Remove storefront geo strategy and privacy controls

-- Drop triggers
DROP TRIGGER IF EXISTS trigger_update_storefront_products_geo ON storefronts;

-- Drop functions
DROP FUNCTION IF EXISTS update_storefront_products_geo();
DROP FUNCTION IF EXISTS calculate_blurred_location(NUMERIC, NUMERIC, location_privacy_level);

-- Drop materialized view
DROP MATERIALIZED VIEW IF EXISTS map_items_cache;

-- Drop indexes
DROP INDEX IF EXISTS idx_storefront_products_individual_location;
DROP INDEX IF EXISTS idx_storefront_products_show_on_map;
DROP INDEX IF EXISTS idx_storefront_products_privacy;
DROP INDEX IF EXISTS idx_storefronts_geo_strategy;
DROP INDEX IF EXISTS idx_storefronts_coordinates;

-- Remove columns from storefront_products
ALTER TABLE storefront_products DROP COLUMN IF EXISTS has_individual_location;
ALTER TABLE storefront_products DROP COLUMN IF EXISTS individual_address;
ALTER TABLE storefront_products DROP COLUMN IF EXISTS individual_latitude;
ALTER TABLE storefront_products DROP COLUMN IF EXISTS individual_longitude;
ALTER TABLE storefront_products DROP COLUMN IF EXISTS location_privacy;
ALTER TABLE storefront_products DROP COLUMN IF EXISTS show_on_map;

-- Remove columns from storefronts
ALTER TABLE storefronts DROP COLUMN IF EXISTS geo_strategy;
ALTER TABLE storefronts DROP COLUMN IF EXISTS default_privacy_level;
ALTER TABLE storefronts DROP COLUMN IF EXISTS address;
ALTER TABLE storefronts DROP COLUMN IF EXISTS latitude;
ALTER TABLE storefronts DROP COLUMN IF EXISTS longitude;
ALTER TABLE storefronts DROP COLUMN IF EXISTS formatted_address;
ALTER TABLE storefronts DROP COLUMN IF EXISTS address_verified;

-- Remove columns from unified_geo
ALTER TABLE unified_geo DROP COLUMN IF EXISTS privacy_level;
ALTER TABLE unified_geo DROP COLUMN IF EXISTS original_location;
ALTER TABLE unified_geo DROP COLUMN IF EXISTS blur_radius_meters;

-- Drop enums
DROP TYPE IF EXISTS location_privacy_level;
DROP TYPE IF EXISTS storefront_geo_strategy;
