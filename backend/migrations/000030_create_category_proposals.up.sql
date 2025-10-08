-- Migration: Create category_proposals table
-- Description: Stores proposals for new marketplace categories from storefront imports
-- Author: AI Category Analyzer
-- Date: 2025-10-06

-- Create category_proposals table
-- Note: User IDs are managed externally via Auth Service, no FK constraints
CREATE TABLE IF NOT EXISTS category_proposals (
    id SERIAL PRIMARY KEY,
    proposed_by_user_id INT NOT NULL, -- User ID from Auth Service
    storefront_id INT REFERENCES storefronts(id) ON DELETE SET NULL,
    name VARCHAR(255) NOT NULL,
    name_translations JSONB DEFAULT '{}'::jsonb, -- Translations: {"ru": "...", "en": "...", "sr": "..."}
    parent_category_id INT REFERENCES marketplace_categories(id) ON DELETE SET NULL,
    description TEXT,
    reasoning TEXT, -- AI reasoning for why this category should be created
    expected_products INT DEFAULT 0, -- Expected number of products in this category
    external_category_source VARCHAR(255), -- External category path that triggered this proposal
    similar_categories INT[] DEFAULT ARRAY[]::INT[], -- Related category IDs
    tags TEXT[] DEFAULT ARRAY[]::TEXT[], -- Tags for categorization
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'approved', 'rejected')),
    reviewed_by_user_id INT, -- User ID from Auth Service
    reviewed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for fast queries
CREATE INDEX idx_category_proposals_status ON category_proposals(status);
CREATE INDEX idx_category_proposals_proposed_by ON category_proposals(proposed_by_user_id);
CREATE INDEX idx_category_proposals_storefront ON category_proposals(storefront_id);
CREATE INDEX idx_category_proposals_created_at ON category_proposals(created_at DESC);

-- Create trigger to update updated_at
CREATE OR REPLACE FUNCTION update_category_proposals_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_category_proposals_updated_at
    BEFORE UPDATE ON category_proposals
    FOR EACH ROW
    EXECUTE FUNCTION update_category_proposals_updated_at();

-- Add comments for documentation
COMMENT ON TABLE category_proposals IS 'Proposals for new marketplace categories from AI analysis of storefront imports';
COMMENT ON COLUMN category_proposals.name IS 'Proposed category name (original language)';
COMMENT ON COLUMN category_proposals.name_translations IS 'Category name translations in JSON format: {"ru": "...", "en": "...", "sr": "..."}';
COMMENT ON COLUMN category_proposals.parent_category_id IS 'Suggested parent category for this new category';
COMMENT ON COLUMN category_proposals.reasoning IS 'AI-generated explanation for why this category should exist';
COMMENT ON COLUMN category_proposals.expected_products IS 'Number of products expected to use this category';
COMMENT ON COLUMN category_proposals.external_category_source IS 'Original external category path from import file';
COMMENT ON COLUMN category_proposals.similar_categories IS 'Array of related category IDs for reference';
COMMENT ON COLUMN category_proposals.status IS 'Proposal status: pending, approved, or rejected';
