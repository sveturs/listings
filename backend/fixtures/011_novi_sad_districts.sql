-- Добавление районов Нови Сада
-- Сначала получаем ID города
DO $$
DECLARE
    novi_sad_id uuid;
BEGIN
    SELECT id INTO novi_sad_id FROM cities WHERE slug = 'novi-sad' LIMIT 1;
    
    IF novi_sad_id IS NOT NULL THEN
        -- Лиман - популярный жилой район
        INSERT INTO districts (name, city_id, country_code, boundary, center_point, population, area_km2, postal_codes)
        VALUES (
            'Лиман',
            novi_sad_id,
            'RS',
            ST_GeomFromText('POLYGON((19.8200 45.2500, 19.8400 45.2500, 19.8400 45.2300, 19.8200 45.2300, 19.8200 45.2500))', 4326),
            ST_GeomFromText('POINT(19.8300 45.2400)', 4326),
            35000,
            4.5,
            ARRAY['21000']
        );

        -- Ново насеље
        INSERT INTO districts (name, city_id, country_code, boundary, center_point, population, area_km2, postal_codes)
        VALUES (
            'Ново насеље',
            novi_sad_id,
            'RS',
            ST_GeomFromText('POLYGON((19.8100 45.2700, 19.8300 45.2700, 19.8300 45.2500, 19.8100 45.2500, 19.8100 45.2700))', 4326),
            ST_GeomFromText('POINT(19.8200 45.2600)', 4326),
            28000,
            3.8,
            ARRAY['21000']
        );

        -- Детелинара
        INSERT INTO districts (name, city_id, country_code, boundary, center_point, population, area_km2, postal_codes)
        VALUES (
            'Детелинара',
            novi_sad_id,
            'RS',
            ST_GeomFromText('POLYGON((19.8000 45.2600, 19.8200 45.2600, 19.8200 45.2400, 19.8000 45.2400, 19.8000 45.2600))', 4326),
            ST_GeomFromText('POINT(19.8100 45.2500)', 4326),
            25000,
            3.2,
            ARRAY['21000']
        );

        -- Грбавица
        INSERT INTO districts (name, city_id, country_code, boundary, center_point, population, area_km2, postal_codes)
        VALUES (
            'Грбавица',
            novi_sad_id,
            'RS',
            ST_GeomFromText('POLYGON((19.8300 45.2650, 19.8500 45.2650, 19.8500 45.2450, 19.8300 45.2450, 19.8300 45.2650))', 4326),
            ST_GeomFromText('POINT(19.8400 45.2550)', 4326),
            32000,
            4.1,
            ARRAY['21000']
        );

        -- Стари град
        INSERT INTO districts (name, city_id, country_code, boundary, center_point, population, area_km2, postal_codes)
        VALUES (
            'Стари град',
            novi_sad_id,
            'RS',
            ST_GeomFromText('POLYGON((19.8400 45.2600, 19.8600 45.2600, 19.8600 45.2400, 19.8400 45.2400, 19.8400 45.2600))', 4326),
            ST_GeomFromText('POINT(19.8500 45.2500)', 4326),
            18000,
            2.5,
            ARRAY['21000']
        );
    END IF;
END $$;