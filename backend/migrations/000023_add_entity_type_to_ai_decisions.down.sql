-- Rollback: Remove entity_type support from ai_category_decisions table

-- Drop check constraint
ALTER TABLE ai_category_decisions
DROP CONSTRAINT IF EXISTS ai_category_decisions_entity_type_check;

-- Drop new index
DROP INDEX IF EXISTS idx_ai_decisions_entity_type;

-- Drop new unique constraint
ALTER TABLE ai_category_decisions
DROP CONSTRAINT IF EXISTS ai_category_decisions_unique_hash;

-- Recreate original unique constraint on title_hash only
CREATE UNIQUE INDEX idx_ai_decisions_unique_title_hash
ON ai_category_decisions(title_hash);

-- Drop entity_type column
ALTER TABLE ai_category_decisions
DROP COLUMN IF EXISTS entity_type;