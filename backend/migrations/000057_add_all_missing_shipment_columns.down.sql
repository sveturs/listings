-- Remove added columns and tracking events table
DROP TABLE IF EXISTS post_express_tracking_events CASCADE;

ALTER TABLE post_express_shipments
DROP COLUMN IF EXISTS pod_url,
DROP COLUMN IF EXISTS registered_at,
DROP COLUMN IF EXISTS picked_up_at,
DROP COLUMN IF EXISTS delivered_at,
DROP COLUMN IF EXISTS failed_at,
DROP COLUMN IF EXISTS returned_at,
DROP COLUMN IF EXISTS internal_notes,
DROP COLUMN IF EXISTS failed_reason;