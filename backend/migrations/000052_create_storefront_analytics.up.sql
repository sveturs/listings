-- Таблица для хранения аналитики витрин
CREATE TABLE IF NOT EXISTS storefront_analytics (
    id SERIAL PRIMARY KEY,
    storefront_id INT NOT NULL REFERENCES marketplace_storefronts(id) ON DELETE CASCADE,
    date DATE NOT NULL,
    
    -- Трафик
    page_views INT DEFAULT 0,
    unique_visitors INT DEFAULT 0,
    bounce_rate DECIMAL(5,2) DEFAULT 0,
    avg_session_time INT DEFAULT 0, -- в секундах
    
    -- Продажи
    orders_count INT DEFAULT 0,
    revenue DECIMAL(10,2) DEFAULT 0,
    avg_order_value DECIMAL(10,2) DEFAULT 0,
    conversion_rate DECIMAL(5,2) DEFAULT 0,
    
    -- JSON поля для детальной информации
    payment_methods_usage JSONB DEFAULT '{}',
    product_views INT DEFAULT 0,
    add_to_cart_count INT DEFAULT 0,
    checkout_count INT DEFAULT 0,
    traffic_sources JSONB DEFAULT '{}',
    top_products JSONB DEFAULT '[]',
    top_categories JSONB DEFAULT '[]',
    orders_by_city JSONB DEFAULT '{}',
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    -- Уникальный индекс для предотвращения дубликатов
    UNIQUE(storefront_id, date)
);

-- Индексы для быстрых запросов
CREATE INDEX idx_storefront_analytics_storefront_date ON storefront_analytics(storefront_id, date DESC);
CREATE INDEX idx_storefront_analytics_date ON storefront_analytics(date);

-- Таблица для детальных событий (для real-time аналитики)
CREATE TABLE IF NOT EXISTS storefront_events (
    id SERIAL PRIMARY KEY,
    storefront_id INT NOT NULL REFERENCES marketplace_storefronts(id) ON DELETE CASCADE,
    event_type VARCHAR(50) NOT NULL, -- page_view, product_view, add_to_cart, checkout, order
    event_data JSONB DEFAULT '{}',
    user_id INT,
    session_id VARCHAR(100),
    ip_address INET,
    user_agent TEXT,
    referrer TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Индексы для событий
CREATE INDEX idx_storefront_events_storefront ON storefront_events(storefront_id);
CREATE INDEX idx_storefront_events_type ON storefront_events(event_type);
CREATE INDEX idx_storefront_events_created ON storefront_events(created_at);
CREATE INDEX idx_storefront_events_session ON storefront_events(session_id);

-- Функция для обновления updated_at
CREATE OR REPLACE FUNCTION update_storefront_analytics_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language plpgsql;

-- Триггер для автоматического обновления updated_at
CREATE TRIGGER update_storefront_analytics_updated_at 
    BEFORE UPDATE ON storefront_analytics
    FOR EACH ROW 
    EXECUTE FUNCTION update_storefront_analytics_updated_at();