-- Migration Rollback: 20251217030015_fix_category_attributes_fk
-- Description: Revert category_id type back to INTEGER
-- Author: Elite Full-Stack Architect (Claude AI)
-- Date: 2025-12-17

-- Drop FK constraint
ALTER TABLE category_attributes DROP CONSTRAINT IF EXISTS category_attributes_category_id_fkey;

-- Drop indexes
DROP INDEX IF EXISTS idx_category_attrs_covering;
DROP INDEX IF EXISTS idx_category_attrs_composite;
DROP INDEX IF EXISTS idx_category_attributes_category;

-- Change back to INTEGER
ALTER TABLE category_attributes ALTER COLUMN category_id TYPE INTEGER USING NULL;

-- Recreate old indexes
CREATE INDEX idx_category_attributes_category ON category_attributes(category_id);
