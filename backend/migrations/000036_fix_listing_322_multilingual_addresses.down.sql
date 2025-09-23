-- Откат к сербским значениям (которые были ошибочно сохранены)

-- Английские переводы - возврат к сербским
UPDATE translations
SET translated_text = 'Васе Стајића 18, Нови Сад, Јужнобачки управни округ'
WHERE entity_type = 'listing'
  AND entity_id = 322
  AND field_name = 'location'
  AND language = 'en';

UPDATE translations
SET translated_text = 'Нови Сад'
WHERE entity_type = 'listing'
  AND entity_id = 322
  AND field_name = 'city'
  AND language = 'en';

UPDATE translations
SET translated_text = 'Србија'
WHERE entity_type = 'listing'
  AND entity_id = 322
  AND field_name = 'country'
  AND language = 'en';

-- Русские переводы - возврат к сербским
UPDATE translations
SET translated_text = 'Васе Стајића 18, Нови Сад, Јужнобачки управни округ'
WHERE entity_type = 'listing'
  AND entity_id = 322
  AND field_name = 'location'
  AND language = 'ru';

UPDATE translations
SET translated_text = 'Нови Сад'
WHERE entity_type = 'listing'
  AND entity_id = 322
  AND field_name = 'city'
  AND language = 'ru';

UPDATE translations
SET translated_text = 'Србија'
WHERE entity_type = 'listing'
  AND entity_id = 322
  AND field_name = 'country'
  AND language = 'ru';