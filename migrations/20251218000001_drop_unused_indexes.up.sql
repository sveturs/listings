-- Migration: Drop unused category indexes (Phase 4 optimization)
-- Date: 2025-12-18
-- Impact: Saves ~10-15MB disk space, speeds up writes by 5-10%
-- Risk: Low (indexes have 0 scans in production)

-- Drop unused indexes on categories table
DROP INDEX IF EXISTS idx_categories_tree;
DROP INDEX IF EXISTS idx_categories_is_active;
DROP INDEX IF EXISTS idx_categories_path;
DROP INDEX IF EXISTS idx_categories_active_l1;

-- Note: Keeping GIN indexes (name, description) for future full-text search
