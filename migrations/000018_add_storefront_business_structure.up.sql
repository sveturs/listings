-- Migration: Add business structure fields to storefronts
-- Created: 2026-01-18
-- Purpose: Separate legal entity type from business category

-- Add new columns
ALTER TABLE storefronts
ADD COLUMN legal_entity_type VARCHAR(50) DEFAULT 'individual' NOT NULL,
ADD COLUMN business_category VARCHAR(50) DEFAULT 'retail' NOT NULL,
ADD COLUMN full_legal_name VARCHAR(500),
ADD COLUMN registration_number VARCHAR(100),
ADD COLUMN legal_representative_name VARCHAR(255),
ADD COLUMN legal_representative_position VARCHAR(100);

-- Add constraints for valid values
ALTER TABLE storefronts
ADD CONSTRAINT chk_legal_entity_type CHECK (
  legal_entity_type IN ('individual', 'preduzetnik', 'doo', 'ad')
),
ADD CONSTRAINT chk_business_category CHECK (
  business_category IN ('retail', 'service', 'restaurant', 'grocery', 'other')
);

-- Create indexes for filtering
CREATE INDEX idx_storefronts_legal_entity_type ON storefronts(legal_entity_type);
CREATE INDEX idx_storefronts_business_category ON storefronts(business_category);

-- Comments for documentation
COMMENT ON COLUMN storefronts.legal_entity_type IS 'Legal entity type: individual (fizičko lice), preduzetnik, doo, ad';
COMMENT ON COLUMN storefronts.business_category IS 'Business activity category: retail, service, restaurant, grocery, other';
COMMENT ON COLUMN storefronts.full_legal_name IS 'Full legal name of the company (only for legal entities)';
COMMENT ON COLUMN storefronts.registration_number IS 'Matični broj (MB) - 8 digits registration number';
COMMENT ON COLUMN storefronts.legal_representative_name IS 'Name of legal representative (for legal entities)';
COMMENT ON COLUMN storefronts.legal_representative_position IS 'Position of legal representative (direktor, generalni direktor, etc.)';
