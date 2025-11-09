-- Down migration for 000011_restore_categories_table
--
-- WARNING: DO NOT drop c2c_categories table!
-- It's required by:
-- 1. listings.category_id FK constraint
-- 2. Categories API
-- 3. Multiple unit tests
--
-- This down migration is intentionally kept minimal to prevent accidental drops.

BEGIN;

-- Only remove indexes (table and data remain)
DROP INDEX IF EXISTS idx_c2c_categories_parent_id;
DROP INDEX IF EXISTS idx_c2c_categories_slug;
DROP INDEX IF EXISTS idx_c2c_categories_level;

COMMIT;

-- If you absolutely must drop (NOT RECOMMENDED):
-- DROP TABLE IF EXISTS c2c_categories CASCADE;
