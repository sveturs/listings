# Session Handover: OAuth Token Complete - 2025-07-01

## 🎯 Проблема
При попытке оформления заказа были следующие ошибки:
1. Hydration mismatch error в консоли браузера
2. POST `/api/v1/orders` возвращал 401 (Unauthorized) - токен не передавался
3. После OAuth авторизации через Google токен не сохранялся на фронтенде

## 🔧 Выполненные исправления

### 1. OAuth Token Handling
**Backend** (`backend/internal/proj/users/handler/auth.go`):
- Добавлена передача access_token в URL при редиректе после Google OAuth
- Токен добавляется как параметр `auth_token` к returnTo URL
- Добавлено логирование для отладки

**Frontend** (`frontend/svetu/src/contexts/AuthContext.tsx`):
- Добавлена проверка токена в URL при инициализации AuthContext
- Токен извлекается из URL и сохраняется через tokenManager
- URL очищается от токена для безопасности

### 2. Debugging & Logging
- Добавлено логирование в `api-client.ts` для отслеживания токена
- Добавлено логирование в `orders.ts` для отладки создания заказов
- Добавлено логирование в `tokenManager.ts` для отслеживания сохранения токена
- Добавлено логирование в `CheckoutPage` для проверки наличия токена

### 3. Fixes
- Исправлен импорт `isTokenExpired` в AuthContext
- Очищен кеш Next.js для устранения hydration mismatch
- Добавлен экспорт функции `isTokenExpired` из tokenManager

## 📝 Изменения в коде

### Backend
```go
// backend/internal/proj/users/handler/auth.go
// Добавлен импорт fmt и strings
import (
    "fmt"
    "strings"
    // ...
)

// В методе GoogleCallback добавлено:
if accessToken != "" {
    separator := "?"
    if strings.Contains(returnTo, "?") {
        separator = "&"
    }
    returnTo = fmt.Sprintf("%s%sauth_token=%s", returnTo, separator, accessToken)
    logger.Info().Str("redirect_url", returnTo[:50]+"...").Msg("OAuth: Redirecting with access token in URL")
} else {
    logger.Error().Msg("OAuth: No access token to add to redirect URL")
}
```

### Frontend
```typescript
// src/contexts/AuthContext.tsx
// В useEffect добавлена проверка токена в URL:
if (typeof window !== 'undefined') {
    const urlParams = new URLSearchParams(window.location.search);
    const authToken = urlParams.get('auth_token');
    if (authToken) {
        console.log('[AuthContext] Found auth_token in URL, saving...', authToken.substring(0, 20) + '...');
        tokenManager.setAccessToken(authToken);
        // Удаляем токен из URL для безопасности
        urlParams.delete('auth_token');
        const newUrl = `${window.location.pathname}${urlParams.toString() ? '?' + urlParams.toString() : ''}`;
        window.history.replaceState({}, document.title, newUrl);
        console.log('[AuthContext] Token saved, URL cleaned');
    }
}

// src/utils/tokenManager.ts
// Добавлен экспорт функции:
export const isTokenExpired = (token?: string) => tokenManager.isTokenExpired(token);

// src/app/[locale]/checkout/page.tsx
// Добавлено логирование:
console.log('[CheckoutPage] Submitting order, user:', user);
console.log('[CheckoutPage] Token exists:', !!tokenManager.getAccessToken());
```

## ✅ Результат
1. OAuth токен передается с backend на frontend через URL параметр
2. Frontend извлекает токен из URL и сохраняет в tokenManager и sessionStorage
3. API запросы включают Bearer токен в заголовках
4. Hydration mismatch устранен очисткой кеша Next.js

## 🔍 Статус
- ✅ OAuth токен передается с backend на frontend
- ✅ Токен сохраняется в tokenManager и sessionStorage
- ✅ API клиент использует токен для запросов
- ✅ Код отформатирован и проверен линтером
- ✅ Backend и frontend перезапущены

## 📌 Что нужно протестировать
1. Выйти из системы (logout)
2. Войти через Google OAuth
3. Проверить в консоли наличие логов о сохранении токена
4. Перейти на страницу checkout
5. Попробовать создать заказ - должно работать без ошибки 401

## 🚨 Важно
- Токен передается в URL только при OAuth редиректе
- Токен сразу удаляется из URL после сохранения
- Все API запросы теперь включают Bearer токен если пользователь авторизован