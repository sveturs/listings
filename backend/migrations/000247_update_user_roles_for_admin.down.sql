-- Возвращаем роль пользователей на обычного пользователя
UPDATE users 
SET role_id = (SELECT id FROM roles WHERE name = 'user')
WHERE email IN ('boxmail386@gmail.com', 'admin@svetu.rs');