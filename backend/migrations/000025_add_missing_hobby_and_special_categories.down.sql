-- Удаление переводов
DELETE FROM translations WHERE entity_type = 'category' AND entity_id BETWEEN 2001 AND 2063;

-- Удаление AI маппингов для новых категорий
DELETE FROM category_ai_mappings WHERE category_id BETWEEN 2001 AND 2063;

-- Удаление категорий (каскадно удалятся подкатегории)
DELETE FROM marketplace_categories WHERE id BETWEEN 2001 AND 2063;