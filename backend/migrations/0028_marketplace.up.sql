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
('Превоз', 'transport', 'car'),
('Некретнине', 'real-estate', 'home'),
('Електроника', 'electronics', 'smartphone'),
('Одећа и обућа', 'clothing-and-shoes', 'tshirt'),
('Кућа и башта', 'home-and-garden', 'couch'),
('Пољопривреда', 'agriculture', 'tractor'),
('Послови', 'jobs', 'briefcase'),
('Лични предмети', 'personal-items', 'watch'),
('Хоби и разонода', 'hobby-and-leisure', 'camera'),
('Кућни љубимци', 'pets', 'paw'),
('Услуге', 'services', 'toolbox'),
('Бизнис и индустрија', 'business-and-industry', 'building');

-- Подкатегории для "Превоз"
INSERT INTO marketplace_categories (name, slug, parent_id, icon) 
SELECT 'Аутомобили', 'cars', id, 'car-side' FROM marketplace_categories WHERE slug = 'transport';

INSERT INTO marketplace_categories (name, slug, parent_id, icon) 
SELECT 'Мотоцикли', 'motorcycles', id, 'motorcycle' FROM marketplace_categories WHERE slug = 'transport';

INSERT INTO marketplace_categories (name, slug, parent_id, icon) 
SELECT 'Електрична возила', 'electric-vehicles', id, 'car-battery' FROM marketplace_categories WHERE slug = 'transport';

INSERT INTO marketplace_categories (name, slug, parent_id, icon) 
SELECT 'Теретна возила', 'trucks', id, 'truck' FROM marketplace_categories WHERE slug = 'transport';

INSERT INTO marketplace_categories (name, slug, parent_id, icon) 
SELECT 'Делови и опрема', 'parts-and-accessories', id, 'wrench' FROM marketplace_categories WHERE slug = 'transport';

-- Подкатегории для "Електрична возила"
INSERT INTO marketplace_categories (name, slug, parent_id, icon) 
SELECT 'Електрични аутомобили', 'electric-cars', id, 'car-electric' FROM marketplace_categories WHERE slug = 'electric-vehicles';

INSERT INTO marketplace_categories (name, slug, parent_id, icon) 
SELECT 'Електрични тротинети', 'electric-scooters', id, 'scooter' FROM marketplace_categories WHERE slug = 'electric-vehicles';

INSERT INTO marketplace_categories (name, slug, parent_id, icon) 
SELECT 'Електрични бицикли', 'electric-bikes', id, 'bicycle-electric' FROM marketplace_categories WHERE slug = 'electric-vehicles';

-- Подкатегории для "Некретнине"
INSERT INTO marketplace_categories (name, slug, parent_id, icon) 
SELECT 'Издавање', 'rent', id, 'key' FROM marketplace_categories WHERE slug = 'real-estate';

INSERT INTO marketplace_categories (name, slug, parent_id, icon) 
SELECT 'Продаја', 'sale', id, 'apartment' FROM marketplace_categories WHERE slug = 'real-estate';

INSERT INTO marketplace_categories (name, slug, parent_id, icon) 
SELECT 'Гараже и паркинг', 'garages-and-parking', id, 'parking' FROM marketplace_categories WHERE slug = 'real-estate';
-- Подкатегории для "Електроника"
INSERT INTO marketplace_categories (name, slug, parent_id, icon) 
SELECT 'Смартфони и опрема', 'smartphones', id, 'mobile' FROM marketplace_categories WHERE slug = 'electronics';

INSERT INTO marketplace_categories (name, slug, parent_id, icon) 
SELECT 'Рачунари и лаптопови', 'computers', id, 'laptop' FROM marketplace_categories WHERE slug = 'electronics';

INSERT INTO marketplace_categories (name, slug, parent_id, icon) 
SELECT 'Паметни уређаји', 'smart-devices', id, 'plug' FROM marketplace_categories WHERE slug = 'electronics';

-- Подкатегории для "Паметни уређаји"
INSERT INTO marketplace_categories (name, slug, parent_id, icon) 
SELECT 'Паметни сатови', 'smart-watches', id, 'watch' FROM marketplace_categories WHERE slug = 'smart-devices';

INSERT INTO marketplace_categories (name, slug, parent_id, icon) 
SELECT 'Паметни звучници', 'smart-speakers', id, 'speaker' FROM marketplace_categories WHERE slug = 'smart-devices';

-- Подкатегории для "Кућни љубимци"
INSERT INTO marketplace_categories (name, slug, parent_id, icon) 
SELECT 'Пси', 'dogs', id, 'dog' FROM marketplace_categories WHERE slug = 'pets';

INSERT INTO marketplace_categories (name, slug, parent_id, icon) 
SELECT 'Мачке', 'cats', id, 'cat' FROM marketplace_categories WHERE slug = 'pets';

INSERT INTO marketplace_categories (name, slug, parent_id, icon) 
SELECT 'Птице', 'birds', id, 'dove' FROM marketplace_categories WHERE slug = 'pets';

-- Подкатегории для "Пси"
INSERT INTO marketplace_categories (name, slug, parent_id, icon) 
SELECT 'Штенци', 'puppies', id, 'dog' FROM marketplace_categories WHERE slug = 'dogs';

INSERT INTO marketplace_categories (name, slug, parent_id, icon) 
SELECT 'Опрема за псе', 'dog-accessories', id, 'bone' FROM marketplace_categories WHERE slug = 'dogs';

-- Подкатегории для "Мачке"
INSERT INTO marketplace_categories (name, slug, parent_id, icon) 
SELECT 'Опрема за мачке', 'cat-accessories', id, 'paw' FROM marketplace_categories WHERE slug = 'cats';

-- Подкатегории для "Пољопривреда"
INSERT INTO marketplace_categories (name, slug, parent_id, icon) 
SELECT 'Пољопривредне машине', 'agricultural-machinery', id, 'tractor' FROM marketplace_categories WHERE slug = 'agriculture';

INSERT INTO marketplace_categories (name, slug, parent_id, icon) 
SELECT 'Домаће животиње', 'farm-animals', id, 'cow' FROM marketplace_categories WHERE slug = 'agriculture';

INSERT INTO marketplace_categories (name, slug, parent_id, icon) 
SELECT 'Пољопривредни производи', 'agricultural-products', id, 'apple-alt' FROM marketplace_categories WHERE slug = 'agriculture';

-- Подкатегории для "Пољопривредне машине"
INSERT INTO marketplace_categories (name, slug, parent_id, icon) 
SELECT 'Трактори', 'tractors', id, 'tractor' FROM marketplace_categories WHERE slug = 'agricultural-machinery';

INSERT INTO marketplace_categories (name, slug, parent_id, icon) 
SELECT 'Комбајни', 'harvesters', id, 'seedling' FROM marketplace_categories WHERE slug = 'agricultural-machinery';

INSERT INTO marketplace_categories (name, slug, parent_id, icon) 
SELECT 'Плугови и дрљаче', 'plows-and-harrows', id, 'tools' FROM marketplace_categories WHERE slug = 'agricultural-machinery';

INSERT INTO marketplace_categories (name, slug, parent_id, icon) 
SELECT 'Сејалице', 'seeding-equipment', id, 'corn' FROM marketplace_categories WHERE slug = 'agricultural-machinery';

INSERT INTO marketplace_categories (name, slug, parent_id, icon) 
SELECT 'Опрема за наводњавање', 'irrigation-equipment', id, 'water' FROM marketplace_categories WHERE slug = 'agricultural-machinery';

-- Подкатегории для "Домаће животиње"
INSERT INTO marketplace_categories (name, slug, parent_id, icon) 
SELECT 'Краве', 'cows', id, 'cow' FROM marketplace_categories WHERE slug = 'farm-animals';

INSERT INTO marketplace_categories (name, slug, parent_id, icon) 
SELECT 'Свиње', 'pigs', id, 'pig' FROM marketplace_categories WHERE slug = 'farm-animals';

INSERT INTO marketplace_categories (name, slug, parent_id, icon) 
SELECT 'Козе и овце', 'goats-and-sheep', id, 'sheep' FROM marketplace_categories WHERE slug = 'farm-animals';

INSERT INTO marketplace_categories (name, slug, parent_id, icon) 
SELECT 'Живина', 'poultry', id, 'egg' FROM marketplace_categories WHERE slug = 'farm-animals';

INSERT INTO marketplace_categories (name, slug, parent_id, icon) 
SELECT 'Сточна храна', 'animal-feed', id, 'hay' FROM marketplace_categories WHERE slug = 'farm-animals';

-- Подкатегории для "Пољопривредни производи"
INSERT INTO marketplace_categories (name, slug, parent_id, icon) 
SELECT 'Поврће', 'vegetables', id, 'carrot' FROM marketplace_categories WHERE slug = 'agricultural-products';

INSERT INTO marketplace_categories (name, slug, parent_id, icon) 
SELECT 'Воће', 'fruits', id, 'apple-alt' FROM marketplace_categories WHERE slug = 'agricultural-products';

INSERT INTO marketplace_categories (name, slug, parent_id, icon) 
SELECT 'Житарице', 'grains', id, 'wheat' FROM marketplace_categories WHERE slug = 'agricultural-products';

INSERT INTO marketplace_categories (name, slug, parent_id, icon) 
SELECT 'Млечни производи', 'dairy-products', id, 'cheese' FROM marketplace_categories WHERE slug = 'agricultural-products';

INSERT INTO marketplace_categories (name, slug, parent_id, icon) 
SELECT 'Месо и месни производи', 'meat-products', id, 'drumstick-bite' FROM marketplace_categories WHERE slug = 'agricultural-products';

INSERT INTO marketplace_categories (name, slug, parent_id, icon) 
SELECT 'Мед и пчеларски производи', 'honey-and-beekeeping', id, 'honey' FROM marketplace_categories WHERE slug = 'agricultural-products';

-- Подкатегории для "Живина"
INSERT INTO marketplace_categories (name, slug, parent_id, icon) 
SELECT 'Кокошке', 'chickens', id, 'chicken' FROM marketplace_categories WHERE slug = 'poultry';

INSERT INTO marketplace_categories (name, slug, parent_id, icon) 
SELECT 'Ћурке', 'turkeys', id, 'turkey' FROM marketplace_categories WHERE slug = 'poultry';

INSERT INTO marketplace_categories (name, slug, parent_id, icon) 
SELECT 'Патке и гуске', 'ducks-and-geese', id, 'duck' FROM marketplace_categories WHERE slug = 'poultry';

-- Добавляем тестовые объявления
INSERT INTO marketplace_listings 
(user_id, category_id, title, description, price, condition, status, location, latitude, longitude, address_city, address_country) 
VALUES
-- Транспорт: Автомобили
(1, (SELECT id FROM marketplace_categories WHERE slug = 'cars'), 
'Toyota Corolla 2018', 
'Продајем Toyota Corolla 2018 годиште, 80.000 км, одлично стање.', 
1150000, 'used', 'active', 
'Нови Сад, Србија', 45.2671, 19.8335, 'Нови Сад', 'Србија'),

-- Электроника: Смартфоны
(1, (SELECT id FROM marketplace_categories WHERE slug = 'smartphones'), 
'Samsung Galaxy S21 Ultra', 
'Samsung Galaxy S21 Ultra 5G, 12GB RAM, 256GB. Под гаранцијом, као нов.', 
90000, 'used', 'active', 
'Нови Сад, Србија', 45.2551, 19.8452, 'Нови Сад', 'Србија'),

-- Электроника: Компьютеры
(1, (SELECT id FROM marketplace_categories WHERE slug = 'computers'), 
'Гејмерски рачунар Ryzen 5', 
'Конфигурација: Ryzen 5 5600X, RTX 3060, 16GB RAM, 1TB SSD. Идеално за игрице.', 
200000, 'used', 'active', 
'Нови Сад, Србија', 45.2541, 19.8401, 'Нови Сад', 'Србија'),

-- Одежда и обувь
(1, (SELECT id FROM marketplace_categories WHERE slug = 'clothing-and-shoes'), 
'Кожна јакна', 
'Мушка кожна јакна, величина L, одлично стање.', 
10000, 'used', 'active', 
'Нови Сад, Србија', 45.2460, 19.8235, 'Нови Сад', 'Србија'),

-- Дом и сад
(1, (SELECT id FROM marketplace_categories WHERE slug = 'home-and-garden'), 
'Кауч на развлачење', 
'Удобан кауч на развлачење, претвара се у кревет. У добром стању.', 
25000, 'used', 'active', 
'Нови Сад, Србија', 45.2701, 19.8500, 'Нови Сад', 'Србија'),

-- Животные
(1, (SELECT id FROM marketplace_categories WHERE slug = 'pets'), 
'Штене бигла', 
'Продајем штенце бигла. Старост 2 месеца, вакцинисани и здрави.', 
20000, 'new', 'active', 
'Нови Сад, Србија', 45.2600, 19.8400, 'Нови Сад', 'Србија'),

-- Хобби и отдых
(1, (SELECT id FROM marketplace_categories WHERE slug = 'hobby-and-leisure'), 
'Трака за трчање', 
'Електрична трака за трчање, склопива, са дисплејом за брзину и раздаљину.', 
45000, 'used', 'active', 
'Нови Сад, Србија', 45.2450, 19.8350, 'Нови Сад', 'Србија');

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