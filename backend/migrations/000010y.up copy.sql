-- Файл исправления для несоответствий ID в категориях

-- 1. Проверка отсутствующих переводов (категории в marketplace_categories, но отсутствуют в translations)
WITH missing_translations AS (
    SELECT mc.id, mc.name, mc.slug
    FROM marketplace_categories mc
    LEFT JOIN translations t ON t.entity_type = 'category' AND t.entity_id = mc.id AND t.language = 'en'
    WHERE t.id IS NULL
)
SELECT * FROM missing_translations;

-- 2. Проверка несуществующих категорий в переводах (переводы есть, но нет категории)
WITH nonexistent_categories AS (
    SELECT t.entity_id, t.language, t.translated_text
    FROM translations t
    LEFT JOIN marketplace_categories mc ON t.entity_type = 'category' AND t.entity_id = mc.id
    WHERE t.entity_type = 'category' AND mc.id IS NULL
)
SELECT * FROM nonexistent_categories;

-- 3. Добавление отсутствующих категорий из paste.txt, которые есть в документе но отсутствуют в БД

-- Сельхозтехника: подкатегории из paste.txt
INSERT INTO marketplace_categories (id, name, slug, parent_id, icon, created_at) VALUES
(501, 'Комбайны', 'combines', 25, 'combine', NOW()),
(502, 'Телескопические погрузчики', 'telescopic-loaders', 25, 'loader', NOW()),
(503, 'Сеялки', 'seeders', 25, 'seeder', NOW()),
(504, 'Культиваторы', 'cultivators', 25, 'cultivator', NOW()),
(505, 'Плуги', 'plows', 25, 'plow', NOW()),
(506, 'Опрыскиватели', 'sprayers', 25, 'sprayer', NOW());

-- Добавляем переводы для новых категорий
INSERT INTO translations (entity_type, entity_id, language, field_name, translated_text, is_machine_translated, is_verified, created_at, updated_at) VALUES
-- English translations
('category', 501, 'en', 'name', 'Combine Harvesters', true, true, NOW(), NOW()),
('category', 502, 'en', 'name', 'Telescopic Loaders', true, true, NOW(), NOW()),
('category', 503, 'en', 'name', 'Seeders', true, true, NOW(), NOW()),
('category', 504, 'en', 'name', 'Cultivators', true, true, NOW(), NOW()),
('category', 505, 'en', 'name', 'Plows', true, true, NOW(), NOW()),
('category', 506, 'en', 'name', 'Sprayers', true, true, NOW(), NOW()),

-- Russian translations
('category', 501, 'ru', 'name', 'Комбайны', true, true, NOW(), NOW()),
('category', 502, 'ru', 'name', 'Телескопические погрузчики', true, true, NOW(), NOW()),
('category', 503, 'ru', 'name', 'Сеялки', true, true, NOW(), NOW()),
('category', 504, 'ru', 'name', 'Культиваторы', true, true, NOW(), NOW()),
('category', 505, 'ru', 'name', 'Плуги', true, true, NOW(), NOW()),
('category', 506, 'ru', 'name', 'Опрыскиватели', true, true, NOW(), NOW()),

-- Serbian translations
('category', 501, 'sr', 'name', 'Kombajni', true, true, NOW(), NOW()),
('category', 502, 'sr', 'name', 'Teleskopski utovarivači', true, true, NOW(), NOW()),
('category', 503, 'sr', 'name', 'Sejalice', true, true, NOW(), NOW()),
('category', 504, 'sr', 'name', 'Kultivatori', true, true, NOW(), NOW()),
('category', 505, 'sr', 'name', 'Plugovi', true, true, NOW(), NOW()),
('category', 506, 'sr', 'name', 'Prskalice', true, true, NOW(), NOW());

-- 4. Исправление несоответствий между файлами Real Estate
-- В русском переводе категории 13-21 не соответствуют английским и сербским
-- В SQL-файле RU идут: Квартира, Комната, Дом/дача, и т.д.
-- А в EN и SR идут: Apartments, Houses, Commercial Space и т.д.

-- Исправляем русские переводы для соответствия структуре
UPDATE translations 
SET translated_text = 'Квартиры'
WHERE entity_type = 'category' AND entity_id = 13 AND language = 'ru';

UPDATE translations 
SET translated_text = 'Дома'
WHERE entity_type = 'category' AND entity_id = 14 AND language = 'ru';

UPDATE translations 
SET translated_text = 'Коммерческие помещения'
WHERE entity_type = 'category' AND entity_id = 15 AND language = 'ru';

UPDATE translations 
SET translated_text = 'Земельные участки'
WHERE entity_type = 'category' AND entity_id = 16 AND language = 'ru';

UPDATE translations 
SET translated_text = 'Гаражи'
WHERE entity_type = 'category' AND entity_id = 17 AND language = 'ru';

UPDATE translations 
SET translated_text = 'Комнаты'
WHERE entity_type = 'category' AND entity_id = 18 AND language = 'ru';

UPDATE translations 
SET translated_text = 'Недвижимость за рубежом'
WHERE entity_type = 'category' AND entity_id = 19 AND language = 'ru';

UPDATE translations 
SET translated_text = 'Продажа коммерческих объектов'
WHERE entity_type = 'category' AND entity_id = 20 AND language = 'ru';

UPDATE translations 
SET translated_text = 'Аренда коммерческих объектов'
WHERE entity_type = 'category' AND entity_id = 21 AND language = 'ru';

-- 5. Исправление ID sequence для корректного автоинкремента
SELECT setval('marketplace_categories_id_seq', (SELECT MAX(id) FROM marketplace_categories), true);
SELECT setval('translations_id_seq', (SELECT MAX(id) FROM translations), true);

-- 6. Добавляем недостающие категории из dokumenta paste.txt

-- Товары для детей и игрушки
INSERT INTO marketplace_categories (id, name, slug, parent_id, icon, created_at) VALUES
(507, 'Товары для детей и игрушки', 'children-goods-toys', NULL, 'toy', NOW());

-- Подкатегории для Товаров для детей и игрушек
INSERT INTO marketplace_categories (id, name, slug, parent_id, icon, created_at) VALUES
(508, 'Детские коляски', 'baby-strollers', 507, 'stroller', NOW()),
(509, 'Детская мебель', 'children-furniture', 507, 'crib', NOW()),
(510, 'Велосипеды и самокаты', 'children-bikes-scooters', 507, 'bike', NOW()),
(511, 'Товары для кормления', 'feeding-goods', 507, 'bottle', NOW()),
(512, 'Автомобильные кресла', 'car-seats', 507, 'car-seat', NOW()),
(513, 'Игрушки', 'toys', 507, 'toy', NOW()),
(514, 'Постельные принадлежности', 'children-bedding', 507, 'blanket', NOW()),
(515, 'Товары для купания', 'bath-goods', 507, 'bath', NOW()),
(516, 'Товары для школы', 'school-supplies', 507, 'backpack', NOW());

-- Добавляем переводы для этих категорий
INSERT INTO translations (entity_type, entity_id, language, field_name, translated_text, is_machine_translated, is_verified, created_at, updated_at) VALUES
-- English translations
('category', 507, 'en', 'name', 'Children Goods and Toys', true, true, NOW(), NOW()),
('category', 508, 'en', 'name', 'Baby Strollers', true, true, NOW(), NOW()),
('category', 509, 'en', 'name', 'Children Furniture', true, true, NOW(), NOW()),
('category', 510, 'en', 'name', 'Bicycles and Scooters', true, true, NOW(), NOW()),
('category', 511, 'en', 'name', 'Feeding Goods', true, true, NOW(), NOW()),
('category', 512, 'en', 'name', 'Car Seats', true, true, NOW(), NOW()),
('category', 513, 'en', 'name', 'Toys', true, true, NOW(), NOW()),
('category', 514, 'en', 'name', 'Bedding', true, true, NOW(), NOW()),
('category', 515, 'en', 'name', 'Bath Goods', true, true, NOW(), NOW()),
('category', 516, 'en', 'name', 'School Supplies', true, true, NOW(), NOW()),

-- Russian translations
('category', 507, 'ru', 'name', 'Товары для детей и игрушки', true, true, NOW(), NOW()),
('category', 508, 'ru', 'name', 'Детские коляски', true, true, NOW(), NOW()),
('category', 509, 'ru', 'name', 'Детская мебель', true, true, NOW(), NOW()),
('category', 510, 'ru', 'name', 'Велосипеды и самокаты', true, true, NOW(), NOW()),
('category', 511, 'ru', 'name', 'Товары для кормления', true, true, NOW(), NOW()),
('category', 512, 'ru', 'name', 'Автомобильные кресла', true, true, NOW(), NOW()),
('category', 513, 'ru', 'name', 'Игрушки', true, true, NOW(), NOW()),
('category', 514, 'ru', 'name', 'Постельные принадлежности', true, true, NOW(), NOW()),
('category', 515, 'ru', 'name', 'Товары для купания', true, true, NOW(), NOW()),
('category', 516, 'ru', 'name', 'Товары для школы', true, true, NOW(), NOW()),

-- Serbian translations
('category', 507, 'sr', 'name', 'Dečja roba i igračke', true, true, NOW(), NOW()),
('category', 508, 'sr', 'name', 'Dečja kolica', true, true, NOW(), NOW()),
('category', 509, 'sr', 'name', 'Dečji nameštaj', true, true, NOW(), NOW()),
('category', 510, 'sr', 'name', 'Bicikli i trotineti', true, true, NOW(), NOW()),
('category', 511, 'sr', 'name', 'Pribor za hranjenje', true, true, NOW(), NOW()),
('category', 512, 'sr', 'name', 'Auto sedišta', true, true, NOW(), NOW()),
('category', 513, 'sr', 'name', 'Igračke', true, true, NOW(), NOW()),
('category', 514, 'sr', 'name', 'Posteljina', true, true, NOW(), NOW()),
('category', 515, 'sr', 'name', 'Pribor za kupanje', true, true, NOW(), NOW()),
('category', 516, 'sr', 'name', 'Školski pribor', true, true, NOW(), NOW());

-- Обновление sequence после добавления новых категорий
SELECT setval('marketplace_categories_id_seq', (SELECT MAX(id) FROM marketplace_categories), true);
SELECT setval('translations_id_seq', (SELECT MAX(id) FROM translations), true);