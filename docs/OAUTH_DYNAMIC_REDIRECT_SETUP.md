# OAuth Dynamic Redirect Setup

## Проблема
При доступе к приложению через VPN или локальную сеть (например, `http://100.88.44.15`) происходит редирект на `localhost:3000/auth/google/callback`, который не работает на мобильных устройствах и удаленных машинах.

## Решение
Реализована поддержка динамического определения callback URL на основе хоста, с которого пришел запрос.

## Изменения в коде

### 1. Backend Handler (`backend/internal/proj/users/handler/auth.go`)
- Добавлено определение хоста из заголовков запроса
- Автоматическое формирование callback URL на основе текущего хоста
- Передача динамического callback URL в сервис авторизации

### 2. Auth Service (`backend/internal/proj/users/service/auth.go`)
- Новый метод `GetGoogleAuthURLWithCallback` для работы с динамическими redirect URL
- Создание временной конфигурации OAuth2 с нужным redirect URL
- Fallback на стандартный метод при отсутствии динамического хоста

## Настройка Google OAuth Console

### Добавление разрешенных redirect URI

1. Перейдите в [Google Cloud Console](https://console.cloud.google.com/)
2. Выберите ваш проект
3. Перейдите в **APIs & Services** → **Credentials**
4. Найдите ваш OAuth 2.0 Client ID
5. Нажмите на него для редактирования
6. В разделе **Authorized redirect URIs** добавьте следующие URL:

```
# Локальная разработка
http://localhost:3000/auth/google/callback
http://localhost:3001/auth/google/callback

# VPN и локальная сеть (примеры)
http://100.88.44.15:3000/auth/google/callback
http://192.168.1.100:3000/auth/google/callback
http://10.0.0.100:3000/auth/google/callback

# Tailscale VPN (если используется)
http://your-machine-name:3000/auth/google/callback
https://your-machine-name.tailnet:3000/auth/google/callback

# Dev и Production
https://dev.svetu.rs/auth/google/callback
https://svetu.rs/auth/google/callback
```

### Важные моменты

1. **Безопасность**: Добавляйте только те URL, которые вы действительно используете
2. **Порты**: Backend всегда работает на порту 3000, frontend на 3001
3. **Протокол**: Используйте HTTPS для production, HTTP для локальной разработки
4. **Wildcard**: Google OAuth не поддерживает wildcard в redirect URI, нужно добавлять каждый URL отдельно

## Как это работает

1. Пользователь заходит на сайт через VPN: `http://100.88.44.15:3001`
2. При клике на "Sign in with Google" запрос идет на backend
3. Backend определяет хост из заголовков и формирует callback URL: `http://100.88.44.15:3000/auth/google/callback`
4. Google OAuth использует этот динамический URL для возврата
5. После успешной авторизации пользователь возвращается на frontend

## Переменные окружения

В `.env` файле можно оставить стандартный redirect URL как fallback:
```env
GOOGLE_OAUTH_REDIRECT_URL=http://localhost:3000/auth/google/callback
```

Он будет использоваться только если не удалось определить хост из запроса.

## Тестирование

1. **Localhost**: http://localhost:3001 - должен работать как обычно
2. **VPN**: http://[VPN_IP]:3001 - должен корректно редиректить
3. **Production**: https://svetu.rs - использует production callback URL

## Логирование

В логах backend можно увидеть какой callback URL используется:
```
GoogleAuth: processing OAuth request host_header=100.88.44.15:3001 callback_host=http://100.88.44.15:3000
GetGoogleAuthURLWithCallback: using dynamic redirect URL redirect_url=http://100.88.44.15:3000/auth/google/callback
```

## Troubleshooting

### Ошибка "redirect_uri_mismatch"
- Убедитесь, что URL добавлен в Google Console
- Проверьте логи, какой именно URL используется
- Убедитесь, что порт правильный (3000 для backend)

### Редирект на localhost вместо VPN адреса
- Проверьте, что frontend делает запрос на правильный backend URL
- Убедитесь, что в nginx/proxy настройках передаются правильные заголовки

### Cookie не сохраняются
- Проверьте настройки SameSite в cookies
- Для cross-origin запросов может потребоваться настройка CORS