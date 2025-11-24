-- Migration: 000023_create_attributes_schema
-- Description: Create unified attributes system for listings microservice
-- Migrates from monolith: unified_attributes, unified_category_attributes, unified_attribute_values
-- Created: 2025-11-13

-- Enable required extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ============================================================================
-- TABLE 1: attributes - Attribute Metadata (Core table)
-- ============================================================================
CREATE TABLE IF NOT EXISTS attributes (
    id SERIAL PRIMARY KEY,
    code VARCHAR(100) NOT NULL UNIQUE,

    -- i18n fields (JSONB format: {"en": "...", "ru": "...", "sr": "..."})
    name JSONB NOT NULL DEFAULT '{}'::jsonb,
    display_name JSONB NOT NULL DEFAULT '{}'::jsonb,

    -- Attribute configuration
    attribute_type VARCHAR(50) NOT NULL CHECK (
        attribute_type IN ('text', 'textarea', 'number', 'boolean', 'select', 'multiselect', 'date', 'color', 'size')
    ),
    purpose VARCHAR(20) NOT NULL DEFAULT 'regular' CHECK (
        purpose IN ('regular', 'variant', 'both')
    ),

    -- JSONB configuration fields
    options JSONB DEFAULT '{}'::jsonb,           -- For select/multiselect: [{"value": "xl", "label": {"en": "XL", "ru": "XL", "sr": "XL"}}]
    validation_rules JSONB DEFAULT '{}'::jsonb,  -- {"min": 0, "max": 100, "pattern": "regex", "required_if": {...}}
    ui_settings JSONB DEFAULT '{}'::jsonb,       -- {"placeholder": {...}, "helpText": {...}, "icon": "...", "show_in_card": true}

    -- Behavior flags
    is_searchable BOOLEAN NOT NULL DEFAULT false,
    is_filterable BOOLEAN NOT NULL DEFAULT false,
    is_required BOOLEAN NOT NULL DEFAULT false,
    is_variant_compatible BOOLEAN NOT NULL DEFAULT false,
    affects_stock BOOLEAN NOT NULL DEFAULT false,
    affects_price BOOLEAN NOT NULL DEFAULT false,
    show_in_card BOOLEAN NOT NULL DEFAULT false,

    -- Metadata
    is_active BOOLEAN NOT NULL DEFAULT true,
    sort_order INTEGER NOT NULL DEFAULT 0,

    -- Legacy mapping (for migration tracking)
    legacy_category_attribute_id INTEGER,
    legacy_product_variant_attribute_id INTEGER,
    icon VARCHAR(255) DEFAULT '',

    -- Full-text search
    search_vector TSVECTOR,

    -- Timestamps
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for attributes
CREATE INDEX idx_attributes_code ON attributes(code);
CREATE INDEX idx_attributes_active ON attributes(is_active) WHERE is_active = true;
CREATE INDEX idx_attributes_purpose ON attributes(purpose);
CREATE INDEX idx_attributes_active_filterable ON attributes(is_active, is_filterable) WHERE is_active = true AND is_filterable = true;
CREATE INDEX idx_attributes_active_searchable ON attributes(is_active, is_searchable) WHERE is_active = true AND is_searchable = true;
CREATE INDEX idx_attributes_active_sort ON attributes(is_active, sort_order, (name->>'en')) WHERE is_active = true;
CREATE INDEX idx_attributes_show_in_card ON attributes(show_in_card) WHERE show_in_card = true;
CREATE INDEX idx_attributes_search_vector ON attributes USING GIN(search_vector);

-- Comment
COMMENT ON TABLE attributes IS 'Unified attributes metadata - core definitions for all attribute types';

-- ============================================================================
-- TABLE 2: category_attributes - Category-Attribute Relationships
-- ============================================================================
CREATE TABLE IF NOT EXISTS category_attributes (
    id SERIAL PRIMARY KEY,
    category_id INTEGER NOT NULL,
    attribute_id INTEGER NOT NULL REFERENCES attributes(id) ON DELETE CASCADE,

    -- Category-specific overrides (NULL = inherit from attributes table)
    is_enabled BOOLEAN DEFAULT true,
    is_required BOOLEAN,
    is_searchable BOOLEAN,
    is_filterable BOOLEAN,
    sort_order INTEGER NOT NULL DEFAULT 0,

    -- Category-specific configuration (overrides attributes.options and attributes.ui_settings)
    category_specific_options JSONB,
    custom_validation_rules JSONB,
    custom_ui_settings JSONB,

    -- Metadata
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    UNIQUE(category_id, attribute_id)
);

-- Indexes for category_attributes
CREATE INDEX idx_category_attributes_category ON category_attributes(category_id);
CREATE INDEX idx_category_attributes_attribute ON category_attributes(attribute_id);
CREATE INDEX idx_category_attributes_enabled ON category_attributes(is_enabled);
CREATE INDEX idx_category_attrs_composite ON category_attributes(category_id, attribute_id, is_enabled, sort_order) WHERE is_enabled = true;
CREATE INDEX idx_category_attrs_covering ON category_attributes(category_id, is_enabled, attribute_id, sort_order, is_required) WHERE is_enabled = true;

-- Comment
COMMENT ON TABLE category_attributes IS 'Category-specific attribute configurations and overrides';

-- ============================================================================
-- TABLE 3: listing_attribute_values - Listing Attribute Values
-- ============================================================================
CREATE TABLE IF NOT EXISTS listing_attribute_values (
    id SERIAL PRIMARY KEY,
    listing_id INTEGER NOT NULL REFERENCES listings(id) ON DELETE CASCADE,
    attribute_id INTEGER NOT NULL REFERENCES attributes(id) ON DELETE CASCADE,

    -- Polymorphic value storage (only one should be set based on attribute_type)
    value_text TEXT,
    value_number DECIMAL(20, 4),
    value_boolean BOOLEAN,
    value_date DATE,
    value_json JSONB,  -- For multiselect, complex objects, arrays

    -- Timestamps
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    UNIQUE(listing_id, attribute_id)
);

-- Indexes for listing_attribute_values
CREATE INDEX idx_listing_attr_values_listing ON listing_attribute_values(listing_id);
CREATE INDEX idx_listing_attr_values_attribute ON listing_attribute_values(attribute_id);
CREATE INDEX idx_listing_attr_values_text ON listing_attribute_values(value_text) WHERE value_text IS NOT NULL;
CREATE INDEX idx_listing_attr_values_number ON listing_attribute_values(value_number) WHERE value_number IS NOT NULL;
CREATE INDEX idx_listing_attr_values_numeric_ranges ON listing_attribute_values(attribute_id, value_number) WHERE value_number IS NOT NULL;
CREATE INDEX idx_listing_attr_values_json ON listing_attribute_values USING GIN(value_json);
CREATE INDEX idx_listing_attr_values_entity_with_type ON listing_attribute_values(listing_id, attribute_id)
    INCLUDE (value_text, value_number, value_boolean, value_date);

-- Comment
COMMENT ON TABLE listing_attribute_values IS 'Attribute values for listings - polymorphic storage';

-- ============================================================================
-- TABLE 4: category_variant_attributes - Variant Attribute Definitions
-- ============================================================================
CREATE TABLE IF NOT EXISTS category_variant_attributes (
    id SERIAL PRIMARY KEY,
    category_id INTEGER NOT NULL,
    attribute_id INTEGER NOT NULL REFERENCES attributes(id) ON DELETE CASCADE,

    -- Variant configuration
    is_required BOOLEAN NOT NULL DEFAULT false,
    affects_price BOOLEAN NOT NULL DEFAULT false,
    affects_stock BOOLEAN NOT NULL DEFAULT true,
    sort_order INTEGER NOT NULL DEFAULT 0,

    -- Display settings
    display_as VARCHAR(50) DEFAULT 'dropdown' CHECK (
        display_as IN ('dropdown', 'buttons', 'swatches', 'radio')
    ),

    -- Metadata
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    UNIQUE(category_id, attribute_id)
);

-- Indexes for category_variant_attributes
CREATE INDEX idx_category_variant_attrs_category ON category_variant_attributes(category_id);
CREATE INDEX idx_category_variant_attrs_attribute ON category_variant_attributes(attribute_id);
CREATE INDEX idx_category_variant_attrs_active ON category_variant_attributes(is_active, category_id) WHERE is_active = true;

-- Comment
COMMENT ON TABLE category_variant_attributes IS 'Variant-specific attribute configurations per category';

-- ============================================================================
-- TABLE 5: variant_attribute_values - Variant Attribute Values
-- ============================================================================
CREATE TABLE IF NOT EXISTS variant_attribute_values (
    id SERIAL PRIMARY KEY,
    variant_id INTEGER NOT NULL,  -- References product_variants (will be created later)
    attribute_id INTEGER NOT NULL REFERENCES attributes(id) ON DELETE CASCADE,

    -- Polymorphic value storage
    value_text TEXT,
    value_number DECIMAL(20, 4),
    value_boolean BOOLEAN,
    value_date DATE,
    value_json JSONB,

    -- Price/stock modifiers
    price_modifier DECIMAL(20, 4) DEFAULT 0.00,
    price_modifier_type VARCHAR(20) DEFAULT 'fixed' CHECK (
        price_modifier_type IN ('fixed', 'percent')
    ),

    -- Timestamps
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    UNIQUE(variant_id, attribute_id)
);

-- Indexes for variant_attribute_values
CREATE INDEX idx_variant_attr_values_variant ON variant_attribute_values(variant_id);
CREATE INDEX idx_variant_attr_values_attribute ON variant_attribute_values(attribute_id);
CREATE INDEX idx_variant_attr_values_text ON variant_attribute_values(value_text) WHERE value_text IS NOT NULL;
CREATE INDEX idx_variant_attr_values_json ON variant_attribute_values USING GIN(value_json);

-- Comment
COMMENT ON TABLE variant_attribute_values IS 'Attribute values for product variants with price modifiers';

-- ============================================================================
-- TABLE 6: attribute_options - Options for Select/Multiselect Attributes
-- ============================================================================
CREATE TABLE IF NOT EXISTS attribute_options (
    id SERIAL PRIMARY KEY,
    attribute_id INTEGER NOT NULL REFERENCES attributes(id) ON DELETE CASCADE,

    -- Option configuration
    option_value VARCHAR(255) NOT NULL,
    option_label JSONB NOT NULL DEFAULT '{}'::jsonb,  -- {"en": "...", "ru": "...", "sr": "..."}

    -- Display settings
    color_hex VARCHAR(7),  -- For color swatches: "#FF0000"
    image_url TEXT,        -- For image-based options
    icon VARCHAR(100),     -- For icon-based options

    -- Behavior
    is_default BOOLEAN NOT NULL DEFAULT false,
    is_active BOOLEAN NOT NULL DEFAULT true,
    sort_order INTEGER NOT NULL DEFAULT 0,

    -- Timestamps
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    UNIQUE(attribute_id, option_value)
);

-- Indexes for attribute_options
CREATE INDEX idx_attribute_options_attribute ON attribute_options(attribute_id);
CREATE INDEX idx_attribute_options_active ON attribute_options(is_active, attribute_id, sort_order) WHERE is_active = true;

-- Comment
COMMENT ON TABLE attribute_options IS 'Predefined options for select/multiselect attributes';

-- ============================================================================
-- TABLE 7: attribute_search_cache - OpenSearch Integration Cache
-- ============================================================================
CREATE TABLE IF NOT EXISTS attribute_search_cache (
    id SERIAL PRIMARY KEY,
    listing_id INTEGER NOT NULL REFERENCES listings(id) ON DELETE CASCADE,

    -- Denormalized attribute data for OpenSearch
    attributes_flat JSONB NOT NULL DEFAULT '{}'::jsonb,  -- {"brand": "Nike", "size": "XL", "color": "Red"}
    attributes_searchable TEXT,  -- Concatenated searchable values: "Nike XL Red"
    attributes_filterable JSONB NOT NULL DEFAULT '{}'::jsonb,  -- Structured filters for facets

    -- Cache metadata
    last_updated TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    cache_version INTEGER NOT NULL DEFAULT 1,

    UNIQUE(listing_id)
);

-- Indexes for attribute_search_cache
CREATE INDEX idx_attr_search_cache_listing ON attribute_search_cache(listing_id);
CREATE INDEX idx_attr_search_cache_updated ON attribute_search_cache(last_updated);
CREATE INDEX idx_attr_search_cache_flat ON attribute_search_cache USING GIN(attributes_flat);
CREATE INDEX idx_attr_search_cache_filterable ON attribute_search_cache USING GIN(attributes_filterable);

-- Comment
COMMENT ON TABLE attribute_search_cache IS 'Denormalized attribute cache for OpenSearch integration';

-- ============================================================================
-- TRIGGERS: Auto-update timestamps
-- ============================================================================

-- Function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_attributes_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Function to update search_vector for attributes
CREATE OR REPLACE FUNCTION update_attributes_search_vector()
RETURNS TRIGGER AS $$
BEGIN
    NEW.search_vector :=
        setweight(to_tsvector('english', COALESCE(NEW.name->>'en', '')), 'A') ||
        setweight(to_tsvector('russian', COALESCE(NEW.name->>'ru', '')), 'A') ||
        setweight(to_tsvector('english', COALESCE(NEW.code, '')), 'B');
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Apply triggers
CREATE TRIGGER trigger_attributes_updated_at
    BEFORE UPDATE ON attributes
    FOR EACH ROW
    EXECUTE FUNCTION update_attributes_timestamp();

CREATE TRIGGER trigger_attributes_search_vector
    BEFORE INSERT OR UPDATE OF name, code ON attributes
    FOR EACH ROW
    EXECUTE FUNCTION update_attributes_search_vector();

CREATE TRIGGER trigger_category_attributes_updated_at
    BEFORE UPDATE ON category_attributes
    FOR EACH ROW
    EXECUTE FUNCTION update_attributes_timestamp();

CREATE TRIGGER trigger_listing_attr_values_updated_at
    BEFORE UPDATE ON listing_attribute_values
    FOR EACH ROW
    EXECUTE FUNCTION update_attributes_timestamp();

CREATE TRIGGER trigger_category_variant_attrs_updated_at
    BEFORE UPDATE ON category_variant_attributes
    FOR EACH ROW
    EXECUTE FUNCTION update_attributes_timestamp();

CREATE TRIGGER trigger_variant_attr_values_updated_at
    BEFORE UPDATE ON variant_attribute_values
    FOR EACH ROW
    EXECUTE FUNCTION update_attributes_timestamp();

CREATE TRIGGER trigger_attribute_options_updated_at
    BEFORE UPDATE ON attribute_options
    FOR EACH ROW
    EXECUTE FUNCTION update_attributes_timestamp();

-- ============================================================================
-- GRANTS: Ensure proper permissions
-- ============================================================================

-- Grant permissions to application user (adjust username as needed)
-- GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO listings_user;
-- GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO listings_user;

-- ============================================================================
-- MIGRATION COMPLETE
-- ============================================================================
-- Tables created: 7
-- - attributes (203 records expected from monolith)
-- - category_attributes (relationships)
-- - listing_attribute_values (listing values)
-- - category_variant_attributes (variant definitions)
-- - variant_attribute_values (variant values)
-- - attribute_options (select/multiselect options)
-- - attribute_search_cache (OpenSearch cache)
--
-- Next steps:
-- 1. Run: migrations/scripts/001_export_monolith_attributes.sh
-- 2. Run: migrations/scripts/002_import_to_listings.sh
-- 3. Run: migrations/scripts/003_validate_migration.sql
-- ============================================================================
