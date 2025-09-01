-- Возвращаем обратно неправильные категории (для отката миграции)
UPDATE marketplace_listings 
SET category_id = 1701, -- Industrial machinery category
    updated_at = NOW()
WHERE id IN (262, 263);