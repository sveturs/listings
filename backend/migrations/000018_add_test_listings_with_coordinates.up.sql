-- Добавляем тестовые объявления с распределенными координатами для демонстрации карты

-- Вставляем новые объявления в районе Белграда
INSERT INTO marketplace_listings (
    user_id, title, description, price, category_id,
    status, latitude, longitude, created_at, updated_at
) VALUES
-- Центр - Tommy Hilfiger Polo majica
(1, 'Tommy Hilfiger Polo majica', 'Оригинальная поло майка Tommy Hilfiger. Размер L, состояние отличное.', 3500, 1003,
'active', 44.8176, 20.4649, NOW(), NOW()),

-- 1км на север - Calvin Klein džins
(1, 'Calvin Klein džins', 'Стильные джинсы Calvin Klein, размер 32/34, как новые.', 4500, 1003,
'active', 44.8266, 20.4649, NOW(), NOW()),

-- 1.5км на восток - Zara haljina
(1, 'Zara haljina', 'Элегантное платье Zara, размер M, надевалось несколько раз.', 2500, 1003,
'active', 44.8176, 20.4829, NOW(), NOW()),

-- 2км на юг - Adidas Ultraboost 22
(1, 'Adidas Ultraboost 22', 'Кроссовки для бега Adidas Ultraboost, размер 42, отличное состояние.', 8500, 1003,
'active', 44.7996, 20.4649, NOW(), NOW()),

-- 1.5км на запад - Nike Air Force 1
(1, 'Nike Air Force 1', 'Классические кроссовки Nike Air Force 1, белые, размер 41.', 9500, 1003,
'active', 44.8176, 20.4469, NOW(), NOW()),

-- 2км на северо-восток - Guess sat za dame
(1, 'Guess sat za dame', 'Женские часы Guess с кристаллами, в подарочной упаковке.', 12000, 1005,
'active', 44.8316, 20.4809, NOW(), NOW()),

-- 1.8км на юго-запад - Ray-Ban Aviator naočare
(1, 'Ray-Ban Aviator naočare', 'Оригинальные солнцезащитные очки Ray-Ban Aviator с чехлом.', 15000, 1005,
'active', 44.8036, 20.4489, NOW(), NOW()),

-- Дополнительные объявления для разнообразия
-- 1.2км на северо-запад - iPhone 13
(1, 'iPhone 13 128GB', 'iPhone 13, 128GB, синий цвет, состояние идеальное, полный комплект.', 75000, 1001,
'active', 44.8256, 20.4529, NOW(), NOW()),

-- 2.5км на юго-восток - Samsung Galaxy S22
(1, 'Samsung Galaxy S22', 'Samsung Galaxy S22, 256GB, черный, с гарантией до 2025.', 65000, 1001,
'active', 44.8026, 20.4869, NOW(), NOW()),

-- 0.8км от центра - MacBook Air M2
(1, 'MacBook Air M2 2023', 'MacBook Air с процессором M2, 8GB RAM, 256GB SSD, как новый.', 135000, 1001,
'active', 44.8216, 20.4709, NOW(), NOW());