-- Rollback: Remove social_links column

DROP INDEX IF EXISTS idx_storefronts_social_links_gin;

ALTER TABLE storefronts
DROP COLUMN IF EXISTS social_links;
