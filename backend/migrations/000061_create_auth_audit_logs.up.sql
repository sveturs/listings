-- Добавление недостающих колонок в существующую таблицу auth.audit_logs
ALTER TABLE auth.audit_logs 
    ADD COLUMN IF NOT EXISTS entity_type VARCHAR(50),
    ADD COLUMN IF NOT EXISTS entity_id VARCHAR(100),
    ADD COLUMN IF NOT EXISTS ip_address INET,
    ADD COLUMN IF NOT EXISTS user_agent TEXT,
    ADD COLUMN IF NOT EXISTS details JSONB;

-- Изменение типов существующих колонок для совместимости
ALTER TABLE auth.audit_logs 
    ALTER COLUMN user_id TYPE BIGINT USING user_id::BIGINT,
    ALTER COLUMN action SET NOT NULL,
    ALTER COLUMN created_at TYPE TIMESTAMP WITH TIME ZONE USING created_at AT TIME ZONE 'UTC';

-- Создание индексов (если они не существуют)
CREATE INDEX IF NOT EXISTS idx_audit_logs_user_id ON auth.audit_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_audit_logs_action ON auth.audit_logs(action);
CREATE INDEX IF NOT EXISTS idx_audit_logs_entity ON auth.audit_logs(entity_type, entity_id);
CREATE INDEX IF NOT EXISTS idx_audit_logs_created_at ON auth.audit_logs(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_audit_logs_ip_address ON auth.audit_logs(ip_address);

-- Комментарии к таблице
COMMENT ON TABLE auth.audit_logs IS 'Аудит лог для отслеживания действий пользователей в Auth Service';
COMMENT ON COLUMN auth.audit_logs.user_id IS 'ID пользователя, выполнившего действие';
COMMENT ON COLUMN auth.audit_logs.action IS 'Тип действия (login, logout, token_refresh, password_change и т.д.)';
COMMENT ON COLUMN auth.audit_logs.entity_type IS 'Тип сущности (user, token, session и т.д.)';
COMMENT ON COLUMN auth.audit_logs.entity_id IS 'ID сущности';
COMMENT ON COLUMN auth.audit_logs.ip_address IS 'IP адрес пользователя';
COMMENT ON COLUMN auth.audit_logs.user_agent IS 'User Agent браузера';
COMMENT ON COLUMN auth.audit_logs.details IS 'Дополнительные детали в формате JSON';