-- Rollback: Remove business structure fields from storefronts

-- Drop indexes
DROP INDEX IF EXISTS idx_storefronts_business_category;
DROP INDEX IF EXISTS idx_storefronts_legal_entity_type;

-- Drop columns
ALTER TABLE storefronts
DROP COLUMN IF EXISTS legal_representative_position,
DROP COLUMN IF EXISTS legal_representative_name,
DROP COLUMN IF EXISTS registration_number,
DROP COLUMN IF EXISTS full_legal_name,
DROP COLUMN IF EXISTS business_category,
DROP COLUMN IF EXISTS legal_entity_type;
