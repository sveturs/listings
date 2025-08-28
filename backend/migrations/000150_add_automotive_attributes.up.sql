-- Обновляем значения для атрибута brand (2003) для автомобилей
UPDATE category_attributes 
SET options = '{"values": ["Volkswagen", "BMW", "Mercedes", "Audi", "Opel", "Peugeot", "Renault", "Fiat", "Ford", "Toyota", "Hyundai", "Kia", "Nissan", "Honda", "Mazda", "Skoda", "Seat", "Citroen", "Volvo", "Other"]}'::jsonb
WHERE id = 2003;

-- Обновляем значения для атрибута color (2004) - добавляем grey и brown
UPDATE category_attributes 
SET options = '{"values": ["black", "white", "silver", "grey", "blue", "red", "green", "yellow", "brown", "purple", "other"]}'::jsonb
WHERE id = 2004;

-- Добавляем атрибут model (которого еще нет)
INSERT INTO category_attributes (id, name, display_name, attribute_type, is_required) VALUES
(3105, 'model', 'Model', 'text', false);

-- Привязываем атрибуты к категории automotive (1003)
INSERT INTO category_attribute_mapping (category_id, attribute_id) VALUES
(1003, 3105), -- model (новый)
(1003, 2202), -- year (существующий)
(1003, 2203), -- mileage (существующий)
(1003, 2206), -- body_type (существующий)
(1003, 3001); -- engine_size (существующий)