-- Rollback migration: Remove tools and measuring devices AI mappings

-- Удалить переводы категории (используем правильную структуру)
DELETE FROM translations
WHERE entity_type = 'category'
AND entity_id IN (
    SELECT id FROM marketplace_categories
    WHERE slug = 'construction-tools'
);

-- Получить ID категории строительных инструментов
DO $$
DECLARE
    tools_category_id INTEGER;
BEGIN
    SELECT id INTO tools_category_id
    FROM marketplace_categories
    WHERE slug = 'construction-tools';

    IF tools_category_id IS NOT NULL THEN
        -- Удалить keyword weights
        DELETE FROM category_keyword_weights
        WHERE category_id = tools_category_id;

        -- Удалить AI mappings
        DELETE FROM category_ai_mappings
        WHERE category_id = tools_category_id;

        -- Удалить категорию (только если она была создана этой миграцией)
        DELETE FROM marketplace_categories
        WHERE id = tools_category_id AND slug = 'construction-tools';
    END IF;
END $$;