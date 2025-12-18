-- Migration rollback: Drop product_variants table
-- Author: Claude (Elite Architect)
-- Date: 2025-12-17

BEGIN;

-- Drop triggers
DROP TRIGGER IF EXISTS trigger_variants_single_default ON product_variants;
DROP TRIGGER IF EXISTS trigger_variants_auto_status ON product_variants;
DROP TRIGGER IF EXISTS trigger_variants_updated_at ON product_variants;

-- Drop functions
DROP FUNCTION IF EXISTS enforce_single_default_variant();
DROP FUNCTION IF EXISTS auto_update_variant_status();

-- Drop table (CASCADE will drop dependent objects)
DROP TABLE IF EXISTS product_variants CASCADE;

COMMIT;
