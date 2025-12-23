-- Category Detections tracking table
-- Для отслеживания результатов детекции и обучения алгоритма

CREATE TABLE IF NOT EXISTS category_detections (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Входные данные
    input_title TEXT NOT NULL,
    input_description TEXT,
    input_language VARCHAR(5) NOT NULL DEFAULT 'sr',

    -- Результат детекции
    detected_category_id UUID REFERENCES categories(id),
    confidence_score DECIMAL(4,3),  -- 0.000 - 1.000
    detection_method VARCHAR(50) NOT NULL,  -- ai_claude, keyword_match, similarity, fallback
    matched_keywords TEXT[],

    -- Альтернативы (JSONB для гибкости)
    alternatives JSONB DEFAULT '[]'::jsonb,

    -- Подтверждение пользователя
    user_confirmed BOOLEAN,
    user_selected_category_id UUID REFERENCES categories(id),

    -- Метаданные
    processing_time_ms INTEGER,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Индексы для аналитики
CREATE INDEX idx_category_detections_created_at ON category_detections(created_at DESC);
CREATE INDEX idx_category_detections_method ON category_detections(detection_method);
CREATE INDEX idx_category_detections_detected_category ON category_detections(detected_category_id);
CREATE INDEX idx_category_detections_user_confirmed ON category_detections(user_confirmed) WHERE user_confirmed IS NOT NULL;

-- Индекс для FTS поиска по categories.meta_keywords (если ещё не существует)
CREATE INDEX IF NOT EXISTS idx_categories_meta_keywords_sr_fts
ON categories USING GIN (to_tsvector('simple', COALESCE(meta_keywords->>'sr', '')));

CREATE INDEX IF NOT EXISTS idx_categories_meta_keywords_en_fts
ON categories USING GIN (to_tsvector('simple', COALESCE(meta_keywords->>'en', '')));

CREATE INDEX IF NOT EXISTS idx_categories_meta_keywords_ru_fts
ON categories USING GIN (to_tsvector('simple', COALESCE(meta_keywords->>'ru', '')));

-- Индекс для similarity поиска по названиям (trigram)
CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE INDEX IF NOT EXISTS idx_categories_name_sr_trgm
ON categories USING GIN ((name->>'sr') gin_trgm_ops);

CREATE INDEX IF NOT EXISTS idx_categories_name_en_trgm
ON categories USING GIN ((name->>'en') gin_trgm_ops);

COMMENT ON TABLE category_detections IS 'Tracking детекций категорий для анализа и улучшения алгоритма';
