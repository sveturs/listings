# SQL Миграции: Атрибуты для категории "Računarske komponente"

## 📋 Обзор

Этот набор миграций создаёт **84+ атрибута** для 9 подкатегорий компьютерных компонентов:

1. **Grafičke kartice** (Видеокарты) - 10 специфичных атрибутов
2. **Procesori** (Процессоры) - 9 специфичных атрибутов
3. **RAM memorija** (Оперативная память) - 8 специфичных атрибутов
4. **SSD nakopitelji** (SSD накопители) - 8 специфичных атрибутов
5. **HDD nakopitelji** (HDD накопители) - 6 специфичных атрибутов
6. **Matične ploče** (Материнские платы) - 11 специфичных атрибутов
7. **Napajanja** (Блоки питания) - 6 специфичных атрибутов
8. **Kućišta** (Корпуса) - 8 специфичных атрибутов
9. **Hlađenje** (Охлаждение) - 6 специфичных атрибутов

**Плюс 6 общих атрибутов** для всех подкатегорий.

## 📂 Файлы миграций

| Миграция | Описание | Статус |
|----------|----------|--------|
| `000022_create_racunarske_komponente_attributes.up.sql` | Создание 84+ атрибутов | ✅ Готова |
| `000022_create_racunarske_komponente_attributes.down.sql` | Откат создания атрибутов | ✅ Готова |
| `000023_link_racunarske_komponente_attributes_to_categories.up.sql` | Привязка атрибутов к категориям | ⚠️ Требует ФАЗУ 0 |
| `000023_link_racunarske_komponente_attributes_to_categories.down.sql` | Откат привязки атрибутов | ✅ Готова |

## 🚀 Порядок применения миграций

### ⚠️ ВАЖНО: Последовательность действий

```
ФАЗА 0 (КРИТИЧНО!) → ФАЗА 1 (эти миграции) → ФАЗА 2 (привязка)
```

**Эти миграции НЕ могут быть применены раньше ФАЗЫ 0!**

### ФАЗА 0: Переструктурирование категорий (сначала!)

**Перед применением этих миграций необходимо:**

1. Удалить 20 старых конкретных подкатегорий:
   - `graficke-rtx-4000`, `graficke-rtx-3000`, `procesor-intel-i9`, `ram-ddr5-32gb`, etc.

2. Создать 9 новых общих подкатегорий:
   - `graficke-kartice`, `procesori`, `ram-memorija`, `ssd-nakopitelji`, `hdd-nakopitelji`, `maticne-ploce`, `napajanja`, `kucista`, `hladjenje`

📖 **Подробнее:** `/p/github.com/vondi-global/passport/category-audit-racunarske-komponente-action-plan.md` (раздел ФАЗА 0)

---

### ФАЗА 1: Создание атрибутов

**После завершения ФАЗЫ 0 применить:**

```bash
# Миграция 000022 - Создание атрибутов
psql "postgres://listings_user:listings_secret@localhost:35434/listings_dev_db?sslmode=disable" \
  -f /p/github.com/vondi-global/listings/migrations/000022_create_racunarske_komponente_attributes.up.sql
```

**Что будет создано:**
- 6 общих атрибутов (pc_brand, pc_model, pc_condition, pc_warranty, pc_color, pc_rgb_lighting)
- 10 атрибутов для видеокарт (gpu_series, gpu_chip, gpu_vram, и т.д.)
- 9 атрибутов для процессоров (cpu_series, cpu_socket, cpu_cores, и т.д.)
- 8 атрибутов для RAM (ram_type, ram_capacity, ram_speed, и т.д.)
- 8 атрибутов для SSD (ssd_capacity, ssd_interface, ssd_read_speed, и т.д.)
- 6 атрибутов для HDD (hdd_capacity, hdd_rpm, hdd_cache, и т.д.)
- 11 атрибутов для материнских плат (mb_socket, mb_chipset, mb_form_factor, и т.д.)
- 6 атрибутов для PSU (psu_wattage, psu_efficiency, psu_modular, и т.д.)
- 8 атрибутов для корпусов (case_form_factor, case_max_gpu_length, и т.д.)
- 6 атрибутов для охлаждения (cooler_type, cooler_radiator, cooler_fan_size, и т.д.)

**Проверка:**
```sql
-- Должно быть 84+ атрибутов
SELECT COUNT(*) FROM attributes WHERE code LIKE 'pc_%'
   OR code LIKE 'gpu_%'
   OR code LIKE 'cpu_%'
   OR code LIKE 'ram_%'
   OR code LIKE 'ssd_%'
   OR code LIKE 'hdd_%'
   OR code LIKE 'mb_%'
   OR code LIKE 'psu_%'
   OR code LIKE 'case_%'
   OR code LIKE 'cooler_%';
```

---

### ФАЗА 2: Привязка атрибутов к категориям

**После применения миграции 000022 применить:**

```bash
# Миграция 000023 - Привязка атрибутов к категориям
psql "postgres://listings_user:listings_secret@localhost:35434/listings_dev_db?sslmode=disable" \
  -f /p/github.com/vondi-global/listings/migrations/000023_link_racunarske_komponente_attributes_to_categories.up.sql
```

**Что будет создано:**
- Связи `category_attributes` для регулярных атрибутов
- Связи `category_variant_attributes` для вариативных атрибутов (14 вариантов: gpu_vram, cpu_cores, ram_type, ram_capacity, ssd_capacity, ssd_interface, hdd_capacity, hdd_rpm, mb_chipset, mb_form_factor, psu_wattage, psu_efficiency, case_form_factor, cooler_type)

**Проверка:**
```sql
-- Проверить количество связей для каждой категории
SELECT
    c.slug,
    c.name->>'sr' as category_name,
    COUNT(ca.id) as attributes_count
FROM categories c
LEFT JOIN category_attributes ca ON ca.category_id = c.id
WHERE c.slug IN (
    'graficke-kartice', 'procesori', 'ram-memorija',
    'ssd-nakopitelji', 'hdd-nakopitelji', 'maticne-ploce',
    'napajanja', 'kucista', 'hladjenje'
)
GROUP BY c.id, c.slug, c.name
ORDER BY c.slug;

-- Проверить вариативные атрибуты
SELECT
    c.slug,
    a.code,
    cva.is_required,
    cva.affects_price,
    cva.affects_stock
FROM category_variant_attributes cva
JOIN categories c ON cva.category_id = c.id
JOIN attributes a ON cva.attribute_id = a.id
WHERE c.slug IN (
    'graficke-kartice', 'procesori', 'ram-memorija',
    'ssd-nakopitelji', 'hdd-nakopitelji', 'maticne-ploce',
    'napajanja', 'kucista', 'hladjenje'
)
ORDER BY c.slug, a.code;
```

---

## 🔄 Откат миграций (Rollback)

### Откат ФАЗЫ 2 (привязка атрибутов)

```bash
psql "postgres://listings_user:listings_secret@localhost:35434/listings_dev_db?sslmode=disable" \
  -f /p/github.com/vondi-global/listings/migrations/000023_link_racunarske_komponente_attributes_to_categories.down.sql
```

### Откат ФАЗЫ 1 (атрибуты)

```bash
psql "postgres://listings_user:listings_secret@localhost:35434/listings_dev_db?sslmode=disable" \
  -f /p/github.com/vondi-global/listings/migrations/000022_create_racunarske_komponente_attributes.down.sql
```

**⚠️ ВАЖНО:** Откат удаляет атрибуты и все их связи, но НЕ удаляет данные из `listing_attribute_values`! Если листинги уже используют эти атрибуты, их значения будут удалены.

---

## 📊 Детали атрибутов

### 1. Общие атрибуты (6)

| Code | Name | Type | Required | Filterable | Searchable |
|------|------|------|----------|------------|------------|
| `pc_brand` | Brand / Brend | select | ✅ | ✅ | ✅ |
| `pc_model` | Model / Model | text | ✅ | ❌ | ✅ |
| `pc_condition` | Condition / Stanje | select | ✅ | ✅ | ❌ |
| `pc_warranty` | Warranty / Garancija | select | ❌ | ✅ | ❌ |
| `pc_color` | Color / Boja | select | ❌ | ✅ | ❌ |
| `pc_rgb_lighting` | RGB Lighting | boolean | ❌ | ✅ | ❌ |

---

### 2. Видеокарты (10 атрибутов)

| Code | Name | Type | Variant | Required | Filterable |
|------|------|------|---------|----------|------------|
| `gpu_series` | GPU Series / GPU serija | select | ❌ | ✅ | ✅ |
| `gpu_chip` | GPU Chip / GPU čip | select | ❌ | ❌ | ✅ |
| `gpu_vram` | VRAM / Video memorija | select | ✅ | ✅ | ✅ |
| `gpu_memory_type` | Memory Type / Tip memorije | select | ❌ | ❌ | ✅ |
| `gpu_memory_bus` | Memory Bus / Magistrala | select | ❌ | ❌ | ✅ |
| `gpu_cooling` | Cooling Type / Tip hlađenja | select | ❌ | ❌ | ✅ |
| `gpu_power` | Power Connectors / Konektori | text | ❌ | ❌ | ❌ |
| `gpu_tdp` | TDP | select | ❌ | ❌ | ✅ |
| `gpu_length` | Length (mm) / Dužina | number | ❌ | ❌ | ✅ |
| `gpu_ray_tracing` | Ray Tracing | boolean | ❌ | ❌ | ✅ |

**Опции для select:**
- `gpu_series`: RTX 4000/3000/2000, GTX 1000, RX 7000/6000/5000, Intel Arc
- `gpu_chip`: RTX 4090/4080/4070 Ti/4070/4060 Ti, RTX 3090 Ti/3090/3080 Ti/3080, RX 7900 XTX/XT, RX 6900 XT/6800 XT, Other
- `gpu_vram`: 6GB, 8GB, 12GB, 16GB, 24GB
- `gpu_memory_type`: GDDR6, GDDR6X
- `gpu_memory_bus`: 128-bit, 192-bit, 256-bit, 384-bit
- `gpu_cooling`: Air, Hybrid, Water
- `gpu_tdp`: 150W, 200W, 250W, 300W, 350W, 450W

---

### 3. Процессоры (9 атрибутов)

| Code | Name | Type | Variant | Required | Filterable |
|------|------|------|---------|----------|------------|
| `cpu_series` | Processor Series / Serija procesora | select | ❌ | ✅ | ✅ |
| `cpu_socket` | Socket / Soket | select | ❌ | ✅ | ✅ |
| `cpu_cores` | Cores / Jezgra | number | ✅ | ✅ | ✅ |
| `cpu_threads` | Threads / Niti | number | ❌ | ❌ | ✅ |
| `cpu_base_clock` | Base Clock / Osnovna frekvencija | text | ❌ | ❌ | ❌ |
| `cpu_boost_clock` | Boost Clock / Boost frekvencija | text | ❌ | ❌ | ❌ |
| `cpu_tdp` | TDP | select | ❌ | ❌ | ✅ |
| `cpu_igpu` | Integrated Graphics / Integrisana grafika | boolean | ❌ | ❌ | ✅ |
| `cpu_generation` | Generation / Generacija | select | ❌ | ❌ | ✅ |

**Опции для select:**
- `cpu_series`: Intel Core i9/i7/i5/i3, AMD Ryzen 9/7/5/3, AMD Threadripper
- `cpu_socket`: LGA1700, LGA1200, AM5, AM4
- `cpu_tdp`: 65W, 95W, 105W, 125W, 170W
- `cpu_generation`: 14th Gen, 13th Gen, 12th Gen, Zen 4, Zen 3

---

### 4. RAM (8 атрибутов)

| Code | Name | Type | Variant | Required | Filterable |
|------|------|------|---------|----------|------------|
| `ram_type` | Memory Type / Tip memorije | select | ✅ | ✅ | ✅ |
| `ram_capacity` | Capacity / Kapacitet | select | ✅ | ✅ | ✅ |
| `ram_kit` | Kit Configuration / Konfiguracija kita | select | ❌ | ❌ | ✅ |
| `ram_speed` | Speed / Brzina | select | ❌ | ❌ | ✅ |
| `ram_cas_latency` | CAS Latency / CAS latencija | text | ❌ | ❌ | ❌ |
| `ram_voltage` | Voltage / Napon | text | ❌ | ❌ | ❌ |
| `ram_ecc` | ECC Support / ECC podrška | boolean | ❌ | ❌ | ✅ |
| `ram_heatspreader` | Heat Spreader / Hladnjak | boolean | ❌ | ❌ | ✅ |

**Опции для select:**
- `ram_type`: DDR4, DDR5
- `ram_capacity`: 8GB, 16GB, 32GB, 64GB, 128GB
- `ram_kit`: 1x8GB, 2x8GB, 2x16GB, 4x16GB
- `ram_speed`: 2400MHz, 3200MHz, 3600MHz, 4800MHz, 5600MHz, 6000MHz

---

### 5. SSD (8 атрибутов)

| Code | Name | Type | Variant | Required | Filterable |
|------|------|------|---------|----------|------------|
| `ssd_capacity` | Capacity / Kapacitet | select | ✅ | ✅ | ✅ |
| `ssd_interface` | Interface / Interfejs | select | ✅ | ✅ | ✅ |
| `ssd_form_factor` | Form Factor / Forma | select | ❌ | ❌ | ✅ |
| `ssd_read_speed` | Read Speed / Brzina čitanja | text | ❌ | ❌ | ❌ |
| `ssd_write_speed` | Write Speed / Brzina pisanja | text | ❌ | ❌ | ❌ |
| `ssd_nand_type` | NAND Type / Tip NAND | select | ❌ | ❌ | ✅ |
| `ssd_dram_cache` | DRAM Cache / DRAM keš | boolean | ❌ | ❌ | ✅ |
| `ssd_endurance` | Endurance (TBW) / Izdržljivost | number | ❌ | ❌ | ✅ |

**Опции для select:**
- `ssd_capacity`: 128GB, 256GB, 500GB, 1TB, 2TB, 4TB
- `ssd_interface`: SATA III, M.2 NVMe PCIe 3.0, M.2 NVMe PCIe 4.0, M.2 NVMe PCIe 5.0
- `ssd_form_factor`: 2.5", M.2 2280, M.2 2242
- `ssd_nand_type`: TLC, QLC, MLC

---

### 6. HDD (6 атрибутов)

| Code | Name | Type | Variant | Required | Filterable |
|------|------|------|---------|----------|------------|
| `hdd_capacity` | Capacity / Kapacitet | select | ✅ | ✅ | ✅ |
| `hdd_rpm` | RPM | select | ✅ | ✅ | ✅ |
| `hdd_cache` | Cache / Keš | select | ❌ | ❌ | ✅ |
| `hdd_interface` | Interface / Interfejs | select | ❌ | ❌ | ✅ |
| `hdd_form_factor` | Form Factor / Forma | select | ❌ | ❌ | ✅ |
| `hdd_usage` | Usage Type / Tip upotrebe | select | ❌ | ❌ | ✅ |

**Опции для select:**
- `hdd_capacity`: 500GB, 1TB, 2TB, 3TB, 4TB, 6TB, 8TB, 10TB+
- `hdd_rpm`: 5400, 7200
- `hdd_cache`: 64MB, 128MB, 256MB
- `hdd_interface`: SATA III
- `hdd_form_factor`: 3.5", 2.5"
- `hdd_usage`: Desktop, NAS, Surveillance, Enterprise

---

### 7. Материнские платы (11 атрибутов)

| Code | Name | Type | Variant | Required | Filterable |
|------|------|------|---------|----------|------------|
| `mb_socket` | Socket / Soket | select | ❌ | ✅ | ✅ |
| `mb_chipset` | Chipset / Čipset | select | ✅ | ✅ | ✅ |
| `mb_form_factor` | Form Factor / Forma | select | ✅ | ✅ | ✅ |
| `mb_memory_type` | Memory Type / Tip memorije | select | ❌ | ✅ | ✅ |
| `mb_memory_slots` | Memory Slots / Slotovi memorije | number | ❌ | ❌ | ✅ |
| `mb_max_memory` | Max Memory / Maksimalna memorija | select | ❌ | ❌ | ✅ |
| `mb_m2_slots` | M.2 Slots / M.2 slotovi | number | ❌ | ❌ | ✅ |
| `mb_sata_ports` | SATA Ports / SATA portovi | number | ❌ | ❌ | ✅ |
| `mb_wifi` | WiFi | boolean | ❌ | ❌ | ✅ |
| `mb_bluetooth` | Bluetooth | boolean | ❌ | ❌ | ✅ |

**Опции для select:**
- `mb_socket`: LGA1700, LGA1200, AM5, AM4
- `mb_chipset`: Z790, B760, H610, X670E, B650, A620
- `mb_form_factor`: ATX, Micro-ATX, Mini-ITX, E-ATX
- `mb_memory_type`: DDR4, DDR5
- `mb_max_memory`: 64GB, 128GB, 192GB

---

### 8. PSU (6 атрибутов)

| Code | Name | Type | Variant | Required | Filterable |
|------|------|------|---------|----------|------------|
| `psu_wattage` | Wattage / Snaga | select | ✅ | ✅ | ✅ |
| `psu_efficiency` | Efficiency Rating / Ocena efikasnosti | select | ✅ | ✅ | ✅ |
| `psu_modular` | Modular Type / Tip modularnosti | select | ❌ | ❌ | ✅ |
| `psu_form_factor` | Form Factor / Forma | select | ❌ | ❌ | ✅ |
| `psu_pcie5` | PCIe 5.0 Ready | boolean | ❌ | ❌ | ✅ |
| `psu_fan_size` | Fan Size / Veličina ventilatora | select | ❌ | ❌ | ✅ |

**Опции для select:**
- `psu_wattage`: 500W, 650W, 750W, 850W, 1000W, 1200W+
- `psu_efficiency`: 80+ Bronze, 80+ Silver, 80+ Gold, 80+ Platinum, 80+ Titanium
- `psu_modular`: Non-Modular, Semi-Modular, Fully Modular
- `psu_form_factor`: ATX, SFX, SFX-L
- `psu_fan_size`: 120mm, 140mm

---

### 9. Корпуса (8 атрибутов)

| Code | Name | Type | Variant | Required | Filterable |
|------|------|------|---------|----------|------------|
| `case_form_factor` | Form Factor / Forma | select | ✅ | ✅ | ✅ |
| `case_max_gpu_length` | Max GPU Length / Maksimalna dužina GPU | number | ❌ | ❌ | ✅ |
| `case_max_cpu_height` | Max CPU Cooler Height / Maksimalna visina hladnjaka | number | ❌ | ❌ | ✅ |
| `case_fans_included` | Fans Included / Uključeni ventilatori | number | ❌ | ❌ | ✅ |
| `case_max_fans` | Max Fans Supported / Maksimalno ventilatora | number | ❌ | ❌ | ✅ |
| `case_tempered_glass` | Tempered Glass / Kaljeno staklo | boolean | ❌ | ❌ | ✅ |
| `case_rgb_fans` | RGB Fans / RGB ventilatori | boolean | ❌ | ❌ | ✅ |
| `case_dust_filters` | Dust Filters / Filteri za prašinu | boolean | ❌ | ❌ | ✅ |

**Опции для select:**
- `case_form_factor`: Full Tower, Mid Tower, Mini Tower, Mini-ITX

---

### 10. Охлаждение (6 атрибутов)

| Code | Name | Type | Variant | Required | Filterable |
|------|------|------|---------|----------|------------|
| `cooler_type` | Cooler Type / Tip hlađenja | select | ✅ | ✅ | ✅ |
| `cooler_radiator` | Radiator Size / Veličina radijatora | select | ❌ | ❌ | ✅ |
| `cooler_fan_size` | Fan Size / Veličina ventilatora | select | ❌ | ❌ | ✅ |
| `cooler_max_tdp` | Max TDP / Maksimalni TDP | select | ❌ | ❌ | ✅ |
| `cooler_rgb` | RGB Lighting / RGB osvetljenje | boolean | ❌ | ❌ | ✅ |
| `cooler_noise` | Noise Level (dB) / Nivo buke | number | ❌ | ❌ | ✅ |

**Опции для select:**
- `cooler_type`: Air Tower, Low Profile, AIO 120mm, AIO 240mm, AIO 280mm, AIO 360mm
- `cooler_radiator`: 120mm, 240mm, 280mm, 360mm, 420mm
- `cooler_fan_size`: 120mm, 140mm
- `cooler_max_tdp`: 150W, 200W, 250W, 300W+

---

## 🎯 Вариативные атрибуты (14 штук)

**Атрибуты с `purpose: 'variant'` или связи в `category_variant_attributes`:**

| Category | Attribute Code | Affects Price | Affects Stock |
|----------|----------------|---------------|---------------|
| Grafičke kartice | `gpu_vram` | ✅ | ✅ |
| Procesori | `cpu_cores` | ✅ | ✅ |
| RAM memorija | `ram_type` | ✅ | ✅ |
| RAM memorija | `ram_capacity` | ✅ | ✅ |
| SSD nakopitelji | `ssd_capacity` | ✅ | ✅ |
| SSD nakopitelji | `ssd_interface` | ✅ | ✅ |
| HDD nakopitelji | `hdd_capacity` | ✅ | ✅ |
| HDD nakopitelji | `hdd_rpm` | ✅ | ✅ |
| Matične ploče | `mb_chipset` | ✅ | ✅ |
| Matične ploče | `mb_form_factor` | ✅ | ✅ |
| Napajanja | `psu_wattage` | ✅ | ✅ |
| Napajanja | `psu_efficiency` | ✅ | ✅ |
| Kućišta | `case_form_factor` | ✅ | ✅ |
| Hlađenje | `cooler_type` | ✅ | ✅ |

**Это означает:** Продавец может создать несколько вариантов товара (например, RTX 4080 с 12GB и 16GB VRAM) с разными ценами и остатками.

---

## 📖 Связанные документы

- **Action Plan:** `/p/github.com/vondi-global/passport/category-audit-racunarske-komponente-action-plan.md`
- **Listings Service CLAUDE.md:** `/p/github.com/vondi-global/listings/CLAUDE.md`
- **Production Quick Reference:** `/p/github.com/vondi-global/passport/PRODUCTION_QUICK_REFERENCE.md`

---

## ✅ Итоговый чеклист

- [ ] ФАЗА 0 завершена (категории переструктурированы)
- [ ] Миграция 000022 применена (атрибуты созданы)
- [ ] Миграция 000023 применена (атрибуты привязаны к категориям)
- [ ] Проверка: атрибуты видны в БД
- [ ] Проверка: фильтры работают на frontend
- [ ] Тестовые листинги созданы для каждой подкатегории
- [ ] OpenSearch индекс обновлён (если используется)

---

**Дата создания:** 2026-01-26
**Автор:** Claude Sonnet 4.5
**Версия:** 1.0
