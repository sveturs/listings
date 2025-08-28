-- Откат изменений: отключаем ID атрибуты от категории автомобилей

-- Удаляем связи ID атрибутов с категорией
DELETE FROM category_attribute_mapping cam
USING category_attributes ca
WHERE cam.attribute_id = ca.id 
  AND cam.category_id = 1301
  AND ca.name IN ('car_make_id', 'car_model_id', 'car_generation_id');

-- Возвращаем обязательность для текстовых атрибутов car_make и car_model
UPDATE category_attribute_mapping cam
SET is_required = true
FROM category_attributes ca
WHERE cam.attribute_id = ca.id 
  AND cam.category_id = 1301
  AND ca.name IN ('car_make', 'car_model');

-- Возвращаем data_source для car_generation_id на 'manual'
UPDATE category_attributes
SET data_source = 'manual'
WHERE name = 'car_generation_id';