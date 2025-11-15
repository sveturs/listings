-- Add missing customer and admin columns to orders table

ALTER TABLE orders
ADD COLUMN IF NOT EXISTS customer_name VARCHAR(255),
ADD COLUMN IF NOT EXISTS customer_email VARCHAR(255),
ADD COLUMN IF NOT EXISTS customer_phone VARCHAR(50),
ADD COLUMN IF NOT EXISTS customer_notes TEXT,
ADD COLUMN IF NOT EXISTS admin_notes TEXT;

-- Add indexes for commonly queried fields
CREATE INDEX IF NOT EXISTS idx_orders_customer_email
ON orders(customer_email)
WHERE customer_email IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_orders_customer_phone
ON orders(customer_phone)
WHERE customer_phone IS NOT NULL;

-- Add comments
COMMENT ON COLUMN orders.customer_name IS 'Customer full name from shipping address';
COMMENT ON COLUMN orders.customer_email IS 'Customer email for order notifications';
COMMENT ON COLUMN orders.customer_phone IS 'Customer phone for delivery coordination';
COMMENT ON COLUMN orders.customer_notes IS 'Notes/instructions from customer';
COMMENT ON COLUMN orders.admin_notes IS 'Internal admin notes about the order';
