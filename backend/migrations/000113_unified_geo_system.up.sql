-- Migration: Unified geo system for marketplace_listings and storefront_products
-- This migration creates a unified geographic data system that supports both listing types

-- Create enum for source types
CREATE TYPE geo_source_type AS ENUM ('marketplace_listing', 'storefront_product');

-- Create unified geo table
CREATE TABLE IF NOT EXISTS unified_geo (
    id BIGSERIAL PRIMARY KEY,
    source_type geo_source_type NOT NULL,
    source_id BIGINT NOT NULL,
    location geography(Point, 4326) NOT NULL,
    geohash VARCHAR(12) NOT NULL,
    is_precise BOOLEAN NOT NULL DEFAULT true,
    blur_radius NUMERIC(10, 2) DEFAULT 0,
    address_components JSONB,
    formatted_address TEXT,
    geocoding_confidence NUMERIC(3, 2),
    address_verified BOOLEAN DEFAULT false,
    input_method VARCHAR(50) DEFAULT 'manual',
    location_privacy VARCHAR(20) DEFAULT 'exact',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    -- Ensure one geo record per source
    CONSTRAINT uk_unified_geo_source UNIQUE (source_type, source_id)
);

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_unified_geo_location ON unified_geo USING GIST (location);
CREATE INDEX IF NOT EXISTS idx_unified_geo_geohash ON unified_geo (geohash);
CREATE INDEX IF NOT EXISTS idx_unified_geo_source_type ON unified_geo (source_type);
CREATE INDEX IF NOT EXISTS idx_unified_geo_source_id ON unified_geo (source_id);
CREATE INDEX IF NOT EXISTS idx_unified_geo_composite ON unified_geo (source_type, source_id);

-- Create function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_unified_geo_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger for updated_at
CREATE TRIGGER trigger_update_unified_geo_updated_at
    BEFORE UPDATE ON unified_geo
    FOR EACH ROW
    EXECUTE FUNCTION update_unified_geo_updated_at();

-- Migrate existing data from listings_geo
INSERT INTO unified_geo (
    source_type, source_id, location, geohash, is_precise, blur_radius, created_at, updated_at
)
SELECT 
    'marketplace_listing'::geo_source_type,
    listing_id,
    location,
    geohash,
    is_precise,
    blur_radius,
    created_at,
    updated_at
FROM listings_geo
ON CONFLICT (source_type, source_id) DO NOTHING;

-- Create view for backward compatibility with existing listings_geo queries
CREATE OR REPLACE VIEW listings_geo_view AS
SELECT 
    ug.id,
    ug.source_id as listing_id,
    ug.location,
    ug.geohash,
    ug.is_precise,
    ug.blur_radius,
    ug.created_at,
    ug.updated_at
FROM unified_geo ug
WHERE ug.source_type = 'marketplace_listing';

-- Create view for storefront products geo data
CREATE OR REPLACE VIEW storefront_products_geo AS
SELECT 
    ug.id,
    ug.source_id as product_id,
    ug.location,
    ug.geohash,
    ug.is_precise,
    ug.blur_radius,
    ug.formatted_address,
    ug.created_at,
    ug.updated_at
FROM unified_geo ug
WHERE ug.source_type = 'storefront_product';

-- Create function to automatically add geo data for new storefront products
CREATE OR REPLACE FUNCTION auto_geocode_storefront_product()
RETURNS TRIGGER AS $$
DECLARE
    storefront_location TEXT;
    storefront_lat NUMERIC;
    storefront_lng NUMERIC;
    calculated_geohash VARCHAR(12);
BEGIN
    -- Get storefront location data
    SELECT s.address, s.latitude, s.longitude
    INTO storefront_location, storefront_lat, storefront_lng
    FROM storefronts s
    WHERE s.id = NEW.storefront_id;
    
    -- If storefront has coordinates, use them for the product
    IF storefront_lat IS NOT NULL AND storefront_lng IS NOT NULL THEN
        -- Calculate geohash (simplified version)
        calculated_geohash := substring(md5(storefront_lat::text || storefront_lng::text), 1, 12);
        
        -- Insert geo data for the new product
        INSERT INTO unified_geo (
            source_type, source_id, location, geohash, 
            formatted_address, created_at, updated_at
        ) VALUES (
            'storefront_product',
            NEW.id,
            ST_SetSRID(ST_MakePoint(storefront_lng, storefront_lat), 4326)::geography,
            calculated_geohash,
            storefront_location,
            NOW(),
            NOW()
        ) ON CONFLICT (source_type, source_id) DO UPDATE SET
            location = EXCLUDED.location,
            geohash = EXCLUDED.geohash,
            formatted_address = EXCLUDED.formatted_address,
            updated_at = EXCLUDED.updated_at;
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger for auto-geocoding storefront products
CREATE TRIGGER trigger_auto_geocode_storefront_product
    AFTER INSERT OR UPDATE ON storefront_products
    FOR EACH ROW
    EXECUTE FUNCTION auto_geocode_storefront_product();

-- Create function to clean up geo data when products are deleted
CREATE OR REPLACE FUNCTION cleanup_unified_geo()
RETURNS TRIGGER AS $$
BEGIN
    DELETE FROM unified_geo 
    WHERE source_type = 'storefront_product' AND source_id = OLD.id;
    RETURN OLD;
END;
$$ LANGUAGE plpgsql;

-- Create trigger for cleanup
CREATE TRIGGER trigger_cleanup_storefront_product_geo
    AFTER DELETE ON storefront_products
    FOR EACH ROW
    EXECUTE FUNCTION cleanup_unified_geo();

-- Add indexes for common query patterns
CREATE INDEX IF NOT EXISTS idx_unified_geo_location_bounds ON unified_geo 
USING GIST (location) 
WHERE location IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_unified_geo_marketplace_active ON unified_geo (source_id)
WHERE source_type = 'marketplace_listing';

CREATE INDEX IF NOT EXISTS idx_unified_geo_storefront_active ON unified_geo (source_id)
WHERE source_type = 'storefront_product';

-- Create materialized view for fast map queries
CREATE MATERIALIZED VIEW IF NOT EXISTS map_items_cache AS
SELECT 
    'marketplace_listing' as item_type,
    ml.id,
    ml.title as name,
    ml.description,
    ml.price,
    ml.category_id,
    mc.name as category_name,
    ml.user_id,
    ml.status,
    ml.views_count,
    COALESCE(rc.average_rating, 0) as rating,
    ST_Y(ug.location::geometry) as latitude,
    ST_X(ug.location::geometry) as longitude,
    ug.geohash,
    ml.created_at,
    ml.updated_at,
    COALESCE(mi.images, '[]'::jsonb) as images
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

SELECT 
    'storefront_product' as item_type,
    sp.id,
    sp.name,
    sp.description,
    sp.price,
    sp.category_id,
    mc.name as category_name,
    s.user_id,
    CASE WHEN sp.is_active THEN 'active' ELSE 'inactive' END as status,
    sp.view_count as views_count,
    0 as rating, -- TODO: Add rating system for storefront products
    ST_Y(ug.location::geometry) as latitude,
    ST_X(ug.location::geometry) as longitude,
    ug.geohash,
    sp.created_at,
    sp.updated_at,
    COALESCE(spi.images, '[]'::jsonb) as images
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
WHERE sp.is_active = true;

-- Create indexes on materialized view
CREATE INDEX IF NOT EXISTS idx_map_items_cache_location ON map_items_cache (latitude, longitude);
CREATE INDEX IF NOT EXISTS idx_map_items_cache_type ON map_items_cache (item_type);
CREATE INDEX IF NOT EXISTS idx_map_items_cache_category ON map_items_cache (category_id);
CREATE INDEX IF NOT EXISTS idx_map_items_cache_geohash ON map_items_cache (geohash);
CREATE INDEX IF NOT EXISTS idx_map_items_cache_status ON map_items_cache (status);

-- Create function to refresh materialized view
CREATE OR REPLACE FUNCTION refresh_map_items_cache()
RETURNS void AS $$
BEGIN
    REFRESH MATERIALIZED VIEW CONCURRENTLY map_items_cache;
END;
$$ LANGUAGE plpgsql;

-- Create function to refresh cache on data changes
CREATE OR REPLACE FUNCTION trigger_refresh_map_cache()
RETURNS TRIGGER AS $$
BEGIN
    -- Use pg_notify to trigger async refresh
    PERFORM pg_notify('refresh_map_cache', '');
    RETURN COALESCE(NEW, OLD);
END;
$$ LANGUAGE plpgsql;

-- Create triggers to refresh cache when data changes
CREATE TRIGGER trigger_marketplace_listings_cache_refresh
    AFTER INSERT OR UPDATE OR DELETE ON marketplace_listings
    FOR EACH ROW
    EXECUTE FUNCTION trigger_refresh_map_cache();

CREATE TRIGGER trigger_storefront_products_cache_refresh
    AFTER INSERT OR UPDATE OR DELETE ON storefront_products
    FOR EACH ROW
    EXECUTE FUNCTION trigger_refresh_map_cache();

CREATE TRIGGER trigger_unified_geo_cache_refresh
    AFTER INSERT OR UPDATE OR DELETE ON unified_geo
    FOR EACH ROW
    EXECUTE FUNCTION trigger_refresh_map_cache();
