-- Создание таблицы для логирования поисковых запросов
CREATE TABLE IF NOT EXISTS search_logs (
    id BIGSERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE SET NULL,
    session_id VARCHAR(255),
    query_text TEXT NOT NULL,
    filters JSONB,
    category_id INTEGER REFERENCES marketplace_categories(id) ON DELETE SET NULL,
    location JSONB, -- {country, region, city, lat, lon}
    results_count INTEGER NOT NULL DEFAULT 0,
    response_time_ms INTEGER NOT NULL,
    page INTEGER NOT NULL DEFAULT 1,
    per_page INTEGER NOT NULL DEFAULT 20,
    sort_by VARCHAR(50),
    user_agent TEXT,
    ip_address INET,
    referer TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Индексы для search_logs
CREATE INDEX idx_search_logs_created_at ON search_logs(created_at DESC);
CREATE INDEX idx_search_logs_user_id ON search_logs(user_id) WHERE user_id IS NOT NULL;
CREATE INDEX idx_search_logs_session_id ON search_logs(session_id);
CREATE INDEX idx_search_logs_category_id ON search_logs(category_id) WHERE category_id IS NOT NULL;
CREATE INDEX idx_search_logs_query_text_gin ON search_logs USING gin(to_tsvector('russian', query_text));
CREATE INDEX idx_search_logs_filters ON search_logs USING gin(filters);
CREATE INDEX idx_search_logs_location ON search_logs USING gin(location);

-- Таблица для агрегированной аналитики поисковых запросов
CREATE TABLE IF NOT EXISTS search_analytics (
    id BIGSERIAL PRIMARY KEY,
    date DATE NOT NULL,
    hour INTEGER NOT NULL CHECK (hour >= 0 AND hour <= 23),
    query_text TEXT NOT NULL,
    category_id INTEGER REFERENCES marketplace_categories(id) ON DELETE SET NULL,
    location_country VARCHAR(100),
    location_region VARCHAR(100),
    location_city VARCHAR(100),
    search_count INTEGER NOT NULL DEFAULT 0,
    unique_users_count INTEGER NOT NULL DEFAULT 0,
    unique_sessions_count INTEGER NOT NULL DEFAULT 0,
    avg_results_count NUMERIC(10,2),
    avg_response_time_ms NUMERIC(10,2),
    zero_results_count INTEGER NOT NULL DEFAULT 0,
    click_through_rate NUMERIC(5,2), -- процент кликов по результатам
    conversion_rate NUMERIC(5,2), -- процент конверсии в контакты/покупки
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    UNIQUE(date, hour, query_text, category_id, location_country, location_region, location_city)
);

-- Индексы для search_analytics
CREATE INDEX idx_search_analytics_date ON search_analytics(date DESC);
CREATE INDEX idx_search_analytics_date_hour ON search_analytics(date DESC, hour);
CREATE INDEX idx_search_analytics_query_text ON search_analytics(query_text);
CREATE INDEX idx_search_analytics_category_id ON search_analytics(category_id) WHERE category_id IS NOT NULL;
CREATE INDEX idx_search_analytics_location ON search_analytics(location_country, location_region, location_city);
CREATE INDEX idx_search_analytics_search_count ON search_analytics(search_count DESC);

-- Таблица для отслеживания кликов по результатам поиска
CREATE TABLE IF NOT EXISTS search_result_clicks (
    id BIGSERIAL PRIMARY KEY,
    search_log_id BIGINT NOT NULL REFERENCES search_logs(id) ON DELETE CASCADE,
    listing_id INTEGER NOT NULL REFERENCES marketplace_listings(id) ON DELETE CASCADE,
    position INTEGER NOT NULL, -- позиция в результатах поиска
    clicked_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Индексы для search_result_clicks
CREATE INDEX idx_search_result_clicks_search_log_id ON search_result_clicks(search_log_id);
CREATE INDEX idx_search_result_clicks_listing_id ON search_result_clicks(listing_id);
CREATE INDEX idx_search_result_clicks_clicked_at ON search_result_clicks(clicked_at DESC);

-- Таблица для популярных поисковых запросов
CREATE TABLE IF NOT EXISTS search_trending_queries (
    id BIGSERIAL PRIMARY KEY,
    query_text TEXT NOT NULL,
    category_id INTEGER REFERENCES marketplace_categories(id) ON DELETE SET NULL,
    location_country VARCHAR(100),
    trend_score NUMERIC(10,2) NOT NULL DEFAULT 0, -- рассчитывается на основе частоты и новизны
    search_count_24h INTEGER NOT NULL DEFAULT 0,
    search_count_7d INTEGER NOT NULL DEFAULT 0,
    search_count_30d INTEGER NOT NULL DEFAULT 0,
    first_seen_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    last_seen_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    UNIQUE(query_text, category_id, location_country)
);

-- Индексы для search_trending_queries
CREATE INDEX idx_search_trending_queries_trend_score ON search_trending_queries(trend_score DESC);
CREATE INDEX idx_search_trending_queries_category_id ON search_trending_queries(category_id) WHERE category_id IS NOT NULL;
CREATE INDEX idx_search_trending_queries_location_country ON search_trending_queries(location_country);
CREATE INDEX idx_search_trending_queries_updated_at ON search_trending_queries(updated_at DESC);

-- Функция для обновления updated_at
CREATE OR REPLACE FUNCTION update_search_analytics_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Триггер для search_analytics
CREATE TRIGGER update_search_analytics_updated_at_trigger
    BEFORE UPDATE ON search_analytics
    FOR EACH ROW
    EXECUTE FUNCTION update_search_analytics_updated_at();

-- Триггер для search_trending_queries
CREATE TRIGGER update_search_trending_queries_updated_at_trigger
    BEFORE UPDATE ON search_trending_queries
    FOR EACH ROW
    EXECUTE FUNCTION update_search_analytics_updated_at();

-- Партиционирование search_logs по месяцам для лучшей производительности
-- (опционально, можно включить позже при большом объеме данных)
-- CREATE TABLE search_logs_2025_01 PARTITION OF search_logs
-- FOR VALUES FROM ('2025-01-01') TO ('2025-02-01');