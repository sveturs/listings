-- Migration: Create Inventory Reservations Table
-- Phase: Phase 17 - Orders Migration (Day 3-4)
-- Date: 2025-11-14
--
-- Purpose: Create inventory_reservations table for temporary stock holds
--          Prevents overselling during checkout process (cart → payment → order)
--
-- Features:
-- - Temporary holds: reserve stock for pending orders (TTL: 30 minutes)
-- - Status lifecycle: active → committed (payment success) or released (timeout/cancel)
-- - Auto-expiry: expires_at timestamp for cleanup jobs
-- - Variant support: optional variant_id for products with variants
-- - Order linkage: tracks which order this reservation belongs to

-- =====================================================
-- INVENTORY_RESERVATIONS TABLE
-- =====================================================

CREATE TABLE IF NOT EXISTS inventory_reservations (
    id BIGSERIAL PRIMARY KEY,

    -- Product references
    listing_id BIGINT NOT NULL,            -- FK to listings
    variant_id BIGINT NULL,                -- FK to b2c_product_variants (optional, for variant products)

    -- Order association (NULL until order created)
    order_id BIGINT NULL,                  -- FK to orders (SET NULL on delete)

    -- Reservation details
    quantity INTEGER NOT NULL,             -- Number of items reserved
    status VARCHAR(20) NOT NULL,           -- active, committed, released, expired

    -- Expiry tracking (TTL mechanism)
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL, -- When reservation expires (default: NOW() + 30 minutes)

    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    -- Foreign Keys
    CONSTRAINT fk_inventory_reservations_listing FOREIGN KEY (listing_id)
        REFERENCES listings(id) ON DELETE CASCADE,

    CONSTRAINT fk_inventory_reservations_variant FOREIGN KEY (variant_id)
        REFERENCES b2c_product_variants(id) ON DELETE CASCADE,

    CONSTRAINT fk_inventory_reservations_order FOREIGN KEY (order_id)
        REFERENCES orders(id) ON DELETE SET NULL,

    -- Business Logic Constraints
    CONSTRAINT chk_inventory_reservations_quantity_positive CHECK (quantity > 0),
    CONSTRAINT chk_inventory_reservations_status CHECK (
        status IN ('active', 'committed', 'released', 'expired')
    )
);

-- =====================================================
-- INDEXES
-- =====================================================

-- Primary lookup: find active reservations for a listing
CREATE INDEX idx_inventory_reservations_listing_status
    ON inventory_reservations(listing_id, status)
    WHERE status = 'active';

-- Primary lookup: find active reservations for a variant
CREATE INDEX idx_inventory_reservations_variant_status
    ON inventory_reservations(variant_id, status)
    WHERE status = 'active' AND variant_id IS NOT NULL;

-- Cleanup job: find expired reservations
CREATE INDEX idx_inventory_reservations_expires_at
    ON inventory_reservations(expires_at, status)
    WHERE status = 'active';

-- Lookup by order (find reservation for an order)
CREATE INDEX idx_inventory_reservations_order_id
    ON inventory_reservations(order_id)
    WHERE order_id IS NOT NULL;

-- General listing lookup (for analytics)
CREATE INDEX idx_inventory_reservations_listing_id
    ON inventory_reservations(listing_id);

-- =====================================================
-- AUTO-UPDATE TRIGGER
-- =====================================================

CREATE OR REPLACE FUNCTION update_inventory_reservations_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_inventory_reservations_updated_at
BEFORE UPDATE ON inventory_reservations
FOR EACH ROW
EXECUTE FUNCTION update_inventory_reservations_updated_at();

-- =====================================================
-- COMMENTS
-- =====================================================

COMMENT ON TABLE inventory_reservations IS
    'Temporary stock holds for pending orders. Prevents overselling during checkout. Auto-expires after 30 min.';

COMMENT ON COLUMN inventory_reservations.listing_id IS
    'Product listing being reserved. Always required.';

COMMENT ON COLUMN inventory_reservations.variant_id IS
    'Product variant (optional). For products with size/color/etc variants. NULL for simple products.';

COMMENT ON COLUMN inventory_reservations.order_id IS
    'Order this reservation belongs to. NULL until order created. SET NULL if order deleted.';

COMMENT ON COLUMN inventory_reservations.quantity IS
    'Number of items reserved. Must be > 0. Deducted from available stock.';

COMMENT ON COLUMN inventory_reservations.status IS
    'Lifecycle: active (pending) → committed (paid) or released (timeout/cancel) or expired (cleanup job).';

COMMENT ON COLUMN inventory_reservations.expires_at IS
    'When reservation expires. Default: NOW() + 30 minutes. Cleanup job runs every 5 minutes.';

COMMENT ON COLUMN inventory_reservations.created_at IS
    'When reservation was created (e.g., when user starts checkout).';

COMMENT ON COLUMN inventory_reservations.updated_at IS
    'Last status change (active → committed/released/expired).';

-- =====================================================
-- HELPER FUNCTION: Calculate Available Stock
-- =====================================================

COMMENT ON TABLE inventory_reservations IS
    'Temporary stock holds for pending orders. Prevents overselling during checkout. Auto-expires after 30 min.

Available stock calculation:
  available = listings.stock - SUM(reservations.quantity WHERE status = ''active'' AND expires_at > NOW())

Example cleanup job (run every 5 minutes):
  UPDATE inventory_reservations
  SET status = ''expired''
  WHERE status = ''active'' AND expires_at < NOW();
';
