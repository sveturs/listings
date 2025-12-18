-- Migration: Create new categories table with JSONB multilingual support
-- Date: 2025-12-16
-- Purpose: Phase 1 - Foundation for Vondi Marketplace category system

-- Drop existing categories table if exists (cleaned in previous migration)
DROP TABLE IF EXISTS categories CASCADE;

-- Create categories table with JSONB for multilingual support
CREATE TABLE categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Slug (unique identifier for URLs)
    slug VARCHAR(255) NOT NULL UNIQUE,

    -- Hierarchy fields
    parent_id UUID REFERENCES categories(id) ON DELETE CASCADE,
    level INTEGER NOT NULL CHECK (level >= 1 AND level <= 3),
    path VARCHAR(1000) NOT NULL, -- Full path like 'clothing/mens/tshirts'
    sort_order INTEGER NOT NULL DEFAULT 0,

    -- Multilingual JSONB fields
    name JSONB NOT NULL DEFAULT '{}'::jsonb, -- {"sr": "Odeća", "en": "Clothing", "ru": "Одежда"}
    description JSONB DEFAULT '{}'::jsonb, -- {"sr": "...", "en": "...", "ru": "..."}

    -- SEO JSONB fields (per locale)
    meta_title JSONB DEFAULT '{}'::jsonb,
    meta_description JSONB DEFAULT '{}'::jsonb,
    meta_keywords JSONB DEFAULT '{}'::jsonb,

    -- Display settings
    icon VARCHAR(50), -- Emoji or icon name
    image_url VARCHAR(500), -- Category image

    -- Status
    is_active BOOLEAN NOT NULL DEFAULT true,

    -- External mappings
    google_category_id INTEGER, -- Google Product Taxonomy
    facebook_category_id VARCHAR(100), -- Facebook Product Category

    -- Timestamps
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Add comment to table
COMMENT ON TABLE categories IS 'Marketplace categories with multilingual JSONB support (sr/en/ru)';

-- Add comments to JSONB columns
COMMENT ON COLUMN categories.name IS 'Multilingual category name: {"sr": "...", "en": "...", "ru": "..."}';
COMMENT ON COLUMN categories.description IS 'Multilingual category description: {"sr": "...", "en": "...", "ru": "..."}';
COMMENT ON COLUMN categories.meta_title IS 'SEO meta title per locale: {"sr": "...", "en": "...", "ru": "..."}';
COMMENT ON COLUMN categories.meta_description IS 'SEO meta description per locale: {"sr": "...", "en": "...", "ru": "..."}';
COMMENT ON COLUMN categories.path IS 'Full hierarchical path (slug-based): clothing/mens/tshirts';
COMMENT ON COLUMN categories.level IS 'Hierarchy level: 1 (L1 top), 2 (L2 subcategory), 3 (L3 leaf)';

-- Add constraint: L1 categories must have parent_id = NULL
CREATE OR REPLACE FUNCTION check_category_level_constraint()
RETURNS TRIGGER AS $$
BEGIN
    -- L1 categories (level=1) must have parent_id = NULL
    IF NEW.level = 1 AND NEW.parent_id IS NOT NULL THEN
        RAISE EXCEPTION 'Level 1 categories must have parent_id = NULL';
    END IF;

    -- L2/L3 categories must have parent_id
    IF NEW.level > 1 AND NEW.parent_id IS NULL THEN
        RAISE EXCEPTION 'Level % categories must have a parent_id', NEW.level;
    END IF;

    -- Verify parent's level is exactly (current level - 1)
    IF NEW.parent_id IS NOT NULL THEN
        DECLARE
            parent_level INTEGER;
        BEGIN
            SELECT level INTO parent_level FROM categories WHERE id = NEW.parent_id;

            IF parent_level IS NULL THEN
                RAISE EXCEPTION 'Parent category not found';
            END IF;

            IF parent_level != NEW.level - 1 THEN
                RAISE EXCEPTION 'Parent level must be % for level % category', NEW.level - 1, NEW.level;
            END IF;
        END;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger for level constraint
CREATE TRIGGER trg_check_category_level
BEFORE INSERT OR UPDATE ON categories
FOR EACH ROW EXECUTE FUNCTION check_category_level_constraint();

-- Create trigger for updated_at
CREATE TRIGGER trg_categories_updated_at
BEFORE UPDATE ON categories
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Create function if doesn't exist
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create indexes for performance
CREATE INDEX idx_categories_slug ON categories(slug);
CREATE INDEX idx_categories_parent_id ON categories(parent_id) WHERE parent_id IS NOT NULL;
CREATE INDEX idx_categories_level ON categories(level);
CREATE INDEX idx_categories_path ON categories USING btree(path varchar_pattern_ops);
CREATE INDEX idx_categories_is_active ON categories(is_active) WHERE is_active = true;
CREATE INDEX idx_categories_sort_order ON categories(parent_id, sort_order) WHERE parent_id IS NOT NULL;

-- GIN indexes for JSONB full-text search
CREATE INDEX idx_categories_name_gin ON categories USING gin(name jsonb_path_ops);
CREATE INDEX idx_categories_description_gin ON categories USING gin(description jsonb_path_ops);

-- Partial indexes for active L1 categories (main menu)
CREATE INDEX idx_categories_active_l1 ON categories(sort_order)
WHERE level = 1 AND is_active = true;

-- Composite index for tree queries
CREATE INDEX idx_categories_tree ON categories(parent_id, level, sort_order)
WHERE is_active = true;
