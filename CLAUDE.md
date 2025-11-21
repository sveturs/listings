# CLAUDE.md - Listings Microservice

## 🎯 О микросервисе

**Listings Service** - микросервис для управления объявлениями, заказами, корзиной, избранным и чатами.

- **Порт gRPC:** 50053
- **Порт HTTP:** 8086
- **База данных:** `listings_dev_db` (PostgreSQL порт 35434)
- **Redis:** порт 36380
- **Директория:** `/p/github.com/sveturs/listings`

---

## 🔴 КРИТИЧЕСКИ ВАЖНО: База данных

**Микросервис использует ОТДЕЛЬНУЮ БД, а НЕ монолитную!**

### ✅ ПРАВИЛЬНАЯ конфигурация (.env):

```bash
# Database - Отдельная БД микросервиса (НЕ монолит!)
SVETULISTINGS_DB_HOST=localhost
SVETULISTINGS_DB_PORT=35434              # НЕ 5433!
SVETULISTINGS_DB_USER=listings_user      # НЕ postgres!
SVETULISTINGS_DB_PASSWORD=listings_secret
SVETULISTINGS_DB_NAME=listings_dev_db    # НЕ svetubd!
SVETULISTINGS_DB_SSLMODE=disable
```

### ❌ НЕПРАВИЛЬНАЯ конфигурация:

```bash
# НЕ ДЕЛАЙ ТАК - это монолитная БД!
SVETULISTINGS_DB_PORT=5433     # ❌ Это монолит!
SVETULISTINGS_DB_NAME=svetubd  # ❌ Это монолит!
SVETULISTINGS_DB_USER=postgres # ❌ Это монолит!
```

### 🐳 Docker контейнер БД:

```bash
# Проверить контейнер
docker ps | grep listings_postgres
# Вывод: listings_postgres   postgres:15-alpine   0.0.0.0:35434->5432/tcp

# Подключиться к БД
psql "postgres://listings_user:listings_secret@localhost:35434/listings_dev_db"

# Проверить таблицы
psql "postgres://listings_user:listings_secret@localhost:35434/listings_dev_db" -c "\dt"
```

---

## 🚀 Запуск и остановка

### Быстрый запуск:

```bash
# Запустить микросервис
/home/dim/.local/bin/start-listings-microservice.sh

# Остановить микросервис
/home/dim/.local/bin/stop-listings-microservice.sh

# Проверить статус
netstat -tlnp | grep ":50053"
tail -f /tmp/listings-microservice.log
```

### Ручной запуск:

```bash
# 1. Остановить старые процессы
/home/dim/.local/bin/kill-port-50053.sh

# 2. Закрыть screen сессии
screen -ls | grep listings-microservice | awk '{print $1}' | xargs -I {} screen -S {} -X quit
screen -wipe

# 3. Запустить
cd /p/github.com/sveturs/listings
screen -dmS listings-microservice-50053 bash -c 'go run ./cmd/server/main.go 2>&1 | tee /tmp/listings-microservice.log'

# 4. Проверить
netstat -tlnp | grep ":50053"
```

---

## 📋 Структура БД

### Основные таблицы:

- **listings** - унифицированная таблица объявлений (C2C + B2C)
- **listing_favorites** - избранное пользователей
- **listing_images** - изображения объявлений
- **listing_locations** - геолокация объявлений
- **listing_attributes** - атрибуты объявлений
- **chats** - чаты по объявлениям
- **messages** - сообщения в чатах
- **chat_attachments** - вложения в сообщениях
- **orders** - заказы
- **cart_items** - корзина покупок
- **storefronts** - витрины магазинов

### Проверка данных:

```bash
# Количество объявлений
psql "postgres://listings_user:listings_secret@localhost:35434/listings_dev_db" \
  -c "SELECT COUNT(*) FROM listings;"

# Количество избранных
psql "postgres://listings_user:listings_secret@localhost:35434/listings_dev_db" \
  -c "SELECT COUNT(*) FROM listing_favorites;"

# Список таблиц
psql "postgres://listings_user:listings_secret@localhost:35434/listings_dev_db" \
  -c "\dt"
```

---

## 🔧 Feature Flags (в монолите)

Включение микросервиса в монолите через переменные окружения:

```bash
# В /p/github.com/sveturs/svetu/backend/.env
USE_LISTINGS_MICROSERVICE=true
USE_ORDERS_MICROSERVICE=true
USE_SEARCH_MICROSERVICE=true
USE_ANALYTICS_MICROSERVICE=true
USE_CHAT_MICROSERVICE=true

LISTINGS_GRPC_URL=localhost:50053
LISTINGS_GRPC_TIMEOUT=10s
```

---

## 🧪 Тестирование

### Проверка доступности:

```bash
# Health check
curl http://localhost:8086/health

# Метрики
curl http://localhost:8086/metrics
```

### Проверка избранного:

```bash
# Получить токен
TOKEN=$(cat /tmp/token)

# Получить список избранного
curl -s -H "Authorization: Bearer $TOKEN" \
  http://localhost:3000/api/v1/marketplace/favorites | jq '.'

# Добавить в избранное
curl -s -X POST -H "Authorization: Bearer $TOKEN" \
  http://localhost:3000/api/v1/marketplace/favorites/11 | jq '.'

# Удалить из избранного
curl -s -X DELETE -H "Authorization: Bearer $TOKEN" \
  http://localhost:3000/api/v1/marketplace/favorites/11 | jq '.'
```

---

## 🐛 Troubleshooting

### Проблема: "listing not found"

**Причина:** Микросервис подключен к монолитной БД вместо своей.

**Решение:**
1. Проверь `.env`:
   ```bash
   cat .env | grep DB_PORT
   # Должно быть: SVETULISTINGS_DB_PORT=35434
   ```

2. Исправь конфигурацию (см. раздел "База данных" выше)

3. Перезапусти микросервис:
   ```bash
   /home/dim/.local/bin/stop-listings-microservice.sh
   /home/dim/.local/bin/start-listings-microservice.sh
   ```

### Проблема: "relation listing_favorites does not exist"

**Причина:** Неправильная БД или миграции не применены.

**Решение:**
1. Проверь подключение к правильной БД (порт 35434)
2. Проверь наличие таблиц:
   ```bash
   psql "postgres://listings_user:listings_secret@localhost:35434/listings_dev_db" \
     -c "\dt listing*"
   ```

---

## 📚 Документация

- **Architecture:** [docs/DATABASE_ARCHITECTURE.md](docs/DATABASE_ARCHITECTURE.md)
- **Migration Plan:** [/p/github.com/sveturs/svetu/docs/migration/MIGRATION_PLAN_TO_MICROSERVICE.md](/p/github.com/sveturs/svetu/docs/migration/MIGRATION_PLAN_TO_MICROSERVICE.md)
- **Chat Design:** [/p/github.com/sveturs/CHAT_MICROSERVICE_DESIGN.md](/p/github.com/sveturs/CHAT_MICROSERVICE_DESIGN.md)

---

## ✅ Чеклист перед запуском

- [ ] `.env` указывает на порт 35434 (НЕ 5433)
- [ ] `.env` указывает на БД `listings_dev_db` (НЕ `svetubd`)
- [ ] Docker container `listings_postgres` запущен
- [ ] Таблица `listings` содержит данные
- [ ] Таблица `listing_favorites` существует
- [ ] Redis доступен на порту 36380
- [ ] OpenSearch доступен на порту 9200
- [ ] MinIO доступен на `s3.svetu.rs`

---

**Последнее обновление:** 2025-11-21
