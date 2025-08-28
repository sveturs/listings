-- Возвращаем пустые опции для атрибута "rooms"
UPDATE category_attributes 
SET options = '{}'::jsonb
WHERE name = 'rooms' AND attribute_type = 'select';