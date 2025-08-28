-- Add Serbian translations for tire attributes

INSERT INTO translations (entity_type, entity_id, field_name, language, translated_text) VALUES
-- Serbian translations
('attribute', (SELECT id FROM category_attributes WHERE name = 'tire_width'), 'display_name', 'sr', 'Ширина гуме'),
('attribute', (SELECT id FROM category_attributes WHERE name = 'tire_profile'), 'display_name', 'sr', 'Профил гуме'),
('attribute', (SELECT id FROM category_attributes WHERE name = 'tire_diameter'), 'display_name', 'sr', 'Пречник'),
('attribute', (SELECT id FROM category_attributes WHERE name = 'tire_season'), 'display_name', 'sr', 'Сезона'),
('attribute', (SELECT id FROM category_attributes WHERE name = 'tire_brand'), 'display_name', 'sr', 'Произвођач'),
('attribute', (SELECT id FROM category_attributes WHERE name = 'tire_condition'), 'display_name', 'sr', 'Стање'),
('attribute', (SELECT id FROM category_attributes WHERE name = 'tread_depth'), 'display_name', 'sr', 'Дубина шаре'),
('attribute', (SELECT id FROM category_attributes WHERE name = 'tire_year'), 'display_name', 'sr', 'Година производње'),
('attribute', (SELECT id FROM category_attributes WHERE name = 'tire_quantity'), 'display_name', 'sr', 'Количина');

-- Also ensure attribute option translations for Serbian are complete
-- Check if these already exist before inserting
INSERT INTO attribute_option_translations (attribute_name, option_value, ru_translation, sr_translation) VALUES
-- Season translations
('tire_season', 'summer', 'Летние', 'Летње'),
('tire_season', 'winter', 'Зимние', 'Зимске'),
('tire_season', 'all-season', 'Всесезонные', 'Целогодишње'),

-- Condition translations
('tire_condition', 'new', 'Новые', 'Нове'),
('tire_condition', 'used', 'Б/У', 'Половне'),

-- Quantity translations
('tire_quantity', '1', '1 шт', '1 ком'),
('tire_quantity', '2', '2 шт', '2 ком'),
('tire_quantity', '3', '3 шт', '3 ком'),
('tire_quantity', '4', '4 шт', '4 ком'),
('tire_quantity', 'set', 'Комплект (4 шт)', 'Комплет (4 ком)')
ON CONFLICT (attribute_name, option_value) DO UPDATE 
SET 
    ru_translation = EXCLUDED.ru_translation,
    sr_translation = EXCLUDED.sr_translation;