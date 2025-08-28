-- Исправление подключения car ID атрибутов к категории автомобилей

-- Удаляем возможные дубликаты (если есть)
DELETE FROM category_attribute_mapping cam
USING category_attributes ca
WHERE cam.attribute_id = ca.id 
  AND cam.category_id = 1301
  AND ca.name IN ('car_make_id', 'car_model_id', 'car_generation_id');

-- Подключаем car_make_id к категории
INSERT INTO category_attribute_mapping (category_id, attribute_id, is_required, sort_order)
SELECT 
    1301 as category_id,
    ca.id as attribute_id,
    true as is_required,
    1 as sort_order
FROM category_attributes ca
WHERE ca.name = 'car_make_id';

-- Подключаем car_model_id к категории
INSERT INTO category_attribute_mapping (category_id, attribute_id, is_required, sort_order)
SELECT 
    1301 as category_id,
    ca.id as attribute_id,
    true as is_required,
    2 as sort_order
FROM category_attributes ca
WHERE ca.name = 'car_model_id';

-- Подключаем car_generation_id к категории (необязательный)
INSERT INTO category_attribute_mapping (category_id, attribute_id, is_required, sort_order)
SELECT 
    1301 as category_id,
    ca.id as attribute_id,
    false as is_required,
    3 as sort_order
FROM category_attributes ca
WHERE ca.name = 'car_generation_id';

-- Делаем текстовые атрибуты car_make и car_model необязательными
UPDATE category_attribute_mapping cam
SET is_required = false
FROM category_attributes ca
WHERE cam.attribute_id = ca.id 
  AND cam.category_id = 1301
  AND ca.name IN ('car_make', 'car_model');

-- Обновляем data_source для car_generation_id на 'api_external'
UPDATE category_attributes
SET data_source = 'api_external'
WHERE name = 'car_generation_id';