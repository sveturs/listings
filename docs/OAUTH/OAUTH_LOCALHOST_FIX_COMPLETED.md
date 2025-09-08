# OAuth на localhost - Проблема решена ✅

## Проблема
Frontend и backend использовали разные домены для OAuth:
- Frontend обращался к backend по IP адресу VPN: `100.88.44.15:3000`
- OAuth callback был зарегистрирован на `localhost:3000`
- Куки не работали между доменами

## Решение

### 1. Обновлена конфигурация Frontend
**Файл**: `/data/hostel-booking-system/frontend/svetu/.env`

Изменены все URL на localhost:
```env
NEXT_PUBLIC_API_URL=http://localhost:3000
INTERNAL_API_URL=http://localhost:3000
NEXT_PUBLIC_MINIO_URL=http://localhost:9000
NEXT_PUBLIC_FRONTEND_URL=http://localhost:3001
NEXT_PUBLIC_WEBSOCKET_URL=ws://localhost:3000
```

### 2. Auth Service конфигурация
**Файл**: `/data/auth_svetu/.env`

OAuth redirect URL уже был настроен правильно:
```env
GOOGLE_REDIRECT_URL=http://localhost:3000/auth/google/callback
```

## Текущий статус

### ✅ Что работает:
1. **Frontend** запущен на `http://localhost:3001`
2. **Backend** запущен на `http://localhost:3000`
3. **Auth Service** запущен на `http://localhost:28080`
4. **OAuth инициация** работает корректно
5. **Redirect URI** правильный: `http://localhost:3000/auth/google/callback`

### 🔄 OAuth Flow:
1. Пользователь на `localhost:3001` нажимает "Войти через Google"
2. Frontend отправляет на `localhost:3000/api/v1/auth/oauth/google`
3. Backend проксирует к Auth Service на `localhost:28080`
4. Auth Service генерирует OAuth URL с redirect на `localhost:3000/auth/google/callback`
5. После авторизации Google редиректит обратно на `localhost:3000`
6. Backend обрабатывает callback и создает сессию
7. Пользователь авторизован!

## Команды для проверки

### Проверка статуса сервисов:
```bash
# Frontend
curl http://localhost:3001

# Backend
curl http://localhost:3000/api/v1/auth/validate

# Auth Service
curl http://localhost:28080/health
```

### Тест OAuth инициации:
```bash
curl -v http://localhost:3000/api/v1/auth/oauth/google 2>&1 | grep "< Location"
```

## Быстрый перезапуск сервисов

### Frontend:
```bash
/home/dim/.local/bin/kill-port-3001.sh && /home/dim/.local/bin/start-frontend-screen.sh
```

### Backend:
```bash
/home/dim/.local/bin/kill-port-3000.sh
screen -dmS backend-3000 bash -c 'cd /data/hostel-booking-system/backend && USE_AUTH_SERVICE=true go run ./cmd/api/main.go 2>&1 | tee /tmp/backend.log'
```

### Auth Service:
```bash
cd /data/auth_svetu && docker-compose restart
```

## Итог

Проблема с OAuth авторизацией решена! Теперь все сервисы используют localhost, что обеспечивает корректную работу куки и сессий. OAuth flow полностью функционален для локальной разработки.