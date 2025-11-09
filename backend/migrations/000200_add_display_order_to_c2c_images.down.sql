-- Rollback display_order column addition

DROP INDEX IF EXISTS idx_c2c_images_listing_display_order;

ALTER TABLE c2c_images
DROP COLUMN IF EXISTS display_order;
