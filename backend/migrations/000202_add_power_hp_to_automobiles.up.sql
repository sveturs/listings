-- Добавляем атрибут power_hp к категории Автомобили (ID: 1003)
INSERT INTO category_attribute_mapping (category_id, attribute_id)
SELECT 1003, id FROM category_attributes WHERE name = 'power_hp'
ON CONFLICT (category_id, attribute_id) DO NOTHING;