-- Обновляем опции для атрибута "rooms" используя формат с "values"
UPDATE category_attributes 
SET options = '{
  "values": ["studio", "1", "2", "3", "4", "5", "6+"]
}'::jsonb
WHERE name = 'rooms' AND attribute_type = 'select';