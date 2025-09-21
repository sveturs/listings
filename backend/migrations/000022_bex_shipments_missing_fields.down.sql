-- Revert status column back to varchar
ALTER TABLE bex_shipments ADD COLUMN IF NOT EXISTS status_varchar VARCHAR(50);

UPDATE bex_shipments SET
  status_varchar = CASE
    WHEN status = 1 THEN 'pending'
    WHEN status = 2 THEN 'in_transit'
    WHEN status = 3 THEN 'delivered'
    WHEN status = 4 THEN 'failed'
    ELSE 'pending'
  END
WHERE status_varchar IS NULL;

ALTER TABLE bex_shipments DROP COLUMN IF EXISTS status;
ALTER TABLE bex_shipments RENAME COLUMN status_varchar TO status;

-- Drop added columns
ALTER TABLE bex_shipments
DROP COLUMN IF EXISTS marketplace_order_id,
DROP COLUMN IF EXISTS storefront_order_id,
DROP COLUMN IF EXISTS status_text,
DROP COLUMN IF EXISTS registered_at,
DROP COLUMN IF EXISTS delivered_at,
DROP COLUMN IF EXISTS failed_reason,
DROP COLUMN IF EXISTS status_history;

-- Drop indexes
DROP INDEX IF EXISTS idx_bex_shipments_marketplace_order_id;
DROP INDEX IF EXISTS idx_bex_shipments_storefront_order_id;