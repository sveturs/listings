-- Migration: 20251217030010_attribute_values
-- Description: Create attribute_values table for predefined attribute values (Phase 2)
-- Author: Elite Full-Stack Architect (Claude AI)
-- Date: 2025-12-17
-- Note: Работает с существующей таблицей attributes (INTEGER PK)

-- ============================================================================
-- TABLE: attribute_values
-- ============================================================================

CREATE TABLE IF NOT EXISTS attribute_values (
    -- Primary Key
    id SERIAL PRIMARY KEY,

    -- Foreign key to existing attributes table (INTEGER)
    attribute_id INTEGER NOT NULL REFERENCES attributes(id) ON DELETE CASCADE,

    -- Value identifier (e.g., "black", "xl", "new")
    value VARCHAR(255) NOT NULL,

    -- Localized labels (JSONB)
    label JSONB NOT NULL DEFAULT '{}'::jsonb,

    -- Optional metadata (e.g., {"hex": "#000000"} for color)
    metadata JSONB DEFAULT '{}'::jsonb,

    -- Sort order for display
    sort_order INT DEFAULT 0,

    -- Is this value active?
    is_active BOOLEAN DEFAULT true,

    -- Timestamps
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),

    -- Unique constraint
    CONSTRAINT uq_attribute_value UNIQUE (attribute_id, value)
);

-- ============================================================================
-- INDEXES
-- ============================================================================

CREATE INDEX IF NOT EXISTS idx_attribute_values_attribute_id ON attribute_values(attribute_id);
CREATE INDEX IF NOT EXISTS idx_attribute_values_active ON attribute_values(is_active) WHERE is_active = true;
CREATE INDEX IF NOT EXISTS idx_attribute_values_sort ON attribute_values(attribute_id, sort_order, value);
CREATE INDEX IF NOT EXISTS idx_attribute_values_label_gin ON attribute_values USING GIN (label);
CREATE INDEX IF NOT EXISTS idx_attribute_values_metadata_gin ON attribute_values USING GIN (metadata);

-- ============================================================================
-- TRIGGER: updated_at
-- ============================================================================

CREATE OR REPLACE FUNCTION update_attribute_values_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_attribute_values_updated_at
    BEFORE UPDATE ON attribute_values
    FOR EACH ROW
    EXECUTE FUNCTION update_attribute_values_updated_at();

-- ============================================================================
-- COMMENTS
-- ============================================================================

COMMENT ON TABLE attribute_values IS 'Predefined values for select/multiselect attributes';
COMMENT ON COLUMN attribute_values.attribute_id IS 'Foreign key to attributes table';
COMMENT ON COLUMN attribute_values.value IS 'Value identifier (e.g., "black", "xl", "new")';
COMMENT ON COLUMN attribute_values.label IS 'Localized labels: {"sr": "Crna", "en": "Black", "ru": "Чёрный"}';
COMMENT ON COLUMN attribute_values.metadata IS 'Optional metadata (e.g., {"hex": "#000000"})';

-- ============================================================================
-- END OF MIGRATION
-- ============================================================================
