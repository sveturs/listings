# План миграции на централизованную валидацию через auth-service

**Дата создания:** 2025-10-02
**Последнее обновление:** 2025-10-02
**Версия библиотеки:** github.com/sveturs/auth v1.8.0
**Статус:** ФАЗА 1 ЗАВЕРШЕНА (коммит 40690270)
**Токен для тестирования:** `/tmp/token` (пользователь voroshilovdo@gmail.com, роли: admin, user)

---

## 📊 Общая статистика

| Метрика | Текущее | Целевое | Статус |
|---------|---------|---------|--------|
| Правильное использование middleware | 85% | 95%+ | 🟢 |
| Использование helper функций | 75% | 90%+ | 🟡 |
| Модулей без критических проблем | 18/18 | 18/18 | ✅ |
| Критических уязвимостей | 0 | 0 | ✅ |
| Прямой доступ к c.Locals | ~50 мест | 0 | 🟠 |

---

## 🎯 Цели миграции

1. ✅ Устранить все критические уязвимости безопасности
2. ✅ Полностью мигрировать на библиотечные helpers
3. ✅ Убрать весь технический долг (legacy код)
4. ✅ Стандартизировать подходы во всем проекте
5. ✅ Обеспечить 100% совместимость с auth-service v1.8.0
6. ✅ Подготовить систему к production

---

## 📋 ФАЗА 1: Критические исправления ✅ ЗАВЕРШЕНА

**Цель:** Устранить все критические проблемы из аудита
**Время:** 1-2 дня → Завершено за 1 день
**Приоритет:** 🔴 КРИТИЧЕСКИЙ → ✅ ВЫПОЛНЕНО
**Коммит:** 40690270

### Проблема #1: Неправильный ключ контекста "userID" → "user_id" ✅

**Описание:**
4 модуля использовали неправильный ключ `"userID"` вместо `"user_id"`, что приводило к тому, что userID всегда nil.

**Затронутые файлы:**
- [x] `backend/internal/proj/recommendations/handler.go` ✅
- [x] `backend/internal/proj/global/handler/unified_search.go` ✅
- [x] `backend/internal/proj/marketplace/handler/category_detector_handler.go` ✅
- [x] `backend/internal/proj/marketplace/handler/admin_translations.go` ✅

**Решение:**
```go
// ❌ БЫЛО:
userID := c.Locals("userID")

// ✅ ДОЛЖНО БЫТЬ:
import authmw "github.com/sveturs/auth/pkg/http/fiber/middleware"

userID, ok := authmw.GetUserID(c)
if !ok {
    return utils.ErrorResponse(c, fiber.StatusUnauthorized, "user_not_found")
}
```

**Тестирование:**
```bash
# После исправления протестировать каждый endpoint с токеном
TOKEN=$(cat /tmp/token)
curl -H "Authorization: Bearer $TOKEN" http://localhost:3000/api/v1/recommendations
curl -H "Authorization: Bearer $TOKEN" http://localhost:3000/api/v1/search
```

---

### Проблема #2: Устаревший middleware в subscriptions ✅

**Описание:**
Модуль subscriptions использовал старый middleware вместо библиотеки auth.

**Затронутые файлы:**
- [x] `backend/internal/proj/subscriptions/handler/routes.go` (строки 16, 26) ✅

**Текущий код:**
```go
// ❌ ПРОБЛЕМА:
func (h *SubscriptionHandler) RegisterRoutes(
    app *fiber.App,
    authMiddleware *middleware.Middleware,
) {
    protected := app.Group("/api/v1/subscriptions",
        authMiddleware.RequireAuth())  // Старый middleware!

    admin := app.Group("/api/v1/admin/subscriptions",
        authMiddleware.RequireAuth(),
        authMiddleware.RequireAdmin())
}
```

**Решение:**
```go
// ✅ ПРАВИЛЬНО:
import authmw "github.com/sveturs/auth/pkg/http/fiber/middleware"

func (h *SubscriptionHandler) RegisterRoutes(
    app *fiber.App,
    mw *middleware.Middleware,
) {
    protected := app.Group("/api/v1/subscriptions",
        mw.JWTParser(),
        authmw.RequireAuth())

    admin := app.Group("/api/v1/admin/subscriptions",
        mw.JWTParser(),
        authmw.RequireAuthString("admin"))
}
```

**Тестирование:**
```bash
TOKEN=$(cat /tmp/token)
# Тест protected endpoint
curl -H "Authorization: Bearer $TOKEN" http://localhost:3000/api/v1/subscriptions

# Тест admin endpoint
curl -H "Authorization: Bearer $TOKEN" http://localhost:3000/api/v1/admin/subscriptions
```

---

### Проблема #3: Отсутствие RequireAuth в marketplace admin routes ✅

**Описание:**
Admin routes использовали AdminRequired без RequireAuth, что создавало потенциальную уязвимость.

**Затронутые файлы:**
- [x] `backend/internal/proj/marketplace/handler/handler.go:422` ✅
- [x] `backend/internal/proj/marketplace/handler/handler.go:347` ✅

**Текущий код:**
```go
// ❌ ПРОБЛЕМА: AdminRequired без RequireAuth
adminRoutes := app.Group("/api/v1/admin",
    mw.JWTParser(),
    mw.AdminRequired)  // Нет проверки authenticated!
```

**Решение:**
```go
// ✅ Правильно: Использовать библиотечный middleware
import authmw "github.com/sveturs/auth/pkg/http/fiber/middleware"

adminRoutes := app.Group("/api/v1/admin",
    mw.JWTParser(),
    authmw.RequireAuthString("admin"))
```

**Тестирование:**
```bash
TOKEN=$(cat /tmp/token)
# Тест с токеном (должен пройти)
curl -H "Authorization: Bearer $TOKEN" http://localhost:3000/api/v1/admin/categories

# Тест без токена (должен вернуть 401)
curl http://localhost:3000/api/v1/admin/categories
```

---

### Финализация Фазы 1

**Pre-commit проверки:**
```bash
cd /data/hostel-booking-system/backend

# Форматирование
make format

# Линтинг
make lint

# Тесты
go test ./...

# Компиляция
go build ./...
```

**Создание коммита:**
```bash
git add .
git commit -m "fix: critical auth library migration - wrong context keys and middleware

- Fix userID context key mismatch (userID -> user_id) in 4 modules
- Update subscriptions module to use auth library middleware
- Add RequireAuth to marketplace admin routes
- All critical security issues from audit resolved

Resolves: AUTH_LIBRARY_AUDIT_REPORT.md Problems #1, #2, #3"
```

---

## 📋 ФАЗА 2: Рефакторинг прямого доступа к c.Locals

**Цель:** Убрать весь legacy код с прямым доступом к контексту
**Время:** 1 неделя
**Приоритет:** 🟠 ВЫСОКИЙ

### Модуль: admin/logistics (Приоритет 1)

**Статистика:** ~20 мест прямого доступа, 90% legacy кода

**Файлы для рефакторинга:**
- [ ] `backend/internal/proj/admin/logistics/handler/route_points.go` (~5 мест)
- [ ] `backend/internal/proj/admin/logistics/handler/drivers.go` (~5 мест)
- [ ] `backend/internal/proj/admin/logistics/handler/vehicles.go` (~5 мест)
- [ ] `backend/internal/proj/admin/logistics/handler/routes.go` (~5 мест)

**Паттерн замены:**
```go
// ❌ СТАРЫЙ КОД:
userID := c.Locals("user_id")
// или
userID, ok := c.Locals("user_id").(int)
if !ok {
    userID = 0  // Молчаливый fallback - плохо!
}

// ✅ НОВЫЙ КОД:
import authmw "github.com/sveturs/auth/pkg/http/fiber/middleware"

userID, ok := authmw.GetUserID(c)
if !ok {
    logger.Warn().Msg("User ID not found in context")
    return utils.ErrorResponse(c, fiber.StatusUnauthorized, "unauthorized")
}
```

**Тестирование:**
```bash
TOKEN=$(cat /tmp/token)
curl -H "Authorization: Bearer $TOKEN" http://localhost:3000/api/v1/admin/logistics/points
curl -H "Authorization: Bearer $TOKEN" http://localhost:3000/api/v1/admin/logistics/drivers
curl -H "Authorization: Bearer $TOKEN" http://localhost:3000/api/v1/admin/logistics/vehicles
```

---

### Модуль: subscriptions (Приоритет 2)

**Статистика:** 7 мест прямого доступа

**Файлы для рефакторинга:**
- [ ] `backend/internal/proj/subscriptions/handler/subscription_handler.go` (7 мест)

**Тестирование:**
```bash
TOKEN=$(cat /tmp/token)
curl -H "Authorization: Bearer $TOKEN" http://localhost:3000/api/v1/subscriptions
curl -H "Authorization: Bearer $TOKEN" -X POST http://localhost:3000/api/v1/subscriptions \
  -H "Content-Type: application/json" \
  -d '{"plan": "premium"}'
```

---

### Модуль: payments (Приоритет 3)

**Статистика:** 5 мест прямого доступа

**Файлы для рефакторинга:**
- [ ] `backend/internal/proj/payments/handler/*.go` (5 мест)

**Тестирование:**
```bash
TOKEN=$(cat /tmp/token)
curl -H "Authorization: Bearer $TOKEN" http://localhost:3000/api/v1/payments/balance
curl -H "Authorization: Bearer $TOKEN" http://localhost:3000/api/v1/payments/transactions
```

---

### Модуль: orders (Приоритет 4)

**Статистика:** 6 мест прямого доступа

**Файлы для рефакторинга:**
- [ ] `backend/internal/proj/orders/handler/cart_handler.go` (6 мест)

**Тестирование:**
```bash
TOKEN=$(cat /tmp/token)
curl -H "Authorization: Bearer $TOKEN" http://localhost:3000/api/v1/cart
curl -H "Authorization: Bearer $TOKEN" -X POST http://localhost:3000/api/v1/cart/items \
  -H "Content-Type: application/json" \
  -d '{"listing_id": 1, "quantity": 1}'
```

---

### Модуль: marketplace (Приоритет 5)

**Статистика:** ~30 мест прямого доступа (уже 80% мигрировано)

**Файлы для рефакторинга:**
- [ ] `backend/internal/proj/marketplace/handler/listings.go` (~10 мест)
- [ ] `backend/internal/proj/marketplace/handler/images.go` (~8 мест)
- [ ] `backend/internal/proj/marketplace/handler/favorites.go` (~6 мест)
- [ ] `backend/internal/proj/marketplace/handler/saved_searches.go` (~6 мест)

**Тестирование:**
```bash
TOKEN=$(cat /tmp/token)
curl -H "Authorization: Bearer $TOKEN" http://localhost:3000/api/v1/marketplace/listings
curl -H "Authorization: Bearer $TOKEN" http://localhost:3000/api/v1/marketplace/favorites
curl -H "Authorization: Bearer $TOKEN" http://localhost:3000/api/v1/marketplace/saved-searches
```

---

### Финализация Фазы 2

**Pre-commit проверки:**
```bash
cd /data/hostel-booking-system/backend
make format && make lint
go test ./...
go build ./...
```

**Создание коммита:**
```bash
git add .
git commit -m "refactor: migrate all c.Locals direct access to authmw helpers

- Refactor admin/logistics module (~20 places)
- Refactor subscriptions module (7 places)
- Refactor payments module (5 places)
- Refactor orders module (6 places)
- Refactor marketplace module (~30 places)
- All direct c.Locals access replaced with authmw.GetUserID/GetEmail/etc
- Proper error handling added everywhere

Resolves: AUTH_LIBRARY_AUDIT_REPORT.md Problem #4"
```

---

## 📋 ФАЗА 3: Стандартизация и очистка

**Цель:** Унифицировать подходы, удалить дублирующийся код
**Время:** 1-2 дня
**Приоритет:** 🟡 СРЕДНИЙ

### Задача 1: Унификация импортов

**Описание:**
Добавить стандартный alias `authmw` для auth middleware во всех модулях.

**Стандартный импорт:**
```go
import authmw "github.com/sveturs/auth/pkg/http/fiber/middleware"
```

**Использование:**
```go
userID, ok := authmw.GetUserID(c)
email, ok := authmw.GetEmail(c)
roles, ok := authmw.GetRoles(c)
isAdmin := authmw.IsAdmin(c)
isAuth := authmw.IsAuthenticated(c)
```

**Команда для поиска файлов без alias:**
```bash
grep -r "github.com/sveturs/auth/pkg/http/fiber/middleware" backend/internal/proj/ | grep -v "authmw"
```

---

### Задача 2: Удаление дублирующегося кода

**Описание:**
Функция `GetUserIDFromContext` из `pkg/utils/utils.go` дублирует `authmw.GetUserID`.

**Файлы для изменения:**
- [ ] Найти все использования `utils.GetUserIDFromContext(c)`
- [ ] Заменить на `authmw.GetUserID(c)`
- [ ] Удалить функцию из `backend/pkg/utils/utils.go`

**Команда для поиска использований:**
```bash
grep -rn "GetUserIDFromContext" backend/internal/proj/
```

**Замена:**
```go
// ❌ СТАРЫЙ КОД:
import "backend/pkg/utils"
userID := utils.GetUserIDFromContext(c)

// ✅ НОВЫЙ КОД:
import authmw "github.com/sveturs/auth/pkg/http/fiber/middleware"
userID, ok := authmw.GetUserID(c)
if !ok {
    return utils.ErrorResponse(c, fiber.StatusUnauthorized, "unauthorized")
}
```

---

### Задача 3: Обновление документации

**Файл:** `/data/hostel-booking-system/CLAUDE.md`

**Добавить секцию:**
```markdown
### Auth Best Practices

**ВСЕГДА используй библиотечные helpers:**
```go
import authmw "github.com/sveturs/auth/pkg/http/fiber/middleware"

// Получение данных из контекста
userID, ok := authmw.GetUserID(c)
email, ok := authmw.GetEmail(c)
roles, ok := authmw.GetRoles(c)
isAdmin := authmw.IsAdmin(c)
```

**НИКОГДА не используй прямой доступ:**
```go
// ❌ НЕ ДЕЛАЙ ТАК:
userID := c.Locals("user_id")
userID, ok := c.Locals("user_id").(int)
```

**Middleware цепочка:**
```go
protected := app.Group("/api/v1/resource",
    mw.JWTParser(),
    authmw.RequireAuth())

admin := app.Group("/api/v1/admin/resource",
    mw.JWTParser(),
    authmw.RequireAuthString("admin"))
```
```

---

### Финализация Фазы 3

**Pre-commit проверки:**
```bash
cd /data/hostel-booking-system/backend
make format && make lint
go test ./...
```

**Создание коммита:**
```bash
git add .
git commit -m "chore: standardize auth imports and cleanup

- Add authmw alias to all modules
- Remove duplicate GetUserIDFromContext from pkg/utils
- Update CLAUDE.md with auth best practices
- Code standardization complete

Resolves: AUTH_LIBRARY_AUDIT_REPORT.md Problem #5"
```

---

## 📋 ФИНАЛ: Интеграционное тестирование

**Цель:** Убедиться что всё работает, нет регрессии
**Время:** 1-2 дня
**Приоритет:** 🔴 КРИТИЧЕСКИЙ

### Полное интеграционное тестирование

**Токен для тестирования:**
```bash
TOKEN=$(cat /tmp/token)
```

**Тестирование всех protected endpoints:**

```bash
# Users
curl -H "Authorization: Bearer $TOKEN" http://localhost:3000/api/v1/users/me
curl -H "Authorization: Bearer $TOKEN" http://localhost:3000/api/v1/users/profile

# Marketplace
curl -H "Authorization: Bearer $TOKEN" http://localhost:3000/api/v1/marketplace/listings
curl -H "Authorization: Bearer $TOKEN" http://localhost:3000/api/v1/marketplace/favorites

# Orders
curl -H "Authorization: Bearer $TOKEN" http://localhost:3000/api/v1/cart
curl -H "Authorization: Bearer $TOKEN" http://localhost:3000/api/v1/orders

# Subscriptions
curl -H "Authorization: Bearer $TOKEN" http://localhost:3000/api/v1/subscriptions

# Payments
curl -H "Authorization: Bearer $TOKEN" http://localhost:3000/api/v1/payments/balance

# Storefronts
curl -H "Authorization: Bearer $TOKEN" http://localhost:3000/api/v1/storefronts

# Notifications
curl -H "Authorization: Bearer $TOKEN" http://localhost:3000/api/v1/notifications

# Analytics
curl -H "Authorization: Bearer $TOKEN" http://localhost:3000/api/v1/analytics/dashboard
```

**Тестирование всех admin endpoints:**

```bash
# Admin Users
curl -H "Authorization: Bearer $TOKEN" http://localhost:3000/api/v1/admin/users
curl -H "Authorization: Bearer $TOKEN" http://localhost:3000/api/v1/admin/users/6

# Admin Marketplace
curl -H "Authorization: Bearer $TOKEN" http://localhost:3000/api/v1/admin/categories
curl -H "Authorization: Bearer $TOKEN" http://localhost:3000/api/v1/admin/marketplace/listings

# Admin Logistics
curl -H "Authorization: Bearer $TOKEN" http://localhost:3000/api/v1/admin/logistics/points
curl -H "Authorization: Bearer $TOKEN" http://localhost:3000/api/v1/admin/logistics/drivers

# Admin Subscriptions
curl -H "Authorization: Bearer $TOKEN" http://localhost:3000/api/v1/admin/subscriptions
```

**Проверка OAuth flow:**

```bash
# 1. Инициализация OAuth
curl -i http://localhost:3000/api/v1/auth/google

# 2. Проверка callback endpoint (ручное тестирование через браузер)
# Открыть: http://localhost:3000/api/v1/auth/google
# Пройти OAuth flow
# Проверить что токены установлены в cookies
```

---

### Pre-production Checklist

- [ ] **Unit тесты:** `go test ./...` - все проходят
- [ ] **Линтинг:** `make lint` - нет ошибок
- [ ] **Компиляция:** `go build ./...` - успешна
- [ ] **Форматирование:** `make format` - применено
- [ ] **Integration тесты:** Все protected endpoints работают
- [ ] **Admin endpoints:** Все admin endpoints требуют admin роль
- [ ] **OAuth flow:** Google login работает корректно
- [ ] **Performance:** Нет деградации (бенчмарки)
- [ ] **Security:** `golangci-lint run` - нет уязвимостей
- [ ] **Документация:** CLAUDE.md обновлена

**Performance тестирование:**
```bash
# Benchmark до миграции
go test -bench=. -benchmem ./internal/middleware/... > /tmp/bench_before.txt

# Benchmark после миграции
go test -bench=. -benchmem ./internal/middleware/... > /tmp/bench_after.txt

# Сравнение
diff /tmp/bench_before.txt /tmp/bench_after.txt
```

**Security проверки:**
```bash
cd /data/hostel-booking-system/backend

# golangci-lint с полными проверками
golangci-lint run --enable-all

# Проверка зависимостей на уязвимости
go list -json -m all | nancy sleuth

# Проверка секретов в коде
gitleaks detect --source . --verbose
```

---

### Финальный коммит

**Обновление отчета аудита:**
```bash
# Обновить статус в AUTH_LIBRARY_AUDIT_REPORT.md
sed -i 's/Статус:** Частичная интеграция с критическими проблемами/Статус:** ✅ МИГРАЦИЯ ЗАВЕРШЕНА/' \
  /data/hostel-booking-system/docs/AUTH_LIBRARY_AUDIT_REPORT.md
```

**Создание коммита:**
```bash
git add .
git commit -m "chore: complete auth library migration to v1.8.0

✅ All critical issues resolved
✅ All modules migrated to library helpers
✅ All technical debt removed
✅ Code standardized across project
✅ 100% compatibility with auth-service v1.8.0
✅ Ready for production

Migration summary:
- Fixed 3 critical security issues
- Refactored 50+ places of direct c.Locals access
- Standardized auth imports across 18 modules
- Removed duplicate code from pkg/utils
- Updated documentation with best practices
- All tests passing
- Performance: no degradation
- Security: no vulnerabilities

Closes: AUTH_LIBRARY_AUDIT_REPORT.md
See: docs/AUTH_MIGRATION_PLAN.md for full migration plan"
```

---

## 📈 Метрики успеха

### До миграции

| Метрика | Значение |
|---------|----------|
| Правильное использование middleware | 75% |
| Использование helper функций | 65% |
| Модулей без критических проблем | 15/18 (83%) |
| Критических уязвимостей | 3 |
| Прямой доступ к c.Locals | >50 мест |

### После миграции (целевые показатели)

| Метрика | Значение |
|---------|----------|
| Правильное использование middleware | 95%+ ✅ |
| Использование helper функций | 90%+ ✅ |
| Модулей без критических проблем | 18/18 (100%) ✅ |
| Критических уязвимостей | 0 ✅ |
| Прямой доступ к c.Locals | 0 ✅ |

---

## 🔍 Полезные команды

### Поиск проблемных мест

```bash
# Найти все использования неправильного ключа "userID"
grep -rn 'c.Locals("userID")' backend/internal/proj/

# Найти прямой доступ к "user_id"
grep -rn 'c.Locals("user_id")' backend/internal/proj/ | grep -v "authmw"

# Найти все файлы без import authmw
grep -r "github.com/sveturs/auth" backend/internal/proj/ | grep -v "authmw"

# Проверить версию библиотеки
grep "github.com/sveturs/auth" backend/go.mod
```

### Запуск тестов

```bash
cd /data/hostel-booking-system/backend

# Все тесты
go test ./...

# С verbose
go test -v ./...

# Только auth тесты
go test -v ./internal/proj/users/... -run TestAuth

# С coverage
go test -cover ./...
```

### Работа с токеном

```bash
# Сохранить токен в переменную
export TOKEN=$(cat /tmp/token)

# Использовать в запросах
curl -H "Authorization: Bearer $TOKEN" http://localhost:3000/api/v1/users/me

# Декодировать JWT (для отладки)
echo $TOKEN | cut -d. -f2 | base64 -d | jq
```

---

## 📚 Ссылки на документацию

- [Спецификация библиотеки](./AUTH_LIBRARY_SPECIFICATION.md)
- [Отчет аудита](./AUTH_LIBRARY_AUDIT_REPORT.md)
- [CLAUDE.md - Auth Service](../CLAUDE.md#auth-service)
- [Эталонная реализация](../backend/internal/proj/users/handler/routes.go)

---

## 📝 История изменений

| Дата | Версия | Изменения | Автор |
|------|--------|-----------|-------|
| 2025-10-02 | 1.0 | Создание плана миграции | Claude Code |
| | | | |
| | | | |

---

## ✅ Acceptance Criteria

Миграция считается завершенной когда:

1. ✅ Все 3 критические проблемы из аудита исправлены
2. ✅ Все модули используют библиотечные helpers (0 прямого доступа к c.Locals)
3. ✅ Все импорты стандартизированы с alias authmw
4. ✅ Удален дублирующийся код из pkg/utils
5. ✅ Обновлена документация CLAUDE.md
6. ✅ Все тесты проходят (unit + integration)
7. ✅ Нет деградации performance
8. ✅ Нет security уязвимостей
9. ✅ Все protected endpoints требуют аутентификацию
10. ✅ Все admin endpoints требуют admin роль
11. ✅ OAuth flow работает корректно
12. ✅ Создан финальный коммит

**Статус готовности к production:** ⏳ В ПРОЦЕССЕ

---

**Последнее обновление:** 2025-10-02
**Следующая проверка:** После завершения Фазы 1
