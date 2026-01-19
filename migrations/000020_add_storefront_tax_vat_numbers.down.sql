-- Rollback: Remove tax_number and vat_number

-- Drop constraints
ALTER TABLE storefronts
DROP CONSTRAINT IF EXISTS chk_tax_number_format,
DROP CONSTRAINT IF EXISTS chk_vat_number_format;

-- Drop indexes
DROP INDEX IF EXISTS idx_storefronts_tax_number;
DROP INDEX IF EXISTS idx_storefronts_vat_number;

-- Drop columns
ALTER TABLE storefronts
DROP COLUMN IF EXISTS tax_number,
DROP COLUMN IF EXISTS vat_number;
