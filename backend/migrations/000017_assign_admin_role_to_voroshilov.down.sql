-- Откат: удаление роли администратора у пользователя voroshilovdo@gmail.com
UPDATE users
SET role_id = NULL
WHERE email = 'voroshilovdo@gmail.com';

-- Удаляем запись из user_roles
DELETE FROM user_roles
WHERE user_id = (SELECT id FROM users WHERE email = 'voroshilovdo@gmail.com')
  AND role_id = 2;