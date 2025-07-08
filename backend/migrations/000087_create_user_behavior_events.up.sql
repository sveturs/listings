-- Создание таблицы для сбора всех поведенческих событий пользователей
-- Эта таблица используется для отслеживания действий пользователей: поиск, клики, просмотры, покупки
-- Данные из этой таблицы используются для улучшения поисковой релевантности и персонализации

CREATE TABLE IF NOT EXISTS user_behavior_events (
    id BIGSERIAL PRIMARY KEY,
    
    -- Тип события: search_performed, result_clicked, item_viewed, item_purchased, item_added_to_cart, etc.
    event_type VARCHAR(50) NOT NULL,
    
    -- Ссылка на пользователя (может быть NULL для анонимных пользователей)
    user_id INTEGER REFERENCES users(id) ON DELETE SET NULL,
    
    -- Идентификатор сессии для группировки событий одного визита
    session_id VARCHAR(100) NOT NULL,
    
    -- Поисковый запрос (заполняется для событий поиска и кликов по результатам)
    search_query TEXT,
    
    -- ID элемента в формате ml_123 (marketplace) или sp_456 (storefront)
    item_id VARCHAR(50),
    
    -- Тип элемента: marketplace или storefront
    item_type VARCHAR(20) CHECK (item_type IN ('marketplace', 'storefront', NULL)),
    
    -- Позиция элемента в результатах поиска (для событий клика)
    position INTEGER,
    
    -- Дополнительные данные события в JSON формате
    -- Может содержать: device_type, browser, referer, search_filters, time_on_page, etc.
    metadata JSONB DEFAULT '{}',
    
    -- Время создания события
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Индексы для быстрого поиска и аналитики
CREATE INDEX idx_user_behavior_events_event_type ON user_behavior_events(event_type);
CREATE INDEX idx_user_behavior_events_user_id ON user_behavior_events(user_id) WHERE user_id IS NOT NULL;
CREATE INDEX idx_user_behavior_events_session_id ON user_behavior_events(session_id);
CREATE INDEX idx_user_behavior_events_created_at ON user_behavior_events(created_at);

-- Составной индекс для поиска событий по запросу и типу
CREATE INDEX idx_user_behavior_events_search_query_type ON user_behavior_events(search_query, event_type) 
    WHERE search_query IS NOT NULL;

-- Индекс для быстрого поиска событий по элементам
CREATE INDEX idx_user_behavior_events_item ON user_behavior_events(item_id, item_type) 
    WHERE item_id IS NOT NULL;

-- Комментарии к таблице и колонкам
COMMENT ON TABLE user_behavior_events IS 'Таблица для хранения всех поведенческих событий пользователей';
COMMENT ON COLUMN user_behavior_events.event_type IS 'Тип события: search_performed, result_clicked, item_viewed, item_purchased, etc.';
COMMENT ON COLUMN user_behavior_events.session_id IS 'Идентификатор сессии для группировки событий одного визита';
COMMENT ON COLUMN user_behavior_events.position IS 'Позиция элемента в результатах поиска (для событий клика)';
COMMENT ON COLUMN user_behavior_events.metadata IS 'Дополнительные данные события: device_type, browser, referer, search_filters, etc.';