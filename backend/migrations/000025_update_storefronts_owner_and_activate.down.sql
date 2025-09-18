-- Откат миграции 000025

-- 1. Деактивируем витрины которые были активированы
UPDATE storefronts
SET is_active = false
WHERE id IN (36, 38, 39);

-- 2. Удаляем добавленные колонки (если они были добавлены этой миграцией)
-- Проверяем существование колонок перед удалением
ALTER TABLE storefronts
DROP COLUMN IF EXISTS followers_count;

ALTER TABLE storefronts
DROP COLUMN IF EXISTS rating;

-- Примечание: откат владельцев невозможен, так как мы не знаем исходных владельцев