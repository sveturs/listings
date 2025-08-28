-- Откатываем изменения
UPDATE marketplace_listings ml
SET storefront_id = NULL
FROM storefront_products sp
WHERE ml.id = sp.id
  AND ml.storefront_id = sp.storefront_id;

-- Удаляем комментарий
COMMENT ON COLUMN marketplace_listings.storefront_id IS NULL;