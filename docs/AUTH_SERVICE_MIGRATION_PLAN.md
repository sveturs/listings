# 🚀 План миграции системы авторизации на микросервисную архитектуру

**Дата создания:** 06.09.2025  
**Статус:** К исполнению  
**Время выполнения:** 5 недель  
**Подход:** Без обратной совместимости (проект не в продакшене)

## 📋 Текущее состояние

### Проблемы существующей системы

1. **Избыточность механизмов**
   - Session tokens (legacy) работают параллельно с JWT
   - 4 места проверки аутентификации в middleware
   - Дублирование логики между компонентами

2. **Архитектурные недостатки**
   - Session tokens в sync.Map (не масштабируется)
   - Отсутствие кеширования
   - Монолитная структура
   - Путаница с provider field (google/email/password/jwt)

3. **Производительность**
   - Частые запросы к БД без кеша
   - Последовательные проверки множества источников токенов
   - Отсутствие connection pooling для auth запросов

## 🎯 Целевая архитектура

### Принципы
- **JWT-only** - никаких session tokens
- **Микросервис** - полностью независимый auth-service
- **Stateless** для access tokens, **Stateful** для refresh tokens
- **Кеширование** через Redis
- **gRPC** для внутренних вызовов, **REST** для внешних

### Технологический стек
- **Язык:** Go 1.22+
- **База данных:** PostgreSQL 15 (основная) + Redis 7 (кеш)
- **Протоколы:** gRPC (внутренний), REST (внешний)
- **Токены:** JWT RS256
- **Контейнеризация:** Docker + Kubernetes
- **Мониторинг:** Prometheus + Grafana
- **Трассировка:** OpenTelemetry + Jaeger

### Структура микросервиса

```
auth-service/
├── cmd/
│   ├── grpc/          # gRPC сервер
│   └── http/          # REST API gateway
├── internal/
│   ├── domain/        # Доменные модели
│   ├── service/       # Бизнес-логика
│   │   ├── auth/      # Аутентификация
│   │   ├── token/     # Управление токенами
│   │   ├── user/      # Управление пользователями
│   │   └── oauth/     # OAuth провайдеры
│   ├── repository/    # Работа с данными
│   │   ├── postgres/  # PostgreSQL репозиторий
│   │   └── redis/     # Redis кеш
│   ├── transport/     # API слой
│   │   ├── grpc/      # gRPC handlers
│   │   └── http/      # REST handlers
│   └── middleware/    # Middleware компоненты
├── pkg/
│   ├── jwt/          # JWT утилиты
│   ├── crypto/       # Криптография
│   └── validator/    # Валидация
└── migrations/       # SQL миграции
```

## 📅 План реализации

### Неделя 1: Очистка текущего кода

#### День 1-2: Аудит и документация
- [ ] Составить полный список всех мест использования session_token
- [ ] Документировать все API endpoints связанные с auth
- [ ] Создать список всех зависимостей от auth в других модулях
- [ ] Подготовить тестовые данные

#### День 3-4: Удаление legacy кода
- [ ] Удалить всю логику session tokens из backend
- [ ] Упростить AuthMiddleware - оставить только JWT проверку
- [ ] Удалить session_token из frontend
- [ ] Очистить базу данных от неиспользуемых полей

#### День 5: Унификация
- [ ] Изменить provider field: оставить только "google" и "local"
- [ ] Обновить все references на новые значения
- [ ] Создать миграцию для обновления существующих данных
- [ ] Провести тестирование

### Неделя 2: Разработка auth-service

#### День 1-2: Инициализация проекта
```bash
# Структура репозитория
mkdir -p auth-service/{cmd,internal,pkg,migrations,deployments,scripts}

# Основные компоненты
- cmd/grpc/main.go         # gRPC сервер
- cmd/http/main.go         # REST gateway
- internal/config/         # Конфигурация
- internal/domain/         # Модели User, Token, Session
- pkg/jwt/                # JWT генерация и валидация
```

#### День 3-4: Реализация core функционала

**Основные сервисы:**
```go
// AuthService
- Login(email, password) -> (accessToken, refreshToken)
- LoginWithGoogle(code) -> (accessToken, refreshToken)
- Register(email, password, name) -> (accessToken, refreshToken)
- Logout(refreshToken) -> error
- RefreshTokens(refreshToken) -> (newAccessToken, newRefreshToken)

// TokenService
- GenerateAccessToken(userID, email) -> token
- GenerateRefreshToken(userID) -> token
- ValidateAccessToken(token) -> claims
- RevokeRefreshToken(token) -> error
- RevokeAllUserTokens(userID) -> error

// UserService
- GetUserByID(id) -> user
- GetUserByEmail(email) -> user
- UpdateUser(user) -> error
- CheckPassword(user, password) -> bool
```

#### День 5: База данных и кеширование

**PostgreSQL схема:**
```sql
-- Таблица пользователей
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    password_hash VARCHAR(255),
    google_id VARCHAR(255),
    provider VARCHAR(20) NOT NULL CHECK (provider IN ('local', 'google')),
    picture_url TEXT,
    email_verified BOOLEAN DEFAULT FALSE,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Таблица refresh токенов
CREATE TABLE refresh_tokens (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id),
    token_hash VARCHAR(255) UNIQUE NOT NULL,
    family_id UUID NOT NULL, -- для rotation detection
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    ip_address INET,
    user_agent TEXT,
    is_revoked BOOLEAN DEFAULT FALSE,
    revoked_at TIMESTAMP,
    INDEX idx_user_tokens (user_id, is_revoked),
    INDEX idx_token_expires (expires_at)
);

-- Таблица для отзыва access токенов (blacklist)
CREATE TABLE revoked_access_tokens (
    jti VARCHAR(255) PRIMARY KEY, -- JWT ID
    expires_at TIMESTAMP NOT NULL,
    revoked_at TIMESTAMP DEFAULT NOW(),
    INDEX idx_expires (expires_at)
);
```

**Redis структура:**
```redis
# Кеш пользователей
user:{id} -> JSON user data (TTL: 5 min)

# Rate limiting
rate:login:{email} -> counter (TTL: 15 min)
rate:refresh:{user_id} -> counter (TTL: 1 min)

# Активные сессии пользователя
sessions:{user_id} -> SET of refresh_token_ids

# Blacklist для access tokens (если нужен instant revoke)
blacklist:{jti} -> 1 (TTL: до expires_at)
```

### Неделя 3: API и интеграции

#### День 1-2: gRPC API

**Proto определения:**
```proto
service AuthService {
    rpc Login(LoginRequest) returns (AuthResponse);
    rpc Register(RegisterRequest) returns (AuthResponse);
    rpc RefreshToken(RefreshTokenRequest) returns (AuthResponse);
    rpc ValidateToken(ValidateTokenRequest) returns (ValidateTokenResponse);
    rpc Logout(LogoutRequest) returns (Empty);
    rpc RevokeAllTokens(RevokeAllTokensRequest) returns (Empty);
}

message AuthResponse {
    string access_token = 1;
    string refresh_token = 2;
    int32 expires_in = 3;
    User user = 4;
}

message ValidateTokenResponse {
    bool valid = 1;
    int32 user_id = 2;
    string email = 3;
    repeated string roles = 4;
}
```

#### День 3-4: REST API Gateway

**Endpoints:**
```yaml
POST   /api/v1/auth/register
POST   /api/v1/auth/login
POST   /api/v1/auth/logout
POST   /api/v1/auth/refresh
GET    /api/v1/auth/validate
POST   /api/v1/auth/google
GET    /api/v1/auth/google/callback
POST   /api/v1/auth/forgot-password
POST   /api/v1/auth/reset-password
GET    /api/v1/auth/verify-email
```

#### День 5: OAuth интеграция

**Google OAuth flow:**
```go
// 1. Инициация
GET /auth/google
-> Redirect to Google OAuth

// 2. Callback
GET /auth/google/callback?code=...
-> Exchange code for Google token
-> Get user info from Google
-> Create/update user in DB
-> Generate JWT tokens
-> Redirect to frontend with tokens
```

### Неделя 4: Безопасность и оптимизация

#### День 1-2: Безопасность

**Реализовать:**
- [ ] JWT с RS256 (асимметричные ключи)
- [ ] Refresh token rotation
- [ ] Device fingerprinting
- [ ] Обнаружение аномалий (множественные локации)
- [ ] Rate limiting через Redis
- [ ] CSRF защита для web endpoints
- [ ] Secure headers (HSTS, CSP, etc.)

**Защита от атак:**
```go
// Token rotation при обнаружении reuse
if tokenAlreadyUsed {
    RevokeTokenFamily(familyID) // Отзываем всю цепочку
    return ErrTokenReuse
}

// Rate limiting
if rateLimiter.Exceeded(email) {
    return ErrTooManyAttempts
}

// Suspicious activity detection
if DetectSuspiciousActivity(userID, ip, userAgent) {
    NotifyUser(userID)
    RequireMFA(userID)
}
```

#### День 3-4: Производительность

**Оптимизации:**
- [ ] Connection pooling для PostgreSQL
- [ ] Prepared statements для частых запросов
- [ ] Batch операции для массовых проверок
- [ ] Кеширование валидных токенов в Redis
- [ ] Graceful shutdown с завершением активных запросов

**Метрики производительности:**
```go
// Prometheus метрики
auth_login_duration_seconds
auth_token_validation_duration_seconds
auth_active_sessions_total
auth_failed_attempts_total
auth_token_refresh_total
```

#### День 5: Мониторинг и логирование

**Настроить:**
- [ ] Structured logging (JSON)
- [ ] Distributed tracing (OpenTelemetry)
- [ ] Метрики (Prometheus)
- [ ] Дашборды (Grafana)
- [ ] Алерты (критические ошибки)

### Неделя 5: Интеграция и миграция

#### День 1-2: Подготовка основного сервиса

**Backend изменения:**
```go
// Заменить internal auth на gRPC клиент
type AuthClient interface {
    ValidateToken(ctx context.Context, token string) (*User, error)
    RefreshToken(ctx context.Context, refreshToken string) (*Tokens, error)
}

// Новый middleware
func AuthMiddleware(authClient AuthClient) fiber.Handler {
    return func(c *fiber.Ctx) error {
        token := extractToken(c)
        user, err := authClient.ValidateToken(c.Context(), token)
        if err != nil {
            return c.Status(401).JSON(errorResponse)
        }
        c.Locals("user", user)
        return c.Next()
    }
}
```

#### День 3: Frontend изменения

**Обновить AuthService:**
```typescript
class AuthService {
  // Новые endpoints
  private readonly AUTH_API = process.env.NEXT_PUBLIC_AUTH_SERVICE_URL;
  
  async login(email: string, password: string) {
    const response = await fetch(`${this.AUTH_API}/api/v1/auth/login`, {
      method: 'POST',
      body: JSON.stringify({ email, password })
    });
    
    const data = await response.json();
    this.tokenManager.setTokens(data.access_token, data.refresh_token);
    return data.user;
  }
  
  // Удалить всё связанное с session_token
  // Упростить логику работы только с JWT
}
```

#### День 4: Миграция данных

**Скрипт миграции:**
```sql
-- 1. Обновить provider field
UPDATE users 
SET provider = 'local' 
WHERE provider IN ('email', 'password', 'jwt');

-- 2. Перенести активные refresh tokens
INSERT INTO auth_service.refresh_tokens 
SELECT * FROM main_db.refresh_tokens 
WHERE NOT is_revoked AND expires_at > NOW();

-- 3. Очистить старые данные
DROP TABLE IF EXISTS session_tokens;
ALTER TABLE users DROP COLUMN IF EXISTS session_data;
```

#### День 5: Тестирование и развертывание

**Чеклист развертывания:**
- [ ] Unit тесты (coverage > 80%)
- [ ] Integration тесты
- [ ] Load тесты (target: 10K RPS)
- [ ] Security scan
- [ ] Docker образы
- [ ] Kubernetes манифесты
- [ ] CI/CD pipeline
- [ ] Rollback план

## 🎯 Ожидаемые результаты

### Производительность
- **Latency:** < 50ms для валидации токена
- **Throughput:** 10,000+ RPS
- **Cache hit rate:** > 90%
- **Uptime:** 99.99%

### Безопасность
- Полное соответствие OWASP
- Защита от token replay attacks
- Automated anomaly detection
- Instant token revocation

### Масштабируемость
- Горизонтальное масштабирование без ограничений
- Stateless архитектура
- Независимое развертывание

### Разработка
- Чистая архитектура без legacy
- Простота поддержки
- Переиспользуемость для других проектов

## 📊 Метрики успеха

1. **Сокращение кодовой базы** на 30% за счет удаления legacy
2. **Увеличение производительности** на 200%
3. **Снижение времени на добавление новых auth методов** с дней до часов
4. **Полная независимость** auth логики от основного сервиса

## 🚦 Риски и митигация

| Риск | Вероятность | Влияние | Митигация |
|------|------------|---------|-----------|
| Потеря данных при миграции | Низкая | Высокое | Backup перед миграцией, поэтапный перенос |
| Проблемы с производительностью | Средняя | Среднее | Load testing, постепенный rollout |
| Баги в новой системе | Средняя | Высокое | Extensive testing, feature flags |
| Задержка в разработке | Низкая | Низкое | Буфер времени, параллельная работа |

## ✅ Критерии завершения

- [ ] Все session-based код удален
- [ ] Auth-service развернут и работает
- [ ] Frontend использует только новые endpoints
- [ ] Все тесты проходят
- [ ] Документация обновлена
- [ ] Мониторинг настроен
- [ ] Производительность соответствует требованиям
- [ ] Security audit пройден

---

**Следующие шаги:**
1. Утвердить план
2. Создать репозиторий auth-service
3. Начать с недели 1 (очистка legacy)
4. Еженедельные синки по прогрессу