-- Откат изменений enum типа location_privacy_level
-- Возвращаем старые значения: exact, approximate, city_only, hidden

-- Шаг 1: Создаем временный enum со старыми значениями
CREATE TYPE location_privacy_level_old AS ENUM ('exact', 'approximate', 'city_only', 'hidden');

-- Шаг 2: Обновляем все таблицы обратно

-- Обновляем таблицу storefronts (колонка default_privacy_level)
ALTER TABLE storefronts 
    ALTER COLUMN default_privacy_level DROP DEFAULT,
    ALTER COLUMN default_privacy_level TYPE location_privacy_level_old 
    USING CASE 
        WHEN default_privacy_level = 'exact' THEN 'exact'::location_privacy_level_old
        WHEN default_privacy_level = 'street' THEN 'approximate'::location_privacy_level_old
        WHEN default_privacy_level = 'district' THEN 'city_only'::location_privacy_level_old
        WHEN default_privacy_level = 'city' THEN 'city_only'::location_privacy_level_old
    END,
    ALTER COLUMN default_privacy_level SET DEFAULT 'exact'::location_privacy_level_old;

-- Обновляем таблицу storefront_products
ALTER TABLE storefront_products 
    ALTER COLUMN location_privacy DROP DEFAULT,
    ALTER COLUMN location_privacy TYPE location_privacy_level_old 
    USING CASE 
        WHEN location_privacy = 'exact' THEN 'exact'::location_privacy_level_old
        WHEN location_privacy = 'street' THEN 'approximate'::location_privacy_level_old
        WHEN location_privacy = 'district' THEN 'city_only'::location_privacy_level_old
        WHEN location_privacy = 'city' THEN 'city_only'::location_privacy_level_old
    END,
    ALTER COLUMN location_privacy SET DEFAULT 'exact'::location_privacy_level_old;

-- Обновляем таблицу unified_geo (колонка privacy_level)
ALTER TABLE unified_geo 
    ALTER COLUMN privacy_level DROP DEFAULT,
    ALTER COLUMN privacy_level TYPE location_privacy_level_old 
    USING CASE 
        WHEN privacy_level = 'exact' THEN 'exact'::location_privacy_level_old
        WHEN privacy_level = 'street' THEN 'approximate'::location_privacy_level_old
        WHEN privacy_level = 'district' THEN 'city_only'::location_privacy_level_old
        WHEN privacy_level = 'city' THEN 'city_only'::location_privacy_level_old
    END,
    ALTER COLUMN privacy_level SET DEFAULT 'exact'::location_privacy_level_old;

-- Восстанавливаем функцию calculate_blurred_location с старым типом
CREATE OR REPLACE FUNCTION calculate_blurred_location(
    exact_location geography,
    privacy_level location_privacy_level_old
) RETURNS geography AS $$
BEGIN
    CASE privacy_level
        WHEN 'exact' THEN
            RETURN exact_location;
        WHEN 'approximate' THEN
            -- Размываем до ~200 метров
            RETURN ST_Buffer(
                ST_Centroid(exact_location::geometry)::geography,
                200 + (random() * 100 - 50) -- 200м ± 50м
            );
        WHEN 'city_only' THEN
            -- Размываем до ~1000 метров
            RETURN ST_Buffer(
                ST_Centroid(exact_location::geometry)::geography,
                1000 + (random() * 200 - 100) -- 1км ± 100м
            );
        WHEN 'hidden' THEN
            RETURN NULL;
        ELSE
            RETURN NULL;
    END CASE;
END;
$$ LANGUAGE plpgsql IMMUTABLE;

-- Не пересоздаем материализованное представление, так как оно будет пересоздано в последующих миграциях

-- Шаг 3: Удаляем новый тип и переименовываем старый
DROP TYPE location_privacy_level;
ALTER TYPE location_privacy_level_old RENAME TO location_privacy_level;

-- Шаг 4: Восстанавливаем комментарий к типу
COMMENT ON TYPE location_privacy_level IS 'Privacy level for location display: exact, approximate, city_only, hidden';