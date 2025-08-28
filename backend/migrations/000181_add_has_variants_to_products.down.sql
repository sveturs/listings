-- Удаляем индекс
DROP INDEX IF EXISTS idx_storefront_products_has_variants;

-- Удаляем колонку has_variants
ALTER TABLE storefront_products
DROP COLUMN IF EXISTS has_variants;