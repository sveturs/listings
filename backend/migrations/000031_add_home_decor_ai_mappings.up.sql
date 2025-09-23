-- Добавляем маппинги для категории "Дом и сад" (home-garden)
-- Категория ID 1005 = "Dom i bašta" (home-garden)

-- Вазы и декоративные предметы
INSERT INTO category_ai_mappings (ai_domain, product_type, category_id, weight) VALUES
('home decor', 'vase', 1005, 0.95),
('home decor', 'ceramic', 1005, 0.90),
('home decor', 'decorative item', 1005, 0.90),
('home decor', 'ornament', 1005, 0.85),
('home decor', 'sculpture', 1005, 0.85),
('home decor', 'figurine', 1005, 0.85),
('home decor', 'planter', 1005, 0.85),
('home decor', 'pot', 1005, 0.85),
('home decor', 'candle holder', 1005, 0.85),
('home decor', 'wall art', 1005, 0.85)
ON CONFLICT (ai_domain, product_type, category_id) DO UPDATE
SET weight = EXCLUDED.weight,
    updated_at = NOW();

-- Домашний декор общий
INSERT INTO category_ai_mappings (ai_domain, product_type, category_id, weight) VALUES
('home', 'vase', 1005, 0.90),
('home', 'decoration', 1005, 0.85),
('home', 'ornament', 1005, 0.85),
('home', 'ceramic', 1005, 0.85),
('interior', 'decoration', 1005, 0.90),
('interior', 'vase', 1005, 0.90),
('interior', 'ornament', 1005, 0.85)
ON CONFLICT (ai_domain, product_type, category_id) DO UPDATE
SET weight = EXCLUDED.weight,
    updated_at = NOW();

-- Садовый декор
INSERT INTO category_ai_mappings (ai_domain, product_type, category_id, weight) VALUES
('garden', 'decoration', 1005, 0.90),
('garden', 'ornament', 1005, 0.85),
('garden', 'planter', 1005, 0.90),
('garden', 'pot', 1005, 0.90),
('outdoor', 'decoration', 1005, 0.85),
('outdoor', 'furniture', 1005, 0.80)
ON CONFLICT (ai_domain, product_type, category_id) DO UPDATE
SET weight = EXCLUDED.weight,
    updated_at = NOW();

-- Мебель и предметы интерьера
INSERT INTO category_ai_mappings (ai_domain, product_type, category_id, weight) VALUES
('furniture', 'chair', 1005, 0.90),
('furniture', 'table', 1005, 0.90),
('furniture', 'shelf', 1005, 0.90),
('furniture', 'cabinet', 1005, 0.90),
('furniture', 'sofa', 1005, 0.90),
('furniture', 'bed', 1005, 0.90),
('furniture', 'wardrobe', 1005, 0.90),
('furniture', 'desk', 1005, 0.90)
ON CONFLICT (ai_domain, product_type, category_id) DO UPDATE
SET weight = EXCLUDED.weight,
    updated_at = NOW();

-- Текстиль для дома
INSERT INTO category_ai_mappings (ai_domain, product_type, category_id, weight) VALUES
('home textiles', 'curtains', 1005, 0.90),
('home textiles', 'bedding', 1005, 0.90),
('home textiles', 'towels', 1005, 0.90),
('home textiles', 'carpet', 1005, 0.90),
('home textiles', 'rug', 1005, 0.90),
('home textiles', 'pillow', 1005, 0.85),
('home textiles', 'blanket', 1005, 0.85)
ON CONFLICT (ai_domain, product_type, category_id) DO UPDATE
SET weight = EXCLUDED.weight,
    updated_at = NOW();

-- Кухонные принадлежности
INSERT INTO category_ai_mappings (ai_domain, product_type, category_id, weight) VALUES
('kitchen', 'utensils', 1005, 0.85),
('kitchen', 'cookware', 1005, 0.85),
('kitchen', 'dishes', 1005, 0.85),
('kitchen', 'appliance', 1005, 0.80),
('kitchen', 'container', 1005, 0.80)
ON CONFLICT (ai_domain, product_type, category_id) DO UPDATE
SET weight = EXCLUDED.weight,
    updated_at = NOW();

-- Освещение
INSERT INTO category_ai_mappings (ai_domain, product_type, category_id, weight) VALUES
('lighting', 'lamp', 1005, 0.90),
('lighting', 'chandelier', 1005, 0.90),
('lighting', 'light fixture', 1005, 0.90),
('lighting', 'led', 1005, 0.85),
('lighting', 'bulb', 1005, 0.80)
ON CONFLICT (ai_domain, product_type, category_id) DO UPDATE
SET weight = EXCLUDED.weight,
    updated_at = NOW();

-- Обновляем статистику для новых маппингов
UPDATE category_ai_mappings
SET success_count = 10,
    failure_count = 1
WHERE ai_domain IN ('home decor', 'home', 'interior', 'garden', 'furniture', 'home textiles', 'kitchen', 'lighting')
  AND category_id = 1005
  AND success_count = 0;