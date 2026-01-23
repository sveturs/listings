-- Add missing columns to product_variants for compatibility with repository code

ALTER TABLE product_variants
ADD COLUMN IF NOT EXISTS cost_price NUMERIC(12,2),
ADD COLUMN IF NOT EXISTS stock_status VARCHAR(20) NOT NULL DEFAULT 'in_stock',
ADD COLUMN IF NOT EXISTS low_stock_threshold INTEGER,
ADD COLUMN IF NOT EXISTS variant_attributes JSONB DEFAULT '{}',
ADD COLUMN IF NOT EXISTS weight NUMERIC(8,3),
ADD COLUMN IF NOT EXISTS dimensions JSONB,
ADD COLUMN IF NOT EXISTS is_active BOOLEAN NOT NULL DEFAULT true,
ADD COLUMN IF NOT EXISTS view_count INTEGER NOT NULL DEFAULT 0,
ADD COLUMN IF NOT EXISTS sold_count INTEGER NOT NULL DEFAULT 0;

-- Add check constraints
ALTER TABLE product_variants
DROP CONSTRAINT IF EXISTS valid_cost_price,
ADD CONSTRAINT valid_cost_price CHECK (cost_price IS NULL OR cost_price >= 0);

ALTER TABLE product_variants
DROP CONSTRAINT IF EXISTS product_variants_stock_status_check,
ADD CONSTRAINT product_variants_stock_status_check 
CHECK (stock_status IN ('in_stock', 'out_of_stock', 'low_stock', 'discontinued'));

-- Add indexes for new columns
CREATE INDEX IF NOT EXISTS idx_product_variants_stock_status ON product_variants(stock_status) WHERE is_active = true;
CREATE INDEX IF NOT EXISTS idx_product_variants_is_active ON product_variants(is_active);

COMMENT ON COLUMN product_variants.cost_price IS 'Cost price for margin calculations';
COMMENT ON COLUMN product_variants.stock_status IS 'Stock availability status (auto-updated by trigger)';
COMMENT ON COLUMN product_variants.low_stock_threshold IS 'Threshold for low stock alerts';
COMMENT ON COLUMN product_variants.variant_attributes IS 'JSONB attributes specific to this variant (color, size, etc.)';
COMMENT ON COLUMN product_variants.weight IS 'Weight in grams';
COMMENT ON COLUMN product_variants.dimensions IS 'JSONB dimensions (length, width, height, unit)';
COMMENT ON COLUMN product_variants.is_active IS 'Is variant actively available for sale';
COMMENT ON COLUMN product_variants.view_count IS 'Number of times this variant was viewed';
COMMENT ON COLUMN product_variants.sold_count IS 'Number of times this variant was sold';
