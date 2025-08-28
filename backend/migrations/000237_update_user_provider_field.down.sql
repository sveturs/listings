-- Откатываем изменения - очищаем поле provider
UPDATE users 
SET provider = '',
    updated_at = CURRENT_TIMESTAMP
WHERE provider IN ('google', 'email');

-- Удаляем комментарий
COMMENT ON COLUMN users.provider IS NULL;