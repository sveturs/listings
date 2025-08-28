-- Rollback migration: Remove additional demo listings
DELETE FROM marketplace_images WHERE listing_id IN (
  SELECT id FROM marketplace_listings 
  WHERE user_id = 7
  AND created_at > NOW() - INTERVAL '20 days'
);

DELETE FROM marketplace_listings 
WHERE user_id = 7
AND created_at > NOW() - INTERVAL '20 days';