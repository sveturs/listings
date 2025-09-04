-- Откат миграции 000037_fix_category_levels
-- Удаляем ограничение
ALTER TABLE marketplace_categories 
DROP CONSTRAINT IF EXISTS check_root_categories_level;

-- Примечание: откат значений level не производится, 
-- так как предыдущие значения были некорректными