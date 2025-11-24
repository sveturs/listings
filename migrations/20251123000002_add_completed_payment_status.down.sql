-- Migration: Revert 'completed' payment status
-- Date: 2025-11-23
--
-- WARNING: This will fail if any orders have payment_status = 'completed'
-- Before running, ensure: SELECT COUNT(*) FROM orders WHERE payment_status = 'completed';

-- Drop new constraint
ALTER TABLE orders DROP CONSTRAINT IF EXISTS chk_orders_payment_status;

-- Restore original constraint without 'completed'
ALTER TABLE orders ADD CONSTRAINT chk_orders_payment_status CHECK (
    payment_status IN ('pending', 'paid', 'failed', 'refunded', 'partially_refunded')
);

-- Restore original comment
COMMENT ON COLUMN orders.payment_status IS
    'Payment tracking: pending → paid (or failed/refunded/partially_refunded).';
