-- Fix product_variants.product_id to use UUID instead of bigint
-- This aligns with ProductVariantV2 domain model and listings.uuid

-- Drop foreign key constraint
ALTER TABLE product_variants DROP CONSTRAINT IF EXISTS fk_product_variants_listing;

-- Drop indexes referencing product_id
DROP INDEX IF EXISTS idx_variants_product;
DROP INDEX IF EXISTS idx_variants_default;
DROP INDEX IF EXISTS idx_variants_stock_status;
DROP INDEX IF EXISTS idx_variants_available;

-- Add new column for UUID
ALTER TABLE product_variants ADD COLUMN product_uuid UUID;

-- Migrate data: product_id (bigint) -> product_uuid (from listings.uuid)
UPDATE product_variants pv
SET product_uuid = l.uuid
FROM listings l
WHERE l.id = pv.product_id;

-- Drop old product_id column
ALTER TABLE product_variants DROP COLUMN product_id;

-- Rename product_uuid to product_id
ALTER TABLE product_variants RENAME COLUMN product_uuid TO product_id;

-- Add NOT NULL constraint
ALTER TABLE product_variants ALTER COLUMN product_id SET NOT NULL;

-- Recreate indexes
CREATE INDEX idx_variants_product ON product_variants(product_id, position);
CREATE INDEX idx_variants_default ON product_variants(product_id) WHERE is_default = true;
CREATE INDEX idx_variants_stock_status ON product_variants(product_id, status, stock_quantity) WHERE status = 'active';
CREATE INDEX idx_variants_available ON product_variants(product_id, (stock_quantity - reserved_quantity)) 
  WHERE status = 'active' AND (stock_quantity - reserved_quantity) > 0;

-- Add foreign key to listings.uuid
ALTER TABLE product_variants
ADD CONSTRAINT fk_product_variants_listing FOREIGN KEY (product_id) REFERENCES listings(uuid) ON DELETE CASCADE;

COMMENT ON COLUMN product_variants.product_id IS 'References listings.uuid (parent product UUID)';
