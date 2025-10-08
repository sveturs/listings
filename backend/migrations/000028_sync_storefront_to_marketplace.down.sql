-- Migration: 000028_sync_storefront_to_marketplace (ROLLBACK)
-- Description: Удаляет триггер синхронизации товаров из storefront_products в marketplace_listings
-- Author: Claude Code
-- Date: 2025-10-06

-- Удалить триггер
DROP TRIGGER IF EXISTS trigger_sync_storefront_to_marketplace ON storefront_products;

-- Удалить функцию
DROP FUNCTION IF EXISTS sync_storefront_product_to_marketplace();

-- Примечание: listings, созданные через синхронизацию, НЕ удаляются,
-- т.к. могут иметь связанные данные (отзывы, чаты, избранное и т.д.)
-- Если нужно удалить все синхронизированные listings, выполнить:
-- DELETE FROM marketplace_listings WHERE metadata->>'source' = 'storefront';
