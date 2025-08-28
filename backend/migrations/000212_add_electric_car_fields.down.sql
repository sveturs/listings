-- Удаляем индексы
DROP INDEX IF EXISTS idx_car_models_is_electric;
DROP INDEX IF EXISTS idx_car_models_battery_capacity;

-- Удаляем колонки для электромобилей
ALTER TABLE car_models 
DROP COLUMN IF EXISTS is_electric,
DROP COLUMN IF EXISTS battery_capacity_kwh,
DROP COLUMN IF EXISTS battery_capacity_net_kwh,
DROP COLUMN IF EXISTS electric_range_km,
DROP COLUMN IF EXISTS electric_range_wltp_km,
DROP COLUMN IF EXISTS electric_range_standard,
DROP COLUMN IF EXISTS charging_time_0_100,
DROP COLUMN IF EXISTS charging_time_10_80,
DROP COLUMN IF EXISTS fast_charging_power_kw,
DROP COLUMN IF EXISTS onboard_charger_kw;