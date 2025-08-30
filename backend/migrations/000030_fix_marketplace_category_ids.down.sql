-- Откат миграции - возвращаем category_id в 0 для всех обновленных записей
-- и удаляем созданные категории

-- 1. Возвращаем category_id в 0 для всех товаров, которые были обновлены
UPDATE marketplace_listings 
SET category_id = 0
WHERE category_id IN (
    SELECT id FROM marketplace_categories 
    WHERE slug IN ('real-estate', 'cars', 'services', 'other')
);

-- 2. Удаляем созданные категории
DELETE FROM marketplace_categories 
WHERE slug IN ('real-estate', 'cars', 'services', 'other');

-- Логируем откат
DO $$
DECLARE
    total_reverted INTEGER;
BEGIN
    SELECT COUNT(*) INTO total_reverted FROM marketplace_listings WHERE category_id = 0;
    RAISE NOTICE 'Reverted % listings to category_id = 0', total_reverted;
END $$;