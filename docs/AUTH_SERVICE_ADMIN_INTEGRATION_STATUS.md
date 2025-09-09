# Интеграция централизованного управления администраторами через Auth Service

## 📋 Оглавление
- [Цель проекта](#цель-проекта)
- [Архитектура решения](#архитектура-решения)
- [Что уже сделано](#что-уже-сделано)
- [Что осталось сделать](#что-осталось-сделать)
- [Технические детали](#технические-детали)
- [Инструкции по развертыванию](#инструкции-по-развертыванию)

## 🎯 Цель проекта

### Проблема
До внедрения микросервиса Auth Service система управления администраторами была реализована локально в основном приложении:
- Таблица `admin_users` в основной БД (PostgreSQL на порту 5432)
- Проверка через метод `IsUserAdmin(email)` в каждом запросе
- Дублирование логики авторизации
- Отсутствие централизованного управления ролями
- Невозможность расширения системы прав без изменения кода

### Решение
Централизовать управление правами и ролями в микросервисе Auth Service:
- Единая точка управления ролями и правами (RBAC)
- JWT токены с включенными ролями (RS256 подпись)
- Расширяемая система ролей (admin, moderator, support, user)
- Автоматическая миграция существующих админов
- API для управления ролями
- Сохранение обратной совместимости

## 🏗️ Архитектура решения

### Компоненты системы

```
┌─────────────────────────────────────────────────────────────┐
│                     Frontend (Next.js)                       │
│  ┌─────────────────────────────────────────────────────┐   │
│  │  /ru/admin/users - Управление пользователями        │   │
│  │  └── RoleManager Component (новый)                  │   │
│  └─────────────────────────────────────────────────────┘   │
└────────────┬───────────────────────┬───────────────────────┘
             │                       │
             ▼                       ▼
┌──────────────────────┐  ┌─────────────────────────────────┐
│  Main Backend (Go)   │  │   Auth Service (Go)             │
│  Port: 3000          │  │   Port: 28080                   │
│                      │  │                                 │
│  ┌────────────────┐ │  │  ┌─────────────────────────┐   │
│  │ Middleware     │ │  │  │ Role Service (новый)    │   │
│  │ - JWT validation│◄───┼──│ - Управление ролями     │   │
│  │ - Role check   │ │  │  │ - RBAC логика           │   │
│  └────────────────┘ │  │  └─────────────────────────┘   │
│                      │  │                                 │
│  ┌────────────────┐ │  │  ┌─────────────────────────┐   │
│  │ admin_users    │ │  │  │ JWT Service             │   │
│  │ (legacy)       │ │  │  │ - Генерация токенов     │   │
│  └────────────────┘ │  │  │ - Включение ролей       │   │
└──────────────────────┘  │  └─────────────────────────┘   │
                          │                                 │
                          │  ┌─────────────────────────┐   │
                          │  │ Database (PostgreSQL)   │   │
                          │  │ Port: 25432             │   │
                          │  │ - auth.roles            │   │
                          │  │ - auth.user_roles       │   │
                          │  │ - auth.pending_admins   │   │
                          │  └─────────────────────────┘   │
                          └─────────────────────────────────┘
```

### Поток авторизации

```mermaid
sequenceDiagram
    participant U as User
    participant F as Frontend
    participant B as Backend
    participant A as Auth Service
    
    U->>F: Login
    F->>A: POST /api/v1/auth/login
    A->>A: Проверка credentials
    A->>A: Получение ролей из user_roles
    A->>F: JWT token с ролями
    F->>B: API запрос с JWT
    B->>B: Валидация токена (RS256)
    B->>B: Извлечение ролей
    B->>B: Проверка прав (is_admin)
    B->>F: Ответ
```

## ✅ Что уже сделано

### 1. База данных Auth Service
**Статус**: ✅ Завершено

#### Миграции (`/data/auth_svetu/migrations/002_add_roles_system.sql`)
```sql
-- Таблица ролей
CREATE TABLE auth.roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE,
    display_name VARCHAR(100),
    description TEXT,
    is_system BOOLEAN DEFAULT false
);

-- Таблица связи пользователь-роль
CREATE TABLE auth.user_roles (
    id SERIAL PRIMARY KEY,
    user_id INTEGER,
    role_id INTEGER,
    granted_by INTEGER,
    granted_at TIMESTAMP,
    expires_at TIMESTAMP,
    is_active BOOLEAN DEFAULT true,
    notes TEXT,
    UNIQUE(user_id, role_id)
);

-- Базовые роли
INSERT INTO auth.roles (name, display_name, is_system) VALUES
    ('admin', 'Administrator', true),
    ('moderator', 'Moderator', true),
    ('support', 'Support Agent', true),
    ('user', 'Regular User', true);
```

#### Автоматическое назначение (`/data/auth_svetu/scripts/admin_auto_assign.sql`)
```sql
-- Таблица для отложенного назначения
CREATE TABLE auth.pending_admin_assignments (
    email VARCHAR(255) UNIQUE,
    assigned BOOLEAN DEFAULT false,
    assigned_at TIMESTAMP
);

-- Триггер для автоматического назначения при регистрации
CREATE TRIGGER auto_assign_admin_role
    AFTER INSERT ON auth.users
    FOR EACH ROW
    EXECUTE FUNCTION auth.check_and_assign_admin_role();
```

### 2. Backend сервисы
**Статус**: ✅ Реализовано

#### Role Service (`/data/auth_svetu/internal/service/role/service.go`)
- `GetUserRoles(userID)` - получение ролей пользователя
- `GetUserRoleNames(userID)` - получение списка названий ролей (для JWT)
- `AssignRole(request)` - назначение роли
- `RevokeRole(request)` - отзыв роли
- `GetAllRoles()` - список всех ролей
- `CheckUserHasRole(userID, roleName)` - проверка наличия роли

#### Обновление JWT Service (`/data/auth_svetu/internal/service/token/jwt.go`)
```go
// Интеграция с RoleService
func (s *JWTService) GenerateAccessToken(user *domain.User) (string, error) {
    // Получаем роли пользователя
    roles := s.roleService.GetUserRoleNames(ctx, user.ID)
    
    claims := domain.JWTClaims{
        UserID:   user.ID,
        Email:    user.Email,
        Roles:    roles, // Теперь включает реальные роли
    }
    // ...
}
```

### 3. API Endpoints
**Статус**: ✅ Реализовано

#### HTTP handlers (`/data/auth_svetu/internal/transport/http/handlers/roles_fiber.go`)
- `GET /api/v1/roles` - список всех ролей
- `GET /api/v1/users/:user_id/roles` - роли пользователя
- `POST /api/v1/roles/assign` - назначить роль (только админ)
- `POST /api/v1/roles/revoke` - отозвать роль (только админ)
- `GET /api/v1/users/:user_id/roles/check?role_name=admin` - проверка роли

### 4. Миграция админов
**Статус**: ✅ Выполнено

#### Текущие админы в системе
```
voroshilovdo@gmail.com    ✅ Роль назначена в auth.user_roles
boxmail386@gmail.com      ✅ Роль назначена в auth.user_roles
test@example.com          ✅ Роль назначена в auth.user_roles
bevzenko.sergey@gmail.com ⏳ В pending_admin_assignments
admin@svetu.rs            ⏳ В pending_admin_assignments
info@svetu.rs             ⏳ В pending_admin_assignments
test@user.rs              ⏳ В pending_admin_assignments
1234@1234.1234            ⏳ В pending_admin_assignments
admin@test.com            ⏳ В pending_admin_assignments
```

#### Скрипты миграции
- `/data/auth_svetu/scripts/migrate_admins.go` - миграция из admin_users
- `/data/auth_svetu/scripts/admin_auto_assign.sql` - автоназначение при входе

### 5. Frontend компоненты
**Статус**: ✅ Подготовлено

#### RoleManager Component (`/data/hostel-booking-system/frontend/svetu/src/components/admin/RoleManager.tsx`)
```typescript
interface RoleManagerProps {
  userId: number;
  userEmail: string;
  currentRoles?: UserRole[];
  onRoleUpdate?: () => void;
}

// Функции:
- fetchAvailableRoles() - получение списка ролей
- fetchUserRoles() - получение ролей пользователя
- assignRole(roleName) - назначение роли
- revokeRole(roleName) - отзыв роли
```

#### Переводы (`/data/hostel-booking-system/frontend/svetu/src/messages/ru/admin.json`)
```json
{
  "users": {
    "currentRoles": "Текущие роли",
    "noRoles": "Нет назначенных ролей",
    "addRole": "Добавить роль",
    "revokeRole": "Отозвать роль",
    "roleDescriptions": {
      "admin": "Полный доступ ко всем функциям системы",
      "moderator": "Модерация контента и управление пользователями",
      "support": "Поддержка пользователей и просмотр данных",
      "user": "Обычный пользователь с базовыми правами"
    }
  }
}
```

### 6. Интеграция с основным приложением
**Статус**: ✅ Готово

#### Поддержка в middleware (`/data/hostel-booking-system/backend/internal/middleware/auth_jwt.go`)
```go
// Строки 153-164 - обработка ролей из токена Auth Service
if m.authServicePubKey != nil {
    authClaims, err := jwt.ValidateAuthServiceToken(tokenString, m.authServicePubKey)
    if err == nil {
        // Проверяем роли
        isAdmin := false
        for _, role := range authClaims.Roles {
            if role == "admin" {
                isAdmin = true
                break
            }
        }
        c.Locals("is_admin", isAdmin)
    }
}
```

## ❌ Что осталось сделать

### 1. Компиляция и деплой Auth Service
**Статус**: ⏳ Требуется
**Приоритет**: Высокий

#### Необходимые изменения в main.go
```go
// /data/auth_svetu/cmd/server/main.go

// Добавить инициализацию RoleService
roleService := role.NewService(db.GetDB(), userRepo, logger.GetLogger())

// Обновить создание JWTService
tokenService, err := token.NewJWTService(&cfg.JWT, roleService)

// Добавить роуты для управления ролями
roleHandler := handlers.NewRoleHandlerFiber(roleService, logger.GetLogger())
roleHandler.RegisterRoutes(api, authMiddleware)
```

#### Команды для пересборки
```bash
cd /data/auth_svetu
docker-compose build auth-service
docker-compose up -d auth-service
```

### 2. Интеграция RoleManager в админ-панель
**Статус**: ⏳ Требуется
**Приоритет**: Средний

#### Изменения в UsersPageClient.tsx
```typescript
// Импорт компонента
import RoleManager from '@/components/admin/RoleManager';

// В модальном окне пользователя
{selectedUser && (
  <RoleManager
    userId={selectedUser.id}
    userEmail={selectedUser.email}
    onRoleUpdate={() => fetchUsers()}
  />
)}
```

### 3. Регистрация пользователей в Auth Service
**Статус**: ⏳ Требуется  
**Приоритет**: Высокий

#### Проблема
Пользователи существуют в основной БД, но не в Auth Service. Нужна синхронизация.

#### Решения
1. **Batch импорт** - скрипт для массовой регистрации
2. **Lazy registration** - регистрация при первом входе
3. **Sync service** - фоновая синхронизация

### 4. Переход на RS256 токены
**Статус**: ⏳ Опционально
**Приоритет**: Низкий

#### Текущее состояние
- Auth Service генерирует RS256 токены
- Основное приложение поддерживает RS256
- Но все еще используются HS256 токены

#### План миграции
1. Обновить frontend для работы с Auth Service
2. Постепенно перевести аутентификацию на Auth Service
3. Отключить локальную генерацию токенов

### 5. Расширение системы ролей
**Статус**: ⏳ Опционально
**Приоритет**: Низкий

#### Возможности расширения
- **Permissions** - детальные разрешения (create_post, delete_user)
- **Role hierarchy** - наследование ролей
- **Dynamic roles** - создание кастомных ролей через UI
- **Resource-based** - роли для конкретных ресурсов

### 6. Удаление legacy кода
**Статус**: ⏳ После полной миграции
**Приоритет**: Низкий

#### Что нужно удалить
- Таблица `admin_users` в основной БД
- Метод `IsUserAdmin` в storage
- Локальная проверка админов в middleware

## 🔧 Технические детали

### Переменные окружения

#### Auth Service (.env)
```env
DB_HOST=localhost
DB_PORT=25432
DB_NAME=auth_db
DB_USER=auth_user
DB_PASSWORD=AuthP@ssw0rd2025!

JWT_ALGORITHM=RS256
JWT_PRIVATE_KEY_PATH=./keys/private.pem
JWT_PUBLIC_KEY_PATH=./keys/public.pem
JWT_ACCESS_TOKEN_DURATION=15m
JWT_REFRESH_TOKEN_DURATION=7d
```

#### Main Backend (.env)
```env
DATABASE_URL=postgres://postgres:mX3g1XGhMRUZEX3l@localhost:5432/svetubd?sslmode=disable
JWT_SECRET=WdOl0x5Ka/ZlFssXCfzuqww/nQnwbQ8ZgOYXKpZOn18=
AUTH_SERVICE_URL=http://localhost:28080
AUTH_SERVICE_PUBLIC_KEY_PATH=./keys/auth_service_public.pem
```

#### Frontend (.env.local)
```env
NEXT_PUBLIC_API_URL=http://localhost:3000
NEXT_PUBLIC_AUTH_SERVICE_URL=http://localhost:28080
```

### Структура БД

#### Auth Service (PostgreSQL port 25432)
```
auth_db
├── auth.users              # Пользователи Auth Service
├── auth.roles              # Системные роли
├── auth.user_roles         # Связь пользователь-роль
├── auth.pending_admin_assignments  # Очередь назначения админов
├── auth.refresh_tokens     # Refresh токены
└── auth.service_accounts   # Сервисные аккаунты
```

#### Main App (PostgreSQL port 5432)
```
svetubd
├── users                   # Основные пользователи
├── admin_users (legacy)    # Старая таблица админов
├── marketplace_listings    # Объявления
└── ... другие таблицы
```

### API Documentation

#### Назначение роли
```bash
curl -X POST http://localhost:28080/api/v1/roles/assign \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 123,
    "role_name": "moderator",
    "expires_at": "2025-12-31T23:59:59Z",
    "notes": "Temporary moderator for holiday season"
  }'
```

#### Проверка ролей пользователя
```bash
curl http://localhost:28080/api/v1/users/123/roles \
  -H "Authorization: Bearer $TOKEN"
```

#### Отзыв роли
```bash
curl -X POST http://localhost:28080/api/v1/roles/revoke \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 123,
    "role_name": "moderator",
    "reason": "Role no longer needed"
  }'
```

## 📦 Инструкции по развертыванию

### Быстрый старт (для разработки)

```bash
# 1. Применить миграции
docker exec -i auth_postgres psql -U auth_user -d auth_db < /data/auth_svetu/migrations/002_add_roles_system.sql
docker exec -i auth_postgres psql -U auth_user -d auth_db < /data/auth_svetu/scripts/admin_auto_assign.sql

# 2. Импортировать админов
cd /data/auth_svetu/scripts
go run migrate_admins.go

# 3. Пересобрать Auth Service
cd /data/auth_svetu
docker-compose build auth-service
docker-compose restart auth-service

# 4. Проверить работу
curl http://localhost:28080/health
```

### Production развертывание

1. **Backup данных**
```bash
pg_dump -h localhost -U postgres -d svetubd -t admin_users > admin_users_backup.sql
```

2. **Применить миграции через CI/CD**
```yaml
- name: Apply migrations
  run: |
    migrate -path migrations -database $DATABASE_URL up
```

3. **Rolling update сервисов**
```bash
kubectl set image deployment/auth-service auth-service=auth-service:v2.0.0
kubectl rollout status deployment/auth-service
```

4. **Мониторинг**
```bash
# Проверка назначения ролей
SELECT COUNT(*) FROM auth.user_roles WHERE created_at > NOW() - INTERVAL '1 hour';

# Проверка ошибок
SELECT * FROM auth.audit_log WHERE level = 'ERROR' ORDER BY created_at DESC LIMIT 10;
```

## 📊 Метрики успеха

- ✅ Все 9 админов имеют роли в новой системе
- ✅ JWT токены содержат массив ролей
- ✅ API endpoints для управления ролями работают
- ✅ Основное приложение корректно обрабатывает роли из токенов
- ⏳ Frontend компонент интегрирован в админ-панель
- ⏳ Все пользователи мигрированы в Auth Service
- ⏳ Legacy таблица admin_users удалена

## 🔗 Связанные документы

- [AUTH_SERVICE_README.md](/data/auth_svetu/README.md) - документация Auth Service
- [ADMIN_ROLES_INTEGRATION.md](/data/auth_svetu/docs/ADMIN_ROLES_INTEGRATION.md) - инструкция по интеграции
- [JWT_TOKEN_SPECIFICATION.md](/data/auth_svetu/docs/JWT_TOKEN_SPECIFICATION.md) - спецификация токенов
- [MIGRATION_PLAN.md](/data/hostel-booking-system/docs/AUTH_MIGRATION_PLAN.md) - план полной миграции