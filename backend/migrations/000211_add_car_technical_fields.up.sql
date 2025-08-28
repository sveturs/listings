-- Добавляем технические поля для автомобилей
ALTER TABLE car_models 
ADD COLUMN IF NOT EXISTS engine_type VARCHAR(50),
ADD COLUMN IF NOT EXISTS engine_power_kw INTEGER,
ADD COLUMN IF NOT EXISTS engine_power_hp INTEGER,
ADD COLUMN IF NOT EXISTS engine_torque_nm INTEGER,
ADD COLUMN IF NOT EXISTS fuel_type VARCHAR(50),
ADD COLUMN IF NOT EXISTS fuel_consumption_city DECIMAL(4,2),
ADD COLUMN IF NOT EXISTS fuel_consumption_highway DECIMAL(4,2),
ADD COLUMN IF NOT EXISTS fuel_consumption_combined DECIMAL(4,2),
ADD COLUMN IF NOT EXISTS co2_emissions INTEGER,
ADD COLUMN IF NOT EXISTS euro_standard VARCHAR(10),
ADD COLUMN IF NOT EXISTS transmission_type VARCHAR(50),
ADD COLUMN IF NOT EXISTS transmission_gears INTEGER,
ADD COLUMN IF NOT EXISTS drive_type VARCHAR(20),
ADD COLUMN IF NOT EXISTS length_mm INTEGER,
ADD COLUMN IF NOT EXISTS width_mm INTEGER,
ADD COLUMN IF NOT EXISTS height_mm INTEGER,
ADD COLUMN IF NOT EXISTS wheelbase_mm INTEGER,
ADD COLUMN IF NOT EXISTS trunk_volume_l INTEGER,
ADD COLUMN IF NOT EXISTS fuel_tank_l INTEGER,
ADD COLUMN IF NOT EXISTS weight_kg INTEGER,
ADD COLUMN IF NOT EXISTS max_speed_kmh INTEGER,
ADD COLUMN IF NOT EXISTS acceleration_0_100 DECIMAL(3,1),
ADD COLUMN IF NOT EXISTS seats INTEGER,
ADD COLUMN IF NOT EXISTS doors INTEGER;

-- Добавляем поля для отслеживания популярности в Сербии
ALTER TABLE car_models
ADD COLUMN IF NOT EXISTS serbia_popularity_score INTEGER DEFAULT 0,
ADD COLUMN IF NOT EXISTS serbia_average_price_eur DECIMAL(10,2),
ADD COLUMN IF NOT EXISTS serbia_listings_count INTEGER DEFAULT 0;

-- Индексы для быстрого поиска
CREATE INDEX IF NOT EXISTS idx_car_models_engine_type ON car_models(engine_type);
CREATE INDEX IF NOT EXISTS idx_car_models_fuel_type ON car_models(fuel_type);
CREATE INDEX IF NOT EXISTS idx_car_models_transmission_type ON car_models(transmission_type);
CREATE INDEX IF NOT EXISTS idx_car_models_drive_type ON car_models(drive_type);
CREATE INDEX IF NOT EXISTS idx_car_models_body_type ON car_models(body_type);
CREATE INDEX IF NOT EXISTS idx_car_models_serbia_popularity ON car_models(serbia_popularity_score DESC);

-- Комментарии к полям
COMMENT ON COLUMN car_models.engine_type IS 'Тип двигателя (бензин, дизель, гибрид, электро)';
COMMENT ON COLUMN car_models.engine_power_kw IS 'Мощность двигателя в кВт';
COMMENT ON COLUMN car_models.engine_power_hp IS 'Мощность двигателя в л.с.';
COMMENT ON COLUMN car_models.engine_torque_nm IS 'Крутящий момент в Нм';
COMMENT ON COLUMN car_models.fuel_type IS 'Тип топлива (бензин, дизель, гибрид, электро, газ)';
COMMENT ON COLUMN car_models.fuel_consumption_city IS 'Расход топлива в городе (л/100км)';
COMMENT ON COLUMN car_models.fuel_consumption_highway IS 'Расход топлива на трассе (л/100км)';
COMMENT ON COLUMN car_models.fuel_consumption_combined IS 'Смешанный расход топлива (л/100км)';
COMMENT ON COLUMN car_models.co2_emissions IS 'Выбросы CO2 (г/км)';
COMMENT ON COLUMN car_models.euro_standard IS 'Экологический стандарт (Euro 4, Euro 5, Euro 6)';
COMMENT ON COLUMN car_models.transmission_type IS 'Тип трансмиссии (механика, автомат, робот, вариатор)';
COMMENT ON COLUMN car_models.transmission_gears IS 'Количество передач';
COMMENT ON COLUMN car_models.drive_type IS 'Тип привода (передний, задний, полный)';
COMMENT ON COLUMN car_models.serbia_popularity_score IS 'Индекс популярности модели в Сербии';
COMMENT ON COLUMN car_models.serbia_average_price_eur IS 'Средняя цена модели в Сербии (EUR)';
COMMENT ON COLUMN car_models.serbia_listings_count IS 'Количество объявлений этой модели в Сербии';