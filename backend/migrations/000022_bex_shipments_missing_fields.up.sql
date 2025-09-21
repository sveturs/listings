-- Add missing fields to bex_shipments table
ALTER TABLE bex_shipments
ADD COLUMN IF NOT EXISTS marketplace_order_id INTEGER,
ADD COLUMN IF NOT EXISTS storefront_order_id BIGINT,
ADD COLUMN IF NOT EXISTS status_text VARCHAR(255),
ADD COLUMN IF NOT EXISTS registered_at TIMESTAMP WITH TIME ZONE,
ADD COLUMN IF NOT EXISTS delivered_at TIMESTAMP WITH TIME ZONE,
ADD COLUMN IF NOT EXISTS failed_reason TEXT,
ADD COLUMN IF NOT EXISTS status_history JSONB;

-- Add indexes for new fields
CREATE INDEX IF NOT EXISTS idx_bex_shipments_marketplace_order_id ON bex_shipments(marketplace_order_id);
CREATE INDEX IF NOT EXISTS idx_bex_shipments_storefront_order_id ON bex_shipments(storefront_order_id);

-- Update status column to be integer for compatibility with queries
-- First, create a new column
ALTER TABLE bex_shipments ADD COLUMN IF NOT EXISTS status_int INTEGER DEFAULT 1;

-- Map existing string statuses to integers if they exist
UPDATE bex_shipments SET
  status_int = CASE
    WHEN status = 'pending' THEN 1
    WHEN status = 'in_transit' THEN 2
    WHEN status = 'delivered' THEN 3
    WHEN status = 'failed' THEN 4
    ELSE 1
  END
WHERE status_int IS NULL;

-- Drop old status column and rename new one
ALTER TABLE bex_shipments DROP COLUMN IF EXISTS status;
ALTER TABLE bex_shipments RENAME COLUMN status_int TO status;

-- Add some initial status_text values based on status
UPDATE bex_shipments SET
  status_text = CASE
    WHEN status = 1 THEN 'Ожидает отправки'
    WHEN status = 2 THEN 'В пути'
    WHEN status = 3 THEN 'Доставлено'
    WHEN status = 4 THEN 'Возврат'
    ELSE 'Неизвестно'
  END
WHERE status_text IS NULL;