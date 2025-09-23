-- Исправление мультиязычных адресов для объявления 324 (Paper Cutting Knife)

-- Добавляем английские переводы
INSERT INTO translations (entity_type, entity_id, field_name, language, translated_text, created_at)
VALUES
    ('listing', 324, 'location', 'en', 'Vase Stajića, Novi Sad, South Bačka District', NOW()),
    ('listing', 324, 'city', 'en', 'Novi Sad', NOW()),
    ('listing', 324, 'country', 'en', 'Serbia', NOW())
ON CONFLICT (entity_type, entity_id, field_name, language)
DO UPDATE SET translated_text = EXCLUDED.translated_text;

-- Добавляем русские переводы
INSERT INTO translations (entity_type, entity_id, field_name, language, translated_text, created_at)
VALUES
    ('listing', 324, 'location', 'ru', 'Васе Стайича, Нови-Сад, Южно-Бачский округ', NOW()),
    ('listing', 324, 'city', 'ru', 'Нови-Сад', NOW()),
    ('listing', 324, 'country', 'ru', 'Сербия', NOW())
ON CONFLICT (entity_type, entity_id, field_name, language)
DO UPDATE SET translated_text = EXCLUDED.translated_text;