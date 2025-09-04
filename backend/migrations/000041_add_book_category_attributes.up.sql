-- Миграция 000041: Добавление атрибутов для категории Books & Stationery
-- Дата: 03.09.2025
-- Цель: Расширить функциональность для книг и канцелярии

-- Добавляем специфичные атрибуты для книг
INSERT INTO unified_attributes (code, name, display_name, attribute_type, purpose, is_active, created_at) VALUES
('book_author', 'Author', 'Автор', 'text', 'regular', true, CURRENT_TIMESTAMP),
('book_isbn', 'ISBN', 'ISBN', 'text', 'regular', true, CURRENT_TIMESTAMP),
('book_year', 'Publication Year', 'Год издания', 'number', 'regular', true, CURRENT_TIMESTAMP),
('book_language', 'Language', 'Язык', 'select', 'regular', true, CURRENT_TIMESTAMP),
('book_pages', 'Pages', 'Количество страниц', 'number', 'regular', true, CURRENT_TIMESTAMP),
('book_publisher', 'Publisher', 'Издательство', 'text', 'regular', true, CURRENT_TIMESTAMP),
('book_genre', 'Genre', 'Жанр', 'select', 'regular', true, CURRENT_TIMESTAMP),
('book_format', 'Format', 'Формат', 'select', 'regular', true, CURRENT_TIMESTAMP),
('book_binding', 'Binding', 'Переплёт', 'select', 'regular', true, CURRENT_TIMESTAMP),
('book_edition', 'Edition', 'Издание', 'text', 'regular', true, CURRENT_TIMESTAMP)
ON CONFLICT (code) DO NOTHING;

-- Обновляем опции для атрибутов типа select (хранятся в JSONB поле options)
UPDATE unified_attributes
SET options = CASE code
    WHEN 'book_language' THEN '[
        {"value": "ru", "label": "Русский"},
        {"value": "en", "label": "Английский"},
        {"value": "sr", "label": "Сербский"},
        {"value": "de", "label": "Немецкий"},
        {"value": "fr", "label": "Французский"},
        {"value": "es", "label": "Испанский"},
        {"value": "it", "label": "Итальянский"},
        {"value": "other", "label": "Другой"}
    ]'::jsonb
    WHEN 'book_genre' THEN '[
        {"value": "fiction", "label": "Художественная литература"},
        {"value": "non-fiction", "label": "Научно-популярная"},
        {"value": "science", "label": "Научная"},
        {"value": "textbook", "label": "Учебник"},
        {"value": "children", "label": "Детская"},
        {"value": "poetry", "label": "Поэзия"},
        {"value": "biography", "label": "Биография"},
        {"value": "business", "label": "Бизнес"},
        {"value": "self-help", "label": "Саморазвитие"},
        {"value": "cooking", "label": "Кулинария"}
    ]'::jsonb
    WHEN 'book_format' THEN '[
        {"value": "hardcover", "label": "Твёрдая обложка"},
        {"value": "paperback", "label": "Мягкая обложка"},
        {"value": "ebook", "label": "Электронная книга"},
        {"value": "audiobook", "label": "Аудиокнига"}
    ]'::jsonb
    WHEN 'book_binding' THEN '[
        {"value": "hard", "label": "Твёрдый"},
        {"value": "soft", "label": "Мягкий"},
        {"value": "spiral", "label": "На пружине"},
        {"value": "glued", "label": "Клеевой"},
        {"value": "stitched", "label": "Прошитый"}
    ]'::jsonb
    ELSE options
END
WHERE code IN ('book_language', 'book_genre', 'book_format', 'book_binding');

-- Привязываем атрибуты к категории Books & Stationery и её подкатегориям
INSERT INTO unified_category_attributes (category_id, attribute_id, is_required, is_enabled, sort_order)
SELECT 
    c.id as category_id,
    ua.id as attribute_id,
    CASE 
        WHEN ua.code IN ('book_author', 'book_year', 'book_language') THEN true
        ELSE false
    END as is_required,
    true as is_enabled,
    ROW_NUMBER() OVER (PARTITION BY c.id ORDER BY 
        CASE ua.code
            WHEN 'book_author' THEN 1
            WHEN 'book_year' THEN 2
            WHEN 'book_language' THEN 3
            WHEN 'book_genre' THEN 4
            WHEN 'book_publisher' THEN 5
            WHEN 'book_isbn' THEN 6
            WHEN 'book_pages' THEN 7
            WHEN 'book_format' THEN 8
            WHEN 'book_binding' THEN 9
            WHEN 'book_edition' THEN 10
            ELSE 11
        END
    ) as sort_order
FROM marketplace_categories c
CROSS JOIN unified_attributes ua
WHERE (c.id = 1012 OR c.parent_id = 1012) -- Books & Stationery и её подкатегории
  AND ua.code LIKE 'book_%'
ON CONFLICT (category_id, attribute_id) DO NOTHING;

-- Добавляем переводы для новых атрибутов
INSERT INTO translations (entity_type, entity_id, language, field_name, translated_text) 
SELECT 'unified_attribute', ua.id, t.lang, 'display_name', t.translation
FROM unified_attributes ua
CROSS JOIN (VALUES
    ('book_author', 'en', 'Author'),
    ('book_author', 'sr', 'Autor'),
    ('book_isbn', 'en', 'ISBN'),
    ('book_isbn', 'sr', 'ISBN'),
    ('book_year', 'en', 'Publication Year'),
    ('book_year', 'sr', 'Godina izdanja'),
    ('book_language', 'en', 'Language'),
    ('book_language', 'sr', 'Jezik'),
    ('book_pages', 'en', 'Pages'),
    ('book_pages', 'sr', 'Broj stranica'),
    ('book_publisher', 'en', 'Publisher'),
    ('book_publisher', 'sr', 'Izdavač'),
    ('book_genre', 'en', 'Genre'),
    ('book_genre', 'sr', 'Žanr'),
    ('book_format', 'en', 'Format'),
    ('book_format', 'sr', 'Format'),
    ('book_binding', 'en', 'Binding'),
    ('book_binding', 'sr', 'Povez'),
    ('book_edition', 'en', 'Edition'),
    ('book_edition', 'sr', 'Izdanje')
) AS t(attr_code, lang, translation)
WHERE ua.code = t.attr_code
ON CONFLICT DO NOTHING;

-- Логирование результата
DO $$
DECLARE
    attr_count INTEGER;
    cat_count INTEGER;
BEGIN
    SELECT COUNT(*) INTO attr_count 
    FROM unified_attributes 
    WHERE code LIKE 'book_%';
    
    SELECT COUNT(DISTINCT category_id) INTO cat_count
    FROM unified_category_attributes uca
    JOIN unified_attributes ua ON uca.attribute_id = ua.id
    WHERE ua.code LIKE 'book_%';
    
    RAISE NOTICE 'Добавлено % атрибутов для книг, привязано к % категориям', attr_count, cat_count;
END $$;