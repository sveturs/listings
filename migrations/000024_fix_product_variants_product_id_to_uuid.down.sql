-- Rollback product_variants.product_id from UUID to bigint

-- Drop foreign key
ALTER TABLE product_variants DROP CONSTRAINT IF EXISTS fk_product_variants_listing;

-- Drop indexes
DROP INDEX IF EXISTS idx_variants_product;
DROP INDEX IF EXISTS idx_variants_default;
DROP INDEX IF EXISTS idx_variants_stock_status;
DROP INDEX IF EXISTS idx_variants_available;

-- Add new column for bigint
ALTER TABLE product_variants ADD COLUMN product_id_bigint BIGINT;

-- Migrate data: product_id (UUID) -> product_id_bigint (from listings.id)
UPDATE product_variants pv
SET product_id_bigint = l.id
FROM listings l
WHERE l.uuid = pv.product_id;

-- Drop UUID column
ALTER TABLE product_variants DROP COLUMN product_id;

-- Rename bigint column
ALTER TABLE product_variants RENAME COLUMN product_id_bigint TO product_id;

-- Add NOT NULL constraint
ALTER TABLE product_variants ALTER COLUMN product_id SET NOT NULL;

-- Recreate indexes
CREATE INDEX idx_variants_product ON product_variants(product_id, position);
CREATE INDEX idx_variants_default ON product_variants(product_id) WHERE is_default = true;
CREATE INDEX idx_variants_stock_status ON product_variants(product_id, status, stock_quantity) WHERE status = 'active';
CREATE INDEX idx_variants_available ON product_variants(product_id, (stock_quantity - reserved_quantity)) 
  WHERE status = 'active' AND (stock_quantity - reserved_quantity) > 0;

-- Add foreign key to listings.id
ALTER TABLE product_variants
ADD CONSTRAINT fk_product_variants_listing FOREIGN KEY (product_id) REFERENCES listings(id) ON DELETE CASCADE;

COMMENT ON COLUMN product_variants.product_id IS 'References listings.id (parent product)';
