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
-- Главные категории
INSERT INTO marketplace_categories (name, slug, icon) VALUES 
('Транспорт', 'transport', 'car'),
('Недвижимость', 'real-estate', 'home'),
('Электроника', 'electronics', 'smartphone'),
('Одежда и обувь', 'clothing-and-shoes', 'tshirt'),
('Дом и сад', 'home-and-garden', 'couch'),
('Сельское хозяйство', 'agriculture', 'tractor'),
('Работа', 'jobs', 'briefcase'),
('Личные вещи', 'personal-items', 'watch'),
('Хобби и отдых', 'hobby-and-leisure', 'camera'),
('Домашние животные', 'pets', 'paw'),
('Услуги', 'services', 'toolbox'),
('Бизнес и промышленность', 'business-and-industry', 'building');

-- Подкатегории для "Транспорт"
INSERT INTO marketplace_categories (name, slug, parent_id, icon) 
SELECT 'Автомобили', 'cars', id, 'car-side' FROM marketplace_categories WHERE slug = 'transport';

INSERT INTO marketplace_categories (name, slug, parent_id, icon) 
SELECT 'Мотоциклы', 'motorcycles', id, 'motorcycle' FROM marketplace_categories WHERE slug = 'transport';

INSERT INTO marketplace_categories (name, slug, parent_id, icon) 
SELECT 'Электротранспорт', 'electric-vehicles', id, 'car-battery' FROM marketplace_categories WHERE slug = 'transport';

INSERT INTO marketplace_categories (name, slug, parent_id, icon) 
SELECT 'Грузовые автомобили', 'trucks', id, 'truck' FROM marketplace_categories WHERE slug = 'transport';

INSERT INTO marketplace_categories (name, slug, parent_id, icon) 
SELECT 'Запчасти и аксессуары', 'parts-and-accessories', id, 'wrench' FROM marketplace_categories WHERE slug = 'transport';

-- Подкатегории для "Электротранспорт"
INSERT INTO marketplace_categories (name, slug, parent_id, icon) 
SELECT 'Электромобили', 'electric-cars', id, 'car-electric' FROM marketplace_categories WHERE slug = 'electric-vehicles';

INSERT INTO marketplace_categories (name, slug, parent_id, icon) 
SELECT 'Электросамокаты', 'electric-scooters', id, 'scooter' FROM marketplace_categories WHERE slug = 'electric-vehicles';

INSERT INTO marketplace_categories (name, slug, parent_id, icon) 
SELECT 'Электровелосипеды', 'electric-bikes', id, 'bicycle-electric' FROM marketplace_categories WHERE slug = 'electric-vehicles';

INSERT INTO marketplace_categories (name, slug, parent_id, icon) 
SELECT 'Аксессуары для электротранспорта', 'electric-vehicle-accessories', id, 'plug' FROM marketplace_categories WHERE slug = 'electric-vehicles';

-- Подкатегории для "Недвижимость"
INSERT INTO marketplace_categories (name, slug, parent_id, icon) 
SELECT 'Аренда', 'rent', id, 'key' FROM marketplace_categories WHERE slug = 'real-estate';

INSERT INTO marketplace_categories (name, slug, parent_id, icon) 
SELECT 'Продажа', 'sale', id, 'apartment' FROM marketplace_categories WHERE slug = 'real-estate';

INSERT INTO marketplace_categories (name, slug, parent_id, icon) 
SELECT 'Гаражи и парковки', 'garages-and-parking', id, 'parking' FROM marketplace_categories WHERE slug = 'real-estate';

-- Подкатегории для "Электроника"
INSERT INTO marketplace_categories (name, slug, parent_id, icon) 
SELECT 'Смартфоны и аксессуары', 'smartphones', id, 'mobile' FROM marketplace_categories WHERE slug = 'electronics';

INSERT INTO marketplace_categories (name, slug, parent_id, icon) 
SELECT 'Компьютеры и ноутбуки', 'computers', id, 'laptop' FROM marketplace_categories WHERE slug = 'electronics';

INSERT INTO marketplace_categories (name, slug, parent_id, icon) 
SELECT 'Умные устройства', 'smart-devices', id, 'plug' FROM marketplace_categories WHERE slug = 'electronics';

-- Подкатегории для "Умные устройства"
INSERT INTO marketplace_categories (name, slug, parent_id, icon) 
SELECT 'Умные часы', 'smart-watches', id, 'watch' FROM marketplace_categories WHERE slug = 'smart-devices';

INSERT INTO marketplace_categories (name, slug, parent_id, icon) 
SELECT 'Умные колонки', 'smart-speakers', id, 'speaker' FROM marketplace_categories WHERE slug = 'smart-devices';

-- Подкатегории для "Животные"
INSERT INTO marketplace_categories (name, slug, parent_id, icon) 
SELECT 'Собаки', 'dogs', id, 'dog' FROM marketplace_categories WHERE slug = 'pets';

INSERT INTO marketplace_categories (name, slug, parent_id, icon) 
SELECT 'Кошки', 'cats', id, 'cat' FROM marketplace_categories WHERE slug = 'pets';

INSERT INTO marketplace_categories (name, slug, parent_id, icon) 
SELECT 'Птицы', 'birds', id, 'dove' FROM marketplace_categories WHERE slug = 'pets';

-- Подкатегории для "Собаки"
INSERT INTO marketplace_categories (name, slug, parent_id, icon) 
SELECT 'Щенки', 'puppies', id, 'dog' FROM marketplace_categories WHERE slug = 'dogs';

INSERT INTO marketplace_categories (name, slug, parent_id, icon) 
SELECT 'Аксессуары для собак', 'dog-accessories', id, 'bone' FROM marketplace_categories WHERE slug = 'dogs';

-- Подкатегории для "Кошки"
INSERT INTO marketplace_categories (name, slug, parent_id, icon) 
SELECT 'Аксессуары для кошек', 'cat-accessories', id, 'paw' FROM marketplace_categories WHERE slug = 'cats';




-- Подкатегории для "Сельское хозяйство"
INSERT INTO marketplace_categories (name, slug, parent_id, icon) 
SELECT 'Сельхозтехника', 'agricultural-machinery', id, 'tractor' FROM marketplace_categories WHERE slug = 'agriculture';

INSERT INTO marketplace_categories (name, slug, parent_id, icon) 
SELECT 'Сельхозживотные', 'farm-animals', id, 'cow' FROM marketplace_categories WHERE slug = 'agriculture';

INSERT INTO marketplace_categories (name, slug, parent_id, icon) 
SELECT 'Продукты сельского хозяйства', 'agricultural-products', id, 'apple-alt' FROM marketplace_categories WHERE slug = 'agriculture';

-- Подкатегории для "Сельхозтехника"
INSERT INTO marketplace_categories (name, slug, parent_id, icon) 
SELECT 'Тракторы', 'tractors', id, 'tractor' FROM marketplace_categories WHERE slug = 'agricultural-machinery';

INSERT INTO marketplace_categories (name, slug, parent_id, icon) 
SELECT 'Комбайны', 'harvesters', id, 'seedling' FROM marketplace_categories WHERE slug = 'agricultural-machinery';

INSERT INTO marketplace_categories (name, slug, parent_id, icon) 
SELECT 'Плуги и бороны', 'plows-and-harrows', id, 'tools' FROM marketplace_categories WHERE slug = 'agricultural-machinery';

INSERT INTO marketplace_categories (name, slug, parent_id, icon) 
SELECT 'Посевная техника', 'seeding-equipment', id, 'corn' FROM marketplace_categories WHERE slug = 'agricultural-machinery';

INSERT INTO marketplace_categories (name, slug, parent_id, icon) 
SELECT 'Оборудование для орошения', 'irrigation-equipment', id, 'water' FROM marketplace_categories WHERE slug = 'agricultural-machinery';

-- Подкатегории для "Сельхозживотные"
INSERT INTO marketplace_categories (name, slug, parent_id, icon) 
SELECT 'Коровы', 'cows', id, 'cow' FROM marketplace_categories WHERE slug = 'farm-animals';

INSERT INTO marketplace_categories (name, slug, parent_id, icon) 
SELECT 'Свиньи', 'pigs', id, 'pig' FROM marketplace_categories WHERE slug = 'farm-animals';

INSERT INTO marketplace_categories (name, slug, parent_id, icon) 
SELECT 'Козы и овцы', 'goats-and-sheep', id, 'sheep' FROM marketplace_categories WHERE slug = 'farm-animals';

INSERT INTO marketplace_categories (name, slug, parent_id, icon) 
SELECT 'Птицы', 'poultry', id, 'egg' FROM marketplace_categories WHERE slug = 'farm-animals';

INSERT INTO marketplace_categories (name, slug, parent_id, icon) 
SELECT 'Корма для животных', 'animal-feed', id, 'hay' FROM marketplace_categories WHERE slug = 'farm-animals';

-- Подкатегории для "Продукты сельского хозяйства"
INSERT INTO marketplace_categories (name, slug, parent_id, icon) 
SELECT 'Овощи', 'vegetables', id, 'carrot' FROM marketplace_categories WHERE slug = 'agricultural-products';

INSERT INTO marketplace_categories (name, slug, parent_id, icon) 
SELECT 'Фрукты', 'fruits', id, 'apple-alt' FROM marketplace_categories WHERE slug = 'agricultural-products';

INSERT INTO marketplace_categories (name, slug, parent_id, icon) 
SELECT 'Зерновые культуры', 'grains', id, 'wheat' FROM marketplace_categories WHERE slug = 'agricultural-products';

INSERT INTO marketplace_categories (name, slug, parent_id, icon) 
SELECT 'Молочная продукция', 'dairy-products', id, 'cheese' FROM marketplace_categories WHERE slug = 'agricultural-products';

INSERT INTO marketplace_categories (name, slug, parent_id, icon) 
SELECT 'Мясо и мясные продукты', 'meat-products', id, 'drumstick-bite' FROM marketplace_categories WHERE slug = 'agricultural-products';

INSERT INTO marketplace_categories (name, slug, parent_id, icon) 
SELECT 'Мёд и продукты пчеловодства', 'honey-and-beekeeping', id, 'honey' FROM marketplace_categories WHERE slug = 'agricultural-products';

-- Подкатегории для "Птицы" в сельхозживотных
INSERT INTO marketplace_categories (name, slug, parent_id, icon) 
SELECT 'Куры', 'chickens', id, 'chicken' FROM marketplace_categories WHERE slug = 'poultry';

INSERT INTO marketplace_categories (name, slug, parent_id, icon) 
SELECT 'Индейки', 'turkeys', id, 'turkey' FROM marketplace_categories WHERE slug = 'poultry';

INSERT INTO marketplace_categories (name, slug, parent_id, icon) 
SELECT 'Утки и гуси', 'ducks-and-geese', id, 'duck' FROM marketplace_categories WHERE slug = 'poultry';

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


