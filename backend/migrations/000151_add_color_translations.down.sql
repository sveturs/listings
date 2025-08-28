-- Удаляем переводы для значений цветов
DELETE FROM attribute_option_translations 
WHERE attribute_id = 2004 
AND option_value IN ('grey', 'silver', 'black', 'white', 'blue', 'red', 'green', 'yellow', 'brown', 'purple', 'other');