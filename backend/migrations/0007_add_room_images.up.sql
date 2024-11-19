-- 0007_add_room_images.up.sql
CREATE TABLE room_images (
    id SERIAL PRIMARY KEY,
    room_id INT NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
    file_path VARCHAR(255) NOT NULL,
    file_name VARCHAR(255) NOT NULL,
    file_size INT NOT NULL,
    content_type VARCHAR(100) NOT NULL,
    is_main BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_main_image_per_room UNIQUE (room_id, is_main) 
        DEFERRABLE INITIALLY DEFERRED
);

CREATE INDEX idx_room_images_room_id ON room_images(room_id);