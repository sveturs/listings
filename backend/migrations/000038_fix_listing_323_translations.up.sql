-- Исправление мультиязычных адресов для объявления 323 (Canon принтер)

-- Исправляем английскую локализацию (проблема с кодировкой ć)
UPDATE translations
SET translated_text = 'Vase Stajića, Novi Sad, South Bačka District'
WHERE entity_type = 'listing'
  AND entity_id = 323
  AND field_name = 'location'
  AND language = 'en';

UPDATE translations
SET translated_text = 'Serbia'
WHERE entity_type = 'listing'
  AND entity_id = 323
  AND field_name = 'country'
  AND language = 'en';

-- Добавляем недостающий русский перевод
INSERT INTO translations (entity_type, entity_id, field_name, language, translated_text, created_at)
VALUES
    ('listing', 323, 'location', 'ru', 'Васе Стайича, Нови-Сад, Южно-Бачский округ', NOW()),
    ('listing', 323, 'city', 'ru', 'Нови-Сад', NOW()),
    ('listing', 323, 'country', 'ru', 'Сербия', NOW())
ON CONFLICT (entity_type, entity_id, field_name, language)
DO UPDATE SET translated_text = EXCLUDED.translated_text;