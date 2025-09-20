-- Drop triggers
DROP TRIGGER IF EXISTS update_bex_shipments_updated_at ON bex_shipments;
DROP TRIGGER IF EXISTS update_bex_configuration_updated_at ON bex_configuration;

-- Drop function
DROP FUNCTION IF EXISTS update_bex_updated_at();

-- Drop tables
DROP TABLE IF EXISTS bex_tracking_events;
DROP TABLE IF EXISTS bex_shipments;
DROP TABLE IF EXISTS bex_configuration;

-- Drop indexes (they will be dropped with tables, but explicit for clarity)
DROP INDEX IF EXISTS idx_bex_shipments_tracking_number;
DROP INDEX IF EXISTS idx_bex_shipments_order_id;
DROP INDEX IF EXISTS idx_bex_shipments_status;
DROP INDEX IF EXISTS idx_bex_shipments_created_at;
DROP INDEX IF EXISTS idx_bex_tracking_events_shipment_id;
DROP INDEX IF EXISTS idx_bex_tracking_events_date;