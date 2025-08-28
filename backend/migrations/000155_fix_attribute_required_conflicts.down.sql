-- Откат изменений is_required к предыдущим значениям

-- 1. Возвращаем condition как обязательный для автомобилей
UPDATE category_attribute_mapping
SET is_required = true
WHERE attribute_id = (SELECT id FROM category_attributes WHERE name = 'condition')
AND category_id = 1301;

-- 2. Возвращаем brand как обязательный где был
UPDATE category_attribute_mapping
SET is_required = true
WHERE attribute_id IN (
    SELECT id FROM category_attributes WHERE name = 'brand'
)
AND category_id IN (1101, 1102, 1601);

UPDATE category_attribute_mapping
SET is_required = true
WHERE attribute_id = (SELECT id FROM category_attributes WHERE name = 'car_make')
AND category_id = 1301;

UPDATE category_attribute_mapping
SET is_required = true
WHERE attribute_id = (SELECT id FROM category_attributes WHERE name = 'motorcycle_make')
AND category_id = 1302;

UPDATE category_attribute_mapping
SET is_required = true
WHERE attribute_id = (SELECT id FROM category_attributes WHERE name = 'truck_make')
AND category_id = 1303;

UPDATE category_attribute_mapping
SET is_required = true
WHERE attribute_id = (SELECT id FROM category_attributes WHERE name = 'boat_make')
AND category_id = 1304;

-- 3. Возвращаем fuel_type и transmission как необязательные
UPDATE category_attribute_mapping
SET is_required = false
WHERE attribute_id IN (
    SELECT id FROM category_attributes 
    WHERE name IN ('fuel_type', 'transmission')
)
AND category_id IN (1302, 1303, 1304);

-- 4. Возвращаем year как необязательный для сельхозтехники
UPDATE category_attribute_mapping
SET is_required = false
WHERE attribute_id = (SELECT id FROM category_attributes WHERE name = 'year')
AND category_id = 1601;