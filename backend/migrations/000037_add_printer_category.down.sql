-- Откат добавления категории принтеров

-- Возвращаем товар обратно в категорию Photo & Video
UPDATE marketplace_listings
SET category_id = 1106
WHERE id = 323;

-- Удаляем ключевые слова
DELETE FROM category_keywords
WHERE category_id IN (2007, 2008);

-- Удаляем переводы
DELETE FROM translations
WHERE entity_type = 'category'
  AND entity_id IN (1109, 2007, 2008);

-- Удаляем категории
DELETE FROM marketplace_categories
WHERE id IN (2007, 2008);

DELETE FROM marketplace_categories
WHERE id = 1109;