-- BEX Express Integration Tables

-- Настройки интеграции BEX Express
CREATE TABLE IF NOT EXISTS bex_settings (
    id SERIAL PRIMARY KEY,
    auth_token VARCHAR(255) NOT NULL,
    client_id VARCHAR(100) NOT NULL,
    api_endpoint VARCHAR(255) NOT NULL DEFAULT 'https://api.bex.rs:62502',
    
    -- Данные отправителя (магазина)
    sender_client_id VARCHAR(100) NOT NULL,
    sender_name VARCHAR(255) NOT NULL,
    sender_address VARCHAR(500) NOT NULL,
    sender_city VARCHAR(100) NOT NULL,
    sender_postal_code VARCHAR(20) NOT NULL,
    sender_phone VARCHAR(50) NOT NULL,
    sender_email VARCHAR(255) NOT NULL,
    
    -- Настройки
    enabled BOOLEAN DEFAULT true,
    test_mode BOOLEAN DEFAULT false,
    auto_print_labels BOOLEAN DEFAULT true,
    auto_track_shipments BOOLEAN DEFAULT true,
    use_address_lookup BOOLEAN DEFAULT true,
    
    -- Уведомления
    notify_on_pickup BOOLEAN DEFAULT true,
    notify_on_delivery BOOLEAN DEFAULT true,
    notify_on_failed_delivery BOOLEAN DEFAULT true,
    
    -- Статистика
    total_shipments INTEGER DEFAULT 0,
    successful_deliveries INTEGER DEFAULT 0,
    failed_deliveries INTEGER DEFAULT 0,
    
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Справочник муниципалитетов BEX
CREATE TABLE IF NOT EXISTS bex_municipalities (
    id SERIAL PRIMARY KEY,
    bex_id INTEGER UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    name_cyrillic VARCHAR(255) NOT NULL,
    code VARCHAR(50),
    region VARCHAR(100),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_bex_municipalities_bex_id ON bex_municipalities(bex_id);
CREATE INDEX idx_bex_municipalities_name ON bex_municipalities(name);
CREATE INDEX idx_bex_municipalities_name_cyrillic ON bex_municipalities(name_cyrillic);
CREATE INDEX idx_bex_municipalities_active ON bex_municipalities(is_active) WHERE is_active = true;

-- Справочник населенных пунктов BEX
CREATE TABLE IF NOT EXISTS bex_places (
    id SERIAL PRIMARY KEY,
    bex_id VARCHAR(100) UNIQUE NOT NULL,
    municipality_id INTEGER REFERENCES bex_municipalities(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    name_cyrillic VARCHAR(255) NOT NULL,
    postal_code VARCHAR(20),
    latitude DECIMAL(10, 8),
    longitude DECIMAL(11, 8),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_bex_places_bex_id ON bex_places(bex_id);
CREATE INDEX idx_bex_places_municipality_id ON bex_places(municipality_id);
CREATE INDEX idx_bex_places_name ON bex_places(name);
CREATE INDEX idx_bex_places_name_cyrillic ON bex_places(name_cyrillic);
CREATE INDEX idx_bex_places_postal_code ON bex_places(postal_code);
CREATE INDEX idx_bex_places_active ON bex_places(is_active) WHERE is_active = true;

-- Справочник улиц BEX
CREATE TABLE IF NOT EXISTS bex_streets (
    id SERIAL PRIMARY KEY,
    bex_id VARCHAR(100) UNIQUE NOT NULL,
    place_id INTEGER REFERENCES bex_places(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    name_cyrillic VARCHAR(255) NOT NULL,
    street_type VARCHAR(50),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_bex_streets_bex_id ON bex_streets(bex_id);
CREATE INDEX idx_bex_streets_place_id ON bex_streets(place_id);
CREATE INDEX idx_bex_streets_name ON bex_streets(name);
CREATE INDEX idx_bex_streets_name_cyrillic ON bex_streets(name_cyrillic);
CREATE INDEX idx_bex_streets_active ON bex_streets(is_active) WHERE is_active = true;

-- Full text search для адресов
CREATE INDEX idx_bex_streets_name_trgm ON bex_streets USING gin (name gin_trgm_ops);
CREATE INDEX idx_bex_streets_name_cyrillic_trgm ON bex_streets USING gin (name_cyrillic gin_trgm_ops);
CREATE INDEX idx_bex_places_name_trgm ON bex_places USING gin (name gin_trgm_ops);
CREATE INDEX idx_bex_places_name_cyrillic_trgm ON bex_places USING gin (name_cyrillic gin_trgm_ops);

-- Пункты выдачи BEX (Parcel Shops)
CREATE TABLE IF NOT EXISTS bex_parcel_shops (
    id SERIAL PRIMARY KEY,
    bex_id INTEGER UNIQUE NOT NULL,
    code VARCHAR(50) NOT NULL,
    name VARCHAR(255) NOT NULL,
    address VARCHAR(500) NOT NULL,
    city VARCHAR(100) NOT NULL,
    postal_code VARCHAR(20),
    phone VARCHAR(50),
    working_hours JSONB,
    latitude DECIMAL(10, 8),
    longitude DECIMAL(11, 8),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_bex_parcel_shops_bex_id ON bex_parcel_shops(bex_id);
CREATE INDEX idx_bex_parcel_shops_city ON bex_parcel_shops(city);
CREATE INDEX idx_bex_parcel_shops_active ON bex_parcel_shops(is_active) WHERE is_active = true;

-- Отправления BEX
CREATE TABLE IF NOT EXISTS bex_shipments (
    id SERIAL PRIMARY KEY,
    marketplace_order_id INTEGER REFERENCES marketplace_orders(id) ON DELETE SET NULL,
    storefront_order_id BIGINT REFERENCES storefront_orders(id) ON DELETE SET NULL,
    
    -- BEX идентификаторы
    bex_shipment_id INTEGER,
    tracking_number VARCHAR(100),
    
    -- Отправитель
    sender_name VARCHAR(255) NOT NULL,
    sender_address VARCHAR(500) NOT NULL,
    sender_city VARCHAR(100) NOT NULL,
    sender_postal_code VARCHAR(20) NOT NULL,
    sender_phone VARCHAR(50) NOT NULL,
    sender_email VARCHAR(255),
    
    -- Получатель
    recipient_name VARCHAR(255) NOT NULL,
    recipient_address VARCHAR(500) NOT NULL,
    recipient_city VARCHAR(100) NOT NULL,
    recipient_postal_code VARCHAR(20) NOT NULL,
    recipient_phone VARCHAR(50) NOT NULL,
    recipient_email VARCHAR(255),
    
    -- Параметры посылки
    shipment_type INTEGER NOT NULL DEFAULT 1,
    shipment_category INTEGER NOT NULL,
    shipment_contents INTEGER NOT NULL DEFAULT 3,
    weight_kg DECIMAL(10, 3) NOT NULL,
    total_packages INTEGER NOT NULL DEFAULT 1,
    
    -- Услуги
    pay_type INTEGER NOT NULL DEFAULT 6,
    cod_amount DECIMAL(10, 2),
    insurance_amount DECIMAL(10, 2),
    personal_delivery BOOLEAN DEFAULT false,
    return_signed_invoices BOOLEAN DEFAULT false,
    return_signed_confirmation BOOLEAN DEFAULT false,
    return_package BOOLEAN DEFAULT false,
    
    -- Комментарии
    comment_public TEXT,
    comment_private TEXT,
    delivery_instructions TEXT,
    
    -- Статус
    status INTEGER NOT NULL DEFAULT 3,
    status_text TEXT,
    failed_reason TEXT,
    
    -- Документы
    label_base64 TEXT,
    label_url VARCHAR(500),
    
    -- Временные метки
    registered_at TIMESTAMP,
    picked_up_at TIMESTAMP,
    delivered_at TIMESTAMP,
    failed_at TIMESTAMP,
    returned_at TIMESTAMP,
    
    -- История статусов
    status_history JSONB,
    
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_bex_shipments_marketplace_order ON bex_shipments(marketplace_order_id);
CREATE INDEX idx_bex_shipments_storefront_order ON bex_shipments(storefront_order_id);
CREATE INDEX idx_bex_shipments_bex_id ON bex_shipments(bex_shipment_id);
CREATE INDEX idx_bex_shipments_tracking ON bex_shipments(tracking_number);
CREATE INDEX idx_bex_shipments_status ON bex_shipments(status);
CREATE INDEX idx_bex_shipments_created ON bex_shipments(created_at DESC);
CREATE INDEX idx_bex_shipments_recipient_phone ON bex_shipments(recipient_phone);

-- История событий отслеживания BEX
CREATE TABLE IF NOT EXISTS bex_tracking_events (
    id SERIAL PRIMARY KEY,
    shipment_id INTEGER REFERENCES bex_shipments(id) ON DELETE CASCADE,
    event_code VARCHAR(50) NOT NULL,
    event_description TEXT NOT NULL,
    event_location VARCHAR(255),
    event_timestamp TIMESTAMP NOT NULL,
    additional_info JSONB,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_bex_tracking_shipment ON bex_tracking_events(shipment_id);
CREATE INDEX idx_bex_tracking_timestamp ON bex_tracking_events(event_timestamp DESC);

-- Тарифы BEX
CREATE TABLE IF NOT EXISTS bex_rates (
    id SERIAL PRIMARY KEY,
    weight_from DECIMAL(10, 3) NOT NULL,
    weight_to DECIMAL(10, 3) NOT NULL,
    base_price DECIMAL(10, 2) NOT NULL,
    
    insurance_included_up_to DECIMAL(10, 2) DEFAULT 0,
    insurance_rate_percent DECIMAL(5, 2) DEFAULT 0.5,
    cod_fee DECIMAL(10, 2) DEFAULT 150,
    
    max_length_cm INTEGER,
    max_width_cm INTEGER,
    max_height_cm INTEGER,
    max_dimensions_sum_cm INTEGER,
    
    delivery_days_min INTEGER DEFAULT 1,
    delivery_days_max INTEGER DEFAULT 3,
    
    is_active BOOLEAN DEFAULT true,
    is_special_offer BOOLEAN DEFAULT false,
    
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_bex_rates_weight ON bex_rates(weight_from, weight_to);
CREATE INDEX idx_bex_rates_active ON bex_rates(is_active) WHERE is_active = true;

-- Добавляем поле delivery_provider в таблицы заказов если его еще нет
DO $$ 
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns 
                   WHERE table_name = 'marketplace_orders' 
                   AND column_name = 'delivery_provider') THEN
        ALTER TABLE marketplace_orders 
        ADD COLUMN delivery_provider VARCHAR(50) DEFAULT 'standard',
        ADD COLUMN delivery_tracking_number VARCHAR(100),
        ADD COLUMN delivery_status VARCHAR(50),
        ADD COLUMN delivery_metadata JSONB;
        
        CREATE INDEX idx_marketplace_orders_delivery_provider ON marketplace_orders(delivery_provider);
        CREATE INDEX idx_marketplace_orders_delivery_tracking ON marketplace_orders(delivery_tracking_number);
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns 
                   WHERE table_name = 'storefront_orders' 
                   AND column_name = 'delivery_provider') THEN
        ALTER TABLE storefront_orders 
        ADD COLUMN delivery_provider VARCHAR(50) DEFAULT 'standard',
        ADD COLUMN delivery_tracking_number VARCHAR(100),
        ADD COLUMN delivery_status VARCHAR(50),
        ADD COLUMN delivery_metadata JSONB;
        
        CREATE INDEX idx_storefront_orders_delivery_provider ON storefront_orders(delivery_provider);
        CREATE INDEX idx_storefront_orders_delivery_tracking ON storefront_orders(delivery_tracking_number);
    END IF;
END $$;

-- Вставляем начальные настройки BEX из переменных окружения
-- Эти данные будут обновлены администратором через UI
INSERT INTO bex_settings (
    auth_token,
    client_id,
    api_endpoint,
    sender_client_id,
    sender_name,
    sender_address,
    sender_city,
    sender_postal_code,
    sender_phone,
    sender_email,
    enabled,
    test_mode
) VALUES (
    'd50261-18wo-8539-ee5a-67uu3tu79', -- Будет заменен из ENV
    '326166',
    'https://api.bex.rs:62502',
    '326166',
    'Sve Tu Platform d.o.o.',
    'Vase Stajića 18',
    'Novi Sad',
    '21000',
    '+381 21 123456',
    'info@svetu.rs',
    true,
    false
) ON CONFLICT DO NOTHING;

-- Вставляем базовые тарифы BEX
INSERT INTO bex_rates (weight_from, weight_to, base_price, delivery_days_min, delivery_days_max) VALUES
(0, 0.5, 250, 1, 2),      -- Документы до 0.5кг
(0.5, 1, 350, 1, 3),       -- Пакет до 1кг
(1, 2, 450, 1, 3),         -- Пакет до 2кг
(2, 5, 550, 1, 3),         -- Пакет до 5кг
(5, 10, 750, 2, 4),        -- Пакет до 10кг
(10, 20, 950, 2, 4),       -- Пакет до 20кг
(20, 50, 1500, 2, 5)       -- Пакет до 50кг
ON CONFLICT DO NOTHING;

-- Создаем функцию для автоматического обновления updated_at
CREATE OR REPLACE FUNCTION update_bex_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Создаем триггеры для автоматического обновления updated_at
CREATE TRIGGER update_bex_settings_updated_at BEFORE UPDATE ON bex_settings
    FOR EACH ROW EXECUTE FUNCTION update_bex_updated_at();

CREATE TRIGGER update_bex_municipalities_updated_at BEFORE UPDATE ON bex_municipalities
    FOR EACH ROW EXECUTE FUNCTION update_bex_updated_at();

CREATE TRIGGER update_bex_places_updated_at BEFORE UPDATE ON bex_places
    FOR EACH ROW EXECUTE FUNCTION update_bex_updated_at();

CREATE TRIGGER update_bex_streets_updated_at BEFORE UPDATE ON bex_streets
    FOR EACH ROW EXECUTE FUNCTION update_bex_updated_at();

CREATE TRIGGER update_bex_parcel_shops_updated_at BEFORE UPDATE ON bex_parcel_shops
    FOR EACH ROW EXECUTE FUNCTION update_bex_updated_at();

CREATE TRIGGER update_bex_shipments_updated_at BEFORE UPDATE ON bex_shipments
    FOR EACH ROW EXECUTE FUNCTION update_bex_updated_at();

CREATE TRIGGER update_bex_rates_updated_at BEFORE UPDATE ON bex_rates
    FOR EACH ROW EXECUTE FUNCTION update_bex_updated_at();

-- Комментарии к таблицам
COMMENT ON TABLE bex_settings IS 'Настройки интеграции с BEX Express';
COMMENT ON TABLE bex_municipalities IS 'Справочник муниципалитетов BEX';
COMMENT ON TABLE bex_places IS 'Справочник населенных пунктов BEX';
COMMENT ON TABLE bex_streets IS 'Справочник улиц BEX';
COMMENT ON TABLE bex_parcel_shops IS 'Пункты выдачи BEX';
COMMENT ON TABLE bex_shipments IS 'Отправления через BEX Express';
COMMENT ON TABLE bex_tracking_events IS 'События отслеживания посылок BEX';
COMMENT ON TABLE bex_rates IS 'Тарифы доставки BEX';