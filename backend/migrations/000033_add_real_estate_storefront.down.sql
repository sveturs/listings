-- Удаляем переводы для товаров
DELETE FROM translations 
WHERE entity_type = 'storefront_product' 
  AND entity_id IN (
    SELECT sp.id 
    FROM storefront_products sp
    JOIN storefronts s ON sp.storefront_id = s.id
    WHERE s.slug = 'belgrade-real-estate'
  );

-- Удаляем переводы для витрины
DELETE FROM translations 
WHERE entity_type = 'storefront' 
  AND entity_id IN (
    SELECT id FROM storefronts WHERE slug = 'belgrade-real-estate'
  );

-- Удаляем товары витрины
DELETE FROM storefront_products 
WHERE storefront_id IN (
    SELECT id FROM storefronts WHERE slug = 'belgrade-real-estate'
);

-- Удаляем витрину
DELETE FROM storefronts WHERE slug = 'belgrade-real-estate';