-- Удаление категории "Прочее" и связанных данных

-- Удаляем ключевые слова
DELETE FROM category_keywords WHERE category_id = 9999;

-- Удаляем переводы
DELETE FROM translations WHERE entity_type = 'category' AND entity_id = 9999;

-- Удаляем саму категорию
DELETE FROM marketplace_categories WHERE id = 9999;