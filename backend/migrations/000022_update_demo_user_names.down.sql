-- Откат изменений имен пользователей
-- Возвращаем Demo User для тестовых аккаунтов

UPDATE users
SET name = 'Demo User'
WHERE email IN ('demo@svetu.rs', 'test@example.com');