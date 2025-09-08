# 📝 Google Console OAuth - Настройка Redirect URIs

## 🔗 Authorized redirect URIs для Google Console

Добавьте следующие URIs в настройки OAuth 2.0 Client ID в Google Cloud Console:

### Для разработки (localhost):
```
http://localhost:3001/auth/oauth/google/callback
http://localhost:3001/en/auth/oauth/google/callback
http://localhost:3001/ru/auth/oauth/google/callback
http://localhost:3001/sr/auth/oauth/google/callback
```

### Для production (svetu.rs):
```
https://svetu.rs/auth/oauth/google/callback
https://svetu.rs/en/auth/oauth/google/callback
https://svetu.rs/ru/auth/oauth/google/callback
https://svetu.rs/sr/auth/oauth/google/callback
https://www.svetu.rs/auth/oauth/google/callback
https://www.svetu.rs/en/auth/oauth/google/callback
https://www.svetu.rs/ru/auth/oauth/google/callback
https://www.svetu.rs/sr/auth/oauth/google/callback
```

## 🔄 Как работает OAuth flow:

1. **Инициация OAuth:**
   - Пользователь нажимает "Login with Google"
   - Frontend вызывает `AuthService.loginWithGoogle()`
   - Происходит редирект на backend: `/api/v1/auth/google?redirect_uri=...`
   - Backend проксирует на микросервис
   - Микросервис редиректит на Google OAuth

2. **Google авторизация:**
   - Пользователь авторизуется в Google
   - Google редиректит на `http://localhost:3001/{locale}/auth/oauth/google/callback?code=...&state=...`

3. **Обработка callback:**
   - Frontend страница `/[locale]/auth/oauth/google/callback/page.tsx` получает код
   - Страница редиректит на backend: `/api/v1/auth/google/callback?code=...&state=...`
   - Backend проксирует на микросервис
   - Микросервис обменивает код на токены
   - Микросервис редиректит на frontend с токеном

## ⚙️ Как добавить в Google Console:

1. Откройте [Google Cloud Console](https://console.cloud.google.com/)
2. Выберите ваш проект
3. Перейдите в **APIs & Services** → **Credentials**
4. Найдите ваш **OAuth 2.0 Client ID**
5. Нажмите на него для редактирования
6. В секции **Authorized redirect URIs** добавьте все URIs из списка выше
7. Нажмите **Save**

## ⚠️ Важные замечания:

1. **URIs должны точно совпадать** - включая протокол (http/https), домен, порт и путь
2. **Локали обязательны** - Next.js добавляет локаль в URL автоматически
3. **Без trailing slash** - не добавляйте `/` в конец URI
4. **Production требует HTTPS** - Google не принимает HTTP для production доменов

## 🔍 Отладка:

Если получаете ошибку "redirect_uri_mismatch":
1. Проверьте точное совпадение URI в консоли и в приложении
2. Убедитесь, что добавили URI для всех локалей
3. Подождите 5-10 минут после сохранения изменений в Google Console

## 📋 Текущая конфигурация:

**Client ID:** `917315728307-au9ga5fl7o3bbid9nv7e4l92gut194pq.apps.googleusercontent.com`
**Client Secret:** `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`

---
*Документ создан: 07.09.2025*