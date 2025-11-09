-- Add display_order column to c2c_images for Phase 10.1.3
-- Allows reordering of images in listings

ALTER TABLE c2c_images
ADD COLUMN IF NOT EXISTS display_order INTEGER DEFAULT 0;

-- Create index for efficient ordering queries
CREATE INDEX IF NOT EXISTS idx_c2c_images_listing_display_order
ON c2c_images(listing_id, display_order, is_main DESC, created_at ASC);

-- Update existing rows: set display_order based on creation date
-- Main images get display_order=1, others get sequential numbers
UPDATE c2c_images
SET display_order = subquery.row_num
FROM (
    SELECT
        id,
        ROW_NUMBER() OVER (
            PARTITION BY listing_id
            ORDER BY is_main DESC NULLS LAST, created_at ASC
        ) as row_num
    FROM c2c_images
) AS subquery
WHERE c2c_images.id = subquery.id;

COMMENT ON COLUMN c2c_images.display_order IS 'Order for displaying images, allows manual reordering (Phase 10.1.3)';
