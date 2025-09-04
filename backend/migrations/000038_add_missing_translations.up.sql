-- Миграция для добавления недостающих переводов категорий
-- Автор: System Architect
-- Дата: 03.09.2025
-- Задача: Добавить переводы для категорий, у которых они отсутствуют

-- Функция для транслитерации и базового перевода
CREATE OR REPLACE FUNCTION translate_category_name(category_name TEXT, target_lang TEXT)
RETURNS TEXT AS $$
DECLARE
    translated TEXT;
BEGIN
    translated := category_name;
    
    -- Базовые переводы для распространенных терминов
    IF target_lang = 'ru' THEN
        translated := REPLACE(translated, ' & ', ' и ');
        translated := REPLACE(translated, 'Accessories', 'Аксессуары');
        translated := REPLACE(translated, 'Electronics', 'Электроника');
        translated := REPLACE(translated, 'Clothing', 'Одежда');
        translated := REPLACE(translated, 'Sports', 'Спорт');
        translated := REPLACE(translated, 'Kids', 'Детские');
        translated := REPLACE(translated, 'Baby', 'Малыши');
        translated := REPLACE(translated, 'Home', 'Дом');
        translated := REPLACE(translated, 'Garden', 'Сад');
        translated := REPLACE(translated, 'Photo', 'Фото');
        translated := REPLACE(translated, 'Video', 'Видео');
        translated := REPLACE(translated, 'Gaming', 'Игровые');
        translated := REPLACE(translated, 'Consoles', 'Консоли');
        translated := REPLACE(translated, 'Smart', 'Умный');
        translated := REPLACE(translated, 'Watches', 'Часы');
        translated := REPLACE(translated, 'Bags', 'Сумки');
        
        -- Автомобильные термины
        translated := REPLACE(translated, 'Gume i točkovi', 'Шины и колеса');
        translated := REPLACE(translated, 'Motor i delovi motora', 'Двигатель и детали двигателя');
        translated := REPLACE(translated, 'Karoserija i delovi', 'Кузов и детали');
        translated := REPLACE(translated, 'Električni i elektronski delovi', 'Электрические и электронные детали');
        translated := REPLACE(translated, 'Sistem za kočenje', 'Тормозная система');
        translated := REPLACE(translated, 'Sistem vešanja', 'Подвеска');
        translated := REPLACE(translated, 'Sistem hlađenja', 'Система охлаждения');
        translated := REPLACE(translated, 'Transmisija i delovi', 'Трансмиссия и детали');
        translated := REPLACE(translated, 'Unutrašnjost', 'Интерьер');
        translated := REPLACE(translated, 'Dodatna oprema', 'Дополнительное оборудование');
        translated := REPLACE(translated, 'Zimske gume', 'Зимние шины');
        translated := REPLACE(translated, 'Letnje gume', 'Летние шины');
        translated := REPLACE(translated, 'Celogodišnje gume', 'Всесезонные шины');
        translated := REPLACE(translated, 'Felne', 'Диски');
        translated := REPLACE(translated, 'Kompletni točkovi', 'Колеса в сборе');
        translated := REPLACE(translated, 'Ratkapne', 'Колпаки');
        translated := REPLACE(translated, 'Vijci za točkove', 'Болты для колес');
        translated := REPLACE(translated, 'Filtri', 'Фильтры');
        translated := REPLACE(translated, 'Remeni i lančanici', 'Ремни и цепи');
        translated := REPLACE(translated, 'wifi-routery', 'WiFi роутеры');
        
    ELSIF target_lang = 'en' THEN
        -- Сербские термины на английский
        translated := REPLACE(translated, 'Gume i točkovi', 'Tires and Wheels');
        translated := REPLACE(translated, 'Motor i delovi motora', 'Engine and Engine Parts');
        translated := REPLACE(translated, 'Karoserija i delovi', 'Body and Parts');
        translated := REPLACE(translated, 'Električni i elektronski delovi', 'Electrical and Electronic Parts');
        translated := REPLACE(translated, 'Sistem za kočenje', 'Braking System');
        translated := REPLACE(translated, 'Sistem vešanja', 'Suspension System');
        translated := REPLACE(translated, 'Sistem hlađenja', 'Cooling System');
        translated := REPLACE(translated, 'Transmisija i delovi', 'Transmission and Parts');
        translated := REPLACE(translated, 'Unutrašnjost', 'Interior');
        translated := REPLACE(translated, 'Dodatna oprema', 'Additional Equipment');
        translated := REPLACE(translated, 'Zimske gume', 'Winter Tires');
        translated := REPLACE(translated, 'Letnje gume', 'Summer Tires');
        translated := REPLACE(translated, 'Celogodišnje gume', 'All-Season Tires');
        translated := REPLACE(translated, 'Felne', 'Rims');
        translated := REPLACE(translated, 'Kompletni točkovi', 'Complete Wheels');
        translated := REPLACE(translated, 'Ratkapne', 'Hubcaps');
        translated := REPLACE(translated, 'Vijci za točkove', 'Wheel Bolts');
        translated := REPLACE(translated, 'Filtri', 'Filters');
        translated := REPLACE(translated, 'Remeni i lančanici', 'Belts and Chains');
        translated := REPLACE(translated, 'wifi-routery', 'WiFi Routers');
        translated := REPLACE(translated, 'photo', 'Photo');
        
    ELSIF target_lang = 'sr' THEN
        -- Английские термины на сербский (если есть)
        translated := REPLACE(translated, 'photo', 'foto');
        translated := REPLACE(translated, 'wifi-routery', 'wifi-ruteri');
        translated := REPLACE(translated, 'Smart Home', 'Pametna kuća');
        translated := REPLACE(translated, 'Gaming Consoles', 'Igračke konzole');
    END IF;
    
    RETURN translated;
END;
$$ LANGUAGE plpgsql;

-- Добавляем недостающие переводы для русского языка
INSERT INTO translations (entity_type, entity_id, language, field_name, translated_text)
SELECT 
    'category' as entity_type,
    mc.id as entity_id,
    'ru' as language,
    'name' as field_name,
    translate_category_name(mc.name, 'ru') as translated_text
FROM marketplace_categories mc
WHERE NOT EXISTS (
    SELECT 1 FROM translations t 
    WHERE t.entity_type = 'category' 
    AND t.entity_id = mc.id 
    AND t.language = 'ru'
    AND t.field_name = 'name'
)
ON CONFLICT DO NOTHING;

-- Добавляем недостающие переводы для английского языка
INSERT INTO translations (entity_type, entity_id, language, field_name, translated_text)
SELECT 
    'category' as entity_type,
    mc.id as entity_id,
    'en' as language,
    'name' as field_name,
    translate_category_name(mc.name, 'en') as translated_text
FROM marketplace_categories mc
WHERE NOT EXISTS (
    SELECT 1 FROM translations t 
    WHERE t.entity_type = 'category' 
    AND t.entity_id = mc.id 
    AND t.language = 'en'
    AND t.field_name = 'name'
)
ON CONFLICT DO NOTHING;

-- Добавляем недостающие переводы для сербского языка
INSERT INTO translations (entity_type, entity_id, language, field_name, translated_text)
SELECT 
    'category' as entity_type,
    mc.id as entity_id,
    'sr' as language,
    'name' as field_name,
    translate_category_name(mc.name, 'sr') as translated_text
FROM marketplace_categories mc
WHERE NOT EXISTS (
    SELECT 1 FROM translations t 
    WHERE t.entity_type = 'category' 
    AND t.entity_id = mc.id 
    AND t.language = 'sr'
    AND t.field_name = 'name'
)
ON CONFLICT DO NOTHING;

-- Выводим статистику
DO $$
DECLARE
    total_categories INTEGER;
    translated_ru INTEGER;
    translated_en INTEGER;
    translated_sr INTEGER;
BEGIN
    SELECT COUNT(*) INTO total_categories FROM marketplace_categories;
    
    SELECT COUNT(DISTINCT entity_id) INTO translated_ru 
    FROM translations 
    WHERE entity_type = 'category' AND language = 'ru' AND field_name = 'name';
    
    SELECT COUNT(DISTINCT entity_id) INTO translated_en
    FROM translations 
    WHERE entity_type = 'category' AND language = 'en' AND field_name = 'name';
    
    SELECT COUNT(DISTINCT entity_id) INTO translated_sr
    FROM translations 
    WHERE entity_type = 'category' AND language = 'sr' AND field_name = 'name';
    
    RAISE NOTICE 'Добавление переводов завершено';
    RAISE NOTICE 'Всего категорий: %', total_categories;
    RAISE NOTICE 'Переводов на русский: %', translated_ru;
    RAISE NOTICE 'Переводов на английский: %', translated_en;
    RAISE NOTICE 'Переводов на сербский: %', translated_sr;
END $$;

-- Удаляем временную функцию
DROP FUNCTION IF EXISTS translate_category_name(TEXT, TEXT);