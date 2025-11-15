-- Add payment_completed_at column to orders table
-- This column tracks when payment was completed for an order

ALTER TABLE orders
ADD COLUMN IF NOT EXISTS payment_completed_at TIMESTAMP WITH TIME ZONE;

-- Add index for querying orders by payment completion time
CREATE INDEX IF NOT EXISTS idx_orders_payment_completed_at
ON orders(payment_completed_at)
WHERE payment_completed_at IS NOT NULL;

-- Add comment
COMMENT ON COLUMN orders.payment_completed_at IS 'Timestamp when payment was completed';
