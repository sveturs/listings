-- Migration: Unify table names (remove c2c_/b2c_ prefixes)
-- Created: 2025-11-11 23:00:00
-- Description:
--   1. Drop empty/backup tables (b2c_product_variants, c2c_categories_backup_20251110)
--   2. Rename c2c_favorites -> listing_favorites
--   3. Rename c2c_categories -> categories (with self-referencing FK fix)
--   4. Add source_type CHECK constraint to listings table

-- =============================================================================
-- Step 1: Drop backup tables
-- =============================================================================

-- Drop backup table created on 2025-11-10
DROP TABLE IF EXISTS c2c_categories_backup_20251110 CASCADE;

-- Note: b2c_product_variants is actively used - DO NOT DROP

-- =============================================================================
-- Step 2: Rename c2c_favorites -> listing_favorites
-- =============================================================================
-- All FK constraints, indexes, and primary key will be automatically renamed
-- with table, preserving structure completely

ALTER TABLE IF EXISTS c2c_favorites RENAME TO listing_favorites;

-- Rename constraints to match new table name
-- PK: c2c_favorites_pkey -> listing_favorites_pkey
ALTER INDEX IF EXISTS c2c_favorites_pkey RENAME TO listing_favorites_pkey;

-- Indexes
ALTER INDEX IF EXISTS c2c_favorites_listing_id_idx RENAME TO listing_favorites_listing_id_idx;
ALTER INDEX IF EXISTS c2c_favorites_user_id_created_at_idx RENAME TO listing_favorites_user_id_created_at_idx;
ALTER INDEX IF EXISTS idx_c2c_favorites_listing_id RENAME TO idx_listing_favorites_listing_id;
ALTER INDEX IF EXISTS idx_c2c_favorites_unique RENAME TO idx_listing_favorites_unique;
ALTER INDEX IF EXISTS idx_c2c_favorites_user_id RENAME TO idx_listing_favorites_user_id;

-- Foreign key constraint
ALTER TABLE IF EXISTS listing_favorites
  RENAME CONSTRAINT fk_c2c_favorites_listing_id TO fk_listing_favorites_listing_id;

-- =============================================================================
-- Step 3: Rename c2c_categories -> categories
-- =============================================================================
-- This is more complex due to self-referencing FK and references from listings

-- First, rename the table
ALTER TABLE IF EXISTS c2c_categories RENAME TO categories;

-- Rename sequence
ALTER SEQUENCE IF EXISTS c2c_categories_id_seq RENAME TO categories_id_seq;

-- Rename primary key
ALTER INDEX IF EXISTS c2c_categories_pkey RENAME TO categories_pkey;

-- Rename unique constraint
ALTER INDEX IF EXISTS c2c_categories_slug_key RENAME TO categories_slug_key;

-- Rename regular indexes
ALTER INDEX IF EXISTS idx_c2c_categories_level RENAME TO idx_categories_level;
ALTER INDEX IF EXISTS idx_c2c_categories_parent_id RENAME TO idx_categories_parent_id;
ALTER INDEX IF EXISTS idx_c2c_categories_slug RENAME TO idx_categories_slug;
ALTER INDEX IF EXISTS idx_c2c_categories_title_en RENAME TO idx_categories_title_en;
ALTER INDEX IF EXISTS idx_c2c_categories_title_ru RENAME TO idx_categories_title_ru;
ALTER INDEX IF EXISTS idx_c2c_categories_title_sr RENAME TO idx_categories_title_sr;

-- Rename self-referencing FK constraint
ALTER TABLE IF EXISTS categories
  RENAME CONSTRAINT c2c_categories_parent_id_fkey TO categories_parent_id_fkey;

-- Rename check constraint
ALTER TABLE IF EXISTS categories
  RENAME CONSTRAINT check_root_categories_level TO check_categories_root_level;

-- Update FK from listings table to reference new constraint name
-- Note: The FK itself doesn't need to be recreated, just renamed
ALTER TABLE IF EXISTS listings
  RENAME CONSTRAINT fk_listings_category_id TO fk_listings_category_id_new;
ALTER TABLE IF EXISTS listings
  RENAME CONSTRAINT fk_listings_category_id_new TO fk_listings_category_id;

-- =============================================================================
-- Step 4: Add source_type CHECK constraint to listings
-- =============================================================================
-- Ensure source_type can only be 'c2c', 'b2c', or 'storefront'

DO $$
BEGIN
  -- Check if constraint doesn't exist before adding
  IF NOT EXISTS (
    SELECT 1 FROM pg_constraint
    WHERE conname = 'listings_source_type_check'
    AND conrelid = 'listings'::regclass
  ) THEN
    ALTER TABLE listings
      ADD CONSTRAINT listings_source_type_check
      CHECK (source_type IN ('c2c', 'b2c', 'storefront'));
  END IF;
END $$;

-- =============================================================================
-- Verification queries (commented out, uncomment to check)
-- =============================================================================
-- SELECT tablename FROM pg_tables WHERE schemaname = 'public' AND tablename LIKE '%favorites%';
-- SELECT tablename FROM pg_tables WHERE schemaname = 'public' AND tablename LIKE '%categories%';
-- SELECT conname FROM pg_constraint WHERE conrelid = 'listings'::regclass AND conname LIKE '%source_type%';
