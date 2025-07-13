-- Добавление районов Белграда с упрощенными полигонами
-- Координаты представляют примерные границы районов

-- Стари Град (Старый город) - центральный район
INSERT INTO districts (name, country_code, boundary, center_point, population, area_km2, postal_codes)
VALUES (
    'Стари Град',
    'RS',
    ST_GeomFromText('POLYGON((20.4500 44.8200, 20.4650 44.8200, 20.4650 44.8100, 20.4500 44.8100, 20.4500 44.8200))', 4326),
    ST_GeomFromText('POINT(20.4575 44.8150)', 4326),
    45000,
    3.5,
    ARRAY['11000']
);

-- Врачар - густонаселенный район
INSERT INTO districts (name, country_code, boundary, center_point, population, area_km2, postal_codes)
VALUES (
    'Врачар',
    'RS',
    ST_GeomFromText('POLYGON((20.4650 44.8100, 20.4900 44.8100, 20.4900 44.7900, 20.4650 44.7900, 20.4650 44.8100))', 4326),
    ST_GeomFromText('POINT(20.4775 44.8000)', 4326),
    55000,
    3.0,
    ARRAY['11000', '11010']
);

-- Савски Венац - престижный район
INSERT INTO districts (name, country_code, boundary, center_point, population, area_km2, postal_codes)
VALUES (
    'Савски венац',
    'RS',
    ST_GeomFromText('POLYGON((20.4300 44.8000, 20.4600 44.8000, 20.4600 44.7700, 20.4300 44.7700, 20.4300 44.8000))', 4326),
    ST_GeomFromText('POINT(20.4450 44.7850)', 4326),
    40000,
    14.0,
    ARRAY['11000', '11040']
);

-- Нови Београд - современный район на левом берегу Савы
INSERT INTO districts (name, country_code, boundary, center_point, population, area_km2, postal_codes)
VALUES (
    'Нови Београд',
    'RS',
    ST_GeomFromText('POLYGON((20.3800 44.8200, 20.4300 44.8200, 20.4300 44.7800, 20.3800 44.7800, 20.3800 44.8200))', 4326),
    ST_GeomFromText('POINT(20.4050 44.8000)', 4326),
    220000,
    40.7,
    ARRAY['11070', '11073', '11077']
);

-- Земун - исторический район
INSERT INTO districts (name, country_code, boundary, center_point, population, area_km2, postal_codes)
VALUES (
    'Земун',
    'RS',
    ST_GeomFromText('POLYGON((20.3500 44.8500, 20.4200 44.8500, 20.4200 44.8200, 20.3500 44.8200, 20.3500 44.8500))', 4326),
    ST_GeomFromText('POINT(20.3850 44.8350)', 4326),
    170000,
    154.0,
    ARRAY['11080', '11185']
);

-- Звездара - восточный район
INSERT INTO districts (name, country_code, boundary, center_point, population, area_km2, postal_codes)
VALUES (
    'Звездара',
    'RS',
    ST_GeomFromText('POLYGON((20.4900 44.8100, 20.5400 44.8100, 20.5400 44.7700, 20.4900 44.7700, 20.4900 44.8100))', 4326),
    ST_GeomFromText('POINT(20.5150 44.7900)', 4326),
    150000,
    32.0,
    ARRAY['11050', '11120', '11160']
);

-- Палилула - крупнейший район по площади
INSERT INTO districts (name, country_code, boundary, center_point, population, area_km2, postal_codes)
VALUES (
    'Палилула',
    'RS',
    ST_GeomFromText('POLYGON((20.4650 44.8200, 20.5500 44.8200, 20.5500 44.7800, 20.4650 44.7800, 20.4650 44.8200))', 4326),
    ST_GeomFromText('POINT(20.5075 44.8000)', 4326),
    180000,
    451.0,
    ARRAY['11000', '11060', '11090']
);

-- Раковица - юго-западный район
INSERT INTO districts (name, country_code, boundary, center_point, population, area_km2, postal_codes)
VALUES (
    'Раковица',
    'RS',
    ST_GeomFromText('POLYGON((20.4200 44.7600, 20.4700 44.7600, 20.4700 44.7200, 20.4200 44.7200, 20.4200 44.7600))', 4326),
    ST_GeomFromText('POINT(20.4450 44.7400)', 4326),
    110000,
    31.0,
    ARRAY['11090', '11275']
);

-- Чукарица - западный район
INSERT INTO districts (name, country_code, boundary, center_point, population, area_km2, postal_codes)
VALUES (
    'Чукарица',
    'RS',
    ST_GeomFromText('POLYGON((20.3700 44.7700, 20.4300 44.7700, 20.4300 44.7300, 20.3700 44.7300, 20.3700 44.7700))', 4326),
    ST_GeomFromText('POINT(20.4000 44.7500)', 4326),
    180000,
    156.0,
    ARRAY['11030', '11250', '11253']
);

-- Вождовац - южный район
INSERT INTO districts (name, country_code, boundary, center_point, population, area_km2, postal_codes)
VALUES (
    'Вождовац',
    'RS',
    ST_GeomFromText('POLYGON((20.4600 44.7700, 20.5100 44.7700, 20.5100 44.7300, 20.4600 44.7300, 20.4600 44.7700))', 4326),
    ST_GeomFromText('POINT(20.4850 44.7500)', 4326),
    160000,
    148.0,
    ARRAY['11010', '11040', '11210']
);

-- Гроцка - пригородный район
INSERT INTO districts (name, country_code, boundary, center_point, population, area_km2, postal_codes)
VALUES (
    'Гроцка',
    'RS',
    ST_GeomFromText('POLYGON((20.6000 44.7200, 20.7500 44.7200, 20.7500 44.6200, 20.6000 44.6200, 20.6000 44.7200))', 4326),
    ST_GeomFromText('POINT(20.6750 44.6700)', 4326),
    85000,
    289.0,
    ARRAY['11306', '11309']
);

-- Младеновац - южный пригород
INSERT INTO districts (name, country_code, boundary, center_point, population, area_km2, postal_codes)
VALUES (
    'Младеновац',
    'RS',
    ST_GeomFromText('POLYGON((20.6000 44.5000, 20.7500 44.5000, 20.7500 44.3500, 20.6000 44.3500, 20.6000 44.5000))', 4326),
    ST_GeomFromText('POINT(20.6750 44.4250)', 4326),
    53000,
    339.0,
    ARRAY['11400']
);

-- Сурчин - западный пригород
INSERT INTO districts (name, country_code, boundary, center_point, population, area_km2, postal_codes)
VALUES (
    'Сурчин',
    'RS',
    ST_GeomFromText('POLYGON((20.2000 44.8000, 20.3500 44.8000, 20.3500 44.7000, 20.2000 44.7000, 20.2000 44.8000))', 4326),
    ST_GeomFromText('POINT(20.2750 44.7500)', 4326),
    45000,
    285.0,
    ARRAY['11271', '11272', '11273']
);

-- Барајево - юго-западный пригород
INSERT INTO districts (name, country_code, boundary, center_point, population, area_km2, postal_codes)
VALUES (
    'Барајево',
    'RS',
    ST_GeomFromText('POLYGON((20.3000 44.6500, 20.4500 44.6500, 20.4500 44.5500, 20.3000 44.5500, 20.3000 44.6500))', 4326),
    ST_GeomFromText('POINT(20.3750 44.6000)', 4326),
    27000,
    213.0,
    ARRAY['11460']
);

-- Лазаревац - южный пригород с углем
INSERT INTO districts (name, country_code, boundary, center_point, population, area_km2, postal_codes)
VALUES (
    'Лазаревац',
    'RS',
    ST_GeomFromText('POLYGON((20.3000 44.4500, 20.5000 44.4500, 20.5000 44.3000, 20.3000 44.3000, 20.3000 44.4500))', 4326),
    ST_GeomFromText('POINT(20.4000 44.3750)', 4326),
    58000,
    384.0,
    ARRAY['11550']
);

-- Обреновац - западный промышленный пригород
INSERT INTO districts (name, country_code, boundary, center_point, population, area_km2, postal_codes)
VALUES (
    'Обреновац',
    'RS',
    ST_GeomFromText('POLYGON((20.1000 44.7000, 20.2500 44.7000, 20.2500 44.5500, 20.1000 44.5500, 20.1000 44.7000))', 4326),
    ST_GeomFromText('POINT(20.1750 44.6250)', 4326),
    72000,
    411.0,
    ARRAY['11500']
);

-- Сопот - южный сельский пригород
INSERT INTO districts (name, country_code, boundary, center_point, population, area_km2, postal_codes)
VALUES (
    'Сопот',
    'RS',
    ST_GeomFromText('POLYGON((20.5000 44.6000, 20.6500 44.6000, 20.6500 44.4500, 20.5000 44.4500, 20.5000 44.6000))', 4326),
    ST_GeomFromText('POINT(20.5750 44.5250)', 4326),
    21000,
    271.0,
    ARRAY['11450']
);