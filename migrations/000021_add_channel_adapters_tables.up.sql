-- Channel Adapters tables (moved from channel-adapters service)
-- These tables store external sales channel integrations (eBay, Shopify, Amazon, Etsy)

-- Channel Integrations table
CREATE TABLE IF NOT EXISTS channel_integrations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id INTEGER NOT NULL,
    channel_type VARCHAR(50) NOT NULL,
    channel_name VARCHAR(255) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    credentials JSONB NOT NULL DEFAULT '{}',
    settings JSONB NOT NULL DEFAULT '{}',
    last_sync_at TIMESTAMP WITH TIME ZONE,
    last_error TEXT,
    webhook_url VARCHAR(500),
    webhook_secret VARCHAR(255),
    auto_sync_enabled BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    storefront_id INTEGER,

    CONSTRAINT channel_integrations_channel_type_check
        CHECK (channel_type IN ('ebay', 'shopify', 'amazon', 'etsy')),
    CONSTRAINT channel_integrations_status_check
        CHECK (status IN ('active', 'inactive', 'error', 'pending', 'suspended'))
);

-- Indexes for channel_integrations
CREATE INDEX IF NOT EXISTS idx_channel_integrations_user_id ON channel_integrations(user_id);
CREATE INDEX IF NOT EXISTS idx_channel_integrations_channel_type ON channel_integrations(channel_type);
CREATE INDEX IF NOT EXISTS idx_channel_integrations_status ON channel_integrations(status);
CREATE INDEX IF NOT EXISTS idx_channel_integrations_auto_sync ON channel_integrations(auto_sync_enabled)
    WHERE auto_sync_enabled = true;
CREATE INDEX IF NOT EXISTS idx_channel_integrations_storefront_id ON channel_integrations(storefront_id);

-- Channel Orders table
CREATE TABLE IF NOT EXISTS channel_orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    integration_id UUID NOT NULL REFERENCES channel_integrations(id) ON DELETE CASCADE,
    channel_order_id VARCHAR(255) NOT NULL,
    channel_type VARCHAR(50) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    customer_name VARCHAR(255),
    customer_email VARCHAR(255),
    shipping_address JSONB NOT NULL DEFAULT '{}',
    items JSONB NOT NULL DEFAULT '[]',
    total_amount DECIMAL(12, 2) NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'RSD',
    raw_data JSONB NOT NULL DEFAULT '{}',
    vondi_order_id UUID,
    synced_to_vondi BOOLEAN NOT NULL DEFAULT false,
    last_sync_attempt_at TIMESTAMP WITH TIME ZONE,
    sync_error TEXT,
    order_date TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),

    CONSTRAINT channel_orders_channel_type_check
        CHECK (channel_type IN ('ebay', 'shopify', 'amazon', 'etsy')),
    CONSTRAINT channel_orders_status_check
        CHECK (status IN ('pending', 'processing', 'shipped', 'delivered', 'cancelled', 'refunded')),
    CONSTRAINT channel_orders_unique_channel_order
        UNIQUE (integration_id, channel_order_id)
);

-- Indexes for channel_orders
CREATE INDEX IF NOT EXISTS idx_channel_orders_integration_id ON channel_orders(integration_id);
CREATE INDEX IF NOT EXISTS idx_channel_orders_channel_order_id ON channel_orders(channel_order_id);
CREATE INDEX IF NOT EXISTS idx_channel_orders_channel_type ON channel_orders(channel_type);
CREATE INDEX IF NOT EXISTS idx_channel_orders_status ON channel_orders(status);
CREATE INDEX IF NOT EXISTS idx_channel_orders_synced ON channel_orders(synced_to_vondi)
    WHERE synced_to_vondi = false;
CREATE INDEX IF NOT EXISTS idx_channel_orders_vondi_order_id ON channel_orders(vondi_order_id)
    WHERE vondi_order_id IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_channel_orders_order_date ON channel_orders(order_date DESC);

-- Triggers for updated_at (function already exists in listings DB)
DROP TRIGGER IF EXISTS update_channel_integrations_updated_at ON channel_integrations;
CREATE TRIGGER update_channel_integrations_updated_at
    BEFORE UPDATE ON channel_integrations
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS update_channel_orders_updated_at ON channel_orders;
CREATE TRIGGER update_channel_orders_updated_at
    BEFORE UPDATE ON channel_orders
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Comments
COMMENT ON TABLE channel_integrations IS 'External sales channel integrations (eBay, Shopify, Amazon, Etsy)';
COMMENT ON TABLE channel_orders IS 'Orders imported from external sales channels';
COMMENT ON COLUMN channel_integrations.user_id IS 'User ID (INTEGER) from monolith users table';
COMMENT ON COLUMN channel_integrations.storefront_id IS 'B2C Storefront ID (INTEGER) - integrations are scoped to storefronts, not users';
COMMENT ON COLUMN channel_integrations.credentials IS 'Encrypted OAuth tokens and API keys';
COMMENT ON COLUMN channel_integrations.settings IS 'Channel-specific configuration';
COMMENT ON COLUMN channel_orders.raw_data IS 'Original order data from external channel';
COMMENT ON COLUMN channel_orders.items IS 'Array of order line items';
