-- Обновляем роль для пользователя boxmail386@gmail.com на admin
UPDATE users 
SET role_id = (SELECT id FROM roles WHERE name = 'admin')
WHERE email = 'boxmail386@gmail.com';

-- Также обновляем роль для admin@svetu.rs на super_admin
UPDATE users 
SET role_id = (SELECT id FROM roles WHERE name = 'super_admin')
WHERE email = 'admin@svetu.rs';