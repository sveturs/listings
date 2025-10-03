# Решение проблемы 404 на защищённых роутах

## Проблема

При обращении к защищённому endpoint'у `/api/v1/review-permission/listing/:id` возникала ошибка 401 "Требуется авторизация", хотя JWT токен был валидным.

## Симптомы

1. Прямой curl запрос с `Authorization: Bearer <token>` возвращал 401
2. Запрос через BFF proxy также возвращал 401
3. В логах backend виден запрос, но middleware отклоняет его
4. Другие защищённые endpoints (например `/api/v1/auth/me`) работали корректно

## Причина

**Неправильное использование middleware аутентификации!**

В `backend/internal/proj/reviews/handler/handler.go` использовался **локальный** middleware `mw.AuthRequiredJWT`, который не работает с библиотекой `github.com/sveturs/auth`.

Правильный подход - использовать middleware из auth библиотеки:
- `jwtParserMW` - парсит JWT токен и сохраняет данные в контексте
- `authMiddleware.RequireAuth()` - проверяет что пользователь аутентифицирован

## Решение

### 1. Обновить imports в handler.go

```go
import (
    "github.com/gofiber/fiber/v2"
    authMiddleware "github.com/sveturs/auth/pkg/http/fiber/middleware"  // ← Добавить

    "backend/internal/middleware"
    globalService "backend/internal/proj/global/service"
)
```

### 2. Обновить структуру Handler

```go
type Handler struct {
    Review      *ReviewHandler
    jwtParserMW fiber.Handler  // ← Добавить поле
}

func NewHandler(services globalService.ServicesInterface, jwtParserMW fiber.Handler) *Handler {
    return &Handler{
        Review:      NewReviewHandler(services),
        jwtParserMW: jwtParserMW,  // ← Сохранить middleware
    }
}
```

### 3. Обновить роуты на правильный middleware

**Было (неправильно):**
```go
app.Get("/api/v1/review-permission/listing/:id", mw.AuthRequiredJWT, h.Review.CanReviewListing)
```

**Стало (правильно):**
```go
app.Get("/api/v1/review-permission/listing/:id", h.jwtParserMW, authMiddleware.RequireAuth(), h.Review.CanReviewListing)
```

### 4. Обновить server.go для передачи middleware

В `backend/internal/server/server.go`:

```go
// Создание jwtParserMW (уже есть для users handler)
jwtParserMW := authMiddleware.JWTParser(authServiceInstance)

// Передать middleware в reviews handler
reviewHandler := reviewHandler.NewHandler(services, jwtParserMW)  // ← Добавить параметр
```

## Правильная цепочка middleware для защищённых роутов

### GET запросы (только аутентификация, БЕЗ CSRF)
```go
app.Get("/api/v1/protected/resource",
    h.jwtParserMW,                    // Парсит JWT
    authMiddleware.RequireAuth(),     // Проверяет аутентификацию
    h.Handler)
```

### POST/PUT/DELETE запросы (аутентификация + CSRF)
```go
app.Post("/api/v1/protected/resource",
    h.jwtParserMW,                    // Парсит JWT
    authMiddleware.RequireAuth(),     // Проверяет аутентификацию
    mw.CSRFProtection(),              // Проверяет CSRF токен
    h.Handler)
```

### Admin endpoints (требуют роль)
```go
app.Get("/api/v1/admin/resource",
    h.jwtParserMW,
    authMiddleware.RequireAuth("admin"),  // Требует роль admin
    h.Handler)
```

## Получение данных пользователя в handler

После правильного middleware, данные доступны через библиотечные функции:

```go
import authMiddleware "github.com/sveturs/auth/pkg/http/fiber/middleware"

func (h *Handler) MyHandler(c *fiber.Ctx) error {
    // Получить user ID
    userID, ok := authMiddleware.GetUserID(c)
    if !ok || userID == 0 {
        return utils.ErrorResponse(c, fiber.StatusUnauthorized, "auth.error.unauthorized")
    }

    // Получить email
    email, ok := authMiddleware.GetEmail(c)

    // Получить роли
    roles, ok := authMiddleware.GetRoles(c)

    // Использовать данные
    // ...
}
```

## Примеры правильной регистрации роутов

### Reviews handler (пример из исправленного кода)

```go
func (h *Handler) RegisterRoutes(app *fiber.App, mw *middleware.Middleware) error {
    // Публичные endpoints (без middleware)
    app.Get("/api/v1/reviews", h.Review.GetReviews)
    app.Get("/api/v1/reviews/:id", h.Review.GetReviewByID)

    // Protected GET endpoints (только Auth, БЕЗ CSRF)
    app.Get("/api/v1/review-permission/listing/:id",
        h.jwtParserMW,
        authMiddleware.RequireAuth(),
        h.Review.CanReviewListing)

    // Protected POST/PUT/DELETE endpoints (Auth + CSRF)
    app.Post("/api/v1/reviews/draft",
        h.jwtParserMW,
        authMiddleware.RequireAuth(),
        mw.CSRFProtection(),
        h.Review.CreateDraftReview)

    return nil
}
```

### Users handler (эталонный пример)

См. `backend/internal/proj/users/handler/routes.go` - там все сделано правильно с самого начала.

## Проверка что всё работает

### 1. Получить JWT токен
```bash
ssh svetu@svetu.rs "cd /opt/svetu-authpreprod && sed 's|/data/auth_svetu/keys/private.pem|./keys/private.pem|g' scripts/create_admin_jwt.go > /tmp/create_jwt_fixed.go && go run /tmp/create_jwt_fixed.go" > /tmp/jwt_token.txt
```

### 2. Протестировать endpoint напрямую
```bash
bash -c 'TOKEN=$(cat /tmp/jwt_token.txt); curl -s -H "Authorization: Bearer $TOKEN" http://localhost:3000/api/v1/review-permission/listing/351'
```

Должен вернуть 200 с данными:
```json
{
  "data": {
    "can_review": false,
    "reason": "Вы не можете оставить отзыв на свой товар",
    "has_existing_review": false
  },
  "success": true
}
```

### 3. Протестировать через BFF proxy
```bash
curl -s -b "access_token=$(cat /tmp/jwt_token.txt)" http://localhost:3001/api/v2/review-permission/listing/351
```

Также должен вернуть 200 с данными.

## Частые ошибки

### ❌ Использование локального middleware
```go
app.Get("/api/v1/protected", mw.AuthRequiredJWT, handler)  // НЕ работает с auth library!
```

### ❌ Забыли добавить jwtParserMW в структуру
```go
type Handler struct {
    Review *ReviewHandler
    // Забыли добавить jwtParserMW fiber.Handler
}
```

### ❌ Забыли обновить constructor
```go
func NewHandler(services globalService.ServicesInterface) *Handler {  // Нет параметра jwtParserMW!
    return &Handler{Review: NewReviewHandler(services)}
}
```

### ❌ Забыли обновить server.go
```go
reviewHandler := reviewHandler.NewHandler(services)  // Забыли передать jwtParserMW!
```

## Ключевые файлы для проверки

При добавлении нового защищённого модуля, проверь:

1. **Handler structure** (`internal/proj/{module}/handler/handler.go`):
   - Есть поле `jwtParserMW fiber.Handler`
   - Constructor принимает `jwtParserMW` параметр
   - Import `authMiddleware "github.com/sveturs/auth/pkg/http/fiber/middleware"`

2. **Routes** (`RegisterRoutes` метод):
   - GET: `h.jwtParserMW, authMiddleware.RequireAuth(), handler`
   - POST/PUT/DELETE: `h.jwtParserMW, authMiddleware.RequireAuth(), mw.CSRFProtection(), handler`

3. **Server initialization** (`internal/server/server.go`):
   - Передаётся `jwtParserMW` при создании handler
   - `jwtParserMW := authMiddleware.JWTParser(authServiceInstance)` уже создан

4. **Handler methods**:
   - Используют `authMiddleware.GetUserID(c)` для получения user ID
   - Не используют устаревшие методы типа `c.Locals("userID")`

## Дата решения

2025-10-02

## Время потраченное на решение

~2 часа (множество попыток с неправильными подходами)

## Автор исправления

Claude (с помощью Dmitry)
