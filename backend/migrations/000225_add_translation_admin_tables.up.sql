-- Migration: Add tables for translation admin panel
-- Description: Creates tables for version control, sync conflicts, audit log, and AI providers

-- 1. История версий переводов
CREATE TABLE IF NOT EXISTS translation_versions (
    id SERIAL PRIMARY KEY,
    translation_id INTEGER REFERENCES translations(id) ON DELETE CASCADE,
    version_number INTEGER NOT NULL,
    entity_type VARCHAR(50) NOT NULL,
    entity_id INTEGER NOT NULL,
    language VARCHAR(10) NOT NULL,
    field_name VARCHAR(50) NOT NULL,
    translated_text TEXT NOT NULL,
    changed_by INTEGER REFERENCES users(id) ON DELETE SET NULL,
    changed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    change_comment TEXT,
    metadata JSONB DEFAULT '{}',
    UNIQUE(translation_id, version_number)
);

-- 2. Конфликты синхронизации
CREATE TABLE IF NOT EXISTS translation_sync_conflicts (
    id SERIAL PRIMARY KEY,
    source_type VARCHAR(50) NOT NULL, -- 'frontend', 'database', 'opensearch'
    target_type VARCHAR(50) NOT NULL,
    entity_identifier TEXT NOT NULL,
    source_value TEXT,
    target_value TEXT,
    conflict_type VARCHAR(50), -- 'missing', 'different', 'outdated'
    resolved BOOLEAN DEFAULT FALSE,
    resolved_by INTEGER REFERENCES users(id) ON DELETE SET NULL,
    resolved_at TIMESTAMP,
    resolution_type VARCHAR(50), -- 'use_source', 'use_target', 'manual'
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 3. Аудит лог действий
CREATE TABLE IF NOT EXISTS translation_audit_log (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE SET NULL,
    action VARCHAR(100) NOT NULL,
    entity_type VARCHAR(50),
    entity_id INTEGER,
    old_value TEXT,
    new_value TEXT,
    ip_address INET,
    user_agent TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 4. Настройки AI провайдеров
CREATE TABLE IF NOT EXISTS translation_providers (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    provider_type VARCHAR(50) NOT NULL, -- 'openai', 'google', 'deepl'
    api_key TEXT, -- Will be encrypted at application level
    settings JSONB DEFAULT '{}',
    usage_limit INTEGER,
    usage_current INTEGER DEFAULT 0,
    is_active BOOLEAN DEFAULT TRUE,
    priority INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 5. Статистика качества переводов
CREATE TABLE IF NOT EXISTS translation_quality_metrics (
    id SERIAL PRIMARY KEY,
    translation_id INTEGER REFERENCES translations(id) ON DELETE CASCADE,
    quality_score DECIMAL(3,2), -- 0.00 to 1.00
    character_count INTEGER,
    word_count INTEGER,
    has_placeholders BOOLEAN DEFAULT FALSE,
    has_html_tags BOOLEAN DEFAULT FALSE,
    checked_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    checked_by VARCHAR(50), -- 'system', 'manual', 'ai'
    issues JSONB DEFAULT '[]'
);

-- 6. Задачи на перевод
CREATE TABLE IF NOT EXISTS translation_tasks (
    id SERIAL PRIMARY KEY,
    task_type VARCHAR(50) NOT NULL, -- 'single', 'batch', 'module'
    status VARCHAR(50) NOT NULL DEFAULT 'pending', -- 'pending', 'processing', 'completed', 'failed'
    source_language VARCHAR(10),
    target_languages TEXT[], -- массив языков для перевода
    entity_references JSONB DEFAULT '[]', -- список ссылок на сущности
    provider_id INTEGER REFERENCES translation_providers(id),
    created_by INTEGER REFERENCES users(id) ON DELETE SET NULL,
    assigned_to INTEGER REFERENCES users(id) ON DELETE SET NULL,
    started_at TIMESTAMP,
    completed_at TIMESTAMP,
    error_message TEXT,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Индексы для производительности
CREATE INDEX idx_translation_versions_lookup ON translation_versions(translation_id, version_number);
CREATE INDEX idx_translation_versions_entity ON translation_versions(entity_type, entity_id);
CREATE INDEX idx_sync_conflicts_unresolved ON translation_sync_conflicts(resolved) WHERE resolved = FALSE;
CREATE INDEX idx_sync_conflicts_dates ON translation_sync_conflicts(created_at DESC);
CREATE INDEX idx_audit_log_user_date ON translation_audit_log(user_id, created_at DESC);
CREATE INDEX idx_audit_log_entity ON translation_audit_log(entity_type, entity_id);
CREATE INDEX idx_providers_active ON translation_providers(is_active, priority) WHERE is_active = TRUE;
CREATE INDEX idx_quality_metrics_score ON translation_quality_metrics(quality_score);
CREATE INDEX idx_quality_metrics_issues ON translation_quality_metrics(has_placeholders, has_html_tags);
CREATE INDEX idx_translation_tasks_status ON translation_tasks(status, created_at) WHERE status != 'completed';
CREATE INDEX idx_translation_tasks_assigned ON translation_tasks(assigned_to, status);

-- Триггер для обновления updated_at в translation_providers
CREATE OR REPLACE FUNCTION update_translation_providers_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_translation_providers_timestamp
BEFORE UPDATE ON translation_providers
FOR EACH ROW
EXECUTE FUNCTION update_translation_providers_updated_at();

-- Добавляем колонку version в основную таблицу translations если её нет
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns 
                   WHERE table_name = 'translations' 
                   AND column_name = 'version') THEN
        ALTER TABLE translations ADD COLUMN version INTEGER DEFAULT 1;
    END IF;
END $$;

-- Функция для создания версии при изменении перевода
CREATE OR REPLACE FUNCTION create_translation_version()
RETURNS TRIGGER AS $$
DECLARE
    next_version INTEGER;
BEGIN
    -- Только если изменился текст перевода
    IF OLD.translated_text != NEW.translated_text THEN
        -- Получаем следующий номер версии
        SELECT COALESCE(MAX(version_number), 0) + 1 INTO next_version
        FROM translation_versions
        WHERE translation_id = NEW.id;
        
        -- Создаем запись версии
        INSERT INTO translation_versions (
            translation_id, version_number, entity_type, entity_id,
            language, field_name, translated_text, metadata
        ) VALUES (
            NEW.id, next_version, NEW.entity_type, NEW.entity_id,
            NEW.language, NEW.field_name, OLD.translated_text, OLD.metadata
        );
        
        -- Обновляем номер версии в основной таблице
        NEW.version = next_version;
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Триггер для версионирования переводов
CREATE TRIGGER translation_versioning
BEFORE UPDATE ON translations
FOR EACH ROW
EXECUTE FUNCTION create_translation_version();

-- Добавляем начальные данные для провайдеров
INSERT INTO translation_providers (name, provider_type, settings, priority, is_active) VALUES
('OpenAI GPT-4', 'openai', '{"model": "gpt-4-turbo-preview", "temperature": 0.3}', 1, true),
('OpenAI GPT-3.5', 'openai', '{"model": "gpt-3.5-turbo", "temperature": 0.3}', 2, true),
('Google Translate', 'google', '{"api_version": "v3"}', 3, true),
('Manual Translation', 'manual', '{}', 10, true)
ON CONFLICT (name) DO NOTHING;