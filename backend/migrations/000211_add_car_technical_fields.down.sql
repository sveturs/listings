-- Удаляем индексы
DROP INDEX IF EXISTS idx_car_models_engine_type;
DROP INDEX IF EXISTS idx_car_models_fuel_type;
DROP INDEX IF EXISTS idx_car_models_transmission_type;
DROP INDEX IF EXISTS idx_car_models_drive_type;
DROP INDEX IF EXISTS idx_car_models_body_type;
DROP INDEX IF EXISTS idx_car_models_serbia_popularity;

-- Удаляем добавленные колонки
ALTER TABLE car_models 
DROP COLUMN IF EXISTS engine_type,
DROP COLUMN IF EXISTS engine_power_kw,
DROP COLUMN IF EXISTS engine_power_hp,
DROP COLUMN IF EXISTS engine_torque_nm,
DROP COLUMN IF EXISTS fuel_type,
DROP COLUMN IF EXISTS fuel_consumption_city,
DROP COLUMN IF EXISTS fuel_consumption_highway,
DROP COLUMN IF EXISTS fuel_consumption_combined,
DROP COLUMN IF EXISTS co2_emissions,
DROP COLUMN IF EXISTS euro_standard,
DROP COLUMN IF EXISTS transmission_type,
DROP COLUMN IF EXISTS transmission_gears,
DROP COLUMN IF EXISTS drive_type,
DROP COLUMN IF EXISTS length_mm,
DROP COLUMN IF EXISTS width_mm,
DROP COLUMN IF EXISTS height_mm,
DROP COLUMN IF EXISTS wheelbase_mm,
DROP COLUMN IF EXISTS trunk_volume_l,
DROP COLUMN IF EXISTS fuel_tank_l,
DROP COLUMN IF EXISTS weight_kg,
DROP COLUMN IF EXISTS max_speed_kmh,
DROP COLUMN IF EXISTS acceleration_0_100,
DROP COLUMN IF EXISTS seats,
DROP COLUMN IF EXISTS doors,
DROP COLUMN IF EXISTS serbia_popularity_score,
DROP COLUMN IF EXISTS serbia_average_price_eur,
DROP COLUMN IF EXISTS serbia_listings_count;