-- Добавление районов Ниша
-- Сначала получаем ID города
DO $$
DECLARE
    nis_id uuid;
BEGIN
    SELECT id INTO nis_id FROM cities WHERE slug = 'nis' LIMIT 1;
    
    IF nis_id IS NOT NULL THEN
        -- Медијана - центральный район
        INSERT INTO districts (name, city_id, country_code, boundary, center_point, population, area_km2, postal_codes)
        VALUES (
            'Медијана',
            nis_id,
            'RS',
            ST_GeomFromText('POLYGON((21.8800 43.3200, 21.9100 43.3200, 21.9100 43.3000, 21.8800 43.3000, 21.8800 43.3200))', 4326),
            ST_GeomFromText('POINT(21.8950 43.3100)', 4326),
            85000,
            16.0,
            ARRAY['18000']
        );

        -- Палилула
        INSERT INTO districts (name, city_id, country_code, boundary, center_point, population, area_km2, postal_codes)
        VALUES (
            'Палилула',
            nis_id,
            'RS',
            ST_GeomFromText('POLYGON((21.9100 43.3300, 21.9400 43.3300, 21.9400 43.3100, 21.9100 43.3100, 21.9100 43.3300))', 4326),
            ST_GeomFromText('POINT(21.9250 43.3200)', 4326),
            73000,
            117.0,
            ARRAY['18000']
        );

        -- Црвени Крст
        INSERT INTO districts (name, city_id, country_code, boundary, center_point, population, area_km2, postal_codes)
        VALUES (
            'Црвени Крст',
            nis_id,
            'RS',
            ST_GeomFromText('POLYGON((21.8700 43.3100, 21.9000 43.3100, 21.9000 43.2900, 21.8700 43.2900, 21.8700 43.3100))', 4326),
            ST_GeomFromText('POINT(21.8850 43.3000)', 4326),
            32000,
            180.0,
            ARRAY['18000']
        );

        -- Пантелеј
        INSERT INTO districts (name, city_id, country_code, boundary, center_point, population, area_km2, postal_codes)
        VALUES (
            'Пантелеј',
            nis_id,
            'RS',
            ST_GeomFromText('POLYGON((21.8600 43.3400, 21.8900 43.3400, 21.8900 43.3200, 21.8600 43.3200, 21.8600 43.3400))', 4326),
            ST_GeomFromText('POINT(21.8750 43.3300)', 4326),
            53000,
            142.0,
            ARRAY['18000']
        );

        -- Нишка Бања
        INSERT INTO districts (name, city_id, country_code, boundary, center_point, population, area_km2, postal_codes)
        VALUES (
            'Нишка Бања',
            nis_id,
            'RS',
            ST_GeomFromText('POLYGON((21.9200 43.3000, 21.9500 43.3000, 21.9500 43.2800, 21.9200 43.2800, 21.9200 43.3000))', 4326),
            ST_GeomFromText('POINT(21.9350 43.2900)', 4326),
            15000,
            145.0,
            ARRAY['18205']
        );
    END IF;
END $$;