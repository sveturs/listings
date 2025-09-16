# 🔐 ЗАДАНИЕ ДЛЯ AUTH SERVICE - Исправление Cookie-based Refresh
## Дата: 2025-09-16
## Приоритет: КРИТИЧЕСКИЙ

---

## 📋 КОНТЕКСТ ПРОБЛЕМЫ

После OAuth авторизации на dev.svetu.rs пользователи теряют сессию при обновлении страницы. Проблема в несоответствии методов хранения refresh токена между Frontend и Auth Service.

### Текущая ситуация:
1. **Auth Service** отправляет refresh_token в HTTP-only cookie
2. **Frontend** ожидает refresh_token в localStorage
3. При обновлении страницы frontend не может использовать cookie для refresh

---

## ✅ ВЫПОЛНЕННЫЕ ИЗМЕНЕНИЯ НА FRONTEND

### 1. **TokenManager** (`/frontend/svetu/src/utils/tokenManager.ts`)
- Добавлен метод `performCookieRefresh()` для refresh через cookie
- Модифицирован `performRefresh()` - сначала пробует cookie, потом localStorage
- Поддержка двух методов авторизации (OAuth через cookie, email через localStorage)

### 2. **OAuthProcessor** (`/frontend/svetu/src/app/[locale]/auth/oauth/google/callback/OAuthProcessor.tsx`)
- Добавлен вызов `/api/v1/auth/session` после получения токена
- Это устанавливает session cookies для последующих запросов

### 3. **AuthService** (`/frontend/svetu/src/services/auth.ts`)
- Обновлен `restoreSession()` для работы с cookie-based refresh
- Автоматический fallback между методами

---

## 🎯 ТРЕБУЕМЫЕ ИЗМЕНЕНИЯ В AUTH SERVICE

### 1. Endpoint `/api/v1/auth/refresh` должен поддерживать ДВА режима:

#### A. Cookie-based refresh (для OAuth):
```go
// Если в запросе нет тела или Authorization header, но есть refresh_token cookie
if refreshCookie != nil && request.Body == nil {
    // Используем refresh_token из cookie
    newAccessToken, newRefreshToken := RefreshTokens(refreshCookie.Value)

    // Возвращаем новый access token в теле ответа
    return JSONResponse{
        "access_token": newAccessToken,
        "refresh_token": newRefreshToken, // опционально
    }
}
```

#### B. Token-based refresh (для email auth):
```go
// Существующая логика для refresh через тело запроса или header
if request.RefreshToken != "" || request.Header["Authorization"] != "" {
    // Используем refresh_token из запроса
    // ... существующая логика
}
```

### 2. CORS настройки должны включать:
```go
cors.Config{
    AllowOrigins:     []string{"https://dev.svetu.rs", "http://localhost:3001"},
    AllowCredentials: true,  // КРИТИЧНО для cookies
    AllowHeaders:     []string{"Authorization", "Content-Type"},
    ExposeHeaders:    []string{"Set-Cookie"},
}
```

### 3. Cookie настройки - КРИТИЧЕСКИ ВАЖНО:

⚠️ **ОБЯЗАТЕЛЬНО изменить SameSite во ВСЕХ местах установки cookies!**

Файл: `internal/transport/http/handlers/auth.go`

Нужно найти и заменить ВСЕ вхождения:
- `SameSite: "Strict"` → `SameSite: "None"`
- `SameSite: "Lax"` → `SameSite: "None"`

**Конкретные строки для замены:**
```go
// Строки ~113, ~162, ~279, ~306:
SameSite: "None", // Changed from Strict to None for cross-origin requests

// Строки ~495, ~866, ~1001:
SameSite: "None", // Changed from Lax to None for cross-origin requests
```

**Правильная конфигурация cookie:**
```go
http.Cookie{
    Name:     "refresh_token",
    Value:    refreshToken,
    Path:     "/",
    Domain:   ".svetu.rs",  // Важно: с точкой для поддоменов
    Secure:   true,          // HTTPS only
    HttpOnly: true,          // Защита от XSS
    SameSite: "None",        // КРИТИЧНО: должно быть "None" для cross-origin!
    MaxAge:   30 * 24 * 60 * 60,     // 30 дней
}
```

**Почему это критично:**
- `SameSite: "Strict"` или `"Lax"` блокирует отправку cookies при cross-origin запросах
- Frontend на dev.svetu.rs и Auth Service на другом домене/порту
- Без `SameSite: "None"` cookies не будут отправляться, и refresh не будет работать

---

## 📊 ТЕСТИРОВАНИЕ

### Сценарий 1: OAuth авторизация
1. Пользователь авторизуется через Google OAuth
2. Получает access_token в URL и refresh_token в cookie
3. При обновлении страницы:
   - Frontend отправляет POST `/api/v1/auth/refresh` с `credentials: include`
   - Auth Service читает refresh_token из cookie
   - Возвращает новый access_token в JSON
   - Frontend сохраняет access_token и продолжает работу

### Сценарий 2: Email авторизация
1. Пользователь логинится через email/password
2. Получает оба токена в теле ответа
3. Frontend сохраняет их в localStorage
4. При refresh отправляет refresh_token в теле запроса
5. Работает как раньше

---

## 🚨 КРИТИЧЕСКИЕ МОМЕНТЫ

1. **ОБЯЗАТЕЛЬНО изменить SameSite во ВСЕХ cookies**
   - Найти ВСЕ места где устанавливаются cookies
   - Заменить `SameSite: "Strict"` и `"Lax"` на `"None"`
   - Это КРИТИЧНО для работы cross-origin запросов

2. **НЕ ломать существующую email авторизацию**
   - Endpoint должен поддерживать ОБА метода
   - Приоритет: сначала проверяем cookie, потом тело/header

3. **Правильные CORS headers**
   - `Access-Control-Allow-Credentials: true`
   - `Access-Control-Allow-Origin` должен быть точным (не *)

4. **Domain для cookies**
   - Для dev: `.svetu.rs` (с точкой)
   - Для prod: `.svetu.rs` (с точкой)
   - Для localhost: пустой или `localhost`

---

## 📝 ПРИМЕРНЫЙ КОД ДЛЯ `/api/v1/auth/refresh`

```go
func (h *Handler) RefreshToken(c *gin.Context) {
    var request RefreshRequest

    // 1. Пробуем получить refresh token из cookie
    refreshCookie, err := c.Cookie("refresh_token")

    // 2. Если нет cookie, пробуем из тела запроса
    if err != nil || refreshCookie == "" {
        if err := c.ShouldBindJSON(&request); err != nil {
            // 3. Если нет тела, пробуем из header
            authHeader := c.GetHeader("Authorization")
            if authHeader != "" {
                request.RefreshToken = strings.TrimPrefix(authHeader, "Bearer ")
            }
        }
    } else {
        request.RefreshToken = refreshCookie
    }

    // 4. Проверяем что есть хоть какой-то refresh token
    if request.RefreshToken == "" {
        c.JSON(400, gin.H{"error": "refresh_token required"})
        return
    }

    // 5. Обновляем токены
    newTokens, err := h.service.RefreshTokens(request.RefreshToken)
    if err != nil {
        c.JSON(401, gin.H{"error": "invalid refresh token"})
        return
    }

    // 6. Устанавливаем новый refresh token в cookie (если это OAuth сессия)
    if refreshCookie != "" {
        c.SetCookie(
            "refresh_token",
            newTokens.RefreshToken,
            30*24*60*60,
            "/",
            ".svetu.rs",
            true,
            true,
        )
    }

    // 7. Возвращаем токены в теле ответа
    c.JSON(200, gin.H{
        "access_token": newTokens.AccessToken,
        "refresh_token": newTokens.RefreshToken,
    })
}
```

---

## ✅ КРИТЕРИИ УСПЕХА

1. OAuth пользователи НЕ теряют сессию при обновлении страницы на dev.svetu.rs
2. Email авторизация продолжает работать как раньше
3. Cookies правильно устанавливаются для поддоменов
4. CORS не блокирует запросы с credentials

---

## 📞 КОНТАКТЫ

При вопросах обращайтесь к разработчикам Frontend.
Frontend изменения уже готовы и ждут обновления Auth Service.

---

**Статус Frontend**: ✅ ГОТОВ
**Ожидаем**: Обновление Auth Service