-- Создание таблицы для отслеживания сессий оптимизации весов поиска
-- Эта таблица хранит информацию о запущенных процессах оптимизации

CREATE TABLE IF NOT EXISTS search_optimization_sessions (
    id BIGSERIAL PRIMARY KEY,
    
    -- Статус сессии: 'running', 'completed', 'failed', 'cancelled'
    status VARCHAR(20) NOT NULL DEFAULT 'running'
        CHECK (status IN ('running', 'completed', 'failed', 'cancelled')),
    
    -- Время начала и окончания оптимизации
    start_time TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    end_time TIMESTAMP WITH TIME ZONE,
    
    -- Количество полей для обработки
    total_fields INTEGER NOT NULL DEFAULT 0,
    processed_fields INTEGER NOT NULL DEFAULT 0,
    
    -- Результаты оптимизации в JSON формате
    results JSONB,
    
    -- Сообщение об ошибке (если статус 'failed')
    error_message TEXT,
    
    -- Кто запустил оптимизацию
    created_by INTEGER NOT NULL REFERENCES admin_users(id) ON DELETE RESTRICT,
    
    -- Временные метки
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Индексы для быстрого поиска
CREATE INDEX IF NOT EXISTS idx_search_optimization_sessions_status ON search_optimization_sessions(status);
CREATE INDEX IF NOT EXISTS idx_search_optimization_sessions_start_time ON search_optimization_sessions(start_time);
CREATE INDEX IF NOT EXISTS idx_search_optimization_sessions_created_by ON search_optimization_sessions(created_by);

-- Триггер для автоматического обновления updated_at
CREATE OR REPLACE FUNCTION update_search_optimization_sessions_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trigger_update_search_optimization_sessions_updated_at ON search_optimization_sessions;
CREATE TRIGGER trigger_update_search_optimization_sessions_updated_at
    BEFORE UPDATE ON search_optimization_sessions
    FOR EACH ROW
    EXECUTE FUNCTION update_search_optimization_sessions_updated_at();

-- Комментарии к таблице и колонкам
COMMENT ON TABLE search_optimization_sessions IS 'Сессии оптимизации весов поиска';
COMMENT ON COLUMN search_optimization_sessions.status IS 'Статус сессии: running, completed, failed, cancelled';
COMMENT ON COLUMN search_optimization_sessions.total_fields IS 'Общее количество полей для оптимизации';
COMMENT ON COLUMN search_optimization_sessions.processed_fields IS 'Количество обработанных полей';
COMMENT ON COLUMN search_optimization_sessions.results IS 'Результаты оптимизации в JSON формате';
COMMENT ON COLUMN search_optimization_sessions.error_message IS 'Сообщение об ошибке (если статус failed)';