-- Migration: Create storefront_category_mappings table
-- Purpose: Map external category paths (from CSV imports) to internal marketplace categories
-- Date: 2025-10-06

CREATE TABLE IF NOT EXISTS storefront_category_mappings (
    id SERIAL PRIMARY KEY,

    -- Storefront reference
    storefront_id INTEGER NOT NULL,

    -- External category path from import source
    -- Example: "Electronics/Phones/Apple" or "Clothing > Men > Shirts"
    source_category_path TEXT NOT NULL,

    -- Target category in our marketplace
    target_category_id INTEGER NOT NULL,

    -- Mapping metadata
    is_manual BOOLEAN NOT NULL DEFAULT false,  -- true if created manually, false if via AI
    confidence_score FLOAT,                     -- AI confidence score (0.0-1.0), NULL for manual mappings

    -- Timestamps
    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,

    -- Foreign keys
    CONSTRAINT fk_storefront_category_mappings_storefront
        FOREIGN KEY (storefront_id) REFERENCES storefronts(id) ON DELETE CASCADE,
    CONSTRAINT fk_storefront_category_mappings_category
        FOREIGN KEY (target_category_id) REFERENCES marketplace_categories(id) ON DELETE CASCADE,

    -- Unique constraint: one mapping per storefront + category path
    CONSTRAINT uq_storefront_category_mapping
        UNIQUE (storefront_id, source_category_path)
);

-- Indexes for performance
CREATE INDEX idx_storefront_category_mappings_storefront
    ON storefront_category_mappings(storefront_id);

CREATE INDEX idx_storefront_category_mappings_target_category
    ON storefront_category_mappings(target_category_id);

-- Index for quick lookup by source category path
CREATE INDEX idx_storefront_category_mappings_source_path
    ON storefront_category_mappings(storefront_id, source_category_path);

-- Trigger for updated_at
CREATE OR REPLACE FUNCTION update_storefront_category_mappings_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_storefront_category_mappings_updated_at
    BEFORE UPDATE ON storefront_category_mappings
    FOR EACH ROW
    EXECUTE FUNCTION update_storefront_category_mappings_updated_at();

-- Comments for documentation
COMMENT ON TABLE storefront_category_mappings IS 'Maps external category paths from import sources to internal marketplace categories';
COMMENT ON COLUMN storefront_category_mappings.source_category_path IS 'Category path from external source (e.g., CSV file). Format varies by source.';
COMMENT ON COLUMN storefront_category_mappings.target_category_id IS 'Target category ID in marketplace_categories';
COMMENT ON COLUMN storefront_category_mappings.is_manual IS 'True if mapping was created manually, false if auto-detected via AI';
COMMENT ON COLUMN storefront_category_mappings.confidence_score IS 'AI confidence score (0.0-1.0), NULL for manual mappings';
