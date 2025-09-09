# Исправление OAuth Callback и Refresh Token

## Дата: 09.09.2025

## Проблема

После успешной авторизации через Google OAuth:
1. Frontend пытался обновить токены но получал 400/401 ошибки
2. `/api/v1/auth/refresh` и `/api/v1/auth/session` возвращали 401 
3. Refresh токены не сохранялись и терялись сразу после логина

## Корневая причина

OAuth callback страница на frontend (`/auth/oauth/google/callback`) делала запросы на неправильный URL:
- **Было**: `window.location.origin` (http://localhost:3001)
- **Нужно**: Backend API URL (http://localhost:3000)

Из-за этого:
1. Callback запрос шёл на frontend вместо backend
2. Cookies с refresh токеном устанавливались для localhost:3001
3. Backend на localhost:3000 не мог прочитать эти cookies (same-origin policy)

## Решение

### Изменения в `/frontend/svetu/src/app/[locale]/auth/oauth/google/callback/page.tsx`:

1. **Импортировали configManager**:
```typescript
import { configManager } from '@/config';
```

2. **Исправили URL для callback запроса**:
```typescript
// Было:
const backendCallbackUrl = new URL(
  '/api/v1/auth/google/callback',
  window.location.origin // Неправильно!
);

// Стало:
const backendCallbackUrl = new URL(
  '/api/v1/auth/google/callback',
  configManager.getApiUrl() // Правильно - backend URL
);
```

3. **Исправили URL для refresh запроса**:
```typescript
// Было:
await fetch('/api/v1/auth/refresh', ...)

// Стало:
await fetch(`${configManager.getApiUrl()}/api/v1/auth/refresh`, ...)
```

## Результат

✅ OAuth callback теперь корректно обращается к backend API
✅ Refresh токены должны правильно устанавливаться и использоваться
✅ Аутентификация должна сохраняться между сессиями

## Тестирование

1. Очистите все cookies и storage:
```javascript
localStorage.clear();
sessionStorage.clear();
document.cookie.split(';').forEach(c => {
  document.cookie = c.trim().split('=')[0] + '=;expires=' + new Date(0).toUTCString() + ';path=/'
});
```

2. Откройте новую вкладку в инкогнито

3. Перейдите на `http://localhost:3000/api/v1/auth/google`

4. Авторизуйтесь через Google

5. После редиректа на frontend проверьте:
   - DevTools → Network → посмотрите что callback идёт на localhost:3000
   - DevTools → Application → Cookies → должны быть установлены токены
   - Обновите страницу - аутентификация должна сохраниться

## Дополнительные замечания

### Cross-Origin Cookies
Основная проблема была в том, что Auth Service устанавливал cookies для frontend домена (localhost:3001), но backend API на другом порту (localhost:3000) не мог их прочитать из-за same-origin policy.

### Правильный флоу
1. Google OAuth → Auth Service
2. Auth Service → Redirect на Frontend callback (localhost:3001)
3. Frontend callback → Запрос к Backend API (localhost:3000) с auth code
4. Backend API → Проксирует к Auth Service
5. Auth Service → Возвращает токены
6. Backend → Устанавливает cookies для своего домена
7. Frontend → Сохраняет access token и использует API

## Статус

✅ Проблема решена
⏳ Требуется тестирование полного флоу с реальным Google OAuth