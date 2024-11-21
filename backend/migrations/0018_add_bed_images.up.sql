CREATE TABLE bed_images (
    id SERIAL PRIMARY KEY,
    bed_id INT NOT NULL REFERENCES beds(id) ON DELETE CASCADE,
    file_path VARCHAR(255) NOT NULL,
    file_name VARCHAR(255) NOT NULL,
    file_size INT NOT NULL,
    content_type VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_bed_images_bed_id ON bed_images(bed_id);