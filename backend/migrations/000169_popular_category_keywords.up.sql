-- Добавление ключевых слов для пяти самых популярных категорий
-- Используем ON CONFLICT DO NOTHING для избежания дублирования

-- 1. ЭЛЕКТРОНИКА (Electronics) - ID: 1001
-- Основные ключевые слова
INSERT INTO category_keywords (category_id, keyword, language, weight, keyword_type, is_negative, created_at, updated_at) VALUES
-- Русский
(1001, 'электроника', 'ru', 10.0, 'main', false, NOW(), NOW()),
(1001, 'техника', 'ru', 9.0, 'main', false, NOW(), NOW()),
(1001, 'гаджеты', 'ru', 8.0, 'synonym', false, NOW(), NOW()),
(1001, 'устройства', 'ru', 7.5, 'synonym', false, NOW(), NOW()),
(1001, 'цифровая техника', 'ru', 8.5, 'synonym', false, NOW(), NOW()),
(1001, 'электронные товары', 'ru', 8.0, 'synonym', false, NOW(), NOW()),
(1001, 'бытовая техника', 'ru', 7.0, 'synonym', false, NOW(), NOW()),
-- Английский
(1001, 'electronics', 'en', 10.0, 'main', false, NOW(), NOW()),
(1001, 'gadgets', 'en', 8.5, 'synonym', false, NOW(), NOW()),
(1001, 'devices', 'en', 7.5, 'synonym', false, NOW(), NOW()),
(1001, 'technology', 'en', 8.0, 'synonym', false, NOW(), NOW()),
(1001, 'tech', 'en', 7.5, 'synonym', false, NOW(), NOW()),
(1001, 'digital', 'en', 7.0, 'synonym', false, NOW(), NOW()),
(1001, 'electronic goods', 'en', 8.0, 'synonym', false, NOW(), NOW()),
-- Сербский
(1001, 'elektronika', 'sr', 10.0, 'main', false, NOW(), NOW()),
(1001, 'tehnika', 'sr', 9.0, 'main', false, NOW(), NOW()),
(1001, 'uređaji', 'sr', 7.5, 'synonym', false, NOW(), NOW()),
(1001, 'digitalna tehnika', 'sr', 8.5, 'synonym', false, NOW(), NOW()),
(1001, 'elektronski proizvodi', 'sr', 8.0, 'synonym', false, NOW(), NOW()),
(1001, 'kućni aparati', 'sr', 7.0, 'synonym', false, NOW(), NOW())
ON CONFLICT (category_id, keyword, language) DO NOTHING;

-- Дополнительные ключевые слова для смартфонов (ID: 1101) - пропускаем дублирующиеся
INSERT INTO category_keywords (category_id, keyword, language, weight, keyword_type, is_negative, created_at, updated_at) VALUES
-- Дополнительные бренды и синонимы
(1101, 'huawei', 'ru', 7.0, 'brand', false, NOW(), NOW()),
(1101, 'oppo', 'ru', 6.5, 'brand', false, NOW(), NOW()),
(1101, 'vivo', 'ru', 6.5, 'brand', false, NOW(), NOW()),
(1101, 'смартфоны', 'ru', 9.0, 'synonym', false, NOW(), NOW()),
-- Английский дополнительные
(1101, 'huawei', 'en', 7.0, 'brand', false, NOW(), NOW()),
(1101, 'oppo', 'en', 6.5, 'brand', false, NOW(), NOW()),
(1101, 'vivo', 'en', 6.5, 'brand', false, NOW(), NOW()),
(1101, 'smartphones', 'en', 9.0, 'synonym', false, NOW(), NOW()),
-- Сербский дополнительные
(1101, 'pametni telefoni', 'sr', 9.5, 'synonym', false, NOW(), NOW()),
(1101, 'telefoni', 'sr', 9.0, 'synonym', false, NOW(), NOW())
ON CONFLICT (category_id, keyword, language) DO NOTHING;

-- Компьютеры (ID: 1102) - добавляем ключевые слова
INSERT INTO category_keywords (category_id, keyword, language, weight, keyword_type, is_negative, created_at, updated_at) VALUES
-- Русский
(1102, 'компьютер', 'ru', 10.0, 'main', false, NOW(), NOW()),
(1102, 'ноутбук', 'ru', 9.5, 'main', false, NOW(), NOW()),
(1102, 'пк', 'ru', 8.5, 'synonym', false, NOW(), NOW()),
(1102, 'лэптоп', 'ru', 8.0, 'synonym', false, NOW(), NOW()),
(1102, 'процессор', 'ru', 7.0, 'attribute', false, NOW(), NOW()),
(1102, 'компьютеры', 'ru', 9.0, 'synonym', false, NOW(), NOW()),
-- Английский
(1102, 'computer', 'en', 10.0, 'main', false, NOW(), NOW()),
(1102, 'laptop', 'en', 9.5, 'main', false, NOW(), NOW()),
(1102, 'pc', 'en', 8.5, 'synonym', false, NOW(), NOW()),
(1102, 'desktop', 'en', 8.0, 'synonym', false, NOW(), NOW()),
(1102, 'computers', 'en', 9.0, 'synonym', false, NOW(), NOW()),
-- Сербский
(1102, 'računar', 'sr', 10.0, 'main', false, NOW(), NOW()),
(1102, 'laptop', 'sr', 9.5, 'main', false, NOW(), NOW()),
(1102, 'kompjuter', 'sr', 8.0, 'synonym', false, NOW(), NOW()),
(1102, 'računari', 'sr', 9.0, 'synonym', false, NOW(), NOW())
ON CONFLICT (category_id, keyword, language) DO NOTHING;

-- 2. МОДА (Fashion) - ID: 1002
INSERT INTO category_keywords (category_id, keyword, language, weight, keyword_type, is_negative, created_at, updated_at) VALUES
-- Русский
(1002, 'мода', 'ru', 10.0, 'main', false, NOW(), NOW()),
(1002, 'одежда', 'ru', 9.5, 'main', false, NOW(), NOW()),
(1002, 'стиль', 'ru', 8.0, 'synonym', false, NOW(), NOW()),
(1002, 'модная одежда', 'ru', 8.5, 'synonym', false, NOW(), NOW()),
(1002, 'аксессуары', 'ru', 7.5, 'synonym', false, NOW(), NOW()),
-- Английский
(1002, 'fashion', 'en', 10.0, 'main', false, NOW(), NOW()),
(1002, 'clothing', 'en', 9.5, 'main', false, NOW(), NOW()),
(1002, 'apparel', 'en', 8.5, 'synonym', false, NOW(), NOW()),
(1002, 'style', 'en', 8.0, 'synonym', false, NOW(), NOW()),
(1002, 'wear', 'en', 7.5, 'synonym', false, NOW(), NOW()),
-- Сербский
(1002, 'moda', 'sr', 10.0, 'main', false, NOW(), NOW()),
(1002, 'odeća', 'sr', 9.5, 'main', false, NOW(), NOW()),
(1002, 'stil', 'sr', 8.0, 'synonym', false, NOW(), NOW()),
(1002, 'garderoba', 'sr', 7.5, 'synonym', false, NOW(), NOW())
ON CONFLICT (category_id, keyword, language) DO NOTHING;

-- Мужская одежда (ID: 1201)
INSERT INTO category_keywords (category_id, keyword, language, weight, keyword_type, is_negative, created_at, updated_at) VALUES
-- Русский
(1201, 'мужская одежда', 'ru', 10.0, 'main', false, NOW(), NOW()),
(1201, 'для мужчин', 'ru', 9.0, 'synonym', false, NOW(), NOW()),
(1201, 'мужское', 'ru', 8.5, 'synonym', false, NOW(), NOW()),
(1201, 'рубашки', 'ru', 8.0, 'attribute', false, NOW(), NOW()),
(1201, 'костюмы', 'ru', 7.5, 'attribute', false, NOW(), NOW()),
-- Английский
(1201, 'mens clothing', 'en', 10.0, 'main', false, NOW(), NOW()),
(1201, 'men', 'en', 9.0, 'synonym', false, NOW(), NOW()),
(1201, 'mens wear', 'en', 8.5, 'synonym', false, NOW(), NOW()),
(1201, 'shirts', 'en', 8.0, 'attribute', false, NOW(), NOW()),
(1201, 'suits', 'en', 7.5, 'attribute', false, NOW(), NOW()),
-- Сербский
(1201, 'muška odeća', 'sr', 10.0, 'main', false, NOW(), NOW()),
(1201, 'za muškarce', 'sr', 9.0, 'synonym', false, NOW(), NOW()),
(1201, 'muško', 'sr', 8.5, 'synonym', false, NOW(), NOW())
ON CONFLICT (category_id, keyword, language) DO NOTHING;

-- Женская одежда (ID: 1202)
INSERT INTO category_keywords (category_id, keyword, language, weight, keyword_type, is_negative, created_at, updated_at) VALUES
-- Русский
(1202, 'женская одежда', 'ru', 10.0, 'main', false, NOW(), NOW()),
(1202, 'для женщин', 'ru', 9.0, 'synonym', false, NOW(), NOW()),
(1202, 'женское', 'ru', 8.5, 'synonym', false, NOW(), NOW()),
(1202, 'платья', 'ru', 8.0, 'attribute', false, NOW(), NOW()),
(1202, 'юбки', 'ru', 7.5, 'attribute', false, NOW(), NOW()),
-- Английский
(1202, 'womens clothing', 'en', 10.0, 'main', false, NOW(), NOW()),
(1202, 'women', 'en', 9.0, 'synonym', false, NOW(), NOW()),
(1202, 'ladies wear', 'en', 8.5, 'synonym', false, NOW(), NOW()),
(1202, 'dresses', 'en', 8.0, 'attribute', false, NOW(), NOW()),
(1202, 'skirts', 'en', 7.5, 'attribute', false, NOW(), NOW()),
-- Сербский
(1202, 'ženska odeća', 'sr', 10.0, 'main', false, NOW(), NOW()),
(1202, 'za žene', 'sr', 9.0, 'synonym', false, NOW(), NOW()),
(1202, 'žensko', 'sr', 8.5, 'synonym', false, NOW(), NOW())
ON CONFLICT (category_id, keyword, language) DO NOTHING;

-- 3. АВТОМОБИЛИ (Automotive) - ID: 1003
INSERT INTO category_keywords (category_id, keyword, language, weight, keyword_type, is_negative, created_at, updated_at) VALUES
-- Русский
(1003, 'автомобили', 'ru', 10.0, 'main', false, NOW(), NOW()),
(1003, 'машины', 'ru', 9.5, 'main', false, NOW(), NOW()),
(1003, 'авто', 'ru', 9.0, 'synonym', false, NOW(), NOW()),
(1003, 'транспорт', 'ru', 8.0, 'synonym', false, NOW(), NOW()),
(1003, 'автотранспорт', 'ru', 8.5, 'synonym', false, NOW(), NOW()),
-- Английский
(1003, 'automotive', 'en', 10.0, 'main', false, NOW(), NOW()),
(1003, 'cars', 'en', 9.5, 'main', false, NOW(), NOW()),
(1003, 'vehicles', 'en', 9.0, 'synonym', false, NOW(), NOW()),
(1003, 'auto', 'en', 8.5, 'synonym', false, NOW(), NOW()),
(1003, 'motor', 'en', 7.5, 'synonym', false, NOW(), NOW()),
-- Сербский
(1003, 'automobili', 'sr', 10.0, 'main', false, NOW(), NOW()),
(1003, 'kola', 'sr', 9.5, 'main', false, NOW(), NOW()),
(1003, 'vozila', 'sr', 9.0, 'synonym', false, NOW(), NOW()),
(1003, 'auto', 'sr', 8.5, 'synonym', false, NOW(), NOW())
ON CONFLICT (category_id, keyword, language) DO NOTHING;

-- 4. НЕДВИЖИМОСТЬ (Real Estate) - ID: 1004
INSERT INTO category_keywords (category_id, keyword, language, weight, keyword_type, is_negative, created_at, updated_at) VALUES
-- Русский
(1004, 'недвижимость', 'ru', 10.0, 'main', false, NOW(), NOW()),
(1004, 'жилье', 'ru', 9.0, 'synonym', false, NOW(), NOW()),
(1004, 'квартиры', 'ru', 8.5, 'synonym', false, NOW(), NOW()),
(1004, 'дома', 'ru', 8.5, 'synonym', false, NOW(), NOW()),
(1004, 'продажа', 'ru', 7.5, 'context', false, NOW(), NOW()),
(1004, 'аренда', 'ru', 7.5, 'context', false, NOW(), NOW()),
-- Английский
(1004, 'real estate', 'en', 10.0, 'main', false, NOW(), NOW()),
(1004, 'property', 'en', 9.5, 'main', false, NOW(), NOW()),
(1004, 'housing', 'en', 9.0, 'synonym', false, NOW(), NOW()),
(1004, 'apartments', 'en', 8.5, 'synonym', false, NOW(), NOW()),
(1004, 'houses', 'en', 8.5, 'synonym', false, NOW(), NOW()),
(1004, 'rent', 'en', 7.5, 'context', false, NOW(), NOW()),
-- Сербский
(1004, 'nekretnine', 'sr', 10.0, 'main', false, NOW(), NOW()),
(1004, 'stanovanje', 'sr', 9.0, 'synonym', false, NOW(), NOW()),
(1004, 'stanovi', 'sr', 8.5, 'synonym', false, NOW(), NOW()),
(1004, 'kuće', 'sr', 8.5, 'synonym', false, NOW(), NOW())
ON CONFLICT (category_id, keyword, language) DO NOTHING;

-- Дома (ID: 1402) - добавляем, если такой категории нет, или дополняем
INSERT INTO category_keywords (category_id, keyword, language, weight, keyword_type, is_negative, created_at, updated_at) VALUES
-- Русский
(1402, 'дома', 'ru', 10.0, 'main', false, NOW(), NOW()),
(1402, 'коттеджи', 'ru', 8.5, 'synonym', false, NOW(), NOW()),
(1402, 'частные дома', 'ru', 9.0, 'synonym', false, NOW(), NOW()),
(1402, 'загородные дома', 'ru', 8.0, 'synonym', false, NOW(), NOW()),
-- Английский
(1402, 'houses', 'en', 10.0, 'main', false, NOW(), NOW()),
(1402, 'homes', 'en', 9.5, 'main', false, NOW(), NOW()),
(1402, 'cottages', 'en', 8.5, 'synonym', false, NOW(), NOW()),
(1402, 'villas', 'en', 8.0, 'synonym', false, NOW(), NOW()),
-- Сербский
(1402, 'kuće', 'sr', 10.0, 'main', false, NOW(), NOW()),
(1402, 'vikendice', 'sr', 8.5, 'synonym', false, NOW(), NOW()),
(1402, 'porodične kuće', 'sr', 9.0, 'synonym', false, NOW(), NOW())
ON CONFLICT (category_id, keyword, language) DO NOTHING;

-- 5. ДОМ И САД (Home & Garden) - ID: 1005
INSERT INTO category_keywords (category_id, keyword, language, weight, keyword_type, is_negative, created_at, updated_at) VALUES
-- Русский
(1005, 'дом и сад', 'ru', 10.0, 'main', false, NOW(), NOW()),
(1005, 'для дома', 'ru', 9.0, 'synonym', false, NOW(), NOW()),
(1005, 'садоводство', 'ru', 8.5, 'synonym', false, NOW(), NOW()),
(1005, 'интерьер', 'ru', 8.0, 'synonym', false, NOW(), NOW()),
(1005, 'декор', 'ru', 7.5, 'synonym', false, NOW(), NOW()),
(1005, 'мебель', 'ru', 8.5, 'synonym', false, NOW(), NOW()),
-- Английский
(1005, 'home garden', 'en', 10.0, 'main', false, NOW(), NOW()),
(1005, 'home', 'en', 9.5, 'synonym', false, NOW(), NOW()),
(1005, 'garden', 'en', 9.0, 'synonym', false, NOW(), NOW()),
(1005, 'interior', 'en', 8.0, 'synonym', false, NOW(), NOW()),
(1005, 'decor', 'en', 7.5, 'synonym', false, NOW(), NOW()),
(1005, 'furniture', 'en', 8.5, 'synonym', false, NOW(), NOW()),
-- Сербский
(1005, 'dom i bašta', 'sr', 10.0, 'main', false, NOW(), NOW()),
(1005, 'za dom', 'sr', 9.0, 'synonym', false, NOW(), NOW()),
(1005, 'bašta', 'sr', 8.5, 'synonym', false, NOW(), NOW()),
(1005, 'enterijer', 'sr', 8.0, 'synonym', false, NOW(), NOW()),
(1005, 'nameštaj', 'sr', 8.5, 'synonym', false, NOW(), NOW())
ON CONFLICT (category_id, keyword, language) DO NOTHING;

-- Добавляем негативные ключевые слова для избежания неправильной категоризации
INSERT INTO category_keywords (category_id, keyword, language, weight, keyword_type, is_negative, created_at, updated_at) VALUES
-- Для электроники - исключаем книги и еду
(1001, 'книга', 'ru', 5.0, 'synonym', true, NOW(), NOW()),
(1001, 'еда', 'ru', 5.0, 'synonym', true, NOW(), NOW()),
(1001, 'book', 'en', 5.0, 'synonym', true, NOW(), NOW()),  
(1001, 'food', 'en', 5.0, 'synonym', true, NOW(), NOW()),
(1001, 'knjiga', 'sr', 5.0, 'synonym', true, NOW(), NOW()),

-- Для моды - исключаем электронику
(1002, 'компьютер', 'ru', 5.0, 'synonym', true, NOW(), NOW()),
(1002, 'телефон', 'ru', 5.0, 'synonym', true, NOW(), NOW()),
(1002, 'computer', 'en', 5.0, 'synonym', true, NOW(), NOW()),
(1002, 'phone', 'en', 5.0, 'synonym', true, NOW(), NOW()),

-- Для автомобилей - исключаем недвижимость
(1003, 'квартира', 'ru', 5.0, 'synonym', true, NOW(), NOW()),
(1003, 'apartment', 'en', 5.0, 'synonym', true, NOW(), NOW()),
(1003, 'stan', 'sr', 5.0, 'synonym', true, NOW(), NOW()),

-- Для недвижимости - исключаем автомобили
(1004, 'автомобиль', 'ru', 5.0, 'synonym', true, NOW(), NOW()),
(1004, 'car', 'en', 5.0, 'synonym', true, NOW(), NOW()),
(1004, 'auto', 'sr', 5.0, 'synonym', true, NOW(), NOW())
ON CONFLICT (category_id, keyword, language) DO NOTHING;