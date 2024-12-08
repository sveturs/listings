-- backend/migrations/0025_update_car_features.up.sql
CREATE TABLE IF NOT EXISTS car_features (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    category VARCHAR(50),
    description TEXT
);

CREATE TABLE IF NOT EXISTS car_feature_links (
    car_id INTEGER REFERENCES cars(id) ON DELETE CASCADE,
    feature_id INTEGER REFERENCES car_features(id) ON DELETE CASCADE,
    PRIMARY KEY(car_id, feature_id)
);

-- Предварительное добавление features
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
    ('Электропривод окон', 'Комфорт')
ON CONFLICT (name) DO NOTHING;

CREATE INDEX IF NOT EXISTS idx_car_features_category ON car_features(category);
CREATE INDEX IF NOT EXISTS idx_car_feature_links_car ON car_feature_links(car_id);
CREATE INDEX IF NOT EXISTS idx_car_feature_links_feature ON car_feature_links(feature_id);