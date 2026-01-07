-- Reverse migration - remove added columns

ALTER TABLE product_variants
DROP COLUMN IF EXISTS cost_price,
DROP COLUMN IF EXISTS stock_status,
DROP COLUMN IF EXISTS low_stock_threshold,
DROP COLUMN IF EXISTS variant_attributes,
DROP COLUMN IF EXISTS weight,
DROP COLUMN IF NOT EXISTS dimensions,
DROP COLUMN IF EXISTS is_active,
DROP COLUMN IF EXISTS view_count,
DROP COLUMN IF EXISTS sold_count;

DROP INDEX IF EXISTS idx_product_variants_stock_status;
DROP INDEX IF EXISTS idx_product_variants_is_active;
