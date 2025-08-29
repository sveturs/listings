# 🔧 Руководство по управлению конфигурациями окружений

## Проблемы которые были найдены и исправлены:

### ❌ Захардкоженные URL в коде
- ✅ `frontend/src/services/importApi.ts` - заменены на `configManager.getApiUrl()`
- ✅ `frontend/src/lib/api-client.ts` - заменены на `configManager.getApiUrl()`
- ✅ `frontend/src/lib/api.ts` - заменены на `configManager.getApiUrl()`
- ✅ `backend/cmd/cli/test_behavior_events.go` - используют переменные окружения

### ✅ Правильные конфигурации для всех окружений

## 🌍 Окружения и их конфигурации

### 1. **Локальная разработка** (localhost)
```bash
# Используй стандартный .env
cp frontend/svetu/.env.example frontend/svetu/.env.local
```

**Настройки:**
- API: `http://localhost:3000`
- Frontend: `http://localhost:3001`
- MinIO: `http://localhost:9000`

### 2. **VPN доступ** (100.88.44.15)
```bash
# Используй специальную конфигурацию для Tailscale
cp frontend/svetu/.env.tailscale frontend/svetu/.env.local
```

**Настройки:**
- API: `http://100.88.44.15:3000`
- Frontend: `http://100.88.44.15:3001` 
- MinIO: `http://100.88.44.15:9000`

### 3. **Dev сервер** (dev.svetu.rs)
**Настройки на сервере (/opt/svetu-dev/.env):**
- API: `https://devapi.svetu.rs` (порт 3002 внутри Docker)
- Frontend: `https://dev.svetu.rs` (порт 3003 внутри Docker)
- MinIO: `https://devs3.svetu.rs` (порт 9002 внутри Docker)

### 4. **Production** (svetu.rs)
```bash
# Используй production конфигурацию
cp frontend/svetu/.env.production frontend/svetu/.env.production.local
```

**Настройки:**
- API: `https://api.svetu.rs`
- Frontend: `https://svetu.rs`
- MinIO: `https://s3.svetu.rs`

## 🚀 Простое использование

### ✅ По умолчанию - VPN доступ работает автоматически!
```bash
cd /data/hostel-booking-system/frontend/svetu
yarn dev -p 3001
```

**Доступ:**
- С компьютера: http://localhost:3001/en
- С телефона (VPN): http://100.88.44.15:3001/en

### Только localhost (если нужно)
```bash
cd /data/hostel-booking-system/frontend/svetu
cp .env.localhost .env.local  # переключить на localhost only
yarn dev -p 3001
```

### Возврат к VPN+localhost
```bash
cd /data/hostel-booking-system/frontend/svetu
rm .env.local  # использовать основной .env
yarn dev -p 3001
```

## 📋 Статус dev.svetu.rs

### ✅ Что работает:
- Frontend: https://dev.svetu.rs (контейнер здоровый)
- Backend API: https://devapi.svetu.rs/api/v1/health
- MinIO: https://devs3.svetu.rs
- Nginx: запущен в Docker контейнере
- OpenSearch: http://svetu.rs:9201
- PostgreSQL: port 5433

### ⚠️ Известные проблемы:
- Frontend логи показывают ошибки переводов (`MISSING_MESSAGE`)
- Системный nginx не запущен (но Docker nginx работает)

### 🔧 Полезные команды для мониторинга dev.svetu.rs:
```bash
# Статус всех сервисов
ssh root@svetu.rs "cd /opt/svetu-dev && docker-compose ps"

# Перезапуск проблемных сервисов
ssh root@svetu.rs "cd /opt/svetu-dev && docker-compose restart frontend"
ssh root@svetu.rs "cd /opt/svetu-dev && docker-compose restart backend"

# Логи сервисов
ssh root@svetu.rs "docker logs svetu-dev_frontend_1 --tail=20"
ssh root@svetu.rs "docker logs svetu-dev_backend_1 --tail=20"
```

## 🎯 Рекомендуемый workflow

1. **Для локальной разработки:** используй `.env.example` → `.env.local`
2. **Для тестирования по VPN:** используй `.env.tailscale` → `.env.local`
3. **Для деплоя на dev.svetu.rs:** изменения автоматически подтянутся через Docker
4. **Для production:** используй `.env.production`

## 🔐 Безопасность

- ✅ Все захардкоженные URL заменены на динамические через `configManager`
- ✅ Создан `.env.tailscale` для безопасного VPN доступа
- ✅ `.env.production` настроен для правильного production деплоя
- ❗ API ключи остаются в приватном репозитории (это безопасно)

## 🧪 Тестирование

После изменения конфигурации всегда проверяй:
```bash
# Frontend
cd frontend/svetu
yarn format && yarn lint && yarn build

# Backend
cd backend
make format && make lint && go build ./...
```