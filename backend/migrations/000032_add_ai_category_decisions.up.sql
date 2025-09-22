-- Создаём таблицу для хранения решений AI по категоризации товаров
-- Это позволит кешировать решения и обучаться на них
CREATE TABLE IF NOT EXISTS ai_category_decisions (
    id SERIAL PRIMARY KEY,

    -- Хеш заголовка для быстрого поиска
    title_hash VARCHAR(64) NOT NULL,

    -- Исходные данные товара
    title TEXT NOT NULL,
    description TEXT,

    -- Результат категоризации
    category_id INTEGER NOT NULL REFERENCES marketplace_categories(id) ON DELETE CASCADE,
    confidence DECIMAL(5,4) NOT NULL CHECK (confidence >= 0 AND confidence <= 1),
    reasoning TEXT,

    -- Альтернативные категории (если были предложены AI)
    alternative_category_ids INTEGER[],

    -- Метаданные
    ai_model VARCHAR(100) DEFAULT 'claude-3-haiku-20240307',
    processing_time_ms INTEGER,

    -- AI hints (если были)
    ai_domain VARCHAR(100),
    ai_product_type VARCHAR(100),
    ai_keywords TEXT[],

    -- Обратная связь от пользователя
    user_confirmed BOOLEAN DEFAULT FALSE,
    user_corrected_category_id INTEGER REFERENCES marketplace_categories(id) ON DELETE SET NULL,
    user_feedback_at TIMESTAMP,

    -- Временные метки
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Индексы для быстрого поиска
CREATE INDEX idx_ai_decisions_title_hash ON ai_category_decisions(title_hash);
CREATE INDEX idx_ai_decisions_category_confidence ON ai_category_decisions(category_id, confidence DESC);
CREATE INDEX idx_ai_decisions_created_at ON ai_category_decisions(created_at DESC);
CREATE INDEX idx_ai_decisions_user_confirmed ON ai_category_decisions(user_confirmed) WHERE user_confirmed = TRUE;

-- Индекс для поиска по домену и типу продукта
CREATE INDEX idx_ai_decisions_domain_type ON ai_category_decisions(ai_domain, ai_product_type);

-- Триггер для обновления updated_at
CREATE OR REPLACE FUNCTION update_ai_category_decisions_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_ai_category_decisions_updated_at
    BEFORE UPDATE ON ai_category_decisions
    FOR EACH ROW
    EXECUTE FUNCTION update_ai_category_decisions_updated_at();

-- Комментарии к таблице и колонкам
COMMENT ON TABLE ai_category_decisions IS 'Хранение решений AI по категоризации товаров для кеширования и обучения';
COMMENT ON COLUMN ai_category_decisions.title_hash IS 'SHA256 хеш заголовка для быстрого поиска дубликатов';
COMMENT ON COLUMN ai_category_decisions.confidence IS 'Уверенность AI в выборе категории (0-1)';
COMMENT ON COLUMN ai_category_decisions.reasoning IS 'Объяснение AI почему выбрана эта категория';
COMMENT ON COLUMN ai_category_decisions.alternative_category_ids IS 'Альтернативные категории с confidence > 0.7 от основной';
COMMENT ON COLUMN ai_category_decisions.user_confirmed IS 'Пользователь подтвердил правильность категории';
COMMENT ON COLUMN ai_category_decisions.user_corrected_category_id IS 'Если пользователь исправил категорию - какая правильная';