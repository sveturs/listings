-- Добавление координат центров районов

-- Нови Сад
UPDATE districts 
SET center_point = ST_SetSRID(ST_Point(19.8227, 45.2438), 4326)
WHERE name = 'Лиман' AND city_id = (SELECT id FROM cities WHERE name = 'Нови Сад');

UPDATE districts 
SET center_point = ST_SetSRID(ST_Point(19.7863, 45.2397), 4326)
WHERE name = 'Ново насеље' AND city_id = (SELECT id FROM cities WHERE name = 'Нови Сад');

UPDATE districts 
SET center_point = ST_SetSRID(ST_Point(19.8285, 45.2544), 4326)
WHERE name = 'Грбавица' AND city_id = (SELECT id FROM cities WHERE name = 'Нови Сад');

UPDATE districts 
SET center_point = ST_SetSRID(ST_Point(19.8050, 45.2790), 4326)
WHERE name = 'Детелинара' AND city_id = (SELECT id FROM cities WHERE name = 'Нови Сад');

UPDATE districts 
SET center_point = ST_SetSRID(ST_Point(19.7741, 45.2324), 4326)
WHERE name = 'Клиса' AND city_id = (SELECT id FROM cities WHERE name = 'Нови Сад');

UPDATE districts 
SET center_point = ST_SetSRID(ST_Point(19.8433, 45.2441), 4326)
WHERE name = 'Подбара' AND city_id = (SELECT id FROM cities WHERE name = 'Нови Сад');

UPDATE districts 
SET center_point = ST_SetSRID(ST_Point(19.7912, 45.2203), 4326)
WHERE name = 'Салајка' AND city_id = (SELECT id FROM cities WHERE name = 'Нови Сад');

UPDATE districts 
SET center_point = ST_SetSRID(ST_Point(19.8188, 45.2688), 4326)
WHERE name = 'Роткварија' AND city_id = (SELECT id FROM cities WHERE name = 'Нови Сад');

UPDATE districts 
SET center_point = ST_SetSRID(ST_Point(19.8044, 45.2290), 4326)
WHERE name = 'Телеп' AND city_id = (SELECT id FROM cities WHERE name = 'Нови Сад');

UPDATE districts 
SET center_point = ST_SetSRID(ST_Point(19.8300, 45.2199), 4326)
WHERE name = 'Адице' AND city_id = (SELECT id FROM cities WHERE name = 'Нови Сад');

UPDATE districts 
SET center_point = ST_SetSRID(ST_Point(19.8418, 45.2913), 4326)
WHERE name = 'Ветерник' AND city_id = (SELECT id FROM cities WHERE name = 'Нови Сад');

UPDATE districts 
SET center_point = ST_SetSRID(ST_Point(19.7337, 45.2385), 4326)
WHERE name = 'Футог' AND city_id = (SELECT id FROM cities WHERE name = 'Нови Сад');

UPDATE districts 
SET center_point = ST_SetSRID(ST_Point(19.8950, 45.3202), 4326)
WHERE name = 'Каћ' AND city_id = (SELECT id FROM cities WHERE name = 'Нови Сад');

UPDATE districts 
SET center_point = ST_SetSRID(ST_Point(19.8661, 45.2486), 4326)
WHERE name = 'Петроварадин' AND city_id = (SELECT id FROM cities WHERE name = 'Нови Сад');

UPDATE districts 
SET center_point = ST_SetSRID(ST_Point(19.8932, 45.2173), 4326)
WHERE name = 'Сремска Каменица' AND city_id = (SELECT id FROM cities WHERE name = 'Нови Сад');

-- Ниш
UPDATE districts 
SET center_point = ST_SetSRID(ST_Point(21.9089, 43.3035), 4326)
WHERE name = 'Медијана' AND city_id = (SELECT id FROM cities WHERE name = 'Ниш');

UPDATE districts 
SET center_point = ST_SetSRID(ST_Point(21.8737, 43.3374), 4326)
WHERE name = 'Палилула' AND city_id = (SELECT id FROM cities WHERE name = 'Ниш');

UPDATE districts 
SET center_point = ST_SetSRID(ST_Point(21.9245, 43.3437), 4326)
WHERE name = 'Пантелеј' AND city_id = (SELECT id FROM cities WHERE name = 'Ниш');

UPDATE districts 
SET center_point = ST_SetSRID(ST_Point(21.8778, 43.3088), 4326)
WHERE name = 'Црвени Крст' AND city_id = (SELECT id FROM cities WHERE name = 'Ниш');

UPDATE districts 
SET center_point = ST_SetSRID(ST_Point(21.9558, 43.2827), 4326)
WHERE name = 'Нишка Бања' AND city_id = (SELECT id FROM cities WHERE name = 'Ниш');

-- Крагујевац
UPDATE districts 
SET center_point = ST_SetSRID(ST_Point(20.9114, 44.0165), 4326)
WHERE name = 'Центар' AND city_id = (SELECT id FROM cities WHERE name = 'Крагујевац');

UPDATE districts 
SET center_point = ST_SetSRID(ST_Point(20.9375, 44.0092), 4326)
WHERE name = 'Аеродром' AND city_id = (SELECT id FROM cities WHERE name = 'Крагујевац');

UPDATE districts 
SET center_point = ST_SetSRID(ST_Point(20.9226, 44.0356), 4326)
WHERE name = 'Станово' AND city_id = (SELECT id FROM cities WHERE name = 'Крагујевац');

UPDATE districts 
SET center_point = ST_SetSRID(ST_Point(20.9071, 44.0180), 4326)
WHERE name = 'Стара чаршија' AND city_id = (SELECT id FROM cities WHERE name = 'Крагујевац');

UPDATE districts 
SET center_point = ST_SetSRID(ST_Point(20.8934, 43.9950), 4326)
WHERE name = 'Илина вода' AND city_id = (SELECT id FROM cities WHERE name = 'Крагујевац');

UPDATE districts 
SET center_point = ST_SetSRID(ST_Point(20.8775, 44.0487), 4326)
WHERE name = 'Бресница' AND city_id = (SELECT id FROM cities WHERE name = 'Крагујевац');

UPDATE districts 
SET center_point = ST_SetSRID(ST_Point(20.9273, 44.0039), 4326)
WHERE name = 'Пивара' AND city_id = (SELECT id FROM cities WHERE name = 'Крагујевац');

-- Суботица
UPDATE districts 
SET center_point = ST_SetSRID(ST_Point(19.6677, 46.1008), 4326)
WHERE name = 'Центар' AND city_id = (SELECT id FROM cities WHERE name = 'Суботица');

UPDATE districts 
SET center_point = ST_SetSRID(ST_Point(19.6542, 46.0766), 4326)
WHERE name = 'Нови град' AND city_id = (SELECT id FROM cities WHERE name = 'Суботица');

UPDATE districts 
SET center_point = ST_SetSRID(ST_Point(19.7642, 46.0951), 4326)
WHERE name = 'Палић' AND city_id = (SELECT id FROM cities WHERE name = 'Суботица');

UPDATE districts 
SET center_point = ST_SetSRID(ST_Point(19.6338, 46.0618), 4326)
WHERE name = 'Дудова шума' AND city_id = (SELECT id FROM cities WHERE name = 'Суботица');

UPDATE districts 
SET center_point = ST_SetSRID(ST_Point(19.6892, 46.1158), 4326)
WHERE name = 'Кер' AND city_id = (SELECT id FROM cities WHERE name = 'Суботица');

UPDATE districts 
SET center_point = ST_SetSRID(ST_Point(19.6408, 46.0938), 4326)
WHERE name = 'Зорка' AND city_id = (SELECT id FROM cities WHERE name = 'Суботица');

UPDATE districts 
SET center_point = ST_SetSRID(ST_Point(19.6595, 46.1251), 4326)
WHERE name = 'Пешчара' AND city_id = (SELECT id FROM cities WHERE name = 'Суботица');

UPDATE districts 
SET center_point = ST_SetSRID(ST_Point(19.6803, 46.1374), 4326)
WHERE name = 'Прозивка' AND city_id = (SELECT id FROM cities WHERE name = 'Суботица');