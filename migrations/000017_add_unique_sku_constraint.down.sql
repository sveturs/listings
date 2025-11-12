-- Migration: 000017_add_unique_sku_constraint (ROLLBACK)
-- Description: Remove UNIQUE constraint on SKU
-- Date: 2025-11-09

-- Drop UNIQUE constraint on SKU
ALTER TABLE b2c_product_variants
DROP CONSTRAINT IF EXISTS b2c_product_variants_sku_unique;
