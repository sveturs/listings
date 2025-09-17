-- Удаляем триггер
DROP TRIGGER IF EXISTS update_storefront_views_count_trigger ON storefront_products;

-- Удаляем функцию
DROP FUNCTION IF EXISTS update_storefront_views_count();

-- Удаляем индекс
DROP INDEX IF EXISTS idx_storefront_products_storefront_id_view_count;

-- Обнуляем счетчики просмотров (опционально)
UPDATE storefronts SET views_count = 0;