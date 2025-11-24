-- Migration: Create Orders Table
-- Phase: Phase 17 - Orders Migration (Day 3-4)
-- Date: 2025-11-14
--
-- Purpose: Create orders table for completed purchases
--          Supports full e-commerce lifecycle: payment, shipping, escrow, cancellation
--
-- Features:
-- - Unique order number (e.g., ORD-2025-001234)
-- - Guest checkout support (user_id can be NULL)
-- - Financial tracking: subtotal, tax, shipping, discount, commission
-- - Payment integration: payment_method, payment_status, transaction_id
-- - Escrow system: hold payment for N days (default 3)
-- - Shipping integration: tracking_number, shipping_provider, shipment_id
-- - JSONB addresses: flexible schema for different countries
-- - Cancellation tracking: cancelled_at, cancellation_reason
-- - Status workflow: pending → confirmed → processing → shipped → delivered

-- =====================================================
-- ORDERS TABLE
-- =====================================================

CREATE TABLE IF NOT EXISTS orders (
    id BIGSERIAL PRIMARY KEY,

    -- Order identification
    order_number VARCHAR(50) UNIQUE NOT NULL, -- e.g., ORD-2025-001234

    -- User association (NULL for guest orders)
    user_id BIGINT NULL,                   -- FK to auth service (authenticated checkout)

    -- Storefront association
    storefront_id BIGINT NOT NULL,         -- FK to storefronts

    -- Order status workflow
    status VARCHAR(20) NOT NULL,           -- pending, confirmed, processing, shipped, delivered, cancelled, refunded, failed

    -- Payment details
    payment_status VARCHAR(20) NOT NULL,   -- pending, paid, failed, refunded, partially_refunded
    payment_method VARCHAR(50),            -- cash, card, bank_transfer, paypal, keks_pay, ips, etc.
    payment_transaction_id VARCHAR(255),   -- Payment gateway transaction ID

    -- Financial breakdown
    subtotal NUMERIC(10,2) NOT NULL,       -- Sum of all items (price * quantity)
    tax NUMERIC(10,2) DEFAULT 0.00,        -- VAT or sales tax
    shipping NUMERIC(10,2) DEFAULT 0.00,   -- Shipping cost
    discount NUMERIC(10,2) DEFAULT 0.00,   -- Discount amount (coupons, promotions)
    total NUMERIC(10,2) NOT NULL,          -- Final amount: subtotal + tax + shipping - discount

    -- Platform commission (marketplace fee)
    commission NUMERIC(10,2) DEFAULT 0.00, -- Platform fee (e.g., 5% of subtotal)
    seller_amount NUMERIC(10,2) NOT NULL,  -- Amount seller receives: total - commission

    -- Currency
    currency VARCHAR(3) DEFAULT 'RSD',     -- ISO 4217 currency code (RSD, EUR, USD)

    -- Address details (JSONB for flexibility)
    shipping_address JSONB,                -- { name, street, city, postal_code, country, phone }
    billing_address JSONB,                 -- { name, street, city, postal_code, country, phone }

    -- Shipping details
    shipping_method VARCHAR(100),          -- e.g., "Standard Shipping", "Express Delivery"
    shipping_provider VARCHAR(100),        -- e.g., "posta_srbije", "aks", "bex", "d_express"
    tracking_number VARCHAR(255),          -- Tracking number from shipping provider

    -- Escrow system (hold payment until delivery confirmed)
    escrow_release_date TIMESTAMP WITH TIME ZONE, -- When to release payment to seller
    escrow_days INTEGER DEFAULT 3,         -- Number of days to hold payment (default: 3)

    -- Integration with Delivery Service
    shipment_id BIGINT NULL,               -- FK to Delivery Service (external microservice)

    -- Additional information
    notes TEXT,                            -- Customer notes, special instructions

    -- Cancellation tracking
    cancelled_at TIMESTAMP WITH TIME ZONE,
    cancellation_reason TEXT,

    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    -- Foreign Keys
    CONSTRAINT fk_orders_storefront FOREIGN KEY (storefront_id)
        REFERENCES storefronts(id) ON DELETE RESTRICT,

    -- Business Logic Constraints
    CONSTRAINT chk_orders_status CHECK (
        status IN ('pending', 'confirmed', 'processing', 'shipped', 'delivered', 'cancelled', 'refunded', 'failed')
    ),
    CONSTRAINT chk_orders_payment_status CHECK (
        payment_status IN ('pending', 'paid', 'failed', 'refunded', 'partially_refunded')
    ),
    CONSTRAINT chk_orders_subtotal_positive CHECK (subtotal >= 0),
    CONSTRAINT chk_orders_total_positive CHECK (total >= 0),
    CONSTRAINT chk_orders_seller_amount_positive CHECK (seller_amount >= 0),
    CONSTRAINT chk_orders_escrow_days_positive CHECK (escrow_days >= 0)
);

-- =====================================================
-- INDEXES
-- =====================================================

-- Primary lookup: find orders by user
CREATE INDEX idx_orders_user_id
    ON orders(user_id)
    WHERE user_id IS NOT NULL;

-- Primary lookup: find orders by storefront (seller dashboard)
CREATE INDEX idx_orders_storefront_id
    ON orders(storefront_id);

-- Unique lookup: find order by order_number
CREATE INDEX idx_orders_order_number
    ON orders(order_number);

-- Filter by status (admin/seller dashboards)
CREATE INDEX idx_orders_status
    ON orders(status);

-- Filter by payment status (finance tracking)
CREATE INDEX idx_orders_payment_status
    ON orders(payment_status);

-- Sort by creation date (recent orders)
CREATE INDEX idx_orders_created_at_desc
    ON orders(created_at DESC);

-- Escrow release job: find orders where escrow_release_date has passed
CREATE INDEX idx_orders_escrow_release_date
    ON orders(escrow_release_date)
    WHERE escrow_release_date IS NOT NULL AND payment_status = 'paid';

-- Delivery integration: find order by shipment_id
CREATE INDEX idx_orders_shipment_id
    ON orders(shipment_id)
    WHERE shipment_id IS NOT NULL;

-- Composite index: user orders sorted by date
CREATE INDEX idx_orders_user_created_at
    ON orders(user_id, created_at DESC)
    WHERE user_id IS NOT NULL;

-- Composite index: storefront orders sorted by date
CREATE INDEX idx_orders_storefront_created_at
    ON orders(storefront_id, created_at DESC);

-- =====================================================
-- AUTO-UPDATE TRIGGER
-- =====================================================

CREATE OR REPLACE FUNCTION update_orders_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_orders_updated_at
BEFORE UPDATE ON orders
FOR EACH ROW
EXECUTE FUNCTION update_orders_updated_at();

-- =====================================================
-- COMMENTS
-- =====================================================

COMMENT ON TABLE orders IS
    'Completed purchase orders. Supports authenticated and guest checkout, escrow, shipping tracking.';

COMMENT ON COLUMN orders.order_number IS
    'Unique order identifier shown to customer (e.g., ORD-2025-001234).';

COMMENT ON COLUMN orders.user_id IS
    'Authenticated user ID (FK to auth service). NULL for guest checkout.';

COMMENT ON COLUMN orders.status IS
    'Order lifecycle: pending → confirmed → processing → shipped → delivered (or cancelled/refunded/failed).';

COMMENT ON COLUMN orders.payment_status IS
    'Payment tracking: pending → paid (or failed/refunded/partially_refunded).';

COMMENT ON COLUMN orders.commission IS
    'Platform commission (marketplace fee). Deducted from seller_amount.';

COMMENT ON COLUMN orders.seller_amount IS
    'Amount seller receives: total - commission. Released after escrow period.';

COMMENT ON COLUMN orders.shipping_address IS
    'JSONB: { name, street, city, postal_code, country, phone, ... }. Flexible schema.';

COMMENT ON COLUMN orders.billing_address IS
    'JSONB: { name, street, city, postal_code, country, phone, ... }. Flexible schema.';

COMMENT ON COLUMN orders.escrow_release_date IS
    'When to release payment to seller. Calculated as: order created_at + escrow_days.';

COMMENT ON COLUMN orders.escrow_days IS
    'Number of days to hold payment (default: 3). Protects buyer, allows returns/disputes.';

COMMENT ON COLUMN orders.shipment_id IS
    'FK to Delivery Service microservice. NULL if not using delivery service.';

COMMENT ON COLUMN orders.cancelled_at IS
    'When order was cancelled. NULL if not cancelled.';

COMMENT ON COLUMN orders.cancellation_reason IS
    'Why order was cancelled (customer request, out of stock, fraud, etc).';
