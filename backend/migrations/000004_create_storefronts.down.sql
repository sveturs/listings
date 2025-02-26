-- Убираем индекс и внешний ключ
DROP INDEX IF EXISTS idx_marketplace_listings_storefront;
ALTER TABLE marketplace_listings DROP COLUMN IF EXISTS storefront_id;

-- Удаляем таблицы в обратном порядке
DROP TABLE IF EXISTS import_history;
DROP TABLE IF EXISTS import_sources;
DROP TABLE IF EXISTS user_storefronts;