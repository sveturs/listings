-- Migration: Add 'completed' and 'cod_pending' to payment_status CHECK constraint
-- Date: 2025-11-23
--
-- Purpose: Add payment statuses to support proper order lifecycle:
--   - 'completed': Payment was successfully received (for online payments)
--   - 'cod_pending': Cash on Delivery - order confirmed, payment will be collected at delivery
--
-- COD orders use 'cod_pending' (NOT 'completed') because payment hasn't happened yet!
-- Only when delivery is complete and cash is collected, status changes to 'completed'.
--
-- Current allowed values: 'pending', 'paid', 'failed', 'refunded', 'partially_refunded'
-- New allowed values: 'pending', 'paid', 'completed', 'cod_pending', 'failed', 'refunded', 'partially_refunded'

-- Drop existing constraint
ALTER TABLE orders DROP CONSTRAINT IF EXISTS chk_orders_payment_status;

-- Add new constraint with 'completed' and 'cod_pending' values
ALTER TABLE orders ADD CONSTRAINT chk_orders_payment_status CHECK (
    payment_status IN ('pending', 'paid', 'completed', 'cod_pending', 'failed', 'refunded', 'partially_refunded')
);

-- Update comment
COMMENT ON COLUMN orders.payment_status IS
    'Payment tracking: pending → paid/completed/cod_pending (or failed/refunded/partially_refunded). ''cod_pending'' is for COD orders awaiting payment at delivery.';
