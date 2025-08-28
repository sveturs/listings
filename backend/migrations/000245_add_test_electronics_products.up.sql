-- Добавление тестовых товаров электроники для демонстрации поиска

-- Добавляем категорию Electronics если её нет
INSERT INTO marketplace_categories (id, parent_id, name, slug, icon, created_at)
VALUES (2001, NULL, 'Electronics', 'electronics', 'laptop', NOW())
ON CONFLICT (id) DO NOTHING;

-- Добавляем подкатегории
INSERT INTO marketplace_categories (id, parent_id, name, slug, icon, created_at)
VALUES 
    (2002, 2001, 'Laptops & Computers', 'laptops-computers', 'laptop', NOW()),
    (2003, 2001, 'Smartphones', 'smartphones', 'smartphone', NOW()),
    (2004, 2001, 'Tablets', 'tablets', 'tablet', NOW()),
    (2005, 2001, 'Audio & Headphones', 'audio-headphones', 'headphones', NOW())
ON CONFLICT (id) DO NOTHING;

-- Добавляем переводы для категорий
INSERT INTO translations (entity_type, entity_id, field_name, language, translated_text, is_verified, created_at, updated_at)
VALUES
    ('category', 2001, 'name', 'ru', 'Электроника', true, NOW(), NOW()),
    ('category', 2001, 'name', 'en', 'Electronics', true, NOW(), NOW()),
    ('category', 2002, 'name', 'ru', 'Ноутбуки и компьютеры', true, NOW(), NOW()),
    ('category', 2002, 'name', 'en', 'Laptops & Computers', true, NOW(), NOW()),
    ('category', 2003, 'name', 'ru', 'Смартфоны', true, NOW(), NOW()),
    ('category', 2003, 'name', 'en', 'Smartphones', true, NOW(), NOW()),
    ('category', 2004, 'name', 'ru', 'Планшеты', true, NOW(), NOW()),
    ('category', 2004, 'name', 'en', 'Tablets', true, NOW(), NOW()),
    ('category', 2005, 'name', 'ru', 'Аудио и наушники', true, NOW(), NOW()),
    ('category', 2005, 'name', 'en', 'Audio & Headphones', true, NOW(), NOW())
ON CONFLICT DO NOTHING;

-- Добавляем тестовые товары - ноутбуки
INSERT INTO marketplace_listings (
    user_id, category_id, title, description, price, 
    condition, status, address_city, address_country, location, show_on_map,
    latitude, longitude, created_at, updated_at
)
VALUES
    -- Ноутбуки
    (1, 2002, 'MacBook Pro 16" M3 Max', 
     'Brand new MacBook Pro 16-inch with M3 Max chip, 36GB RAM, 1TB SSD. Perfect for professional work, video editing, and development. Includes original packaging and warranty.',
     450000, 'new', 'active', 'Novi Sad', 'Serbia', 'Liman 4', true,
     45.2481, 19.8335, NOW(), NOW()),
    
    (1, 2002, 'Gaming Laptop ASUS ROG Strix G15', 
     'Powerful gaming laptop with RTX 4070, AMD Ryzen 9 7940HS, 32GB DDR5 RAM, 2TB NVMe SSD. 165Hz display, RGB keyboard. Excellent condition, used for 3 months.',
     280000, 'like_new', 'active', 'Belgrade', 'Serbia', 'Vračar', true,
     44.7989, 20.4651, NOW(), NOW()),
    
    (2, 2002, 'Dell XPS 15 9530 Laptop', 
     'Dell XPS 15 with Intel Core i7-13700H, 16GB RAM, 512GB SSD, NVIDIA RTX 4050. 15.6" OLED display. Perfect for work and creative tasks. Minor scratches on the lid.',
     180000, 'good', 'active', 'Novi Sad', 'Serbia', 'Centar', true,
     45.2551, 19.8451, NOW(), NOW()),
    
    (2, 2002, 'Lenovo ThinkPad X1 Carbon Gen 11', 
     'Business ultrabook laptop. Intel Core i7-1365U, 16GB RAM, 1TB SSD. 14" display, excellent battery life. Ideal for business professionals.',
     165000, 'like_new', 'active', 'Subotica', 'Serbia', 'Centar', true,
     46.1001, 19.6650, NOW(), NOW()),
    
    (1, 2002, 'HP Pavilion 15 Student Laptop', 
     'Affordable laptop for students. Intel Core i5-1235U, 8GB RAM, 256GB SSD. 15.6" Full HD display. Good for studying, office work, and light gaming.',
     65000, 'used', 'active', 'Novi Sad', 'Serbia', 'Novo Naselje', true,
     45.2671, 19.8135, NOW(), NOW()),

    -- Смартфоны
    (2, 2003, 'iPhone 15 Pro Max 256GB', 
     'Latest iPhone 15 Pro Max, Titanium Blue, 256GB storage. Excellent camera system, A17 Pro chip. Includes original box and accessories.',
     180000, 'new', 'active', 'Belgrade', 'Serbia', 'Stari Grad', true,
     44.8176, 20.4633, NOW(), NOW()),
    
    (1, 2003, 'Samsung Galaxy S24 Ultra', 
     'Samsung flagship with S Pen, 512GB storage, 12GB RAM. Amazing camera with 200MP sensor. Perfect condition with warranty.',
     165000, 'like_new', 'active', 'Novi Sad', 'Serbia', 'Liman 3', true,
     45.2431, 19.8285, NOW(), NOW()),

    -- Планшеты
    (2, 2004, 'iPad Pro 12.9" M2 with Magic Keyboard', 
     'iPad Pro 12.9-inch with M2 chip, 256GB, WiFi + Cellular. Includes Magic Keyboard and Apple Pencil 2. Perfect for creative work.',
     195000, 'like_new', 'active', 'Belgrade', 'Serbia', 'Dedinje', true,
     44.7489, 20.4451, NOW(), NOW()),
    
    (1, 2004, 'Samsung Galaxy Tab S9+ 5G', 
     'Premium Android tablet with S Pen, 12.4" AMOLED display, 256GB storage. Great for productivity and entertainment.',
     110000, 'new', 'active', 'Novi Sad', 'Serbia', 'Podbara', true,
     45.2621, 19.8585, NOW(), NOW()),

    -- Аудио
    (2, 2005, 'Sony WH-1000XM5 Headphones', 
     'Premium noise-cancelling headphones. Industry-leading ANC, 30-hour battery life. Includes carrying case and cables.',
     45000, 'like_new', 'active', 'Belgrade', 'Serbia', 'Zemun', true,
     44.8433, 20.4111, NOW(), NOW()),
    
    (1, 2005, 'AirPods Pro 2nd Generation', 
     'Apple AirPods Pro with USB-C charging case. Active noise cancellation, transparency mode. Excellent sound quality.',
     35000, 'new', 'active', 'Novi Sad', 'Serbia', 'Grbavica', true,
     45.2471, 19.8185, NOW(), NOW());

-- Добавляем переводы для товаров
INSERT INTO translations (entity_type, entity_id, field_name, language, translated_text, is_verified, created_at, updated_at)
SELECT 
    'listing' as entity_type,
    id as entity_id,
    'title' as field_name,
    'ru' as language,
    CASE 
        WHEN title LIKE '%MacBook%' THEN 'MacBook Pro 16" M3 Max - Новейший ноутбук Apple'
        WHEN title LIKE '%ASUS ROG%' THEN 'Игровой ноутбук ASUS ROG Strix G15'
        WHEN title LIKE '%Dell XPS%' THEN 'Ноутбук Dell XPS 15 9530'
        WHEN title LIKE '%ThinkPad%' THEN 'Ультрабук Lenovo ThinkPad X1 Carbon Gen 11'
        WHEN title LIKE '%HP Pavilion%' THEN 'Ноутбук HP Pavilion 15 для студентов'
        WHEN title LIKE '%iPhone%' THEN 'iPhone 15 Pro Max 256ГБ'
        WHEN title LIKE '%Galaxy S24%' THEN 'Samsung Galaxy S24 Ultra'
        WHEN title LIKE '%iPad Pro%' THEN 'iPad Pro 12.9" M2 с Magic Keyboard'
        WHEN title LIKE '%Galaxy Tab%' THEN 'Планшет Samsung Galaxy Tab S9+ 5G'
        WHEN title LIKE '%Sony%' THEN 'Наушники Sony WH-1000XM5'
        WHEN title LIKE '%AirPods%' THEN 'AirPods Pro 2-го поколения'
        ELSE title
    END as translated_text,
    true as is_verified,
    NOW() as created_at,
    NOW() as updated_at
FROM marketplace_listings
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

-- Добавляем изображения для новых товаров (заглушки)
INSERT INTO marketplace_images (listing_id, file_path, file_name, file_size, content_type, public_url, storage_type, is_main, created_at)
SELECT 
    id as listing_id,
    CONCAT(id, '/placeholder-', id, '.jpg') as file_path,
    CONCAT('placeholder-', id, '.jpg') as file_name,
    1024 as file_size,
    'image/jpeg' as content_type,
    CONCAT('/listings/', id, '/placeholder-', id, '.jpg') as public_url,
    'minio' as storage_type,
    true as is_main,
    NOW() as created_at
FROM marketplace_listings
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