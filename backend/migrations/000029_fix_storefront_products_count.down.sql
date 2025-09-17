-- Удаляем триггер
DROP TRIGGER IF EXISTS trigger_update_storefront_products_count ON storefront_products;

-- Удаляем функцию
DROP FUNCTION IF EXISTS update_storefront_products_count();

-- Сбрасываем счётчики в 0 (опционально)
UPDATE storefronts SET products_count = 0;

-- Удаляем комментарий
COMMENT ON COLUMN storefronts.products_count IS NULL;