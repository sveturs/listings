-- Удаляем переводы
DELETE FROM translations 
WHERE entity_type = 'category_attribute' 
AND entity_id IN (
    SELECT id FROM category_attributes WHERE name IN ('car_make_id', 'car_model_id')
);

-- Удаляем связи с категориями
DELETE FROM category_attribute_mapping 
WHERE attribute_id IN (
    SELECT id FROM category_attributes WHERE name IN ('car_make_id', 'car_model_id')
);

-- Удаляем атрибуты
DELETE FROM category_attributes 
WHERE name IN ('car_make_id', 'car_model_id');