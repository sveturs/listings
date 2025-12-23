-- Migration: Add payment gateway fields to orders table
-- Description: Add payment_provider, payment_session_id, payment_intent_id, payment_idempotency_key
-- Author: Claude (Elite Architect)
-- Date: 2025-12-20
-- Phase: Payment Integration
-- Task: Payment Gateway Support

BEGIN;

-- ============================================================================
-- ALTER TABLE: orders - Add payment gateway fields
-- ============================================================================

-- Payment provider (stripe, allsecure, null for offline/COD)
ALTER TABLE orders
    ADD COLUMN IF NOT EXISTS payment_provider VARCHAR(50);

COMMENT ON COLUMN orders.payment_provider IS 'Payment gateway provider: stripe, allsecure, or null for offline/COD';

-- Checkout session ID from payment gateway
ALTER TABLE orders
    ADD COLUMN IF NOT EXISTS payment_session_id VARCHAR(255);

COMMENT ON COLUMN orders.payment_session_id IS 'Checkout session ID from payment gateway (e.g., Stripe checkout session)';

-- Payment intent ID after successful payment
ALTER TABLE orders
    ADD COLUMN IF NOT EXISTS payment_intent_id VARCHAR(255);

COMMENT ON COLUMN orders.payment_intent_id IS 'Payment intent ID from gateway after successful payment';

-- Idempotency key for duplicate prevention
ALTER TABLE orders
    ADD COLUMN IF NOT EXISTS payment_idempotency_key VARCHAR(255);

COMMENT ON COLUMN orders.payment_idempotency_key IS 'Idempotency key for duplicate payment prevention';

-- ============================================================================
-- INDEXES for payment fields
-- ============================================================================

-- Index on payment_session_id for webhook lookups
CREATE INDEX IF NOT EXISTS idx_orders_payment_session
    ON orders(payment_session_id)
    WHERE payment_session_id IS NOT NULL;

COMMENT ON INDEX idx_orders_payment_session IS 'Fast lookup by payment session ID for webhook processing';

-- Unique index on idempotency key for duplicate prevention
CREATE UNIQUE INDEX IF NOT EXISTS idx_orders_idempotency_key
    ON orders(payment_idempotency_key)
    WHERE payment_idempotency_key IS NOT NULL;

COMMENT ON INDEX idx_orders_idempotency_key IS 'Unique constraint on idempotency key for duplicate payment prevention';

-- Composite index for payment provider + session
CREATE INDEX IF NOT EXISTS idx_orders_provider_session
    ON orders(payment_provider, payment_session_id)
    WHERE payment_provider IS NOT NULL AND payment_session_id IS NOT NULL;

COMMENT ON INDEX idx_orders_provider_session IS 'Fast lookup by provider and session for payment reconciliation';

-- Index on payment_intent_id for payment tracking
CREATE INDEX IF NOT EXISTS idx_orders_payment_intent
    ON orders(payment_intent_id)
    WHERE payment_intent_id IS NOT NULL;

COMMENT ON INDEX idx_orders_payment_intent IS 'Fast lookup by payment intent ID for payment tracking';

COMMIT;
