-- Rollback migration: Remove demo listings
DELETE FROM marketplace_images WHERE listing_id IN (
  SELECT id FROM marketplace_listings WHERE user_id = 7
);

DELETE FROM translations WHERE entity_type = 'marketplace_listing' AND entity_id IN (
  SELECT id FROM marketplace_listings WHERE user_id = 7
);

DELETE FROM marketplace_listings WHERE user_id = 7;

-- Keep user 7 for future use