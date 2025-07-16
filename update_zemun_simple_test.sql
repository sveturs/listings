-- Тестовое обновление границ района Земун с простым прямоугольником для проверки
-- Это поможет убедиться, что система отображения работает правильно
UPDATE districts 
SET boundary = ST_GeomFromGeoJSON('{
  "type": "Polygon",
  "coordinates": [[
    [20.32, 44.86],
    [20.41, 44.86],
    [20.41, 44.91],
    [20.32, 44.91],
    [20.32, 44.86]
  ]]
}')
WHERE id = 'a916dcfd-a916-4c43-83e0-e3a3e2c88853';

-- Проверяем результат
SELECT id, name, ST_AsGeoJSON(boundary) as boundary_geojson 
FROM districts 
WHERE id = 'a916dcfd-a916-4c43-83e0-e3a3e2c88853';
