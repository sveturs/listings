-- Откат: удаляем добавленные ключевые слова

DELETE FROM category_keywords
WHERE category_id IN (1301, 1302, 1303)
AND created_at >= NOW() - INTERVAL '1 hour';