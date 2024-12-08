-- backend/migrations/0024_add_car.up.sql
CREATE TABLE cars (
    id SERIAL PRIMARY KEY,
    make VARCHAR(50) NOT NULL,
    model VARCHAR(50) NOT NULL,
    year INT NOT NULL,
    price_per_day NUMERIC(10, 2) NOT NULL,
    availability BOOLEAN DEFAULT TRUE,
    location VARCHAR(100) NOT NULL,
    latitude DECIMAL(10, 8),
    longitude DECIMAL(11, 8),
    description TEXT,
    seats INT DEFAULT 4,
    transmission VARCHAR(20) CHECK (transmission IN ('manual', 'automatic')),
    fuel_type VARCHAR(20) CHECK (fuel_type IN ('petrol', 'diesel', 'electric', 'hybrid')),
    daily_mileage_limit INT,
    insurance_included BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Таблица категорий
CREATE TABLE car_categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    description TEXT
);

-- Таблица особенностей
CREATE TABLE car_features (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    category VARCHAR(50),  -- Для группировки особенностей (комфорт, безопасность и т.д.)
    description TEXT
);

-- Связующая таблица автомобиль-особенности
CREATE TABLE car_feature_links (
    car_id INT REFERENCES cars(id) ON DELETE CASCADE,
    feature_id INT REFERENCES car_features(id) ON DELETE CASCADE,
    PRIMARY KEY (car_id, feature_id)
);

-- Таблица изображений
CREATE TABLE car_images (
    id SERIAL PRIMARY KEY,
    car_id INT REFERENCES cars(id) ON DELETE CASCADE,
    file_path VARCHAR(255) NOT NULL,
    file_name VARCHAR(255) NOT NULL,
    file_size INT NOT NULL,
    content_type VARCHAR(100) NOT NULL,
    is_main BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Таблица бронирований
CREATE TABLE car_bookings (
    id SERIAL PRIMARY KEY,
    car_id INT REFERENCES cars(id) ON DELETE CASCADE,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    start_date TIMESTAMP NOT NULL,
    end_date TIMESTAMP NOT NULL,
    pickup_location TEXT,
    dropoff_location TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'pending'
        CHECK (status IN ('pending', 'confirmed', 'cancelled', 'completed')),
    total_price DECIMAL(10, 2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Индексы
CREATE INDEX idx_cars_location ON cars(location);
CREATE INDEX idx_cars_availability ON cars(availability);
CREATE INDEX idx_car_bookings_dates ON car_bookings(start_date, end_date);
CREATE INDEX idx_car_bookings_status ON car_bookings(status);
CREATE INDEX idx_car_images_car_id ON car_images(car_id);
CREATE UNIQUE INDEX unique_main_image_per_car ON car_images (car_id) WHERE is_main = true;

-- Добавляем предопределенные особенности
INSERT INTO car_features (name, category) VALUES
    ('Кондиционер', 'Климат'),
    ('Климат-контроль', 'Климат'),
    ('Круиз-контроль', 'Комфорт'),
    ('Парктроники', 'Безопасность'),
    ('Камера заднего вида', 'Безопасность'),
    ('Навигация', 'Мультимедиа'),
    ('Bluetooth', 'Мультимедиа'),
    ('USB', 'Мультимедиа'),
    ('AUX', 'Мультимедиа'),
    ('MP3', 'Мультимедиа'),
    ('CD', 'Мультимедиа'),
    ('Кожаный салон', 'Интерьер'),
    ('Люк', 'Интерьер'),
    ('Панорамная крыша', 'Интерьер'),
    ('Подогрев сидений', 'Комфорт'),
    ('Электропривод сидений', 'Комфорт'),
    ('Электропривод зеркал', 'Комфорт'),
    ('Электропривод окон', 'Комфорт');

-- Добавляем категории автомобилей
INSERT INTO car_categories (name, description) VALUES
    ('Эконом', 'Бюджетные автомобили с хорошей топливной эффективностью'),
    ('Комфорт', 'Автомобили среднего класса с улучшенным комфортом'),
    ('Премиум', 'Люксовые автомобили с премиальными удобствами'),
    ('Внедорожник', 'Автомобили повышенной проходимости для путешествий'),
    ('Минивэн', 'Просторные автомобили для групповых поездок');