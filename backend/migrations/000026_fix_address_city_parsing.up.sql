-- Миграция для исправления неправильного парсинга города из полного адреса
-- Проблема: address_city содержит улицу вместо города
-- Формат адреса: "Улица номер, Город почтовый_код, Регион, Страна"
-- Нужно извлечь Город из второй части (index 1) и убрать почтовый код

-- Функция для извлечения города из полного адреса
CREATE OR REPLACE FUNCTION extract_city_from_location(full_address TEXT)
RETURNS TEXT AS $$
DECLARE
    parts TEXT[];
    city_part TEXT;
BEGIN
    IF full_address IS NULL OR full_address = '' THEN
        RETURN NULL;
    END IF;

    -- Разбиваем адрес по запятой
    parts := string_to_array(full_address, ',');

    -- Если только одна часть - возвращаем как есть (скорее всего уже город)
    IF array_length(parts, 1) = 1 THEN
        RETURN trim(parts[1]);
    END IF;

    -- Если больше одной части - берем вторую (index 2 в PostgreSQL, т.к. массивы с 1)
    IF array_length(parts, 1) > 1 THEN
        city_part := trim(parts[2]);
        -- Убираем почтовый код (5 цифр в конце): "Нови Сад 21101" → "Нови Сад"
        city_part := regexp_replace(city_part, '\s+\d{5}.*$', '');
        RETURN trim(city_part);
    END IF;

    RETURN NULL;
END;
$$ LANGUAGE plpgsql IMMUTABLE;

-- Функция для извлечения страны из полного адреса
CREATE OR REPLACE FUNCTION extract_country_from_location(full_address TEXT)
RETURNS TEXT AS $$
DECLARE
    parts TEXT[];
BEGIN
    IF full_address IS NULL OR full_address = '' THEN
        RETURN 'Србија';
    END IF;

    -- Разбиваем адрес по запятой
    parts := string_to_array(full_address, ',');

    -- Если только одна часть - возвращаем дефолт
    IF array_length(parts, 1) = 1 THEN
        RETURN 'Србија';
    END IF;

    -- Берем последнюю часть
    RETURN trim(parts[array_length(parts, 1)]);
END;
$$ LANGUAGE plpgsql IMMUTABLE;

-- Обновляем address_city и address_country для всех объявлений, где location заполнен
UPDATE marketplace_listings
SET
    address_city = extract_city_from_location(location),
    address_country = extract_country_from_location(location)
WHERE location IS NOT NULL AND location != '';

-- Обновляем переводы city и country в таблице translations
-- Используем значение из location для каждого языка
UPDATE translations t
SET translated_text = extract_city_from_location(
    COALESCE(
        (SELECT translated_text FROM translations t2
         WHERE t2.entity_type = 'listing'
         AND t2.entity_id = t.entity_id
         AND t2.language = t.language
         AND t2.field_name = 'location'
         LIMIT 1),
        (SELECT location FROM marketplace_listings WHERE id = t.entity_id)
    )
)
WHERE t.entity_type = 'listing'
AND t.field_name = 'city'
AND EXISTS (
    SELECT 1 FROM marketplace_listings ml
    WHERE ml.id = t.entity_id
    AND ml.location IS NOT NULL
    AND ml.location != ''
);

UPDATE translations t
SET translated_text = extract_country_from_location(
    COALESCE(
        (SELECT translated_text FROM translations t2
         WHERE t2.entity_type = 'listing'
         AND t2.entity_id = t.entity_id
         AND t2.language = t.language
         AND t2.field_name = 'location'
         LIMIT 1),
        (SELECT location FROM marketplace_listings WHERE id = t.entity_id)
    )
)
WHERE t.entity_type = 'listing'
AND t.field_name = 'country'
AND EXISTS (
    SELECT 1 FROM marketplace_listings ml
    WHERE ml.id = t.entity_id
    AND ml.location IS NOT NULL
    AND ml.location != ''
);

-- Логируем результат
DO $$
DECLARE
    updated_listings INTEGER;
    updated_translations_city INTEGER;
    updated_translations_country INTEGER;
BEGIN
    SELECT COUNT(*) INTO updated_listings
    FROM marketplace_listings
    WHERE location IS NOT NULL AND location != '';

    SELECT COUNT(*) INTO updated_translations_city
    FROM translations
    WHERE entity_type = 'listing' AND field_name = 'city';

    SELECT COUNT(*) INTO updated_translations_country
    FROM translations
    WHERE entity_type = 'listing' AND field_name = 'country';

    RAISE NOTICE 'Updated % listings with corrected city/country parsing', updated_listings;
    RAISE NOTICE 'Updated % city translations', updated_translations_city;
    RAISE NOTICE 'Updated % country translations', updated_translations_country;
END $$;
