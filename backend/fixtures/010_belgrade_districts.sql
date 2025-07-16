-- Добавление районов Белграда с упрощенными полигонами
-- Координаты представляют примерные границы районов

DO $$
DECLARE
    belgrade_id uuid;
BEGIN
    SELECT id INTO belgrade_id FROM cities WHERE slug = 'belgrade' LIMIT 1;
    
    IF belgrade_id IS NOT NULL THEN
        -- Стари Град (Старый город) - центральный район
        INSERT INTO districts (name, city_id, country_code, boundary, center_point, population, area_km2, postal_codes)
        VALUES (
            'Стари Град',
            belgrade_id,
            'RS',
            ST_GeomFromText('POLYGON((20.4500 44.8200, 20.4650 44.8200, 20.4650 44.8100, 20.4500 44.8100, 20.4500 44.8200))', 4326),
            ST_GeomFromText('POINT(20.4575 44.8150)', 4326),
            45000,
            3.5,
            ARRAY['11000']
        );

        -- Врачар - густонаселенный район
        INSERT INTO districts (name, city_id, country_code, boundary, center_point, population, area_km2, postal_codes)
        VALUES (
            'Врачар',
            belgrade_id,
            'RS',
            ST_GeomFromText('POLYGON((20.4650 44.8100, 20.4900 44.8100, 20.4900 44.7900, 20.4650 44.7900, 20.4650 44.8100))', 4326),
            ST_GeomFromText('POINT(20.4775 44.8000)', 4326),
            55000,
            3.0,
            ARRAY['11000', '11010']
        );

        -- Савски Венац - престижный район
        INSERT INTO districts (name, city_id, country_code, boundary, center_point, population, area_km2, postal_codes)
        VALUES (
            'Савски венац',
            belgrade_id,
            'RS',
            ST_GeomFromText('POLYGON((20.4300 44.8000, 20.4600 44.8000, 20.4600 44.7700, 20.4300 44.7700, 20.4300 44.8000))', 4326),
            ST_GeomFromText('POINT(20.4450 44.7850)', 4326),
            40000,
            14.0,
            ARRAY['11000', '11040']
        );

        -- Нови Београд - современный район на левом берегу Савы
        INSERT INTO districts (name, city_id, country_code, boundary, center_point, population, area_km2, postal_codes)
        VALUES (
            'Нови Београд',
            belgrade_id,
            'RS',
            ST_GeomFromText('POLYGON((20.3800 44.8200, 20.4300 44.8200, 20.4300 44.7800, 20.3800 44.7800, 20.3800 44.8200))', 4326),
            ST_GeomFromText('POINT(20.4050 44.8000)', 4326),
            220000,
            40.7,
            ARRAY['11070', '11073', '11077']
        );

        -- Земун - исторический район
        INSERT INTO districts (name, city_id, country_code, boundary, center_point, population, area_km2, postal_codes)
        VALUES (
            'Земун',
            belgrade_id,
            'RS',
            ST_GeomFromText('POLYGON((20.3500 44.8500, 20.4200 44.8500, 20.4200 44.8200, 20.3500 44.8200, 20.3500 44.8500))', 4326),
            ST_GeomFromText('POINT(20.3850 44.8350)', 4326),
            170000,
            154.0,
            ARRAY['11080', '11185']
        );

        -- Звездара - восточный район
        INSERT INTO districts (name, city_id, country_code, boundary, center_point, population, area_km2, postal_codes)
        VALUES (
            'Звездара',
            belgrade_id,
            'RS',
            ST_GeomFromText('POLYGON((20.4900 44.8100, 20.5400 44.8100, 20.5400 44.7700, 20.4900 44.7700, 20.4900 44.8100))', 4326),
            ST_GeomFromText('POINT(20.5150 44.7900)', 4326),
            150000,
            32.0,
            ARRAY['11050', '11120', '11160']
        );

        -- Палилула - крупнейший район по площади
        INSERT INTO districts (name, city_id, country_code, boundary, center_point, population, area_km2, postal_codes)
        VALUES (
            'Палилула',
            belgrade_id,
            'RS',
            ST_GeomFromText('POLYGON((20.4650 44.8200, 20.5500 44.8200, 20.5500 44.7800, 20.4650 44.7800, 20.4650 44.8200))', 4326),
            ST_GeomFromText('POINT(20.5075 44.8000)', 4326),
            180000,
            451.0,
            ARRAY['11000', '11060', '11090']
        );

        -- Раковица - юго-западный район
        INSERT INTO districts (name, city_id, country_code, boundary, center_point, population, area_km2, postal_codes)
        VALUES (
            'Раковица',
            belgrade_id,
            'RS',
            ST_GeomFromText('POLYGON((20.4100 44.7700, 20.4500 44.7700, 20.4500 44.7400, 20.4100 44.7400, 20.4100 44.7700))', 4326),
            ST_GeomFromText('POINT(20.4300 44.7550)', 4326),
            110000,
            29.0,
            ARRAY['11090', '11283']
        );

        -- Чукарица - юго-западный район
        INSERT INTO districts (name, city_id, country_code, boundary, center_point, population, area_km2, postal_codes)
        VALUES (
            'Чукарица',
            belgrade_id,
            'RS',
            ST_GeomFromText('POLYGON((20.3700 44.7800, 20.4300 44.7800, 20.4300 44.7400, 20.3700 44.7400, 20.3700 44.7800))', 4326),
            ST_GeomFromText('POINT(20.4000 44.7600)', 4326),
            180000,
            156.0,
            ARRAY['11030', '11250', '11253']
        );

        -- Вождовац - юго-восточный район
        INSERT INTO districts (name, city_id, country_code, boundary, center_point, population, area_km2, postal_codes)
        VALUES (
            'Вождовац',
            belgrade_id,
            'RS',
            ST_GeomFromText('POLYGON((20.4600 44.7700, 20.5100 44.7700, 20.5100 44.7300, 20.4600 44.7300, 20.4600 44.7700))', 4326),
            ST_GeomFromText('POINT(20.4850 44.7500)', 4326),
            160000,
            148.0,
            ARRAY['11010', '11210', '11211']
        );
    END IF;
END $$;