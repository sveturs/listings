-- Исправление функции auto_geocode_storefront_product для использования правильной сигнатуры calculate_blurred_location

CREATE OR REPLACE FUNCTION public.auto_geocode_storefront_product()
 RETURNS trigger
 LANGUAGE plpgsql
AS $function$
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
    -- Используем правильную сигнатуру функции с типом numeric
    product_location := calculate_blurred_location(
        ST_Y(original_location::geometry)::numeric,
        ST_X(original_location::geometry)::numeric,
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
$function$;