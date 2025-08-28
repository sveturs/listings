-- Удаляем атрибут power_hp из категории Автомобили (ID: 1003)
DELETE FROM category_attribute_mapping 
WHERE category_id = 1003 
AND attribute_id = (SELECT id FROM category_attributes WHERE name = 'power_hp');