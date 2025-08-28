-- Обновляем поле provider для пользователей с Google ID
UPDATE users 
SET provider = 'google',
    updated_at = CURRENT_TIMESTAMP
WHERE google_id IS NOT NULL 
  AND google_id != ''
  AND (provider IS NULL OR provider = '');

-- Обновляем поле provider для пользователей без Google ID
UPDATE users 
SET provider = 'email',
    updated_at = CURRENT_TIMESTAMP  
WHERE (google_id IS NULL OR google_id = '')
  AND (provider IS NULL OR provider = '');

-- Добавляем комментарий к полю для ясности
COMMENT ON COLUMN users.provider IS 'Способ регистрации/авторизации: google, email, facebook и т.д.';