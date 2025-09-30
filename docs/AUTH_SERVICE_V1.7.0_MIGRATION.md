# Миграция на auth-service v1.7.0

## Дата: 2025-09-30

## Обзор

Успешно выполнена миграция svetu backend на использование централизованного auth-service v1.7.0.
Таблицы `users` и `user_roles` удалены из локальной базы данных.
Все операции с пользователями теперь проходят через auth-service.

## Основные изменения

### 1. Обновление зависимости

```bash
github.com/sveturs/auth v1.5.0 => v1.7.0
```

### 2. Архитектурные изменения

**До (v1.5.0):**
- Один сервис `AuthService` для всех операций
- Прямая работа с PostgreSQL через `users/storage/postgres`
- Таблицы `users` и `user_roles` в локальной БД

**После (v1.7.0):**
- Два отдельных сервиса:
  - `AuthService` - аутентификация (login, register, logout, validate)
  - `UserService` - управление пользователями (CRUD, roles, admin)
- Использование только `entity` типов из `github.com/sveturs/auth/pkg/http/entity`
- Middleware из `github.com/sveturs/auth/pkg/http/fiber/middleware`
- Все данные пользователей в auth-service

### 3. Измененные файлы

**Созданы:**
- `backend/internal/proj/users/service/converter.go` - конвертеры entity ↔ models
- `backend/migrations/000020_remove_users.up.sql` - миграция удаления таблиц
- `backend/migrations/000020_remove_users.down.sql` - необратимая операция

**Удалены:**
- `backend/internal/proj/users/storage/postgres/` (весь пакет)
- `backend/internal/proj/users/storage/interface.go`
- `backend/fixtures/000047_users.up.sql`
- `backend/fixtures/000060_user_roles.up.sql`

**Обновлены:**
- `backend/go.mod` - зависимость auth v1.7.0
- `backend/internal/server/server.go` - создание AuthService и UserService
- `backend/internal/proj/global/service/service.go` - прием обоих сервисов
- `backend/internal/proj/users/service/service.go` - использование обоих сервисов
- `backend/internal/proj/users/service/user.go` - полностью переписан
- `backend/internal/proj/users/handler/auth.go` - использование entity типов
- `backend/internal/storage/postgres/db.go` - стабы методов
- `backend/internal/storage/postgres/admin_methods.go` - стабы методов
- `backend/fixtures/000078_enable_triggers.up.sql` - закомментированы REFRESH views

## Функциональность auth-service v1.7.0

### Реализованные endpoint'ы

#### Аутентификация (AuthService)
- `POST /api/v1/auth/register` - регистрация
- `POST /api/v1/auth/login` - вход
- `POST /api/v1/auth/logout` - выход
- `POST /api/v1/auth/refresh` - обновление токена
- `GET /api/v1/auth/validate` - валидация токена

#### Управление пользователями (UserService)
- `GET /api/v1/users/all` - все пользователи
- `GET /api/v1/users/:id` - пользователь по ID
- `GET /api/v1/users/by-email?email={email}` - поиск по email
- `PATCH /api/v1/users/:id` - обновление профиля
- `PATCH /api/v1/users/:id/status` - обновление статуса
- `DELETE /api/v1/users/:id?permanent=true` - удаление

#### Управление ролями
- `GET /api/v1/roles` - все роли
- `GET /api/v1/users/:id/roles` - роли пользователя
- `POST /api/v1/users/:id/roles` - добавить роль
- `DELETE /api/v1/users/:id/roles/:role` - удалить роль
- `GET /api/v1/users/:id/is-admin` - проверка админа
- `GET /api/v1/users?role=admin` - пользователи по роли

### Предустановленные роли

- `user` - обычный пользователь
- `admin` - администратор
- `moderator` - модератор
- `support` - поддержка
- `superadmin` - супер-администратор

## Изменения в коде

### Инициализация сервисов (server.go)

```go
authClient, err := authclient.NewClientWithResponses(cfg.AuthServiceURL)

zerologLogger := *logger.Get()
authServiceInstance := authService.NewAuthService(authClient, zerologLogger)
userServiceInstance := authService.NewUserService(authClient, zerologLogger)
oauthServiceInstance := authService.NewOAuthService(authClient)

services := globalService.NewService(ctx, db, cfg, translationService,
                                    authServiceInstance, userServiceInstance)
```

### Использование entity типов (handler)

```go
// До
var req authclient.EntityUserRegistrationRequest

// После
var req entity.UserRegistrationRequest
```

### Примеры использования UserService

```go
// Получение пользователя
user, err := s.userService.GetUser(ctx, userID)

// Поиск по email
user, err := s.userService.GetUserByEmail(ctx, email)

// Обновление профиля
req := entity.UpdateProfileRequest{
    Name: &newName,
    Bio:  &newBio,
}
profile, err := s.userService.UpdateUserProfile(ctx, userID, req)

// Проверка админа
adminResp, err := s.userService.IsUserAdmin(ctx, userID)

// Добавление роли
req := entity.AddRoleRequest{Role: "admin"}
_, err := s.userService.AddUserRole(ctx, userID, req)
```

## Backward Compatibility

### Совместимость сохранена для:

1. HTTP API endpoints - все endpoint'ы работают без изменений
2. Response структуры - конвертеры обеспечивают совместимость
3. Middleware - интерфейс тот же

### Требуют адаптации:

1. Прямые SQL запросы - используйте UserService
2. User storage - удален, все через auth-service
3. JWT validation - через JWTParser middleware

## Ссылки

- Auth Service: https://github.com/sveturs/auth
- Auth Service v1.7.0: https://github.com/sveturs/auth/releases/tag/v1.7.0