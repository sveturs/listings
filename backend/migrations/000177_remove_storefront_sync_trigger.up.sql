-- Migration: 000177_remove_storefront_sync_trigger
-- Description: Удаляет триггер sync_storefront_to_marketplace и дубликаты B2C товаров из c2c_listings
-- Reason: Костыль с дублированием данных приводит к багам (отсутствующие изображения)
-- New approach: Unified listings через UNION запросы без дублирования в БД
-- Date: 2025-10-11

BEGIN;

-- Шаг 1: Удалить триггер (если существует)
DROP TRIGGER IF EXISTS trigger_sync_storefront_to_marketplace ON b2c_products;

-- Шаг 2: Удалить функцию триггера (если существует)
DROP FUNCTION IF EXISTS sync_storefront_product_to_marketplace();

-- Шаг 3: Удалить дубликаты B2C товаров из c2c_listings
-- Это товары, которые были автоматически скопированы из b2c_products
-- Определяем их по наличию storefront_id
DELETE FROM c2c_listings
WHERE storefront_id IS NOT NULL
  AND metadata->>'source' = 'storefront';

-- Шаг 4: Комментарий о новом подходе
COMMENT ON TABLE c2c_listings IS
'C2C listings (customer-to-customer объявления).
ВАЖНО: С миграции 000177 больше не содержит дубликаты B2C товаров.
Для получения unified списка используйте UNION запросы с b2c_products.';

COMMENT ON TABLE b2c_products IS
'B2C products (business-to-customer товары витрин).
ВАЖНО: С миграции 000177 больше не синхронизируются в c2c_listings.
Для получения unified списка используйте UNION запросы с c2c_listings.';

COMMIT;
