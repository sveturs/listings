-- Migration: Drop category_proposals table
-- Description: Rollback for category proposals table creation
-- Date: 2025-10-06

-- Drop trigger and function
DROP TRIGGER IF EXISTS trigger_update_category_proposals_updated_at ON category_proposals;
DROP FUNCTION IF EXISTS update_category_proposals_updated_at();

-- Drop indexes
DROP INDEX IF EXISTS idx_category_proposals_created_at;
DROP INDEX IF EXISTS idx_category_proposals_storefront;
DROP INDEX IF EXISTS idx_category_proposals_proposed_by;
DROP INDEX IF EXISTS idx_category_proposals_status;

-- Drop table
DROP TABLE IF EXISTS category_proposals;
