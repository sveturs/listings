-- Migration Rollback: 20251217030011_enhance_category_attributes
-- Description: Remove Phase 2 enhancements from category_attributes
-- Author: Elite Full-Stack Architect (Claude AI)
-- Date: 2025-12-17

-- Drop helper function
DROP FUNCTION IF EXISTS get_category_attributes_with_inheritance(INTEGER, VARCHAR);

-- Drop indexes
DROP INDEX IF EXISTS idx_category_attributes_filterable;
DROP INDEX IF EXISTS idx_category_attributes_searchable;

-- Drop columns
ALTER TABLE category_attributes DROP COLUMN IF EXISTS custom_ui_settings;
ALTER TABLE category_attributes DROP COLUMN IF EXISTS custom_validation;
ALTER TABLE category_attributes DROP COLUMN IF EXISTS category_options;
ALTER TABLE category_attributes DROP COLUMN IF EXISTS is_filterable;
ALTER TABLE category_attributes DROP COLUMN IF EXISTS is_searchable;
