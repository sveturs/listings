-- Add status transition timestamps to orders table
ALTER TABLE orders ADD COLUMN confirmed_at TIMESTAMP WITH TIME ZONE;
ALTER TABLE orders ADD COLUMN shipped_at TIMESTAMP WITH TIME ZONE;
ALTER TABLE orders ADD COLUMN delivered_at TIMESTAMP WITH TIME ZONE;

-- Add indexes for performance (status timestamps are often queried)
CREATE INDEX idx_orders_confirmed_at ON orders(confirmed_at) WHERE confirmed_at IS NOT NULL;
CREATE INDEX idx_orders_shipped_at ON orders(shipped_at) WHERE shipped_at IS NOT NULL;
CREATE INDEX idx_orders_delivered_at ON orders(delivered_at) WHERE delivered_at IS NOT NULL;

-- Add comments
COMMENT ON COLUMN orders.confirmed_at IS 'Timestamp when order was confirmed';
COMMENT ON COLUMN orders.shipped_at IS 'Timestamp when order was shipped';
COMMENT ON COLUMN orders.delivered_at IS 'Timestamp when order was delivered';
