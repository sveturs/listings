-- Проверяем существование столбца storefront_id и добавляем его, если он отсутствует
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 
        FROM information_schema.columns 
        WHERE table_name = 'marketplace_listings' AND column_name = 'storefront_id'
    ) THEN
        ALTER TABLE marketplace_listings ADD COLUMN storefront_id INT REFERENCES user_storefronts(id) ON DELETE SET NULL;
        CREATE INDEX IF NOT EXISTS idx_marketplace_listings_storefront ON marketplace_listings(storefront_id);
    END IF;
END$$;

-- Оптимизированное материализованное представление
DO $$
BEGIN
    -- Проверяем, существует ли представление
    IF NOT EXISTS (
        SELECT 1 
        FROM pg_matviews 
        WHERE matviewname = 'category_listing_counts'
    ) THEN
        CREATE MATERIALIZED VIEW category_listing_counts AS
        WITH RECURSIVE category_tree AS (
            SELECT 
                id,
                ARRAY[id] as category_path,
                name,
                1 as depth,
                (
                    SELECT COUNT(*) 
                    FROM marketplace_listings ml 
                    WHERE ml.category_id = marketplace_categories.id 
                    AND ml.status = 'active'
                ) as direct_count
            FROM marketplace_categories
            WHERE parent_id IS NULL
            
            UNION ALL
            
            SELECT 
                c.id,
                ct.category_path || c.id,
                c.name,
                ct.depth + 1,
                (
                    SELECT COUNT(*) 
                    FROM marketplace_listings ml 
                    WHERE ml.category_id = c.id 
                    AND ml.status = 'active'
                ) as direct_count
            FROM marketplace_categories c
            INNER JOIN category_tree ct ON c.parent_id = ct.id
            WHERE ct.depth < 10
        )
        SELECT 
            ct.id as category_id,
            ct.direct_count + COALESCE((
                SELECT SUM(ch.direct_count)
                FROM category_tree ch
                WHERE ch.category_path[1:array_length(ct.category_path, 1)] = ct.category_path
                AND ch.id != ct.id
            ), 0) as listing_count,
            MAX(ct.depth) as category_depth
        FROM category_tree ct
        GROUP BY ct.id, ct.direct_count, ct.category_path;

        -- Создаем уникальный индекс для возможности CONCURRENT обновления
        CREATE UNIQUE INDEX category_listing_counts_idx ON category_listing_counts(category_id);
    ELSE
        -- Если представление уже существует, обновляем его
        REFRESH MATERIALIZED VIEW category_listing_counts;
    END IF;
END $$;

-- Создаем функцию обновления, если она не существует или заменяем, если существует
CREATE OR REPLACE FUNCTION refresh_category_listing_counts()
RETURNS TRIGGER AS $$
DECLARE
    current_ts TIMESTAMP;
    last_refresh TIMESTAMP;
BEGIN
    -- Проверяем, не обновлялось ли представление в последние N секунд
    SELECT INTO last_refresh COALESCE(
        (SELECT obj_description('category_listing_counts'::regclass)::timestamp),
        '1970-01-01'::timestamp
    );
    
    current_ts := CURRENT_TIMESTAMP;
    
    IF current_ts - last_refresh > interval '5 seconds' THEN
        -- Обновляем представление
        REFRESH MATERIALIZED VIEW category_listing_counts;
        
        -- Сохраняем время последнего обновления
        EXECUTE format(
            'COMMENT ON MATERIALIZED VIEW category_listing_counts IS %L',
            current_ts::text
        );
    END IF;
    
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- Создаем триггеры только если они не существуют
DO $$
BEGIN
    -- Для INSERT
    IF NOT EXISTS (
        SELECT 1 
        FROM pg_trigger 
        WHERE tgname = 'refresh_category_counts_insert'
    ) THEN
        CREATE TRIGGER refresh_category_counts_insert
            AFTER INSERT ON marketplace_listings
            FOR EACH ROW
            EXECUTE FUNCTION refresh_category_listing_counts();
    END IF;

    -- Для UPDATE
    IF NOT EXISTS (
        SELECT 1 
        FROM pg_trigger 
        WHERE tgname = 'refresh_category_counts_update'
    ) THEN
        CREATE TRIGGER refresh_category_counts_update
            AFTER UPDATE ON marketplace_listings
            FOR EACH ROW
            WHEN (OLD.status IS DISTINCT FROM NEW.status)
            EXECUTE FUNCTION refresh_category_listing_counts();
    END IF;

    -- Для DELETE
    IF NOT EXISTS (
        SELECT 1 
        FROM pg_trigger 
        WHERE tgname = 'refresh_category_counts_delete'
    ) THEN
        CREATE TRIGGER refresh_category_counts_delete
            AFTER DELETE ON marketplace_listings
            FOR EACH ROW
            EXECUTE FUNCTION refresh_category_listing_counts();
    END IF;
END $$;

-- Делаем начальное обновление
REFRESH MATERIALIZED VIEW CONCURRENTLY category_listing_counts;