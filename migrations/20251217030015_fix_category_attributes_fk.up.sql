-- Migration: 20251217030015_fix_category_attributes_fk
-- Description: Fix category_attributes.category_id type mismatch (INTEGER → UUID) and add FK constraint
-- Author: Elite Full-Stack Architect (Claude AI)
-- Date: 2025-12-17
-- Note: This is a prerequisite for BE-2.7 (seed category_attributes)

-- ============================================================================
-- FIX TYPE MISMATCH
-- ============================================================================

-- Drop existing indexes first
DROP INDEX IF EXISTS idx_category_attributes_category;
DROP INDEX IF EXISTS idx_category_attrs_composite;
DROP INDEX IF EXISTS idx_category_attrs_covering;

-- Change category_id from INTEGER to UUID
ALTER TABLE category_attributes
ALTER COLUMN category_id TYPE UUID USING NULL;  -- Clear data first since table is empty

-- Recreate indexes
CREATE INDEX IF NOT EXISTS idx_category_attributes_category ON category_attributes(category_id);
CREATE INDEX IF NOT EXISTS idx_category_attrs_composite ON category_attributes(category_id, attribute_id, is_enabled, sort_order) WHERE is_enabled = true;
CREATE INDEX IF NOT EXISTS idx_category_attrs_covering ON category_attributes(category_id, is_enabled, attribute_id, sort_order, is_required) WHERE is_enabled = true;

-- ============================================================================
-- ADD FOREIGN KEY CONSTRAINT
-- ============================================================================

ALTER TABLE category_attributes
ADD CONSTRAINT category_attributes_category_id_fkey
FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE;

-- ============================================================================
-- COMMENTS
-- ============================================================================

COMMENT ON COLUMN category_attributes.category_id IS 'Foreign key to categories.id (UUID)';

-- ============================================================================
-- END OF MIGRATION
-- ============================================================================
