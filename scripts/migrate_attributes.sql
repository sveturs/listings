-- ============================================================================
-- Migration Script: unified_attributes → attributes
-- ============================================================================
-- Source: svetubd.unified_attributes (localhost:5433)
-- Target: listings_dev_db.attributes (localhost:35434)
--
-- Converts VARCHAR i18n fields to JSONB format
-- Idempotent: skips existing records based on 'code' field
-- ============================================================================

BEGIN;

-- Create temporary table to hold source data
CREATE TEMP TABLE temp_attributes_migration (
    id INTEGER,
    code VARCHAR(100),
    name VARCHAR(100),
    display_name VARCHAR(200),
    attribute_type VARCHAR(50),
    purpose VARCHAR(20),
    options JSONB,
    validation_rules JSONB,
    ui_settings JSONB,
    is_searchable BOOLEAN,
    is_filterable BOOLEAN,
    is_required BOOLEAN,
    affects_stock BOOLEAN,
    affects_price BOOLEAN,
    sort_order INTEGER,
    is_active BOOLEAN,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    legacy_category_attribute_id INTEGER,
    legacy_product_variant_attribute_id INTEGER,
    is_variant_compatible BOOLEAN,
    icon VARCHAR(255),
    show_in_card BOOLEAN
);

-- Note: This script expects data to be imported into temp table first
-- See accompanying migration guide for data import steps

-- Insert into target table with JSONB conversion
INSERT INTO attributes (
    id,
    code,
    name,
    display_name,
    attribute_type,
    purpose,
    options,
    validation_rules,
    ui_settings,
    is_searchable,
    is_filterable,
    is_required,
    affects_stock,
    affects_price,
    sort_order,
    is_active,
    created_at,
    updated_at,
    legacy_category_attribute_id,
    legacy_product_variant_attribute_id,
    is_variant_compatible,
    icon,
    show_in_card
)
SELECT
    id,
    code,
    -- Convert VARCHAR to JSONB i18n format
    jsonb_build_object(
        'en', name,
        'ru', name,
        'sr', name
    ) AS name,
    jsonb_build_object(
        'en', display_name,
        'ru', display_name,
        'sr', display_name
    ) AS display_name,
    attribute_type,
    purpose,
    COALESCE(options, '{}'::jsonb),
    COALESCE(validation_rules, '{}'::jsonb),
    COALESCE(ui_settings, '{}'::jsonb),
    COALESCE(is_searchable, false),
    COALESCE(is_filterable, false),
    COALESCE(is_required, false),
    COALESCE(affects_stock, false),
    COALESCE(affects_price, false),
    COALESCE(sort_order, 0),
    COALESCE(is_active, true),
    COALESCE(created_at, CURRENT_TIMESTAMP),
    COALESCE(updated_at, CURRENT_TIMESTAMP),
    legacy_category_attribute_id,
    legacy_product_variant_attribute_id,
    COALESCE(is_variant_compatible, false),
    COALESCE(icon, ''),
    COALESCE(show_in_card, false)
FROM temp_attributes_migration
WHERE NOT EXISTS (
    -- Idempotent: skip if code already exists
    SELECT 1 FROM attributes WHERE attributes.code = temp_attributes_migration.code
);

-- Update sequence to continue from max ID
SELECT setval('attributes_id_seq', (SELECT COALESCE(MAX(id), 0) FROM attributes));

-- Clean up
DROP TABLE temp_attributes_migration;

COMMIT;

-- ============================================================================
-- Validation Queries (run after migration)
-- ============================================================================

-- Check record count
-- SELECT COUNT(*) as migrated_count FROM attributes;

-- Verify JSONB structure
-- SELECT id, code, name, display_name
-- FROM attributes
-- ORDER BY id
-- LIMIT 5;

-- Check for missing records (run on source DB and compare counts)
-- SELECT COUNT(*) FROM unified_attributes;
