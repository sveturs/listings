-- Исправляем коллизию ID между marketplace_listings и storefront_products
-- Создаем общую последовательность для всех продуктов

-- 1. Находим максимальный ID среди всех продуктов
DO $$
DECLARE
    max_marketplace_id INTEGER;
    max_storefront_id INTEGER;
    max_id INTEGER;
BEGIN
    SELECT COALESCE(MAX(id), 0) INTO max_marketplace_id FROM marketplace_listings;
    SELECT COALESCE(MAX(id), 0) INTO max_storefront_id FROM storefront_products;
    
    max_id := GREATEST(max_marketplace_id, max_storefront_id);
    
    RAISE NOTICE 'Max marketplace ID: %, Max storefront ID: %, Starting new sequence from: %', 
        max_marketplace_id, max_storefront_id, max_id + 1;
END $$;

-- 2. Создаем новую общую последовательность
CREATE SEQUENCE IF NOT EXISTS global_product_id_seq;

-- 3. Устанавливаем начальное значение последовательности
-- Берем максимум из обеих таблиц и добавляем 1
SELECT setval('global_product_id_seq', 
    GREATEST(
        COALESCE((SELECT MAX(id) FROM marketplace_listings), 0),
        COALESCE((SELECT MAX(id) FROM storefront_products), 0)
    ) + 1
);

-- 4. Изменяем дефолтное значение для marketplace_listings
ALTER TABLE marketplace_listings 
    ALTER COLUMN id SET DEFAULT nextval('global_product_id_seq');

-- 5. Изменяем дефолтное значение для storefront_products  
ALTER TABLE storefront_products
    ALTER COLUMN id SET DEFAULT nextval('global_product_id_seq');

-- 6. Удаляем старые последовательности (они больше не нужны)
-- Но сначала убедимся, что они не используются
ALTER TABLE marketplace_listings 
    ALTER COLUMN id DROP DEFAULT;
ALTER TABLE storefront_products
    ALTER COLUMN id DROP DEFAULT;

-- Теперь можно безопасно удалить старые последовательности
DROP SEQUENCE IF EXISTS marketplace_listings_id_seq CASCADE;
DROP SEQUENCE IF EXISTS storefront_products_id_seq CASCADE;

-- 7. Восстанавливаем дефолтные значения с новой последовательностью
ALTER TABLE marketplace_listings 
    ALTER COLUMN id SET DEFAULT nextval('global_product_id_seq');
ALTER TABLE storefront_products
    ALTER COLUMN id SET DEFAULT nextval('global_product_id_seq');

-- 8. Логирование
DO $$
BEGIN
    RAISE NOTICE 'Successfully created global product ID sequence. All new products will have unique IDs across both tables.';
END $$;