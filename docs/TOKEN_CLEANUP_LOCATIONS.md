# 🔑 МЕСТА ОЧИСТКИ ТОКЕНОВ АВТОРИЗАЦИИ
## Дата: 2025-09-16

---

## 📍 ОСНОВНЫЕ ТОЧКИ ОЧИСТКИ

### 1. **AuthContext.tsx** - главный контекст авторизации
```typescript
// Строки 544-548: Основная функция logout
tokenManager.clearTokens();
localStorage.removeItem('svetu_user');
localStorage.removeItem('svetu_access_token');
localStorage.removeItem('svetu_refresh_token');

// Строка 153: При переполнении квоты
localStorage.clear();

// Строка 357: При обнаружении флага logout
localStorage.removeItem('svetu_logout_flag');
```

### 2. **tokenManager.ts** - менеджер токенов
```typescript
// Строки 159-168: Метод clearTokens()
clearTokens() {
  this.accessToken = null;
  localStorage.removeItem('svetu_access_token');
  localStorage.removeItem('svetu_refresh_token');
  this.clearRefreshTimer();
  this.refreshAttempts = 0;
  this.lastRefreshAttempt = 0;
  this.rateLimitedUntil = 0;
}
```

### 3. **AuthStateManager.tsx** - управление состоянием
```typescript
// Строка 42: Полная очистка localStorage (сохраняя корзину и локаль)
const cart = localStorage.getItem('svetu_cart');
localStorage.clear();
if (locale) localStorage.setItem('NEXT_LOCALE', locale);
if (cart) localStorage.setItem('svetu_cart', cart);

// Строка 38: Очистка sessionStorage
keysToRemove.forEach(key => sessionStorage.removeItem(key));
```

---

## 🧹 ВСПОМОГАТЕЛЬНЫЕ УТИЛИТЫ ОЧИСТКИ

### 4. **forceTokenCleanup.ts** - принудительная очистка
```typescript
// Удаляет старые токены по паттернам:
- access_token
- refresh_token
- auth_token
- jwt_token
- Все HS256 токены
- Невалидные JWT
```

### 5. **tokenMigration.ts** - миграция токенов
```typescript
// Строки 121, 134: Очистка старых ключей
localStorageKeysToRemove.forEach(key => localStorage.removeItem(key));
sessionStorageKeysToRemove.forEach(key => sessionStorage.removeItem(key));
```

### 6. **clearLargeHeaders.ts** - очистка больших данных
```typescript
// Удаляет элементы больше 8KB
if (value.length > MAX_HEADER_SIZE) {
  localStorage.removeItem(key);
  sessionStorage.removeItem(key);
}
```

---

## 🗺️ КАРТА ТОКЕНОВ

### localStorage токены:
- `svetu_access_token` - основной access токен
- `svetu_refresh_token` - refresh токен
- `svetu_user` - данные пользователя
- `svetu_logout_flag` - флаг выхода

### sessionStorage токены:
- `svetu_access_token` - дубликат в сессии
- `svetu_user` - данные пользователя в сессии
- `client_id` - ID клиента

### Устаревшие (удаляются при миграции):
- `access_token`
- `refresh_token`
- `auth_token`
- `jwt_token`
- `user`

---

## 🔄 СЦЕНАРИИ ОЧИСТКИ

### 1. **Обычный Logout**
- Вызов: `AuthContext.logout()`
- Действия:
  1. `tokenManager.clearTokens()`
  2. Удаление `svetu_user`, `svetu_access_token`, `svetu_refresh_token`
  3. Установка `svetu_logout_flag`
  4. Вызов API `/auth/logout`

### 2. **Автоматическая очистка при ошибке**
- При 401 ошибке
- При невалидном токене
- При истечении refresh токена

### 3. **Принудительная очистка**
- Скрипт `force-relogin.js`
- Утилита `forceTokenCleanup.ts`
- При переполнении localStorage

### 4. **Очистка при инициализации**
- `clearLargeHeaders.ts` - при загрузке приложения
- `tokenMigration.ts` - миграция старых токенов

---

## ⚠️ ПОТЕНЦИАЛЬНЫЕ ПРОБЛЕМЫ

### 1. **Неполная очистка**
- **Проблема**: Некоторые токены могут остаться в sessionStorage
- **Решение**: Использовать `AuthStateManager` для полной очистки

### 2. **Race conditions**
- **Проблема**: Одновременная очистка из разных мест
- **Решение**: Централизовать через `tokenManager.clearTokens()`

### 3. **Кэш браузера**
- **Проблема**: HTTP Only cookies не очищаются через JS
- **Решение**: Вызов backend `/auth/logout` для очистки cookies

### 4. **Множественные вкладки**
- **Проблема**: Очистка в одной вкладке не влияет на другие
- **Решение**: Использовать storage events для синхронизации

---

## 🛠️ РЕКОМЕНДАЦИИ

### Для разработчиков:
1. **Всегда используйте** `tokenManager.clearTokens()` для очистки
2. **Не очищайте** токены напрямую через `localStorage.removeItem()`
3. **Проверяйте** наличие токенов перед использованием
4. **Логируйте** все операции очистки для отладки

### Для тестирования:
```javascript
// Полная очистка всех токенов (для отладки)
function clearAllAuthData() {
  // localStorage
  ['svetu_access_token', 'svetu_refresh_token', 'svetu_user',
   'svetu_logout_flag', 'access_token', 'refresh_token',
   'auth_token', 'jwt_token', 'user'].forEach(key =>
    localStorage.removeItem(key)
  );

  // sessionStorage
  Object.keys(sessionStorage).forEach(key => {
    if (key.includes('token') || key.includes('auth') ||
        key.includes('user') || key.includes('svetu')) {
      sessionStorage.removeItem(key);
    }
  });

  // Очистка cookies через API
  fetch('/api/v1/auth/logout', {
    method: 'POST',
    credentials: 'include'
  });
}
```

---

## 📊 СТАТИСТИКА ИСПОЛЬЗОВАНИЯ

| Файл | Количество вызовов очистки | Типы очистки |
|------|----------------------------|---------------|
| AuthContext.tsx | 7 | logout, error handling, init |
| tokenManager.ts | 5 | clearTokens, setters |
| AuthStateManager.tsx | 2 | full clear |
| forceTokenCleanup.ts | 6 | migration, cleanup |
| tokenMigration.ts | 2 | migration |
| Остальные | 15+ | различные |

---

## 🔐 БЕЗОПАСНОСТЬ

### Важные моменты:
1. **HTTP Only Cookies** не могут быть очищены через JavaScript
2. **Refresh токены** должны удаляться через backend API
3. **При logout** всегда вызывайте backend endpoint
4. **Синхронизация вкладок** через storage events критична

### Правильный порядок очистки:
1. Вызов backend `/auth/logout`
2. Очистка localStorage токенов
3. Очистка sessionStorage
4. Установка флага logout
5. Редирект на страницу входа

---

**Дата аудита**: 2025-09-16
**Статус**: Система имеет множественные точки очистки, требуется унификация