-- Add attributes for tire categories

-- First, create tire-specific attributes in category_attributes
INSERT INTO category_attributes (name, display_name, attribute_type, is_required, options, sort_order) VALUES
-- Размер шин
('tire_width', 'Širina gume', 'text', true, NULL, 1),
('tire_profile', 'Profil gume', 'text', true, NULL, 2),
('tire_diameter', 'Prečnik', 'text', true, NULL, 3),
-- Тип шин
('tire_season', 'Sezona', 'select', true, '{"values": ["summer", "winter", "all-season"]}', 4),
-- Бренд
('tire_brand', 'Proizvođač', 'text', false, NULL, 5),
-- Состояние
('tire_condition', 'Stanje', 'select', true, '{"values": ["new", "used"]}', 6),
-- Глубина протектора
('tread_depth', 'Dubina šare', 'text', false, NULL, 7),
-- Год производства
('tire_year', 'Godina proizvodnje', 'number', false, NULL, 8),
-- Количество
('tire_quantity', 'Količina', 'select', true, '{"values": ["1", "2", "3", "4", "set"]}', 9);

-- Map attributes to tire category (1304 - Gume i točkovi)
INSERT INTO category_attribute_mapping (category_id, attribute_id, is_required, sort_order)
SELECT 1304, id, is_required, sort_order 
FROM category_attributes 
WHERE name IN (
    'tire_width', 'tire_profile', 'tire_diameter', 'tire_season',
    'tire_brand', 'tire_condition', 'tread_depth', 'tire_year',
    'tire_quantity'
);

-- Also map to subcategories
-- Summer tires (1314)
INSERT INTO category_attribute_mapping (category_id, attribute_id, is_required, sort_order)
SELECT 1314, id, is_required, sort_order 
FROM category_attributes 
WHERE name IN (
    'tire_width', 'tire_profile', 'tire_diameter', 'tire_season',
    'tire_brand', 'tire_condition', 'tread_depth', 'tire_year',
    'tire_quantity'
);

-- Winter tires (1315)
INSERT INTO category_attribute_mapping (category_id, attribute_id, is_required, sort_order)
SELECT 1315, id, is_required, sort_order 
FROM category_attributes 
WHERE name IN (
    'tire_width', 'tire_profile', 'tire_diameter', 'tire_season',
    'tire_brand', 'tire_condition', 'tread_depth', 'tire_year',
    'tire_quantity'
);

-- All-season tires (1316)
INSERT INTO category_attribute_mapping (category_id, attribute_id, is_required, sort_order)
SELECT 1316, id, is_required, sort_order 
FROM category_attributes 
WHERE name IN (
    'tire_width', 'tire_profile', 'tire_diameter', 'tire_season',
    'tire_brand', 'tire_condition', 'tread_depth', 'tire_year',
    'tire_quantity'
);

-- Add translations for new attributes
INSERT INTO translations (entity_type, entity_id, field_name, language, translated_text) VALUES
-- English translations
('attribute', (SELECT id FROM category_attributes WHERE name = 'tire_width'), 'display_name', 'en', 'Tire Width'),
('attribute', (SELECT id FROM category_attributes WHERE name = 'tire_profile'), 'display_name', 'en', 'Tire Profile'),
('attribute', (SELECT id FROM category_attributes WHERE name = 'tire_diameter'), 'display_name', 'en', 'Diameter'),
('attribute', (SELECT id FROM category_attributes WHERE name = 'tire_season'), 'display_name', 'en', 'Season'),
('attribute', (SELECT id FROM category_attributes WHERE name = 'tire_brand'), 'display_name', 'en', 'Brand'),
('attribute', (SELECT id FROM category_attributes WHERE name = 'tire_condition'), 'display_name', 'en', 'Condition'),
('attribute', (SELECT id FROM category_attributes WHERE name = 'tread_depth'), 'display_name', 'en', 'Tread Depth'),
('attribute', (SELECT id FROM category_attributes WHERE name = 'tire_year'), 'display_name', 'en', 'Production Year'),
('attribute', (SELECT id FROM category_attributes WHERE name = 'tire_quantity'), 'display_name', 'en', 'Quantity'),

-- Russian translations
('attribute', (SELECT id FROM category_attributes WHERE name = 'tire_width'), 'display_name', 'ru', 'Ширина шины'),
('attribute', (SELECT id FROM category_attributes WHERE name = 'tire_profile'), 'display_name', 'ru', 'Профиль шины'),
('attribute', (SELECT id FROM category_attributes WHERE name = 'tire_diameter'), 'display_name', 'ru', 'Диаметр'),
('attribute', (SELECT id FROM category_attributes WHERE name = 'tire_season'), 'display_name', 'ru', 'Сезон'),
('attribute', (SELECT id FROM category_attributes WHERE name = 'tire_brand'), 'display_name', 'ru', 'Производитель'),
('attribute', (SELECT id FROM category_attributes WHERE name = 'tire_condition'), 'display_name', 'ru', 'Состояние'),
('attribute', (SELECT id FROM category_attributes WHERE name = 'tread_depth'), 'display_name', 'ru', 'Глубина протектора'),
('attribute', (SELECT id FROM category_attributes WHERE name = 'tire_year'), 'display_name', 'ru', 'Год производства'),
('attribute', (SELECT id FROM category_attributes WHERE name = 'tire_quantity'), 'display_name', 'ru', 'Количество');

-- Add attribute option translations using the existing structure
INSERT INTO attribute_option_translations (attribute_name, option_value, ru_translation, sr_translation) VALUES
-- Season translations
('tire_season', 'summer', 'Летние', 'Letnje'),
('tire_season', 'winter', 'Зимние', 'Zimske'),
('tire_season', 'all-season', 'Всесезонные', 'Celogodišnje'),

-- Condition translations
('tire_condition', 'new', 'Новые', 'Nove'),
('tire_condition', 'used', 'Б/У', 'Polovne'),

-- Quantity translations
('tire_quantity', '1', '1 шт', '1 kom'),
('tire_quantity', '2', '2 шт', '2 kom'),
('tire_quantity', '3', '3 шт', '3 kom'),
('tire_quantity', '4', '4 шт', '4 kom'),
('tire_quantity', 'set', 'Комплект (4 шт)', 'Komplet (4 kom)');