-- Удаляем триггер и функцию
DROP TRIGGER IF EXISTS sync_storefront_products_to_marketplace ON storefront_products;
DROP FUNCTION IF EXISTS sync_storefront_to_marketplace();

-- Удаляем записи из marketplace_listings, которые относятся к товарам витрин
DELETE FROM marketplace_listings 
WHERE storefront_id IS NOT NULL;

-- Логирование
DO $$
BEGIN
    RAISE NOTICE 'Removed storefront sync trigger and cleaned up marketplace_listings';
END $$;