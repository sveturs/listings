-- Проверяем материализованные представления
\echo 'Checking materialized views...'

-- Проверка user_ratings
\echo '=== user_ratings ==='
SELECT * FROM user_ratings WHERE user_id = 25;

-- Проверка user_rating_distribution
\echo '=== user_rating_distribution ==='
SELECT * FROM user_rating_distribution WHERE user_id = 25;

-- Проверка исходных данных в reviews
\echo '=== reviews for user 25 ==='
SELECT id, entity_type, entity_id, entity_origin_type, entity_origin_id, rating, status
FROM reviews 
WHERE entity_origin_type = 'user' 
AND entity_origin_id = 25 
AND status = 'published';

-- Проверка последнего созданного отзыва
\echo '=== Latest review ==='
SELECT id, user_id, entity_type, entity_id, entity_origin_type, entity_origin_id, rating, comment, status, created_at
FROM reviews 
ORDER BY created_at DESC 
LIMIT 1;

-- Проверка структуры материализованных представлений
\echo '=== Structure of user_ratings ==='
\d user_ratings

\echo '=== Structure of user_rating_distribution ==='
\d user_rating_distribution