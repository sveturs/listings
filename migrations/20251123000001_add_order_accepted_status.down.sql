-- Rollback: Remove accepted status columns

DROP INDEX IF EXISTS idx_orders_accepted_at;

ALTER TABLE orders DROP COLUMN IF EXISTS accepted_at;
ALTER TABLE orders DROP COLUMN IF EXISTS label_url;
ALTER TABLE orders DROP COLUMN IF EXISTS seller_notes;

-- Restore original CHECK constraint without 'accepted'
ALTER TABLE orders DROP CONSTRAINT IF EXISTS orders_status_check;
ALTER TABLE orders ADD CONSTRAINT orders_status_check
CHECK (status IN ('unspecified', 'pending', 'confirmed', 'processing', 'shipped', 'delivered', 'cancelled', 'refunded', 'failed'));
