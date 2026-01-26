-- ==========================================
-- Migration: Unlink Computer Components Attributes from Categories
-- Категория: Računarske komponente
-- Дата: 2026-01-26
-- Описание: Удаление связей атрибутов с 9 подкатегориями
-- ==========================================

BEGIN;

-- ==========================================
-- Удаление всех связей категория-атрибут
-- ==========================================

-- Удалить из category_attributes
DELETE FROM category_attributes
WHERE category_id IN (
    SELECT id FROM categories WHERE slug IN (
        'graficke-kartice',
        'procesori',
        'ram-memorija',
        'ssd-nakopitelji',
        'hdd-nakopitelji',
        'maticne-ploce',
        'napajanja',
        'kucista',
        'hladjenje'
    )
)
AND attribute_id IN (
    SELECT id FROM attributes WHERE code LIKE ANY(ARRAY[
        'pc_%',
        'gpu_%',
        'cpu_%',
        'ram_%',
        'ssd_%',
        'hdd_%',
        'mb_%',
        'psu_%',
        'case_%',
        'cooler_%'
    ])
);

-- Удалить из category_variant_attributes
DELETE FROM category_variant_attributes
WHERE category_id IN (
    SELECT id FROM categories WHERE slug IN (
        'graficke-kartice',
        'procesori',
        'ram-memorija',
        'ssd-nakopitelji',
        'hdd-nakopitelji',
        'maticne-ploce',
        'napajanja',
        'kucista',
        'hladjenje'
    )
)
AND attribute_id IN (
    SELECT id FROM attributes WHERE code IN (
        'gpu_vram',           -- Variant
        'cpu_cores',          -- Variant
        'ram_type',           -- Variant
        'ram_capacity',       -- Variant
        'ssd_capacity',       -- Variant
        'ssd_interface',      -- Variant
        'hdd_capacity',       -- Variant
        'hdd_rpm',            -- Variant
        'mb_chipset',         -- Variant
        'mb_form_factor',     -- Variant
        'psu_wattage',        -- Variant
        'psu_efficiency',     -- Variant
        'case_form_factor',   -- Variant
        'cooler_type'         -- Variant
    )
);

COMMIT;

-- Итоговое сообщение
DO $$
BEGIN
    RAISE NOTICE '✅ Rollback completed successfully!';
    RAISE NOTICE '';
    RAISE NOTICE 'All attribute links have been removed from computer components categories.';
END $$;
