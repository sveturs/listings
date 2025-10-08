-- Migration: Drop storefront_category_mappings table
-- Purpose: Rollback migration 000029
-- Date: 2025-10-06

DROP TRIGGER IF EXISTS trigger_update_storefront_category_mappings_updated_at ON storefront_category_mappings;
DROP FUNCTION IF EXISTS update_storefront_category_mappings_updated_at() CASCADE;
DROP INDEX IF EXISTS idx_storefront_category_mappings_source_path;
DROP INDEX IF EXISTS idx_storefront_category_mappings_target_category;
DROP INDEX IF EXISTS idx_storefront_category_mappings_storefront;
DROP TABLE IF EXISTS storefront_category_mappings CASCADE;
