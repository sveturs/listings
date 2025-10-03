-- Add variant attributes for photo-video category
-- This migration adds common variant attributes for Photo & Video category

-- Insert variant attribute mappings for photo-video category (ID: 1106)
-- Using unified_attributes:
-- - color (ID: 144)
-- - size (ID: 146)
-- - storage (ID: 112)

INSERT INTO variant_attribute_mappings (variant_attribute_id, category_id, sort_order, is_required)
VALUES
    (144, 1106, 1, false),  -- color
    (146, 1106, 2, false),  -- size
    (112, 1106, 3, false)   -- storage
ON CONFLICT (variant_attribute_id, category_id) DO NOTHING;
