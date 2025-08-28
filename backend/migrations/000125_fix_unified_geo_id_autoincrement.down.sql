-- Откат изменений для автоинкремента id в таблице unified_geo

-- Удаляем значение по умолчанию
ALTER TABLE unified_geo 
    ALTER COLUMN id DROP DEFAULT;

-- Удаляем последовательность
DROP SEQUENCE IF EXISTS unified_geo_id_seq;