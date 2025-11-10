-- Add unique constraint to b2c_stores.slug to prevent duplicates
-- First ensure all slugs are unique and non-empty
UPDATE b2c_stores
SET slug = 'store-' || id
WHERE slug = '' OR slug IS NULL;

-- Add unique constraint
ALTER TABLE b2c_stores
    ADD CONSTRAINT b2c_stores_slug_unique UNIQUE (slug);

-- Create index for performance
CREATE INDEX IF NOT EXISTS idx_b2c_stores_slug ON b2c_stores(slug);
