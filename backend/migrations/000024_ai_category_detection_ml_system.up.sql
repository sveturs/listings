-- Таблица для хранения обратной связи и обучения AI модели
CREATE TABLE IF NOT EXISTS category_detection_feedback (
    id BIGSERIAL PRIMARY KEY,
    listing_id INTEGER REFERENCES marketplace_listings(id) ON DELETE SET NULL,
    detected_category_id INTEGER REFERENCES marketplace_categories(id),
    correct_category_id INTEGER REFERENCES marketplace_categories(id),
    ai_hints JSONB NOT NULL DEFAULT '{}',
    keywords TEXT[] DEFAULT '{}',
    confidence_score DECIMAL(5,4) DEFAULT 0.0000,
    user_confirmed BOOLEAN DEFAULT FALSE,
    algorithm_version VARCHAR(50) DEFAULT 'stable_v1',
    processing_time_ms INTEGER,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Таблица для хранения паттернов AI domain -> категория
CREATE TABLE IF NOT EXISTS category_ai_mappings (
    id SERIAL PRIMARY KEY,
    ai_domain VARCHAR(100) NOT NULL,
    product_type VARCHAR(100) NOT NULL,
    category_id INTEGER REFERENCES marketplace_categories(id) ON DELETE CASCADE,
    weight DECIMAL(3,2) DEFAULT 1.00,
    success_count INTEGER DEFAULT 0,
    failure_count INTEGER DEFAULT 0,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(ai_domain, product_type, category_id)
);

-- Таблица для хранения весов ключевых слов
CREATE TABLE IF NOT EXISTS category_keyword_weights (
    id BIGSERIAL PRIMARY KEY,
    keyword VARCHAR(255) NOT NULL,
    category_id INTEGER REFERENCES marketplace_categories(id) ON DELETE CASCADE,
    weight DECIMAL(5,4) DEFAULT 1.0000,
    occurrence_count INTEGER DEFAULT 1,
    success_rate DECIMAL(5,4) DEFAULT 0.5000,
    language VARCHAR(10) DEFAULT 'ru',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(keyword, category_id, language)
);

-- Таблица для A/B тестирования алгоритмов
CREATE TABLE IF NOT EXISTS category_detection_experiments (
    id SERIAL PRIMARY KEY,
    experiment_name VARCHAR(100) NOT NULL UNIQUE,
    algorithm_a VARCHAR(50) NOT NULL,
    algorithm_b VARCHAR(50) NOT NULL,
    traffic_split DECIMAL(3,2) DEFAULT 0.10,
    total_requests_a INTEGER DEFAULT 0,
    successful_a INTEGER DEFAULT 0,
    total_requests_b INTEGER DEFAULT 0,
    successful_b INTEGER DEFAULT 0,
    is_active BOOLEAN DEFAULT TRUE,
    started_at TIMESTAMP DEFAULT NOW(),
    ended_at TIMESTAMP
);

-- Таблица статистики для мониторинга точности
CREATE TABLE IF NOT EXISTS category_detection_stats (
    id BIGSERIAL PRIMARY KEY,
    date DATE NOT NULL,
    hour INTEGER NOT NULL CHECK (hour >= 0 AND hour < 24),
    algorithm_version VARCHAR(50) NOT NULL,
    total_detections INTEGER DEFAULT 0,
    successful_detections INTEGER DEFAULT 0,
    failed_detections INTEGER DEFAULT 0,
    accuracy_percent DECIMAL(5,2) DEFAULT 0.00,
    avg_confidence_score DECIMAL(5,4) DEFAULT 0.0000,
    median_processing_time_ms INTEGER,
    p95_processing_time_ms INTEGER,
    p99_processing_time_ms INTEGER,
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(date, hour, algorithm_version)
);

-- Таблица для кэширования результатов определения категорий
CREATE TABLE IF NOT EXISTS category_detection_cache (
    id BIGSERIAL PRIMARY KEY,
    cache_key VARCHAR(255) NOT NULL UNIQUE,
    category_id INTEGER REFERENCES marketplace_categories(id),
    confidence_score DECIMAL(5,4),
    ai_hints JSONB,
    keywords TEXT[],
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Индексы для производительности
CREATE INDEX idx_feedback_keywords ON category_detection_feedback USING GIN(keywords);
CREATE INDEX idx_feedback_ai_hints ON category_detection_feedback USING GIN(ai_hints);
CREATE INDEX idx_feedback_created_at ON category_detection_feedback(created_at DESC);
CREATE INDEX idx_feedback_user_confirmed ON category_detection_feedback(user_confirmed) WHERE user_confirmed = TRUE;

CREATE INDEX idx_ai_mappings_domain_type ON category_ai_mappings(ai_domain, product_type);
CREATE INDEX idx_ai_mappings_category ON category_ai_mappings(category_id);
CREATE INDEX idx_ai_mappings_active ON category_ai_mappings(is_active) WHERE is_active = TRUE;

CREATE INDEX idx_keyword_weights_keyword ON category_keyword_weights(keyword);
CREATE INDEX idx_keyword_weights_category ON category_keyword_weights(category_id);
CREATE INDEX idx_keyword_weights_weight ON category_keyword_weights(weight DESC);

CREATE INDEX idx_detection_stats_date ON category_detection_stats(date DESC, hour DESC);
CREATE INDEX idx_detection_stats_algorithm ON category_detection_stats(algorithm_version);

CREATE INDEX idx_detection_cache_key ON category_detection_cache(cache_key);
CREATE INDEX idx_detection_cache_expires ON category_detection_cache(expires_at);

-- View для статистики точности
CREATE OR REPLACE VIEW category_detection_accuracy AS
SELECT
    DATE(created_at) as date,
    algorithm_version,
    COUNT(*) as total_detections,
    SUM(CASE WHEN user_confirmed THEN 1 ELSE 0 END) as confirmed,
    ROUND(100.0 * SUM(CASE WHEN user_confirmed THEN 1 ELSE 0 END) / NULLIF(COUNT(*), 0), 2) as accuracy_percent,
    AVG(confidence_score) as avg_confidence,
    PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY processing_time_ms) as median_processing_ms
FROM category_detection_feedback
WHERE created_at > NOW() - INTERVAL '30 days'
GROUP BY DATE(created_at), algorithm_version
ORDER BY date DESC, algorithm_version;

-- View для топ ошибочных категорий
CREATE OR REPLACE VIEW category_detection_errors AS
SELECT
    dc.name as detected_category,
    cc.name as correct_category,
    COUNT(*) as error_count,
    ARRAY_AGG(DISTINCT jsonb_extract_path_text(f.ai_hints, 'domain')) as ai_domains,
    ARRAY_AGG(DISTINCT jsonb_extract_path_text(f.ai_hints, 'productType')) as product_types
FROM category_detection_feedback f
LEFT JOIN marketplace_categories dc ON dc.id = f.detected_category_id
LEFT JOIN marketplace_categories cc ON cc.id = f.correct_category_id
WHERE f.user_confirmed = FALSE
  AND f.detected_category_id != f.correct_category_id
  AND f.created_at > NOW() - INTERVAL '7 days'
GROUP BY dc.name, cc.name
ORDER BY error_count DESC
LIMIT 20;

-- Функция для автоматического обновления updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Триггеры для обновления updated_at
CREATE TRIGGER update_category_detection_feedback_updated_at
    BEFORE UPDATE ON category_detection_feedback
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_category_ai_mappings_updated_at
    BEFORE UPDATE ON category_ai_mappings
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_category_keyword_weights_updated_at
    BEFORE UPDATE ON category_keyword_weights
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Начальные данные для AI domain маппинга (используем только существующие категории)
INSERT INTO category_ai_mappings (ai_domain, product_type, category_id, weight) VALUES
-- Electronics (категория 1001)
('electronics', 'laptop', 1102, 1.00),
('electronics', 'smartphone', 1101, 1.00),
('electronics', 'tablet', 1101, 0.95),
('electronics', 'router', 1103, 1.00),
('electronics', 'monitor', 1102, 0.90),
('electronics', 'keyboard', 1102, 0.85),
('electronics', 'mouse', 1102, 0.85),
('electronics', 'headphones', 1103, 1.00),
('electronics', 'smartwatch', 1101, 0.90),
('electronics', 'camera', 1103, 1.00),
('electronics', 'gaming', 1105, 1.00),

-- Fashion (категория 1002)
('fashion', 'clothing', 1202, 1.00),
('fashion', 'shoes', 1002, 1.00),
('fashion', 'bag', 1002, 1.00),
('fashion', 'accessories', 1002, 1.00),
('fashion', 'watch', 1002, 0.95),
('fashion', 'jewelry', 1002, 0.95),

-- Automotive (категория 1003)
('automotive', 'car', 1301, 1.00),
('automotive', 'motorcycle', 1302, 1.00),
('automotive', 'parts', 1303, 1.00),
('automotive', 'tires', 1303, 0.95),
('automotive', 'tools', 1303, 0.85),

-- Real Estate (категория 1004)
('real-estate', 'apartment', 1401, 1.00),
('real-estate', 'house', 1402, 1.00),
('real-estate', 'land', 1403, 1.00),
('real-estate', 'commercial', 1404, 1.00),
('real-estate', 'room', 1401, 0.85),

-- Home & Garden (категория 1005)
('home-garden', 'furniture', 1501, 1.00),
('home-garden', 'appliance', 1005, 1.00),
('home-garden', 'decoration', 1005, 1.00),
('home-garden', 'garden-tools', 1005, 1.00),
('home-garden', 'plants', 1005, 0.95),

-- Agriculture (категория 1006)
('agriculture', 'machinery', 1601, 1.00),
('agriculture', 'seeds', 1602, 1.00),
('agriculture', 'livestock', 1603, 1.00),
('agriculture', 'produce', 1604, 1.00),

-- Entertainment (категория 1015)
('entertainment', 'puzzle', 1015, 1.00),
('entertainment', 'game', 1015, 1.00),
('entertainment', 'toy', 1013, 1.00),
('entertainment', 'book', 1012, 1.00),
('entertainment', 'board-game', 1015, 0.95),
('entertainment', 'collectible', 1017, 0.90),
('entertainment', 'music', 1016, 1.00),

-- Sports & Recreation (категория 1010)
('sports-recreation', 'bicycle', 1010, 1.00),
('sports-recreation', 'fitness', 1010, 1.00),
('sports-recreation', 'outdoor', 1010, 1.00),
('sports-recreation', 'camping', 1010, 0.95),
('sports-recreation', 'sports-equipment', 1010, 1.00),

-- Services (категория 1009)
('services', 'repair', 1009, 1.00),
('services', 'delivery', 1009, 1.00),
('services', 'cleaning', 1009, 1.00),
('services', 'education', 1019, 1.00),
('services', 'consulting', 1009, 1.00),

-- Industrial/Construction (категория 1007)
('construction', 'materials', 1504, 1.00),
('construction', 'tools', 1007, 1.00),
('construction', 'sand', 1504, 0.95),
('construction', 'cement', 1504, 0.95),
('construction', 'bricks', 1504, 0.95),
('industrial', 'equipment', 1007, 1.00),
('industrial', 'machinery', 1007, 1.00),

-- Food & Beverages (категория 1008)
('food-beverages', 'food', 1008, 1.00),
('food-beverages', 'beverages', 1008, 1.00),
('food-beverages', 'organic', 1008, 0.95),

-- Nature/Crafts
('nature', 'natural-materials', 1005, 1.00),
('nature', 'plants', 1005, 1.00),
('nature', 'seeds', 1602, 0.95),
('nature', 'wood', 1504, 0.90),

-- Antiques & Art (категория 1017)
('antiques', 'vintage', 1017, 1.00),
('antiques', 'collectible', 1017, 1.00),
('antiques', 'art', 1017, 0.95),
('antiques', 'coins', 1017, 0.90),

-- Pets (категория 1011)
('pets', 'pet-supplies', 1011, 1.00),
('pets', 'pet-food', 1011, 0.95),

-- Kids & Baby (категория 1013)
('kids', 'toys', 1013, 1.00),
('kids', 'clothing', 1013, 0.95),
('kids', 'baby-items', 1013, 1.00),

-- Health & Beauty (категория 1014)
('health-beauty', 'cosmetics', 1014, 1.00),
('health-beauty', 'health-products', 1014, 1.00),

-- Other/General (fallback на Electronics как наиболее общую)
('other', 'miscellaneous', 1001, 0.50),
('other', 'unknown', 1001, 0.25)
ON CONFLICT (ai_domain, product_type, category_id) DO NOTHING;

-- Функция для получения рекомендованной категории по AI domain
CREATE OR REPLACE FUNCTION get_category_by_ai_hints(
    p_domain VARCHAR,
    p_product_type VARCHAR
) RETURNS TABLE(
    category_id INTEGER,
    confidence DECIMAL(5,4)
) AS $$
BEGIN
    RETURN QUERY
    SELECT
        m.category_id,
        m.weight * (1.0 + (m.success_count::DECIMAL / GREATEST(m.success_count + m.failure_count, 1)) * 0.2) as confidence
    FROM category_ai_mappings m
    WHERE m.ai_domain = p_domain
      AND m.product_type = p_product_type
      AND m.is_active = TRUE
    ORDER BY confidence DESC
    LIMIT 1;
END;
$$ LANGUAGE plpgsql;

-- Функция для обновления статистики успешности маппинга
CREATE OR REPLACE FUNCTION update_mapping_stats(
    p_domain VARCHAR,
    p_product_type VARCHAR,
    p_category_id INTEGER,
    p_success BOOLEAN
) RETURNS VOID AS $$
BEGIN
    IF p_success THEN
        UPDATE category_ai_mappings
        SET success_count = success_count + 1,
            weight = LEAST(weight * 1.01, 1.0)
        WHERE ai_domain = p_domain
          AND product_type = p_product_type
          AND category_id = p_category_id;
    ELSE
        UPDATE category_ai_mappings
        SET failure_count = failure_count + 1,
            weight = GREATEST(weight * 0.99, 0.1)
        WHERE ai_domain = p_domain
          AND product_type = p_product_type
          AND category_id = p_category_id;
    END IF;
END;
$$ LANGUAGE plpgsql;