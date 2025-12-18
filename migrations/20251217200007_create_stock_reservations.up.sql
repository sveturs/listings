-- Migration: Create stock_reservations table (Phase 3: Variants System)
-- Description: Stock reservation system for order processing and inventory management
-- Author: Claude (Elite Architect)
-- Date: 2025-12-17
-- Phase: 3 (Variants)
-- Task: BE-3.3

BEGIN;

-- ============================================================================
-- TABLE: stock_reservations - Temporary stock reservations for orders
-- ============================================================================
CREATE TABLE IF NOT EXISTS stock_reservations (
    -- Primary key
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Foreign keys
    variant_id UUID NOT NULL,
    order_id UUID NOT NULL,

    -- Reservation details
    quantity INTEGER NOT NULL CHECK (quantity > 0),

    -- Expiration
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,

    -- Status
    status VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (
        status IN ('active', 'confirmed', 'cancelled', 'expired')
    ),

    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,

    -- Constraints
    CONSTRAINT positive_quantity CHECK (quantity > 0)
);

-- ============================================================================
-- INDEXES for stock_reservations
-- ============================================================================

-- Primary lookups
CREATE INDEX idx_reservations_variant ON stock_reservations(variant_id, status);
CREATE INDEX idx_reservations_order ON stock_reservations(order_id);

-- Expiration management
CREATE INDEX idx_reservations_expires ON stock_reservations(expires_at, status)
WHERE status = 'active';

-- Active reservations per variant
CREATE INDEX idx_reservations_active ON stock_reservations(variant_id, status, quantity)
WHERE status = 'active';

-- Composite index for reservation queries
CREATE INDEX idx_reservations_composite ON stock_reservations(variant_id, order_id, status);

-- ============================================================================
-- TRIGGERS
-- ============================================================================

-- Auto-update updated_at timestamp
CREATE TRIGGER trigger_reservations_updated_at
    BEFORE UPDATE ON stock_reservations
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Auto-expire old reservations
CREATE OR REPLACE FUNCTION auto_expire_reservations()
RETURNS TRIGGER AS $$
BEGIN
    -- Mark reservation as expired if past expiration time
    IF NEW.status = 'active' AND NEW.expires_at < CURRENT_TIMESTAMP THEN
        NEW.status = 'expired';
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_reservations_auto_expire
    BEFORE UPDATE ON stock_reservations
    FOR EACH ROW
    WHEN (NEW.status = 'active')
    EXECUTE FUNCTION auto_expire_reservations();

-- Update variant reserved_quantity when reservation is created/updated
CREATE OR REPLACE FUNCTION sync_variant_reserved_quantity()
RETURNS TRIGGER AS $$
DECLARE
    old_quantity INTEGER := 0;
    new_quantity INTEGER := 0;
BEGIN
    -- Handle INSERT
    IF (TG_OP = 'INSERT') THEN
        IF NEW.status = 'active' THEN
            UPDATE product_variants
            SET reserved_quantity = reserved_quantity + NEW.quantity
            WHERE id = NEW.variant_id;
        END IF;
        RETURN NEW;
    END IF;

    -- Handle UPDATE
    IF (TG_OP = 'UPDATE') THEN
        -- Calculate quantity delta
        IF OLD.status = 'active' THEN
            old_quantity := OLD.quantity;
        END IF;

        IF NEW.status = 'active' THEN
            new_quantity := NEW.quantity;
        END IF;

        -- Update variant reserved_quantity
        IF old_quantity != new_quantity THEN
            UPDATE product_variants
            SET reserved_quantity = reserved_quantity - old_quantity + new_quantity
            WHERE id = NEW.variant_id;
        END IF;

        RETURN NEW;
    END IF;

    -- Handle DELETE
    IF (TG_OP = 'DELETE') THEN
        IF OLD.status = 'active' THEN
            UPDATE product_variants
            SET reserved_quantity = reserved_quantity - OLD.quantity
            WHERE id = OLD.variant_id;
        END IF;
        RETURN OLD;
    END IF;

    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_reservations_sync_quantity
    AFTER INSERT OR UPDATE OR DELETE ON stock_reservations
    FOR EACH ROW
    EXECUTE FUNCTION sync_variant_reserved_quantity();

-- ============================================================================
-- COMMENTS
-- ============================================================================

COMMENT ON TABLE stock_reservations IS 'Temporary stock reservations for order processing';
COMMENT ON COLUMN stock_reservations.variant_id IS 'Reference to product_variants.id';
COMMENT ON COLUMN stock_reservations.order_id IS 'Reference to order (UUID)';
COMMENT ON COLUMN stock_reservations.quantity IS 'Reserved quantity (must be positive)';
COMMENT ON COLUMN stock_reservations.expires_at IS 'Reservation expiration timestamp';
COMMENT ON COLUMN stock_reservations.status IS 'Reservation status: active, confirmed, cancelled, expired';

-- ============================================================================
-- HELPER FUNCTION: Cleanup expired reservations
-- ============================================================================

CREATE OR REPLACE FUNCTION cleanup_expired_reservations()
RETURNS TABLE(cleaned_count BIGINT) AS $$
DECLARE
    result BIGINT;
BEGIN
    -- Update expired reservations
    WITH expired AS (
        UPDATE stock_reservations
        SET status = 'expired'
        WHERE status = 'active' AND expires_at < CURRENT_TIMESTAMP
        RETURNING id
    )
    SELECT COUNT(*) INTO result FROM expired;

    RETURN QUERY SELECT result;
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION cleanup_expired_reservations() IS 'Cleanup expired stock reservations (run as cron job)';

-- ============================================================================
-- FOREIGN KEY to product_variants (will be added after variant table exists)
-- ============================================================================
-- NOTE: This FK is commented out here and will be enabled in a separate migration
-- after product_variants table is fully populated
-- ALTER TABLE stock_reservations
--     ADD CONSTRAINT stock_reservations_variant_fk
--     FOREIGN KEY (variant_id) REFERENCES product_variants(id) ON DELETE CASCADE;

COMMIT;
