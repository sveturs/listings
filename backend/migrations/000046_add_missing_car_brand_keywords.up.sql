-- Добавляем полные названия брендов с дефисами и другие варианты написания
-- для улучшения определения категорий

INSERT INTO category_keywords (category_id, keyword, language, weight, keyword_type, is_negative, created_at)
VALUES
    -- Полные названия брендов с дефисами
    (1301, 'mercedes-benz', 'en', 0.8, 'brand', false, NOW()),
    (1301, 'mercedes benz', 'en', 0.8, 'brand', false, NOW()),
    (1301, 'alfa-romeo', 'en', 0.7, 'brand', false, NOW()),
    (1301, 'alfa romeo', 'en', 0.7, 'brand', false, NOW()),
    (1301, 'land-rover', 'en', 0.7, 'brand', false, NOW()),
    (1301, 'land rover', 'en', 0.7, 'brand', false, NOW()),
    (1301, 'rolls-royce', 'en', 0.7, 'brand', false, NOW()),
    (1301, 'rolls royce', 'en', 0.7, 'brand', false, NOW()),

    -- Популярные модели как ключевые слова
    (1301, 'e-class', 'en', 0.6, 'attribute', false, NOW()),
    (1301, 'c-class', 'en', 0.6, 'attribute', false, NOW()),
    (1301, 's-class', 'en', 0.6, 'attribute', false, NOW()),
    (1301, 'a-class', 'en', 0.6, 'attribute', false, NOW()),
    (1301, 'g-class', 'en', 0.6, 'attribute', false, NOW()),
    (1301, 'gle', 'en', 0.6, 'attribute', false, NOW()),
    (1301, 'glc', 'en', 0.6, 'attribute', false, NOW()),
    (1301, 'gla', 'en', 0.6, 'attribute', false, NOW()),
    (1301, 'glb', 'en', 0.6, 'attribute', false, NOW()),
    (1301, 'gls', 'en', 0.6, 'attribute', false, NOW()),

    -- BMW модели
    (1301, 'x1', 'en', 0.6, 'attribute', false, NOW()),
    (1301, 'x2', 'en', 0.6, 'attribute', false, NOW()),
    (1301, 'x3', 'en', 0.6, 'attribute', false, NOW()),
    (1301, 'x4', 'en', 0.6, 'attribute', false, NOW()),
    (1301, 'x5', 'en', 0.6, 'attribute', false, NOW()),
    (1301, 'x6', 'en', 0.6, 'attribute', false, NOW()),
    (1301, 'x7', 'en', 0.6, 'attribute', false, NOW()),

    -- Audi модели
    (1301, 'a1', 'en', 0.6, 'attribute', false, NOW()),
    (1301, 'a3', 'en', 0.6, 'attribute', false, NOW()),
    (1301, 'a4', 'en', 0.6, 'attribute', false, NOW()),
    (1301, 'a5', 'en', 0.6, 'attribute', false, NOW()),
    (1301, 'a6', 'en', 0.6, 'attribute', false, NOW()),
    (1301, 'a7', 'en', 0.6, 'attribute', false, NOW()),
    (1301, 'a8', 'en', 0.6, 'attribute', false, NOW()),
    (1301, 'q2', 'en', 0.6, 'attribute', false, NOW()),
    (1301, 'q3', 'en', 0.6, 'attribute', false, NOW()),
    (1301, 'q5', 'en', 0.6, 'attribute', false, NOW()),
    (1301, 'q7', 'en', 0.6, 'attribute', false, NOW()),
    (1301, 'q8', 'en', 0.6, 'attribute', false, NOW()),

    -- Volkswagen модели
    (1301, 'golf', 'en', 0.6, 'attribute', false, NOW()),
    (1301, 'passat', 'en', 0.6, 'attribute', false, NOW()),
    (1301, 'tiguan', 'en', 0.6, 'attribute', false, NOW()),
    (1301, 'touran', 'en', 0.6, 'attribute', false, NOW()),
    (1301, 'touareg', 'en', 0.6, 'attribute', false, NOW()),
    (1301, 'polo', 'en', 0.6, 'attribute', false, NOW()),
    (1301, 'arteon', 'en', 0.6, 'attribute', false, NOW()),
    (1301, 't-roc', 'en', 0.6, 'attribute', false, NOW()),
    (1301, 't-cross', 'en', 0.6, 'attribute', false, NOW()),

    -- Добавляем типы двигателей и привода (часто упоминаются в описаниях)
    (1301, 'tdi', 'en', 0.5, 'attribute', false, NOW()),
    (1301, 'tfsi', 'en', 0.5, 'attribute', false, NOW()),
    (1301, 'tsi', 'en', 0.5, 'attribute', false, NOW()),
    (1301, 'cdi', 'en', 0.5, 'attribute', false, NOW()),
    (1301, 'bluemotion', 'en', 0.5, 'attribute', false, NOW()),
    (1301, 'xdrive', 'en', 0.5, 'attribute', false, NOW()),
    (1301, 'quattro', 'en', 0.5, 'attribute', false, NOW()),
    (1301, '4matic', 'en', 0.5, 'attribute', false, NOW()),
    (1301, '4motion', 'en', 0.5, 'attribute', false, NOW())
ON CONFLICT (category_id, keyword, language) DO UPDATE SET
    weight = EXCLUDED.weight,
    keyword_type = EXCLUDED.keyword_type,
    is_negative = EXCLUDED.is_negative,
    updated_at = NOW();