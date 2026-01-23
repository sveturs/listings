-- Migration: Add tax_number and vat_number to storefronts
-- Created: 2026-01-18
-- Purpose: Add PIB (tax number) and PDV broj (VAT number) fields

ALTER TABLE storefronts
ADD COLUMN tax_number VARCHAR(9),   -- PIB (9 digits)
ADD COLUMN vat_number VARCHAR(11);  -- PDV broj (RS + 9 digits)

-- Create indexes
CREATE INDEX idx_storefronts_tax_number ON storefronts(tax_number) WHERE tax_number IS NOT NULL;
CREATE INDEX idx_storefronts_vat_number ON storefronts(vat_number) WHERE vat_number IS NOT NULL;

-- Add constraints
ALTER TABLE storefronts
ADD CONSTRAINT chk_tax_number_format CHECK (tax_number IS NULL OR tax_number ~ '^\d{9}$'),
ADD CONSTRAINT chk_vat_number_format CHECK (vat_number IS NULL OR vat_number ~ '^RS\d{9}$');

-- Comments
COMMENT ON COLUMN storefronts.tax_number IS 'PIB (Poreski identifikacioni broj) - 9 digits tax identification number';
COMMENT ON COLUMN storefronts.vat_number IS 'PDV broj (VAT number) - format: RS + 9 digits (optional)';
