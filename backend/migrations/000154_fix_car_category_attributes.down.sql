-- Восстанавливаем связь с общим brand для автомобилей
DELETE FROM category_attribute_mapping 
WHERE category_id = 1301 AND attribute_id = 2210;

INSERT INTO category_attribute_mapping (category_id, attribute_id, is_enabled, is_required, sort_order)
VALUES (1301, 2003, true, true, 3);

UPDATE category_attribute_mapping 
SET sort_order = sort_order - 1
WHERE category_id = 1301 AND sort_order > 3;

-- Восстанавливаем для мотоциклов
DELETE FROM category_attribute_mapping 
WHERE category_id = 1302 AND attribute_id = 2211;

INSERT INTO category_attribute_mapping (category_id, attribute_id, is_enabled, is_required, sort_order)
VALUES (1302, 2003, true, true, 3);

UPDATE category_attribute_mapping 
SET sort_order = sort_order - 1
WHERE category_id = 1302 AND sort_order > 3;

-- Восстанавливаем для грузовиков
DELETE FROM category_attribute_mapping 
WHERE category_id = 1303 AND attribute_id = 2212;

INSERT INTO category_attribute_mapping (category_id, attribute_id, is_enabled, is_required, sort_order)
VALUES (1303, 2003, true, true, 3);

UPDATE category_attribute_mapping 
SET sort_order = sort_order - 1
WHERE category_id = 1303 AND sort_order > 3;

-- Восстанавливаем для водного транспорта
DELETE FROM category_attribute_mapping 
WHERE category_id = 1304 AND attribute_id = 2213;

INSERT INTO category_attribute_mapping (category_id, attribute_id, is_enabled, is_required, sort_order)
VALUES (1304, 2003, true, true, 3);

UPDATE category_attribute_mapping 
SET sort_order = sort_order - 1
WHERE category_id = 1304 AND sort_order > 3;

-- Удаляем созданные атрибуты
DELETE FROM category_attributes WHERE id IN (2210, 2211, 2212, 2213);