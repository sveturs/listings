-- Восстановление удаленных таблиц

-- 1. Восстановление таблицы car_trims
CREATE TABLE IF NOT EXISTS car_trims (
    id SERIAL PRIMARY KEY,
    model_id INTEGER NOT NULL REFERENCES car_models(id) ON DELETE CASCADE,
    generation_id INTEGER REFERENCES car_generations(id) ON DELETE SET NULL,
    year INTEGER NOT NULL,
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL,
    engine_type VARCHAR(50),
    engine_displacement DECIMAL(3,1),
    engine_cylinders INTEGER,
    engine_configuration VARCHAR(50),
    power_hp INTEGER,
    power_kw INTEGER,
    torque_nm INTEGER,
    torque_rpm INTEGER,
    transmission_type VARCHAR(50),
    transmission_gears INTEGER,
    drive_type VARCHAR(20),
    fuel_type VARCHAR(30),
    fuel_economy_city DECIMAL(4,2),
    fuel_economy_highway DECIMAL(4,2),
    fuel_economy_combined DECIMAL(4,2),
    fuel_tank_capacity INTEGER,
    battery_capacity DECIMAL(5,1),
    electric_range_km INTEGER,
    length_mm INTEGER,
    width_mm INTEGER,
    height_mm INTEGER,
    wheelbase_mm INTEGER,
    ground_clearance_mm INTEGER,
    cargo_volume_l INTEGER,
    cargo_volume_max_l INTEGER,
    curb_weight_kg INTEGER,
    gross_weight_kg INTEGER,
    acceleration_0_100 DECIMAL(3,1),
    top_speed_kmh INTEGER,
    seating_capacity INTEGER,
    doors INTEGER,
    body_type VARCHAR(50),
    external_id VARCHAR(100),
    carapi_trim_id INTEGER,
    is_active BOOLEAN DEFAULT true,
    metadata JSONB DEFAULT '{}',
    features JSONB DEFAULT '[]',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_sync_at TIMESTAMP
);

-- Восстановление индексов для car_trims
CREATE INDEX idx_car_trims_model_year ON car_trims(model_id, year);
CREATE INDEX idx_car_trims_generation_year ON car_trims(generation_id, year);
CREATE INDEX idx_car_trims_slug ON car_trims(slug);
CREATE INDEX idx_car_trims_external_id ON car_trims(external_id);
CREATE INDEX idx_car_trims_carapi_id ON car_trims(carapi_trim_id);
CREATE INDEX idx_car_trims_fuel_type ON car_trims(fuel_type);
CREATE INDEX idx_car_trims_body_type ON car_trims(body_type);
CREATE UNIQUE INDEX idx_car_trims_unique ON car_trims(model_id, year, slug);

-- Триггер для car_trims
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

-- 2. Восстановление таблицы car_sync_log
CREATE TABLE IF NOT EXISTS car_sync_log (
    id SERIAL PRIMARY KEY,
    sync_type VARCHAR(50) NOT NULL,
    entity_type VARCHAR(50) NOT NULL,
    entity_id INTEGER,
    external_id VARCHAR(100),
    action VARCHAR(20) NOT NULL,
    details JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_car_sync_log_type_date ON car_sync_log(sync_type, created_at DESC);
CREATE INDEX idx_car_sync_log_entity ON car_sync_log(entity_type, entity_id);

-- 3. Восстановление таблицы carapi_usage
CREATE TABLE IF NOT EXISTS carapi_usage (
    id SERIAL PRIMARY KEY,
    date DATE NOT NULL UNIQUE,
    requests_count INTEGER DEFAULT 0,
    last_request_at TIMESTAMP,
    daily_limit INTEGER DEFAULT 1500,
    metadata JSONB DEFAULT '{}'
);

-- 4. Восстановление таблицы vin_decode_cache
CREATE TABLE IF NOT EXISTS vin_decode_cache (
    id SERIAL PRIMARY KEY,
    vin VARCHAR(17) UNIQUE NOT NULL,
    make_id INTEGER REFERENCES car_makes(id),
    model_id INTEGER REFERENCES car_models(id),
    generation_id INTEGER REFERENCES car_generations(id),
    trim_id INTEGER REFERENCES car_trims(id),
    year INTEGER,
    make_name VARCHAR(100),
    model_name VARCHAR(100),
    trim_name VARCHAR(255),
    body_type VARCHAR(50),
    engine_description VARCHAR(255),
    transmission_description VARCHAR(255),
    drive_type VARCHAR(20),
    decoded_data JSONB NOT NULL,
    source VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP NOT NULL,
    last_accessed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_vin_decode_cache_vin ON vin_decode_cache(vin);
CREATE INDEX idx_vin_decode_cache_expires ON vin_decode_cache(expires_at);
CREATE INDEX idx_vin_decode_cache_make_model ON vin_decode_cache(make_id, model_id);

-- 5. Восстановление таблицы category_variant_attributes
CREATE TABLE IF NOT EXISTS category_variant_attributes (
    id SERIAL PRIMARY KEY,
    category_id INTEGER NOT NULL REFERENCES marketplace_categories(id) ON UPDATE CASCADE ON DELETE CASCADE,
    variant_attribute_name VARCHAR(100) NOT NULL,
    sort_order INTEGER DEFAULT 0,
    is_required BOOLEAN DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_category_variant_attributes_category ON category_variant_attributes(category_id);
CREATE INDEX idx_category_variant_attributes_sort_order ON category_variant_attributes(category_id, sort_order);
CREATE UNIQUE INDEX idx_category_variant_attributes_unique ON category_variant_attributes(category_id, variant_attribute_name);
CREATE INDEX idx_category_variant_attributes_variant_name ON category_variant_attributes(variant_attribute_name);

-- Восстановление комментариев
COMMENT ON TABLE car_trims IS 'Комплектации автомобилей с техническими характеристиками из CarAPI';
COMMENT ON TABLE vin_decode_cache IS 'Кеш результатов VIN декодирования';
COMMENT ON TABLE car_sync_log IS 'Лог синхронизации данных с внешними API';
COMMENT ON TABLE carapi_usage IS 'Отслеживание использования лимитов CarAPI';

-- Очистка external_id обратно (для тех записей, которые были обновлены)
UPDATE car_generations SET external_id = NULL WHERE external_id LIKE 'local_%';
UPDATE car_models SET external_id = NULL WHERE external_id LIKE 'local_%';
UPDATE car_makes SET external_id = NULL WHERE external_id LIKE 'local_%';