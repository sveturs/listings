# OAuth через VPN/LAN с использованием localhost redirect

## Проблема
Google OAuth не позволяет использовать IP адреса в качестве redirect URI. Требуется домен с публичным TLD (.com, .org и т.д.).

## Решение
Реализован механизм прокси-редиректа через localhost с сохранением контекста оригинального хоста.

## Как это работает

### 1. Инициация OAuth (VPN → localhost)
Когда пользователь с VPN адреса (например, `100.88.44.15:3001`) нажимает "Sign in with Google":

1. Frontend отправляет запрос на backend: `http://100.88.44.15:3000/auth/google`
2. Backend сохраняет реальный хост в cookie `oauth_real_host`
3. Google OAuth использует стандартный `localhost:3000` redirect URI
4. Пользователь авторизуется в Google

### 2. OAuth Callback (localhost → VPN)
После успешной авторизации в Google:

1. Google редиректит на `http://localhost:3000/auth/google/callback?code=...`
2. Backend проверяет cookie `oauth_real_host`
3. Если cookie существует, происходит редирект на реальный хост:
   ```
   http://100.88.44.15:3000/auth/google/callback?code=...&state=...
   ```
4. На реальном хосте код обрабатывается и создается сессия
5. Пользователь возвращается на frontend с токеном

## Реализация

### Backend изменения

#### 1. GoogleAuth handler (`/auth/google`)
```go
// Сохраняем реальный хост для non-localhost запросов
if host != "" && !strings.HasPrefix(host, "localhost") {
    realBackendURL := fmt.Sprintf("%s://%s", protocol, host)
    c.Cookie(&fiber.Cookie{
        Name:     "oauth_real_host",
        Value:    realBackendURL,
        Path:     "/",
        MaxAge:   300, // 5 минут
        Secure:   false, // Важно для HTTP
        HTTPOnly: true,
        SameSite: "Lax",
    })
}
```

#### 2. GoogleCallback handler (`/auth/google/callback`)
```go
// Проверяем сохраненный реальный хост
realHost := c.Cookies("oauth_real_host")
if realHost != "" {
    // Редиректим на реальный хост с кодом
    redirectURL := fmt.Sprintf("%s/auth/google/callback?code=%s&state=%s", 
                               realHost, code, state)
    // Удаляем cookie
    c.Cookie(&fiber.Cookie{
        Name:   "oauth_real_host",
        Value:  "",
        Path:   "/",
        MaxAge: -1,
    })
    return c.Redirect(redirectURL)
}
```

## Настройка

### 1. Google Console
В Google Cloud Console нужен только один redirect URI:
```
http://localhost:3000/auth/google/callback
```

### 2. Frontend конфигурация
Frontend должен делать запросы на тот же backend хост, откуда загружена страница:
- При доступе через `http://100.88.44.15:3001` → запросы на `http://100.88.44.15:3000`
- При доступе через `http://localhost:3001` → запросы на `http://localhost:3000`

### 3. Cookie настройки
Важные параметры для работы через HTTP:
- `Secure: false` - для работы через HTTP
- `SameSite: "Lax"` - для кросс-сайтовых редиректов
- `HTTPOnly: true` - для безопасности

## Преимущества

1. **Не требует изменений в Google Console** - используется только localhost redirect
2. **Работает с любыми IP адресами** - VPN, LAN, Docker и т.д.
3. **Сохраняет контекст** - пользователь возвращается на тот же хост
4. **Безопасность** - cookie с коротким временем жизни и HTTPOnly

## Ограничения

1. **Требует доступа к localhost:3000** - браузер должен иметь возможность открыть localhost
2. **Два редиректа** - сначала на localhost, потом на реальный хост
3. **Cookie должны быть включены** - решение зависит от cookies

## Альтернативные решения

### 1. SSH туннель
```bash
# На клиенте
ssh -L 3000:localhost:3000 user@vpn-server
```
Затем используйте `http://localhost:3000` для доступа

### 2. Использование доменного имени
Настройте DNS или `/etc/hosts`:
```
100.88.44.15 myapp.local
```
Но Google OAuth все равно не примет `.local` домен

### 3. Ngrok или другие туннели
```bash
ngrok http 3000
```
Получите публичный URL типа `https://abc123.ngrok.io`

## Тестирование

### Проверка через curl:
```bash
# 1. Инициация OAuth с VPN хоста
curl -s -H "Host: 100.88.44.15:3000" \
     -H "Origin: http://100.88.44.15:3001" \
     http://localhost:3000/auth/google -I

# Должен вернуть:
# - Location: https://accounts.google.com/...redirect_uri=localhost...
# - Set-Cookie: oauth_real_host=http://100.88.44.15:3000

# 2. Эмуляция callback
curl -s -b "oauth_real_host=http://100.88.44.15:3000" \
     http://localhost:3000/auth/google/callback?code=test -I

# Должен вернуть:
# - Location: http://100.88.44.15:3000/auth/google/callback?code=test...
```

## Логирование

В логах backend можно увидеть процесс:
```
GoogleAuth: Saving real host for non-localhost access real_backend_url=http://100.88.44.15:3000
GoogleCallback: Found saved real host, will redirect there real_host=http://100.88.44.15:3000
```

## FAQ

**Q: Почему не работает с телефона через VPN?**
A: Телефон не имеет доступа к localhost:3000. Нужно использовать SSH туннель или ngrok.

**Q: Можно ли использовать HTTPS?**
A: Да, но нужен валидный SSL сертификат для IP адреса, что сложно получить.

**Q: Работает ли с Docker?**
A: Да, если контейнер проброшен на localhost:3000 хоста.