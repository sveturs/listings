-- VIN декодер и история проверок
-- Создано: 2025-09-28
-- Описание: Таблицы для хранения результатов VIN декодирования и истории проверок

-- Таблица для хранения декодированной информации VIN
CREATE TABLE IF NOT EXISTS vin_decode_cache (
    id BIGSERIAL PRIMARY KEY,
    vin VARCHAR(17) UNIQUE NOT NULL,
    -- Основная информация
    make VARCHAR(100),
    model VARCHAR(100),
    year INTEGER,
    engine_type VARCHAR(100),
    engine_displacement VARCHAR(50),
    transmission_type VARCHAR(100),
    drivetrain VARCHAR(50),
    body_type VARCHAR(100),
    fuel_type VARCHAR(50),

    -- Технические характеристики
    doors INTEGER,
    seats INTEGER,
    color_exterior VARCHAR(100),
    color_interior VARCHAR(100),

    -- Информация о производителе
    manufacturer VARCHAR(200),
    country_of_origin VARCHAR(100),
    assembly_plant VARCHAR(200),

    -- Дополнительная информация
    vehicle_class VARCHAR(100),
    vehicle_type VARCHAR(100),
    gross_vehicle_weight VARCHAR(50),

    -- Результат декодирования
    decode_status VARCHAR(50) DEFAULT 'success', -- success, partial, failed
    error_message TEXT,
    raw_response JSONB, -- Полный ответ от API для будущих расширений

    -- Метаданные
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,

    -- Индекс для быстрого поиска
    CONSTRAINT vin_length CHECK (LENGTH(vin) = 17)
);

-- Таблица истории проверок VIN пользователями
CREATE TABLE IF NOT EXISTS vin_check_history (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT,
    vin VARCHAR(17) NOT NULL,
    listing_id BIGINT, -- Связь с объявлением, если проверка делается для объявления

    -- Результаты проверки
    decode_success BOOLEAN DEFAULT true,
    decode_cache_id BIGINT REFERENCES vin_decode_cache(id),

    -- Дополнительная информация о проверке
    check_type VARCHAR(50) DEFAULT 'manual', -- manual, auto_fill, verification
    ip_address INET,
    user_agent TEXT,

    -- Время проверки
    checked_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,

    -- Внешние ключи
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL,
    FOREIGN KEY (listing_id) REFERENCES marketplace_listings(id) ON DELETE SET NULL
);

-- Таблица для хранения информации об отзывах и проблемах (будущее расширение)
CREATE TABLE IF NOT EXISTS vin_recalls (
    id BIGSERIAL PRIMARY KEY,
    vin VARCHAR(17) NOT NULL,
    recall_id VARCHAR(100),
    campaign_number VARCHAR(100),

    -- Информация об отзыве
    component VARCHAR(500),
    summary TEXT,
    consequence TEXT,
    remedy TEXT,

    -- Даты
    report_date DATE,
    remedy_date DATE,

    -- Статус
    status VARCHAR(50), -- open, completed, in_progress

    -- Метаданные
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (vin) REFERENCES vin_decode_cache(vin) ON DELETE CASCADE
);

-- Таблица для хранения истории владельцев (будущее расширение)
CREATE TABLE IF NOT EXISTS vin_ownership_history (
    id BIGSERIAL PRIMARY KEY,
    vin VARCHAR(17) NOT NULL,

    -- Информация о владельце
    owner_number INTEGER,
    ownership_type VARCHAR(50), -- personal, fleet, rental, lease

    -- Период владения
    purchase_date DATE,
    sale_date DATE,

    -- Местоположение
    state VARCHAR(100),
    city VARCHAR(100),

    -- Дополнительная информация
    mileage_at_purchase INTEGER,
    mileage_at_sale INTEGER,

    -- Метаданные
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (vin) REFERENCES vin_decode_cache(vin) ON DELETE CASCADE
);

-- Таблица для хранения истории происшествий (будущее расширение)
CREATE TABLE IF NOT EXISTS vin_accident_history (
    id BIGSERIAL PRIMARY KEY,
    vin VARCHAR(17) NOT NULL,

    -- Информация о происшествии
    accident_date DATE,
    accident_type VARCHAR(100), -- collision, theft, flood, fire, vandalism
    severity VARCHAR(50), -- minor, moderate, severe, total_loss

    -- Детали повреждений
    damage_areas TEXT[], -- массив поврежденных зон
    airbag_deployed BOOLEAN,
    structural_damage BOOLEAN,

    -- Ремонт
    repair_cost DECIMAL(10,2),
    repair_date DATE,

    -- Источник информации
    data_source VARCHAR(100),
    report_number VARCHAR(100),

    -- Метаданные
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (vin) REFERENCES vin_decode_cache(vin) ON DELETE CASCADE
);

-- Индексы для оптимизации
CREATE INDEX idx_vin_decode_cache_vin ON vin_decode_cache(vin);
CREATE INDEX idx_vin_decode_cache_make_model ON vin_decode_cache(make, model);
CREATE INDEX idx_vin_decode_cache_year ON vin_decode_cache(year);
CREATE INDEX idx_vin_check_history_user_id ON vin_check_history(user_id);
CREATE INDEX idx_vin_check_history_vin ON vin_check_history(vin);
CREATE INDEX idx_vin_check_history_listing_id ON vin_check_history(listing_id);
CREATE INDEX idx_vin_check_history_checked_at ON vin_check_history(checked_at DESC);
CREATE INDEX idx_vin_recalls_vin ON vin_recalls(vin);
CREATE INDEX idx_vin_ownership_history_vin ON vin_ownership_history(vin);
CREATE INDEX idx_vin_accident_history_vin ON vin_accident_history(vin);

-- Триггер для обновления updated_at
CREATE OR REPLACE FUNCTION update_vin_decode_cache_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_vin_decode_cache_timestamp
    BEFORE UPDATE ON vin_decode_cache
    FOR EACH ROW
    EXECUTE FUNCTION update_vin_decode_cache_updated_at();

-- Комментарии к таблицам
COMMENT ON TABLE vin_decode_cache IS 'Кэш декодированной информации VIN номеров';
COMMENT ON TABLE vin_check_history IS 'История проверок VIN номеров пользователями';
COMMENT ON TABLE vin_recalls IS 'Информация об отзывах производителя для VIN';
COMMENT ON TABLE vin_ownership_history IS 'История владельцев автомобиля';
COMMENT ON TABLE vin_accident_history IS 'История происшествий и ремонтов';