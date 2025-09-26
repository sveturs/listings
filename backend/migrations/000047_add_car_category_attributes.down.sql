-- Удаление атрибутов для автомобильных категорий
DELETE FROM category_attribute_mapping
WHERE category_id IN (1003, 1301)
  AND attribute_id IN (86, 87, 91, 113, 140, 148, 149);

-- Удаление атрибутов для подкатегорий
DELETE FROM category_attribute_mapping
WHERE category_id IN (
  SELECT id FROM marketplace_categories WHERE parent_id = 1301
)
AND attribute_id IN (86, 87, 91, 113, 140, 148, 149);