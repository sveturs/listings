-- Добавление недостающих категорий для высокой точности AI определения

-- Получаем максимальный ID для новых категорий
DO $$
DECLARE
    max_id INTEGER;
BEGIN
    SELECT COALESCE(MAX(id), 2000) INTO max_id FROM marketplace_categories;

    -- Сбросим sequence если нужно
    PERFORM setval('marketplace_categories_id_seq', max_id + 1, false);
END $$;

-- Игрушки и хобби (подкатегории для 1015)
INSERT INTO marketplace_categories (id, name, slug, parent_id, icon, is_active, created_at) VALUES
(2001, 'Игрушки', 'toys', 1015, 'toy-brick', true, NOW()),
(2002, 'Пазлы', 'puzzles', 1015, 'puzzle-piece', true, NOW()),
(2003, 'Настольные игры', 'board-games', 1015, 'dice', true, NOW()),
(2004, 'Коллекционирование', 'collectibles', 1015, 'star', true, NOW()),
(2005, 'Конструкторы', 'constructors', 1015, 'cube', true, NOW()),
(2006, 'Развивающие игры', 'educational-games', 1015, 'book-open', true, NOW()),
(2007, 'Модели и сборка', 'models', 1015, 'wrench', true, NOW())
ON CONFLICT (id) DO NOTHING;

-- Строительные материалы (подкатегории для 1007)
INSERT INTO marketplace_categories (id, name, slug, parent_id, icon, is_active, created_at) VALUES
(2010, 'Строительные материалы', 'construction-materials', 1007, 'hammer', true, NOW()),
(2011, 'Сыпучие материалы', 'bulk-materials', 1007, 'cube', true, NOW()),
(2012, 'Инструменты', 'tools', 1007, 'wrench', true, NOW()),
(2013, 'Краски и лаки', 'paints', 1007, 'palette', true, NOW()),
(2014, 'Сантехника', 'plumbing', 1007, 'droplet', true, NOW()),
(2015, 'Электрика', 'electrical', 1007, 'zap', true, NOW())
ON CONFLICT (id) DO NOTHING;

-- Природные материалы (новая основная категория)
INSERT INTO marketplace_categories (id, name, slug, parent_id, icon, is_active, created_at) VALUES
(2020, 'Природные материалы', 'natural-materials', NULL, 'tree', true, NOW()),
(2021, 'Дерево и пиломатериалы', 'wood-materials', 2020, 'tree', true, NOW()),
(2022, 'Камни и минералы', 'stones-minerals', 2020, 'gem', true, NOW()),
(2023, 'Растения и семена', 'plants-seeds', 2020, 'flower', true, NOW()),
(2024, 'Природный декор', 'natural-decor', 2020, 'leaf', true, NOW())
ON CONFLICT (id) DO NOTHING;

-- Декор и рукоделие (новая основная категория)
INSERT INTO marketplace_categories (id, name, slug, parent_id, icon, is_active, created_at) VALUES
(2030, 'Декор и рукоделие', 'crafts', NULL, 'palette', true, NOW()),
(2031, 'Материалы для творчества', 'craft-materials', 2030, 'brush', true, NOW()),
(2032, 'Готовые изделия', 'handmade', 2030, 'heart', true, NOW()),
(2033, 'Швейные принадлежности', 'sewing', 2030, 'scissors', true, NOW()),
(2034, 'Вязание и пряжа', 'knitting', 2030, 'circle', true, NOW())
ON CONFLICT (id) DO NOTHING;

-- Специальные категории
INSERT INTO marketplace_categories (id, name, slug, parent_id, icon, is_active, created_at) VALUES
(2040, 'Антиквариат', 'antiques', NULL, 'clock', true, NOW()),
(2041, 'Винтаж', 'vintage', 2040, 'archive', true, NOW()),
(2042, 'Монеты и банкноты', 'coins', 2040, 'dollar-sign', true, NOW()),
(2043, 'Марки', 'stamps', 2040, 'mail', true, NOW()),
(2044, 'Искусство', 'art', 2040, 'image', true, NOW())
ON CONFLICT (id) DO NOTHING;

-- Авиация и военные товары
INSERT INTO marketplace_categories (id, name, slug, parent_id, icon, is_active, created_at) VALUES
(2050, 'Авиация', 'aviation', NULL, 'plane', true, NOW()),
(2051, 'Модели самолетов', 'aircraft-models', 2050, 'plane', true, NOW()),
(2052, 'Авиационные запчасти', 'aircraft-parts', 2050, 'settings', true, NOW()),
(2053, 'Военные товары', 'military', NULL, 'shield', true, NOW()),
(2054, 'Военная форма', 'military-uniform', 2053, 'shirt', true, NOW()),
(2055, 'Военное снаряжение', 'military-equipment', 2053, 'shield', true, NOW())
ON CONFLICT (id) DO NOTHING;

-- Животные и зоотовары (расширение существующей категории 1012)
INSERT INTO marketplace_categories (id, name, slug, parent_id, icon, is_active, created_at) VALUES
(2060, 'Корма', 'pet-food', 1012, 'package', true, NOW()),
(2061, 'Аксессуары для животных', 'pet-accessories', 1012, 'heart', true, NOW()),
(2062, 'Аквариумистика', 'aquarium', 1012, 'droplet', true, NOW()),
(2063, 'Ветеринарные товары', 'veterinary', 1012, 'activity', true, NOW())
ON CONFLICT (id) DO NOTHING;

-- Добавляем переводы для новых категорий
INSERT INTO translations (entity_type, entity_id, language, field_name, field_value, created_at) VALUES
-- Игрушки и хобби
('category', 2001, 'ru', 'name', 'Игрушки', NOW()),
('category', 2001, 'en', 'name', 'Toys', NOW()),
('category', 2001, 'sr', 'name', 'Igračke', NOW()),

('category', 2002, 'ru', 'name', 'Пазлы', NOW()),
('category', 2002, 'en', 'name', 'Puzzles', NOW()),
('category', 2002, 'sr', 'name', 'Slagalice', NOW()),

('category', 2003, 'ru', 'name', 'Настольные игры', NOW()),
('category', 2003, 'en', 'name', 'Board Games', NOW()),
('category', 2003, 'sr', 'name', 'Društvene igre', NOW()),

('category', 2004, 'ru', 'name', 'Коллекционирование', NOW()),
('category', 2004, 'en', 'name', 'Collectibles', NOW()),
('category', 2004, 'sr', 'name', 'Kolekcionarstvo', NOW()),

-- Строительные материалы
('category', 2010, 'ru', 'name', 'Строительные материалы', NOW()),
('category', 2010, 'en', 'name', 'Construction Materials', NOW()),
('category', 2010, 'sr', 'name', 'Građevinski materijal', NOW()),

('category', 2011, 'ru', 'name', 'Сыпучие материалы', NOW()),
('category', 2011, 'en', 'name', 'Bulk Materials', NOW()),
('category', 2011, 'sr', 'name', 'Rasuti materijal', NOW()),

-- Природные материалы
('category', 2020, 'ru', 'name', 'Природные материалы', NOW()),
('category', 2020, 'en', 'name', 'Natural Materials', NOW()),
('category', 2020, 'sr', 'name', 'Prirodni materijali', NOW()),

-- Декор и рукоделие
('category', 2030, 'ru', 'name', 'Декор и рукоделие', NOW()),
('category', 2030, 'en', 'name', 'Crafts and Decor', NOW()),
('category', 2030, 'sr', 'name', 'Rukotvorine i dekor', NOW()),

-- Антиквариат
('category', 2040, 'ru', 'name', 'Антиквариат', NOW()),
('category', 2040, 'en', 'name', 'Antiques', NOW()),
('category', 2040, 'sr', 'name', 'Antikviteti', NOW()),

-- Авиация
('category', 2050, 'ru', 'name', 'Авиация', NOW()),
('category', 2050, 'en', 'name', 'Aviation', NOW()),
('category', 2050, 'sr', 'name', 'Avijacija', NOW()),

-- Военные товары
('category', 2053, 'ru', 'name', 'Военные товары', NOW()),
('category', 2053, 'en', 'name', 'Military Goods', NOW()),
('category', 2053, 'sr', 'name', 'Vojna oprema', NOW())
ON CONFLICT (entity_type, entity_id, language, field_name) DO NOTHING;

-- Обновляем AI маппинги для новых категорий
INSERT INTO category_ai_mappings (ai_domain, product_type, category_id, weight) VALUES
-- Пазлы и игрушки
('entertainment', 'puzzle', 2002, 1.00),
('entertainment', 'toy', 2001, 1.00),
('entertainment', 'board-game', 2003, 1.00),
('entertainment', 'collectible', 2004, 1.00),
('entertainment', 'constructor', 2005, 1.00),

-- Строительные материалы
('construction', 'sand', 2011, 1.00),
('construction', 'gravel', 2011, 1.00),
('construction', 'cement', 2010, 1.00),
('construction', 'bricks', 2010, 1.00),
('construction', 'tools', 2012, 1.00),

-- Природные материалы
('nature', 'acorn', 2024, 1.00),
('nature', 'wood', 2021, 1.00),
('nature', 'stones', 2022, 1.00),
('nature', 'plants', 2023, 1.00),
('nature', 'natural-decor', 2024, 1.00),

-- Антиквариат
('antiques', 'vintage-item', 2041, 1.00),
('antiques', 'coins', 2042, 1.00),
('antiques', 'stamps', 2043, 1.00),
('antiques', 'art', 2044, 1.00),

-- Авиация
('aviation', 'aircraft', 2050, 1.00),
('aviation', 'plane-model', 2051, 1.00),
('aviation', 'aircraft-parts', 2052, 1.00),

-- Военные товары
('military', 'uniform', 2054, 1.00),
('military', 'equipment', 2055, 1.00)
ON CONFLICT (ai_domain, product_type, category_id) DO NOTHING;