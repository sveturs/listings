-- Тестовое обновление границ района Земун
BEGIN;

-- Сохраняем текущее состояние
CREATE TEMP TABLE zemun_backup AS 
SELECT id, name, boundary, ST_AsGeoJSON(boundary) as boundary_json, area_km2
FROM districts 
WHERE id = 'a916dcfd-a916-4c43-83e0-e3a3e2c88853';

-- Показываем текущее состояние
SELECT 'BEFORE UPDATE:' as stage, name, 
       ST_NPoints(boundary) as points,
       area_km2,
       ST_AsGeoJSON(boundary)::json->'coordinates' as coordinates
FROM zemun_backup;

-- Обновляем границы
UPDATE districts 
SET boundary = ST_GeomFromGeoJSON('{"type": "Polygon", "coordinates": [[[20.32, 44.86], [20.33, 44.865], [20.34, 44.87], [20.35, 44.875], [20.36, 44.88], [20.37, 44.885], [20.38, 44.887], [20.39, 44.888], [20.4, 44.887], [20.408, 44.885], [20.41, 44.88], [20.409, 44.875], [20.407, 44.87], [20.404, 44.865], [20.4, 44.86], [20.395, 44.855], [20.39, 44.852], [20.385, 44.85], [20.38, 44.849], [20.375, 44.85], [20.37, 44.852], [20.365, 44.854], [20.36, 44.856], [20.355, 44.858], [20.35, 44.859], [20.345, 44.8595], [20.34, 44.8598], [20.335, 44.8599], [20.33, 44.86], [20.325, 44.86], [20.32, 44.86]]]}'),
    area_km2 = ST_Area(ST_GeomFromGeoJSON('{"type": "Polygon", "coordinates": [[[20.32, 44.86], [20.33, 44.865], [20.34, 44.87], [20.35, 44.875], [20.36, 44.88], [20.37, 44.885], [20.38, 44.887], [20.39, 44.888], [20.4, 44.887], [20.408, 44.885], [20.41, 44.88], [20.409, 44.875], [20.407, 44.87], [20.404, 44.865], [20.4, 44.86], [20.395, 44.855], [20.39, 44.852], [20.385, 44.85], [20.38, 44.849], [20.375, 44.85], [20.37, 44.852], [20.365, 44.854], [20.36, 44.856], [20.355, 44.858], [20.35, 44.859], [20.345, 44.8595], [20.34, 44.8598], [20.335, 44.8599], [20.33, 44.86], [20.325, 44.86], [20.32, 44.86]]]}')::geography) / 1000000,
    updated_at = NOW()
WHERE id = 'a916dcfd-a916-4c43-83e0-e3a3e2c88853';

-- Проверяем результат
SELECT 'AFTER UPDATE:' as stage, d.name, 
       ST_IsValid(d.boundary) as is_valid,
       ST_NPoints(d.boundary) as points,
       d.area_km2,
       ST_AsGeoJSON(d.boundary)::json->'type' as geom_type
FROM districts d
WHERE id = 'a916dcfd-a916-4c43-83e0-e3a3e2c88853';

-- Проверяем, что новые границы содержат центральную точку
SELECT 'CONTAINS CENTER:' as check,
       ST_Contains(boundary, center_point) as contains_center
FROM districts 
WHERE id = 'a916dcfd-a916-4c43-83e0-e3a3e2c88853';

-- Сравнение площадей
SELECT 'AREA COMPARISON:' as info,
       z.area_km2 as old_area_km2,
       d.area_km2 as new_area_km2,
       (d.area_km2 - z.area_km2) as difference_km2,
       ((d.area_km2 - z.area_km2) / z.area_km2 * 100)::numeric(5,2) as percent_change
FROM districts d, zemun_backup z
WHERE d.id = 'a916dcfd-a916-4c43-83e0-e3a3e2c88853';

-- Откатываем изменения для теста
ROLLBACK;