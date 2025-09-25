-- Добавляем ключевые слова для улучшения определения категорий

-- Очищаем старые ключевые слова для автомобильной категории
DELETE FROM category_keywords WHERE category_id IN (1301, 1302, 1303);

-- Ключевые слова для автомобилей (1301)
INSERT INTO category_keywords (category_id, keyword, language, weight, keyword_type, is_negative, created_at)
VALUES
    -- Русский язык
    (1301, 'автомобиль', 'ru', 1.0, 'main', false, NOW()),
    (1301, 'машина', 'ru', 1.0, 'main', false, NOW()),
    (1301, 'авто', 'ru', 0.9, 'main', false, NOW()),
    (1301, 'легковой', 'ru', 0.8, 'synonym', false, NOW()),
    (1301, 'седан', 'ru', 0.9, 'attribute', false, NOW()),
    (1301, 'хэтчбек', 'ru', 0.9, 'attribute', false, NOW()),
    (1301, 'универсал', 'ru', 0.9, 'attribute', false, NOW()),
    (1301, 'минивэн', 'ru', 0.9, 'attribute', false, NOW()),
    (1301, 'кроссовер', 'ru', 0.9, 'attribute', false, NOW()),
    (1301, 'внедорожник', 'ru', 0.9, 'attribute', false, NOW()),
    (1301, 'купе', 'ru', 0.9, 'attribute', false, NOW()),
    (1301, 'кабриолет', 'ru', 0.9, 'attribute', false, NOW()),
    (1301, 'пикап', 'ru', 0.9, 'attribute', false, NOW()),

    -- Английский язык
    (1301, 'car', 'en', 1.0, 'main', false, NOW()),
    (1301, 'automobile', 'en', 1.0, 'main', false, NOW()),
    (1301, 'vehicle', 'en', 0.8, 'main', false, NOW()),
    (1301, 'sedan', 'en', 0.9, 'attribute', false, NOW()),
    (1301, 'hatchback', 'en', 0.9, 'attribute', false, NOW()),
    (1301, 'wagon', 'en', 0.9, 'attribute', false, NOW()),
    (1301, 'minivan', 'en', 0.9, 'attribute', false, NOW()),
    (1301, 'suv', 'en', 0.9, 'attribute', false, NOW()),
    (1301, 'crossover', 'en', 0.9, 'attribute', false, NOW()),
    (1301, 'coupe', 'en', 0.9, 'attribute', false, NOW()),
    (1301, 'convertible', 'en', 0.9, 'attribute', false, NOW()),
    (1301, 'pickup', 'en', 0.9, 'attribute', false, NOW()),

    -- Сербский язык
    (1301, 'automobil', 'sr', 1.0, 'main', false, NOW()),
    (1301, 'kola', 'sr', 1.0, 'main', false, NOW()),
    (1301, 'vozilo', 'sr', 0.8, 'main', false, NOW()),
    (1301, 'limuzina', 'sr', 0.9, 'attribute', false, NOW()),
    (1301, 'karavan', 'sr', 0.9, 'attribute', false, NOW()),
    (1301, 'dzip', 'sr', 0.9, 'attribute', false, NOW()),
    (1301, 'terenac', 'sr', 0.9, 'attribute', false, NOW()),

    -- Бренды автомобилей (универсальные - используем код 'en' для всех языков)
    (1301, 'volkswagen', 'en', 0.7, 'brand', false, NOW()),
    (1301, 'toyota', 'en', 0.7, 'brand', false, NOW()),
    (1301, 'mercedes', 'en', 0.7, 'brand', false, NOW()),
    (1301, 'bmw', 'en', 0.7, 'brand', false, NOW()),
    (1301, 'audi', 'en', 0.7, 'brand', false, NOW()),
    (1301, 'ford', 'en', 0.7, 'brand', false, NOW()),
    (1301, 'opel', 'en', 0.7, 'brand', false, NOW()),
    (1301, 'renault', 'en', 0.7, 'brand', false, NOW()),
    (1301, 'peugeot', 'en', 0.7, 'brand', false, NOW()),
    (1301, 'citroen', 'en', 0.7, 'brand', false, NOW()),
    (1301, 'skoda', 'en', 0.7, 'brand', false, NOW()),
    (1301, 'honda', 'en', 0.7, 'brand', false, NOW()),
    (1301, 'mazda', 'en', 0.7, 'brand', false, NOW()),
    (1301, 'nissan', 'en', 0.7, 'brand', false, NOW()),
    (1301, 'hyundai', 'en', 0.7, 'brand', false, NOW()),
    (1301, 'kia', 'en', 0.7, 'brand', false, NOW()),
    (1301, 'fiat', 'en', 0.7, 'brand', false, NOW()),
    (1301, 'alfa', 'en', 0.7, 'brand', false, NOW()),
    (1301, 'volvo', 'en', 0.7, 'brand', false, NOW()),
    (1301, 'seat', 'en', 0.7, 'brand', false, NOW())
ON CONFLICT (category_id, keyword, language) DO UPDATE SET
    weight = EXCLUDED.weight,
    keyword_type = EXCLUDED.keyword_type,
    is_negative = EXCLUDED.is_negative;

-- Ключевые слова для мотоциклов (1302)
INSERT INTO category_keywords (category_id, keyword, language, weight, keyword_type, is_negative, created_at)
VALUES
    -- Русский язык
    (1302, 'мотоцикл', 'ru', 1.0, 'main', false, NOW()),
    (1302, 'мотобайк', 'ru', 1.0, 'main', false, NOW()),
    (1302, 'скутер', 'ru', 0.9, 'attribute', false, NOW()),
    (1302, 'мопед', 'ru', 0.9, 'attribute', false, NOW()),

    -- Английский язык
    (1302, 'motorcycle', 'en', 1.0, 'main', false, NOW()),
    (1302, 'motorbike', 'en', 1.0, 'main', false, NOW()),
    (1302, 'scooter', 'en', 0.9, 'attribute', false, NOW()),
    (1302, 'moped', 'en', 0.9, 'attribute', false, NOW()),

    -- Сербский язык
    (1302, 'motocikl', 'sr', 1.0, 'main', false, NOW()),
    (1302, 'motor', 'sr', 0.9, 'main', false, NOW()),
    (1302, 'skuter', 'sr', 0.9, 'attribute', false, NOW())
ON CONFLICT (category_id, keyword, language) DO UPDATE SET
    weight = EXCLUDED.weight,
    keyword_type = EXCLUDED.keyword_type,
    is_negative = EXCLUDED.is_negative;

-- Ключевые слова для автозапчастей (1303)
INSERT INTO category_keywords (category_id, keyword, language, weight, keyword_type, is_negative, created_at)
VALUES
    -- Русский язык
    (1303, 'запчасти', 'ru', 1.0, 'main', false, NOW()),
    (1303, 'автозапчасти', 'ru', 1.0, 'main', false, NOW()),
    (1303, 'детали', 'ru', 0.9, 'main', false, NOW()),
    (1303, 'шины', 'ru', 0.9, 'attribute', false, NOW()),
    (1303, 'резина', 'ru', 0.8, 'attribute', false, NOW()),
    (1303, 'диски', 'ru', 0.9, 'attribute', false, NOW()),
    (1303, 'колеса', 'ru', 0.8, 'attribute', false, NOW()),
    (1303, 'двигатель', 'ru', 0.9, 'attribute', false, NOW()),
    (1303, 'мотор', 'ru', 0.9, 'attribute', false, NOW()),
    (1303, 'коробка', 'ru', 0.8, 'attribute', false, NOW()),
    (1303, 'тормоза', 'ru', 0.8, 'attribute', false, NOW()),

    -- Английский язык
    (1303, 'parts', 'en', 0.8, 'main', false, NOW()),
    (1303, 'spare parts', 'en', 1.0, 'main', false, NOW()),
    (1303, 'auto parts', 'en', 1.0, 'main', false, NOW()),
    (1303, 'tires', 'en', 0.9, 'attribute', false, NOW()),
    (1303, 'wheels', 'en', 0.9, 'attribute', false, NOW()),
    (1303, 'engine', 'en', 0.9, 'attribute', false, NOW()),
    (1303, 'transmission', 'en', 0.8, 'attribute', false, NOW()),
    (1303, 'brakes', 'en', 0.8, 'attribute', false, NOW()),

    -- Сербский язык
    (1303, 'delovi', 'sr', 1.0, 'main', false, NOW()),
    (1303, 'auto delovi', 'sr', 1.0, 'main', false, NOW()),
    (1303, 'gume', 'sr', 0.9, 'attribute', false, NOW()),
    (1303, 'felne', 'sr', 0.9, 'attribute', false, NOW()),
    (1303, 'motor', 'sr', 0.8, 'attribute', false, NOW()),
    (1303, 'menjač', 'sr', 0.8, 'attribute', false, NOW()),
    (1303, 'kočnice', 'sr', 0.8, 'attribute', false, NOW())
ON CONFLICT (category_id, keyword, language) DO UPDATE SET
    weight = EXCLUDED.weight,
    keyword_type = EXCLUDED.keyword_type,
    is_negative = EXCLUDED.is_negative;