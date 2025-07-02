# Паспорт API Endpoints: Auth (Аутентификация)

## 📋 Метаданные
- **Группа API**: Authentication
- **Базовый путь**: `/api/v1/auth`
- **Handler**: `backend/internal/proj/users/handler/routes.go`
- **Количество endpoints**: 8
- **Безопасность**: CSRF защита, rate limiting, session management

## 🎯 Назначение
Управление аутентификацией и авторизацией пользователей через:
- Локальная регистрация/вход (email/password)
- OAuth Google интеграция
- JWT токены и refresh tokens
- Session management

## 📡 Endpoints

### 🔓 Публичные (без авторизации)

#### POST `/api/v1/auth/register`
**Назначение**: Регистрация нового пользователя
- **Handler**: `h.Auth.Register`
- **Body**: `{"email": "string", "password": "string", "name": "string"}`
- **Response**: JWT токен + refresh token
- **Валидация**: Email уникальность, password requirements

#### POST `/api/v1/auth/login`
**Назначение**: Вход в систему
- **Handler**: `h.Auth.Login`
- **Body**: `{"email": "string", "password": "string"}`
- **Response**: JWT токен + refresh token
- **Особенности**: Rate limiting для защиты от brute force

#### POST `/api/v1/auth/refresh`
**Назначение**: Обновление JWT токена
- **Handler**: `h.Auth.RefreshToken`
- **Body**: `{"refresh_token": "string"}`
- **Response**: Новый JWT токен
- **Security**: Проверка валидности refresh token

#### GET `/api/v1/auth/google`
**Назначение**: Инициация OAuth авторизации Google
- **Handler**: `h.Auth.GoogleAuth`
- **Response**: Редирект на Google OAuth
- **Params**: `state`, `redirect_uri`

#### GET `/api/v1/auth/google/callback`
**Назначение**: Callback для Google OAuth
- **Handler**: `h.Auth.GoogleCallback`
- **Params**: `code`, `state`
- **Response**: Редирект с токенами или ошибкой

### 🔒 Защищенные (требуют авторизации)

#### GET `/api/v1/auth/session`
**Назначение**: Получение информации о текущей сессии
- **Handler**: `h.Auth.GetSession`
- **Headers**: `Authorization: Bearer <token>`
- **Response**: User profile + session details
- **Используется**: AuthContext на фронтенде

#### POST `/api/v1/auth/logout`
**Назначение**: Выход из системы
- **Handler**: `h.Auth.Logout`
- **Effect**: Инвалидация JWT и refresh токенов
- **Response**: Success message

### 🔍 Служебные

#### GET `/api/v1/admin-check/:email`
**Назначение**: Проверка административных прав
- **Handler**: `h.User.IsAdminPublic`
- **Params**: `email` - email пользователя
- **Response**: `{"is_admin": boolean}`
- **Используется**: Frontend для показа админ-интерфейса

## 🔐 Модель безопасности

### JWT Токены
```typescript
interface JWTPayload {
  user_id: string;
  email: string;
  name: string;
  role: "user" | "admin";
  exp: number; // 15 минут
}

interface RefreshToken {
  token_id: string;
  user_id: string;
  expires_at: string; // 30 дней
}
```

### Session Management
- **Access Token**: 15 минут lifetime
- **Refresh Token**: 30 дней lifetime, stored in database
- **CSRF Protection**: Required for state-changing operations
- **Rate Limiting**: Login attempts limited per IP

## 🔄 Интеграции

### Google OAuth 2.0
```typescript
interface GoogleOAuthConfig {
  client_id: string;
  client_secret: string;
  redirect_uri: string;
  scopes: ["openid", "email", "profile"];
}
```

### Database Tables
- `users` - основная информация пользователей
- `refresh_tokens` - активные refresh токены
- `admin_users` - список администраторов

## 🎭 Типы данных

### Запросы
```typescript
interface RegisterRequest {
  email: string;        // валидный email
  password: string;     // мин 8 символов
  name: string;         // 2-50 символов
  terms_accepted: boolean;
}

interface LoginRequest {
  email: string;
  password: string;
  remember_me?: boolean;
}
```

### Ответы
```typescript
interface AuthResponse {
  access_token: string;
  refresh_token: string;
  user: {
    id: string;
    email: string;
    name: string;
    avatar_url?: string;
    role: "user" | "admin";
    created_at: string;
  };
}

interface SessionResponse {
  user: UserProfile;
  session: {
    created_at: string;
    expires_at: string;
    ip_address: string;
  };
}
```

## 🌐 Frontend интеграция

### AuthContext использует:
- `POST /auth/login` - для входа
- `POST /auth/logout` - для выхода
- `GET /auth/session` - для проверки сессии
- `POST /auth/refresh` - для автообновления токенов

### Google OAuth Flow:
1. User clicks "Sign in with Google"
2. Redirect to `/auth/google`
3. Google redirects to `/auth/google/callback`
4. Backend creates session and redirects to frontend
5. Frontend calls `/auth/session` to get user data

## ⚠️ Известные особенности

### Security Features
- **CORS**: Настроен для домена svetu.rs
- **CSRF**: Токены обязательны для POST/PUT/DELETE
- **Rate Limiting**: 5 попыток входа в минуту с IP
- **Password Policy**: Минимум 8 символов, спецсимволы

### Error Handling
- Все ошибки возвращают локализованные placeholders
- Реальные ошибки логируются в backend
- 401/403 ошибки перенаправляют на login

### Session Security
- Refresh токены храняться в БД с возможностью отзыва
- JWT содержат только публичную информацию
- Автоматическое обновление токенов в AuthContext

## 🔧 Настройки окружения

```env
# JWT Settings
JWT_SECRET=your-secret-key
JWT_EXPIRY=15m
REFRESH_TOKEN_EXPIRY=30d

# Google OAuth
GOOGLE_CLIENT_ID=your-client-id
GOOGLE_CLIENT_SECRET=your-secret
GOOGLE_REDIRECT_URI=https://svetu.rs/auth/google/callback

# Security
CSRF_SECRET=your-csrf-secret
RATE_LIMIT_LOGIN=5/min
```

## 🧪 Примеры использования

### Регистрация
```bash
curl -X POST /api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"securepass123","name":"John Doe"}'
```

### Проверка сессии
```bash
curl -X GET /api/v1/auth/session \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

### Обновление токена
```bash
curl -X POST /api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{"refresh_token":"uuid-refresh-token-here"}'
```