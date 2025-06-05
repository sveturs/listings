-- Создаем таблицу для refresh токенов
CREATE TABLE IF NOT EXISTS refresh_tokens (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token VARCHAR(255) UNIQUE NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    -- Информация об устройстве/браузере
    user_agent TEXT,
    ip VARCHAR(45), -- Поддержка IPv6
    device_name VARCHAR(100),
    
    -- Для отзыва токенов
    is_revoked BOOLEAN DEFAULT FALSE,
    revoked_at TIMESTAMP
);

-- Индексы для быстрого поиска
CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens(user_id) WHERE NOT is_revoked;
CREATE INDEX idx_refresh_tokens_token ON refresh_tokens(token) WHERE NOT is_revoked;
CREATE INDEX idx_refresh_tokens_expires_at ON refresh_tokens(expires_at) WHERE NOT is_revoked;

-- Функция для автоматической очистки истекших токенов
CREATE OR REPLACE FUNCTION cleanup_expired_refresh_tokens()
RETURNS void AS $$
BEGIN
    DELETE FROM refresh_tokens 
    WHERE expires_at < CURRENT_TIMESTAMP 
    OR (is_revoked = TRUE AND revoked_at < CURRENT_TIMESTAMP - INTERVAL '30 days');
END;
$$ LANGUAGE plpgsql;

-- Комментарии для документации
COMMENT ON TABLE refresh_tokens IS 'Хранит refresh токены для JWT аутентификации';
COMMENT ON COLUMN refresh_tokens.device_name IS 'Опциональное имя устройства для управления сессиями';
COMMENT ON COLUMN refresh_tokens.is_revoked IS 'Флаг для принудительного отзыва токена';