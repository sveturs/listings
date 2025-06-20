-- Создание тестового пользователя для проверки UI
-- Email: test@example.com
-- Password: test123

-- Проверяем, существует ли пользователь
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM users WHERE email = 'test@example.com') THEN
        INSERT INTO users (name, email, password, provider)
        VALUES ('Test User', 'test@example.com', '$2a$10$YourHashedPasswordHere', 'email');
    END IF;
END $$;

-- Получаем ID пользователя
WITH test_user AS (
    SELECT id FROM users WHERE email = 'test@example.com'
)
-- Создаем витрину для этого пользователя, если её нет
INSERT INTO user_storefronts (user_id, name, slug, description, city, country, status)
SELECT 
    tu.id,
    'Test Electronics Store',
    'test-electronics-store',
    'A test store for UI testing',
    'Belgrade', 
    'Serbia',
    'active'
FROM test_user tu
WHERE NOT EXISTS (
    SELECT 1 FROM user_storefronts 
    WHERE user_id = tu.id AND slug = 'test-electronics-store'
);

-- Также обновим таблицу storefronts для совместимости
WITH test_user AS (
    SELECT id FROM users WHERE email = 'test@example.com'
)
UPDATE storefronts 
SET is_active = true, 
    name = 'Test Electronics Store',
    description = 'A test store for UI testing',
    city = 'Belgrade'
WHERE user_id = (SELECT id FROM test_user);

-- Добавим несколько тестовых товаров для витрины
WITH test_data AS (
    SELECT 
        u.id as user_id,
        us.id as storefront_id
    FROM users u
    JOIN user_storefronts us ON u.id = us.user_id
    WHERE u.email = 'test@example.com' 
    AND us.slug = 'test-electronics-store'
)
INSERT INTO marketplace_listings (
    user_id, storefront_id, category_id, title, description, 
    price, condition, status, location, address_city, address_country
)
SELECT 
    td.user_id,
    td.storefront_id,
    1, -- Electronics category
    'Test Product ' || generate_series,
    'This is test product number ' || generate_series,
    100.00 * generate_series,
    'new',
    'active',
    'Belgrade',
    'Belgrade',
    'Serbia'
FROM test_data td
CROSS JOIN generate_series(1, 3)
WHERE NOT EXISTS (
    SELECT 1 FROM marketplace_listings 
    WHERE storefront_id = td.storefront_id
);

-- Вывести информацию о созданной витрине
SELECT 
    u.id as user_id,
    u.email,
    u.name,
    us.id as storefront_id,
    us.name as storefront_name,
    us.slug,
    (SELECT COUNT(*) FROM marketplace_listings WHERE storefront_id = us.id) as product_count
FROM users u
JOIN user_storefronts us ON u.id = us.user_id
WHERE u.email = 'test@example.com';