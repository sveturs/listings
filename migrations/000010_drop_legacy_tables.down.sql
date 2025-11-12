-- Migration Rollback: 000010_drop_legacy_tables
-- Description: Cannot rollback - tables and data have been dropped
-- Created: 2025-11-06
-- Phase: 11.5

BEGIN;

-- =============================================================================
-- ROLLBACK IMPOSSIBLE
-- =============================================================================
--
-- This migration drops tables and their data permanently.
-- To rollback, you must restore from database backup:
--
--   /tmp/listings_dev_db_before_phase_11_5_*.sql
--
-- Restore procedure:
--   PGPASSWORD=listings_secret psql -h localhost -p 35434 -U listings_user \
--     -d listings_dev_db < /tmp/listings_dev_db_before_phase_11_5_YYYYMMDD_HHMMSS.sql
--
-- Dropped tables (cannot be recreated without data):
--   - c2c_listings (4 records lost)
--   - b2c_products (7 records lost)
--   - c2c_images (1 record lost)
--   - b2c_inventory_movements (3 records lost)
--   - c2c_categories (77 records lost)
--   - c2c_listing_variants (empty)
--   - c2c_orders (empty)
--   - b2c_product_variants (empty)
--   - c2c_favorites_backup_phase_11_4 (temporary backup)
--
-- NOTE: All this data was migrated to unified tables (listings, listing_images,
--       inventory_movements) before dropping, so the actual data is preserved
--       in the new schema.
--
-- =============================================================================

RAISE EXCEPTION 'Cannot rollback migration 000010: Tables and data have been dropped. Restore from backup: /tmp/listings_dev_db_before_phase_11_5_*.sql';

ROLLBACK;
