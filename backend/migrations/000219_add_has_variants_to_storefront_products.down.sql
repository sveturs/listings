-- Remove index
DROP INDEX IF EXISTS idx_storefront_products_has_variants;

-- Remove has_variants column
ALTER TABLE storefront_products 
DROP COLUMN IF EXISTS has_variants;