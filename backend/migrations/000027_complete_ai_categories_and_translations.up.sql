-- Полная миграция для достижения 99% точности AI детекции категорий
-- Добавление недостающих категорий, переводов и маппингов

-- ===================================================================
-- 1. ДОБАВЛЕНИЕ НЕДОСТАЮЩИХ ПОДКАТЕГОРИЙ
-- ===================================================================

-- Подкатегории для "Природные материалы" (10207)
INSERT INTO marketplace_categories (name, slug, parent_id, icon, is_active, level, created_at) VALUES
('Дерево и пиломатериалы', 'wood-materials', 10207, 'tree', true, 1, NOW()),
('Камни и минералы', 'stones-minerals', 10207, 'gem', true, 1, NOW()),
('Растения и семена', 'plants-seeds', 10207, 'flower', true, 1, NOW()),
('Природный декор', 'natural-decor', 10207, 'leaf', true, 1, NOW())
ON CONFLICT (slug) DO NOTHING;

-- Дополнительные подкатегории для "Строительные материалы" (10206)
INSERT INTO marketplace_categories (name, slug, parent_id, icon, is_active, level, created_at) VALUES
('Сыпучие материалы', 'bulk-materials', 10206, 'cube', true, 1, NOW()),
('Инструменты', 'tools', 10206, 'wrench', true, 1, NOW()),
('Краски и лаки', 'paints', 10206, 'palette', true, 1, NOW()),
('Сантехника', 'plumbing', 10206, 'droplet', true, 1, NOW()),
('Электрика', 'electrical', 10206, 'zap', true, 1, NOW())
ON CONFLICT (slug) DO NOTHING;

-- Новые основные категории
INSERT INTO marketplace_categories (name, slug, parent_id, icon, is_active, level, created_at) VALUES
('Декор и рукоделие', 'crafts', NULL, 'palette', true, 0, NOW()),
('Антиквариат', 'antiques', NULL, 'clock', true, 0, NOW()),
('Авиация', 'aviation', NULL, 'plane', true, 0, NOW()),
('Военные товары', 'military', NULL, 'shield', true, 0, NOW()),
('Разное', 'miscellaneous', NULL, 'box', true, 0, NOW())
ON CONFLICT (slug) DO NOTHING;

-- Подкатегории для "Декор и рукоделие"
INSERT INTO marketplace_categories (name, slug, parent_id, icon, is_active, level, created_at)
SELECT
    'Материалы для творчества', 'craft-materials', id, 'brush', true, 1, NOW()
FROM marketplace_categories WHERE slug = 'crafts'
ON CONFLICT (slug) DO NOTHING;

INSERT INTO marketplace_categories (name, slug, parent_id, icon, is_active, level, created_at)
SELECT
    'Готовые изделия', 'handmade', id, 'heart', true, 1, NOW()
FROM marketplace_categories WHERE slug = 'crafts'
ON CONFLICT (slug) DO NOTHING;

INSERT INTO marketplace_categories (name, slug, parent_id, icon, is_active, level, created_at)
SELECT
    'Швейные принадлежности', 'sewing', id, 'scissors', true, 1, NOW()
FROM marketplace_categories WHERE slug = 'crafts'
ON CONFLICT (slug) DO NOTHING;

INSERT INTO marketplace_categories (name, slug, parent_id, icon, is_active, level, created_at)
SELECT
    'Вязание и пряжа', 'knitting', id, 'circle', true, 1, NOW()
FROM marketplace_categories WHERE slug = 'crafts'
ON CONFLICT (slug) DO NOTHING;

-- Дополнительные подкатегории для "Игрушки и развлечения" (1015) если еще нет
INSERT INTO marketplace_categories (name, slug, parent_id, icon, is_active, level, created_at) VALUES
('Настольные игры', 'board-games', 1015, 'dice', true, 1, NOW()),
('Конструкторы', 'constructors', 1015, 'cube', true, 1, NOW()),
('Развивающие игры', 'educational-games', 1015, 'book-open', true, 1, NOW()),
('Модели и сборка', 'models', 1015, 'wrench', true, 1, NOW())
ON CONFLICT (slug) DO NOTHING;

-- ===================================================================
-- 2. ДОБАВЛЕНИЕ ПЕРЕВОДОВ ДЛЯ ВСЕХ НОВЫХ КАТЕГОРИЙ
-- ===================================================================

-- Функция для добавления переводов категории
CREATE OR REPLACE FUNCTION add_category_translations(
    p_category_slug TEXT,
    p_name_en TEXT,
    p_name_ru TEXT,
    p_name_sr TEXT
) RETURNS VOID AS $$
DECLARE
    v_category_id INTEGER;
BEGIN
    SELECT id INTO v_category_id FROM marketplace_categories WHERE slug = p_category_slug;

    IF v_category_id IS NOT NULL THEN
        -- Английский
        INSERT INTO translations (entity_type, entity_id, field_name, language, translated_text, created_at)
        VALUES ('category', v_category_id, 'name', 'en', p_name_en, NOW())
        ON CONFLICT (entity_type, entity_id, field_name, language)
        DO UPDATE SET translated_text = EXCLUDED.translated_text, updated_at = NOW();

        -- Русский
        INSERT INTO translations (entity_type, entity_id, field_name, language, translated_text, created_at)
        VALUES ('category', v_category_id, 'name', 'ru', p_name_ru, NOW())
        ON CONFLICT (entity_type, entity_id, field_name, language)
        DO UPDATE SET translated_text = EXCLUDED.translated_text, updated_at = NOW();

        -- Сербский
        INSERT INTO translations (entity_type, entity_id, field_name, language, translated_text, created_at)
        VALUES ('category', v_category_id, 'name', 'sr', p_name_sr, NOW())
        ON CONFLICT (entity_type, entity_id, field_name, language)
        DO UPDATE SET translated_text = EXCLUDED.translated_text, updated_at = NOW();
    END IF;
END;
$$ LANGUAGE plpgsql;

-- Добавляем переводы для природных материалов
SELECT add_category_translations('natural-materials', 'Natural Materials', 'Природные материалы', 'Prirodni materijali');
SELECT add_category_translations('wood-materials', 'Wood & Lumber', 'Дерево и пиломатериалы', 'Drvo i rezana građa');
SELECT add_category_translations('stones-minerals', 'Stones & Minerals', 'Камни и минералы', 'Kamenje i minerali');
SELECT add_category_translations('plants-seeds', 'Plants & Seeds', 'Растения и семена', 'Biljke i seme');
SELECT add_category_translations('natural-decor', 'Natural Decor', 'Природный декор', 'Prirodni dekor');

-- Переводы для строительных материалов
SELECT add_category_translations('construction-materials', 'Construction Materials', 'Строительные материалы', 'Građevinski materijal');
SELECT add_category_translations('bulk-materials', 'Bulk Materials', 'Сыпучие материалы', 'Rasuti materijali');
SELECT add_category_translations('tools', 'Tools', 'Инструменты', 'Alati');
SELECT add_category_translations('paints', 'Paints & Varnishes', 'Краски и лаки', 'Boje i lakovi');
SELECT add_category_translations('plumbing', 'Plumbing', 'Сантехника', 'Vodoinstalaterski materijal');
SELECT add_category_translations('electrical', 'Electrical', 'Электрика', 'Električni materijal');

-- Переводы для декора и рукоделия
SELECT add_category_translations('crafts', 'Crafts & Handmade', 'Декор и рукоделие', 'Zanati i ručni rad');
SELECT add_category_translations('craft-materials', 'Craft Materials', 'Материалы для творчества', 'Materijali za kreativnost');
SELECT add_category_translations('handmade', 'Handmade Items', 'Готовые изделия', 'Ručno rađeni proizvodi');
SELECT add_category_translations('sewing', 'Sewing Supplies', 'Швейные принадлежности', 'Pribor za šivenje');
SELECT add_category_translations('knitting', 'Knitting & Yarn', 'Вязание и пряжа', 'Pletenje i pređa');

-- Переводы для игрушек и хобби
SELECT add_category_translations('toys', 'Toys', 'Игрушки', 'Igračke');
SELECT add_category_translations('puzzles', 'Puzzles', 'Пазлы', 'Slagalice');
SELECT add_category_translations('board-games', 'Board Games', 'Настольные игры', 'Društvene igre');
SELECT add_category_translations('collectibles', 'Collectibles', 'Коллекционирование', 'Kolekcionarstvo');
SELECT add_category_translations('constructors', 'Building Sets', 'Конструкторы', 'Konstruktori');
SELECT add_category_translations('educational-games', 'Educational Games', 'Развивающие игры', 'Edukativne igre');
SELECT add_category_translations('models', 'Models & Assembly', 'Модели и сборка', 'Modeli i sklapanje');

-- Переводы для специальных категорий
SELECT add_category_translations('antiques', 'Antiques', 'Антиквариат', 'Antikviteti');
SELECT add_category_translations('aviation', 'Aviation', 'Авиация', 'Avijacija');
SELECT add_category_translations('military', 'Military Items', 'Военные товары', 'Vojni predmeti');
SELECT add_category_translations('miscellaneous', 'Miscellaneous', 'Разное', 'Razno');

-- ===================================================================
-- 3. УЛУЧШЕННЫЕ AI МАППИНГИ С УЧЕТОМ НОВЫХ КАТЕГОРИЙ
-- ===================================================================

-- Сначала получаем ID категорий для маппингов
DO $$
DECLARE
    v_toys_id INTEGER;
    v_puzzles_id INTEGER;
    v_board_games_id INTEGER;
    v_collectibles_id INTEGER;
    v_natural_materials_id INTEGER;
    v_construction_materials_id INTEGER;
    v_crafts_id INTEGER;
    v_antiques_id INTEGER;
    v_aviation_id INTEGER;
    v_military_id INTEGER;
    v_misc_id INTEGER;
BEGIN
    SELECT id INTO v_toys_id FROM marketplace_categories WHERE slug = 'toys';
    SELECT id INTO v_puzzles_id FROM marketplace_categories WHERE slug = 'puzzles';
    SELECT id INTO v_board_games_id FROM marketplace_categories WHERE slug = 'board-games';
    SELECT id INTO v_collectibles_id FROM marketplace_categories WHERE slug = 'collectibles';
    SELECT id INTO v_natural_materials_id FROM marketplace_categories WHERE slug = 'natural-materials';
    SELECT id INTO v_construction_materials_id FROM marketplace_categories WHERE slug = 'construction-materials';
    SELECT id INTO v_crafts_id FROM marketplace_categories WHERE slug = 'crafts';
    SELECT id INTO v_antiques_id FROM marketplace_categories WHERE slug = 'antiques';
    SELECT id INTO v_aviation_id FROM marketplace_categories WHERE slug = 'aviation';
    SELECT id INTO v_military_id FROM marketplace_categories WHERE slug = 'military';
    SELECT id INTO v_misc_id FROM marketplace_categories WHERE slug = 'miscellaneous';

    -- Добавляем маппинги для развлечений
    IF v_toys_id IS NOT NULL THEN
        INSERT INTO category_ai_mappings (ai_domain, product_type, category_id, weight, is_active) VALUES
        ('entertainment', 'toy', v_toys_id, 0.95, true),
        ('entertainment', 'doll', v_toys_id, 0.95, true),
        ('entertainment', 'action_figure', v_toys_id, 0.95, true)
        ON CONFLICT (ai_domain, product_type, category_id) DO UPDATE SET weight = EXCLUDED.weight;
    END IF;

    IF v_puzzles_id IS NOT NULL THEN
        INSERT INTO category_ai_mappings (ai_domain, product_type, category_id, weight, is_active) VALUES
        ('entertainment', 'puzzle', v_puzzles_id, 0.99, true),
        ('entertainment', 'jigsaw', v_puzzles_id, 0.99, true)
        ON CONFLICT (ai_domain, product_type, category_id) DO UPDATE SET weight = EXCLUDED.weight;
    END IF;

    IF v_board_games_id IS NOT NULL THEN
        INSERT INTO category_ai_mappings (ai_domain, product_type, category_id, weight, is_active) VALUES
        ('entertainment', 'game', v_board_games_id, 0.85, true),
        ('entertainment', 'board_game', v_board_games_id, 0.95, true),
        ('entertainment', 'card_game', v_board_games_id, 0.90, true)
        ON CONFLICT (ai_domain, product_type, category_id) DO UPDATE SET weight = EXCLUDED.weight;
    END IF;

    -- Маппинги для природных материалов
    IF v_natural_materials_id IS NOT NULL THEN
        INSERT INTO category_ai_mappings (ai_domain, product_type, category_id, weight, is_active) VALUES
        ('nature', 'acorn', v_natural_materials_id, 0.90, true),
        ('nature', 'wood', v_natural_materials_id, 0.85, true),
        ('nature', 'stone', v_natural_materials_id, 0.85, true),
        ('nature', 'mineral', v_natural_materials_id, 0.85, true),
        ('nature', 'seed', v_natural_materials_id, 0.85, true),
        ('nature', 'plant', v_natural_materials_id, 0.85, true)
        ON CONFLICT (ai_domain, product_type, category_id) DO UPDATE SET weight = EXCLUDED.weight;
    END IF;

    -- Маппинги для строительных материалов
    IF v_construction_materials_id IS NOT NULL THEN
        INSERT INTO category_ai_mappings (ai_domain, product_type, category_id, weight, is_active) VALUES
        ('construction', 'sand', v_construction_materials_id, 0.90, true),
        ('construction', 'cement', v_construction_materials_id, 0.95, true),
        ('construction', 'brick', v_construction_materials_id, 0.95, true),
        ('construction', 'tool', v_construction_materials_id, 0.85, true),
        ('construction', 'paint', v_construction_materials_id, 0.85, true)
        ON CONFLICT (ai_domain, product_type, category_id) DO UPDATE SET weight = EXCLUDED.weight;
    END IF;

    -- Маппинги для специальных категорий
    IF v_aviation_id IS NOT NULL THEN
        INSERT INTO category_ai_mappings (ai_domain, product_type, category_id, weight, is_active) VALUES
        ('aviation', 'airplane', v_aviation_id, 0.95, true),
        ('aviation', 'aircraft', v_aviation_id, 0.95, true),
        ('aviation', 'helicopter', v_aviation_id, 0.95, true)
        ON CONFLICT (ai_domain, product_type, category_id) DO UPDATE SET weight = EXCLUDED.weight;
    END IF;

    IF v_military_id IS NOT NULL THEN
        INSERT INTO category_ai_mappings (ai_domain, product_type, category_id, weight, is_active) VALUES
        ('military', 'uniform', v_military_id, 0.90, true),
        ('military', 'equipment', v_military_id, 0.85, true),
        ('military', 'medal', v_military_id, 0.90, true)
        ON CONFLICT (ai_domain, product_type, category_id) DO UPDATE SET weight = EXCLUDED.weight;
    END IF;

    -- Fallback маппинг для неопределенных товаров
    IF v_misc_id IS NOT NULL THEN
        INSERT INTO category_ai_mappings (ai_domain, product_type, category_id, weight, is_active) VALUES
        ('other', 'unknown', v_misc_id, 0.50, true),
        ('other', 'misc', v_misc_id, 0.50, true)
        ON CONFLICT (ai_domain, product_type, category_id) DO UPDATE SET weight = EXCLUDED.weight;
    END IF;
END $$;

-- ===================================================================
-- 4. ДОБАВЛЕНИЕ КЛЮЧЕВЫХ СЛОВ ДЛЯ УЛУЧШЕНИЯ ДЕТЕКЦИИ
-- ===================================================================

-- Функция для добавления ключевых слов
CREATE OR REPLACE FUNCTION add_category_keywords(
    p_category_slug TEXT,
    p_keywords TEXT[]
) RETURNS VOID AS $$
DECLARE
    v_category_id INTEGER;
    v_keyword TEXT;
BEGIN
    SELECT id INTO v_category_id FROM marketplace_categories WHERE slug = p_category_slug;

    IF v_category_id IS NOT NULL THEN
        FOREACH v_keyword IN ARRAY p_keywords
        LOOP
            INSERT INTO category_keyword_weights (
                category_id, keyword, language, weight, keyword_type, success_rate
            ) VALUES (
                v_category_id, LOWER(v_keyword), 'ru', 0.8, 'primary', 0.5
            ) ON CONFLICT (category_id, keyword, language) DO UPDATE
            SET weight = GREATEST(category_keyword_weights.weight, EXCLUDED.weight);
        END LOOP;
    END IF;
END;
$$ LANGUAGE plpgsql;

-- Добавляем ключевые слова для категорий
SELECT add_category_keywords('puzzles', ARRAY['пазл', 'пазлы', 'головоломка', 'собирать', 'картинка']);
SELECT add_category_keywords('toys', ARRAY['игрушка', 'игрушки', 'кукла', 'машинка', 'мягкая']);
SELECT add_category_keywords('natural-materials', ARRAY['желудь', 'шишка', 'ракушка', 'камень', 'дерево']);
SELECT add_category_keywords('construction-materials', ARRAY['песок', 'цемент', 'кирпич', 'плитка', 'гипс']);
SELECT add_category_keywords('aviation', ARRAY['самолет', 'вертолет', 'авиация', 'крыло', 'пропеллер']);
SELECT add_category_keywords('military', ARRAY['военный', 'армейский', 'форма', 'медаль', 'знак']);

-- Удаляем временную функцию
DROP FUNCTION IF EXISTS add_category_translations(TEXT, TEXT, TEXT, TEXT);
DROP FUNCTION IF EXISTS add_category_keywords(TEXT, TEXT[]);

-- ===================================================================
-- 5. ОБНОВЛЕНИЕ СТАТУСА КАТЕГОРИЙ
-- ===================================================================

-- Активируем все новые категории
UPDATE marketplace_categories
SET is_active = true
WHERE slug IN (
    'natural-materials', 'wood-materials', 'stones-minerals', 'plants-seeds', 'natural-decor',
    'construction-materials', 'bulk-materials', 'tools', 'paints', 'plumbing', 'electrical',
    'crafts', 'craft-materials', 'handmade', 'sewing', 'knitting',
    'toys', 'puzzles', 'board-games', 'collectibles', 'constructors', 'educational-games', 'models',
    'antiques', 'aviation', 'military', 'miscellaneous'
);

-- Обновляем счетчики уровней для правильной иерархии
UPDATE marketplace_categories
SET level = 0
WHERE parent_id IS NULL;

UPDATE marketplace_categories
SET level = 1
WHERE parent_id IN (SELECT id FROM marketplace_categories WHERE parent_id IS NULL);

UPDATE marketplace_categories
SET level = 2
WHERE parent_id IN (SELECT id FROM marketplace_categories WHERE level = 1);