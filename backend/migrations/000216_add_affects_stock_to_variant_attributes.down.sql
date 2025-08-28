-- Удаляем индекс
DROP INDEX IF EXISTS idx_product_variant_attributes_affects_stock;

-- Удаляем поле affects_stock из product_variant_attributes
ALTER TABLE product_variant_attributes 
DROP COLUMN IF EXISTS affects_stock;