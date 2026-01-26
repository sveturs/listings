# ФАЗА 0: Переструктурирование категории "Računarske komponente"

**Дата:** 2026-01-26
**Статус:** ✅ Готово к применению
**Миграции:** 000025 (категории) + 000023 (атрибуты)

---

## 📋 Что делают миграции

### Миграция 000025: Переструктурирование категорий

**Удаляет 20 конкретных подкатегорий:**
- graficke-rtx-4000, graficke-rtx-3000, graficke-gtx-1000
- graficke-amd-rx-7000, graficke-amd-rx-6000
- procesor-intel-i9, procesor-intel-i7, procesor-intel-i5
- procesor-amd-ryzen-9, procesor-amd-ryzen-7, procesor-amd-ryzen-5
- ram-ddr5-32gb, ram-ddr5-16gb, ram-ddr4-16gb, ram-ddr4-8gb
- ssd-nvme-1tb, ssd-nvme-500gb, ssd-sata-1tb
- maticna-ploca-intel-z790, maticna-ploca-amd-x670

**Создаёт 9 общих подкатегорий:**
1. **graficke-kartice** - Grafičke kartice (Видеокарты) 🎮
2. **procesori** - Procesori (Процессоры) ⚙️
3. **ram-memorija** - RAM memorija (Оперативная память) 🧠
4. **ssd-nakopitelji** - SSD nakopitelji (SSD накопители) 💿
5. **hdd-nakopitelji** - HDD nakopitelji (HDD накопители) 💾
6. **maticne-ploce** - Matične ploče (Материнские платы) 🔌
7. **napajanja** - Napajanja (PSU) (Блоки питания) ⚡
8. **kucista** - Kućišta (Корпуса) 📦
9. **hladjenje** - Hlađenje (Охлаждение) ❄️

**Для каждой категории:**
- ✅ Полные переводы (en/ru/sr)
- ✅ SEO метаданные (title, description, keywords)
- ✅ Icon (emoji)
- ✅ Path и hierarchy

---

### Миграция 000023: Привязка атрибутов к категориям

**Привязывает 77 атрибутов к 9 категориям:**

**Общие атрибуты (для всех 9 категорий):**
- pc_brand, pc_model, pc_condition, pc_warranty, pc_color, pc_rgb_lighting

**Специфичные атрибуты:**
- graficke-kartice: 10 атрибутов (gpu_series, gpu_vram, gpu_chip, etc.)
- procesori: 9 атрибутов (cpu_series, cpu_cores, cpu_socket, etc.)
- ram-memorija: 8 атрибутов (ram_type, ram_capacity, ram_speed, etc.)
- ssd-nakopitelji: 8 атрибутов (ssd_capacity, ssd_interface, etc.)
- hdd-nakopitelji: 6 атрибутов (hdd_capacity, hdd_rpm, etc.)
- maticne-ploce: 10 атрибутов (mb_socket, mb_chipset, etc.)
- napajanja: 6 атрибутов (psu_wattage, psu_efficiency, etc.)
- kucista: 8 атрибутов (case_form_factor, etc.)
- hladjenje: 6 атрибутов (cooler_type, etc.)

**Результат:**
- ~200 связей category_attributes
- ~40 связей category_variant_attributes (вариативные атрибуты для создания вариантов)

---

## 🚀 Как применить (автоматически через PR)

### ⚠️ ВАЖНО: НЕ применять миграции вручную!

**Правильный способ - через PR:**

```bash
# 1. Создать feature ветку
cd /p/github.com/vondi-global/listings
git checkout main && git pull vondi main
git checkout -b feature/phase0-categories-restructure

# 2. Добавить миграции (уже созданы локально)
git add migrations/000025_*.sql
git add migrations/000023_*.sql
git add migrations/README_PHASE0_CATEGORIES.md

# 3. Обновить CHANGELOG.md
echo "### Added - 2026-01-26 (COMMIT_HASH)

**ФАЗА 0: Переструктурирование категорий Računarske komponente**

#### Миграция 000025
- Удалено 20 конкретных подкатегорий
- Создано 9 общих подкатегорий

#### Миграция 000023
- Привязано 77 атрибутов к 9 категориям
- ~200 связей category_attributes
- ~40 вариативных атрибутов

#### Файлы
- migrations/000025_restructure_racunarske_komponente_categories.{up,down}.sql
- migrations/000023_link_racunarske_komponente_attributes_to_categories.{up,down}.sql
" >> CHANGELOG.md

# 4. Commit + Push
git commit -m "feat: ФАЗА 0 - переструктурирование категорий (9 общих вместо 20)"
git push vondi feature/phase0-categories-restructure

# 5. Создать PR
gh pr create --repo vondi-global/listings \
  --base main \
  --head feature/phase0-categories-restructure \
  --title "feat: ФАЗА 0 - категории Računarske komponente (9 общих вместо 20)"

# 6. Дождаться CI checks
gh pr checks --watch

# 7. Merge PR
gh pr merge --merge
```

---

## 🔄 Автоматическое применение миграций

**После merge PR в main:**

```
1. Deploy workflow запускается автоматически
   ↓
2. Build job: собирает Docker image
   ↓
3. Migrate job (АВТОМАТИЧЕСКИ):
   - kubectl cp migrations/ в postgres pod
   - kubectl exec: /tmp/migrate -path /tmp/migrations -database '...' up
   - Применяет 000025 → удаляет 20 категорий, создаёт 9 новых
   - Применяет 000023 → привязывает атрибуты
   ↓
4. Deploy job: обновляет k8s-configs (GitOps)
   ↓
5. ArgoCD синхронизирует изменения
   ↓
6. Listings Service перезапускается с новыми категориями
```

**Workflow:** `.github/workflows/deploy-production.yml` (строки 77-169)

**Время deployment:** ~5 минут (Build + Migrate + Deploy)

---

## ✅ Проверка после deployment

**1. Проверить категории в production:**

```bash
ssh vondi "kubectl exec -n production postgres-listings-0 -- psql -U listings_user -d listings_db -c \"
SELECT slug, name->>'sr' as name, icon
FROM categories
WHERE parent_id = 'fe64130f-deea-4767-a937-6f9d584a4395'::uuid
ORDER BY sort_order;
\""
```

**Ожидаемый результат:**
```
      slug       |       name       | icon
-----------------+------------------+------
 graficke-kartice | Grafičke kartice | 🎮
 procesori        | Procesori        | ⚙️
 ram-memorija     | RAM memorija     | 🧠
 ssd-nakopitelji  | SSD nakopitelji  | 💿
 hdd-nakopitelji  | HDD nakopitelji  | 💾
 maticne-ploce    | Matične ploče    | 🔌
 napajanja        | Napajanja (PSU)  | ⚡
 kucista          | Kućišta          | 📦
 hladjenje        | Hlađenje         | ❄️
(9 rows)
```

---

**2. Проверить привязку атрибутов:**

```bash
ssh vondi "kubectl exec -n production postgres-listings-0 -- psql -U listings_user -d listings_db -c \"
SELECT
  c.slug,
  COUNT(ca.id) as attributes_count,
  COUNT(cva.id) as variant_attrs_count
FROM categories c
LEFT JOIN category_attributes ca ON ca.category_id = c.id
LEFT JOIN category_variant_attributes cva ON cva.category_id::uuid = c.id
WHERE c.slug IN ('graficke-kartice', 'procesori', 'ram-memorija')
GROUP BY c.id, c.slug
ORDER BY c.slug;
\""
```

**Ожидаемый результат:**
```
      slug       | attributes_count | variant_attrs_count
-----------------+------------------+--------------------
 graficke-kartice |       16         |          1
 procesori        |       15         |          1
 ram-memorija     |       14         |          2
```

---

**3. Проверить на frontend:**

Открыть: https://vondi.rs/sr/categories/racunarske-komponente

**Должны отображаться 9 подкатегорий вместо 20.**

---

## 📊 Статистика миграций

| Миграция | Операция | Записей | Время |
|----------|----------|---------|-------|
| **000025** | DELETE categories | 20 | <1s |
| **000025** | INSERT categories | 9 | <1s |
| **000023** | INSERT category_attributes | ~200 | ~2s |
| **000023** | INSERT category_variant_attributes | ~40 | <1s |
| **ИТОГО** | - | ~270 | **~5s** |

---

## 🎯 Преимущества новой структуры

**До (НЕПРАВИЛЬНО):**
```
Grafičke RTX 4000 serija
  └── RTX 4090 24GB (1 товар)
Grafičke RTX 3000 serija
  └── RTX 3090 24GB (1 товар)
```
❌ 2 категории, 2 товара (должен быть 1 товар с вариантами!)

**После (ПРАВИЛЬНО):**
```
Grafičke kartice (ОБЩАЯ категория)
  └── NVIDIA GeForce RTX 4090 (1 товар)
      ├── Вариант: 12GB GDDR6X - 1299€
      ├── Вариант: 16GB GDDR6X - 1499€
      └── Вариант: 24GB GDDR6X - 1699€

      Фильтры:
      - GPU Series: RTX 4000, RTX 3000, RX 7000, etc.
      - GPU Chip: RTX 4090, RTX 4080, RTX 4070, etc.
      - VRAM: 6GB, 8GB, 12GB, 16GB, 24GB
```
✅ 1 категория, 1 товар, 3 варианта (правильная архитектура!)

---

**Дата создания:** 2026-01-26
**Автор:** Claude Sonnet 4.5
**Статус:** Ready for PR
