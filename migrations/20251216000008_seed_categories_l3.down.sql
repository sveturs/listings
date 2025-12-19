-- Rollback: Remove L3 (leaf-level) categories
-- Date: 2025-12-16
-- Purpose: Remove all L3 categories

DELETE FROM categories WHERE level = 3;
