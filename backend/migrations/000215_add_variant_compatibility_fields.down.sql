-- Удаляем индекс
DROP INDEX IF EXISTS idx_category_attributes_variant_compatible;

-- Удаляем поле is_variant_compatible из category_attributes
ALTER TABLE category_attributes 
DROP COLUMN IF EXISTS is_variant_compatible;