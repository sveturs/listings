-- Migration: Storefront geo strategy and privacy controls
-- This migration adds geo strategy and privacy controls for storefront products

-- Create enum for geo strategy
CREATE TYPE storefront_geo_strategy AS ENUM ('storefront_location', 'individual_location');

-- Create enum for location privacy levels
CREATE TYPE location_privacy_level AS ENUM ('exact', 'approximate', 'city_only', 'hidden');

-- Add geo strategy to storefronts table
ALTER TABLE storefronts ADD COLUMN IF NOT EXISTS geo_strategy storefront_geo_strategy DEFAULT 'storefront_location';
ALTER TABLE storefronts ADD COLUMN IF NOT EXISTS default_privacy_level location_privacy_level DEFAULT 'exact';

-- Add location fields to storefront_products table
ALTER TABLE storefront_products ADD COLUMN IF NOT EXISTS has_individual_location BOOLEAN DEFAULT false;
ALTER TABLE storefront_products ADD COLUMN IF NOT EXISTS individual_address TEXT;
ALTER TABLE storefront_products ADD COLUMN IF NOT EXISTS individual_latitude NUMERIC(10, 8);
ALTER TABLE storefront_products ADD COLUMN IF NOT EXISTS individual_longitude NUMERIC(11, 8);
ALTER TABLE storefront_products ADD COLUMN IF NOT EXISTS location_privacy location_privacy_level DEFAULT 'exact';
ALTER TABLE storefront_products ADD COLUMN IF NOT EXISTS show_on_map BOOLEAN DEFAULT true;

-- Add geo metadata to storefronts
ALTER TABLE storefronts ADD COLUMN IF NOT EXISTS address TEXT;
ALTER TABLE storefronts ADD COLUMN IF NOT EXISTS latitude NUMERIC(10, 8);
ALTER TABLE storefronts ADD COLUMN IF NOT EXISTS longitude NUMERIC(11, 8);
ALTER TABLE storefronts ADD COLUMN IF NOT EXISTS formatted_address TEXT;
ALTER TABLE storefronts ADD COLUMN IF NOT EXISTS address_verified BOOLEAN DEFAULT false;

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_storefront_products_individual_location ON storefront_products (has_individual_location);
CREATE INDEX IF NOT EXISTS idx_storefront_products_show_on_map ON storefront_products (show_on_map);
CREATE INDEX IF NOT EXISTS idx_storefront_products_privacy ON storefront_products (location_privacy);
CREATE INDEX IF NOT EXISTS idx_storefronts_geo_strategy ON storefronts (geo_strategy);
CREATE INDEX IF NOT EXISTS idx_storefronts_coordinates ON storefronts (latitude, longitude) WHERE latitude IS NOT NULL AND longitude IS NOT NULL;

-- Update unified_geo table to support privacy levels
ALTER TABLE unified_geo ADD COLUMN IF NOT EXISTS privacy_level location_privacy_level DEFAULT 'exact';
ALTER TABLE unified_geo ADD COLUMN IF NOT EXISTS original_location geography(Point, 4326);
ALTER TABLE unified_geo ADD COLUMN IF NOT EXISTS blur_radius_meters INTEGER DEFAULT 0;

-- Create function to calculate blurred coordinates based on privacy level
CREATE OR REPLACE FUNCTION calculate_blurred_location(
    original_lat NUMERIC,
    original_lng NUMERIC,
    privacy_level location_privacy_level
) RETURNS geography AS $$
DECLARE
    blur_radius INTEGER;
    random_angle NUMERIC;
    random_distance NUMERIC;
    blurred_lat NUMERIC;
    blurred_lng NUMERIC;
BEGIN
    -- Determine blur radius based on privacy level
    CASE privacy_level
        WHEN 'exact' THEN
            RETURN ST_SetSRID(ST_MakePoint(original_lng, original_lat), 4326)::geography;
        WHEN 'approximate' THEN
            blur_radius := 500; -- 500 meters
        WHEN 'city_only' THEN
            blur_radius := 5000; -- 5 km
        WHEN 'hidden' THEN
            RETURN NULL; -- Don't show on map at all
    END CASE;
    
    -- Generate random offset within blur radius
    random_angle := random() * 2 * pi();
    random_distance := random() * blur_radius;
    
    -- Calculate offset in degrees (approximate)
    -- 1 degree latitude ≈ 111,000 meters
    -- 1 degree longitude ≈ 111,000 * cos(latitude) meters
    blurred_lat := original_lat + (random_distance * cos(random_angle)) / 111000.0;
    blurred_lng := original_lng + (random_distance * sin(random_angle)) / (111000.0 * cos(radians(original_lat)));
    
    RETURN ST_SetSRID(ST_MakePoint(blurred_lng, blurred_lat), 4326)::geography;
END;
$$ LANGUAGE plpgsql;

-- Update auto_geocode_storefront_product function to handle new logic
CREATE OR REPLACE FUNCTION auto_geocode_storefront_product()
RETURNS TRIGGER AS $$
DECLARE
    storefront_rec RECORD;
    product_location geography(Point, 4326);
    original_location geography(Point, 4326);
    calculated_geohash VARCHAR(12);
    privacy_level location_privacy_level;
BEGIN
    -- Get storefront data
    SELECT s.geo_strategy, s.address, s.latitude, s.longitude, s.default_privacy_level
    INTO storefront_rec
    FROM storefronts s
    WHERE s.id = NEW.storefront_id;
    
    -- Determine location and privacy level
    IF NEW.has_individual_location AND NEW.individual_latitude IS NOT NULL AND NEW.individual_longitude IS NOT NULL THEN
        -- Product has individual location
        original_location := ST_SetSRID(ST_MakePoint(NEW.individual_longitude, NEW.individual_latitude), 4326)::geography;
        privacy_level := COALESCE(NEW.location_privacy, storefront_rec.default_privacy_level);
    ELSIF storefront_rec.latitude IS NOT NULL AND storefront_rec.longitude IS NOT NULL THEN
        -- Use storefront location
        original_location := ST_SetSRID(ST_MakePoint(storefront_rec.longitude, storefront_rec.latitude), 4326)::geography;
        privacy_level := 'exact'; -- Storefront location is always exact
    ELSE
        -- No location available
        RETURN NEW;
    END IF;
    
    -- Calculate display location based on privacy level
    product_location := calculate_blurred_location(
        ST_Y(original_location::geometry),
        ST_X(original_location::geometry),
        privacy_level
    );
    
    -- Skip if location should be hidden
    IF product_location IS NULL OR NOT NEW.show_on_map THEN
        -- Remove from geo table if exists
        DELETE FROM unified_geo 
        WHERE source_type = 'storefront_product' AND source_id = NEW.id;
        RETURN NEW;
    END IF;
    
    -- Calculate geohash
    calculated_geohash := substring(md5(ST_Y(product_location::geometry)::text || ST_X(product_location::geometry)::text), 1, 12);
    
    -- Insert or update geo data
    INSERT INTO unified_geo (
        source_type, source_id, location, original_location, geohash, 
        privacy_level, blur_radius_meters,
        formatted_address, created_at, updated_at
    ) VALUES (
        'storefront_product',
        NEW.id,
        product_location,
        original_location,
        calculated_geohash,
        privacy_level,
        CASE privacy_level
            WHEN 'approximate' THEN 500
            WHEN 'city_only' THEN 5000
            ELSE 0
        END,
        COALESCE(NEW.individual_address, storefront_rec.address),
        NOW(),
        NOW()
    ) ON CONFLICT (source_type, source_id) DO UPDATE SET
        location = EXCLUDED.location,
        original_location = EXCLUDED.original_location,
        geohash = EXCLUDED.geohash,
        privacy_level = EXCLUDED.privacy_level,
        blur_radius_meters = EXCLUDED.blur_radius_meters,
        formatted_address = EXCLUDED.formatted_address,
        updated_at = EXCLUDED.updated_at;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create function to handle storefront location updates
CREATE OR REPLACE FUNCTION update_storefront_products_geo()
RETURNS TRIGGER AS $$
BEGIN
    -- If storefront coordinates changed, update all products that use storefront location
    IF (OLD.latitude IS DISTINCT FROM NEW.latitude OR OLD.longitude IS DISTINCT FROM NEW.longitude) THEN
        -- Trigger re-geocoding for all products that don't have individual locations
        UPDATE storefront_products 
        SET updated_at = NOW()
        WHERE storefront_id = NEW.id AND has_individual_location = false;
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger for storefront location updates
CREATE TRIGGER trigger_update_storefront_products_geo
    AFTER UPDATE ON storefronts
    FOR EACH ROW
    EXECUTE FUNCTION update_storefront_products_geo();

-- Update materialized view to handle new geo strategy
DROP MATERIALIZED VIEW IF EXISTS map_items_cache;

CREATE MATERIALIZED VIEW map_items_cache AS
SELECT 
    'marketplace_listing' as item_type,
    ml.id,
    ml.title as name,
    ml.description,
    ml.price,
    ml.category_id,
    mc.name as category_name,
    ml.user_id,
    NULL::INTEGER as storefront_id,
    ml.status,
    ml.views_count,
    COALESCE(rc.average_rating, 0) as rating,
    ST_Y(ug.location::geometry) as latitude,
    ST_X(ug.location::geometry) as longitude,
    ug.geohash,
    ug.privacy_level,
    ug.blur_radius_meters,
    ml.created_at,
    ml.updated_at,
    COALESCE(mi.images, '[]'::jsonb) as images,
    'individual' as display_strategy
FROM marketplace_listings ml
JOIN unified_geo ug ON ug.source_type = 'marketplace_listing' AND ug.source_id = ml.id
LEFT JOIN marketplace_categories mc ON ml.category_id = mc.id
LEFT JOIN rating_cache rc ON rc.entity_type = 'listing' AND rc.entity_id = ml.id
LEFT JOIN (
    SELECT listing_id, jsonb_agg(
        jsonb_build_object(
            'id', id,
            'image_url', public_url,
            'is_main', is_main
        ) ORDER BY is_main DESC, created_at
    ) as images
    FROM marketplace_images
    GROUP BY listing_id
) mi ON mi.listing_id = ml.id
WHERE ml.status = 'active'

UNION ALL

-- Storefront products with individual locations
SELECT 
    'storefront_product' as item_type,
    sp.id,
    sp.name,
    sp.description,
    sp.price,
    sp.category_id,
    mc.name as category_name,
    s.user_id,
    sp.storefront_id,
    CASE WHEN sp.is_active THEN 'active' ELSE 'inactive' END as status,
    sp.view_count as views_count,
    0 as rating,
    ST_Y(ug.location::geometry) as latitude,
    ST_X(ug.location::geometry) as longitude,
    ug.geohash,
    ug.privacy_level,
    ug.blur_radius_meters,
    sp.created_at,
    sp.updated_at,
    COALESCE(spi.images, '[]'::jsonb) as images,
    CASE 
        WHEN sp.has_individual_location THEN 'individual'
        ELSE 'storefront_grouped'
    END as display_strategy
FROM storefront_products sp
JOIN unified_geo ug ON ug.source_type = 'storefront_product' AND ug.source_id = sp.id
JOIN storefronts s ON sp.storefront_id = s.id
LEFT JOIN marketplace_categories mc ON sp.category_id = mc.id
LEFT JOIN (
    SELECT 
        spv.product_id,
        jsonb_agg(
            jsonb_build_object(
                'id', spvi.id,
                'image_url', spvi.image_url,
                'is_main', spvi.is_main
            ) ORDER BY spvi.is_main DESC, spvi.display_order
        ) as images
    FROM storefront_product_variants spv
    JOIN storefront_product_variant_images spvi ON spv.id = spvi.variant_id
    WHERE spv.is_active = true
    GROUP BY spv.product_id
) spi ON spi.product_id = sp.id
WHERE sp.is_active = true AND sp.show_on_map = true

UNION ALL

-- Storefronts as grouped markers (for products without individual locations)
SELECT 
    'storefront' as item_type,
    s.id,
    s.name,
    s.description,
    NULL as price,
    NULL as category_id,
    NULL as category_name,
    s.user_id,
    s.id as storefront_id,
    CASE WHEN s.is_active THEN 'active' ELSE 'inactive' END as status,
    0 as views_count,
    0 as rating,
    s.latitude,
    s.longitude,
    substring(md5(s.latitude::text || s.longitude::text), 1, 12) as geohash,
    'exact'::location_privacy_level as privacy_level,
    0 as blur_radius_meters,
    s.created_at,
    s.updated_at,
    '[]'::jsonb as images,
    'storefront_grouped' as display_strategy
FROM storefronts s
WHERE s.is_active = true 
  AND s.latitude IS NOT NULL 
  AND s.longitude IS NOT NULL
  AND s.geo_strategy = 'storefront_location'
  AND EXISTS (
      SELECT 1 FROM storefront_products sp 
      WHERE sp.storefront_id = s.id 
        AND sp.is_active = true 
        AND sp.show_on_map = true
        AND sp.has_individual_location = false
  );

-- Recreate indexes on materialized view
CREATE INDEX IF NOT EXISTS idx_map_items_cache_location ON map_items_cache (latitude, longitude);
CREATE INDEX IF NOT EXISTS idx_map_items_cache_type ON map_items_cache (item_type);
CREATE INDEX IF NOT EXISTS idx_map_items_cache_category ON map_items_cache (category_id);
CREATE INDEX IF NOT EXISTS idx_map_items_cache_geohash ON map_items_cache (geohash);
CREATE INDEX IF NOT EXISTS idx_map_items_cache_status ON map_items_cache (status);
CREATE INDEX IF NOT EXISTS idx_map_items_cache_strategy ON map_items_cache (display_strategy);
CREATE INDEX IF NOT EXISTS idx_map_items_cache_storefront ON map_items_cache (storefront_id) WHERE storefront_id IS NOT NULL;
