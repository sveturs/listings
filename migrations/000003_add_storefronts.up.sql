-- Phase 7.5: Add Storefronts (B2C Stores) table
-- Migrating b2c_stores from main svetu database to listings microservice

-- Enable required extensions for geospatial queries
CREATE EXTENSION IF NOT EXISTS cube;
CREATE EXTENSION IF NOT EXISTS earthdistance;

-- =============================================================================
-- storefronts (B2C stores)
-- =============================================================================
CREATE TABLE IF NOT EXISTS storefronts (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    slug VARCHAR(100) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    description TEXT,

    -- Branding
    logo_url VARCHAR(500),
    banner_url VARCHAR(500),
    theme JSONB DEFAULT '{"layout": "grid", "primaryColor": "#1976d2"}'::jsonb,

    -- Contact information
    phone VARCHAR(50),
    email VARCHAR(255),
    website VARCHAR(255),

    -- Address & location
    address TEXT,
    city VARCHAR(100),
    postal_code VARCHAR(20),
    country VARCHAR(2) DEFAULT 'RS',
    latitude NUMERIC(10,8),
    longitude NUMERIC(11,8),
    formatted_address TEXT,
    geo_strategy VARCHAR(50) DEFAULT 'storefront_location',
    default_privacy_level VARCHAR(20) DEFAULT 'exact',
    address_verified BOOLEAN DEFAULT FALSE,

    -- Settings & configuration
    settings JSONB DEFAULT '{}'::jsonb,
    seo_meta JSONB DEFAULT '{}'::jsonb,

    -- Status flags
    is_active BOOLEAN DEFAULT TRUE,
    is_verified BOOLEAN DEFAULT FALSE,
    verification_date TIMESTAMP,

    -- Statistics
    rating NUMERIC(3,2) DEFAULT 0.00,
    reviews_count INTEGER DEFAULT 0,
    products_count INTEGER DEFAULT 0,
    sales_count INTEGER DEFAULT 0,
    views_count INTEGER DEFAULT 0,
    followers_count INTEGER DEFAULT 0,

    -- Subscription management
    subscription_plan VARCHAR(50) DEFAULT 'starter',
    subscription_expires_at TIMESTAMP,
    subscription_id INTEGER,
    is_subscription_active BOOLEAN DEFAULT TRUE,
    commission_rate NUMERIC(5,2) DEFAULT 3.00,

    -- Features
    ai_agent_enabled BOOLEAN DEFAULT FALSE,
    ai_agent_config JSONB DEFAULT '{}'::jsonb,
    live_shopping_enabled BOOLEAN DEFAULT FALSE,
    group_buying_enabled BOOLEAN DEFAULT FALSE,

    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for storefronts
CREATE INDEX IF NOT EXISTS idx_storefronts_user_id ON storefronts(user_id);
CREATE INDEX IF NOT EXISTS idx_storefronts_slug ON storefronts(slug);
CREATE INDEX IF NOT EXISTS idx_storefronts_city ON storefronts(city) WHERE city IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_storefronts_is_active ON storefronts(is_active) WHERE is_active = TRUE;
CREATE INDEX IF NOT EXISTS idx_storefronts_is_verified ON storefronts(is_verified) WHERE is_verified = TRUE;
CREATE INDEX IF NOT EXISTS idx_storefronts_location ON storefronts USING GIST (
    ll_to_earth(latitude::float, longitude::float)
) WHERE latitude IS NOT NULL AND longitude IS NOT NULL;

-- Update trigger for updated_at
CREATE OR REPLACE FUNCTION update_storefronts_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_storefronts_updated_at
BEFORE UPDATE ON storefronts
FOR EACH ROW
EXECUTE FUNCTION update_storefronts_updated_at();

-- Add storefront_id to c2c_listings if not exists
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name = 'c2c_listings' AND column_name = 'storefront_id'
    ) THEN
        ALTER TABLE c2c_listings ADD COLUMN storefront_id INTEGER;
        CREATE INDEX idx_c2c_listings_storefront_id ON c2c_listings(storefront_id) WHERE storefront_id IS NOT NULL;
    END IF;
END $$;

COMMENT ON TABLE storefronts IS 'B2C storefronts (business stores) for marketplace';
COMMENT ON COLUMN storefronts.geo_strategy IS 'Location strategy: storefront_location, product_locations, or both';
COMMENT ON COLUMN storefronts.default_privacy_level IS 'Default privacy level for location: exact, approximate, or hidden';
