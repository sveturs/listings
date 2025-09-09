# 📋 План полной миграции на Auth Service микросервис

## 🔴 Текущие проблемы

### 1. **Не передается `is_admin` флаг в JWT токен**
- Auth Service не включает поле `is_admin` в ответ `/api/v1/auth/session`
- Frontend ожидает `user.is_admin`, но получает `null`
- Админские меню не отображаются в UI

### 2. **Двойная система авторизации**
- Старая: таблица `admin_users` в основной БД
- Новая: таблица `auth.user_roles` в Auth Service БД
- Несинхронизированные данные между системами

## 📊 Детальный план миграции

### Фаза 1: Исправление критических проблем (1-2 дня)

#### 1.1 Добавить `is_admin` в Auth Service
```go
// В Auth Service добавить в UserResponse:
type UserResponse struct {
    ID       int      `json:"id"`
    Email    string   `json:"email"`
    Name     string   `json:"name"`
    Roles    []string `json:"roles"`
    IsAdmin  bool     `json:"is_admin"` // NEW!
}

// При формировании ответа:
response.IsAdmin = containsRole(user.Roles, "admin")
```

#### 1.2 Синхронизировать админов между системами
```sql
-- Миграция для переноса админов из admin_users в Auth Service
INSERT INTO auth.user_roles (user_id, role_id, granted_at, is_active)
SELECT 
    u.id as user_id,
    r.id as role_id,
    NOW() as granted_at,
    true as is_active
FROM auth.users u
JOIN public.admin_users au ON u.email = au.email
JOIN auth.roles r ON r.name = 'admin'
WHERE NOT EXISTS (
    SELECT 1 FROM auth.user_roles ur 
    WHERE ur.user_id = u.id AND ur.role_id = r.id
);
```

### Фаза 2: Унификация системы авторизации (3-5 дней)

#### 2.1 Обновить backend middleware
- [ ] Удалить проверку `IsUserAdmin()` из таблицы `admin_users`
- [ ] Полагаться только на JWT claims от Auth Service
- [ ] Обновить все эндпоинты использующие `utils.IsAdmin()`

#### 2.2 Создать единый источник правды
```go
// backend/internal/middleware/auth_jwt.go
// Убрать двойную проверку:
isAdmin := false
for _, role := range authClaims.Roles {
    if role == "admin" {
        isAdmin = true
        break
    }
}
// УДАЛИТЬ эту строку:
// isAdmin, _ = m.services.User().IsUserAdmin(c.Context(), authClaims.Email)
```

### Фаза 3: Миграция функциональности (5-7 дней)

#### 3.1 Перенести управление админами в Auth Service
- [ ] API для добавления/удаления админов
- [ ] UI для управления ролями
- [ ] Аудит лог изменений ролей

#### 3.2 Обновить Frontend
- [ ] Добавить страницу управления ролями `/admin/roles`
- [ ] Интегрировать с Auth Service API
- [ ] Добавить индикатор роли в профиле

### Фаза 4: Расширенная авторизация (7-10 дней)

#### 4.1 Добавить систему разрешений (Permissions)
```sql
-- Новые таблицы в Auth Service
CREATE TABLE auth.permissions (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    resource VARCHAR(100) NOT NULL,
    action VARCHAR(50) NOT NULL
);

CREATE TABLE auth.role_permissions (
    role_id INT REFERENCES auth.roles(id),
    permission_id INT REFERENCES auth.permissions(id),
    PRIMARY KEY (role_id, permission_id)
);
```

#### 4.2 Implement RBAC (Role-Based Access Control)
- [ ] Создать permissions: `listings.delete`, `users.ban`, `payments.refund`
- [ ] Привязать permissions к ролям
- [ ] Проверять permissions в middleware

### Фаза 5: Очистка и оптимизация (2-3 дня)

#### 5.1 Удалить legacy код
- [ ] Удалить таблицу `admin_users` (после проверки)
- [ ] Удалить `IsUserAdmin()` методы
- [ ] Удалить дублирующие проверки

#### 5.2 Документация
- [ ] API документация для Auth Service
- [ ] Руководство по управлению ролями
- [ ] Migration guide для разработчиков

## 🚀 Быстрое решение (HOTFIX)

Для немедленного исправления админских прав:

### Вариант 1: Патч в backend (временное решение)
```go
// backend/internal/middleware/auth_jwt.go
// После строки 187 добавить:

// ВРЕМЕННЫЙ ХОТФИКС: всегда проверяем admin_users
if !isAdmin {
    isAdmin, _ = m.services.User().IsUserAdmin(c.Context(), authClaims.Email)
}

// И передаем флаг в user данные
if user != nil {
    user.IsAdmin = isAdmin // Добавить это поле в модель User
}
```

### Вариант 2: Обновить Auth Service (правильное решение)
1. Добавить поле `is_admin` в JWT claims
2. Включать его в `/api/v1/auth/session` response
3. Обновить Frontend типы для поддержки `is_admin`

## 📅 Timeline

| Фаза | Срок | Приоритет | Статус |
|------|------|-----------|---------|
| Фаза 1 | 1-2 дня | 🔴 Критический | ⏳ Не начато |
| Фаза 2 | 3-5 дней | 🟠 Высокий | ⏳ Не начато |
| Фаза 3 | 5-7 дней | 🟡 Средний | ⏳ Не начато |
| Фаза 4 | 7-10 дней | 🟢 Низкий | ⏳ Не начато |
| Фаза 5 | 2-3 дня | 🔵 Очистка | ⏳ Не начато |

## 🎯 Критерии успеха

1. ✅ Админские функции отображаются в UI
2. ✅ Единый источник правды для ролей (Auth Service)
3. ✅ Нет дублирования логики авторизации
4. ✅ Поддержка RBAC с permissions
5. ✅ Полная документация и тесты

## 🔧 Команда для проверки текущего статуса

```bash
# Проверить роли в Auth Service
docker exec auth_postgres psql -U auth_user -d auth_db -c "
SELECT u.email, r.name as role 
FROM auth.users u 
JOIN auth.user_roles ur ON u.id = ur.user_id 
JOIN auth.roles r ON ur.role_id = r.id 
WHERE u.email = 'boxmail386@gmail.com';"

# Проверить админов в основной БД
PGPASSWORD=mX3g1XGhMRUZEX3l psql -h localhost -U postgres -d svetubd -c "
SELECT * FROM admin_users WHERE email = 'boxmail386@gmail.com';"

# Проверить что возвращает Auth Service
curl -s http://localhost:28080/api/v1/auth/session \
  -H "Authorization: Bearer YOUR_TOKEN" | jq '.user'
```

## 📝 Заметки

- Auth Service создан 8 сентября 2025 (очень свежий)
- Необходима координация между backend и Auth Service командами
- Критично сохранить обратную совместимость во время миграции