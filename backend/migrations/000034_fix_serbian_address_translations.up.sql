-- Исправляем переводы городов на сербский язык
UPDATE translations
SET translated_text = 'Београд'
WHERE entity_type = 'listing'
  AND language = 'sr'
  AND field_name IN ('location', 'city')
  AND translated_text IN ('Белград', 'Belgrade');

-- Исправляем перевод страны на сербский язык
UPDATE translations
SET translated_text = 'Србија'
WHERE entity_type = 'listing'
  AND language = 'sr'
  AND field_name = 'country'
  AND translated_text IN ('Сербия', 'Serbia');

-- Исправляем другие распространенные города Сербии
UPDATE translations
SET translated_text = 'Нови Сад'
WHERE entity_type = 'listing'
  AND language = 'sr'
  AND field_name IN ('location', 'city')
  AND translated_text IN ('Новый Сад', 'Novi Sad');

UPDATE translations
SET translated_text = 'Ниш'
WHERE entity_type = 'listing'
  AND language = 'sr'
  AND field_name IN ('location', 'city')
  AND translated_text IN ('Ниш', 'Niš', 'Nis');

UPDATE translations
SET translated_text = 'Крагујевац'
WHERE entity_type = 'listing'
  AND language = 'sr'
  AND field_name IN ('location', 'city')
  AND translated_text IN ('Крагуевац', 'Kragujevac');