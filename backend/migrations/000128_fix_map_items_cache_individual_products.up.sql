-- Удаляем старую materialized view
DROP MATERIALIZED VIEW IF EXISTS map_items_cache;

-- Создаем обновленную materialized view с поддержкой индивидуальных локаций товаров
CREATE MATERIALIZED VIEW map_items_cache AS
WITH combined_items AS (
    -- Marketplace listings (без изменений)
    SELECT 
        ml.id,
        ml.title AS name,
        ml.description,
        ml.price,
        c.name AS category_name,
        ug.location,
        ST_Y(ug.location::geometry) AS latitude,
        ST_X(ug.location::geometry) AS longitude,
        ug.formatted_address,
        ml.user_id,
        ml.storefront_id,
        ml.status,
        ml.created_at,
        ml.updated_at,
        ml.views_count,
        0 AS rating,
        'marketplace_listing'::text AS item_type,
        COALESCE(ug.privacy_level, 'exact'::location_privacy_level) AS privacy_level,
        COALESCE(ug.blur_radius_meters, 0) AS blur_radius_meters,
        'individual'::text AS display_strategy
    FROM marketplace_listings ml
    JOIN unified_geo ug ON ug.source_type = 'marketplace_listing'::geo_source_type AND ug.source_id = ml.id
    LEFT JOIN marketplace_categories c ON ml.category_id = c.id
    WHERE ml.status = 'active' AND ml.show_on_map = true

    UNION ALL

    -- Storefronts (витрины без индивидуальных товаров на карте)
    SELECT 
        s.id,
        s.name,
        s.description,
        0 AS price,
        'Витрина'::varchar AS category_name,
        ug.location,
        ST_Y(ug.location::geometry) AS latitude,
        ST_X(ug.location::geometry) AS longitude,
        ug.formatted_address,
        s.user_id,
        NULL::integer AS storefront_id,
        CASE WHEN s.is_active THEN 'active'::text ELSE 'inactive'::text END AS status,
        s.created_at,
        s.updated_at,
        s.views_count,
        0 AS rating,
        'storefront'::text AS item_type,
        COALESCE(ug.privacy_level, 'exact'::location_privacy_level) AS privacy_level,
        COALESCE(ug.blur_radius_meters, 0) AS blur_radius_meters,
        'grouped'::text AS display_strategy
    FROM storefronts s
    JOIN unified_geo ug ON ug.source_type = 'storefront'::geo_source_type AND ug.source_id = s.id
    WHERE s.is_active = true

    UNION ALL

    -- Storefront products with individual locations (товары с индивидуальными адресами)
    SELECT 
        sp.id,
        sp.name,
        sp.description,
        COALESCE(spv.price, sp.price) AS price,
        spc.name AS category_name,
        ug.location,
        ST_Y(ug.location::geometry) AS latitude,
        ST_X(ug.location::geometry) AS longitude,
        ug.formatted_address,
        s.user_id,
        sp.storefront_id,
        CASE WHEN sp.is_active AND s.is_active THEN 'active'::text ELSE 'inactive'::text END AS status,
        sp.created_at,
        sp.updated_at,
        sp.view_count AS views_count,
        0 AS rating,
        'storefront_product'::text AS item_type,
        COALESCE(ug.privacy_level, 'exact'::location_privacy_level) AS privacy_level,
        COALESCE(ug.blur_radius_meters, 0) AS blur_radius_meters,
        'individual'::text AS display_strategy
    FROM storefront_products sp
    JOIN storefronts s ON sp.storefront_id = s.id
    JOIN unified_geo ug ON ug.source_type = 'storefront_product'::geo_source_type AND ug.source_id = sp.id
    LEFT JOIN storefront_product_variants spv ON sp.id = spv.product_id AND spv.is_default = true
    LEFT JOIN marketplace_categories spc ON sp.category_id = spc.id
    WHERE sp.is_active = true 
      AND s.is_active = true 
      AND sp.has_individual_location = true
      AND sp.show_on_map = true

    UNION ALL

    -- Storefront products without individual locations (товары без индивидуальных адресов - используют адрес витрины)
    SELECT 
        sp.id,
        sp.name,
        sp.description,
        COALESCE(spv.price, sp.price) AS price,
        spc.name AS category_name,
        ug.location,
        ST_Y(ug.location::geometry) AS latitude,
        ST_X(ug.location::geometry) AS longitude,
        ug.formatted_address,
        s.user_id,
        sp.storefront_id,
        CASE WHEN sp.is_active AND s.is_active THEN 'active'::text ELSE 'inactive'::text END AS status,
        sp.created_at,
        sp.updated_at,
        sp.view_count AS views_count,
        0 AS rating,
        'storefront_product'::text AS item_type,
        COALESCE(ug.privacy_level, 'exact'::location_privacy_level) AS privacy_level,
        COALESCE(ug.blur_radius_meters, 0) AS blur_radius_meters,
        'grouped'::text AS display_strategy
    FROM storefront_products sp
    JOIN storefronts s ON sp.storefront_id = s.id
    JOIN unified_geo ug ON ug.source_type = 'storefront'::geo_source_type AND ug.source_id = s.id
    LEFT JOIN storefront_product_variants spv ON sp.id = spv.product_id AND spv.is_default = true
    LEFT JOIN marketplace_categories spc ON sp.category_id = spc.id
    WHERE sp.is_active = true 
      AND s.is_active = true 
      AND (sp.has_individual_location = false OR sp.has_individual_location IS NULL)
      AND sp.show_on_map = true
)
SELECT * FROM combined_items;

-- Создаем индексы для materialized view
CREATE INDEX idx_map_items_cache_location ON map_items_cache USING GIST(location);
CREATE INDEX idx_map_items_cache_item_type ON map_items_cache(item_type);
CREATE INDEX idx_map_items_cache_status ON map_items_cache(status);
CREATE INDEX idx_map_items_cache_user_id ON map_items_cache(user_id);
CREATE INDEX idx_map_items_cache_storefront_id ON map_items_cache(storefront_id);
CREATE INDEX idx_map_items_cache_category ON map_items_cache(category_name);
CREATE INDEX idx_map_items_cache_price ON map_items_cache(price);

-- Создаем уникальный индекс для поддержки CONCURRENTLY refresh
CREATE UNIQUE INDEX idx_map_items_cache_unique ON map_items_cache(item_type, id);

-- Обновляем данные
REFRESH MATERIALIZED VIEW map_items_cache;

-- Добавляем комментарий
COMMENT ON MATERIALIZED VIEW map_items_cache IS 'Кэш элементов для отображения на карте с поддержкой товаров витрин с индивидуальными адресами';