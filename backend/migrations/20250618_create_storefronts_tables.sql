-- Migration: Create storefronts tables
-- Description: Initial storefront system tables with support for future killer features

-- Main storefronts table
CREATE TABLE IF NOT EXISTS storefronts (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    slug VARCHAR(100) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    
    -- Branding
    logo_url VARCHAR(500),
    banner_url VARCHAR(500),
    theme JSONB DEFAULT '{"primaryColor": "#1976d2", "layout": "grid"}',
    
    -- Contact information
    phone VARCHAR(50),
    email VARCHAR(255),
    website VARCHAR(255),
    
    -- Location
    address TEXT,
    city VARCHAR(100),
    postal_code VARCHAR(20),
    country VARCHAR(2) DEFAULT 'RS',
    latitude DECIMAL(10, 8),
    longitude DECIMAL(11, 8),
    
    -- Business settings
    settings JSONB DEFAULT '{}',
    seo_meta JSONB DEFAULT '{}',
    
    -- Status and stats
    is_active BOOLEAN DEFAULT false,
    is_verified BOOLEAN DEFAULT false,
    verification_date TIMESTAMP,
    rating DECIMAL(3, 2) DEFAULT 0.00,
    reviews_count INT DEFAULT 0,
    products_count INT DEFAULT 0,
    sales_count INT DEFAULT 0,
    views_count INT DEFAULT 0,
    
    -- Subscription (for monetization)
    subscription_plan VARCHAR(50) DEFAULT 'starter', -- starter, professional, business, enterprise
    subscription_expires_at TIMESTAMP,
    commission_rate DECIMAL(5, 2) DEFAULT 3.00, -- percentage
    
    -- AI and killer features preparation
    ai_agent_enabled BOOLEAN DEFAULT false,
    ai_agent_config JSONB DEFAULT '{}',
    live_shopping_enabled BOOLEAN DEFAULT false,
    group_buying_enabled BOOLEAN DEFAULT false,
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    -- Indexes
    INDEX idx_storefronts_user_id (user_id),
    INDEX idx_storefronts_slug (slug),
    INDEX idx_storefronts_city (city),
    INDEX idx_storefronts_is_active (is_active),
    INDEX idx_storefronts_rating (rating DESC)
);

-- Storefront staff/employees
CREATE TABLE IF NOT EXISTS storefront_staff (
    id SERIAL PRIMARY KEY,
    storefront_id INT NOT NULL REFERENCES storefronts(id) ON DELETE CASCADE,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role VARCHAR(50) NOT NULL DEFAULT 'staff', -- owner, manager, support, moderator
    permissions JSONB DEFAULT '{}', -- granular permissions
    
    -- Activity tracking
    last_active_at TIMESTAMP,
    actions_count INT DEFAULT 0,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    UNIQUE(storefront_id, user_id),
    INDEX idx_staff_storefront_id (storefront_id),
    INDEX idx_staff_user_id (user_id)
);

-- Working hours
CREATE TABLE IF NOT EXISTS storefront_hours (
    id SERIAL PRIMARY KEY,
    storefront_id INT NOT NULL REFERENCES storefronts(id) ON DELETE CASCADE,
    day_of_week INT NOT NULL CHECK (day_of_week >= 0 AND day_of_week <= 6), -- 0=Sunday, 6=Saturday
    open_time TIME,
    close_time TIME,
    is_closed BOOLEAN DEFAULT false,
    
    -- Support for special hours
    special_date DATE, -- for holidays/special days
    special_note VARCHAR(255), -- e.g., "Closed for holiday"
    
    UNIQUE(storefront_id, day_of_week, special_date),
    INDEX idx_hours_storefront_id (storefront_id)
);

-- Payment methods
CREATE TABLE IF NOT EXISTS storefront_payment_methods (
    id SERIAL PRIMARY KEY,
    storefront_id INT NOT NULL REFERENCES storefronts(id) ON DELETE CASCADE,
    method_type VARCHAR(50) NOT NULL, -- cash, card, bank_transfer, paypal, crypto, postanska
    is_enabled BOOLEAN DEFAULT true,
    
    -- Provider specific settings
    provider VARCHAR(50), -- postanska_stedionica, stripe, paypal, etc.
    settings JSONB DEFAULT '{}', -- API keys, merchant IDs, etc.
    
    -- Fees and limits
    transaction_fee DECIMAL(5, 2) DEFAULT 0.00, -- percentage
    min_amount DECIMAL(10, 2),
    max_amount DECIMAL(10, 2),
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    INDEX idx_payment_storefront_id (storefront_id),
    INDEX idx_payment_is_enabled (is_enabled)
);

-- Delivery options
CREATE TABLE IF NOT EXISTS storefront_delivery_options (
    id SERIAL PRIMARY KEY,
    storefront_id INT NOT NULL REFERENCES storefronts(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    
    -- Pricing
    base_price DECIMAL(10, 2) NOT NULL DEFAULT 0.00,
    price_per_km DECIMAL(10, 2) DEFAULT 0.00,
    price_per_kg DECIMAL(10, 2) DEFAULT 0.00,
    free_above_amount DECIMAL(10, 2), -- free delivery above this amount
    
    -- Delivery constraints
    min_order_amount DECIMAL(10, 2),
    max_weight_kg DECIMAL(10, 2),
    max_distance_km DECIMAL(10, 2),
    estimated_days_min INT DEFAULT 1,
    estimated_days_max INT DEFAULT 3,
    
    -- Zones and availability
    zones JSONB DEFAULT '[]', -- array of {name, postal_codes[], price_modifier}
    available_days JSONB DEFAULT '[1,2,3,4,5]', -- array of day numbers
    cutoff_time TIME, -- orders after this time go next day
    
    -- Provider integration
    provider VARCHAR(50), -- posta_srbije, aks, bex, city_express, self
    provider_config JSONB DEFAULT '{}',
    
    is_active BOOLEAN DEFAULT true,
    display_order INT DEFAULT 0,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    INDEX idx_delivery_storefront_id (storefront_id),
    INDEX idx_delivery_is_active (is_active)
);

-- Link listings to storefronts
ALTER TABLE marketplace_listings 
ADD COLUMN IF NOT EXISTS storefront_id INT REFERENCES storefronts(id) ON DELETE SET NULL,
ADD COLUMN IF NOT EXISTS is_storefront_featured BOOLEAN DEFAULT false,
ADD INDEX IF NOT EXISTS idx_listings_storefront_id (storefront_id);

-- Analytics and metrics (for future killer features)
CREATE TABLE IF NOT EXISTS storefront_analytics (
    id SERIAL PRIMARY KEY,
    storefront_id INT NOT NULL REFERENCES storefronts(id) ON DELETE CASCADE,
    date DATE NOT NULL,
    
    -- Basic metrics
    views INT DEFAULT 0,
    unique_visitors INT DEFAULT 0,
    products_viewed INT DEFAULT 0,
    add_to_cart_count INT DEFAULT 0,
    checkout_started INT DEFAULT 0,
    orders_completed INT DEFAULT 0,
    
    -- Financial metrics
    revenue DECIMAL(12, 2) DEFAULT 0.00,
    average_order_value DECIMAL(10, 2) DEFAULT 0.00,
    
    -- Engagement metrics (for killer features)
    live_stream_minutes INT DEFAULT 0,
    live_stream_viewers INT DEFAULT 0,
    group_buying_participants INT DEFAULT 0,
    ai_chat_interactions INT DEFAULT 0,
    social_shares INT DEFAULT 0,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    UNIQUE(storefront_id, date),
    INDEX idx_analytics_storefront_date (storefront_id, date DESC)
);

-- Triggers for updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_storefronts_updated_at BEFORE UPDATE ON storefronts
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_storefront_staff_updated_at BEFORE UPDATE ON storefront_staff
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_storefront_delivery_updated_at BEFORE UPDATE ON storefront_delivery_options
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Function to generate unique slug
CREATE OR REPLACE FUNCTION generate_unique_slug(base_name VARCHAR, table_name VARCHAR)
RETURNS VARCHAR AS $$
DECLARE
    slug VARCHAR;
    counter INT := 0;
BEGIN
    -- Convert to lowercase, replace spaces and special chars
    slug := LOWER(REGEXP_REPLACE(base_name, '[^a-zA-Z0-9]+', '-', 'g'));
    slug := TRIM(BOTH '-' FROM slug);
    
    -- Check if slug exists
    LOOP
        IF counter = 0 THEN
            -- First try without number
            EXIT WHEN NOT EXISTS (
                SELECT 1 FROM storefronts WHERE slug = slug
            );
        ELSE
            -- Add counter
            EXIT WHEN NOT EXISTS (
                SELECT 1 FROM storefronts WHERE slug = slug || '-' || counter
            );
            slug := slug || '-' || counter;
        END IF;
        counter := counter + 1;
    END LOOP;
    
    RETURN slug;
END;
$$ LANGUAGE plpgsql;

-- Insert some default data for testing
INSERT INTO storefront_payment_methods (storefront_id, method_type, provider, is_enabled)
SELECT id, 'cash', 'cash_on_delivery', true FROM storefronts
ON CONFLICT DO NOTHING;

-- Comments for future features
COMMENT ON COLUMN storefronts.ai_agent_config IS 'Configuration for AI sales agent (Prodavac) - personality, knowledge base, price negotiation rules';
COMMENT ON COLUMN storefronts.live_shopping_enabled IS 'Enable live shopping shows feature';
COMMENT ON COLUMN storefronts.group_buying_enabled IS 'Enable group buying campaigns';
COMMENT ON COLUMN storefront_analytics.live_stream_minutes IS 'Total minutes of live shopping broadcasts';
COMMENT ON COLUMN storefront_analytics.ai_chat_interactions IS 'Number of AI agent conversations';