-- backend/migrations/0040_add_demo_listings_with_translations.up.sql
-- Сначала очистим существующие демо-данные
DELETE FROM marketplace_images WHERE listing_id IN (SELECT id FROM marketplace_listings WHERE user_id = 1);
DELETE FROM translations WHERE entity_type = 'listing' AND entity_id IN (SELECT id FROM marketplace_listings WHERE user_id = 1);
DELETE FROM marketplace_listings WHERE user_id = 1;

-- Добавляем объявления на разных языках
INSERT INTO marketplace_listings 
(user_id, category_id, title, description, price, condition, status, location, latitude, longitude, address_city, address_country, original_language)
VALUES
-- Объявление на сербском
(1, (SELECT id FROM marketplace_categories WHERE slug = 'cars'),
'Toyota Corolla 2018',
'Продајем Toyota Corolla 2018 годиште, 80.000 км, одлично стање. Први власник, редовно одржавање, сва документација доступна.',
1150000, 'used', 'active',
'Нови Сад, Србија', 45.2671, 19.8335, 'Нови Сад', 'Србија', 'sr'),

-- Объявление на английском
(1, (SELECT id FROM marketplace_categories WHERE slug = 'smartphones'),
'mobile Samsung Galaxy S21',
'Selling Samsung Galaxy S21, 256GB, Deep Purple. Perfect condition, complete set with original box and accessories. AppleCare+ until 2024.',
120000, 'used', 'active',
'Novi Sad, Serbia', 45.2551, 19.8452, 'Novi Sad', 'Serbia', 'en'),

-- Объявление на русском
(1, (SELECT id FROM marketplace_categories WHERE slug = 'computers'),
'Игровой компьютер RTX 4080',
'Продаю мощный игровой ПК: Intel Core i9-13900K, RTX 4080, 64GB RAM, 2TB NVMe SSD. Идеален для любых игр и тяжелых задач.',
350000, 'used', 'active',
'Нови-Сад, Сербия', 45.2541, 19.8401, 'Нови-Сад', 'Сербия', 'ru');

-- Добавляем переводы для первого объявления (сербский оригинал)
INSERT INTO translations (entity_type, entity_id, language, field_name, translated_text, is_machine_translated, is_verified)
VALUES
-- Сербский оригинал
('listing', (SELECT id FROM marketplace_listings WHERE title = 'Toyota Corolla 2018'), 'sr', 'title', 
'Toyota Corolla 2018', false, true),
('listing', (SELECT id FROM marketplace_listings WHERE title = 'Toyota Corolla 2018'), 'sr', 'description',
'Продајем Toyota Corolla 2018 годиште, 80.000 км, одлично стање. Први власник, редовно одржавање, сва документација доступна.',
false, true),
-- Английский перевод
('listing', (SELECT id FROM marketplace_listings WHERE title = 'Toyota Corolla 2018'), 'en', 'title',
'Toyota Corolla 2018', true, false),
('listing', (SELECT id FROM marketplace_listings WHERE title = 'Toyota Corolla 2018'), 'en', 'description',
'Selling Toyota Corolla 2018, 80,000 km, excellent condition. First owner, regular maintenance, all documentation available.',
true, false),
-- Русский перевод
('listing', (SELECT id FROM marketplace_listings WHERE title = 'Toyota Corolla 2018'), 'ru', 'title',
'Toyota Corolla 2018', true, false),
('listing', (SELECT id FROM marketplace_listings WHERE title = 'Toyota Corolla 2018'), 'ru', 'description',
'Продаю Toyota Corolla 2018 года, 80.000 км, отличное состояние. Первый владелец, регулярное обслуживание, вся документация в наличии.',
true, false);

-- Добавляем переводы для второго объявления (английский оригинал)
INSERT INTO translations (entity_type, entity_id, language, field_name, translated_text, is_machine_translated, is_verified)
VALUES
-- Английский оригинал
('listing', (SELECT id FROM marketplace_listings WHERE title = 'mobile Samsung Galaxy S21'), 'en', 'title',
'mobile Samsung Galaxy S21', false, true),
('listing', (SELECT id FROM marketplace_listings WHERE title = 'mobile Samsung Galaxy S21'), 'en', 'description',
'Selling Samsung Galaxy S21 Ultra 5G, 12 GB RAM, 256 GB. Guaranteed, like new.',
false, true),
-- Сербский перевод
('listing', (SELECT id FROM marketplace_listings WHERE title = 'mobile Samsung Galaxy S21'), 'sr', 'title',
'мобилни телефон Samsung Galaxy S21', true, false),
('listing', (SELECT id FROM marketplace_listings WHERE title = 'mobile Samsung Galaxy S21'), 'sr', 'description',
'Samsung Galaxy S21 Ultra 5G, 12GB RAM, 256GB. Под гаранцијом, као нов.',
true, false),
-- Русский перевод
('listing', (SELECT id FROM marketplace_listings WHERE title = 'mobile Samsung Galaxy S21'), 'ru', 'title',
'мобильник Samsung Galaxy S21', true, false),
('listing', (SELECT id FROM marketplace_listings WHERE title = 'mobile Samsung Galaxy S21'), 'ru', 'description',
'Продаю Samsung Galaxy S21 Ultra 5G, 12 ГБ ОЗУ, 256 ГБ. На гарантии, как новый.',
true, false);

-- Добавляем переводы для третьего объявления (русский оригинал)
INSERT INTO translations (entity_type, entity_id, language, field_name, translated_text, is_machine_translated, is_verified)
VALUES
-- Русский оригинал
('listing', (SELECT id FROM marketplace_listings WHERE title = 'Игровой компьютер RTX 4080'), 'ru', 'title',
'Игровой компьютер RTX 4080', false, true),
('listing', (SELECT id FROM marketplace_listings WHERE title = 'Игровой компьютер RTX 4080'), 'ru', 'description',
'Продаю мощный игровой ПК: Intel Core i9-13900K, RTX 4080, 64GB RAM, 2TB NVMe SSD. Идеален для любых игр и тяжелых задач.',
false, true),
-- Сербский перевод
('listing', (SELECT id FROM marketplace_listings WHERE title = 'Игровой компьютер RTX 4080'), 'sr', 'title',
'Гејмерски рачунар RTX 4080', true, false),
('listing', (SELECT id FROM marketplace_listings WHERE title = 'Игровой компьютер RTX 4080'), 'sr', 'description',
'Продајем моћан гејмерски рачунар: Intel Core i9-13900K, RTX 4080, 64GB RAM, 2TB NVMe SSD. Идеалан за све игре и захтевне задатке.',
true, false),
-- Английский перевод
('listing', (SELECT id FROM marketplace_listings WHERE title = 'Игровой компьютер RTX 4080'), 'en', 'title',
'Gaming PC RTX 4080', true, false),
('listing', (SELECT id FROM marketplace_listings WHERE title = 'Игровой компьютер RTX 4080'), 'en', 'description',
'Selling powerful gaming PC: Intel Core i9-13900K, RTX 4080, 64GB RAM, 2TB NVMe SSD. Perfect for any games and demanding tasks.',
true, false);

-- Добавляем изображения
INSERT INTO marketplace_images
(listing_id, file_path, file_name, file_size, content_type, is_main)
VALUES
-- Изображения для Toyota Corolla
((SELECT id FROM marketplace_listings WHERE title = 'Toyota Corolla 2018'),
'toyota_1.jpg', 'toyota_1.jpg', 1024, 'image/jpeg', true),
((SELECT id FROM marketplace_listings WHERE title = 'Toyota Corolla 2018'),
'toyota_2.jpg', 'toyota_2.jpg', 1024, 'image/jpeg', false),

-- Изображения для iPhone
((SELECT id FROM marketplace_listings WHERE title = 'mobile Samsung Galaxy S21'),
'galaxy_s21_1.jpg', 'galaxy_s21_1.jpg', 1024, 'image/jpeg', true),
((SELECT id FROM marketplace_listings WHERE title = 'mobile Samsung Galaxy S21'),
'galaxy_s21_2.jpg', 'galaxy_s21_2.jpg', 1024, 'image/jpeg', false),

-- Изображения для игрового ПК
((SELECT id FROM marketplace_listings WHERE title = 'Игровой компьютер RTX 4080'),
'gaming_pc_1.jpg', 'gaming_pc_1.jpg', 1024, 'image/jpeg', true),
((SELECT id FROM marketplace_listings WHERE title = 'Игровой компьютер RTX 4080'),
'gaming_pc_2.jpg', 'gaming_pc_2.jpg', 1024, 'image/jpeg', false);






