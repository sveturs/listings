-- Migration: Drop duplicate and unused indexes
-- Created: 2025-10-31
-- Sprint: 2.1 - Database Cleanup
-- Purpose: Remove duplicate and unused indexes to improve write performance

-- 1. Drop duplicate UNIQUE constraint/index on same column
-- ai_category_decisions has two UNIQUE indexes on title_hash
-- First drop the CONSTRAINT (not index) since index is tied to constraint
ALTER TABLE ai_category_decisions DROP CONSTRAINT IF EXISTS ai_category_decisions_unique_hash;
-- Keep: idx_ai_decisions_unique_title_hash (newer, more explicit name)

-- 2. Drop duplicate index on title_hash (non-unique version conflicts with unique)
DROP INDEX IF EXISTS idx_ai_decisions_title_hash;
-- Keep: idx_ai_decisions_unique_title_hash (UNIQUE constraint is sufficient)

-- 3. Drop duplicate index on b2c_products.storefront_id
-- b2c_products_storefront_id_idx covers same column as idx_b2c_products_storefront
DROP INDEX IF EXISTS b2c_products_storefront_id_idx;
-- Keep: idx_b2c_products_storefront (composite: storefront_id, is_active) - more useful

-- 4. Drop duplicate index on b2c_product_variants.sku
-- Both b2c_product_variants_sku_idx and b2c_product_variants_sku_key index sku column
DROP INDEX IF EXISTS b2c_product_variants_sku_idx;
-- Keep: b2c_product_variants_sku_key (UNIQUE constraint)

-- 5. Drop unused indexes (idx_scan = 0 from analysis)
-- These indexes are never used in queries but slow down INSERT/UPDATE

-- Note: translations_pkey is PRIMARY KEY - cannot be dropped
-- Skipping to avoid errors

-- Car models unused indexes (no queries use these filters)
DROP INDEX IF EXISTS idx_car_models_drive_type;
DROP INDEX IF EXISTS idx_car_models_transmission_type;
DROP INDEX IF EXISTS idx_car_models_engine_type;
DROP INDEX IF EXISTS idx_car_models_serbia_popularity;

-- Category keywords unused indexes
DROP INDEX IF EXISTS idx_category_weight;
-- Keep: idx_unique_keyword_category and idx_keyword_lower for actual queries

-- Summary:
-- - Removed 1 duplicate UNIQUE constraint (ai_category_decisions)
-- - Removed 2 duplicate indexes
-- - Removed 5 unused indexes (idx_scan = 0)
-- - Total: 8 indexes/constraints dropped
-- - Expected improvement: ~5-10% faster INSERT/UPDATE on affected tables
