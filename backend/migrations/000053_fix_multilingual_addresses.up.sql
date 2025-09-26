-- Исправляем мультиязычные адреса для существующих объявлений
-- Для английского языка - транслитерация и перевод районов
UPDATE marketplace_listings
SET address_multilingual = jsonb_build_object(
    'en', CASE
        WHEN location LIKE '%Васе Стајића%' THEN 'Vase Stajica 20, Novi Sad 21101, South Bačka District, Serbia'
        WHEN location LIKE '%Vase Stajica%' THEN location
        ELSE location
    END,
    'ru', CASE
        WHEN location LIKE '%Васе Стајића%' THEN 'Васе Стаича 20, Нови-Сад 21101, Южно-Бачский округ, Сербия'
        WHEN location LIKE '%Vase Stajica%' THEN 'Васе Стаича, Нови-Сад, Южно-Бачский округ'
        ELSE location
    END,
    'sr', CASE
        WHEN location LIKE '%Vase Stajica%' THEN 'Васе Стајића, Нови Сад, Јужнобачки округ'
        ELSE location
    END
)
WHERE address_multilingual IS NOT NULL
AND (location LIKE '%Васе Стајића%' OR location LIKE '%Vase Stajica%');

-- Обновляем адреса для конкретных объявлений с правильными переводами
UPDATE marketplace_listings
SET address_multilingual = jsonb_build_object(
    'en', 'Vase Stajica 20, Novi Sad 21101, South Bačka District, Serbia',
    'ru', 'Васе Стаича 20, Нови-Сад 21101, Южно-Бачский округ, Сербия',
    'sr', 'Васе Стајића 20, Нови Сад 21101, Јужнобачки управни округ, Србија'
)
WHERE id = 328;

UPDATE marketplace_listings
SET address_multilingual = jsonb_build_object(
    'en', 'Vase Stajica, Novi Sad, South Bačka District',
    'ru', 'Васе Стаича, Нови-Сад, Южно-Бачский округ',
    'sr', 'Васе Стајића, Нови Сад, Јужнобачки округ'
)
WHERE id IN (325, 326, 327);