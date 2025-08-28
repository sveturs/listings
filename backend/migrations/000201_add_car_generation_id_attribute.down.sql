-- Удаляем атрибут car_generation_id из категории Автомобили
DELETE FROM category_attribute_mapping 
WHERE category_id = 1003 
AND attribute_id = (SELECT id FROM category_attributes WHERE name = 'car_generation_id');

-- Удаляем атрибут car_generation_id
DELETE FROM category_attributes WHERE name = 'car_generation_id';