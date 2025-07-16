-- Добавление основных городов Сербии с координатами

-- Белград
INSERT INTO cities (name, slug, country_code, center_point, has_districts, priority, population, area_km2)
VALUES (
    'Belgrade',
    'belgrade',
    'RS',
    ST_GeomFromText('POINT(20.4489 44.7866)', 4326),
    true,
    1,
    1378682,
    359.96
) ON CONFLICT (slug) DO NOTHING;

-- Нови Сад
INSERT INTO cities (name, slug, country_code, center_point, has_districts, priority, population, area_km2)
VALUES (
    'Novi Sad',
    'novi-sad',
    'RS',
    ST_GeomFromText('POINT(19.8335 45.2671)', 4326),
    true,
    2,
    341625,
    699.0
) ON CONFLICT (slug) DO NOTHING;

-- Ниш
INSERT INTO cities (name, slug, country_code, center_point, has_districts, priority, population, area_km2)
VALUES (
    'Niš',
    'nis',
    'RS',
    ST_GeomFromText('POINT(21.8958 43.3209)', 4326),
    true,
    3,
    260237,
    597.0
) ON CONFLICT (slug) DO NOTHING;

-- Крагујевац
INSERT INTO cities (name, slug, country_code, center_point, has_districts, priority, population, area_km2)
VALUES (
    'Kragujevac',
    'kragujevac',
    'RS',
    ST_GeomFromText('POINT(20.9167 44.0167)', 4326),
    true,
    4,
    179417,
    835.0
) ON CONFLICT (slug) DO NOTHING;

-- Суботица
INSERT INTO cities (name, slug, country_code, center_point, has_districts, priority, population, area_km2)
VALUES (
    'Subotica',
    'subotica',
    'RS',
    ST_GeomFromText('POINT(19.6669 46.1002)', 4326),
    true,
    5,
    141554,
    1007.0
) ON CONFLICT (slug) DO NOTHING;