-- Migration: Add slug and expires_at fields to listings table
-- Phase: 13.1.15.4 - Repository Layer Merge
-- Purpose: Support SEO-friendly URLs and listing expiration

-- Add slug column (nullable initially for data migration)
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name = 'listings' AND column_name = 'slug'
    ) THEN
        ALTER TABLE listings ADD COLUMN slug VARCHAR(255);
        RAISE NOTICE 'Added slug column to listings table';
    ELSE
        RAISE NOTICE 'Column slug already exists, skipping';
    END IF;
END $$;

-- Add expires_at column (nullable, optional field for C2C listings)
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name = 'listings' AND column_name = 'expires_at'
    ) THEN
        ALTER TABLE listings ADD COLUMN expires_at TIMESTAMP WITH TIME ZONE;
        RAISE NOTICE 'Added expires_at column to listings table';
    ELSE
        RAISE NOTICE 'Column expires_at already exists, skipping';
    END IF;
END $$;

-- Generate slugs for existing listings without one
-- Slug format: lowercase, alphanumeric + hyphens, no special chars
UPDATE listings
SET slug = LOWER(
    TRIM(
        BOTH '-' FROM
        REGEXP_REPLACE(
            REGEXP_REPLACE(
                REGEXP_REPLACE(title, '[^a-zA-Z0-9\s-]', '', 'g'),  -- Remove special chars
                '\s+', '-', 'g'                                      -- Replace spaces with hyphens
            ),
            '-+', '-', 'g'                                           -- Replace multiple hyphens with single
        )
    )
)
WHERE slug IS NULL OR slug = '';

-- Handle duplicate slugs by appending row number
WITH ranked_slugs AS (
    SELECT
        id,
        slug,
        ROW_NUMBER() OVER (PARTITION BY slug ORDER BY id) as rn
    FROM listings
    WHERE is_deleted = false AND slug IS NOT NULL
)
UPDATE listings l
SET slug = CASE
    WHEN rs.rn > 1 THEN CONCAT(rs.slug, '-', rs.rn)
    ELSE rs.slug
END
FROM ranked_slugs rs
WHERE l.id = rs.id AND rs.rn > 1;

-- Create function to auto-generate slug from title (для фикстур и тестов)
CREATE OR REPLACE FUNCTION generate_slug_from_title()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.slug IS NULL OR NEW.slug = '' THEN
        NEW.slug := LOWER(
            TRIM(
                BOTH '-' FROM
                REGEXP_REPLACE(
                    REGEXP_REPLACE(
                        REGEXP_REPLACE(NEW.title, '[^a-zA-Z0-9\s-]', '', 'g'),
                        '\s+', '-', 'g'
                    ),
                    '-+', '-', 'g'
                )
            )
        );
        -- Handle potential duplicates by appending random suffix
        IF EXISTS (SELECT 1 FROM listings WHERE slug = NEW.slug AND id != COALESCE(NEW.id, 0) AND is_deleted = false) THEN
            NEW.slug := NEW.slug || '-' || COALESCE(NEW.id, (EXTRACT(EPOCH FROM NOW()) * 1000)::BIGINT);
        END IF;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create BEFORE INSERT/UPDATE trigger
DROP TRIGGER IF EXISTS trigger_generate_slug ON listings;
CREATE TRIGGER trigger_generate_slug
    BEFORE INSERT OR UPDATE ON listings
    FOR EACH ROW
    EXECUTE FUNCTION generate_slug_from_title();

-- DON'T add NOT NULL constraint - trigger will ensure slug is always generated
-- This allows fixtures to work without specifying slug manually
-- ALTER TABLE listings ALTER COLUMN slug SET NOT NULL;

-- Create unique index on slug (excluding deleted listings)
CREATE UNIQUE INDEX IF NOT EXISTS idx_listings_slug
ON listings(slug) WHERE is_deleted = false;

-- Add regular index for all slugs (for lookups)
CREATE INDEX IF NOT EXISTS idx_listings_slug_all
ON listings(slug);

-- Add index for expires_at queries (find expiring/expired listings)
CREATE INDEX IF NOT EXISTS idx_listings_expires_at
ON listings(expires_at) WHERE expires_at IS NOT NULL AND is_deleted = false;

-- Add compound index for active listings with expiration
CREATE INDEX IF NOT EXISTS idx_listings_status_expires_at
ON listings(status, expires_at) WHERE is_deleted = false AND expires_at IS NOT NULL;

-- Add comment for documentation
COMMENT ON COLUMN listings.slug IS 'SEO-friendly URL slug, unique across active listings';
COMMENT ON COLUMN listings.expires_at IS 'Optional expiration timestamp for C2C listings (auto-archive after this time)';
