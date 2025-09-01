-- Добавление адресов и координат товарам витрин на основе данных их витрин
-- Все товары витрин наследуют местоположение от своих витрин если у них нет индивидуального местоположения

-- Обновляем товары витрин, у которых нет индивидуального местоположения
-- Устанавливаем адрес и координаты из их витрин
UPDATE storefront_products sp
SET 
    individual_address = COALESCE(s.address, '') || 
                        CASE 
                            WHEN s.city IS NOT NULL AND s.city != '' THEN ', ' || s.city 
                            ELSE '' 
                        END ||
                        CASE 
                            WHEN s.country IS NOT NULL AND s.country != '' THEN ', ' || s.country 
                            ELSE '' 
                        END,
    individual_latitude = s.latitude,
    individual_longitude = s.longitude,
    has_individual_location = true,
    location_privacy = 'exact',
    show_on_map = true
FROM storefronts s
WHERE sp.storefront_id = s.id
    AND sp.has_individual_location = false
    AND s.latitude IS NOT NULL
    AND s.longitude IS NOT NULL;

-- Для товаров витрин без координат используем случайные координаты в Сербии
-- Создаем временную таблицу с локациями
CREATE TEMP TABLE IF NOT EXISTS temp_serbian_locations (
    location_name VARCHAR(100),
    lat NUMERIC(10,8),
    lng NUMERIC(11,8)
);

INSERT INTO temp_serbian_locations VALUES
    ('Белград, Стари Град', 44.8125, 20.4612),
    ('Белград, Нови Београд', 44.8176, 20.4190),
    ('Белград, Земун', 44.8433, 20.4011),
    ('Белград, Палилула', 44.8150, 20.4750),
    ('Белград, Звездара', 44.8047, 20.5080),
    ('Белград, Вождовац', 44.7730, 20.4790),
    ('Белград, Чукарица', 44.7580, 20.4170),
    ('Белград, Раковица', 44.7570, 20.4440),
    ('Нови Сад, Центр', 45.2551, 19.8451),
    ('Нови Сад, Лиман', 45.2380, 19.8360),
    ('Нови Сад, Ново Насеље', 45.2480, 19.8130),
    ('Нови Сад, Детелинара', 45.2590, 19.8180),
    ('Ниш, Центр', 43.3209, 21.8958),
    ('Крагујевац, Центр', 44.0128, 20.9114),
    ('Суботица, Центр', 46.1011, 19.6651),
    ('Зрењанин, Центр', 45.3639, 20.3861),
    ('Панчево, Центр', 44.8736, 20.6403),
    ('Смедерево, Центр', 44.6650, 20.9289),
    ('Вршац, Центр', 45.1167, 21.3033),
    ('Шабац, Центр', 44.7463, 19.6908);

-- Обновляем товары без координат, распределяя их по локациям
UPDATE storefront_products sp
SET 
    individual_address = tsl.location_name,
    individual_latitude = tsl.lat + (random() - 0.5) * 0.01,
    individual_longitude = tsl.lng + (random() - 0.5) * 0.01,
    has_individual_location = true,
    location_privacy = 'exact',
    show_on_map = true
FROM (
    SELECT sp2.id,
           tsl2.location_name,
           tsl2.lat,
           tsl2.lng,
           ROW_NUMBER() OVER (ORDER BY sp2.id) as rn
    FROM storefront_products sp2
    CROSS JOIN temp_serbian_locations tsl2
    WHERE sp2.has_individual_location = false
        OR sp2.individual_latitude IS NULL
        OR sp2.individual_longitude IS NULL
    ORDER BY sp2.id, random()
) AS distribution
JOIN temp_serbian_locations tsl ON tsl.location_name = distribution.location_name
WHERE sp.id = distribution.id
    AND distribution.rn = (
        SELECT MIN(d2.rn) 
        FROM (
            SELECT sp3.id, ROW_NUMBER() OVER (ORDER BY sp3.id) as rn
            FROM storefront_products sp3
            WHERE sp3.has_individual_location = false
                OR sp3.individual_latitude IS NULL
                OR sp3.individual_longitude IS NULL
        ) d2 
        WHERE d2.id = sp.id
    );

DROP TABLE IF EXISTS temp_serbian_locations;

-- Добавляем адреса объявлениям marketplace, у которых их еще нет
UPDATE marketplace_listings ml
SET 
    location = CASE 
        WHEN ml.location IS NULL OR ml.location = '' THEN
            COALESCE(ml.address_city, '') ||
            CASE 
                WHEN ml.address_country IS NOT NULL AND ml.address_country != '' THEN ', ' || ml.address_country 
                ELSE ', Сербия' 
            END
        ELSE ml.location
    END
WHERE (ml.location IS NULL OR ml.location = '')
    AND ml.latitude IS NOT NULL
    AND ml.longitude IS NOT NULL;

-- Проставляем координаты для объявлений без них
CREATE TEMP TABLE IF NOT EXISTS temp_cities (
    city VARCHAR(100),
    country VARCHAR(100),
    lat NUMERIC(10,8),
    lng NUMERIC(11,8)
);

INSERT INTO temp_cities VALUES
    ('Белград', 'Сербия', 44.8125, 20.4612),
    ('Нови Сад', 'Сербия', 45.2551, 19.8451),
    ('Ниш', 'Сербия', 43.3209, 21.8958),
    ('Крагујевац', 'Сербия', 44.0128, 20.9114),
    ('Суботица', 'Сербия', 46.1011, 19.6651);

UPDATE marketplace_listings ml
SET 
    latitude = tc.lat + (random() - 0.5) * 0.05,
    longitude = tc.lng + (random() - 0.5) * 0.05,
    address_city = COALESCE(ml.address_city, tc.city),
    address_country = COALESCE(ml.address_country, tc.country),
    location = COALESCE(ml.location, tc.city || ', ' || tc.country)
FROM temp_cities tc
WHERE (ml.latitude IS NULL OR ml.longitude IS NULL)
    AND tc.city = (
        SELECT city FROM temp_cities 
        ORDER BY random() 
        LIMIT 1
    );

DROP TABLE IF EXISTS temp_cities;