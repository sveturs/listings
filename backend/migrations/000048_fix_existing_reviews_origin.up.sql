-- Обновляем entity_origin для существующих отзывов на товары
UPDATE reviews r
SET 
    entity_origin_type = CASE 
        WHEN l.storefront_id IS NOT NULL AND l.storefront_id > 0 THEN 'storefront'
        ELSE 'user'
    END,
    entity_origin_id = CASE 
        WHEN l.storefront_id IS NOT NULL AND l.storefront_id > 0 THEN l.storefront_id
        ELSE l.user_id
    END
FROM marketplace_listings l
WHERE r.entity_type = 'listing' 
AND r.entity_id = l.id
AND r.entity_origin_type IS NULL;

-- Обновляем entity_origin для отзывов на пользователей
UPDATE reviews
SET entity_origin_type = 'user',
    entity_origin_id = entity_id
WHERE entity_type = 'user' 
AND entity_origin_type IS NULL;

-- Обновляем entity_origin для отзывов на витрины
UPDATE reviews
SET entity_origin_type = 'storefront',
    entity_origin_id = entity_id
WHERE entity_type = 'storefront' 
AND entity_origin_type IS NULL;

-- Обновляем материализованные представления
REFRESH MATERIALIZED VIEW CONCURRENTLY user_ratings;
REFRESH MATERIALIZED VIEW CONCURRENTLY storefront_ratings;