# Спецификация библиотеки github.com/sveturs/auth v1.8.0

**Дата создания:** 2025-10-02
**Версия библиотеки:** v1.8.0
**Статус:** Полная избыточная самодостаточная спецификация

---

## 📋 Содержание

1. [Обзор](#обзор)
2. [Архитектура библиотеки](#архитектура-библиотеки)
3. [Основные компоненты](#основные-компоненты)
4. [HTTP Services](#http-services)
5. [Fiber Middleware](#fiber-middleware)
6. [Entity типы](#entity-типы)
7. [Роли и права доступа](#роли-и-права-доступа)
8. [OAuth интеграция](#oauth-интеграция)
9. [Примеры использования](#примеры-использования)
10. [Best Practices](#best-practices)
11. [Миграция и совместимость](#миграция-и-совместимость)

---

## Обзор

### Назначение

Библиотека `github.com/sveturs/auth` предоставляет клиентский SDK для интеграции Go-приложений с Auth Service - централизованным микросервисом аутентификации и авторизации.

### Ключевые особенности

- ✅ **Централизованная аутентификация** - все через auth-service микросервис
- ✅ **Fiber middleware** - готовые middleware для Fiber framework
- ✅ **OAuth 2.0** - полная поддержка Google OAuth (расширяемо)
- ✅ **Типизированные роли** - 30+ предопределенных ролей с permissions
- ✅ **User management** - CRUD операции с пользователями
- ✅ **Role management** - гибкая система управления ролями
- ✅ **OpenAPI client** - автоматически сгенерированный HTTP клиент

### Версия и зависимости

```go
// go.mod
require github.com/sveturs/auth v1.8.0

// Основные зависимости
require (
    github.com/gofiber/fiber/v2  // Web framework
    github.com/rs/zerolog        // Structured logging
    github.com/golang-jwt/jwt/v5 // JWT parsing
)
```

---

## Архитектура библиотеки

### Структура пакетов

```
github.com/sveturs/auth@v1.8.0/
├── pkg/                          # Публичные пакеты (используются клиентами)
│   ├── http/                     # HTTP клиенты и сервисы
│   │   ├── client/               # OpenAPI сгенерированный клиент
│   │   ├── entity/               # Доменные типы и структуры
│   │   ├── fiber/                # Fiber интеграция
│   │   │   └── middleware/       # Fiber middleware (JWTParser, RequireAuth)
│   │   └── service/              # Бизнес-логика клиента
│   │       ├── auth.go           # AuthService
│   │       ├── user.go           # UserService
│   │       └── oauth.go          # OAuthService
│   └── proto/                    # gRPC определения (опционально)
└── internal/                     # Внутренние пакеты auth-service
    └── ...                       # (не используются клиентами)
```

### Паттерны использования

1. **Инициализация клиента** → HTTP client к auth-service
2. **Создание сервисов** → AuthService, UserService, OAuthService
3. **Настройка middleware** → JWTParser + RequireAuth
4. **Регистрация роутов** → Применение middleware к эндпоинтам

---

## Основные компоненты

### 1. HTTP Client

Автоматически сгенерированный OpenAPI клиент для взаимодействия с auth-service.

```go
package client

type ClientWithResponsesInterface interface {
    // Auth endpoints
    PostApiV1AuthRegisterWithResponse(ctx context.Context, body RegisterRequest) (*PostApiV1AuthRegisterResponse, error)
    PostApiV1AuthLoginWithResponse(ctx context.Context, body LoginRequest) (*PostApiV1AuthLoginResponse, error)
    PostApiV1AuthRefreshWithResponse(ctx context.Context, body RefreshTokenRequest) (*PostApiV1AuthRefreshResponse, error)
    PostApiV1AuthLogoutWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*PostApiV1AuthLogoutResponse, error)
    GetApiV1AuthValidateWithResponse(ctx context.Context, params *GetApiV1AuthValidateParams) (*GetApiV1AuthValidateResponse, error)

    // User endpoints
    GetApiV1UsersAllWithResponse(ctx context.Context) (*GetApiV1UsersAllResponse, error)
    GetApiV1UsersIdWithResponse(ctx context.Context, id int) (*GetApiV1UsersIdResponse, error)
    GetApiV1UsersByEmailWithResponse(ctx context.Context, params *GetApiV1UsersByEmailParams) (*GetApiV1UsersByEmailResponse, error)
    PatchApiV1UsersIdWithResponse(ctx context.Context, id int, body UpdateProfileRequest) (*PatchApiV1UsersIdResponse, error)

    // Role endpoints
    GetApiV1UsersIdRolesWithResponse(ctx context.Context, id int) (*GetApiV1UsersIdRolesResponse, error)
    PostApiV1UsersIdRolesWithResponse(ctx context.Context, id int, body AddRoleRequest) (*PostApiV1UsersIdRolesResponse, error)
    DeleteApiV1UsersIdRolesRoleWithResponse(ctx context.Context, id int, role string) (*DeleteApiV1UsersIdRolesRoleResponse, error)

    // OAuth endpoints
    PostApiV1AuthOauthProviderUrlWithResponse(ctx context.Context, provider string, body GenerateOAuthURLRequest) (*PostApiV1AuthOauthProviderUrlResponse, error)
    PostApiV1AuthOauthProviderExchangeWithResponse(ctx context.Context, provider string, body ExchangeOAuthCodeRequest) (*PostApiV1AuthOauthProviderExchangeResponse, error)
}
```

**Создание клиента:**

```go
import authclient "github.com/sveturs/auth/pkg/http/client"

// Простой клиент
client, err := authclient.NewClientWithResponses("http://localhost:28080")

// С настройками
httpClient := &http.Client{Timeout: 30 * time.Second}
client, err := authclient.NewClientWithResponses(
    "http://localhost:28080",
    authclient.WithHTTPClient(httpClient),
)
```

---

## HTTP Services

### AuthService

Основной сервис для аутентификации и валидации токенов.

#### Создание

```go
import (
    authservice "github.com/sveturs/auth/pkg/http/service"
    "github.com/rs/zerolog"
)

// Создание AuthService с подключением к микросервису
authSvc := authservice.NewAuthService(client, logger)
```

#### Методы

```go
type AuthService struct {
    client    authclient.ClientWithResponsesInterface
    logger    zerolog.Logger
}

// ValidateToken - валидация JWT токена через auth-service
func (s *AuthService) ValidateToken(ctx context.Context, token string) (*entity.TokenValidationResponse, error)

// Register - регистрация нового пользователя
func (s *AuthService) Register(ctx context.Context, req entity.UserRegistrationRequest) (*authclient.PostApiV1AuthRegisterResponse, error)

// Login - аутентификация пользователя
func (s *AuthService) Login(ctx context.Context, req entity.UserLoginRequest) (*authclient.PostApiV1AuthLoginResponse, error)

// Logout - выход пользователя (с Authorization header)
func (s *AuthService) Logout(ctx context.Context, authHeader string) (*authclient.PostApiV1AuthLogoutResponse, error)

// RefreshToken - обновление токенов
func (s *AuthService) RefreshToken(ctx context.Context, req entity.RefreshTokenRequest) (*authclient.PostApiV1AuthRefreshResponse, error)

// GetClient - получение базового HTTP клиента
func (s *AuthService) GetClient() authclient.ClientWithResponsesInterface
```

#### Валидация через микросервис

Все валидации токенов происходят через централизованный auth-service:
- ✅ Проверяет актуальный статус токена
- ✅ Учитывает revocation (logout)
- ✅ Единый источник правды для всех сервисов
- ✅ Централизованное управление безопасностью

### UserService

Сервис для управления пользователями.

#### Создание

```go
userSvc := authservice.NewUserService(client, logger)
```

#### Методы

```go
type UserService struct {
    client authclient.ClientWithResponsesInterface
    logger zerolog.Logger
}

// GetAllUsers - получить всех пользователей
func (s *UserService) GetAllUsers(ctx context.Context) (*entity.UsersListResponse, error)

// GetUser - получить пользователя по ID
func (s *UserService) GetUser(ctx context.Context, userID int) (*entity.UserProfile, error)

// GetUserByEmail - получить пользователя по email
func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*entity.UserProfile, error)

// UpdateUserProfile - обновить профиль пользователя
func (s *UserService) UpdateUserProfile(ctx context.Context, userID int, req entity.UpdateProfileRequest) (*entity.UserProfile, error)

// UpdateUserStatus - изменить статус пользователя
func (s *UserService) UpdateUserStatus(ctx context.Context, userID int, req entity.UpdateStatusRequest) error

// DeleteUser - удалить пользователя (soft или permanent)
func (s *UserService) DeleteUser(ctx context.Context, userID int, permanent bool) (*entity.DeleteUserResponse, error)

// IsUserAdmin - проверить админ права
func (s *UserService) IsUserAdmin(ctx context.Context, userID int) (*entity.IsAdminResponse, error)

// GetUserRoles - получить роли пользователя
func (s *UserService) GetUserRoles(ctx context.Context, userID int) (*entity.UserRolesResponse, error)

// AddUserRole - добавить роль пользователю
func (s *UserService) AddUserRole(ctx context.Context, userID int, req entity.AddRoleRequest) (*entity.UserRolesResponse, error)

// RemoveUserRole - удалить роль у пользователя
func (s *UserService) RemoveUserRole(ctx context.Context, userID int, role string) (*entity.UserRolesResponse, error)

// GetAllRoles - получить все доступные роли
func (s *UserService) GetAllRoles(ctx context.Context) (*entity.AllRolesResponse, error)

// GetUsersByRole - получить пользователей с определенной ролью
func (s *UserService) GetUsersByRole(ctx context.Context, role string) (*entity.UsersListResponse, error)
```

### OAuthService

Сервис для OAuth 2.0 интеграции.

#### Создание

```go
oauthSvc := authservice.NewOAuthService(client)
```

#### Методы

```go
type OAuthService struct {
    mu     sync.Mutex
    states map[string]entity.OAuthState
    client authclient.ClientWithResponsesInterface
}

// GenerateState - генерация случайного state для CSRF защиты
func (s *OAuthService) GenerateState() string

// StoreState - сохранение state и метаданных
func (s *OAuthService) StoreState(stateID, provider, redirectURI, locale, returnPath string)

// ValidateState - проверка и удаление использованного state
func (s *OAuthService) ValidateState(stateID, provider string) (*entity.OAuthState, error)

// StartGoogleOAuth - начало Google OAuth flow
func (s *OAuthService) StartGoogleOAuth(ctx context.Context, redirectURI, locale, returnPath string) (string, error)

// CompleteGoogleOAuth - завершение Google OAuth и получение токенов
func (s *OAuthService) CompleteGoogleOAuth(ctx context.Context, code, state string) (*OAuthResult, error)
```

#### OAuthResult структура

```go
type OAuthResult struct {
    AccessToken  string
    RefreshToken string
    Email        string
    Locale       string
    ReturnPath   string
}
```

---

## Fiber Middleware

### JWTParser

Извлекает и валидирует JWT токен, сохраняя информацию о пользователе в контекст.

**ВАЖНО:** Этот middleware НЕ требует аутентификации - он только парсит токен если он есть.

#### Использование

```go
import authmiddleware "github.com/sveturs/auth/pkg/http/fiber/middleware"

// Создание middleware
jwtParser := authmiddleware.JWTParser(authService)

// Применение глобально
app.Use(jwtParser)

// Или к группе роутов
api := app.Group("/api", jwtParser)
```

#### Что делает

1. Извлекает токен из `Authorization: Bearer <token>` header
2. Валидирует токен через authService (обращение к микросервису)
3. Сохраняет в `fiber.Ctx.Locals()`:
   - `user_id` (int)
   - `email` (string)
   - `roles` ([]string)
   - `is_admin` (bool)
   - `authenticated` (bool)
   - `token` (string)
   - `name`, `term_accepted`, `email_verified`, `two_factor_enabled` (из claims)

4. **НЕ блокирует** запросы без токена или с невалидным токеном

#### Helper функции

```go
// GetUserID - извлечь user ID из контекста
userID, ok := authmiddleware.GetUserID(c)

// GetEmail - извлечь email
email, ok := authmiddleware.GetEmail(c)

// GetRoles - извлечь роли
roles, ok := authmiddleware.GetRoles(c)

// IsAuthenticated - проверить аутентификацию
if authmiddleware.IsAuthenticated(c) {
    // пользователь аутентифицирован
}

// IsAdmin - проверить админ роль
if authmiddleware.IsAdmin(c) {
    // пользователь админ
}

// GetToken - получить оригинальный JWT токен
token, ok := authmiddleware.GetToken(c)
```

### RequireAuth

Требует наличия валидного аутентифицированного пользователя и опционально проверяет роли.

**ВАЖНО:** Должен использоваться ПОСЛЕ JWTParser middleware.

#### Сигнатуры

```go
// RequireAuth - типизированные роли (рекомендуется)
func RequireAuth(roles ...entity.Role) fiber.Handler

// RequireAuthString - строковые роли (backward compatibility)
func RequireAuthString(roles ...string) fiber.Handler
```

#### Использование

```go
import (
    authmiddleware "github.com/sveturs/auth/pkg/http/fiber/middleware"
    "github.com/sveturs/auth/pkg/http/entity"
)

// Только аутентификация (любой пользователь)
app.Get("/protected",
    jwtParser,
    authmiddleware.RequireAuth(),
    handler,
)

// Требуется роль admin
app.Get("/admin/dashboard",
    jwtParser,
    authmiddleware.RequireAuth(entity.RoleAdmin),
    handler,
)

// Требуется одна из ролей (admin ИЛИ moderator)
app.Get("/moderate",
    jwtParser,
    authmiddleware.RequireAuth(entity.RoleAdmin, entity.RoleModerator),
    handler,
)

// Строковые роли (старый способ)
app.Get("/admin",
    jwtParser,
    authmiddleware.RequireAuthString("admin"),
    handler,
)
```

#### Ответы при ошибках

**401 Unauthorized** - нет аутентификации:
```json
{
  "error": "unauthorized",
  "message": "Authentication required"
}
```

**403 Forbidden** - недостаточно прав:
```json
{
  "error": "forbidden",
  "message": "Insufficient permissions"
}
```

---

## Entity типы

### Roles (Роли)

```go
package entity

type Role string

func (r Role) String() string

// Административные роли
const (
    RoleSuperAdmin Role = "super_admin"
    RoleAdmin      Role = "admin"
)

// Модерационные роли
const (
    RoleModerator        Role = "moderator"
    RoleContentModerator Role = "content_moderator"
    RoleReviewModerator  Role = "review_moderator"
    RoleChatModerator    Role = "chat_moderator"
    RoleDisputeManager   Role = "dispute_manager"
)

// Бизнес роли
const (
    RoleVendorManager    Role = "vendor_manager"
    RoleCategoryManager  Role = "category_manager"
    RoleMarketingManager Role = "marketing_manager"
    RoleFinancialManager Role = "financial_manager"
)

// Операционные роли
const (
    RoleWarehouseManager Role = "warehouse_manager"
    RoleWarehouseWorker  Role = "warehouse_worker"
    RolePickupManager    Role = "pickup_manager"
    RolePickupWorker     Role = "pickup_worker"
    RoleCourier          Role = "courier"
)

// Поддержка
const (
    RoleSupportL1 Role = "support_l1"
    RoleSupportL2 Role = "support_l2"
    RoleSupportL3 Role = "support_l3"
)

// Комплаенс
const (
    RoleLegalAdvisor      Role = "legal_advisor"
    RoleComplianceOfficer Role = "compliance_officer"
)

// Продавцы
const (
    RoleProfessionalVendor Role = "professional_vendor"
    RoleVendor             Role = "vendor"
    RoleIndividualSeller   Role = "individual_seller"
    RoleStorefrontStaff    Role = "storefront_staff"
)

// Покупатели
const (
    RoleVerifiedBuyer Role = "verified_buyer"
    RoleVIPCustomer   Role = "vip_customer"
    RoleUser          Role = "user"
)

// Аналитика
const (
    RoleDataAnalyst     Role = "data_analyst"
    RoleBusinessAnalyst Role = "business_analyst"
)
```

### Permissions (Разрешения)

Более 70 предопределенных permissions для точного контроля доступа.

```go
type Permission string

func (p Permission) String() string

// Примеры permissions
const (
    // User Management
    PermUsersView       Permission = "users.view"
    PermUsersEdit       Permission = "users.edit"
    PermUsersDelete     Permission = "users.delete"

    // Admin Panel
    PermAdminAccess       Permission = "admin.access"
    PermAdminCategories   Permission = "admin.categories"
    PermAdminTranslations Permission = "admin.translations"

    // Listings
    PermListingsCreate    Permission = "listings.create"
    PermListingsEditOwn   Permission = "listings.edit_own"
    PermListingsEditAny   Permission = "listings.edit_any"
    PermListingsModerate  Permission = "listings.moderate"

    // Orders
    PermOrdersViewAll Permission = "orders.view_all"
    PermOrdersProcess Permission = "orders.process"
    PermOrdersRefund  Permission = "orders.refund"

    // Payments
    PermPaymentsProcess Permission = "payments.process"
    PermPaymentsRefund  Permission = "payments.refund"

    // ... еще ~60 permissions
)
```

### Request/Response типы

#### Authentication

```go
// Регистрация
type UserRegistrationRequest struct {
    Email         string `json:"email" validate:"required,email"`
    Password      string `json:"password" validate:"required,min=8"`
    Name          string `json:"name" validate:"required,min=2,max=100"`
    TermsAccepted bool   `json:"terms_accepted" validate:"required"`
}

type RegisterResponse struct {
    AccessToken  string       `json:"access_token"`
    RefreshToken string       `json:"refresh_token"`
    TokenType    string       `json:"token_type"`
    ExpiresIn    int          `json:"expires_in"`
    User         *UserProfile `json:"user"`
}

// Логин
type UserLoginRequest struct {
    Email      string `json:"email" validate:"required,email"`
    Password   string `json:"password" validate:"required"`
    DeviceID   string `json:"device_id,omitempty"`
    DeviceName string `json:"device_name,omitempty"`
}

type LoginResponse struct {
    AccessToken  string       `json:"access_token"`
    RefreshToken string       `json:"refresh_token"`
    TokenType    string       `json:"token_type"`
    ExpiresIn    int          `json:"expires_in"`
    User         *UserProfile `json:"user"`
}

// Refresh
type RefreshTokenRequest struct {
    RefreshToken string `json:"refresh_token"`
}

type RefreshTokenResponse struct {
    AccessToken  string       `json:"access_token"`
    RefreshToken string       `json:"refresh_token"`
    TokenType    string       `json:"token_type"`
    ExpiresIn    int          `json:"expires_in"`
    User         *UserProfile `json:"user"`
}

// Валидация
type TokenValidationResponse struct {
    Valid  bool                   `json:"valid"`
    UserID int                    `json:"user_id,omitempty"`
    Email  string                 `json:"email,omitempty"`
    Roles  []string               `json:"roles,omitempty"`
    Claims map[string]interface{} `json:"claims,omitempty"`
    Error  string                 `json:"error,omitempty"`
}
```

#### User Management

```go
type UserProfile struct {
    ID               int       `json:"id"`
    Email            string    `json:"email"`
    Name             string    `json:"name"`
    PictureURL       string    `json:"picture_url,omitempty"`
    Phone            string    `json:"phone,omitempty"`
    PhoneVerified    bool      `json:"phone_verified"`
    Bio              string    `json:"bio,omitempty"`
    Timezone         string    `json:"timezone"`
    City             string    `json:"city,omitempty"`
    Country          string    `json:"country,omitempty"`
    EmailVerified    bool      `json:"email_verified"`
    TwoFactorEnabled bool      `json:"two_factor_enabled"`
    IsAdmin          bool      `json:"is_admin"`
    CreatedAt        time.Time `json:"created_at"`
    UpdatedAt        time.Time `json:"updated_at"`
}

type UpdateProfileRequest struct {
    Name     *string `json:"name,omitempty" validate:"omitempty,min=2,max=100"`
    Phone    *string `json:"phone,omitempty"`
    Bio      *string `json:"bio,omitempty" validate:"omitempty,max=500"`
    Timezone *string `json:"timezone,omitempty"`
    City     *string `json:"city,omitempty" validate:"omitempty,max=100"`
    Country  *string `json:"country,omitempty" validate:"omitempty,max=100"`
}

type UpdateStatusRequest struct {
    Status string `json:"status" validate:"required,oneof=active suspended banned deleted"`
}

type UsersListResponse struct {
    Users []*UserProfile `json:"users"`
}
```

#### Role Management

```go
type AddRoleRequest struct {
    Role string `json:"role" validate:"required"`
}

type UserRolesResponse struct {
    UserID int      `json:"user_id"`
    Roles  []string `json:"roles"`
}

type IsAdminResponse struct {
    UserID     int      `json:"user_id"`
    IsAdmin    bool     `json:"is_admin"`
    AdminRoles []string `json:"admin_roles,omitempty"`
}

type AllRolesResponse struct {
    Roles []RoleInfo `json:"roles"`
}

type RoleInfo struct {
    Name        string `json:"name"`
    Description string `json:"description"`
}
```

#### OAuth

```go
type OAuthState struct {
    Provider    string
    RedirectURI string
    Locale      string
    ReturnPath  string
    CreatedAt   time.Time
}

type GenerateOAuthURLRequest struct {
    RedirectURI string `json:"redirect_uri"`
    State       string `json:"state,omitempty"`
}

type GenerateOAuthURLResponse struct {
    URL string `json:"url"`
}

type ExchangeOAuthCodeRequest struct {
    Code        string `json:"code"`
    RedirectURI string `json:"redirect_uri"`
    State       string `json:"state,omitempty"`
}

type OAuthExchangeResponse struct {
    AccessToken  string   `json:"access_token"`
    RefreshToken string   `json:"refresh_token"`
    User         UserInfo `json:"user"`
}

type UserInfo struct {
    ID      int      `json:"id"`
    Email   string   `json:"email"`
    Name    string   `json:"name,omitempty"`
    Picture string   `json:"picture,omitempty"`
    Roles   []string `json:"roles"`
    IsAdmin bool     `json:"is_admin"`
}
```

---

## Роли и права доступа

### Иерархия приоритетов

```go
type Priority int

const (
    PrioritySuperAdmin       Priority = 1   // Наивысший
    PriorityAdmin            Priority = 10
    PriorityModerator        Priority = 20
    PriorityManager          Priority = 30
    PrioritySupport          Priority = 40
    PriorityVendor           Priority = 50
    PriorityVerifiedCustomer Priority = 60
    PriorityUser             Priority = 100 // Базовый
)
```

### Группировка ролей

#### Административный уровень
- `super_admin` - полный доступ ко всему
- `admin` - административный доступ

#### Модерация
- `moderator` - общая модерация
- `content_moderator` - модерация контента
- `review_moderator` - модерация отзывов
- `chat_moderator` - модерация чата
- `dispute_manager` - управление спорами

#### Бизнес управление
- `vendor_manager` - управление продавцами
- `category_manager` - управление категориями
- `marketing_manager` - маркетинг
- `financial_manager` - финансы

#### Операции и логистика
- `warehouse_manager`, `warehouse_worker` - склад
- `pickup_manager`, `pickup_worker` - пункты выдачи
- `courier` - доставка

#### Поддержка клиентов
- `support_l1`, `support_l2`, `support_l3` - уровни поддержки

#### Продавцы и клиенты
- `professional_vendor`, `vendor`, `individual_seller` - типы продавцов
- `storefront_staff` - персонал витрины
- `verified_buyer`, `vip_customer`, `user` - покупатели

### Проверка ролей в коде

```go
// В middleware
authmiddleware.RequireAuth(entity.RoleAdmin)

// В handler
func (h *Handler) AdminOnly(c *fiber.Ctx) error {
    if !authmiddleware.IsAdmin(c) {
        return fiber.ErrForbidden
    }

    roles, _ := authmiddleware.GetRoles(c)

    // Проверка конкретной роли
    hasModeratorRole := false
    for _, role := range roles {
        if role == entity.RoleModerator.String() {
            hasModeratorRole = true
            break
        }
    }

    // ...
}
```

---

## OAuth интеграция

### Google OAuth Flow

#### 1. Инициализация OAuth

```go
import authservice "github.com/sveturs/auth/pkg/http/service"

oauthSvc := authservice.NewOAuthService(authClient)

func (h *Handler) GoogleAuth(c *fiber.Ctx) error {
    // Параметры из query
    locale := c.Query("locale", "en")
    returnPath := c.Query("return_path", "/")

    // Построить redirect URI
    redirectURI := fmt.Sprintf("%s/api/v1/auth/google/callback", h.backendURL)

    // Получить OAuth URL
    authURL, err := h.oauthSvc.StartGoogleOAuth(
        c.Context(),
        redirectURI,
        locale,
        returnPath,
    )
    if err != nil {
        return err
    }

    // Редирект на Google
    return c.Redirect(authURL, fiber.StatusTemporaryRedirect)
}
```

#### 2. Обработка callback

```go
func (h *Handler) GoogleCallback(c *fiber.Ctx) error {
    code := c.Query("code")
    state := c.Query("state")

    // Обмен code на токены
    result, err := h.oauthSvc.CompleteGoogleOAuth(
        c.Context(),
        code,
        state,
    )
    if err != nil {
        // Редирект на frontend с ошибкой
        return c.Redirect(fmt.Sprintf(
            "%s/auth/error?message=%s",
            h.frontendURL,
            url.QueryEscape(err.Error()),
        ))
    }

    // Установка cookies
    c.Cookie(&fiber.Cookie{
        Name:     "access_token",
        Value:    result.AccessToken,
        HTTPOnly: true,
        Secure:   true,
        SameSite: "Lax",
        MaxAge:   15 * 60, // 15 минут
    })

    c.Cookie(&fiber.Cookie{
        Name:     "refresh_token",
        Value:    result.RefreshToken,
        HTTPOnly: true,
        Secure:   true,
        SameSite: "Lax",
        MaxAge:   30 * 24 * 60 * 60, // 30 дней
    })

    // Редирект на frontend
    redirectURL := fmt.Sprintf(
        "%s%s?locale=%s",
        h.frontendURL,
        result.ReturnPath,
        result.Locale,
    )

    return c.Redirect(redirectURL)
}
```

### State Management

OAuth service автоматически управляет state для CSRF защиты:

- ✅ Генерация криптографически безопасного state
- ✅ Хранение метаданных (provider, redirectURI, locale, returnPath)
- ✅ Автоматическая очистка старых state (>10 минут)
- ✅ One-time use (state удаляется после валидации)

---

## Примеры использования

### Полная инициализация в main.go

```go
package main

import (
    "github.com/gofiber/fiber/v2"
    "github.com/rs/zerolog"

    authclient "github.com/sveturs/auth/pkg/http/client"
    authmiddleware "github.com/sveturs/auth/pkg/http/fiber/middleware"
    authservice "github.com/sveturs/auth/pkg/http/service"
)

func main() {
    // Logger
    logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

    // Auth client
    authClient, err := authclient.NewClientWithResponses("http://localhost:28080")
    if err != nil {
        logger.Fatal().Err(err).Msg("Failed to create auth client")
    }

    // Services
    authSvc := authservice.NewAuthServiceWithLocalValidation(authClient, logger)
    userSvc := authservice.NewUserService(authClient, logger)
    oauthSvc := authservice.NewOAuthService(authClient)

    // Middleware
    jwtParser := authmiddleware.JWTParser(authSvc)

    // Fiber app
    app := fiber.New()

    // Global JWT parser
    app.Use(jwtParser)

    // Public routes
    app.Post("/api/v1/auth/login", loginHandler)
    app.Post("/api/v1/auth/register", registerHandler)

    // Protected routes (any authenticated user)
    protected := app.Group("/api/v1/protected", authmiddleware.RequireAuth())
    protected.Get("/profile", getProfileHandler)

    // Admin routes
    admin := app.Group("/api/v1/admin",
        authmiddleware.RequireAuth(entity.RoleAdmin),
    )
    admin.Get("/users", getUsersHandler)
    admin.Post("/users/:id/roles", addRoleHandler)

    // Vendor routes
    vendor := app.Group("/api/v1/vendor",
        authmiddleware.RequireAuth(entity.RoleVendor, entity.RoleProfessionalVendor),
    )
    vendor.Get("/dashboard", vendorDashboardHandler)

    app.Listen(":3000")
}
```

### Handler примеры

```go
// Простой protected handler
func getProfileHandler(c *fiber.Ctx) error {
    userID, ok := authmiddleware.GetUserID(c)
    if !ok {
        return fiber.ErrUnauthorized
    }

    // Получить профиль из БД
    profile, err := db.GetUserProfile(userID)
    if err != nil {
        return err
    }

    return c.JSON(profile)
}

// Admin handler с проверкой ролей
func getUsersHandler(c *fiber.Ctx) error {
    // Middleware уже проверил роль admin

    users, err := userSvc.GetAllUsers(c.Context())
    if err != nil {
        return err
    }

    return c.JSON(users)
}

// Handler с условной логикой по ролям
func createListingHandler(c *fiber.Ctx) error {
    userID, _ := authmiddleware.GetUserID(c)
    roles, _ := authmiddleware.GetRoles(c)

    // Разные лимиты для разных ролей
    var listingLimit int
    switch {
    case containsRole(roles, entity.RoleProfessionalVendor):
        listingLimit = 1000
    case containsRole(roles, entity.RoleVendor):
        listingLimit = 100
    default:
        listingLimit = 10
    }

    // Проверка лимита
    count, _ := db.CountUserListings(userID)
    if count >= listingLimit {
        return fiber.NewError(fiber.StatusForbidden, "Listing limit reached")
    }

    // Создание объявления
    // ...
}

func containsRole(roles []string, role entity.Role) bool {
    for _, r := range roles {
        if r == role.String() {
            return true
        }
    }
    return false
}
```

### Использование UserService

```go
func adminGetUsersHandler(c *fiber.Ctx) error {
    // Получить всех пользователей
    users, err := userSvc.GetAllUsers(c.Context())
    if err != nil {
        return err
    }

    return c.JSON(users)
}

func adminUpdateUserHandler(c *fiber.Ctx) error {
    userID, err := c.ParamsInt("id")
    if err != nil {
        return fiber.ErrBadRequest
    }

    var req entity.UpdateProfileRequest
    if err := c.BodyParser(&req); err != nil {
        return fiber.ErrBadRequest
    }

    profile, err := userSvc.UpdateUserProfile(c.Context(), userID, req)
    if err != nil {
        return err
    }

    return c.JSON(profile)
}

func adminAddRoleHandler(c *fiber.Ctx) error {
    userID, _ := c.ParamsInt("id")

    var req entity.AddRoleRequest
    if err := c.BodyParser(&req); err != nil {
        return fiber.ErrBadRequest
    }

    roles, err := userSvc.AddUserRole(c.Context(), userID, req)
    if err != nil {
        return err
    }

    return c.JSON(roles)
}
```

---

## Best Practices

### 1. Подключайтесь к auth-service микросервису

```go
// ✅ Правильно - валидация через централизованный сервис
authSvc := authservice.NewAuthService(client, logger)
```

### 2. Применяйте JWTParser глобально

```go
// ✅ Хорошо
app.Use(jwtParser)

// ❌ Плохо (дублирование)
app.Get("/route1", jwtParser, authmiddleware.RequireAuth(), handler1)
app.Get("/route2", jwtParser, authmiddleware.RequireAuth(), handler2)
```

### 3. Используйте типизированные роли

```go
// ✅ Хорошо
authmiddleware.RequireAuth(entity.RoleAdmin)

// ❌ Плохо (магические строки)
authmiddleware.RequireAuthString("admin")
```

### 4. Группируйте роуты по уровням доступа

```go
// ✅ Хорошо
admin := app.Group("/admin", authmiddleware.RequireAuth(entity.RoleAdmin))
admin.Get("/users", handler1)
admin.Get("/settings", handler2)

vendor := app.Group("/vendor", authmiddleware.RequireAuth(entity.RoleVendor))
vendor.Get("/products", handler3)

// ❌ Плохо
app.Get("/admin/users", authmiddleware.RequireAuth(entity.RoleAdmin), handler1)
app.Get("/admin/settings", authmiddleware.RequireAuth(entity.RoleAdmin), handler2)
```

### 5. Обрабатывайте ошибки валидации

```go
// ✅ Хорошо
validation, err := authSvc.ValidateToken(ctx, token)
if err != nil {
    logger.Error().Err(err).Msg("Token validation failed")
    return err
}
if !validation.Valid {
    return fiber.ErrUnauthorized
}

// ❌ Плохо
validation, _ := authSvc.ValidateToken(ctx, token)
// Игнорируем ошибки
```

### 6. Логируйте важные события

```go
// ✅ Хорошо
func loginHandler(c *fiber.Ctx) error {
    resp, err := authSvc.Login(ctx, req)
    if err != nil {
        logger.Warn().
            Str("email", req.Email).
            Err(err).
            Msg("Login failed")
        return err
    }

    logger.Info().
        Str("email", req.Email).
        Msg("User logged in successfully")

    return c.JSON(resp)
}
```

### 7. Используйте контекст с таймаутами

```go
// ✅ Хорошо
ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
defer cancel()

users, err := userSvc.GetAllUsers(ctx)

// ❌ Плохо
users, err := userSvc.GetAllUsers(context.Background())
```

### 8. Не храните sensitive данные в claims

```go
// ✅ Хорошо - только IDs и роли
userID, _ := authmiddleware.GetUserID(c)
profile, _ := db.GetUserProfile(userID) // Полная информация из БД

// ❌ Плохо - полагаться на claims для критичных данных
claims, _ := authmiddleware.GetClaims(c)
balance := claims["balance"] // Может быть устаревшим
```

### 9. Проверяйте доступ на уровне handler

```go
// ✅ Хорошо - двойная проверка
func updateListingHandler(c *fiber.Ctx) error {
    listingID := c.ParamsInt("id")
    userID, _ := authmiddleware.GetUserID(c)

    listing, _ := db.GetListing(listingID)

    // Проверяем владельца
    if listing.OwnerID != userID && !authmiddleware.IsAdmin(c) {
        return fiber.ErrForbidden
    }

    // Обновление
    // ...
}
```

### 10. Правильная обработка OAuth

```go
// ✅ Хорошо
func googleCallback(c *fiber.Ctx) error {
    result, err := oauthSvc.CompleteGoogleOAuth(c.Context(), code, state)
    if err != nil {
        // Редирект на frontend с ошибкой
        errorURL := fmt.Sprintf("%s/auth/error?message=%s",
            frontendURL,
            url.QueryEscape(err.Error()),
        )
        return c.Redirect(errorURL)
    }

    // Установка secure cookies
    setCookie(c, "access_token", result.AccessToken, true, true)
    setCookie(c, "refresh_token", result.RefreshToken, true, true)

    return c.Redirect(fmt.Sprintf("%s%s", frontendURL, result.ReturnPath))
}

// ❌ Плохо
func googleCallback(c *fiber.Ctx) error {
    result, _ := oauthSvc.CompleteGoogleOAuth(c.Context(), code, state)
    // Игнорируем ошибки

    // Токены в query параметрах (небезопасно!)
    return c.Redirect(fmt.Sprintf(
        "%s?access_token=%s&refresh_token=%s",
        frontendURL,
        result.AccessToken,
        result.RefreshToken,
    ))
}
```

---

## Миграция и совместимость

### Версионность

- **v1.8.0** (текущая) - стабильная версия с улучшениями
- **v1.7.x** - добавлены role management методы
- **v1.6.x** - добавлен UserService
- **v1.5.x** - первая стабильная версия с OAuth

### Breaking changes

#### v1.8.0
- Улучшена производительность валидации
- Обновлены зависимости
- Backward compatible

#### v1.7.0
- Расширены типы ролей до 30+
- Добавлены permissions
- Backward compatible

### Миграция с v1.6.x на v1.8.0

```go
// v1.6.x и v1.8.0 - совместимы
authSvc := authservice.NewAuthService(client, logger)

// Никаких изменений не требуется
```

### Обновление зависимости

```bash
# Обновить до последней версии
go get github.com/sveturs/auth@latest

# Или конкретная версия
go get github.com/sveturs/auth@v1.8.0

# Обновить go.mod и go.sum
go mod tidy
```

---

## Заключение

### Ключевые преимущества библиотеки

✅ **Простая интеграция** - минимальный boilerplate код
✅ **Централизация** - все через auth-service микросервис
✅ **Безопасность** - проверенные middleware и OAuth flow
✅ **Масштабируемость** - 30+ ролей и 70+ permissions
✅ **Типобезопасность** - строгая типизация во всех API
✅ **Расширяемость** - легко добавить новые OAuth провайдеры

### Дополнительные ресурсы

- [Auth Service README](https://github.com/sveturs/auth/blob/main/README.md)
- [OpenAPI спецификация](https://github.com/sveturs/auth/blob/main/swagger/openapi3.yaml)
- [Примеры использования](https://github.com/sveturs/auth/tree/main/examples)
- [Changelog](https://github.com/sveturs/auth/blob/main/CHANGELOG.md)

### Поддержка

При возникновении вопросов или проблем:
1. Проверьте документацию
2. Посмотрите примеры в `/examples`
3. Создайте issue в репозитории

---

**Версия документа:** 1.0
**Автор:** Claude Code Analysis
**Дата:** 2025-10-02
