-- Migration: Add Listing Variants Table
-- Date: 2025-12-16
-- Purpose: Add variants support for listings (size, color, etc.)
--          Required by cart_repository.go which JOINs with this table

-- =====================================================
-- 1. LISTING_VARIANTS TABLE
-- =====================================================

CREATE TABLE listing_variants (
    id              BIGSERIAL PRIMARY KEY,
    listing_id      BIGINT NOT NULL,
    sku             VARCHAR(100),

    -- Variant attributes (e.g., {"color": "red", "size": "M"})
    attributes      JSONB NOT NULL DEFAULT '{}',

    -- Pricing (NULL means use listing price)
    price           DECIMAL(12, 2),

    -- Stock management
    stock           INT NOT NULL DEFAULT 0,

    -- Media
    image_url       VARCHAR(500),

    -- Status
    is_active       BOOLEAN NOT NULL DEFAULT true,

    -- Timestamps
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Foreign Keys
    CONSTRAINT fk_listing_variants_listing FOREIGN KEY (listing_id)
        REFERENCES listings(id) ON DELETE CASCADE
);

-- =====================================================
-- 2. INDEXES
-- =====================================================

CREATE INDEX idx_listing_variants_listing ON listing_variants(listing_id);
CREATE INDEX idx_listing_variants_sku ON listing_variants(sku) WHERE sku IS NOT NULL;
CREATE INDEX idx_listing_variants_active ON listing_variants(listing_id) WHERE is_active = true;
CREATE INDEX idx_listing_variants_attributes ON listing_variants USING gin(attributes);

-- =====================================================
-- 3. UNIQUE CONSTRAINT
-- =====================================================

-- Each listing can only have one variant with the same attributes
CREATE UNIQUE INDEX idx_listing_variants_unique_attrs ON listing_variants(listing_id, md5(attributes::text));

-- =====================================================
-- 4. UPDATE TRIGGER
-- =====================================================

CREATE OR REPLACE FUNCTION update_listing_variants_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_listing_variants_updated_at
    BEFORE UPDATE ON listing_variants
    FOR EACH ROW
    EXECUTE FUNCTION update_listing_variants_updated_at();

-- =====================================================
-- 5. COMMENTS
-- =====================================================

COMMENT ON TABLE listing_variants IS 'Product variants for listings (size, color, etc.)';
COMMENT ON COLUMN listing_variants.attributes IS 'JSON object with variant attributes, e.g., {"color": "red", "size": "M"}';
COMMENT ON COLUMN listing_variants.price IS 'Variant-specific price, NULL means use listing base price';
COMMENT ON COLUMN listing_variants.stock IS 'Available stock for this specific variant';
COMMENT ON COLUMN listing_variants.image_url IS 'Optional variant-specific image';
