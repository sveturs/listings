-- Удаление координат центров районов

-- Обновляем все районы, устанавливая center_point в NULL
UPDATE districts 
SET center_point = NULL
WHERE city_id IN (
    SELECT id FROM cities 
    WHERE name IN ('Нови Сад', 'Ниш', 'Крагујевац', 'Суботица')
);