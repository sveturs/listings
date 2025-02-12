
-- Complete migration that combines all previous migrations
-- First, drop all existing tables if they exist (in correct order due to dependencies)
DROP TABLE IF EXISTS marketplace_messages CASCADE;
DROP TABLE IF EXISTS marketplace_chats CASCADE;
DROP TABLE IF EXISTS marketplace_favorites CASCADE;
DROP TABLE IF EXISTS marketplace_images CASCADE;
DROP TABLE IF EXISTS translations CASCADE;
DROP TABLE IF EXISTS notification_settings CASCADE;
DROP TABLE IF EXISTS notifications CASCADE;
DROP TABLE IF EXISTS review_votes CASCADE;
DROP TABLE IF EXISTS review_responses CASCADE;
DROP TABLE IF EXISTS reviews CASCADE;
DROP TABLE IF EXISTS marketplace_listings CASCADE;
DROP TABLE IF EXISTS marketplace_categories CASCADE;
DROP TABLE IF EXISTS user_telegram_connections CASCADE;
DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS schema_migrations CASCADE;

-- Drop existing functions
DROP FUNCTION IF EXISTS update_marketplace_chats_updated_at CASCADE;
DROP FUNCTION IF EXISTS update_updated_at_column CASCADE;
DROP FUNCTION IF EXISTS update_user_updated_at CASCADE;
DROP FUNCTION IF EXISTS update_notification_settings_updated_at CASCADE;
DROP FUNCTION IF EXISTS update_translations_updated_at CASCADE;
DROP FUNCTION IF EXISTS calculate_entity_rating CASCADE;

-- Create base tables
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(150) NOT NULL UNIQUE,
    google_id VARCHAR(255) UNIQUE,
    picture_url TEXT,
    phone VARCHAR(20),
    bio TEXT,
    notification_email BOOLEAN DEFAULT true,
    timezone VARCHAR(50) DEFAULT 'UTC',
    last_seen TIMESTAMP,
    account_status VARCHAR(20) DEFAULT 'active',
    settings JSONB DEFAULT '{}',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT users_account_status_check CHECK (account_status IN ('active', 'inactive', 'suspended'))
);

CREATE TABLE marketplace_categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(100) NOT NULL UNIQUE,
    parent_id INT REFERENCES marketplace_categories(id),
    icon VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE marketplace_listings (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id),
    category_id INT REFERENCES marketplace_categories(id),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    price DECIMAL(12,2),
    condition VARCHAR(50),
    status VARCHAR(20) DEFAULT 'active',
    location VARCHAR(255),
    latitude DECIMAL(10,8),
    longitude DECIMAL(11,8),
    address_city VARCHAR(100),
    address_country VARCHAR(100),
    views_count INT DEFAULT 0,
    show_on_map BOOLEAN NOT NULL DEFAULT true,
    original_language VARCHAR(10) DEFAULT 'sr',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

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

CREATE TABLE marketplace_favorites (
    user_id INT REFERENCES users(id),
    listing_id INT REFERENCES marketplace_listings(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, listing_id)
);

CREATE TABLE marketplace_chats (
    id SERIAL PRIMARY KEY,
    listing_id INT REFERENCES marketplace_listings(id) ON DELETE CASCADE,
    buyer_id INT REFERENCES users(id),
    seller_id INT REFERENCES users(id),
    last_message_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_archived BOOLEAN DEFAULT false,
    UNIQUE(listing_id, buyer_id, seller_id)
);

CREATE TABLE marketplace_messages (
    id SERIAL PRIMARY KEY,
    chat_id INT REFERENCES marketplace_chats(id) ON DELETE CASCADE,
    listing_id INT REFERENCES marketplace_listings(id) ON DELETE CASCADE,
    sender_id INT REFERENCES users(id),
    receiver_id INT REFERENCES users(id),
    content TEXT NOT NULL,
    is_read BOOLEAN DEFAULT false,
    original_language VARCHAR(2) DEFAULT 'en',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE reviews (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    entity_type VARCHAR(50) NOT NULL,
    entity_id INT NOT NULL,
    rating INT NOT NULL CHECK (rating >= 1 AND rating <= 5),
    comment TEXT,
    pros TEXT,
    cons TEXT,
    photos TEXT[],
    likes_count INT DEFAULT 0,
    helpful_votes INT DEFAULT 0,
    not_helpful_votes INT DEFAULT 0,
    is_verified_purchase BOOLEAN DEFAULT false,
    status VARCHAR(20) DEFAULT 'published' CHECK (status IN ('draft', 'published', 'hidden')),
    original_language VARCHAR(2) DEFAULT 'en',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE review_responses (
    id SERIAL PRIMARY KEY,
    review_id INT NOT NULL REFERENCES reviews(id) ON DELETE CASCADE,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    response TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE review_votes (
    review_id INT NOT NULL REFERENCES reviews(id) ON DELETE CASCADE,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    vote_type VARCHAR(20) NOT NULL CHECK (vote_type IN ('helpful', 'not_helpful')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (review_id, user_id)
);

CREATE TABLE translations (
    id SERIAL PRIMARY KEY,
    entity_type VARCHAR(50) NOT NULL,
    entity_id INTEGER NOT NULL,
    language VARCHAR(10) NOT NULL,
    field_name VARCHAR(50) NOT NULL,
    translated_text TEXT NOT NULL,
    is_machine_translated BOOLEAN DEFAULT true,
    is_verified BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(entity_type, entity_id, language, field_name)
);

CREATE TABLE user_telegram_connections (
    user_id INT PRIMARY KEY REFERENCES users(id),
    telegram_chat_id VARCHAR(100) NOT NULL,
    telegram_username VARCHAR(100),
    connected_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE notification_settings (
    user_id INT NOT NULL REFERENCES users(id),
    notification_type VARCHAR(50) NOT NULL,
    telegram_enabled BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, notification_type)
);

CREATE TABLE notifications (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id),
    type VARCHAR(50) NOT NULL,
    title TEXT NOT NULL,
    message TEXT NOT NULL,
    data JSONB,
    is_read BOOLEAN DEFAULT false,
    delivered_to JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE schema_migrations (
    version bigint NOT NULL PRIMARY KEY,
    dirty boolean NOT NULL
);

-- Create indexes
CREATE INDEX idx_users_phone ON users(phone);
CREATE INDEX idx_users_status ON users(account_status);

CREATE INDEX idx_marketplace_messages_chat ON marketplace_messages(chat_id);
CREATE INDEX idx_marketplace_messages_listing ON marketplace_messages(listing_id);
CREATE INDEX idx_marketplace_messages_sender ON marketplace_messages(sender_id);
CREATE INDEX idx_marketplace_messages_receiver ON marketplace_messages(receiver_id);
CREATE INDEX idx_marketplace_messages_created ON marketplace_messages(created_at);

CREATE INDEX idx_marketplace_chats_buyer ON marketplace_chats(buyer_id);
CREATE INDEX idx_marketplace_chats_seller ON marketplace_chats(seller_id);
CREATE INDEX idx_marketplace_chats_updated ON marketplace_chats(updated_at);

CREATE INDEX idx_marketplace_listings_status ON marketplace_listings(status);

CREATE INDEX idx_translations_lookup ON translations(entity_type, entity_id, language);

CREATE INDEX idx_reviews_entity ON reviews(entity_type, entity_id);
CREATE INDEX idx_reviews_user ON reviews(user_id);
CREATE INDEX idx_reviews_rating ON reviews(rating);
CREATE INDEX idx_reviews_status ON reviews(status);

CREATE INDEX idx_notifications_user ON notifications(user_id);
CREATE INDEX idx_notifications_type ON notifications(type);
CREATE INDEX idx_notifications_created ON notifications(created_at);

-- Create trigger functions
CREATE OR REPLACE FUNCTION update_marketplace_chats_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE OR REPLACE FUNCTION update_user_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE OR REPLACE FUNCTION update_notification_settings_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE OR REPLACE FUNCTION update_translations_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE OR REPLACE FUNCTION calculate_entity_rating(p_entity_type VARCHAR, p_entity_id INT)
RETURNS NUMERIC AS $$
DECLARE
    avg_rating NUMERIC;
BEGIN
    SELECT COALESCE(AVG(rating)::NUMERIC(3,2), 0)
    INTO avg_rating
    FROM reviews
    WHERE entity_type = p_entity_type 
    AND entity_id = p_entity_id 
    AND status = 'published';
    RETURN avg_rating;
END;
$$ LANGUAGE plpgsql;

-- Create triggers
CREATE TRIGGER update_marketplace_chats_timestamp
    BEFORE UPDATE ON marketplace_chats
    FOR EACH ROW
    EXECUTE FUNCTION update_marketplace_chats_updated_at();

CREATE TRIGGER update_marketplace_messages_timestamp
    BEFORE UPDATE ON marketplace_messages
    FOR EACH ROW
    EXECUTE FUNCTION update_marketplace_chats_updated_at();

CREATE TRIGGER update_reviews_updated_at
    BEFORE UPDATE ON reviews
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_review_responses_updated_at
    BEFORE UPDATE ON review_responses
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_user_updated_at();

CREATE TRIGGER update_notification_settings_timestamp
    BEFORE UPDATE ON notification_settings
    FOR EACH ROW
    EXECUTE FUNCTION update_notification_settings_updated_at();

CREATE TRIGGER update_translations_timestamp
    BEFORE UPDATE ON translations
    FOR EACH ROW
    EXECUTE FUNCTION update_translations_updated_at();

-- Insert initial data
INSERT INTO users (id, name, email, created_at, google_id, picture_url, phone, bio, notification_email, timezone, last_seen, account_status, settings, updated_at) VALUES
(1, 'Demo User', 'test@example.com', '2025-02-07 07:13:52.804162', NULL, NULL, NULL, NULL, true, 'UTC', NULL, 'active', '{}', '2025-02-07 07:13:52.887303'),
(2, 'Dmitry Voroshilov', 'voroshilovdo@gmail.com', '2025-02-07 07:14:48.873609', '102440686443518161778', 'https://lh3.googleusercontent.com/a/ACg8ocI7sRI0jzgpBiZmqITMtqT_BavTMXgiBJcZ5Qy1GlqyvN3BdlI=s96-c', NULL, NULL, true, 'UTC', NULL, 'active', '{}', '2025-02-07 07:14:48.873609'),
(3, 'Mail Box', 'boxmail386@gmail.com', '2025-02-07 07:48:08.691308', '105865404629114315097', 'https://lh3.googleusercontent.com/a/ACg8ocKlutGBMU-dy0DBOKpu78fL5o3b-PxMVaOwtlyNqEuEXNJRfg=s96-c', NULL, NULL, true, 'UTC', NULL, 'active', '{}', '2025-02-07 07:48:08.691308'),
(4, 'Azamat Salakhov', 'azamat777salakhov@gmail.com', '2025-02-08 11:50:47.665263', '112906345377302904578', 'https://lh3.googleusercontent.com/a/ACg8ocKQ_MwCSDPuSujNEC7UVdGWfEO6CG0D2RtaxLJ80JJcVrDKfA=s96-c', NULL, NULL, true, 'UTC', NULL, 'active', '{}', '2025-02-08 11:50:47.665263');
SELECT setval('users_id_seq', 3, true);

INSERT INTO user_telegram_connections (user_id, telegram_chat_id, telegram_username, connected_at) VALUES
(2, '158107689', 'Dmvool', '2025-02-07 07:34:33.511245'),
(4, '888561089', 'A777A1982', '2025-02-08 11:31:36.516123');


INSERT INTO notification_settings (user_id, notification_type, telegram_enabled, created_at, updated_at) VALUES
(2, 'new_message', true, '2025-02-07 07:14:49.238702', '2025-02-07 07:34:33.513959'),
(2, 'new_review', true, '2025-02-07 07:14:49.240252', '2025-02-07 07:34:33.519069'),
(2, 'review_vote', true, '2025-02-07 07:14:49.240911', '2025-02-07 07:34:33.519634'),
(2, 'review_response', true, '2025-02-07 07:14:49.241442', '2025-02-07 07:34:33.522901'),
(2, 'listing_status', true, '2025-02-07 07:14:49.241945', '2025-02-07 07:34:33.523478'),
(2, 'favorite_price', true, '2025-02-07 07:14:49.243803', '2025-02-07 07:34:33.523998'),
(4, 'new_message', true, '2025-02-08 11:38:09.274778', '2025-02-07 07:48:09.274778'),
(4, 'new_review', true, '2025-02-08 11:38:09.27687', '2025-02-07 07:48:09.27687'),
(4, 'review_vote', true, '2025-02-08 11:38:09.277562', '2025-02-07 07:48:09.277562'),
(4, 'review_response', true, '2025-02-08 11:38:09.278246', '2025-02-07 07:48:09.278246'),
(4, 'listing_status', true, '2025-02-08 11:38:09.279014', '2025-02-07 07:48:09.279014'),
(4, 'favorite_price', true, '2025-02-08 11:38:09.279632', '2025-02-07 07:48:09.279632');

-- Insert marketplace categories
INSERT INTO marketplace_categories (id, name, slug, parent_id, icon, created_at) VALUES
(1, 'Превоз', 'transport', NULL, 'car', '2025-02-07 07:13:52.823283'),
(2, 'Некретнине', 'real-estate', NULL, 'home', '2025-02-07 07:13:52.823283'),
(3, 'Електроника', 'electronics', NULL, 'smartphone', '2025-02-07 07:13:52.823283'),
(4, 'Одећа и обућа', 'clothing-and-shoes', NULL, 'tshirt', '2025-02-07 07:13:52.823283'),
(5, 'Кућа и башта', 'home-and-garden', NULL, 'couch', '2025-02-07 07:13:52.823283'),
(6, 'Пољопривреда', 'agriculture', NULL, 'tractor', '2025-02-07 07:13:52.823283'),
(7, 'Послови', 'jobs', NULL, 'briefcase', '2025-02-07 07:13:52.823283'),
(8, 'Лични предмети', 'personal-items', NULL, 'watch', '2025-02-07 07:13:52.823283'),
(9, 'Хоби и разонода', 'hobby-and-leisure', NULL, 'camera', '2025-02-07 07:13:52.823283'),
(10, 'Кућни љубимци', 'pets', NULL, 'paw', '2025-02-07 07:13:52.823283'),
(11, 'Услуге', 'services', NULL, 'toolbox', '2025-02-07 07:13:52.823283'),
(12, 'Бизнис и индустрија', 'business-and-industry', NULL, 'building', '2025-02-07 07:13:52.823283'),
(13, 'Аутомобили', 'cars', 1, 'car-side', '2025-02-07 07:13:52.823283'),
(14, 'Мотоцикли', 'motorcycles', 1, 'motorcycle', '2025-02-07 07:13:52.823283'),
(15, 'Електрична возила', 'electric-vehicles', 1, 'car-battery', '2025-02-07 07:13:52.823283'),
(16, 'Теретна возила', 'trucks', 1, 'truck', '2025-02-07 07:13:52.823283'),
(17, 'Делови и опрема', 'parts-and-accessories', 1, 'wrench', '2025-02-07 07:13:52.823283'),
(18, 'Електрични аутомобили', 'electric-cars', 15, 'car-electric', '2025-02-07 07:13:52.823283'),
(19, 'Електрични тротинети', 'electric-scooters', 15, 'scooter', '2025-02-07 07:13:52.823283'),
(20, 'Електрични бицикли', 'electric-bikes', 15, 'bicycle-electric', '2025-02-07 07:13:52.823283'),
(21, 'Издавање', 'rent', 2, 'key', '2025-02-07 07:13:52.823283'),
(22, 'Продаја', 'sale', 2, 'apartment', '2025-02-07 07:13:52.823283'),
(23, 'Гараже и паркинг', 'garages-and-parking', 2, 'parking', '2025-02-07 07:13:52.823283'),
(24, 'Смартфони и опрема', 'smartphones', 3, 'mobile', '2025-02-07 07:13:52.823283'),
(25, 'Рачунари и лаптопови', 'computers', 3, 'laptop', '2025-02-07 07:13:52.823283'),
(26, 'Паметни уређаји', 'smart-devices', 3, 'plug', '2025-02-07 07:13:52.823283'),
(27, 'Паметни сатови', 'smart-watches', 26, 'watch', '2025-02-07 07:13:52.823283'),
(28, 'Паметни звучници', 'smart-speakers', 26, 'speaker', '2025-02-07 07:13:52.823283'),
(29, 'Пси', 'dogs', 10, 'dog', '2025-02-07 07:13:52.823283'),
(30, 'Мачке', 'cats', 10, 'cat', '2025-02-07 07:13:52.823283'),
(31, 'Птице', 'birds', 10, 'dove', '2025-02-07 07:13:52.823283'),
(32, 'Штенци', 'puppies', 29, 'dog', '2025-02-07 07:13:52.823283'),
(33, 'Опрема за псе', 'dog-accessories', 29, 'bone', '2025-02-07 07:13:52.823283'),
(34, 'Опрема за мачке', 'cat-accessories', 30, 'paw', '2025-02-07 07:13:52.823283'),
(35, 'Пољопривредне машине', 'agricultural-machinery', 6, 'tractor', '2025-02-07 07:13:52.823283'),
(36, 'Домаће животиње', 'farm-animals', 6, 'cow', '2025-02-07 07:13:52.823283'),
(37, 'Пољопривредни производи', 'agricultural-products', 6, 'apple-alt', '2025-02-07 07:13:52.823283'),
(38, 'Трактори', 'tractors', 35, 'tractor', '2025-02-07 07:13:52.823283'),
(39, 'Комбајни', 'harvesters', 35, 'seedling', '2025-02-07 07:13:52.823283'),
(40, 'Плугови и дрљаче', 'plows-and-harrows', 35, 'tools', '2025-02-07 07:13:52.823283'),
(41, 'Сејалице', 'seeding-equipment', 35, 'corn', '2025-02-07 07:13:52.823283'),
(42, 'Опрема за наводњавање', 'irrigation-equipment', 35, 'water', '2025-02-07 07:13:52.823283'),
(43, 'Краве', 'cows', 36, 'cow', '2025-02-07 07:13:52.823283'),
(44, 'Свиње', 'pigs', 36, 'pig', '2025-02-07 07:13:52.823283'),
(45, 'Козе и овце', 'goats-and-sheep', 36, 'sheep', '2025-02-07 07:13:52.823283'),
(46, 'Живина', 'poultry', 36, 'egg', '2025-02-07 07:13:52.823283'),
(47, 'Сточна храна', 'animal-feed', 36, 'hay', '2025-02-07 07:13:52.823283'),
(48, 'Поврће', 'vegetables', 37, 'carrot', '2025-02-07 07:13:52.823283'),
(49, 'Воће', 'fruits', 37, 'apple-alt', '2025-02-07 07:13:52.823283'),
(50, 'Житарице', 'grains', 37, 'wheat', '2025-02-07 07:13:52.823283'),
(51, 'Млечни производи', 'dairy-products', 37, 'cheese', '2025-02-07 07:13:52.823283'),
(52, 'Месо и месни производи', 'meat-products', 37, 'drumstick-bite', '2025-02-07 07:13:52.823283'),
(53, 'Мед и пчеларски производи', 'honey-and-beekeeping', 37, 'honey', '2025-02-07 07:13:52.823283'),
(54, 'Кокошке', 'chickens', 46, 'chicken', '2025-02-07 07:13:52.823283'),
(55, 'Ћурке', 'turkeys', 46, 'turkey', '2025-02-07 07:13:52.823283'),
(56, 'Патке и гуске', 'ducks-and-geese', 46, 'duck', '2025-02-07 07:13:52.823283');

SELECT setval('marketplace_categories_id_seq', 56, true);

-- Insert marketplace listings
INSERT INTO marketplace_listings (id, user_id, category_id, title, description, price, condition, status, location, latitude, longitude, address_city, address_country, views_count, created_at, updated_at, show_on_map, original_language) VALUES
(8, 2, 13, 'Toyota Corolla 2018', 'Продајем Toyota Corolla 2018 годиште, 80.000 км, одлично стање. Први власник, редовно одржавање, сва документација доступна.', 1150000.00, 'used', 'active', 'Нови Сад, Србија', 45.26710000, 19.83350000, 'Нови Сад', 'Србија', 0, '2025-02-07 07:13:52.973909', '2025-02-07 07:13:52.973909', true, 'sr'),
(9, 3, 24, 'mobile Samsung Galaxy S21', 'Selling Samsung Galaxy S21, 256GB, Deep Purple. Perfect condition, complete set with original box and accessories. AppleCare+ until 2024.', 120000.00, 'used', 'active', 'Novi Sad, Serbia', 45.25510000, 19.84520000, 'Novi Sad', 'Serbia', 0, '2025-02-07 07:13:52.973909', '2025-02-07 07:13:52.973909', true, 'en'),
(10, 4, 25, 'Игровой компьютер RTX 4080', 'Продаю мощный игровой ПК: Intel Core i9-13900K, RTX 4080, 64GB RAM, 2TB NVMe SSD. Идеален для любых игр и тяжелых задач.', 350000.00, 'used', 'active', 'Нови-Сад, Сербия', 45.25410000, 19.84010000, 'Нови-Сад', 'Сербия', 0, '2025-02-07 07:13:52.973909', '2025-02-07 07:13:52.973909', true, 'ru'),
(12, 2, 13, 'автомобиль Toyota Corolla 2018', 'Продаю Toyota Corolla 2018 года, 80.000 км, отличное состояние. Первый владелец, регулярное обслуживание, вся документация в наличии.', 1475000.00, 'used', 'active', 'Косте Мајинског 4, Ветерник, Сербия', 45.24755670, 19.76878366, 'Ветерник', 'Сербия', 0, '2025-02-07 17:33:27.680035', '2025-02-07 17:40:23.957971', true, 'ru');

SELECT setval('marketplace_listings_id_seq', 12, true);

-- Insert marketplace images
INSERT INTO marketplace_images (id, listing_id, file_path, file_name, file_size, content_type, is_main, created_at) VALUES
(15, 8, 'toyota_1.jpg', 'toyota_1.jpg', 1024, 'image/jpeg', true, '2025-02-07 07:13:52.973909'),
(16, 8, 'toyota_2.jpg', 'toyota_2.jpg', 1024, 'image/jpeg', false, '2025-02-07 07:13:52.973909'),
(17, 9, 'galaxy_s21_1.jpg', 'galaxy_s21_1.jpg', 1024, 'image/jpeg', true, '2025-02-07 07:13:52.973909'),
(18, 9, 'galaxy_s21_2.jpg', 'galaxy_s21_2.jpg', 1024, 'image/jpeg', false, '2025-02-07 07:13:52.973909'),
(19, 10, 'gaming_pc_1.jpg', 'gaming_pc_1.jpg', 1024, 'image/jpeg', true, '2025-02-07 07:13:52.973909'),
(20, 10, 'gaming_pc_2.jpg', 'gaming_pc_2.jpg', 1024, 'image/jpeg', false, '2025-02-07 07:13:52.973909'),
(21, 12, 'toyota_1.jpg', 'toyota_1.jpg', 454842, 'image/jpeg', true, '2025-02-07 17:35:09.579393'),
(22, 12, 'toyota_2.jpg', 'toyota_2.jpg', 398035, 'image/jpeg', true, '2025-02-07 17:40:24.397595');

SELECT setval('marketplace_images_id_seq', 22, true);

-- Insert reviews and related data
INSERT INTO reviews (id, user_id, entity_type, entity_id, rating, comment, pros, cons, photos, likes_count, is_verified_purchase, status, created_at, updated_at, helpful_votes, not_helpful_votes, original_language) VALUES
(1, 2, 'listing', 8, 5, 'норм', NULL, NULL, NULL, 0, true, 'published', '2025-02-07 07:47:17.001726', '2025-02-07 14:25:23.586871', 0, 1, 'ru');

SELECT setval('reviews_id_seq', 1, true);

INSERT INTO review_responses (id, review_id, user_id, response, created_at, updated_at) VALUES
(1, 1, 3, 'ok', '2025-02-07 07:49:14.935308', '2025-02-07 07:49:14.935308');

SELECT setval('review_responses_id_seq', 1, true);

INSERT INTO review_votes (review_id, user_id, vote_type, created_at) VALUES
(1, 3, 'not_helpful', '2025-02-07 07:48:11.709016');

-- Insert translations
INSERT INTO translations (id, entity_type, entity_id, language, field_name, translated_text, is_machine_translated, is_verified, created_at, updated_at) VALUES
(1, 'listing', 8, 'sr', 'title', 'Toyota Corolla 2018', false, true, '2025-02-07 07:13:52.973909', '2025-02-07 07:13:52.973909'),
(2, 'listing', 8, 'sr', 'description', 'Продајем Toyota Corolla 2018 годиште, 80.000 км, одлично стање. Први власник, редовно одржавање, сва документација доступна.', false, true, '2025-02-07 07:13:52.973909', '2025-02-07 07:13:52.973909'),
(3, 'listing', 8, 'en', 'title', 'Toyota Corolla 2018', true, false, '2025-02-07 07:13:52.973909', '2025-02-07 07:13:52.973909'),
(4, 'listing', 8, 'en', 'description', 'Selling Toyota Corolla 2018, 80,000 km, excellent condition. First owner, regular maintenance, all documentation available.', true, false, '2025-02-07 07:13:52.973909', '2025-02-07 07:13:52.973909'),
(5, 'listing', 8, 'ru', 'title', 'Toyota Corolla 2018', true, false, '2025-02-07 07:13:52.973909', '2025-02-07 07:13:52.973909'),
(6, 'listing', 8, 'ru', 'description', 'Продаю Toyota Corolla 2018 года, 80.000 км, отличное состояние. Первый владелец, регулярное обслуживание, вся документация в наличии.', true, false, '2025-02-07 07:13:52.973909', '2025-02-07 07:13:52.973909'),
(7, 'listing', 9, 'en', 'title', 'mobile Samsung Galaxy S21', false, true, '2025-02-07 07:13:52.973909', '2025-02-07 07:13:52.973909'),
(8, 'listing', 9, 'en', 'description', 'Selling Samsung Galaxy S21 Ultra 5G, 12 GB RAM, 256 GB. Guaranteed, like new.', false, true, '2025-02-07 07:13:52.973909', '2025-02-07 07:13:52.973909'),
(9, 'listing', 9, 'sr', 'title', 'мобилни телефон Samsung Galaxy S21', true, false, '2025-02-07 07:13:52.973909', '2025-02-07 07:13:52.973909'),
(10, 'listing', 9, 'sr', 'description', 'Samsung Galaxy S21 Ultra 5G, 12GB RAM, 256GB. Под гаранцијом, као нов.', true, false, '2025-02-07 07:13:52.973909', '2025-02-07 07:13:52.973909'),
(11, 'listing', 9, 'ru', 'title', 'мобильник Samsung Galaxy S21', true, false, '2025-02-07 07:13:52.973909', '2025-02-07 07:13:52.973909'),
(12, 'listing', 9, 'ru', 'description', 'Продаю Samsung Galaxy S21 Ultra 5G, 12 ГБ ОЗУ, 256 ГБ. На гарантии, как новый.', true, false, '2025-02-07 07:13:52.973909', '2025-02-07 07:13:52.973909'),
(13, 'listing', 10, 'ru', 'title', 'Игровой компьютер RTX 4080', false, true, '2025-02-07 07:13:52.973909', '2025-02-07 07:13:52.973909'),
(14, 'listing', 10, 'ru', 'description', 'Продаю мощный игровой ПК: Intel Core i9-13900K, RTX 4080, 64GB RAM, 2TB NVMe SSD. Идеален для любых игр и тяжелых задач.', false, true, '2025-02-07 07:13:52.973909', '2025-02-07 07:13:52.973909'),
(15, 'listing', 10, 'sr', 'title', 'Гејмерски рачунар RTX 4080', true, false, '2025-02-07 07:13:52.973909', '2025-02-07 07:13:52.973909'),
(16, 'listing', 10, 'sr', 'description', 'Продајем моћан гејмерски рачунар: Intel Core i9-13900K, RTX 4080, 64GB RAM, 2TB NVMe SSD. Идеалан за све игре и захтевне задатке.', true, false, '2025-02-07 07:13:52.973909', '2025-02-07 07:13:52.973909'),
(17, 'listing', 10, 'en', 'title', 'Gaming PC RTX 4080', true, false, '2025-02-07 07:13:52.973909', '2025-02-07 07:13:52.973909'),
(18, 'listing', 10, 'en', 'description', 'Selling powerful gaming PC: Intel Core i9-13900K, RTX 4080, 64GB RAM, 2TB NVMe SSD. Perfect for any games and demanding tasks.', true, false, '2025-02-07 07:13:52.973909', '2025-02-07 07:13:52.973909'),
(171, 'listing', 11, 'ru', 'title', 'скамейка', false, true, '2025-02-07 07:43:15.187606', '2025-02-07 07:43:15.187606'),
(172, 'listing', 11, 'ru', 'description', 'продам шкаф.\n100х70х1000\nсамовывоз, звонить в нерабочее.', false, true, '2025-02-07 07:43:15.188679', '2025-02-07 07:43:15.188679'),
(173, 'listing', 11, 'en', 'title', 'bench', true, false, '2025-02-07 07:43:16.321152', '2025-02-07 07:43:16.321152'),
(174, 'listing', 11, 'en', 'description', 'I will sell a wardrobe.\n100x70x1000\nPick up only, call during non-working hours.', true, false, '2025-02-07 07:43:18.194747', '2025-02-07 07:43:18.194747'),
(175, 'listing', 11, 'sr', 'title', 'klupa', true, false, '2025-02-07 07:43:19.142196', '2025-02-07 07:43:19.142196'),
(176, 'listing', 11, 'sr', 'description', 'Продајем плакар.\n100х70х1000\nсамо преузимање, зовите у слободно време.', true, false, '2025-02-07 07:43:20.634479', '2025-02-07 07:43:20.634479'),
(177, 'review', 1, 'ru', 'comment', 'норм', false, true, '2025-02-07 07:47:17.001726', '2025-02-07 07:47:17.001726'),
(178, 'review', 1, 'en', 'comment', 'Okay.', true, false, '2025-02-07 07:47:17.001726', '2025-02-07 07:47:17.001726'),
(179, 'review', 1, 'sr', 'comment', 'uredu', true, false, '2025-02-07 07:47:17.001726', '2025-02-07 07:47:17.001726'),
(180, 'listing', 12, 'ru', 'title', 'автомобиль Toyota Corolla 2018', false, true, '2025-02-07 17:33:27.693733', '2025-02-07 17:33:27.693733'),
(181, 'listing', 12, 'ru', 'description', 'Продаю Toyota Corolla 2018 года, 80.000 км, отличное состояние. Первый владелец, регулярное обслуживание, вся документация в наличии.', false, true, '2025-02-07 17:33:27.694721', '2025-02-07 17:33:27.694721'),
(182, 'listing', 12, 'en', 'title', 'Toyota Corolla 2018 car', true, false, '2025-02-07 17:33:28.912973', '2025-02-07 17:33:28.912973'),
(183, 'listing', 12, 'en', 'description', 'I''m selling a Toyota Corolla 2018, 80,000 km, excellent condition. First owner, regular maintenance, all documentation available.', true, false, '2025-02-07 17:33:30.50245', '2025-02-07 17:33:30.50245'),
(184, 'listing', 12, 'sr', 'title', 'Automobil Toyota Corolla 2018.', true, false, '2025-02-07 17:33:31.762252', '2025-02-07 17:33:31.762252'),
(185, 'listing', 12, 'sr', 'description', 'Prodajem Toyota Corolla 2018. godište, 80.000 km, odlično stanje. Prvi vlasnik, redovno održavanje, sva dokumentacija prisutna.', true, false, '2025-02-07 17:33:33.866847', '2025-02-07 17:33:33.866847');

-- Set sequence value for translations
SELECT setval('translations_id_seq', 185, true);

-- Insert translations for categories (English translations)
INSERT INTO translations (entity_type, entity_id, language, field_name, translated_text, is_machine_translated, is_verified, created_at, updated_at) VALUES
('category', 1, 'en', 'name', 'Transport', true, true, '2025-02-07 07:13:52.985682', '2025-02-07 07:13:52.985682'),
('category', 2, 'en', 'name', 'Real Estate', true, true, '2025-02-07 07:13:52.985682', '2025-02-07 07:13:52.985682'),
('category', 3, 'en', 'name', 'Electronics', true, true, '2025-02-07 07:13:52.985682', '2025-02-07 07:13:52.985682'),
('category', 4, 'en', 'name', 'Clothing and Shoes', true, true, '2025-02-07 07:13:52.985682', '2025-02-07 07:13:52.985682'),
('category', 5, 'en', 'name', 'Home and Garden', true, true, '2025-02-07 07:13:52.985682', '2025-02-07 07:13:52.985682'),
('category', 6, 'en', 'name', 'Agriculture', true, true, '2025-02-07 07:13:52.985682', '2025-02-07 07:13:52.985682'),
('category', 7, 'en', 'name', 'Jobs', true, true, '2025-02-07 07:13:52.985682', '2025-02-07 07:13:52.985682'),
('category', 8, 'en', 'name', 'Personal Items', true, true, '2025-02-07 07:13:52.985682', '2025-02-07 07:13:52.985682'),
('category', 9, 'en', 'name', 'Hobbies and Leisure', true, true, '2025-02-07 07:13:52.985682', '2025-02-07 07:13:52.985682'),
('category', 10, 'en', 'name', 'Pets', true, true, '2025-02-07 07:13:52.985682', '2025-02-07 07:13:52.985682'),
('category', 11, 'en', 'name', 'Services', true, true, '2025-02-07 07:13:52.985682', '2025-02-07 07:13:52.985682'),
('category', 12, 'en', 'name', 'Business and Industry', true, true, '2025-02-07 07:13:52.985682', '2025-02-07 07:13:52.985682'),
('category', 13, 'en', 'name', 'Cars', true, true, '2025-02-07 07:13:52.985682', '2025-02-07 07:13:52.985682'),
('category', 14, 'en', 'name', 'Motorcycles', true, true, '2025-02-07 07:13:52.985682', '2025-02-07 07:13:52.985682'),
('category', 15, 'en', 'name', 'Electric Vehicles', true, true, '2025-02-07 07:13:52.985682', '2025-02-07 07:13:52.985682'),
('category', 16, 'en', 'name', 'Trucks', true, true, NOW(), NOW()),
('category', 17, 'en', 'name', 'Parts and Accessories', true, true, NOW(), NOW()),
('category', 18, 'en', 'name', 'Electric Cars', true, true, NOW(), NOW()),
('category', 19, 'en', 'name', 'Electric Scooters', true, true, NOW(), NOW()),
('category', 20, 'en', 'name', 'Electric Bikes', true, true, NOW(), NOW()),
('category', 21, 'en', 'name', 'Rent', true, true, NOW(), NOW()),
('category', 22, 'en', 'name', 'Sale', true, true, NOW(), NOW()),
('category', 23, 'en', 'name', 'Garages and Parking', true, true, NOW(), NOW()),
('category', 24, 'en', 'name', 'Smartphones and Equipment', true, true, NOW(), NOW()),
('category', 25, 'en', 'name', 'Computers and Laptops', true, true, NOW(), NOW()),
('category', 26, 'en', 'name', 'Smart Devices', true, true, NOW(), NOW()),
('category', 27, 'en', 'name', 'Smart Watches', true, true, NOW(), NOW()),
('category', 28, 'en', 'name', 'Smart Speakers', true, true, NOW(), NOW()),
('category', 29, 'en', 'name', 'Dogs', true, true, NOW(), NOW()),
('category', 30, 'en', 'name', 'Cats', true, true, NOW(), NOW()),
('category', 31, 'en', 'name', 'Birds', true, true, NOW(), NOW()),
('category', 32, 'en', 'name', 'Puppies', true, true, NOW(), NOW()),
('category', 33, 'en', 'name', 'Dog Accessories', true, true, NOW(), NOW()),
('category', 34, 'en', 'name', 'Cat Accessories', true, true, NOW(), NOW()),
('category', 35, 'en', 'name', 'Agricultural Machinery', true, true, NOW(), NOW()),
('category', 36, 'en', 'name', 'Farm Animals', true, true, NOW(), NOW()),
('category', 37, 'en', 'name', 'Agricultural Products', true, true, NOW(), NOW()),
('category', 38, 'en', 'name', 'Tractors', true, true, NOW(), NOW()),
('category', 39, 'en', 'name', 'Harvesters', true, true, NOW(), NOW()),
('category', 40, 'en', 'name', 'Plows and Harrows', true, true, NOW(), NOW()),
('category', 41, 'en', 'name', 'Seeding Equipment', true, true, NOW(), NOW()),
('category', 42, 'en', 'name', 'Irrigation Equipment', true, true, NOW(), NOW()),
('category', 43, 'en', 'name', 'Cows', true, true, NOW(), NOW()),
('category', 44, 'en', 'name', 'Pigs', true, true, NOW(), NOW()),
('category', 45, 'en', 'name', 'Goats and Sheep', true, true, NOW(), NOW()),
('category', 46, 'en', 'name', 'Poultry', true, true, NOW(), NOW()),
('category', 47, 'en', 'name', 'Animal Feed', true, true, NOW(), NOW()),
('category', 48, 'en', 'name', 'Vegetables', true, true, NOW(), NOW()),
('category', 49, 'en', 'name', 'Fruits', true, true, NOW(), NOW()),
('category', 50, 'en', 'name', 'Grains', true, true, NOW(), NOW()),
('category', 51, 'en', 'name', 'Dairy Products', true, true, NOW(), NOW()),
('category', 52, 'en', 'name', 'Meat Products', true, true, NOW(), NOW()),
('category', 53, 'en', 'name', 'Honey and Beekeeping', true, true, NOW(), NOW()),
('category', 54, 'en', 'name', 'Chickens', true, true, NOW(), NOW()),
('category', 55, 'en', 'name', 'Turkeys', true, true, NOW(), NOW()),
('category', 56, 'en', 'name', 'Ducks and Geese', true, true, NOW(), NOW()),
('category', 1, 'ru', 'name', 'Транспорт', true, true, NOW(), NOW()),
('category', 2, 'ru', 'name', 'Недвижимость', true, true, NOW(), NOW()),
('category', 3, 'ru', 'name', 'Электроника', true, true, NOW(), NOW()),
('category', 4, 'ru', 'name', 'Одежда и обувь', true, true, NOW(), NOW()),
('category', 5, 'ru', 'name', 'Дом и сад', true, true, NOW(), NOW()),
('category', 6, 'ru', 'name', 'Сельское хозяйство', true, true, NOW(), NOW()),
('category', 7, 'ru', 'name', 'Работа', true, true, NOW(), NOW()),
('category', 8, 'ru', 'name', 'Личные вещи', true, true, NOW(), NOW()),
('category', 9, 'ru', 'name', 'Хобби и отдых', true, true, NOW(), NOW()),
('category', 10, 'ru', 'name', 'Домашние животные', true, true, NOW(), NOW()),
('category', 11, 'ru', 'name', 'Услуги', true, true, NOW(), NOW()),
('category', 12, 'ru', 'name', 'Бизнес и промышленность', true, true, NOW(), NOW()),
('category', 13, 'ru', 'name', 'Автомобили', true, true, NOW(), NOW()),
('category', 14, 'ru', 'name', 'Мотоциклы', true, true, NOW(), NOW()),
('category', 15, 'ru', 'name', 'Электротранспорт', true, true, NOW(), NOW()),
('category', 16, 'ru', 'name', 'Грузовые автомобили', true, true, NOW(), NOW()),
('category', 17, 'ru', 'name', 'Запчасти и аксессуары', true, true, NOW(), NOW()),
('category', 18, 'ru', 'name', 'Электромобили', true, true, NOW(), NOW()),
('category', 19, 'ru', 'name', 'Электросамокаты', true, true, NOW(), NOW()),
('category', 20, 'ru', 'name', 'Электровелосипеды', true, true, NOW(), NOW()),
('category', 21, 'ru', 'name', 'Аренда', true, true, NOW(), NOW()),
('category', 22, 'ru', 'name', 'Продажа', true, true, NOW(), NOW()),
('category', 23, 'ru', 'name', 'Гаражи и парковки', true, true, NOW(), NOW()),
('category', 24, 'ru', 'name', 'Смартфоны и аксессуары', true, true, NOW(), NOW()),
('category', 25, 'ru', 'name', 'Компьютеры и ноутбуки', true, true, NOW(), NOW()),
('category', 26, 'ru', 'name', 'Умные устройства', true, true, NOW(), NOW()),
('category', 27, 'ru', 'name', 'Умные часы', true, true, NOW(), NOW()),
('category', 28, 'ru', 'name', 'Умные колонки', true, true, NOW(), NOW()),
('category', 29, 'ru', 'name', 'Собаки', true, true, NOW(), NOW()),
('category', 30, 'ru', 'name', 'Кошки', true, true, NOW(), NOW()),
('category', 31, 'ru', 'name', 'Птицы', true, true, NOW(), NOW()),
('category', 32, 'ru', 'name', 'Щенки', true, true, NOW(), NOW()),
('category', 33, 'ru', 'name', 'Аксессуары для собак', true, true, NOW(), NOW()),
('category', 34, 'ru', 'name', 'Аксессуары для кошек', true, true, NOW(), NOW()),
('category', 35, 'ru', 'name', 'Сельхозтехника', true, true, NOW(), NOW()),
('category', 36, 'ru', 'name', 'Сельскохозяйственные животные', true, true, NOW(), NOW()),
('category', 37, 'ru', 'name', 'Сельхозпродукция', true, true, NOW(), NOW()),
('category', 38, 'ru', 'name', 'Тракторы', true, true, NOW(), NOW()),
('category', 39, 'ru', 'name', 'Комбайны', true, true, NOW(), NOW()),
('category', 40, 'ru', 'name', 'Плуги и бороны', true, true, NOW(), NOW()),
('category', 41, 'ru', 'name', 'Сеялки', true, true, NOW(), NOW()),
('category', 42, 'ru', 'name', 'Оборудование для полива', true, true, NOW(), NOW()),
('category', 43, 'ru', 'name', 'Коровы', true, true, NOW(), NOW()),
('category', 44, 'ru', 'name', 'Свиньи', true, true, NOW(), NOW()),
('category', 45, 'ru', 'name', 'Козы и овцы', true, true, NOW(), NOW()),
('category', 46, 'ru', 'name', 'Птица', true, true, NOW(), NOW()),
('category', 47, 'ru', 'name', 'Корма', true, true, NOW(), NOW()),
('category', 48, 'ru', 'name', 'Овощи', true, true, NOW(), NOW()),
('category', 49, 'ru', 'name', 'Фрукты', true, true, NOW(), NOW()),
('category', 50, 'ru', 'name', 'Зерновые', true, true, NOW(), NOW()),
('category', 51, 'ru', 'name', 'Молочные продукты', true, true, NOW(), NOW()),
('category', 52, 'ru', 'name', 'Мясные продукты', true, true, NOW(), NOW()),
('category', 53, 'ru', 'name', 'Мёд и пчеловодство', true, true, NOW(), NOW()),
('category', 54, 'ru', 'name', 'Куры', true, true, NOW(), NOW()),
('category', 55, 'ru', 'name', 'Индейки', true, true, NOW(), NOW()),
('category', 56, 'ru', 'name', 'Утки и гуси', true, true, NOW(), NOW());


-- Set updated sequence value
SELECT setval('translations_id_seq', (SELECT MAX(id) FROM translations), true);


-- Set initial schema version
INSERT INTO schema_migrations (version, dirty) VALUES (41, false);

-- Create final indexes for optimization
CREATE INDEX IF NOT EXISTS idx_marketplace_listings_status ON marketplace_listings(status);
CREATE INDEX IF NOT EXISTS idx_translations_lookup ON translations(entity_type, entity_id, language);
CREATE INDEX IF NOT EXISTS idx_reviews_entity ON reviews(entity_type, entity_id);
CREATE INDEX IF NOT EXISTS idx_reviews_user ON reviews(user_id);
CREATE INDEX IF NOT EXISTS idx_reviews_rating ON reviews(rating);
CREATE INDEX IF NOT EXISTS idx_reviews_status ON reviews(status);
CREATE INDEX IF NOT EXISTS idx_notifications_user ON notifications(user_id);
CREATE INDEX IF NOT EXISTS idx_notifications_type ON notifications(type);
CREATE INDEX IF NOT EXISTS idx_notifications_created ON notifications(created_at);

CREATE INDEX IF NOT EXISTS idx_marketplace_categories_parent ON marketplace_categories(parent_id);