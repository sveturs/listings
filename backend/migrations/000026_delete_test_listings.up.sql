-- Удаление тестовых данных

-- Сначала удаляем данные из зависимых таблиц, чтобы избежать нарушения ограничений внешнего ключа

-- Удаление голосов за отзывы
DELETE FROM review_votes
WHERE review_id IN (SELECT id FROM reviews WHERE entity_type = 'listing' AND entity_id IN (8, 9, 10, 12));

-- Удаление ответов на отзывы
DELETE FROM review_responses
WHERE review_id IN (SELECT id FROM reviews WHERE entity_type = 'listing' AND entity_id IN (8, 9, 10, 12));

-- Удаление отзывов
DELETE FROM reviews
WHERE entity_type = 'listing' AND entity_id IN (8, 9, 10, 12);

-- Удаление изображений объявлений
DELETE FROM marketplace_images
WHERE listing_id IN (8, 9, 10, 12);

-- Удаление самих объявлений
DELETE FROM marketplace_listings
WHERE id IN (8, 9, 10, 12);