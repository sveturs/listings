-- Migration: 20251222120000_category_ai_mapping (rollback)
-- Purpose: Drop AI category mapping table

DROP TRIGGER IF EXISTS trg_ai_mapping_updated_at ON category_ai_mapping;
DROP FUNCTION IF EXISTS update_ai_mapping_updated_at();
DROP TABLE IF EXISTS category_ai_mapping CASCADE;
