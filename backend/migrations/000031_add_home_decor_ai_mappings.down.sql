-- Удаляем маппинги для категории "Дом и сад"
DELETE FROM category_ai_mappings
WHERE category_id = 1005
  AND ai_domain IN ('home decor', 'home', 'interior', 'garden', 'furniture', 'home textiles', 'kitchen', 'lighting');