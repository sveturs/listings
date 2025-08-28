-- Удалить триггер и функцию
DROP TRIGGER IF EXISTS trigger_update_variant_stock_status ON storefront_product_variants;
DROP FUNCTION IF EXISTS update_variant_stock_status();

-- Удалить индексы
DROP INDEX IF EXISTS idx_storefront_product_variants_low_stock;
DROP INDEX IF EXISTS idx_storefront_product_variants_stock_status;
DROP INDEX IF EXISTS idx_storefront_product_variants_unique_attributes;

-- Удалить добавленные колонки
ALTER TABLE storefront_product_variants 
DROP COLUMN IF EXISTS available_quantity,
DROP COLUMN IF EXISTS reserved_quantity;