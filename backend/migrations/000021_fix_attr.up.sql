-- Обновляем существующие данные для различных числовых атрибутов
UPDATE listing_attribute_values
SET unit = 'm²'
WHERE attribute_id IN (SELECT id FROM category_attributes WHERE name = 'area')
AND unit IS NULL;

UPDATE listing_attribute_values
SET unit = 'ar'
WHERE attribute_id IN (SELECT id FROM category_attributes WHERE name = 'land_area')
AND unit IS NULL;

UPDATE listing_attribute_values
SET unit = 'km'
WHERE attribute_id IN (SELECT id FROM category_attributes WHERE name = 'mileage')
AND unit IS NULL;

UPDATE listing_attribute_values
SET unit = 'l'
WHERE attribute_id IN (SELECT id FROM category_attributes WHERE name = 'engine_capacity')
AND unit IS NULL;

UPDATE listing_attribute_values
SET unit = 'ks'
WHERE attribute_id IN (SELECT id FROM category_attributes WHERE name = 'power')
AND unit IS NULL;

UPDATE listing_attribute_values
SET unit = 'inč'
WHERE attribute_id IN (SELECT id FROM category_attributes WHERE name = 'screen_size')
AND unit IS NULL;

UPDATE listing_attribute_values
SET unit = 'soba'
WHERE attribute_id IN (SELECT id FROM category_attributes WHERE name = 'rooms')
AND unit IS NULL;

-- Обновляем информацию о единицах измерения в category_attributes
UPDATE category_attributes
SET options = jsonb_build_object('type', 'number', 'unit', 'm²')
WHERE name = 'area' AND (options IS NULL OR options = '{}' OR options = 'null');

UPDATE category_attributes
SET options = jsonb_build_object('type', 'number', 'unit', 'ar')
WHERE name = 'land_area' AND (options IS NULL OR options = '{}' OR options = 'null');

UPDATE category_attributes
SET options = jsonb_build_object('type', 'number', 'unit', 'km')
WHERE name = 'mileage' AND (options IS NULL OR options = '{}' OR options = 'null');

UPDATE category_attributes
SET options = jsonb_build_object('type', 'number', 'unit', 'l')
WHERE name = 'engine_capacity' AND (options IS NULL OR options = '{}' OR options = 'null');

UPDATE category_attributes
SET options = jsonb_build_object('type', 'number', 'unit', 'ks')
WHERE name = 'power' AND (options IS NULL OR options = '{}' OR options = 'null');

UPDATE category_attributes
SET options = jsonb_build_object('type', 'number', 'unit', 'inč')
WHERE name = 'screen_size' AND (options IS NULL OR options = '{}' OR options = 'null');

UPDATE category_attributes
SET options = jsonb_build_object('type', 'number', 'unit', 'soba')
WHERE name = 'rooms' AND (options IS NULL OR options = '{}' OR options = 'null');