-- Удаление связей атрибутов с категориями
DELETE FROM category_attribute_mapping WHERE attribute_id >= 2601;

-- Удаление новых атрибутов
DELETE FROM category_attributes WHERE id >= 2601;

-- Восстановление пустых опций для атрибутов
UPDATE category_attributes SET options = '{}' WHERE name IN ('body_type', 'service_type', 'availability', 'service_area');

-- Удаление новых подкатегорий
DELETE FROM marketplace_categories WHERE id IN (
    1105, 1106, 1107, 1108, -- электроника
    1205, 1206, 1207, 1208, -- мода
    1505, 1506, 1507, 1508, -- дом и сад
    2011, 2012, 2013, 2014  -- спорт
);

-- Восстанавливаем тестовую категорию
INSERT INTO marketplace_categories (id, slug, parent_id, icon, created_at, name) VALUES
(2005, 'test', 1103, '🚌', CURRENT_TIMESTAMP, 'Test');

-- Удаление новых основных категорий
DELETE FROM marketplace_categories WHERE id BETWEEN 1011 AND 1020;