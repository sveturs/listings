-- ==========================================
-- Migration: Link Computer Components Attributes to Categories
-- Категория: Računarske komponente
-- Дата: 2026-01-26
-- Описание: Привязка атрибутов к 9 подкатегориям компьютерных компонентов
-- ВАЖНО: Эта миграция должна быть применена ПОСЛЕ создания подкатегорий (Phase 0)
-- ==========================================

BEGIN;

-- ==========================================
-- ПРОВЕРКА: Убедиться что категории существуют
-- ==========================================

DO $$
DECLARE
    missing_categories TEXT[];
    required_categories TEXT[] := ARRAY[
        'graficke-kartice',
        'procesori',
        'ram-memorija',
        'ssd-nakopitelji',
        'hdd-nakopitelji',
        'maticne-ploce',
        'napajanja',
        'kucista',
        'hladjenje'
    ];
    cat TEXT;
    cat_exists INTEGER;
BEGIN
    -- Проверка существования категорий
    FOREACH cat IN ARRAY required_categories
    LOOP
        SELECT COUNT(*) INTO cat_exists FROM categories WHERE slug = cat;
        IF cat_exists = 0 THEN
            missing_categories := array_append(missing_categories, cat);
        END IF;
    END LOOP;

    -- Если категории отсутствуют - вывести ошибку
    IF array_length(missing_categories, 1) > 0 THEN
        RAISE EXCEPTION 'Cannot link attributes: Missing categories: %. Please run Phase 0 migration first to create subcategories.', array_to_string(missing_categories, ', ');
    END IF;

    RAISE NOTICE '✅ All required categories exist. Proceeding with attribute linking...';
END $$;

-- ==========================================
-- 1. ПРИВЯЗКА ОБЩИХ АТРИБУТОВ (ко всем 9 категориям)
-- ==========================================

-- Привязать общие атрибуты ко всем подкатегориям
INSERT INTO category_attributes (category_id, attribute_id, is_required, is_filterable, is_searchable, sort_order)
SELECT
    c.id,
    a.id,
    CASE WHEN a.code IN ('pc_brand', 'pc_model', 'pc_condition') THEN true ELSE false END,
    CASE WHEN a.code IN ('pc_brand', 'pc_condition', 'pc_warranty', 'pc_color', 'pc_rgb_lighting') THEN true ELSE false END,
    CASE WHEN a.code IN ('pc_brand', 'pc_model') THEN true ELSE false END,
    CASE
        WHEN a.code = 'pc_brand' THEN 10
        WHEN a.code = 'pc_model' THEN 20
        WHEN a.code = 'pc_condition' THEN 30
        WHEN a.code = 'pc_warranty' THEN 40
        WHEN a.code = 'pc_color' THEN 50
        WHEN a.code = 'pc_rgb_lighting' THEN 60
    END
FROM categories c
CROSS JOIN attributes a
WHERE c.slug IN (
    'graficke-kartice', 'procesori', 'ram-memorija',
    'ssd-nakopitelji', 'hdd-nakopitelji', 'maticne-ploce',
    'napajanja', 'kucista', 'hladjenje'
)
AND a.code IN ('pc_brand', 'pc_model', 'pc_condition', 'pc_warranty', 'pc_color', 'pc_rgb_lighting')
ON CONFLICT (category_id, attribute_id) DO NOTHING;

-- ==========================================
-- 2. ПРИВЯЗКА АТРИБУТОВ ДЛЯ ВИДЕОКАРТ
-- ==========================================

INSERT INTO category_attributes (category_id, attribute_id, is_required, is_filterable, is_searchable, sort_order)
SELECT
    (SELECT id FROM categories WHERE slug = 'graficke-kartice'),
    a.id,
    CASE WHEN a.code IN ('gpu_series', 'gpu_vram') THEN true ELSE false END,
    CASE WHEN a.code IN ('gpu_series', 'gpu_chip', 'gpu_vram', 'gpu_memory_type', 'gpu_cooling', 'gpu_tdp', 'gpu_ray_tracing') THEN true ELSE false END,
    CASE WHEN a.code IN ('gpu_chip') THEN true ELSE false END,
    CASE
        WHEN a.code = 'gpu_series' THEN 100
        WHEN a.code = 'gpu_chip' THEN 110
        WHEN a.code = 'gpu_vram' THEN 120
        WHEN a.code = 'gpu_memory_type' THEN 130
        WHEN a.code = 'gpu_memory_bus' THEN 140
        WHEN a.code = 'gpu_cooling' THEN 150
        WHEN a.code = 'gpu_power' THEN 160
        WHEN a.code = 'gpu_tdp' THEN 170
        WHEN a.code = 'gpu_length' THEN 180
        WHEN a.code = 'gpu_ray_tracing' THEN 190
    END
FROM attributes a
WHERE a.code LIKE 'gpu_%'
ON CONFLICT (category_id, attribute_id) DO NOTHING;

-- Добавить вариативные атрибуты для видеокарт
INSERT INTO category_variant_attributes (category_id, attribute_id, is_required, affects_price, affects_stock, sort_order)
VALUES
    ((SELECT id FROM categories WHERE slug = 'graficke-kartice'),
     (SELECT id FROM attributes WHERE code = 'gpu_vram'),
     true, true, true, 1)
ON CONFLICT (category_id, attribute_id) DO NOTHING;

-- ==========================================
-- 3. ПРИВЯЗКА АТРИБУТОВ ДЛЯ ПРОЦЕССОРОВ
-- ==========================================

INSERT INTO category_attributes (category_id, attribute_id, is_required, is_filterable, is_searchable, sort_order)
SELECT
    (SELECT id FROM categories WHERE slug = 'procesori'),
    a.id,
    CASE WHEN a.code IN ('cpu_series', 'cpu_socket', 'cpu_cores') THEN true ELSE false END,
    CASE WHEN a.code IN ('cpu_series', 'cpu_socket', 'cpu_cores', 'cpu_tdp', 'cpu_igpu', 'cpu_generation') THEN true ELSE false END,
    false,
    CASE
        WHEN a.code = 'cpu_series' THEN 200
        WHEN a.code = 'cpu_socket' THEN 210
        WHEN a.code = 'cpu_cores' THEN 220
        WHEN a.code = 'cpu_threads' THEN 230
        WHEN a.code = 'cpu_base_clock' THEN 240
        WHEN a.code = 'cpu_boost_clock' THEN 250
        WHEN a.code = 'cpu_tdp' THEN 260
        WHEN a.code = 'cpu_igpu' THEN 270
        WHEN a.code = 'cpu_generation' THEN 280
    END
FROM attributes a
WHERE a.code LIKE 'cpu_%'
ON CONFLICT (category_id, attribute_id) DO NOTHING;

-- Добавить вариативные атрибуты для процессоров
INSERT INTO category_variant_attributes (category_id, attribute_id, is_required, affects_price, affects_stock, sort_order)
VALUES
    ((SELECT id FROM categories WHERE slug = 'procesori'),
     (SELECT id FROM attributes WHERE code = 'cpu_cores'),
     true, true, true, 1)
ON CONFLICT (category_id, attribute_id) DO NOTHING;

-- ==========================================
-- 4. ПРИВЯЗКА АТРИБУТОВ ДЛЯ RAM
-- ==========================================

INSERT INTO category_attributes (category_id, attribute_id, is_required, is_filterable, is_searchable, sort_order)
SELECT
    (SELECT id FROM categories WHERE slug = 'ram-memorija'),
    a.id,
    CASE WHEN a.code IN ('ram_type', 'ram_capacity') THEN true ELSE false END,
    CASE WHEN a.code IN ('ram_type', 'ram_capacity', 'ram_kit', 'ram_speed', 'ram_ecc', 'ram_heatspreader') THEN true ELSE false END,
    false,
    CASE
        WHEN a.code = 'ram_type' THEN 300
        WHEN a.code = 'ram_capacity' THEN 310
        WHEN a.code = 'ram_kit' THEN 320
        WHEN a.code = 'ram_speed' THEN 330
        WHEN a.code = 'ram_cas_latency' THEN 340
        WHEN a.code = 'ram_voltage' THEN 350
        WHEN a.code = 'ram_ecc' THEN 360
        WHEN a.code = 'ram_heatspreader' THEN 370
    END
FROM attributes a
WHERE a.code LIKE 'ram_%'
ON CONFLICT (category_id, attribute_id) DO NOTHING;

-- Добавить вариативные атрибуты для RAM
INSERT INTO category_variant_attributes (category_id, attribute_id, is_required, affects_price, affects_stock, sort_order)
VALUES
    ((SELECT id FROM categories WHERE slug = 'ram-memorija'),
     (SELECT id FROM attributes WHERE code = 'ram_type'),
     true, true, true, 1),
    ((SELECT id FROM categories WHERE slug = 'ram-memorija'),
     (SELECT id FROM attributes WHERE code = 'ram_capacity'),
     true, true, true, 2)
ON CONFLICT (category_id, attribute_id) DO NOTHING;

-- ==========================================
-- 5. ПРИВЯЗКА АТРИБУТОВ ДЛЯ SSD
-- ==========================================

INSERT INTO category_attributes (category_id, attribute_id, is_required, is_filterable, is_searchable, sort_order)
SELECT
    (SELECT id FROM categories WHERE slug = 'ssd-nakopitelji'),
    a.id,
    CASE WHEN a.code IN ('ssd_capacity', 'ssd_interface') THEN true ELSE false END,
    CASE WHEN a.code IN ('ssd_capacity', 'ssd_interface', 'ssd_form_factor', 'ssd_nand_type', 'ssd_dram_cache') THEN true ELSE false END,
    false,
    CASE
        WHEN a.code = 'ssd_capacity' THEN 400
        WHEN a.code = 'ssd_interface' THEN 410
        WHEN a.code = 'ssd_form_factor' THEN 420
        WHEN a.code = 'ssd_read_speed' THEN 430
        WHEN a.code = 'ssd_write_speed' THEN 440
        WHEN a.code = 'ssd_nand_type' THEN 450
        WHEN a.code = 'ssd_dram_cache' THEN 460
        WHEN a.code = 'ssd_endurance' THEN 470
    END
FROM attributes a
WHERE a.code LIKE 'ssd_%'
ON CONFLICT (category_id, attribute_id) DO NOTHING;

-- Добавить вариативные атрибуты для SSD
INSERT INTO category_variant_attributes (category_id, attribute_id, is_required, affects_price, affects_stock, sort_order)
VALUES
    ((SELECT id FROM categories WHERE slug = 'ssd-nakopitelji'),
     (SELECT id FROM attributes WHERE code = 'ssd_capacity'),
     true, true, true, 1),
    ((SELECT id FROM categories WHERE slug = 'ssd-nakopitelji'),
     (SELECT id FROM attributes WHERE code = 'ssd_interface'),
     true, true, true, 2)
ON CONFLICT (category_id, attribute_id) DO NOTHING;

-- ==========================================
-- 6. ПРИВЯЗКА АТРИБУТОВ ДЛЯ HDD
-- ==========================================

INSERT INTO category_attributes (category_id, attribute_id, is_required, is_filterable, is_searchable, sort_order)
SELECT
    (SELECT id FROM categories WHERE slug = 'hdd-nakopitelji'),
    a.id,
    CASE WHEN a.code IN ('hdd_capacity', 'hdd_rpm') THEN true ELSE false END,
    CASE WHEN a.code IN ('hdd_capacity', 'hdd_rpm', 'hdd_cache', 'hdd_form_factor', 'hdd_usage') THEN true ELSE false END,
    false,
    CASE
        WHEN a.code = 'hdd_capacity' THEN 500
        WHEN a.code = 'hdd_rpm' THEN 510
        WHEN a.code = 'hdd_cache' THEN 520
        WHEN a.code = 'hdd_interface' THEN 530
        WHEN a.code = 'hdd_form_factor' THEN 540
        WHEN a.code = 'hdd_usage' THEN 550
    END
FROM attributes a
WHERE a.code LIKE 'hdd_%'
ON CONFLICT (category_id, attribute_id) DO NOTHING;

-- Добавить вариативные атрибуты для HDD
INSERT INTO category_variant_attributes (category_id, attribute_id, is_required, affects_price, affects_stock, sort_order)
VALUES
    ((SELECT id FROM categories WHERE slug = 'hdd-nakopitelji'),
     (SELECT id FROM attributes WHERE code = 'hdd_capacity'),
     true, true, true, 1),
    ((SELECT id FROM categories WHERE slug = 'hdd-nakopitelji'),
     (SELECT id FROM attributes WHERE code = 'hdd_rpm'),
     true, true, true, 2)
ON CONFLICT (category_id, attribute_id) DO NOTHING;

-- ==========================================
-- 7. ПРИВЯЗКА АТРИБУТОВ ДЛЯ МАТЕРИНСКИХ ПЛАТ
-- ==========================================

INSERT INTO category_attributes (category_id, attribute_id, is_required, is_filterable, is_searchable, sort_order)
SELECT
    (SELECT id FROM categories WHERE slug = 'maticne-ploce'),
    a.id,
    CASE WHEN a.code IN ('mb_socket', 'mb_chipset', 'mb_form_factor', 'mb_memory_type') THEN true ELSE false END,
    CASE WHEN a.code IN ('mb_socket', 'mb_chipset', 'mb_form_factor', 'mb_memory_type', 'mb_memory_slots', 'mb_wifi', 'mb_bluetooth') THEN true ELSE false END,
    false,
    CASE
        WHEN a.code = 'mb_socket' THEN 600
        WHEN a.code = 'mb_chipset' THEN 610
        WHEN a.code = 'mb_form_factor' THEN 620
        WHEN a.code = 'mb_memory_type' THEN 630
        WHEN a.code = 'mb_memory_slots' THEN 640
        WHEN a.code = 'mb_max_memory' THEN 650
        WHEN a.code = 'mb_m2_slots' THEN 660
        WHEN a.code = 'mb_sata_ports' THEN 670
        WHEN a.code = 'mb_wifi' THEN 680
        WHEN a.code = 'mb_bluetooth' THEN 690
    END
FROM attributes a
WHERE a.code LIKE 'mb_%'
ON CONFLICT (category_id, attribute_id) DO NOTHING;

-- Добавить вариативные атрибуты для материнских плат
INSERT INTO category_variant_attributes (category_id, attribute_id, is_required, affects_price, affects_stock, sort_order)
VALUES
    ((SELECT id FROM categories WHERE slug = 'maticne-ploce'),
     (SELECT id FROM attributes WHERE code = 'mb_chipset'),
     true, true, true, 1),
    ((SELECT id FROM categories WHERE slug = 'maticne-ploce'),
     (SELECT id FROM attributes WHERE code = 'mb_form_factor'),
     true, true, true, 2)
ON CONFLICT (category_id, attribute_id) DO NOTHING;

-- ==========================================
-- 8. ПРИВЯЗКА АТРИБУТОВ ДЛЯ PSU
-- ==========================================

INSERT INTO category_attributes (category_id, attribute_id, is_required, is_filterable, is_searchable, sort_order)
SELECT
    (SELECT id FROM categories WHERE slug = 'napajanja'),
    a.id,
    CASE WHEN a.code IN ('psu_wattage', 'psu_efficiency') THEN true ELSE false END,
    CASE WHEN a.code IN ('psu_wattage', 'psu_efficiency', 'psu_modular', 'psu_form_factor', 'psu_pcie5') THEN true ELSE false END,
    false,
    CASE
        WHEN a.code = 'psu_wattage' THEN 700
        WHEN a.code = 'psu_efficiency' THEN 710
        WHEN a.code = 'psu_modular' THEN 720
        WHEN a.code = 'psu_form_factor' THEN 730
        WHEN a.code = 'psu_pcie5' THEN 740
        WHEN a.code = 'psu_fan_size' THEN 750
    END
FROM attributes a
WHERE a.code LIKE 'psu_%'
ON CONFLICT (category_id, attribute_id) DO NOTHING;

-- Добавить вариативные атрибуты для PSU
INSERT INTO category_variant_attributes (category_id, attribute_id, is_required, affects_price, affects_stock, sort_order)
VALUES
    ((SELECT id FROM categories WHERE slug = 'napajanja'),
     (SELECT id FROM attributes WHERE code = 'psu_wattage'),
     true, true, true, 1),
    ((SELECT id FROM categories WHERE slug = 'napajanja'),
     (SELECT id FROM attributes WHERE code = 'psu_efficiency'),
     true, true, true, 2)
ON CONFLICT (category_id, attribute_id) DO NOTHING;

-- ==========================================
-- 9. ПРИВЯЗКА АТРИБУТОВ ДЛЯ КОРПУСОВ
-- ==========================================

INSERT INTO category_attributes (category_id, attribute_id, is_required, is_filterable, is_searchable, sort_order)
SELECT
    (SELECT id FROM categories WHERE slug = 'kucista'),
    a.id,
    CASE WHEN a.code IN ('case_form_factor') THEN true ELSE false END,
    CASE WHEN a.code IN ('case_form_factor', 'case_max_gpu_length', 'case_fans_included', 'case_tempered_glass', 'case_rgb_fans', 'case_dust_filters') THEN true ELSE false END,
    false,
    CASE
        WHEN a.code = 'case_form_factor' THEN 800
        WHEN a.code = 'case_max_gpu_length' THEN 810
        WHEN a.code = 'case_max_cpu_height' THEN 820
        WHEN a.code = 'case_fans_included' THEN 830
        WHEN a.code = 'case_max_fans' THEN 840
        WHEN a.code = 'case_tempered_glass' THEN 850
        WHEN a.code = 'case_rgb_fans' THEN 860
        WHEN a.code = 'case_dust_filters' THEN 870
    END
FROM attributes a
WHERE a.code LIKE 'case_%'
ON CONFLICT (category_id, attribute_id) DO NOTHING;

-- Добавить вариативные атрибуты для корпусов
INSERT INTO category_variant_attributes (category_id, attribute_id, is_required, affects_price, affects_stock, sort_order)
VALUES
    ((SELECT id FROM categories WHERE slug = 'kucista'),
     (SELECT id FROM attributes WHERE code = 'case_form_factor'),
     true, true, true, 1)
ON CONFLICT (category_id, attribute_id) DO NOTHING;

-- ==========================================
-- 10. ПРИВЯЗКА АТРИБУТОВ ДЛЯ ОХЛАЖДЕНИЯ
-- ==========================================

INSERT INTO category_attributes (category_id, attribute_id, is_required, is_filterable, is_searchable, sort_order)
SELECT
    (SELECT id FROM categories WHERE slug = 'hladjenje'),
    a.id,
    CASE WHEN a.code IN ('cooler_type') THEN true ELSE false END,
    CASE WHEN a.code IN ('cooler_type', 'cooler_radiator', 'cooler_fan_size', 'cooler_max_tdp', 'cooler_rgb') THEN true ELSE false END,
    false,
    CASE
        WHEN a.code = 'cooler_type' THEN 900
        WHEN a.code = 'cooler_radiator' THEN 910
        WHEN a.code = 'cooler_fan_size' THEN 920
        WHEN a.code = 'cooler_max_tdp' THEN 930
        WHEN a.code = 'cooler_rgb' THEN 940
        WHEN a.code = 'cooler_noise' THEN 950
    END
FROM attributes a
WHERE a.code LIKE 'cooler_%'
ON CONFLICT (category_id, attribute_id) DO NOTHING;

-- Добавить вариативные атрибуты для охлаждения
INSERT INTO category_variant_attributes (category_id, attribute_id, is_required, affects_price, affects_stock, sort_order)
VALUES
    ((SELECT id FROM categories WHERE slug = 'hladjenje'),
     (SELECT id FROM attributes WHERE code = 'cooler_type'),
     true, true, true, 1)
ON CONFLICT (category_id, attribute_id) DO NOTHING;

COMMIT;

-- ==========================================
-- ИТОГОВАЯ СТАТИСТИКА
-- ==========================================

DO $$
DECLARE
    total_links INTEGER;
    total_variant_links INTEGER;
    graficke_kartice_count INTEGER;
    procesori_count INTEGER;
    ram_count INTEGER;
    ssd_count INTEGER;
    hdd_count INTEGER;
    maticne_ploce_count INTEGER;
    napajanja_count INTEGER;
    kucista_count INTEGER;
    hladjenje_count INTEGER;
BEGIN
    -- Подсчет связей категория-атрибут
    SELECT COUNT(*) INTO graficke_kartice_count
    FROM category_attributes ca
    JOIN categories c ON ca.category_id = c.id
    WHERE c.slug = 'graficke-kartice';

    SELECT COUNT(*) INTO procesori_count
    FROM category_attributes ca
    JOIN categories c ON ca.category_id = c.id
    WHERE c.slug = 'procesori';

    SELECT COUNT(*) INTO ram_count
    FROM category_attributes ca
    JOIN categories c ON ca.category_id = c.id
    WHERE c.slug = 'ram-memorija';

    SELECT COUNT(*) INTO ssd_count
    FROM category_attributes ca
    JOIN categories c ON ca.category_id = c.id
    WHERE c.slug = 'ssd-nakopitelji';

    SELECT COUNT(*) INTO hdd_count
    FROM category_attributes ca
    JOIN categories c ON ca.category_id = c.id
    WHERE c.slug = 'hdd-nakopitelji';

    SELECT COUNT(*) INTO maticne_ploce_count
    FROM category_attributes ca
    JOIN categories c ON ca.category_id = c.id
    WHERE c.slug = 'maticne-ploce';

    SELECT COUNT(*) INTO napajanja_count
    FROM category_attributes ca
    JOIN categories c ON ca.category_id = c.id
    WHERE c.slug = 'napajanja';

    SELECT COUNT(*) INTO kucista_count
    FROM category_attributes ca
    JOIN categories c ON ca.category_id = c.id
    WHERE c.slug = 'kucista';

    SELECT COUNT(*) INTO hladjenje_count
    FROM category_attributes ca
    JOIN categories c ON ca.category_id = c.id
    WHERE c.slug = 'hladjenje';

    total_links := graficke_kartice_count + procesori_count + ram_count + ssd_count +
                   hdd_count + maticne_ploce_count + napajanja_count + kucista_count + hladjenje_count;

    -- Подсчет вариативных связей
    SELECT COUNT(*) INTO total_variant_links FROM category_variant_attributes cva
    JOIN categories c ON cva.category_id = c.id
    WHERE c.slug IN (
        'graficke-kartice', 'procesori', 'ram-memorija',
        'ssd-nakopitelji', 'hdd-nakopitelji', 'maticne-ploce',
        'napajanja', 'kucista', 'hladjenje'
    );

    RAISE NOTICE '✅ Attributes linked successfully!';
    RAISE NOTICE '';
    RAISE NOTICE 'Total attribute links created: %', total_links;
    RAISE NOTICE '  - Grafičke kartice: % attributes', graficke_kartice_count;
    RAISE NOTICE '  - Procesori: % attributes', procesori_count;
    RAISE NOTICE '  - RAM memorija: % attributes', ram_count;
    RAISE NOTICE '  - SSD nakopitelji: % attributes', ssd_count;
    RAISE NOTICE '  - HDD nakopitelji: % attributes', hdd_count;
    RAISE NOTICE '  - Matične ploče: % attributes', maticne_ploce_count;
    RAISE NOTICE '  - Napajanja: % attributes', napajanja_count;
    RAISE NOTICE '  - Kućišta: % attributes', kucista_count;
    RAISE NOTICE '  - Hlađenje: % attributes', hladjenje_count;
    RAISE NOTICE '';
    RAISE NOTICE 'Total variant attribute links: %', total_variant_links;
    RAISE NOTICE '';
    RAISE NOTICE '✅ Phase 1 completed! Attributes are now linked to categories.';
END $$;
