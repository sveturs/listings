-- Удаление таблицы category_keywords
DROP TRIGGER IF EXISTS update_category_keywords_updated_at ON category_keywords;
DROP TABLE IF EXISTS category_keywords;