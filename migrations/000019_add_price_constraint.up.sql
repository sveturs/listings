-- Migration: 000019_add_price_constraint
-- Description: Add CHECK constraint for price, compare_at_price, and cost_price to prevent negative values
-- Date: 2025-11-10
-- Author: Database schema improvement

-- Add CHECK constraints for price fields
ALTER TABLE b2c_product_variants
ADD CONSTRAINT b2c_product_variants_price_check CHECK (price IS NULL OR price >= 0);

ALTER TABLE b2c_product_variants
ADD CONSTRAINT b2c_product_variants_compare_at_price_check CHECK (compare_at_price IS NULL OR compare_at_price >= 0);

ALTER TABLE b2c_product_variants
ADD CONSTRAINT b2c_product_variants_cost_price_check CHECK (cost_price IS NULL OR cost_price >= 0);

-- Add comment
COMMENT ON CONSTRAINT b2c_product_variants_price_check ON b2c_product_variants IS 'Ensure price is non-negative';
COMMENT ON CONSTRAINT b2c_product_variants_compare_at_price_check ON b2c_product_variants IS 'Ensure compare_at_price is non-negative';
COMMENT ON CONSTRAINT b2c_product_variants_cost_price_check ON b2c_product_variants IS 'Ensure cost_price is non-negative';
