-- backend/migrations/0022_update_bed_images.down.sql
DROP INDEX IF EXISTS unique_main_image_per_bed;
ALTER TABLE bed_images DROP COLUMN IF EXISTS is_main;