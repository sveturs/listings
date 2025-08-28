-- Сначала создаем недостающие записи в user_storefronts на основе storefronts
INSERT INTO user_storefronts (id, user_id, name, description, slug, status, created_at, updated_at)
SELECT s.id, s.user_id, s.name, s.description, s.slug, 'active', s.created_at, s.updated_at
FROM storefronts s
WHERE NOT EXISTS (SELECT 1 FROM user_storefronts us WHERE us.id = s.id);

-- Обновляем поле storefront_id в marketplace_listings для товаров витрин
UPDATE marketplace_listings ml
SET storefront_id = sp.storefront_id
FROM storefront_products sp
WHERE ml.id = sp.id
  AND ml.storefront_id IS NULL
  AND sp.storefront_id IS NOT NULL
  AND EXISTS (SELECT 1 FROM user_storefronts us WHERE us.id = sp.storefront_id);

-- Добавляем комментарий для понимания структуры
COMMENT ON COLUMN marketplace_listings.storefront_id IS 'ID витрины, если это товар витрины. Используется для определения типа товара и правильной навигации';