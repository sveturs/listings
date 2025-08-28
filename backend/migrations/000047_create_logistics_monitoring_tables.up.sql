-- Таблица для хранения агрегированных метрик логистики
CREATE TABLE IF NOT EXISTS logistics_metrics (
    id SERIAL PRIMARY KEY,
    date DATE NOT NULL UNIQUE,
    
    -- Общие метрики
    total_shipments INT DEFAULT 0,
    delivered INT DEFAULT 0,
    in_transit INT DEFAULT 0,
    problems INT DEFAULT 0,
    returns INT DEFAULT 0,
    cancelled INT DEFAULT 0,
    
    -- Метрики по курьерским службам
    bex_shipments INT DEFAULT 0,
    bex_delivered INT DEFAULT 0,
    postexpress_shipments INT DEFAULT 0,
    postexpress_delivered INT DEFAULT 0,
    
    -- Временные метрики (в часах)
    avg_delivery_time_hours FLOAT,
    min_delivery_time_hours FLOAT,
    max_delivery_time_hours FLOAT,
    
    -- Финансовые метрики
    total_shipping_cost DECIMAL(10, 2) DEFAULT 0,
    total_cod_collected DECIMAL(10, 2) DEFAULT 0,
    total_return_cost DECIMAL(10, 2) DEFAULT 0,
    
    -- Метаданные
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Индекс для быстрого поиска по дате
CREATE INDEX idx_logistics_metrics_date ON logistics_metrics(date DESC);

-- Таблица для отслеживания проблемных отправлений
CREATE TABLE IF NOT EXISTS problem_shipments (
    id SERIAL PRIMARY KEY,
    
    -- Связь с отправлением
    shipment_id INT NOT NULL,
    shipment_type VARCHAR(50) NOT NULL, -- 'bex', 'postexpress', 'marketplace_order'
    tracking_number VARCHAR(100),
    
    -- Информация о проблеме
    problem_type VARCHAR(50) NOT NULL, -- 'delayed', 'lost', 'damaged', 'return_requested', 'wrong_address', 'complaint'
    severity VARCHAR(20) DEFAULT 'medium', -- 'low', 'medium', 'high', 'critical'
    description TEXT,
    
    -- Статус решения
    status VARCHAR(50) DEFAULT 'open', -- 'open', 'investigating', 'waiting_response', 'resolved', 'closed'
    assigned_to INT REFERENCES users(id),
    resolution TEXT,
    
    -- Связанные данные
    order_id INT,
    user_id INT REFERENCES users(id),
    
    -- Временные метки
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    resolved_at TIMESTAMP,
    
    -- Дополнительные данные в JSON
    metadata JSONB DEFAULT '{}'
);

-- Индексы для проблемных отправлений
CREATE INDEX idx_problem_shipments_status ON problem_shipments(status);
CREATE INDEX idx_problem_shipments_type ON problem_shipments(problem_type);
CREATE INDEX idx_problem_shipments_shipment ON problem_shipments(shipment_id, shipment_type);
CREATE INDEX idx_problem_shipments_created ON problem_shipments(created_at DESC);

-- Таблица для логирования действий администраторов
CREATE TABLE IF NOT EXISTS logistics_admin_logs (
    id SERIAL PRIMARY KEY,
    
    -- Кто выполнил действие
    admin_id INT NOT NULL REFERENCES users(id),
    admin_email VARCHAR(255),
    
    -- Над чем выполнено действие
    entity_type VARCHAR(50) NOT NULL, -- 'shipment', 'problem', 'metric', 'report'
    entity_id INT,
    
    -- Действие
    action VARCHAR(100) NOT NULL, -- 'view', 'update_status', 'assign', 'resolve', 'export', 'notify'
    
    -- Детали действия
    details JSONB DEFAULT '{}',
    
    -- IP и user agent для аудита
    ip_address INET,
    user_agent TEXT,
    
    created_at TIMESTAMP DEFAULT NOW()
);

-- Индексы для логов
CREATE INDEX idx_logistics_admin_logs_admin ON logistics_admin_logs(admin_id);
CREATE INDEX idx_logistics_admin_logs_entity ON logistics_admin_logs(entity_type, entity_id);
CREATE INDEX idx_logistics_admin_logs_created ON logistics_admin_logs(created_at DESC);

-- Таблица для хранения настроек мониторинга
CREATE TABLE IF NOT EXISTS logistics_monitoring_settings (
    id SERIAL PRIMARY KEY,
    
    -- Пороговые значения для алертов
    delay_threshold_hours INT DEFAULT 72, -- Считать задержкой если > 72 часов
    critical_delay_hours INT DEFAULT 168, -- Критическая задержка > 7 дней
    
    -- Настройки уведомлений
    notify_on_delays BOOLEAN DEFAULT true,
    notify_on_returns BOOLEAN DEFAULT true,
    notify_on_problems BOOLEAN DEFAULT true,
    notification_emails TEXT[], -- Массив email для уведомлений
    
    -- Настройки автоматизации
    auto_create_problems BOOLEAN DEFAULT true, -- Автоматически создавать проблемы для задержек
    auto_assign_problems BOOLEAN DEFAULT false, -- Автоматически назначать проблемы
    
    -- Настройки отчётов
    daily_report_enabled BOOLEAN DEFAULT false,
    weekly_report_enabled BOOLEAN DEFAULT true,
    report_recipients TEXT[],
    
    updated_at TIMESTAMP DEFAULT NOW(),
    updated_by INT REFERENCES users(id)
);

-- Вставляем настройки по умолчанию
INSERT INTO logistics_monitoring_settings (id) VALUES (1) ON CONFLICT DO NOTHING;

-- Таблица для кеширования данных dashboard
CREATE TABLE IF NOT EXISTS logistics_dashboard_cache (
    id SERIAL PRIMARY KEY,
    cache_key VARCHAR(100) UNIQUE NOT NULL,
    cache_data JSONB NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Индекс для очистки устаревшего кеша
CREATE INDEX idx_logistics_dashboard_cache_expires ON logistics_dashboard_cache(expires_at);

-- Функция для автоматического обновления updated_at
CREATE OR REPLACE FUNCTION update_logistics_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Триггеры для автоматического обновления updated_at
CREATE TRIGGER update_logistics_metrics_updated_at
    BEFORE UPDATE ON logistics_metrics
    FOR EACH ROW
    EXECUTE FUNCTION update_logistics_updated_at();

CREATE TRIGGER update_problem_shipments_updated_at
    BEFORE UPDATE ON problem_shipments
    FOR EACH ROW
    EXECUTE FUNCTION update_logistics_updated_at();

CREATE TRIGGER update_logistics_monitoring_settings_updated_at
    BEFORE UPDATE ON logistics_monitoring_settings
    FOR EACH ROW
    EXECUTE FUNCTION update_logistics_updated_at();

-- Комментарии к таблицам
COMMENT ON TABLE logistics_metrics IS 'Агрегированные метрики логистики по дням';
COMMENT ON TABLE problem_shipments IS 'Отслеживание проблемных отправлений';
COMMENT ON TABLE logistics_admin_logs IS 'Аудит действий администраторов в логистике';
COMMENT ON TABLE logistics_monitoring_settings IS 'Настройки системы мониторинга логистики';
COMMENT ON TABLE logistics_dashboard_cache IS 'Кеш данных для dashboard логистики';