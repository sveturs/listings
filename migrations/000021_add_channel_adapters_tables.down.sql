-- Drop channel adapters tables

-- Drop triggers
DROP TRIGGER IF EXISTS update_channel_orders_updated_at ON channel_orders;
DROP TRIGGER IF EXISTS update_channel_integrations_updated_at ON channel_integrations;

-- Drop tables (order matters due to FK)
DROP TABLE IF EXISTS channel_orders;
DROP TABLE IF EXISTS channel_integrations;
