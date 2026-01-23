-- Rollback: Restore 'individual' legal entity type

-- Drop current constraint
ALTER TABLE storefronts
DROP CONSTRAINT chk_legal_entity_type;

-- Add old constraint with 'individual'
ALTER TABLE storefronts
ADD CONSTRAINT chk_legal_entity_type CHECK (
  legal_entity_type IN ('individual', 'preduzetnik', 'doo', 'ad')
);

-- Restore default value
ALTER TABLE storefronts
ALTER COLUMN legal_entity_type SET DEFAULT 'individual';

-- Restore comment
COMMENT ON COLUMN storefronts.legal_entity_type IS 'Legal entity type: individual (fizičko lice), preduzetnik, doo, ad';
