-- Migration: Remove 'individual' legal entity type
-- Created: 2026-01-18
-- Purpose: Only business entities (preduzetnik, doo, ad) can have storefronts, not individuals

-- Migrate existing 'individual' records to 'preduzetnik'
UPDATE storefronts
SET legal_entity_type = 'preduzetnik'
WHERE legal_entity_type = 'individual';

-- Drop old constraint
ALTER TABLE storefronts
DROP CONSTRAINT chk_legal_entity_type;

-- Add new constraint without 'individual'
ALTER TABLE storefronts
ADD CONSTRAINT chk_legal_entity_type CHECK (
  legal_entity_type IN ('preduzetnik', 'doo', 'ad')
);

-- Update default value
ALTER TABLE storefronts
ALTER COLUMN legal_entity_type SET DEFAULT 'preduzetnik';

-- Update comment
COMMENT ON COLUMN storefronts.legal_entity_type IS 'Legal entity type: preduzetnik (sole proprietor), doo (LLC), ad (JSC)';
