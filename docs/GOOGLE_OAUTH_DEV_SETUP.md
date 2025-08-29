# 🔐 Настройка Google OAuth для dev.svetu.rs

## Проблема
Google OAuth работает локально и на production (svetu.rs), но не работает на dev.svetu.rs из-за неправильных redirect URLs.

## Решение: Настройка в Google Cloud Console

### 1. Откройте Google Cloud Console
1. Перейдите в [Google Cloud Console](https://console.cloud.google.com/)
2. Выберите проект: **neat-environs-140712**
3. Перейдите в **APIs & Services** → **Credentials**

### 2. Найдите OAuth 2.0 Client ID
Найдите клиента с ID: `917315728307-au9ga5fl7o3bbid9nv7e4l92gut194pq.apps.googleusercontent.com`

### 3. Добавьте Authorized redirect URIs

В разделе **Authorized redirect URIs** добавьте новые записи:
```
https://devapi.svetu.rs/auth/google/callback
https://devapi.svetu.rs/api/v1/auth/google/callback
```

**Существующие URI (оставить как есть):**
```
http://localhost:3000/auth/google/callback
https://api.svetu.rs/auth/google/callback
https://api.svetu.rs/api/v1/auth/google/callback
```

### 4. Проверьте Authorized JavaScript origins

Убедитесь что есть следующие домены:
```
http://localhost:3000
http://localhost:3001
https://svetu.rs
https://dev.svetu.rs
https://devapi.svetu.rs
```

### 5. Конфигурация backend на dev сервере

Проверьте настройки в `/opt/svetu-dev/.env`:
```bash
# Google OAuth Configuration
GOOGLE_CLIENT_ID=917315728307-au9ga5fl7o3bbid9nv7e4l92gut194pq.apps.googleusercontent.com
GOOGLE_CLIENT_SECRET=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
GOOGLE_OAUTH_REDIRECT_URL=https://devapi.svetu.rs/auth/google/callback
FRONTEND_URL=https://dev.svetu.rs
```

### 6. Конфигурация frontend на dev сервере

В файле окружения frontend должно быть:
```bash
NEXT_PUBLIC_GOOGLE_CLIENT_ID=917315728307-au9ga5fl7o3bbid9nv7e4l92gut194pq.apps.googleusercontent.com
```

## 🔧 Проверка настроек на сервере

```bash
# Подключиться к dev серверу
ssh root@svetu.rs

# Проверить backend настройки
cd /opt/svetu-dev
grep -E "GOOGLE|FRONTEND" .env

# Перезапустить сервисы после изменения .env
docker-compose restart backend frontend
```

## 🧪 Тестирование OAuth

### 1. Откройте https://dev.svetu.rs
### 2. Нажмите кнопку "Sign in with Google"
### 3. Должен произойти редирект на Google OAuth
### 4. После авторизации должен вернуться на https://dev.svetu.rs

## ⏱️ Время вступления в силу

Изменения в Google Console вступают в силу через **5-10 минут**.

## 🔍 Отладка проблем

### Проверка логов backend:
```bash
ssh root@svetu.rs "docker logs svetu-dev_backend_1 --tail=50"
```

### Проверка Network tab в браузере:
1. Откройте Developer Tools (F12)
2. Вкладка Network
3. Попробуйте авторизацию
4. Ищите ошибки 400/403/404 при редиректе

### Частые ошибки:
- **redirect_uri_mismatch** - неправильный URI в Google Console
- **invalid_client** - неправильный CLIENT_ID
- **CORS errors** - проблемы с доменами в JavaScript origins

## 📋 Checklist настройки

- [ ] Добавлены redirect URIs в Google Console
- [ ] Добавлены JavaScript origins в Google Console  
- [ ] Обновлена конфигурация backend на dev сервере
- [ ] Обновлена конфигурация frontend на dev сервере
- [ ] Перезапущены сервисы на dev сервере
- [ ] Протестирована авторизация через браузер
- [ ] Подождали 5-10 минут для вступления изменений в силу

## 🎯 Результат

После выполнения всех шагов Google OAuth должен работать на всех окружениях:
- ✅ Локально: http://localhost:3001
- ✅ Production: https://svetu.rs  
- ✅ Dev: https://dev.svetu.rs

---
*Создано: 28.08.2025*  
*Автор: Система настройки OAuth*  
*Статус: Готово к применению*