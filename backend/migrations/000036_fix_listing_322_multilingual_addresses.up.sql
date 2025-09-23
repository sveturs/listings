-- Исправление мультиязычных адресов для объявления 322, созданного через AI
-- Это объявление было создано до исправления мультиязычной поддержки

-- Английские переводы
UPDATE translations
SET translated_text = 'Vase Stajića 18, Novi Sad, South Bačka District'
WHERE entity_type = 'listing'
  AND entity_id = 322
  AND field_name = 'location'
  AND language = 'en';

UPDATE translations
SET translated_text = 'Novi Sad'
WHERE entity_type = 'listing'
  AND entity_id = 322
  AND field_name = 'city'
  AND language = 'en';

UPDATE translations
SET translated_text = 'Serbia'
WHERE entity_type = 'listing'
  AND entity_id = 322
  AND field_name = 'country'
  AND language = 'en';

-- Русские переводы
UPDATE translations
SET translated_text = 'Васе Стаича 18, Нови-Сад, Южно-Бачский округ'
WHERE entity_type = 'listing'
  AND entity_id = 322
  AND field_name = 'location'
  AND language = 'ru';

UPDATE translations
SET translated_text = 'Нови-Сад'
WHERE entity_type = 'listing'
  AND entity_id = 322
  AND field_name = 'city'
  AND language = 'ru';

UPDATE translations
SET translated_text = 'Сербия'
WHERE entity_type = 'listing'
  AND entity_id = 322
  AND field_name = 'country'
  AND language = 'ru';