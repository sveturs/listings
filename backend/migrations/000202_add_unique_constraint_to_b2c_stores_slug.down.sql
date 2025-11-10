-- Remove unique constraint and index from b2c_stores.slug
DROP INDEX IF EXISTS idx_b2c_stores_slug;

ALTER TABLE b2c_stores
    DROP CONSTRAINT IF EXISTS b2c_stores_slug_unique;
