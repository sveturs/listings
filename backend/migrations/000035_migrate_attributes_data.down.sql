-- Откат миграции данных
-- ВАЖНО: Удаляем только мигрированные данные, старые таблицы остаются

BEGIN;

-- Удаление значений атрибутов
DELETE FROM unified_attribute_values;

-- Удаление связей с категориями
DELETE FROM unified_category_attributes;

-- Удаление атрибутов
DELETE FROM unified_attributes;

-- Сброс последовательностей
ALTER SEQUENCE unified_attributes_id_seq RESTART WITH 1;
ALTER SEQUENCE unified_category_attributes_id_seq RESTART WITH 1;
ALTER SEQUENCE unified_attribute_values_id_seq RESTART WITH 1;

COMMIT;