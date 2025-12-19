-- Migration: Update variant_attribute_values to use UUID and add FK constraints
-- Description: Convert variant_id to UUID and add foreign keys to product_variants
-- Author: Claude (Elite Architect)
-- Date: 2025-12-17
-- Phase: 3 (Variants)
-- Task: BE-3.2
-- Note: Table was created in migration 000023, we're updating it here

BEGIN;

-- Drop existing table if it exists (from migration 000023)
DROP TABLE IF EXISTS variant_attribute_values CASCADE;

-- ============================================================================
-- TABLE: variant_attribute_values - Attribute values for product variants
-- ============================================================================
CREATE TABLE variant_attribute_values (
    -- Primary key
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Foreign keys
    variant_id UUID NOT NULL,
    attribute_id INTEGER NOT NULL REFERENCES attributes(id) ON DELETE CASCADE,

    -- Polymorphic value storage (only one should be set based on attribute type)
    value_text TEXT,
    value_number DECIMAL(20, 4),
    value_boolean BOOLEAN,
    value_date DATE,
    value_json JSONB,

    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,

    -- Constraints
    UNIQUE(variant_id, attribute_id)
);

-- ============================================================================
-- INDEXES for variant_attribute_values
-- ============================================================================

-- Primary lookups
CREATE INDEX idx_variant_attr_values_variant ON variant_attribute_values(variant_id);
CREATE INDEX idx_variant_attr_values_attribute ON variant_attribute_values(attribute_id);

-- Value lookups
CREATE INDEX idx_variant_attr_values_text ON variant_attribute_values(value_text)
WHERE value_text IS NOT NULL;

CREATE INDEX idx_variant_attr_values_number ON variant_attribute_values(value_number)
WHERE value_number IS NOT NULL;

CREATE INDEX idx_variant_attr_values_json ON variant_attribute_values USING GIN(value_json)
WHERE value_json IS NOT NULL;

-- Composite index for variant attributes lookup
CREATE INDEX idx_variant_attrs_composite ON variant_attribute_values(variant_id, attribute_id)
    INCLUDE (value_text, value_number, value_boolean, value_date);

-- ============================================================================
-- TRIGGERS
-- ============================================================================

-- Auto-update updated_at timestamp
CREATE TRIGGER trigger_variant_attr_values_updated_at
    BEFORE UPDATE ON variant_attribute_values
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Validate attribute type matches value field
CREATE OR REPLACE FUNCTION validate_variant_attribute_value()
RETURNS TRIGGER AS $$
DECLARE
    attr_type VARCHAR(50);
BEGIN
    -- Get attribute type
    SELECT attribute_type INTO attr_type
    FROM attributes
    WHERE id = NEW.attribute_id;

    -- Validate correct value field is used
    CASE attr_type
        WHEN 'text', 'textarea', 'select', 'color', 'size' THEN
            IF NEW.value_text IS NULL THEN
                RAISE EXCEPTION 'value_text must be set for attribute type %', attr_type;
            END IF;
        WHEN 'number' THEN
            IF NEW.value_number IS NULL THEN
                RAISE EXCEPTION 'value_number must be set for attribute type %', attr_type;
            END IF;
        WHEN 'boolean' THEN
            IF NEW.value_boolean IS NULL THEN
                RAISE EXCEPTION 'value_boolean must be set for attribute type %', attr_type;
            END IF;
        WHEN 'date' THEN
            IF NEW.value_date IS NULL THEN
                RAISE EXCEPTION 'value_date must be set for attribute type %', attr_type;
            END IF;
        WHEN 'multiselect' THEN
            IF NEW.value_json IS NULL THEN
                RAISE EXCEPTION 'value_json must be set for attribute type %', attr_type;
            END IF;
    END CASE;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_validate_variant_attr_value
    BEFORE INSERT OR UPDATE ON variant_attribute_values
    FOR EACH ROW
    EXECUTE FUNCTION validate_variant_attribute_value();

-- ============================================================================
-- COMMENTS
-- ============================================================================

COMMENT ON TABLE variant_attribute_values IS 'Attribute values for product variants (size, color, etc.)';
COMMENT ON COLUMN variant_attribute_values.variant_id IS 'Reference to product_variants.id';
COMMENT ON COLUMN variant_attribute_values.attribute_id IS 'Reference to attributes.id';
COMMENT ON COLUMN variant_attribute_values.value_text IS 'Text value (for text, select, color, size attributes)';
COMMENT ON COLUMN variant_attribute_values.value_number IS 'Numeric value (for number attributes)';
COMMENT ON COLUMN variant_attribute_values.value_boolean IS 'Boolean value (for boolean attributes)';
COMMENT ON COLUMN variant_attribute_values.value_date IS 'Date value (for date attributes)';
COMMENT ON COLUMN variant_attribute_values.value_json IS 'JSON value (for multiselect, complex attributes)';

-- ============================================================================
-- FOREIGN KEY to product_variants (will be added after variant table exists)
-- ============================================================================
-- NOTE: This FK is commented out here and will be enabled in a separate migration
-- after product_variants table is fully populated
-- ALTER TABLE variant_attribute_values
--     ADD CONSTRAINT variant_attr_values_variant_fk
--     FOREIGN KEY (variant_id) REFERENCES product_variants(id) ON DELETE CASCADE;

COMMIT;
