-- Удаляем добавленные колонки из таблицы reviews
ALTER TABLE reviews
DROP COLUMN IF EXISTS seller_confirmed,
DROP COLUMN IF EXISTS has_active_dispute;

-- Удаляем таблицы в обратном порядке создания
DROP TABLE IF EXISTS review_dispute_messages;
DROP TABLE IF EXISTS review_disputes;
DROP TABLE IF EXISTS review_confirmations;

-- Удаляем уникальный индекс
DROP INDEX IF EXISTS idx_reviews_user_entity_unique;

-- Удаляем индексы
DROP INDEX IF EXISTS idx_reviews_seller_confirmed;
DROP INDEX IF EXISTS idx_reviews_has_dispute;