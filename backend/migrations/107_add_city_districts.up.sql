-- Миграция для добавления районов крупных городов Сербии

-- Сначала обновляем флаг has_districts для городов с районами
UPDATE cities 
SET has_districts = true 
WHERE name IN ('Нови Сад', 'Ниш', 'Крагујевац', 'Суботица');

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

-- Добавляем районы Нови Сада
INSERT INTO districts (city_id, name) VALUES
(get_city_id_by_name('Нови Сад'), 'Лиман'),
(get_city_id_by_name('Нови Сад'), 'Ново насеље'),
(get_city_id_by_name('Нови Сад'), 'Грбавица'),
(get_city_id_by_name('Нови Сад'), 'Детелинара'),
(get_city_id_by_name('Нови Сад'), 'Клиса'),
(get_city_id_by_name('Нови Сад'), 'Подбара'),
(get_city_id_by_name('Нови Сад'), 'Салајка'),
(get_city_id_by_name('Нови Сад'), 'Роткварија'),
(get_city_id_by_name('Нови Сад'), 'Телеп'),
(get_city_id_by_name('Нови Сад'), 'Адице'),
(get_city_id_by_name('Нови Сад'), 'Ветерник'),
(get_city_id_by_name('Нови Сад'), 'Футог'),
(get_city_id_by_name('Нови Сад'), 'Каћ'),
(get_city_id_by_name('Нови Сад'), 'Петроварадин'),
(get_city_id_by_name('Нови Сад'), 'Сремска Каменица');

-- Добавляем районы Ниша
INSERT INTO districts (city_id, name) VALUES
(get_city_id_by_name('Ниш'), 'Медијана'),
(get_city_id_by_name('Ниш'), 'Палилула'),
(get_city_id_by_name('Ниш'), 'Пантелеј'),
(get_city_id_by_name('Ниш'), 'Црвени Крст'),
(get_city_id_by_name('Ниш'), 'Нишка Бања');

-- Добавляем районы Крагујевца
INSERT INTO districts (city_id, name) VALUES
(get_city_id_by_name('Крагујевац'), 'Аеродром'),
(get_city_id_by_name('Крагујевац'), 'Бресница'),
(get_city_id_by_name('Крагујевац'), 'Илина вода'),
(get_city_id_by_name('Крагујевац'), 'Пивара'),
(get_city_id_by_name('Крагујевац'), 'Станово'),
(get_city_id_by_name('Крагујевац'), 'Стара чаршија'),
(get_city_id_by_name('Крагујевац'), 'Центар');

-- Добавляем районы Суботицы
INSERT INTO districts (city_id, name) VALUES
(get_city_id_by_name('Суботица'), 'Центар'),
(get_city_id_by_name('Суботица'), 'Нови град'),
(get_city_id_by_name('Суботица'), 'Дудова шума'),
(get_city_id_by_name('Суботица'), 'Кер'),
(get_city_id_by_name('Суботица'), 'Зорка'),
(get_city_id_by_name('Суботица'), 'Пешчара'),
(get_city_id_by_name('Суботица'), 'Прозивка'),
(get_city_id_by_name('Суботица'), 'Палић');

-- Удаляем временную функцию
DROP FUNCTION IF EXISTS get_city_id_by_name(TEXT);