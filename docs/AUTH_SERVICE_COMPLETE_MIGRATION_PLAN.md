# План полной миграции на централизованную систему авторизации через Auth Service

## 📋 Оглавление
- [Принципы миграции](#принципы-миграции)
- [Фаза 1: Полный переход на RS256](#фаза-1-полный-переход-на-rs256)
- [Фаза 2: Миграция всех пользователей](#фаза-2-миграция-всех-пользователей)
- [Фаза 3: Удаление legacy кода](#фаза-3-удаление-legacy-кода)
- [Фаза 4: Расширение системы прав](#фаза-4-расширение-системы-прав)
- [Чек-лист завершения](#чек-лист-завершения)

## 🎯 Принципы миграции

### Главный принцип: **Никакого технического долга!**

- ❌ НЕ оставляем временные решения
- ❌ НЕ поддерживаем две системы параллельно
- ❌ НЕ откладываем рефакторинг "на потом"
- ✅ Делаем сразу правильно
- ✅ Удаляем старый код немедленно
- ✅ Документируем все изменения

## 🔄 Фаза 1: Полный переход на RS256
**Срок**: 1-2 дня  
**Downtime**: 1-2 часа (планируемый)

### 1.1 Генерация ключей и настройка Auth Service

```bash
# Генерируем производственные ключи RSA 4096 bit
cd /data/auth_svetu
mkdir -p keys
openssl genrsa -out keys/private.pem 4096
openssl rsa -in keys/private.pem -pubout -out keys/public.pem

# Обновляем конфигурацию Auth Service
cat >> .env << EOF
JWT_ALGORITHM=RS256
JWT_PRIVATE_KEY_PATH=./keys/private.pem
JWT_PUBLIC_KEY_PATH=./keys/public.pem
EOF
```

### 1.2 Полная компиляция Auth Service с Role Service

```go
// /data/auth_svetu/cmd/server/main.go

import (
    "github.com/svetu/auth-service/internal/service/role"
    // ...
)

func main() {
    // ...
    
    // Инициализация Role Service
    roleService := role.NewService(db.GetDB(), userRepo, logger.GetLogger())
    
    // JWT Service с поддержкой ролей
    tokenService, err := token.NewJWTService(&cfg.JWT, roleService)
    if err != nil {
        logger.Fatal("Failed to initialize JWT service", map[string]interface{}{"error": err.Error()})
    }
    
    // HTTP handlers для ролей
    roleHandler := handlers.NewRoleHandlerFiber(roleService, logger.GetLogger())
    
    // Регистрация роутов
    api := app.Group("/api/v1")
    roleHandler.RegisterRoutes(api, middleware.AuthMiddleware(authService))
    
    // ...
}
```

### 1.3 Отключение HS256 в основном приложении

```go
// /data/hostel-booking-system/backend/internal/middleware/auth_jwt.go

func (m *Middleware) AuthRequiredJWT(c *fiber.Ctx) error {
    // УДАЛЯЕМ всю логику HS256
    // ОСТАВЛЯЕМ только RS256 валидацию
    
    authHeader := c.Get("Authorization")
    if authHeader == "" {
        return utils.ErrorResponse(c, fiber.StatusUnauthorized, "users.auth.error.authentication_required")
    }
    
    parts := strings.Split(authHeader, " ")
    if len(parts) != 2 || parts[0] != "Bearer" {
        return utils.ErrorResponse(c, fiber.StatusUnauthorized, "users.auth.error.invalid_token")
    }
    
    // ТОЛЬКО RS256 проверка
    authClaims, err := jwt.ValidateAuthServiceToken(parts[1], m.authServicePubKey)
    if err != nil {
        return utils.ErrorResponse(c, fiber.StatusUnauthorized, "users.auth.error.invalid_token")
    }
    
    // Устанавливаем данные пользователя
    c.Locals("user_id", authClaims.UserID)
    c.Locals("user_email", authClaims.Email)
    c.Locals("roles", authClaims.Roles)
    
    // Проверяем админа
    isAdmin := false
    for _, role := range authClaims.Roles {
        if role == "admin" {
            isAdmin = true
            break
        }
    }
    c.Locals("is_admin", isAdmin)
    
    return c.Next()
}
```

### 1.4 Удаление JWT_SECRET из всех .env файлов

```bash
# Backend
sed -i '/JWT_SECRET=/d' /data/hostel-booking-system/backend/.env

# Добавляем путь к публичному ключу
echo "AUTH_SERVICE_PUBLIC_KEY_PATH=/data/auth_svetu/keys/public.pem" >> /data/hostel-booking-system/backend/.env

# Frontend - удаляем все упоминания JWT
sed -i '/JWT/d' /data/hostel-booking-system/frontend/svetu/.env.local

# Добавляем только Auth Service URL
echo "NEXT_PUBLIC_AUTH_SERVICE_URL=http://localhost:28080" >> /data/hostel-booking-system/frontend/svetu/.env.local
```

### 1.5 Удаление локальной генерации токенов

```bash
# Удаляем все скрипты генерации токенов
rm /data/hostel-booking-system/backend/scripts/create_test_jwt.go
rm /data/hostel-booking-system/backend/scripts/create_admin_jwt.go

# Удаляем пакет JWT из backend (оставляем только валидацию RS256)
rm -rf /data/hostel-booking-system/backend/pkg/jwt/generate.go
```

## 🚀 Фаза 2: Миграция всех пользователей
**Срок**: 1 день  
**Downtime**: 30 минут

### 2.1 Полный дамп пользователей из основной БД

```sql
-- Создаем полный экспорт пользователей
COPY (
    SELECT 
        u.id,
        u.email,
        u.name,
        u.phone,
        u.created_at,
        u.google_id,
        u.facebook_id,
        u.account_status,
        CASE 
            WHEN au.email IS NOT NULL THEN 'admin'
            ELSE 'user'
        END as role
    FROM users u
    LEFT JOIN admin_users au ON u.email = au.email
) TO '/tmp/users_export.csv' WITH CSV HEADER;
```

### 2.2 Массовый импорт в Auth Service

```go
// /data/auth_svetu/scripts/mass_import_users.go

package main

import (
    "encoding/csv"
    "os"
    "database/sql"
    "time"
)

func main() {
    // Читаем CSV
    file, _ := os.Open("/tmp/users_export.csv")
    defer file.Close()
    
    reader := csv.NewReader(file)
    records, _ := reader.ReadAll()
    
    // Подключаемся к Auth Service DB
    authDB, _ := sql.Open("postgres", "postgres://auth_user:AuthP@ssw0rd2025!@localhost:25432/auth_db")
    
    // Batch insert
    tx, _ := authDB.Begin()
    
    for _, record := range records[1:] { // Skip header
        // Вставляем пользователя
        var userID int
        err := tx.QueryRow(`
            INSERT INTO auth.users (email, name, provider, created_at)
            VALUES ($1, $2, 'migrated', $3)
            ON CONFLICT (email) DO UPDATE SET name = EXCLUDED.name
            RETURNING id
        `, record[1], record[2], record[4]).Scan(&userID)
        
        if err == nil && record[9] == "admin" {
            // Назначаем роль админа
            tx.Exec(`
                INSERT INTO auth.user_roles (user_id, role_id, granted_at, is_active, notes)
                SELECT $1, r.id, NOW(), true, 'Migrated from main DB'
                FROM auth.roles r WHERE r.name = 'admin'
                ON CONFLICT DO NOTHING
            `, userID)
        }
    }
    
    tx.Commit()
    
    println("Migration completed!")
}
```

### 2.3 Переключение Frontend на Auth Service

```typescript
// /data/hostel-booking-system/frontend/svetu/src/services/auth.ts

class AuthService {
    // УДАЛЯЕМ все методы работы с backend auth
    // ТОЛЬКО Auth Service
    
    async login(email: string, password: string) {
        const response = await fetch(`${config.getAuthServiceUrl()}/api/v1/auth/login`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ email, password })
        });
        
        const data = await response.json();
        if (data.access_token) {
            tokenManager.setTokens(data.access_token, data.refresh_token);
        }
        return data;
    }
    
    async register(userData: RegisterData) {
        // Только через Auth Service
        return fetch(`${config.getAuthServiceUrl()}/api/v1/auth/register`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(userData)
        });
    }
    
    // Удаляем методы:
    // - loginWithBackend()
    // - registerWithBackend()
    // - validateWithBackend()
}
```

### 2.4 Отключение auth endpoints в Backend

```go
// /data/hostel-booking-system/backend/internal/server/routes.go

func SetupRoutes(app *fiber.App) {
    // УДАЛЯЕМ все auth endpoints
    // ❌ api.Post("/auth/login", handlers.Login)
    // ❌ api.Post("/auth/register", handlers.Register)
    // ❌ api.Post("/auth/refresh", handlers.RefreshToken)
    // ❌ api.Get("/auth/validate", handlers.ValidateToken)
    
    // Оставляем только бизнес-логику
    api := app.Group("/api/v1")
    
    // Marketplace
    api.Get("/marketplace/listings", middleware.AuthRequiredJWT, handlers.GetListings)
    api.Post("/marketplace/listings", middleware.AuthRequiredJWT, handlers.CreateListing)
    
    // Admin (проверка ролей через токен)
    admin := api.Group("/admin", middleware.AuthRequiredJWT, middleware.RequireAdmin)
    admin.Get("/users", handlers.AdminGetUsers)
    admin.Put("/users/:id/status", handlers.AdminUpdateUserStatus)
}
```

## 🗑️ Фаза 3: Удаление legacy кода
**Срок**: 1 день  
**Downtime**: 0

### 3.1 Удаление таблицы admin_users

```sql
-- Сначала сохраняем backup
CREATE TABLE admin_users_backup_2025_01_09 AS SELECT * FROM admin_users;

-- Удаляем все зависимости
ALTER TABLE admin_users DROP CONSTRAINT IF EXISTS admin_users_created_by_fkey;
DROP INDEX IF EXISTS idx_admin_users_email;

-- Удаляем таблицу
DROP TABLE admin_users;

-- Удаляем из миграций упоминания
DELETE FROM schema_migrations WHERE version IN (
    SELECT version FROM schema_migrations 
    WHERE version LIKE '%admin_users%'
);
```

### 3.2 Удаление кода проверки админов

```go
// Удаляем файлы:
rm /data/hostel-booking-system/backend/internal/proj/users/storage/postgres/admin_users.go
rm /data/hostel-booking-system/backend/internal/storage/postgres/admin_methods.go

// Удаляем методы из интерфейсов:
// - IsUserAdmin()
// - GetAllAdmins()
// - AddAdmin()
// - RemoveAdmin()
```

### 3.3 Удаление сессий и cookie auth

```go
// Удаляем из middleware:
// - Проверку session_token cookie
// - Fallback на сессии
// - Генерацию JWT из сессий

// Удаляем таблицы:
DROP TABLE user_sessions;
DROP TABLE session_tokens;
```

### 3.4 Очистка Frontend

```typescript
// Удаляем:
// - Компоненты локальной авторизации
// - Хуки для session-based auth
// - Cookie менеджеры для JWT
// - Fallback логику

// Оставляем только:
// - Работу с Auth Service
// - Token manager для access/refresh токенов
```

## 🚀 Фаза 4: Расширение системы прав
**Срок**: 3-5 дней  
**Downtime**: 0

### 4.1 Добавление детальных permissions

```sql
-- Таблица разрешений
CREATE TABLE auth.permissions (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    resource VARCHAR(50) NOT NULL,
    action VARCHAR(50) NOT NULL,
    description TEXT
);

-- Связь роль-разрешение
CREATE TABLE auth.role_permissions (
    role_id INTEGER REFERENCES auth.roles(id),
    permission_id INTEGER REFERENCES auth.permissions(id),
    PRIMARY KEY (role_id, permission_id)
);

-- Базовые разрешения
INSERT INTO auth.permissions (name, resource, action) VALUES
    ('users.read', 'users', 'read'),
    ('users.write', 'users', 'write'),
    ('users.delete', 'users', 'delete'),
    ('listings.moderate', 'listings', 'moderate'),
    ('payments.manage', 'payments', 'manage'),
    ('analytics.view', 'analytics', 'view');

-- Назначаем разрешения ролям
INSERT INTO auth.role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM auth.roles r
CROSS JOIN auth.permissions p
WHERE r.name = 'admin'; -- Админ получает все

INSERT INTO auth.role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM auth.roles r
CROSS JOIN auth.permissions p
WHERE r.name = 'moderator' 
  AND p.name IN ('users.read', 'listings.moderate');
```

### 4.2 Добавление resource-based permissions

```go
// Проверка прав на конкретный ресурс
type ResourcePermission struct {
    UserID     int
    Resource   string  // "listing:123"
    Permission string  // "edit"
}

func CheckResourcePermission(ctx context.Context, req ResourcePermission) bool {
    // Проверяем:
    // 1. Общие права роли
    // 2. Специфичные права на ресурс
    // 3. Ownership (владелец ресурса)
}
```

### 4.3 Интеграция с Frontend

```typescript
// Компонент проверки прав
function Can({ 
    permission, 
    resource, 
    children, 
    fallback = null 
}: CanProps) {
    const { user, permissions } = useAuth();
    
    const hasPermission = checkPermission(user, permission, resource);
    
    return hasPermission ? children : fallback;
}

// Использование
<Can permission="listings.moderate">
    <button onClick={moderateListing}>Модерировать</button>
</Can>

<Can permission="users.delete" resource={`user:${userId}`}>
    <button onClick={deleteUser}>Удалить пользователя</button>
</Can>
```

### 4.4 Аудит и логирование

```sql
-- Таблица аудита действий
CREATE TABLE auth.audit_log (
    id BIGSERIAL PRIMARY KEY,
    user_id INTEGER,
    action VARCHAR(100),
    resource VARCHAR(200),
    details JSONB,
    ip_address INET,
    user_agent TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Индексы для быстрого поиска
CREATE INDEX idx_audit_user_id ON auth.audit_log(user_id);
CREATE INDEX idx_audit_action ON auth.audit_log(action);
CREATE INDEX idx_audit_created_at ON auth.audit_log(created_at);
```

## ✅ Чек-лист завершения

### Фаза 1 - RS256
- [ ] Ключи RSA 4096 сгенерированы
- [ ] Auth Service использует только RS256
- [ ] Backend проверяет только RS256 токены
- [ ] JWT_SECRET удален из всех .env
- [ ] Локальная генерация токенов удалена
- [ ] Все сервисы перезапущены

### Фаза 2 - Миграция пользователей
- [ ] Все пользователи экспортированы
- [ ] Пользователи импортированы в Auth Service
- [ ] Роли назначены корректно
- [ ] Frontend использует только Auth Service
- [ ] Backend auth endpoints удалены
- [ ] Тесты обновлены

### Фаза 3 - Удаление legacy
- [ ] Таблица admin_users удалена
- [ ] Код проверки админов удален
- [ ] Session auth удален
- [ ] Cookie auth удален
- [ ] Старые миграции очищены
- [ ] Документация обновлена

### Фаза 4 - Расширение
- [ ] Permissions система добавлена
- [ ] Resource-based проверки работают
- [ ] Frontend компоненты обновлены
- [ ] Аудит логирование настроено
- [ ] Monitoring и alerting настроены
- [ ] Load testing проведен

## 📊 Метрики успеха после миграции

1. **Производительность**
   - Время валидации токена < 10ms
   - Время генерации токена < 50ms
   - Zero downtime при ротации ключей

2. **Безопасность**
   - 0 секретных ключей в коде
   - Все токены RS256
   - Аудит лог 100% действий

3. **Чистота кода**
   - 0 legacy auth кода
   - 0 TODO комментариев
   - 100% покрытие тестами

4. **Масштабируемость**
   - Auth Service может обрабатывать 10k RPS
   - Горизонтальное масштабирование работает
   - Cache слой для токенов

## 🎉 Результат

После выполнения всех фаз:
- **Единая** система авторизации
- **Безопасная** RS256 криптография  
- **Чистый** код без legacy
- **Расширяемая** система прав
- **Готовность** к production

**Никакого технического долга!** 🚀