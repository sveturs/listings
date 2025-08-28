-- Обновляем опции для атрибута "rooms" (количество комнат)
UPDATE category_attributes 
SET options = '{
  "options": [
    {"value": "studio", "label": "Studio"},
    {"value": "1", "label": "1"},
    {"value": "2", "label": "2"},
    {"value": "3", "label": "3"},
    {"value": "4", "label": "4"},
    {"value": "5", "label": "5"},
    {"value": "6+", "label": "6+"}
  ]
}'::jsonb
WHERE name = 'rooms' AND attribute_type = 'select';