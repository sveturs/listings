-- Удаляем тестовые объявления
DELETE FROM marketplace_listings
WHERE title IN (
    'Tommy Hilfiger Polo majica',
    'Calvin Klein džins',
    'Zara haljina',
    'Adidas Ultraboost 22',
    'Nike Air Force 1',
    'Guess sat za dame',
    'Ray-Ban Aviator naočare',
    'iPhone 13 128GB',
    'Samsung Galaxy S22',
    'MacBook Air M2 2023'
) AND user_id = 1;