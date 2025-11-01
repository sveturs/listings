-- Migration: Drop legacy and unused tables
-- Created: 2025-10-31
-- Sprint: 2.1 - Database Cleanup
-- Purpose: Remove technical debt - unused/backup tables that are not referenced in code

-- 1. Drop backup tables (no FK dependencies, not used in code)
DROP TABLE IF EXISTS districts_leskovac_backup CASCADE;
DROP TABLE IF EXISTS districts_novi_sad_backup_20250715 CASCADE;

-- 2. Drop listing_attribute_values (empty table, replaced by unified_attribute_values)
-- Has FK to unified_attributes but no data (0 rows)
-- No other tables reference it
DROP TABLE IF EXISTS listing_attribute_values CASCADE;

-- Note: category_variant_attributes was considered but is still in use
-- by variant_attributes.go (has 4 rows), will be removed in future sprint
-- after code migration to unified_attributes
