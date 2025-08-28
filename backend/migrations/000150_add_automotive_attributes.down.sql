-- Удаляем привязки к категории
DELETE FROM category_attribute_mapping WHERE category_id = 1003 AND attribute_id IN (3105, 2202, 2203, 2206, 3001);

-- Удаляем атрибут model
DELETE FROM category_attributes WHERE id = 3105;

-- Возвращаем старые значения для brand
UPDATE category_attributes 
SET options = '{"values": ["Apple", "Samsung", "Sony", "LG", "Dell", "HP", "Lenovo", "Asus", "Acer", "Other"]}'::jsonb
WHERE id = 2003;