-- Таблица для статистики определения категорий
CREATE TABLE category_detection_stats (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    session_id VARCHAR(100),
    
    -- Метод определения
    method VARCHAR(50) NOT NULL CHECK (method IN ('similarity', 'keywords', 'combined', 'manual')),
    
    -- Результаты AI анализа
    ai_keywords TEXT[], -- массив ключевых слов из AI
    ai_attributes JSONB, -- атрибуты из AI
    ai_domain VARCHAR(100), -- домен из AI (automotive, electronics и т.д.)
    ai_product_type VARCHAR(100), -- тип продукта из AI
    
    -- Категории
    ai_suggested_category_id INTEGER REFERENCES marketplace_categories(id),
    final_category_id INTEGER REFERENCES marketplace_categories(id),
    alternative_categories JSONB, -- альтернативные категории с их scores
    
    -- Метрики
    confidence_score FLOAT CHECK (confidence_score >= 0.0 AND confidence_score <= 1.0),
    similarity_score FLOAT, -- если использовался similarity search
    keyword_score FLOAT, -- если использовались keywords
    
    -- Детали similarity search
    similar_listings_found INTEGER DEFAULT 0,
    top_similar_listing_id INTEGER REFERENCES marketplace_listings(id),
    top_similarity_score FLOAT,
    
    -- Детали keyword matching
    matched_keywords TEXT[],
    matched_negative_keywords TEXT[],
    
    -- Производительность
    processing_time_ms INTEGER,
    
    -- Результат
    user_confirmed BOOLEAN, -- подтвердил ли пользователь выбор
    user_selected_category_id INTEGER REFERENCES marketplace_categories(id), -- если пользователь изменил
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Индексы для анализа
CREATE INDEX idx_detection_stats_created ON category_detection_stats(created_at DESC);
CREATE INDEX idx_detection_stats_method ON category_detection_stats(method);
CREATE INDEX idx_detection_stats_user ON category_detection_stats(user_id) WHERE user_id IS NOT NULL;
CREATE INDEX idx_detection_stats_categories ON category_detection_stats(ai_suggested_category_id, final_category_id);
CREATE INDEX idx_detection_stats_confidence ON category_detection_stats(confidence_score);

-- View для анализа эффективности методов
CREATE VIEW category_detection_effectiveness AS
SELECT 
    method,
    COUNT(*) as total_detections,
    AVG(CASE WHEN ai_suggested_category_id = final_category_id THEN 1.0 ELSE 0.0 END) as accuracy,
    AVG(confidence_score) as avg_confidence,
    AVG(processing_time_ms) as avg_time_ms,
    COUNT(CASE WHEN user_confirmed = true THEN 1 END) as user_confirmed_count,
    COUNT(CASE WHEN user_selected_category_id IS NOT NULL THEN 1 END) as user_corrected_count
FROM category_detection_stats
WHERE created_at > NOW() - INTERVAL '30 days'
GROUP BY method;

-- View для анализа проблемных категорий
CREATE VIEW problematic_categories AS
SELECT 
    c.id,
    c.name,
    c.slug,
    COUNT(DISTINCT s.id) as detection_attempts,
    AVG(CASE WHEN s.ai_suggested_category_id = s.final_category_id THEN 1.0 ELSE 0.0 END) as success_rate,
    AVG(s.confidence_score) as avg_confidence,
    COUNT(CASE WHEN s.user_selected_category_id IS NOT NULL THEN 1 END) as correction_count
FROM marketplace_categories c
LEFT JOIN category_detection_stats s ON c.id = s.ai_suggested_category_id
WHERE s.created_at > NOW() - INTERVAL '30 days'
GROUP BY c.id, c.name, c.slug
HAVING COUNT(DISTINCT s.id) > 5
ORDER BY AVG(CASE WHEN s.ai_suggested_category_id = s.final_category_id THEN 1.0 ELSE 0.0 END) ASC
LIMIT 20;

-- Комментарии
COMMENT ON TABLE category_detection_stats IS 'Статистика определения категорий для анализа и улучшения алгоритма';
COMMENT ON COLUMN category_detection_stats.method IS 'similarity - по похожим объявлениям, keywords - по ключевым словам, combined - комбинированный';
COMMENT ON VIEW category_detection_effectiveness IS 'Эффективность разных методов определения категорий';
COMMENT ON VIEW problematic_categories IS 'Категории с низким процентом правильного определения';