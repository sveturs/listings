-- Откат к предыдущей версии VIEW без полей city и country

DROP VIEW IF EXISTS map_items_view CASCADE;

CREATE VIEW map_items_view AS
SELECT 
    ml.id,
    ml.title,
    ml.description,
    ml.price,
    ml.condition,
    ml.location,
    ml.latitude,
    ml.longitude,
    ml.status,
    ml.created_at,
    ml.updated_at,
    ml.user_id,
    ml.category_id,
    mc.name as category_name,
    mc.slug as category_slug,
    (
        SELECT mi.public_url 
        FROM marketplace_images mi 
        WHERE mi.listing_id = ml.id 
          AND mi.is_main = true 
        LIMIT 1
    ) as main_image_url,
    u.name as user_name,
    ml.show_on_map
FROM marketplace_listings ml
LEFT JOIN marketplace_categories mc ON ml.category_id = mc.id
LEFT JOIN users u ON ml.user_id = u.id
WHERE ml.status = 'active' 
  AND ml.show_on_map = true
  AND ml.latitude IS NOT NULL 
  AND ml.longitude IS NOT NULL;