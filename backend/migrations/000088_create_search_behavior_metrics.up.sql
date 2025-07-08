-- Создание таблицы для агрегированных метрик поведения при поиске
-- Эта таблица содержит предрассчитанные метрики для каждого поискового запроса
-- Используется для быстрого анализа эффективности поиска и оптимизации результатов

CREATE TABLE IF NOT EXISTS search_behavior_metrics (
    id BIGSERIAL PRIMARY KEY,
    
    -- Поисковый запрос
    search_query TEXT NOT NULL,
    
    -- Общее количество выполненных поисков
    total_searches INTEGER DEFAULT 0,
    
    -- Общее количество кликов по результатам
    total_clicks INTEGER DEFAULT 0,
    
    -- Click-Through Rate (отношение кликов к поискам)
    ctr FLOAT DEFAULT 0,
    
    -- Средняя позиция, по которой кликают пользователи
    avg_click_position FLOAT DEFAULT 0,
    
    -- Количество конверсий (покупок после поиска)
    conversions INTEGER DEFAULT 0,
    
    -- Коэффициент конверсии (отношение покупок к поискам)
    conversion_rate FLOAT DEFAULT 0,
    
    -- Период агрегации данных
    period_start TIMESTAMP WITH TIME ZONE NOT NULL,
    period_end TIMESTAMP WITH TIME ZONE NOT NULL,
    
    -- Временные метки
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Уникальный индекс для предотвращения дублирования метрик за один период
CREATE UNIQUE INDEX IF NOT EXISTS idx_search_behavior_metrics_unique ON search_behavior_metrics(search_query, period_start);

-- Индексы для быстрого поиска
CREATE INDEX IF NOT EXISTS idx_search_behavior_metrics_query ON search_behavior_metrics(search_query);
CREATE INDEX IF NOT EXISTS idx_search_behavior_metrics_period ON search_behavior_metrics(period_start, period_end);
CREATE INDEX IF NOT EXISTS idx_search_behavior_metrics_ctr ON search_behavior_metrics(ctr DESC);
CREATE INDEX IF NOT EXISTS idx_search_behavior_metrics_conversions ON search_behavior_metrics(conversions DESC);

-- Триггер для автоматического обновления updated_at
CREATE OR REPLACE FUNCTION update_search_behavior_metrics_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER IF NOT EXISTS trigger_update_search_behavior_metrics_updated_at
    BEFORE UPDATE ON search_behavior_metrics
    FOR EACH ROW
    EXECUTE FUNCTION update_search_behavior_metrics_updated_at();

-- Комментарии к таблице и колонкам
COMMENT ON TABLE search_behavior_metrics IS 'Агрегированные метрики поведения пользователей при поиске';
COMMENT ON COLUMN search_behavior_metrics.search_query IS 'Поисковый запрос';
COMMENT ON COLUMN search_behavior_metrics.total_searches IS 'Общее количество выполненных поисков';
COMMENT ON COLUMN search_behavior_metrics.total_clicks IS 'Общее количество кликов по результатам';
COMMENT ON COLUMN search_behavior_metrics.ctr IS 'Click-Through Rate (отношение кликов к поискам)';
COMMENT ON COLUMN search_behavior_metrics.avg_click_position IS 'Средняя позиция, по которой кликают пользователи';
COMMENT ON COLUMN search_behavior_metrics.conversions IS 'Количество конверсий (покупок после поиска)';
COMMENT ON COLUMN search_behavior_metrics.conversion_rate IS 'Коэффициент конверсии (отношение покупок к поискам)';