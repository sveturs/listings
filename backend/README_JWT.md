# JWT Authentication and Fiber Middleware Guide

## Проблема: Middleware "утекает" на публичные маршруты

### Симптомы
- Публичные endpoints (например `/api/v1/storefronts`) возвращают 401 Unauthorized
- В консоли браузера флудят ошибки 401 при загрузке главной страницы
- Middleware `AuthRequiredJWT` вызывается для публичных маршрутов

### Причина
В **Fiber v2.52.6** создание группы с middleware и широким префиксом применяет этот middleware ко **ВСЕМ** маршрутам, начинающимся с этого префикса, даже если они зарегистрированы отдельно.

### ❌ Неправильный подход

```go
// ПЛОХО: Создание группы /api/v1 с middleware
authedAPIGroup := app.Group("/api/v1", mw.AuthRequiredJWT)

marketplaceProtected := authedAPIGroup.Group("/marketplace")
marketplaceProtected.Post("/listings", h.CreateListing)

// ⚠️ Этот код применит AuthRequiredJWT ко ВСЕМ маршрутам /api/v1/*
// Включая публичные /api/v1/storefronts, /api/v1/marketplace/search и т.д.!
```

### ✅ Правильный подход

```go
// ХОРОШО: Создание группы ТОЛЬКО для конкретного модуля
marketplaceProtected := app.Group("/api/v1/marketplace", mw.AuthRequiredJWT)
marketplaceProtected.Post("/listings", h.CreateListing)

// Middleware применится ТОЛЬКО к /api/v1/marketplace/*
// Публичные маршруты /api/v1/storefronts остаются доступными
```

### Альтернативный подход: Inline middleware

Если нужно применить разные middleware к разным маршрутам в одной группе:

```go
api := app.Group("/api/v1")

// Публичные маршруты - БЕЗ middleware
api.Get("/storefronts", h.ListStorefronts)
api.Get("/storefronts/search", h.SearchStorefronts)

// Защищенные маршруты - с inline middleware
api.Post("/storefronts", mw.AuthRequiredJWT, h.CreateStorefront)
api.Put("/storefronts/:id", mw.AuthRequiredJWT, h.UpdateStorefront)
```

## Конкретный пример исправления

### Проблемный код в `marketplace/handler/handler.go`

```go
// ❌ БЫЛО (line 368)
authedAPIGroup := app.Group("/api/v1", mw.AuthRequiredJWT)

marketplaceProtected := authedAPIGroup.Group("/marketplace")
marketplaceProtected.Post("/listings", h.Listings.CreateListing)
// ... остальные protected маршруты

chat := authedAPIGroup.Group("/marketplace/chat")
chat.Get("/", h.Chat.GetChats)
```

**Результат:** Middleware применялся ко ВСЕМ `/api/v1/*`, включая публичные `/api/v1/storefronts`.

### Исправленный код

```go
// ✅ СТАЛО
// ВАЖНО: НЕ создаём глобальную группу /api/v1 с middleware!
// Вместо этого создаём КОНКРЕТНУЮ подгруппу только для marketplace
marketplaceProtected := app.Group("/api/v1/marketplace", mw.AuthRequiredJWT)
marketplaceProtected.Post("/listings", h.Listings.CreateListing)
// ... остальные protected маршруты

// Чат требует аутентификацию
chat := app.Group("/api/v1/marketplace/chat", mw.AuthRequiredJWT)
chat.Get("/", h.Chat.GetChats)
```

**Результат:** Middleware применяется ТОЛЬКО к `/api/v1/marketplace/*`, публичные маршруты работают корректно.

## Диагностика проблемы

### 1. Проверка middleware вызовов

Добавьте временное логирование в `AuthRequiredJWT`:

```go
func (m *Middleware) AuthRequiredJWT(c *fiber.Ctx) error {
    // DEBUG
    pkglogger.GetLogger().Info("🔐 AuthRequiredJWT called", "path", c.Path())

    // ... остальная логика
}
```

Если видите логи для публичных маршрутов - middleware "утекает".

### 2. Проверка всех Group() регистраций

Найдите все места, где создаются группы с широким префиксом:

```bash
# Поиск всех Group() с /api/v1
cd backend
grep -rn 'Group("/api/v1"' internal/proj/
```

Проверьте, не применяется ли middleware к слишком широкой группе.

### 3. Список зарегистрированных маршрутов

Добавьте в `server.go` логирование всех маршрутов:

```go
func (s *Server) RegisterRoutes() {
    // ... регистрация маршрутов

    // DEBUG: логируем все маршруты
    stack := s.app.Stack()
    for _, route := range stack {
        for _, r := range route {
            if strings.Contains(r.Path, "storefronts") {
                logger.Info().
                    Str("method", r.Method).
                    Str("path", r.Path).
                    Msg("Registered storefront route")
            }
        }
    }
}
```

## Рекомендации по архитектуре

### 1. Префиксы групп должны быть специфичными

```go
// ❌ Слишком широко
app.Group("/api/v1", middleware)

// ✅ Специфично
app.Group("/api/v1/marketplace", middleware)
app.Group("/api/v1/admin", middleware)
app.Group("/api/v1/storefronts/:id/products", middleware)
```

### 2. Используйте inline middleware для смешанных маршрутов

Если в одной логической группе есть и публичные, и защищенные маршруты:

```go
api := app.Group("/api/v1/storefronts")

// Публичные
api.Get("/", h.List)           // БЕЗ middleware
api.Get("/:id", h.Get)         // БЕЗ middleware

// Защищенные
api.Post("/", mw.Auth, h.Create)       // С inline middleware
api.Put("/:id", mw.Auth, h.Update)     // С inline middleware
```

### 3. Избегайте вложенных групп с родительским middleware

```go
// ❌ ПЛОХО - родительский middleware применится к ВСЕМ дочерним группам
parent := app.Group("/api/v1", parentMiddleware)
child1 := parent.Group("/marketplace")  // наследует parentMiddleware
child2 := parent.Group("/storefronts")  // тоже наследует!

// ✅ ХОРОШО - независимые группы
marketplace := app.Group("/api/v1/marketplace", marketplaceMiddleware)
storefronts := app.Group("/api/v1/storefronts", storefrontsMiddleware)
```

## Порядок регистрации маршрутов

Fiber использует **Last-In-First-Out** для matching маршрутов, но middleware применяется на этапе регистрации групп.

**Важно:** Порядок регистрации маршрутов НЕ ВЛИЯЕТ на применение middleware от групп. Middleware применяется на основе prefix matching при создании группы.

### Пример с конфликтующими префиксами

```go
// Эти группы будут конфликтовать, если используют один префикс
orders := app.Group("/api/v1/storefronts/:id/orders", mw.RequireAuth())  // ❌
storefronts := app.Group("/api/v1/storefronts")  // Тоже получит RequireAuth!

// Правильно - используйте прямую регистрацию
app.Get("/api/v1/storefronts/:id/orders", mw.RequireAuth(), h.GetOrders)
app.Get("/api/v1/storefronts", h.ListStorefronts)  // БЕЗ middleware
```

## Troubleshooting чеклист

Если публичный endpoint возвращает 401:

1. ✅ Проверьте, нет ли родительской группы `/api/v1` с middleware
2. ✅ Проверьте, нет ли другого модуля, создающего группу с перекрывающимся префиксом
3. ✅ Убедитесь, что маршрут зарегистрирован напрямую через `app.Get()` или в группе БЕЗ middleware
4. ✅ Проверьте orders/cart маршруты - часто они создают `/api/v1/storefronts/:id/...` группы с middleware
5. ✅ Добавьте временное логирование в `AuthRequiredJWT` для отладки
6. ✅ Проверьте логи - handler должен выполняться, а не останавливаться на middleware

## Связанные файлы

- `internal/middleware/middleware.go` - AuthRequiredJWT, RequireAuth, OptionalAuth
- `internal/proj/marketplace/handler/handler.go` - регистрация marketplace маршрутов
- `internal/proj/storefronts/module.go` - регистрация storefronts маршрутов
- `internal/proj/orders/handler/routes.go` - регистрация cart/orders маршрутов
- `internal/server/server.go` - инициализация и регистрация всех модулей

## JWT Microservice Architecture (2025-09-30)

### Auth Service Integration
**Auth Microservice URL**: https://authpreprod.svetu.rs (preprod), https://auth.svetu.rs (prod)
- **Algorithm**: RS256 (RSA asymmetric encryption)
- **Private Key**: Stored ONLY in auth service
- **Public Key**: Distributed to all backend services for token validation
- **Library**: `github.com/sveturs/auth`

### JWT Token Structure
```json
{
  "iss": "https://auth.svetu.rs",
  "sub": "6",
  "aud": ["https://svetu.rs"],
  "exp": 1759266831,
  "user_id": 6,
  "email": "user@example.com",
  "name": "User Name",
  "roles": ["admin", "user"],
  "provider": "google",
  "email_verified": true
}
```

### Cookie-Based Authentication
- **Access Token**: `access_token` cookie, HttpOnly, 15 minutes
- **Refresh Token**: `refresh_token` cookie, HttpOnly, 720 hours (30 days)

## Middleware Chain

### ✅ Current Pattern (2025-09-30)
```go
import (
    authMiddleware "github.com/sveturs/auth/pkg/http/fiber/middleware"
)

// 1. JWTParserWithCookies - reads token from cookies
// 2. RequireAuthString("admin") - validates role from JWT
adminGroup := app.Group("/api/v1/admin/categories",
    mw.JWTParserWithCookies(),
    authMiddleware.RequireAuthString("admin"),
)
```

### ❌ Deprecated Pattern (DO NOT USE)
```go
// OLD - uses local middleware
mw.AdminRequired
mw.AuthRequiredJWT
```

## Валидация JWT токенов

JWT токены валидируются автоматически через middleware из библиотеки `github.com/sveturs/auth/pkg/http/fiber/middleware`.

**Важно:** Библиотека сама управляет получением и кэшированием публичного ключа от Auth Service.
Вам **не нужно** вручную конфигурировать путь к публичному ключу или загружать его.

### Как это работает:
1. Middleware `JWTParser()` автоматически получает публичный ключ от Auth Service
2. Ключ кэшируется для последующих запросов
3. Токен валидируется с использованием RS256 алгоритма
4. Claims из токена (user_id, email, roles) становятся доступны в контексте запроса

## История изменений

### 2025-10-01 - Исправление middleware leak в marketplace handler

**Проблема**: После миграции на sveturs/auth микросервис, публичный endpoint `/api/v1/storefronts` возвращал 401 Unauthorized.

**Root Cause**: В `internal/proj/marketplace/handler/handler.go:369` была создана глобальная группа:
```go
authedAPIGroup := app.Group("/api/v1", mw.JWTParser(), authMiddleware.RequireAuth())
```

Эта группа применяла middleware **КО ВСЕМ** маршрутам с префиксом `/api/v1/*`, включая:
- `/api/v1/storefronts` (должен быть публичным)
- `/api/v1/marketplace/search` (должен быть публичным)
- И другие публичные endpoints

**Решение**:
1. Удалена глобальная группа `/api/v1` с middleware
2. Все защищенные маршруты marketplace переведены на inline регистрацию:
```go
// Создаем массив middleware для переиспользования
authMW := []fiber.Handler{mw.JWTParser(), authMiddleware.RequireAuth()}

// Регистрируем каждый защищенный маршрут отдельно
app.Post("/api/v1/marketplace/listings", append(authMW, h.Listings.CreateListing)...)
app.Put("/api/v1/marketplace/listings/:id", append(authMW, h.Listings.UpdateListing)...)
// и т.д.
```

3. Chat routes переведены на узкий префикс:
```go
chat := app.Group("/api/v1/marketplace/chat", mw.JWTParser(), authMiddleware.RequireAuth())
```

4. Orders модуль (`internal/proj/orders/handler/routes.go`) также исправлен - удалены группы с широким префиксом `/api/v1/storefronts/:storefront_id/cart`.

**Файлы изменены**:
- `internal/proj/marketplace/handler/handler.go` - удалена глобальная группа `/api/v1`
- `internal/proj/orders/handler/routes.go` - заменены группы на inline middleware
- 15+ других файлов - замена `AuthRequiredJWT` на `JWTParser() + RequireAuth()`

**Результат**: Все публичные endpoints теперь работают без авторизации, защищенные - с правильной проверкой JWT токенов.

---

### 2025-09-30 - Миграция на microservice auth с RS256 JWT

**Проблема:**
1. Публичный endpoint `/api/v1/storefronts` возвращал 401
2. Admin endpoints возвращали 401/400 после валидного ответа handler'а
3. Old `AdminRequired` middleware конфликтовал с новым из библиотеки

**Решение:**
1. **Исправлена утечка middleware**:
   - Заменили `app.Group("/api/v1", mw.AuthRequiredJWT)` на `app.Group("/api/v1/marketplace", mw.AuthRequiredJWT)`
   - Заменили `authedAPIGroup.Group("/marketplace/chat")` на `app.Group("/api/v1/marketplace/chat", mw.AuthRequiredJWT)`

2. **Мигрированы все admin routes на библиотечный middleware**:
   - Заменили `mw.AdminRequired` → `authMiddleware.RequireAuthString("admin")`
   - Добавили `mw.JWTParserWithCookies()` для чтения токенов из cookies
   - Обновлено 10 модулей: marketplace, search_admin, logistics, delivery, analytics, search_optimization, behavior_tracking, translation_admin, subscriptions

3. **Исправлена Fiber route inheritance проблема**:
   - Subscriptions создавал `/api/v1/admin` группу, которая загрязняла все `/api/v1/admin/*` routes
   - Изменено на `/api/v1/admin/subscriptions` для изоляции middleware

**Файлы:**
- `internal/proj/marketplace/handler/handler.go` - Categories, Attributes, Listings admin
- `internal/proj/search_admin/handler/routes.go` - Search admin
- `internal/proj/admin/logistics/module.go` - Logistics admin
- `internal/proj/delivery/module.go` - Delivery admin
- `internal/proj/analytics/routes/routes.go` - Analytics
- `internal/proj/search_optimization/module.go` - Search optimization
- `internal/proj/behavior_tracking/module.go` - Behavior tracking
- `internal/proj/translation_admin/module.go` - Translation admin
- `internal/proj/subscriptions/handler/routes.go` - Subscriptions admin (route path fix)
- `internal/middleware/middleware.go` - JWTParserWithCookies, loadAuthServicePublicKey

**Результат:**
- ✅ Публичные маршруты работают без аутентификации
- ✅ Admin endpoints корректно валидируют JWT токены через RS256
- ✅ Cookie-based authentication работает для фронтенда
- ✅ Нет конфликтов middleware между модулями
