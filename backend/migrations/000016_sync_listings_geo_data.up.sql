-- Синхронизация геоданных из marketplace_listings в listings_geo
-- Заполнение таблицы listings_geo данными из marketplace_listings

-- Вставка геоданных из marketplace_listings
INSERT INTO listings_geo (
    listing_id,
    location,
    geohash,
    is_precise,
    blur_radius,
    address_components,
    geocoding_confidence,
    address_verified,
    input_method,
    location_privacy,
    blurred_location,
    formatted_address,
    created_at,
    updated_at
)
SELECT
    ml.id as listing_id,
    ST_SetSRID(ST_MakePoint(ml.longitude, ml.latitude)::geography, 4326) as location,
    -- Генерация geohash с точностью 7 символов (~150м)
    ST_GeoHash(ST_SetSRID(ST_MakePoint(ml.longitude, ml.latitude), 4326), 7) as geohash,
    -- Приватность зависит от типа объявления (есть ли storefront_id)
    CASE
        WHEN ml.storefront_id IS NOT NULL THEN true  -- Бизнесы показывают точное расположение
        ELSE false  -- Частные лица скрывают точное расположение
    END as is_precise,
    -- Радиус размытия для приватности
    CASE
        WHEN ml.storefront_id IS NOT NULL THEN 0  -- Без размытия для бизнесов
        ELSE 500  -- 500м размытие для частных лиц
    END as blur_radius,
    -- Компоненты адреса в JSON
    jsonb_build_object(
        'city', ml.address_city,
        'country', ml.address_country,
        'formatted', CONCAT_WS(', ',
            NULLIF(ml.address_city, ''),
            COALESCE(ml.address_country, 'Serbia')
        )
    ) as address_components,
    -- Уверенность в геокодировании (предполагаем высокую для существующих данных)
    0.95 as geocoding_confidence,
    -- Считаем адреса проверенными для существующих данных
    true as address_verified,
    -- Метод ввода - предполагаем ручной для существующих
    'manual' as input_method,
    -- Уровень приватности
    CASE
        WHEN ml.storefront_id IS NOT NULL THEN 'exact'  -- Точное для бизнесов
        ELSE 'street'  -- Уровень улицы для частных
    END as location_privacy,
    -- Размытые координаты для приватности
    CASE
        WHEN ml.storefront_id IS NOT NULL THEN
            ST_SetSRID(ST_MakePoint(ml.longitude, ml.latitude)::geography, 4326)
        ELSE
            -- Добавляем случайное смещение до 500м для приватности
            ST_SetSRID(
                ST_MakePoint(
                    ml.longitude + (random() - 0.5) * 0.009,  -- ~500м в долготе
                    ml.latitude + (random() - 0.5) * 0.009    -- ~500м в широте
                )::geography,
                4326
            )
    END as blurred_location,
    -- Форматированный адрес
    CONCAT_WS(', ',
        NULLIF(ml.address_city, ''),
        COALESCE(ml.address_country, 'Serbia')
    ) as formatted_address,
    ml.created_at,
    ml.updated_at
FROM marketplace_listings ml
LEFT JOIN users u ON ml.user_id = u.id
WHERE ml.latitude IS NOT NULL
  AND ml.longitude IS NOT NULL
  AND NOT EXISTS (
      -- Проверяем, что записи еще нет в listings_geo
      SELECT 1 FROM listings_geo lg WHERE lg.listing_id = ml.id
  );

-- Создание индекса для оптимизации поиска по privacy_level (если его нет)
CREATE INDEX IF NOT EXISTS idx_listings_geo_privacy_level ON listings_geo(location_privacy);

-- Создание индекса для оптимизации поиска по listing_id
CREATE INDEX IF NOT EXISTS idx_listings_geo_listing_id ON listings_geo(listing_id);

-- Анализ таблиц для обновления статистики
ANALYZE listings_geo;