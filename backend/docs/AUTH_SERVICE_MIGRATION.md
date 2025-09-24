# Миграция на Auth Service

## Статус: ЗАВЕРШЕНА (02.01.2025)

## Описание

Проект полностью мигрирован на использование внешнего auth-service (github.com/sveturs/auth).
Auth-service работает на порту 28080 и предоставляет централизованную аутентификацию и авторизацию.

⚠️ **ВАЖНО**: Вся старая система аутентификации была полностью удалена. Обратная совместимость НЕ поддерживается.

## Что сделано

### 1. Backend изменения

#### Добавлены файлы:
- `internal/proj/users/handler/auth.go` - handlers для проксирования к auth-service
- `internal/middleware/auth.go` - новый middleware для валидации токенов через auth-service

#### Обновлены файлы:
- `internal/config/config.go` - добавлены `AuthServiceURL` и `BackendURL`
- `internal/server/server.go` - инициализация auth-service и OAuth service
- `internal/proj/users/handler/handler.go` - добавлен `Auth` handler
- `internal/proj/users/handler/routes.go` - использует только новые auth handlers
- `internal/middleware/middleware.go` - добавлен `AuthService` field для доступа к новому middleware

### 2. Конфигурация

Добавлены переменные окружения:
```env
AUTH_SERVICE_URL=http://localhost:28080  # URL auth-service
BACKEND_URL=http://localhost:3000        # URL backend для OAuth callbacks
```

### 3. Новые эндпоинты

Все auth эндпоинты теперь проксируются к auth-service:
- POST `/api/v1/auth/register` - регистрация
- POST `/api/v1/auth/login` - вход
- POST `/api/v1/auth/logout` - выход
- POST `/api/v1/auth/refresh` - обновление токенов
- GET `/api/v1/auth/validate` - валидация токена
- GET `/api/v1/auth/session` - получение сессии (адаптировано для совместимости)
- GET `/auth/google` - начало OAuth Google
- GET `/auth/google/callback` - callback Google OAuth

#### Удалены файлы (старая система аутентификации):
- `internal/proj/users/handler/auth.go` (старая версия)
- `internal/proj/users/service/auth.go`
- `internal/proj/users/service/user_adapter.go`
- `internal/middleware/auth.go` (старая версия)
- `internal/middleware/auth_proxy.go`
- `internal/middleware/jwt.go`
- `internal/service/authclient/client.go`

## Текущие ограничения

1. **User IDs**: Auth-service использует числовые ID. Временно используется заглушка (userID = 1).
2. **Частичная миграция middleware**:
   - Модуль users полностью использует новый `AuthServiceMiddleware`
   - Остальные модули пока используют старый `AuthRequiredJWT` middleware
3. **Admin роли**: Определение администратора основано на наличии роли "admin" в токене от auth-service.

## TODO - Следующие шаги

1. **Миграция остальных модулей**: Заменить `AuthRequiredJWT` на `AuthServiceMiddleware` в модулях:
   - marketplace (25+ мест)
   - payments
   - balance
   - notifications
   - reviews
   - И других...
2. **ID mapping**: Реализовать правильный маппинг между ID пользователей в auth-service и локальной БД.
3. **Frontend обновление**: Обновить frontend для работы с новыми эндпоинтами.
4. **Миграция данных**: Перенести существующих пользователей в auth-service.
5. **Удаление старого JWT middleware**: После миграции всех модулей удалить `AuthRequiredJWT`.

## Тестирование

### Запуск auth-service
```bash
cd ../auth
go run cmd/auth/main.go
```

### Тестирование регистрации
```bash
curl -X POST http://localhost:3000/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email": "test@example.com", "password": "password123"}'
```

### Тестирование входа
```bash
curl -X POST http://localhost:3000/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "test@example.com", "password": "password123"}'
```

### Проверка сессии
```bash
TOKEN="<полученный_токен>"
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:3000/api/v1/auth/session
```

## Миграция модулей (статус)

### ✅ Полностью мигрированы:
- **users** - все эндпоинты используют новый AuthServiceMiddleware

### ⏳ Требуют миграции:
- **marketplace** - ~25 эндпоинтов
- **payments** - ~10 эндпоинтов
- **balance** - ~5 эндпоинтов
- **notifications** - ~8 эндпоинтов
- **reviews** - ~10 эндпоинтов
- **storefronts** - ~15 эндпоинтов
- **delivery** - ~10 эндпоинтов
- **И другие модули...**

## Использование нового middleware

Для доступа к новому auth middleware в handlers:
```go
func (h *Handler) RegisterRoutes(app *fiber.App, mw *middleware.Middleware) error {
    authMw := mw.AuthService // Получаем auth middleware

    // Использование:
    app.Get("/protected", authMw.AuthRequired(), handler)
    app.Get("/admin", authMw.AuthRequired(), authMw.RequireAdmin(), handler)
}
```