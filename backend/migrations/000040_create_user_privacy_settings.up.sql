-- Триггер для автоматического обновления updated_at
CREATE OR REPLACE FUNCTION update_user_privacy_settings_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS update_user_privacy_settings_updated_at ON update_user_privacy_settings_updated_at;
CREATE TRIGGER update_user_privacy_settings_updated_at
    BEFORE UPDATE ON user_privacy_settings
    FOR EACH ROW
    EXECUTE FUNCTION update_user_privacy_settings_updated_at();

-- Вставляем дефолтные настройки для существующих пользователей
INSERT INTO user_privacy_settings (user_id)
SELECT id FROM users
ON CONFLICT (user_id) DO NOTHING;
