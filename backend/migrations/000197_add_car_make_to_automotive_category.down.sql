-- Убираем атрибуты car_make, car_make_id, car_model_id из категории 1003
DELETE FROM category_attribute_mapping 
WHERE category_id = 1003 
AND attribute_id IN (
    SELECT id FROM category_attributes 
    WHERE name IN ('car_make', 'car_make_id', 'car_model_id')
);