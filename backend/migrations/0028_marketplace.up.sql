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
INSERT INTO marketplace_categories (name, slug, icon) VALUES 
('Транспорт', 'transport', 'car'),
('Недвижимость', 'real-estate', 'home'),
('Электроника', 'electronics', 'smartphone'),
('Одежда и обувь', 'clothing-and-shoes', 'tshirt'),
('Дом и сад', 'home-and-garden', 'couch'),
('Работа', 'jobs', 'briefcase'),
('Личные вещи', 'personal-items', 'watch'),
('Хобби и отдых', 'hobby-and-leisure', 'camera'),
('Животные', 'pets', 'paw'),
('Услуги', 'services', 'toolbox'),
('Бизнес и промышленность', 'business-and-industry', 'building');

-- Подкатегории для "Транспорт"
INSERT INTO marketplace_categories (name, slug, parent_id, icon) 
SELECT 'Автомобили', 'cars', id, 'car-side' FROM marketplace_categories WHERE slug = 'transport';

INSERT INTO marketplace_categories (name, slug, parent_id, icon) 
SELECT 'Мотоциклы', 'motorcycles', id, 'motorcycle' FROM marketplace_categories WHERE slug = 'transport';

INSERT INTO marketplace_categories (name, slug, parent_id, icon) 
SELECT 'Велосипеды', 'bicycles', id, 'bicycle' FROM marketplace_categories WHERE slug = 'transport';

INSERT INTO marketplace_categories (name, slug, parent_id, icon) 
SELECT 'Запчасти и аксессуары', 'parts-and-accessories', id, 'wrench' FROM marketplace_categories WHERE slug = 'transport';

INSERT INTO marketplace_categories (name, slug, parent_id, icon) 
SELECT 'Водный транспорт', 'water-transport', id, 'ship' FROM marketplace_categories WHERE slug = 'transport';

-- Подкатегории для "Недвижимость"
INSERT INTO marketplace_categories (name, slug, parent_id, icon) 
SELECT 'Аренда', 'rent', id, 'key' FROM marketplace_categories WHERE slug = 'real-estate';

INSERT INTO marketplace_categories (name, slug, parent_id, icon) 
SELECT 'Продажа', 'sale', id, 'apartment' FROM marketplace_categories WHERE slug = 'real-estate';

INSERT INTO marketplace_categories (name, slug, parent_id, icon) 
SELECT 'Коммерческая недвижимость', 'commercial', id, 'office-building' FROM marketplace_categories WHERE slug = 'real-estate';

INSERT INTO marketplace_categories (name, slug, parent_id, icon) 
SELECT 'Земельные участки', 'land', id, 'tree' FROM marketplace_categories WHERE slug = 'real-estate';

-- Подкатегории для "Электроника"
INSERT INTO marketplace_categories (name, slug, parent_id, icon) 
SELECT 'Смартфоны и аксессуары', 'smartphones', id, 'mobile' FROM marketplace_categories WHERE slug = 'electronics';

INSERT INTO marketplace_categories (name, slug, parent_id, icon) 
SELECT 'Компьютеры и ноутбуки', 'computers', id, 'laptop' FROM marketplace_categories WHERE slug = 'electronics';

INSERT INTO marketplace_categories (name, slug, parent_id, icon) 
SELECT 'ТВ и видео', 'tv-and-video', id, 'tv' FROM marketplace_categories WHERE slug = 'electronics';

INSERT INTO marketplace_categories (name, slug, parent_id, icon) 
SELECT 'Игровые консоли', 'gaming-consoles', id, 'gamepad' FROM marketplace_categories WHERE slug = 'electronics';

INSERT INTO marketplace_categories (name, slug, parent_id, icon) 
SELECT 'Аудио и наушники', 'audio', id, 'headphones' FROM marketplace_categories WHERE slug = 'electronics';

-- Подкатегории для "Одежда и обувь"
INSERT INTO marketplace_categories (name, slug, parent_id, icon) 
SELECT 'Мужская', 'mens-clothing', id, 'male' FROM marketplace_categories WHERE slug = 'clothing-and-shoes';

INSERT INTO marketplace_categories (name, slug, parent_id, icon) 
SELECT 'Женская', 'womens-clothing', id, 'female' FROM marketplace_categories WHERE slug = 'clothing-and-shoes';

INSERT INTO marketplace_categories (name, slug, parent_id, icon) 
SELECT 'Детская', 'kids-clothing', id, 'child' FROM marketplace_categories WHERE slug = 'clothing-and-shoes';

-- Подкатегории для "Дом и сад"
INSERT INTO marketplace_categories (name, slug, parent_id, icon) 
SELECT 'Мебель', 'furniture', id, 'chair' FROM marketplace_categories WHERE slug = 'home-and-garden';

INSERT INTO marketplace_categories (name, slug, parent_id, icon) 
SELECT 'Освещение', 'lighting', id, 'lightbulb' FROM marketplace_categories WHERE slug = 'home-and-garden';

INSERT INTO marketplace_categories (name, slug, parent_id, icon) 
SELECT 'Инструменты и ремонт', 'tools-and-repair', id, 'hammer' FROM marketplace_categories WHERE slug = 'home-and-garden';

INSERT INTO marketplace_categories (name, slug, parent_id, icon) 
SELECT 'Декор', 'decor', id, 'picture' FROM marketplace_categories WHERE slug = 'home-and-garden';
-- Добавляем тестовые объявления
-- Добавляем тестовые объявления
INSERT INTO marketplace_listings 
(user_id, category_id, title, description, price, condition, status, location, latitude, longitude, address_city, address_country) 
VALUES
-- Транспорт: Автомобили
(1, (SELECT id FROM marketplace_categories WHERE slug = 'cars'), 
'Toyota Corolla 2018', 
'Продаю Toyota Corolla 2018 года, 80,000 км пробега, отличное состояние.', 
1150000, 'used', 'active', 
'Novi Sad, Serbia', 45.2671, 19.8335, 'Novi Sad', 'Serbia'),

-- Электроника: Смартфоны
(1, (SELECT id FROM marketplace_categories WHERE slug = 'smartphones'), 
'Samsung Galaxy S21 Ultra', 
'Samsung Galaxy S21 Ultra 5G, 12GB RAM, 256GB. На гарантии, как новый.', 
90000, 'used', 'active', 
'Novi Sad, Serbia', 45.2551, 19.8452, 'Novi Sad', 'Serbia'),

-- Электроника: Компьютеры
(1, (SELECT id FROM marketplace_categories WHERE slug = 'computers'), 
'Геймерский ПК Ryzen 5', 
'Сборка: Ryzen 5 5600X, RTX 3060, 16GB RAM, 1TB SSD. Идеально для игр.', 
200000, 'used', 'active', 
'Novi Sad, Serbia', 45.2541, 19.8401, 'Novi Sad', 'Serbia'),

-- Одежда и обувь: Мужская
(1, (SELECT id FROM marketplace_categories WHERE slug = 'mens-clothing'), 
'Кожаная куртка', 
'Мужская кожаная куртка, размер L, состояние отличное.', 
10000, 'used', 'active', 
'Novi Sad, Serbia', 45.2460, 19.8235, 'Novi Sad', 'Serbia'),

-- Дом и сад: Мебель
(1, (SELECT id FROM marketplace_categories WHERE slug = 'furniture'), 
'Диван-кровать', 
'Удобный диван-кровать, трансформируется в полноценное спальное место. В хорошем состоянии.', 
25000, 'used', 'active', 
'Novi Sad, Serbia', 45.2701, 19.8500, 'Novi Sad', 'Serbia'),

-- Животные: Собаки
(1, (SELECT id FROM marketplace_categories WHERE slug = 'pets'), 
'Щенок бигля', 
'Продаются щенки бигля. Возраст 2 месяца, привиты и здоровы.', 
20000, 'new', 'active', 
'Novi Sad, Serbia', 45.2600, 19.8400, 'Novi Sad', 'Serbia'),

-- Хобби и отдых: Спорт и фитнес
(1, (SELECT id FROM marketplace_categories WHERE slug = 'hobby-and-leisure'), 
'Беговая дорожка', 
'Электрическая беговая дорожка, складная, с монитором скорости и расстояния.', 
45000, 'used', 'active', 
'Novi Sad, Serbia', 45.2450, 19.8350, 'Novi Sad', 'Serbia');

-- Добавляем изображения для объявлений
INSERT INTO marketplace_images 
(listing_id, file_path, file_name, file_size, content_type, is_main)
VALUES
-- Изображения для Toyota Corolla
(1, 'toyota_1.jpg', 'toyota_1.jpg', 1024, 'image/jpeg', true),
(1, 'toyota_2.jpg', 'toyota_2.jpg', 1024, 'image/jpeg', false),

-- Изображения для Samsung Galaxy S21 Ultra
(2, 'galaxy_s21_1.jpg', 'galaxy_s21_1.jpg', 1024, 'image/jpeg', true),
(2, 'galaxy_s21_2.jpg', 'galaxy_s21_2.jpg', 1024, 'image/jpeg', false),

-- Изображения для Геймерского ПК
(3, 'gaming_pc_1.jpg', 'gaming_pc_1.jpg', 1024, 'image/jpeg', true),
(3, 'gaming_pc_2.jpg', 'gaming_pc_2.jpg', 1024, 'image/jpeg', false),

-- Изображения для кожаной куртки
(4, 'leather_jacket_1.jpg', 'leather_jacket_1.jpg', 1024, 'image/jpeg', true),
(4, 'leather_jacket_2.jpg', 'leather_jacket_2.jpg', 1024, 'image/jpeg', false),

-- Изображения для дивана
(5, 'sofa_1.jpg', 'sofa_1.jpg', 1024, 'image/jpeg', true),
(5, 'sofa_2.jpg', 'sofa_2.jpg', 1024, 'image/jpeg', false),

-- Изображения для щенка бигля
(6, 'beagle_puppy_1.jpg', 'beagle_puppy_1.jpg', 1024, 'image/jpeg', true),
(6, 'beagle_puppy_2.jpg', 'beagle_puppy_2.jpg', 1024, 'image/jpeg', false),

-- Изображения для беговой дорожки
(7, 'treadmill_1.jpg', 'treadmill_1.jpg', 1024, 'image/jpeg', true),
(7, 'treadmill_2.jpg', 'treadmill_2.jpg', 1024, 'image/jpeg', false);


-- Добавляем несколько записей в избранное
INSERT INTO marketplace_favorites (user_id, listing_id)
VALUES
(1, 2),
(1, 3);

-- После уже существующих INSERT для категорий добавим:

-- Добавляем тестовые объявления
INSERT INTO marketplace_listings 
(user_id, category_id, title, description, price, condition, status, location, latitude, longitude, address_city, address_country) 
VALUES
-- Смартфон
(1, 
(SELECT id FROM marketplace_categories WHERE slug = 'phones'), 
'iPhone 14 Pro', 
'Продаю iPhone 14 Pro 256GB. Цвет - космический черный. На гарантии еще 9 месяцев. Полный комплект.',
95000, 
'used', 
'active', 
'Novi Sad, Serbia', 
45.2671, 
19.8335, 
'Novi Sad', 
'Serbia'),

-- Ноутбук
(1, 
(SELECT id FROM marketplace_categories WHERE slug = 'laptops'),
'MacBook Pro 16" M2',
'Новый MacBook Pro 16" с чипом M2 Pro. 32GB RAM, 1TB SSD. Максимальная комплектация.',
250000,
'new',
'active',
'Novi Sad, Serbia',
45.2551,
19.8452,
'Novi Sad',
'Serbia'),

-- Квартира
(1,
(SELECT id FROM marketplace_categories WHERE slug = 'apartments'),
'3-х комнатная квартира в центре',
'Просторная квартира с отличным ремонтом. 85м². Подземный паркинг. 2 санузла. Вид на реку.',
12500000,
'new',
'active',
'Novi Sad, Serbia',
45.2541,
19.8401,
'Novi Sad',
'Serbia');

