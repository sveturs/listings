-- Migration rollback: 000019_add_price_constraint
-- Description: Remove CHECK constraints for price fields

-- Drop CHECK constraints
ALTER TABLE b2c_product_variants DROP CONSTRAINT IF EXISTS b2c_product_variants_price_check;
ALTER TABLE b2c_product_variants DROP CONSTRAINT IF EXISTS b2c_product_variants_compare_at_price_check;
ALTER TABLE b2c_product_variants DROP CONSTRAINT IF EXISTS b2c_product_variants_cost_price_check;
