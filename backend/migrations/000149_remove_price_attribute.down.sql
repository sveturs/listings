-- Восстанавливаем атрибут price

-- 1. Восстанавливаем атрибут price
INSERT INTO category_attributes (id, name, display_name, attribute_type, is_required, sort_order, is_active, created_at, is_searchable, show_in_filters, show_in_card, use_in_similar, validation_rules) 
VALUES (2001, 'price', 'Cena', 'number', false, 10, true, CURRENT_TIMESTAMP, false, false, true, false, NULL);

-- 2. Восстанавливаем переводы
INSERT INTO translations (entity_type, entity_id, field_name, language_code, translated_text)
VALUES 
    ('attribute', 2001, 'display_name', 'ru', 'Цена'),
    ('attribute', 2001, 'display_name', 'en', 'Price'),
    ('attribute', 2001, 'display_name', 'sr', 'Cena');

-- Примечание: значения атрибутов для существующих объявлений не восстанавливаются