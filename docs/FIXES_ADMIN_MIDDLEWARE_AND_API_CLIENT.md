# Исправление AdminRequired Middleware и ApiClient JWT

## Дата: 2025-10-01

## Проблемы

### 1. AdminRequired Middleware возвращал 404 после успешного выполнения

**Симптомы:**
- Backend логи показывали успешное выполнение handler (status=200)
- Но сразу после этого Fiber ErrorHandler возвращал 404 "Cannot GET /api/v1/admin/..."
- Затрагивало **ВСЕ** endpoints в `/api/v1/admin/*` группе

**Пример из логов:**
```
3:40PM INF RESPONSE duration=232.599595 status=200
3:40PM ERR Error in handler error="Cannot GET /api/v1/admin/categories-all"
```

**Корневая причина:**
В `middleware.go` функция `AdminRequired` использовала `return c.Next()` который возвращает результат следующего handler в цепочке. В Fiber это приводит к тому, что middleware пытается обработать результат handler как ошибку.

**Код ДО (неправильно):**
```go
// middleware.go:167
if role == "admin" {
    logger.Info().Msg("Access granted")
    c.Locals("admin_id", userID)
    return c.Next()  // ❌ BUG: Возвращает результат c.Next()
}
```

**Код ПОСЛЕ (правильно):**
```go
// middleware.go:167-168
if role == "admin" {
    logger.Info().Msg("Access granted")
    c.Locals("admin_id", userID)
    // ИСПРАВЛЕНИЕ: не возвращаем результат c.Next(), а просто вызываем его
    return nil  // ✅ FIX: Возвращаем nil для успеха
}
```

**Изменения:**
- `backend/internal/middleware/middleware.go:167` - JWT token admin check
- `backend/internal/middleware/middleware.go:181` - Hardcoded user ID check

---

### 2. ApiClient не передавал JWT токен в запросах

**Симптомы:**
- Спам ошибок 401 Unauthorized для `/api/v1/admin/marketplace-translations/status`
- Backend логи показывали `user_id=0` вместо реального ID
- Frontend консоль: 75 параллельных запросов с ошибкой 401

**Пример из логов:**
```
3:46PM INF AdminRequired: checking user ID user_id=0 user_id_ok=false
3:46PM INF AdminRequired: User ID not found or invalid
3:46PM INF RESPONSE status=401
```

**Корневая причина:**
Класс `ApiClient` в `frontend/svetu/src/services/api-client.ts` не добавлял JWT токен в заголовок Authorization. Он использовался в `adminApi.getTranslationStatus()`, но не имел логики для добавления токена.

**Решение:**
Добавлен импорт `tokenManager` и логика добавления JWT токена в метод `request()`:

```typescript
// api-client.ts:1-3
import configManager from '@/config';
import { logger } from '@/utils/logger';
import { tokenManager } from '@/utils/tokenManager';  // ✅ Добавлено

// api-client.ts:142-148
// Добавляем JWT токен из tokenManager
if (typeof window !== 'undefined') {
  const token = tokenManager.getAccessToken();
  if (token && !headers.has('Authorization')) {
    headers.set('Authorization', `Bearer ${token}`);
  }
}
```

**Изменения:**
- `frontend/svetu/src/services/api-client.ts:3` - Добавлен импорт tokenManager
- `frontend/svetu/src/services/api-client.ts:142-148` - Добавлена логика добавления JWT токена

---

### 3. TranslationStatus компонент - подавление 401 ошибок

**Дополнительное улучшение:**
Добавлена логика игнорирования 401 ошибок в консоли, так как это нормально если пользователь не залогинен при первом рендере компонента.

```typescript
// TranslationStatus.tsx:52-56
catch (error) {
  // Игнорируем ошибки 401 - это нормально если пользователь не залогинен
  if ((error as any)?.status !== 401) {
    console.error('Failed to fetch translation status:', error);
  }
}
```

**Изменения:**
- `frontend/svetu/src/components/attributes/TranslationStatus.tsx:52-56`

---

### 4. Fiber Route Conflict - категории endpoint

**Бонусное исправление:**
Изменён путь с `/categories/all` на `/categories-all` чтобы избежать конфликта с параметром `:id` в Fiber routing.

**Изменения:**
- `backend/internal/proj/marketplace/handler/handler.go:440` - Изменён маршрут
- `frontend/svetu/src/services/admin.ts:183` - Обновлён URL на frontend

---

## Как диагностировать подобные проблемы в будущем

### Признаки проблемы с `return c.Next()`

1. **Backend логи показывают успешный ответ, но клиент получает 404:**
   ```
   RESPONSE status=200
   ERR Error in handler error="Cannot GET /path"
   ```

2. **Handler выполняется корректно, но Fiber ErrorHandler срабатывает:**
   - Смотрите на middleware которые используют `return c.Next()`
   - Правильный паттерн: `return nil` после `c.Locals()` или других операций

### Признаки проблемы с отсутствующим JWT токеном

1. **Backend логи показывают `user_id=0` для защищённых endpoints:**
   ```
   AdminRequired: checking user ID user_id=0 user_id_ok=false
   ```

2. **Множественные 401 ошибки в browser console:**
   - Проверьте что API client добавляет `Authorization: Bearer <token>` заголовок
   - Проверьте что `tokenManager.getAccessToken()` возвращает токен

---

## Правильные паттерны

### Middleware в Fiber

```go
// ✅ ПРАВИЛЬНО: Возвращаем nil для успеха
func (m *Middleware) MyMiddleware(c *fiber.Ctx) error {
    // Проверки...
    if authorized {
        c.Locals("key", value)
        return nil  // Middleware успешен, продолжаем цепочку
    }
    return fiber.ErrUnauthorized
}

// ❌ НЕПРАВИЛЬНО: Возвращаем результат c.Next()
func (m *Middleware) MyMiddleware(c *fiber.Ctx) error {
    if authorized {
        return c.Next()  // BUG: может вернуть неожиданный результат
    }
    return fiber.ErrUnauthorized
}
```

### API Client с JWT

```typescript
// ✅ ПРАВИЛЬНО: Добавляем токен перед запросом
const token = tokenManager.getAccessToken();
if (token) {
  headers.set('Authorization', `Bearer ${token}`);
}

// ❌ НЕПРАВИЛЬНО: Забываем добавить токен
// Просто делаем fetch без заголовка Authorization
```

---

## Тестирование

### Проверка middleware

```bash
# Тест с валидным JWT токеном
TOKEN=$(cat /tmp/jwt_token.txt)
curl -H "Authorization: Bearer $TOKEN" http://localhost:3000/api/v1/admin/categories-all

# Должен вернуть 200 с данными, НЕ 404
```

### Проверка API client

```bash
# Проверить логи backend на наличие user_id > 0
tail -f /tmp/backend.log | grep "user_id"

# Должно быть: user_id=6 (или другой реальный ID)
# НЕ должно быть: user_id=0
```

---

## Связанные файлы

### Backend
- `backend/internal/middleware/middleware.go` - AdminRequired middleware
- `backend/internal/proj/marketplace/handler/handler.go` - Route registration

### Frontend
- `frontend/svetu/src/services/api-client.ts` - API client с JWT
- `frontend/svetu/src/services/admin.ts` - Admin API methods
- `frontend/svetu/src/components/attributes/TranslationStatus.tsx` - Компонент с запросами
- `frontend/svetu/src/utils/tokenManager.ts` - Управление JWT токенами

---

## Ссылки

- [Fiber Middleware Guide](https://docs.gofiber.io/guide/routing#middleware)
- [JWT Token Management](/data/hostel-booking-system/frontend/svetu/src/utils/tokenManager.ts)
- [Auth Service Integration](ssh://svetu@svetu.rs/opt/svetu-authpreprod/MARKETPLACE_INTEGRATION_SPEC.md)
