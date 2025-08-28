-- Удаление функции
DROP FUNCTION IF EXISTS increment_keyword_usage(INTEGER, TEXT[], VARCHAR(2));

-- Удаление начальных данных
DELETE FROM category_keywords WHERE source = 'manual';