# 🧹 Legacy Code Cleanup Plan
## Удаление устаревшего кода системы атрибутов

*Дата создания: 04.09.2025*
*Период выполнения: Дни 16-20 проекта унификации*
*Приоритет: HIGH*
*Риск: MEDIUM*

---

## 📋 Executive Summary

После успешного deployment unified attributes системы, необходимо провести планомерное удаление legacy кода для упрощения поддержки и уменьшения технического долга.

### Объем работ:
- **14 таблиц БД** для удаления/архивации
- **~8,500 строк** backend кода
- **~2,600 строк** frontend кода
- **3 параллельные системы** атрибутов

---

## ⚠️ Критические правила безопасности

### ОБЯЗАТЕЛЬНО:
1. ✅ Сохранить полный backup перед удалением
2. ✅ Проверить отсутствие зависимостей
3. ✅ Тестировать после каждого этапа
4. ✅ Мониторить ошибки в production
5. ✅ Документировать все изменения

### ЗАПРЕЩЕНО:
1. ❌ Удалять код без проверки использования
2. ❌ Удалять данные без архивации
3. ❌ Вносить изменения в пятницу
4. ❌ Удалять более 1000 строк за раз

---

## 📊 Инвентаризация Legacy компонентов

### Database Tables (PostgreSQL)

#### Для архивации (уже не используются):
```sql
-- Marketplace система (старая)
category_attributes          -- 85 записей
listing_attributes           -- 0 записей
category_attribute_values    -- 15 записей

-- Admin панель (дубликат)
admin_category_attributes    -- 0 записей
admin_attribute_values       -- 0 записей
admin_listing_attributes     -- 0 записей

-- Автомобильная система
automotive_makes             -- 0 записей
automotive_models            -- 0 записей
automotive_attributes        -- 0 записей
automotive_listings          -- 0 записей
vehicle_attributes           -- 0 записей
vehicle_types               -- 0 записей
vehicle_features            -- 0 записей
vehicle_conditions          -- 0 записей
```

#### Индексы для удаления:
```sql
-- 17 неиспользуемых индексов
idx_automotive_*
idx_vehicle_*
idx_admin_attributes_*
```

### Backend Code (Go)

#### Полностью удалить:
```
backend/internal/proj/admin/attributes/    -- ~2,500 строк
backend/internal/proj/automotive/          -- ~3,000 строк
backend/internal/storage/postgres/
  - attributes_old.go                       -- ~800 строк
  - automotive.go                           -- ~1,200 строк
```

#### Частично очистить:
```
backend/internal/proj/marketplace/handler/
  - attributes.go (старые endpoints)        -- ~500 строк
backend/internal/domain/models/
  - attributes_legacy.go                    -- ~300 строк
```

### Frontend Code (React/TypeScript)

#### Компоненты для удаления:
```
frontend/svetu/src/components/
  - AttributeSelector_OLD.tsx               -- ~450 строк
  - CategoryAttributes_BACKUP.tsx          -- ~380 строк
  - admin/AttributeManager.tsx             -- ~620 строк
  - automotive/VehicleAttributes.tsx       -- ~890 строк
```

#### Сервисы для удаления:
```
frontend/svetu/src/services/
  - attributeService_old.ts                -- ~320 строк
  - automotiveService.ts                   -- ~450 строк
```

---

## 📅 План выполнения по дням

### День 16: Подготовка и анализ
**Цель:** Полная инвентаризация и подготовка

**Задачи:**
1. ✅ Создать полный backup БД
2. ✅ Проанализировать зависимости через grep/ast
3. ✅ Создать скрипты для поиска использования
4. ✅ Подготовить rollback план

**Скрипты:**
```bash
# Поиск использования в коде
grep -r "category_attributes" --exclude-dir=node_modules
grep -r "automotive_" --exclude-dir=vendor
grep -r "AttributeSelector_OLD" --exclude="*.log"
```

### День 17: Архивация БД
**Цель:** Безопасная архивация неиспользуемых таблиц

**SQL миграция:**
```sql
-- 000036_archive_legacy_tables.up.sql
BEGIN;

-- Создать архивную схему
CREATE SCHEMA IF NOT EXISTS archive_legacy;

-- Переместить таблицы
ALTER TABLE category_attributes SET SCHEMA archive_legacy;
ALTER TABLE listing_attributes SET SCHEMA archive_legacy;
ALTER TABLE category_attribute_values SET SCHEMA archive_legacy;

-- Автомобильные таблицы
ALTER TABLE automotive_makes SET SCHEMA archive_legacy;
ALTER TABLE automotive_models SET SCHEMA archive_legacy;
-- ... остальные таблицы

-- Добавить метаданные
COMMENT ON SCHEMA archive_legacy IS 'Archived legacy attribute tables - Day 17 unified attributes project';

COMMIT;
```

**Валидация:**
- Проверить работу приложения
- Убедиться в отсутствии ошибок
- Мониторить логи 2 часа

### День 18: Очистка Backend
**Цель:** Удаление неиспользуемого Go кода

**Этапы:**
1. Удалить automotive модуль:
```bash
rm -rf backend/internal/proj/automotive/
rm backend/internal/storage/postgres/automotive.go
```

2. Удалить admin attributes:
```bash
rm -rf backend/internal/proj/admin/attributes/
```

3. Очистить старые модели:
```bash
rm backend/internal/domain/models/attributes_legacy.go
```

4. Обновить imports и зависимости:
```bash
go mod tidy
go test ./...
```

**Проверки:**
- `go build` успешно
- Все тесты проходят
- API endpoints работают

### День 19: Очистка Frontend
**Цель:** Удаление legacy React компонентов

**Этапы:**
1. Удалить старые компоненты:
```bash
rm frontend/svetu/src/components/AttributeSelector_OLD.tsx
rm frontend/svetu/src/components/CategoryAttributes_BACKUP.tsx
rm -rf frontend/svetu/src/components/automotive/
```

2. Удалить неиспользуемые сервисы:
```bash
rm frontend/svetu/src/services/attributeService_old.ts
rm frontend/svetu/src/services/automotiveService.ts
```

3. Обновить импорты:
```bash
npm run lint:fix
npm run build
```

**Проверки:**
- Build успешно завершается
- Нет broken imports
- UI функционирует корректно

### День 20: Финализация и документация
**Цель:** Завершение cleanup и обновление документации

**Задачи:**
1. Удалить неиспользуемые индексы:
```sql
DROP INDEX IF EXISTS idx_automotive_makes_name;
DROP INDEX IF EXISTS idx_vehicle_attributes_listing_id;
-- ... остальные индексы
```

2. Очистить конфигурационные файлы:
- Удалить упоминания legacy систем из .env.example
- Обновить docker-compose.yml
- Очистить nginx конфигурацию

3. Обновить документацию:
- README.md - удалить упоминания старых систем
- API документация - убрать deprecated endpoints
- Архитектурная документация

4. Финальные проверки:
- Full regression testing
- Performance benchmarks
- Security scan

---

## 🔍 Инструменты для анализа зависимостей

### Backend (Go)
```bash
# Найти все импорты модуля
go list -f '{{.ImportPath}} {{.Imports}}' ./... | grep automotive

# Проверить использование структур
grep -r "CategoryAttribute" --include="*.go" | grep -v unified

# AST анализ
go vet ./...
staticcheck ./...
```

### Frontend (TypeScript)
```bash
# Найти импорты компонента
grep -r "AttributeSelector_OLD" --include="*.tsx" --include="*.ts"

# Проверить зависимости
npm ls | grep automotive

# TypeScript проверка
npx tsc --noEmit
```

### Database
```sql
-- Проверить foreign keys
SELECT 
    tc.table_name, 
    kcu.column_name, 
    ccu.table_name AS referenced_table
FROM information_schema.table_constraints AS tc 
JOIN information_schema.key_column_usage AS kcu
    ON tc.constraint_name = kcu.constraint_name
JOIN information_schema.constraint_column_usage AS ccu
    ON ccu.constraint_name = tc.constraint_name
WHERE tc.constraint_type = 'FOREIGN KEY' 
    AND ccu.table_name IN ('category_attributes', 'automotive_makes');
```

---

## ⚠️ Риски и митигация

### Risk Matrix

| Риск | Вероятность | Влияние | Митигация |
|------|-------------|---------|-----------|
| Broken dependencies | Medium | High | Thorough grep analysis |
| Data loss | Low | Critical | Full backup + archive |
| Performance degradation | Low | Medium | Benchmark before/after |
| Missing functionality | Medium | High | Feature flag fallback |
| User complaints | Low | Medium | Gradual rollout |

### Rollback Strategy

#### Level 1: Code rollback (5 минут)
```bash
git revert HEAD
git push origin main
kubectl rollout restart deployment/backend
```

#### Level 2: Database rollback (15 минут)
```sql
-- Восстановить таблицы из архива
ALTER TABLE archive_legacy.category_attributes SET SCHEMA public;
ALTER TABLE archive_legacy.listing_attributes SET SCHEMA public;
```

#### Level 3: Full restoration (30 минут)
```bash
# Восстановить из backup
pg_restore -d svetubd /backups/pre_cleanup_backup.dump

# Deploy предыдущую версию
kubectl set image deployment/backend backend=backend:v1.9.0
kubectl set image deployment/frontend frontend=frontend:v1.9.0
```

---

## 📊 Метрики успеха

### Технические метрики:
- ✅ Уменьшение codebase на ~11,000 строк
- ✅ Удаление 14 неиспользуемых таблиц
- ✅ Уменьшение размера БД на ~50MB
- ✅ Ускорение CI/CD на 20%
- ✅ Уменьшение памяти приложения на 10%

### Бизнес метрики:
- ✅ 0 инцидентов во время cleanup
- ✅ 0 жалоб пользователей
- ✅ Сохранение 100% функциональности
- ✅ Улучшение maintainability score

---

## 📝 Чеклист для каждого этапа

### Перед удалением:
- [ ] Backup создан и проверен
- [ ] Dependencies проанализированы
- [ ] Tests написаны/обновлены
- [ ] Team уведомлена
- [ ] Monitoring настроен

### После удаления:
- [ ] Build успешен
- [ ] Tests проходят
- [ ] No errors в логах
- [ ] Performance не деградировала
- [ ] Documentation обновлена

---

## 🏁 Финальная проверка (День 20)

### Automated checks:
```bash
# Backend
go build ./...
go test ./...
golangci-lint run
go mod verify

# Frontend
npm run build
npm run test
npm run lint
npm audit

# Database
psql -c "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public';"
psql -c "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'archive_legacy';"
```

### Manual verification:
1. Create new listing with attributes ✓
2. Search with attribute filters ✓
3. Edit existing listing attributes ✓
4. Admin panel functionality ✓
5. API response times ✓

---

## 📚 Documentation Updates Required

1. **README.md** - Remove legacy system mentions
2. **API.md** - Remove deprecated endpoints
3. **ARCHITECTURE.md** - Update system diagram
4. **DEPLOYMENT.md** - Update configuration
5. **CONTRIBUTING.md** - Update development setup

---

## 🎯 Expected Outcomes

### After completion:
- **Codebase:** -40% complexity
- **Maintenance:** -60% effort
- **Performance:** +10% speed
- **Developer Experience:** Significantly improved
- **Technical Debt:** Substantially reduced

---

## 📅 Timeline Summary

| День | Фаза | Риск | Время |
|------|------|------|-------|
| 16 | Preparation & Analysis | Low | 4h |
| 17 | Database Archival | Medium | 6h |
| 18 | Backend Cleanup | Medium | 8h |
| 19 | Frontend Cleanup | Medium | 6h |
| 20 | Finalization | Low | 4h |

**Total effort:** 28 hours

---

## ✅ Success Criteria

The cleanup is considered successful when:
1. All legacy code is removed or archived
2. Zero production incidents
3. All tests passing
4. Performance improved or stable
5. Documentation updated
6. Team satisfied with results

---

**Document Status:** READY FOR EXECUTION
**Version:** 1.0.0
**Author:** System Architect
**Next Review:** Day 16 before execution

---