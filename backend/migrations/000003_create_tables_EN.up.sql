-- Английские переводы основных категорий
INSERT INTO translations (entity_type, entity_id, language, field_name, translated_text, is_machine_translated, is_verified, created_at, updated_at) VALUES
-- Основные категории
('category', 1, 'en', 'name', 'Real Estate', true, true, NOW(), NOW()),
('category', 2, 'en', 'name', 'Transport', true, true, NOW(), NOW()),
('category', 3, 'en', 'name', 'Electronics', true, true, NOW(), NOW()),
('category', 4, 'en', 'name', 'Home & Apartment', true, true, NOW(), NOW()),
('category', 5, 'en', 'name', 'Garden', true, true, NOW(), NOW()),
('category', 6, 'en', 'name', 'Hobbies & Leisure', true, true, NOW(), NOW()),
('category', 7, 'en', 'name', 'Animals', true, true, NOW(), NOW()),
('category', 8, 'en', 'name', 'Business & Equipment', true, true, NOW(), NOW()),
('category', 9, 'en', 'name', 'Other', true, true, NOW(), NOW()),
('category', 10, 'en', 'name', 'Jobs', true, true, NOW(), NOW()),
('category', 11, 'en', 'name', 'Clothing & Accessories', true, true, NOW(), NOW()),
('category', 12, 'en', 'name', 'Kids Clothing & Accessories', true, true, NOW(), NOW()),

-- Подкатегории Недвижимости
('category', 13, 'en', 'name', 'Apartment', true, true, NOW(), NOW()),
('category', 14, 'en', 'name', 'Room', true, true, NOW(), NOW()),
('category', 15, 'en', 'name', 'House, Cottage', true, true, NOW(), NOW()),
('category', 16, 'en', 'name', 'Land Plot', true, true, NOW(), NOW()),
('category', 17, 'en', 'name', 'Garage & Parking', true, true, NOW(), NOW()),
('category', 18, 'en', 'name', 'Commercial Property', true, true, NOW(), NOW()),
('category', 19, 'en', 'name', 'Foreign Property', true, true, NOW(), NOW()),
('category', 20, 'en', 'name', 'Hotel', true, true, NOW(), NOW()),
('category', 21, 'en', 'name', 'Apartments', true, true, NOW(), NOW()),

-- Подкатегории Транспорт
('category', 22, 'en', 'name', 'Passenger Cars', true, true, NOW(), NOW()),
('category', 23, 'en', 'name', 'Commercial Vehicles', true, true, NOW(), NOW()),
('category', 24, 'en', 'name', 'Special Equipment', true, true, NOW(), NOW()),
('category', 25, 'en', 'name', 'Agricultural Vehicles', true, true, NOW(), NOW()),
('category', 26, 'en', 'name', 'Vehicle & Equipment Rental', true, true, NOW(), NOW()),
('category', 27, 'en', 'name', 'Motorcycles & Moto Equipment', true, true, NOW(), NOW()),
('category', 28, 'en', 'name', 'Water Transport', true, true, NOW(), NOW()),
('category', 29, 'en', 'name', 'Parts & Accessories', true, true, NOW(), NOW());

-- Обновляем sequence для translations
SELECT setval('translations_id_seq', (SELECT MAX(id) FROM translations), true);