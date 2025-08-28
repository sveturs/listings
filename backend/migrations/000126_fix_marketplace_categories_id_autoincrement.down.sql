-- Удаляем значение по умолчанию
ALTER TABLE marketplace_categories 
    ALTER COLUMN id DROP DEFAULT;

-- Удаляем последовательность
DROP SEQUENCE IF EXISTS marketplace_categories_id_seq;