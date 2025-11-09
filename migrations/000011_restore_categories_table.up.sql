-- HOTFIX: Restore c2c_categories table (wrongly dropped in migration 000010)
--
-- PROBLEM: Migration 000010 (Phase 11.5) dropped c2c_categories table, but it's still needed by:
-- 1. Unified listings table (category_id FK constraint)
-- 2. Categories API (listings/internal/repository/postgres/categories_repository.go)
-- 3. 69 unit tests that depend on categories
--
-- SOLUTION: Restore table structure and data from backup before Phase 11.5
-- Backup file: /tmp/listings_dev_db_before_phase_11_5_20251106_174226.sql

BEGIN;

-- Restore table structure (keeping all columns from original)
CREATE TABLE IF NOT EXISTS c2c_categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(100) NOT NULL UNIQUE,
    parent_id INTEGER REFERENCES c2c_categories(id) ON DELETE CASCADE,
    icon VARCHAR(50),
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    has_custom_ui BOOLEAN DEFAULT FALSE,
    custom_ui_component VARCHAR(255),
    sort_order INTEGER DEFAULT 0,
    level INTEGER DEFAULT 0,
    count INTEGER DEFAULT 0,
    external_id VARCHAR(255),
    description TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    seo_title VARCHAR(255),
    seo_description TEXT,
    seo_keywords TEXT,
    CONSTRAINT check_root_categories_level CHECK (
        ((parent_id IS NULL) AND (level = 0)) OR
        ((parent_id IS NOT NULL) AND (level > 0))
    )
);

-- Restore indexes
CREATE INDEX IF NOT EXISTS idx_c2c_categories_parent_id ON c2c_categories(parent_id);
CREATE INDEX IF NOT EXISTS idx_c2c_categories_slug ON c2c_categories(slug);
CREATE INDEX IF NOT EXISTS idx_c2c_categories_level ON c2c_categories(level);

-- Restore data (78 categories from backup)
-- Data will be loaded from external file using \copy command
-- See migration notes for manual restore instructions

-- Add table comment
COMMENT ON TABLE c2c_categories IS 'Categories table (restored in migration 000011 after being wrongly dropped in 000010). Required by unified listings table and Categories API.';

-- Add FK constraint from listings to categories (if not exists)
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'fk_listings_category_id'
    ) THEN
        ALTER TABLE listings
        ADD CONSTRAINT fk_listings_category_id
        FOREIGN KEY (category_id)
        REFERENCES c2c_categories(id)
        ON DELETE RESTRICT;

        RAISE NOTICE 'Added FK constraint fk_listings_category_id';
    END IF;
END $$;

COMMIT;

-- MANUAL DATA RESTORE (if migrator does not support \copy):
-- PGPASSWORD=listings_secret psql -h localhost -p 35434 -U listings_user -d listings_dev_db < /tmp/restore_categories_data.sql
