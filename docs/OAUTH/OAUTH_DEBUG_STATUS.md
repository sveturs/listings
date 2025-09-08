# OAuth Отладка - Текущий Статус

## Что было исправлено

### 1. Конфигурация URL ✅
- Frontend `.env.local` изменен с IP адреса на `localhost`
- Backend `.env` изменен с IP адреса на `localhost:3001`
- Все сервисы теперь работают на localhost

### 2. Обработка OAuth токена ✅
- Удалена дублирующая обработка токена из `HomePageClient.tsx`
- Вся логика обработки токена теперь в `AuthContext.tsx`
- Добавлено расширенное логирование для отладки

### 3. Ключевые изменения в коде

#### AuthContext.tsx
- Добавлено автоматическое обнаружение `auth_token` в URL при загрузке
- Токен сохраняется в `tokenManager` и сразу вызывается `refreshSession()`
- Добавлено детальное логирование всех этапов

#### AuthService.ts
- Добавлено логирование при добавлении токена в заголовки запросов
- Улучшена обработка ошибок при восстановлении сессии

## Текущий OAuth Flow

1. **Пользователь нажимает "Войти через Google"**
   - Браузер перенаправляется на: `http://localhost:3000/api/v1/auth/oauth/google`

2. **Backend перенаправляет на Google**
   - Google OAuth consent screen

3. **Google возвращает на backend callback**
   - URL: `http://localhost:3000/auth/google/callback?code=...&state=...`

4. **Backend обрабатывает callback**
   - Обменивает код на токены
   - Создает JWT токен
   - Редиректит на frontend: `http://localhost:3001?auth_token=<JWT>`

5. **Frontend обрабатывает токен**
   - `AuthContext` обнаруживает токен в URL
   - Сохраняет токен через `tokenManager`
   - Вызывает `refreshSession()` для получения данных пользователя
   - Удаляет токен из URL для безопасности

## Как проверить работу OAuth

### 1. Подготовка
```bash
# Очистите кэш браузера и куки
# Откройте инкогнито режим в браузере
# Откройте консоль разработчика (F12)
```

### 2. Тестирование
1. Перейдите на http://localhost:3001
2. Нажмите кнопку "Войти"
3. Выберите "Войти через Google"
4. Авторизуйтесь в Google
5. После редиректа проверьте консоль браузера

### 3. Что должно быть в консоли при успехе
```
[AuthContext] Found OAuth token in URL: eyJhbGciOiJIUzI1NiIsInR...
[AuthContext] Token length: 193
[AuthContext] Token saved to tokenManager
[AuthContext] Verification - token retrieved: Success
[AuthContext] Token removed from URL for security
[AuthContext] Starting session refresh with new token...
[AuthService] Adding token to headers: eyJhbGciOiJIUzI1NiIsInR...
[AuthContext] Attempting to restore session via JWT...
[AuthContext] JWT session restored successfully
```

## Возможные проблемы и решения

### Проблема: 401 Unauthorized при вызове /api/v1/auth/session
**Причина**: Токен не передается в заголовках или токен невалидный
**Решение**: 
- Проверьте логи в консоли - есть ли сообщение "Adding token to headers"
- Проверьте Network tab - есть ли заголовок Authorization в запросе
- Проверьте что токен начинается с "Bearer "

### Проблема: Токен не обнаруживается в URL
**Причина**: Backend редиректит не на тот URL или без токена
**Решение**:
- Проверьте backend логи: `tail -f /tmp/backend.log`
- Убедитесь что FRONTEND_URL в backend .env = `http://localhost:3001`
- Проверьте что backend редиректит с параметром `auth_token`

### Проблема: Session refresh fails
**Причина**: Токен истек или невалидный
**Решение**:
- Проверьте JWT_SECRET одинаковый в backend и auth service
- Проверьте срок действия токена (должен быть 24 часа)
- Попробуйте сгенерировать новый токен

## Логи для мониторинга

### Frontend логи
```bash
tail -f /tmp/frontend.log
```

### Backend логи
```bash
tail -f /tmp/backend.log | grep -E "OAuth|JWT|auth"
```

### Проверка в браузере
1. Откройте Network tab в DevTools
2. Фильтр по "auth" или "session"
3. Проверьте заголовки запросов на наличие Authorization

## Текущий статус: ТРЕБУЕТ ТЕСТИРОВАНИЯ

OAuth flow полностью настроен со следующими улучшениями:
- ✅ URL конфигурация исправлена (localhost вместо IP)
- ✅ Удалено дублирование обработки токена
- ✅ Добавлено детальное логирование
- ⏳ Требует тестирования пользователем

Пожалуйста, протестируйте OAuth авторизацию и сообщите о результатах!