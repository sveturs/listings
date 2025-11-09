-- Migration: Fix b2c_products schema compatibility
-- Phase: 13.1.7 - Repository/Service Layer Migration Completion
-- Date: 2025-11-08
--
-- Purpose: Add all missing columns from b2c_products that are used by repository code
--          This migration ensures full compatibility between old b2c_products code
--          and new unified listings table
--
-- Changes:
-- 1. Rename views_count → view_count (align with b2c_products naming)
-- 2. Add sold_count (sales counter)
-- 3. Add location fields (has_individual_location, individual_address, etc.)
-- 4. Add show_on_map (map visibility flag)
-- 5. Add has_variants (product variants support flag)

-- 1. Rename views_count to view_count (align with b2c_products)
ALTER TABLE listings RENAME COLUMN views_count TO view_count;

-- 2. Add sold_count column
ALTER TABLE listings
ADD COLUMN IF NOT EXISTS sold_count INTEGER DEFAULT 0 NOT NULL;

-- 3. Add individual location support fields
ALTER TABLE listings
ADD COLUMN IF NOT EXISTS has_individual_location BOOLEAN DEFAULT FALSE,
ADD COLUMN IF NOT EXISTS individual_address TEXT,
ADD COLUMN IF NOT EXISTS individual_latitude NUMERIC(10,8),
ADD COLUMN IF NOT EXISTS individual_longitude NUMERIC(11,8),
ADD COLUMN IF NOT EXISTS location_privacy VARCHAR(20) DEFAULT 'exact'
    CHECK (location_privacy IN ('exact', 'approximate', 'hidden')),
ADD COLUMN IF NOT EXISTS show_on_map BOOLEAN DEFAULT TRUE;

-- 4. Add product variants support flag
ALTER TABLE listings
ADD COLUMN IF NOT EXISTS has_variants BOOLEAN DEFAULT FALSE;

-- 5. Create indexes for new columns
CREATE INDEX IF NOT EXISTS idx_listings_view_count
ON listings(view_count DESC) WHERE is_deleted = false;

CREATE INDEX IF NOT EXISTS idx_listings_sold_count
ON listings(sold_count DESC) WHERE is_deleted = false;

CREATE INDEX IF NOT EXISTS idx_listings_location
ON listings(individual_latitude, individual_longitude)
WHERE individual_latitude IS NOT NULL
  AND individual_longitude IS NOT NULL
  AND has_individual_location = TRUE;

CREATE INDEX IF NOT EXISTS idx_listings_has_variants
ON listings(has_variants) WHERE has_variants = TRUE;

-- 6. Add comments
COMMENT ON COLUMN listings.view_count IS 'Number of times this listing has been viewed (renamed from views_count for b2c compatibility)';
COMMENT ON COLUMN listings.sold_count IS 'Number of times this listing has been sold';
COMMENT ON COLUMN listings.has_individual_location IS 'If true, product has its own location independent of storefront';
COMMENT ON COLUMN listings.individual_address IS 'Individual address for products with custom location';
COMMENT ON COLUMN listings.individual_latitude IS 'Latitude for individual product location';
COMMENT ON COLUMN listings.individual_longitude IS 'Longitude for individual product location';
COMMENT ON COLUMN listings.location_privacy IS 'Privacy level for location display (exact, approximate, hidden)';
COMMENT ON COLUMN listings.show_on_map IS 'Whether to show this listing on map interface';
COMMENT ON COLUMN listings.has_variants IS 'Whether this product has variants (sizes, colors, etc.)';
