-- Таблица весов полей поиска
CREATE TABLE IF NOT EXISTS search_weights (
    id BIGSERIAL PRIMARY KEY,
    field_name VARCHAR(100) NOT NULL UNIQUE,
    weight DECIMAL(3,2) NOT NULL DEFAULT 1.0 CHECK (weight > 0 AND weight <= 10),
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Таблица синонимов
CREATE TABLE IF NOT EXISTS search_synonyms_config (
    id BIGSERIAL PRIMARY KEY,
    term VARCHAR(255) NOT NULL,
    synonyms TEXT[] NOT NULL,
    language VARCHAR(10) NOT NULL DEFAULT 'ru',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(term, language)
);

-- Таблица правил транслитерации
CREATE TABLE IF NOT EXISTS transliteration_rules (
    id BIGSERIAL PRIMARY KEY,
    from_script VARCHAR(50) NOT NULL,
    to_script VARCHAR(50) NOT NULL,
    from_pattern VARCHAR(100) NOT NULL,
    to_pattern VARCHAR(100) NOT NULL,
    priority INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(from_script, to_script, from_pattern)
);

-- Таблица статистики поиска
CREATE TABLE IF NOT EXISTS search_statistics (
    id BIGSERIAL PRIMARY KEY,
    query TEXT NOT NULL,
    results_count INTEGER NOT NULL DEFAULT 0,
    search_duration_ms BIGINT NOT NULL,
    user_id BIGINT REFERENCES users(id) ON DELETE SET NULL,
    search_filters JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Таблица общей конфигурации поиска
CREATE TABLE IF NOT EXISTS search_config (
    id BIGSERIAL PRIMARY KEY,
    min_search_length INTEGER NOT NULL DEFAULT 2,
    max_suggestions INTEGER NOT NULL DEFAULT 10,
    fuzzy_enabled BOOLEAN NOT NULL DEFAULT true,
    fuzzy_max_edits INTEGER NOT NULL DEFAULT 2,
    synonyms_enabled BOOLEAN NOT NULL DEFAULT true,
    transliteration_enabled BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Индексы для оптимизации
CREATE INDEX IF NOT EXISTS idx_search_statistics_query ON search_statistics(query);
CREATE INDEX IF NOT EXISTS idx_search_statistics_created_at ON search_statistics(created_at);
CREATE INDEX IF NOT EXISTS idx_search_statistics_user_id ON search_statistics(user_id);
CREATE INDEX IF NOT EXISTS idx_search_synonyms_config_term ON search_synonyms_config(term);
CREATE INDEX IF NOT EXISTS idx_transliteration_rules_scripts ON transliteration_rules(from_script, to_script);

-- Вставка начальных данных
INSERT INTO search_config (min_search_length, max_suggestions, fuzzy_enabled, fuzzy_max_edits, synonyms_enabled, transliteration_enabled)
VALUES (2, 10, true, 2, true, true);

-- Веса по умолчанию
INSERT INTO search_weights (field_name, weight, description) VALUES
('title', 3.0, 'Заголовок объявления'),
('description', 1.5, 'Описание объявления'),
('category', 2.0, 'Категория'),
('location', 1.2, 'Местоположение'),
('tags', 1.0, 'Теги');

-- Синонимы по умолчанию
INSERT INTO search_synonyms_config (term, synonyms, language) VALUES
('квартира', ARRAY['апартаменты', 'жилье', 'комната'], 'ru'),
('машина', ARRAY['автомобиль', 'авто', 'тачка'], 'ru'),
('телефон', ARRAY['смартфон', 'мобильный', 'сотовый'], 'ru');

-- Правила транслитерации кириллица -> латиница
INSERT INTO transliteration_rules (from_script, to_script, from_pattern, to_pattern, priority) VALUES
('cyrillic', 'latin', 'а', 'a', 1),
('cyrillic', 'latin', 'б', 'b', 1),
('cyrillic', 'latin', 'в', 'v', 1),
('cyrillic', 'latin', 'г', 'g', 1),
('cyrillic', 'latin', 'д', 'd', 1),
('cyrillic', 'latin', 'е', 'e', 1),
('cyrillic', 'latin', 'ё', 'yo', 1),
('cyrillic', 'latin', 'ж', 'zh', 1),
('cyrillic', 'latin', 'з', 'z', 1),
('cyrillic', 'latin', 'и', 'i', 1),
('cyrillic', 'latin', 'й', 'y', 1),
('cyrillic', 'latin', 'к', 'k', 1),
('cyrillic', 'latin', 'л', 'l', 1),
('cyrillic', 'latin', 'м', 'm', 1),
('cyrillic', 'latin', 'н', 'n', 1),
('cyrillic', 'latin', 'о', 'o', 1),
('cyrillic', 'latin', 'п', 'p', 1),
('cyrillic', 'latin', 'р', 'r', 1),
('cyrillic', 'latin', 'с', 's', 1),
('cyrillic', 'latin', 'т', 't', 1),
('cyrillic', 'latin', 'у', 'u', 1),
('cyrillic', 'latin', 'ф', 'f', 1),
('cyrillic', 'latin', 'х', 'h', 1),
('cyrillic', 'latin', 'ц', 'ts', 1),
('cyrillic', 'latin', 'ч', 'ch', 1),
('cyrillic', 'latin', 'ш', 'sh', 1),
('cyrillic', 'latin', 'щ', 'sch', 1),
('cyrillic', 'latin', 'ъ', '', 1),
('cyrillic', 'latin', 'ы', 'y', 1),
('cyrillic', 'latin', 'ь', '', 1),
('cyrillic', 'latin', 'э', 'e', 1),
('cyrillic', 'latin', 'ю', 'yu', 1),
('cyrillic', 'latin', 'я', 'ya', 1);