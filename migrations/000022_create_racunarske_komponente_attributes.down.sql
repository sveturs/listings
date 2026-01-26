-- ==========================================
-- Migration: Rollback Computer Components Attributes
-- Категория: Računarske komponente
-- Дата: 2026-01-26
-- Описание: Удаление 84+ атрибутов компьютерных компонентов
-- ==========================================

BEGIN;

-- ==========================================
-- Удаление всех атрибутов компьютерных компонентов
-- ==========================================

-- Сначала удалить все связи с категориями (если существуют)
DELETE FROM category_attributes
WHERE attribute_id IN (
    SELECT id FROM attributes WHERE code IN (
        -- Общие атрибуты
        'pc_brand', 'pc_model', 'pc_condition', 'pc_warranty', 'pc_color', 'pc_rgb_lighting',

        -- Видеокарты (10 атрибутов)
        'gpu_series', 'gpu_chip', 'gpu_vram', 'gpu_memory_type', 'gpu_memory_bus',
        'gpu_cooling', 'gpu_power', 'gpu_tdp', 'gpu_length', 'gpu_ray_tracing',

        -- Процессоры (9 атрибутов)
        'cpu_series', 'cpu_socket', 'cpu_cores', 'cpu_threads', 'cpu_base_clock',
        'cpu_boost_clock', 'cpu_tdp', 'cpu_igpu', 'cpu_generation',

        -- RAM (8 атрибутов)
        'ram_type', 'ram_capacity', 'ram_kit', 'ram_speed', 'ram_cas_latency',
        'ram_voltage', 'ram_ecc', 'ram_heatspreader',

        -- SSD (8 атрибутов)
        'ssd_capacity', 'ssd_interface', 'ssd_form_factor', 'ssd_read_speed',
        'ssd_write_speed', 'ssd_nand_type', 'ssd_dram_cache', 'ssd_endurance',

        -- HDD (6 атрибутов)
        'hdd_capacity', 'hdd_rpm', 'hdd_cache', 'hdd_interface',
        'hdd_form_factor', 'hdd_usage',

        -- Материнские платы (11 атрибутов)
        'mb_socket', 'mb_chipset', 'mb_form_factor', 'mb_memory_type',
        'mb_memory_slots', 'mb_max_memory', 'mb_m2_slots', 'mb_sata_ports',
        'mb_wifi', 'mb_bluetooth',

        -- PSU (6 атрибутов)
        'psu_wattage', 'psu_efficiency', 'psu_modular', 'psu_form_factor',
        'psu_pcie5', 'psu_fan_size',

        -- Корпуса (8 атрибутов)
        'case_form_factor', 'case_max_gpu_length', 'case_max_cpu_height',
        'case_fans_included', 'case_max_fans', 'case_tempered_glass',
        'case_rgb_fans', 'case_dust_filters',

        -- Охлаждение (6 атрибутов)
        'cooler_type', 'cooler_radiator', 'cooler_fan_size', 'cooler_max_tdp',
        'cooler_rgb', 'cooler_noise'
    )
);

-- Удалить все вариативные атрибуты (если существуют)
DELETE FROM category_variant_attributes
WHERE attribute_id IN (
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

-- Удалить все значения атрибутов в листингах
DELETE FROM listing_attribute_values
WHERE attribute_id IN (
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

-- Удалить сами атрибуты
DELETE FROM attributes WHERE code IN (
    -- Общие атрибуты
    'pc_brand', 'pc_model', 'pc_condition', 'pc_warranty', 'pc_color', 'pc_rgb_lighting',

    -- Видеокарты (10 атрибутов)
    'gpu_series', 'gpu_chip', 'gpu_vram', 'gpu_memory_type', 'gpu_memory_bus',
    'gpu_cooling', 'gpu_power', 'gpu_tdp', 'gpu_length', 'gpu_ray_tracing',

    -- Процессоры (9 атрибутов)
    'cpu_series', 'cpu_socket', 'cpu_cores', 'cpu_threads', 'cpu_base_clock',
    'cpu_boost_clock', 'cpu_tdp', 'cpu_igpu', 'cpu_generation',

    -- RAM (8 атрибутов)
    'ram_type', 'ram_capacity', 'ram_kit', 'ram_speed', 'ram_cas_latency',
    'ram_voltage', 'ram_ecc', 'ram_heatspreader',

    -- SSD (8 атрибутов)
    'ssd_capacity', 'ssd_interface', 'ssd_form_factor', 'ssd_read_speed',
    'ssd_write_speed', 'ssd_nand_type', 'ssd_dram_cache', 'ssd_endurance',

    -- HDD (6 атрибутов)
    'hdd_capacity', 'hdd_rpm', 'hdd_cache', 'hdd_interface',
    'hdd_form_factor', 'hdd_usage',

    -- Материнские платы (11 атрибутов)
    'mb_socket', 'mb_chipset', 'mb_form_factor', 'mb_memory_type',
    'mb_memory_slots', 'mb_max_memory', 'mb_m2_slots', 'mb_sata_ports',
    'mb_wifi', 'mb_bluetooth',

    -- PSU (6 атрибутов)
    'psu_wattage', 'psu_efficiency', 'psu_modular', 'psu_form_factor',
    'psu_pcie5', 'psu_fan_size',

    -- Корпуса (8 атрибутов)
    'case_form_factor', 'case_max_gpu_length', 'case_max_cpu_height',
    'case_fans_included', 'case_max_fans', 'case_tempered_glass',
    'case_rgb_fans', 'case_dust_filters',

    -- Охлаждение (6 атрибутов)
    'cooler_type', 'cooler_radiator', 'cooler_fan_size', 'cooler_max_tdp',
    'cooler_rgb', 'cooler_noise'
);

COMMIT;

-- Итоговое сообщение
DO $$
BEGIN
    RAISE NOTICE '✅ Rollback completed successfully!';
    RAISE NOTICE '';
    RAISE NOTICE 'All computer components attributes have been removed from the database.';
END $$;
