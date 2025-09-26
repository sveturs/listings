-- Исправляем неправильные маппинги для автомобилей
-- Заменяем категорию 1002 (fashion) на правильную 1301 (cars)

-- Удаляем неправильные маппинги
DELETE FROM category_ai_mappings
WHERE ai_domain = 'automotive' AND category_id = 1002;

-- Добавляем правильные маппинги для автомобильной категории
INSERT INTO category_ai_mappings (ai_domain, product_type, category_id, weight, success_count, failure_count, is_active)
VALUES
    -- Основные типы автомобилей
    ('automotive', 'car', 1301, 0.95, 0, 0, true),
    ('automotive', 'vehicle', 1301, 0.90, 0, 0, true),
    ('automotive', 'automobile', 1301, 0.95, 0, 0, true),
    ('automotive', 'auto', 1301, 0.90, 0, 0, true),

    -- Типы кузовов
    ('automotive', 'sedan', 1301, 0.95, 0, 0, true),
    ('automotive', 'suv', 1301, 0.95, 0, 0, true),
    ('automotive', 'minivan', 1301, 0.95, 0, 0, true),
    ('automotive', 'hatchback', 1301, 0.95, 0, 0, true),
    ('automotive', 'coupe', 1301, 0.95, 0, 0, true),
    ('automotive', 'wagon', 1301, 0.95, 0, 0, true),
    ('automotive', 'pickup', 1301, 0.95, 0, 0, true),
    ('automotive', 'crossover', 1301, 0.95, 0, 0, true),
    ('automotive', 'convertible', 1301, 0.95, 0, 0, true),
    ('automotive', 'van', 1301, 0.90, 0, 0, true),

    -- Специализированные категории
    ('automotive', 'electric-car', 10170, 0.95, 0, 0, true), -- Электромобили
    ('automotive', 'hybrid', 10171, 0.95, 0, 0, true),      -- Гибриды
    ('automotive', 'luxury-car', 10172, 0.95, 0, 0, true),  -- Люксовые авто
    ('automotive', 'sports-car', 10173, 0.95, 0, 0, true),  -- Спортивные авто
    ('automotive', 'city-car', 10176, 0.95, 0, 0, true),    -- Городские авто

    -- Мотоциклы
    ('automotive', 'motorcycle', 1302, 0.95, 0, 0, true),
    ('automotive', 'scooter', 1302, 0.95, 0, 0, true),
    ('automotive', 'moped', 1302, 0.90, 0, 0, true),

    -- Автозапчасти (включая шины и диски - все в одной категории)
    ('automotive', 'auto-parts', 1303, 0.95, 0, 0, true),
    ('automotive', 'spare-parts', 1303, 0.95, 0, 0, true),
    ('automotive', 'car-parts', 1303, 0.95, 0, 0, true),
    ('automotive', 'tires', 1303, 0.95, 0, 0, true),
    ('automotive', 'wheels', 1303, 0.95, 0, 0, true),
    ('automotive', 'tires-wheels', 1303, 0.95, 0, 0, true)
ON CONFLICT (ai_domain, product_type, category_id) DO UPDATE SET
    weight = EXCLUDED.weight,
    is_active = EXCLUDED.is_active,
    updated_at = CURRENT_TIMESTAMP;

-- Добавляем маппинги для других основных доменов, чтобы система работала корректно
INSERT INTO category_ai_mappings (ai_domain, product_type, category_id, weight, success_count, failure_count, is_active)
VALUES
    -- Электроника
    ('electronics', 'smartphone', 1001, 0.95, 0, 0, true),
    ('electronics', 'laptop', 1001, 0.95, 0, 0, true),
    ('electronics', 'computer', 1001, 0.95, 0, 0, true),
    ('electronics', 'tablet', 1001, 0.95, 0, 0, true),
    ('electronics', 'tv', 1001, 0.95, 0, 0, true),
    ('electronics', 'camera', 1001, 0.95, 0, 0, true),
    ('electronics', 'headphones', 1001, 0.95, 0, 0, true),
    ('electronics', 'smartwatch', 1001, 0.95, 0, 0, true),
    ('electronics', 'gaming-console', 1001, 0.95, 0, 0, true),

    -- Одежда и мода
    ('fashion', 'clothing', 1002, 0.95, 0, 0, true),
    ('fashion', 'shoes', 1002, 0.95, 0, 0, true),
    ('fashion', 'accessories', 1002, 0.95, 0, 0, true),
    ('fashion', 'bags', 1002, 0.95, 0, 0, true),
    ('fashion', 'jewelry', 1002, 0.95, 0, 0, true),
    ('fashion', 'watches', 1002, 0.95, 0, 0, true),

    -- Недвижимость
    ('real-estate', 'apartment', 1004, 0.95, 0, 0, true),
    ('real-estate', 'house', 1004, 0.95, 0, 0, true),
    ('real-estate', 'land', 1004, 0.95, 0, 0, true),
    ('real-estate', 'commercial', 1004, 0.95, 0, 0, true),
    ('real-estate', 'garage', 1004, 0.95, 0, 0, true),

    -- Дом и сад
    ('home-garden', 'furniture', 1005, 0.95, 0, 0, true),
    ('home-garden', 'appliances', 1005, 0.95, 0, 0, true),
    ('home-garden', 'tools', 1005, 0.95, 0, 0, true),
    ('home-garden', 'garden', 1005, 0.95, 0, 0, true),
    ('home-garden', 'decor', 1005, 0.95, 0, 0, true),

    -- Спорт и отдых
    ('sports', 'equipment', 1010, 0.95, 0, 0, true),
    ('sports', 'fitness', 1010, 0.95, 0, 0, true),
    ('sports', 'outdoor', 1010, 0.95, 0, 0, true),
    ('sports', 'bicycle', 1010, 0.95, 0, 0, true),
    ('sports', 'camping', 1010, 0.95, 0, 0, true)
ON CONFLICT (ai_domain, product_type, category_id) DO UPDATE SET
    weight = EXCLUDED.weight,
    is_active = EXCLUDED.is_active,
    updated_at = CURRENT_TIMESTAMP;