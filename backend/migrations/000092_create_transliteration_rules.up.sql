-- Создание таблицы для хранения правил транслитерации
CREATE TABLE IF NOT EXISTS transliteration_rules (
    id SERIAL PRIMARY KEY,
    
    -- Основные поля
    source_char VARCHAR(10) NOT NULL,  -- Исходный символ (может быть многосимвольным, например "дж")
    target_char VARCHAR(20) NOT NULL,  -- Целевой символ/строка (может быть пустой для удаления)
    language VARCHAR(2) NOT NULL,      -- Язык: 'ru', 'sr', 'bg' и т.д.
    
    -- Статус и приоритет
    enabled BOOLEAN DEFAULT TRUE,
    priority INTEGER DEFAULT 0,       -- Приоритет применения (больше = раньше)
    
    -- Метаданные
    description TEXT,                  -- Описание правила
    rule_type VARCHAR(20) DEFAULT 'custom', -- 'builtin', 'custom', 'generated'
    
    -- Временные метки
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    -- Уникальный индекс для пары язык-символ
    UNIQUE(language, source_char)
);

-- Создание индексов для быстрого поиска
CREATE INDEX IF NOT EXISTS idx_transliteration_rules_language ON transliteration_rules(language);
CREATE INDEX IF NOT EXISTS idx_transliteration_rules_enabled ON transliteration_rules(enabled);
CREATE INDEX IF NOT EXISTS idx_transliteration_rules_priority ON transliteration_rules(priority DESC);
CREATE INDEX IF NOT EXISTS idx_transliteration_rules_type ON transliteration_rules(rule_type);

-- Создание составного индекса для быстрого поиска активных правил по языку
CREATE INDEX IF NOT EXISTS idx_transliteration_rules_active ON transliteration_rules(language, enabled, priority DESC);

-- Триггер для автоматического обновления updated_at
CREATE OR REPLACE FUNCTION update_transliteration_rules_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trigger_update_transliteration_rules_updated_at ON trigger_update_transliteration_rules_updated_at;
CREATE TRIGGER trigger_update_transliteration_rules_updated_at
    BEFORE UPDATE ON transliteration_rules
    FOR EACH ROW
    EXECUTE FUNCTION update_transliteration_rules_updated_at();

-- Вставка базовых правил для русского языка (встроенные)
INSERT INTO transliteration_rules (source_char, target_char, language, rule_type, description, priority) VALUES
    -- Основная кириллица
    ('а', 'a', 'ru', 'builtin', 'Русская буква А', 100),
    ('б', 'b', 'ru', 'builtin', 'Русская буква Б', 100),
    ('в', 'v', 'ru', 'builtin', 'Русская буква В', 100),
    ('г', 'g', 'ru', 'builtin', 'Русская буква Г', 100),
    ('д', 'd', 'ru', 'builtin', 'Русская буква Д', 100),
    ('е', 'e', 'ru', 'builtin', 'Русская буква Е', 100),
    ('ё', 'yo', 'ru', 'builtin', 'Русская буква Ё', 100),
    ('ж', 'zh', 'ru', 'builtin', 'Русская буква Ж', 100),
    ('з', 'z', 'ru', 'builtin', 'Русская буква З', 100),
    ('и', 'i', 'ru', 'builtin', 'Русская буква И', 100),
    ('й', 'y', 'ru', 'builtin', 'Русская буква Й', 100),
    ('к', 'k', 'ru', 'builtin', 'Русская буква К', 100),
    ('л', 'l', 'ru', 'builtin', 'Русская буква Л', 100),
    ('м', 'm', 'ru', 'builtin', 'Русская буква М', 100),
    ('н', 'n', 'ru', 'builtin', 'Русская буква Н', 100),
    ('о', 'o', 'ru', 'builtin', 'Русская буква О', 100),
    ('п', 'p', 'ru', 'builtin', 'Русская буква П', 100),
    ('р', 'r', 'ru', 'builtin', 'Русская буква Р', 100),
    ('с', 's', 'ru', 'builtin', 'Русская буква С', 100),
    ('т', 't', 'ru', 'builtin', 'Русская буква Т', 100),
    ('у', 'u', 'ru', 'builtin', 'Русская буква У', 100),
    ('ф', 'f', 'ru', 'builtin', 'Русская буква Ф', 100),
    ('х', 'h', 'ru', 'builtin', 'Русская буква Х', 100),
    ('ц', 'ts', 'ru', 'builtin', 'Русская буква Ц', 100),
    ('ч', 'ch', 'ru', 'builtin', 'Русская буква Ч', 100),
    ('ш', 'sh', 'ru', 'builtin', 'Русская буква Ш', 100),
    ('щ', 'sch', 'ru', 'builtin', 'Русская буква Щ', 100),
    ('ъ', '', 'ru', 'builtin', 'Твердый знак (удаляется)', 100),
    ('ы', 'y', 'ru', 'builtin', 'Русская буква Ы', 100),
    ('ь', '', 'ru', 'builtin', 'Мягкий знак (удаляется)', 100),
    ('э', 'e', 'ru', 'builtin', 'Русская буква Э', 100),
    ('ю', 'yu', 'ru', 'builtin', 'Русская буква Ю', 100),
    ('я', 'ya', 'ru', 'builtin', 'Русская буква Я', 100);

-- Вставка базовых правил для сербского языка (встроенные)
INSERT INTO transliteration_rules (source_char, target_char, language, rule_type, description, priority) VALUES
    -- Специфические сербские символы (кириллица)
    ('ђ', 'đ', 'sr', 'builtin', 'Сербская буква Ђ', 200),
    ('ј', 'j', 'sr', 'builtin', 'Сербская буква Ј', 200),
    ('љ', 'lj', 'sr', 'builtin', 'Сербская буква Љ', 200),
    ('њ', 'nj', 'sr', 'builtin', 'Сербская буква Њ', 200),
    ('ћ', 'ć', 'sr', 'builtin', 'Сербская буква Ћ', 200),
    ('џ', 'dž', 'sr', 'builtin', 'Сербская буква Џ', 200),
    -- Сербские варианты стандартных букв
    ('ш', 'š', 'sr', 'builtin', 'Сербская буква Ш', 200),
    ('ж', 'ž', 'sr', 'builtin', 'Сербская буква Ж', 200),
    ('ч', 'č', 'sr', 'builtin', 'Сербская буква Ч', 200),
    -- Заглавные буквы
    ('Ђ', 'Đ', 'sr', 'builtin', 'Сербская заглавная буква Ђ', 200),
    ('Ј', 'J', 'sr', 'builtin', 'Сербская заглавная буква Ј', 200),
    ('Љ', 'LJ', 'sr', 'builtin', 'Сербская заглавная буква Љ', 200),
    ('Њ', 'NJ', 'sr', 'builtin', 'Сербская заглавная буква Њ', 200),
    ('Ћ', 'Ć', 'sr', 'builtin', 'Сербская заглавная буква Ћ', 200),
    ('Џ', 'DŽ', 'sr', 'builtin', 'Сербская заглавная буква Џ', 200),
    ('Ш', 'Š', 'sr', 'builtin', 'Сербская заглавная буква Ш', 200),
    ('Ж', 'Ž', 'sr', 'builtin', 'Сербская заглавная буква Ж', 200),
    ('Ч', 'Č', 'sr', 'builtin', 'Сербская заглавная буква Ч', 200);

-- Комментарий к таблице
COMMENT ON TABLE transliteration_rules IS 'Правила транслитерации для создания URL-friendly slug из кириллических названий';
COMMENT ON COLUMN transliteration_rules.source_char IS 'Исходный символ или группа символов для замены';
COMMENT ON COLUMN transliteration_rules.target_char IS 'Целевой символ или строка для замены (может быть пустой)';
COMMENT ON COLUMN transliteration_rules.language IS 'Двухбуквенный код языка (ISO 639-1)';
COMMENT ON COLUMN transliteration_rules.priority IS 'Приоритет применения правила (больше = раньше)';
COMMENT ON COLUMN transliteration_rules.rule_type IS 'Тип правила: builtin (встроенное), custom (пользовательское), generated (сгенерированное)';