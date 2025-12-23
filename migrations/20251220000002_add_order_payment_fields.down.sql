-- Migration rollback: Remove payment gateway fields from orders table
-- Author: Claude (Elite Architect)
-- Date: 2025-12-20

BEGIN;

-- Drop indexes
DROP INDEX IF EXISTS idx_orders_payment_intent;
DROP INDEX IF EXISTS idx_orders_provider_session;
DROP INDEX IF EXISTS idx_orders_idempotency_key;
DROP INDEX IF EXISTS idx_orders_payment_session;

-- Drop columns
ALTER TABLE orders DROP COLUMN IF EXISTS payment_idempotency_key;
ALTER TABLE orders DROP COLUMN IF EXISTS payment_intent_id;
ALTER TABLE orders DROP COLUMN IF EXISTS payment_session_id;
ALTER TABLE orders DROP COLUMN IF EXISTS payment_provider;

COMMIT;
