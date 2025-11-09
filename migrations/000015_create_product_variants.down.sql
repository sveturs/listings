-- Migration: 000015_create_product_variants (ROLLBACK)
-- Description: Drop b2c_product_variants table
-- Date: 2025-11-09

-- Drop indexes
DROP INDEX IF EXISTS idx_b2c_product_variants_stock_status;
DROP INDEX IF EXISTS idx_b2c_product_variants_is_active;
DROP INDEX IF EXISTS idx_b2c_product_variants_product_id;

-- Drop table (will cascade drop foreign key constraint)
DROP TABLE IF EXISTS b2c_product_variants CASCADE;

-- Drop sequence
DROP SEQUENCE IF EXISTS b2c_product_variants_id_seq;
