-- Миграция для полной интеграции данных из CarAPI
-- Добавляем недостающие поля и новые таблицы

-- 1. Расширяем таблицу car_makes
ALTER TABLE car_makes
ADD COLUMN IF NOT EXISTS external_id VARCHAR(100),
ADD COLUMN IF NOT EXISTS manufacturer_id VARCHAR(100),
ADD COLUMN IF NOT EXISTS last_sync_at TIMESTAMP,
ADD COLUMN IF NOT EXISTS metadata JSONB DEFAULT '{}';

-- Индекс для быстрого поиска по external_id
CREATE INDEX IF NOT EXISTS idx_car_makes_external_id ON car_makes(external_id);

-- 2. Расширяем таблицу car_models
ALTER TABLE car_models
ADD COLUMN IF NOT EXISTS external_id VARCHAR(100),
ADD COLUMN IF NOT EXISTS body_type VARCHAR(50),
ADD COLUMN IF NOT EXISTS segment VARCHAR(20),
ADD COLUMN IF NOT EXISTS years_range INT4RANGE,
ADD COLUMN IF NOT EXISTS metadata JSONB DEFAULT '{}',
ADD COLUMN IF NOT EXISTS last_sync_at TIMESTAMP;

-- Индекс для external_id
CREATE INDEX IF NOT EXISTS idx_car_models_external_id ON car_models(external_id);

-- 3. Расширяем таблицу car_generations
ALTER TABLE car_generations
ADD COLUMN IF NOT EXISTS external_id VARCHAR(100),
ADD COLUMN IF NOT EXISTS platform VARCHAR(100),
ADD COLUMN IF NOT EXISTS production_country VARCHAR(100),
ADD COLUMN IF NOT EXISTS metadata JSONB DEFAULT '{}',
ADD COLUMN IF NOT EXISTS last_sync_at TIMESTAMP;

-- 4. Создаем таблицу car_trims (комплектации)
CREATE TABLE IF NOT EXISTS car_trims (
    id SERIAL PRIMARY KEY,
    model_id INTEGER NOT NULL REFERENCES car_models(id) ON DELETE CASCADE,
    generation_id INTEGER REFERENCES car_generations(id) ON DELETE SET NULL,
    year INTEGER NOT NULL,
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL,
    
    -- Технические характеристики двигателя
    engine_type VARCHAR(50), -- gasoline, diesel, hybrid, electric
    engine_displacement DECIMAL(3,1), -- объем в литрах
    engine_cylinders INTEGER,
    engine_configuration VARCHAR(50), -- inline, V, W, boxer
    power_hp INTEGER,
    power_kw INTEGER,
    torque_nm INTEGER,
    torque_rpm INTEGER,
    
    -- Трансмиссия
    transmission_type VARCHAR(50), -- manual, automatic, cvt, dct
    transmission_gears INTEGER,
    drive_type VARCHAR(20), -- fwd, rwd, awd, 4wd
    
    -- Топливо и экономичность
    fuel_type VARCHAR(30), -- gasoline, diesel, hybrid, electric, lpg, cng
    fuel_economy_city DECIMAL(4,2), -- л/100км
    fuel_economy_highway DECIMAL(4,2),
    fuel_economy_combined DECIMAL(4,2),
    fuel_tank_capacity INTEGER, -- литры
    battery_capacity DECIMAL(5,1), -- кВт*ч для электромобилей
    electric_range_km INTEGER, -- запас хода для электромобилей
    
    -- Размеры и вес
    length_mm INTEGER,
    width_mm INTEGER,
    height_mm INTEGER,
    wheelbase_mm INTEGER,
    ground_clearance_mm INTEGER,
    cargo_volume_l INTEGER,
    cargo_volume_max_l INTEGER, -- со сложенными сиденьями
    curb_weight_kg INTEGER,
    gross_weight_kg INTEGER,
    
    -- Производительность
    acceleration_0_100 DECIMAL(3,1), -- секунды
    top_speed_kmh INTEGER,
    
    -- Дополнительные характеристики
    seating_capacity INTEGER,
    doors INTEGER,
    body_type VARCHAR(50), -- sedan, hatchback, suv, etc
    
    -- Метаданные и синхронизация
    external_id VARCHAR(100),
    carapi_trim_id INTEGER, -- ID в системе CarAPI
    is_active BOOLEAN DEFAULT true,
    metadata JSONB DEFAULT '{}', -- дополнительные данные из API
    features JSONB DEFAULT '[]', -- список опций и функций
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_sync_at TIMESTAMP
);

-- Индексы для производительности
CREATE INDEX idx_car_trims_model_year ON car_trims(model_id, year);
CREATE INDEX idx_car_trims_generation_year ON car_trims(generation_id, year);
CREATE INDEX idx_car_trims_slug ON car_trims(slug);
CREATE INDEX idx_car_trims_external_id ON car_trims(external_id);
CREATE INDEX idx_car_trims_carapi_id ON car_trims(carapi_trim_id);
CREATE INDEX idx_car_trims_fuel_type ON car_trims(fuel_type);
CREATE INDEX idx_car_trims_body_type ON car_trims(body_type);

-- Уникальный индекс для предотвращения дубликатов
CREATE UNIQUE INDEX idx_car_trims_unique ON car_trims(model_id, year, slug);

-- 5. Таблица для VIN декодирования
CREATE TABLE IF NOT EXISTS vin_decode_cache (
    id SERIAL PRIMARY KEY,
    vin VARCHAR(17) UNIQUE NOT NULL,
    make_id INTEGER REFERENCES car_makes(id),
    model_id INTEGER REFERENCES car_models(id),
    generation_id INTEGER REFERENCES car_generations(id),
    trim_id INTEGER REFERENCES car_trims(id),
    year INTEGER,
    
    -- Декодированные данные
    make_name VARCHAR(100),
    model_name VARCHAR(100),
    trim_name VARCHAR(255),
    body_type VARCHAR(50),
    engine_description VARCHAR(255),
    transmission_description VARCHAR(255),
    drive_type VARCHAR(20),
    
    -- Полные данные из API
    decoded_data JSONB NOT NULL,
    source VARCHAR(50) NOT NULL, -- carapi, nhtsa, etc
    
    -- Метаданные
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP NOT NULL, -- кеш на 30 дней
    last_accessed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Индексы для VIN кеша
CREATE INDEX idx_vin_decode_cache_vin ON vin_decode_cache(vin);
CREATE INDEX idx_vin_decode_cache_expires ON vin_decode_cache(expires_at);
CREATE INDEX idx_vin_decode_cache_make_model ON vin_decode_cache(make_id, model_id);

-- 6. Таблица для логирования синхронизации
CREATE TABLE IF NOT EXISTS car_sync_log (
    id SERIAL PRIMARY KEY,
    sync_type VARCHAR(50) NOT NULL, -- makes, models, trims, vin
    entity_type VARCHAR(50) NOT NULL, -- make, model, generation, trim
    entity_id INTEGER,
    external_id VARCHAR(100),
    action VARCHAR(20) NOT NULL, -- create, update, skip, error
    details JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_car_sync_log_type_date ON car_sync_log(sync_type, created_at DESC);
CREATE INDEX idx_car_sync_log_entity ON car_sync_log(entity_type, entity_id);

-- 7. Таблица для отслеживания лимитов API
CREATE TABLE IF NOT EXISTS carapi_usage (
    id SERIAL PRIMARY KEY,
    date DATE NOT NULL UNIQUE,
    requests_count INTEGER DEFAULT 0,
    last_request_at TIMESTAMP,
    daily_limit INTEGER DEFAULT 1500, -- Base план
    metadata JSONB DEFAULT '{}'
);

-- Триггер для обновления updated_at
CREATE OR REPLACE FUNCTION update_car_trims_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_car_trims_updated_at
    BEFORE UPDATE ON car_trims
    FOR EACH ROW
    EXECUTE FUNCTION update_car_trims_updated_at();

-- Комментарии к таблицам
COMMENT ON TABLE car_trims IS 'Комплектации автомобилей с техническими характеристиками из CarAPI';
COMMENT ON TABLE vin_decode_cache IS 'Кеш результатов VIN декодирования';
COMMENT ON TABLE car_sync_log IS 'Лог синхронизации данных с внешними API';
COMMENT ON TABLE carapi_usage IS 'Отслеживание использования лимитов CarAPI';