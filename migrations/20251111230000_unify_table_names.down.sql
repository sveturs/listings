-- Migration rollback: Unify table names
-- Created: 2025-11-11 23:00:00
-- Description: Reverts table renames back to c2c_ prefixed names
-- Note: Does NOT restore deleted empty tables (b2c_product_variants, c2c_categories_backup_20251110)

-- =============================================================================
-- Step 1: Remove source_type CHECK constraint from listings
-- =============================================================================

ALTER TABLE IF EXISTS listings
  DROP CONSTRAINT IF EXISTS listings_source_type_check;

-- =============================================================================
-- Step 2: Revert categories -> c2c_categories
-- =============================================================================

-- Revert FK from listings table
ALTER TABLE IF EXISTS listings
  RENAME CONSTRAINT fk_listings_category_id TO fk_listings_category_id_old;
ALTER TABLE IF EXISTS listings
  RENAME CONSTRAINT fk_listings_category_id_old TO fk_listings_category_id;

-- Revert check constraint name
ALTER TABLE IF EXISTS categories
  RENAME CONSTRAINT check_categories_root_level TO check_root_categories_level;

-- Revert self-referencing FK
ALTER TABLE IF EXISTS categories
  RENAME CONSTRAINT categories_parent_id_fkey TO c2c_categories_parent_id_fkey;

-- Revert indexes
ALTER INDEX IF EXISTS idx_categories_title_sr RENAME TO idx_c2c_categories_title_sr;
ALTER INDEX IF EXISTS idx_categories_title_ru RENAME TO idx_c2c_categories_title_ru;
ALTER INDEX IF EXISTS idx_categories_title_en RENAME TO idx_c2c_categories_title_en;
ALTER INDEX IF EXISTS idx_categories_slug RENAME TO idx_c2c_categories_slug;
ALTER INDEX IF EXISTS idx_categories_parent_id RENAME TO idx_c2c_categories_parent_id;
ALTER INDEX IF EXISTS idx_categories_level RENAME TO idx_c2c_categories_level;

-- Revert unique constraint
ALTER INDEX IF EXISTS categories_slug_key RENAME TO c2c_categories_slug_key;

-- Revert primary key
ALTER INDEX IF EXISTS categories_pkey RENAME TO c2c_categories_pkey;

-- Revert sequence
ALTER SEQUENCE IF EXISTS categories_id_seq RENAME TO c2c_categories_id_seq;

-- Revert table name
ALTER TABLE IF EXISTS categories RENAME TO c2c_categories;

-- =============================================================================
-- Step 3: Revert listing_favorites -> c2c_favorites
-- =============================================================================

-- Revert FK constraint
ALTER TABLE IF EXISTS listing_favorites
  RENAME CONSTRAINT fk_listing_favorites_listing_id TO fk_c2c_favorites_listing_id;

-- Revert indexes
ALTER INDEX IF EXISTS idx_listing_favorites_user_id RENAME TO idx_c2c_favorites_user_id;
ALTER INDEX IF EXISTS idx_listing_favorites_unique RENAME TO idx_c2c_favorites_unique;
ALTER INDEX IF EXISTS idx_listing_favorites_listing_id RENAME TO idx_c2c_favorites_listing_id;
ALTER INDEX IF EXISTS listing_favorites_user_id_created_at_idx RENAME TO c2c_favorites_user_id_created_at_idx;
ALTER INDEX IF EXISTS listing_favorites_listing_id_idx RENAME TO c2c_favorites_listing_id_idx;

-- Revert primary key
ALTER INDEX IF EXISTS listing_favorites_pkey RENAME TO c2c_favorites_pkey;

-- Revert table name
ALTER TABLE IF EXISTS listing_favorites RENAME TO c2c_favorites;

-- =============================================================================
-- Step 4: Note about dropped tables
-- =============================================================================
-- The following tables were dropped in UP migration and are NOT restored:
--   - b2c_product_variants (was empty)
--   - c2c_categories_backup_20251110 (was a backup)
-- If restoration is needed, recover from database backups.

-- =============================================================================
-- Verification queries (commented out, uncomment to check)
-- =============================================================================
-- SELECT tablename FROM pg_tables WHERE schemaname = 'public' AND tablename LIKE '%favorites%';
-- SELECT tablename FROM pg_tables WHERE schemaname = 'public' AND tablename LIKE '%categories%';
-- SELECT conname FROM pg_constraint WHERE conrelid = 'listings'::regclass AND conname LIKE '%source_type%';
