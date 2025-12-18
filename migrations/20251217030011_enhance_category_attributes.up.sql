-- Migration: 20251217030011_enhance_category_attributes
-- Description: Add additional fields to category_attributes for Phase 2 (filters, customization)
-- Author: Elite Full-Stack Architect (Claude AI)
-- Date: 2025-12-17

-- ============================================================================
-- ADD COLUMNS TO category_attributes
-- ============================================================================

-- Add is_searchable (if not exists)
ALTER TABLE category_attributes
ADD COLUMN IF NOT EXISTS is_searchable BOOLEAN DEFAULT false;

-- Add is_filterable (if not exists)
ALTER TABLE category_attributes
ADD COLUMN IF NOT EXISTS is_filterable BOOLEAN DEFAULT true;

-- Add category-specific options override
ALTER TABLE category_attributes
ADD COLUMN IF NOT EXISTS category_options JSONB;

-- Add category-specific validation rules
ALTER TABLE category_attributes
ADD COLUMN IF NOT EXISTS custom_validation JSONB;

-- Add category-specific UI settings
ALTER TABLE category_attributes
ADD COLUMN IF NOT EXISTS custom_ui_settings JSONB;

-- ============================================================================
-- CREATE ADDITIONAL INDEXES
-- ============================================================================

CREATE INDEX IF NOT EXISTS idx_category_attributes_searchable
ON category_attributes(category_id, is_searchable)
WHERE is_searchable = true;

CREATE INDEX IF NOT EXISTS idx_category_attributes_filterable
ON category_attributes(category_id, is_filterable)
WHERE is_filterable = true;

-- ============================================================================
-- HELPER FUNCTION: Get attributes with inheritance
-- ============================================================================

CREATE OR REPLACE FUNCTION get_category_attributes_with_inheritance(
    p_category_id INTEGER,
    p_locale VARCHAR DEFAULT 'sr'
)
RETURNS TABLE (
    attribute_id INTEGER,
    attribute_code VARCHAR,
    attribute_name TEXT,
    attribute_type VARCHAR,
    purpose VARCHAR,
    is_required BOOLEAN,
    is_filterable BOOLEAN,
    is_searchable BOOLEAN,
    sort_order INT,
    options JSONB,
    validation_rules JSONB,
    ui_settings JSONB,
    source_category_id INTEGER,
    is_inherited BOOLEAN
) AS $$
BEGIN
    RETURN QUERY
    WITH RECURSIVE category_path AS (
        -- Start with the given category
        SELECT
            c.id,
            c.parent_id,
            c.level,
            1 as distance
        FROM categories c
        WHERE c.id = p_category_id

        UNION ALL

        -- Recursively get parent categories
        SELECT
            c.id,
            c.parent_id,
            c.level,
            cp.distance + 1
        FROM categories c
        INNER JOIN category_path cp ON c.id = cp.parent_id
    ),
    all_attributes AS (
        -- Get category-specific attributes (including inherited)
        SELECT DISTINCT ON (a.id)
            a.id as attribute_id,
            a.code as attribute_code,
            a.name->>p_locale as attribute_name,
            a.attribute_type,
            a.purpose,
            COALESCE(ca.is_required, a.is_required) as is_required,
            COALESCE(ca.is_filterable, a.is_filterable) as is_filterable,
            COALESCE(ca.is_searchable, a.is_searchable) as is_searchable,
            COALESCE(ca.sort_order, a.sort_order) as sort_order,
            COALESCE(ca.category_options, a.options) as options,
            COALESCE(ca.custom_validation, a.validation_rules) as validation_rules,
            COALESCE(ca.custom_ui_settings, a.ui_settings) as ui_settings,
            ca.category_id as source_category_id,
            (ca.category_id != p_category_id) as is_inherited,
            cp.distance as priority
        FROM category_path cp
        INNER JOIN category_attributes ca ON cp.id = ca.category_id
        INNER JOIN attributes a ON ca.attribute_id = a.id
        WHERE ca.is_enabled = true
          AND a.is_active = true
        ORDER BY a.id, priority ASC
    )
    SELECT
        aa.attribute_id,
        aa.attribute_code,
        aa.attribute_name,
        aa.attribute_type,
        aa.purpose,
        aa.is_required,
        aa.is_filterable,
        aa.is_searchable,
        aa.sort_order,
        aa.options,
        aa.validation_rules,
        aa.ui_settings,
        aa.source_category_id,
        aa.is_inherited
    FROM all_attributes aa
    ORDER BY aa.sort_order, aa.attribute_code;
END;
$$ LANGUAGE plpgsql STABLE;

-- ============================================================================
-- COMMENTS
-- ============================================================================

COMMENT ON COLUMN category_attributes.is_searchable IS 'Is this attribute searchable in this category?';
COMMENT ON COLUMN category_attributes.is_filterable IS 'Is this attribute filterable in this category?';
COMMENT ON COLUMN category_attributes.category_options IS 'Category-specific options override (NULL = use global)';
COMMENT ON COLUMN category_attributes.custom_validation IS 'Category-specific validation rules';
COMMENT ON COLUMN category_attributes.custom_ui_settings IS 'Category-specific UI settings';

COMMENT ON FUNCTION get_category_attributes_with_inheritance IS 'Get all attributes for a category including inherited from parents';

-- ============================================================================
-- END OF MIGRATION
-- ============================================================================
