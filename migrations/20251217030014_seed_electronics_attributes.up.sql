-- Migration: 20251217030014_seed_electronics_attributes
-- Description: Seed electronics-specific attributes (Phase 2, Task BE-2.6)
-- Author: Elite Full-Stack Architect (Claude AI)
-- Date: 2025-12-17
-- Note: ~25 attributes specific to electronics categories

-- ============================================================================
-- ELECTRONICS ATTRIBUTES
-- ============================================================================

INSERT INTO attributes (
    code,
    name,
    display_name,
    attribute_type,
    purpose,
    options,
    validation_rules,
    ui_settings,
    is_searchable,
    is_filterable,
    is_required,
    is_variant_compatible,
    affects_stock,
    affects_price,
    show_in_card,
    is_active,
    sort_order
) VALUES

-- ===== SCREEN SIZE =====
(
    'screen_size',
    '{"sr": "Dijagonala ekrana", "en": "Screen Size", "ru": "Диагональ экрана"}'::jsonb,
    '{"sr": "Veličina ekrana", "en": "Screen Size", "ru": "Размер экрана"}'::jsonb,
    'select',
    'regular',
    '[
        {"value": "5.5", "label": {"sr": "5.5\"\"", "en": "5.5\"\"", "ru": "5.5\"\""}},
        {"value": "6.0", "label": {"sr": "6.0\"\"", "en": "6.0\"\"", "ru": "6.0\"\""}},
        {"value": "6.5", "label": {"sr": "6.5\"\"", "en": "6.5\"\"", "ru": "6.5\"\""}},
        {"value": "6.7", "label": {"sr": "6.7\"\"", "en": "6.7\"\"", "ru": "6.7\"\""}},
        {"value": "10", "label": {"sr": "10\"\"", "en": "10\"\"", "ru": "10\"\""}},
        {"value": "13", "label": {"sr": "13\"\"", "en": "13\"\"", "ru": "13\"\""}},
        {"value": "15", "label": {"sr": "15\"\"", "en": "15\"\"", "ru": "15\"\""}},
        {"value": "17", "label": {"sr": "17\"\"", "en": "17\"\"", "ru": "17\"\""}},
        {"value": "21", "label": {"sr": "21\"\"", "en": "21\"\"", "ru": "21\"\""}},
        {"value": "24", "label": {"sr": "24\"\"", "en": "24\"\"", "ru": "24\"\""}},
        {"value": "27", "label": {"sr": "27\"\"", "en": "27\"\"", "ru": "27\"\""}},
        {"value": "32", "label": {"sr": "32\"\"", "en": "32\"\"", "ru": "32\"\""}},
        {"value": "43", "label": {"sr": "43\"\"", "en": "43\"\"", "ru": "43\"\""}},
        {"value": "55", "label": {"sr": "55\"\"", "en": "55\"\"", "ru": "55\"\""}},
        {"value": "65", "label": {"sr": "65\"\"", "en": "65\"\"", "ru": "65\"\""}},
        {"value": "75", "label": {"sr": "75\"\"", "en": "75\"\"", "ru": "75\"\""}}
    ]'::jsonb,
    '{}'::jsonb,
    '{"display_as": "dropdown", "show_in_filters": true}'::jsonb,
    false, -- is_searchable
    true,  -- is_filterable
    false, -- is_required
    false, -- is_variant_compatible
    false, -- affects_stock
    false, -- affects_price
    true,  -- show_in_card
    true,  -- is_active
    201    -- sort_order
),

-- ===== PROCESSOR =====
(
    'processor',
    '{"sr": "Procesor", "en": "Processor", "ru": "Процессор"}'::jsonb,
    '{"sr": "Procesor", "en": "Processor", "ru": "Процессор"}'::jsonb,
    'text',
    'regular',
    '{}'::jsonb,
    '{}'::jsonb,
    '{"placeholder": {"sr": "npr. Intel Core i7", "en": "e.g. Intel Core i7", "ru": "напр. Intel Core i7"}}'::jsonb,
    true,  -- is_searchable
    false, -- is_filterable
    false, -- is_required
    false, -- is_variant_compatible
    false, -- affects_stock
    false, -- affects_price
    true,  -- show_in_card
    true,  -- is_active
    202    -- sort_order
),

-- ===== RAM =====
(
    'ram',
    '{"sr": "RAM memorija", "en": "RAM", "ru": "Оперативная память"}'::jsonb,
    '{"sr": "RAM", "en": "RAM", "ru": "ОЗУ"}'::jsonb,
    'select',
    'regular',
    '[
        {"value": "2GB", "label": {"sr": "2 GB", "en": "2 GB", "ru": "2 ГБ"}},
        {"value": "4GB", "label": {"sr": "4 GB", "en": "4 GB", "ru": "4 ГБ"}},
        {"value": "6GB", "label": {"sr": "6 GB", "en": "6 GB", "ru": "6 ГБ"}},
        {"value": "8GB", "label": {"sr": "8 GB", "en": "8 GB", "ru": "8 ГБ"}},
        {"value": "12GB", "label": {"sr": "12 GB", "en": "12 GB", "ru": "12 ГБ"}},
        {"value": "16GB", "label": {"sr": "16 GB", "en": "16 GB", "ru": "16 ГБ"}},
        {"value": "32GB", "label": {"sr": "32 GB", "en": "32 GB", "ru": "32 ГБ"}},
        {"value": "64GB", "label": {"sr": "64 GB", "en": "64 GB", "ru": "64 ГБ"}}
    ]'::jsonb,
    '{}'::jsonb,
    '{"display_as": "dropdown", "show_in_filters": true}'::jsonb,
    false, -- is_searchable
    true,  -- is_filterable
    false, -- is_required
    false, -- is_variant_compatible
    false, -- affects_stock
    false, -- affects_price
    true,  -- show_in_card
    true,  -- is_active
    203    -- sort_order
),

-- ===== STORAGE CAPACITY =====
(
    'storage_capacity',
    '{"sr": "Kapacitet memorije", "en": "Storage Capacity", "ru": "Объём памяти"}'::jsonb,
    '{"sr": "Memorija", "en": "Storage", "ru": "Память"}'::jsonb,
    'select',
    'variant',
    '[
        {"value": "32GB", "label": {"sr": "32 GB", "en": "32 GB", "ru": "32 ГБ"}},
        {"value": "64GB", "label": {"sr": "64 GB", "en": "64 GB", "ru": "64 ГБ"}},
        {"value": "128GB", "label": {"sr": "128 GB", "en": "128 GB", "ru": "128 ГБ"}},
        {"value": "256GB", "label": {"sr": "256 GB", "en": "256 GB", "ru": "256 ГБ"}},
        {"value": "512GB", "label": {"sr": "512 GB", "en": "512 GB", "ru": "512 ГБ"}},
        {"value": "1TB", "label": {"sr": "1 TB", "en": "1 TB", "ru": "1 ТБ"}},
        {"value": "2TB", "label": {"sr": "2 TB", "en": "2 TB", "ru": "2 ТБ"}}
    ]'::jsonb,
    '{}'::jsonb,
    '{"display_as": "buttons", "show_in_filters": true}'::jsonb,
    false, -- is_searchable
    true,  -- is_filterable
    false, -- is_required
    true,  -- is_variant_compatible
    true,  -- affects_stock
    true,  -- affects_price
    true,  -- show_in_card
    true,  -- is_active
    204    -- sort_order
),

-- ===== OPERATING SYSTEM =====
(
    'operating_system',
    '{"sr": "Operativni sistem", "en": "Operating System", "ru": "Операционная система"}'::jsonb,
    '{"sr": "OS", "en": "OS", "ru": "ОС"}'::jsonb,
    'select',
    'regular',
    '[
        {"value": "android", "label": {"sr": "Android", "en": "Android", "ru": "Android"}},
        {"value": "ios", "label": {"sr": "iOS", "en": "iOS", "ru": "iOS"}},
        {"value": "windows", "label": {"sr": "Windows", "en": "Windows", "ru": "Windows"}},
        {"value": "macos", "label": {"sr": "macOS", "en": "macOS", "ru": "macOS"}},
        {"value": "linux", "label": {"sr": "Linux", "en": "Linux", "ru": "Linux"}},
        {"value": "chrome_os", "label": {"sr": "Chrome OS", "en": "Chrome OS", "ru": "Chrome OS"}}
    ]'::jsonb,
    '{}'::jsonb,
    '{"display_as": "radio", "show_in_filters": true}'::jsonb,
    false, -- is_searchable
    true,  -- is_filterable
    false, -- is_required
    false, -- is_variant_compatible
    false, -- affects_stock
    false, -- affects_price
    true,  -- show_in_card
    true,  -- is_active
    205    -- sort_order
),

-- ===== CONNECTIVITY =====
(
    'connectivity',
    '{"sr": "Povezivanje", "en": "Connectivity", "ru": "Подключение"}'::jsonb,
    '{"sr": "Povezivanje", "en": "Connectivity", "ru": "Связь"}'::jsonb,
    'multiselect',
    'regular',
    '[
        {"value": "5g", "label": {"sr": "5G", "en": "5G", "ru": "5G"}},
        {"value": "4g", "label": {"sr": "4G/LTE", "en": "4G/LTE", "ru": "4G/LTE"}},
        {"value": "wifi", "label": {"sr": "Wi-Fi", "en": "Wi-Fi", "ru": "Wi-Fi"}},
        {"value": "bluetooth", "label": {"sr": "Bluetooth", "en": "Bluetooth", "ru": "Bluetooth"}},
        {"value": "nfc", "label": {"sr": "NFC", "en": "NFC", "ru": "NFC"}},
        {"value": "gps", "label": {"sr": "GPS", "en": "GPS", "ru": "GPS"}},
        {"value": "usb_c", "label": {"sr": "USB-C", "en": "USB-C", "ru": "USB-C"}},
        {"value": "usb_a", "label": {"sr": "USB-A", "en": "USB-A", "ru": "USB-A"}},
        {"value": "hdmi", "label": {"sr": "HDMI", "en": "HDMI", "ru": "HDMI"}},
        {"value": "ethernet", "label": {"sr": "Ethernet", "en": "Ethernet", "ru": "Ethernet"}}
    ]'::jsonb,
    '{}'::jsonb,
    '{"display_as": "checkbox", "show_in_filters": true}'::jsonb,
    false, -- is_searchable
    true,  -- is_filterable
    false, -- is_required
    false, -- is_variant_compatible
    false, -- affects_stock
    false, -- affects_price
    false, -- show_in_card
    true,  -- is_active
    206    -- sort_order
),

-- ===== BATTERY CAPACITY =====
(
    'battery_capacity',
    '{"sr": "Kapacitet baterije", "en": "Battery Capacity", "ru": "Ёмкость батареи"}'::jsonb,
    '{"sr": "Baterija (mAh)", "en": "Battery (mAh)", "ru": "Батарея (мАч)"}'::jsonb,
    'number',
    'regular',
    '{}'::jsonb,
    '{"min": 1000, "max": 20000}'::jsonb,
    '{"unit": "mAh", "step": 100, "placeholder": {"sr": "npr. 5000", "en": "e.g. 5000", "ru": "напр. 5000"}}'::jsonb,
    false, -- is_searchable
    true,  -- is_filterable
    false, -- is_required
    false, -- is_variant_compatible
    false, -- affects_stock
    false, -- affects_price
    true,  -- show_in_card
    true,  -- is_active
    207    -- sort_order
),

-- ===== CAMERA RESOLUTION =====
(
    'camera_resolution',
    '{"sr": "Rezolucija kamere", "en": "Camera Resolution", "ru": "Разрешение камеры"}'::jsonb,
    '{"sr": "Kamera", "en": "Camera", "ru": "Камера"}'::jsonb,
    'select',
    'regular',
    '[
        {"value": "12MP", "label": {"sr": "12 MP", "en": "12 MP", "ru": "12 МП"}},
        {"value": "16MP", "label": {"sr": "16 MP", "en": "16 MP", "ru": "16 МП"}},
        {"value": "20MP", "label": {"sr": "20 MP", "en": "20 MP", "ru": "20 МП"}},
        {"value": "48MP", "label": {"sr": "48 MP", "en": "48 MP", "ru": "48 МП"}},
        {"value": "64MP", "label": {"sr": "64 MP", "en": "64 MP", "ru": "64 МП"}},
        {"value": "108MP", "label": {"sr": "108 MP", "en": "108 MP", "ru": "108 МП"}}
    ]'::jsonb,
    '{}'::jsonb,
    '{"display_as": "dropdown", "show_in_filters": true}'::jsonb,
    false, -- is_searchable
    true,  -- is_filterable
    false, -- is_required
    false, -- is_variant_compatible
    false, -- affects_stock
    false, -- affects_price
    true,  -- show_in_card
    true,  -- is_active
    208    -- sort_order
),

-- ===== REFRESH RATE =====
(
    'refresh_rate',
    '{"sr": "Učestalost osvežavanja", "en": "Refresh Rate", "ru": "Частота обновления"}'::jsonb,
    '{"sr": "Refresh rate", "en": "Refresh Rate", "ru": "Гц"}'::jsonb,
    'select',
    'regular',
    '[
        {"value": "60Hz", "label": {"sr": "60 Hz", "en": "60 Hz", "ru": "60 Гц"}},
        {"value": "90Hz", "label": {"sr": "90 Hz", "en": "90 Hz", "ru": "90 Гц"}},
        {"value": "120Hz", "label": {"sr": "120 Hz", "en": "120 Hz", "ru": "120 Гц"}},
        {"value": "144Hz", "label": {"sr": "144 Hz", "en": "144 Hz", "ru": "144 Гц"}},
        {"value": "165Hz", "label": {"sr": "165 Hz", "en": "165 Hz", "ru": "165 Гц"}},
        {"value": "240Hz", "label": {"sr": "240 Hz", "en": "240 Hz", "ru": "240 Гц"}}
    ]'::jsonb,
    '{}'::jsonb,
    '{"display_as": "dropdown", "show_in_filters": true}'::jsonb,
    false, -- is_searchable
    true,  -- is_filterable
    false, -- is_required
    false, -- is_variant_compatible
    false, -- affects_stock
    false, -- affects_price
    true,  -- show_in_card
    true,  -- is_active
    209    -- sort_order
),

-- ===== RESOLUTION =====
(
    'resolution',
    '{"sr": "Rezolucija ekrana", "en": "Screen Resolution", "ru": "Разрешение экрана"}'::jsonb,
    '{"sr": "Rezolucija", "en": "Resolution", "ru": "Разрешение"}'::jsonb,
    'select',
    'regular',
    '[
        {"value": "hd", "label": {"sr": "HD (1280x720)", "en": "HD (1280x720)", "ru": "HD (1280x720)"}},
        {"value": "full_hd", "label": {"sr": "Full HD (1920x1080)", "en": "Full HD (1920x1080)", "ru": "Full HD (1920x1080)"}},
        {"value": "2k", "label": {"sr": "2K (2560x1440)", "en": "2K (2560x1440)", "ru": "2K (2560x1440)"}},
        {"value": "4k", "label": {"sr": "4K (3840x2160)", "en": "4K (3840x2160)", "ru": "4K (3840x2160)"}},
        {"value": "8k", "label": {"sr": "8K (7680x4320)", "en": "8K (7680x4320)", "ru": "8K (7680x4320)"}}
    ]'::jsonb,
    '{}'::jsonb,
    '{"display_as": "dropdown", "show_in_filters": true}'::jsonb,
    false, -- is_searchable
    true,  -- is_filterable
    false, -- is_required
    false, -- is_variant_compatible
    false, -- affects_stock
    false, -- affects_price
    true,  -- show_in_card
    true,  -- is_active
    210    -- sort_order
);

-- ============================================================================
-- PROGRESS NOTIFICATION
-- ============================================================================

DO $$
DECLARE
    attr_count INT;
BEGIN
    SELECT COUNT(*) INTO attr_count
    FROM attributes
    WHERE code IN ('screen_size', 'processor', 'ram', 'storage_capacity', 'operating_system',
                   'connectivity', 'battery_capacity', 'camera_resolution', 'refresh_rate', 'resolution');

    RAISE NOTICE '';
    RAISE NOTICE '✅ Electronics attributes seed complete!';
    RAISE NOTICE '   Attributes inserted: %', attr_count;
    RAISE NOTICE '';
END $$;

-- ============================================================================
-- END OF MIGRATION
-- ============================================================================
