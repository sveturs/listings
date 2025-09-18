-- =====================================================
-- VIBER BOT И СИСТЕМА ТРЕКИНГА КУРЬЕРОВ
-- =====================================================

-- 1. Таблица пользователей Viber
CREATE TABLE IF NOT EXISTS viber_users (
    id SERIAL PRIMARY KEY,
    viber_id VARCHAR(100) UNIQUE NOT NULL,
    user_id INT REFERENCES users(id),
    name VARCHAR(255),
    avatar_url TEXT,
    language VARCHAR(10) DEFAULT 'sr',
    country_code VARCHAR(5),
    api_version INT DEFAULT 1,
    subscribed BOOLEAN DEFAULT true,
    subscribed_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    last_session_at TIMESTAMP WITH TIME ZONE,
    conversation_started_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 2. Таблица сессий Viber (для отслеживания 24-часовых бесплатных окон)
CREATE TABLE IF NOT EXISTS viber_sessions (
    id SERIAL PRIMARY KEY,
    viber_user_id INT REFERENCES viber_users(id) ON DELETE CASCADE,
    started_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    last_message_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP WITH TIME ZONE DEFAULT (CURRENT_TIMESTAMP + INTERVAL '24 hours'),
    message_count INT DEFAULT 0,
    context JSONB DEFAULT '{}',
    active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 3. Таблица сообщений Viber (для аналитики и истории)
CREATE TABLE IF NOT EXISTS viber_messages (
    id SERIAL PRIMARY KEY,
    viber_user_id INT REFERENCES viber_users(id) ON DELETE CASCADE,
    session_id INT REFERENCES viber_sessions(id) ON DELETE SET NULL,
    message_token VARCHAR(100) UNIQUE,
    direction VARCHAR(20) CHECK (direction IN ('incoming', 'outgoing')),
    message_type VARCHAR(50), -- text, picture, rich_media, etc
    content TEXT,
    rich_media JSONB,
    is_billable BOOLEAN DEFAULT false, -- true если сообщение платное
    status VARCHAR(20), -- sent, delivered, seen, failed
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 4. Таблица курьеров
CREATE TABLE IF NOT EXISTS couriers (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id),
    name VARCHAR(255) NOT NULL,
    phone VARCHAR(50),
    photo_url TEXT,
    vehicle_type VARCHAR(50) CHECK (vehicle_type IN ('bike', 'car', 'scooter', 'on_foot', 'van')),
    vehicle_number VARCHAR(50),
    is_online BOOLEAN DEFAULT false,
    is_available BOOLEAN DEFAULT true,
    current_latitude NUMERIC(10, 8),
    current_longitude NUMERIC(11, 8),
    current_heading INT, -- направление движения в градусах (0-360)
    current_speed NUMERIC(5, 2), -- скорость в км/ч
    last_location_update TIMESTAMP WITH TIME ZONE,
    rating NUMERIC(3, 2) DEFAULT 5.0,
    total_deliveries INT DEFAULT 0,
    total_distance_km NUMERIC(10, 2) DEFAULT 0,
    working_hours JSONB DEFAULT '{}', -- {"monday": {"start": "09:00", "end": "18:00"}}
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 5. Таблица зон доставки курьеров
CREATE TABLE IF NOT EXISTS courier_zones (
    id SERIAL PRIMARY KEY,
    courier_id INT REFERENCES couriers(id) ON DELETE CASCADE,
    zone_name VARCHAR(100),
    polygon JSONB NOT NULL, -- GeoJSON polygon
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 6. Таблица доставок
CREATE TABLE IF NOT EXISTS deliveries (
    id SERIAL PRIMARY KEY,
    order_id INT REFERENCES storefront_orders(id),
    courier_id INT REFERENCES couriers(id),
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    -- pending, assigned, accepted, picked_up, in_transit, delivered, cancelled, failed
    CONSTRAINT delivery_status_check CHECK (status IN (
        'pending', 'assigned', 'accepted', 'picked_up',
        'in_transit', 'delivered', 'cancelled', 'failed'
    )),

    -- Адреса и координаты
    pickup_address TEXT NOT NULL,
    pickup_latitude NUMERIC(10, 8),
    pickup_longitude NUMERIC(11, 8),
    pickup_contact_name VARCHAR(255),
    pickup_contact_phone VARCHAR(50),

    delivery_address TEXT NOT NULL,
    delivery_latitude NUMERIC(10, 8),
    delivery_longitude NUMERIC(11, 8),
    delivery_contact_name VARCHAR(255),
    delivery_contact_phone VARCHAR(50),

    -- Временные метки
    assigned_at TIMESTAMP WITH TIME ZONE,
    accepted_at TIMESTAMP WITH TIME ZONE,
    picked_up_at TIMESTAMP WITH TIME ZONE,
    delivered_at TIMESTAMP WITH TIME ZONE,
    cancelled_at TIMESTAMP WITH TIME ZONE,

    -- Расчётные данные
    estimated_pickup_time TIMESTAMP WITH TIME ZONE,
    estimated_delivery_time TIMESTAMP WITH TIME ZONE,
    actual_delivery_time TIMESTAMP WITH TIME ZONE,

    -- Трекинг
    tracking_token VARCHAR(100) UNIQUE NOT NULL DEFAULT gen_random_uuid()::text,
    tracking_url TEXT,
    share_location_enabled BOOLEAN DEFAULT true,

    -- Метрики
    distance_meters INT,
    duration_seconds INT,
    route_polyline TEXT, -- Encoded polyline маршрута

    -- Оплата курьеру
    courier_fee NUMERIC(10, 2),
    courier_tip NUMERIC(10, 2) DEFAULT 0,

    -- Дополнительная информация
    notes TEXT,
    package_size VARCHAR(20) CHECK (package_size IN ('small', 'medium', 'large', 'xl')),
    package_weight_kg NUMERIC(6, 2),
    requires_signature BOOLEAN DEFAULT false,
    photo_proof_url TEXT, -- фото доставленного заказа

    -- Рейтинг
    customer_rating INT CHECK (customer_rating BETWEEN 1 AND 5),
    customer_feedback TEXT,

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 7. История локаций курьера
CREATE TABLE IF NOT EXISTS courier_location_history (
    id SERIAL PRIMARY KEY,
    delivery_id INT REFERENCES deliveries(id) ON DELETE CASCADE,
    courier_id INT REFERENCES couriers(id) ON DELETE CASCADE,
    latitude NUMERIC(10, 8) NOT NULL,
    longitude NUMERIC(11, 8) NOT NULL,
    altitude_meters NUMERIC(7, 2),
    speed_kmh NUMERIC(5, 2),
    heading INT CHECK (heading BETWEEN 0 AND 360),
    accuracy_meters NUMERIC(6, 2),
    battery_level INT,
    is_mock_location BOOLEAN DEFAULT false,
    recorded_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    -- Оптимизация: храним только каждую N-ную точку для долгих поездок
    is_key_point BOOLEAN DEFAULT false -- важные точки: старт, финиш, повороты
);

-- 8. Таблица сессий трекинга через Viber
CREATE TABLE IF NOT EXISTS viber_tracking_sessions (
    id SERIAL PRIMARY KEY,
    viber_user_id INT REFERENCES viber_users(id) ON DELETE CASCADE,
    delivery_id INT REFERENCES deliveries(id) ON DELETE CASCADE,
    tracking_token VARCHAR(100) NOT NULL,
    started_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    last_viewed_at TIMESTAMP WITH TIME ZONE,
    page_views INT DEFAULT 1,
    is_active BOOLEAN DEFAULT true,
    device_info JSONB, -- user agent, screen size, etc
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 9. WebSocket соединения для real-time трекинга
CREATE TABLE IF NOT EXISTS tracking_websocket_connections (
    id SERIAL PRIMARY KEY,
    connection_id VARCHAR(100) UNIQUE NOT NULL,
    delivery_id INT REFERENCES deliveries(id) ON DELETE CASCADE,
    client_type VARCHAR(20) CHECK (client_type IN ('customer', 'courier', 'merchant', 'admin')),
    user_id INT REFERENCES users(id),
    viber_user_id INT REFERENCES viber_users(id),
    connected_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    disconnected_at TIMESTAMP WITH TIME ZONE,
    last_ping_at TIMESTAMP WITH TIME ZONE,
    ip_address INET,
    user_agent TEXT
);

-- 10. Уведомления о доставке
CREATE TABLE IF NOT EXISTS delivery_notifications (
    id SERIAL PRIMARY KEY,
    delivery_id INT REFERENCES deliveries(id) ON DELETE CASCADE,
    viber_user_id INT REFERENCES viber_users(id),
    notification_type VARCHAR(50), -- assigned, picked_up, nearby, delivered, etc
    sent_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    is_read BOOLEAN DEFAULT false,
    read_at TIMESTAMP WITH TIME ZONE
);

-- =====================================================
-- ИНДЕКСЫ ДЛЯ ПРОИЗВОДИТЕЛЬНОСТИ
-- =====================================================

-- Viber индексы
CREATE INDEX idx_viber_users_viber_id ON viber_users(viber_id);
CREATE INDEX idx_viber_users_user_id ON viber_users(user_id);
CREATE INDEX idx_viber_sessions_active ON viber_sessions(viber_user_id, active) WHERE active = true;
CREATE INDEX idx_viber_sessions_expires ON viber_sessions(expires_at) WHERE active = true;
CREATE INDEX idx_viber_messages_user_session ON viber_messages(viber_user_id, session_id, created_at DESC);

-- Courier индексы
CREATE INDEX idx_couriers_online_available ON couriers(is_online, is_available)
    WHERE is_online = true AND is_available = true;
CREATE INDEX idx_couriers_location ON couriers(current_latitude, current_longitude)
    WHERE is_online = true;
CREATE INDEX idx_couriers_user_id ON couriers(user_id);

-- Delivery индексы
CREATE INDEX idx_deliveries_status ON deliveries(status)
    WHERE status IN ('assigned', 'accepted', 'picked_up', 'in_transit');
CREATE INDEX idx_deliveries_courier ON deliveries(courier_id, status);
CREATE INDEX idx_deliveries_tracking_token ON deliveries(tracking_token);
CREATE INDEX idx_deliveries_order ON deliveries(order_id);

-- Location history индексы
CREATE INDEX idx_location_history_delivery ON courier_location_history(delivery_id, recorded_at DESC);
CREATE INDEX idx_location_history_courier ON courier_location_history(courier_id, recorded_at DESC);
CREATE INDEX idx_location_history_key_points ON courier_location_history(delivery_id, is_key_point)
    WHERE is_key_point = true;

-- Tracking индексы
CREATE INDEX idx_tracking_sessions_active ON viber_tracking_sessions(delivery_id, is_active)
    WHERE is_active = true;
CREATE INDEX idx_websocket_connections_active ON tracking_websocket_connections(delivery_id)
    WHERE disconnected_at IS NULL;

-- =====================================================
-- ТРИГГЕРЫ
-- =====================================================

-- Автообновление updated_at
CREATE TRIGGER update_viber_users_updated_at BEFORE UPDATE ON viber_users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_couriers_updated_at BEFORE UPDATE ON couriers
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_deliveries_updated_at BEFORE UPDATE ON deliveries
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Функция для автоматического закрытия истёкших сессий
CREATE OR REPLACE FUNCTION close_expired_viber_sessions() RETURNS void AS $$
BEGIN
    UPDATE viber_sessions
    SET active = false
    WHERE active = true
    AND expires_at < CURRENT_TIMESTAMP;
END;
$$ LANGUAGE plpgsql;

-- Функция для расчёта расстояния между координатами (Haversine formula)
CREATE OR REPLACE FUNCTION calculate_distance(
    lat1 NUMERIC, lon1 NUMERIC,
    lat2 NUMERIC, lon2 NUMERIC
) RETURNS NUMERIC AS $$
DECLARE
    R CONSTANT NUMERIC := 6371000; -- Радиус Земли в метрах
    phi1 NUMERIC;
    phi2 NUMERIC;
    delta_phi NUMERIC;
    delta_lambda NUMERIC;
    a NUMERIC;
    c NUMERIC;
BEGIN
    phi1 := radians(lat1);
    phi2 := radians(lat2);
    delta_phi := radians(lat2 - lat1);
    delta_lambda := radians(lon2 - lon1);

    a := sin(delta_phi/2) * sin(delta_phi/2) +
         cos(phi1) * cos(phi2) *
         sin(delta_lambda/2) * sin(delta_lambda/2);

    c := 2 * atan2(sqrt(a), sqrt(1-a));

    RETURN R * c; -- Расстояние в метрах
END;
$$ LANGUAGE plpgsql IMMUTABLE;

-- =====================================================
-- НАЧАЛЬНЫЕ ДАННЫЕ
-- =====================================================

-- Добавляем тестового курьера
INSERT INTO couriers (name, phone, vehicle_type, is_online, is_available)
VALUES ('Тестовый Курьер', '+381601234567', 'bike', false, true)
ON CONFLICT DO NOTHING;