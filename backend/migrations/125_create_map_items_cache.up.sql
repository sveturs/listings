-- Создаем materialized view для объединения всех geo-элементов
CREATE MATERIALIZED VIEW IF NOT EXISTS map_items_cache AS
WITH combined_items AS (
    -- Объявления marketplace
    SELECT 
        ml.id,
        ml.title as name,
        ml.description,
        ml.price,
        c.name as category_name,
        ug.location,
        ST_Y(ug.location::geometry) as latitude,
        ST_X(ug.location::geometry) as longitude,
        ug.formatted_address,
        ml.user_id,
        ml.storefront_id,
        ml.status,
        ml.created_at,
        ml.updated_at,
        ml.views_count,
        0 as rating,
        'marketplace_listing' as item_type,
        COALESCE(ug.privacy_level, 'exact') as privacy_level,
        COALESCE(ug.blur_radius_meters, 0) as blur_radius_meters,
        'individual' as display_strategy
    FROM marketplace_listings ml
    INNER JOIN unified_geo ug ON ug.source_type = 'marketplace_listing' AND ug.source_id = ml.id
    LEFT JOIN marketplace_categories c ON ml.category_id = c.id
    WHERE ml.status = 'active' AND ml.show_on_map = true
    
    UNION ALL
    
    -- Витрины
    SELECT 
        s.id,
        s.name,
        s.description,
        0 as price, -- У витрин нет единой цены
        'Витрина' as category_name,
        ug.location,
        ST_Y(ug.location::geometry) as latitude,
        ST_X(ug.location::geometry) as longitude,
        ug.formatted_address,
        s.user_id,
        NULL as storefront_id,
        CASE WHEN s.is_active THEN 'active' ELSE 'inactive' END as status,
        s.created_at,
        s.updated_at,
        s.views_count,
        0 as rating,
        'storefront' as item_type,
        COALESCE(ug.privacy_level, 'exact') as privacy_level,
        COALESCE(ug.blur_radius_meters, 0) as blur_radius_meters,
        'grouped' as display_strategy
    FROM storefronts s
    INNER JOIN unified_geo ug ON ug.source_type = 'storefront' AND ug.source_id = s.id
    WHERE s.is_active = true
    
    UNION ALL
    
    -- Товары витрин
    SELECT 
        sp.id,
        sp.name,
        sp.description,
        COALESCE(spv.price, sp.price) as price,
        spc.name as category_name,
        ug.location,
        ST_Y(ug.location::geometry) as latitude,
        ST_X(ug.location::geometry) as longitude,
        ug.formatted_address,
        s.user_id,
        sp.storefront_id,
        CASE WHEN sp.is_active AND s.is_active THEN 'active' ELSE 'inactive' END as status,
        sp.created_at,
        sp.updated_at,
        sp.view_count as views_count,
        0 as rating,
        'storefront_product' as item_type,
        COALESCE(ug.privacy_level, 'exact') as privacy_level,
        COALESCE(ug.blur_radius_meters, 0) as blur_radius_meters,
        'grouped' as display_strategy
    FROM storefront_products sp
    INNER JOIN storefronts s ON sp.storefront_id = s.id
    INNER JOIN unified_geo ug ON ug.source_type = 'storefront' AND ug.source_id = s.id
    LEFT JOIN storefront_product_variants spv ON sp.id = spv.product_id AND spv.is_default = true
    LEFT JOIN marketplace_categories spc ON sp.category_id = spc.id
    WHERE sp.is_active = true AND s.is_active = true
)
SELECT * FROM combined_items;

-- Создаем индексы для производительности
CREATE INDEX IF NOT EXISTS idx_map_items_cache_location ON map_items_cache USING GIST(location);
CREATE INDEX IF NOT EXISTS idx_map_items_cache_item_type ON map_items_cache(item_type);
CREATE INDEX IF NOT EXISTS idx_map_items_cache_status ON map_items_cache(status);
CREATE INDEX IF NOT EXISTS idx_map_items_cache_user_id ON map_items_cache(user_id);
CREATE INDEX IF NOT EXISTS idx_map_items_cache_storefront_id ON map_items_cache(storefront_id);
CREATE INDEX IF NOT EXISTS idx_map_items_cache_category ON map_items_cache(category_name);
CREATE INDEX IF NOT EXISTS idx_map_items_cache_price ON map_items_cache(price);
CREATE INDEX IF NOT EXISTS idx_map_items_cache_created_at ON map_items_cache(created_at);

-- Создаем функцию для автоматического обновления materialized view
CREATE OR REPLACE FUNCTION refresh_map_items_cache()
RETURNS void AS $$
BEGIN
    REFRESH MATERIALIZED VIEW CONCURRENTLY map_items_cache;
END;
$$ LANGUAGE plpgsql;

-- Комментарий для администраторов
COMMENT ON MATERIALIZED VIEW map_items_cache IS 'Кеш всех элементов для отображения на карте. Обновляется через REFRESH MATERIALIZED VIEW CONCURRENTLY map_items_cache;';