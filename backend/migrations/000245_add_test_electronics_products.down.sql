-- Удаление тестовых товаров электроники

-- Удаляем переводы товаров
DELETE FROM translations 
WHERE entity_type = 'listing' 
AND entity_id IN (
    SELECT id FROM marketplace_listings 
    WHERE title IN (
        'MacBook Pro 16" M3 Max',
        'Gaming Laptop ASUS ROG Strix G15',
        'Dell XPS 15 9530 Laptop',
        'Lenovo ThinkPad X1 Carbon Gen 11',
        'HP Pavilion 15 Student Laptop',
        'iPhone 15 Pro Max 256GB',
        'Samsung Galaxy S24 Ultra',
        'iPad Pro 12.9" M2 with Magic Keyboard',
        'Samsung Galaxy Tab S9+ 5G',
        'Sony WH-1000XM5 Headphones',
        'AirPods Pro 2nd Generation'
    )
);

-- Удаляем изображения товаров
DELETE FROM marketplace_images 
WHERE listing_id IN (
    SELECT id FROM marketplace_listings 
    WHERE title IN (
        'MacBook Pro 16" M3 Max',
        'Gaming Laptop ASUS ROG Strix G15',
        'Dell XPS 15 9530 Laptop',
        'Lenovo ThinkPad X1 Carbon Gen 11',
        'HP Pavilion 15 Student Laptop',
        'iPhone 15 Pro Max 256GB',
        'Samsung Galaxy S24 Ultra',
        'iPad Pro 12.9" M2 with Magic Keyboard',
        'Samsung Galaxy Tab S9+ 5G',
        'Sony WH-1000XM5 Headphones',
        'AirPods Pro 2nd Generation'
    )
);

-- Удаляем товары
DELETE FROM marketplace_listings 
WHERE title IN (
    'MacBook Pro 16" M3 Max',
    'Gaming Laptop ASUS ROG Strix G15',
    'Dell XPS 15 9530 Laptop',
    'Lenovo ThinkPad X1 Carbon Gen 11',
    'HP Pavilion 15 Student Laptop',
    'iPhone 15 Pro Max 256GB',
    'Samsung Galaxy S24 Ultra',
    'iPad Pro 12.9" M2 with Magic Keyboard',
    'Samsung Galaxy Tab S9+ 5G',
    'Sony WH-1000XM5 Headphones',
    'AirPods Pro 2nd Generation'
);

-- Удаляем переводы категорий
DELETE FROM translations 
WHERE entity_type = 'category' 
AND entity_id IN (2001, 2002, 2003, 2004, 2005);

-- Удаляем категории
DELETE FROM marketplace_categories 
WHERE id IN (2002, 2003, 2004, 2005);

DELETE FROM marketplace_categories 
WHERE id = 2001;