-- Создаем универсальную таблицу для истории просмотров
CREATE TABLE IF NOT EXISTS user_view_history (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    listing_id INTEGER REFERENCES marketplace_listings(id) ON DELETE CASCADE,
    category_id INTEGER REFERENCES marketplace_categories(id),
    listing_type VARCHAR(50) DEFAULT 'marketplace', -- 'marketplace', 'storefront', 'auction', etc.
    session_id VARCHAR(100),
    viewed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    view_duration_seconds INTEGER,
    interaction_type VARCHAR(50) DEFAULT 'view', -- 'view', 'click_phone', 'add_favorite', 'contact_seller'
    device_type VARCHAR(50), -- 'mobile', 'tablet', 'desktop'
    referrer VARCHAR(255),
    -- Дополнительные метаданные
    ip_address INET,
    user_agent TEXT,
    viewport_width INTEGER,
    viewport_height INTEGER,
    -- Для аналитики
    page_depth INTEGER, -- Глубина просмотра (количество страниц в сессии)
    is_return_visit BOOLEAN DEFAULT FALSE,
    source VARCHAR(100), -- 'direct', 'search', 'social', 'email', 'referral'
    medium VARCHAR(100), -- 'organic', 'cpc', 'banner', etc.
    campaign VARCHAR(200),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Индексы для производительности
CREATE INDEX idx_user_view_history_user_id ON user_view_history(user_id);
CREATE INDEX idx_user_view_history_listing_id ON user_view_history(listing_id);
CREATE INDEX idx_user_view_history_category_id ON user_view_history(category_id);
CREATE INDEX idx_user_view_history_session_id ON user_view_history(session_id);
CREATE INDEX idx_user_view_history_viewed_at ON user_view_history(viewed_at DESC);
CREATE INDEX idx_user_view_history_interaction_type ON user_view_history(interaction_type);
CREATE INDEX idx_user_view_history_device_type ON user_view_history(device_type);

-- Композитные индексы для аналитики
CREATE INDEX idx_user_view_history_user_listing ON user_view_history(user_id, listing_id, viewed_at DESC);
CREATE INDEX idx_user_view_history_category_time ON user_view_history(category_id, viewed_at DESC);
CREATE INDEX idx_user_view_history_analytics ON user_view_history(source, medium, viewed_at DESC);

-- Мигрируем данные из user_car_view_history если таблица существует
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.tables
               WHERE table_schema = 'public'
               AND table_name = 'user_car_view_history') THEN

        INSERT INTO user_view_history (
            user_id,
            listing_id,
            category_id,
            listing_type,
            session_id,
            viewed_at,
            view_duration_seconds,
            interaction_type,
            device_type,
            referrer,
            created_at
        )
        SELECT
            ucvh.user_id,
            ucvh.listing_id,
            ml.category_id, -- Получаем category_id из marketplace_listings
            'marketplace',
            ucvh.session_id,
            ucvh.viewed_at,
            ucvh.view_duration_seconds,
            'view', -- По умолчанию считаем просмотром
            ucvh.device_type,
            ucvh.referrer,
            ucvh.viewed_at
        FROM user_car_view_history ucvh
        LEFT JOIN marketplace_listings ml ON ml.id = ucvh.listing_id
        WHERE ml.id IS NOT NULL; -- Только если объявление существует

        RAISE NOTICE 'Migrated % records from user_car_view_history', (SELECT COUNT(*) FROM user_car_view_history);
    END IF;
END $$;

-- Создаем таблицу для агрегированной статистики просмотров (для быстрых запросов)
CREATE TABLE IF NOT EXISTS view_statistics (
    id SERIAL PRIMARY KEY,
    listing_id INTEGER REFERENCES marketplace_listings(id) ON DELETE CASCADE,
    category_id INTEGER REFERENCES marketplace_categories(id),
    date DATE NOT NULL,
    views_count INTEGER DEFAULT 0,
    unique_users_count INTEGER DEFAULT 0,
    unique_sessions_count INTEGER DEFAULT 0,
    avg_view_duration DECIMAL(10,2),
    mobile_views INTEGER DEFAULT 0,
    desktop_views INTEGER DEFAULT 0,
    tablet_views INTEGER DEFAULT 0,
    contact_clicks INTEGER DEFAULT 0,
    favorite_adds INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(listing_id, date)
);

-- Индексы для статистики
CREATE INDEX idx_view_statistics_listing_date ON view_statistics(listing_id, date DESC);
CREATE INDEX idx_view_statistics_category_date ON view_statistics(category_id, date DESC);

-- Функция для обновления статистики (можно вызывать периодически через cron)
CREATE OR REPLACE FUNCTION update_view_statistics(target_date DATE DEFAULT CURRENT_DATE)
RETURNS void AS $$
BEGIN
    -- Удаляем старую статистику за эту дату
    DELETE FROM view_statistics WHERE date = target_date;

    -- Вставляем новую агрегированную статистику
    INSERT INTO view_statistics (
        listing_id,
        category_id,
        date,
        views_count,
        unique_users_count,
        unique_sessions_count,
        avg_view_duration,
        mobile_views,
        desktop_views,
        tablet_views,
        contact_clicks,
        favorite_adds
    )
    SELECT
        uvh.listing_id,
        ml.category_id,
        target_date,
        COUNT(*) as views_count,
        COUNT(DISTINCT uvh.user_id) as unique_users_count,
        COUNT(DISTINCT uvh.session_id) as unique_sessions_count,
        AVG(uvh.view_duration_seconds) as avg_view_duration,
        COUNT(*) FILTER (WHERE uvh.device_type = 'mobile') as mobile_views,
        COUNT(*) FILTER (WHERE uvh.device_type = 'desktop') as desktop_views,
        COUNT(*) FILTER (WHERE uvh.device_type = 'tablet') as tablet_views,
        COUNT(*) FILTER (WHERE uvh.interaction_type = 'click_phone') as contact_clicks,
        COUNT(*) FILTER (WHERE uvh.interaction_type = 'add_favorite') as favorite_adds
    FROM user_view_history uvh
    JOIN marketplace_listings ml ON ml.id = uvh.listing_id
    WHERE DATE(uvh.viewed_at) = target_date
    GROUP BY uvh.listing_id, ml.category_id;
END;
$$ LANGUAGE plpgsql;

-- Комментарии для документации
COMMENT ON TABLE user_view_history IS 'Универсальная таблица истории просмотров для всех категорий товаров';
COMMENT ON TABLE view_statistics IS 'Агрегированная статистика просмотров для быстрых запросов';
COMMENT ON FUNCTION update_view_statistics IS 'Обновляет агрегированную статистику просмотров за указанную дату';

-- Создаем индекс для быстрого поиска популярных товаров (без условия WHERE из-за ограничений PostgreSQL)
CREATE INDEX idx_user_view_history_popular ON user_view_history(listing_id, viewed_at DESC);

-- Триггер для автоматического определения is_return_visit
CREATE OR REPLACE FUNCTION check_return_visit()
RETURNS TRIGGER AS $$
BEGIN
    -- Проверяем, был ли этот пользователь раньше
    IF NEW.user_id IS NOT NULL THEN
        NEW.is_return_visit := EXISTS (
            SELECT 1 FROM user_view_history
            WHERE user_id = NEW.user_id
            AND listing_id = NEW.listing_id
            AND id < NEW.id
            LIMIT 1
        );
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_check_return_visit
    BEFORE INSERT ON user_view_history
    FOR EACH ROW
    EXECUTE FUNCTION check_return_visit();