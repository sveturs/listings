# Fix: "Cannot GET /auth/login" Error

**Дата:** 2025-10-19
**Статус:** ✅ РЕШЕНО
**Приоритет:** CRITICAL

---

## 🐛 Проблема

E2E тесты падали с ошибкой:
```
6:33PM ERR Error in handler error="Cannot GET /auth/login"
```

При попытке перехода на страницу логина, запрос уходил на backend API вместо Next.js frontend, что приводило к ошибке 404.

## 🔍 Корневая причина

В `frontend/svetu/next.config.ts` были неправильные rewrite правила, которые перенаправляли ВСЕ запросы к `/auth/*` на backend API:

```typescript
// ПРОБЛЕМНЫЕ ПРАВИЛА (УДАЛЕНЫ):
{
  source: '/:locale/auth/:path((?!callback|oauth/google/callback).*)',
  destination: `${apiUrl}/auth/:path*`,
},
{
  source: '/auth/:path((?!callback).*)',
  destination: `${apiUrl}/auth/:path*`,
}
```

### Почему это было проблемой:

1. **Frontend pages** (`/en/auth/login`, `/en/auth/register`) должны отображаться Next.js
2. **API requests** (`/api/v1/auth/login`) уже проксируются через правило `/api/:path*`
3. Эти rewrite rules создавали **двойное проксирование** - страницы логина пытались загрузиться с backend вместо Next.js

## ✅ Решение

### 1. Удалены проблемные rewrite rules

**Файл:** `frontend/svetu/next.config.ts`

**Изменение:**
```typescript
// ДО (неправильно):
async rewrites() {
  return [
    // ... другие правила
    {
      source: '/:locale/auth/:path((?!callback|oauth/google/callback).*)',
      destination: `${apiUrl}/auth/:path*`,
    },
    {
      source: '/auth/:path((?!callback).*)',
      destination: `${apiUrl}/auth/:path*`,
    },
    // ...
  ];
}

// ПОСЛЕ (правильно):
async rewrites() {
  return [
    // ... другие правила
    // УДАЛЕНЫ auth rewrite rules
    // Теперь /en/auth/login обрабатывается Next.js
    // API запросы /api/v1/auth/* проксируются через общее правило /api/:path*
  ];
}
```

### 2. Обновлены E2E тесты

**Файлы:**
- `e2e/user-journey-create-listing.spec.ts`
- `e2e/search.spec.ts`
- `e2e/axe/a11y-keyboard-navigation.spec.ts`
- `e2e/axe/a11y-wcag-compliance.spec.ts`

**Изменения:**
1. Кнопка "Login" → "Sign In" (актуальное название)
2. Селекторы `input[type="email"]` → `input[type="email"]:visible` (игнорируют Google One Tap скрытые поля)
3. Добавлена работа с модальным окном логина через `[role="dialog"]`

### 3. Создан упрощённый тест проверки auth routing

**Файл:** `e2e/simple-auth-test.spec.ts`

Этот тест проверяет главное:
- ✅ Страница `/en/auth/login` загружается с кодом 200
- ✅ Редирект `/auth/login` → `/en/auth/login` работает
- ✅ Нет ошибки "Cannot GET /auth/login"

## 📊 Результаты

### До исправления:
```
6:33PM INF REQUEST method=GET path=/auth/login
6:33PM ERR Error in handler error="Cannot GET /auth/login"
```

### После исправления:
```
GET /en/auth/login 200 in 181ms
```

## 🔧 Архитектура Auth Routing

### Правильная структура:

1. **Frontend Pages (Next.js):**
   - `/en/auth/login` - страница логина (Next.js page)
   - `/en/auth/register` - страница регистрации
   - `/en/auth/callback` - OAuth callback страница

2. **API Endpoints (Backend):**
   - `/api/v1/auth/login` - POST запрос для логина
   - `/api/v1/auth/register` - POST запрос для регистрации
   - `/api/v1/auth/refresh` - обновление токена

3. **BFF Proxy (Next.js API Routes):**
   - `/api/v2/auth/*` - BFF proxy для auth API (с httpOnly cookies)

### Как работает routing:

```
User navigates to /auth/login
  → Next.js middleware добавляет локаль
  → Redirect to /en/auth/login
  → Next.js отображает страницу login
  → User submits form
  → POST request to /api/v2/auth/login (BFF proxy)
  → BFF proxy forwards to /api/v1/auth/login (backend)
  → Backend validates credentials
  → Returns JWT token
  → BFF proxy sets httpOnly cookie
  → Redirect to /en/marketplace or /en/admin
```

## ⚠️ Известные проблемы

### Google One Tap interference

Google One Tap автоматически активируется на странице `/en/auth/login` и может блокировать E2E тесты.

**Решение для E2E тестов:**
1. Удалить Google One Tap iframe через `evaluate()`
2. Или добавить env variable `NEXT_PUBLIC_DISABLE_GOOGLE_ONE_TAP=true` для тестов
3. Использовать модальное окно логина вместо отдельной страницы

### Модальное окно vs Страница логина

В приложении есть ДВА способа входа:
1. **Модальное окно** - открывается при клике на "Sign In" в header (остаётся на текущей странице)
2. **Отдельная страница** - `/en/auth/login` (полноценная страница)

E2E тесты должны учитывать оба варианта.

## 📝 Чек-лист для проверки

- [x] Удалены проблемные rewrite rules из `next.config.ts`
- [x] Frontend перезапущен с новой конфигурацией
- [x] Страница `/en/auth/login` загружается с кодом 200
- [x] Редирект `/auth/login` → `/en/auth/login` работает
- [x] Нет ошибки "Cannot GET /auth/login" в логах
- [x] Обновлены E2E тесты с правильными селекторами
- [x] Создан упрощённый тест проверки routing

## 🔗 Связанные файлы

**Конфигурация:**
- `frontend/svetu/next.config.ts` - удалены rewrite rules
- `frontend/svetu/src/middleware.ts` - автоматическое добавление локали
- `frontend/svetu/playwright.config.ts` - конфигурация E2E тестов

**E2E тесты:**
- `e2e/simple-auth-test.spec.ts` - новый тест проверки routing
- `e2e/user-journey-create-listing.spec.ts` - обновлён
- `e2e/search.spec.ts` - обновлён
- `e2e/axe/a11y-*.spec.ts` - обновлены

**Документация:**
- `docs/E2E_HEADLESS_SETUP.md` - настройка headless mode
- `docs/FIX_AUTH_LOGIN_ROUTING.md` - этот файл

## 👥 Авторы

- Исправление: Claude Code
- Дата: 2025-10-19
- Версия: 0.2.4

---

**✅ ПРОБЛЕМА ПОЛНОСТЬЮ РЕШЕНА**
