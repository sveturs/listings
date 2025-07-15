-- SQL скрипт для обновления границ районов из GeoJSON файлов
-- Сгенерировано: 2025-07-15T14:00:53.687616

BEGIN;

-- Создаем резервную копию текущих границ
CREATE TABLE IF NOT EXISTS districts_boundaries_backup_20250715 AS 
SELECT id, name, boundary, updated_at FROM districts;

-- Обновления для города Београд

-- Район: Стари Град (источник: Стари Град)
UPDATE districts 
SET boundary = ST_GeomFromGeoJSON('{"type": "Polygon", "coordinates": [[[20.45, 44.82], [20.455, 44.825], [20.465, 44.82], [20.47, 44.815], [20.465, 44.81], [20.455, 44.805], [20.45, 44.81], [20.445, 44.815], [20.45, 44.82]]]}'),
    area_km2 = ST_Area(ST_GeomFromGeoJSON('{"type": "Polygon", "coordinates": [[[20.45, 44.82], [20.455, 44.825], [20.465, 44.82], [20.47, 44.815], [20.465, 44.81], [20.455, 44.805], [20.45, 44.81], [20.445, 44.815], [20.45, 44.82]]]}')::geography) / 1000000,
    updated_at = NOW()
WHERE id = '2f2c2f9a-0ecd-4df5-a5a4-ba626288a9d0';


-- Район: Врачар (источник: Врачар)
UPDATE districts 
SET boundary = ST_GeomFromGeoJSON('{"type": "Polygon", "coordinates": [[[20.465, 44.81], [20.47, 44.815], [20.48, 44.81], [20.49, 44.805], [20.49, 44.795], [20.48, 44.79], [20.47, 44.795], [20.465, 44.8], [20.465, 44.81]]]}'),
    area_km2 = ST_Area(ST_GeomFromGeoJSON('{"type": "Polygon", "coordinates": [[[20.465, 44.81], [20.47, 44.815], [20.48, 44.81], [20.49, 44.805], [20.49, 44.795], [20.48, 44.79], [20.47, 44.795], [20.465, 44.8], [20.465, 44.81]]]}')::geography) / 1000000,
    updated_at = NOW()
WHERE id = 'dce8dbfc-ef0d-492e-8d9c-ffb28e4e41c1';


-- Всего обновлений: 2

-- Проверка валидности геометрии
SELECT c.name as city, d.name as district, 
       ST_IsValid(d.boundary) as is_valid,
       ST_Area(d.boundary::geography) / 1000000 as area_km2,
       ST_NPoints(d.boundary) as num_points
FROM districts d
JOIN cities c ON d.city_id = c.id
WHERE d.updated_at > NOW() - INTERVAL '1 minute'
ORDER BY c.name, d.name;

-- Исправление невалидных геометрий
UPDATE districts
SET boundary = ST_MakeValid(boundary)
WHERE NOT ST_IsValid(boundary) AND boundary IS NOT NULL;

COMMIT;