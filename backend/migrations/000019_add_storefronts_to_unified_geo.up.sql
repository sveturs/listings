-- Добавляем активные витрины в unified_geo для отображения на карте
INSERT INTO unified_geo (source_type, source_id, location, geohash, formatted_address, privacy_level)
SELECT
    'storefront'::geo_source_type,
    id,
    ST_MakePoint(longitude, latitude)::geography,
    substring(ST_GeoHash(ST_MakePoint(longitude, latitude)::geography) from 1 for 8),
    COALESCE(address, ''),
    'exact'::location_privacy_level
FROM storefronts
WHERE is_active = true
  AND latitude IS NOT NULL
  AND longitude IS NOT NULL
  AND latitude != 0
  AND longitude != 0
ON CONFLICT (source_type, source_id) DO UPDATE
SET
    location = EXCLUDED.location,
    geohash = EXCLUDED.geohash,
    formatted_address = EXCLUDED.formatted_address,
    updated_at = NOW();

-- Также добавляем товары витрин
INSERT INTO unified_geo (source_type, source_id, location, geohash, formatted_address, privacy_level)
SELECT
    'storefront_product'::geo_source_type,
    sp.id,
    ST_MakePoint(s.longitude, s.latitude)::geography,
    substring(ST_GeoHash(ST_MakePoint(s.longitude, s.latitude)::geography) from 1 for 8),
    COALESCE(s.address, ''),
    'exact'::location_privacy_level
FROM storefront_products sp
JOIN storefronts s ON sp.storefront_id = s.id
WHERE sp.is_active = true
  AND s.is_active = true
  AND s.latitude IS NOT NULL
  AND s.longitude IS NOT NULL
  AND s.latitude != 0
  AND s.longitude != 0
ON CONFLICT (source_type, source_id) DO UPDATE
SET
    location = EXCLUDED.location,
    geohash = EXCLUDED.geohash,
    formatted_address = EXCLUDED.formatted_address,
    updated_at = NOW();