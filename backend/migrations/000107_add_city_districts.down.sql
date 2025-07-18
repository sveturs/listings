-- Откат миграции для удаления районов крупных городов Сербии

-- Функция для получения city_id по имени
CREATE OR REPLACE FUNCTION get_city_id_by_name(city_name TEXT) 
RETURNS UUID AS $$
DECLARE
    city_id UUID;
BEGIN
    SELECT id INTO city_id FROM cities WHERE name = city_name LIMIT 1;
    RETURN city_id;
END;
$$ LANGUAGE plpgsql;

-- Удаляем районы Нови Сада
DELETE FROM districts WHERE city_id = get_city_id_by_name('Нови Сад') 
AND name IN ('Лиман', 'Ново насеље', 'Грбавица', 'Детелинара', 'Клиса', 
             'Подбара', 'Салајка', 'Роткварија', 'Телеп', 'Адице', 
             'Ветерник', 'Футог', 'Каћ', 'Петроварадин', 'Сремска Каменица');

-- Удаляем районы Ниша
DELETE FROM districts WHERE city_id = get_city_id_by_name('Ниш') 
AND name IN ('Медијана', 'Палилула', 'Пантелеј', 'Црвени Крст', 'Нишка Бања');

-- Удаляем районы Крагујевца
DELETE FROM districts WHERE city_id = get_city_id_by_name('Крагујевац') 
AND name IN ('Аеродром', 'Бресница', 'Илина вода', 'Пивара', 
             'Станово', 'Стара чаршија', 'Центар');

-- Удаляем районы Суботицы
DELETE FROM districts WHERE city_id = get_city_id_by_name('Суботица') 
AND name IN ('Центар', 'Нови град', 'Дудова шума', 'Кер', 
             'Зорка', 'Пешчара', 'Прозивка', 'Палић');

-- Возвращаем флаг has_districts в false для этих городов
UPDATE cities 
SET has_districts = false 
WHERE name IN ('Нови Сад', 'Ниш', 'Крагујевац', 'Суботица');

-- Удаляем временную функцию
DROP FUNCTION IF EXISTS get_city_id_by_name(TEXT);