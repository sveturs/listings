-- Скрипт для проверки производительности запросов с новыми индексами
-- Запустите этот скрипт после применения миграции 000112

-- Включаем вывод времени выполнения
\timing on

-- 1. Проверка запроса рекурсивного получения категорий с переводами
EXPLAIN (ANALYZE, BUFFERS) 
WITH RECURSIVE category_tree AS (
    SELECT
        c.id,
        c.name,
        c.slug,
        c.icon,
        c.parent_id,
        1 as level,
        (SELECT COUNT(*) FROM marketplace_categories sc WHERE sc.parent_id = c.id) as children_count
    FROM marketplace_categories c
    WHERE c.parent_id IS NULL AND c.is_active = true
    
    UNION ALL
    
    SELECT
        c.id,
        c.name,
        c.slug,
        c.icon,
        c.parent_id,
        ct.level + 1,
        (SELECT COUNT(*) FROM marketplace_categories sc WHERE sc.parent_id = c.id)
    FROM marketplace_categories c
    INNER JOIN category_tree ct ON ct.id = c.parent_id
    WHERE ct.level < 10 AND c.is_active = true
),
categories_with_translations AS (
    SELECT
        ct.*,
        COALESCE(
            jsonb_object_agg(
                t.language,
                t.translated_text
            ) FILTER (WHERE t.language IS NOT NULL),
            '{}'::jsonb
        ) as translations
    FROM category_tree ct
    LEFT JOIN translations t ON
        t.entity_type = 'category'
        AND t.entity_id = ct.id
        AND t.field_name = 'name'
    GROUP BY
        ct.id, ct.name, ct.slug, ct.icon, ct.parent_id,
        ct.level, ct.children_count
)
SELECT * FROM categories_with_translations
ORDER BY level, name;

-- 2. Проверка запроса получения атрибутов категории с переводами
EXPLAIN (ANALYZE, BUFFERS)
WITH attribute_translations AS (
    SELECT
        entity_id,
        jsonb_object_agg(language, translated_text) as translations
    FROM translations
    WHERE entity_type = 'attribute'
        AND field_name = 'display_name'
    GROUP BY entity_id
)
SELECT DISTINCT ON (a.id)
    a.id, a.name, a.display_name, a.attribute_type, a.icon,
    a.options, a.validation_rules, a.is_searchable, a.is_filterable,
    a.is_required, a.sort_order, a.created_at, a.custom_component,
    COALESCE(at.translations, '{}'::jsonb) as translations
FROM category_attribute_mapping m
JOIN category_attributes a ON m.attribute_id = a.id
LEFT JOIN attribute_translations at ON a.id = at.entity_id
WHERE m.category_id = 1 AND m.is_enabled = true
ORDER BY a.id, a.sort_order, a.display_name;

-- 3. Проверка запроса поиска атрибутов по типу
EXPLAIN (ANALYZE, BUFFERS)
SELECT * FROM category_attributes
WHERE attribute_type = 'select'
  AND (is_searchable = true OR is_filterable = true)
ORDER BY sort_order;

-- 4. Проверка запроса значений атрибутов для листинга
EXPLAIN (ANALYZE, BUFFERS)
SELECT 
    lav.*, 
    ca.name as attribute_name,
    ca.display_name,
    ca.attribute_type
FROM listing_attribute_values lav
JOIN category_attributes ca ON lav.attribute_id = ca.id
WHERE lav.listing_id = 100
ORDER BY ca.sort_order;

-- 5. Проверка запроса для получения всех атрибутов с группами
EXPLAIN (ANALYZE, BUFFERS)
WITH all_attributes AS (
    SELECT DISTINCT
        a.id, a.name, a.display_name, a.attribute_type, a.icon, a.options,
        a.validation_rules, a.is_searchable, a.is_filterable, a.is_required,
        a.created_at, a.custom_component,
        cam.sort_order as effective_sort_order,
        'direct' as source
    FROM category_attributes a
    JOIN category_attribute_mapping cam ON a.id = cam.attribute_id
    WHERE cam.category_id = 1 AND cam.is_enabled = true
    
    UNION
    
    SELECT DISTINCT
        a.id, a.name, a.display_name, a.attribute_type, a.icon, a.options,
        a.validation_rules, a.is_searchable, a.is_filterable, a.is_required,
        a.created_at, a.custom_component,
        (cag.sort_order * 1000 + agi.sort_order) as effective_sort_order,
        'group' as source
    FROM category_attributes a
    JOIN attribute_group_items agi ON a.id = agi.attribute_id
    JOIN attribute_groups ag ON agi.group_id = ag.id
    JOIN category_attribute_groups cag ON ag.id = cag.group_id
    WHERE cag.category_id = 1 AND cag.is_active = true
)
SELECT * FROM all_attributes
ORDER BY effective_sort_order, display_name;

-- 6. Проверка запроса переводов для множества сущностей
EXPLAIN (ANALYZE, BUFFERS)
SELECT * FROM translations
WHERE entity_type = 'category'
  AND entity_id IN (1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
  AND field_name = 'name'
  AND language IN ('en', 'ru', 'sr');

-- 7. Проверка рекурсивного запроса для получения всех подкатегорий
EXPLAIN (ANALYZE, BUFFERS)
WITH RECURSIVE category_tree AS (
    SELECT id FROM marketplace_categories WHERE id = 1
    UNION ALL
    SELECT c.id FROM marketplace_categories c
    JOIN category_tree t ON c.parent_id = t.id
    WHERE c.is_active = true
)
SELECT string_agg(id::text, ',') as category_ids FROM category_tree;

-- 8. Проверка запроса для поиска категорий по parent_id
EXPLAIN (ANALYZE, BUFFERS)
SELECT 
    c.*,
    COUNT(sc.id) as subcategory_count,
    COUNT(DISTINCT l.id) as listing_count
FROM marketplace_categories c
LEFT JOIN marketplace_categories sc ON sc.parent_id = c.id AND sc.is_active = true
LEFT JOIN marketplace_listings l ON l.category_id = c.id AND l.is_active = true
WHERE c.parent_id = 1 AND c.is_active = true
GROUP BY c.id
ORDER BY c.sort_order, c.name;

-- Статистика по таблицам
SELECT 
    schemaname,
    tablename,
    n_live_tup as row_count,
    n_dead_tup as dead_rows,
    last_vacuum,
    last_analyze
FROM pg_stat_user_tables
WHERE tablename IN (
    'marketplace_categories',
    'category_attributes', 
    'category_attribute_mapping',
    'translations',
    'listing_attribute_values',
    'attribute_groups',
    'attribute_group_items',
    'category_attribute_groups'
)
ORDER BY tablename;

-- Информация об индексах и их использовании
SELECT 
    schemaname,
    tablename,
    indexname,
    idx_scan as index_scans,
    idx_tup_read as tuples_read,
    idx_tup_fetch as tuples_fetched,
    pg_size_pretty(pg_relation_size(indexrelid)) as index_size
FROM pg_stat_user_indexes
WHERE tablename IN (
    'marketplace_categories',
    'category_attributes', 
    'category_attribute_mapping',
    'translations',
    'listing_attribute_values'
)
ORDER BY tablename, indexname;