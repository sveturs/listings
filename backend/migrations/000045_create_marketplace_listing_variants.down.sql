-- 000045_create_marketplace_listing_variants.down.sql

-- Удаление триггера
DROP TRIGGER IF EXISTS marketplace_listing_variants_updated_at ON marketplace_listing_variants;

-- Удаление функции триггера
DROP FUNCTION IF EXISTS update_marketplace_listing_variants_updated_at();

-- Удаление таблицы вариантов
DROP TABLE IF EXISTS marketplace_listing_variants;