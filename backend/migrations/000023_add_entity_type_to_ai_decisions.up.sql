-- Add entity_type column to ai_category_decisions table to support both marketplace listings and storefront products

-- Add entity_type column with default 'listing'
ALTER TABLE ai_category_decisions
ADD COLUMN entity_type VARCHAR(20) DEFAULT 'listing' NOT NULL;

-- Drop existing unique constraint on title_hash
ALTER TABLE ai_category_decisions
DROP CONSTRAINT IF EXISTS idx_ai_decisions_unique_title_hash;

-- Create new unique constraint on (title_hash, entity_type)
ALTER TABLE ai_category_decisions
ADD CONSTRAINT ai_category_decisions_unique_hash
UNIQUE (title_hash, entity_type);

-- Create index for efficient querying by entity_type and created_at
CREATE INDEX idx_ai_decisions_entity_type
ON ai_category_decisions(entity_type, created_at DESC);

-- Add check constraint to ensure valid entity_type values
ALTER TABLE ai_category_decisions
ADD CONSTRAINT ai_category_decisions_entity_type_check
CHECK (entity_type IN ('listing', 'product'));