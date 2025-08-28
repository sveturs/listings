-- Добавляем поля для электромобилей
ALTER TABLE car_models 
ADD COLUMN IF NOT EXISTS is_electric BOOLEAN DEFAULT FALSE,
ADD COLUMN IF NOT EXISTS battery_capacity_kwh DECIMAL(10,2),
ADD COLUMN IF NOT EXISTS battery_capacity_net_kwh DECIMAL(10,2),
ADD COLUMN IF NOT EXISTS electric_range_km INTEGER,
ADD COLUMN IF NOT EXISTS electric_range_wltp_km INTEGER,
ADD COLUMN IF NOT EXISTS electric_range_standard VARCHAR(20),
ADD COLUMN IF NOT EXISTS charging_time_0_100 DECIMAL(10,2),
ADD COLUMN IF NOT EXISTS charging_time_10_80 DECIMAL(10,2),
ADD COLUMN IF NOT EXISTS fast_charging_power_kw DECIMAL(10,2),
ADD COLUMN IF NOT EXISTS onboard_charger_kw DECIMAL(10,2);

-- Создаем индексы для поиска электромобилей
CREATE INDEX IF NOT EXISTS idx_car_models_is_electric ON car_models(is_electric) WHERE is_electric = true;
CREATE INDEX IF NOT EXISTS idx_car_models_battery_capacity ON car_models(battery_capacity_kwh) WHERE battery_capacity_kwh IS NOT NULL;

-- Обновляем существующие электромобили на основе типа топлива
UPDATE car_models 
SET is_electric = true 
WHERE fuel_type IN ('Electricity', 'Electric', 'EV', 'BEV')
   OR metadata->>'fuel' = 'Electricity';

COMMENT ON COLUMN car_models.is_electric IS 'Флаг для электромобилей';
COMMENT ON COLUMN car_models.battery_capacity_kwh IS 'Общая ёмкость батареи в кВт⋅ч';
COMMENT ON COLUMN car_models.battery_capacity_net_kwh IS 'Полезная ёмкость батареи в кВт⋅ч';
COMMENT ON COLUMN car_models.electric_range_km IS 'Запас хода на электротяге в км';
COMMENT ON COLUMN car_models.electric_range_wltp_km IS 'Запас хода по стандарту WLTP в км';
COMMENT ON COLUMN car_models.electric_range_standard IS 'Стандарт измерения запаса хода (WLTP, EPA, CLTC и т.д.)';
COMMENT ON COLUMN car_models.charging_time_0_100 IS 'Время зарядки от 0 до 100% в часах';
COMMENT ON COLUMN car_models.charging_time_10_80 IS 'Время быстрой зарядки от 10 до 80% в минутах';
COMMENT ON COLUMN car_models.fast_charging_power_kw IS 'Максимальная мощность быстрой зарядки DC в кВт';
COMMENT ON COLUMN car_models.onboard_charger_kw IS 'Мощность встроенного зарядного устройства AC в кВт';