-- Migration: Create all storefront related tables
-- Description: Complete storefront system with products, analytics, and import functionality

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
    subscription_plan VARCHAR(50) DEFAULT 'starter',
    subscription_expires_at TIMESTAMP,
    commission_rate DECIMAL(5, 2) DEFAULT 3.00,

    -- AI and killer features preparation
    ai_agent_enabled BOOLEAN DEFAULT false,
    ai_agent_config JSONB DEFAULT '{}',
    live_shopping_enabled BOOLEAN DEFAULT false,
    group_buying_enabled BOOLEAN DEFAULT false,

    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for storefronts
CREATE INDEX idx_storefronts_user_id ON storefronts(user_id);
CREATE INDEX idx_storefronts_slug ON storefronts(slug);
CREATE INDEX idx_storefronts_city ON storefronts(city);
CREATE INDEX idx_storefronts_is_active ON storefronts(is_active);
CREATE INDEX idx_storefronts_rating ON storefronts(rating DESC);

-- Storefront staff/employees
CREATE TABLE IF NOT EXISTS storefront_staff (
    id SERIAL PRIMARY KEY,
    storefront_id INT NOT NULL REFERENCES storefronts(id) ON DELETE CASCADE,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role VARCHAR(50) NOT NULL DEFAULT 'staff',
    permissions JSONB DEFAULT '{}',
    last_active_at TIMESTAMP,
    actions_count INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(storefront_id, user_id)
);

CREATE INDEX idx_staff_storefront_id ON storefront_staff(storefront_id);
CREATE INDEX idx_staff_user_id ON storefront_staff(user_id);

-- Working hours
CREATE TABLE IF NOT EXISTS storefront_hours (
    id SERIAL PRIMARY KEY,
    storefront_id INT NOT NULL REFERENCES storefronts(id) ON DELETE CASCADE,
    day_of_week INT NOT NULL CHECK (day_of_week >= 0 AND day_of_week <= 6),
    open_time TIME,
    close_time TIME,
    is_closed BOOLEAN DEFAULT false,
    special_date DATE,
    special_note VARCHAR(255),
    UNIQUE(storefront_id, day_of_week, special_date)
);

CREATE INDEX idx_hours_storefront_id ON storefront_hours(storefront_id);

-- Payment methods
CREATE TABLE IF NOT EXISTS storefront_payment_methods (
    id SERIAL PRIMARY KEY,
    storefront_id INT NOT NULL REFERENCES storefronts(id) ON DELETE CASCADE,
    method_type VARCHAR(50) NOT NULL,
    is_enabled BOOLEAN DEFAULT true,
    provider VARCHAR(50),
    settings JSONB DEFAULT '{}',
    transaction_fee DECIMAL(5, 2) DEFAULT 0.00,
    min_amount DECIMAL(10, 2),
    max_amount DECIMAL(10, 2),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_payment_storefront_id ON storefront_payment_methods(storefront_id);
CREATE INDEX idx_payment_is_enabled ON storefront_payment_methods(is_enabled);

-- Delivery options
CREATE TABLE IF NOT EXISTS storefront_delivery_options (
    id SERIAL PRIMARY KEY,
    storefront_id INT NOT NULL REFERENCES storefronts(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    base_price DECIMAL(10, 2) NOT NULL DEFAULT 0.00,
    price_per_km DECIMAL(10, 2) DEFAULT 0.00,
    price_per_kg DECIMAL(10, 2) DEFAULT 0.00,
    free_above_amount DECIMAL(10, 2),
    min_order_amount DECIMAL(10, 2),
    max_weight_kg DECIMAL(10, 2),
    max_distance_km DECIMAL(10, 2),
    estimated_days_min INT DEFAULT 1,
    estimated_days_max INT DEFAULT 3,
    zones JSONB DEFAULT '[]',
    available_days JSONB DEFAULT '[1,2,3,4,5]',
    cutoff_time TIME,
    provider VARCHAR(50),
    provider_config JSONB DEFAULT '{}',
    is_active BOOLEAN DEFAULT true,
    display_order INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_delivery_storefront_id ON storefront_delivery_options(storefront_id);
CREATE INDEX idx_delivery_is_active ON storefront_delivery_options(is_active);

-- Analytics and metrics
CREATE TABLE IF NOT EXISTS storefront_analytics (
    id SERIAL PRIMARY KEY,
    storefront_id INT NOT NULL REFERENCES storefronts(id) ON DELETE CASCADE,
    date DATE NOT NULL,
    views INT DEFAULT 0,
    unique_visitors INT DEFAULT 0,
    products_viewed INT DEFAULT 0,
    add_to_cart_count INT DEFAULT 0,
    checkout_started INT DEFAULT 0,
    orders_completed INT DEFAULT 0,
    revenue DECIMAL(12, 2) DEFAULT 0.00,
    average_order_value DECIMAL(10, 2) DEFAULT 0.00,
    live_stream_minutes INT DEFAULT 0,
    live_stream_viewers INT DEFAULT 0,
    group_buying_participants INT DEFAULT 0,
    ai_chat_interactions INT DEFAULT 0,
    social_shares INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(storefront_id, date)
);

CREATE INDEX idx_analytics_storefront_date ON storefront_analytics(storefront_id, date DESC);

-- Create storefront_products table
CREATE TABLE IF NOT EXISTS storefront_products (
    id SERIAL PRIMARY KEY,
    storefront_id INTEGER NOT NULL REFERENCES storefronts(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    price NUMERIC(15, 2) NOT NULL CHECK (price >= 0),
    currency CHAR(3) NOT NULL DEFAULT 'USD',
    category_id INTEGER NOT NULL REFERENCES marketplace_categories(id),
    sku VARCHAR(100),
    barcode VARCHAR(100),
    external_id VARCHAR(100),
    stock_quantity INTEGER NOT NULL DEFAULT 0 CHECK (stock_quantity >= 0),
    stock_status VARCHAR(20) NOT NULL DEFAULT 'in_stock' CHECK (stock_status IN ('in_stock', 'low_stock', 'out_of_stock')),
    is_active BOOLEAN NOT NULL DEFAULT true,
    attributes JSONB DEFAULT '{}',
    view_count INTEGER NOT NULL DEFAULT 0,
    sold_count INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for storefront_products
CREATE INDEX idx_storefront_products_storefront_id ON storefront_products(storefront_id);
CREATE INDEX idx_storefront_products_category_id ON storefront_products(category_id);
CREATE INDEX idx_storefront_products_stock_status ON storefront_products(stock_status);
CREATE INDEX idx_storefront_products_is_active ON storefront_products(is_active);
CREATE INDEX idx_storefront_products_sku ON storefront_products(sku) WHERE sku IS NOT NULL;
CREATE INDEX idx_storefront_products_barcode ON storefront_products(barcode) WHERE barcode IS NOT NULL;
CREATE INDEX idx_storefront_products_external_id ON storefront_products(external_id) WHERE external_id IS NOT NULL;
CREATE INDEX idx_storefront_products_name_gin ON storefront_products USING gin(to_tsvector('simple', name));

-- Create storefront_product_images table
CREATE TABLE IF NOT EXISTS storefront_product_images (
    id SERIAL PRIMARY KEY,
    storefront_product_id INTEGER NOT NULL REFERENCES storefront_products(id) ON DELETE CASCADE,
    image_url TEXT NOT NULL,
    thumbnail_url TEXT NOT NULL,
    display_order INTEGER NOT NULL DEFAULT 0,
    is_default BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_storefront_product_images_product_id ON storefront_product_images(storefront_product_id);
CREATE INDEX idx_storefront_product_images_display_order ON storefront_product_images(display_order);

-- Create storefront_product_variants table
CREATE TABLE IF NOT EXISTS storefront_product_variants (
    id SERIAL PRIMARY KEY,
    storefront_product_id INTEGER NOT NULL REFERENCES storefront_products(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    sku VARCHAR(100),
    price DECIMAL(10, 2) NOT NULL CHECK (price >= 0),
    stock_quantity INTEGER NOT NULL DEFAULT 0 CHECK (stock_quantity >= 0),
    attributes JSONB DEFAULT '{}',
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_storefront_product_variants_product_id ON storefront_product_variants(storefront_product_id);
CREATE INDEX idx_storefront_product_variants_sku ON storefront_product_variants(sku) WHERE sku IS NOT NULL;
CREATE INDEX idx_storefront_product_variants_is_active ON storefront_product_variants(is_active);

-- Create storefront_inventory_movements table
CREATE TABLE IF NOT EXISTS storefront_inventory_movements (
    id SERIAL PRIMARY KEY,
    storefront_product_id INTEGER NOT NULL REFERENCES storefront_products(id) ON DELETE CASCADE,
    variant_id INTEGER REFERENCES storefront_product_variants(id) ON DELETE CASCADE,
    type VARCHAR(20) NOT NULL CHECK (type IN ('in', 'out', 'adjustment')),
    quantity INTEGER NOT NULL,
    reason VARCHAR(50) NOT NULL,
    order_id INTEGER,
    notes TEXT,
    user_id INTEGER NOT NULL REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_storefront_inventory_movements_product_id ON storefront_inventory_movements(storefront_product_id);
CREATE INDEX idx_storefront_inventory_movements_variant_id ON storefront_inventory_movements(variant_id) WHERE variant_id IS NOT NULL;
CREATE INDEX idx_storefront_inventory_movements_type ON storefront_inventory_movements(type);
CREATE INDEX idx_storefront_inventory_movements_created_at ON storefront_inventory_movements(created_at);

-- Add unique constraints
CREATE UNIQUE INDEX unique_storefront_product_sku ON storefront_products (storefront_id, sku) WHERE sku IS NOT NULL;
CREATE UNIQUE INDEX unique_storefront_product_barcode ON storefront_products (storefront_id, barcode) WHERE barcode IS NOT NULL;
CREATE UNIQUE INDEX unique_variant_sku ON storefront_product_variants (storefront_product_id, sku) WHERE sku IS NOT NULL;

-- Extended analytics fields
ALTER TABLE storefront_analytics
ADD COLUMN IF NOT EXISTS page_views INT DEFAULT 0,
ADD COLUMN IF NOT EXISTS bounce_rate DECIMAL(5,2) DEFAULT 0,
ADD COLUMN IF NOT EXISTS avg_session_time INT DEFAULT 0,
ADD COLUMN IF NOT EXISTS orders_count INT DEFAULT 0,
ADD COLUMN IF NOT EXISTS avg_order_value DECIMAL(10,2) DEFAULT 0,
ADD COLUMN IF NOT EXISTS conversion_rate DECIMAL(5,2) DEFAULT 0,
ADD COLUMN IF NOT EXISTS payment_methods_usage JSONB DEFAULT '{}',
ADD COLUMN IF NOT EXISTS product_views INT DEFAULT 0,
ADD COLUMN IF NOT EXISTS checkout_count INT DEFAULT 0,
ADD COLUMN IF NOT EXISTS traffic_sources JSONB DEFAULT '{}',
ADD COLUMN IF NOT EXISTS top_products JSONB DEFAULT '[]',
ADD COLUMN IF NOT EXISTS top_categories JSONB DEFAULT '[]',
ADD COLUMN IF NOT EXISTS orders_by_city JSONB DEFAULT '{}',
ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP;

-- Table for detailed events
CREATE TABLE IF NOT EXISTS storefront_events (
    id SERIAL PRIMARY KEY,
    storefront_id INT NOT NULL REFERENCES storefronts(id) ON DELETE CASCADE,
    event_type VARCHAR(50) NOT NULL,
    event_data JSONB DEFAULT '{}',
    user_id INT,
    session_id VARCHAR(100),
    ip_address INET,
    user_agent TEXT,
    referrer TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_storefront_events_storefront ON storefront_events(storefront_id);
CREATE INDEX idx_storefront_events_type ON storefront_events(event_type);
CREATE INDEX idx_storefront_events_created ON storefront_events(created_at);
CREATE INDEX idx_storefront_events_session ON storefront_events(session_id);

-- Import jobs table
CREATE TABLE import_jobs (
    id SERIAL PRIMARY KEY,
    storefront_id INTEGER NOT NULL REFERENCES storefronts(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    file_name VARCHAR(255) NOT NULL,
    file_type VARCHAR(10) NOT NULL CHECK (file_type IN ('xml', 'csv', 'zip')),
    file_url TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'processing', 'completed', 'failed')),
    total_records INTEGER DEFAULT 0,
    processed_records INTEGER DEFAULT 0,
    successful_records INTEGER DEFAULT 0,
    failed_records INTEGER DEFAULT 0,
    error_message TEXT,
    started_at TIMESTAMP,
    completed_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_import_jobs_storefront_id ON import_jobs(storefront_id);
CREATE INDEX idx_import_jobs_status ON import_jobs(status);
CREATE INDEX idx_import_jobs_created_at ON import_jobs(created_at);

-- Import errors table
CREATE TABLE import_errors (
    id SERIAL PRIMARY KEY,
    job_id INTEGER NOT NULL REFERENCES import_jobs(id) ON DELETE CASCADE,
    line_number INTEGER NOT NULL,
    field_name VARCHAR(100) NOT NULL,
    error_message TEXT NOT NULL,
    raw_data TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_import_errors_job_id ON import_errors(job_id);
CREATE INDEX idx_import_errors_line_number ON import_errors(line_number);

-- Category mappings table
CREATE TABLE category_mappings (
    id SERIAL PRIMARY KEY,
    storefront_id INTEGER NOT NULL REFERENCES storefronts(id) ON DELETE CASCADE,
    import_category1 VARCHAR(255) NOT NULL,
    import_category2 VARCHAR(255),
    import_category3 VARCHAR(255),
    local_category_id INTEGER NOT NULL REFERENCES marketplace_categories(id),
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(storefront_id, import_category1, import_category2, import_category3)
);

CREATE INDEX idx_category_mappings_storefront_id ON category_mappings(storefront_id);
CREATE INDEX idx_category_mappings_import_categories ON category_mappings(import_category1, import_category2, import_category3);

-- Link listings to storefronts
ALTER TABLE marketplace_listings
ADD COLUMN IF NOT EXISTS storefront_id INT REFERENCES storefronts(id) ON DELETE SET NULL,
ADD COLUMN IF NOT EXISTS is_storefront_featured BOOLEAN DEFAULT false;

CREATE INDEX IF NOT EXISTS idx_listings_storefront_id ON marketplace_listings(storefront_id);

-- Comments for documentation
COMMENT ON TABLE storefronts IS 'Main storefront entities';
COMMENT ON TABLE storefront_products IS 'Products sold in storefronts';
COMMENT ON TABLE storefront_product_images IS 'Images for storefront products';
COMMENT ON TABLE storefront_product_variants IS 'Product variants (e.g., sizes, colors)';
COMMENT ON TABLE storefront_inventory_movements IS 'Track inventory changes for audit';
COMMENT ON COLUMN storefronts.ai_agent_config IS 'Configuration for AI sales agent (Prodavac)';
COMMENT ON COLUMN storefronts.live_shopping_enabled IS 'Enable live shopping shows feature';
COMMENT ON COLUMN storefronts.group_buying_enabled IS 'Enable group buying campaigns';
COMMENT ON COLUMN storefront_analytics.live_stream_minutes IS 'Total minutes of live shopping broadcasts';
COMMENT ON COLUMN storefront_analytics.ai_chat_interactions IS 'Number of AI agent conversations';
COMMENT ON COLUMN storefront_products.stock_status IS 'Calculated based on stock_quantity';
COMMENT ON COLUMN storefront_inventory_movements.type IS 'Type of movement: in, out, adjustment';
COMMENT ON COLUMN storefront_inventory_movements.reason IS 'Reason: sale, return, damage, restock, adjustment';

-- Create trigger to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_import_tables_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_import_jobs_updated_at
  BEFORE UPDATE ON import_jobs
  FOR EACH ROW
  EXECUTE FUNCTION update_import_tables_updated_at();

CREATE TRIGGER trigger_category_mappings_updated_at
  BEFORE UPDATE ON category_mappings
  FOR EACH ROW
  EXECUTE FUNCTION update_import_tables_updated_at();
