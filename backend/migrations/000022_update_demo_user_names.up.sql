-- Обновляем имена пользователей, которые имеют "Demo User" на более осмысленные
-- основанные на их email адресах

UPDATE users
SET name = CASE
    WHEN email = 'demo@svetu.rs' THEN 'Demo Account'
    WHEN email = 'test@example.com' THEN 'Test User'
    ELSE SPLIT_PART(email, '@', 1) -- Используем часть до @ как имя
END
WHERE name = 'Demo User' OR name IS NULL OR name = '';

-- Для пользователей с реальными email делаем имя более читаемым
UPDATE users
SET name = INITCAP(REPLACE(SPLIT_PART(email, '@', 1), '.', ' '))
WHERE name = SPLIT_PART(email, '@', 1)
  AND email NOT IN ('demo@svetu.rs', 'test@example.com');