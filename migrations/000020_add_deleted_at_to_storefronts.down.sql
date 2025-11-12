-- Remove index first
DROP INDEX IF EXISTS idx_storefronts_deleted_at;

-- Remove column
ALTER TABLE storefronts
DROP COLUMN IF EXISTS deleted_at;
