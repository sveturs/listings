# Настройка Google OAuth Console

## Текущая проблема
OAuth flow настроен, но Google не принимает redirect URI `http://localhost:3001/auth/oauth/callback`, так как он не добавлен в консоль Google.

## Необходимые действия в Google Cloud Console

1. Перейти в [Google Cloud Console](https://console.cloud.google.com/)
2. Выбрать проект (или создать новый)
3. Перейти в **APIs & Services** → **Credentials**
4. Найти OAuth 2.0 Client ID: `917315728307-au9ga5fl7o3bbid9nv7e4l92gut194pq.apps.googleusercontent.com`
5. Нажать на него для редактирования

### Добавить Authorized redirect URIs:
```
http://localhost:3001/ru/auth/oauth/callback
http://localhost:3001/en/auth/oauth/callback
http://localhost:3001/sr/auth/oauth/callback
http://localhost:3000/ru/auth/oauth/callback
http://localhost:3000/en/auth/oauth/callback
http://localhost:3000/sr/auth/oauth/callback
https://svetu.rs/ru/auth/oauth/callback
https://svetu.rs/en/auth/oauth/callback
https://svetu.rs/sr/auth/oauth/callback
https://www.svetu.rs/ru/auth/oauth/callback
https://www.svetu.rs/en/auth/oauth/callback
https://www.svetu.rs/sr/auth/oauth/callback
```

### Проверить Authorized JavaScript origins:
```
http://localhost:3001
http://localhost:3000
https://svetu.rs
https://www.svetu.rs
```

## Текущие credentials
- **Client ID**: `917315728307-au9ga5fl7o3bbid9nv7e4l92gut194pq.apps.googleusercontent.com`
- **Client Secret**: `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`

## После настройки
Изменения применяются мгновенно, перезапуск сервисов не требуется.

## Проверка
1. Открыть http://localhost:3001
2. Нажать "Войти" → "Войти через Google"
3. Должен произойти редирект на Google
4. После авторизации должен вернуться на http://localhost:3001/ru/auth/oauth/callback (или другую локаль)
5. Frontend должен обработать код и получить токены

## Важно!
Frontend использует i18n маршрутизацию, поэтому все OAuth callback URL должны включать локаль (/ru/, /en/, /sr/).
Auth Service настроен на использование русской локали по умолчанию.