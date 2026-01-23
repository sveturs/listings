-- Fix product_variants.product_id to use bigint instead of UUID to match listings.id

-- Drop indexes that reference product_id
DROP INDEX IF EXISTS idx_variants_available;
DROP INDEX IF EXISTS idx_variants_default;
DROP INDEX IF EXISTS idx_variants_product;
DROP INDEX IF EXISTS idx_variants_stock_status;

-- Change product_id from UUID to bigint
ALTER TABLE product_variants 
ALTER COLUMN product_id TYPE bigint USING product_id::text::bigint;

-- Recreate indexes
CREATE INDEX idx_variants_product ON product_variants(product_id, position);
CREATE INDEX idx_variants_default ON product_variants(product_id) WHERE is_default = true;
CREATE INDEX idx_variants_stock_status ON product_variants(product_id, status, stock_quantity) WHERE status = 'active';
CREATE INDEX idx_variants_available ON product_variants(product_id, (stock_quantity - reserved_quantity)) 
  WHERE status = 'active' AND (stock_quantity - reserved_quantity) > 0;

-- Add foreign key constraint to listings
ALTER TABLE product_variants
ADD CONSTRAINT fk_product_variants_listing FOREIGN KEY (product_id) REFERENCES listings(id) ON DELETE CASCADE;

COMMENT ON COLUMN product_variants.product_id IS 'References listings.id (parent product)';
