-- Core attributes for Serbia marketplace categories

-- Create basic attributes first
INSERT INTO category_attributes (id, name, display_name, attribute_type, is_required, sort_order, validation_rules) VALUES
-- Common attributes
(2001, 'price', 'Cena', 'number', true, 1, '{"min": 0, "max": 999999999}'),
(2002, 'condition', 'Stanje', 'select', false, 2, '{}'),
(2003, 'brand', 'Brend', 'select', false, 3, '{}'),
(2004, 'color', 'Boja', 'select', false, 4, '{}'),

-- Electronics specific
(2101, 'storage', 'Memorija', 'select', false, 5, '{}'),
(2102, 'operating_system', 'Operativni sistem', 'select', false, 6, '{}'),
(2103, 'processor', 'Procesor', 'text', false, 7, '{}'),
(2104, 'ram', 'RAM memorija', 'select', false, 8, '{}'),
(2105, 'storage_type', 'Tip skladišta', 'select', false, 9, '{}'),

-- Automotive specific  
(2201, 'car_model', 'Model', 'text', true, 10, '{}'),
(2202, 'year', 'Godište', 'number', true, 11, '{"min": 1950, "max": 2025}'),
(2203, 'mileage', 'Kilometraža', 'number', false, 12, '{"min": 0, "max": 999999}'),
(2204, 'fuel_type', 'Gorivo', 'select', true, 13, '{}'),
(2205, 'transmission', 'Menjač', 'select', true, 14, '{}'),
(2206, 'body_type', 'Tip karoserije', 'select', false, 15, '{}'),

-- Real Estate specific
(2301, 'area', 'Površina (m²)', 'number', true, 16, '{"min": 10, "max": 1000}'),
(2302, 'rooms', 'Broj soba', 'select', true, 17, '{}'),
(2303, 'floor', 'Sprat', 'number', false, 18, '{"min": -2, "max": 50}'),
(2304, 'furnished', 'Namešteno', 'boolean', false, 19, '{}'),
(2305, 'parking', 'Parking', 'boolean', false, 20, '{}'),
(2306, 'balcony', 'Balkon/terasa', 'boolean', false, 21, '{}'),

-- House specific
(2307, 'house_area', 'Površina kuće (m²)', 'number', true, 22, '{"min": 50, "max": 2000}'),
(2308, 'land_area', 'Površina placa (m²)', 'number', false, 23, '{"min": 100, "max": 50000}'),
(2309, 'bathrooms', 'Broj kupatila', 'number', false, 24, '{"min": 1, "max": 10}'),
(2310, 'garden', 'Bašta', 'boolean', false, 25, '{}'),
(2311, 'garage', 'Garaža', 'boolean', false, 26, '{}'),

-- Agriculture specific
(2401, 'working_hours', 'Radnih sati', 'number', false, 27, '{"min": 0, "max": 50000}'),
(2402, 'power_hp', 'Snaga (KS)', 'number', false, 28, '{"min": 10, "max": 1000}'),

-- Services specific
(2501, 'service_type', 'Tip usluge', 'multiselect', true, 29, '{}'),
(2502, 'experience_years', 'Godine iskustva', 'number', false, 30, '{"min": 0, "max": 50}'),
(2503, 'availability', 'Dostupnost', 'select', false, 31, '{}'),
(2504, 'service_area', 'Oblast rada', 'multiselect', false, 32, '{}');

-- Map attributes to categories
-- Electronics categories (price, condition, brand)
INSERT INTO category_attribute_mapping (category_id, attribute_id, is_required, sort_order) VALUES
-- Main Electronics category
(1001, 2001, true, 1),  -- price
(1001, 2002, false, 2), -- condition
(1001, 2003, false, 3), -- brand

-- Smartphones
(1101, 2001, true, 1),   -- price
(1101, 2002, false, 2),  -- condition  
(1101, 2003, true, 3),   -- brand
(1101, 2004, false, 4),  -- color
(1101, 2101, false, 5),  -- storage
(1101, 2102, false, 6),  -- operating_system

-- Computers
(1102, 2001, true, 1),   -- price
(1102, 2002, false, 2),  -- condition
(1102, 2003, true, 3),   -- brand
(1102, 2103, false, 4),  -- processor
(1102, 2104, false, 5),  -- ram
(1102, 2105, false, 6),  -- storage_type

-- Automotive categories
-- Cars
(1301, 2001, true, 1),   -- price
(1301, 2002, true, 2),   -- condition
(1301, 2003, true, 3),   -- brand
(1301, 2201, true, 4),   -- car_model
(1301, 2202, true, 5),   -- year
(1301, 2203, false, 6),  -- mileage
(1301, 2204, true, 7),   -- fuel_type
(1301, 2205, true, 8),   -- transmission
(1301, 2206, false, 9),  -- body_type

-- Real Estate categories
-- Apartments
(1401, 2001, true, 1),   -- price
(1401, 2301, true, 2),   -- area
(1401, 2302, true, 3),   -- rooms
(1401, 2303, false, 4),  -- floor
(1401, 2304, false, 5),  -- furnished
(1401, 2305, false, 6),  -- parking
(1401, 2306, false, 7),  -- balcony

-- Houses
(1402, 2001, true, 1),   -- price
(1402, 2307, true, 2),   -- house_area
(1402, 2308, false, 3),  -- land_area
(1402, 2302, true, 4),   -- rooms
(1402, 2309, false, 5),  -- bathrooms
(1402, 2310, false, 6),  -- garden
(1402, 2311, false, 7),  -- garage

-- Agriculture
-- Farm machinery
(1601, 2001, true, 1),   -- price
(1601, 2002, false, 2),  -- condition
(1601, 2003, true, 3),   -- brand
(1601, 2202, false, 4),  -- year
(1601, 2401, false, 5),  -- working_hours
(1601, 2402, false, 6),  -- power_hp

-- Services
-- Construction services  
(1901, 2001, true, 1),   -- price
(1901, 2501, true, 2),   -- service_type
(1901, 2502, false, 3),  -- experience_years
(1901, 2503, false, 4),  -- availability
(1901, 2504, false, 5);  -- service_area

-- Create attribute options
INSERT INTO attribute_option_translations (id, attribute_name, option_value, sr_translation, ru_translation) VALUES
-- Condition options
(3001, 'condition', 'new', 'Novo', 'Новое'),
(3002, 'condition', 'used', 'Korišćeno', 'Б/у'),
(3003, 'condition', 'refurbished', 'Obnovljeno', 'Восстановленное'),

-- Phone brands
(3101, 'brand', 'apple', 'Apple', 'Apple'),
(3102, 'brand', 'samsung', 'Samsung', 'Samsung'),
(3103, 'brand', 'huawei', 'Huawei', 'Huawei'),
(3104, 'brand', 'xiaomi', 'Xiaomi', 'Xiaomi'),

-- Storage options
(3201, 'storage', '64gb', '64GB', '64ГБ'),
(3202, 'storage', '128gb', '128GB', '128ГБ'),
(3203, 'storage', '256gb', '256GB', '256ГБ'),
(3204, 'storage', '512gb', '512GB', '512ГБ'),

-- Fuel type options
(3301, 'fuel_type', 'gasoline', 'Benzin', 'Бензин'),
(3302, 'fuel_type', 'diesel', 'Dizel', 'Дизель'),
(3303, 'fuel_type', 'hybrid', 'Hibrid', 'Гибрид'),
(3304, 'fuel_type', 'electric', 'Električni', 'Электрический'),

-- Transmission options
(3401, 'transmission', 'manual', 'Manuelni', 'Механика'),
(3402, 'transmission', 'automatic', 'Automatski', 'Автомат'),

-- Room count options
(3501, 'rooms', '1', '1 soba', '1 комната'),
(3502, 'rooms', '2', '2 sobe', '2 комнаты'),
(3503, 'rooms', '3', '3 sobe', '3 комнаты'),
(3504, 'rooms', '4', '4 sobe', '4 комнаты'),
(3505, 'rooms', '5+', '5+ soba', '5+ комнат');

-- Translations for attributes
INSERT INTO translations (entity_type, entity_id, field_name, language, translated_text, is_machine_translated, is_verified) VALUES
-- Price translations
('attribute', 2001, 'display_name', 'sr', 'Цена', false, true),
('attribute', 2001, 'display_name', 'ru', 'Цена', true, false),
('attribute', 2001, 'display_name', 'en', 'Price', true, false),

-- Condition translations
('attribute', 2002, 'display_name', 'sr', 'Стање', false, true),
('attribute', 2002, 'display_name', 'ru', 'Состояние', true, false),
('attribute', 2002, 'display_name', 'en', 'Condition', true, false),

-- Brand translations
('attribute', 2003, 'display_name', 'sr', 'Бренд', false, true),
('attribute', 2003, 'display_name', 'ru', 'Бренд', true, false),
('attribute', 2003, 'display_name', 'en', 'Brand', true, false),

-- Area translations
('attribute', 2301, 'display_name', 'sr', 'Површина (м²)', false, true),
('attribute', 2301, 'display_name', 'ru', 'Площадь (м²)', true, false),
('attribute', 2301, 'display_name', 'en', 'Area (m²)', true, false);

-- Reset sequences
SELECT setval('category_attributes_id_seq', 3000, true);