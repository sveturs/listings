-- Оптимизированное материализованное представление
CREATE MATERIALIZED VIEW category_listing_counts AS
WITH RECURSIVE category_tree AS (
    -- Базовый случай: корневые категории
    SELECT 
        id,
        ARRAY[id] as category_path,
        name,
        1 as depth
    FROM marketplace_categories
    WHERE parent_id IS NULL
    
    UNION ALL
    
    -- Рекурсивная часть: дочерние категории
    SELECT 
        c.id,
        ct.category_path || c.id,
        c.name,
        ct.depth + 1
    FROM marketplace_categories c
    INNER JOIN category_tree ct ON c.parent_id = ct.id
    WHERE ct.depth < 10  -- Ограничиваем глубину рекурсии
)
SELECT 
    ct.id as category_id,
    COUNT(DISTINCT l.id) as listing_count,
    MAX(ct.depth) as category_depth
FROM category_tree ct
LEFT JOIN marketplace_listings l ON (
    l.category_id = ct.id 
    AND l.status = 'active'
)
GROUP BY ct.id;

-- Создаем уникальный индекс для возможности CONCURRENT обновления
CREATE UNIQUE INDEX category_listing_counts_idx ON category_listing_counts(category_id);

-- Оптимизированная функция обновления
CREATE OR REPLACE FUNCTION refresh_category_listing_counts()
RETURNS TRIGGER AS $$
DECLARE
    current_ts TIMESTAMP;  -- Изменили имя переменной
    last_refresh TIMESTAMP;
BEGIN
    -- Проверяем, не обновлялось ли представление в последние N секунд
    SELECT INTO last_refresh COALESCE(
        (SELECT obj_description('category_listing_counts'::regclass)::timestamp),
        '1970-01-01'::timestamp
    );
    
    current_ts := CURRENT_TIMESTAMP;  -- Используем новое имя переменной вместо current_time
    
    IF current_ts - last_refresh > interval '5 seconds' THEN
        -- Обновляем представление
        REFRESH MATERIALIZED VIEW category_listing_counts;
        
        -- Сохраняем время последнего обновления
        EXECUTE format(
            'COMMENT ON MATERIALIZED VIEW category_listing_counts IS %L',
            current_ts::text  -- Используем новое имя переменной
        );
    END IF;
    
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- Создаем триггеры с оптимизированной функцией
CREATE TRIGGER refresh_category_counts_insert
    AFTER INSERT ON marketplace_listings
    FOR EACH ROW
    EXECUTE FUNCTION refresh_category_listing_counts();

CREATE TRIGGER refresh_category_counts_update
    AFTER UPDATE ON marketplace_listings
    FOR EACH ROW
    WHEN (OLD.status IS DISTINCT FROM NEW.status)
    EXECUTE FUNCTION refresh_category_listing_counts();

CREATE TRIGGER refresh_category_counts_delete
    AFTER DELETE ON marketplace_listings
    FOR EACH ROW
    EXECUTE FUNCTION refresh_category_listing_counts();

-- Делаем начальное обновление
REFRESH MATERIALIZED VIEW category_listing_counts;