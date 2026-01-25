-- ==========================================
-- Migration: Create Computer Components Attributes
-- Категория: Računarske komponente
-- Дата: 2026-01-26
-- Описание: Создание 84+ атрибутов для 9 подкатегорий компьютерных компонентов
-- ==========================================

BEGIN;

-- ==========================================
-- 1. ОБЩИЕ АТРИБУТЫ (для всех подкатегорий)
-- ==========================================

-- Brand (Производитель)
INSERT INTO attributes (code, name, display_name, attribute_type, purpose, options, is_filterable, is_searchable, is_required, is_active, sort_order)
VALUES (
    'pc_brand',
    '{"en": "Brand", "ru": "Производитель", "sr": "Brend"}',
    '{"en": "Brand", "ru": "Производитель", "sr": "Brend"}',
    'select',
    'regular',
    '{
        "allowed_values": [
            "AMD", "Intel", "NVIDIA",
            "ASUS", "MSI", "Gigabyte", "ASRock", "EVGA",
            "Corsair", "G.Skill", "Kingston", "Crucial",
            "Samsung", "Western Digital", "Seagate",
            "Noctua", "be quiet!", "Cooler Master",
            "Fractal Design", "NZXT", "Thermaltake", "Other"
        ]
    }',
    true,
    true,
    true,
    true,
    10
) ON CONFLICT (code) DO NOTHING;

-- Model (Модель)
INSERT INTO attributes (code, name, display_name, attribute_type, purpose, validation_rules, is_searchable, is_required, is_active, sort_order)
VALUES (
    'pc_model',
    '{"en": "Model", "ru": "Модель", "sr": "Model"}',
    '{"en": "Model", "ru": "Модель", "sr": "Model"}',
    'text',
    'regular',
    '{"max_length": 100}',
    true,
    true,
    true,
    20
) ON CONFLICT (code) DO NOTHING;

-- Condition (Состояние)
INSERT INTO attributes (code, name, display_name, attribute_type, purpose, options, is_filterable, is_required, is_active, sort_order)
VALUES (
    'pc_condition',
    '{"en": "Condition", "ru": "Состояние", "sr": "Stanje"}',
    '{"en": "Condition", "ru": "Состояние", "sr": "Stanje"}',
    'select',
    'regular',
    '{
        "allowed_values": [
            {"value": "new", "label": {"en": "New", "ru": "Новое", "sr": "Novo"}},
            {"value": "used_like_new", "label": {"en": "Used - Like New", "ru": "Б/У - как новое", "sr": "Polovno - kao novo"}},
            {"value": "used_good", "label": {"en": "Used - Good", "ru": "Б/У - хорошее", "sr": "Polovno - dobro"}},
            {"value": "used_fair", "label": {"en": "Used - Fair", "ru": "Б/У - удовлетворительное", "sr": "Polovno - zadovoljavajuće"}},
            {"value": "refurbished", "label": {"en": "Refurbished", "ru": "Восстановленное", "sr": "Obnovljeno"}}
        ]
    }',
    true,
    true,
    true,
    30
) ON CONFLICT (code) DO NOTHING;

-- Warranty (Гарантия)
INSERT INTO attributes (code, name, display_name, attribute_type, purpose, options, is_filterable, is_active, sort_order)
VALUES (
    'pc_warranty',
    '{"en": "Warranty", "ru": "Гарантия", "sr": "Garancija"}',
    '{"en": "Warranty", "ru": "Гарантия", "sr": "Garancija"}',
    'select',
    'regular',
    '{
        "allowed_values": [
            {"value": "no_warranty", "label": {"en": "No warranty", "ru": "Без гарантии", "sr": "Bez garancije"}},
            {"value": "6_months", "label": {"en": "6 months", "ru": "6 месяцев", "sr": "6 meseci"}},
            {"value": "1_year", "label": {"en": "1 year", "ru": "1 год", "sr": "1 godina"}},
            {"value": "2_years", "label": {"en": "2 years", "ru": "2 года", "sr": "2 godine"}},
            {"value": "3_years_plus", "label": {"en": "3+ years", "ru": "3+ года", "sr": "3+ godine"}}
        ]
    }',
    true,
    true,
    40
) ON CONFLICT (code) DO NOTHING;

-- Color (Цвет)
INSERT INTO attributes (code, name, display_name, attribute_type, purpose, options, is_filterable, is_active, sort_order)
VALUES (
    'pc_color',
    '{"en": "Color", "ru": "Цвет", "sr": "Boja"}',
    '{"en": "Color", "ru": "Цвет", "sr": "Boja"}',
    'select',
    'regular',
    '{
        "allowed_values": ["Black", "White", "Silver", "RGB", "Red", "Blue", "Other"]
    }',
    true,
    true,
    50
) ON CONFLICT (code) DO NOTHING;

-- RGB Lighting (RGB подсветка)
INSERT INTO attributes (code, name, display_name, attribute_type, purpose, is_filterable, is_active, sort_order)
VALUES (
    'pc_rgb_lighting',
    '{"en": "RGB Lighting", "ru": "RGB подсветка", "sr": "RGB osvetljenje"}',
    '{"en": "RGB Lighting", "ru": "RGB подсветка", "sr": "RGB osvetljenje"}',
    'boolean',
    'regular',
    true,
    true,
    60
) ON CONFLICT (code) DO NOTHING;

-- ==========================================
-- 2. АТРИБУТЫ ДЛЯ ВИДЕОКАРТ (Grafičke kartice)
-- ==========================================

-- GPU Series
INSERT INTO attributes (code, name, display_name, attribute_type, purpose, options, is_filterable, is_required, is_active, sort_order)
VALUES (
    'gpu_series',
    '{"en": "GPU Series", "ru": "Серия GPU", "sr": "GPU serija"}',
    '{"en": "GPU Series", "ru": "Серия GPU", "sr": "GPU serija"}',
    'select',
    'regular',
    '{
        "allowed_values": [
            "NVIDIA RTX 4000 Series",
            "NVIDIA RTX 3000 Series",
            "NVIDIA RTX 2000 Series",
            "NVIDIA GTX 1000 Series",
            "AMD RX 7000 Series",
            "AMD RX 6000 Series",
            "AMD RX 5000 Series",
            "Intel Arc"
        ]
    }',
    true,
    true,
    true,
    100
) ON CONFLICT (code) DO NOTHING;

-- GPU Chip
INSERT INTO attributes (code, name, display_name, attribute_type, purpose, options, is_filterable, is_searchable, is_active, sort_order)
VALUES (
    'gpu_chip',
    '{"en": "GPU Chip", "ru": "Чип GPU", "sr": "GPU čip"}',
    '{"en": "GPU Chip", "ru": "Чип GPU", "sr": "GPU čip"}',
    'select',
    'regular',
    '{
        "allowed_values": [
            "RTX 4090", "RTX 4080", "RTX 4070 Ti", "RTX 4070", "RTX 4060 Ti",
            "RTX 3090 Ti", "RTX 3090", "RTX 3080 Ti", "RTX 3080",
            "RX 7900 XTX", "RX 7900 XT", "RX 6900 XT", "RX 6800 XT", "Other"
        ]
    }',
    true,
    true,
    true,
    110
) ON CONFLICT (code) DO NOTHING;

-- GPU VRAM
INSERT INTO attributes (code, name, display_name, attribute_type, purpose, options, is_filterable, is_required, is_active, sort_order)
VALUES (
    'gpu_vram',
    '{"en": "VRAM", "ru": "Видеопамять", "sr": "Video memorija"}',
    '{"en": "VRAM", "ru": "Видеопамять", "sr": "Video memorija"}',
    'select',
    'variant',
    '{
        "allowed_values": ["6GB", "8GB", "12GB", "16GB", "24GB"]
    }',
    true,
    true,
    true,
    120
) ON CONFLICT (code) DO NOTHING;

-- GPU Memory Type
INSERT INTO attributes (code, name, display_name, attribute_type, purpose, options, is_filterable, is_active, sort_order)
VALUES (
    'gpu_memory_type',
    '{"en": "Memory Type", "ru": "Тип памяти", "sr": "Tip memorije"}',
    '{"en": "Memory Type", "ru": "Тип памяти", "sr": "Tip memorije"}',
    'select',
    'regular',
    '{
        "allowed_values": ["GDDR6", "GDDR6X"]
    }',
    true,
    true,
    130
) ON CONFLICT (code) DO NOTHING;

-- GPU Memory Bus
INSERT INTO attributes (code, name, display_name, attribute_type, purpose, options, is_filterable, is_active, sort_order)
VALUES (
    'gpu_memory_bus',
    '{"en": "Memory Bus", "ru": "Шина памяти", "sr": "Magistrala memorije"}',
    '{"en": "Memory Bus", "ru": "Шина памяти", "sr": "Magistrala memorije"}',
    'select',
    'regular',
    '{
        "allowed_values": ["128-bit", "192-bit", "256-bit", "384-bit"]
    }',
    true,
    true,
    140
) ON CONFLICT (code) DO NOTHING;

-- GPU Cooling Type
INSERT INTO attributes (code, name, display_name, attribute_type, purpose, options, is_filterable, is_active, sort_order)
VALUES (
    'gpu_cooling',
    '{"en": "Cooling Type", "ru": "Тип охлаждения", "sr": "Tip hlađenja"}',
    '{"en": "Cooling Type", "ru": "Тип охлаждения", "sr": "Tip hlađenja"}',
    'select',
    'regular',
    '{
        "allowed_values": ["Air", "Hybrid", "Water"]
    }',
    true,
    true,
    150
) ON CONFLICT (code) DO NOTHING;

-- GPU Power Connectors
INSERT INTO attributes (code, name, display_name, attribute_type, purpose, validation_rules, is_active, sort_order)
VALUES (
    'gpu_power',
    '{"en": "Power Connectors", "ru": "Разъёмы питания", "sr": "Konektori napajanja"}',
    '{"en": "Power Connectors", "ru": "Разъёмы питания", "sr": "Konektori napajanja"}',
    'text',
    'regular',
    '{"max_length": 100, "placeholder": "2x 8-pin"}',
    true,
    160
) ON CONFLICT (code) DO NOTHING;

-- GPU TDP
INSERT INTO attributes (code, name, display_name, attribute_type, purpose, options, is_filterable, is_active, sort_order)
VALUES (
    'gpu_tdp',
    '{"en": "TDP", "ru": "TDP", "sr": "TDP"}',
    '{"en": "TDP", "ru": "TDP", "sr": "TDP"}',
    'select',
    'regular',
    '{
        "allowed_values": ["150W", "200W", "250W", "300W", "350W", "450W"]
    }',
    true,
    true,
    170
) ON CONFLICT (code) DO NOTHING;

-- GPU Length
INSERT INTO attributes (code, name, display_name, attribute_type, purpose, validation_rules, is_filterable, is_active, sort_order)
VALUES (
    'gpu_length',
    '{"en": "Length (mm)", "ru": "Длина (мм)", "sr": "Dužina (mm)"}',
    '{"en": "Length (mm)", "ru": "Длина (мм)", "sr": "Dužina (mm)"}',
    'number',
    'regular',
    '{"min": 150, "max": 400, "unit": "mm"}',
    true,
    true,
    180
) ON CONFLICT (code) DO NOTHING;

-- GPU Ray Tracing
INSERT INTO attributes (code, name, display_name, attribute_type, purpose, is_filterable, is_active, sort_order)
VALUES (
    'gpu_ray_tracing',
    '{"en": "Ray Tracing", "ru": "Трассировка лучей", "sr": "Ray Tracing"}',
    '{"en": "Ray Tracing", "ru": "Трассировка лучей", "sr": "Ray Tracing"}',
    'boolean',
    'regular',
    true,
    true,
    190
) ON CONFLICT (code) DO NOTHING;

-- ==========================================
-- 3. АТРИБУТЫ ДЛЯ ПРОЦЕССОРОВ (Procesori)
-- ==========================================

-- Processor Series
INSERT INTO attributes (code, name, display_name, attribute_type, purpose, options, is_filterable, is_required, is_active, sort_order)
VALUES (
    'cpu_series',
    '{"en": "Processor Series", "ru": "Серия процессора", "sr": "Serija procesora"}',
    '{"en": "Processor Series", "ru": "Серия процессора", "sr": "Serija procesora"}',
    'select',
    'regular',
    '{
        "allowed_values": [
            "Intel Core i9", "Intel Core i7", "Intel Core i5", "Intel Core i3",
            "AMD Ryzen 9", "AMD Ryzen 7", "AMD Ryzen 5", "AMD Ryzen 3",
            "AMD Threadripper"
        ]
    }',
    true,
    true,
    true,
    200
) ON CONFLICT (code) DO NOTHING;

-- CPU Socket
INSERT INTO attributes (code, name, display_name, attribute_type, purpose, options, is_filterable, is_required, is_active, sort_order)
VALUES (
    'cpu_socket',
    '{"en": "Socket", "ru": "Сокет", "sr": "Soket"}',
    '{"en": "Socket", "ru": "Сокет", "sr": "Soket"}',
    'select',
    'regular',
    '{
        "allowed_values": ["LGA1700", "LGA1200", "AM5", "AM4"]
    }',
    true,
    true,
    true,
    210
) ON CONFLICT (code) DO NOTHING;

-- CPU Cores
INSERT INTO attributes (code, name, display_name, attribute_type, purpose, validation_rules, is_filterable, is_required, is_active, sort_order)
VALUES (
    'cpu_cores',
    '{"en": "Cores", "ru": "Ядра", "sr": "Jezgra"}',
    '{"en": "Cores", "ru": "Ядра", "sr": "Jezgra"}',
    'number',
    'variant',
    '{"min": 4, "max": 64}',
    true,
    true,
    true,
    220
) ON CONFLICT (code) DO NOTHING;

-- CPU Threads
INSERT INTO attributes (code, name, display_name, attribute_type, purpose, validation_rules, is_filterable, is_active, sort_order)
VALUES (
    'cpu_threads',
    '{"en": "Threads", "ru": "Потоки", "sr": "Niti"}',
    '{"en": "Threads", "ru": "Потоки", "sr": "Niti"}',
    'number',
    'regular',
    '{"min": 4, "max": 128}',
    true,
    true,
    230
) ON CONFLICT (code) DO NOTHING;

-- CPU Base Clock
INSERT INTO attributes (code, name, display_name, attribute_type, purpose, validation_rules, is_active, sort_order)
VALUES (
    'cpu_base_clock',
    '{"en": "Base Clock", "ru": "Базовая частота", "sr": "Osnovna frekvencija"}',
    '{"en": "Base Clock", "ru": "Базовая частота", "sr": "Osnovna frekvencija"}',
    'text',
    'regular',
    '{"max_length": 50, "placeholder": "3.6 GHz"}',
    true,
    240
) ON CONFLICT (code) DO NOTHING;

-- CPU Boost Clock
INSERT INTO attributes (code, name, display_name, attribute_type, purpose, validation_rules, is_active, sort_order)
VALUES (
    'cpu_boost_clock',
    '{"en": "Boost Clock", "ru": "Частота разгона", "sr": "Boost frekvencija"}',
    '{"en": "Boost Clock", "ru": "Частота разгона", "sr": "Boost frekvencija"}',
    'text',
    'regular',
    '{"max_length": 50, "placeholder": "5.0 GHz"}',
    true,
    250
) ON CONFLICT (code) DO NOTHING;

-- CPU TDP
INSERT INTO attributes (code, name, display_name, attribute_type, purpose, options, is_filterable, is_active, sort_order)
VALUES (
    'cpu_tdp',
    '{"en": "TDP", "ru": "TDP", "sr": "TDP"}',
    '{"en": "TDP", "ru": "TDP", "sr": "TDP"}',
    'select',
    'regular',
    '{
        "allowed_values": ["65W", "95W", "105W", "125W", "170W"]
    }',
    true,
    true,
    260
) ON CONFLICT (code) DO NOTHING;

-- CPU Integrated Graphics
INSERT INTO attributes (code, name, display_name, attribute_type, purpose, is_filterable, is_active, sort_order)
VALUES (
    'cpu_igpu',
    '{"en": "Integrated Graphics", "ru": "Встроенная графика", "sr": "Integrisana grafika"}',
    '{"en": "Integrated Graphics", "ru": "Встроенная графика", "sr": "Integrisana grafika"}',
    'boolean',
    'regular',
    true,
    true,
    270
) ON CONFLICT (code) DO NOTHING;

-- CPU Generation
INSERT INTO attributes (code, name, display_name, attribute_type, purpose, options, is_filterable, is_active, sort_order)
VALUES (
    'cpu_generation',
    '{"en": "Generation", "ru": "Поколение", "sr": "Generacija"}',
    '{"en": "Generation", "ru": "Поколение", "sr": "Generacija"}',
    'select',
    'regular',
    '{
        "allowed_values": ["14th Gen", "13th Gen", "12th Gen", "Zen 4", "Zen 3"]
    }',
    true,
    true,
    280
) ON CONFLICT (code) DO NOTHING;

-- ==========================================
-- 4. АТРИБУТЫ ДЛЯ RAM (RAM memorija)
-- ==========================================

-- RAM Type
INSERT INTO attributes (code, name, display_name, attribute_type, purpose, options, is_filterable, is_required, is_active, sort_order)
VALUES (
    'ram_type',
    '{"en": "Memory Type", "ru": "Тип памяти", "sr": "Tip memorije"}',
    '{"en": "Memory Type", "ru": "Тип памяти", "sr": "Tip memorije"}',
    'select',
    'variant',
    '{
        "allowed_values": ["DDR4", "DDR5"]
    }',
    true,
    true,
    true,
    300
) ON CONFLICT (code) DO NOTHING;

-- RAM Capacity
INSERT INTO attributes (code, name, display_name, attribute_type, purpose, options, is_filterable, is_required, is_active, sort_order)
VALUES (
    'ram_capacity',
    '{"en": "Capacity", "ru": "Объём", "sr": "Kapacitet"}',
    '{"en": "Capacity", "ru": "Объём", "sr": "Kapacitet"}',
    'select',
    'variant',
    '{
        "allowed_values": ["8GB", "16GB", "32GB", "64GB", "128GB"]
    }',
    true,
    true,
    true,
    310
) ON CONFLICT (code) DO NOTHING;

-- RAM Kit Configuration
INSERT INTO attributes (code, name, display_name, attribute_type, purpose, options, is_filterable, is_active, sort_order)
VALUES (
    'ram_kit',
    '{"en": "Kit Configuration", "ru": "Конфигурация комплекта", "sr": "Konfiguracija kita"}',
    '{"en": "Kit Configuration", "ru": "Конфигурация комплекта", "sr": "Konfiguracija kita"}',
    'select',
    'regular',
    '{
        "allowed_values": ["1x8GB", "2x8GB", "2x16GB", "4x16GB"]
    }',
    true,
    true,
    320
) ON CONFLICT (code) DO NOTHING;

-- RAM Speed
INSERT INTO attributes (code, name, display_name, attribute_type, purpose, options, is_filterable, is_active, sort_order)
VALUES (
    'ram_speed',
    '{"en": "Speed", "ru": "Частота", "sr": "Brzina"}',
    '{"en": "Speed", "ru": "Частота", "sr": "Brzina"}',
    'select',
    'regular',
    '{
        "allowed_values": ["2400MHz", "3200MHz", "3600MHz", "4800MHz", "5600MHz", "6000MHz"]
    }',
    true,
    true,
    330
) ON CONFLICT (code) DO NOTHING;

-- RAM CAS Latency
INSERT INTO attributes (code, name, display_name, attribute_type, purpose, validation_rules, is_active, sort_order)
VALUES (
    'ram_cas_latency',
    '{"en": "CAS Latency", "ru": "CAS задержка", "sr": "CAS latencija"}',
    '{"en": "CAS Latency", "ru": "CAS задержка", "sr": "CAS latencija"}',
    'text',
    'regular',
    '{"max_length": 20, "placeholder": "CL16"}',
    true,
    340
) ON CONFLICT (code) DO NOTHING;

-- RAM Voltage
INSERT INTO attributes (code, name, display_name, attribute_type, purpose, validation_rules, is_active, sort_order)
VALUES (
    'ram_voltage',
    '{"en": "Voltage", "ru": "Напряжение", "sr": "Napon"}',
    '{"en": "Voltage", "ru": "Напряжение", "sr": "Napon"}',
    'text',
    'regular',
    '{"max_length": 20, "placeholder": "1.35V"}',
    true,
    350
) ON CONFLICT (code) DO NOTHING;

-- RAM ECC Support
INSERT INTO attributes (code, name, display_name, attribute_type, purpose, is_filterable, is_active, sort_order)
VALUES (
    'ram_ecc',
    '{"en": "ECC Support", "ru": "Поддержка ECC", "sr": "ECC podrška"}',
    '{"en": "ECC Support", "ru": "Поддержка ECC", "sr": "ECC podrška"}',
    'boolean',
    'regular',
    true,
    true,
    360
) ON CONFLICT (code) DO NOTHING;

-- RAM Heat Spreader
INSERT INTO attributes (code, name, display_name, attribute_type, purpose, is_filterable, is_active, sort_order)
VALUES (
    'ram_heatspreader',
    '{"en": "Heat Spreader", "ru": "Радиатор", "sr": "Hladnjak"}',
    '{"en": "Heat Spreader", "ru": "Радиатор", "sr": "Hladnjak"}',
    'boolean',
    'regular',
    true,
    true,
    370
) ON CONFLICT (code) DO NOTHING;

-- ==========================================
-- 5. АТРИБУТЫ ДЛЯ SSD (SSD nakopitelji)
-- ==========================================

-- SSD Capacity
INSERT INTO attributes (code, name, display_name, attribute_type, purpose, options, is_filterable, is_required, is_active, sort_order)
VALUES (
    'ssd_capacity',
    '{"en": "Capacity", "ru": "Объём", "sr": "Kapacitet"}',
    '{"en": "Capacity", "ru": "Объём", "sr": "Kapacitet"}',
    'select',
    'variant',
    '{
        "allowed_values": ["128GB", "256GB", "500GB", "1TB", "2TB", "4TB"]
    }',
    true,
    true,
    true,
    400
) ON CONFLICT (code) DO NOTHING;

-- SSD Interface
INSERT INTO attributes (code, name, display_name, attribute_type, purpose, options, is_filterable, is_required, is_active, sort_order)
VALUES (
    'ssd_interface',
    '{"en": "Interface", "ru": "Интерфейс", "sr": "Interfejs"}',
    '{"en": "Interface", "ru": "Интерфейс", "sr": "Interfejs"}',
    'select',
    'variant',
    '{
        "allowed_values": [
            "SATA III",
            "M.2 NVMe PCIe 3.0",
            "M.2 NVMe PCIe 4.0",
            "M.2 NVMe PCIe 5.0"
        ]
    }',
    true,
    true,
    true,
    410
) ON CONFLICT (code) DO NOTHING;

-- SSD Form Factor
INSERT INTO attributes (code, name, display_name, attribute_type, purpose, options, is_filterable, is_active, sort_order)
VALUES (
    'ssd_form_factor',
    '{"en": "Form Factor", "ru": "Форм-фактор", "sr": "Forma"}',
    '{"en": "Form Factor", "ru": "Форм-фактор", "sr": "Forma"}',
    'select',
    'regular',
    '{
        "allowed_values": ["2.5\"", "M.2 2280", "M.2 2242"]
    }',
    true,
    true,
    420
) ON CONFLICT (code) DO NOTHING;

-- SSD Read Speed
INSERT INTO attributes (code, name, display_name, attribute_type, purpose, validation_rules, is_active, sort_order)
VALUES (
    'ssd_read_speed',
    '{"en": "Read Speed", "ru": "Скорость чтения", "sr": "Brzina čitanja"}',
    '{"en": "Read Speed", "ru": "Скорость чтения", "sr": "Brzina čitanja"}',
    'text',
    'regular',
    '{"max_length": 50, "placeholder": "3500 MB/s"}',
    true,
    430
) ON CONFLICT (code) DO NOTHING;

-- SSD Write Speed
INSERT INTO attributes (code, name, display_name, attribute_type, purpose, validation_rules, is_active, sort_order)
VALUES (
    'ssd_write_speed',
    '{"en": "Write Speed", "ru": "Скорость записи", "sr": "Brzina pisanja"}',
    '{"en": "Write Speed", "ru": "Скорость записи", "sr": "Brzina pisanja"}',
    'text',
    'regular',
    '{"max_length": 50, "placeholder": "3000 MB/s"}',
    true,
    440
) ON CONFLICT (code) DO NOTHING;

-- SSD NAND Type
INSERT INTO attributes (code, name, display_name, attribute_type, purpose, options, is_filterable, is_active, sort_order)
VALUES (
    'ssd_nand_type',
    '{"en": "NAND Type", "ru": "Тип NAND", "sr": "Tip NAND"}',
    '{"en": "NAND Type", "ru": "Тип NAND", "sr": "Tip NAND"}',
    'select',
    'regular',
    '{
        "allowed_values": ["TLC", "QLC", "MLC"]
    }',
    true,
    true,
    450
) ON CONFLICT (code) DO NOTHING;

-- SSD DRAM Cache
INSERT INTO attributes (code, name, display_name, attribute_type, purpose, is_filterable, is_active, sort_order)
VALUES (
    'ssd_dram_cache',
    '{"en": "DRAM Cache", "ru": "DRAM кэш", "sr": "DRAM keš"}',
    '{"en": "DRAM Cache", "ru": "DRAM кэш", "sr": "DRAM keš"}',
    'boolean',
    'regular',
    true,
    true,
    460
) ON CONFLICT (code) DO NOTHING;

-- SSD Endurance (TBW)
INSERT INTO attributes (code, name, display_name, attribute_type, purpose, validation_rules, is_active, sort_order)
VALUES (
    'ssd_endurance',
    '{"en": "Endurance (TBW)", "ru": "Ресурс (TBW)", "sr": "Izdržljivost (TBW)"}',
    '{"en": "Endurance (TBW)", "ru": "Ресурс (TBW)", "sr": "Izdržljivost (TBW)"}',
    'number',
    'regular',
    '{"min": 100, "max": 5000, "unit": "TBW"}',
    true,
    470
) ON CONFLICT (code) DO NOTHING;

-- ==========================================
-- 6. АТРИБУТЫ ДЛЯ HDD (HDD nakopitelji)
-- ==========================================

-- HDD Capacity
INSERT INTO attributes (code, name, display_name, attribute_type, purpose, options, is_filterable, is_required, is_active, sort_order)
VALUES (
    'hdd_capacity',
    '{"en": "Capacity", "ru": "Объём", "sr": "Kapacitet"}',
    '{"en": "Capacity", "ru": "Объём", "sr": "Kapacitet"}',
    'select',
    'variant',
    '{
        "allowed_values": ["500GB", "1TB", "2TB", "3TB", "4TB", "6TB", "8TB", "10TB+"]
    }',
    true,
    true,
    true,
    500
) ON CONFLICT (code) DO NOTHING;

-- HDD RPM
INSERT INTO attributes (code, name, display_name, attribute_type, purpose, options, is_filterable, is_required, is_active, sort_order)
VALUES (
    'hdd_rpm',
    '{"en": "RPM", "ru": "Скорость вращения", "sr": "RPM"}',
    '{"en": "RPM", "ru": "Скорость вращения", "sr": "RPM"}',
    'select',
    'variant',
    '{
        "allowed_values": ["5400", "7200"]
    }',
    true,
    true,
    true,
    510
) ON CONFLICT (code) DO NOTHING;

-- HDD Cache
INSERT INTO attributes (code, name, display_name, attribute_type, purpose, options, is_filterable, is_active, sort_order)
VALUES (
    'hdd_cache',
    '{"en": "Cache", "ru": "Кэш", "sr": "Keš"}',
    '{"en": "Cache", "ru": "Кэш", "sr": "Keš"}',
    'select',
    'regular',
    '{
        "allowed_values": ["64MB", "128MB", "256MB"]
    }',
    true,
    true,
    520
) ON CONFLICT (code) DO NOTHING;

-- HDD Interface
INSERT INTO attributes (code, name, display_name, attribute_type, purpose, options, is_filterable, is_active, sort_order)
VALUES (
    'hdd_interface',
    '{"en": "Interface", "ru": "Интерфейс", "sr": "Interfejs"}',
    '{"en": "Interface", "ru": "Интерфейс", "sr": "Interfejs"}',
    'select',
    'regular',
    '{
        "allowed_values": ["SATA III"]
    }',
    true,
    true,
    530
) ON CONFLICT (code) DO NOTHING;

-- HDD Form Factor
INSERT INTO attributes (code, name, display_name, attribute_type, purpose, options, is_filterable, is_active, sort_order)
VALUES (
    'hdd_form_factor',
    '{"en": "Form Factor", "ru": "Форм-фактор", "sr": "Forma"}',
    '{"en": "Form Factor", "ru": "Форм-фактор", "sr": "Forma"}',
    'select',
    'regular',
    '{
        "allowed_values": ["3.5\"", "2.5\""]
    }',
    true,
    true,
    540
) ON CONFLICT (code) DO NOTHING;

-- HDD Usage Type
INSERT INTO attributes (code, name, display_name, attribute_type, purpose, options, is_filterable, is_active, sort_order)
VALUES (
    'hdd_usage',
    '{"en": "Usage Type", "ru": "Тип использования", "sr": "Tip upotrebe"}',
    '{"en": "Usage Type", "ru": "Тип использования", "sr": "Tip upotrebe"}',
    'select',
    'regular',
    '{
        "allowed_values": ["Desktop", "NAS", "Surveillance", "Enterprise"]
    }',
    true,
    true,
    550
) ON CONFLICT (code) DO NOTHING;

-- ==========================================
-- 7. АТРИБУТЫ ДЛЯ МАТЕРИНСКИХ ПЛАТ (Matične ploče)
-- ==========================================

-- Motherboard Socket
INSERT INTO attributes (code, name, display_name, attribute_type, purpose, options, is_filterable, is_required, is_active, sort_order)
VALUES (
    'mb_socket',
    '{"en": "Socket", "ru": "Сокет", "sr": "Soket"}',
    '{"en": "Socket", "ru": "Сокет", "sr": "Soket"}',
    'select',
    'regular',
    '{
        "allowed_values": ["LGA1700", "LGA1200", "AM5", "AM4"]
    }',
    true,
    true,
    true,
    600
) ON CONFLICT (code) DO NOTHING;

-- Motherboard Chipset
INSERT INTO attributes (code, name, display_name, attribute_type, purpose, options, is_filterable, is_required, is_active, sort_order)
VALUES (
    'mb_chipset',
    '{"en": "Chipset", "ru": "Чипсет", "sr": "Čipset"}',
    '{"en": "Chipset", "ru": "Чипсет", "sr": "Čipset"}',
    'select',
    'variant',
    '{
        "allowed_values": ["Z790", "B760", "H610", "X670E", "B650", "A620"]
    }',
    true,
    true,
    true,
    610
) ON CONFLICT (code) DO NOTHING;

-- Motherboard Form Factor
INSERT INTO attributes (code, name, display_name, attribute_type, purpose, options, is_filterable, is_required, is_active, sort_order)
VALUES (
    'mb_form_factor',
    '{"en": "Form Factor", "ru": "Форм-фактор", "sr": "Forma"}',
    '{"en": "Form Factor", "ru": "Форм-фактор", "sr": "Forma"}',
    'select',
    'variant',
    '{
        "allowed_values": ["ATX", "Micro-ATX", "Mini-ITX", "E-ATX"]
    }',
    true,
    true,
    true,
    620
) ON CONFLICT (code) DO NOTHING;

-- Motherboard Memory Type
INSERT INTO attributes (code, name, display_name, attribute_type, purpose, options, is_filterable, is_required, is_active, sort_order)
VALUES (
    'mb_memory_type',
    '{"en": "Memory Type", "ru": "Тип памяти", "sr": "Tip memorije"}',
    '{"en": "Memory Type", "ru": "Тип памяти", "sr": "Tip memorije"}',
    'select',
    'regular',
    '{
        "allowed_values": ["DDR4", "DDR5"]
    }',
    true,
    true,
    true,
    630
) ON CONFLICT (code) DO NOTHING;

-- Motherboard Memory Slots
INSERT INTO attributes (code, name, display_name, attribute_type, purpose, validation_rules, is_filterable, is_active, sort_order)
VALUES (
    'mb_memory_slots',
    '{"en": "Memory Slots", "ru": "Слоты памяти", "sr": "Slotovi memorije"}',
    '{"en": "Memory Slots", "ru": "Слоты памяти", "sr": "Slotovi memorije"}',
    'number',
    'regular',
    '{"min": 2, "max": 8}',
    true,
    true,
    640
) ON CONFLICT (code) DO NOTHING;

-- Motherboard Max Memory
INSERT INTO attributes (code, name, display_name, attribute_type, purpose, options, is_filterable, is_active, sort_order)
VALUES (
    'mb_max_memory',
    '{"en": "Max Memory", "ru": "Макс. память", "sr": "Maksimalna memorija"}',
    '{"en": "Max Memory", "ru": "Макс. память", "sr": "Maksimalna memorija"}',
    'select',
    'regular',
    '{
        "allowed_values": ["64GB", "128GB", "192GB"]
    }',
    true,
    true,
    650
) ON CONFLICT (code) DO NOTHING;

-- Motherboard M.2 Slots
INSERT INTO attributes (code, name, display_name, attribute_type, purpose, validation_rules, is_filterable, is_active, sort_order)
VALUES (
    'mb_m2_slots',
    '{"en": "M.2 Slots", "ru": "M.2 слоты", "sr": "M.2 slotovi"}',
    '{"en": "M.2 Slots", "ru": "M.2 слоты", "sr": "M.2 slotovi"}',
    'number',
    'regular',
    '{"min": 0, "max": 6}',
    true,
    true,
    660
) ON CONFLICT (code) DO NOTHING;

-- Motherboard SATA Ports
INSERT INTO attributes (code, name, display_name, attribute_type, purpose, validation_rules, is_filterable, is_active, sort_order)
VALUES (
    'mb_sata_ports',
    '{"en": "SATA Ports", "ru": "Порты SATA", "sr": "SATA portovi"}',
    '{"en": "SATA Ports", "ru": "Порты SATA", "sr": "SATA portovi"}',
    'number',
    'regular',
    '{"min": 2, "max": 12}',
    true,
    true,
    670
) ON CONFLICT (code) DO NOTHING;

-- Motherboard WiFi
INSERT INTO attributes (code, name, display_name, attribute_type, purpose, is_filterable, is_active, sort_order)
VALUES (
    'mb_wifi',
    '{"en": "WiFi", "ru": "WiFi", "sr": "WiFi"}',
    '{"en": "WiFi", "ru": "WiFi", "sr": "WiFi"}',
    'boolean',
    'regular',
    true,
    true,
    680
) ON CONFLICT (code) DO NOTHING;

-- Motherboard Bluetooth
INSERT INTO attributes (code, name, display_name, attribute_type, purpose, is_filterable, is_active, sort_order)
VALUES (
    'mb_bluetooth',
    '{"en": "Bluetooth", "ru": "Bluetooth", "sr": "Bluetooth"}',
    '{"en": "Bluetooth", "ru": "Bluetooth", "sr": "Bluetooth"}',
    'boolean',
    'regular',
    true,
    true,
    690
) ON CONFLICT (code) DO NOTHING;

-- ==========================================
-- 8. АТРИБУТЫ ДЛЯ PSU (Napajanja)
-- ==========================================

-- PSU Wattage
INSERT INTO attributes (code, name, display_name, attribute_type, purpose, options, is_filterable, is_required, is_active, sort_order)
VALUES (
    'psu_wattage',
    '{"en": "Wattage", "ru": "Мощность", "sr": "Snaga"}',
    '{"en": "Wattage", "ru": "Мощность", "sr": "Snaga"}',
    'select',
    'variant',
    '{
        "allowed_values": ["500W", "650W", "750W", "850W", "1000W", "1200W+"]
    }',
    true,
    true,
    true,
    700
) ON CONFLICT (code) DO NOTHING;

-- PSU Efficiency Rating
INSERT INTO attributes (code, name, display_name, attribute_type, purpose, options, is_filterable, is_required, is_active, sort_order)
VALUES (
    'psu_efficiency',
    '{"en": "Efficiency Rating", "ru": "Эффективность", "sr": "Ocena efikasnosti"}',
    '{"en": "Efficiency Rating", "ru": "Эффективность", "sr": "Ocena efikasnosti"}',
    'select',
    'variant',
    '{
        "allowed_values": ["80+ Bronze", "80+ Silver", "80+ Gold", "80+ Platinum", "80+ Titanium"]
    }',
    true,
    true,
    true,
    710
) ON CONFLICT (code) DO NOTHING;

-- PSU Modular Type
INSERT INTO attributes (code, name, display_name, attribute_type, purpose, options, is_filterable, is_active, sort_order)
VALUES (
    'psu_modular',
    '{"en": "Modular Type", "ru": "Модульность", "sr": "Tip modularnosti"}',
    '{"en": "Modular Type", "ru": "Модульность", "sr": "Tip modularnosti"}',
    'select',
    'regular',
    '{
        "allowed_values": ["Non-Modular", "Semi-Modular", "Fully Modular"]
    }',
    true,
    true,
    720
) ON CONFLICT (code) DO NOTHING;

-- PSU Form Factor
INSERT INTO attributes (code, name, display_name, attribute_type, purpose, options, is_filterable, is_active, sort_order)
VALUES (
    'psu_form_factor',
    '{"en": "Form Factor", "ru": "Форм-фактор", "sr": "Forma"}',
    '{"en": "Form Factor", "ru": "Форм-фактор", "sr": "Forma"}',
    'select',
    'regular',
    '{
        "allowed_values": ["ATX", "SFX", "SFX-L"]
    }',
    true,
    true,
    730
) ON CONFLICT (code) DO NOTHING;

-- PSU PCIe 5.0 Ready
INSERT INTO attributes (code, name, display_name, attribute_type, purpose, is_filterable, is_active, sort_order)
VALUES (
    'psu_pcie5',
    '{"en": "PCIe 5.0 Ready", "ru": "PCIe 5.0 Ready", "sr": "PCIe 5.0 spremnost"}',
    '{"en": "PCIe 5.0 Ready (12VHPWR)", "ru": "PCIe 5.0 Ready (12VHPWR)", "sr": "PCIe 5.0 spremnost (12VHPWR)"}',
    'boolean',
    'regular',
    true,
    true,
    740
) ON CONFLICT (code) DO NOTHING;

-- PSU Fan Size
INSERT INTO attributes (code, name, display_name, attribute_type, purpose, options, is_filterable, is_active, sort_order)
VALUES (
    'psu_fan_size',
    '{"en": "Fan Size", "ru": "Размер вентилятора", "sr": "Veličina ventilatora"}',
    '{"en": "Fan Size", "ru": "Размер вентилятора", "sr": "Veličina ventilatora"}',
    'select',
    'regular',
    '{
        "allowed_values": ["120mm", "140mm"]
    }',
    true,
    true,
    750
) ON CONFLICT (code) DO NOTHING;

-- ==========================================
-- 9. АТРИБУТЫ ДЛЯ КОРПУСОВ (Kućišta)
-- ==========================================

-- Case Form Factor
INSERT INTO attributes (code, name, display_name, attribute_type, purpose, options, is_filterable, is_required, is_active, sort_order)
VALUES (
    'case_form_factor',
    '{"en": "Form Factor", "ru": "Форм-фактор", "sr": "Forma"}',
    '{"en": "Form Factor", "ru": "Форм-фактор", "sr": "Forma"}',
    'select',
    'variant',
    '{
        "allowed_values": ["Full Tower", "Mid Tower", "Mini Tower", "Mini-ITX"]
    }',
    true,
    true,
    true,
    800
) ON CONFLICT (code) DO NOTHING;

-- Case Max GPU Length
INSERT INTO attributes (code, name, display_name, attribute_type, purpose, validation_rules, is_filterable, is_active, sort_order)
VALUES (
    'case_max_gpu_length',
    '{"en": "Max GPU Length", "ru": "Макс. длина GPU", "sr": "Maksimalna dužina GPU"}',
    '{"en": "Max GPU Length (mm)", "ru": "Макс. длина GPU (мм)", "sr": "Maksimalna dužina GPU (mm)"}',
    'number',
    'regular',
    '{"min": 250, "max": 450, "unit": "mm"}',
    true,
    true,
    810
) ON CONFLICT (code) DO NOTHING;

-- Case Max CPU Cooler Height
INSERT INTO attributes (code, name, display_name, attribute_type, purpose, validation_rules, is_filterable, is_active, sort_order)
VALUES (
    'case_max_cpu_height',
    '{"en": "Max CPU Cooler Height", "ru": "Макс. высота кулера", "sr": "Maksimalna visina hladnjaka"}',
    '{"en": "Max CPU Cooler Height (mm)", "ru": "Макс. высота кулера (мм)", "sr": "Maksimalna visina hladnjaka (mm)"}',
    'number',
    'regular',
    '{"min": 140, "max": 200, "unit": "mm"}',
    true,
    true,
    820
) ON CONFLICT (code) DO NOTHING;

-- Case Fans Included
INSERT INTO attributes (code, name, display_name, attribute_type, purpose, validation_rules, is_filterable, is_active, sort_order)
VALUES (
    'case_fans_included',
    '{"en": "Fans Included", "ru": "Вентиляторов в комплекте", "sr": "Uključeni ventilatori"}',
    '{"en": "Fans Included", "ru": "Вентиляторов в комплекте", "sr": "Uključeni ventilatori"}',
    'number',
    'regular',
    '{"min": 0, "max": 10}',
    true,
    true,
    830
) ON CONFLICT (code) DO NOTHING;

-- Case Max Fans Supported
INSERT INTO attributes (code, name, display_name, attribute_type, purpose, validation_rules, is_filterable, is_active, sort_order)
VALUES (
    'case_max_fans',
    '{"en": "Max Fans Supported", "ru": "Макс. вентиляторов", "sr": "Maksimalno podržanih ventilatora"}',
    '{"en": "Max Fans Supported", "ru": "Макс. вентиляторов", "sr": "Maksimalno podržanih ventilatora"}',
    'number',
    'regular',
    '{"min": 4, "max": 15}',
    true,
    true,
    840
) ON CONFLICT (code) DO NOTHING;

-- Case Tempered Glass
INSERT INTO attributes (code, name, display_name, attribute_type, purpose, is_filterable, is_active, sort_order)
VALUES (
    'case_tempered_glass',
    '{"en": "Tempered Glass", "ru": "Закалённое стекло", "sr": "Kaljeno staklo"}',
    '{"en": "Tempered Glass", "ru": "Закалённое стекло", "sr": "Kaljeno staklo"}',
    'boolean',
    'regular',
    true,
    true,
    850
) ON CONFLICT (code) DO NOTHING;

-- Case RGB Fans
INSERT INTO attributes (code, name, display_name, attribute_type, purpose, is_filterable, is_active, sort_order)
VALUES (
    'case_rgb_fans',
    '{"en": "RGB Fans", "ru": "RGB вентиляторы", "sr": "RGB ventilatori"}',
    '{"en": "RGB Fans", "ru": "RGB вентиляторы", "sr": "RGB ventilatori"}',
    'boolean',
    'regular',
    true,
    true,
    860
) ON CONFLICT (code) DO NOTHING;

-- Case Dust Filters
INSERT INTO attributes (code, name, display_name, attribute_type, purpose, is_filterable, is_active, sort_order)
VALUES (
    'case_dust_filters',
    '{"en": "Dust Filters", "ru": "Пылевые фильтры", "sr": "Filteri za prašinu"}',
    '{"en": "Dust Filters", "ru": "Пылевые фильтры", "sr": "Filteri za prašinu"}',
    'boolean',
    'regular',
    true,
    true,
    870
) ON CONFLICT (code) DO NOTHING;

-- ==========================================
-- 10. АТРИБУТЫ ДЛЯ ОХЛАЖДЕНИЯ (Hlađenje)
-- ==========================================

-- Cooler Type
INSERT INTO attributes (code, name, display_name, attribute_type, purpose, options, is_filterable, is_required, is_active, sort_order)
VALUES (
    'cooler_type',
    '{"en": "Cooler Type", "ru": "Тип охлаждения", "sr": "Tip hlađenja"}',
    '{"en": "Cooler Type", "ru": "Тип охлаждения", "sr": "Tip hlađenja"}',
    'select',
    'variant',
    '{
        "allowed_values": [
            "Air Tower",
            "Low Profile",
            "AIO 120mm",
            "AIO 240mm",
            "AIO 280mm",
            "AIO 360mm"
        ]
    }',
    true,
    true,
    true,
    900
) ON CONFLICT (code) DO NOTHING;

-- Cooler Radiator Size
INSERT INTO attributes (code, name, display_name, attribute_type, purpose, options, is_filterable, is_active, sort_order)
VALUES (
    'cooler_radiator',
    '{"en": "Radiator Size", "ru": "Размер радиатора", "sr": "Veličina radijatora"}',
    '{"en": "Radiator Size", "ru": "Размер радиатора", "sr": "Veličina radijatora"}',
    'select',
    'regular',
    '{
        "allowed_values": ["120mm", "240mm", "280mm", "360mm", "420mm"]
    }',
    true,
    true,
    910
) ON CONFLICT (code) DO NOTHING;

-- Cooler Fan Size
INSERT INTO attributes (code, name, display_name, attribute_type, purpose, options, is_filterable, is_active, sort_order)
VALUES (
    'cooler_fan_size',
    '{"en": "Fan Size", "ru": "Размер вентилятора", "sr": "Veličina ventilatora"}',
    '{"en": "Fan Size", "ru": "Размер вентилятора", "sr": "Veličina ventilatora"}',
    'select',
    'regular',
    '{
        "allowed_values": ["120mm", "140mm"]
    }',
    true,
    true,
    920
) ON CONFLICT (code) DO NOTHING;

-- Cooler Max TDP
INSERT INTO attributes (code, name, display_name, attribute_type, purpose, options, is_filterable, is_active, sort_order)
VALUES (
    'cooler_max_tdp',
    '{"en": "Max TDP", "ru": "Макс. TDP", "sr": "Maksimalni TDP"}',
    '{"en": "Max TDP", "ru": "Макс. TDP", "sr": "Maksimalni TDP"}',
    'select',
    'regular',
    '{
        "allowed_values": ["150W", "200W", "250W", "300W+"]
    }',
    true,
    true,
    930
) ON CONFLICT (code) DO NOTHING;

-- Cooler RGB Lighting
INSERT INTO attributes (code, name, display_name, attribute_type, purpose, is_filterable, is_active, sort_order)
VALUES (
    'cooler_rgb',
    '{"en": "RGB Lighting", "ru": "RGB подсветка", "sr": "RGB osvetljenje"}',
    '{"en": "RGB Lighting", "ru": "RGB подсветка", "sr": "RGB osvetljenje"}',
    'boolean',
    'regular',
    true,
    true,
    940
) ON CONFLICT (code) DO NOTHING;

-- Cooler Noise Level
INSERT INTO attributes (code, name, display_name, attribute_type, purpose, validation_rules, is_filterable, is_active, sort_order)
VALUES (
    'cooler_noise',
    '{"en": "Noise Level", "ru": "Уровень шума", "sr": "Nivo buke"}',
    '{"en": "Noise Level (dB)", "ru": "Уровень шума (дБ)", "sr": "Nivo buke (dB)"}',
    'number',
    'regular',
    '{"min": 10, "max": 50, "unit": "dB"}',
    true,
    true,
    950
) ON CONFLICT (code) DO NOTHING;

COMMIT;

-- ==========================================
-- ИТОГОВАЯ СТАТИСТИКА
-- ==========================================

-- Подсчет созданных атрибутов по типам
DO $$
DECLARE
    total_count INTEGER;
    general_count INTEGER;
    gpu_count INTEGER;
    cpu_count INTEGER;
    ram_count INTEGER;
    ssd_count INTEGER;
    hdd_count INTEGER;
    mb_count INTEGER;
    psu_count INTEGER;
    case_count INTEGER;
    cooler_count INTEGER;
BEGIN
    -- Общие атрибуты
    SELECT COUNT(*) INTO general_count FROM attributes WHERE code LIKE 'pc_%';

    -- Видеокарты
    SELECT COUNT(*) INTO gpu_count FROM attributes WHERE code LIKE 'gpu_%';

    -- Процессоры
    SELECT COUNT(*) INTO cpu_count FROM attributes WHERE code LIKE 'cpu_%';

    -- RAM
    SELECT COUNT(*) INTO ram_count FROM attributes WHERE code LIKE 'ram_%';

    -- SSD
    SELECT COUNT(*) INTO ssd_count FROM attributes WHERE code LIKE 'ssd_%';

    -- HDD
    SELECT COUNT(*) INTO hdd_count FROM attributes WHERE code LIKE 'hdd_%';

    -- Материнские платы
    SELECT COUNT(*) INTO mb_count FROM attributes WHERE code LIKE 'mb_%';

    -- PSU
    SELECT COUNT(*) INTO psu_count FROM attributes WHERE code LIKE 'psu_%';

    -- Корпуса
    SELECT COUNT(*) INTO case_count FROM attributes WHERE code LIKE 'case_%';

    -- Охлаждение
    SELECT COUNT(*) INTO cooler_count FROM attributes WHERE code LIKE 'cooler_%';

    total_count := general_count + gpu_count + cpu_count + ram_count + ssd_count +
                   hdd_count + mb_count + psu_count + case_count + cooler_count;

    RAISE NOTICE '✅ Migration completed successfully!';
    RAISE NOTICE '';
    RAISE NOTICE 'Created % attributes for computer components:', total_count;
    RAISE NOTICE '  - General (pc_*): %', general_count;
    RAISE NOTICE '  - Graphics Cards (gpu_*): %', gpu_count;
    RAISE NOTICE '  - Processors (cpu_*): %', cpu_count;
    RAISE NOTICE '  - RAM (ram_*): %', ram_count;
    RAISE NOTICE '  - SSD (ssd_*): %', ssd_count;
    RAISE NOTICE '  - HDD (hdd_*): %', hdd_count;
    RAISE NOTICE '  - Motherboards (mb_*): %', mb_count;
    RAISE NOTICE '  - Power Supplies (psu_*): %', psu_count;
    RAISE NOTICE '  - Cases (case_*): %', case_count;
    RAISE NOTICE '  - Cooling (cooler_*): %', cooler_count;
    RAISE NOTICE '';
    RAISE NOTICE '⚠️  IMPORTANT: These attributes are ready for use once the category restructuring (Phase 0) is completed.';
    RAISE NOTICE 'After creating the proper subcategories (graficke-kartice, procesori, etc.),';
    RAISE NOTICE 'you will need to link these attributes to categories via category_attributes table.';
END $$;
