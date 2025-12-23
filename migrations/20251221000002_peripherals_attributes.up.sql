-- ============================================================================
-- PERIPHERALS ATTRIBUTES
-- Add proper attributes for Electronics > Peripherals category
-- ============================================================================

-- Category ID: 2c1f8391-95e3-4a37-baec-f923b4c9e5a1 (periferija / Peripherals)

-- Step 1: Remove incorrect attributes (smartphone attributes) from Peripherals
DELETE FROM category_attributes
WHERE category_id = '2c1f8391-95e3-4a37-baec-f923b4c9e5a1'
  AND attribute_id IN (21, 22, 23, 24, 25, 27, 28);
-- Keeping: 26 (connectivity) - it's universal and useful for peripherals

-- Step 2: Create peripheral-specific attributes

-- 101: Connection Type (multiselect)
INSERT INTO attributes (id, code, attribute_type, purpose, name, display_name, options, is_required, is_filterable, is_searchable, is_active, sort_order)
VALUES (
    101,
    'connection_type',
    'multiselect',
    'regular',
    '{"en": "Connection Type", "ru": "Тип подключения", "sr": "Tip povezivanja"}',
    '{"en": "Connection Type", "ru": "Тип подключения", "sr": "Tip povezivanja"}',
    '[
        {"value": "usb_a", "label": {"en": "USB-A", "ru": "USB-A", "sr": "USB-A"}},
        {"value": "usb_c", "label": {"en": "USB-C", "ru": "USB-C", "sr": "USB-C"}},
        {"value": "bluetooth", "label": {"en": "Bluetooth", "ru": "Bluetooth", "sr": "Bluetooth"}},
        {"value": "wireless_2_4ghz", "label": {"en": "2.4GHz Wireless", "ru": "Беспроводной 2.4GHz", "sr": "Bežični 2.4GHz"}},
        {"value": "ps2", "label": {"en": "PS/2", "ru": "PS/2", "sr": "PS/2"}},
        {"value": "aux_3_5mm", "label": {"en": "3.5mm Audio Jack", "ru": "Аудио разъём 3.5мм", "sr": "3.5mm audio priključak"}},
        {"value": "optical", "label": {"en": "Optical (S/PDIF)", "ru": "Оптический (S/PDIF)", "sr": "Optički (S/PDIF)"}}
    ]'::jsonb,
    false,
    true,
    false,
    true,
    1
) ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    display_name = EXCLUDED.display_name,
    options = EXCLUDED.options;

-- 102: Keyboard Layout (select)
INSERT INTO attributes (id, code, attribute_type, purpose, name, display_name, options, is_required, is_filterable, is_searchable, is_active, sort_order)
VALUES (
    102,
    'keyboard_layout',
    'select',
    'regular',
    '{"en": "Keyboard Layout", "ru": "Раскладка клавиатуры", "sr": "Raspored tastature"}',
    '{"en": "Layout", "ru": "Раскладка", "sr": "Raspored"}',
    '[
        {"value": "qwerty_us", "label": {"en": "QWERTY (US)", "ru": "QWERTY (США)", "sr": "QWERTY (US)"}},
        {"value": "qwerty_uk", "label": {"en": "QWERTY (UK)", "ru": "QWERTY (Британия)", "sr": "QWERTY (UK)"}},
        {"value": "qwertz_de", "label": {"en": "QWERTZ (German)", "ru": "QWERTZ (Немецкая)", "sr": "QWERTZ (Nemački)"}},
        {"value": "azerty_fr", "label": {"en": "AZERTY (French)", "ru": "AZERTY (Французская)", "sr": "AZERTY (Francuski)"}},
        {"value": "serbian_latin", "label": {"en": "Serbian Latin", "ru": "Сербская латиница", "sr": "Srpska latinica"}},
        {"value": "serbian_cyrillic", "label": {"en": "Serbian Cyrillic", "ru": "Сербская кириллица", "sr": "Srpska ćirilica"}},
        {"value": "russian", "label": {"en": "Russian", "ru": "Русская", "sr": "Ruska"}}
    ]'::jsonb,
    false,
    true,
    false,
    true,
    2
) ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    display_name = EXCLUDED.display_name,
    options = EXCLUDED.options;

-- 103: Key Switch Type (select)
INSERT INTO attributes (id, code, attribute_type, purpose, name, display_name, options, is_required, is_filterable, is_searchable, is_active, sort_order)
VALUES (
    103,
    'key_switch_type',
    'select',
    'regular',
    '{"en": "Key Switch Type", "ru": "Тип переключателей", "sr": "Tip prekidača"}',
    '{"en": "Switches", "ru": "Переключатели", "sr": "Prekidači"}',
    '[
        {"value": "mechanical_cherry", "label": {"en": "Mechanical (Cherry MX)", "ru": "Механические (Cherry MX)", "sr": "Mehanički (Cherry MX)"}},
        {"value": "mechanical_gateron", "label": {"en": "Mechanical (Gateron)", "ru": "Механические (Gateron)", "sr": "Mehanički (Gateron)"}},
        {"value": "mechanical_other", "label": {"en": "Mechanical (Other)", "ru": "Механические (Другие)", "sr": "Mehanički (Ostalo)"}},
        {"value": "membrane", "label": {"en": "Membrane", "ru": "Мембранные", "sr": "Membranski"}},
        {"value": "scissor", "label": {"en": "Scissor", "ru": "Ножничные", "sr": "Makaze"}},
        {"value": "optical", "label": {"en": "Optical", "ru": "Оптические", "sr": "Optički"}},
        {"value": "hybrid", "label": {"en": "Hybrid", "ru": "Гибридные", "sr": "Hibridni"}}
    ]'::jsonb,
    false,
    true,
    false,
    true,
    3
) ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    display_name = EXCLUDED.display_name,
    options = EXCLUDED.options;

-- 104: Mouse DPI (number)
INSERT INTO attributes (id, code, attribute_type, purpose, name, display_name, validation_rules, is_required, is_filterable, is_searchable, is_active, sort_order)
VALUES (
    104,
    'mouse_dpi',
    'number',
    'regular',
    '{"en": "Mouse DPI", "ru": "DPI мыши", "sr": "DPI miša"}',
    '{"en": "DPI", "ru": "DPI", "sr": "DPI"}',
    '{"min": 100, "max": 32000, "step": 100}'::jsonb,
    false,
    true,
    false,
    true,
    4
) ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    display_name = EXCLUDED.display_name,
    validation_rules = EXCLUDED.validation_rules;

-- 105: Headphone Type (select)
INSERT INTO attributes (id, code, attribute_type, purpose, name, display_name, options, is_required, is_filterable, is_searchable, is_active, sort_order)
VALUES (
    105,
    'headphone_type',
    'select',
    'regular',
    '{"en": "Headphone Type", "ru": "Тип наушников", "sr": "Tip slušalica"}',
    '{"en": "Type", "ru": "Тип", "sr": "Tip"}',
    '[
        {"value": "over_ear", "label": {"en": "Over-ear (Full-size)", "ru": "Полноразмерные", "sr": "Preko ušiju (Full-size)"}},
        {"value": "on_ear", "label": {"en": "On-ear", "ru": "Накладные", "sr": "Na ušima"}},
        {"value": "in_ear", "label": {"en": "In-ear", "ru": "Внутриканальные", "sr": "U ušima"}},
        {"value": "true_wireless", "label": {"en": "True Wireless (TWS)", "ru": "Полностью беспроводные (TWS)", "sr": "Potpuno bežične (TWS)"}},
        {"value": "earbuds", "label": {"en": "Earbuds", "ru": "Вкладыши", "sr": "Slušalice"}},
        {"value": "bone_conduction", "label": {"en": "Bone Conduction", "ru": "Костной проводимости", "sr": "Koštana provodljivost"}}
    ]'::jsonb,
    false,
    true,
    false,
    true,
    5
) ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    display_name = EXCLUDED.display_name,
    options = EXCLUDED.options;

-- 106: Has Microphone (boolean)
INSERT INTO attributes (id, code, attribute_type, purpose, name, display_name, is_required, is_filterable, is_searchable, is_active, sort_order)
VALUES (
    106,
    'has_microphone',
    'boolean',
    'regular',
    '{"en": "Has Microphone", "ru": "Есть микрофон", "sr": "Ima mikrofon"}',
    '{"en": "Microphone", "ru": "Микрофон", "sr": "Mikrofon"}',
    false,
    true,
    false,
    true,
    6
) ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    display_name = EXCLUDED.display_name;

-- 107: RGB Lighting (boolean)
INSERT INTO attributes (id, code, attribute_type, purpose, name, display_name, is_required, is_filterable, is_searchable, is_active, sort_order)
VALUES (
    107,
    'has_rgb',
    'boolean',
    'regular',
    '{"en": "RGB Lighting", "ru": "RGB подсветка", "sr": "RGB osvetljenje"}',
    '{"en": "RGB", "ru": "RGB", "sr": "RGB"}',
    false,
    true,
    false,
    true,
    7
) ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    display_name = EXCLUDED.display_name;

-- 108: Wireless (boolean)
INSERT INTO attributes (id, code, attribute_type, purpose, name, display_name, is_required, is_filterable, is_searchable, is_active, sort_order)
VALUES (
    108,
    'is_wireless',
    'boolean',
    'regular',
    '{"en": "Wireless", "ru": "Беспроводное", "sr": "Bežično"}',
    '{"en": "Wireless", "ru": "Беспроводное", "sr": "Bežično"}',
    false,
    true,
    false,
    true,
    8
) ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    display_name = EXCLUDED.display_name;

-- 109: Battery Life (number - hours)
INSERT INTO attributes (id, code, attribute_type, purpose, name, display_name, validation_rules, is_required, is_filterable, is_searchable, is_active, sort_order)
VALUES (
    109,
    'battery_life_hours',
    'number',
    'regular',
    '{"en": "Battery Life (hours)", "ru": "Время работы от батареи (часы)", "sr": "Trajanje baterije (sati)"}',
    '{"en": "Battery Life", "ru": "Батарея", "sr": "Baterija"}',
    '{"min": 1, "max": 1000, "step": 1, "unit": "h"}'::jsonb,
    false,
    true,
    false,
    true,
    9
) ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    display_name = EXCLUDED.display_name,
    validation_rules = EXCLUDED.validation_rules;

-- 110: Polling Rate (select)
INSERT INTO attributes (id, code, attribute_type, purpose, name, display_name, options, is_required, is_filterable, is_searchable, is_active, sort_order)
VALUES (
    110,
    'polling_rate',
    'select',
    'regular',
    '{"en": "Polling Rate", "ru": "Частота опроса", "sr": "Brzina osvežavanja"}',
    '{"en": "Polling Rate", "ru": "Частота опроса", "sr": "Polling Rate"}',
    '[
        {"value": "125hz", "label": {"en": "125 Hz", "ru": "125 Гц", "sr": "125 Hz"}},
        {"value": "250hz", "label": {"en": "250 Hz", "ru": "250 Гц", "sr": "250 Hz"}},
        {"value": "500hz", "label": {"en": "500 Hz", "ru": "500 Гц", "sr": "500 Hz"}},
        {"value": "1000hz", "label": {"en": "1000 Hz (1ms)", "ru": "1000 Гц (1мс)", "sr": "1000 Hz (1ms)"}},
        {"value": "2000hz", "label": {"en": "2000 Hz", "ru": "2000 Гц", "sr": "2000 Hz"}},
        {"value": "4000hz", "label": {"en": "4000 Hz", "ru": "4000 Гц", "sr": "4000 Hz"}},
        {"value": "8000hz", "label": {"en": "8000 Hz", "ru": "8000 Гц", "sr": "8000 Hz"}}
    ]'::jsonb,
    false,
    true,
    false,
    true,
    10
) ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    display_name = EXCLUDED.display_name,
    options = EXCLUDED.options;

-- 111: Driver Size (for headphones, number in mm)
INSERT INTO attributes (id, code, attribute_type, purpose, name, display_name, validation_rules, is_required, is_filterable, is_searchable, is_active, sort_order)
VALUES (
    111,
    'driver_size_mm',
    'number',
    'regular',
    '{"en": "Driver Size (mm)", "ru": "Размер драйвера (мм)", "sr": "Veličina drajvera (mm)"}',
    '{"en": "Driver Size", "ru": "Драйвер", "sr": "Drajver"}',
    '{"min": 5, "max": 100, "step": 1, "unit": "mm"}'::jsonb,
    false,
    true,
    false,
    true,
    11
) ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    display_name = EXCLUDED.display_name,
    validation_rules = EXCLUDED.validation_rules;

-- 112: Webcam Resolution (select)
INSERT INTO attributes (id, code, attribute_type, purpose, name, display_name, options, is_required, is_filterable, is_searchable, is_active, sort_order)
VALUES (
    112,
    'webcam_resolution',
    'select',
    'regular',
    '{"en": "Webcam Resolution", "ru": "Разрешение веб-камеры", "sr": "Rezolucija veb kamere"}',
    '{"en": "Resolution", "ru": "Разрешение", "sr": "Rezolucija"}',
    '[
        {"value": "720p", "label": {"en": "HD (720p)", "ru": "HD (720p)", "sr": "HD (720p)"}},
        {"value": "1080p", "label": {"en": "Full HD (1080p)", "ru": "Full HD (1080p)", "sr": "Full HD (1080p)"}},
        {"value": "2k", "label": {"en": "2K (1440p)", "ru": "2K (1440p)", "sr": "2K (1440p)"}},
        {"value": "4k", "label": {"en": "4K (2160p)", "ru": "4K (2160p)", "sr": "4K (2160p)"}}
    ]'::jsonb,
    false,
    true,
    false,
    true,
    12
) ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    display_name = EXCLUDED.display_name,
    options = EXCLUDED.options;

-- Step 3: Link new attributes to Peripherals category
-- Category: 2c1f8391-95e3-4a37-baec-f923b4c9e5a1

INSERT INTO category_attributes (category_id, attribute_id, is_enabled, is_required, is_filterable, is_searchable, sort_order)
VALUES
    -- Connection Type
    ('2c1f8391-95e3-4a37-baec-f923b4c9e5a1', 101, true, false, true, false, 1),
    -- Keyboard Layout
    ('2c1f8391-95e3-4a37-baec-f923b4c9e5a1', 102, true, false, true, false, 2),
    -- Key Switch Type
    ('2c1f8391-95e3-4a37-baec-f923b4c9e5a1', 103, true, false, true, false, 3),
    -- Mouse DPI
    ('2c1f8391-95e3-4a37-baec-f923b4c9e5a1', 104, true, false, true, false, 4),
    -- Headphone Type
    ('2c1f8391-95e3-4a37-baec-f923b4c9e5a1', 105, true, false, true, false, 5),
    -- Has Microphone
    ('2c1f8391-95e3-4a37-baec-f923b4c9e5a1', 106, true, false, true, false, 6),
    -- RGB Lighting
    ('2c1f8391-95e3-4a37-baec-f923b4c9e5a1', 107, true, false, true, false, 7),
    -- Wireless
    ('2c1f8391-95e3-4a37-baec-f923b4c9e5a1', 108, true, false, true, false, 8),
    -- Battery Life
    ('2c1f8391-95e3-4a37-baec-f923b4c9e5a1', 109, true, false, true, false, 9),
    -- Polling Rate
    ('2c1f8391-95e3-4a37-baec-f923b4c9e5a1', 110, true, false, true, false, 10),
    -- Driver Size
    ('2c1f8391-95e3-4a37-baec-f923b4c9e5a1', 111, true, false, true, false, 11),
    -- Webcam Resolution
    ('2c1f8391-95e3-4a37-baec-f923b4c9e5a1', 112, true, false, true, false, 12)
ON CONFLICT (category_id, attribute_id) DO UPDATE SET
    is_enabled = EXCLUDED.is_enabled,
    is_filterable = EXCLUDED.is_filterable,
    sort_order = EXCLUDED.sort_order;

-- Also keep brand (1), condition (2), warranty (7) as they are universal
INSERT INTO category_attributes (category_id, attribute_id, is_enabled, is_required, is_filterable, is_searchable, sort_order)
VALUES
    ('2c1f8391-95e3-4a37-baec-f923b4c9e5a1', 1, true, false, true, true, 0),   -- Brand
    ('2c1f8391-95e3-4a37-baec-f923b4c9e5a1', 2, true, false, true, false, 0),  -- Condition
    ('2c1f8391-95e3-4a37-baec-f923b4c9e5a1', 7, true, false, true, false, 0)   -- Warranty
ON CONFLICT (category_id, attribute_id) DO UPDATE SET
    is_enabled = EXCLUDED.is_enabled,
    is_filterable = EXCLUDED.is_filterable;
