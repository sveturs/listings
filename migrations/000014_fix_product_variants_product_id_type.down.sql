-- Reverse migration

ALTER TABLE product_variants
DROP CONSTRAINT IF EXISTS fk_product_variants_listing;

DROP INDEX IF EXISTS idx_variants_available;
DROP INDEX IF EXISTS idx_variants_default;
DROP INDEX IF EXISTS idx_variants_product;
DROP INDEX IF EXISTS idx_variants_stock_status;

ALTER TABLE product_variants 
ALTER COLUMN product_id TYPE uuid USING gen_random_uuid();

-- Recreate original indexes
CREATE INDEX idx_variants_product ON product_variants(product_id, position);
CREATE INDEX idx_variants_default ON product_variants(product_id) WHERE is_default = true;
CREATE INDEX idx_variants_stock_status ON product_variants(product_id, status, stock_quantity) WHERE status = 'active';
CREATE INDEX idx_variants_available ON product_variants(product_id, (stock_quantity - reserved_quantity)) 
  WHERE status::text = 'active'::text AND (stock_quantity - reserved_quantity) > 0;
