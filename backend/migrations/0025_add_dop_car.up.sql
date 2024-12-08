CREATE TABLE car_categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    description TEXT
);

ALTER TABLE cars
    ADD COLUMN description TEXT,
    ADD COLUMN category_id INT REFERENCES car_categories(id),
    ADD COLUMN seats INT NOT NULL DEFAULT 5,
    ADD COLUMN transmission VARCHAR(20) CHECK (transmission IN ('manual', 'automatic')),
    ADD COLUMN fuel_type VARCHAR(20) CHECK (fuel_type IN ('petrol', 'diesel', 'electric', 'hybrid')),
    ADD COLUMN features TEXT[],
    ADD COLUMN latitude DECIMAL(10, 8),
    ADD COLUMN longitude DECIMAL(11, 8),
    ADD COLUMN address TEXT,
    ADD COLUMN daily_mileage_limit INT,
    ADD COLUMN insurance_included BOOLEAN DEFAULT false;

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

-- Indexes for performance
CREATE INDEX idx_cars_category ON cars(category_id);
CREATE INDEX idx_car_bookings_dates ON car_bookings(start_date, end_date);
CREATE INDEX idx_car_bookings_status ON car_bookings(status);
CREATE INDEX idx_car_images_car_id ON car_images(car_id);
CREATE UNIQUE INDEX unique_main_image_per_car ON car_images (car_id) WHERE is_main = true;

-- Insert some default categories
INSERT INTO car_categories (name, description) VALUES
('Economy', 'Budget-friendly cars with good fuel efficiency'),
('Comfort', 'Mid-size cars with enhanced comfort features'),
('Premium', 'Luxury vehicles with premium amenities'),
('SUV', 'Sport utility vehicles for adventure and space'),
('Van', 'Spacious vehicles for group travel');