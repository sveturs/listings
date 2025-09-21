-- BEX Express shipments table
CREATE TABLE IF NOT EXISTS bex_shipments (
    id SERIAL PRIMARY KEY,
    tracking_number VARCHAR(100) UNIQUE NOT NULL,
    order_id INTEGER,
    provider_order_id VARCHAR(100),

    -- Sender information
    sender_name VARCHAR(255) NOT NULL,
    sender_phone VARCHAR(50),
    sender_email VARCHAR(255),
    sender_address TEXT,
    sender_city VARCHAR(100),
    sender_postal_code VARCHAR(20),

    -- Recipient information
    recipient_name VARCHAR(255) NOT NULL,
    recipient_phone VARCHAR(50),
    recipient_email VARCHAR(255),
    recipient_address TEXT,
    recipient_city VARCHAR(100),
    recipient_postal_code VARCHAR(20),

    -- Package details
    weight DECIMAL(10, 3), -- in kg
    width DECIMAL(10, 2),  -- in cm
    height DECIMAL(10, 2), -- in cm
    length DECIMAL(10, 2), -- in cm
    package_type VARCHAR(50),
    content_description TEXT,
    content_value DECIMAL(10, 2),

    -- Status tracking
    status VARCHAR(50) DEFAULT 'pending',
    status_description TEXT,
    last_status_update TIMESTAMP WITH TIME ZONE,

    -- Delivery details
    delivery_type VARCHAR(50),
    delivery_date DATE,
    estimated_delivery TIMESTAMP WITH TIME ZONE,
    actual_delivery TIMESTAMP WITH TIME ZONE,

    -- Cost information
    delivery_cost DECIMAL(10, 2),
    cod_amount DECIMAL(10, 2), -- Cash on delivery
    insurance_amount DECIMAL(10, 2),
    total_cost DECIMAL(10, 2),

    -- API response data
    bex_shipment_id VARCHAR(100),
    bex_barcode VARCHAR(100),
    label_url TEXT,
    raw_response JSONB,

    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Index for tracking number lookup
CREATE INDEX IF NOT EXISTS idx_bex_shipments_tracking_number ON bex_shipments(tracking_number);

-- Index for order lookup
CREATE INDEX IF NOT EXISTS idx_bex_shipments_order_id ON bex_shipments(order_id);

-- Index for status filtering
CREATE INDEX IF NOT EXISTS idx_bex_shipments_status ON bex_shipments(status);

-- Index for date filtering
CREATE INDEX IF NOT EXISTS idx_bex_shipments_created_at ON bex_shipments(created_at DESC);

-- BEX tracking events table
CREATE TABLE IF NOT EXISTS bex_tracking_events (
    id SERIAL PRIMARY KEY,
    shipment_id INTEGER REFERENCES bex_shipments(id) ON DELETE CASCADE,
    event_date TIMESTAMP WITH TIME ZONE NOT NULL,
    event_type VARCHAR(100),
    event_code VARCHAR(50),
    description TEXT,
    location VARCHAR(255),
    raw_data JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Index for shipment events lookup
CREATE INDEX IF NOT EXISTS idx_bex_tracking_events_shipment_id ON bex_tracking_events(shipment_id);

-- Index for event date
CREATE INDEX IF NOT EXISTS idx_bex_tracking_events_date ON bex_tracking_events(event_date DESC);

-- BEX configuration table
CREATE TABLE IF NOT EXISTS bex_configuration (
    id SERIAL PRIMARY KEY,
    key VARCHAR(100) UNIQUE NOT NULL,
    value TEXT,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Insert default BEX configuration
INSERT INTO bex_configuration (key, value, description) VALUES
    ('api_url', 'https://api.bexexpress.rs', 'BEX Express API URL'),
    ('api_username', '', 'BEX Express API username'),
    ('api_password', '', 'BEX Express API password (encrypted)'),
    ('default_package_type', 'STANDARD', 'Default package type'),
    ('max_weight', '30', 'Maximum weight in kg'),
    ('max_dimension', '200', 'Maximum dimension in cm'),
    ('webhook_secret', '', 'Webhook verification secret')
ON CONFLICT (key) DO NOTHING;

-- Add trigger to update the updated_at timestamp
CREATE OR REPLACE FUNCTION update_bex_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_bex_shipments_updated_at
    BEFORE UPDATE ON bex_shipments
    FOR EACH ROW
    EXECUTE FUNCTION update_bex_updated_at();

CREATE TRIGGER update_bex_configuration_updated_at
    BEFORE UPDATE ON bex_configuration
    FOR EACH ROW
    EXECUTE FUNCTION update_bex_updated_at();