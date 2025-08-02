-- Добавляем модели для сербских марок автомобилей

-- Получаем ID марок
DO $$
DECLARE
    zastava_id INTEGER;
    yugo_id INTEGER;
    fap_id INTEGER;
    imt_id INTEGER;
BEGIN
    -- Получаем ID марок
    SELECT id INTO zastava_id FROM car_makes WHERE name = 'Zastava';
    SELECT id INTO yugo_id FROM car_makes WHERE name = 'Yugo';
    SELECT id INTO fap_id FROM car_makes WHERE name = 'FAP';
    SELECT id INTO imt_id FROM car_makes WHERE name = 'IMT';

    -- Добавляем модели для Zastava (если их еще нет)
    INSERT INTO car_models (make_id, name, slug, production_start, production_end, is_active)
    SELECT zastava_id, name, slug, prod_start, prod_end, true
    FROM (VALUES
        ('750 (Fića)', '750-fica', 1955, 1985),
        ('600 (Fićo)', '600-fico', 1960, 1986),
        ('101 (Stojadin)', '101-stojadin', 1971, 2008),
        ('128 (Skala)', '128-skala', 1988, 2008),
        ('Florida', 'florida', 1988, 2008),
        ('Koral (Yugo Koral)', 'koral', 1980, 2008),
        ('10', '10', 2005, 2008),
        ('1100', '1100', 1961, 1969),
        ('1300', '1300', 1961, 1979),
        ('Campagnola AR', 'campagnola-ar', 1953, 1962)
    ) AS models(name, slug, prod_start, prod_end)
    WHERE zastava_id IS NOT NULL
    AND NOT EXISTS (
        SELECT 1 FROM car_models 
        WHERE make_id = zastava_id 
        AND LOWER(name) = LOWER(models.name)
    );

    -- Добавляем модели для Yugo (если их еще нет)
    INSERT INTO car_models (make_id, name, slug, production_start, production_end, is_active)
    SELECT yugo_id, name, slug, prod_start, prod_end, true
    FROM (VALUES
        ('45', '45', 1980, 1992),
        ('55', '55', 1983, 1992),
        ('60', '60', 1985, 1991),
        ('65', '65', 1985, 1991),
        ('GV', 'gv', 1985, 1992),
        ('GVX', 'gvx', 1987, 1991),
        ('GVL', 'gvl', 1988, 1990),
        ('GV Plus', 'gv-plus', 1988, 1991),
        ('Tempo', 'tempo', 1988, 1991),
        ('Cabrio', 'cabrio', 1990, 1991)
    ) AS models(name, slug, prod_start, prod_end)
    WHERE yugo_id IS NOT NULL
    AND NOT EXISTS (
        SELECT 1 FROM car_models 
        WHERE make_id = yugo_id 
        AND LOWER(name) = LOWER(models.name)
    );

    -- Добавляем модели для FAP (грузовики)
    INSERT INTO car_models (make_id, name, slug, is_active)
    SELECT fap_id, name, slug, true
    FROM (VALUES
        ('1314', '1314'),
        ('1414', '1414'),
        ('1620', '1620'),
        ('1921', '1921'),
        ('2023', '2023'),
        ('2025', '2025'),
        ('2628', '2628'),
        ('2632', '2632'),
        ('3232', '3232'),
        ('2228', '2228')
    ) AS models(name, slug)
    WHERE fap_id IS NOT NULL
    AND NOT EXISTS (
        SELECT 1 FROM car_models 
        WHERE make_id = fap_id 
        AND LOWER(name) = LOWER(models.name)
    );

    -- Добавляем модели для IMT (тракторы)
    INSERT INTO car_models (make_id, name, slug, is_active)
    SELECT imt_id, name, slug, true
    FROM (VALUES
        ('533', '533'),
        ('539', '539'),
        ('542', '542'),
        ('549', '549'),
        ('558', '558'),
        ('560', '560'),
        ('565', '565'),
        ('577', '577'),
        ('579', '579'),
        ('590', '590')
    ) AS models(name, slug)
    WHERE imt_id IS NOT NULL
    AND NOT EXISTS (
        SELECT 1 FROM car_models 
        WHERE make_id = imt_id 
        AND LOWER(name) = LOWER(models.name)
    );

    RAISE NOTICE 'Serbian car models added successfully';
END $$;

-- Обновляем счетчик популярности для сербских марок
UPDATE car_makes 
SET popularity_rs = CASE 
    WHEN name = 'Zastava' THEN 100
    WHEN name = 'Yugo' THEN 90
    WHEN name = 'FAP' THEN 70
    WHEN name = 'IMT' THEN 60
    WHEN name = 'IMK' THEN 50
    WHEN name = 'IDA-Opel' THEN 40
    ELSE popularity_rs
END
WHERE is_domestic = true;

-- Проверяем результат
SELECT m.name as make, COUNT(md.id) as models
FROM car_makes m
LEFT JOIN car_models md ON md.make_id = m.id
WHERE m.is_domestic = true
GROUP BY m.id, m.name
ORDER BY m.popularity_rs DESC;