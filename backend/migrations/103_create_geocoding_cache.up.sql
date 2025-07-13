-- Миграция 103: Создание таблицы кэша геокодирования
-- Для оптимизации производительности и экономии API запросов

CREATE TABLE IF NOT EXISTS geocoding_cache (
    id BIGSERIAL PRIMARY KEY,
    input_address TEXT NOT NULL,
    normalized_address TEXT NOT NULL,
    location geography(Point, 4326) NOT NULL,
    address_components JSONB NOT NULL,
    formatted_address TEXT NOT NULL,
    confidence NUMERIC(3, 2) NOT NULL,
    provider VARCHAR(50) NOT NULL DEFAULT 'mapbox',
    language VARCHAR(5) DEFAULT 'en',
    country_code VARCHAR(2),
    cache_hits BIGINT DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP NOT NULL DEFAULT (CURRENT_TIMESTAMP + INTERVAL '30 days'),
    
    -- Уникальный индекс для предотвращения дублирования
    CONSTRAINT uk_geocoding_cache_normalized 
        UNIQUE (normalized_address, language, country_code)
);

-- Индексы для быстрого поиска
CREATE INDEX IF NOT EXISTS idx_geocoding_cache_input_address ON geocoding_cache (input_address);
CREATE INDEX IF NOT EXISTS idx_geocoding_cache_normalized ON geocoding_cache (normalized_address);
CREATE INDEX IF NOT EXISTS idx_geocoding_cache_location ON geocoding_cache USING GIST (location);
CREATE INDEX IF NOT EXISTS idx_geocoding_cache_expires_at ON geocoding_cache (expires_at);
CREATE INDEX IF NOT EXISTS idx_geocoding_cache_provider ON geocoding_cache (provider);
CREATE INDEX IF NOT EXISTS idx_geocoding_cache_confidence ON geocoding_cache (confidence);
CREATE INDEX IF NOT EXISTS idx_geocoding_cache_cache_hits ON geocoding_cache (cache_hits);
CREATE INDEX IF NOT EXISTS idx_geocoding_cache_country_lang ON geocoding_cache (country_code, language);

-- GIN индекс для поиска по компонентам адреса
CREATE INDEX IF NOT EXISTS idx_geocoding_cache_address_components ON geocoding_cache USING GIN (address_components);

-- Триггер для автоматического обновления updated_at
CREATE OR REPLACE FUNCTION update_geocoding_cache_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_geocoding_cache_updated_at
    BEFORE UPDATE ON geocoding_cache
    FOR EACH ROW
    EXECUTE FUNCTION update_geocoding_cache_updated_at();

-- Функция для очистки устаревшего кэша
CREATE OR REPLACE FUNCTION cleanup_expired_geocoding_cache()
RETURNS INTEGER AS $$
DECLARE
    deleted_count INTEGER;
BEGIN
    DELETE FROM geocoding_cache WHERE expires_at < CURRENT_TIMESTAMP;
    GET DIAGNOSTICS deleted_count = ROW_COUNT;
    RETURN deleted_count;
END;
$$ LANGUAGE plpgsql;

-- Автоматический триггер очистки при каждой вставке (с ограничением частоты)
CREATE OR REPLACE FUNCTION trigger_cleanup_geocoding_cache()
RETURNS TRIGGER AS $$
BEGIN
    -- Очищаем кэш только в 1% случаев для избежания частых очисток
    IF random() < 0.01 THEN
        PERFORM cleanup_expired_geocoding_cache();
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_cleanup_geocoding_cache
    AFTER INSERT ON geocoding_cache
    FOR EACH STATEMENT
    EXECUTE FUNCTION trigger_cleanup_geocoding_cache();

-- Функция для получения статистики кэша
CREATE OR REPLACE FUNCTION get_geocoding_cache_stats()
RETURNS TABLE (
    total_entries BIGINT,
    active_entries BIGINT,
    expired_entries BIGINT,
    total_cache_hits BIGINT,
    avg_confidence NUMERIC,
    top_providers TEXT[]
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        COUNT(*) as total_entries,
        COUNT(*) FILTER (WHERE expires_at > CURRENT_TIMESTAMP) as active_entries,
        COUNT(*) FILTER (WHERE expires_at <= CURRENT_TIMESTAMP) as expired_entries,
        COALESCE(SUM(cache_hits), 0) as total_cache_hits,
        ROUND(AVG(confidence), 2) as avg_confidence,
        ARRAY_AGG(DISTINCT provider) as top_providers
    FROM geocoding_cache;
END;
$$ LANGUAGE plpgsql;

-- Комментарии к таблице и полям
COMMENT ON TABLE geocoding_cache IS 'Кэш результатов геокодирования для оптимизации производительности';
COMMENT ON COLUMN geocoding_cache.input_address IS 'Оригинальный введенный адрес';
COMMENT ON COLUMN geocoding_cache.normalized_address IS 'Нормализованный адрес для поиска дубликатов';
COMMENT ON COLUMN geocoding_cache.location IS 'Геокоординаты результата геокодирования';
COMMENT ON COLUMN geocoding_cache.address_components IS 'JSON с компонентами адреса (страна, город, улица и т.д.)';
COMMENT ON COLUMN geocoding_cache.formatted_address IS 'Отформатированный адрес от провайдера';
COMMENT ON COLUMN geocoding_cache.confidence IS 'Уровень доверия результату геокодирования (0.0-1.0)';
COMMENT ON COLUMN geocoding_cache.provider IS 'Провайдер геокодирования (mapbox, nominatim, etc.)';
COMMENT ON COLUMN geocoding_cache.language IS 'Язык запроса (en, ru, etc.)';
COMMENT ON COLUMN geocoding_cache.country_code IS 'Код страны для фильтрации (RS, HR, etc.)';
COMMENT ON COLUMN geocoding_cache.cache_hits IS 'Количество попаданий в кэш для аналитики';
COMMENT ON COLUMN geocoding_cache.expires_at IS 'Время истечения кэша (по умолчанию 30 дней)';

COMMENT ON FUNCTION cleanup_expired_geocoding_cache() IS 'Функция для очистки устаревших записей кэша';
COMMENT ON FUNCTION get_geocoding_cache_stats() IS 'Функция для получения статистики использования кэша';