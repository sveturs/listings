-- Create Post Express settings table
CREATE TABLE IF NOT EXISTS post_express_settings (
    id SERIAL PRIMARY KEY,
    api_username VARCHAR(255) NOT NULL,
    api_password VARCHAR(255) NOT NULL,
    api_endpoint VARCHAR(500) NOT NULL DEFAULT 'https://wsp.postexpress.rs/api',
    
    -- Sender information
    sender_name VARCHAR(255) NOT NULL,
    sender_address VARCHAR(500) NOT NULL,
    sender_city VARCHAR(255) NOT NULL,
    sender_postal_code VARCHAR(20) NOT NULL,
    sender_phone VARCHAR(50) NOT NULL,
    sender_email VARCHAR(255),
    
    -- Settings
    enabled BOOLEAN DEFAULT true,
    test_mode BOOLEAN DEFAULT true,
    auto_print_labels BOOLEAN DEFAULT false,
    auto_track_shipments BOOLEAN DEFAULT false,
    
    -- Notifications
    notify_on_pickup BOOLEAN DEFAULT false,
    notify_on_delivery BOOLEAN DEFAULT true,
    notify_on_failed_delivery BOOLEAN DEFAULT true,
    
    -- Statistics
    total_shipments INTEGER DEFAULT 0,
    successful_deliveries INTEGER DEFAULT 0,
    failed_deliveries INTEGER DEFAULT 0,
    
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Create Post Express locations table
CREATE TABLE IF NOT EXISTS post_express_locations (
    id SERIAL PRIMARY KEY,
    post_express_id INTEGER UNIQUE,
    name VARCHAR(255) NOT NULL,
    name_cyrillic VARCHAR(255),
    postal_code VARCHAR(20),
    municipality VARCHAR(255),
    
    -- Coordinates
    latitude DOUBLE PRECISION,
    longitude DOUBLE PRECISION,
    
    -- Administrative division
    region VARCHAR(255),
    district VARCHAR(255),
    delivery_zone VARCHAR(50),
    
    -- Capabilities
    is_active BOOLEAN DEFAULT true,
    supports_cod BOOLEAN DEFAULT true,
    supports_express BOOLEAN DEFAULT false,
    
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Create Post Express offices table  
CREATE TABLE IF NOT EXISTS post_express_offices (
    id SERIAL PRIMARY KEY,
    office_code VARCHAR(50) UNIQUE NOT NULL,
    location_id INTEGER REFERENCES post_express_locations(id),
    
    name VARCHAR(255) NOT NULL,
    address VARCHAR(500) NOT NULL,
    phone VARCHAR(50),
    email VARCHAR(255),
    
    working_hours JSONB,
    
    -- Coordinates
    latitude DOUBLE PRECISION,
    longitude DOUBLE PRECISION,
    
    -- Capabilities
    accepts_packages BOOLEAN DEFAULT true,
    issues_packages BOOLEAN DEFAULT true,
    has_atm BOOLEAN DEFAULT false,
    has_parking BOOLEAN DEFAULT false,
    wheelchair_accessible BOOLEAN DEFAULT false,
    
    -- Status
    is_active BOOLEAN DEFAULT true,
    temporary_closed BOOLEAN DEFAULT false,
    closed_until TIMESTAMPTZ,
    
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Create Post Express rates table
CREATE TABLE IF NOT EXISTS post_express_rates (
    id SERIAL PRIMARY KEY,
    weight_from DECIMAL(10,3) NOT NULL,
    weight_to DECIMAL(10,3) NOT NULL,
    base_price DECIMAL(10,2) NOT NULL,
    
    -- Insurance
    insurance_included_up_to DECIMAL(10,2) DEFAULT 15000,
    insurance_rate_percent DECIMAL(5,2) DEFAULT 1.0,
    cod_fee DECIMAL(10,2) DEFAULT 45,
    
    -- Size limits
    max_length_cm INTEGER DEFAULT 60,
    max_width_cm INTEGER DEFAULT 60,
    max_height_cm INTEGER DEFAULT 60,
    max_dimensions_sum_cm INTEGER DEFAULT 180,
    
    -- Delivery time
    delivery_days_min INTEGER DEFAULT 1,
    delivery_days_max INTEGER DEFAULT 3,
    
    -- Status
    is_active BOOLEAN DEFAULT true,
    is_special_offer BOOLEAN DEFAULT false,
    
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT unique_weight_range UNIQUE (weight_from, weight_to)
);

-- Create Post Express shipments table
CREATE TABLE IF NOT EXISTS post_express_shipments (
    id SERIAL PRIMARY KEY,
    marketplace_order_id INTEGER,
    storefront_order_id BIGINT,
    
    -- Post Express identifiers
    tracking_number VARCHAR(255) UNIQUE,
    barcode VARCHAR(255),
    post_express_id VARCHAR(255),
    
    -- Sender information
    sender_name VARCHAR(255) NOT NULL,
    sender_address VARCHAR(500) NOT NULL,
    sender_city VARCHAR(255) NOT NULL,
    sender_postal_code VARCHAR(20) NOT NULL,
    sender_phone VARCHAR(50) NOT NULL,
    sender_email VARCHAR(255),
    
    -- Receiver information
    receiver_name VARCHAR(255) NOT NULL,
    receiver_address VARCHAR(500) NOT NULL,
    receiver_city VARCHAR(255) NOT NULL,
    receiver_postal_code VARCHAR(20) NOT NULL,
    receiver_phone VARCHAR(50) NOT NULL,
    receiver_email VARCHAR(255),
    
    -- Package details
    weight_kg DECIMAL(10,3) NOT NULL,
    length_cm INTEGER,
    width_cm INTEGER,
    height_cm INTEGER,
    package_contents TEXT,
    declared_value DECIMAL(10,2),
    
    -- Service options
    service_type VARCHAR(50) DEFAULT 'standard',
    cod_amount DECIMAL(10,2),
    insurance_amount DECIMAL(10,2),
    express_delivery BOOLEAN DEFAULT false,
    office_pickup BOOLEAN DEFAULT false,
    office_code VARCHAR(50),
    
    -- Status and tracking
    status VARCHAR(50) DEFAULT 'created',
    status_description TEXT,
    last_tracking_update TIMESTAMPTZ,
    pickup_date TIMESTAMPTZ,
    delivery_date TIMESTAMPTZ,
    
    -- Label and documents
    label_url TEXT,
    label_printed_at TIMESTAMPTZ,
    receipt_url TEXT,
    
    -- History
    status_history JSONB DEFAULT '[]'::jsonb,
    
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes
CREATE INDEX idx_post_express_locations_postal_code ON post_express_locations(postal_code);
CREATE INDEX idx_post_express_locations_name ON post_express_locations(name);
CREATE INDEX idx_post_express_offices_location_id ON post_express_offices(location_id);
CREATE INDEX idx_post_express_offices_is_active ON post_express_offices(is_active);
CREATE INDEX idx_post_express_shipments_tracking_number ON post_express_shipments(tracking_number);
CREATE INDEX idx_post_express_shipments_status ON post_express_shipments(status);
CREATE INDEX idx_post_express_shipments_marketplace_order_id ON post_express_shipments(marketplace_order_id);
CREATE INDEX idx_post_express_shipments_storefront_order_id ON post_express_shipments(storefront_order_id);

-- Insert default settings for testing
INSERT INTO post_express_settings (
    api_username,
    api_password,
    api_endpoint,
    sender_name,
    sender_address,
    sender_city,
    sender_postal_code,
    sender_phone,
    sender_email,
    enabled,
    test_mode
) VALUES (
    'testuser',
    'testpassword',
    'https://wsp.postexpress.rs/api',
    'Sve Tu Test Store',
    'Bulevar oslobođenja 127',
    'Novi Sad',
    '21000',
    '+381601234567',
    'test@svetu.rs',
    true,
    true
);

-- Insert some test rates based on Post Express pricing
INSERT INTO post_express_rates (weight_from, weight_to, base_price, delivery_days_min, delivery_days_max) VALUES
(0, 2, 340, 1, 2),
(2, 5, 450, 1, 2),
(5, 10, 580, 1, 3),
(10, 20, 790, 2, 3),
(20, 30, 1100, 2, 4);

-- Insert some test locations
INSERT INTO post_express_locations (post_express_id, name, name_cyrillic, postal_code, municipality, supports_cod, supports_express) VALUES
(11000, 'Beograd', 'Београд', '11000', 'Beograd', true, true),
(21000, 'Novi Sad', 'Нови Сад', '21000', 'Novi Sad', true, true),
(18000, 'Niš', 'Ниш', '18000', 'Niš', true, true),
(34000, 'Kragujevac', 'Крагујевац', '34000', 'Kragujevac', true, false),
(24000, 'Subotica', 'Суботица', '24000', 'Subotica', true, false);