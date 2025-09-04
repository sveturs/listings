-- Миграция 000042: Добавление атрибутов для категории Pets
-- Дата: 03.09.2025
-- Цель: Расширить функциональность для животных и зоотоваров

-- Добавляем специфичные атрибуты для животных
INSERT INTO unified_attributes (code, name, display_name, attribute_type, purpose, is_active, created_at) VALUES
('pet_type', 'Pet Type', 'Тип животного', 'select', 'regular', true, CURRENT_TIMESTAMP),
('pet_breed', 'Breed', 'Порода', 'text', 'regular', true, CURRENT_TIMESTAMP),
('pet_age', 'Age', 'Возраст', 'text', 'regular', true, CURRENT_TIMESTAMP),
('pet_gender', 'Gender', 'Пол', 'select', 'regular', true, CURRENT_TIMESTAMP),
('pet_vaccinated', 'Vaccinated', 'Вакцинирован', 'boolean', 'regular', true, CURRENT_TIMESTAMP),
('pet_sterilized', 'Sterilized', 'Стерилизован', 'boolean', 'regular', true, CURRENT_TIMESTAMP),
('pet_chipped', 'Microchipped', 'Чипирован', 'boolean', 'regular', true, CURRENT_TIMESTAMP),
('pet_pedigree', 'Pedigree', 'Родословная', 'boolean', 'regular', true, CURRENT_TIMESTAMP),
('pet_color', 'Color', 'Окрас', 'text', 'regular', true, CURRENT_TIMESTAMP),
('pet_weight', 'Weight', 'Вес (кг)', 'number', 'regular', true, CURRENT_TIMESTAMP)
ON CONFLICT (code) DO NOTHING;

-- Обновляем опции для атрибутов типа select (хранятся в JSONB поле options)
UPDATE unified_attributes
SET options = CASE code
    WHEN 'pet_type' THEN '[
        {"value": "dog", "label": "Собака"},
        {"value": "cat", "label": "Кошка"},
        {"value": "bird", "label": "Птица"},
        {"value": "fish", "label": "Рыбка"},
        {"value": "rodent", "label": "Грызун"},
        {"value": "reptile", "label": "Рептилия"},
        {"value": "rabbit", "label": "Кролик"},
        {"value": "horse", "label": "Лошадь"},
        {"value": "farm", "label": "Сельхоз животное"},
        {"value": "exotic", "label": "Экзотическое"},
        {"value": "other", "label": "Другое"}
    ]'::jsonb
    WHEN 'pet_gender' THEN '[
        {"value": "male", "label": "Самец"},
        {"value": "female", "label": "Самка"},
        {"value": "unknown", "label": "Не определен"}
    ]'::jsonb
    ELSE options
END
WHERE code IN ('pet_type', 'pet_gender');

-- Привязываем атрибуты к категории Pets и её подкатегориям
INSERT INTO unified_category_attributes (category_id, attribute_id, is_required, is_enabled, sort_order)
SELECT 
    c.id as category_id,
    ua.id as attribute_id,
    CASE 
        WHEN ua.code IN ('pet_type', 'pet_age') THEN true
        ELSE false
    END as is_required,
    true as is_enabled,
    ROW_NUMBER() OVER (PARTITION BY c.id ORDER BY 
        CASE ua.code
            WHEN 'pet_type' THEN 1
            WHEN 'pet_breed' THEN 2
            WHEN 'pet_age' THEN 3
            WHEN 'pet_gender' THEN 4
            WHEN 'pet_color' THEN 5
            WHEN 'pet_weight' THEN 6
            WHEN 'pet_vaccinated' THEN 7
            WHEN 'pet_sterilized' THEN 8
            WHEN 'pet_chipped' THEN 9
            WHEN 'pet_pedigree' THEN 10
            ELSE 11
        END
    ) as sort_order
FROM marketplace_categories c
CROSS JOIN unified_attributes ua
WHERE (c.id = 1011 OR c.parent_id = 1011) -- Pets и её подкатегории
  AND ua.code LIKE 'pet_%'
ON CONFLICT (category_id, attribute_id) DO NOTHING;

-- Добавляем переводы для новых атрибутов
INSERT INTO translations (entity_type, entity_id, language, field_name, translated_text) 
SELECT 'unified_attribute', ua.id, t.lang, 'display_name', t.translation
FROM unified_attributes ua
CROSS JOIN (VALUES
    ('pet_type', 'en', 'Pet Type'),
    ('pet_type', 'sr', 'Tip životinje'),
    ('pet_breed', 'en', 'Breed'),
    ('pet_breed', 'sr', 'Rasa'),
    ('pet_age', 'en', 'Age'),
    ('pet_age', 'sr', 'Uzrast'),
    ('pet_gender', 'en', 'Gender'),
    ('pet_gender', 'sr', 'Pol'),
    ('pet_vaccinated', 'en', 'Vaccinated'),
    ('pet_vaccinated', 'sr', 'Vakcinisan'),
    ('pet_sterilized', 'en', 'Sterilized'),
    ('pet_sterilized', 'sr', 'Sterilisan'),
    ('pet_chipped', 'en', 'Microchipped'),
    ('pet_chipped', 'sr', 'Čipovan'),
    ('pet_pedigree', 'en', 'Pedigree'),
    ('pet_pedigree', 'sr', 'Rodovnik'),
    ('pet_color', 'en', 'Color'),
    ('pet_color', 'sr', 'Boja'),
    ('pet_weight', 'en', 'Weight (kg)'),
    ('pet_weight', 'sr', 'Težina (kg)')
) AS t(attr_code, lang, translation)
WHERE ua.code = t.attr_code
ON CONFLICT DO NOTHING;

-- Переводы для опций хранятся в самом JSONB, дополнительные переводы не требуются

-- Логирование результата
DO $$
DECLARE
    attr_count INTEGER;
    cat_count INTEGER;
BEGIN
    SELECT COUNT(*) INTO attr_count 
    FROM unified_attributes 
    WHERE code LIKE 'pet_%';
    
    SELECT COUNT(DISTINCT category_id) INTO cat_count
    FROM unified_category_attributes uca
    JOIN unified_attributes ua ON uca.attribute_id = ua.id
    WHERE ua.code LIKE 'pet_%';
    
    RAISE NOTICE 'Добавлено % атрибутов для животных, привязано к % категориям', attr_count, cat_count;
END $$;