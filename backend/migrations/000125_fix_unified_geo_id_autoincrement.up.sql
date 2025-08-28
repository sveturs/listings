-- Добавляем автоинкремент для поля id в таблице unified_geo

-- Создаем последовательность
CREATE SEQUENCE IF NOT EXISTS unified_geo_id_seq;

-- Устанавливаем владельца последовательности
ALTER SEQUENCE unified_geo_id_seq OWNED BY unified_geo.id;

-- Устанавливаем значение по умолчанию для колонки id
ALTER TABLE unified_geo 
    ALTER COLUMN id SET DEFAULT nextval('unified_geo_id_seq');

-- Устанавливаем текущее значение последовательности на максимальный существующий id + 1
SELECT setval('unified_geo_id_seq', COALESCE((SELECT MAX(id) FROM unified_geo), 0) + 1, false);