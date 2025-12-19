-- Migration rollback: Drop variant_attribute_values updates
-- Author: Claude (Elite Architect)
-- Date: 2025-12-17

BEGIN;

-- Drop triggers
DROP TRIGGER IF EXISTS trigger_validate_variant_attr_value ON variant_attribute_values;
DROP TRIGGER IF EXISTS trigger_variant_attr_values_updated_at ON variant_attribute_values;

-- Drop function
DROP FUNCTION IF EXISTS validate_variant_attribute_value();

-- Drop table (will be recreated by original migration 000023 if needed)
DROP TABLE IF EXISTS variant_attribute_values CASCADE;

COMMIT;
