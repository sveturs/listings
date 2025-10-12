-- Удаление дублирующихся индексов
-- Эти regular индексы дублируют UNIQUE constraints на тех же колонках
-- Оставляем только UNIQUE constraints, так как они выполняют обе функции

-- 1. BEX Shipments
DROP INDEX IF EXISTS idx_bex_shipments_tracking_number;
-- Остаётся: bex_shipments_tracking_number_key (UNIQUE)

-- 2. Car Makes
DROP INDEX IF EXISTS idx_car_makes_slug;
-- Остаётся: car_makes_slug_key (UNIQUE)

-- 3. Car Models (пример из задачи)
DROP INDEX IF EXISTS idx_car_models_make_slug;
-- Остаётся: car_models_make_id_slug_key (UNIQUE)

-- 4. Cities
DROP INDEX IF EXISTS idx_cities_slug;
-- Остаётся: cities_slug_key (UNIQUE)

-- 5. Custom UI Components
DROP INDEX IF EXISTS idx_custom_ui_components_name;
-- Остаётся: custom_ui_components_name_key (UNIQUE)

-- 6. Custom UI Templates
DROP INDEX IF EXISTS idx_custom_ui_templates_name;
-- Остаётся: custom_ui_templates_name_key (UNIQUE)

-- 7. Deliveries
DROP INDEX IF EXISTS idx_deliveries_tracking_token;
-- Остаётся: deliveries_tracking_token_key (UNIQUE)

-- 8. Payment Transactions
DROP INDEX IF EXISTS idx_payment_transactions_order_reference;
-- Остаётся: payment_transactions_order_reference_key (UNIQUE)

-- 9. Post Express Shipments
DROP INDEX IF EXISTS idx_post_express_shipments_tracking_number;
-- Остаётся: post_express_shipments_tracking_number_key (UNIQUE)

-- 10. Unified Attributes
DROP INDEX IF EXISTS idx_unified_attributes_code;
-- Остаётся: unified_attributes_code_key (UNIQUE)

-- 11. Viber Users
DROP INDEX IF EXISTS idx_viber_users_viber_id;
-- Остаётся: viber_users_viber_id_key (UNIQUE)

-- 12. VIN Decode Cache
DROP INDEX IF EXISTS idx_vin_decode_cache_vin;
-- Остаётся: vin_decode_cache_vin_key (UNIQUE)
