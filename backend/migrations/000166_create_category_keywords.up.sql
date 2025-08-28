-- Таблица для хранения ключевых слов категорий
CREATE TABLE category_keywords (
    id SERIAL PRIMARY KEY,
    category_id INTEGER NOT NULL REFERENCES marketplace_categories(id) ON DELETE CASCADE,
    keyword VARCHAR(100) NOT NULL,
    language VARCHAR(2) DEFAULT 'en',
    weight FLOAT DEFAULT 1.0 CHECK (weight >= 0.0 AND weight <= 10.0),
    
    -- Тип ключевого слова
    keyword_type VARCHAR(20) DEFAULT 'general' CHECK (keyword_type IN ('main', 'synonym', 'brand', 'attribute', 'context', 'pattern')),
    
    -- Исключающее слово (если true, то наличие этого слова исключает категорию)
    is_negative BOOLEAN DEFAULT FALSE,
    
    -- Источник добавления
    source VARCHAR(50) DEFAULT 'manual' CHECK (source IN ('manual', 'ai_extracted', 'user_confirmed', 'auto_learned')),
    
    -- Статистика использования
    usage_count INTEGER DEFAULT 0,
    success_rate FLOAT DEFAULT 0.0 CHECK (success_rate >= 0.0 AND success_rate <= 1.0),
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Индексы для быстрого поиска
CREATE INDEX idx_keyword_lower ON category_keywords(LOWER(keyword));
CREATE INDEX idx_category_weight ON category_keywords(category_id, weight DESC);
CREATE INDEX idx_language ON category_keywords(language);
CREATE INDEX idx_keyword_type ON category_keywords(keyword_type);
CREATE UNIQUE INDEX idx_unique_keyword_category ON category_keywords(category_id, keyword, language);

-- Триггер для обновления updated_at
CREATE TRIGGER update_category_keywords_updated_at 
    BEFORE UPDATE ON category_keywords 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- Комментарии к таблице
COMMENT ON TABLE category_keywords IS 'Ключевые слова для семантического поиска категорий';
COMMENT ON COLUMN category_keywords.keyword IS 'Ключевое слово или фраза';
COMMENT ON COLUMN category_keywords.weight IS 'Вес ключевого слова (0.0 - 10.0)';
COMMENT ON COLUMN category_keywords.keyword_type IS 'Тип: main - основное, synonym - синоним, brand - бренд, attribute - атрибут, context - контекст, pattern - паттерн';
COMMENT ON COLUMN category_keywords.is_negative IS 'Если true, наличие этого слова исключает категорию';
COMMENT ON COLUMN category_keywords.usage_count IS 'Сколько раз использовалось при поиске';
COMMENT ON COLUMN category_keywords.success_rate IS 'Процент успешных определений категории';