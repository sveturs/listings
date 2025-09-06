-- Add all missing columns to post_express_shipments
ALTER TABLE post_express_shipments 
ADD COLUMN IF NOT EXISTS pod_url TEXT,
ADD COLUMN IF NOT EXISTS registered_at TIMESTAMPTZ,
ADD COLUMN IF NOT EXISTS picked_up_at TIMESTAMPTZ,
ADD COLUMN IF NOT EXISTS delivered_at TIMESTAMPTZ,
ADD COLUMN IF NOT EXISTS failed_at TIMESTAMPTZ,
ADD COLUMN IF NOT EXISTS returned_at TIMESTAMPTZ,
ADD COLUMN IF NOT EXISTS internal_notes TEXT,
ADD COLUMN IF NOT EXISTS failed_reason TEXT;

-- Create tracking events table if not exists
CREATE TABLE IF NOT EXISTS post_express_tracking_events (
    id SERIAL PRIMARY KEY,
    shipment_id INTEGER REFERENCES post_express_shipments(id) ON DELETE CASCADE,
    event_code VARCHAR(50),
    event_description TEXT,
    event_location VARCHAR(255),
    event_timestamp TIMESTAMPTZ,
    additional_info JSONB,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Create index for tracking events
CREATE INDEX IF NOT EXISTS idx_post_express_tracking_events_shipment_id ON post_express_tracking_events(shipment_id);
CREATE INDEX IF NOT EXISTS idx_post_express_tracking_events_timestamp ON post_express_tracking_events(event_timestamp);