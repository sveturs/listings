-- Восстановление дублирующихся индексов
-- (на случай если потребуется откатить изменения)

-- 1. BEX Shipments
CREATE INDEX IF NOT EXISTS idx_bex_shipments_tracking_number ON bex_shipments(tracking_number);

-- 2. Car Makes
CREATE INDEX IF NOT EXISTS idx_car_makes_slug ON car_makes(slug);

-- 3. Car Models
CREATE INDEX IF NOT EXISTS idx_car_models_make_slug ON car_models(make_id, slug);

-- 4. Cities
CREATE INDEX IF NOT EXISTS idx_cities_slug ON cities(slug);

-- 5. Custom UI Components
CREATE INDEX IF NOT EXISTS idx_custom_ui_components_name ON custom_ui_components(name);

-- 6. Custom UI Templates
CREATE INDEX IF NOT EXISTS idx_custom_ui_templates_name ON custom_ui_templates(name);

-- 7. Deliveries
CREATE INDEX IF NOT EXISTS idx_deliveries_tracking_token ON deliveries(tracking_token);

-- 8. Payment Transactions
CREATE INDEX IF NOT EXISTS idx_payment_transactions_order_reference ON payment_transactions(order_reference);

-- 9. Post Express Shipments
CREATE INDEX IF NOT EXISTS idx_post_express_shipments_tracking_number ON post_express_shipments(tracking_number);

-- 10. Unified Attributes
CREATE INDEX IF NOT EXISTS idx_unified_attributes_code ON unified_attributes(code);

-- 11. Viber Users
CREATE INDEX IF NOT EXISTS idx_viber_users_viber_id ON viber_users(viber_id);

-- 12. VIN Decode Cache
CREATE INDEX IF NOT EXISTS idx_vin_decode_cache_vin ON vin_decode_cache(vin);
