-- Миграция для исправления уровней категорий
-- Автор: System Architect
-- Дата: 03.09.2025
-- Задача: Исправить некорректные значения в поле level для всех категорий

-- Используем рекурсивный CTE для правильного вычисления уровней
WITH RECURSIVE category_hierarchy AS (
    -- Шаг 1: Выбираем корневые категории (parent_id IS NULL)
    -- Они должны иметь level = 0
    SELECT 
        id, 
        parent_id, 
        0 as correct_level
    FROM marketplace_categories
    WHERE parent_id IS NULL
    
    UNION ALL
    
    -- Шаг 2: Рекурсивно обходим все дочерние категории
    -- Увеличиваем level на 1 для каждого уровня вложенности
    SELECT 
        c.id, 
        c.parent_id, 
        ch.correct_level + 1 as correct_level
    FROM marketplace_categories c
    INNER JOIN category_hierarchy ch ON c.parent_id = ch.id
)
-- Обновляем таблицу marketplace_categories правильными значениями level
UPDATE marketplace_categories mc
SET 
    level = ch.correct_level
FROM category_hierarchy ch
WHERE mc.id = ch.id
  AND mc.level != ch.correct_level; -- Обновляем только если значение изменилось

-- Проверяем результат и выводим статистику
DO $$
DECLARE
    updated_count INTEGER;
    root_count INTEGER;
    level1_count INTEGER;
    level2_count INTEGER;
    level3_count INTEGER;
BEGIN
    -- Подсчитываем количество обновленных записей
    GET DIAGNOSTICS updated_count = ROW_COUNT;
    
    -- Подсчитываем категории по уровням
    SELECT COUNT(*) INTO root_count FROM marketplace_categories WHERE level = 0 AND parent_id IS NULL;
    SELECT COUNT(*) INTO level1_count FROM marketplace_categories WHERE level = 1;
    SELECT COUNT(*) INTO level2_count FROM marketplace_categories WHERE level = 2;
    SELECT COUNT(*) INTO level3_count FROM marketplace_categories WHERE level = 3;
    
    RAISE NOTICE 'Исправление уровней категорий завершено';
    RAISE NOTICE 'Обновлено категорий: %', updated_count;
    RAISE NOTICE 'Корневых категорий (level 0): %', root_count;
    RAISE NOTICE 'Категорий уровня 1: %', level1_count;
    RAISE NOTICE 'Категорий уровня 2: %', level2_count;
    RAISE NOTICE 'Категорий уровня 3: %', level3_count;
END $$;

-- Добавляем проверочное ограничение для предотвращения будущих ошибок
-- Корневые категории должны иметь level = 0
ALTER TABLE marketplace_categories 
ADD CONSTRAINT check_root_categories_level 
CHECK (
    (parent_id IS NULL AND level = 0) OR 
    (parent_id IS NOT NULL AND level > 0)
);