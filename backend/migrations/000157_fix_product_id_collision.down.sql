-- Откат к раздельным последовательностям

-- 1. Создаем обратно старые последовательности
CREATE SEQUENCE IF NOT EXISTS marketplace_listings_id_seq;
CREATE SEQUENCE IF NOT EXISTS storefront_products_id_seq;

-- 2. Устанавливаем текущие значения для последовательностей
SELECT setval('marketplace_listings_id_seq', 
    COALESCE((SELECT MAX(id) FROM marketplace_listings), 1)
);
SELECT setval('storefront_products_id_seq',
    COALESCE((SELECT MAX(id) FROM storefront_products), 1)
);

-- 3. Изменяем дефолтные значения обратно на старые последовательности
ALTER TABLE marketplace_listings 
    ALTER COLUMN id SET DEFAULT nextval('marketplace_listings_id_seq');
ALTER TABLE storefront_products
    ALTER COLUMN id SET DEFAULT nextval('storefront_products_id_seq');

-- 4. Удаляем общую последовательность
DROP SEQUENCE IF EXISTS global_product_id_seq CASCADE;

-- 5. Логирование
DO $$
BEGIN
    RAISE NOTICE 'Reverted to separate ID sequences for marketplace_listings and storefront_products';
END $$;