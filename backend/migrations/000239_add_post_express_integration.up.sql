-- Миграция для интеграции с Post Express
-- Добавление поддержки доставки через национального почтового оператора Сербии

-- =====================================================
-- 1. ТАБЛИЦА НАСТРОЕК POST EXPRESS
-- =====================================================
CREATE TABLE post_express_settings (
    id SERIAL PRIMARY KEY,
    
    -- Учетные данные API
    api_username VARCHAR(100) NOT NULL,
    api_password VARCHAR(255) NOT NULL,  -- Зашифрованный пароль
    api_endpoint VARCHAR(255) NOT NULL DEFAULT 'https://wsp.postexpress.rs/api/Transakcija',
    
    -- Настройки отправителя по умолчанию
    sender_name VARCHAR(200) NOT NULL DEFAULT 'Sve Tu Platform',
    sender_address VARCHAR(500) NOT NULL DEFAULT 'Улица Микија Манојловића 53',
    sender_city VARCHAR(100) NOT NULL DEFAULT 'Нови Сад',
    sender_postal_code VARCHAR(20) NOT NULL DEFAULT '21000',
    sender_phone VARCHAR(50) NOT NULL DEFAULT '+381 21 XXX-XXXX',
    sender_email VARCHAR(200) NOT NULL DEFAULT 'shipping@svetu.rs',
    
    -- Флаги и настройки
    enabled BOOLEAN DEFAULT false,
    test_mode BOOLEAN DEFAULT true,  -- Тестовый режим для отладки
    auto_print_labels BOOLEAN DEFAULT false,
    auto_track_shipments BOOLEAN DEFAULT true,
    
    -- Настройки уведомлений
    notify_on_pickup BOOLEAN DEFAULT true,
    notify_on_delivery BOOLEAN DEFAULT true,
    notify_on_failed_delivery BOOLEAN DEFAULT true,
    
    -- Статистика
    total_shipments INTEGER DEFAULT 0,
    successful_deliveries INTEGER DEFAULT 0,
    failed_deliveries INTEGER DEFAULT 0,
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- =====================================================
-- 2. ТАБЛИЦА НАСЕЛЕННЫХ ПУНКТОВ POST EXPRESS
-- =====================================================
CREATE TABLE post_express_locations (
    id SERIAL PRIMARY KEY,
    post_express_id INTEGER NOT NULL UNIQUE,  -- ID из системы Post Express
    
    -- Основные данные
    name VARCHAR(200) NOT NULL,
    name_cyrillic VARCHAR(200),
    postal_code VARCHAR(20),
    municipality VARCHAR(200),
    
    -- Географические данные
    latitude DECIMAL(10, 8),
    longitude DECIMAL(11, 8),
    
    -- Дополнительная информация
    region VARCHAR(100),
    district VARCHAR(100),
    delivery_zone VARCHAR(50),  -- Зона доставки для расчета цены
    
    -- Флаги
    is_active BOOLEAN DEFAULT true,
    supports_cod BOOLEAN DEFAULT true,  -- Поддержка наложенного платежа
    supports_express BOOLEAN DEFAULT true,
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_post_express_locations_name ON post_express_locations(name);
CREATE INDEX idx_post_express_locations_postal_code ON post_express_locations(postal_code);

-- =====================================================
-- 3. ТАБЛИЦА ОТДЕЛЕНИЙ POST EXPRESS
-- =====================================================
CREATE TABLE post_express_offices (
    id SERIAL PRIMARY KEY,
    office_code VARCHAR(50) NOT NULL UNIQUE,
    location_id INTEGER REFERENCES post_express_locations(id),
    
    -- Данные отделения
    name VARCHAR(200) NOT NULL,
    address VARCHAR(500) NOT NULL,
    phone VARCHAR(50),
    email VARCHAR(200),
    
    -- График работы (JSON формат)
    working_hours JSONB,  -- {"monday": {"open": "08:00", "close": "20:00"}, ...}
    
    -- Координаты
    latitude DECIMAL(10, 8),
    longitude DECIMAL(11, 8),
    
    -- Возможности отделения
    accepts_packages BOOLEAN DEFAULT true,
    issues_packages BOOLEAN DEFAULT true,
    has_atm BOOLEAN DEFAULT false,
    has_parking BOOLEAN DEFAULT false,
    wheelchair_accessible BOOLEAN DEFAULT false,
    
    -- Статус
    is_active BOOLEAN DEFAULT true,
    temporary_closed BOOLEAN DEFAULT false,
    closed_until DATE,
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_post_express_offices_location ON post_express_offices(location_id);
CREATE INDEX idx_post_express_offices_active ON post_express_offices(is_active);

-- =====================================================
-- 4. ТАБЛИЦА ТАРИФОВ ДОСТАВКИ
-- =====================================================
CREATE TABLE post_express_rates (
    id SERIAL PRIMARY KEY,
    
    -- Весовая категория
    weight_from DECIMAL(10,3) NOT NULL,  -- кг
    weight_to DECIMAL(10,3) NOT NULL,    -- кг
    
    -- Цены в RSD (без НДС)
    base_price DECIMAL(12,2) NOT NULL,   -- Базовая цена
    
    -- Дополнительные услуги
    insurance_included_up_to DECIMAL(12,2) DEFAULT 15000,  -- Включенная страховка
    insurance_rate_percent DECIMAL(5,2) DEFAULT 1.0,  -- % за доп. страховку
    cod_fee DECIMAL(12,2) DEFAULT 45,  -- Комиссия за наложенный платеж
    
    -- Ограничения размеров
    max_length_cm INTEGER DEFAULT 150,
    max_width_cm INTEGER DEFAULT 60,
    max_height_cm INTEGER DEFAULT 50,
    max_dimensions_sum_cm INTEGER DEFAULT 300,
    
    -- Срок доставки
    delivery_days_min INTEGER DEFAULT 1,
    delivery_days_max INTEGER DEFAULT 2,
    
    -- Флаги
    is_active BOOLEAN DEFAULT true,
    is_special_offer BOOLEAN DEFAULT false,
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Вставляем стандартные тарифы из коммерческого предложения
INSERT INTO post_express_rates (weight_from, weight_to, base_price) VALUES
    (0, 2, 340.00),
    (2, 5, 450.00),
    (5, 10, 580.00),
    (10, 20, 790.00);

-- =====================================================
-- 5. ТАБЛИЦА ОТПРАВЛЕНИЙ
-- =====================================================
CREATE TABLE post_express_shipments (
    id SERIAL PRIMARY KEY,
    
    -- Связь с заказами
    marketplace_order_id INTEGER REFERENCES marketplace_orders(id),
    storefront_order_id BIGINT REFERENCES storefront_orders(id),
    
    -- Идентификаторы Post Express
    tracking_number VARCHAR(100) UNIQUE,  -- Трекинг номер
    barcode VARCHAR(100),                 -- Штрих-код
    post_express_id VARCHAR(100),         -- ID в системе Post Express
    
    -- Отправитель (может отличаться от дефолтного)
    sender_name VARCHAR(200) NOT NULL,
    sender_address VARCHAR(500) NOT NULL,
    sender_city VARCHAR(100) NOT NULL,
    sender_postal_code VARCHAR(20) NOT NULL,
    sender_phone VARCHAR(50) NOT NULL,
    sender_email VARCHAR(200),
    sender_location_id INTEGER REFERENCES post_express_locations(id),
    
    -- Получатель
    recipient_name VARCHAR(200) NOT NULL,
    recipient_address VARCHAR(500) NOT NULL,
    recipient_city VARCHAR(100) NOT NULL,
    recipient_postal_code VARCHAR(20) NOT NULL,
    recipient_phone VARCHAR(50) NOT NULL,
    recipient_email VARCHAR(200),
    recipient_location_id INTEGER REFERENCES post_express_locations(id),
    
    -- Параметры посылки
    weight_kg DECIMAL(10,3) NOT NULL,
    length_cm INTEGER,
    width_cm INTEGER,
    height_cm INTEGER,
    declared_value DECIMAL(12,2),
    
    -- Услуги
    service_type VARCHAR(50) DEFAULT 'danas_za_sutra',  -- danas_za_sutra, standard
    cod_amount DECIMAL(12,2),  -- Сумма наложенного платежа
    cod_reference VARCHAR(100),  -- Референс для COD
    insurance_amount DECIMAL(12,2),
    
    -- Расчет стоимости
    base_price DECIMAL(12,2) NOT NULL,
    insurance_fee DECIMAL(12,2) DEFAULT 0,
    cod_fee DECIMAL(12,2) DEFAULT 0,
    total_price DECIMAL(12,2) NOT NULL,
    
    -- Статусы
    status VARCHAR(50) DEFAULT 'created',  -- created, registered, picked_up, in_transit, delivered, failed, returned
    delivery_status VARCHAR(100),  -- Детальный статус от Post Express
    
    -- Документы
    label_url VARCHAR(500),  -- URL этикетки
    invoice_url VARCHAR(500),  -- URL накладной
    pod_url VARCHAR(500),  -- URL подтверждения доставки (Proof of Delivery)
    
    -- Временные метки
    registered_at TIMESTAMP WITH TIME ZONE,  -- Зарегистрировано в Post Express
    picked_up_at TIMESTAMP WITH TIME ZONE,   -- Забрано курьером
    delivered_at TIMESTAMP WITH TIME ZONE,   -- Доставлено
    failed_at TIMESTAMP WITH TIME ZONE,      -- Неудачная доставка
    returned_at TIMESTAMP WITH TIME ZONE,    -- Возвращено отправителю
    
    -- История статусов (JSON)
    status_history JSONB DEFAULT '[]',
    
    -- Дополнительная информация
    notes TEXT,
    internal_notes TEXT,
    delivery_instructions TEXT,
    failed_reason VARCHAR(500),
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_post_express_shipments_tracking ON post_express_shipments(tracking_number);
CREATE INDEX idx_post_express_shipments_status ON post_express_shipments(status);
CREATE INDEX idx_post_express_shipments_marketplace_order ON post_express_shipments(marketplace_order_id);
CREATE INDEX idx_post_express_shipments_storefront_order ON post_express_shipments(storefront_order_id);
CREATE INDEX idx_post_express_shipments_created ON post_express_shipments(created_at);

-- =====================================================
-- 6. ТАБЛИЦА ЛОГОВ ВЗАИМОДЕЙСТВИЯ С API
-- =====================================================
CREATE TABLE post_express_api_logs (
    id SERIAL PRIMARY KEY,
    
    -- Транзакция
    transaction_id VARCHAR(100) NOT NULL,  -- GUID транзакции
    transaction_type INTEGER NOT NULL,     -- Тип транзакции (3, 10, 63, etc.)
    
    -- Запрос/Ответ
    request_data JSONB,
    response_data JSONB,
    
    -- Статус
    status VARCHAR(50),  -- success, error, timeout
    error_message TEXT,
    
    -- Метаданные
    shipment_id INTEGER REFERENCES post_express_shipments(id),
    execution_time_ms INTEGER,
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_post_express_api_logs_transaction ON post_express_api_logs(transaction_id);
CREATE INDEX idx_post_express_api_logs_shipment ON post_express_api_logs(shipment_id);
CREATE INDEX idx_post_express_api_logs_created ON post_express_api_logs(created_at);

-- =====================================================
-- 7. ТАБЛИЦА СОБЫТИЙ ОТСЛЕЖИВАНИЯ
-- =====================================================
CREATE TABLE post_express_tracking_events (
    id SERIAL PRIMARY KEY,
    shipment_id INTEGER REFERENCES post_express_shipments(id) NOT NULL,
    
    -- Событие
    event_code VARCHAR(50) NOT NULL,
    event_description VARCHAR(500) NOT NULL,
    event_location VARCHAR(200),
    event_timestamp TIMESTAMP WITH TIME ZONE NOT NULL,
    
    -- Дополнительная информация
    additional_info JSONB,
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_post_express_tracking_events_shipment ON post_express_tracking_events(shipment_id);
CREATE INDEX idx_post_express_tracking_events_timestamp ON post_express_tracking_events(event_timestamp);

-- =====================================================
-- 8. ДОБАВЛЕНИЕ ПОЛЕЙ В СУЩЕСТВУЮЩИЕ ТАБЛИЦЫ
-- =====================================================

-- Добавление методов доставки в заказы маркетплейса
ALTER TABLE marketplace_orders 
ADD COLUMN IF NOT EXISTS delivery_method VARCHAR(50) DEFAULT 'pickup',  -- pickup, post_express, warehouse_pickup
ADD COLUMN IF NOT EXISTS delivery_cost DECIMAL(12,2) DEFAULT 0,
ADD COLUMN IF NOT EXISTS delivery_address JSONB,
ADD COLUMN IF NOT EXISTS delivery_notes TEXT;

-- Добавление методов доставки в заказы витрин
ALTER TABLE storefront_orders 
ADD COLUMN IF NOT EXISTS delivery_method VARCHAR(50) DEFAULT 'pickup',
ADD COLUMN IF NOT EXISTS delivery_cost DECIMAL(12,2) DEFAULT 0,
ADD COLUMN IF NOT EXISTS delivery_address JSONB,
ADD COLUMN IF NOT EXISTS delivery_notes TEXT;

-- =====================================================
-- 9. ТРИГГЕРЫ И ФУНКЦИИ
-- =====================================================

-- Функция для обновления updated_at
CREATE OR REPLACE FUNCTION update_post_express_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Применяем триггер к таблицам
CREATE TRIGGER update_post_express_settings_updated_at 
    BEFORE UPDATE ON post_express_settings 
    FOR EACH ROW EXECUTE FUNCTION update_post_express_updated_at();

CREATE TRIGGER update_post_express_locations_updated_at 
    BEFORE UPDATE ON post_express_locations 
    FOR EACH ROW EXECUTE FUNCTION update_post_express_updated_at();

CREATE TRIGGER update_post_express_offices_updated_at 
    BEFORE UPDATE ON post_express_offices 
    FOR EACH ROW EXECUTE FUNCTION update_post_express_updated_at();

CREATE TRIGGER update_post_express_rates_updated_at 
    BEFORE UPDATE ON post_express_rates 
    FOR EACH ROW EXECUTE FUNCTION update_post_express_updated_at();

CREATE TRIGGER update_post_express_shipments_updated_at 
    BEFORE UPDATE ON post_express_shipments 
    FOR EACH ROW EXECUTE FUNCTION update_post_express_updated_at();

-- Функция для логирования изменений статуса
CREATE OR REPLACE FUNCTION log_shipment_status_change()
RETURNS TRIGGER AS $$
BEGIN
    IF OLD.status IS DISTINCT FROM NEW.status THEN
        NEW.status_history = NEW.status_history || jsonb_build_object(
            'status', NEW.status,
            'timestamp', NOW(),
            'old_status', OLD.status
        );
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER log_post_express_shipment_status 
    BEFORE UPDATE ON post_express_shipments 
    FOR EACH ROW EXECUTE FUNCTION log_shipment_status_change();

-- =====================================================
-- 10. КОММЕНТАРИИ К ТАБЛИЦАМ
-- =====================================================

COMMENT ON TABLE post_express_settings IS 'Настройки интеграции с Post Express API';
COMMENT ON TABLE post_express_locations IS 'Справочник населенных пунктов Post Express';
COMMENT ON TABLE post_express_offices IS 'Справочник почтовых отделений Post Express';
COMMENT ON TABLE post_express_rates IS 'Тарифы на доставку Post Express';
COMMENT ON TABLE post_express_shipments IS 'Отправления через Post Express';
COMMENT ON TABLE post_express_api_logs IS 'Логи взаимодействия с Post Express API';
COMMENT ON TABLE post_express_tracking_events IS 'События отслеживания посылок Post Express';