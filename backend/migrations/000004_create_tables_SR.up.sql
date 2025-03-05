-- Сербские переводы основных категорий
INSERT INTO translations (entity_type, entity_id, language, field_name, translated_text, is_machine_translated, is_verified, created_at, updated_at) VALUES
-- Основные категории
('category', 1, 'sr', 'name', 'Nekretnine', true, true, NOW(), NOW()),
('category', 2, 'sr', 'name', 'Automobili', true, true, NOW(), NOW()),
('category', 3, 'sr', 'name', 'Elektronika', true, true, NOW(), NOW()),
('category', 4, 'sr', 'name', 'Za kuću i stan', true, true, NOW(), NOW()),
('category', 5, 'sr', 'name', 'Za baštu', true, true, NOW(), NOW()),
('category', 6, 'sr', 'name', 'Hobi i razonoda', true, true, NOW(), NOW()),
('category', 7, 'sr', 'name', 'Životinje', true, true, NOW(), NOW()),
('category', 8, 'sr', 'name', 'Poslovi i oprema', true, true, NOW(), NOW()),
('category', 9, 'sr', 'name', 'Ostalo', true, true, NOW(), NOW()),
('category', 10, 'sr', 'name', 'Posao', true, true, NOW(), NOW()),
('category', 11, 'sr', 'name', 'Odeća, obuća, dodaci', true, true, NOW(), NOW()),
('category', 12, 'sr', 'name', 'Dečja odeća i obuća', true, true, NOW(), NOW()),

-- Подкатегории Недвижимости
('category', 13, 'sr', 'name', 'Stan', true, true, NOW(), NOW()),
('category', 14, 'sr', 'name', 'Soba', true, true, NOW(), NOW()),
('category', 15, 'sr', 'name', 'Kuća, vikendica', true, true, NOW(), NOW()),
('category', 16, 'sr', 'name', 'Plac', true, true, NOW(), NOW()),
('category', 17, 'sr', 'name', 'Garaža i parking', true, true, NOW(), NOW()),
('category', 18, 'sr', 'name', 'Poslovni prostor', true, true, NOW(), NOW()),
('category', 19, 'sr', 'name', 'Nekretnine u inostranstvu', true, true, NOW(), NOW()),
('category', 20, 'sr', 'name', 'Hotel', true, true, NOW(), NOW()),
('category', 21, 'sr', 'name', 'Apartmani', true, true, NOW(), NOW()),

-- Подкатегории Транспорт
('category', 22, 'sr', 'name', 'Putnički automobili', true, true, NOW(), NOW()),
('category', 23, 'sr', 'name', 'Teretna vozila', true, true, NOW(), NOW()),
('category', 24, 'sr', 'name', 'Specijalna vozila', true, true, NOW(), NOW()),
('category', 25, 'sr', 'name', 'Poljoprivredna vozila', true, true, NOW(), NOW()),
('category', 26, 'sr', 'name', 'Iznajmljivanje vozila', true, true, NOW(), NOW()),
('category', 27, 'sr', 'name', 'Motocikli i oprema', true, true, NOW(), NOW()),
('category', 28, 'sr', 'name', 'Plovila', true, true, NOW(), NOW()),
('category', 29, 'sr', 'name', 'Delovi i dodaci', true, true, NOW(), NOW());

-- Обновляем sequence для translations
SELECT setval('translations_id_seq', (SELECT MAX(id) FROM translations), true);


-- Insert marketplace listings
INSERT INTO marketplace_listings (id, user_id, category_id, title, description, price, condition, status, location, latitude, longitude, address_city, address_country, views_count, created_at, updated_at, show_on_map, original_language) VALUES
(8, 2, 13, 'Toyota Corolla 2018', 'Продајем Toyota Corolla 2018 годиште, 80.000 км, одлично стање. Први власник, редовно одржавање, сва документација доступна.', 1150000.00, 'used', 'active', 'Нови Сад, Србија', 45.26710000, 19.83350000, 'Нови Сад', 'Србија', 0, '2025-02-07 07:13:52.973909', '2025-02-07 07:13:52.973909', true, 'sr'),
(9, 3, 13, 'mobile Samsung Galaxy S21', 'Selling Samsung Galaxy S21, 256GB, Deep Purple. Perfect condition, complete set with original box and accessories. AppleCare+ until 2024.', 120000.00, 'used', 'active', 'Novi Sad, Serbia', 45.25510000, 19.84520000, 'Novi Sad', 'Serbia', 0, '2025-02-07 07:13:52.973909', '2025-02-07 07:13:52.973909', true, 'en'),
(10, 4, 13, 'Игровой компьютер RTX 4080', 'Продаю мощный игровой ПК: Intel Core i9-13900K, RTX 4080, 64GB RAM, 2TB NVMe SSD. Идеален для любых игр и тяжелых задач.', 350000.00, 'used', 'active', 'Нови-Сад, Сербия', 45.25410000, 19.84010000, 'Нови-Сад', 'Сербия', 0, '2025-02-07 07:13:52.973909', '2025-02-07 07:13:52.973909', true, 'ru'),
(12, 2, 13, 'автомобиль Toyota Corolla 2018', 'Продаю Toyota Corolla 2018 года, 80.000 км, отличное состояние. Первый владелец, регулярное обслуживание, вся документация в наличии.', 1475000.00, 'used', 'active', 'Косте Мајинског 4, Ветерник, Сербия', 45.24755670, 19.76878366, 'Ветерник', 'Сербия', 0, '2025-02-07 17:33:27.680035', '2025-02-07 17:40:23.957971', true, 'ru');

SELECT setval('marketplace_listings_id_seq', 12, true);

-- Insert marketplace images
INSERT INTO marketplace_images (id, listing_id, file_path, file_name, file_size, content_type, is_main, created_at) VALUES
(15, 8, 'toyota_1.jpg', 'toyota_1.jpg', 1024, 'image/jpeg', true, '2025-02-07 07:13:52.973909'),
(16, 8, 'toyota_2.jpg', 'toyota_2.jpg', 1024, 'image/jpeg', false, '2025-02-07 07:13:52.973909'),
(17, 9, 'galaxy_s21_1.jpg', 'galaxy_s21_1.jpg', 1024, 'image/jpeg', true, '2025-02-07 07:13:52.973909'),
(18, 9, 'galaxy_s21_2.jpg', 'galaxy_s21_2.jpg', 1024, 'image/jpeg', false, '2025-02-07 07:13:52.973909'),
(19, 10, 'gaming_pc_1.jpg', 'gaming_pc_1.jpg', 1024, 'image/jpeg', true, '2025-02-07 07:13:52.973909'),
(20, 10, 'gaming_pc_2.jpg', 'gaming_pc_2.jpg', 1024, 'image/jpeg', false, '2025-02-07 07:13:52.973909'),
(21, 12, 'toyota_1.jpg', 'toyota_1.jpg', 454842, 'image/jpeg', true, '2025-02-07 17:35:09.579393'),
(22, 12, 'toyota_2.jpg', 'toyota_2.jpg', 398035, 'image/jpeg', true, '2025-02-07 17:40:24.397595');

SELECT setval('marketplace_images_id_seq', 22, true);

-- Insert reviews and related data
INSERT INTO reviews (id, user_id, entity_type, entity_id, rating, comment, pros, cons, photos, likes_count, is_verified_purchase, status, created_at, updated_at, helpful_votes, not_helpful_votes, original_language) VALUES
(1, 2, 'listing', 8, 5, 'норм', NULL, NULL, NULL, 0, true, 'published', '2025-02-07 07:47:17.001726', '2025-02-07 14:25:23.586871', 0, 1, 'ru');

SELECT setval('reviews_id_seq', 1, true);

INSERT INTO review_responses (id, review_id, user_id, response, created_at, updated_at) VALUES
(1, 1, 3, 'ok', '2025-02-07 07:49:14.935308', '2025-02-07 07:49:14.935308');

SELECT setval('review_responses_id_seq', 1, true);

INSERT INTO review_votes (review_id, user_id, vote_type, created_at) VALUES
(1, 3, 'not_helpful', '2025-02-07 07:48:11.709016');
