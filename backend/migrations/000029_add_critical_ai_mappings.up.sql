-- Критические AI маппинги для повышения точности определения категорий

-- Добавляем AI маппинги для существующих категорий
INSERT INTO category_ai_mappings (ai_domain, product_type, category_id, weight, is_active, created_at)
VALUES
    -- Развлечения (пазлы и игры) - категория 1015
    ('entertainment', 'puzzle', 1015, 0.95, TRUE, NOW()),
    ('entertainment', 'board-game', 1015, 0.95, TRUE, NOW()),
    ('entertainment', 'toy', 1015, 0.90, TRUE, NOW()),
    ('entertainment', 'game', 1015, 0.90, TRUE, NOW()),
    ('entertainment', 'cards', 1015, 0.85, TRUE, NOW()),

    -- Строительные материалы - категория 1504
    ('construction', 'sand', 1504, 0.95, TRUE, NOW()),
    ('construction', 'bulk-materials', 1504, 0.95, TRUE, NOW()),
    ('construction', 'cement', 1504, 0.90, TRUE, NOW()),
    ('construction', 'gravel', 1504, 0.90, TRUE, NOW()),
    ('construction', 'bricks', 1504, 0.90, TRUE, NOW()),

    -- Электроника - категория 1003
    ('electronics', 'smartphone', 1003, 0.95, TRUE, NOW()),
    ('electronics', 'laptop', 1003, 0.95, TRUE, NOW()),
    ('electronics', 'computer', 1003, 0.90, TRUE, NOW()),
    ('electronics', 'tablet', 1003, 0.90, TRUE, NOW()),
    ('electronics', 'phone', 1003, 0.85, TRUE, NOW()),

    -- Автомобили - категория 1002
    ('automotive', 'car', 1002, 0.95, TRUE, NOW()),
    ('automotive', 'vehicle', 1002, 0.90, TRUE, NOW()),
    ('automotive', 'motorcycle', 1002, 0.90, TRUE, NOW()),
    ('automotive', 'truck', 1002, 0.85, TRUE, NOW()),

    -- Недвижимость - категория 1401
    ('real-estate', 'apartment', 1401, 0.95, TRUE, NOW()),
    ('real-estate', 'house', 1401, 0.95, TRUE, NOW()),
    ('real-estate', 'flat', 1401, 0.90, TRUE, NOW()),
    ('real-estate', 'property', 1401, 0.85, TRUE, NOW()),

    -- Природные материалы - используем категорию дом и сад 1005
    ('nature', 'acorn', 1005, 0.85, TRUE, NOW()),
    ('nature', 'seeds', 1005, 0.85, TRUE, NOW()),
    ('nature', 'wood', 1005, 0.80, TRUE, NOW()),
    ('nature', 'plant', 1005, 0.80, TRUE, NOW()),

    -- Антиквариат - используем категорию разное 1001 временно
    ('antiques', 'coin', 1001, 0.75, TRUE, NOW()),
    ('antiques', 'stamp', 1001, 0.75, TRUE, NOW()),
    ('antiques', 'vintage-item', 1001, 0.70, TRUE, NOW()),

    -- Авиация - используем категорию игрушки 1015
    ('aviation', 'aircraft-model', 1015, 0.80, TRUE, NOW()),
    ('aviation', 'airplane-model', 1015, 0.80, TRUE, NOW()),

    -- Военные товары - используем категорию одежда 1004
    ('military', 'uniform', 1004, 0.75, TRUE, NOW()),
    ('military', 'medal', 1001, 0.70, TRUE, NOW()),

    -- Рукоделие - используем категорию дом и сад 1005
    ('crafts', 'craft-kit', 1005, 0.75, TRUE, NOW()),
    ('crafts', 'handmade', 1005, 0.75, TRUE, NOW()),
    ('crafts', 'sewing', 1005, 0.70, TRUE, NOW())
ON CONFLICT (ai_domain, product_type, category_id) DO UPDATE
SET
    weight = EXCLUDED.weight,
    is_active = EXCLUDED.is_active,
    updated_at = NOW();

-- Добавляем ключевые слова для улучшения точности
INSERT INTO category_keyword_weights (keyword, category_id, language, weight, occurrence_count, success_rate)
VALUES
    -- Ключевые слова для пазлов
    ('пазл', 1015, 'ru', 1.0, 1, 1.0),
    ('пазлы', 1015, 'ru', 1.0, 1, 1.0),
    ('puzzle', 1015, 'en', 1.0, 1, 1.0),
    ('деталей', 1015, 'ru', 0.6, 1, 0.6),
    ('головоломка', 1015, 'ru', 0.8, 1, 0.8),

    -- Ключевые слова для строительных материалов
    ('песок', 1504, 'ru', 1.0, 1, 1.0),
    ('цемент', 1504, 'ru', 0.95, 1, 0.95),
    ('мешок', 1504, 'ru', 0.6, 1, 0.6),
    ('строительный', 1504, 'ru', 0.8, 1, 0.8),
    ('sand', 1504, 'en', 1.0, 1, 1.0),
    ('cement', 1504, 'en', 0.95, 1, 0.95),

    -- Ключевые слова для природных материалов
    ('желудь', 1005, 'ru', 0.9, 1, 0.9),
    ('жёлудь', 1005, 'ru', 0.9, 1, 0.9),
    ('семена', 1005, 'ru', 0.85, 1, 0.85),
    ('подсолнух', 1005, 'ru', 0.8, 1, 0.8),
    ('acorn', 1005, 'en', 0.9, 1, 0.9),
    ('seeds', 1005, 'en', 0.85, 1, 0.85),

    -- Ключевые слова для авиации
    ('самолет', 1015, 'ru', 0.7, 1, 0.7),
    ('самолёт', 1015, 'ru', 0.7, 1, 0.7),
    ('модель', 1015, 'ru', 0.6, 1, 0.6),
    ('boeing', 1015, 'all', 0.75, 1, 0.75),
    ('airplane', 1015, 'en', 0.7, 1, 0.7),

    -- Ключевые слова для антиквариата
    ('монета', 1001, 'ru', 0.8, 1, 0.8),
    ('марки', 1001, 'ru', 0.8, 1, 0.8),
    ('старинная', 1001, 'ru', 0.7, 1, 0.7),
    ('ссср', 1001, 'ru', 0.7, 1, 0.7),
    ('coin', 1001, 'en', 0.8, 1, 0.8),
    ('stamp', 1001, 'en', 0.8, 1, 0.8),

    -- Ключевые слова для военных товаров
    ('военная', 1004, 'ru', 0.7, 1, 0.7),
    ('форма', 1004, 'ru', 0.65, 1, 0.65),
    ('камуфляж', 1004, 'ru', 0.75, 1, 0.75),
    ('медаль', 1001, 'ru', 0.7, 1, 0.7),
    ('military', 1004, 'en', 0.7, 1, 0.7),

    -- Ключевые слова для рукоделия
    ('вышивание', 1005, 'ru', 0.8, 1, 0.8),
    ('крестиком', 1005, 'ru', 0.75, 1, 0.75),
    ('набор', 1005, 'ru', 0.5, 1, 0.5),
    ('embroidery', 1005, 'en', 0.8, 1, 0.8),

    -- Ключевые слова для электроники
    ('iphone', 1003, 'all', 1.0, 1, 1.0),
    ('macbook', 1003, 'all', 1.0, 1, 1.0),
    ('смартфон', 1003, 'ru', 0.95, 1, 0.95),
    ('ноутбук', 1003, 'ru', 0.95, 1, 0.95),

    -- Ключевые слова для автомобилей
    ('volkswagen', 1002, 'all', 0.95, 1, 0.95),
    ('golf', 1002, 'all', 0.8, 1, 0.8),
    ('дизель', 1002, 'ru', 0.6, 1, 0.6),

    -- Ключевые слова для недвижимости
    ('квартира', 1401, 'ru', 1.0, 1, 1.0),
    ('комнаты', 1401, 'ru', 0.7, 1, 0.7),
    ('apartment', 1401, 'en', 1.0, 1, 1.0)
ON CONFLICT (keyword, category_id, language) DO UPDATE
SET
    weight = GREATEST(category_keyword_weights.weight, EXCLUDED.weight),
    occurrence_count = category_keyword_weights.occurrence_count + 1,
    success_rate = (category_keyword_weights.success_rate * category_keyword_weights.occurrence_count + EXCLUDED.success_rate) /
                   (category_keyword_weights.occurrence_count + 1);

-- Создаем индексы для улучшения производительности
CREATE INDEX IF NOT EXISTS idx_category_ai_mappings_domain_type
ON category_ai_mappings(ai_domain, product_type);

CREATE INDEX IF NOT EXISTS idx_category_keyword_weights_keyword_lang
ON category_keyword_weights(keyword, language);