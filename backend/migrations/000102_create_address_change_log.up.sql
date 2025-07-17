-- Миграция 102: Создание таблицы логов изменений адресов
-- Для аудита и анализа изменений адресов пользователями

CREATE TABLE IF NOT EXISTS address_change_log (
    id BIGSERIAL PRIMARY KEY,
    listing_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    old_address TEXT,
    new_address TEXT,
    old_location geography(Point, 4326),
    new_location geography(Point, 4326),
    change_reason VARCHAR(100),
    confidence_before NUMERIC(3, 2),
    confidence_after NUMERIC(3, 2),
    ip_address INET,
    user_agent TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    -- Foreign keys
    CONSTRAINT fk_address_log_listing_id 
        FOREIGN KEY (listing_id) 
        REFERENCES marketplace_listings(id) 
        ON DELETE CASCADE,
    
    CONSTRAINT fk_address_log_user_id 
        FOREIGN KEY (user_id) 
        REFERENCES users(id) 
        ON DELETE CASCADE
);

-- Индексы для производительности поиска и аналитики
CREATE INDEX IF NOT EXISTS idx_address_log_listing_id ON address_change_log (listing_id);
CREATE INDEX IF NOT EXISTS idx_address_log_user_id ON address_change_log (user_id);
CREATE INDEX IF NOT EXISTS idx_address_log_created_at ON address_change_log (created_at);
CREATE INDEX IF NOT EXISTS idx_address_log_change_reason ON address_change_log (change_reason);
CREATE INDEX IF NOT EXISTS idx_address_log_confidence_after ON address_change_log (confidence_after);

-- Пространственные индексы для геоаналитики
CREATE INDEX IF NOT EXISTS idx_address_log_old_location ON address_change_log USING GIST (old_location);
CREATE INDEX IF NOT EXISTS idx_address_log_new_location ON address_change_log USING GIST (new_location);

-- Комментарии к таблице и полям
COMMENT ON TABLE address_change_log IS 'Лог изменений адресов объявлений для аудита и анализа';
COMMENT ON COLUMN address_change_log.listing_id IS 'ID объявления, для которого изменился адрес';
COMMENT ON COLUMN address_change_log.user_id IS 'ID пользователя, который изменил адрес';
COMMENT ON COLUMN address_change_log.old_address IS 'Предыдущий адрес (текстовое представление)';
COMMENT ON COLUMN address_change_log.new_address IS 'Новый адрес (текстовое представление)';
COMMENT ON COLUMN address_change_log.old_location IS 'Предыдущие координаты';
COMMENT ON COLUMN address_change_log.new_location IS 'Новые координаты';
COMMENT ON COLUMN address_change_log.change_reason IS 'Причина изменения: geocoded, map_click, manual_correction, etc.';
COMMENT ON COLUMN address_change_log.confidence_before IS 'Уровень доверия до изменения';
COMMENT ON COLUMN address_change_log.confidence_after IS 'Уровень доверия после изменения';
COMMENT ON COLUMN address_change_log.ip_address IS 'IP адрес пользователя для аудита';
COMMENT ON COLUMN address_change_log.user_agent IS 'User-Agent браузера для анализа источников изменений';

-- Функция для автоматической очистки старых логов (опционально)
-- Удаляет записи старше 2 лет для экономии места
CREATE OR REPLACE FUNCTION cleanup_old_address_logs()
RETURNS INTEGER AS $$
DECLARE
    deleted_count INTEGER;
BEGIN
    DELETE FROM address_change_log 
    WHERE created_at < CURRENT_TIMESTAMP - INTERVAL '2 years';
    
    GET DIAGNOSTICS deleted_count = ROW_COUNT;
    
    -- Логируем результат очистки
    INSERT INTO address_change_log (
        listing_id, user_id, change_reason, created_at
    ) VALUES (
        0, 0, 'cleanup_old_logs_' || deleted_count, CURRENT_TIMESTAMP
    );
    
    RETURN deleted_count;
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION cleanup_old_address_logs() IS 'Функция для очистки логов адресов старше 2 лет';