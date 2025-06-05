-- Удаляем функцию очистки
DROP FUNCTION IF EXISTS cleanup_expired_refresh_tokens();

-- Удаляем таблицу refresh_tokens
DROP TABLE IF EXISTS refresh_tokens;