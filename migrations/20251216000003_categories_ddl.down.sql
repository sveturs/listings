-- Rollback: Drop new categories table and related objects
-- Date: 2025-12-16

-- Drop triggers
DROP TRIGGER IF EXISTS trg_categories_updated_at ON categories;
DROP TRIGGER IF EXISTS trg_check_category_level ON categories;

-- Drop functions
DROP FUNCTION IF EXISTS check_category_level_constraint();
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop indexes (will be dropped automatically with table, but explicit for clarity)
DROP INDEX IF EXISTS idx_categories_slug;
DROP INDEX IF EXISTS idx_categories_parent_id;
DROP INDEX IF EXISTS idx_categories_level;
DROP INDEX IF EXISTS idx_categories_path;
DROP INDEX IF EXISTS idx_categories_is_active;
DROP INDEX IF EXISTS idx_categories_sort_order;
DROP INDEX IF EXISTS idx_categories_name_gin;
DROP INDEX IF EXISTS idx_categories_description_gin;
DROP INDEX IF EXISTS idx_categories_active_l1;
DROP INDEX IF EXISTS idx_categories_tree;

-- Drop table
DROP TABLE IF EXISTS categories CASCADE;
