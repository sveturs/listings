-- Удаление триггера и функции
DROP TRIGGER IF EXISTS update_user_privacy_settings_updated_at ON user_privacy_settings;
DROP FUNCTION IF EXISTS update_user_privacy_settings_updated_at();

-- Удаление таблицы
DROP TABLE IF EXISTS user_privacy_settings;