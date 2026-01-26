-- ==========================================
-- ROLLBACK: Откат переструктурирования категории "Računarske komponente"
-- Дата: 2026-01-26
-- Описание: Удаление 9 общих категорий
-- ВАЖНО: Старые 20 конкретных категорий НЕ восстанавливаются (удалены в up миграции)
-- ==========================================

BEGIN;

DO $$
DECLARE
    deleted_cats INTEGER;
BEGIN
    RAISE NOTICE '';
    RAISE NOTICE '⏪ Rolling back Phase 0...';
    RAISE NOTICE '';

    -- Удалить связи атрибутов новых категорий (если были созданы вручную)
    DELETE FROM category_attributes WHERE category_id IN (
        SELECT id FROM categories WHERE slug IN (
            'graficke-kartice', 'procesori', 'ram-memorija',
            'ssd-nakopitelji', 'hdd-nakopitelji', 'maticne-ploce',
            'napajanja', 'kucista', 'hladjenje'
        )
    );

    DELETE FROM category_variant_attributes WHERE category_id IN (
        SELECT id FROM categories WHERE slug IN (
            'graficke-kartice', 'procesori', 'ram-memorija',
            'ssd-nakopitelji', 'hdd-nakopitelji', 'maticne-ploce',
            'napajanja', 'kucista', 'hladjenje'
        )
    );

    -- Удалить 9 новых общих категорий
    DELETE FROM categories WHERE slug IN (
        'graficke-kartice', 'procesori', 'ram-memorija',
        'ssd-nakopitelji', 'hdd-nakopitelji', 'maticne-ploce',
        'napajanja', 'kucista', 'hladjenje'
    );
    GET DIAGNOSTICS deleted_cats = ROW_COUNT;

    RAISE NOTICE '✅ Deleted % new general categories', deleted_cats;
    RAISE NOTICE '';
    RAISE NOTICE '⚠️  NOTE: Old 20 concrete subcategories are NOT restored!';
    RAISE NOTICE '⚠️  They were removed during migration up.';
    RAISE NOTICE '⚠️  If you need them back, restore from database backup.';
    RAISE NOTICE '';
    RAISE NOTICE '✅ Rollback completed.';
    RAISE NOTICE '';
END $$;

COMMIT;
