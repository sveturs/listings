-- Миграция для добавления мониторинга и аналитики точности AI детекции категорий

-- ===================================================================
-- 1. ТАБЛИЦА ДЛЯ A/B ТЕСТИРОВАНИЯ АЛГОРИТМОВ
-- ===================================================================

-- Таблица уже существует, добавляем только недостающие колонки
ALTER TABLE category_detection_experiments
ADD COLUMN IF NOT EXISTS description TEXT,
ADD COLUMN IF NOT EXISTS accuracy_percent DECIMAL(5,2),
ADD COLUMN IF NOT EXISTS avg_confidence_score DECIMAL(3,2),
ADD COLUMN IF NOT EXISTS median_processing_time_ms INTEGER,
ADD COLUMN IF NOT EXISTS created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP;

-- Индекс для быстрого поиска активных экспериментов
CREATE INDEX IF NOT EXISTS idx_experiments_active ON category_detection_experiments(is_active, started_at DESC);

-- ===================================================================
-- 2. ТАБЛИЦА ДЛЯ ДЕТАЛЬНОЙ СТАТИСТИКИ
-- ===================================================================

CREATE TABLE IF NOT EXISTS category_detection_stats (
    id SERIAL PRIMARY KEY,
    date DATE NOT NULL,
    hour INTEGER NOT NULL CHECK (hour >= 0 AND hour < 24),
    algorithm_version VARCHAR(50) NOT NULL,
    total_detections INTEGER DEFAULT 0,
    confirmed_detections INTEGER DEFAULT 0,
    avg_confidence_score DECIMAL(3,2),
    median_processing_time_ms INTEGER,
    p95_processing_time_ms INTEGER,
    p99_processing_time_ms INTEGER,
    unique_users INTEGER DEFAULT 0,
    unique_categories INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(date, hour, algorithm_version)
);

-- Индексы для быстрой аналитики
CREATE INDEX idx_stats_date ON category_detection_stats(date DESC);
CREATE INDEX idx_stats_algorithm ON category_detection_stats(algorithm_version, date DESC);

-- ===================================================================
-- 3. VIEW ДЛЯ АНАЛИТИКИ ТОЧНОСТИ ПО ДНЯМ
-- ===================================================================

CREATE OR REPLACE VIEW category_detection_daily_accuracy AS
SELECT
    DATE(created_at) as date,
    COUNT(*) as total_detections,
    SUM(CASE WHEN user_confirmed THEN 1 ELSE 0 END) as confirmed_detections,
    ROUND(100.0 * SUM(CASE WHEN user_confirmed THEN 1 ELSE 0 END) / NULLIF(COUNT(*), 0), 2) as accuracy_percent,
    AVG(confidence_score) as avg_confidence,
    PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY processing_time_ms) as median_processing_time,
    COUNT(DISTINCT listing_id) as unique_listings
FROM category_detection_feedback
WHERE created_at > CURRENT_DATE - INTERVAL '30 days'
GROUP BY DATE(created_at)
ORDER BY date DESC;

-- ===================================================================
-- 4. VIEW ДЛЯ АНАЛИТИКИ ПО КАТЕГОРИЯМ
-- ===================================================================

CREATE OR REPLACE VIEW category_detection_by_category AS
SELECT
    c.id as category_id,
    c.name as category_name,
    c.slug as category_slug,
    COUNT(f.id) as total_detections,
    SUM(CASE WHEN f.user_confirmed AND f.detected_category_id = f.correct_category_id THEN 1 ELSE 0 END) as correct_detections,
    ROUND(100.0 *
        SUM(CASE WHEN f.user_confirmed AND f.detected_category_id = f.correct_category_id THEN 1 ELSE 0 END) /
        NULLIF(COUNT(f.id), 0), 2
    ) as accuracy_percent,
    AVG(f.confidence_score) as avg_confidence,
    COUNT(DISTINCT f.listing_id) as unique_listings
FROM marketplace_categories c
LEFT JOIN category_detection_feedback f ON c.id = f.detected_category_id
WHERE f.created_at > CURRENT_DATE - INTERVAL '30 days'
GROUP BY c.id, c.name, c.slug
HAVING COUNT(f.id) > 0
ORDER BY total_detections DESC;

-- ===================================================================
-- 5. VIEW ДЛЯ ТОПОВЫХ ОШИБОК
-- ===================================================================

CREATE OR REPLACE VIEW category_detection_top_errors AS
SELECT
    dc.name as detected_category,
    cc.name as correct_category,
    COUNT(*) as error_count,
    AVG(f.confidence_score) as avg_confidence_when_wrong,
    array_agg(DISTINCT f.keywords) as common_keywords
FROM category_detection_feedback f
JOIN marketplace_categories dc ON dc.id = f.detected_category_id
JOIN marketplace_categories cc ON cc.id = f.correct_category_id
WHERE f.user_confirmed = TRUE
  AND f.detected_category_id != f.correct_category_id
  AND f.created_at > CURRENT_DATE - INTERVAL '7 days'
GROUP BY dc.name, cc.name
ORDER BY error_count DESC
LIMIT 20;

-- ===================================================================
-- 6. VIEW ДЛЯ АНАЛИЗА AI МАППИНГОВ
-- ===================================================================

CREATE OR REPLACE VIEW category_ai_mapping_performance AS
SELECT
    m.ai_domain,
    m.product_type,
    c.name as category_name,
    m.weight,
    m.success_count,
    m.failure_count,
    ROUND(100.0 * m.success_count / NULLIF(m.success_count + m.failure_count, 0), 2) as success_rate,
    m.success_count + m.failure_count as total_uses,
    m.updated_at
FROM category_ai_mappings m
JOIN marketplace_categories c ON c.id = m.category_id
WHERE m.is_active = TRUE
ORDER BY total_uses DESC, success_rate DESC;

-- ===================================================================
-- 7. ФУНКЦИЯ ДЛЯ РАСЧЕТА ТОЧНОСТИ В РЕАЛЬНОМ ВРЕМЕНИ
-- ===================================================================

CREATE OR REPLACE FUNCTION get_realtime_accuracy(
    p_hours INTEGER DEFAULT 24
) RETURNS TABLE (
    total_detections BIGINT,
    confirmed_detections BIGINT,
    accuracy_percent NUMERIC,
    avg_confidence NUMERIC,
    median_processing_time NUMERIC,
    unique_categories BIGINT,
    top_algorithm TEXT
) AS $$
BEGIN
    RETURN QUERY
    SELECT
        COUNT(*) as total_detections,
        SUM(CASE WHEN user_confirmed THEN 1 ELSE 0 END) as confirmed_detections,
        ROUND(100.0 * SUM(CASE WHEN user_confirmed THEN 1 ELSE 0 END) / NULLIF(COUNT(*), 0), 2) as accuracy_percent,
        ROUND(AVG(confidence_score), 3) as avg_confidence,
        PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY processing_time_ms) as median_processing_time,
        COUNT(DISTINCT detected_category_id) as unique_categories,
        MODE() WITHIN GROUP (ORDER BY algorithm_version) as top_algorithm
    FROM category_detection_feedback
    WHERE created_at > NOW() - (p_hours || ' hours')::INTERVAL;
END;
$$ LANGUAGE plpgsql;

-- ===================================================================
-- 8. ФУНКЦИЯ ДЛЯ АВТОМАТИЧЕСКОГО ОТКЛЮЧЕНИЯ ПЛОХИХ АЛГОРИТМОВ
-- ===================================================================

CREATE OR REPLACE FUNCTION check_algorithm_performance() RETURNS VOID AS $$
DECLARE
    v_accuracy NUMERIC;
    v_threshold NUMERIC := 70.0; -- Минимальная точность 70%
BEGIN
    -- Получаем точность за последний час
    SELECT accuracy_percent INTO v_accuracy
    FROM get_realtime_accuracy(1);

    -- Если точность упала ниже порога, отключаем эксперименты
    IF v_accuracy < v_threshold THEN
        UPDATE category_detection_experiments
        SET is_active = FALSE,
            ended_at = NOW()
        WHERE is_active = TRUE;

        -- Логируем событие
        INSERT INTO system_alerts (
            alert_type, severity, message, metadata, created_at
        ) VALUES (
            'ai_accuracy_drop', 'critical',
            'AI category detection accuracy dropped below ' || v_threshold || '%',
            jsonb_build_object('accuracy', v_accuracy, 'threshold', v_threshold),
            NOW()
        );
    END IF;
END;
$$ LANGUAGE plpgsql;

-- ===================================================================
-- 9. ТРИГГЕР ДЛЯ ОБНОВЛЕНИЯ СТАТИСТИКИ
-- ===================================================================

CREATE OR REPLACE FUNCTION update_detection_stats() RETURNS TRIGGER AS $$
BEGIN
    -- Обновляем статистику в реальном времени
    INSERT INTO category_detection_stats (
        date, hour, algorithm_version, total_detections,
        avg_confidence_score, median_processing_time_ms
    ) VALUES (
        DATE(NEW.created_at),
        EXTRACT(HOUR FROM NEW.created_at),
        NEW.algorithm_version,
        1,
        NEW.confidence_score,
        NEW.processing_time_ms
    )
    ON CONFLICT (date, hour, algorithm_version)
    DO UPDATE SET
        total_detections = category_detection_stats.total_detections + 1,
        avg_confidence_score =
            (category_detection_stats.avg_confidence_score * category_detection_stats.total_detections + NEW.confidence_score) /
            (category_detection_stats.total_detections + 1);

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_detection_stats
AFTER INSERT ON category_detection_feedback
FOR EACH ROW
EXECUTE FUNCTION update_detection_stats();

-- ===================================================================
-- 10. ИНДЕКСЫ ДЛЯ ОПТИМИЗАЦИИ ЗАПРОСОВ
-- ===================================================================

-- Индексы для feedback таблицы (если еще нет)
CREATE INDEX IF NOT EXISTS idx_feedback_created ON category_detection_feedback(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_feedback_confirmed ON category_detection_feedback(user_confirmed, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_feedback_category ON category_detection_feedback(detected_category_id, correct_category_id);
CREATE INDEX IF NOT EXISTS idx_feedback_algorithm ON category_detection_feedback(algorithm_version, created_at DESC);

-- ===================================================================
-- 11. НАЧАЛЬНЫЕ ДАННЫЕ ДЛЯ ЭКСПЕРИМЕНТОВ
-- ===================================================================

-- Добавляем эксперименты с учетом существующей структуры
INSERT INTO category_detection_experiments (
    experiment_name, algorithm_a, algorithm_b, traffic_split, is_active, description
) VALUES (
    'Baseline_vs_ML',
    'stable_v1',
    'experimental_v2',
    0.10,
    TRUE,
    'Сравнение стабильной версии с экспериментальным ML алгоритмом'
) ON CONFLICT (experiment_name) DO NOTHING;

-- ===================================================================
-- 12. РАСПИСАНИЕ ДЛЯ АВТОМАТИЧЕСКОЙ ПРОВЕРКИ (через cron или pg_cron)
-- ===================================================================

-- Если установлен pg_cron, можно добавить:
-- SELECT cron.schedule('check-ai-accuracy', '*/30 * * * *', 'SELECT check_algorithm_performance();');