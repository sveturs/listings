-- Откат: возвращаем исходные значения
UPDATE marketplace_listings
SET address_multilingual = jsonb_build_object(
    'en', location,
    'ru', location,
    'sr', location
)
WHERE address_multilingual IS NOT NULL;