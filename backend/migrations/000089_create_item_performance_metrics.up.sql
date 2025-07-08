-- Создание таблицы для метрик производительности товаров
-- Эта таблица содержит агрегированные данные о том, как товары показываются в поиске
-- Используется для оптимизации ранжирования и выявления популярных/непопулярных товаров

CREATE TABLE IF NOT EXISTS item_performance_metrics (
    id BIGSERIAL PRIMARY KEY,
    
    -- ID элемента в формате ml_123 (marketplace) или sp_456 (storefront)
    item_id VARCHAR(50) NOT NULL,
    
    -- Тип элемента: marketplace или storefront
    item_type VARCHAR(20) NOT NULL CHECK (item_type IN ('marketplace', 'storefront')),
    
    -- Количество показов товара в результатах поиска
    impressions INTEGER DEFAULT 0,
    
    -- Количество кликов по товару
    clicks INTEGER DEFAULT 0,
    
    -- Click-Through Rate товара (отношение кликов к показам)
    ctr FLOAT DEFAULT 0,
    
    -- Количество конверсий (покупок)
    conversions INTEGER DEFAULT 0,
    
    -- Средняя позиция товара в результатах поиска
    avg_position FLOAT DEFAULT 0,
    
    -- Период агрегации данных
    period_start TIMESTAMP WITH TIME ZONE NOT NULL,
    period_end TIMESTAMP WITH TIME ZONE NOT NULL,
    
    -- Временные метки
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Уникальный индекс для предотвращения дублирования метрик за один период
CREATE UNIQUE INDEX IF NOT EXISTS idx_item_performance_metrics_unique ON item_performance_metrics(item_id, period_start);

-- Индексы для быстрого поиска
CREATE INDEX IF NOT EXISTS idx_item_performance_metrics_item ON item_performance_metrics(item_id, item_type);
CREATE INDEX IF NOT EXISTS idx_item_performance_metrics_period ON item_performance_metrics(period_start, period_end);
CREATE INDEX IF NOT EXISTS idx_item_performance_metrics_ctr ON item_performance_metrics(ctr DESC);
CREATE INDEX IF NOT EXISTS idx_item_performance_metrics_impressions ON item_performance_metrics(impressions DESC);
CREATE INDEX IF NOT EXISTS idx_item_performance_metrics_conversions ON item_performance_metrics(conversions DESC);

-- Триггер для автоматического обновления updated_at
CREATE OR REPLACE FUNCTION update_item_performance_metrics_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trigger_update_item_performance_metrics_updated_at ON trigger_update_item_performance_metrics_updated_at;
CREATE TRIGGER trigger_update_item_performance_metrics_updated_at
    BEFORE UPDATE ON item_performance_metrics
    FOR EACH ROW
    EXECUTE FUNCTION update_item_performance_metrics_updated_at();

-- Комментарии к таблице и колонкам
COMMENT ON TABLE item_performance_metrics IS 'Метрики производительности товаров в поисковой системе';
COMMENT ON COLUMN item_performance_metrics.item_id IS 'ID элемента в формате ml_123 (marketplace) или sp_456 (storefront)';
COMMENT ON COLUMN item_performance_metrics.item_type IS 'Тип элемента: marketplace или storefront';
COMMENT ON COLUMN item_performance_metrics.impressions IS 'Количество показов товара в результатах поиска';
COMMENT ON COLUMN item_performance_metrics.clicks IS 'Количество кликов по товару';
COMMENT ON COLUMN item_performance_metrics.ctr IS 'Click-Through Rate товара (отношение кликов к показам)';
COMMENT ON COLUMN item_performance_metrics.conversions IS 'Количество конверсий (покупок)';
COMMENT ON COLUMN item_performance_metrics.avg_position IS 'Средняя позиция товара в результатах поиска';