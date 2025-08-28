-- Rollback migration: Remove demo storefronts and products
-- First remove product images
DELETE FROM storefront_product_images WHERE product_id IN (
  SELECT sp.id FROM storefront_products sp
  JOIN storefronts s ON sp.storefront_id = s.id
  WHERE s.slug IN (
    'technova-electronics',
    'fashion-house-belgrade',
    'home-garden-center',
    'agroshop-supplies'
  )
);

-- Remove products
DELETE FROM storefront_products WHERE storefront_id IN (
  SELECT id FROM storefronts WHERE slug IN (
    'technova-electronics',
    'fashion-house-belgrade',
    'home-garden-center',
    'agroshop-supplies'
  )
);

-- Remove storefronts
DELETE FROM storefronts WHERE slug IN (
  'technova-electronics',
  'fashion-house-belgrade',
  'home-garden-center',
  'agroshop-supplies'
);