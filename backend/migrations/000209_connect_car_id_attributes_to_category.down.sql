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

-- Очищаем конфигурацию data_source для ID атрибутов
UPDATE category_attributes
SET 
    data_source = NULL,
    data_source_config = NULL
WHERE name IN ('car_make_id', 'car_model_id', 'car_generation_id');

-- Удаляем переводы для ID атрибутов
DELETE FROM translations
WHERE key IN (
    'attribute.car_make_id',
    'attribute.car_model_id', 
    'attribute.car_generation_id'
);