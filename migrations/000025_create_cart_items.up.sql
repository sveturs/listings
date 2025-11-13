-- Migration: Create Cart Items Table
-- Phase: Phase 17 - Orders Migration (Day 3-4)
-- Date: 2025-11-14
--
-- Purpose: Create cart_items table to store individual products in shopping carts
--          Supports both simple listings and variant-based products
--
-- Features:
-- - Price snapshot: stores price at add-to-cart time (for price change detection)
-- - Variant support: optional variant_id for products with variants
-- - Quantity validation: must be > 0
-- - Unique constraint: one item per listing/variant combination per cart
-- - Cascade delete: items deleted when cart is deleted

-- =====================================================
-- CART_ITEMS TABLE
-- =====================================================

CREATE TABLE IF NOT EXISTS cart_items (
    id BIGSERIAL PRIMARY KEY,

    -- Cart association
    cart_id BIGINT NOT NULL,               -- FK to shopping_carts (CASCADE on delete)

    -- Product references
    listing_id BIGINT NOT NULL,            -- FK to listings
    variant_id BIGINT NULL,                -- FK to listing_variants (optional, for variant products)

    -- Item details
    quantity INTEGER NOT NULL,             -- Number of items
    price_snapshot NUMERIC(10,2) NOT NULL, -- Price at add-to-cart time (for comparison)

    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    -- Foreign Keys
    CONSTRAINT fk_cart_items_cart FOREIGN KEY (cart_id)
        REFERENCES shopping_carts(id) ON DELETE CASCADE,

    CONSTRAINT fk_cart_items_listing FOREIGN KEY (listing_id)
        REFERENCES listings(id) ON DELETE CASCADE,

    CONSTRAINT fk_cart_items_variant FOREIGN KEY (variant_id)
        REFERENCES listing_variants(id) ON DELETE CASCADE,

    -- Business Logic Constraints
    CONSTRAINT chk_cart_items_quantity_positive CHECK (quantity > 0)
);

-- =====================================================
-- INDEXES
-- =====================================================

-- Primary lookup: find all items in a cart
CREATE INDEX idx_cart_items_cart_id
    ON cart_items(cart_id);

-- Lookup by listing (for inventory checks, price updates)
CREATE INDEX idx_cart_items_listing_id
    ON cart_items(listing_id);

-- Lookup by variant (for inventory checks)
CREATE INDEX idx_cart_items_variant_id
    ON cart_items(variant_id)
    WHERE variant_id IS NOT NULL;

-- =====================================================
-- UNIQUE CONSTRAINTS
-- =====================================================

-- Prevent duplicate items: one row per listing/variant combination per cart
-- If variant_id is NULL, this ensures one row per listing per cart
-- If variant_id is set, this ensures one row per variant per cart
CREATE UNIQUE INDEX idx_cart_items_unique_per_cart
    ON cart_items(cart_id, listing_id, COALESCE(variant_id, 0));

-- =====================================================
-- AUTO-UPDATE TRIGGER
-- =====================================================

CREATE OR REPLACE FUNCTION update_cart_items_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_cart_items_updated_at
BEFORE UPDATE ON cart_items
FOR EACH ROW
EXECUTE FUNCTION update_cart_items_updated_at();

-- =====================================================
-- COMMENTS
-- =====================================================

COMMENT ON TABLE cart_items IS
    'Individual items in shopping carts. Supports both simple listings and variant-based products.';

COMMENT ON COLUMN cart_items.cart_id IS
    'Shopping cart this item belongs to. Items cascade-deleted when cart is deleted.';

COMMENT ON COLUMN cart_items.listing_id IS
    'Product listing. Always required.';

COMMENT ON COLUMN cart_items.variant_id IS
    'Product variant (optional). For products with size/color/etc variants. NULL for simple products.';

COMMENT ON COLUMN cart_items.quantity IS
    'Number of items. Must be > 0. Updated when user changes quantity.';

COMMENT ON COLUMN cart_items.price_snapshot IS
    'Price per unit at add-to-cart time. Used to detect price changes before checkout.';

COMMENT ON COLUMN cart_items.updated_at IS
    'Auto-updated when quantity or price changes. Last modification timestamp.';
