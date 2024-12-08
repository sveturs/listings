-- backend/migrations/0026_add_car_images.up.sql
CREATE TABLE car_images (
    id SERIAL PRIMARY KEY,
    car_id INT NOT NULL REFERENCES cars(id) ON DELETE CASCADE,
    file_path VARCHAR(255) NOT NULL,
    file_name VARCHAR(255) NOT NULL,
    file_size INT NOT NULL,
    content_type VARCHAR(100) NOT NULL,
    is_main BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_car_images_car_id ON car_images(car_id);
CREATE UNIQUE INDEX unique_main_image_per_car ON car_images (car_id) WHERE is_main = true;