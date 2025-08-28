-- Добавление категории "Прочее" / "Other" / "Ostalo" для объявлений, 
-- которые не удалось автоматически классифицировать

-- Создаем категорию "Прочее" с ID 9999 (специальный ID для легкой идентификации)
INSERT INTO marketplace_categories (id, name, slug, icon, parent_id, sort_order, is_active, level, count)
VALUES (9999, 'Other', 'other', '❓', NULL, 9999, true, 0, 0);

-- Добавляем переводы для категории
INSERT INTO translations (entity_type, entity_id, field_name, language, translated_text) VALUES
-- Английский (основной)
('category', 9999, 'name', 'en', 'Other'),
('category', 9999, 'description', 'en', 'Items that do not fit into other categories'),

-- Русский
('category', 9999, 'name', 'ru', 'Прочее'),
('category', 9999, 'description', 'ru', 'Товары, не подходящие под другие категории'),

-- Сербский (латиница)
('category', 9999, 'name', 'sr', 'Ostalo'),
('category', 9999, 'description', 'sr', 'Proizvodi koji ne spadaju u druge kategorije');

-- Добавляем ключевые слова для категории "Прочее"
-- Эти слова помогут направлять сюда неопределенные товары
INSERT INTO category_keywords (category_id, keyword, language, weight, keyword_type, source) VALUES
-- Английский
(9999, 'other', 'en', 0.5, 'main', 'manual'),
(9999, 'miscellaneous', 'en', 0.5, 'main', 'manual'),
(9999, 'various', 'en', 0.3, 'synonym', 'manual'),
(9999, 'different', 'en', 0.3, 'synonym', 'manual'),
(9999, 'mixed', 'en', 0.3, 'synonym', 'manual'),

-- Русский
(9999, 'прочее', 'ru', 0.5, 'main', 'manual'),
(9999, 'другое', 'ru', 0.5, 'main', 'manual'),
(9999, 'разное', 'ru', 0.5, 'main', 'manual'),
(9999, 'остальное', 'ru', 0.3, 'synonym', 'manual'),
(9999, 'различное', 'ru', 0.3, 'synonym', 'manual'),

-- Сербский
(9999, 'ostalo', 'sr', 0.5, 'main', 'manual'),
(9999, 'drugo', 'sr', 0.5, 'main', 'manual'),
(9999, 'razno', 'sr', 0.5, 'main', 'manual'),
(9999, 'različito', 'sr', 0.3, 'synonym', 'manual');

-- Обновляем последовательность, чтобы избежать конфликтов ID в будущем
SELECT setval('marketplace_categories_id_seq', GREATEST(10000, (SELECT MAX(id) FROM marketplace_categories)));

-- Добавляем комментарий к категории для документации
COMMENT ON COLUMN marketplace_categories.id IS 'Category ID. ID 9999 is reserved for "Other" category for unclassified items';