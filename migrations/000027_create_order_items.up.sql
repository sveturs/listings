-- Migration: Create Order Items Table
-- Phase: Phase 17 - Orders Migration (Day 3-4)
-- Date: 2025-11-14
--
-- Purpose: Create order_items table to store snapshot of purchased products
--          Immutable historical record - preserves product details even if listing is deleted
--
-- Features:
-- - Product snapshot: name, SKU, variant details, attributes stored at order time
-- - Price immutability: price/total frozen at checkout (no retroactive changes)
-- - Variant support: stores variant_data JSONB for size/color/etc
-- - Attributes snapshot: full product attributes at purchase time
-- - Cascade delete: items deleted when order is deleted
-- - Handles deleted listings: references listing_id but data is self-contained

-- =====================================================
-- ORDER_ITEMS TABLE
-- =====================================================

CREATE TABLE IF NOT EXISTS order_items (
    id BIGSERIAL PRIMARY KEY,

    -- Order association
    order_id BIGINT NOT NULL,              -- FK to orders (CASCADE on delete)

    -- Product references (may become invalid if listing deleted)
    listing_id BIGINT NOT NULL,            -- Reference to listing (for analytics, even if deleted)
    variant_id BIGINT NULL,                -- Reference to variant (optional)

    -- Product snapshot (preserved even if listing deleted)
    listing_name VARCHAR(255) NOT NULL,    -- Product name at order time
    sku VARCHAR(100),                      -- Product SKU at order time
    variant_data JSONB,                    -- Variant details: { size: "L", color: "Red", ... }
    attributes JSONB,                      -- Full product attributes at order time

    -- Item details
    quantity INTEGER NOT NULL,             -- Number of items purchased
    price NUMERIC(10,2) NOT NULL,          -- Price per unit at order time
    total NUMERIC(10,2) NOT NULL,          -- Total for this line item: price * quantity

    -- Timestamp (no updated_at - order items are immutable)
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    -- Foreign Keys
    CONSTRAINT fk_order_items_order FOREIGN KEY (order_id)
        REFERENCES orders(id) ON DELETE CASCADE,

    -- Note: NO FK to listings/variants - they might be deleted
    -- We keep listing_id/variant_id for analytics, but rely on snapshot data

    -- Business Logic Constraints
    CONSTRAINT chk_order_items_quantity_positive CHECK (quantity > 0),
    CONSTRAINT chk_order_items_price_non_negative CHECK (price >= 0),
    CONSTRAINT chk_order_items_total_non_negative CHECK (total >= 0)
);

-- =====================================================
-- INDEXES
-- =====================================================

-- Primary lookup: find all items in an order
CREATE INDEX idx_order_items_order_id
    ON order_items(order_id);

-- Analytics: find all orders for a listing (even if listing deleted)
CREATE INDEX idx_order_items_listing_id
    ON order_items(listing_id);

-- Analytics: find all orders for a variant
CREATE INDEX idx_order_items_variant_id
    ON order_items(variant_id)
    WHERE variant_id IS NOT NULL;

-- Reporting: find orders by product name (fuzzy search)
CREATE INDEX idx_order_items_listing_name
    ON order_items USING gin(to_tsvector('simple', listing_name));

-- =====================================================
-- COMMENTS
-- =====================================================

COMMENT ON TABLE order_items IS
    'Snapshot of purchased products. Immutable historical record preserved even if listing is deleted.';

COMMENT ON COLUMN order_items.order_id IS
    'Order this item belongs to. Items cascade-deleted when order is deleted.';

COMMENT ON COLUMN order_items.listing_id IS
    'Reference to listing (for analytics). Listing may be deleted, rely on snapshot data.';

COMMENT ON COLUMN order_items.variant_id IS
    'Reference to variant (optional). Variant may be deleted, rely on variant_data snapshot.';

COMMENT ON COLUMN order_items.listing_name IS
    'Product name at order time. Preserved even if listing renamed or deleted.';

COMMENT ON COLUMN order_items.sku IS
    'Product SKU at order time. For inventory reconciliation and seller tracking.';

COMMENT ON COLUMN order_items.variant_data IS
    'JSONB snapshot of variant attributes: { size: "L", color: "Red", material: "Cotton", ... }';

COMMENT ON COLUMN order_items.attributes IS
    'JSONB snapshot of full product attributes at order time. Preserves exact product state.';

COMMENT ON COLUMN order_items.quantity IS
    'Number of items purchased. Immutable after order creation.';

COMMENT ON COLUMN order_items.price IS
    'Price per unit at order time. Frozen at checkout, not affected by future price changes.';

COMMENT ON COLUMN order_items.total IS
    'Total for this line item: price * quantity. Immutable after order creation.';

COMMENT ON COLUMN order_items.created_at IS
    'When order item was created. No updated_at - order items never change after creation.';
