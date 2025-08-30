-- Миграция для исправления category_id в marketplace_listings
-- Исправляем нулевые category_id на основе названий товаров и категорий

-- 1. Роутеры, телефоны, электроника -> Electronics/TV & Audio
UPDATE marketplace_listings 
SET category_id = (SELECT id FROM marketplace_categories WHERE slug = 'tv-audio' LIMIT 1)
WHERE category_id = 0 
AND (title ILIKE '%роутер%' OR title ILIKE '%телефон%' OR title ILIKE '%router%' OR title ILIKE '%huawei%');

-- 2. Дома, квартиры, недвижимость -> Real Estate (найдем подходящую)
UPDATE marketplace_listings 
SET category_id = (SELECT id FROM marketplace_categories WHERE name ILIKE '%real estate%' OR name ILIKE '%недвижимость%' LIMIT 1)
WHERE category_id = 0 
AND (title ILIKE '%дом%' OR title ILIKE '%квартира%' OR title ILIKE '%пентхаус%' OR title ILIKE '%house%');

-- 3. Если нет подходящей категории недвижимости, создаем новую
INSERT INTO marketplace_categories (name, slug, icon, parent_id, sort_order, created_at, updated_at)
SELECT 'Real Estate', 'real-estate', 'home', NULL, 100, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM marketplace_categories WHERE slug = 'real-estate');

-- Обновляем недвижимость на созданную категорию
UPDATE marketplace_listings 
SET category_id = (SELECT id FROM marketplace_categories WHERE slug = 'real-estate' LIMIT 1)
WHERE category_id = 0 
AND (title ILIKE '%дом%' OR title ILIKE '%квартира%' OR title ILIKE '%пентхаус%' OR title ILIKE '%house%');

-- 4. Мёд, сельхоз продукты -> Agricultural 
UPDATE marketplace_listings 
SET category_id = (SELECT id FROM marketplace_categories WHERE slug = 'farm-machinery' LIMIT 1)
WHERE category_id = 0 
AND (title ILIKE '%мёд%' OR title ILIKE '%honey%' OR title ILIKE '%фермер%');

-- 5. Автомобили -> Cars
UPDATE marketplace_listings 
SET category_id = (SELECT id FROM marketplace_categories WHERE name ILIKE '%car%' OR name ILIKE '%auto%' OR slug LIKE '%car%' LIMIT 1)
WHERE category_id = 0 
AND (title ILIKE '%volkswagen%' OR title ILIKE '%golf%' OR title ILIKE '%mercedes%' OR title ILIKE '%bmw%' OR title ILIKE '%yamaha%');

-- 6. Если категория Cars не найдена, создаем
INSERT INTO marketplace_categories (name, slug, icon, parent_id, sort_order, created_at, updated_at)
SELECT 'Cars', 'cars', 'car', NULL, 110, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM marketplace_categories WHERE slug = 'cars');

-- Обновляем автомобили
UPDATE marketplace_listings 
SET category_id = (SELECT id FROM marketplace_categories WHERE slug = 'cars' LIMIT 1)
WHERE category_id = 0 
AND (title ILIKE '%volkswagen%' OR title ILIKE '%golf%' OR title ILIKE '%mercedes%' OR title ILIKE '%bmw%' OR title ILIKE '%yamaha%');

-- 7. Услуги строительства, web-разработка -> Services
INSERT INTO marketplace_categories (name, slug, icon, parent_id, sort_order, created_at, updated_at)
SELECT 'Services', 'services', 'briefcase', NULL, 120, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM marketplace_categories WHERE slug = 'services');

UPDATE marketplace_listings 
SET category_id = (SELECT id FROM marketplace_categories WHERE slug = 'services' LIMIT 1)
WHERE category_id = 0 
AND (title ILIKE '%web%' OR title ILIKE '%сайт%' OR title ILIKE '%renoviranje%' OR title ILIKE '%строительство%');

-- 8. Для всех остальных товаров с category_id = 0 ставим категорию "Other"
INSERT INTO marketplace_categories (name, slug, icon, parent_id, sort_order, created_at, updated_at)
SELECT 'Other', 'other', 'tag', NULL, 999, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM marketplace_categories WHERE slug = 'other');

-- Обновляем остальные товары
UPDATE marketplace_listings 
SET category_id = (SELECT id FROM marketplace_categories WHERE slug = 'other' LIMIT 1)
WHERE category_id = 0;

-- Логируем результаты
DO $$
DECLARE
    total_fixed INTEGER;
BEGIN
    SELECT COUNT(*) INTO total_fixed FROM marketplace_listings WHERE category_id != 0;
    RAISE NOTICE 'Fixed category_id for % listings', total_fixed;
END $$;