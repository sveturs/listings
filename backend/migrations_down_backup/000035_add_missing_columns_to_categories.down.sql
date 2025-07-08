-- Удаляем индексы
DROP INDEX IF EXISTS idx_marketplace_categories_external_id;
DROP INDEX IF EXISTS idx_marketplace_categories_sort_order;

-- Удаляем колонки
ALTER TABLE marketplace_categories 
    DROP COLUMN IF EXISTS external_id,
    DROP COLUMN IF EXISTS count,
    DROP COLUMN IF EXISTS level,
    DROP COLUMN IF EXISTS sort_order;