-- Тестовый скрипт для проверки создания записей в unified_geo

-- Проверяем последние объявления с геоданными
SELECT 
    ml.id,
    ml.title,
    ml.location,
    ml.latitude,
    ml.longitude,
    ml.address_city,
    ml.address_country,
    ml.created_at
FROM marketplace_listings ml
WHERE ml.latitude IS NOT NULL
ORDER BY ml.created_at DESC
LIMIT 5;

-- Проверяем записи в unified_geo для объявлений
SELECT 
    ug.source_type,
    ug.source_id,
    ST_Y(ug.location::geometry) as lat,
    ST_X(ug.location::geometry) as lng,
    ug.formatted_address,
    ug.address_components,
    ug.created_at
FROM unified_geo ug
WHERE ug.source_type = 'marketplace_listing'
ORDER BY ug.created_at DESC
LIMIT 5;

-- Проверяем последние витрины с геоданными
SELECT 
    s.id,
    s.name,
    s.address,
    s.latitude,
    s.longitude,
    s.city,
    s.country,
    s.created_at
FROM storefronts s
WHERE s.latitude IS NOT NULL
ORDER BY s.created_at DESC
LIMIT 5;

-- Проверяем записи в unified_geo для витрин
SELECT 
    ug.source_type,
    ug.source_id,
    ST_Y(ug.location::geometry) as lat,
    ST_X(ug.location::geometry) as lng,
    ug.formatted_address,
    ug.address_components,
    ug.created_at
FROM unified_geo ug
WHERE ug.source_type = 'storefront'
ORDER BY ug.created_at DESC
LIMIT 5;