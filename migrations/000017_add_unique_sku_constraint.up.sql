-- Migration: 000017_add_unique_sku_constraint
-- Description: Add UNIQUE constraint on SKU in b2c_product_variants table
-- Date: 2025-11-09
-- Author: Test fix
--
-- This migration adds a unique constraint on the SKU field to prevent duplicate SKUs.
-- This is required for the duplicate SKU test and business logic.

-- Add UNIQUE constraint on SKU (allowing NULL since SKU is optional)
-- This constraint only applies to non-NULL SKUs
ALTER TABLE b2c_product_variants
ADD CONSTRAINT b2c_product_variants_sku_unique UNIQUE (sku);

-- Add comment
COMMENT ON CONSTRAINT b2c_product_variants_sku_unique ON b2c_product_variants IS 'Ensures SKU uniqueness when provided (allows NULL)';
