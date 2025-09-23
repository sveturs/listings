-- Откат исправлений для объявления 323

-- Возврат к старым значениям
UPDATE translations
SET translated_text = 'Vase Stajiћa, Novi Sad, South Bacsk District'
WHERE entity_type = 'listing'
  AND entity_id = 323
  AND field_name = 'location'
  AND language = 'en';

UPDATE translations
SET translated_text = 'Srbija'
WHERE entity_type = 'listing'
  AND entity_id = 323
  AND field_name = 'country'
  AND language = 'en';

-- Удаление русских переводов
DELETE FROM translations
WHERE entity_type = 'listing'
  AND entity_id = 323
  AND language = 'ru';