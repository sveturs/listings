-- Используем существующую таблицу marketplace_favorites для избранных автомобилей
-- Она уже имеет структуру: user_id, listing_id, created_at

-- Таблица для сохраненных поисков
CREATE TABLE IF NOT EXISTS saved_searches (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    filters JSONB NOT NULL DEFAULT '{}',
    search_type VARCHAR(50) DEFAULT 'cars', -- тип поиска: cars, marketplace, etc.
    notify_enabled BOOLEAN DEFAULT false,
    notify_frequency VARCHAR(20) DEFAULT 'daily', -- daily, weekly, instant
    last_notified_at TIMESTAMP WITH TIME ZONE,
    results_count INTEGER DEFAULT 0, -- количество результатов при последней проверке
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Индексы для сохраненных поисков
CREATE INDEX idx_saved_searches_user_id ON saved_searches(user_id);
CREATE INDEX idx_saved_searches_search_type ON saved_searches(search_type);
CREATE INDEX idx_saved_searches_notify_enabled ON saved_searches(notify_enabled) WHERE notify_enabled = true;
CREATE INDEX idx_saved_searches_created_at ON saved_searches(created_at DESC);
CREATE INDEX idx_saved_searches_filters ON saved_searches USING gin(filters);

-- Таблица для истории уведомлений о поисках
CREATE TABLE IF NOT EXISTS saved_search_notifications (
    id SERIAL PRIMARY KEY,
    saved_search_id INTEGER NOT NULL REFERENCES saved_searches(id) ON DELETE CASCADE,
    new_listings_count INTEGER NOT NULL DEFAULT 0,
    notification_sent BOOLEAN DEFAULT false,
    sent_at TIMESTAMP WITH TIME ZONE,
    error_message TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Индексы для истории уведомлений
CREATE INDEX idx_saved_search_notifications_saved_search_id ON saved_search_notifications(saved_search_id);
CREATE INDEX idx_saved_search_notifications_created_at ON saved_search_notifications(created_at DESC);

-- Таблица для истории просмотров автомобилей (для рекомендаций)
CREATE TABLE IF NOT EXISTS user_car_view_history (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    listing_id INTEGER NOT NULL REFERENCES marketplace_listings(id) ON DELETE CASCADE,
    session_id VARCHAR(100), -- для анонимных пользователей
    viewed_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    view_duration_seconds INTEGER, -- время просмотра в секундах
    referrer VARCHAR(255), -- откуда пришел пользователь
    device_type VARCHAR(50), -- mobile, tablet, desktop
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Индексы для истории просмотров
CREATE INDEX idx_user_car_view_history_user_id ON user_car_view_history(user_id) WHERE user_id IS NOT NULL;
CREATE INDEX idx_user_car_view_history_listing_id ON user_car_view_history(listing_id);
CREATE INDEX idx_user_car_view_history_session_id ON user_car_view_history(session_id) WHERE session_id IS NOT NULL;
CREATE INDEX idx_user_car_view_history_viewed_at ON user_car_view_history(viewed_at DESC);

-- Триггер для обновления updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_saved_searches_updated_at BEFORE UPDATE ON saved_searches
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Добавляем комментарии к таблицам
COMMENT ON TABLE saved_searches IS 'Сохраненные поиски пользователей с возможностью уведомлений';
COMMENT ON TABLE saved_search_notifications IS 'История уведомлений о новых объявлениях';
COMMENT ON TABLE user_car_view_history IS 'История просмотров автомобилей для персонализации и рекомендаций';