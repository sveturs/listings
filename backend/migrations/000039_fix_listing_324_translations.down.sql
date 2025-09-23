-- Откат исправлений для объявления 324

-- Удаление английских переводов
DELETE FROM translations
WHERE entity_type = 'listing'
  AND entity_id = 324
  AND field_name IN ('location', 'city', 'country')
  AND language = 'en';

-- Удаление русских переводов
DELETE FROM translations
WHERE entity_type = 'listing'
  AND entity_id = 324
  AND field_name IN ('location', 'city', 'country')
  AND language = 'ru';