-- Категории объявлений
CREATE TABLE marketplace_categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(100) NOT NULL UNIQUE,
    parent_id INT REFERENCES marketplace_categories(id),
    icon VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Объявления
CREATE TABLE marketplace_listings (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id),
    category_id INT REFERENCES marketplace_categories(id),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    price DECIMAL(12,2),
    condition VARCHAR(50), -- new, used, etc
    status VARCHAR(20) DEFAULT 'active', -- active, sold, archived
    location VARCHAR(255),
    latitude DECIMAL(10,8),
    longitude DECIMAL(11,8),
    address_city VARCHAR(100),
    address_country VARCHAR(100),
    views_count INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Изображения объявлений
CREATE TABLE marketplace_images (
    id SERIAL PRIMARY KEY,
    listing_id INT REFERENCES marketplace_listings(id) ON DELETE CASCADE,
    file_path VARCHAR(255) NOT NULL,
    file_name VARCHAR(255) NOT NULL,
    file_size INT NOT NULL,
    content_type VARCHAR(100) NOT NULL,
    is_main BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Избранное
CREATE TABLE marketplace_favorites (
    user_id INT REFERENCES users(id),
    listing_id INT REFERENCES marketplace_listings(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, listing_id)
);
-- Добавляем основные категории
INSERT INTO marketplace_categories (name, slug, icon) VALUES 
('Электроника', 'electronics', 'phone'),
('Недвижимость', 'real-estate', 'home'),
('Транспорт', 'vehicles', 'car'),
('Работа', 'jobs', 'briefcase'),
('Для дома и дачи', 'home-and-garden', 'couch'),
('Личные вещи', 'personal', 'tshirt'),
('Хобби и отдых', 'hobby', 'camera'),
('Животные', 'pets', 'paw'),
('Бизнес и оборудование', 'business', 'building');

-- Добавляем подкатегории
INSERT INTO marketplace_categories (name, slug, parent_id, icon) 
SELECT 'Телефоны', 'phones', id, 'smartphone' FROM marketplace_categories WHERE slug = 'electronics';

INSERT INTO marketplace_categories (name, slug, parent_id, icon) 
SELECT 'Ноутбуки', 'laptops', id, 'laptop' FROM marketplace_categories WHERE slug = 'electronics';

INSERT INTO marketplace_categories (name, slug, parent_id, icon) 
SELECT 'Квартиры', 'apartments', id, 'apartment' FROM marketplace_categories WHERE slug = 'real-estate';

INSERT INTO marketplace_categories (name, slug, parent_id, icon) 
SELECT 'Дома', 'houses', id, 'house' FROM marketplace_categories WHERE slug = 'real-estate';
-- Добавляем тестовые объявления
INSERT INTO marketplace_listings 
(user_id, category_id, title, description, price, condition, status, location, latitude, longitude, address_city, address_country) 
VALUES
-- Телефоны
(1, (SELECT id FROM marketplace_categories WHERE slug = 'phones'), 
'iPhone 13 Pro Max', 
'Продаю iPhone 13 Pro Max 256GB в идеальном состоянии. Полный комплект, на гарантии.', 
85000, 'used', 'active', 
'Novi Sad, Serbia', 45.2671, 19.8335, 'Novi Sad', 'Serbia'),

-- Ноутбуки
(1, (SELECT id FROM marketplace_categories WHERE slug = 'laptops'),
'MacBook Pro 14" M1 Pro',
'MacBook Pro 14" (2021) с чипом M1 Pro, 16GB RAM, 512GB SSD. Состояние нового.', 
150000, 'used', 'active',
'Novi Sad, Serbia', 45.2551, 19.8452, 'Novi Sad', 'Serbia'),

-- Квартиры
(1, (SELECT id FROM marketplace_categories WHERE slug = 'apartments'),
'3-комнатная квартира в центре',
'Просторная 3-комнатная квартира в историческом центре города. Свежий ремонт, вся инфраструктура рядом.', 
15000000, 'new', 'active',
'Novi Sad, Serbia', 45.2541, 19.8401, 'Novi Sad', 'Serbia'),

-- Дома
(1, (SELECT id FROM marketplace_categories WHERE slug = 'houses'),
'Современный дом с участком',
'Новый двухэтажный дом 200м² с участком 10 соток. Все коммуникации, готов к проживанию.', 
25000000, 'new', 'active',
'Novi Sad, Serbia', 45.2460, 19.8235, 'Novi Sad', 'Serbia');

-- Добавляем изображения для объявлений
INSERT INTO marketplace_images 
(listing_id, file_path, file_name, file_size, content_type, is_main)
VALUES
-- Изображения для iPhone
(1, 'iphone13_1.jpg', 'iphone13_1.jpg', 1024, 'image/jpeg', true),
(1, 'iphone13_2.jpg', 'iphone13_2.jpg', 1024, 'image/jpeg', false),
(1, 'iphone13_3.jpg', 'iphone13_3.jpg', 1024, 'image/jpeg', false),

-- Изображения для MacBook
(2, 'macbook_1.jpg', 'macbook_1.jpg', 1024, 'image/jpeg', true),
(2, 'macbook_2.jpg', 'macbook_2.jpg', 1024, 'image/jpeg', false),
(2, 'macbook_3.jpg', 'macbook_3.jpg', 1024, 'image/jpeg', false),

-- Изображения для квартиры
(3, 'apartment_1.jpg', 'apartment_1.jpg', 1024, 'image/jpeg', true),
(3, 'apartment_2.jpg', 'apartment_2.jpg', 1024, 'image/jpeg', false),
(3, 'apartment_3.jpg', 'apartment_3.jpg', 1024, 'image/jpeg', false),
(3, 'apartment_4.jpg', 'apartment_4.jpg', 1024, 'image/jpeg', false),

-- Изображения для дома
(4, 'house_1.jpg', 'house_1.jpg', 1024, 'image/jpeg', true),
(4, 'house_2.jpg', 'house_2.jpg', 1024, 'image/jpeg', false),
(4, 'house_3.jpg', 'house_3.jpg', 1024, 'image/jpeg', false),
(4, 'house_4.jpg', 'house_4.jpg', 1024, 'image/jpeg', false);

-- Добавляем несколько записей в избранное
INSERT INTO marketplace_favorites (user_id, listing_id)
VALUES
(1, 2),
(1, 3);