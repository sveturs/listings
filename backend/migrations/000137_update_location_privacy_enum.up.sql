-- Обновление enum типа location_privacy_level для соответствия новой документации
-- Старые значения: exact, approximate, city_only, hidden
-- Новые значения: exact, street, district, city

-- Шаг 1: Создаем временный enum с новыми значениями
CREATE TYPE location_privacy_level_new AS ENUM ('exact', 'street', 'district', 'city');

-- Шаг 2: Сначала удаляем зависимые объекты
DROP MATERIALIZED VIEW IF EXISTS map_items_cache CASCADE;
DROP FUNCTION IF EXISTS calculate_blurred_location(numeric, numeric, location_privacy_level);

-- Шаг 3: Обновляем все таблицы, которые используют этот тип

-- Обновляем таблицу storefronts (колонка default_privacy_level)
ALTER TABLE storefronts 
    ALTER COLUMN default_privacy_level DROP DEFAULT,
    ALTER COLUMN default_privacy_level TYPE location_privacy_level_new 
    USING CASE 
        WHEN default_privacy_level = 'exact' THEN 'exact'::location_privacy_level_new
        WHEN default_privacy_level = 'approximate' THEN 'street'::location_privacy_level_new
        WHEN default_privacy_level = 'city_only' THEN 'city'::location_privacy_level_new
        WHEN default_privacy_level = 'hidden' THEN 'city'::location_privacy_level_new
    END,
    ALTER COLUMN default_privacy_level SET DEFAULT 'exact'::location_privacy_level_new;

-- Обновляем таблицу storefront_products
ALTER TABLE storefront_products 
    ALTER COLUMN location_privacy DROP DEFAULT,
    ALTER COLUMN location_privacy TYPE location_privacy_level_new 
    USING CASE 
        WHEN location_privacy = 'exact' THEN 'exact'::location_privacy_level_new
        WHEN location_privacy = 'approximate' THEN 'street'::location_privacy_level_new
        WHEN location_privacy = 'city_only' THEN 'city'::location_privacy_level_new
        WHEN location_privacy = 'hidden' THEN 'city'::location_privacy_level_new
    END,
    ALTER COLUMN location_privacy SET DEFAULT 'exact'::location_privacy_level_new;

-- Обновляем таблицу unified_geo (колонка privacy_level)
ALTER TABLE unified_geo 
    ALTER COLUMN privacy_level DROP DEFAULT,
    ALTER COLUMN privacy_level TYPE location_privacy_level_new 
    USING CASE 
        WHEN privacy_level = 'exact' THEN 'exact'::location_privacy_level_new
        WHEN privacy_level = 'approximate' THEN 'street'::location_privacy_level_new
        WHEN privacy_level = 'city_only' THEN 'city'::location_privacy_level_new
        WHEN privacy_level = 'hidden' THEN 'city'::location_privacy_level_new
    END,
    ALTER COLUMN privacy_level SET DEFAULT 'exact'::location_privacy_level_new;

-- Обновляем любые функции, которые используют этот тип
-- Пересоздаем функцию calculate_blurred_location с новым типом
CREATE OR REPLACE FUNCTION calculate_blurred_location(
    exact_location geography,
    privacy_level location_privacy_level_new
) RETURNS geography AS $$
BEGIN
    CASE privacy_level
        WHEN 'exact' THEN
            RETURN exact_location;
        WHEN 'street' THEN
            -- Размываем до ~200 метров
            RETURN ST_Buffer(
                ST_Centroid(exact_location::geometry)::geography,
                200 + (random() * 100 - 50) -- 200м ± 50м
            );
        WHEN 'district' THEN
            -- Размываем до ~1000 метров
            RETURN ST_Buffer(
                ST_Centroid(exact_location::geometry)::geography,
                1000 + (random() * 200 - 100) -- 1км ± 100м
            );
        WHEN 'city' THEN
            -- Размываем до ~5000 метров
            RETURN ST_Buffer(
                ST_Centroid(exact_location::geometry)::geography,
                5000 + (random() * 1000 - 500) -- 5км ± 500м
            );
        ELSE
            RETURN NULL;
    END CASE;
END;
$$ LANGUAGE plpgsql IMMUTABLE;

-- Шаг 4: Не пересоздаем материализованное представление, так как оно будет пересоздано в последующих миграциях

-- Шаг 5: Удаляем старый тип и переименовываем новый
DROP TYPE location_privacy_level;
ALTER TYPE location_privacy_level_new RENAME TO location_privacy_level;

-- Шаг 6: Обновляем комментарий к типу
COMMENT ON TYPE location_privacy_level IS 'Privacy level for location display: exact (точный адрес), street (улица), district (район), city (город)';