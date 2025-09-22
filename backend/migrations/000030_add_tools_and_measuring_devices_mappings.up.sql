-- Migration: Add tools and measuring devices AI mappings
-- Fixing critical issue: AI sends "construction tools" but we only have "construction"

-- Найти ID категории для строительных инструментов
-- Если нет - создать новую категорию

-- Проверим существующие строительные категории
INSERT INTO marketplace_categories (name, slug, parent_id, level, icon, is_active)
SELECT 'Строительные инструменты', 'construction-tools', 1007, 2, 'hammer', true
WHERE NOT EXISTS (
    SELECT 1 FROM marketplace_categories
    WHERE slug = 'construction-tools'
);

-- Получим ID созданной категории
DO $$
DECLARE
    tools_category_id INTEGER;
BEGIN
    SELECT id INTO tools_category_id
    FROM marketplace_categories
    WHERE slug = 'construction-tools';

    -- Добавим AI mappings для строительных инструментов
    INSERT INTO category_ai_mappings (ai_domain, product_type, category_id, weight) VALUES
    -- Construction tools domain
    ('construction tools', 'laser measure', tools_category_id, 0.95),
    ('construction tools', 'distance meter', tools_category_id, 0.95),
    ('construction tools', 'measuring device', tools_category_id, 0.95),
    ('construction tools', 'hammer', tools_category_id, 0.95),
    ('construction tools', 'drill', tools_category_id, 0.95),
    ('construction tools', 'saw', tools_category_id, 0.95),
    ('construction tools', 'level', tools_category_id, 0.95),
    ('construction tools', 'screwdriver', tools_category_id, 0.95),
    ('construction tools', 'wrench', tools_category_id, 0.95),
    ('construction tools', 'pliers', tools_category_id, 0.95),

    -- Industrial domain для измерительных приборов
    ('industrial', 'laser measure', tools_category_id, 0.90),
    ('industrial', 'distance meter', tools_category_id, 0.90),
    ('industrial', 'measuring device', tools_category_id, 0.90),
    ('industrial', 'measurement tool', tools_category_id, 0.90),

    -- Electronics domain для лазерных приборов
    ('electronics', 'laser device', tools_category_id, 0.85),
    ('electronics', 'measuring instrument', tools_category_id, 0.85)

    ON CONFLICT (ai_domain, product_type, category_id) DO NOTHING;

    -- Добавим ключевые слова для улучшения детекции
    INSERT INTO category_keyword_weights (category_id, keyword, weight, language) VALUES
    (tools_category_id, 'laser', 0.9, 'en'),
    (tools_category_id, 'measure', 0.9, 'en'),
    (tools_category_id, 'distance', 0.9, 'en'),
    (tools_category_id, 'meter', 0.8, 'en'),
    (tools_category_id, 'bosch', 0.8, 'en'),
    (tools_category_id, 'glm', 0.8, 'en'),
    (tools_category_id, 'construction', 0.7, 'en'),
    (tools_category_id, 'tool', 0.7, 'en'),
    (tools_category_id, 'professional', 0.6, 'en'),
    (tools_category_id, 'измерительный', 0.9, 'ru'),
    (tools_category_id, 'лазерный', 0.9, 'ru'),
    (tools_category_id, 'дальномер', 0.9, 'ru'),
    (tools_category_id, 'инструмент', 0.8, 'ru'),
    (tools_category_id, 'строительный', 0.7, 'ru')
    ON CONFLICT (category_id, keyword, language) DO NOTHING;
END $$;

-- Добавим переводы для новой категории (используем правильную структуру translations)
DO $$
DECLARE
    tools_category_id INTEGER;
BEGIN
    SELECT id INTO tools_category_id
    FROM marketplace_categories
    WHERE slug = 'construction-tools';

    IF tools_category_id IS NOT NULL THEN
        INSERT INTO translations (entity_type, entity_id, language, field_name, translated_text, is_machine_translated, is_verified) VALUES
        ('category', tools_category_id, 'en', 'name', 'Construction Tools', false, true),
        ('category', tools_category_id, 'ru', 'name', 'Строительные инструменты', false, true),
        ('category', tools_category_id, 'sr', 'name', 'Građevinski alati', false, true)
        ON CONFLICT (entity_type, entity_id, language, field_name) DO UPDATE SET
            translated_text = EXCLUDED.translated_text,
            is_verified = true;
    END IF;
END $$;