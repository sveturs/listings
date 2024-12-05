-- backend/migrations/0022_update_bed_images.up.sql
ALTER TABLE bed_images
ADD COLUMN is_main BOOLEAN DEFAULT false;

CREATE UNIQUE INDEX unique_main_image_per_bed ON bed_images (bed_id) WHERE is_main = true;