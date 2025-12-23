-- Таблица событий витрин для аналитики
CREATE TABLE IF NOT EXISTS storefront_events (
    id BIGSERIAL PRIMARY KEY,
    storefront_id INTEGER NOT NULL REFERENCES storefronts(id) ON DELETE CASCADE,
    event_type VARCHAR(50) NOT NULL,
    event_data JSONB DEFAULT '{}',
    user_id INTEGER,
    session_id VARCHAR(100) NOT NULL,
    ip_address VARCHAR(45),
    user_agent TEXT,
    referrer TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Индексы для быстрого поиска
CREATE INDEX idx_storefront_events_storefront_id ON storefront_events(storefront_id);
CREATE INDEX idx_storefront_events_event_type ON storefront_events(event_type);
CREATE INDEX idx_storefront_events_session_id ON storefront_events(session_id);
CREATE INDEX idx_storefront_events_created_at ON storefront_events(created_at);
CREATE INDEX idx_storefront_events_user_id ON storefront_events(user_id) WHERE user_id IS NOT NULL;

-- Составной индекс для аналитических запросов
CREATE INDEX idx_storefront_events_analytics ON storefront_events(storefront_id, event_type, created_at DESC);

COMMENT ON TABLE storefront_events IS 'События витрин для аналитики (page_view, product_view, add_to_cart, checkout, order)';
COMMENT ON COLUMN storefront_events.event_type IS 'Тип события: page_view, product_view, add_to_cart, checkout, order';
COMMENT ON COLUMN storefront_events.event_data IS 'Дополнительные данные события в формате JSON';
