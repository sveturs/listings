-- Откат синхронизации пользователей с auth сервисом
-- ВНИМАНИЕ: Откат может привести к несоответствию с auth сервисом!

BEGIN;

-- 1. Восстанавливаем пользователя ID=10
UPDATE users
SET
    email = REPLACE(email, 'archived_', ''),
    name = REPLACE(name, '[ARCHIVED] ', '')
WHERE id = 10 AND email LIKE 'archived_%';

-- 2. Восстанавливаем старый email для пользователя ID=7
UPDATE users
SET
    email = COALESCE(old_email, 'demo@svetu.rs'),
    name = 'Demo Account'
WHERE id = 7;

-- 3. Переносим обратно все связанные данные с user_id=7 на user_id=10
-- (только те записи, которые были созданы после определенной даты)

-- Примечание: Полный откат сложен, так как невозможно точно определить,
-- какие записи принадлежали какому пользователю изначально.
-- Рекомендуется не выполнять откат этой миграции в продакшене!

-- 4. Удаляем временную колонку если она существует
ALTER TABLE users DROP COLUMN IF EXISTS old_email;

COMMIT;