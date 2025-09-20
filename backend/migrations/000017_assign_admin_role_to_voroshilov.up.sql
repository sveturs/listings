-- Назначение роли администратора пользователю voroshilovdo@gmail.com
UPDATE users
SET role_id = 2 -- admin role
WHERE email = 'voroshilovdo@gmail.com';

-- Также добавим запись в user_roles для полной совместимости
INSERT INTO user_roles (user_id, role_id, assigned_by, assigned_at)
SELECT
    u.id,
    2, -- admin role
    u.id, -- self-assigned for initial setup
    NOW()
FROM users u
WHERE u.email = 'voroshilovdo@gmail.com'
ON CONFLICT (user_id, role_id) DO NOTHING;