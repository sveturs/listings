-- Migration Rollback: 20251217030010_attribute_values
-- Description: Drop attribute_values table
-- Author: Elite Full-Stack Architect (Claude AI)
-- Date: 2025-12-17

DROP TRIGGER IF EXISTS trigger_attribute_values_updated_at ON attribute_values;
DROP FUNCTION IF EXISTS update_attribute_values_updated_at();

DROP INDEX IF EXISTS idx_attribute_values_metadata_gin;
DROP INDEX IF EXISTS idx_attribute_values_label_gin;
DROP INDEX IF EXISTS idx_attribute_values_sort;
DROP INDEX IF EXISTS idx_attribute_values_active;
DROP INDEX IF EXISTS idx_attribute_values_attribute_id;

DROP TABLE IF EXISTS attribute_values CASCADE;
