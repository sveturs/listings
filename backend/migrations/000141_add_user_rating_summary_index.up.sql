-- Добавление уникального индекса для materialized view user_rating_summary
-- Необходимо для поддержки REFRESH MATERIALIZED VIEW CONCURRENTLY

-- Создаем уникальный индекс на user_id
CREATE UNIQUE INDEX IF NOT EXISTS user_rating_summary_user_id_idx 
ON user_rating_summary (user_id);